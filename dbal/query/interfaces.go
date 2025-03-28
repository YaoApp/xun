package query

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun"
)

// Query The database Query interface
type Query interface {

	// defined in the builder.go file
	New() Query
	Clone() Query
	Reset() Query
	Builder() *Builder
	Driver() (string, error)

	// defined in the query.go file
	Table(name string) Query
	Get(v ...interface{}) ([]xun.R, error)
	MustGet(v ...interface{}) []xun.R
	First(v ...interface{}) (xun.R, error)
	MustFirst(v ...interface{}) xun.R
	Find(id interface{}, args ...interface{}) (xun.R, error)
	MustFind(id interface{}, args ...interface{}) xun.R
	Value(column string, v ...interface{}) (interface{}, error)
	MustValue(column string, v ...interface{}) interface{}
	Exists() (bool, error)
	MustExists() bool
	DoesntExist() (bool, error)
	MustDoesntExist() bool
	ToSQL() string
	GetBindings() []interface{}

	// defined in the paginate.go file
	Paginate(perpage int, page int, v ...interface{}) (xun.P, error)
	MustPaginate(perpage int, page int, v ...interface{}) xun.P
	Chunk(size int, callback func(items []interface{}, page int) error, v ...interface{}) error
	MustChunk(size int, callback func(items []interface{}, page int) error, v ...interface{})

	// defined in the connection.go file
	DB(usewrite ...bool) *sqlx.DB
	IsRead() bool
	UseRead() Query
	UseWrite() Query
	IsWrite() bool

	// defined in the aggregate.go file
	Count(columns ...interface{}) (int64, error)
	MustCount(columns ...interface{}) int64
	Min(columns ...interface{}) (xun.N, error)
	MustMin(columns ...interface{}) xun.N
	Max(columns ...interface{}) (xun.N, error)
	MustMax(columns ...interface{}) xun.N
	Sum(columns ...interface{}) (xun.N, error)
	MustSum(columns ...interface{}) xun.N
	Avg(columns ...interface{}) (xun.N, error)
	MustAvg(columns ...interface{}) xun.N

	// defined in the select.go file
	SQL(stmt string, bindings ...interface{}) Query
	Select(columns ...interface{}) Query
	SelectAppend(columns ...interface{}) Query
	SelectRaw(expression string, bindings ...interface{}) Query
	SelectSub(qb interface{}, alias string) Query
	Distinct(args ...interface{}) Query

	// defined in the from.go file
	From(name string) Query
	FromRaw(sql string, bindings ...interface{}) Query
	FromSub(qb interface{}, alias string) Query

	// defined in the union.go file
	Union(query interface{}, all ...bool) Query
	UnionAll(query interface{}) Query

	// defined in the join.go file
	Join(table string, first interface{}, args ...interface{}) Query
	JoinRaw(sql string, bindings ...interface{}) Query
	JoinSub(qb interface{}, alias string, first interface{}, args ...interface{}) Query
	LeftJoin(table string, first interface{}, args ...interface{}) Query
	LeftJoinSub(qb interface{}, alias string, first interface{}, args ...interface{}) Query
	RightJoin(table string, first interface{}, args ...interface{}) Query
	RightJoinSub(qb interface{}, alias string, first interface{}, args ...interface{}) Query
	CrossJoin(table string) Query
	CrossJoinSub(qb interface{}, alias string) Query
	On(first interface{}, args ...interface{}) Query
	OrOn(first interface{}, args ...interface{}) Query

	// defined in the where.go file
	Where(column interface{}, args ...interface{}) Query
	OrWhere(column interface{}, args ...interface{}) Query
	WhereColumn(first interface{}, args ...interface{}) Query
	OrWhereColumn(first interface{}, args ...interface{}) Query
	WhereNull(column interface{}, args ...interface{}) Query
	OrWhereNull(column interface{}) Query
	WhereNotNull(column interface{}, args ...interface{}) Query
	OrWhereNotNull(column interface{}) Query
	WhereRaw(sql string, bindings ...interface{}) Query
	OrWhereRaw(sql string, bindings ...interface{}) Query
	WhereBetween(column interface{}, values interface{}) Query
	OrWhereBetween(column interface{}, values interface{}) Query
	WhereNotBetween(column interface{}, values interface{}) Query
	OrWhereNotBetween(column interface{}, values interface{}) Query
	WhereIn(column interface{}, values interface{}) Query
	OrWhereIn(column interface{}, values interface{}) Query
	WhereNotIn(column interface{}, values interface{}) Query
	OrWhereNotIn(column interface{}, values interface{}) Query
	WhereExists(closure func(qb Query)) Query
	OrWhereExists(closure func(qb Query)) Query
	WhereNotExists(closure func(qb Query)) Query
	OrWhereNotExists(closure func(qb Query)) Query
	WhereDate(column interface{}, args ...interface{}) Query
	OrWhereDate(column interface{}, args ...interface{}) Query
	WhereTime(column interface{}, args ...interface{}) Query
	OrWhereTime(column interface{}, args ...interface{}) Query
	WhereYear(column interface{}, args ...interface{}) Query
	OrWhereYear(column interface{}, args ...interface{}) Query
	WhereMonth(column interface{}, args ...interface{}) Query
	OrWhereMonth(column interface{}, args ...interface{}) Query
	WhereDay(column interface{}, args ...interface{}) Query
	OrWhereDay(column interface{}, args ...interface{}) Query
	When(value bool, callback func(qb Query, value bool), defaults ...func(qb Query, value bool)) Query
	Unless(value bool, callback func(qb Query, value bool), defaults ...func(qb Query, value bool)) Query

	// defined in the group.go file
	GroupBy(groups ...interface{}) Query
	GroupByRaw(expression string, bindings ...interface{}) Query
	Having(column interface{}, args ...interface{}) Query
	OrHaving(column interface{}, args ...interface{}) Query
	HavingBetween(column interface{}, values interface{}, args ...interface{}) Query
	OrHavingBetween(column interface{}, values interface{}, args ...interface{}) Query
	HavingRaw(sql string, bindings ...interface{}) Query
	OrHavingRaw(sql string, bindings ...interface{}) Query

	// defined in the order.go file
	OrderBy(column interface{}, args ...string) Query
	OrderByDesc(column interface{}) Query
	OrderByRaw(sql string, bindings ...interface{}) Query

	// defined in the limit.go file
	Skip(value int) Query
	Offset(value int) Query
	Take(value int) Query
	Limit(value int) Query

	// defined in the lock.go file
	SharedLock() Query
	LockForUpdate() Query

	// defined in the insert.go file
	Insert(v interface{}, columns ...interface{}) error
	MustInsert(v interface{}, columns ...interface{})
	InsertOrIgnore(v interface{}, columns ...interface{}) (int64, error)
	MustInsertOrIgnore(v interface{}, columns ...interface{}) int64
	InsertGetID(v interface{}, args ...interface{}) (int64, error)
	MustInsertGetID(v interface{}, args ...interface{}) int64
	InsertUsing(qb interface{}, columns ...interface{}) (int64, error)
	MustInsertUsing(qb interface{}, columns ...interface{}) int64

	// defined in the update.go file
	Upsert(values interface{}, uniqueBy interface{}, update interface{}, columns ...interface{}) (int64, error)
	MustUpsert(values interface{}, uniqueBy interface{}, update interface{}, columns ...interface{}) int64
	UpdateOrInsert(attributes interface{}, values ...interface{}) (bool, error)
	MustUpdateOrInsert(attributes interface{}, values ...interface{}) bool
	Update(v interface{}) (int64, error)
	MustUpdate(v interface{}) int64
	Increment(column interface{}, amount interface{}, extra ...interface{}) (int64, error)
	MustIncrement(column interface{}, amount interface{}, extra ...interface{}) int64
	Decrement(column interface{}, amount interface{}, extra ...interface{}) (int64, error)
	MustDecrement(column interface{}, amount interface{}, extra ...interface{}) int64

	// defined in the delete.go file
	Delete() (int64, error)
	MustDelete() int64
	Truncate() error
	MustTruncate()

	// defined in the exec.go file
	Exec(sql string, bindings ...interface{}) (sql.Result, error)
	ExecWrite(sql string, bindings ...interface{}) (sql.Result, error)

	// defined in the debug.go file
	DD()
	Dump()
}

