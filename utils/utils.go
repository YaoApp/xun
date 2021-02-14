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

// Set set the value for the pointer
func Set(v interface{}, value interface{}) {
	v = &value
}

// GetInt get the int type pointers value
func GetInt(v *int, defaults ...int) int {
	if v == nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	return *v
}

// GetString get the string type pointers value
func GetString(v *string, defaults ...string) string {
	if v == nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return ""
	}
	return *v
}
