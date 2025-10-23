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

// Params holds parsed query parameters.
type Params struct {
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
	qb := builder.NewQueryBuilder(s.Table())

	if err := parseAndSetFilter(qb, qp.Filter, s); err != nil {
		return nil, err
	}

	if err := validateAndSetFields(qb, qp.Fields, s); err != nil {
		return nil, err
	}

	if err := validateAndSetSort(qb, qp.Sort, s); err != nil {
		return nil, err
	}

	setPagination(qb, qp.Limit, qp.Offset)

	return qb, nil
}

// parseAndSetFilter parses and validates the filter, then sets it in the query builder.
func parseAndSetFilter(qb *builder.QueryBuilder, filter string, s *schema.Schema) error {
	if filter == "" {
		return nil
	}

	parsedFilter, err := parser.ParseFilter(filter)
	if err != nil {
		return err
	}

	if err := s.ValidateFilter(parsedFilter); err != nil {
		return err
	}

	qb.SetFilter(parsedFilter)
	return nil
}

// validateAndSetFields validates the requested fields and sets them in the query builder.
func validateAndSetFields(qb *builder.QueryBuilder, fields []string, s *schema.Schema) error {
	if len(fields) == 0 {
		return nil
	}

	if err := s.ValidateFields(fields); err != nil {
		return err
	}

	qb.SetFields(fields)
	return nil
}

// validateAndSetSort validates the sort fields and sets them in the query builder.
func validateAndSetSort(qb *builder.QueryBuilder, sort []string, s *schema.Schema) error {
	if len(sort) == 0 {
		return nil
	}

	sortFields := make([]string, 0, len(sort))
	for _, sortField := range sort {
		// Extract field name (remove - prefix if present)
		field := strings.TrimPrefix(sortField, "-")

		// Validate field
		if !s.IsFieldAllowed(field) {
			return fmt.Errorf("field '%s' is not allowed. Allowed fields: %v", field, s.AllowedFields())
		}
		sortFields = append(sortFields, sortField)
	}

	qb.SetSort(sortFields)
	return nil
}

// setPagination sets the limit and offset in the query builder.
func setPagination(qb *builder.QueryBuilder, limit, offset int) {
	if limit > 0 {
		qb.SetLimit(limit)
	}
	if offset > 0 {
		qb.SetOffset(offset)
	}
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
func parseQueryParams(params url.Values) *Params {
	return &Params{
		Fields: parseCommaSeparatedList(params.Get("fields")),
		Filter: params.Get("filter"),
		Sort:   parseCommaSeparatedList(params.Get("sort")),
		Limit:  parseIntParam(params, "limit"),
		Offset: parseIntParam(params, "offset"),
	}
}
