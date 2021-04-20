package model

import (
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/query"
)

// Reset Reset()
func (model *Model) Reset() *Model {
	model.Builder.Reset()
	model.withs = []With{}
	model.values = xun.MakeRow()
	return model
}

// UseRead UseRead()
func (model *Model) UseRead() *Model {
	model.Builder.UseRead()
	return model
}

// UseWrite UseWrite()
func (model *Model) UseWrite() *Model {
	model.Builder.UseWrite()
	return model
}

// Select Select()
func (model *Model) Select(columns ...interface{}) *Model {
	model.Builder.Select(columns...)
	return model
}

// SelectRaw Select()
func (model *Model) SelectRaw(expression string, bindings ...interface{}) *Model {
	model.Builder.SelectRaw(expression, bindings...)
	return model
}

// SelectSub SelectSub()
func (model *Model) SelectSub(qb interface{}, alias string) *Model {
	model.Builder.SelectSub(qb, alias)
	return model
}

// Distinct Distinct()
func (model *Model) Distinct(args ...interface{}) *Model {
	model.Builder.Distinct(args...)
	return model
}

// From  From()
func (model *Model) From(name string) *Model {
	model.Builder.From(name)
	return model
}

// FromRaw FromRaw()
func (model *Model) FromRaw(sql string, bindings ...interface{}) *Model {
	model.Builder.FromRaw(sql, bindings...)
	return model
}

// FromSub FromSub()
func (model *Model) FromSub(qb interface{}, alias string) *Model {
	model.Builder.FromSub(qb, alias)
	return model
}

// Union  Union()
func (model *Model) Union(query interface{}, all ...bool) *Model {
	model.Builder.Union(query, all...)
	return model
}

// UnionAll  UnionAll()
func (model *Model) UnionAll(query interface{}) *Model {
	model.Builder.UnionAll(query)
	return model
}

// Join Join()
func (model *Model) Join(table string, first interface{}, args ...interface{}) *Model {
	model.Builder.Join(table, first, args...)
	return model
}

// JoinSub JoinSub()
func (model *Model) JoinSub(qb interface{}, alias string, first interface{}, args ...interface{}) *Model {
	model.Builder.JoinSub(qb, alias, first, args...)
	return model
}

// LeftJoin LeftJoin()
func (model *Model) LeftJoin(table string, first interface{}, args ...interface{}) *Model {
	model.Builder.LeftJoin(table, first, args...)
	return model
}

// LeftJoinSub LeftJoinSub()
func (model *Model) LeftJoinSub(qb interface{}, alias string, first interface{}, args ...interface{}) *Model {
	model.Builder.LeftJoinSub(qb, alias, first, args...)
	return model
}

// RightJoin RightJoin()
func (model *Model) RightJoin(table string, first interface{}, args ...interface{}) *Model {
	model.Builder.RightJoin(table, first, args...)
	return model
}

// RightJoinSub RightJoinSub()
func (model *Model) RightJoinSub(qb interface{}, alias string, first interface{}, args ...interface{}) *Model {
	model.Builder.RightJoinSub(qb, alias, first, args...)
	return model
}

// CrossJoin CrossJoin()
func (model *Model) CrossJoin(table string) *Model {
	model.Builder.CrossJoin(table)
	return model
}

// CrossJoinSub CrossJoinSub()
func (model *Model) CrossJoinSub(qb interface{}, alias string) *Model {
	model.Builder.CrossJoinSub(qb, alias)
	return model
}

// On On()
func (model *Model) On(first interface{}, args ...interface{}) *Model {
	model.Builder.On(first, args...)
	return model
}

// OrOn OrOn()
func (model *Model) OrOn(first interface{}, args ...interface{}) *Model {
	model.Builder.OrOn(first, args...)
	return model
}

// Where Add a basic where clause to the query.
func (model *Model) Where(column interface{}, args ...interface{}) *Model {
	model.Builder.Where(column, args...)
	return model
}

// OrWhere OrWhere()
func (model *Model) OrWhere(column interface{}, args ...interface{}) *Model {
	model.Builder.OrWhere(column, args...)
	return model
}

// WhereColumn WhereColumn()
func (model *Model) WhereColumn(first interface{}, args ...interface{}) *Model {
	model.Builder.WhereColumn(first, args...)
	return model
}

