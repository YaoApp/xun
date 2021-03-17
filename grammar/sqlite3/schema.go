package sqlite3

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/blang/semver/v4"
	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/logger"
	"github.com/yaoapp/xun/utils"
)

// GetVersion get the version of the connection database
func (grammarSQL SQLite3) GetVersion(db *sqlx.DB) (*dbal.Version, error) {
	sql := fmt.Sprintf("SELECT SQLITE_VERSION()")
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	rows := []string{}
	err := db.Select(&rows, sql)
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("Can't get the version")
	}

	ver, err := semver.Make(rows[0])
	if err != nil {
		return nil, err
	}

	return &dbal.Version{
		Version: ver,
		Driver:  grammarSQL.Driver,
	}, nil
}

// Exists the Exists
func (grammarSQL SQLite3) Exists(name string, db *sqlx.DB) bool {
	sql := fmt.Sprintf("SELECT `name` FROM `sqlite_master` WHERE type='table' AND name=%s", grammarSQL.Quoter.VAL(name, db))
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
func (grammarSQL SQLite3) Create(table *dbal.Table, db *sqlx.DB) error {

	name := grammarSQL.Quoter.ID(table.TableName, db)
	sql := fmt.Sprintf("CREATE TABLE %s (\n", name)
	stmts := []string{}

	var primary *dbal.Primary = nil
	columns := []*dbal.Column{}
	indexes := []*dbal.Index{}
	cbCommands := []*dbal.Command{}

	// Commands
	// The commands must be:
	//    AddColumn(column *Column)    for adding a column
	//    ChangeColumn(column *Column) for modifying a colu
	//    RenameColumn(old string,new string)  for renaming a column
	//    DropColumn(name string)  for dropping a column
	//    CreateIndex(index *Index) for creating a index
	//    DropIndex( name string) for  dropping a index
	//    RenameIndex(old string,new string)  for renaming a index
	for _, command := range table.Commands {
		switch command.Name {
		case "AddColumn":
			columns = append(columns, command.Params[0].(*dbal.Column))
			cbCommands = append(cbCommands, command)
			break
		case "CreateIndex":
			indexes = append(indexes, command.Params[0].(*dbal.Index))
			cbCommands = append(cbCommands, command)
			break
		case "CreatePrimary":
			primary = command.Params[0].(*dbal.Primary)
			cbCommands = append(cbCommands, command)
		}
	}

	// Columns
	for _, Column := range columns {
		stmts = append(stmts,
			grammarSQL.SQLAddColumn(db, Column),
		)
	}

	// Primary key
	if primary != nil && len(primary.Columns) > 1 {
		stmts = append(stmts,
			grammarSQL.SQLAddPrimary(db, primary),
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
			grammarSQL.SQLAddIndex(db, index),
		)
	}
	defer logger.Debug(logger.CREATE, indexStmts...).TimeCost(time.Now())
	_, err = db.Exec(strings.Join(indexStmts, ";\n"))

	for _, cmd := range cbCommands {
		cmd.Callback(err)
	}

	if err != nil {
		return err
	}

	return nil
}

// Rename a table on the schema.
func (grammarSQL SQLite3) Rename(old string, new string, db *sqlx.DB) error {
	sql := fmt.Sprintf("ALTER TABLE %s RENAME TO %s", grammarSQL.Quoter.ID(old, db), grammarSQL.Quoter.ID(new, db))
	defer logger.Debug(logger.UPDATE, sql).TimeCost(time.Now())
	_, err := db.Exec(sql)
	return err
}

// Get a table on the schema
func (grammarSQL SQLite3) Get(table *dbal.Table, db *sqlx.DB) error {
	columns, err := grammarSQL.GetColumnListing(table.DBName, table.TableName, db)
	if err != nil {
		return err
	}

	indexes, err := grammarSQL.GetIndexListing(table.DBName, table.TableName, db)
	if err != nil {
		return err
	}

	primaryKeyName := ""

	// attaching columns
	for _, column := range columns {
		column.Indexes = []*dbal.Index{}
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
			index.Columns = []*dbal.Column{}
			column.Indexes = append(column.Indexes, &index)
			table.PushIndex(&index)
		}
		index := table.IndexMap[idx.Name]
		index.Columns = append(index.Columns, column)

		if index.Type == "primary" {
			primaryKeyName = idx.Name
		}
	}

	// attaching primary
	if _, has := table.IndexMap[primaryKeyName]; has {
		idx := table.IndexMap[primaryKeyName]
		table.Primary = &dbal.Primary{
			Name:      idx.Name,
			TableName: idx.TableName,
			DBName:    idx.DBName,
			Table:     idx.Table,
			Columns:   idx.Columns,
		}
		delete(table.IndexMap, idx.Name)
		for _, column := range table.Primary.Columns {
			column.Primary = true
			column.Indexes = []*dbal.Index{}
		}
	}

	return nil
}

