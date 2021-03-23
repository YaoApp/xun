package query

// Query The database Query interface
type Query interface {
	Where()
	Join()

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

}
