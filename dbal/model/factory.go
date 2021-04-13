package model

import (
	"fmt"
	"reflect"

	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/dbal/schema"
)

// modelsRegistered the models have been registered
var modelsRegistered map[interface{}]*Factory = map[interface{}]*Factory{}

// Register register the model
func Register(v interface{}, args ...interface{}) {

	reflectPtr := reflect.ValueOf(v)
	reflectValue := reflect.Indirect(reflectPtr)
	if reflectValue.Kind() == reflect.String {
		name := v.(string)
		schema, flow := prepareRegisterArgs(args...)
		modelsRegistered[name] = &Factory{
			Model:  &Model{},
			Schema: schema,
			Flow:   flow,
		}
		return
	} else if reflectPtr.Kind() == reflect.Ptr && reflectValue.Kind() == reflect.Struct {
		name := reflectPtr.Type().String()
		schema, flow := prepareRegisterArgs(args...)
		modelsRegistered[name] = &Factory{
			Model:  v,
			Schema: schema,
			Flow:   flow,
		}
		return
	}

	panic(fmt.Errorf("The type kind (%s) can't be register", reflectValue.Kind().String()))
}

// Make make a new xun model instance
func Make(query query.Query, schema schema.Schema, v interface{}, flow ...interface{}) *Model {
	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		makeByStruct(query, schema, v)
		return nil
	}
	return makeBySchema(query, schema, v, flow...)
}

func makeByStruct(query query.Query, schema schema.Schema, v interface{}) {
	name := getTypeName(v)
	factory, has := modelsRegistered[name]
	if !has {
		panic(fmt.Errorf("The model (%s) doesn't register", name))
	}
	factory.Clone(v)
}

func makeBySchema(query query.Query, schema schema.Schema, v interface{}, flow ...interface{}) *Model {
	name := fmt.Sprintf("%s%s", v, flow)
	factory, has := modelsRegistered[name]
	if !has {
		args := []interface{}{}
		args = append(args, v)
		args = append(args, flow...)
		Register(name, args...)

		factory, has = modelsRegistered[name]
		if !has {
			panic(fmt.Errorf("the model register failure"))
		}
	}
	return factory.Clone()
}

func getTypeName(v interface{}) string {
	return reflect.TypeOf(v).String()
}

// Clone build a model instance quickly
func (factory *Factory) Clone(v ...interface{}) *Model {

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
