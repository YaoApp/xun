package model

import (
	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/dbal/schema"
)

// Model the xun model
type Model struct {
	Query      query.Query
	Schema     schema.Schema
	Attributes []Attribute
}

// Factory xun model factory
type Factory struct {
	Schema  *Schema
	Flow    *Flow
	Methods []Method
	Model   interface{}
}

// Schema the Xun model schema description struct
type Schema struct {
	Name          string         `json:"name"`
	Table         Table          `json:"table"`
	Fields        []Field        `json:"fields"`
	Indexes       []Index        `json:"indexes"`
	Relationships []Relationship `json:"relationships"`
	Option        SchemaOption   `json:"option"`
}

// Flow the Xun model flow description struct
type Flow struct {
	Name string `json:"name"`
	Node Node   `json:"node"`
}

// Node the flow node struct
type Node struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Children    []Node `json:"children"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Comment     string `json:"comment,omitempty"`
}

// SchemaOption The shcmea option
type SchemaOption struct {
	SoftDeletes bool `json:"soft_deletes"`
	Timestamps  bool `json:"timestamps"`
}

// Field the field description struct
type Field struct {
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Comment     string       `json:"comment,omitempty"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Length      *int         `json:"length,omitempty"`
	Precision   *int         `json:"precision,omitempty"`
	Scale       *int         `json:"scale,omitempty"`
	Nullable    *bool        `json:"nullable,omitempty"`
	Option      []string     `json:"option,omitempty"`
	Default     interface{}  `json:"default,omitempty"`
	Example     interface{}  `json:"example,omitempty"`
	Generate    string       `json:"generate,omitempty"` // Increment, UUID,...
	Encoder     string       `json:"encoder,omitempty"`  // AES-256, AES-128, PASSWORD-HASH, ...
	Decoder     string       `json:"decoder,omitempty"`  // AES-256, AES-128, ...
	Validations []Validation `json:"validations,omitempty"`
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
	Fields  []string `json:"fields,omitempty"`
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
}

// Attribute the model attribute
type Attribute struct {
	Name         string
	Field        *Field
	Relationship *Relationship
}

// Method the method can be exported
type Method struct {
	Name   string
	Path   string
	In     []interface{}
	Out    []interface{}
	Export bool
}

// MakerFunc the function for create a model
type MakerFunc func(v interface{}, flow ...interface{}) *Model
