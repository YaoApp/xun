package sql

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/blang/semver/v4"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/logger"
	"github.com/yaoapp/xun/utils"
)

// GetDatabase get the database name of the current connection
func (grammarSQL SQL) GetDatabase() string {
	return grammarSQL.DatabaseName
}

// GetSchema get the schema name of the current connection
func (grammarSQL SQL) GetSchema() string {
	return grammarSQL.SchemaName
}

// GetVersion get the version of the connection database
func (grammarSQL SQL) GetVersion() (*dbal.Version, error) {
	sql := fmt.Sprintf("SELECT VERSION()")
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	rows := []string{}
	err := grammarSQL.DB.Select(&rows, sql)
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("Can't get the version")
	}

	ver, err := semver.Make(rows[0])
	if err != nil {
		defer logger.Error(500, rows[0]).Write()
		return nil, err
	}

	return &dbal.Version{
		Version: ver,
		Driver:  grammarSQL.Driver,
	}, nil
}

// GetTables Get all of the table names for the database.
func (grammarSQL SQL) GetTables() ([]string, error) {
	sql := "SHOW TABLES"
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	tables := []string{}
	err := grammarSQL.DB.Select(&tables, sql)
	if err != nil {
		return nil, err
	}
	return tables, nil
}

// TableExists check if the table exists
func (grammarSQL SQL) TableExists(name string) (bool, error) {
	sql := fmt.Sprintf("SHOW TABLES like %s", grammarSQL.VAL(name))
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	rows := []string{}
	err := grammarSQL.DB.Select(&rows, sql)
	if err != nil {
		return false, err
	}
	if len(rows) == 0 {
		return false, nil
	}
	return name == fmt.Sprintf("%s", rows[0]), nil
}

// GetTable get a table on the schema
func (grammarSQL SQL) GetTable(name string) (*dbal.Table, error) {

	has, err := grammarSQL.TableExists(name)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, fmt.Errorf("the table %s does not exists", name)
	}

	table := dbal.NewTable(name, grammarSQL.GetSchema(), grammarSQL.GetDatabase())
	columns, err := grammarSQL.GetColumnListing(table.SchemaName, table.TableName)
	if err != nil {
		return nil, err
	}

	indexes, err := grammarSQL.GetIndexListing(table.SchemaName, table.TableName)
	if err != nil {
		return nil, err
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
			return nil, fmt.Errorf("the column does not exists %s", idx.ColumnName)
		}
		column := table.ColumnMap[idx.ColumnName]
		if !table.HasIndex(idx.Name) {
			index := *idx
			index.Columns = []*dbal.Column{}
			column.Indexes = append(column.Indexes, &index)
			table.PushIndex(&index)
		}
		index := table.IndexMap[idx.Name]
		index.AddColumn(column)
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

	return table, nil
}

// GetIndexListing get a table indexes structure
func (grammarSQL SQL) GetIndexListing(dbName string, tableName string) ([]*dbal.Index, error) {
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
		grammarSQL.VAL(dbName),
		grammarSQL.VAL(tableName),
	)
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	indexes := []*dbal.Index{}
	err := grammarSQL.DB.Select(&indexes, sql)
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
func (grammarSQL SQL) GetColumnListing(dbName string, tableName string) ([]*dbal.Column, error) {
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
		"COLUMN_TYPE as `type_name`",
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
		grammarSQL.Quoter.VAL(dbName),
		grammarSQL.Quoter.VAL(tableName),
	)
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	columns := []*dbal.Column{}
	err := grammarSQL.DB.Select(&columns, sql)
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

		if column.Type == "enum" {
			re := regexp.MustCompile(`enum\('(.*)'\)`)
			matched := re.FindStringSubmatch(column.TypeName)
			if len(matched) == 2 {
				options := strings.Split(matched[1], "','")
				column.Option = options
			}
		}

		if utils.StringVal(column.Extra) == "auto_increment" {
			column.Extra = utils.StringPtr("AutoIncrement")
		}
	}
	return columns, nil
}

// CreateTable create a new table on the schema
func (grammarSQL SQL) CreateTable(table *dbal.Table) error {
	name := grammarSQL.ID(table.TableName)
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
	//    CreatePrimary for creating the primary key
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

	// Columns
	for _, Column := range columns {
		stmts = append(stmts,
			grammarSQL.SQLAddColumn(Column),
		)
	}

	// Primary key
	if primary != nil {
		stmts = append(stmts,
			grammarSQL.SQLAddPrimary(primary),
		)
	}

	// indexes
	for _, index := range indexes {
		indexStmt := grammarSQL.SQLAddIndex(index)
		if indexStmt != "" {
			stmts = append(stmts, indexStmt)
		}
	}

	engine := utils.GetIF(table.Engine != "", "ENGINE "+table.Engine, "")
	charset := utils.GetIF(table.Charset != "", "DEFAULT CHARSET "+table.Charset, "")
	collation := utils.GetIF(table.Collation != "", "COLLATE="+table.Collation, "")

	sql = sql + strings.Join(stmts, ",\n")
	sql = sql + fmt.Sprintf(
		"\n) %s %s %s ROW_FORMAT=DYNAMIC",
		engine, charset, collation,
	)
	defer logger.Debug(logger.CREATE, sql).TimeCost(time.Now())
	_, err := grammarSQL.DB.Exec(sql)

	// Callback
	for _, cmd := range cbCommands {
		cmd.Callback(err)
	}

	return err
}