// OrWhereColumn OrWhereColumn()
func (model *Model) OrWhereColumn(first interface{}, args ...interface{}) *Model {
	model.Builder.OrWhereColumn(first, args...)
	return model
}

// WhereNull WhereNull()
func (model *Model) WhereNull(column interface{}, args ...interface{}) *Model {
	model.Builder.WhereNull(column, args...)
	return model
}

// OrWhereNull OrWhereNull()
func (model *Model) OrWhereNull(column interface{}) *Model {
	model.Builder.OrWhereNull(column)
	return model
}

// WhereNotNull WhereNotNull()
func (model *Model) WhereNotNull(column interface{}, args ...interface{}) *Model {
	model.Builder.WhereNotNull(column, args...)
	return model
}

// OrWhereNotNull OrWhereNotNull()
func (model *Model) OrWhereNotNull(column interface{}) *Model {
	model.Builder.OrWhereNotNull(column)
	return model
}

// WhereRaw WhereRaw()
func (model *Model) WhereRaw(sql string, bindings ...interface{}) *Model {
	model.Builder.WhereRaw(sql, bindings...)
	return model
}

// OrWhereRaw OrWhereRaw()
func (model *Model) OrWhereRaw(sql string, bindings ...interface{}) *Model {
	model.Builder.OrWhereRaw(sql, bindings...)
	return model
}

// WhereBetween WhereBetween()
func (model *Model) WhereBetween(column interface{}, values interface{}) *Model {
	model.Builder.WhereBetween(column, values)
	return model
}

// OrWhereBetween OrWhereBetween()
func (model *Model) OrWhereBetween(column interface{}, values interface{}) *Model {
	model.Builder.OrWhereBetween(column, values)
	return model
}

// WhereNotBetween WhereNotBetween()
func (model *Model) WhereNotBetween(column interface{}, values interface{}) *Model {
	model.Builder.WhereNotBetween(column, values)
	return model
}

// OrWhereNotBetween OrWhereNotBetween()
func (model *Model) OrWhereNotBetween(column interface{}, values interface{}) *Model {
	model.Builder.OrWhereNotBetween(column, values)
	return model
}

// WhereIn WhereIn()
func (model *Model) WhereIn(column interface{}, values interface{}) *Model {
	model.Builder.WhereIn(column, values)
	return model
}

// OrWhereIn OrWhereIn()
func (model *Model) OrWhereIn(column interface{}, values interface{}) *Model {
	model.Builder.OrWhereIn(column, values)
	return model
}

// WhereNotIn WhereNotIn()
func (model *Model) WhereNotIn(column interface{}, values interface{}) *Model {
	model.Builder.WhereNotIn(column, values)
	return model
}

// OrWhereNotIn OrWhereNotIn()
func (model *Model) OrWhereNotIn(column interface{}, values interface{}) *Model {
	model.Builder.OrWhereNotIn(column, values)
	return model
}

// WhereExists WhereExists()
func (model *Model) WhereExists(closure func(qb query.Query)) *Model {
	model.Builder.WhereExists(closure)
	return model
}

// OrWhereExists OrWhereExists()
func (model *Model) OrWhereExists(closure func(qb query.Query)) *Model {
	model.Builder.OrWhereExists(closure)
	return model
}

// WhereNotExists WhereNotExists()
func (model *Model) WhereNotExists(closure func(qb query.Query)) *Model {
	model.Builder.WhereNotExists(closure)
	return model
}

// OrWhereNotExists OrWhereNotExists()
func (model *Model) OrWhereNotExists(closure func(qb query.Query)) *Model {
	model.Builder.OrWhereNotExists(closure)
	return model
}

// WhereDate WhereDate()
func (model *Model) WhereDate(column interface{}, args ...interface{}) *Model {
	model.Builder.WhereDate(column, args...)
	return model
}

// OrWhereDate OrWhereDate()
func (model *Model) OrWhereDate(column interface{}, args ...interface{}) *Model {
	model.Builder.OrWhereDate(column, args...)
	return model
}

// WhereTime WhereTime()
func (model *Model) WhereTime(column interface{}, args ...interface{}) *Model {
	model.Builder.WhereTime(column, args...)
	return model
}

// OrWhereTime OrWhereTime()
func (model *Model) OrWhereTime(column interface{}, args ...interface{}) *Model {
	model.Builder.OrWhereTime(column, args...)
	return model
}

