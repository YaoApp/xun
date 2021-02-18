package postgres

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/logger"
	"github.com/yaoapp/xun/utils"
)

// SchemaName Get the schema name of the connection
func (grammarSQL *Postgres) SchemaName() string {
	uinfo, err := url.Parse(grammarSQL.DSN)
	if err != nil {
		return "unknown"
	}
	schema := uinfo.Query().Get("search_path")
	if schema == "" {
		schema = "public"
	}
	return grammarSQL.Schema
}

// Create a new table on the schema
func (grammarSQL Postgres) Create(table *grammar.Table, db *sqlx.DB) error {
	name := grammarSQL.Quoter.ID(table.Name, db)
	sql := fmt.Sprintf("CREATE TABLE %s (\n", name)
	stmts := []string{}

	columns := []*grammar.Column{}
	indexes := []*grammar.Index{}

	// Commands
	// The commands must be:
	//    AddColumn(column *Column)    for adding a column
	//    ModifyColumn(column *Column) for modifying a colu
	//    RenameColumn(old string,new string)  for renaming a column
	//    DropColumn(name string)  for dropping a column
	//    CreateIndex(index *Index) for creating a index
	//    DropIndex( name string) for  dropping a index
	//    RenameIndex(old string,new string)  for renaming a index
	for _, command := range table.Commands {
		switch command.Name {
		case "AddColumn":
			columns = append(columns, command.Params[0].(*grammar.Column))
			break
		case "CreateIndex":
			indexes = append(indexes, command.Params[0].(*grammar.Index))
			break
		}
	}

	// Columns
	for _, Column := range columns {
		stmts = append(stmts,
			grammarSQL.Builder.SQLAddColumn(db, Column, grammarSQL.Types, grammarSQL.Quoter),
		)
	}
	// Primary key
	for _, index := range indexes {
		if index.Type == "primary" {
			stmts = append(stmts,
				grammarSQL.Builder.SQLAddIndex(db, index, grammarSQL.IndexTypes, grammarSQL.Quoter),
			)
		}
	}
	sql = sql + strings.Join(stmts, ",\n")
	sql = sql + fmt.Sprintf("\n)")

	// Create table
	defer logger.Debug(logger.CREATE, sql).TimeCost(time.Now())
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}

	// indexes
	indexStmts := []string{}
	for _, index := range indexes {
		if index.Type == "primary" {
			continue
		}
		indexStmts = append(indexStmts,
			grammarSQL.Builder.SQLAddIndex(db, index, grammarSQL.IndexTypes, grammarSQL.Quoter),
		)
	}
	defer logger.Debug(logger.CREATE, indexStmts...).TimeCost(time.Now())
	_, err = db.Exec(strings.Join(indexStmts, ";\n"))
	if err != nil {
		return err
	}

	return nil
}

// Get a table on the schema
func (grammarSQL Postgres) Get(table *grammar.Table, db *sqlx.DB) error {
	columns, err := grammarSQL.GetColumnListing(table.DBName, table.Name, db)
	if err != nil {
		return err
	}

	indexes, err := grammarSQL.GetIndexListing(table.DBName, table.Name, db)
	if err != nil {
		return err
	}

	// attaching columns
	for _, column := range columns {
		column.Indexes = []*grammar.Index{}
		table.PushColumn(column)
	}

	// attaching indexes
	for i := range indexes {
		idx := indexes[i]
		if !table.HasColumn(idx.ColumnName) {
			return errors.New("the column does not exists" + idx.ColumnName)
		}
		column := table.ColumnMap[idx.ColumnName]
		if !table.HasIndex(idx.Name) {
			index := *idx
			index.Columns = []*grammar.Column{}
			column.Indexes = append(column.Indexes, &index)
			table.PushIndex(&index)
		}
		index := table.IndexMap[idx.Name]
		index.Columns = append(index.Columns, column)
	}

	return nil
}