// DropTable a table from the schema.
func (grammarSQL SQL) DropTable(name string) error {
	sql := fmt.Sprintf("DROP TABLE %s", grammarSQL.ID(name))
	defer logger.Debug(logger.DELETE, sql).TimeCost(time.Now())
	_, err := grammarSQL.DB.Exec(sql)
	return err
}

// DropTableIfExists if the table exists, drop it from the schema.
func (grammarSQL SQL) DropTableIfExists(name string) error {
	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s", grammarSQL.ID(name))
	defer logger.Debug(logger.DELETE, sql).TimeCost(time.Now())
	_, err := grammarSQL.DB.Exec(sql)
	return err
}

// RenameTable rename a table on the schema.
func (grammarSQL SQL) RenameTable(old string, new string) error {
	sql := fmt.Sprintf("ALTER TABLE %s RENAME %s", grammarSQL.ID(old), grammarSQL.ID(new))
	defer logger.Debug(logger.UPDATE, sql).TimeCost(time.Now())
	_, err := grammarSQL.DB.Exec(sql)
	return err
}

// AlterTable alter a table on the schema
func (grammarSQL SQL) AlterTable(table *dbal.Table) error {

	sql := fmt.Sprintf("ALTER TABLE %s ", grammarSQL.ID(table.TableName))
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
			grammarSQL.alterTableAddColumn(table, command, sql, &stmts, &errs)
			break
		case "ChangeColumn":
			grammarSQL.alterTableChangeColumn(table, command, sql, &stmts, &errs)
			break
		case "RenameColumn":
			grammarSQL.alterTableRenameColumn(table, command, sql, &stmts, &errs)
			break
		case "DropColumn":
			grammarSQL.alterTableDropColumn(table, command, sql, &stmts, &errs)
			break
		case "CreateIndex":
			grammarSQL.alterTableCreateIndex(table, command, sql, &stmts, &errs)
			break
		case "RenameIndex":
			grammarSQL.alterTableRenameIndex(table, command, sql, &stmts, &errs)
			break
		case "DropIndex":
			grammarSQL.alterTableDropIndex(table, command, sql, &stmts, &errs)
			break
		case "CreatePrimary":
			grammarSQL.alterTableCreatePrimary(table, command, sql, &stmts, &errs)
			break
		case "DropPrimary":
			grammarSQL.alterTableDropPrimary(table, command, sql, &stmts, &errs)
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

func (grammarSQL SQL) alterTableAddColumn(table *dbal.Table, command *dbal.Command, sql string, stmts *[]string, errs *[]error) {
	column := command.Params[0].(*dbal.Column)
	stmt := "ADD " + grammarSQL.SQLAddColumn(column)
	*stmts = append(*stmts, sql+stmt)
	err := grammarSQL.ExecSQL(table, sql+stmt)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("AddColumn: %s", err))
	}
	command.Callback(err)
}

func (grammarSQL SQL) alterTableChangeColumn(table *dbal.Table, command *dbal.Command, sql string, stmts *[]string, errs *[]error) {
	column := command.Params[0].(*dbal.Column)
	stmt := "MODIFY " + grammarSQL.SQLAddColumn(column)
	*stmts = append(*stmts, sql+stmt)
	err := grammarSQL.ExecSQL(table, sql+stmt)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("ChangeColumn %s: %s", column.Name, err))
	}
	command.Callback(err)
}

func (grammarSQL SQL) alterTableRenameColumn(table *dbal.Table, command *dbal.Command, sql string, stmts *[]string, errs *[]error) {
	old := command.Params[0].(string)
	new := command.Params[1].(string)
	column, has := table.ColumnMap[old]
	var err error
	if !has {
		err = fmt.Errorf("RenameColumn: The column %s not exists", old)
		*errs = append(*errs, fmt.Errorf("RenameColumn: The column %s not exists", old))
	} else {
		column.Name = new
		stmt := fmt.Sprintf("CHANGE COLUMN %s %s",
			grammarSQL.ID(old),
			grammarSQL.SQLAddColumn(column),
		)
		*stmts = append(*stmts, sql+stmt)
		err = grammarSQL.ExecSQL(table, sql+stmt)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("RenameColumn: %s", err))
		}
	}
	command.Callback(err)
}

