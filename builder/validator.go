package builder

import (
	"fmt"
	"strings"

	"github.com/lucasvillarinho/restql/parser"
)

// Validator validates query parameters against configured rules.
type Validator struct {
	qb            *QueryBuilder
	allowedFields map[string]bool
	maxLimit      *int
	maxOffset     *int
}

// ToSQL builds the SQL query after validating all parameters.
// Returns an error if any validation fails.
func (v *Validator) ToSQL() (string, []any, error) {
	// Validate fields (SELECT clause)
	if len(v.qb.fields) > 0 && len(v.allowedFields) > 0 {
		if err := v.validateFields(v.qb.fields); err != nil {
			return "", nil, err
		}
	}

	// Validate filter (WHERE clause)
	if v.qb.filter != nil && len(v.allowedFields) > 0 {
		if err := v.validateFilter(v.qb.filter); err != nil {
			return "", nil, err
		}
	}

	// Validate sort (ORDER BY clause)
	if len(v.qb.sort) > 0 && len(v.allowedFields) > 0 {
		if err := v.validateSort(v.qb.sort); err != nil {
			return "", nil, err
		}
	}

	// Validate limit and offset
	if err := v.validateLimitOffset(); err != nil {
		return "", nil, err
	}

	// If all validations pass, build SQL
	sql, args, err := v.qb.ToSQL()
	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}

// validateFields validates that all fields in the slice are allowed.
func (v *Validator) validateFields(fields []string) error {
	for _, field := range fields {
		if !v.isFieldAllowed(field) {
			return fmt.Errorf("field '%s' is not allowed. Allowed fields: %v", field, v.allowedFieldsList())
		}
	}
	return nil
}

// validateFilter validates all fields used in the filter AST.
func (v *Validator) validateFilter(filter *parser.Filter) error {
	if filter == nil || filter.Expression == nil {
		return nil
	}
	return v.validateOrExpr(filter.Expression)
}

// validateSort validates the sort fields.
func (v *Validator) validateSort(sort []string) error {
	for _, sortField := range sort {
		// Extract field name (remove - prefix if present)
		field := strings.TrimPrefix(sortField, "-")

		if !v.isFieldAllowed(field) {
			return fmt.Errorf("field '%s' is not allowed. Allowed fields: %v", field, v.allowedFieldsList())
		}
	}
	return nil
}

// validateLimitOffset validates limit and offset against configured maximums.
func (v *Validator) validateLimitOffset() error {
	if v.maxLimit != nil && v.qb.limit > *v.maxLimit {
		return fmt.Errorf("limit %d exceeds maximum allowed limit of %d", v.qb.limit, *v.maxLimit)
	}

	if v.maxOffset != nil && v.qb.offset > *v.maxOffset {
		return fmt.Errorf("offset %d exceeds maximum allowed offset of %d", v.qb.offset, *v.maxOffset)
	}

	return nil
}

// validateOrExpr validates OR expressions recursively.
func (v *Validator) validateOrExpr(expr *parser.OrExpr) error {
	if expr == nil {
		return nil
	}
	for _, andExpr := range expr.And {
		if err := v.validateAndExpr(andExpr); err != nil {
			return err
		}
	}
	return nil
}

// validateAndExpr validates AND expressions recursively.
func (v *Validator) validateAndExpr(expr *parser.AndExpr) error {
	if expr == nil {
		return nil
	}
	for _, comp := range expr.Comparison {
		if err := v.validateComparison(comp); err != nil {
			return err
		}
	}
	return nil
}

// validateComparison validates a comparison expression.
func (v *Validator) validateComparison(comp *parser.Comparison) error {
	if comp == nil {
		return nil
	}

	if comp.Left == nil {
		return nil
	}

	// Validate field name
	if comp.Left.Field != "" {
		field := strings.TrimSpace(comp.Left.Field)
		if !v.isFieldAllowed(field) {
			return fmt.Errorf("field '%s' is not allowed. Allowed fields: %v", field, v.allowedFieldsList())
		}
	}

	// Validate subexpression if present
	if comp.Left.SubExpr != nil {
		return v.validateOrExpr(comp.Left.SubExpr)
	}

	return nil
}

// isFieldAllowed checks if a field is in the whitelist.
func (v *Validator) isFieldAllowed(field string) bool {
	if len(v.allowedFields) == 0 {
		// If no allowed fields are configured, allow all
		return true
	}
	return v.allowedFields[field]
}

// allowedFieldsList returns all allowed fields as a slice for error messages.
func (v *Validator) allowedFieldsList() []string {
	fields := make([]string, 0, len(v.allowedFields))
	for field := range v.allowedFields {
		fields = append(fields, field)
	}
	return fields
}