// @todo
// Chunking Results:
// table(`users`).where("weight", ">", 99.00).chunk(100, func( users ){ ... } )
// table(`users`).where("weight", ">", 99.00).chunkById(100, func( users ){ update... } )

// Aggregates:
// table(`users`).where("weight", ">", 99.00).count()
// table(`users`).where("weight", ">", 99.00).max("price")
// table(`users`).where("weight", ">", 99.00).avg("price")

// Determining If Records Exist:
// table(`users`).where("weight", ">", 99.00).exists()

// Select Statements:
// table(`users`).select(`name`, `nickname as user_nickname`)
// table(`users`).distinct()
// table(`users`).addSelect("height")

// Raw Expressions:
// table(`users`).select(dbal.raw(`count(*) as user_count, status`))

// Raw Methods:
// table(`orders`).selectRaw(`weight * ? as price_with_tax`, [1.0825])
// table(`orders`).whereRaw(`price > IF(state = "TX", ?, 100)``, [200])
// table(`orders`).where("price" , ">", 0).orWhereRaw(`price > IF(state = "TX", ?, 100)`, [200])
// table(`orders`).
// 		select(`department`, dbal.raw(`SUM(price) as total_sales`)).
// 		groupBy(`department`).
// 		havingRaw(`SUM(price) > ?`, [2500]).
//		get()
// table(`orders`).
// 		select(`department`, dbal.raw(`SUM(price) as total_sales`)).
// 		groupBy(`department`).
// 		orHavingRaw(`SUM(price) > ?`, [2500]).
//		get()
// table(`orders`).orderByRaw(`updated_at - created_at DESC`)
// table(`orders`).groupByRaw(`city, state`)

