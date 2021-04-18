package model

import (
	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/dbal/schema"
)

// Basic the basic interface of model
type Basic interface {
	ToSQL() string
	Clone() Basic
	BasicSetup(buidler *query.Builder, schema schema.Schema) Basic
	QueryBuilder() *query.Builder
	GetRelationshipLink(name string, ids ...string) (string, string)
	SelectAddColumn(foreign string) Basic
	BasicQueryForRelationship(columns []string, closure func(query.Query)) Basic
	MakeModelForRelationship(name string) Basic
}
