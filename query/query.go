package query

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/lucasvillarinho/restql/builder"
	"github.com/lucasvillarinho/restql/parser"
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
// Validation is optional - use QueryBuilder.Validate() to enable it.
func Parse(params url.Values, table string) (*builder.QueryBuilder, error) {
	// Parse query parameters
	qp := parseQueryParams(params)
	qb := builder.NewQueryBuilder(table)

	// Parse and set filter (no validation)
	if err := parseAndSetFilter(qb, qp.Filter); err != nil {
		return nil, err
	}

	// Set fields (no validation)
	if len(qp.Fields) > 0 {
		qb.SetFields(qp.Fields)
	}

	// Set sort (no validation)
	if len(qp.Sort) > 0 {
		qb.SetSort(qp.Sort)
	}

	// Set pagination
	if qp.Limit > 0 {
		qb.SetLimit(qp.Limit)
	}
	if qp.Offset > 0 {
		qb.SetOffset(qp.Offset)
	}

	return qb, nil
}

// parseAndSetFilter parses the filter and sets it in the query builder.
func parseAndSetFilter(qb *builder.QueryBuilder, filter string) error {
	if filter == "" {
		return nil
	}

	parsedFilter, err := parser.ParseFilter(filter)
	if err != nil {
		return err
	}

	qb.SetFilter(parsedFilter)
	return nil
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
