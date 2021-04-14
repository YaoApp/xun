package model

import (
	"reflect"
	"strings"

	"github.com/yaoapp/xun"
)

// GetFullname get the fullname of model
func (column *Column) set(name string, v string) {
	name = xun.UpperFirst(name)
	reflectValue := reflect.ValueOf(column)
	reflectValue = reflect.Indirect(reflectValue)
	field := reflectValue.FieldByName(name)
	kind := field.Kind()
	if kind != reflect.Invalid {
		if kind == reflect.Slice {
			values := strings.Split(v, ",")
			for i := range values {
				values[i] = strings.Trim(values[i], " ")
			}
			field.Set(reflect.ValueOf(values))
		} else if kind == reflect.String {
			field.Set(reflect.ValueOf(v))
		}
	}
}

func (column *Column) merge(columns ...Column) *Column {
	for _, new := range columns {
		column.setString(&column.Name, new.Name)
		column.setString(&column.Type, new.Type)
		column.setString(&column.Title, new.Title)
		column.setString(&column.Comment, new.Comment)
		column.setString(&column.DefaultRaw, new.DefaultRaw)
		column.setString(&column.Description, new.Description)
		column.setInterface(&column.Default, new.Default)
		column.setStringSlice(&column.Option, new.Option)
	}
	return column
}

func (column *Column) setString(v *string, value string) {
	if value != "" {
		*v = value
	}
}

func (column *Column) setInterface(v *interface{}, value interface{}) {
	if value != nil {
		*v = value
	}
}

func (column *Column) setStringSlice(v *[]string, value []string) {
	if value != nil {
		*v = value
	}
}
