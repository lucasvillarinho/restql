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
<summary><b>ILIKE (case-insensitive)</b></summary>

```go
params, _ := url.ParseQuery("filter=email ILIKE '%@gmail.com'")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE email ILIKE ?
// args: ["%@gmail.com"]
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
