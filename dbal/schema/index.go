package schema

import (
	"errors"
)

var indexTypes = map[string]string{
	"unique":  "UNIQUE",
	"index":   "INDEX",
	"primary": "PRIMARY",
}

// Drop mark as dropped for the index
func (index *Index) Drop() {
	index.dropped = true
}

// Rename mark as renamed with the given name for the index
func (index *Index) Rename(new string) {
	index.renamed = true
	index.newname = new
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
