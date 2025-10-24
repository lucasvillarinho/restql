# `↔️ restQL`

From REST to SQL queries

[![Go Version](https://img.shields.io/github/go-mod/go-version/lucasvillarinho/restql)](https://github.com/lucasvillarinho/restql)
[![Go Report Card](https://goreportcard.com/badge/github.com/lucasvillarinho/restql)](https://goreportcard.com/report/github.com/lucasvillarinho/restql)
[![codecov](https://codecov.io/gh/lucasvillarinho/restql/branch/master/graph/badge.svg)](https://codecov.io/gh/lucasvillarinho/restql)
[![Docs](https://img.shields.io/badge/docs-latest-blue.svg)](https://github.com/lucasvillarinho/restql/tree/master/docs)


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

    // Parse the query parameters
    params, _ := url.ParseQuery("filter=age>18&limit=50")

    // Parse the query parameters and return a QueryBuilder
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

## Operators

RestQL supports a comprehensive set of operators for building complex queries:

- **Comparison**: `=`, `!=`, `<>`, `>`, `<`, `>=`, `<=`
- **Pattern Matching**: `LIKE`, `NOT LIKE`
- **List Operations**: `IN`, `NOT IN`
- **Null Checks**: `IS NULL`, `IS NOT NULL`
- **Logical**: `AND` (`&&`), `OR` (`||`), grouping with `()`

**Examples:**

```go
// Comparison
params, _ := url.ParseQuery("filter=age>18")

// Logical operators
params, _ := url.ParseQuery("filter=age>=18 && status='active'")

// Complex expressions
params, _ := url.ParseQuery("filter=(age>=18 && country='US') || role='admin'")
```

📖 **[View complete operators documentation →](docs/operators.md)**

## Security

RestQL provides built-in security features to protect your application:

```go
// Field whitelisting - prevent unauthorized data access
query.Validate(
    restql.WithAllowedFields([]string{"id", "name", "email"}),
).ToSQL()

// Limit protection - prevent excessive data retrieval
query.Validate(
    restql.WithMaxLimit(100),
    restql.WithMaxOffset(10000),
).ToSQL()
```

🔒 **[View complete security guide →](docs/security.md)**

## Integrations

RestQL works seamlessly with popular Go libraries and frameworks:

**ORMs**: `database/sql` • `GORM` • `sqlx`

**HTTP Frameworks**: `Echo` • `Fiber` • `Chi` • `Gin` • Standard `net/http`

```go
// Example with Echo
e.GET("/users", func(c echo.Context) error {
    sql, args, err := rql.Parse(c.QueryParams(), "users",
        restql.WithAllowedFields([]string{"id", "name", "email"}),
        restql.WithMaxLimit(100),
    ).ToSQL()
    // ... execute query
})
```

🔌 **[View complete integration examples →](docs/integrations.md)**

## License

[![License](https://img.shields.io/github/license/lucasvillarinho/restql)](https://github.com/lucasvillarinho/restql/blob/master/LICENSE)
