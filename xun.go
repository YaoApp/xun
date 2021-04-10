package xun

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

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

// UpperFirst upcase the first letter
func UpperFirst(str string) string {
	if len(str) < 1 {
		return strings.ToUpper(str)
	}
	first := strings.ToUpper(string(str[0]))
	other := str[1:]
	return first + other
}

// MapScan scan the result from sql.Rows
func MapScan(rows *sql.Rows) ([]R, error) {
	defer rows.Close()
	res := []R{}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	numColumns := len(columns)

	values := make([]interface{}, numColumns)
	for i := range values {
		values[i] = new(interface{})
	}

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}

		dest := R{}
		for i, column := range columns {
			reflectValue := reflect.ValueOf(values[i])
			value := reflect.Indirect(reflectValue).Interface()
			switch value.(type) {
			case []byte:
				bytes := ""
				for _, v := range value.([]byte) {
					bytes = fmt.Sprintf("%s%s", bytes, string(v))
				}

				dest[column] = bytes
				if len(bytes) < 20 {
					intv, err := strconv.ParseInt(bytes, 10, 64)
					if err == nil {
						dest[column] = intv
						break
					}

					floatv, err := strconv.ParseFloat(bytes, 64)
					if err == nil {
						dest[column] = floatv
						break
					}
				}

			default:
				dest[column] = value
			}

		}
		res = append(res, dest)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// MakeTime Create a new time struct
func MakeTime(value ...interface{}) T {
	if len(value) == 0 {
		return T{Value: time.Now()}
	}
	return T{
		Value: value[0],
	}
}

// ToTime cast the T to time.Time
func (t T) ToTime(formats ...string) (time.Time, error) {
	if len(formats) == 0 {
		formats = []string{
			"2006-01-02T15:04:05-0700",
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
			"15:04:05",
		}
	}

	switch t.Value.(type) {
	case int, int64, int32, int16, int8, uint8:
		var err error
		var i int64
		var s int64
		strValue := fmt.Sprintf("%v", t.Value)
		if len(strValue) == 10 {
			i, err = strconv.ParseInt(strValue, 10, 64)
			s = 0
			if err != nil {
				return time.Now(), err
			}
		} else if len(strValue) == 13 {
			i, err = strconv.ParseInt(strValue[0:10], 10, 64)
			if err != nil {
				return time.Now(), err
			}
			s, err = strconv.ParseInt(strValue[10:13], 10, 64)
			if err != nil {
				return time.Now(), err
			}
		}
		return time.Unix(i, s), nil

	case string, []byte:
		var err error
		strValue := fmt.Sprintf("%s", t.Value)
		dateValue := time.Now()
		for _, format := range formats {
			dateValue, err = time.Parse(format, strValue)
			if err == nil {
				return dateValue, nil
			}
		}
		if err != nil {
			return dateValue, fmt.Errorf("%s(%s)", err, formats)
		}
		return dateValue, fmt.Errorf("cannot parse %s (%s)", t.Value, formats)
	case time.Time:
		return t.Value.(time.Time), nil
	default:
		return time.Now(), nil
	}
}

// MustToTime cast the T to time.Time
func (t T) MustToTime(formats ...string) time.Time {
	value, err := t.ToTime(formats...)
	if err != nil {
		panic(err)
	}
	return value
}

// MakeR create a new R struct
func MakeR(value ...interface{}) R {
	if len(value) == 0 {
		return R{}
	}
	v := value[0]
	reflectValue := reflect.ValueOf(v)
	reflectValue = reflect.Indirect(reflectValue)
	switch reflectValue.Kind() {
	case reflect.Slice, reflect.Array:
		if reflectValue.Len() > 0 {
			return MakeR(reflectValue.Index(0).Interface())
		}
		return MakeR()
	case reflect.Map:
		return mapToR(v)
	case reflect.Struct:
		return structToR(v)
	}

	panic(fmt.Errorf("The type of given value is %s, should be struct", reflectValue.Type().String()))
}

// Get get the value of the given key
func (row R) Get(key interface{}) interface{} {
	keys := strings.Split(fmt.Sprintf("%v", key), ".")
	nextRow := row
	length := len(keys) - 1
	for i, k := range keys {
		value, has := nextRow[k]
		if !has {
			panic(fmt.Errorf("the key %v does not exists", key))
		}
		if length == i {
			return value
		}
		nextRow = MakeR(value)
	}
	return nil
}

// ToMap cast to map[string]interface{}
func (row R) ToMap() map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range row {
		res[k] = v
	}
	return res
}

// Keys get keys of R
func (row R) Keys() []interface{} {
	keys := []interface{}{}
	for k := range row {
		keys = append(keys, k)
	}
	return keys
}

// KeysString get keys of R
func (row R) KeysString() []string {
	keys := []string{}
	for k := range row {
		keys = append(keys, k)
	}
	return keys
}

// Merge get keys of R
func (row *R) Merge(v ...interface{}) {
	values := MakeRSlice(v...)
	for _, value := range values {
		for k, v := range value {
			(*row)[k] = v
		}
	}
}

// MakeRows alias MakeRSlice
func MakeRows(value ...interface{}) []R {
	return MakeRSlice(value...)
}

// MakeRSlice convert any struct to R slice
func MakeRSlice(value ...interface{}) []R {
	if len(value) == 0 {
		return []R{}
	}

	values := []interface{}{}
	if len(value) == 1 {
		reflectValue := reflect.ValueOf(value[0])
		reflectValue = reflect.Indirect(reflectValue)
		reflectKind := reflectValue.Kind()
		if reflectKind == reflect.Slice || reflectKind == reflect.Array {
			for i := 0; i < reflectValue.Len(); i++ {
				values = append(values, reflectValue.Index(i).Interface())
			}
		} else {
			values = append(values, value[0])
		}
	} else {
		values = value
	}

	res := []R{}
	for _, v := range values {
		res = append(res, MakeR(v))
	}
	return res
}

func mapToR(value interface{}) R {
	r := R{}
	reflectValue := reflect.ValueOf(value)
	reflectValue = reflect.Indirect(reflectValue)
	if reflectValue.Kind() == reflect.Map {
		for _, key := range reflectValue.MapKeys() {
			k := fmt.Sprintf("%v", key)
			v := reflectValue.MapIndex(key).Interface()
			r[k] = v
		}
	}
	return r
}

func structToR(value interface{}) R {
	r := R{}
	reflectValue := reflect.ValueOf(value)
	reflectValue = reflect.Indirect(reflectValue)
	reflectType := reflectValue.Type()
	if reflectValue.Kind() == reflect.Struct {
		for i := 0; i < reflectValue.NumField(); i++ {
			tag := GetTagName(reflectType.Field(i), "json")
			field := reflectValue.Field(i).Interface()
			if tag != "" && tag != "-" {
				kind := reflectType.Field(i).Type.Kind()
				if kind == reflect.Struct {
					r[tag] = structToR(field)
				} else if kind == reflect.Slice || kind == reflect.Array {
					r[tag] = MakeRSlice(field)
				} else {
					r[tag] = field
				}
			}
		}
	}
	return r
}

// GetTagName get the tag name of the reflect.StructField
func GetTagName(field reflect.StructField, name string) string {
	tag := field.Tag.Get(name)
	if tag == "" {
		tag = ToSnakeCase(field.Name)
	}
	return tag
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
