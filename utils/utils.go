package utils

import (
	"fmt"

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
