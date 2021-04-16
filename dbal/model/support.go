package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/dbal/schema"
)

// prepareRegisterNames parse the model name, return (fullname string, namespace string, name string)
func prepareRegisterNames(name string) (string, string, string) {
	sep := "."
	if strings.Contains(name, "/") {
		sep = "/"
	}
	name = strings.ToLower(strings.TrimPrefix(name, "*"))
	namer := strings.Split(name, sep)
	length := len(namer)
	if length <= 1 {
		return name, "", name
	}
	fullname := strings.Join(namer, ".")
	namespace := strings.Join(namer[0:length-1], ".")
	name = namer[length-1]
	return fullname, namespace, name
}

// prepareRegisterArgs parse the params for Register(), return (schema *Schema, flow *Flow)
func prepareRegisterArgs(args ...interface{}) (*Schema, *Flow) {
	var schema *Schema = nil
	var flow *Flow = nil

	if len(args) > 0 {
		content, ok := args[0].([]byte)
		if !ok {
			panic(fmt.Errorf("The schema type is %s, should be []byte", reflect.TypeOf(args[0]).String()))
		}

		schema = &Schema{}
		err := json.Unmarshal(content, schema)
		if err != nil {
			panic(fmt.Errorf("The parse schema error. %s ", err.Error()))
		}

	}

	if len(args) > 1 {
		content, ok := args[1].([]byte)
		if !ok {
			panic(fmt.Errorf("The flow type is %s, should be []byte", reflect.TypeOf(args[1]).String()))
		}

		flow = &Flow{}
		err := json.Unmarshal(content, flow)
		if err != nil {
			panic(fmt.Errorf("The parse flow error. %s ", err.Error()))
		}
	}

	return schema, flow
}

// prepareMigrateArgs parse the params for migrate, return (refresh bool, force bool)
func prepareMigrateArgs(args ...bool) (bool, bool) {
	refresh := false
	force := false
	if len(args) > 0 {
		refresh = args[0]
	}

	if len(args) > 1 {
		force = args[1]
	}

	return refresh, force
}

func prepareBlueprintArgs(method string, column *Column) []reflect.Value {
	in := []reflect.Value{reflect.ValueOf(column.Name)}
	switch method {
	case "String", "Char", "Binary":
		if column.Length > 0 {
			in = append(in, reflect.ValueOf(column.Length))
		}
		break
	case "Decimal", "UnsignedDecimal", "Float", "UnsignedFloat", "Double", "UnsignedDouble":
		args := []int{}
		if column.Precision > 0 {
			args = append(args, column.Precision)
		}
		if column.Scale > 0 {
			if len(args) == 0 {
				args = append(args, 10)
			}
			args = append(args, column.Scale)
		}
		if len(args) > 0 {
			for _, arg := range args {
				in = append(in, reflect.ValueOf(arg))
			}
		}
		break
	case "DateTime", "DateTimeTz", "Time", "TimeTz", "Timestamp", "TimestampTz", "Timestamps", "TimestampsTz", "SoftDeletes", "SoftDeletesTz":
		if column.Precision > 0 {
			in = append(in, reflect.ValueOf(column.Length))
		}
		break
	case "Enum":
		in = append(in, reflect.ValueOf(column.Option))
		break
	}
	return in
}

func getTypeName(v interface{}) string {
	return reflect.TypeOf(v).String()
}

// determine if the interface{} is json schema
func isSchema(reflectValue reflect.Value, args ...interface{}) bool {
	return reflectValue.Kind() == reflect.String && len(args) > 0
}

// determine if the interface{} is golang struct
func isStruct(reflectPtr reflect.Value, reflectValue reflect.Value) bool {
	return reflectPtr.Kind() == reflect.Ptr && reflectValue.Kind() == reflect.Struct && reflectValue.FieldByName("Model").Type() == typeOfModel
}

// register the model by given json schema
func registerSchema(v interface{}, args ...interface{}) {
	origin := v.(string)
	fullname, namespace, name := prepareRegisterNames(origin)
	schema, flow := prepareRegisterArgs(args...)
	model := Model{}
	model.namespace = namespace
	model.name = name
	setupAttributes(&model, schema)

	factory := &Factory{
		Namespace: namespace,
		Name:      name,
		Model:     &model,
		Schema:    schema,
		Flow:      flow,
	}
	modelsRegistered[origin] = factory
	modelsAlias[origin] = modelsRegistered[origin]
	modelsAlias[fullname] = modelsRegistered[origin]
}

