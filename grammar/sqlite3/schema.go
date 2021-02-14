package sqlite3

import (
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
