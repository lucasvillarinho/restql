package restql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFilter(t *testing.T) {
	t.Parallel()

	t.Run("simple equality", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("age=18")

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("greater than", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("age>18")

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("string comparison", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("name='john'")

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("AND expression", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("age>18 && status='active'")

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("OR expression", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("age>18 || role='admin'")

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("complex expression", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("(age>18 && status='active') || role='admin'")

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("LIKE operator", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("name LIKE '%john%'")

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("IN operator", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("status IN ('active', 'pending')")

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("IS NULL", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("deleted_at IS NULL")

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("IS NOT NULL", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("email IS NOT NULL")

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestParseFilter_Empty(t *testing.T) {
	t.Parallel()

	t.Run("empty string input", func(t *testing.T) {
		t.Parallel()

		filter, err := ParseFilter("")

		require.NoError(t, err)
		assert.Nil(t, filter)
	})
}
