# Supported Operators

This document provides a comprehensive guide to all operators supported by RestQL.

## Table of Contents

- [Comparison Operators](#comparison-operators)
  - [Equal (=)](#equal-)
  - [Not Equal (!=, <>)](#not-equal--)
  - [Greater Than (>)](#greater-than-)
  - [Less Than (<)](#less-than-)
  - [Greater Than or Equal (>=)](#greater-than-or-equal-)
  - [Less Than or Equal (<=)](#less-than-or-equal-)
- [Pattern Matching](#pattern-matching)
  - [LIKE (case-sensitive)](#like-case-sensitive)
  - [NOT LIKE](#not-like)
- [List Operations](#list-operations)
  - [IN](#in)
  - [NOT IN](#not-in)
- [Null Checks](#null-checks)
  - [IS NULL](#is-null)
  - [IS NOT NULL](#is-not-null)
- [Logical Operators](#logical-operators)
  - [AND (&&)](#and-)
  - [OR (||)](#or-)
  - [Grouping ()](#grouping-)

## Comparison Operators

### Equal (=)

```go
params, _ := url.ParseQuery("filter=status='active'")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE status = ?
// args: ["active"]
```

### Not Equal (!=, <>)

```go
params, _ := url.ParseQuery("filter=status!='inactive'")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE status != ?
// args: ["inactive"]
```

### Greater Than (>)

```go
params, _ := url.ParseQuery("filter=age>18")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE age > ?
// args: [18]
```

### Less Than (<)

```go
params, _ := url.ParseQuery("filter=price<100")
sql, args, _ := restql.Parse(params, "products").ToSQL()
// SELECT * FROM products WHERE price < ?
// args: [100]
```

### Greater Than or Equal (>=)

```go
params, _ := url.ParseQuery("filter=rating>=4.5")
sql, args, _ := restql.Parse(params, "reviews").ToSQL()
// SELECT * FROM reviews WHERE rating >= ?
// args: [4.5]
```

### Less Than or Equal (<=)

```go
params, _ := url.ParseQuery("filter=stock<=10")
sql, args, _ := restql.Parse(params, "products").ToSQL()
// SELECT * FROM products WHERE stock <= ?
// args: [10]
```

## Pattern Matching

### LIKE (case-sensitive)

```go
params, _ := url.ParseQuery("filter=name LIKE '%John%'")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE name LIKE ?
// args: ["%John%"]
```

### NOT LIKE

```go
params, _ := url.ParseQuery("filter=name NOT LIKE '%test%'")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE name NOT LIKE ?
// args: ["%test%"]
```

## List Operations

### IN

```go
params, _ := url.ParseQuery("filter=status IN ('active','pending','approved')")
sql, args, _ := restql.Parse(params, "orders").ToSQL()
// SELECT * FROM orders WHERE status IN (?, ?, ?)
// args: ["active", "pending", "approved"]
```

### NOT IN

```go
params, _ := url.ParseQuery("filter=role NOT IN ('admin','superadmin')")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE role NOT IN (?, ?)
// args: ["admin", "superadmin"]
```

## Null Checks

### IS NULL

```go
params, _ := url.ParseQuery("filter=deleted_at IS NULL")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE deleted_at IS NULL
// args: []
```

### IS NOT NULL

```go
params, _ := url.ParseQuery("filter=email IS NOT NULL")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE email IS NOT NULL
// args: []
```

## Logical Operators

### AND (&&)

```go
params, _ := url.ParseQuery("filter=age>=18 && status='active'")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE (age >= ? AND status = ?)
// args: [18, "active"]
```

### OR (||)

```go
params, _ := url.ParseQuery("filter=role='admin' || role='moderator'")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE (role = ? OR role = ?)
// args: ["admin", "moderator"]
```

### Grouping ()

```go
params, _ := url.ParseQuery("filter=(age>=18 && country='US') || (age>=21 && country='UK')")
sql, args, _ := restql.Parse(params, "users").ToSQL()
// SELECT * FROM users WHERE ((age >= ? AND country = ?) OR (age >= ? AND country = ?))
// args: [18, "US", 21, "UK"]
```

