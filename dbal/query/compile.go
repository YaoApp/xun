package query

import (
	"fmt"
	"strings"
)

// ToSQL Get the SQL representation of the query.
func (builder *Builder) ToSQL() string {
	return builder.compileSelect()
}

func (builder *Builder) compileSelect() string {

	sqls := map[string]string{}

	// If the query does not have any columns set, we'll set the columns to the
	// * character to just get all of the columns from the database. Then we
	// can build the query and concatenate all the pieces together as one.
	columns := builder.Attr.Columns
	if len(columns) == 0 {
		builder.AddColumn("*")
	}

	// sqls["aggregate"] = builder.compileAggregate()
	// sqls["columns"] = builder.compileColumns()
	// sqls["from"] = builder.compileFrom()
	// sqls["joins"] = builder.compileJoins()
	sqls["wheres"] = builder.compileWheres()
	// sqls["groups"] = builder.compileGroups()
	// sqls["havings"] = builder.compileHavings()
	// sqls["orders"] = builder.compileOrders()
	// sqls["limit"] = builder.compileLimit()
	// sqls["offset"] = builder.compileOffset()
	// sqls["lock"] = builder.compileLock()

	// reset columns
	builder.Attr.Columns = columns
	return fmt.Sprintf("%s", sqls)
}

func (builder *Builder) compileWheres() string {
	// Each type of where clauses has its own compiler function which is responsible
	// for actually creating the where clauses SQL. This helps keep the code nice
	// and maintainable since each clause has a very small method that it uses.
	if builder.Attr.Wheres == nil {
		return ""
	}

	clauses := []string{}
	// If we actually have some where clauses, we will strip off the first boolean
	// operator, which is added by the query builders for convenience so we can
	// avoid checking for the first clauses in each of the compilers methods.
	for _, where := range builder.Attr.Wheres {
		switch where.Type {
		case "basic":
			clauses = append(clauses, fmt.Sprintf("%s %s", where.Boolean, builder.whereBasic(where)))
			break
		case "sub":
			clauses = append(clauses, fmt.Sprintf("%s %s", where.Boolean, builder.whereSub()))
			break
		case "nested":
			clauses = append(clauses, fmt.Sprintf("%s %s", where.Boolean, builder.whereNested()))
			break
		}
	}

	return strings.Join(clauses, " ")
}

func (builder *Builder) whereBasic(where Where) string {
	return fmt.Sprintf("%s %s %v", where.Column, where.Operator, where.Value)
}

func (builder *Builder) whereSub() string {
	return "whereSub"
}

func (builder *Builder) whereNested() string {
	return "whereNested"
}