// GetIndexListing get a table indexes structure
func (grammarSQL SQLite3) GetIndexListing(dbName string, tableName string, db *sqlx.DB) ([]*dbal.Index, error) {
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
	indexes := []*dbal.Index{}
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
func (grammarSQL SQLite3) GetColumnListing(schemaName string, tableName string, db *sqlx.DB) ([]*dbal.Column, error) {
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
	columns := []*dbal.Column{}
	err := db.Select(&columns, sql)
	if err != nil {
		return nil, err
	}

	// Get the table Constraints
	constraints, err := grammarSQL.GetConstraintListing(db, schemaName, tableName)
	if err != nil {
		return nil, err
	}

	// Cast the database data type to DBAL data type
	for _, column := range columns {
		grammarSQL.ParseType(column)
		column.DBName = schemaName
		constraint, has := constraints[column.Name]
		if has {
			column.Constraint = constraint

			// enum
			if column.Type == "text" && constraint.Type == "CHECK" && len(constraint.Args) >= 1 && strings.Contains(constraint.Args[0], "IN ('") {
				re := regexp.MustCompile(`IN \('(.*)'\)`)
				matched := re.FindStringSubmatch(constraint.Args[0])
				if len(matched) == 2 {
					options := strings.Split(matched[1], "','")
					column.Option = options
					column.Type = "enum"
				}
			}
		}

	}
	return columns, nil
}

// Alter a table on the schema
func (grammarSQL SQLite3) Alter(table *dbal.Table, db *sqlx.DB) error {

	sql := fmt.Sprintf("ALTER TABLE %s ", grammarSQL.Quoter.ID(table.TableName, db))
	stmts := []string{}
	errs := []error{}

	// Commands
	// The commands must be:
	//    AddColumn(column *Column)    for adding a column
	//    ChangeColumn(column *Column) for modifying a colu
	//    RenameColumn(old string,new string)  for renaming a column
	//    DropColumn(name string)  for dropping a column
	//    CreateIndex(index *Index) for creating a index
	//    DropIndex(name string) for  dropping a index
	//    RenameIndex(old string,new string)  for renaming a index
	for _, command := range table.Commands {
		switch command.Name {
		case "AddColumn":
			column := command.Params[0].(*dbal.Column)
			stmt := ""
			stmt = sql + "ADD COLUMN " + grammarSQL.SQLAddColumn(db, column)
			stmts = append(stmts, stmt)
			err := grammarSQL.ExecSQL(db, table, stmt)
			if err != nil {
				errs = append(errs, errors.New("SQL: "+stmt+" ERROR: "+err.Error()))
			}
			command.Callback(err)
			break
		case "ChangeColumn":
			logger.Warn(logger.CREATE, "sqlite3 not support ChangeColumn operation").Write()
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
			command.Callback(err)
			break
		case "DropColumn":
			logger.Warn(logger.CREATE, "sqlite3 not support DropColumn operation").Write()
			break
		case "CreateIndex":
			index := command.Params[0].(*dbal.Index)
			stmt := grammarSQL.SQLAddIndex(db, index)
			stmts = append(stmts, stmt)
			err := grammarSQL.ExecSQL(db, table, stmt)
			if err != nil {
				errs = append(errs, errors.New("SQL: "+stmt+" ERROR: "+err.Error()))
			}
			command.Callback(err)
			break
		case "DropIndex":
			name := command.Params[0].(string)
			stmt := fmt.Sprintf("DROP INDEX IF EXISTS %s", grammarSQL.Quoter.ID(name, db))
			stmts = append(stmts, stmt)
			err := grammarSQL.ExecSQL(db, table, stmt)
			if err != nil {
				errs = append(errs, errors.New("SQL: "+stmt+" ERROR: "+err.Error()))
			}
			command.Callback(err)
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

// GetConstraintListing get the constraints of the table
func (grammarSQL SQLite3) GetConstraintListing(db *sqlx.DB, schemaName string, tableName string) (map[string]*dbal.Constraint, error) {
	rows := []string{}
	err := db.Select(&rows, "SELECT `sql` FROM sqlite_master WHERE type='table' and name=?", tableName)
	if err != nil {
		return nil, err
	}

	if len(rows) < 1 {
		return nil, fmt.Errorf("the table %s does not exists", tableName)
	}

	sql := rows[0]
	lines := strings.Split(sql, "\n")
	constraints := map[string]*dbal.Constraint{}
	for _, line := range lines {
		constraint := grammarSQL.parseConstraint(schemaName, tableName, line)
		if constraint != nil {
			constraints[constraint.ColumnName] = constraint
		}
	}
	return constraints, nil
}

func (grammarSQL SQLite3) parseConstraint(schemaName string, tableName string, line string) *dbal.Constraint {
	// fmt.Printf("GetConstraintListing Line: %#v\n", line)
	if strings.Contains(line, "CHECK(") {
		re := regexp.MustCompile("`([0-9a-zA-Z_]+)` .* CHECK\\((.*)\\)")
		matched := re.FindStringSubmatch(line)
		if len(matched) == 3 {
			column := strings.Trim(matched[1], " ")
			exp := strings.Trim(matched[2], " ")
			constraint := dbal.NewConstraint(schemaName, tableName, column)
			constraint.Type = "CHECK"
			constraint.Args = append(constraint.Args, exp)
			return constraint
		}
	} else if strings.Contains(line, "PRIMARY KEY") {

	} else if strings.Contains(line, "UNIQUE") {

	} else if strings.Contains(line, "FOREIGN KEY") {

	}
	return nil
}

// ExecSQL execute sql then update table structure
func (grammarSQL SQLite3) ExecSQL(db *sqlx.DB, table *dbal.Table, sql string) error {
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
func (grammarSQL SQLite3) ParseType(column *dbal.Column) {

	re := regexp.MustCompile(`([A-Z ]+)[\(]*([0-9,]*)[\)]*`)
	matched := re.FindStringSubmatch(strings.ToUpper(column.Type))
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
					column.DateTimePrecision = utils.IntPtr(precision)
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
