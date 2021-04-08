package utils

import (
	"fmt"
	"reflect"

	"github.com/TylerBrock/colorjson"
	jsoniter "github.com/json-iterator/go"
)

// GetIF if the condition is true return the given value, else return the default value
func GetIF(condition bool, value interface{}, defaultValue interface{}) interface{} {
	if condition {
		return value
	}
	return defaultValue
}

// PanicIF if the given value is not nil then panic, else do nothing
func PanicIF(v interface{}) {
	if v != nil {
		panic(v)
	}
}

// Println pretty print var
func Println(v interface{}) {
	f := colorjson.NewFormatter()
	f.Indent = 4
	var res interface{}
	txt, _ := jsoniter.Marshal(v)
	jsoniter.Unmarshal(txt, &res)
	s, _ := f.Marshal(res)
	fmt.Printf("%s\n", s)
}

// MapFilp flipping the given map
func MapFilp(v interface{}) (interface{}, bool) {
	switch v.(type) {
	case map[string]string:
		res := map[string]string{}
		for key, val := range v.(map[string]string) {
			res[val] = key
		}
		return res, true
	case map[string]int:
		res := map[int]string{}
		for key, val := range v.(map[string]int) {
			res[val] = key
		}
		return res, true
	case map[int]string:
		res := map[string]int{}
		for key, val := range v.(map[int]string) {
			res[val] = key
		}
		return res, true
	default:
		return nil, false
	}
}

// IntPtr get the int value for the pointer
func IntPtr(value int) *int {
	return &value
}

// IntVal get the int type pointers value
func IntVal(v *int, defaults ...int) int {
	if v == nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	return *v
}

// StringPtr get the string value for the pointer
func StringPtr(value string) *string {
	return &value
}

// StringVal get the string type pointers value
func StringVal(v *string, defaults ...string) string {
	if v == nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return ""
	}
	return *v
}

// IntHave Contains method for a slice
func IntHave(source []int, one int) bool {
	for _, a := range source {
		if a == one {
			return true
		}
	}
	return false
}

// StringHave Contains method for a slice
func StringHave(source []string, one string) bool {
	for _, a := range source {
		if a == one {
			return true
		}
	}
	return false
}

// IsNil Check if an interface is nil
func IsNil(value interface{}) bool {
	if value == nil {
		return true
	}

	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() == reflect.Ptr {
		return reflectValue.IsNil()
	}

	return false
}

// Flatten flattens a multi-dimensional array into a single level array:
func Flatten(value interface{}) []interface{} {
	reflectValue := reflect.ValueOf(value)
	reflectValue = reflect.Indirect(reflectValue)
	kind := reflectValue.Kind()
	res := []interface{}{}
	if kind == reflect.Array || kind == reflect.Slice {
		for i := 0; i < reflectValue.Len(); i++ {
			value := reflectValue.Index(i)
			valueKind := value.Kind()
			if valueKind == reflect.Array || valueKind == reflect.Slice {
				res = append(res, Flatten(value.Interface()))
			} else {
				res = append(res, value.Interface())
			}
		}
	} else {
		res = append(res, value)
	}

	return res
}
