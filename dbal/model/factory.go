package model

import (
	"fmt"
	"reflect"
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

// Class get the factory by model name
func Class(name interface{}) *Factory {
	factory, has := modelsRegistered[name]
	if !has {
		panic(fmt.Errorf("The model (%s) doesn't register", name))
	}
	return factory
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
