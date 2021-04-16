package model

import (
	"reflect"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/dbal/schema"
)

// Model the xun model
type Model struct {
	namespace   string
	name        string
	attributes  map[string]Attribute
	values      xun.R
	table       *Table
	columns     []*Column
	searchable  []string
	primary     string
	primaryKeys []string
	uniqueKeys  []string
	columnNames []string
	softDeletes bool
	Timestamps  bool
	query       query.Query
	schema      schema.Schema
}

// Factory xun model factory
type Factory struct {
	Schema    *Schema
	Flow      *Flow
	methods   []Method
	Namespace string
	Name      string
	Model     interface{}
}

// Schema the Xun model schema description struct
type Schema struct {
	Name          string         `json:"name"`
	Table         Table          `json:"table,omitempty"`
	Columns       []Column       `json:"columns,omitempty"`
	Indexes       []Index        `json:"indexes,omitempty"`
	Relationships []Relationship `json:"relationships,omitempty"`
	Option        SchemaOption   `json:"option,omitempty"`
	Values        []xun.R        `json:"values,omitempty"`
}

// Flow the Xun model flow description struct
type Flow struct {
	Name string `json:"name"`
	Node Node   `json:"node"`
}

// Node the flow node struct
type Node struct {
	Name        string `json:"name"`
	Method      Method `json:"method"`
	Type        string `json:"type,omitempty"`
	Children    []Node `json:"children,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Comment     string `json:"comment,omitempty"`
}

// SchemaOption The shcmea option
type SchemaOption struct {
	SoftDeletes bool `json:"soft_deletes"`
	Timestamps  bool `json:"timestamps"`
}

// Column the field description struct
type Column struct {
	Name        string       `json:"name"`
	Type        string       `json:"type,omitempty"`
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Comment     string       `json:"comment,omitempty"`
	Length      int          `json:"length,omitempty"`
	Precision   int          `json:"precision,omitempty"`
	Scale       int          `json:"scale,omitempty"`
	Nullable    bool         `json:"nullable,omitempty"`
	Option      []string     `json:"option,omitempty"`
	Default     interface{}  `json:"default,omitempty"`
	DefaultRaw  string       `json:"default_raw,omitempty"`
	Example     interface{}  `json:"example,omitempty"`
	Generate    string       `json:"generate,omitempty"` // Increment, UUID,...
	Encoder     string       `json:"encoder,omitempty"`  // AES-256, AES-128, PASSWORD-HASH, ...
	Decoder     string       `json:"decoder,omitempty"`  // AES-256, AES-128, ...
	Validations []Validation `json:"validations,omitempty"`
	Index       bool         `json:"index,omitempty"`
	Unique      bool         `json:"unique,omitempty"`
	Primary     bool         `json:"primary,omitempty"`
	Field       string       // the field name
}

// Validation the field validation struct
type Validation struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	Pattern     string `json:"pattern,omitempty"`
}

// Index the search index struct
type Index struct {
	Comment string   `json:"comment,omitempty"`
	Name    string   `json:"name,omitempty"`
	Columns []string `json:"columns,omitempty"`
	Type    string   `json:"string"` // primary,unique,index,match
}

// Table the model mapping table in DB
type Table struct {
	Name        string   `json:"name"`
	Comment     string   `json:"comment,omitempty"`
	Engine      string   `json:"engine,omitempty"` // InnoDB,MyISAM ( MySQL Only )
	Collation   string   `json:"collation"`
	Charset     string   `json:"charset"`
	PrimaryKeys []string `json:"primarykeys"`
}

// Relationship types
const (
	RelHasOne         = "hasOne"         // 1 v 1
	RelHasMany        = "hasMany"        // 1 v n
	RelBelongsTo      = "belongsTo"      // inverse  1 v 1 / 1 v n / n v n
	RelHasOneThrough  = "hasOneThrough"  // 1 v 1 ( t1 <-> t2 <-> t3)
	RelHasManyThrough = "hasManyThrough" // 1 v n ( t1 <-> t2 <-> t3)
	RelBelongsToMany  = "belongsToMany"  // 1 v1 / 1 v n / n v n
	RelMorphOne       = "morphOne"
	RelMorphMany      = "morphMany"
	RelMorphToMany    = "morphToMany"
	RelMorphByMany    = "morphByMany"
	RelMorphMap       = "morphMap"
)

// Relationship xun model relationships
type Relationship struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Models []string `json:"models"`
	Links  []string `json:"links,omitempty"` //  M1.Local, (->) M2.Foreign, M2.Local, (->) M3.Foreign ...
}

// Attribute the model attribute
type Attribute struct {
	Name         string
	Column       *Column
	Relationship *Relationship
}

// Method the method can be exported
type Method struct {
	Name    string   `json:"name"`
	Path    string   `json:"path,omitempty"`
	Process string   `json:"process,omitempty"`
	In      []string `json:"in,omitempty"`
	Out     []string `json:"out,omitempty"`
	Export  bool     `json:"export,omitempty"`
}

// MakerFunc the function for create a model
type MakerFunc func(v interface{}, flow ...interface{}) *Model

// StructMapping golang struct mapping to json schema
var StructMapping = map[reflect.Kind]Column{
	reflect.Bool:    {Type: "boolean"},
	reflect.Int8:    {Type: "integer"},
	reflect.Int16:   {Type: "integer"},
	reflect.Int32:   {Type: "integer"},
	reflect.Int64:   {Type: "bigInteger"},
	reflect.Int:     {Type: "bigInteger"},
	reflect.Uint8:   {Type: "unsignedInteger"},
	reflect.Uint16:  {Type: "unsignedInteger"},
	reflect.Uint32:  {Type: "unsignedInteger"},
	reflect.Uint64:  {Type: "unsignedBigInteger"},
	reflect.Uint:    {Type: "unsignedBigInteger"},
	reflect.Float32: {Type: "float", Precision: 8, Scale: 2},
	reflect.Float64: {Type: "float", Precision: 16, Scale: 4},
	reflect.String:  {Type: "string", Length: 200},
	reflect.Struct:  {Type: "json"},
}