// register the model by given golang struct pointer
func registerStruct(reflectPtr reflect.Value, reflectValue reflect.Value, v interface{}, args ...interface{}) {
	origin := reflectPtr.Type().String()
	fullname, namespace, name := prepareRegisterNames(origin)
	schema, flow := prepareRegisterArgs(args...)
	SetModel(v, func(model *Model) {
		model.namespace = namespace
		model.name = name
		setupAttributesStruct(model, schema, reflectValue)
	})

	factory := &Factory{
		Namespace: namespace,
		Name:      name,
		Model:     v,
		Schema:    schema,
		Flow:      flow,
	}
	modelsRegistered[origin] = factory
	modelsAlias[origin] = modelsRegistered[origin]
	modelsAlias[fullname] = modelsRegistered[origin]
}

func setupAttributesStruct(model *Model, schema *Schema, reflectValue reflect.Value) {

	names := []string{}
	columns := []Column{}
	for i := 0; i < reflectValue.NumField(); i++ {
		column := fieldToColumn(reflectValue.Type().Field(i))
		if column != nil {
			columns = append(columns, *column)
		}
	}

	columns = append(columns, schema.Columns...)

	// merge schema
	columnsMap := map[string]Column{}
	for _, column := range columns {
		if col, has := columnsMap[column.Name]; has {
			columnsMap[column.Name] = *col.merge(column)
		} else {
			columnsMap[column.Name] = column
			names = append(names, column.Name)
		}
	}

	schema.Columns = []Column{}
	for _, name := range names {
		schema.Columns = append(schema.Columns, columnsMap[name])
	}

	setupAttributes(model, schema)
}

func fieldToColumn(field reflect.StructField) *Column {
	if field.Type == typeOfModel {
		return nil
	}

	column, has := StructMapping[field.Type.Kind()]
	if !has {
		return nil
	}

	ctag := parseFieldTag(string(field.Tag))
	if ctag != nil {
		column = *column.merge(*ctag)
	}

	if column.Name == "" {
		column.Name = xun.ToSnakeCase(field.Name)
	}

	column.Field = field.Name
	return &column
}

func parseFieldTag(tag string) *Column {
	if !strings.Contains(tag, "x-") {
		return nil
	}

	params := map[string]string{}
	tagarr := strings.Split(tag, "x-")

	for _, tagstr := range tagarr {
		tagr := strings.Split(tagstr, ":")
		if len(tagr) == 2 {
			key := strings.Trim(tagr[0], " ")
			value := strings.Trim(strings.Trim(tagr[1], " "), "\"")
			key = strings.TrimPrefix(key, "x-")
			key = strings.ReplaceAll(key, "-", ".")
			if key == "json" {
				key = "name"
			}
			params[key] = value
		}
	}

	if len(params) == 0 {
		return nil
	}

	column := Column{}
	for name, value := range params {
		column.set(name, value)
	}

	return &column
}

func extendKeys(column Column, primaryKeys *[]string, uniqueKeys *[]string) {

	if column.Primary || column.Type == "ID" {
		*primaryKeys = append(*primaryKeys, column.Name)
	}

	if column.Unique {
		*uniqueKeys = append(*uniqueKeys, column.Name)
	}

}

func setupAttributes(model *Model, schema *Schema) {

	// init
	model.attributes = map[string]Attribute{}
	model.values = xun.MakeRow()
	model.columns = []*Column{}
	model.columnNames = []string{}
	model.searchable = []string{}
	model.uniqueKeys = []string{}
	searchable := map[string]bool{}
	model.primary = ""

	// set Columns
	for i, column := range schema.Columns {
		name := column.Name
		attr := Attribute{
			Name:         column.Name,
			Column:       &schema.Columns[i],
			Relationship: nil,
		}
		model.attributes[name] = attr
		model.columns = append(model.columns, &schema.Columns[i])
		model.columnNames = append(model.columnNames, column.Name)

		// set indexes
		if column.Index || column.Unique || column.Primary || column.Type == "ID" {
			searchable[column.Name] = true
		}

		extendKeys(column, &model.primaryKeys, &model.uniqueKeys)
	}

	// set Relationships
	for i, relation := range schema.Relationships {
		name := relation.Name
		attr := Attribute{
			Name:         relation.Name,
			Relationship: &schema.Relationships[i],
			Column:       nil,
		}
		model.attributes[name] = attr
	}

	// set indexes
	for _, index := range schema.Indexes {
		for _, column := range index.Columns {
			searchable[column] = true
			if index.Type == "primary" {
				model.primaryKeys = append(model.primaryKeys, column)
			}
			if index.Type == "unique" {
				model.uniqueKeys = append(model.uniqueKeys, column)
			}
		}

	}

	// set searchable
	for column := range searchable {
		model.searchable = append(model.searchable, column)
	}

	// set primary
	if len(model.primaryKeys) > 0 {
		model.primary = model.primaryKeys[0]
	}

	// set table
	model.table = &schema.Table

}