// Joins:
// table(`users`).
// 		join(`contacts`, `users.id`, `=`, `contacts.user_id`).
// 		join(`orders`, `users.id`, `=`, `orders.user_id`).
// 		select(`users.*`, `contacts.phone`, `orders.price`)
// table(`users`).leftJoin(`posts`, `users.id`, `=`, `posts.user_id`)
// table(`users`).rightJoin(`posts`, `users.id`, `=`, `posts.user_id`)
// table(`sizes`).crossJoin(`colors`)
// table(`sizes`).join(`contacts`, func(join){ join.on(`users.id`, `=`, `contacts.user_id`).orOn(...)})
// table(`sizes`).join(`contacts`, func(join){ join.on(`users.id`, `=`, `contacts.user_id`).where(`contacts.user_id`, `>`, 5)})
// Subquery Joins:
// 		latestPostsQB := table.("posts").select(`user_id`, DB::raw(`MAX(created_at) as last_post_created_at`)).where(`is_published`, true).groupBy(`user_id`)
// 		table(`users`).joinSub($latestPostsQB, `latest_posts`, function (join){ join.on(`users.id`, `=`, `latest_posts.user_id`)})

// Unions:
// 		firstQB := table(`users`).whereNull(`first_name`)
// 		table(`users`).whereNull(`last_name`).union($firstQB)

// Basic Where Clauses
// Where Clauses:
// 		table("users").where(`votes`, `=`, 100)
// 		table("users").where(`weight`, `>`, 99.00)
// 		table("users").where(`votes`, `>=`, 100)
// 		table("users").where(`votes`, `<>`, 100)
// 		table("users").where(`name`, `like`, `T%`)
// 		table("users").where(`votes`, 100)
// 		table("users").where([]query.W{
//			query.W{"status", "=", "1"},
//			query.W{"subscribed", "<>", "1"},
// 		})
// Or Where Clauses:
//		table("users").where(`votes`, `>=`, 100).orWhere(`name`, `John`)
//		table("users").where(`votes`, `>=`, 100).orWhere( func( qb query.Query ){
//			qb.where(`name`, `Abigail`).
// 			  .where(`votes`, `>`, 50)
// 		})
// JSON Where Clauses:
//		table("users").where(`preferences->dining->meal`, `salad`)
// 		table("users").whereJsonContains(`options->languages`, `en`)
// 		table("users").whereJsonContains(`options->languages`, [`en`, `de`])
// 		table("users").whereJsonLength(`options->languages`, 0)
// 		table("users").whereJsonLength(`options->languages`, `>`, 1)
// Additional Where Clauses:
// 		whereBetween / orWhereBetween
// 		whereNotBetween / orWhereNotBetween
// 		whereIn / whereNotIn / orWhereIn / orWhereNotIn
// 		whereNull / whereNotNull / orWhereNull / orWhereNotNull
// 		whereDate / whereMonth / whereDay / whereYear / whereTime
// 		whereColumn / orWhereColumn
// Logical Grouping:
//		table("users").
// 			.where(`name`, `=`, `John`)
// 			.where(func( qb query.Query ) {
//				qb.where(`votes`, `>`, 100).
// 				   orWhere(`title`, `=`, `Admin`)
// 			})

