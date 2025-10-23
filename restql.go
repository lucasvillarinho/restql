// Package restql provides a REST query parameter to SQL converter.
//
// RestQL allows you to convert HTTP query parameters into SQL queries with
// built-in validation and security features.
package restql

import (
	"github.com/lucasvillarinho/restql/builder"
	"github.com/lucasvillarinho/restql/parser"
	"github.com/lucasvillarinho/restql/query"
	"github.com/lucasvillarinho/restql/schema"
)

type (
	// QueryBuilder builds SQL queries from parsed filter expressions.
	QueryBuilder = builder.QueryBuilder

	// Schema defines the allowed fields and table for a query.
	Schema = schema.Schema

	// Filter represents the root of the filter expression tree.
	Filter = parser.Filter

	// QueryParams holds parsed query parameters.
	QueryParams = query.QueryParams
)

var (
	// NewSchema creates a new schema for the given table.
	NewSchema = schema.NewSchema

	// NewQueryBuilder creates a new query builder for the given table.
	NewQueryBuilder = builder.NewQueryBuilder

	// ParseFilter parses a filter string into an AST.
	ParseFilter = parser.ParseFilter

	// Parse parses URL query parameters and returns a QueryBuilder.
	Parse = query.Parse
)
