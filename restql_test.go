package restql_test

import (
	"net/url"
	"testing"

	"github.com/lucasvillarinho/restql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRestQL(t *testing.T) {
	t.Parallel()

	t.Run("creates instance with no options", func(t *testing.T) {
		t.Parallel()
		rql := restql.NewRestQL()
		assert.NotNil(t, rql)
	})

	// Future: when Option types are implemented, test them here
	// Example:
	// t.Run("creates instance with dialect option", func(t *testing.T) {
	//     rql := restql.NewRestQL(
	//         restql.WithDialect("postgres"),
	//     )
	//     assert.NotNil(t, rql)
	// })
}

func TestRestQL_Parse(t *testing.T) {
	t.Parallel()

	t.Run("parses query without validation when no options", func(t *testing.T) {
		t.Parallel()
		rql := restql.NewRestQL()

		params, err := url.ParseQuery("filter=age>18&limit=50")
		require.NoError(t, err)

		query, err := rql.Parse(params, "users")
		require.NoError(t, err)
		require.NotNil(t, query)

		sql, args, err := query.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "WHERE age > ?")
		assert.Equal(t, []any{18}, args)
		assert.Contains(t, sql, "LIMIT 50")
	})

	t.Run("applies validation options in Parse", func(t *testing.T) {
		t.Parallel()
		rql := restql.NewRestQL()

		// Try to request 200 with max limit of 100 (should error because exceeds max)
		params, err := url.ParseQuery("limit=200")
		require.NoError(t, err)

		query, err := rql.Parse(params, "users",
			restql.WithMaxLimit(100),
		)
		require.NoError(t, err)

		_, _, err = query.ToSQL()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "limit 200 exceeds maximum allowed limit of 100")
	})

	t.Run("applies specific validation options", func(t *testing.T) {
		t.Parallel()
		rql := restql.NewRestQL()

		params, err := url.ParseQuery("filter=age>18&fields=id,name,email")
		require.NoError(t, err)

		// Specific allowed fields
		query, err := rql.Parse(params, "users",
			restql.WithAllowedFields([]string{"id", "name", "email", "age"}),
		)
		require.NoError(t, err)

		sql, args, err := query.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "SELECT id, name, email FROM users")
		assert.Equal(t, []any{18}, args)
	})

	t.Run("applies multiple validation options", func(t *testing.T) {
		t.Parallel()
		rql := restql.NewRestQL()

		params, err := url.ParseQuery("filter=age>18&limit=75")
		require.NoError(t, err)

		// Multiple validation options
		query, err := rql.Parse(params, "users",
			restql.WithAllowedFields([]string{"age"}),
			restql.WithMaxLimit(100),
		)
		require.NoError(t, err)

		sql, _, err := query.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "LIMIT 75")
		assert.Contains(t, sql, "WHERE age > ?")
	})

	t.Run("combines multiple validation options", func(t *testing.T) {
		t.Parallel()
		rql := restql.NewRestQL()

		params, err := url.ParseQuery("filter=age>18&limit=50&offset=500")
		require.NoError(t, err)

		// Multiple validation options
		query, err := rql.Parse(params, "users",
			restql.WithAllowedFields([]string{"age", "name"}),
			restql.WithMaxLimit(100),
			restql.WithMaxOffset(1000),
		)
		require.NoError(t, err)

		sql, args, err := query.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "WHERE age > ?")
		assert.Contains(t, sql, "LIMIT 50")
		assert.Contains(t, sql, "OFFSET 500")
		assert.Equal(t, []any{18}, args)
	})

	t.Run("rejects disallowed fields", func(t *testing.T) {
		t.Parallel()
		rql := restql.NewRestQL()

		// Try to filter on 'password' which is not in allowed fields
		params, err := url.ParseQuery("filter=password='secret'")
		require.NoError(t, err)

		query, err := rql.Parse(params, "users",
			restql.WithAllowedFields([]string{"id", "name"}),
		)
		require.NoError(t, err)

		_, _, err = query.ToSQL()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password")
	})

	t.Run("enforces max offset", func(t *testing.T) {
		t.Parallel()
		rql := restql.NewRestQL()

		// Try to offset by 2000 (should error because exceeds max)
		params, err := url.ParseQuery("offset=2000")
		require.NoError(t, err)

		query, err := rql.Parse(params, "users",
			restql.WithMaxOffset(1000),
		)
		require.NoError(t, err)

		_, _, err = query.ToSQL()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "offset 2000 exceeds maximum allowed offset of 1000")
	})

	t.Run("works with complex filters and validation", func(t *testing.T) {
		t.Parallel()
		rql := restql.NewRestQL()

		// URL encode the filter to properly handle special characters
		filterExpr := url.QueryEscape("(age>=18 && status='active') || role='admin'")
		queryStr := "filter=" + filterExpr + "&limit=50&sort=-created_at"
		params, err := url.ParseQuery(queryStr)
		require.NoError(t, err)

		query, err := rql.Parse(params, "users",
			restql.WithAllowedFields([]string{"age", "status", "role", "created_at"}),
			restql.WithMaxLimit(100),
		)
		require.NoError(t, err)

		sql, args, err := query.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "WHERE ((age >= ? AND status = ?) OR role = ?)")
		assert.Contains(t, sql, "ORDER BY created_at DESC")
		assert.Contains(t, sql, "LIMIT 50")
		assert.Equal(t, []any{18, "active", "admin"}, args)
	})

	t.Run("validates limit within maximum", func(t *testing.T) {
		t.Parallel()
		rql := restql.NewRestQL()

		// Request 50 which is within the max of 100
		params, err := url.ParseQuery("limit=50")
		require.NoError(t, err)

		query, err := rql.Parse(params, "users",
			restql.WithMaxLimit(100),
		)
		require.NoError(t, err)

		sql, _, err := query.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "LIMIT 50")
	})
}

func TestRestQL_Compatibility(t *testing.T) {
	t.Parallel()

	t.Run("produces same output as standalone Parse", func(t *testing.T) {
		t.Parallel()
		params, err := url.ParseQuery("filter=age>18&limit=50&sort=-created_at")
		require.NoError(t, err)

		// Using standalone Parse
		query1, err := restql.Parse(params, "users")
		require.NoError(t, err)
		sql1, args1, err := query1.ToSQL()
		require.NoError(t, err)

		// Using NewRestQL().Parse
		rql := restql.NewRestQL()
		query2, err := rql.Parse(params, "users")
		require.NoError(t, err)
		sql2, args2, err := query2.ToSQL()
		require.NoError(t, err)

		assert.Equal(t, sql1, sql2)
		assert.Equal(t, args1, args2)
	})
}
