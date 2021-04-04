package query

import (
	"github.com/yaoapp/xun/dbal"
)

// Union Add a union statement to the query.
func (builder *Builder) Union(query interface{}, all ...bool) Query {

	isUnionAll := false
	if len(all) > 0 && all[0] == true {
		isUnionAll = true
	}
	var qb *Builder
	switch query.(type) {
	case *Builder:
		qb = query.(*Builder)
		break
	case func(Query):
		callback := query.(func(Query))
		qb = builder.new()
		callback(qb)
		break
	}

	if qb != nil {
		builder.Query.Unions = append(builder.Query.Unions, dbal.Union{
			Query: qb.Query,
			All:   isUnionAll,
		})
		builder.Query.AddBinding("union", qb.GetBindings())
	}
	return builder

}

// UnionAll Add a union all statement to the query.
func (builder *Builder) UnionAll(query interface{}) Query {
	return builder.Union(query, true)
}