func (grammarSQL SQL) alterTableDropColumn(table *dbal.Table, command *dbal.Command, sql string, stmts *[]string, errs *[]error) {
	name := command.Params[0].(string)
	stmt := fmt.Sprintf("DROP COLUMN %s", grammarSQL.ID(name))
	*stmts = append(*stmts, sql+stmt)
	err := grammarSQL.ExecSQL(table, sql+stmt)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("DropColumn: %s", err))
	}
	command.Callback(err)
}

func (grammarSQL SQL) alterTableCreateIndex(table *dbal.Table, command *dbal.Command, sql string, stmts *[]string, errs *[]error) {
	index := command.Params[0].(*dbal.Index)
	stmt := "ADD " + grammarSQL.SQLAddIndex(index)
	*stmts = append(*stmts, sql+stmt)
	err := grammarSQL.ExecSQL(table, sql+stmt)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("CreateIndex: %s", err))
	}
	command.Callback(err)
}

func (grammarSQL SQL) alterTableRenameIndex(table *dbal.Table, command *dbal.Command, sql string, stmts *[]string, errs *[]error) {
	old := command.Params[0].(string)
	new := command.Params[1].(string)
	oldIndex := table.GetIndex(old)
	if oldIndex == nil {
		err := fmt.Errorf("RenameIndex: The index %s not found", old)
		*errs = append(*errs, err)
		command.Callback(err)
		return
	}

	stmt := fmt.Sprintf("DROP INDEX %s", grammarSQL.ID(old))
	*stmts = append(*stmts, sql+stmt)
	err := grammarSQL.ExecSQL(table, sql+stmt)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("RenameIndex: %s", err))
		command.Callback(err)
		return
	}

	newIndex := oldIndex
	newIndex.Name = new
	stmt = "ADD " + grammarSQL.SQLAddIndex(newIndex)
	*stmts = append(*stmts, sql+stmt)
	err = grammarSQL.ExecSQL(table, sql+stmt)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("RenameIndex: %s", err))
	}
	command.Callback(err)

	// stmt := fmt.Sprintf("RENAME INDEX %s TO %s", grammarSQL.Quoter.ID(old, db), grammarSQL.Quoter.ID(new, db))
	// stmts = append(stmts, sql+stmt)
	// err := grammarSQL.ExecSQL(db, table, sql+stmt)
	// if err != nil {
	// 	errs = append(errs, err)
	// }
	// command.Callback(err)
}

func (grammarSQL SQL) alterTableDropIndex(table *dbal.Table, command *dbal.Command, sql string, stmts *[]string, errs *[]error) {
	name := command.Params[0].(string)
	stmt := fmt.Sprintf("DROP INDEX %s", grammarSQL.ID(name))
	*stmts = append(*stmts, sql+stmt)
	err := grammarSQL.ExecSQL(table, sql+stmt)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("DropIndex: %s", err))
	}
	command.Callback(err)
}

func (grammarSQL SQL) alterTableCreatePrimary(table *dbal.Table, command *dbal.Command, sql string, stmts *[]string, errs *[]error) {
	primary := command.Params[0].(*dbal.Primary)
	stmt := "ADD " + grammarSQL.SQLAddPrimary(primary)
	*stmts = append(*stmts, sql+stmt)
	err := grammarSQL.ExecSQL(table, sql+stmt)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("CreateIndex: %s", err))
	}
	command.Callback(err)
}

func (grammarSQL SQL) alterTableDropPrimary(table *dbal.Table, command *dbal.Command, sql string, stmts *[]string, errs *[]error) {
	columns := command.Params[1].([]*dbal.Column)
	for _, column := range columns {
		// remove AutoIncrement
		if utils.StringVal(column.Extra) == "AutoIncrement" {
			column.Extra = nil
			stmt := "MODIFY " + grammarSQL.SQLAddColumn(column)
			*stmts = append(*stmts, sql+stmt)
			err := grammarSQL.ExecSQL(table, sql+stmt)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("DropPrimary: %s", err))
			}
		}
	}
	stmt := "DROP PRIMARY KEY"
	*stmts = append(*stmts, sql+stmt)
	err := grammarSQL.ExecSQL(table, sql+stmt)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("DropPrimary: %s", err))
	}
	command.Callback(err)
}

// ExecSQL execute sql then update table structure
func (grammarSQL SQL) ExecSQL(table *dbal.Table, sql string) error {
	_, err := grammarSQL.DB.Exec(sql)
	if err != nil {
		return err
	}
	// update table structure
	new, err := grammarSQL.GetTable(table.TableName)
	if err != nil {
		return err
	}

	*table = *new
	return nil
}

// GetTypeFromComment Get the type name from comment
func (grammarSQL SQL) GetTypeFromComment(comment *string) string {
	if comment == nil {
		return ""
	}

	lines := strings.Split(*comment, "|")
	if len(lines) < 1 {
		return ""
	}

	re := regexp.MustCompile(`^T:([a-zA-Z]+)`)
	matched := re.FindStringSubmatch(lines[0])
	if len(matched) == 2 {
		return matched[1]
	}

	return ""
}
