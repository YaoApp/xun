package sqlite3

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/logger"
)

// Create a new table on the schema
func (grammarSQL SQLite3) Create(table *grammar.Table, db *sqlx.DB) error {
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
func (grammarSQL SQLite3) Get(table *grammar.Table, db *sqlx.DB) error {
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
func (grammarSQL SQLite3) GetIndexListing(dbName string, tableName string, db *sqlx.DB) ([]*grammar.Index, error) {
	selectColumns := []string{
		"m.`tbl_name` AS `table_name`",
		"il.`name` AS `index_name`",
		"ii.`name` AS `column_name`",
		`CASE 
			WHEN il.origin = 'pk' then 'primary' 
			WHEN il.[unique] = 1  THEN 'unique'
			ELSE 'index'
		END as index_type`,
		`CASE
			WHEN il.[unique] = 1 THEN 1
			WHEN il.[unique] = 0 THEN 0
			ELSE 0
		END AS ` + "`unique`",
		"il.`seq`  AS `seq_in_index`",
		"ii.`seqno` AS  `seq_in_column`",
	}

	sql := fmt.Sprintf(`
			SELECT %s
				FROM sqlite_master AS m,
				pragma_index_list(m.name) AS il,
				pragma_index_info(il.name) AS ii
			WHERE 
				m.type = 'table'
				and m.tbl_name = %s
			GROUP BY
				m.tbl_name,
				il.name,
				ii.name,
				il.origin,
				il.partial,
				il.seq
			UNION
			SELECT 
				%s as table_name, 
				'PRIMARY' as index_name, 
				ti.name as column_name,
				"primary" as index_type,
				1 as `+"`unique`"+`,
				0 as `+"`seq_in_index`"+`,
				0 as `+"`seq_in_column`"+`
			FROM pragma_table_info(%s) AS ti WHERE ti.pk=1
			ORDER BY seq_in_index,index_name,seq_in_column
		`,
		strings.Join(selectColumns, ","),
		grammarSQL.Quoter.VAL(tableName, db),
		grammarSQL.Quoter.VAL(tableName, db),
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
		index.Nullable = true
		index.DBName = dbName
		index.Type = index.IndexType
		// utils.Println(index)
	}
	return indexes, nil
}

// GetColumnListing get a table columns structure
func (grammarSQL SQLite3) GetColumnListing(dbName string, tableName string, db *sqlx.DB) ([]*grammar.Column, error) {
	selectColumns := []string{
		"m.name AS `table_name`",
		"p.name AS `name`",
		"p.cid AS `position`",
		"p.dflt_value AS `default`",
		"UPPER(p.type) as `type`",
		`CASE
			WHEN ` + "p.`notnull`" + ` == 0 THEN 1
			ELSE 0
		END AS ` + "`nullable`",
		`CASE
		   WHEN INSTR(` + "p.`type`" + `, 'UNSIGNED' ) THEN 1
		   WHEN  p.pk = 1 and ` + "p.`type`" + ` = 'INTEGER' THEN 1 
		   ELSE 0
		END AS` + "`unsigned`",
		`CASE
			WHEN p.pk = 1 THEN 1
			ELSE 0
		END AS ` + "`primary`",
		`CASE
			WHEN p.pk = 1 and INSTR(m.sql, 'AUTOINCREMENT' ) THEN "AutoIncrement"
			ELSE ""
		END AS ` + "`extra`",
	}
	sql := fmt.Sprintf(`
			SELECT %s
			FROM sqlite_master m
			LEFT OUTER JOIN pragma_table_info((m.name)) p  ON m.name <> p.name
			WHERE m.type = 'table' and table_name=%s
		`,
		strings.Join(selectColumns, ","),
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
		grammarSQL.ParseType(column)
		column.DBName = dbName
	}
	return columns, nil
}

// Alter a table on the schema
func (grammarSQL SQLite3) Alter(table *grammar.Table, db *sqlx.DB) error {

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
				logger.Warn(logger.CREATE, "sqlite3 not support ModifyColumn operation").Write()
				break
			}
			stmt = sql + "ADD COLUMN " + grammarSQL.Builder.SQLAddColumn(db, column, grammarSQL.Types, grammarSQL.Quoter)
			stmts = append(stmts, stmt)
			err := grammarSQL.ExecSQL(db, table, stmt)
			if err != nil {
				errs = append(errs, errors.New("SQL: "+stmt+" ERROR: "+err.Error()))
			}
			break
		case "RenameColumn":
			old := command.Params[0].(string)
			new := command.Params[1].(string)
			stmt := fmt.Sprintf("%s RENAME COLUMN %s TO %s",
				sql,
				grammarSQL.Quoter.ID(old, db),
				grammarSQL.Quoter.ID(new, db),
			)
			stmts = append(stmts, stmt)
			err := grammarSQL.ExecSQL(db, table, stmt)
			if err != nil {
				errs = append(errs, errors.New("SQL: "+stmt+" ERROR: "+err.Error()))
			}
			break
		case "DropColumn":
			logger.Warn(logger.CREATE, "sqlite3 not support DropColumn operation").Write()
			break
		case "CreateIndex":
			index := command.Params[0].(*grammar.Index)
			stmt := grammarSQL.Builder.SQLAddIndex(db, index, grammarSQL.IndexTypes, grammarSQL.Quoter)
			stmts = append(stmts, stmt)
			err := grammarSQL.ExecSQL(db, table, stmt)
			if err != nil {
				errs = append(errs, errors.New("SQL: "+stmt+" ERROR: "+err.Error()))
			}
			break
		case "DropIndex":
			name := command.Params[0].(string)
			stmt := fmt.Sprintf("DROP INDEX IF EXISTS %s", grammarSQL.Quoter.ID(name, db))
			stmts = append(stmts, stmt)
			err := grammarSQL.ExecSQL(db, table, stmt)
			if err != nil {
				errs = append(errs, errors.New("SQL: "+stmt+" ERROR: "+err.Error()))
			}
			break
		case "RenameIndex":
			logger.Warn(logger.CREATE, "sqlite3 not support RenameIndex operation").Write()
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
func (grammarSQL SQLite3) ExecSQL(db *sqlx.DB, table *grammar.Table, sql string) error {
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