func (factory *Factory) createTable(tableName string, sch schema.Schema) error {
	return sch.CreateTable(tableName, func(table schema.Blueprint) {
		// Columns
		for _, column := range factory.Schema.Columns {
			factory.setColumn(table, column)
		}

		// Indexes
		for _, index := range factory.Schema.Indexes {
			factory.createIndex(table, index)
		}
	})
}

func (factory *Factory) setColumn(table schema.Blueprint, column Column) {

	reflectTable := reflect.ValueOf(table)
	methodName := xun.UpperFirst(column.Type)
	method := reflectTable.MethodByName(methodName)
	if method.Kind() == reflect.Func && column.Name != "" {
		in := prepareBlueprintArgs(methodName, &column)
		out := method.Call(in)
		if len(out) != 1 {
			panic(fmt.Errorf("call %s(%s), return value is error", methodName, column.Name))
		}
		col, ok := out[0].Interface().(*schema.Column)
		if !ok {
			panic(fmt.Errorf("call %s(%s), return value is error", methodName, column.Name))
		}
		if column.Comment != "" {
			col.SetComment(column.Comment)
		}

		if column.Primary {
			col.Primary()
		}

		if column.Index {
			col.Index()
		}

		if column.Unique {
			col.Unique()
		}

		if column.DefaultRaw != "" {
			col.SetDefaultRaw(column.DefaultRaw)
		} else if column.Default != nil {
			col.SetDefault(column.Default)
		}

		if !column.Nullable {
			col.NotNull()
		}
	}
}

func (factory *Factory) createIndex(table schema.Blueprint, index Index) {

	if len(index.Columns) == 0 {
		return
	}

	name := index.Name
	if name == "" {
		name = strings.Join(index.Columns, "_")
	}

	// primary,unique,index,match
	if index.Type == "index" || index.Type == "" {
		table.AddIndex(name, index.Columns...)
	} else if index.Type == "unique" {
		table.AddUnique(name, index.Columns...)
	} else if index.Type == "primary" {
		table.AddPrimary(index.Columns...)
	} else if index.Type == "fulltext" {
		table.AddFulltext(name, index.Columns...)
	}

}

// makeBySchema make a new xun model instance
func makeBySchema(query query.Query, schema schema.Schema, v interface{}, args ...interface{}) *Model {

	name, ok := v.(string)
	if !ok {
		panic(fmt.Errorf("the model name is not string"))
	}

	class, has := modelsRegistered[name]
	if !has {
		Register(name, args...)
		class, has = modelsRegistered[name]
		if !has {
			panic(fmt.Errorf("the model register failure"))
		}
	}
	model := class.New()
	model.schema = schema
	model.query = query
	return model
}

// makeByStruct make a new xun model instance
func makeByStruct(query query.Query, schema schema.Schema, v interface{}) {
	name := getTypeName(v)
	Class(name).New(v)
	SetModel(v, func(model *Model) {
		model.query = query
		model.schema = schema
	})
}

func setFieldValue(v interface{}, field string, value interface{}) {
	reflectPtr := reflect.ValueOf(v)
	if reflectPtr.Kind() != reflect.Ptr {
		panic(fmt.Errorf("v is %s, should be a struct pointer", reflectPtr.Type().String()))
	}

	reflectField := reflectPtr.Elem().FieldByName(field)
	if reflectField.Kind() == reflect.Invalid {
		return
	}

	if !reflectField.CanSet() {
		return
	}

	reflectValue := reflect.ValueOf(value)
	if !xun.CastType(&reflectValue, reflectValue.Kind(), reflectField.Kind()) {
		panic(fmt.Errorf("field %s value type is %s, should be %s", field, reflectValue.Kind().String(), reflectField.Kind().String()))
	}

	reflectField.Set(reflectValue)
}

func getFieldTags(v interface{}) []string {
	reflectPtr := reflect.ValueOf(v)
	structValue := reflectPtr.Elem()
	structType := reflect.TypeOf(structValue.Interface())
	tags := []string{}
	for i := 0; i < structType.NumField(); i++ {
		if !structValue.Field(i).CanInterface() {
			continue
		}
		tag := xun.GetTagName(structType.Field(i), "json")
		if tag != "" && tag != "-" {
			tags = append(tags, tag)
		}
	}
	return tags
}
