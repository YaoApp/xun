package model

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// prepareRegisterArgs parse the params for Register()
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

func getTypeName(v interface{}) string {
	return reflect.TypeOf(v).String()
}
