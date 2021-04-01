package xun

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase convert camel case string to snake case
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// AnyToR convert any struct to R type
func AnyToR(v interface{}) R {
	res := R{}
	if v == nil {
		return res
	}
	typ := reflect.TypeOf(v)
	reflectValue := reflect.ValueOf(v)
	reflectValue = reflect.Indirect(reflectValue)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag == "" {
			tag = ToSnakeCase(typ.Field(i).Name)
		}
		if tag != "" && tag != "-" {
			kind := typ.Field(i).Type.Kind()
			if kind == reflect.Struct {
				res[tag] = AnyToR(field)
			} else if kind == reflect.Slice || kind == reflect.Array {
				res[tag] = AnyToRs(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}

// AnyToRs convert any struct to R slice
func AnyToRs(v interface{}) []R {
	res := []R{}
	typ := reflect.TypeOf(v)
	kind := typ.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return res
	}
	reflectValue := reflect.ValueOf(v)
	for i := 0; i < reflectValue.Len(); i++ {
		res = append(res, AnyToR(reflectValue.Index(i).Interface()))
	}
	return res
}

// AnyToRows convert any inteface to R slice
func AnyToRows(v interface{}) []R {
	values := []R{}
	switch v.(type) {
	case R:
		values = append(values, v.(R))
		break
	case []R:
		values = v.([]R)
		break
	default:
		typ := reflect.TypeOf(v)
		kind := typ.Kind()
		if kind == reflect.Slice || kind == reflect.Array {
			values = AnyToRs(v)
		} else {
			values = append(values, AnyToR(v))
		}
	}
	return values
}

// MapToR map[string]inteface{} to R{}, and cast []int8 to string
func MapToR(row map[string]interface{}) R {
	res := R{}
	for key, value := range row {
		switch value.(type) {
		case []uint8:
			bytes := ""
			for _, v := range value.([]uint8) {
				bytes = fmt.Sprintf("%s%s", bytes, string(v))
			}
			res[key] = bytes
		default:
			res[key] = value
		}
	}
	return res
}