// GetIndexListing get a table indexes structure
func (grammarSQL Postgres) GetIndexListing(dbName string, tableName string, db *sqlx.DB) ([]*grammar.Index, error) {
	selectColumns := []string{
		"\"TABLE_SCHEMA\" AS \"db_name\"",
		"\"TABLE_NAME\" AS \"table_name\"",
		"\"INDEX_NAME\" AS \"index_name\"",
		"\"COLUMN_NAME\" AS \"column_name\"",
		"\"COLLATION\" AS \"collation\"",
		`CASE
			WHEN NULLABLE = 'YES' THEN true
			WHEN NULLABLE = "NO" THEN false
			ELSE false
		END AS \" + "\"nullable\"`,
		`CASE
			WHEN NON_UNIQUE = 0 THEN true
			WHEN NON_UNIQUE = 1 THEN false
			ELSE 0
		END AS \" + "\"unique\"`,
		"\"COMMENT\" AS \"comment\"",
		"\"INDEX_TYPE\" AS \"index_type\"",
		"\"COLUMN_NAME\" AS \"column_name\"",
		"\"SEQ_IN_INDEX\" AS \"seq_in_index\"",
		"\"INDEX_COMMENT\" AS \"index_comment\"",
	}
	sql := fmt.Sprintf(`
			SELECT %s
			FROM INFORMATION_SCHEMA.STATISTICS
			WHERE TABLE_SCHEMA = %s AND TABLE_NAME = %s
			ORDER BY SEQ_IN_INDEX;
			`,
		strings.Join(selectColumns, ","),
		grammarSQL.Quoter.VAL(dbName, db),
		grammarSQL.Quoter.VAL(tableName, db),
	)
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	indexes := []*grammar.Index{}
	err := db.Select(&indexes, sql)
	if err != nil {
		return nil, err
	}

	// counting the type of indexes
	for _, index := range indexes {
		if index.Name == "PRIMARY" {
			index.Type = "primary"
		} else if index.Unique {
			index.Type = "unique"
		} else {
			index.Type = "index"
		}
	}
	return indexes, nil
}

// GetColumnListing get a table columns structure
func (grammarSQL Postgres) GetColumnListing(dbName string, tableName string, db *sqlx.DB) ([]*grammar.Column, error) {
	selectColumns := []string{
		"TABLE_SCHEMA AS \"db_name\"",
		"TABLE_NAME AS \"table_name\"",
		"COLUMN_NAME AS \"name\"",
		"ORDINAL_POSITION AS \"position\"",
		"COLUMN_DEFAULT AS \"default\"",
		`CASE
			WHEN IS_NULLABLE = 'YES' THEN true
			WHEN IS_NULLABLE = "NO" THEN false
			ELSE false
		END AS \" + "\"nullable\"`,
		`CASE
		   WHEN LOCATE('unsigned', COLUMN_TYPE) THEN true
		   ELSE false
		END AS\" + "\"unsigned\"`,
		"UPPER(DATA_TYPE) as \"type\"",
		"CHARACTER_MAXIMUM_LENGTH as \"length\"",
		"CHARACTER_OCTET_LENGTH as \"octet_length\"",
		"NUMERIC_PRECISION as \"precision\"",
		"NUMERIC_SCALE as \"scale\"",
		"DATETIME_PRECISION as \"datetime_precision\"",
		"CHARACTER_SET_NAME as \"charset\"",
		"COLLATION_NAME as \"collation\"",
		"COLUMN_KEY as \"key\"",
		`CASE
			WHEN COLUMN_KEY = 'PRI' THEN true
			ELSE false
		END AS \" + "\"primary\"`,
		"EXTRA as \"extra\"",
		"COLUMN_COMMENT as \"comment\"",
	}
	sql := fmt.Sprintf(`
			SELECT %s
			FROM INFORMATION_SCHEMA.COLUMNS
			WHERE  TABLE_CATALOG = %s AND TABLE_NAME = %s;
			`,
		strings.Join(selectColumns, ","),
		grammarSQL.Quoter.VAL(dbName, db),
		grammarSQL.Quoter.VAL(tableName, db),
	)
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	columns := []*grammar.Column{}
	err := db.Select(&columns, sql)
	if err != nil {
		return nil, err
	}

	// Cast the database data type to DBAL data type
	for _, column := range columns {
		typ, has := grammarSQL.FlipTypes[column.Type]
		if has {
			column.Type = typ
		}

		if utils.StringVal(column.Extra) == "auto_increment" {
			column.Extra = utils.StringPtr("AutoIncrement")
		}
	}
	return columns, nil
}
