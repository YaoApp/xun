package model

import (
	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/dbal/schema"
)

// Model the xun model
type Model struct {
	Query  query.Query
	Schema schema.Schema
}

// Factory xun model factory
type Factory struct {
	Schema *Schema
	Flow   *Flow
	Model  interface{}
}

// Schema the Xun model schema description struct
type Schema struct {
	Name    string  `json:"name"`
	Table   Table   `json:"table"`
	Fields  []Field `json:"fields"`
	Indexes []Index `json:"indexes"`
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
	Name      string `json:"name"`
	Comment   string `json:"comment,omitempty"`
	Engine    string `json:"engine,omitempty"` // InnoDB,MyISAM ( MySQL Only )
	Collation string `json:"collation"`
	Charset   string `json:"charset"`
}
