package model

import (
	"fmt"
	"reflect"
)

// modelsRegistered the models have been registered
var modelsRegistered map[string]*Factory = map[string]*Factory{}
var modelsAlias map[string]*Factory = map[string]*Factory{}

// Register register the model
func Register(v interface{}, args ...interface{}) {

	reflectPtr := reflect.ValueOf(v)
	reflectValue := reflect.Indirect(reflectPtr)
	if reflectValue.Kind() == reflect.String && len(args) > 0 {
		origin := v.(string)
		fullname, namespace, name := prepareRegisterNames(origin)
		schema, flow := prepareRegisterArgs(args...)
		model := Model{}
		model.namespace = namespace
		model.name = name
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
		return
	} else if reflectPtr.Kind() == reflect.Ptr && reflectValue.Kind() == reflect.Struct && reflectValue.FieldByName("Model").Type() == reflect.TypeOf(Model{}) {
		origin := reflectPtr.Type().String()
		fullname, namespace, name := prepareRegisterNames(origin)
		schema, flow := prepareRegisterArgs(args...)
		SetModel(v, func(model *Model) {
			model.namespace = namespace
			model.name = name
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
		if reflectPtr.Kind() == reflect.Ptr && reflectValue.Kind() == reflect.Struct && reflectValue.FieldByName("Model").Type() == reflect.TypeOf(Model{}) {
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
		if reflectPtr.Kind() == reflect.Ptr && reflectValue.Kind() == reflect.Struct && reflectValue.FieldByName("Model").Type() == reflect.TypeOf(Model{}) {
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
