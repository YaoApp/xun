package dbal

import (
	"fmt"
	"reflect"
	"strings"
)

// Grammars loaded grammar driver
var Grammars = map[string]Grammar{}

// BindingKeys the binding key orders
var BindingKeys = []string{
	"select", "from", "join", "where",
	"groupBy", "having",
	"order",
	"union", "unionOrder",
	"sql",
}

// Operators the default Operators
var Operators = []string{
	"=", "<", ">", "<=", ">=", "<>", "!=", "<=>",
	"like", "like binary", "not like", "ilike",
	"&", "|", "^", "<<", ">>",
	"rlike", "not rlike", "regexp", "not regexp",
	"~", "~*", "!~", "!~*", "similar to",
	"not similar to", "not ilike", "~~*", "!~~*",
}

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
		UseWriteConnection: false,
		IsJoinClause:       false,
		Distinct:           false,
		DistinctColumns:    []interface{}{},
		Columns:            []interface{}{},
		Aggregate:          Aggregate{},
		Groups:             []interface{}{},
		Havings:            []Having{},
		Limit:              -1,
		Offset:             -1,
		UnionLimit:         -1,
		UnionOffset:        -1,
		BindingOffset:      0,
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
			"sql":        {},
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

// IsExpression Determine if the given value is a raw expression.
func IsExpression(value interface{}) bool {
	switch value.(type) {
	case Expression:
		return true
	default:
		return false
	}
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

// IsEmpty determine if the from value is empty
func (from From) IsEmpty() bool {
	return from.Name == nil
}

// GetValue Get the value of the expression.
func (expression Expression) GetValue() string {
	return fmt.Sprintf("%v", expression.Value)
}

// Clone clone the query instance
func (query *Query) Clone() *Query {

	new := Query{
		UseWriteConnection: query.UseWriteConnection,    // Whether to use write connection for the select. default is false
		Lock:               query.CopyLock(),            //  Indicates whether row locking is being used.
		From:               query.CopyFrom(),            // The table which the query is targeting.
		Columns:            query.CopyColumns(),         // The columns that should be returned. (Name or Expression)
		Aggregate:          query.CopyAggregate(),       // An aggregate function and column to be run.
		Wheres:             query.CopyWheres(),          // The where constraints for the query.
		Joins:              query.CopyJoins(),           // The table joins for the query.
		Unions:             query.CopyUnions(),          // The query union statements.
		UnionLimit:         query.UnionLimit,            // The maximum number of union records to return.
		UnionOffset:        query.UnionOffset,           // The number of union records to skip.
		UnionOrders:        query.CopyUnionOrders(),     // The orderings for the union query.
		Orders:             query.CopyOrders(),          // The orderings for the query.
		Limit:              query.Limit,                 // The maximum number of records to return.
		Offset:             query.Offset,                // The number of records to skip.
		Groups:             query.CopyGroups(),          // The groupings for the query.
		Havings:            query.CopyHavings(),         // The having constraints for the query.
		Bindings:           query.CopyBindings(),        // The current query value bindings.
		Distinct:           query.Distinct,              // Indicates if the query returns distinct results. Occasionally contains the columns that should be distinct. default is false
		DistinctColumns:    query.CopyDistinctColumns(), // Indicates if the query returns distinct results. Occasionally contains the columns that should be distinct.
		IsJoinClause:       query.IsJoinClause,          // Determine if the query is a join clause.
		BindingOffset:      query.BindingOffset,         // The Binding offset before select
	}

	// // new := NewQuery()
	// // *new = *query
	// // new.Bindings = map[string][]interface{}{}
	// // new.Bindings = *&query.Bindings
	// fmt.Printf("%p\n%p\n", &new.From, &query.From)
	// os.Exit(0)
	return &new
}

// CopyBindings copy Bindings
func (query *Query) CopyBindings() map[string][]interface{} {
	new := map[string][]interface{}{}
	for key, bindings := range query.Bindings {
		new[key] = []interface{}{}
		for _, binding := range bindings {
			new[key] = append(new[key], binding)
		}
	}
	return new
}

// CopyAggregate copy Aggregate
func (query *Query) CopyAggregate() Aggregate {
	new := query.Aggregate
	return new
}

// CopyLock copy Lock
func (query *Query) CopyLock() interface{} {
	new := query.Lock
	return new
}

// CopyFrom copy from
func (query *Query) CopyFrom() From {
	new := query.From
	return new
}

// CopyColumns copy columns
func (query *Query) CopyColumns() []interface{} {
	new := []interface{}{}
	for _, column := range query.Columns {
		new = append(new, column)
	}
	return new
}

// CopyDistinctColumns copy DistinctColumns
func (query *Query) CopyDistinctColumns() []interface{} {
	new := []interface{}{}
	for _, column := range query.DistinctColumns {
		new = append(new, column)
	}
	return new
}

// CopyWheres copy wheres
func (query *Query) CopyWheres() []Where {
	new := []Where{}
	for _, where := range query.Wheres {
		new = append(new, where)
	}
	return new
}

// CopyJoins copy joins
func (query *Query) CopyJoins() []Join {
	new := []Join{}
	for _, join := range query.Joins {
		new = append(new, join)
	}
	return new
}

// CopyUnions copy unions
func (query *Query) CopyUnions() []Union {
	new := []Union{}
	for _, union := range query.Unions {
		new = append(new, union)
	}
	return new
}

// CopyUnionOrders copy UnionOrders
func (query *Query) CopyUnionOrders() []Order {
	new := []Order{}
	for _, order := range query.UnionOrders {
		new = append(new, order)
	}
	return new
}

// CopyOrders copy Orders
func (query *Query) CopyOrders() []Order {
	new := []Order{}
	for _, order := range query.Orders {
		new = append(new, order)
	}
	return new
}

// CopyGroups copy Groups
func (query *Query) CopyGroups() []interface{} {
	new := []interface{}{}
	for _, group := range query.Groups {
		new = append(new, group)
	}
	return new
}

// CopyHavings copy Havings
func (query *Query) CopyHavings() []Having {
	new := []Having{}
	for _, having := range query.Havings {
		new = append(new, having)
	}
	return new
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

// GetBindings Get the current query value bindings in a flattened array.
func (query *Query) GetBindings(key ...string) []interface{} {

	if len(key) > 0 {
		name := key[0]
		values, has := query.Bindings[name]
		if !has {
			panic(fmt.Errorf("The %s of bindings does not exist", name))
		}
		return values
	}

	bindings := []interface{}{}
	for _, name := range BindingKeys {
		values, has := query.Bindings[name]
		if has && len(values) > 0 {
			bindings = append(bindings, values...)
		}
	}
	return bindings
}