// Advanced Where Clauses
// Where Exists Clauses:
// 		table("users").whereExists(func( qb query.Query ) {
//			qb.select(dbal.raw(1)).
// 				from(`orders`).
// 				whereColumn(`orders.user_id`, `users.id`)
// 		})
// Subquery Where Clauses:
//		table("users").where(func( qb query.Query ) {
//			qb.select(`type`).from(`membership`).
// 			  whereColumn(`membership.user_id`, `users.id`).
//			  orderByDesc(`membership.start_date`).
// 			  limit(1)
//		}, "pro")
//
//		table("users").where(func( qb query.Query ) {
//			qb.selectRaw(`avg(i.amount)`).from(`incomes as i`);
// 		})
//

// Ordering, Grouping, Limit & Offset
// Ordering:
//		table(`users`).orderBy(`name`, `desc`)
//		table(`users`).orderBy(`email`, `asc`)
//		table(`users`).latest().first()
// 		table(`users`).oldest().first()
//		table(`users).InRandomOrder().first()
//
// 		qb := table(`users`).orderBy(`name`, `desc`)
//		rows := qb.reorder().get()   // Removing Existing Orderings
//
//		qb := table(`users`).orderBy(`name`, `desc`)
//  	rows := qb.reorder(`email`, `desc`).get() //remove all existing "order by" clauses and apply an entirely new order
// Grouping:
//		table(`users`).
// 			groupBy(`account_id`).
// 			having(`account_id`, `>`, 100)
//
//		table(`users`).
// 			groupBy(`first_name`, `status`).
// 			having(`account_id`, `>`, 100)
// Limit & Offset:
//		table(`users`).skip(10).take(5)
//		table(`users`).offset(10).limit(5)

// Conditional Clauses
// role := "admin"
// table(`users`).when( role != "", func( qb, role){
// 		return $qb.where(`role_id`, role)
// }, nil)
//
// sortBy := "votes"
// table(`users`).when( sort == "votes", func( qb, sortBy){
//		return qb.orderBy("votes")
// }, func( qb, sortBy){
//		return qb.orderBy("name")
// })

// Insert Statements
// table(`users`).insert([ `email` : `kayla@example.com`,`votes` : 0])
// table(`users`).insert(
// 		[`email` : `picard@example.com`, `votes` : 0],
// 		[`email` : `janeway@example.com`, `votes` : 0],
// )
// table(`users`).insertOrIgnore(
// 		[`id` : 1, `email` : `sisko@example.com`],
// 		[`id` : 2, `email` : `archer@example.com`],
// )
// id, err := table(`users`).insertGetId(
// 		[`email` : `john@example.com`, `votes` : 0]
// )

// Update Statements
// table(`users`).where("id", 1).update([`votes` : 1])
// table(`users`).where("id", 1).update([`options->enabled` : true])
// table(`users`).increment(`votes`)
// table(`users`).increment(`votes`, 5)
// table(`users`).increment(`votes`, 1, ["name":"John"])
// table(`users`).decrement(`votes`)
// table(`users`).decrement(`votes`, 5)
// table(`users`).decrement(`votes`, 1, ["name":"John"])

// Update Or Insert Statements
// table(`flights`).upsert([
//     [`departure` : `Oakland`, `destination` : `San Diego`, `price` : 99],
//     [`departure` : `Chicago`, `destination` : `New York`, `price` : 150],
// ], [`departure`, `destination`], [`price`])
// table(`flights`).updateOrInsert(
// 		[`email` : `john@example.com`, `name` : `John`],
// 		[`votes` : `2`],
// )

// Delete Statements
// table(`users`).where("id", 1).delete()
// table(`users`).delete()
// table(`users`).truncate() // When truncating a PostgreSQL database, the CASCADE behavior will be applied. This means that all foreign key related records in other tables will be deleted as well.

// Pessimistic Locking
// table(`user`).
// 		where('votes', '>', 100).
//		sharedLock().
// 		get() // LOCK IN SHARE MODE
// table(`user`).
// 		where('votes', '>', 100).
//		lockForUpdate().  //  FOR UPDATE
// 		get()

// Debugging
// table(`user`).where('votes', '>', 100).DD()
// table(`user`).where('votes', '>', 100).Dump()

// Paginate
// table(`user`).where('votes', '>', 100).Paginate(15)
