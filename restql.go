// Package restql provides a REST query parameter to SQL converter.
//
// RestQL allows you to convert HTTP query parameters into SQL queries with
// optional validation and security features.
package restql

import (
	"github.com/lucasvillarinho/restql/builder"
	"github.com/lucasvillarinho/restql/parser"
	"github.com/lucasvillarinho/restql/query"
)

type (
	// QueryBuilder builds SQL queries from parsed filter expressions.
	QueryBuilder = builder.QueryBuilder

	// Validator validates query parameters against configured rules.
	Validator = builder.Validator

	// ValidateOption is a function that configures a Validator.
	ValidateOption = builder.ValidateOption

	// Filter represents the root of the filter expression tree.
	Filter = parser.Filter

	// QueryParams holds parsed query parameters.
	QueryParams = query.Params
)

var (
	// NewQueryBuilder creates a new query builder for the given table.
	NewQueryBuilder = builder.NewQueryBuilder

	// ParseFilter parses a filter string into an AST.
	ParseFilter = parser.ParseFilter

	// Parse parses URL query parameters and returns a QueryBuilder.
	// Validation is optional - use QueryBuilder.Validate() to enable it.
	Parse = query.Parse

	// WithAllowedFields sets the allowed fields whitelist for validation.
	WithAllowedFields = builder.WithAllowedFields

	// WithMaxLimit sets the maximum allowed limit value.
	WithMaxLimit = builder.WithMaxLimit

	// WithMaxOffset sets the maximum allowed offset value.
	WithMaxOffset = builder.WithMaxOffset
)
