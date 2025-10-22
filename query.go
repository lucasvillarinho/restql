package restql

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
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
func Parse(params url.Values, schema *Schema) (*QueryBuilder, error) {
	// Parse query parameters
	qp := parseQueryParams(params)

	// Create query builder
	qb := NewQueryBuilder(schema.Table())

	// Parse and set filter
	if qp.Filter != "" {
		filter, err := ParseFilter(qp.Filter)
		if err != nil {
			return nil, err
		}

		// Validate filter fields
		if err := schema.ValidateFilter(filter); err != nil {
			return nil, err
		}

		qb.SetFilter(filter)
	}

	// Validate and set fields
	if len(qp.Fields) > 0 {
		if err := schema.ValidateFields(qp.Fields); err != nil {
			return nil, err
		}
		qb.SetFields(qp.Fields)
	}

	// Validate and set sort fields
	if len(qp.Sort) > 0 {
		sortFields := make([]string, 0, len(qp.Sort))
		for _, s := range qp.Sort {
			// Extract field name (remove - prefix if present)
			field := strings.TrimPrefix(s, "-")

			// Validate field
			if !schema.IsFieldAllowed(field) {
				return nil, fmt.Errorf("field '%s' is not allowed. Allowed fields: %v", field, schema.AllowedFields())
			}
			sortFields = append(sortFields, s)
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

// parseQueryParams extracts query parameters from url.Values.
func parseQueryParams(params url.Values) *QueryParams {
	qp := &QueryParams{}

	// Parse fields
	if fields := params.Get("fields"); fields != "" {
		qp.Fields = strings.Split(fields, ",")
		for i := range qp.Fields {
			qp.Fields[i] = strings.TrimSpace(qp.Fields[i])
		}
	}

	// Parse filter
	qp.Filter = params.Get("filter")

	// Parse sort
	if sort := params.Get("sort"); sort != "" {
		qp.Sort = strings.Split(sort, ",")
		for i := range qp.Sort {
			qp.Sort[i] = strings.TrimSpace(qp.Sort[i])
		}
	}

	// Parse limit
	if limit := params.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			qp.Limit = l
		}
	}

	// Parse offset
	if offset := params.Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			qp.Offset = o
		}
	}

	return qp
}
