package query

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/lucasvillarinho/restql/builder"
	"github.com/lucasvillarinho/restql/parser"
	"github.com/lucasvillarinho/restql/schema"
)

// QueryParams holds parsed query parameters.
type QueryParams struct {
	Fields []string
	Filter string
	Sort   []string
	Limit  int
	Offset int
}

// Parse parses URL query parameters and returns a QueryBuilder.
func Parse(params url.Values, s *schema.Schema) (*builder.QueryBuilder, error) {
	// Parse query parameters
	qp := parseQueryParams(params)

	// Create query builder
	qb := builder.NewQueryBuilder(s.Table())

	// Parse and set filter
	if qp.Filter != "" {
		filter, err := parser.ParseFilter(qp.Filter)
		if err != nil {
			return nil, err
		}

		// Validate filter fields
		if err := s.ValidateFilter(filter); err != nil {
			return nil, err
		}

		qb.SetFilter(filter)
	}

	// Validate and set fields
	if len(qp.Fields) > 0 {
		if err := s.ValidateFields(qp.Fields); err != nil {
			return nil, err
		}
		qb.SetFields(qp.Fields)
	}

	// Validate and set sort fields
	if len(qp.Sort) > 0 {
		sortFields := make([]string, 0, len(qp.Sort))
		for _, sortField := range qp.Sort {
			// Extract field name (remove - prefix if present)
			field := strings.TrimPrefix(sortField, "-")

			// Validate field
			if !s.IsFieldAllowed(field) {
				return nil, fmt.Errorf("field '%s' is not allowed. Allowed fields: %v", field, s.AllowedFields())
			}
			sortFields = append(sortFields, sortField)
		}
		qb.SetSort(sortFields)
	}

	// Set limit and offset
	if qp.Limit > 0 {
		qb.SetLimit(qp.Limit)
	}
	if qp.Offset > 0 {
		qb.SetOffset(qp.Offset)
	}

	return qb, nil
}

// parseCommaSeparatedList splits a comma-separated string and trims each value.
func parseCommaSeparatedList(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// parseIntParam parses an integer parameter from url.Values.
func parseIntParam(params url.Values, key string) int {
	if value := params.Get(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return 0
}

// parseQueryParams extracts query parameters from url.Values.
func parseQueryParams(params url.Values) *QueryParams {
	return &QueryParams{
		Fields: parseCommaSeparatedList(params.Get("fields")),
		Filter: params.Get("filter"),
		Sort:   parseCommaSeparatedList(params.Get("sort")),
		Limit:  parseIntParam(params, "limit"),
		Offset: parseIntParam(params, "offset"),
	}
}