// WhereYear WhereYear()
func (model *Model) WhereYear(column interface{}, args ...interface{}) *Model {
	model.Builder.WhereYear(column, args...)
	return model
}

// OrWhereYear OrWhereYear()
func (model *Model) OrWhereYear(column interface{}, args ...interface{}) *Model {
	model.Builder.OrWhereYear(column, args...)
	return model
}

// WhereMonth WhereMonth()
func (model *Model) WhereMonth(column interface{}, args ...interface{}) *Model {
	model.Builder.WhereMonth(column, args...)
	return model
}

// OrWhereMonth OrWhereMonth()
func (model *Model) OrWhereMonth(column interface{}, args ...interface{}) *Model {
	model.Builder.OrWhereMonth(column, args...)
	return model
}

// WhereDay WhereDay()
func (model *Model) WhereDay(column interface{}, args ...interface{}) *Model {
	model.Builder.WhereDay(column, args...)
	return model
}

// OrWhereDay OrWhereDay()
func (model *Model) OrWhereDay(column interface{}, args ...interface{}) *Model {
	model.Builder.OrWhereDay(column, args...)
	return model
}

// When When()
func (model *Model) When(value bool, callback func(qb query.Query, value bool), defaults ...func(qb query.Query, value bool)) *Model {
	model.Builder.When(value, callback, defaults...)
	return model
}

// Unless Unless()
func (model *Model) Unless(value bool, callback func(qb query.Query, value bool), defaults ...func(qb query.Query, value bool)) *Model {
	model.Builder.Unless(value, callback, defaults...)
	return model
}

// GroupBy GroupBy()
func (model *Model) GroupBy(groups ...interface{}) *Model {
	model.Builder.GroupBy(groups...)
	return model
}

// GroupByRaw GroupByRaw()
func (model *Model) GroupByRaw(expression string, bindings ...interface{}) *Model {
	model.Builder.GroupByRaw(expression, bindings...)
	return model
}

// Having Having()
func (model *Model) Having(column interface{}, args ...interface{}) *Model {
	model.Builder.Having(column, args...)
	return model
}

// OrHaving OrHaving()
func (model *Model) OrHaving(column interface{}, args ...interface{}) *Model {
	model.Builder.OrHaving(column, args...)
	return model
}

// HavingBetween HavingBetween()
func (model *Model) HavingBetween(column interface{}, values interface{}, args ...interface{}) *Model {
	model.Builder.HavingBetween(column, values, args...)
	return model
}

// OrHavingBetween OrHavingBetween()
func (model *Model) OrHavingBetween(column interface{}, values interface{}, args ...interface{}) *Model {
	model.Builder.OrHavingBetween(column, values, args...)
	return model
}

// HavingRaw HavingRaw()
func (model *Model) HavingRaw(sql string, bindings ...interface{}) *Model {
	model.Builder.HavingRaw(sql, bindings...)
	return model
}

// OrHavingRaw OrHavingRaw()
func (model *Model) OrHavingRaw(sql string, bindings ...interface{}) *Model {
	model.Builder.OrHavingRaw(sql, bindings...)
	return model
}

// OrderBy OrderBy()
func (model *Model) OrderBy(column interface{}, args ...string) *Model {
	model.Builder.OrderBy(column, args...)
	return model
}

// OrderByDesc OrderByDesc()
func (model *Model) OrderByDesc(column interface{}) *Model {
	model.Builder.OrderByDesc(column)
	return model
}

// OrderByRaw OrderByRaw()
func (model *Model) OrderByRaw(sql string, bindings ...interface{}) *Model {
	model.Builder.OrderByRaw(sql, bindings...)
	return model
}

// Skip Skip()
func (model *Model) Skip(value int) *Model {
	model.Builder.Skip(value)
	return model
}

// Offset Offset()
func (model *Model) Offset(value int) *Model {
	model.Builder.Offset(value)
	return model
}

// Take Take()
func (model *Model) Take(value int) *Model {
	model.Builder.Take(value)
	return model
}

// Limit Limit()
func (model *Model) Limit(value int) *Model {
	model.Builder.Limit(value)
	return model
}

// SharedLock SharedLock()
func (model *Model) SharedLock() *Model {
	model.Builder.SharedLock()
	return model
}

// LockForUpdate LockForUpdate()
func (model *Model) LockForUpdate() *Model {
	model.Builder.LockForUpdate()
	return model
}
