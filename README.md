# Xun Database

[![Build Status](https://travis-ci.com/YaoApp/xun.svg?branch=main)](https://travis-ci.com/YaoApp/xun)
[![codecov](https://codecov.io/gh/YaoApp/xun/branch/main/graph/badge.svg?token=R4FW9PXF01)](https://codecov.io/gh/YaoApp/xun)
[![Go Report Card](https://goreportcard.com/badge/github.com/YaoApp/xun)](https://goreportcard.com/report/github.com/YaoApp/xun)
[![Go Reference](https://pkg.go.dev/badge/github.com/yaoapp/xun.svg)](https://pkg.go.dev/github.com/yaoapp/xun)

Xun Database is an object-relational mapper (ORM), that is written in golang and supports JSON schema. Xun providing `query builder`, `schema builder` and `model builder`, can change the table structure at run time, especially suitable for use in Low-Code application.

The name Xun comes from the Chinese word 巽(xùn). It is one of the eight trigrams, a symbol of wind. it also symbolizes the object filled in everywhere.

## Installation

To install Xun package, you need to install Go and set your Go workspace before.

1. The first need Go installed (version 1.12+ is required), then you can use the below Go command to install Xun.

```bash
$ go get -u github.com/yaoapp/xun
```

2. Import xun in your code:

```golang
import "github.com/yaoapp/xun/capsule"
```

3. Import the grammar driver that your project used.
   Xun package providing `MySQL`, `PostgreSQL` and `SQLite` grammar drivers, you can also using the third-party grammar driver or written by yourself. See [how to write Xun grammar driver](docs/contributing/xun-grammar-driver.md)

`PostgreSQL`:

```golang
import (
    "github.com/yaoapp/xun/capsule"
    _ "github.com/yaoapp/xun/grammar/postgres" // PostgreSQL
)
```

`MySQL` or `MariaDB`:

```golang
import (
    "github.com/yaoapp/xun/capsule"
    _ "github.com/yaoapp/xun/grammar/mysql"    // MySQL
)
```

`SQLite`:

```golang
import (
    "github.com/yaoapp/xun/capsule"
    _ "github.com/yaoapp/xun/grammar/sqlite3"  // SQLite
)
```

if your project used several types database, can also import them together.

```golang
import (
    "github.com/yaoapp/xun/capsule"
    _ "github.com/yaoapp/xun/grammar/postgres"  // PostgreSQL
    _ "github.com/yaoapp/xun/grammar/sqlite3"   // SQLite
    _ "github.com/third/party/clickhouse"       // third-party or yourself driver
)
```

## Quick start

First, create a new "db" manager instance. capsule aims to make configuring the library for usage outside of the Yao framework as easy as possible.

```golang
import (
    "github.com/yaoapp/xun/capsule"
    _ "github.com/yaoapp/xun/grammar/postgres"  // PostgreSQL
)

func main(){

    // Connect to PostgreSQL
    db := capsule.New().AddConn("primary", "postgres",
            "postgres://postgres:123456@127.0.0.1/xun?sslmode=disable&search_path=xun",
        )

    // Get the schema interface
    schema := db.Schema()

    // Get the query interface
    query := db.Query()

    // Get the model interface
    model := db.Model()

}
```

`Connect to MySQL`

```golang
  db := capsule.New().AddConn("primary", "mysql",
            "root:123456@tcp(192.168.31.119:3307)/xun?charset=utf8mb4&parseTime=True&loc=Local",
        )
```

`Connect to SQLite`

```golang
  db := capsule.New().AddConn("primary", "sqlite3", "file:///data/xun.db")
```

```golang
  db := capsule.New().AddConn("primary", "sqlite3", ":memory:")
```

`Multiple connections`

```golang
  db := capsule.New().
        AddConn("primary", "mysql",
            "root:123456@tcp(192.168.31.119:3307)/xun?charset=utf8mb4&parseTime=True&loc=Local",
        ).
        AddReadConn("secondary",  "mysql",
            "readonly:123456@tcp(192.168.31.119:3306)/xun?charset=utf8mb4&parseTime=True&loc=Local",
        ) // Add a readonly connection
```

Read more [Xun Capsule References](docs/capsule.md)

### Using The Schema Interface

```golang
import (
    "github.com/yaoapp/xun/capsule"
    "github.com/yaoapp/xun/dbal/schema"
    _ "github.com/yaoapp/xun/grammar/postgres"  // PostgreSQL
)

func main(){

    db := capsule.New().AddConn("primary",
            "postgres",
            "postgres://postgres:123456@127.0.0.1/xun?sslmode=disable&search_path=xun",
        )

    // Get the schema interface
    schema := db.Schema()

    // Create table
    builder.MustCreateTable("user", func(table schema.Blueprint) {
        table.ID("id")
        table.String("name", 80).Index()
        table.String("nickname", 128).Unique()
        table.String("bio")
        table.TinyInteger("gender").Index()
        table.DateTime("birthday").Index()
        talbe.IpAddress("login_ip").Index()
        table.AddIndex("birthday_gender", "birthday", "gender")
        table.SoftDeletes() // Add deleted_at field
        table.Timestamps()  // Add created_at and updated_at fields
    })

    // Alter table
    builder.MustAlterTable("user", func(table schema.Blueprint) {
        table.Float("BMI", 3, 1).Index() // 20.3 Float(name string, total int, places int)
        table.Float("weight", 5, 2).Index()  // 103.17
        table.SmallInteger("height").Index()
        table.Float("distance").Index()  // 16981.62 The total default value is 10 , the places default value is 2.
        table.Year("birthday")
    })
}
```

Read more [Xun Schema References](docs/schema.md)

### Using The Query Interface

```golang
// comming soon
```

Read more [Xun Query References](docs/query.md)

### Using The Model Interface

```golang
// comming soon
```

Read more [Xun Model References](docs/model.md)
