package model

import (
	"fmt"
	"strings"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/utils"
)

// MakeBasic create a new xun model basic interface
func (model *Model) MakeBasic(name string) Basic {
	return Class(name).
		NewBasic().
		BasicSetup(model.Builder, model.schema)
}

// BasicSetup get the query builder pointer
func (model *Model) BasicSetup(buidler *query.Builder, schema schema.Schema) Basic {
	model.Builder = buidler.NewBuilder()
	model.schema = schema
	return model
}

// SelectAddColumn add a column
func (model *Model) SelectAddColumn(column string, table ...string) Basic {
	columns := model.Query.Columns
	if len(columns) == 0 {
		all := "*"
		if len(table) > 0 {
			all = fmt.Sprintf("%s.*", table[0])
		}
		columns = append(columns, all)
	}
	columns = append(columns, column)
	columns = utils.InterfaceUnique(columns)
	model.Select(columns)
	return model
}

// TableColumnize parse columns name with table (id -> car.id)
func (model *Model) TableColumnize(table string) {
	columns := model.Query.Columns
	for i, column := range columns {
		if name, ok := column.(string); ok && !strings.Contains(name, ".") {
			columns[i] = fmt.Sprintf("%s.%s", table, name)
		}
	}
	model.Query.Columns = columns
}

// MakeModelForRelationship make a new model instance for the relationship query
func (model *Model) MakeModelForRelationship(name string) Basic {
	relFullname := name
	if !strings.Contains(relFullname, ".") {
		relFullname = fmt.Sprintf("%s.%s", model.namespace, relFullname)
	}
	return model.MakeBasic(relFullname)
}

// BasicQueryForRelationship execute basic query for relationship
func (model *Model) BasicQueryForRelationship(columns []string, closure func(Basic)) Basic {

	if closure != nil {
		closure(model)
	} else if columns != nil {
		model.Select(columns)
	}

	model.BasicQuery()
	return model
}

// GetTableName get the table name
func (model *Model) GetTableName() string {
	return model.table.Name
}

// GetRelationshipLink Get the Relationship local and foreign
func (model *Model) GetRelationshipLink(name string, ids ...string) (string, string) {
	foreign := "id"
	local := fmt.Sprintf("%s_id", strings.ToLower(name))
	if len(ids) == 2 {
		local = ids[0]
		foreign = ids[1]
	}
	return local, foreign
}

// QueryBuilder get the query builder pointer
func (model *Model) QueryBuilder() *query.Builder {
	return model.Builder
}

// Clone get the query builder pointer
func (model *Model) Clone() Basic {
	clone := *model
	clone.values = xun.MakeRow()
	return &clone
}
