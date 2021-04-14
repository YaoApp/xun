package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// prepareRegisterNames parse name and return the namesapce, name
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
