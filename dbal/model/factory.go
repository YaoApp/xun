package model

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/yaoapp/xun"
)

// modelsRegistered the models have been registered
var modelsRegistered map[string]*Factory = map[string]*Factory{}
var modelsAlias map[string]*Factory = map[string]*Factory{}
var typeOfModel = reflect.TypeOf(Model{})

// Register register the model
func Register(v interface{}, args ...interface{}) {
	reflectPtr := reflect.ValueOf(v)
	reflectValue := reflect.Indirect(reflectPtr)

	if isSchema(reflectValue, args...) {
		registerSchema(v, args...)
		return
	} else if isStruct(reflectPtr, reflectValue) {
		registerStruct(reflectPtr, reflectValue, v, args...)
		return
	}

	panic(fmt.Errorf("The type kind (%s) can't be register, have %d arguments", reflectPtr.Type().String(), len(args)))
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
		}
	}

	schema.Columns = []Column{}
	for _, column := range columnsMap {
		schema.Columns = append(schema.Columns, column)
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
		column = *ctag.merge(column)

	}

	if column.Name == "" {
		column.Name = xun.ToSnakeCase(field.Name)
	}
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

func setupAttributes(model *Model, schema *Schema) {

	// init
	model.attributes = map[string]Attribute{}

	// set Columns
	for i, column := range schema.Columns {
		name := column.Name
		attr := Attribute{
			Name:         column.Name,
			Column:       &schema.Columns[i],
			Value:        nil,
			Relationship: nil,
		}
		model.attributes[name] = attr
	}

	// set Relationships
	for i, relation := range schema.Relationships {
		name := relation.Name
		attr := Attribute{
			Name:         relation.Name,
			Relationship: &schema.Relationships[i],
			Column:       nil,
			Value:        nil,
		}
		model.attributes[name] = attr
	}
}

// Class get the factory by model name
func Class(name string) *Factory {
	factory, has := modelsAlias[name]
	if !has {
		panic(fmt.Errorf("The model (%s) doesn't register", name))
	}
	return factory
}

// GetModel get the model instance pointer
func GetModel(v interface{}) Model {
	switch v.(type) {
	case *Model:
		return *v.(*Model)
	default:
		reflectPtr := reflect.ValueOf(v)
		reflectValue := reflect.Indirect(reflectPtr)
		if reflectPtr.Kind() == reflect.Ptr && reflectValue.Kind() == reflect.Struct && reflectValue.FieldByName("Model").Type() == typeOfModel {
			return reflectValue.FieldByName("Model").Interface().(Model)
		}
	}
	panic(fmt.Errorf("The type  (%s) can't be register", reflect.TypeOf(v).String()))
}

// SetModel set the model instance pointer
func SetModel(v interface{}, model interface{}) {
	switch model.(type) {
	case Model:
		reflectPtr := reflect.ValueOf(v)
		reflectValue := reflect.Indirect(reflectPtr)
		if reflectPtr.Kind() == reflect.Ptr && reflectValue.Kind() == reflect.Struct && reflectValue.FieldByName("Model").Type() == typeOfModel {
			reflectPtr.Elem().FieldByName("Model").Set(reflect.ValueOf(model))
			return
		}
		break
	case func(model *Model):
		callback := model.(func(model *Model))
		vModel := GetModel(v)
		callback(&vModel)
		SetModel(v, vModel)
		return
	}

	panic(fmt.Errorf("The type (%s) can't be set", reflect.TypeOf(v).String()))
}

// New build a model instance quickly
func (factory *Factory) New(v ...interface{}) *Model {

	if len(v) > 0 {
		ptr := reflect.ValueOf(v[0])
		if ptr.Kind() != reflect.Ptr {
			panic(fmt.Errorf("The model type (%s) must be a pointer", ptr.Kind().String()))
		}
		ptr.Elem().Set(reflect.ValueOf(factory.Model).Elem())
		return nil
	}

	clone := *(factory.Model.(*Model))
	return &clone
}

// Migrate running a database migration automate
func (factory *Factory) Migrate(model string, args ...bool) {
}

// GetMethods get the model methods for auto-generate the APIs
func (factory *Factory) GetMethods(model string, args ...bool) {
}
