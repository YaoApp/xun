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

// Registered get all registered models
func Registered() []string {
	names := []string{}
	for name := range modelsRegistered {
		names = append(names, name)
	}
	return names
}

// IsRegistered determine if the model was registered
func IsRegistered(name string) bool {
	_, has := modelsAlias[name]
	return has
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
		SetModel(factory.Model, func(model *Model) {
			model.values = xun.MakeRow()
			model.relations = factory.Schema.Relations
		})
		new := reflect.ValueOf(factory.Model).Elem()
		ptr.Elem().Set(new)
		return nil
	}

	model, ok := factory.Model.(*Model)
	if !ok {
		panic(fmt.Errorf("The factory.Model is not a model pointer"))
	}

	clone := *model
	clone.values = xun.MakeRow()
	clone.relations = factory.Schema.Relations
	return &clone
}

// NewBasic create a basic interface{}
func (factory *Factory) NewBasic(v ...interface{}) Basic {
	basic, ok := factory.Model.(Basic)
	if !ok {
		panic(fmt.Errorf("The factory.Model is not a model Basic interface"))
	}
	return basic.Clone()
}

// Methods get the model methods for auto-generate the APIs
func (factory *Factory) Methods(args ...bool) []Method {

	if factory.methods != nil {
		return factory.methods
	}

	methods := []Method{}
	reflectModelType := reflect.TypeOf(factory.Model)
	for i := 0; i < reflectModelType.NumMethod(); i++ {
		method := reflectModelType.Method(i)
		name := method.Name
		path := strings.ToLower(name)

		in := []string{}
		for j := 1; j < method.Type.NumIn(); j++ {
			in = append(in, method.Type.In(j).Kind().String())
		}

		out := []string{}
		for j := 0; j < method.Type.NumOut(); j++ {
			out = append(out, method.Type.Out(j).Kind().String())
		}

		methods = append(methods, Method{
			Name:   name,
			Path:   path,
			In:     in,
			Out:    out,
			Export: true,
		})
	}
	factory.methods = methods
	return methods
}

// Migrate running a database migration automate
func (factory *Factory) Migrate(schema schema.Schema, query query.Query, args ...bool) error {

	if factory.Schema.Table.Name == "" {
		return fmt.Errorf("The table name does not set")
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

func (factory *Factory) diffSchema(schema schema.Schema, force bool) {
	panic(fmt.Errorf(`This feature does not support it yet. It working when the first parameter refresh is true.(model.Class("user").Migrate(schema, true))`))
}
