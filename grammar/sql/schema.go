package sql

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/logger"
	"github.com/yaoapp/xun/utils"
)

// DBName Get the database name of the connection
func (grammarSQL *SQL) DBName() string {
	uinfo, err := url.Parse(grammarSQL.DSN)
	if err != nil {
		return "unknown"
	}
	grammarSQL.DB = filepath.Base(uinfo.Path)
	return grammarSQL.DB
}

// SchemaName Get the schema name of the connection
func (grammarSQL *SQL) SchemaName() string {
	schema := grammarSQL.DBName()
	grammarSQL.Schema = schema
	return grammarSQL.Schema
}

// Exists the Exists
func (grammarSQL SQL) Exists(name string, db *sqlx.DB) bool {
	sql := grammarSQL.Builder.SQLTableExists(db, name, grammarSQL.Quoter)
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	row := db.QueryRowx(sql)
	if row.Err() != nil {
		panic(row.Err())
	}
	res, err := row.SliceScan()
	if err != nil {
		return false
	}
	return name == fmt.Sprintf("%s", res[0])
}

// Get a table on the schema
func (grammarSQL SQL) Get(table *grammar.Table, db *sqlx.DB) error {
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
func (grammarSQL SQL) GetIndexListing(dbName string, tableName string, db *sqlx.DB) ([]*grammar.Index, error) {
	selectColumns := []string{
		"`TABLE_SCHEMA` AS `db_name`",
		"`TABLE_NAME` AS `table_name`",
		"`INDEX_NAME` AS `index_name`",
		"`COLUMN_NAME` AS `column_name`",
		"`COLLATION` AS `collation`",
		`CASE
			WHEN NULLABLE = 'YES' THEN true
			WHEN NULLABLE = "NO" THEN false
			ELSE false
		END AS ` + "`nullable`",
		`CASE
			WHEN NON_UNIQUE = 0 THEN true
			WHEN NON_UNIQUE = 1 THEN false
			ELSE 0
		END AS ` + "`unique`",
		"`COMMENT` AS `comment`",
		"`INDEX_TYPE` AS `index_type`",
		"`SEQ_IN_INDEX` AS `seq_in_index`",
		"`INDEX_COMMENT` AS `index_comment`",
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
func (grammarSQL SQL) GetColumnListing(dbName string, tableName string, db *sqlx.DB) ([]*grammar.Column, error) {
	selectColumns := []string{
		"TABLE_SCHEMA AS `db_name`",
		"TABLE_NAME AS `table_name`",
		"COLUMN_NAME AS `name`",
		"ORDINAL_POSITION AS `position`",
		"COLUMN_DEFAULT AS `default`",
		`CASE
			WHEN IS_NULLABLE = 'YES' THEN true
			WHEN IS_NULLABLE = "NO" THEN false
			ELSE false
		END AS ` + "`nullable`",
		`CASE
		   WHEN LOCATE('unsigned', COLUMN_TYPE) THEN true
		   ELSE false
		END AS` + "`unsigned`",
		"UPPER(DATA_TYPE) as `type`",
		"CHARACTER_MAXIMUM_LENGTH as `length`",
		"CHARACTER_OCTET_LENGTH as `octet_length`",
		"NUMERIC_PRECISION as `precision`",
		"NUMERIC_SCALE as `scale`",
		"DATETIME_PRECISION as `datetime_precision`",
		"CHARACTER_SET_NAME as `charset`",
		"COLLATION_NAME as `collation`",
		"COLUMN_KEY as `key`",
		`CASE
			WHEN COLUMN_KEY = 'PRI' THEN true
			ELSE false
		END AS ` + "`primary`",
		"EXTRA as `extra`",
		"COLUMN_COMMENT as `comment`",
	}
	sql := fmt.Sprintf(`
			SELECT %s
			FROM INFORMATION_SCHEMA.COLUMNS
			WHERE TABLE_SCHEMA = %s AND TABLE_NAME = %s;
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

// Create a new table on the schema
func (grammarSQL SQL) Create(table *grammar.Table, db *sqlx.DB) error {
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

	// indexes
	for _, index := range indexes {
		stmts = append(stmts,
			grammarSQL.Builder.SQLAddIndex(db, index, grammarSQL.IndexTypes, grammarSQL.Quoter),
		)
	}

	engine := utils.GetIF(table.Engine != "", "ENGINE "+table.Engine, "")
	charset := utils.GetIF(table.Charset != "", "DEFAULT CHARSET "+table.Charset, "")
	collation := utils.GetIF(table.Collation != "", "COLLATE="+table.Collation, "")

	sql = sql + strings.Join(stmts, ",\n")
	sql = sql + fmt.Sprintf(
		"\n) %s %s %s",
		engine, charset, collation,
	)
	defer logger.Debug(logger.CREATE, sql).TimeCost(time.Now())
	_, err := db.Exec(sql)
	return err
}

// Drop a table from the schema.
func (grammarSQL SQL) Drop(name string, db *sqlx.DB) error {
	sql := fmt.Sprintf("DROP TABLE %s", grammarSQL.Quoter.ID(name, db))
	defer logger.Debug(logger.DELETE, sql).TimeCost(time.Now())
	_, err := db.Exec(sql)
	return err
}

// DropIfExists if the table exists, drop it from the schema.
func (grammarSQL SQL) DropIfExists(name string, db *sqlx.DB) error {
	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s", grammarSQL.Quoter.ID(name, db))
	defer logger.Debug(logger.DELETE, sql).TimeCost(time.Now())
	_, err := db.Exec(sql)
	return err
}

// Rename a table on the schema.
func (grammarSQL SQL) Rename(old string, new string, db *sqlx.DB) error {
	sql := grammarSQL.Builder.SQLRenameTable(db, old, new, grammarSQL.Quoter)
	defer logger.Debug(logger.UPDATE, sql).TimeCost(time.Now())
	_, err := db.Exec(sql)
	return err
}

// Alter a table on the schema
func (grammarSQL SQL) Alter(table *grammar.Table, db *sqlx.DB) error {

	err := grammarSQL.Get(table, db)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf("ALTER TABLE %s ", grammarSQL.Quoter.ID(table.Name, db))
	stmts := []string{}
	errs := []error{}

	// Commands
	// The commands must be:
	//    AddColumn(column *Column)    for adding a column
	//    ModifyColumn(column *Column) for modifying a colu
	//    RenameColumn(old string,new string)  for renaming a column
	//    DropColumn(name string)  for dropping a column
	//    CreateIndex(index *Index) for creating a index
	//    DropIndex(name string) for  dropping a index
	//    RenameIndex(old string,new string)  for renaming a index
	for _, command := range table.Commands {
		switch command.Name {
		case "AddColumn", "ModifyColumn":
			column := command.Params[0].(*grammar.Column)
			stmt := ""
			if table.HasColumn(column.Name) {
				stmt = "MODIFY " + grammarSQL.Builder.SQLAddColumn(db, column, grammarSQL.Types, grammarSQL.Quoter)
			} else {
				stmt = "ADD " + grammarSQL.Builder.SQLAddColumn(db, column, grammarSQL.Types, grammarSQL.Quoter)
			}
			stmts = append(stmts, sql+stmt)
			err := grammarSQL.ExecSQL(db, table, sql+stmt)
			if err != nil {
				errs = append(errs, err)
			}
			break
		case "RenameColumn":
			old := command.Params[0].(string)
			new := command.Params[1].(string)
			column, has := table.ColumnMap[old]
			if !has {
				return errors.New("the column " + old + " not exists")
			}
			column.Name = new
			stmt := fmt.Sprintf("CHANGE COLUMN %s %s",
				grammarSQL.Quoter.ID(old, db),
				grammarSQL.Builder.SQLAddColumn(db, column, grammarSQL.Types, grammarSQL.Quoter),
			)
			stmts = append(stmts, sql+stmt)
			err := grammarSQL.ExecSQL(db, table, sql+stmt)
			if err != nil {
				errs = append(errs, err)
			}
			break
		case "DropColumn":
			name := command.Params[0].(string)
			stmt := fmt.Sprintf("DROP COLUMN %s", grammarSQL.Quoter.ID(name, db))
			stmts = append(stmts, sql+stmt)
			err := grammarSQL.ExecSQL(db, table, sql+stmt)
			if err != nil {
				errs = append(errs, err)
			}
			break
		case "CreateIndex":
			index := command.Params[0].(*grammar.Index)
			stmt := "ADD " + grammarSQL.Builder.SQLAddIndex(db, index, grammarSQL.IndexTypes, grammarSQL.Quoter)
			stmts = append(stmts, sql+stmt)
			err := grammarSQL.ExecSQL(db, table, sql+stmt)
			if err != nil {
				errs = append(errs, err)
			}
			break
		case "DropIndex":
			name := command.Params[0].(string)
			stmt := fmt.Sprintf("DROP INDEX %s", grammarSQL.Quoter.ID(name, db))
			stmts = append(stmts, sql+stmt)
			err := grammarSQL.ExecSQL(db, table, sql+stmt)
			if err != nil {
				errs = append(errs, err)
			}
			break
		case "RenameIndex":
			old := command.Params[0].(string)
			new := command.Params[1].(string)
			stmt := fmt.Sprintf("RENAME INDEX %s TO %s", grammarSQL.Quoter.ID(old, db), grammarSQL.Quoter.ID(new, db))
			stmts = append(stmts, sql+stmt)
			err := grammarSQL.ExecSQL(db, table, sql+stmt)
			if err != nil {
				errs = append(errs, err)
			}
			break
		}
	}

	defer logger.Debug(logger.CREATE, strings.Join(stmts, "\n")).TimeCost(time.Now())

	// Return Errors
	if len(errs) > 0 {
		message := ""
		for _, err := range errs {
			message = message + err.Error() + "\n"
		}
		return errors.New(message)
	}

	return nil
}

// ExecSQL execute sql then update table structure
func (grammarSQL SQL) ExecSQL(db *sqlx.DB, table *grammar.Table, sql string) error {
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	// update table structure
	err = grammarSQL.Get(table, db)
	if err != nil {
		return err
	}
	return nil
}

// ParseType parse type and flip to DBAL
func (grammarSQL SQL) ParseType(column *grammar.Column) {
	typeinfo := strings.Split(strings.ToUpper(column.Type), " ")
	re := regexp.MustCompile(`([A-Z]+)[\(]*([0-9,]*)[\)]*`)
	matched := re.FindStringSubmatch(typeinfo[0])
	if len(matched) == 3 {
		typeName := matched[1]
		typeArgs := strings.Trim(matched[2], " ")
		args := []string{}
		if typeArgs != "" {
			args = strings.Split(strings.Trim(matched[2], " "), ",")
		}
		typ, has := grammarSQL.FlipTypes[typeName]
		if has {
			column.Type = typ
		}
		switch column.Type {
		case "bigInteger", "integer":
			if len(args) > 0 {
				precision, err := strconv.Atoi(args[0])
				if err == nil {
					column.Precision = utils.IntPtr(precision)
				}
			} else if column.IsUnsigned {
				column.Precision = utils.IntPtr(20)
			} else {
				column.Precision = utils.IntPtr(19)
			}
			break
		case "timestamp":
			if len(args) > 0 {
				precision, err := strconv.Atoi(args[0])
				if err == nil {
					column.DatetimePrecision = utils.IntPtr(precision)
				}
			}
			break
		case "float":
			if len(args) > 0 {
				precision, err := strconv.Atoi(args[0])
				if err == nil {
					column.Precision = utils.IntPtr(precision)
				}

				if len(args) > 1 {
					scale, err := strconv.Atoi(args[1])
					if err == nil {
						column.Scale = utils.IntPtr(scale)
					}
				}
			}
			break
		case "string", "text":
			if len(args) > 0 {
				length, err := strconv.Atoi(args[0])
				if err == nil {
					column.Length = utils.IntPtr(length)
				}
			}
			break
		}
	}
	// utils.Println(column)
}
