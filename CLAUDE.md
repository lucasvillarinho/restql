# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

RestQL is a Go library that converts REST query parameters into SQL queries with built-in validation and security features. It provides a type-safe way to expose database filtering, sorting, and pagination through HTTP APIs.

## Core Architecture

The library is organized into four main packages that form a processing pipeline:

1. **parser** (`parser/`): Parses filter expressions into an Abstract Syntax Tree (AST)
   - Uses the `participle/v2` library for parsing
   - Supports operators: `=`, `!=`, `<>`, `>`, `<`, `>=`, `<=`, `LIKE`, `ILIKE`, `NOT LIKE`, `IN`, `NOT IN`, `IS NULL`, `IS NOT NULL`
   - Supports logical operators: `&&` (AND), `||` (OR)
   - Entry point: `ParseFilter(filter string) (*Filter, error)`

2. **schema** (`schema/`): Defines and validates allowed fields for security
   - Whitelists fields to prevent SQL injection and unauthorized data access
   - Validates filter expressions, sort fields, and select fields against the whitelist
   - Entry point: `NewSchema(table string).AllowFields(...fields)`

3. **builder** (`builder/`): Generates SQL queries from the AST
   - Converts AST into SQL WHERE clauses with parameterized queries
   - Handles SELECT, WHERE, ORDER BY, LIMIT, and OFFSET clauses
   - Uses `?` placeholders for parameter binding
   - Entry point: `NewQueryBuilder(table string)`

4. **query** (`query/`): High-level API that orchestrates the pipeline
   - Parses URL query parameters: `filter`, `fields`, `sort`, `limit`, `offset`
   - Orchestrates validation and SQL generation
   - Entry point: `Parse(params url.Values, schema *Schema) (*QueryBuilder, error)`

## Common Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests in a specific package
go test ./parser
go test ./builder
go test ./schema

# Run a specific test
go test -v -run TestParseFilter ./parser
```

### Linting
```bash
# Run golangci-lint
make lint

# Format code
make fmt
```

### Dependencies
```bash
# Install/update dependencies
go mod tidy

# Download dependencies
go get ./...
```

### Running Examples
```bash
# Run basic example
go run examples/basic/main.go

# Run Echo + GORM example (requires setup)
cd examples/echo-gorm && go run .
```

## Development Notes

### Query Parameter Format

The library expects URL query parameters in this format:
- `filter`: Filter expression (e.g., `status='active' && age>=18`)
- `fields`: Comma-separated list of fields to select (e.g., `id,name,email`)
- `sort`: Comma-separated sort fields, prefix with `-` for DESC (e.g., `-created_at,name`)
- `limit`: Maximum number of results (e.g., `10`)
- `offset`: Number of results to skip (e.g., `20`)

### Security Model

The schema whitelist is critical for security. All field access must be explicitly allowed:
```go
schema := restql.NewSchema("users").AllowFields("id", "name", "email")
```

Fields not in the whitelist will be rejected with an error, preventing:
- SQL injection attacks
- Unauthorized column access
- Information disclosure

### SQL Generation

The builder generates parameterized queries using `?` placeholders. Return values from `ToSQL()` and `Where()` are:
- SQL string with `?` placeholders
- Slice of arguments (`[]any`) to bind to the placeholders

This design is compatible with most Go SQL libraries (database/sql, GORM, sqlx, etc.).

### Test Conventions

- Tests use `testify/assert` and `testify/require` for assertions
- All tests run with `t.Parallel()` for performance
- Subtests use `t.Run()` for organization
