package builder

import (
	"fmt"
	"strings"

	"github.com/lucasvillarinho/restql/parser"
)

// QueryBuilder builds SQL queries from parsed filter expressions.
type QueryBuilder struct {
	table            string
	fields           []string
	filter           *parser.Filter
	sort             []string
	limit            int
	offset           int
	args             []any
	placeholderStyle string // Placeholder style: "?", "$1", ":1", etc.
	placeholderCount int    // Counter for numbered placeholders
}

// NewQueryBuilder creates a new query builder for the given table.
func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{
		table:            table,
		args:             make([]any, 0),
		placeholderStyle: "?", // Default to MySQL/SQLite style
	}
}

// SetPlaceholder sets the placeholder style for this query builder.
func (qb *QueryBuilder) SetPlaceholder(style string) *QueryBuilder {
	qb.placeholderStyle = style
	return qb
}

// getPlaceholder returns the next placeholder string based on the configured style.
func (qb *QueryBuilder) getPlaceholder() string {
	if qb.placeholderStyle == "?" {
		return "?"
	}

	// For numbered placeholders like $1, $2, ... or :1, :2, ...
	qb.placeholderCount++
	return fmt.Sprintf("%s%d", qb.placeholderStyle[:1], qb.placeholderCount)
}

// Validate creates a validator for this query with the given options.
// Use this to enable field whitelisting and limit/offset validation.
func (qb *QueryBuilder) Validate(opts ...ValidateOption) *Validator {
	v := &Validator{
		qb:            qb,
		allowedFields: make(map[string]bool),
	}

	for _, opt := range opts {
		opt(v)
	}

	return v
}

// SetFields sets the fields to select.
func (qb *QueryBuilder) SetFields(fields []string) *QueryBuilder {
	qb.fields = fields
	return qb
}

// SetFilter sets the filter expression.
func (qb *QueryBuilder) SetFilter(filter *parser.Filter) *QueryBuilder {
	qb.filter = filter
	return qb
}

// SetSort sets the sort fields.
func (qb *QueryBuilder) SetSort(sort []string) *QueryBuilder {
	qb.sort = sort
	return qb
}

// SetLimit sets the limit.
func (qb *QueryBuilder) SetLimit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

// SetOffset sets the offset.
func (qb *QueryBuilder) SetOffset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

// ToSQL builds the complete SQL query and returns the SQL string and arguments.
// This method does not perform validation. Use Validate().ToSQL() for validated queries.
func (qb *QueryBuilder) ToSQL() (string, []any, error) {
	qb.args = make([]any, 0) // Reset args
	qb.placeholderCount = 0  // Reset placeholder counter

	var sql strings.Builder

	// SELECT clause
	sql.WriteString("SELECT ")
	if len(qb.fields) > 0 {
		sql.WriteString(strings.Join(qb.fields, ", "))
	} else {
		sql.WriteString("*")
	}

	// FROM clause
	sql.WriteString(" FROM ")
	sql.WriteString(qb.table)

	// WHERE clause
	if qb.filter != nil && qb.filter.Expression != nil {
		whereSQL := qb.buildOrExpr(qb.filter.Expression)
		if whereSQL != "" {
			sql.WriteString(" WHERE ")
			sql.WriteString(whereSQL)
		}
	}

	// ORDER BY clause
	if len(qb.sort) > 0 {
		sql.WriteString(" ORDER BY ")
		orderClauses := make([]string, 0, len(qb.sort))
		for _, s := range qb.sort {
			if strings.HasPrefix(s, "-") {
				orderClauses = append(orderClauses, s[1:]+" DESC")
			} else {
				orderClauses = append(orderClauses, s+" ASC")
			}
		}
		sql.WriteString(strings.Join(orderClauses, ", "))
	}

	// LIMIT clause
	if qb.limit > 0 {
		sql.WriteString(fmt.Sprintf(" LIMIT %d", qb.limit))
	}

	// OFFSET clause
	if qb.offset > 0 {
		sql.WriteString(fmt.Sprintf(" OFFSET %d", qb.offset))
	}

	return sql.String(), qb.args, nil
}

// Where builds only the WHERE clause.
func (qb *QueryBuilder) Where() (string, []any) {
	qb.args = make([]any, 0) // Reset args
	qb.placeholderCount = 0  // Reset placeholder counter

	if qb.filter == nil || qb.filter.Expression == nil {
		return "", nil
	}

	whereSQL := qb.buildOrExpr(qb.filter.Expression)
	return whereSQL, qb.args
}

// buildOrExpr builds SQL for OR expressions.
func (qb *QueryBuilder) buildOrExpr(expr *parser.OrExpr) string {
	if expr == nil {
		return ""
	}

	parts := make([]string, 0, len(expr.And))
	for _, andExpr := range expr.And {
		if sql := qb.buildAndExpr(andExpr); sql != "" {
			parts = append(parts, sql)
		}
	}

	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}

	return "(" + strings.Join(parts, " OR ") + ")"
}

// buildAndExpr builds SQL for AND expressions.
func (qb *QueryBuilder) buildAndExpr(expr *parser.AndExpr) string {
	if expr == nil {
		return ""
	}

	parts := make([]string, 0, len(expr.Comparison))
	for _, comp := range expr.Comparison {
		if sql := qb.buildComparison(comp); sql != "" {
			parts = append(parts, sql)
		}
	}

	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}

	return "(" + strings.Join(parts, " AND ") + ")"
}

// buildComparison builds SQL for comparison operations.
func (qb *QueryBuilder) buildComparison(comp *parser.Comparison) string {
	if comp == nil {
		return ""
	}

	// Handle subexpression in parentheses
	if comp.Left != nil && comp.Left.SubExpr != nil {
		return qb.buildOrExpr(comp.Left.SubExpr)
	}

	// Get field name
	field := ""
	if comp.Left != nil {
		field = comp.Left.Field
	}

	if field == "" {
		return ""
	}

	// Handle IS NULL / IS NOT NULL
	if comp.Null != nil {
		if comp.Null.IsNull {
			return field + " IS NULL"
		}
		if comp.Null.IsNotNull {
			return field + " IS NOT NULL"
		}
	}

	// Handle regular operators
	if comp.Op == nil || comp.Right == nil {
		return ""
	}

	operator := comp.Op.String()

	// Handle IN/NOT IN with arrays
	if (comp.Op.In || comp.Op.NotIn) && comp.Right.Array != nil {
		placeholders := make([]string, 0, len(comp.Right.Array.Values))
		for _, val := range comp.Right.Array.Values {
			qb.args = append(qb.args, qb.extractValue(val))
			placeholders = append(placeholders, qb.getPlaceholder())
		}
		return field + " " + operator + " (" + strings.Join(placeholders, ", ") + ")"
	}

	// Handle regular comparison
	value := qb.extractValue(comp.Right)
	qb.args = append(qb.args, value)

	return field + " " + operator + " " + qb.getPlaceholder()
}

// extractValue extracts the actual value from a Value node.
func (qb *QueryBuilder) extractValue(val *parser.Value) any {
	if val == nil {
		return nil
	}

	if val.String != nil {
		// Remove quotes from string
		s := *val.String
		if len(s) >= 2 && (s[0] == '\'' || s[0] == '"') {
			return s[1 : len(s)-1]
		}
		return s
	}

	if val.Int != nil {
		return *val.Int
	}

	if val.Number != nil {
		return *val.Number
	}

	if val.Boolean != nil {
		return val.Boolean.Value()
	}

	return nil
}
