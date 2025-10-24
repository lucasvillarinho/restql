# `↔️ restQL`

From REST filters to SQL queries

[![Go Version](https://img.shields.io/github/go-mod/go-version/lucasvillarinho/restql)](https://github.com/lucasvillarinho/restql)
[![Go Report Card](https://goreportcard.com/badge/github.com/lucasvillarinho/restql)](https://goreportcard.com/report/github.com/lucasvillarinho/restql)
[![codecov](https://codecov.io/gh/lucasvillarinho/restql/branch/master/graph/badge.svg)](https://codecov.io/gh/lucasvillarinho/restql)

RestQL is a Go library that converts REST query parameters into SQL queries with optional validation and security features. It provides a type-safe way to expose database filtering, sorting, and pagination through HTTP APIs.

## Features

- **Filter Expressions**: Parse complex filter expressions with operators (`=`, `!=`, `>`, `<`, `>=`, `<=`, `LIKE`, `IN`, etc.)
- **Optional Validation**: Field whitelisting and limit/offset validation when you need it
- **Fluent API**: Clean, chainable interface for building queries
- **SQL Injection Protection**: Parameterized queries with proper escaping
- **Minimal Dependencies**: Only requires participle/v2 for parsing

## Installation

```bash
go get github.com/lucasvillarinho/restql
```

## Quick Start

### Basic Usage

```go

package main

import (
    "log"
    "net/url"

    "github.com/lucasvillarinho/restql"
)

func main() {
    // Create a RestQL instance
    rql := restql.NewRestQL()

    params, _ := url.ParseQuery("filter=age>18&limit=50")

   
    query, err := rql.Parse(params, "users",
        restql.WithAllowedFields([]string{"id", "name", "email", "age"}),
        restql.WithMaxLimit(100),
        restql.WithMaxOffset(1000),
    )
    if err != nil {
        log.Fatal(err)
    }

    sql, args, err := query.ToSQL()
    if err != nil {
        log.Fatal(err)
    }

    // sql: SELECT * FROM users WHERE age > ? LIMIT 50
    // args: [18]
}

```

## Query Parameters

RestQL supports these URL query parameters:

- `filter` - Filter expression (e.g., `age>18 && status='active'`)
- `fields` - Comma-separated fields to select (e.g., `id,name,email`)
- `sort` - Comma-separated sort fields, prefix with `-` for DESC (e.g., `-created_at,name`)
- `limit` - Maximum number of results
- `offset` - Number of results to skip

## Supported Operators

### Comparison Operators

<details>
<summary><b>Equal (=)</b></summary>

```go
params, _ := url.ParseQuery("filter=status='active'")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE status = ?
// args: ["active"]
```

</details>

<details>
<summary><b>Not Equal (!=, <>)</b></summary>

```go
params, _ := url.ParseQuery("filter=status!='inactive'")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE status != ?
// args: ["inactive"]
```

</details>

<details>
<summary><b>Greater Than (>)</b></summary>

```go
params, _ := url.ParseQuery("filter=age>18")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE age > ?
// args: [18]
```

</details>

<details>
<summary><b>Less Than (<)</b></summary>

```go
params, _ := url.ParseQuery("filter=price<100")
sql, args, _ := restql.Parse(params, "products").ToSQL()
// SELECT * FROM products WHERE price < ?
// args: [100]
```

</details>

<details>
<summary><b>Greater Than or Equal (>=)</b></summary>

```go
params, _ := url.ParseQuery("filter=rating>=4.5")
sql, args, _ := restql.Parse(params, "reviews").ToSQL()
// SELECT * FROM reviews WHERE rating >= ?
// args: [4.5]
```

</details>

<details>
<summary><b>Less Than or Equal (<=)</b></summary>

```go
params, _ := url.ParseQuery("filter=stock<=10")
sql, args, _ := restql.Parse(params, "products").ToSQL()
// SELECT * FROM products WHERE stock <= ?
// args: [10]
```

</details>

### Pattern Matching

<details>
<summary><b>LIKE (case-sensitive)</b></summary>

```go
params, _ := url.ParseQuery("filter=name LIKE '%John%'")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE name LIKE ?
// args: ["%John%"]
```

</details>

<details>
<summary><b>NOT LIKE</b></summary>

```go
params, _ := url.ParseQuery("filter=name NOT LIKE '%test%'")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE name NOT LIKE ?
// args: ["%test%"]
```

</details>

### List Operations

<details>
<summary><b>IN</b></summary>

```go
params, _ := url.ParseQuery("filter=status IN ('active','pending','approved')")
sql, args, _ := restql.Parse(params, "orders").ToSQL()
// SELECT * FROM orders WHERE status IN (?, ?, ?)
// args: ["active", "pending", "approved"]
```

</details>

<details>
<summary><b>NOT IN</b></summary>

```go
params, _ := url.ParseQuery("filter=role NOT IN ('admin','superadmin')")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE role NOT IN (?, ?)
// args: ["admin", "superadmin"]
```

</details>

### Null Checks

<details>
<summary><b>IS NULL</b></summary>

```go
params, _ := url.ParseQuery("filter=deleted_at IS NULL")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE deleted_at IS NULL
// args: []
```

</details>

<details>
<summary><b>IS NOT NULL</b></summary>

```go
params, _ := url.ParseQuery("filter=email IS NOT NULL")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE email IS NOT NULL
// args: []
```

</details>

### Logical Operators

<details>
<summary><b>AND (&&)</b></summary>

```go
params, _ := url.ParseQuery("filter=age>=18 && status='active'")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE (age >= ? AND status = ?)
// args: [18, "active"]
```

</details>

<details>
<summary><b>OR (||)</b></summary>

```go
params, _ := url.ParseQuery("filter=role='admin' || role='moderator'")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE (role = ? OR role = ?)
// args: ["admin", "moderator"]
```

</details>

<details>
<summary><b>Grouping ()</b></summary>

```go
params, _ := url.ParseQuery("filter=(age>=18 && country='US') || (age>=21 && country='UK')")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE ((age >= ? AND country = ?) OR (age >= ? AND country = ?))
// args: [18, "US", 21, "UK"]
```

</details>

## Security

### Field Whitelisting

Always validate fields in production to prevent unauthorized data access:

```go
query.Validate(
    restql.WithAllowedFields([]string{"id", "name", "email"}),
).ToSQL()

// Attempting to filter/select/sort on 'password' will fail
```

### Limit Protection

Prevent excessive data retrieval:

```go
query.Validate(
    restql.WithMaxLimit(100),
    restql.WithMaxOffset(10000),
).ToSQL()

// User can't request limit=999999
```

## Integration Examples

### ORMs

<details>
<summary><b>database/sql</b></summary>

```go
package main

import (
    "database/sql"
    "log"
    "net/url"

    _ "github.com/lib/pq"
    "github.com/lucasvillarinho/restql"
)

func main() {
    db, _ := sql.Open("postgres", "connection-string")
    defer db.Close()

    params, _ := url.ParseQuery("filter=age>=18&sort=-created_at&limit=10")

    sql, args, err := restql.Parse(params, "users").ToSQL()
    if err != nil {
        log.Fatal(err)
    }

    rows, err := db.Query(sql, args...)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    // Process rows...
}
```

</details>

<details>
<summary><b>GORM</b></summary>

```go
package main

import (
    "log"
    "net/url"

    "github.com/lucasvillarinho/restql"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type User struct {
    ID        uint
    Name      string
    Email     string
    Age       int
    Status    string
    CreatedAt time.Time
}

func main() {
    db, _ := gorm.Open(postgres.Open("connection-string"), &gorm.Config{})

    params, _ := url.ParseQuery("filter=(age>=18 && status='active')&fields=id,name,email&limit=50")

    allowedFields := []string{"id", "name", "email", "age", "status", "created_at"}

    sql, args, err := restql.Parse(params, "users").
        Validate(restql.WithAllowedFields(allowedFields)).
        ToSQL()
    if err != nil {
        log.Fatal(err)
    }

    var users []User
    if err := db.Raw(sql, args...).Scan(&users).Error; err != nil {
        log.Fatal(err)
    }

    // Use users...
}
```

</details>

<details>
<summary><b>sqlx</b></summary>

```go
package main

import (
    "log"
    "net/url"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "github.com/lucasvillarinho/restql"
)

type User struct {
    ID        int       `db:"id"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    Age       int       `db:"age"`
    Status    string    `db:"status"`
    CreatedAt time.Time `db:"created_at"`
}

func main() {
    db, _ := sqlx.Connect("postgres", "connection-string")
    defer db.Close()

    params, _ := url.ParseQuery("filter=status='active'&sort=-created_at&limit=20")

    sql, args, err := restql.Parse(params, "users").ToSQL()
    if err != nil {
        log.Fatal(err)
    }

    var users []User
    if err := db.Select(&users, sql, args...); err != nil {
        log.Fatal(err)
    }

    // Use users...
}
```

</details>

### HTTP Frameworks

<details>
<summary><b>Echo Framework (with NewRestQL for reusable config)</b></summary>

```go
package main

import (
    "database/sql"
    "net/http"

    "github.com/labstack/echo/v4"
    _ "github.com/lib/pq"
    "github.com/lucasvillarinho/restql"
)

type User struct {
    ID     int    `json:"id"`
    Name   string `json:"name"`
    Email  string `json:"email"`
    Age    int    `json:"age"`
    Status string `json:"status"`
}

type Product struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Price float64 `json:"price"`
}

func main() {
    e := echo.New()
    db, _ := sql.Open("postgres", "connection-string")
    defer db.Close()

    // Create reusable RestQL instance
    rql := restql.NewRestQL()

    // Users endpoint - with validation
    e.GET("/users", func(c echo.Context) error {
        sql, args, err := rql.Parse(c.QueryParams(), "users",
            restql.WithAllowedFields([]string{"id", "name", "email", "age", "status", "created_at"}),
            restql.WithMaxLimit(100),
            restql.WithMaxOffset(1000),
        ).ToSQL()

        if err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
        }

        rows, err := db.Query(sql, args...)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
        }
        defer rows.Close()

        var users []User
        // Scan rows into users...
        return c.JSON(http.StatusOK, users)
    })

    // Products endpoint - different validation rules
    e.GET("/products", func(c echo.Context) error {
        sql, args, err := rql.Parse(c.QueryParams(), "products",
            restql.WithAllowedFields([]string{"id", "name", "price", "category"}),
            restql.WithMaxLimit(50),  // Different limit for products
        ).ToSQL()

        if err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
        }

        rows, err := db.Query(sql, args...)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
        }
        defer rows.Close()

        var products []Product
        // Scan rows into products...
        return c.JSON(http.StatusOK, products)
    })

    e.Logger.Fatal(e.Start(":8080"))
}
```

</details>

<details>
<summary><b>Fiber</b></summary>

```go
package main

