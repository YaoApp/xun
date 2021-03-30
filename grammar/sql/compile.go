package sql

import (
	"fmt"
	"strings"

	"github.com/yaoapp/xun/dbal"
)

// CompileSelect Compile a select query into SQL.
func (grammarSQL SQL) CompileSelect(query *dbal.Query) string {

	sqls := map[string]string{}

	// If the query does not have any columns set, we'll set the columns to the
	// * character to just get all of the columns from the database. Then we
	// can build the query and concatenate all the pieces together as one.
	columns := query.Columns
	if len(columns) == 0 {
		query.AddColumn("*")
	}

	// sqls["aggregate"] = builder.compileAggregate()
	// sqls["columns"] = builder.compileColumns()
	// sqls["from"] = builder.compileFrom()
	// sqls["joins"] = builder.compileJoins()
	sqls["wheres"] = grammarSQL.compileWheres(query)
	// sqls["groups"] = builder.compileGroups()
	// sqls["havings"] = builder.compileHavings()
	// sqls["orders"] = builder.compileOrders()
	// sqls["limit"] = builder.compileLimit()
	// sqls["offset"] = builder.compileOffset()
	// sqls["lock"] = builder.compileLock()

	// reset columns
	query.Columns = columns
	return fmt.Sprintf("%s", sqls)
}

func (grammarSQL SQL) compileWheres(query *dbal.Query) string {
	// Each type of where clauses has its own compiler function which is responsible
	// for actually creating the where clauses SQL. This helps keep the code nice
	// and maintainable since each clause has a very small method that it uses.
	if query.Wheres == nil {
		return ""
	}

	clauses := []string{}
	// If we actually have some where clauses, we will strip off the first boolean
	// operator, which is added by the query builders for convenience so we can
	// avoid checking for the first clauses in each of the compilers methods.
	for _, where := range query.Wheres {
		switch where.Type {
		case "basic":
			clauses = append(clauses, fmt.Sprintf("%s %s", where.Boolean, grammarSQL.whereBasic(where, "bindings")))
			break
		case "sub":
			clauses = append(clauses, fmt.Sprintf("%s %s", where.Boolean, grammarSQL.whereSub()))
			break
		case "nested":
			clauses = append(clauses, fmt.Sprintf("%s %s", where.Boolean, grammarSQL.whereNested()))
			break
		}
	}

	return strings.Join(clauses, " ")
}

func (grammarSQL SQL) whereBasic(where dbal.Where, bindings interface{}) string {
	return fmt.Sprintf("%s %s %v", where.Column, where.Operator, bindings)
}

func (grammarSQL SQL) whereSub() string {
	return "whereSub"
}

func (grammarSQL SQL) whereNested() string {
	return "whereNested"
}
