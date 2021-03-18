package postgres

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/blang/semver/v4"
	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/logger"
	"github.com/yaoapp/xun/utils"
)

// GetVersion get the version of the connection database
func (grammarSQL Postgres) GetVersion(db *sqlx.DB) (*dbal.Version, error) {
	sql := fmt.Sprintf("SELECT VERSION()")
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	rows := []string{}
	err := db.Select(&rows, sql)
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("Can't get the version")
	}

	verArr := strings.Split(rows[0], " ")
	if len(verArr) < 2 {
		return nil, fmt.Errorf("Can't parse the version: %s", rows[0])
	}
	ver, err := semver.Make(verArr[1])
	if err != nil {
		return nil, err
	}

	return &dbal.Version{
		Version: ver,
		Driver:  grammarSQL.Driver,
	}, nil
}

// Exists the Exists
func (grammarSQL Postgres) Exists(name string, db *sqlx.DB) bool {
	sql := fmt.Sprintf("SELECT table_name AS name FROM information_schema.tables WHERE table_name = %s", grammarSQL.Quoter.VAL(name, db))
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

// CreateType create user defined type
func (grammarSQL Postgres) CreateType(table *dbal.Table, db *sqlx.DB, types map[string][]string) error {
	// Create Types
	for name, option := range types {
		typ := fmt.Sprintf("ENUM('%s')", strings.Join(option, "','"))
		typeSQL := fmt.Sprintf(`
	DO $$ BEGIN
		CREATE TYPE %s.%s AS %s; 
	EXCEPTION  
		WHEN duplicate_object THEN null;
	END $$;
	`, table.SchemaName, name, typ)
		defer logger.Debug(logger.CREATE, typeSQL).TimeCost(time.Now())
		_, err := db.Exec(typeSQL)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create a new table on the schema
func (grammarSQL Postgres) Create(table *dbal.Table, db *sqlx.DB) error {
	name := grammarSQL.Quoter.ID(table.TableName, db)
	sql := fmt.Sprintf("CREATE TABLE %s (\n", name)
	stmts := []string{}
	commentStmts := []string{}

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
			break
		}
	}

	// Enum types
	types := map[string][]string{}

	// Columns
	for _, column := range columns {
		stmts = append(stmts,
			grammarSQL.SQLAddColumn(db, column),
		)

		commentStmt := grammarSQL.SQLAddComment(db, column)
		if commentStmt != "" {
			commentStmts = append(commentStmts, commentStmt)
		}
		if column.Type == "enum" {
			typeName := strings.ToLower("ENUM__" + strings.Join(column.Option, "_EOPT_"))
			types[typeName] = column.Option
		}
	}

	// Create Types
	err := grammarSQL.CreateType(table, db, types)
	if err != nil {
		return err
	}

	// Primary key
	if primary != nil {
		stmts = append(stmts,
			grammarSQL.SQLAddPrimary(db, primary),
		)
	}
	sql = sql + strings.Join(stmts, ",\n")
	sql = sql + fmt.Sprintf("\n)")

	// Create table
	defer logger.Debug(logger.CREATE, sql).TimeCost(time.Now())
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	// indexes
	indexStmts := []string{}
	for _, index := range indexes {
		if index.Type == "primary" {
			continue
		}
		indexStmt := grammarSQL.SQLAddIndex(db, index)
		if indexStmt != "" {
			indexStmts = append(indexStmts, indexStmt)
		}
	}
	if len(indexStmts) > 0 {
		defer logger.Debug(logger.CREATE, indexStmts...).TimeCost(time.Now())
		_, err = db.Exec(strings.Join(indexStmts, ";\n"))
	}

	// Comments
	if len(commentStmts) > 0 {
		defer logger.Debug(logger.CREATE, commentStmts...).TimeCost(time.Now())
		_, err = db.Exec(strings.Join(commentStmts, ";\n"))
	}

	// Callback
	for _, cmd := range cbCommands {
		cmd.Callback(err)
	}

	if err != nil {
		return err
	}

	return nil
}

// Rename a table on the schema.
func (grammarSQL Postgres) Rename(old string, new string, db *sqlx.DB) error {
	sql := fmt.Sprintf("ALTER TABLE %s RENAME TO %s", grammarSQL.Quoter.ID(old, db), grammarSQL.Quoter.ID(new, db))
	defer logger.Debug(logger.UPDATE, sql).TimeCost(time.Now())
	_, err := db.Exec(sql)
	return err
}

// Get a table on the schema
func (grammarSQL Postgres) Get(table *dbal.Table, db *sqlx.DB) error {
	columns, err := grammarSQL.GetColumnListing(table.SchemaName, table.TableName, db)
	if err != nil {
		return err
	}
	indexes, err := grammarSQL.GetIndexListing(table.SchemaName, table.TableName, db)
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

// Alter a table on the schema
func (grammarSQL Postgres) Alter(table *dbal.Table, db *sqlx.DB) error {

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
			stmt := "ADD COLUMN " + grammarSQL.SQLAddColumn(db, column)
			stmts = append(stmts, sql+stmt)
			if column.Type == "enum" {
				typeName := strings.ToLower("ENUM__" + strings.Join(column.Option, "_EOPT_"))
				types := map[string][]string{}
				types[typeName] = column.Option
				err := grammarSQL.CreateType(table, db, types)
				if err != nil {
					return err
				}
			}
			err := grammarSQL.ExecSQL(db, table, sql+stmt)
			if err != nil {
				errs = append(errs, err)
			}

			commentStmt := grammarSQL.SQLAddComment(db, column)
			if commentStmt != "" {
				err := grammarSQL.ExecSQL(db, table, commentStmt)
				if err != nil {
					errs = append(errs, err)
				}
			}

			command.Callback(err)
			break
		case "ChangeColumn":
			column := command.Params[0].(*dbal.Column)
			stmt := "ALTER COLUMN " + grammarSQL.SQLAlterColumnType(db, column)
			stmts = append(stmts, sql+stmt)
			if column.Type == "enum" {
				typeName := strings.ToLower("ENUM__" + strings.Join(column.Option, "_EOPT_"))
				types := map[string][]string{}
				types[typeName] = column.Option
				err := grammarSQL.CreateType(table, db, types)
				if err != nil {
					return err
				}
			}
			err := grammarSQL.ExecSQL(db, table, sql+stmt)
			if err != nil {
				errs = append(errs, err)
			}

			commentStmt := grammarSQL.SQLAddComment(db, column)
			if commentStmt != "" {
				err := grammarSQL.ExecSQL(db, table, commentStmt)
				if err != nil {
					errs = append(errs, err)
				}
			}

			command.Callback(err)
			break
		case "RenameColumn":
			old := command.Params[0].(string)
			new := command.Params[1].(string)
			column, has := table.ColumnMap[old]
			if !has {
				return errors.New("the column " + old + " not exists")
			}
			column.Name = new
			stmt := fmt.Sprintf("RENAME COLUMN %s TO %s",
				grammarSQL.Quoter.ID(old, db),
				grammarSQL.Quoter.ID(new, db),
			)
			stmts = append(stmts, sql+stmt)
			err := grammarSQL.ExecSQL(db, table, sql+stmt)
			if err != nil {
				errs = append(errs, err)
			}
			command.Callback(err)
			break
		case "DropColumn":
			name := command.Params[0].(string)
			stmt := fmt.Sprintf("DROP COLUMN %s", grammarSQL.Quoter.ID(name, db))
			stmts = append(stmts, sql+stmt)
			err := grammarSQL.ExecSQL(db, table, sql+stmt)
			if err != nil {
				errs = append(errs, err)
			}
			command.Callback(err)
			break
		case "CreateIndex":
			index := command.Params[0].(*dbal.Index)
			stmt := grammarSQL.SQLAddIndex(db, index)
			stmts = append(stmts, stmt)
			err := grammarSQL.ExecSQL(db, table, stmt)
			if err != nil {
				errs = append(errs, err)
			}
			break
		case "DropIndex":
			name := command.Params[0].(string)
			stmt := fmt.Sprintf(
				"DROP INDEX %s",
				grammarSQL.Quoter.ID(name, db),
				// grammarSQL.Quoter.ID(table.TableName, db),
			)
			stmts = append(stmts, stmt)
			err := grammarSQL.ExecSQL(db, table, stmt)
			if err != nil {
				errs = append(errs, err)
			}
			command.Callback(err)
			break
		case "RenameIndex":
			old := command.Params[0].(string)
			new := command.Params[1].(string)
			stmt := fmt.Sprintf(
				"ALTER INDEX IF EXISTS %s RENAME TO %s",
				grammarSQL.Quoter.ID(old, db),
				grammarSQL.Quoter.ID(new, db),
				// grammarSQL.Quoter.ID(table.TableName, db),
			)
			stmts = append(stmts, stmt)
			err := grammarSQL.ExecSQL(db, table, stmt)
			if err != nil {
				errs = append(errs, err)
			}
			command.Callback(err)
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
func (grammarSQL Postgres) ExecSQL(db *sqlx.DB, table *dbal.Table, sql string) error {
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

// SQLAlterColumnType return the add column sql for table alter
func (grammarSQL Postgres) SQLAlterColumnType(db *sqlx.DB, Column *dbal.Column) string {
	// `id` bigint(20) unsigned NOT NULL,
	types := grammarSQL.Types
	quoter := grammarSQL.Quoter

	typ, has := types[Column.Type]
	if !has {
		typ = "VARCHAR"
	}
	if Column.Precision != nil && Column.Scale != nil && (typ == "NUMBERIC" || typ == "DECIMAL") {
		typ = fmt.Sprintf("%s(%d,%d)", typ, utils.IntVal(Column.Precision), utils.IntVal(Column.Scale))
	} else if strings.Contains(typ, "TIMESTAMP(%d)") || strings.Contains(typ, "TIME(%d)") {
		DateTimePrecision := utils.IntVal(Column.DateTimePrecision, 0)
		typ = fmt.Sprintf(typ, DateTimePrecision)
	} else if typ == "BYTEA" {
		typ = "BYTEA"
	} else if typ == "ENUM" {
		typ = strings.ToLower("ENUM__" + strings.Join(Column.Option, "_EOPT_"))
	} else if Column.Length != nil {
		typ = fmt.Sprintf("%s(%d)", typ, utils.IntVal(Column.Length))
	} else if typ == "IPADDRESS" { // ipAddress
		typ = "integer"
	}

	if utils.StringVal(Column.Extra) != "" {
		if typ == "BIGINT" {
			typ = "BIGSERIAL"
		} else {
			typ = "SERIAL"
		}
	}

	// sql := fmt.Sprintf(
	// 	"%s SET DATA TYPE %s ",
	// 	quoter.ID(Column.Name, db), typ)

	nameQuoter := quoter.ID(Column.Name, db)
	sql := fmt.Sprintf(
		"%s TYPE %s USING (%s::%s) ",
		nameQuoter, typ, nameQuoter, typ)

	sql = strings.Trim(sql, " ")
	return sql
}

// SQLAlterIndex  return the add index sql for table alter
func (grammarSQL Postgres) SQLAlterIndex(db *sqlx.DB, index *dbal.Index) string {
	indexTypes := grammarSQL.IndexTypes
	quoter := grammarSQL.Quoter
	typ, has := indexTypes[index.Type]
	if !has {
		typ = "KEY"
	}

	// UNIQUE KEY `unionid` (`unionid`) COMMENT 'xxxx'
	columns := []string{}
	for _, Column := range index.Columns {
		columns = append(columns, quoter.ID(Column.Name, db))
	}

	name := quoter.ID(index.Name, db)
	sql := fmt.Sprintf(
		"CREATE %s %s ON %s (%s)",
		typ, name, quoter.ID(index.TableName, db), strings.Join(columns, ","))

	if typ == "PRIMARY KEY" {
		sql = fmt.Sprintf(
			"%s (%s) ",
			typ, strings.Join(columns, ","))
	}
	return sql
}

// GetIndexListing get a table indexes structure
func (grammarSQL Postgres) GetIndexListing(dbName string, tableName string, db *sqlx.DB) ([]*dbal.Index, error) {
	selectColumns := []string{
		"n.nspname as db_name",
		"t.relname as table_name",
		"i.relname as index_name",
		"a.attname as column_name",
		"'' as collation",
		"false as nullable",
		"indisunique as unique",
		`indisprimary as "primary"`,
		"'' as comment",
		"'BTREE' as index_type",
		"a.attnum as seq_in_index",
		"'' as index_comment",
	}
	sql := fmt.Sprintf(`
			SELECT %s 
			FROM
				pg_class t,pg_class i,pg_index ix,pg_attribute a,pg_type as typ,pg_namespace as n
			WHERE 
				t.oid = ix.indrelid
				and n.oid = t.relnamespace
				and i.oid = ix.indexrelid
				and	typ.oid = a.atttypid
				and a.attrelid = t.oid
				and a.attnum = ANY(ix.indkey)
				and t.relkind = 'r'
				and n.nspname = %s
				and t.relname = %s
			ORDER BY
				t.relname, i.relname,a.attnum desc
			`,
		strings.Join(selectColumns, ","),
		grammarSQL.Quoter.VAL(dbName, db),
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
		if index.Primary {
			index.Type = "primary"
			index.Name = "PRIMARY"
		} else if index.Unique {
			index.Type = "unique"
		} else {
			index.Type = "index"
		}
	}
	return indexes, nil
}

// GetColumnListing get a table columns structure
func (grammarSQL Postgres) GetColumnListing(dbName string, tableName string, db *sqlx.DB) ([]*dbal.Column, error) {
	selectColumns := []string{
		"TABLE_SCHEMA AS \"db_name\"",
		"TABLE_NAME AS \"table_name\"",
		"COLUMN_NAME AS \"name\"",
		"ORDINAL_POSITION AS \"position\"",
		"COLUMN_DEFAULT AS \"default\"",
		`CASE
			WHEN IS_NULLABLE = 'YES' THEN true
			WHEN IS_NULLABLE = 'NO' THEN false
			ELSE false
		END AS "nullable"`,
		`CASE
			WHEN (UDT_NAME ~ 'unsigned')  THEN true
			ELSE false
		END AS "unsigned"`,
		"UDT_NAME as \"type_name\"",
		"UPPER(DATA_TYPE) as \"type\"",
		"CHARACTER_MAXIMUM_LENGTH as \"length\"",
		"CHARACTER_OCTET_LENGTH as \"octet_length\"",
		"NUMERIC_PRECISION as \"precision\"",
		"NUMERIC_SCALE as \"scale\"",
		"DATETIME_PRECISION as \"datetime_precision\"",
		"CHARACTER_SET_NAME as \"charset\"",
		"COLLATION_NAME as \"collation\"",
		"'' as \"key\"",
		`false AS "primary"`,
		`CASE 
		 	WHEN (COLUMN_DEFAULT ~ 'nextval\(.*_seq') THEN 'auto_increment'
		 	ELSE ''
		END as "extra"`,
		"pg_catalog.col_description(format('%s.%s',table_schema,table_name)::regclass::oid,ordinal_position)  as \"comment\"",
	}
	sql := fmt.Sprintf(`
			SELECT %s
			FROM INFORMATION_SCHEMA.COLUMNS
			WHERE  TABLE_SCHEMA = %s AND TABLE_NAME = %s;
		`,
		strings.Join(selectColumns, ","),
		grammarSQL.Quoter.VAL(dbName, db),
		grammarSQL.Quoter.VAL(tableName, db),
	)
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	columns := []*dbal.Column{}
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

		if column.Comment != nil {
			typ = grammarSQL.GetTypeFromComment(column.Comment)
			if typ != "" {
				column.Type = typ
			}
		}

		// user defined types
		if column.Type == "USER-DEFINED" {

			// enum options
			enumOptions := map[string][]string{}
			if strings.Contains(column.TypeName, "enum__") {
				column.Type = "enum"
				if _, has := enumOptions[column.TypeName]; !has {
					optionRange := []string{}
					err := db.Select(&optionRange, fmt.Sprintf("select enum_range(null::%s.%s)", dbName, column.TypeName))
					if err != nil {
						return nil, err
					}
					optionStr := strings.TrimLeft(optionRange[0], "{")
					optionStr = strings.TrimRight(optionStr, "}")
					enumOptions[column.TypeName] = strings.Split(optionStr, ",")
				}
				column.Option = enumOptions[column.TypeName]
			}
		}

		if utils.StringVal(column.Extra) == "auto_increment" {
			column.Extra = utils.StringPtr("AutoIncrement")
		}
	}
	return columns, nil
}
