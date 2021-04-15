package model

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/dbal/schema"
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
	panic(fmt.Errorf("v is (%s) not a model", reflect.TypeOf(v).String()))
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

	panic(fmt.Errorf("v is (%s) not a model", reflect.TypeOf(v).String()))
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

// GetMethods get the model methods for auto-generate the APIs
func (factory *Factory) GetMethods(model string, args ...bool) {
}

// Migrate running a database migration automate
func (factory *Factory) Migrate(schema schema.Schema, query query.Query, args ...bool) error {

	if factory.Schema.Table.Name == "" {
		return nil
	}

	table := factory.Schema.Table.Name
	refresh, force := prepareMigrateArgs(args...)
	if refresh {
		err := schema.DropTableIfExists(table)
		if err != nil {
			return err
		}
		err = factory.createTable(table, schema)
		if err != nil {
			return err
		}

		// Insert values
		if factory.Schema.Values == nil {
			return nil
		}

		err = query.Table(table).Insert(factory.Schema.Values)
		if err != nil {
			return err
		}
		return nil
	}

	// @todo
	factory.diffSchema(schema, force)
	return nil
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

func (factory *Factory) diffSchema(schema schema.Schema, force bool) {
	panic(fmt.Errorf(`This feature does not support it yet. It working when the first parameter refresh is true.(model.Class("user").Migrate(schema, true))`))
}
