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

	// To compile the query, we'll spin through each component of the query and
	// see if that component exists. If it does we'll just call the compiler
	// function for the component which is responsible for making the SQL.

	// sqls["aggregate"] = grammarSQL.compileAggregate()
	sqls["columns"] = grammarSQL.compileColumns(query, query.Columns)
	sqls["from"] = grammarSQL.compileFrom(query, query.From)
	// sqls["joins"] = grammarSQL.compileJoins()
	sqls["wheres"] = grammarSQL.compileWheres(query, query.Wheres)
	// sqls["groups"] = grammarSQL.compileGroups()
	// sqls["havings"] = grammarSQL.compileHavings()
	// sqls["orders"] = grammarSQL.compileOrders()
	// sqls["limit"] = grammarSQL.compileLimit()
	// sqls["offset"] = grammarSQL.compileOffset()
	// sqls["lock"] = grammarSQL.compileLock()

	sql := ""
	for _, name := range []string{"aggregate", "columns", "from", "joins", "wheres", "groups", "havings", "orders", "limit", "offset", "lock"} {
		segment, has := sqls[name]
		if has && segment != "" {
			sql = sql + segment + " "
		}
	}

	// reset columns
	query.Columns = columns
	return strings.Trim(sql, " ")
}

// compileColumns Compile the "select *" portion of the query.
func (grammarSQL SQL) compileColumns(query *dbal.Query, columns []dbal.Name) string {

	// If the query is actually performing an aggregating select, we will let that
	// compiler handle the building of the select clauses, as it will need some
	// more syntax that is best handled by that function to keep things neat.
	// if (len(query.Aggregate) > 0 ) {
	//     return;
	// }

	sql := "select"
	if query.Distinct {
		sql = "select distinct"
	}
	return fmt.Sprintf("%s %s", sql, grammarSQL.Columnize(columns))
}

//  Compile the "from" portion of the query.
func (grammarSQL SQL) compileFrom(query *dbal.Query, table dbal.Name) string {
	if table.As() != "" {
		return fmt.Sprintf("from %s as %s", grammarSQL.ID(table.Fullname()), grammarSQL.ID(table.As()))
	}
	return fmt.Sprintf("from %s", grammarSQL.ID(table.Fullname()))
}

func (grammarSQL SQL) compileWheres(query *dbal.Query, wheres []dbal.Where) string {

	// Each type of where clauses has its own compiler function which is responsible
	// for actually creating the where clauses SQL. This helps keep the code nice
	// and maintainable since each clause has a very small method that it uses.
	if wheres == nil {
		return ""
	}

	clauses := []string{}
	// If we actually have some where clauses, we will strip off the first boolean
	// operator, which is added by the query builders for convenience so we can
	// avoid checking for the first clauses in each of the compilers methods.
	for _, where := range wheres {
		// fmt.Printf("CompileSelect: %s %s %s %v\n", where.Boolean, where.Type, where.Operator, where.Value)
		boolen := strings.ToLower(where.Boolean)
		switch where.Type {
		case "basic":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.whereBasic(where)))
			break
		case "sub":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.whereSub()))
			break
		case "nested":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.whereNested(query, where)))
			break
		}
	}

	// $conjunction = $query instanceof JoinClause ? 'on' : 'where'; ( offset 3 / 6)
	conjunction := "where"
	return fmt.Sprintf("%s %s", conjunction, grammarSQL.RemoveLeadingBoolean(strings.Join(clauses, " ")))
}

func (grammarSQL SQL) whereBasic(where dbal.Where) string {
	value := grammarSQL.Parameter(where.Value)
	operator := strings.ReplaceAll(where.Operator, "?", "??")
	return fmt.Sprintf("%s %s %s", grammarSQL.ID(where.Column), operator, value)
}

func (grammarSQL SQL) whereSub() string {
	return "whereSub"
}

func (grammarSQL SQL) whereNested(query *dbal.Query, where dbal.Where) string {

	// $offset = $query instanceof JoinClause ? 3 : 6;
	offset := 6
	sql := grammarSQL.compileWheres(where.Query, where.Query.Wheres)
	end := len(sql)
	if end > offset {
		sql = sql[offset:end]
	}
	return fmt.Sprintf("(%s)", sql)
}

// Utils for compiling

// RemoveLeadingBoolean Remove the leading boolean from a statement.
func (grammarSQL SQL) RemoveLeadingBoolean(value string) string {
	value = strings.TrimLeft(value, "and ")
	value = strings.TrimLeft(value, "or ")
	return value
}
