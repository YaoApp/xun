package model

import (
	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/dbal/schema"
)

// Basic the basic interface of model
type Basic interface {
	ToSQL() string
	Clone() Basic

	// defined in the file model.go
	BasicSetup(buidler *query.Builder, schema schema.Schema) Basic

	GetName() string
	GetNamespace() string
	GetFullname() string
	GetTableName() string

	SelectAddColumn(column string) Basic
	TableColumnize(table string)

	QueryBuilder() *query.Builder
	GetRelationshipLink(name string, ids ...string) (string, string)
	BasicQueryForRelationship(columns []string, closure func(Basic)) Basic
	MakeModelForRelationship(name string) Basic
}
