package sql

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/logger"
	"github.com/yaoapp/xun/utils"
)

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
			grammarSQL.Builder.SQLCreateColumn(db, Column, grammarSQL.Types, grammarSQL.Quoter),
		)
	}

	// indexes
	for _, index := range indexes {
		stmts = append(stmts,
			grammarSQL.Builder.SQLCreateIndex(db, index, grammarSQL.IndexTypes, grammarSQL.Quoter),
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
	// sql := `SELECT xxx
	// 	FROM INFORMATION_SCHEMA.COLUMNS
	// 	WHERE TABLE_SCHEMA = 'xxx' AND TABLE_NAME ='xxx'
	// `
	name := grammarSQL.Quoter.ID(table.Name, db)
	sql := fmt.Sprintf("ALTER TABLE %s \n", name)
	stmts := []string{}

	// Columns
	for _, Column := range table.Columns {
		if Column.Dropped { // Drop
			stmts = append(stmts, fmt.Sprintf("DROP COLUMN %s", grammarSQL.Quoter.ID(Column.Name, db)))
			// } else if Column.Exists { // Modify
			// 	stmts = append(stmts, fmt.Sprintf("MODIFY COLUMN %s column_definition", grammarSQL.Quoter.ID(Column.Name, db)))
		} else if Column.Newname != "" { // Rename
			stmts = append(stmts, fmt.Sprintf("RENAME COLUMN %s TO %s", grammarSQL.Quoter.ID(Column.Name, db), grammarSQL.Quoter.ID(Column.Newname, db)))
		} else { // ADD
			stmts = append(stmts,
				"ADD "+grammarSQL.Builder.SQLCreateColumn(db, Column, grammarSQL.Types, grammarSQL.Quoter),
			)
		}
	}
	sql = sql + strings.Join(stmts, ",\n")
	defer logger.Debug(logger.UPDATE, sql).TimeCost(time.Now())
	return nil
}
