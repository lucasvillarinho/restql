// Package restql provides a REST query parameter to SQL converter.
//
// RestQL allows you to convert HTTP query parameters into SQL queries with
// optional validation and security features.
package restql

import (
	"net/url"

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

// SQLBuilder represents any type that can generate SQL queries.
// Both QueryBuilder and Validator implement this interface.
type SQLBuilder interface {
	ToSQL() (string, []any, error)
}

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

// Option is a function that configures a RestQL instance.
// Used for global application-level settings like SQL dialect, placeholder style, etc.
type Option func(*RestQL)

// RestQL holds global configuration for query parsing.
// Use NewRestQL to create an instance with default options that can be
// reused across multiple Parse calls.
type RestQL struct {
	// Future: placeholder style, SQL dialect, naming strategy, logger, etc.
}

// NewRestQL creates a new RestQL instance with global configuration options.
// These options configure application-level settings (e.g., SQL dialect, placeholder style).
// Validation options are passed per-endpoint via Parse().
//
// Example:
//
//	rql := restql.NewRestQL(
//	    // Future: restql.WithPlaceholder("$1"),
//	    // Future: restql.WithDialect("postgres"),
//	)
//	query, err := rql.Parse(params, "users",
//	    restql.WithAllowedFields([]string{"id", "name", "email"}),
//	    restql.WithMaxLimit(100),
//	)
func NewRestQL(opts ...Option) *RestQL {
	rql := &RestQL{}

	for _, opt := range opts {
		opt(rql)
	}

	return rql
}

// Parse parses URL query parameters and returns a SQLBuilder with optional validation.
// Validation options are passed as arguments and applied to this specific query.
//
// Returns a QueryBuilder if no validation options are provided, or a Validator if
// validation options are provided. Both implement the SQLBuilder interface.
//
// Example:
//
//	rql := restql.NewRestQL()
//	query, err := rql.Parse(params, "users",
//	    restql.WithAllowedFields([]string{"id", "name"}),
//	    restql.WithMaxLimit(100),
//	)
//	sql, args, err := query.ToSQL()
func (r *RestQL) Parse(params url.Values, table string, opts ...ValidateOption) (SQLBuilder, error) {
	// Parse query parameters using the query package
	qb, err := query.Parse(params, table)
	if err != nil {
		return nil, err
	}

	// If validation options are provided, apply them
	if len(opts) > 0 {
		validator := qb.Validate(opts...)
		return validator, nil
	}

	return qb, nil
}
