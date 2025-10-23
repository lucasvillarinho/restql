package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucasvillarinho/restql/parser"
)

func TestQueryBuilder_ToSQL(t *testing.T) {
	t.Parallel()
	t.Run("simple equality", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("age=18")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT * FROM users WHERE age = ?", sql)
		assert.Len(t, args, 1)
	})

	t.Run("greater than", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("age>18")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT * FROM users WHERE age > ?", sql)
		assert.Len(t, args, 1)
	})

	t.Run("AND expression", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("age>18 && status='active'")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT * FROM users WHERE (age > ? AND status = ?)", sql)
		assert.Len(t, args, 2)
	})

	t.Run("OR expression", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("age>18 || role='admin'")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT * FROM users WHERE (age > ? OR role = ?)", sql)
		assert.Len(t, args, 2)
	})

	t.Run("with fields", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("age>18")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)
		qb.SetFields([]string{"id", "name", "age"})

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT id, name, age FROM users WHERE age > ?", sql)
		assert.Len(t, args, 1)
	})

	t.Run("with sort ascending", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("age>18")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)
		qb.SetSort([]string{"name"})

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT * FROM users WHERE age > ? ORDER BY name ASC", sql)
		assert.Len(t, args, 1)
	})

	t.Run("with sort descending", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("age>18")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)
		qb.SetSort([]string{"-created_at"})

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT * FROM users WHERE age > ? ORDER BY created_at DESC", sql)
		assert.Len(t, args, 1)
	})

	t.Run("with limit", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("age>18")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)
		qb.SetLimit(10)

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT * FROM users WHERE age > ? LIMIT 10", sql)
		assert.Len(t, args, 1)
	})

	t.Run("with offset", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("age>18")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)
		qb.SetOffset(20)

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT * FROM users WHERE age > ? OFFSET 20", sql)
		assert.Len(t, args, 1)
	})

	t.Run("full query", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("age>18 && status='active'")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)
		qb.SetFields([]string{"id", "name"})
		qb.SetSort([]string{"-created_at", "name"})
		qb.SetLimit(10)
		qb.SetOffset(20)

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT id, name FROM users WHERE (age > ? AND status = ?) ORDER BY created_at DESC, name ASC LIMIT 10 OFFSET 20", sql)
		assert.Len(t, args, 2)
	})

	t.Run("IS NULL", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("deleted_at IS NULL")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT * FROM users WHERE deleted_at IS NULL", sql)
		assert.Empty(t, args)
	})

	t.Run("IN operator", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("status IN ('active', 'pending')")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT * FROM users WHERE status IN (?, ?)", sql)
		assert.Len(t, args, 2)
	})
}

func TestQueryBuilder_Where(t *testing.T) {
	t.Parallel()

	t.Run("basic AND expression", func(t *testing.T) {
		t.Parallel()

		filter, err := parser.ParseFilter("age>18 && status='active'")
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(filter)

		whereSQL, args := qb.Where()

		assert.Equal(t, "(age > ? AND status = ?)", whereSQL)
		assert.Len(t, args, 2)
	})
}

func TestQueryBuilder_NoFilter(t *testing.T) {
	t.Parallel()

	t.Run("query without filter", func(t *testing.T) {
		t.Parallel()

		qb := NewQueryBuilder("users")
		qb.SetFields([]string{"id", "name"})
		qb.SetSort([]string{"-created_at"})
		qb.SetLimit(10)

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT id, name FROM users ORDER BY created_at DESC LIMIT 10", sql)
		assert.Empty(t, args)
	})
}

func TestQueryBuilder_ComplexNesting(t *testing.T) {
	t.Parallel()

	t.Run("complex nested expression", func(t *testing.T) {
		t.Parallel()

		filter := "(age>18 && status='active') || (role='admin' && verified=true)"

		ast, err := parser.ParseFilter(filter)
		require.NoError(t, err)

		qb := NewQueryBuilder("users")
		qb.SetFilter(ast)

		sql, args := qb.ToSQL()

		assert.Contains(t, sql, "WHERE")
		assert.Len(t, args, 4)
	})
}
