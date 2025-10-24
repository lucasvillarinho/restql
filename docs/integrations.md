# Integration Examples

This guide shows how to integrate RestQL with popular Go database libraries and HTTP frameworks.

## Table of Contents

- [ORMs and Database Libraries](#orms-and-database-libraries)
  - [database/sql](#databasesql)
  - [GORM](#gorm)
  - [sqlx](#sqlx)
- [HTTP Frameworks](#http-frameworks)
  - [Echo Framework](#echo-framework)
  - [Fiber](#fiber)
  - [Chi Router](#chi-router)
- [Tips for Integration](#tips-for-integration)

## ORMs and Database Libraries

### database/sql

Standard library SQL package integration:

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

### GORM

GORM ORM integration with model validation:

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

### sqlx

sqlx library integration with struct scanning:

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

## HTTP Frameworks

### Echo Framework

Echo with reusable RestQL configuration:

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

### Fiber

Fiber framework integration with query parameter conversion:

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

### Chi Router

Chi router integration with standard library patterns:

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

## Tips for Integration

### 1. Reuse RestQL Instances

Create a single RestQL instance and reuse it across handlers:

```go
// Create once at application startup
rql := restql.NewRestQL()

// Reuse in multiple handlers
func handler1() { rql.Parse(...) }
func handler2() { rql.Parse(...) }
```

### 2. Centralize Validation Rules

Define field whitelist and limits in a configuration:

```go
var userAllowedFields = []string{"id", "name", "email", "age"}

func userHandler(c echo.Context) error {
    query, err := rql.Parse(c.QueryParams(), "users",
        restql.WithAllowedFields(userAllowedFields),
        restql.WithMaxLimit(100),
    )
    // ...
}
```

### 3. Error Handling

Always handle validation errors appropriately:

```go
sql, args, err := rql.Parse(params, "users").ToSQL()
if err != nil {
    // Return user-friendly error message
    return c.JSON(http.StatusBadRequest, map[string]string{
        "error": "Invalid query parameters",
        "details": err.Error(),
    })
}
```

### 4. Logging and Monitoring

Log queries for debugging and monitoring:

```go
sql, args, err := rql.Parse(params, "users").ToSQL()
log.Printf("Query: %s, Args: %v", sql, args)
```

