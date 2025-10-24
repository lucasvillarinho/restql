# Security

RestQL provides several security features to protect your application from unauthorized data access and excessive resource usage.

## Table of Contents

- [Field Whitelisting](#field-whitelisting)
- [Limit Protection](#limit-protection)
- [SQL Injection Protection](#sql-injection-protection)
- [Complete Example: Production-Ready Configuration](#complete-example-production-ready-configuration)
- [Best Practices](#best-practices)

## Field Whitelisting

Always validate fields in production to prevent unauthorized data access. This ensures that users can only query, select, or sort by approved fields.

```go
query.Validate(
    restql.WithAllowedFields([]string{"id", "name", "email"}),
).ToSQL()

// Attempting to filter/select/sort on 'password' will fail
```

### Example: Preventing Password Exposure

```go
params, _ := url.ParseQuery("filter=age>18&fields=id,name,password")

allowedFields := []string{"id", "name", "email", "age", "status"}

sql, args, err := restql.Parse(params, "users").
    Validate(restql.WithAllowedFields(allowedFields)).
    ToSQL()

// This will return an error because 'password' is not in the allowed fields
if err != nil {
    // Error: field "password" is not allowed
}
```

## Limit Protection

Prevent excessive data retrieval by setting maximum limits for pagination. This protects your database from performance issues caused by large queries.

```go
query.Validate(
    restql.WithMaxLimit(100),
    restql.WithMaxOffset(10000),
).ToSQL()

// User can't request limit=999999
```

### Example: Enforcing Query Limits

```go
params, _ := url.ParseQuery("filter=status='active'&limit=500&offset=50000")

sql, args, err := restql.Parse(params, "users").
    Validate(
        restql.WithMaxLimit(100),     // Max 100 records per request
        restql.WithMaxOffset(10000),  // Can't skip more than 10000 records
    ).
    ToSQL()

// This will return an error because limit and offset exceed maximums
if err != nil {
    // Error: limit exceeds maximum allowed value
}
```

## SQL Injection Protection

RestQL automatically uses parameterized queries to prevent SQL injection attacks. All user input is properly escaped and passed as arguments.

```go
// User input: filter=name='Robert'; DROP TABLE users;--'
params, _ := url.ParseQuery("filter=name='Robert'; DROP TABLE users;--'")

sql, args, _ := restql.Parse(params, "users").ToSQL()

// Safe output:
// SQL: SELECT * FROM users WHERE name = ?
// Args: ["Robert'; DROP TABLE users;--"]
// The malicious SQL is treated as a string value, not executed
```

## Complete Example: Production-Ready Configuration

```go
package main

import (
    "database/sql"
    "log"
    "net/http"

    "github.com/labstack/echo/v4"
    _ "github.com/lib/pq"
    "github.com/lucasvillarinho/restql"
)

func main() {
    e := echo.New()
    db, _ := sql.Open("postgres", "connection-string")
    defer db.Close()

    // Create RestQL instance with security defaults
    rql := restql.NewRestQL()

    e.GET("/users", func(c echo.Context) error {
        // Define allowed fields (exclude sensitive data)
        allowedFields := []string{
            "id",
            "name",
            "email",
            "age",
            "status",
            "created_at",
            "updated_at",
        }

        // Parse with strict validation
        sql, args, err := rql.Parse(c.QueryParams(), "users",
            restql.WithAllowedFields(allowedFields),
            restql.WithMaxLimit(100),      // Prevent large queries
            restql.WithMaxOffset(10000),   // Prevent deep pagination
        ).ToSQL()

        if err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{
                "error": err.Error(),
            })
        }

        // Execute parameterized query (safe from SQL injection)
        rows, err := db.Query(sql, args...)
        if err != nil {
            log.Printf("Database error: %v", err)
            return c.JSON(http.StatusInternalServerError, map[string]string{
                "error": "Internal server error",
            })
        }
        defer rows.Close()

        // Process and return results...
        return c.JSON(http.StatusOK, results)
    })

    e.Logger.Fatal(e.Start(":8080"))
}
```

## Best Practices

1. **Always use field whitelisting in production** - Never allow users to query arbitrary fields
2. **Set reasonable limits** - Protect your database from expensive queries
3. **Log validation errors** - Monitor for potential abuse attempts
4. **Use HTTPS** - Protect query parameters in transit
5. **Rate limiting** - Combine with API rate limiting for additional protection
6. **Audit logs** - Track who queries what data for compliance

