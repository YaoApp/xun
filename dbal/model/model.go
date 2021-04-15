package model

import (
	"fmt"
	"reflect"

	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/dbal/schema"
)

// Make make a new xun model instance
func Make(query query.Query, schema schema.Schema, v interface{}, args ...interface{}) *Model {
	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		makeByStruct(query, schema, v)
		return nil
	}
	return makeBySchema(query, schema, v, args...)
}

// MakeUsing create model using makeer
func MakeUsing(maker MakerFunc, v interface{}, args ...interface{}) *Model {
	return maker(v, args...)
}

// GetFullname get the fullname of model
func (model *Model) GetFullname() string {
	if model.namespace == "" {
		return model.name
	}
	return fmt.Sprintf("%s.%s", model.namespace, model.name)
}

// GetQuery get the query interface
func (model *Model) GetQuery() query.Query {
	return model.query
}

// GetSchema get the query interface
func (model *Model) GetSchema() schema.Schema {
	return model.schema
}

// GetName get the name of model
func (model *Model) GetName() string {
	return model.name
}

// GetNamespace get the name of model
func (model *Model) GetNamespace() string {
	return model.namespace
}

// GetAttributes get all of attribute
func (model *Model) GetAttributes() []Attribute {
	attrs := []Attribute{}
	for _, attr := range model.attributes {
		attrs = append(attrs, attr)
	}
	return attrs
}

// GetAttributeNames get all of the attribute name
func (model *Model) GetAttributeNames() []string {
	names := []string{}
	for name := range model.attributes {
		names = append(names, name)
	}
	return names
}

// GetAttr get the Attribute by name
func (model *Model) GetAttr(name string) Attribute {
	attr, ok := model.attributes[name]
	if !ok {
		panic(fmt.Errorf("The %s attribute does not exists", name))
	}
	return attr
}

// Get get the Attribute value
func (model *Model) Get(name string) interface{} {
	return model.GetAttr(name).Value
}

// Set set the Attribute value
func (model *Model) Set(name string, value interface{}) *Model {
	attr := model.GetAttr(name)
	attr.Value = value
	model.attributes[name] = attr
	return model
}

// Columns get the columns of model  struct
func (model *Model) Columns() {
}

// Indexes get the indexes of model  struct
func (model *Index) Indexes() {
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

// Destroy deleting an dxisting model by its Primary Key
func (model *Model) Destroy(id interface{}) {
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
