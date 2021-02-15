package sqlite3

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/logger"
	"github.com/yaoapp/xun/utils"
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
		// "CHARACTER_MAXIMUM_LENGTH as `length`",
		// "CHARACTER_OCTET_LENGTH as `octet_length`",
		// "NUMERIC_PRECISION as `precision`",
		// "NUMERIC_SCALE as `scale`",
		// "DATETIME_PRECISION as `datetime_precision`",
		// "CHARACTER_SET_NAME as `charset`",
		// "COLLATION_NAME as `collation`",
		// "COLUMN_KEY as `key`",
		// "COLUMN_COMMENT as `comment`",
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
	}
	return columns, nil
}

// ParseType parse type and flip to DBAL
func (grammarSQL SQLite3) ParseType(column *grammar.Column) {
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
				column.Precision = utils.IntPtr(19)
			} else {
				column.Precision = utils.IntPtr(20)
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

	// fmt.Printf(
	// 	"ParseType %s: %s %d (%d,%d) %d\n",
	// 	column.Name, column.Type,
	// 	utils.IntVal(column.Length),
	// 	utils.IntVal(column.Precision),
	// 	utils.IntVal(column.Scale),
	// 	utils.IntVal(column.DatetimePrecision),
	// )
}
