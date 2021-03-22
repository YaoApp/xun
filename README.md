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

3. Import the database driver that your project used.
   Xun package providing `MySQL`, `Postgres` and `SQLite` grammar drivers, you can also using the third-party grammar driver or written by yourself. See [how to write Xun grammar driver](xun-grammar-driver.md)

`PostgreSQL`:

```golang
import (
    "github.com/yaoapp/xun/capsule"
    _ "github.com/yaoapp/xun/grammar/postgres" // Postgres
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

if your project used several types database server, can also import them together:

```golang
import (
    "github.com/yaoapp/xun/capsule"
    _ "github.com/yaoapp/xun/grammar/postgres"  // Postgres
    _ "github.com/yaoapp/xun/grammar/sqlite3"   // SQLite
    _ "github.com/third/party/clickhouse"       // third-party or yourself driver
)
```

## Quick start

First, create a new "db" manager instance. capsule aims to make configuring the library for usage outside of the Yao framework as easy as possible.

```golang
import (
    "github.com/yaoapp/xun/capsule"
    _ "github.com/yaoapp/xun/grammar/postgres"  // Postgres
)

func main(){

    db := capsule.New().AddConn("primary",
            "postgres",
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

### Using The Schema Interface

[Xun Schema References](docs/schema.md)

```golang

db := capsule.New().AddConn("primary",
        "postgres",
        "postgres://postgres:123456@127.0.0.1/xun?sslmode=disable&search_path=xun",
    )

// Get the schema interface
schema := db.Schema()

// Create table
builder.MustCreateTable("user", func(table Blueprint) {
	table.ID("id")
	table.String("name", 80).Index()
	table.String("nickname", 128).Unique()
    table.TinyInteger("gender").Index()
    table.DateTime("birthday").Index()
    talbe.IpAddress("login_ip").Index()
    table.AddIndex("birthday_gender", "birthday", "gender")
})

// Alter table
builder.MustAlterTable("user", func(table Blueprint) {
	table.SmallInteger("tall").Index()
    table.Year("birthday")
})

```

read more [Xun Schema References](docs/schema.md)

### Using The Query Interface

```golang
// comming soon
```

read more [Xun Query References](docs/query.md)

### Using The Model Interface

```golang
// comming soon
```

read more [Xun Model References](docs/model.md)
