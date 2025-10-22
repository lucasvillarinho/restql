package restql

import (
	"fmt"
	"strings"
)

// Schema defines the allowed fields and table for a query.
type Schema struct {
	table         string
	allowedFields map[string]bool
}

// NewSchema creates a new schema for the given table.
func NewSchema(table string) *Schema {
	return &Schema{
		table:         table,
		allowedFields: make(map[string]bool),
	}
}

// AllowFields adds fields to the whitelist.
func (s *Schema) AllowFields(fields ...string) *Schema {
	for _, field := range fields {
		s.allowedFields[field] = true
	}
	return s
}

// IsFieldAllowed checks if a field is in the whitelist.
func (s *Schema) IsFieldAllowed(field string) bool {
	return s.allowedFields[field]
}

// Table returns the table name.
func (s *Schema) Table() string {
	return s.table
}

// AllowedFields returns all allowed fields as a slice.
func (s *Schema) AllowedFields() []string {
	fields := make([]string, 0, len(s.allowedFields))
	for field := range s.allowedFields {
		fields = append(fields, field)
	}
	return fields
}

// ValidateFields validates that all fields in the slice are allowed.
func (s *Schema) ValidateFields(fields []string) error {
	for _, field := range fields {
		if !s.IsFieldAllowed(field) {
			return fmt.Errorf("field '%s' is not allowed. Allowed fields: %v", field, s.AllowedFields())
		}
	}
	return nil
}

// ValidateFilter validates all fields used in the filter AST.
func (s *Schema) ValidateFilter(filter *Filter) error {
	if filter == nil || filter.Expression == nil {
		return nil
	}
	return s.validateOrExpr(filter.Expression)
}

func (s *Schema) validateOrExpr(expr *OrExpr) error {
	if expr == nil {
		return nil
	}
	for _, andExpr := range expr.And {
		if err := s.validateAndExpr(andExpr); err != nil {
			return err
		}
	}
	return nil
}

func (s *Schema) validateAndExpr(expr *AndExpr) error {
	if expr == nil {
		return nil
	}
	for _, comp := range expr.Comparison {
		if err := s.validateComparison(comp); err != nil {
			return err
		}
	}
	return nil
}

func (s *Schema) validateComparison(comp *Comparison) error {
	if comp == nil {
		return nil
	}

	if comp.Left != nil {
		if comp.Left.Field != "" {
			field := strings.TrimSpace(comp.Left.Field)
			if !s.IsFieldAllowed(field) {
				return fmt.Errorf("field '%s' is not allowed. Allowed fields: %v", field, s.AllowedFields())
			}
		} else if comp.Left.SubExpr != nil {
			if err := s.validateOrExpr(comp.Left.SubExpr); err != nil {
				return err
			}
		}
	}

	return nil
}
