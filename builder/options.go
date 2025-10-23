package builder

// ValidateOption is a function that configures a Validator.
type ValidateOption func(*Validator)

// WithAllowedFields sets the allowed fields whitelist for validation.
// Only fields in this list will be permitted in filters, selects, and sorts.
func WithAllowedFields(fields []string) ValidateOption {
	return func(v *Validator) {
		if v.allowedFields == nil {
			v.allowedFields = make(map[string]bool)
		}
		for _, field := range fields {
			v.allowedFields[field] = true
		}
	}
}

// WithMaxLimit sets the maximum allowed limit value.
// If the query requests a limit greater than this, validation will fail.
func WithMaxLimit(max int) ValidateOption {
	return func(v *Validator) {
		v.maxLimit = &max
	}
}

// WithMaxOffset sets the maximum allowed offset value.
// If the query requests an offset greater than this, validation will fail.
func WithMaxOffset(max int) ValidateOption {
	return func(v *Validator) {
		v.maxOffset = &max
	}
}
