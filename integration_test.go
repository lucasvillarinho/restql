package restql

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	schema := NewSchema("products").
		AllowFields("id", "name", "price", "category", "created_at", "stock")

	t.Run("simple filter", func(t *testing.T) {
		t.Parallel()
		params, err := url.ParseQuery("filter=price>100&fields=id,name,price&sort=-created_at&limit=10")
		require.NoError(t, err)

		qb, err := Parse(params, schema)
		require.NoError(t, err)

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT id, name, price FROM products WHERE price > ? ORDER BY created_at DESC LIMIT 10", sql)
		assert.Len(t, args, 1)
	})

	t.Run("complex filter", func(t *testing.T) {
		t.Parallel()
		params, err := url.ParseQuery("filter=" + url.QueryEscape("(price>100 && category='electronics') || stock<10") + "&fields=id,name&limit=20")
		require.NoError(t, err)

		qb, err := Parse(params, schema)
		require.NoError(t, err)

		_, args := qb.ToSQL()

		assert.Len(t, args, 3)
	})

	t.Run("LIKE filter", func(t *testing.T) {
		t.Parallel()
		params, err := url.ParseQuery("filter=" + url.QueryEscape("name LIKE '%phone%'") + "&fields=id,name,price")
		require.NoError(t, err)

		qb, err := Parse(params, schema)
		require.NoError(t, err)

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT id, name, price FROM products WHERE name LIKE ?", sql)
		assert.Len(t, args, 1)
	})

	t.Run("invalid field in filter", func(t *testing.T) {
		t.Parallel()
		params, err := url.ParseQuery("filter=invalid_field>100")
		require.NoError(t, err)

		_, err = Parse(params, schema)

		assert.Error(t, err)
	})

	t.Run("invalid field in select", func(t *testing.T) {
		t.Parallel()
		params, err := url.ParseQuery("fields=invalid_field")
		require.NoError(t, err)

		_, err = Parse(params, schema)

		assert.Error(t, err)
	})

	t.Run("invalid field in sort", func(t *testing.T) {
		t.Parallel()
		params, err := url.ParseQuery("sort=-invalid_field")
		require.NoError(t, err)

		_, err = Parse(params, schema)

		assert.Error(t, err)
	})

	t.Run("IS NULL filter", func(t *testing.T) {
		t.Parallel()
		params, err := url.ParseQuery("filter=stock IS NULL")
		require.NoError(t, err)

		qb, err := Parse(params, schema)
		require.NoError(t, err)

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT * FROM products WHERE stock IS NULL", sql)
		assert.Empty(t, args)
	})

	t.Run("IN operator", func(t *testing.T) {
		t.Parallel()
		params, err := url.ParseQuery("filter=category IN ('electronics', 'books')&fields=id,name")
		require.NoError(t, err)

		qb, err := Parse(params, schema)
		require.NoError(t, err)

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT id, name FROM products WHERE category IN (?, ?)", sql)
		assert.Len(t, args, 2)
	})
}

func TestParse_NoFilter(t *testing.T) {
	t.Parallel()

	t.Run("query without filter", func(t *testing.T) {
		t.Parallel()

		schema := NewSchema("users").
			AllowFields("id", "name", "email")

		params, _ := url.ParseQuery("fields=id,name&sort=-id&limit=5")

		qb, err := Parse(params, schema)
		require.NoError(t, err)

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT id, name FROM users ORDER BY id DESC LIMIT 5", sql)
		assert.Empty(t, args)
	})
}

func TestParse_EmptyParams(t *testing.T) {
	t.Parallel()

	t.Run("empty parameters", func(t *testing.T) {
		t.Parallel()

		schema := NewSchema("users").
			AllowFields("id", "name")

		params := url.Values{}

		qb, err := Parse(params, schema)
		require.NoError(t, err)

		sql, args := qb.ToSQL()

		assert.Equal(t, "SELECT * FROM users", sql)
		assert.Empty(t, args)
	})
}

func TestSchema_Validation(t *testing.T) {
	t.Parallel()

	t.Run("valid filter", func(t *testing.T) {
		t.Parallel()

		schema := NewSchema("products").
			AllowFields("id", "name", "price")

		filter, _ := ParseFilter("price>100")
		err := schema.ValidateFilter(filter)

		assert.NoError(t, err)
	})

	t.Run("invalid filter", func(t *testing.T) {
		t.Parallel()

		schema := NewSchema("products").
			AllowFields("id", "name", "price")

		filter, _ := ParseFilter("invalid_field>100")
		err := schema.ValidateFilter(filter)

		assert.Error(t, err)
	})

	t.Run("valid fields", func(t *testing.T) {
		t.Parallel()

		schema := NewSchema("products").
			AllowFields("id", "name", "price")

		err := schema.ValidateFields([]string{"id", "name"})

		assert.NoError(t, err)
	})

	t.Run("invalid fields", func(t *testing.T) {
		t.Parallel()

		schema := NewSchema("products").
			AllowFields("id", "name", "price")

		err := schema.ValidateFields([]string{"id", "invalid"})

		assert.Error(t, err)
	})
}
