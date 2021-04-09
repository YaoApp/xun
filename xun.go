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

// Time Create a new time struct
func Time(value interface{}) T {
	return T{
		Value: value,
	}
}

// AnyToR convert any struct to R type
func AnyToR(v interface{}) R {

	if v == nil {
		return R{}
	} else if res, ok := v.(R); ok {
		return res
	} else if res, ok := v.(map[string]interface{}); ok {
		return res
	}

	res := R{}
	reflectValue := reflect.ValueOf(v)
	reflectValue = reflect.Indirect(reflectValue)
	typ := reflectValue.Type()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	kind := typ.Kind()
	if kind == reflect.Array || kind == reflect.Slice {
		if reflectValue.Len() == 0 {
			return R{}
		}
		return AnyToR(reflectValue.Index(0).Interface())
	}

	if typ.Kind() != reflect.Struct {
		panic(fmt.Errorf("The type of given value is %s, should be struct", typ.String()))
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

// Get get the value of the given key
func (row R) Get(key interface{}) interface{} {
	keyStr := fmt.Sprintf("%v", key)
	value, has := row[keyStr]
	if !has {
		return nil
	}
	return value
}

// MustGet get the value of the given key
func (row R) MustGet(key interface{}) interface{} {
	keyStr := fmt.Sprintf("%v", key)
	value, has := row[keyStr]
	if !has {
		panic(fmt.Errorf("the key %v does not exists", key))
	}
	return value
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
	values := AnyToRows(v)
	for _, value := range values {
		for k, v := range value {
			(*row)[k] = v
		}
	}
}
