# Supported Operators

This document provides a comprehensive guide to all operators supported by RestQL.

## Comparison Operators

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

## Pattern Matching

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

## List Operations

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

## Null Checks

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

## Logical Operators

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