import (
    "database/sql"
    "net/url"

    "github.com/gofiber/fiber/v2"
    _ "github.com/lib/pq"
    "github.com/lucasvillarinho/restql"
)

type User struct {
    ID     int    `json:"id"`
    Name   string `json:"name"`
    Email  string `json:"email"`
    Age    int    `json:"age"`
    Status string `json:"status"`
}

func main() {
    app := fiber.New()
    db, _ := sql.Open("postgres", "connection-string")
    defer db.Close()

    app.Get("/users", func(c *fiber.Ctx) error {
        // Convert Fiber query params to url.Values
        params := make(url.Values)
        c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
            params.Add(string(key), string(value))
        })

        allowedFields := []string{"id", "name", "email", "age", "status", "created_at"}

        sql, args, err := restql.Parse(params, "users").
            Validate(
                restql.WithAllowedFields(allowedFields),
                restql.WithMaxLimit(100),
            ).
            ToSQL()

        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": err.Error(),
            })
        }

        rows, err := db.Query(sql, args...)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Database error",
            })
        }
        defer rows.Close()

        var users []User
        // Scan rows into users...

        return c.JSON(users)
    })

    app.Listen(":8080")
}
```

</details>

<details>
<summary><b>Chi Router</b></summary>

```go
package main

import (
    "database/sql"
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"
    _ "github.com/lib/pq"
    "github.com/lucasvillarinho/restql"
)

type User struct {
    ID     int    `json:"id"`
    Name   string `json:"name"`
    Email  string `json:"email"`
    Age    int    `json:"age"`
    Status string `json:"status"`
}

func main() {
    r := chi.NewRouter()
    db, _ := sql.Open("postgres", "connection-string")
    defer db.Close()

    r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
        allowedFields := []string{"id", "name", "email", "age", "status", "created_at"}

        sql, args, err := restql.Parse(r.URL.Query(), "users").
            Validate(
                restql.WithAllowedFields(allowedFields),
                restql.WithMaxLimit(100),
                restql.WithMaxOffset(1000),
            ).
            ToSQL()

        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{
                "error": err.Error(),
            })
            return
        }

        rows, err := db.Query(sql, args...)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{
                "error": "Database error",
            })
            return
        }
        defer rows.Close()

        var users []User
        // Scan rows into users...

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    })

    http.ListenAndServe(":8080", r)
}
```

</details>

## License

[![License](https://img.shields.io/github/license/lucasvillarinho/restql)](https://github.com/lucasvillarinho/restql/blob/master/LICENSE)
