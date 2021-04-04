package xun

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/yaoapp/xun/utils"
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

// ToFixed the return value is the type of float64 and keeps the given decimal places
func (n N) ToFixed(places int) (float64, error) {
	if n.Value == nil {
		return 0, fmt.Errorf("the value is nil")
	}
	format := "%" + fmt.Sprintf(".%df", places)
	return strconv.ParseFloat(fmt.Sprintf(format, n.Value), 64)
}

// MustToFixed the return value is the type of float64 and keeps the given decimal places
func (n N) MustToFixed(places int) float64 {
	value, err := n.ToFixed(places)
	utils.PanicIF(err)
	return value
}

// Int64 the return value is the type of int64 and remove the decimal
func (n N) Int64() (int64, error) {
	if n.Value == nil {
		return 0, fmt.Errorf("the value is nil")
	}
	return strconv.ParseInt(fmt.Sprintf("%v", n.Value), 10, 64)
}

// MustInt64  the return value is the type of int64 and remove the decimal
func (n N) MustInt64() int64 {
	value, err := n.Int64()
	utils.PanicIF(err)
	return value
}

// Int32 the return value is the type of int64 and remove the decimal
func (n N) Int32() (int32, error) {
	if n.Value == nil {
		return 0, fmt.Errorf("the value is nil")
	}
	value, err := strconv.ParseInt(fmt.Sprintf("%v", n.Value), 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(value), nil
}

// MustInt32  the return value is the type of int64 and remove the decimal
func (n N) MustInt32() int32 {
	value, err := n.Int32()
	utils.PanicIF(err)
	return value
}

// Int the return value is the type of int and remove the decimal
func (n N) Int() (int, error) {
	if n.Value == nil {
		return 0, fmt.Errorf("the value is nil")
	}
	value, err := strconv.ParseInt(fmt.Sprintf("%v", n.Value), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(value), nil
}

// MustInt  the return value is the type of int and remove the decimal
func (n N) MustInt() int {
	value, err := n.Int()
	utils.PanicIF(err)
	return value
}
