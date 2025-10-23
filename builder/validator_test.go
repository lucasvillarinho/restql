package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucasvillarinho/restql/parser"
)

func TestValidator_AllowedFields(t *testing.T) {
	t.Parallel()

	t.Run("filter with allowed fields succeeds", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("age>18 && status='active'")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)

		sql, args, err := qb.Validate(
			WithAllowedFields([]string{"age", "status"}),
		).ToSQL()

		require.NoError(t, err)
		assert.Equal(t, "SELECT * FROM users WHERE (age > ? AND status = ?)", sql)
		assert.Len(t, args, 2)
	})

	t.Run("filter with disallowed field fails", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("age>18 && password='secret'")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)

		_, _, err = qb.Validate(
			WithAllowedFields([]string{"age", "status"}),
		).ToSQL()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "password")
		assert.Contains(t, err.Error(), "not allowed")
	})

	t.Run("select with allowed fields succeeds", func(t *testing.T) {
		t.Parallel()

		qb := NewQueryBuilder("users")
		qb.SetFields([]string{"id", "name", "email"})

		sql, args, err := qb.Validate(
			WithAllowedFields([]string{"id", "name", "email", "age"}),
		).ToSQL()

		require.NoError(t, err)
		assert.Equal(t, "SELECT id, name, email FROM users", sql)
		assert.Empty(t, args)
	})

	t.Run("select with disallowed field fails", func(t *testing.T) {
		t.Parallel()

		qb := NewQueryBuilder("users")
		qb.SetFields([]string{"id", "password"})

		_, _, err := qb.Validate(
			WithAllowedFields([]string{"id", "name", "email"}),
		).ToSQL()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "password")
		assert.Contains(t, err.Error(), "not allowed")
	})

	t.Run("sort with allowed fields succeeds", func(t *testing.T) {
		t.Parallel()

		qb := NewQueryBuilder("users")
		qb.SetSort([]string{"-created_at", "name"})

		sql, args, err := qb.Validate(
			WithAllowedFields([]string{"created_at", "name", "id"}),
		).ToSQL()

		require.NoError(t, err)
		assert.Equal(t, "SELECT * FROM users ORDER BY created_at DESC, name ASC", sql)
		assert.Empty(t, args)
	})

	t.Run("sort with disallowed field fails", func(t *testing.T) {
		t.Parallel()

		qb := NewQueryBuilder("users")
		qb.SetSort([]string{"-password"})

		_, _, err := qb.Validate(
			WithAllowedFields([]string{"id", "name"}),
		).ToSQL()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "password")
		assert.Contains(t, err.Error(), "not allowed")
	})

	t.Run("no allowed fields means all fields allowed", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("password='secret'")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)
		qb.SetFields([]string{"password"})

		sql, args, err := qb.Validate().ToSQL()

		require.NoError(t, err)
		assert.Equal(t, "SELECT password FROM users WHERE password = ?", sql)
		assert.Len(t, args, 1)
	})
}

func TestValidator_MaxLimit(t *testing.T) {
	t.Parallel()

	t.Run("limit within max succeeds", func(t *testing.T) {
		t.Parallel()

		qb := NewQueryBuilder("users")
		qb.SetLimit(50)

		sql, args, err := qb.Validate(
			WithMaxLimit(100),
		).ToSQL()

		require.NoError(t, err)
		assert.Equal(t, "SELECT * FROM users LIMIT 50", sql)
		assert.Empty(t, args)
	})

	t.Run("limit equal to max succeeds", func(t *testing.T) {
		t.Parallel()

		qb := NewQueryBuilder("users")
		qb.SetLimit(100)

		sql, args, err := qb.Validate(
			WithMaxLimit(100),
		).ToSQL()

		require.NoError(t, err)
		assert.Equal(t, "SELECT * FROM users LIMIT 100", sql)
		assert.Empty(t, args)
	})

	t.Run("limit exceeds max fails", func(t *testing.T) {
		t.Parallel()

		qb := NewQueryBuilder("users")
		qb.SetLimit(1000)

		_, _, err := qb.Validate(
			WithMaxLimit(100),
		).ToSQL()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "limit")
		assert.Contains(t, err.Error(), "exceeds")
	})

	t.Run("no limit with max succeeds", func(t *testing.T) {
		t.Parallel()

		qb := NewQueryBuilder("users")

		sql, args, err := qb.Validate(
			WithMaxLimit(100),
		).ToSQL()

		require.NoError(t, err)
		assert.Equal(t, "SELECT * FROM users", sql)
		assert.Empty(t, args)
	})
}

func TestValidator_MaxOffset(t *testing.T) {
	t.Parallel()

	t.Run("offset within max succeeds", func(t *testing.T) {
		t.Parallel()

		qb := NewQueryBuilder("users")
		qb.SetOffset(500)

		sql, args, err := qb.Validate(
			WithMaxOffset(1000),
		).ToSQL()

		require.NoError(t, err)
		assert.Equal(t, "SELECT * FROM users OFFSET 500", sql)
		assert.Empty(t, args)
	})

	t.Run("offset exceeds max fails", func(t *testing.T) {
		t.Parallel()

		qb := NewQueryBuilder("users")
		qb.SetOffset(2000)

		_, _, err := qb.Validate(
			WithMaxOffset(1000),
		).ToSQL()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "offset")
		assert.Contains(t, err.Error(), "exceeds")
	})
}

func TestValidator_Combined(t *testing.T) {
	t.Parallel()

	t.Run("all validations pass", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("age>18 && status='active'")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)
		qb.SetFields([]string{"id", "name", "age"})
		qb.SetSort([]string{"-created_at"})
		qb.SetLimit(50)
		qb.SetOffset(100)

		sql, args, err := qb.Validate(
			WithAllowedFields([]string{"id", "name", "age", "status", "created_at"}),
			WithMaxLimit(100),
			WithMaxOffset(1000),
		).ToSQL()

		require.NoError(t, err)
		assert.Contains(t, sql, "SELECT id, name, age FROM users")
		assert.Contains(t, sql, "WHERE (age > ? AND status = ?)")
		assert.Contains(t, sql, "ORDER BY created_at DESC")
		assert.Contains(t, sql, "LIMIT 50")
		assert.Contains(t, sql, "OFFSET 100")
		assert.Len(t, args, 2)
	})

	t.Run("field validation fails", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("age>18")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)
		qb.SetLimit(50)

		_, _, err = qb.Validate(
			WithAllowedFields([]string{"name", "email"}), // age not allowed
			WithMaxLimit(100),
		).ToSQL()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "age")
	})

	t.Run("limit validation fails", func(t *testing.T) {
		t.Parallel()

		qb := NewQueryBuilder("users")
		qb.SetFields([]string{"id", "name"})
		qb.SetLimit(200)

		_, _, err := qb.Validate(
			WithAllowedFields([]string{"id", "name"}),
			WithMaxLimit(100), // limit too high
		).ToSQL()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "limit")
	})
}

func TestValidator_ComplexFilter(t *testing.T) {
	t.Parallel()

	t.Run("nested filter with all allowed fields", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("(age>18 && status='active') || (role='admin' && verified=true)")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)

		sql, args, err := qb.Validate(
			WithAllowedFields([]string{"age", "status", "role", "verified"}),
		).ToSQL()

		require.NoError(t, err)
		assert.Contains(t, sql, "WHERE")
		assert.Len(t, args, 4)
	})

	t.Run("nested filter with disallowed field", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("(age>18 && password='secret') || role='admin'")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)

		_, _, err = qb.Validate(
			WithAllowedFields([]string{"age", "role"}),
		).ToSQL()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "password")
	})
}
