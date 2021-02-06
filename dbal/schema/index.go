package schema

import (
	"errors"
	"fmt"
	"strings"
)

var indexTypes = map[string]string{
	"unique":  "UNIQUE",
	"index":   "INDEX",
	"primary": "PRIMARY",
}

func (index *Index) sqlCreate() string {
	// UNIQUE KEY `unionid` (`unionid`)
	columns := []string{}
	for _, column := range index.Columns {
		columns = append(columns, column.Name)
	}
	sql := fmt.Sprintf("%s KEY `%s` (`%s`)", GetIndexType(index.Type), index.Name, strings.Join(columns, "`,`"))
	return sql
}

// GetIndexType get the index type
func GetIndexType(name string) string {
	if _, has := indexTypes[name]; has {
		return indexTypes[name]
	}
	return "INDEX"
}

func (index *Index) validate() *Index {
	if index.Name == "" {
		err := errors.New("the index name must be set")
		panic(err)
	}

	if len(index.Columns) == 0 {
		err := errors.New("the index " + index.Name + " must have one column at least")
		panic(err)
	}

	if index.Table == nil {
		err := errors.New("the index " + index.Name + "not bind table")
		panic(err)
	}

	return index
}
