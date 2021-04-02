package dbal

import (
	"fmt"
	"reflect"
	"strings"
)

// Grammars loaded grammar driver
var Grammars = map[string]Grammar{}

// Register register the grammar driver
func Register(name string, grammar Grammar) {
	Grammars[name] = grammar
}

// GetName get the table name
func (table *Table) GetName() string {
	return table.TableName
}

// NewConstraint make a new constraint intstance
func NewConstraint(schemaName string, tableName string, columnName string) *Constraint {
	return &Constraint{
		TableName:  tableName,
		SchemaName: schemaName,
		ColumnName: columnName,
		Args:       []string{},
	}
}

// NewTable make a grammar table
func NewTable(name string, schemaName string, dbName string) *Table {
	return &Table{
		DBName:     dbName,
		SchemaName: schemaName,
		TableName:  name,
		Primary:    nil,
		Columns:    []*Column{},
		ColumnMap:  map[string]*Column{},
		Indexes:    []*Index{},
		IndexMap:   map[string]*Index{},
		Commands:   []*Command{},
	}
}

// NewName make a new Name instance
func NewName(fullname string, prefix ...string) Name {
	name := Name{}
	if len(prefix) > 0 {
		name.Prefix = prefix[0]
	}
	namer := strings.Split(strings.ToLower(fullname), " as ")
	if len(namer) == 2 {
		name.Name = strings.Trim(namer[0], " ")
		name.Alias = strings.Trim(namer[1], " ")
		return name
	}
	name.Name = strings.Trim(fullname, " ")
	return name
}

// NewQuery make a new Query instance
func NewQuery() *Query {
	query := Query{
		IsJoinClause: false,
		Distinct:     false,
		Operators: []string{
			"=", "<", ">", "<=", ">=", "<>", "!=", "<=>",
			"like", "like binary", "not like", "ilike",
			"&", "|", "^", "<<", ">>",
			"rlike", "not rlike", "regexp", "not regexp",
			"~", "~*", "!~", "!~*", "similar to",
			"not similar to", "not ilike", "~~*", "!~~*",
		},
		BindingKeys: []string{
			"select", "from", "join", "where",
			"groupBy", "having",
			"order",
			"union", "unionOrder",
		},
		Bindings: map[string][]interface{}{
			"select":     {},
			"from":       {},
			"join":       {},
			"where":      {},
			"groupBy":    {},
			"having":     {},
			"order":      {},
			"union":      {},
			"unionOrder": {},
		},
	}
	return &query
}

// NewExpression make a new expression instance
func NewExpression(value interface{}) Expression {
	return Expression{
		Value: value,
	}
}

// Raw make a new expression instance
func Raw(value interface{}) Expression {
	return NewExpression(value)
}

// NewPrimary create a new primary intstance
func (table *Table) NewPrimary(name string, columns ...*Column) *Primary {
	return &Primary{
		DBName:    table.DBName,
		TableName: table.TableName,
		Table:     table,
		Name:      name,
		Columns:   columns,
	}
}

// GetPrimary  get the primary key instance
func (table *Table) GetPrimary(name string, columns ...*Column) *Primary {
	return table.Primary
}

// NewColumn create a new column intstance
func (table *Table) NewColumn(name string) *Column {
	return &Column{
		DBName:            table.DBName,
		TableName:         table.TableName,
		Table:             table,
		Name:              name,
		Length:            nil,
		OctetLength:       nil,
		Precision:         nil,
		Scale:             nil,
		DateTimePrecision: nil,
		Charset:           nil,
		Collation:         nil,
		Key:               nil,
		Extra:             nil,
		Comment:           nil,
	}
}

// PushColumn push a column instance to the table columns
func (table *Table) PushColumn(column *Column) *Table {
	table.ColumnMap[column.Name] = column
	table.Columns = append(table.Columns, column)
	return table
}

// HasColumn checking if the given name column exists
func (table *Table) HasColumn(name string) bool {
	_, has := table.ColumnMap[name]
	return has
}

// GetColumn get the given name column instance
func (table *Table) GetColumn(name string) *Column {
	return table.ColumnMap[name]
}

// NewIndex create a new index intstance
func (table *Table) NewIndex(name string, columns ...*Column) *Index {
	return &Index{
		DBName:    table.DBName,
		TableName: table.TableName,
		Table:     table,
		Name:      name,
		Type:      "index",
		Columns:   columns,
	}
}

// PushIndex push an index instance to the table indexes
func (table *Table) PushIndex(index *Index) *Table {
	table.IndexMap[index.Name] = index
	table.Indexes = append(table.Indexes, index)
	return table
}

// HasIndex checking if the given name index exists
func (table *Table) HasIndex(name string) bool {
	_, has := table.IndexMap[name]
	return has
}

// GetIndex get the given name index instance
func (table *Table) GetIndex(name string) *Index {
	return table.IndexMap[name]
}

// AddCommand Add a new command to the table.
//
// The commands must be:
//    AddColumn(column *Column)    for adding a column
//    ModifyColumn(column *Column) for modifying a colu
//    RenameColumn(old string,new string)  for renaming a column
//    DropColumn(name string)  for dropping a column
//    CreateIndex(index *Index) for creating a index
//    DropIndex( name string) for  dropping a index
//    RenameIndex(old string,new string)  for renaming a index
func (table *Table) AddCommand(name string, success func(), fail func(), params ...interface{}) {
	table.Commands = append(table.Commands, &Command{
		Name:    name,
		Params:  params,
		Success: success,
		Fail:    fail,
	})
}

// Callback run the callback code
func (command *Command) Callback(err error) {
	if err == nil && command.Success != nil {
		command.Success()
	} else if err != nil && command.Fail != nil {
		command.Fail()
	}
}

// AddColumn add column to index
func (index *Index) AddColumn(column *Column) {
	for _, col := range index.Columns {
		if col.Name == column.Name {
			return
		}
	}
	index.Columns = append(index.Columns, column)
}

// Fullname get the name name with prefix
func (name Name) Fullname() string {
	return fmt.Sprintf("%s%s", name.Prefix, name.Name)
}

// As get the alias of name
func (name Name) As() string {
	return name.Alias
}

// GetValue Get the value of the expression.
func (expression Expression) GetValue() string {
	return fmt.Sprintf("%v", expression.Value)
}

// Clone clone the query instance
func (query *Query) Clone() *Query {
	new := *query
	return &new
}

// AddColumn add a column to query
func (query *Query) AddColumn(column interface{}) *Query {
	switch column.(type) {
	case Expression:
		query.Columns = append(query.Columns, column)
	case string:
		query.Columns = append(query.Columns, NewName(column.(string)))
	}
	return query
}

// AddBinding Add a binding to the query.
func (query *Query) AddBinding(typ string, value interface{}) {
	if _, has := query.Bindings[typ]; !has {
		panic(fmt.Errorf("Invalid binding type: %s", typ))
	}

	valueKind := reflect.TypeOf(value).Kind()
	if valueKind == reflect.Array || valueKind == reflect.Slice {
		reflectValue := reflect.ValueOf(value)
		reflectValue = reflect.Indirect(reflectValue)
		for i := 0; i < reflectValue.Len(); i++ {
			value = reflectValue.Index(i).Interface()
			query.Bindings[typ] = append(query.Bindings[typ], value)
		}
	} else {
		query.Bindings[typ] = append(query.Bindings[typ], value)
	}
}
