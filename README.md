# 🔁 restQL

> From REST filters to SQL queries

[![License](https://img.shields.io/github/license/lucasvillarinho/noproblem)](https://github.com/lucasvillarinho/restql/blob/master/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/lucasvillarinho/restql)](https://github.com/lucasvillarinho/restql)
[![Go Report Card](https://goreportcard.com/badge/github.com/lucasvillarinho/restql)](https://goreportcard.com/report/github.com/lucasvillarinho/restql)

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

### Basic Usage (Without Validation)

```go
package main

import (
    "log"
    "net/url"

    "github.com/lucasvillarinho/restql"
)

func main() {
    // Parse query string
    params, _ := url.ParseQuery("filter=(age>=18)&sort=-created_at&limit=10")

    // Build SQL query
    query, err := restql.Parse(params, "users")
    if err != nil {
        log.Fatal(err)
    }

    sql, args, err := query.ToSQL()
    if err != nil {
        log.Fatal(err)
    }

    // Use with your database
    // db.Query(sql, args...)
    // Output: SELECT * FROM users WHERE age >= ? ORDER BY created_at DESC LIMIT 10
}
```

### With Validation

```go
// Parse and validate in one fluent chain
params, _ := url.ParseQuery("filter=(age>=18 && status='active')&limit=1000")

sql, args, err := restql.Parse(params, "users").
    Validate(
        restql.WithAllowedFields([]string{"age", "status", "name", "email"}),
        restql.WithMaxLimit(100),
        restql.WithMaxOffset(1000),
    ).
    ToSQL()

if err != nil {
    // Validation failed - field not allowed or limit exceeded
    log.Fatal(err)
}

// sql: SELECT * FROM users WHERE (age >= ? AND status = ?) LIMIT 100
// args: [18, "active"]
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

- `=` - Equal
- `!=`, `<>` - Not equal
- `>` - Greater than
- `<` - Less than
- `>=` - Greater than or equal
- `<=` - Less than or equal

### Pattern Matching

- `LIKE` - Pattern matching (case-sensitive)
- `ILIKE` - Pattern matching (case-insensitive)
- `NOT LIKE` - Negated pattern matching

### List Operations

- `IN` - Value in list
- `NOT IN` - Value not in list

### Null Checks

- `IS NULL` - Value is null
- `IS NOT NULL` - Value is not null

### Logical Operators

- `&&` - AND
- `||` - OR
- `()` - Grouping

## Examples

### Complex Filters

```go
// Nested conditions with OR and AND
filter := "(age>=18 && status='active') || (role='admin' && verified=true)"
params, _ := url.ParseQuery("filter=" + url.QueryEscape(filter))

query, _ := restql.Parse(params, "users")
sql, args, _ := query.ToSQL()

// SELECT * FROM users WHERE ((age >= ? AND status = ?) OR (role = ? AND verified = ?))
```

### Field Selection

```go
params, _ := url.ParseQuery("fields=id,name,email&filter=age>18")

sql, args, _ := restql.Parse(params, "users").ToSQL()

// SELECT id, name, email FROM users WHERE age > ?
```

### Sorting

```go
params, _ := url.ParseQuery("sort=-created_at,name&limit=10")

sql, args, _ := restql.Parse(params, "users").ToSQL()

// SELECT * FROM users ORDER BY created_at DESC, name ASC LIMIT 10
```

### Full Example with Validation

```go
params, _ := url.ParseQuery(
    "filter=(age>=21 && country='US')&" +
    "fields=id,name,email&" +
    "sort=-created_at&" +
    "limit=50&" +
    "offset=100",
)

sql, args, err := restql.Parse(params, "users").
    Validate(
        restql.WithAllowedFields([]string{
            "id", "name", "email", "age", "country", "created_at",
        }),
        restql.WithMaxLimit(100),
        restql.WithMaxOffset(1000),
    ).
    ToSQL()

// SELECT id, name, email FROM users
// WHERE (age >= ? AND country = ?)
// ORDER BY created_at DESC
// LIMIT 50 OFFSET 100
```

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

### With database/sql

```go
sql, args, err := restql.Parse(params, "users").ToSQL()
if err != nil {
    return err
}

rows, err := db.Query(sql, args...)
```

### With GORM

```go
sql, args, err := restql.Parse(params, "users").
    Validate(restql.WithAllowedFields(allowedFields)).
    ToSQL()
if err != nil {
    return err
}

var users []User
db.Raw(sql, args...).Scan(&users)
```

### With sqlx

```go
sql, args, err := restql.Parse(params, "users").ToSQL()
if err != nil {
    return err
}

var users []User
err = db.Select(&users, sql, args...)
```

## License

MIT License - see [LICENSE](LICENSE) for details
