package model

import (
	"fmt"
	"reflect"

	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/dbal/schema"
)

// Make make a new xun model instance
func Make(query query.Query, schema schema.Schema, v interface{}, flow ...interface{}) *Model {
	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		makeByStruct(query, schema, v)
		return nil
	}
	return makeBySchema(query, schema, v, flow...)
}

// MakeUsing create model using makeer
func MakeUsing(maker MakerFunc, v interface{}, flow ...interface{}) *Model {
	return maker(v, flow...)
}

// Fill to fill attributes into model
func (model *Model) Fill(attributes interface{}) {
}

// Create to create one model
func (model *Model) Create(attributes interface{}) {
}

// Save to create or update one model
func (model *Model) Save() {
}

// Destory deleting an dxisting model by its Primary Key
func (model *Model) Destory(id interface{}) {
}

// Restore To restore a soft deleted model,
func (model *Model) Restore() {
}

// With where the array key is a relationship name and the array value is a closure that adds additional constraints to the eager loading query
func (model *Model) With() {
}

// Where same as the query where, return the query builder
func (model *Model) Where() {
}

// Insert same as the query insert
func (model *Model) Insert(v interface{}, columns ...interface{}) {
}

// Update  same as the query update
func (model *Model) Update() {
}

// Upsert same as the query upsert
func (model *Model) Upsert() {
}

// UpdateOrInsert same as the query UpdateOrInsert
func (model *Model) UpdateOrInsert() {
}

// Delete same as the query Delete
func (model *Model) Delete() {
}

// Truncate same as the query Truncate
func (model *Model) Truncate() {
}

// Chunk same as the query Chunk
func (model *Model) Chunk() {
}

// Paginate same as the query Paginate
func (model *Model) Paginate() {
}

// Search search by given params
func (model *Model) Search() interface{} {
	return nil
}

// Find find by primary key
func (model *Model) Find(v ...interface{}) interface{} {
	if len(v) > 0 {
		return v[0]
	}
	return model
}

// Export export data
func (model *Model) Export() {
}

// Import import data
func (model *Model) Import() {
}

// Flow process a flow by the given flow name and return the result
func (model *Model) Flow(name string) interface{} {
	return nil
}

// FlowRaw process a flow by the given json description file and return the result
func (model *Model) FlowRaw(flow []byte) interface{} {
	return nil
}

// makeBySchema make a new xun model instance
func makeBySchema(query query.Query, schema schema.Schema, v interface{}, flow ...interface{}) *Model {
	name := fmt.Sprintf("%s%s", v, flow)
	class, has := modelsRegistered[name]
	if !has {
		args := []interface{}{}
		args = append(args, v)
		args = append(args, flow...)
		Register(name, args...)

		class, has = modelsRegistered[name]
		if !has {
			panic(fmt.Errorf("the model register failure"))
		}
	}
	return class.New()
}

// makeByStruct make a new xun model instance
func makeByStruct(query query.Query, schema schema.Schema, v interface{}) {
	name := getTypeName(v)
	Class(name).New(v)
}
