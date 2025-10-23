package query

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucasvillarinho/restql/builder"
)

func TestParse(t *testing.T) {
	t.Parallel()

	t.Run("full query with all parameters", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("filter=age>18&fields=id,name,email&sort=-created_at,name&limit=10&offset=20")

		qb, err := Parse(params, "users")

		require.NoError(t, err)
		require.NotNil(t, qb)

		sql, args, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "SELECT id, name, email FROM users")
		assert.Contains(t, sql, "WHERE age > ?")
		assert.Contains(t, sql, "ORDER BY created_at DESC, name ASC")
		assert.Contains(t, sql, "LIMIT 10")
		assert.Contains(t, sql, "OFFSET 20")
		assert.Equal(t, []any{18}, args)
	})

	t.Run("query with only filter", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("filter=status='active'")

		qb, err := Parse(params, "orders")

		require.NoError(t, err)
		sql, args, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "SELECT * FROM orders")
		assert.Contains(t, sql, "WHERE status = ?")
		assert.Equal(t, []any{"active"}, args)
	})

	t.Run("query with only fields", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("fields=id,name,email")

		qb, err := Parse(params, "users")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "SELECT id, name, email FROM users")
		assert.NotContains(t, sql, "WHERE")
	})

	t.Run("query with only sort", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("sort=-created_at,name")

		qb, err := Parse(params, "posts")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "SELECT * FROM posts")
		assert.Contains(t, sql, "ORDER BY created_at DESC, name ASC")
	})

	t.Run("query with only limit", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("limit=50")

		qb, err := Parse(params, "products")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "SELECT * FROM products")
		assert.Contains(t, sql, "LIMIT 50")
		assert.NotContains(t, sql, "OFFSET")
	})

	t.Run("query with only offset", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("offset=100")

		qb, err := Parse(params, "items")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "SELECT * FROM items")
		assert.Contains(t, sql, "OFFSET 100")
	})

	t.Run("query with limit and offset", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("limit=25&offset=50")

		qb, err := Parse(params, "records")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "LIMIT 25")
		assert.Contains(t, sql, "OFFSET 50")
	})

	t.Run("empty query parameters", func(t *testing.T) {
		t.Parallel()
		params := url.Values{}

		qb, err := Parse(params, "users")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Equal(t, "SELECT * FROM users", sql)
	})

	t.Run("query with invalid filter syntax", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("filter=age >> 18")

		qb, err := Parse(params, "users")

		assert.Error(t, err)
		assert.Nil(t, qb)
		assert.Contains(t, err.Error(), "invalid filter syntax")
	})

	t.Run("query with zero limit is ignored", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("limit=0")

		qb, err := Parse(params, "users")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.NotContains(t, sql, "LIMIT")
	})

	t.Run("query with zero offset is ignored", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("offset=0")

		qb, err := Parse(params, "users")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.NotContains(t, sql, "OFFSET")
	})

	t.Run("complex filter with multiple conditions", func(t *testing.T) {
		t.Parallel()
		// Use url.Values to avoid URL encoding issues with quotes
		params := url.Values{}
		params.Set("filter", "(age>=18 && status='active') || role='admin'")
		params.Set("fields", "id,name")
		params.Set("limit", "100")

		qb, err := Parse(params, "users")

		require.NoError(t, err)
		sql, args, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "SELECT id, name FROM users")
		assert.Contains(t, sql, "WHERE")
		assert.Contains(t, sql, "LIMIT 100")
		assert.Len(t, args, 3) // age, status, role
	})

	t.Run("query with special characters in filter", func(t *testing.T) {
		t.Parallel()
		// Use url.Values directly to avoid URL encoding issues with %
		params := url.Values{}
		params.Set("filter", "email LIKE '%@example.com'")

		qb, err := Parse(params, "users")

		require.NoError(t, err)
		sql, args, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "WHERE email LIKE ?")
		assert.Equal(t, []any{"%@example.com"}, args)
	})

	t.Run("fields with spaces are trimmed", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("fields=id, name , email")

		qb, err := Parse(params, "users")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "SELECT id, name, email FROM users")
	})

	t.Run("sort with spaces are trimmed", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("sort=-created_at , name")

		qb, err := Parse(params, "posts")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "ORDER BY created_at DESC, name ASC")
	})
}

func TestParseCommaSeparatedList(t *testing.T) {
	t.Parallel()

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		result := parseCommaSeparatedList("")
		assert.Nil(t, result)
	})

	t.Run("single value", func(t *testing.T) {
		t.Parallel()
		result := parseCommaSeparatedList("id")
		assert.Equal(t, []string{"id"}, result)
	})

	t.Run("multiple values", func(t *testing.T) {
		t.Parallel()
		result := parseCommaSeparatedList("id,name,email")
		assert.Equal(t, []string{"id", "name", "email"}, result)
	})

	t.Run("values with spaces are trimmed", func(t *testing.T) {
		t.Parallel()
		result := parseCommaSeparatedList("id, name , email")
		assert.Equal(t, []string{"id", "name", "email"}, result)
	})

	t.Run("values with leading and trailing spaces", func(t *testing.T) {
		t.Parallel()
		result := parseCommaSeparatedList("  id  ,  name  ,  email  ")
		assert.Equal(t, []string{"id", "name", "email"}, result)
	})

	t.Run("single value with spaces", func(t *testing.T) {
		t.Parallel()
		result := parseCommaSeparatedList("  id  ")
		assert.Equal(t, []string{"id"}, result)
	})

	t.Run("many values", func(t *testing.T) {
		t.Parallel()
		result := parseCommaSeparatedList("a,b,c,d,e,f,g,h,i,j")
		assert.Len(t, result, 10)
		assert.Equal(t, "a", result[0])
		assert.Equal(t, "j", result[9])
	})

	t.Run("values with underscores and numbers", func(t *testing.T) {
		t.Parallel()
		result := parseCommaSeparatedList("user_id,created_at,field123")
		assert.Equal(t, []string{"user_id", "created_at", "field123"}, result)
	})
}

func TestParseIntParam(t *testing.T) {
	t.Parallel()

	t.Run("valid positive integer", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("limit=100")
		result := parseIntParam(params, "limit")
		assert.Equal(t, 100, result)
	})

	t.Run("valid zero", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("limit=0")
		result := parseIntParam(params, "limit")
		assert.Equal(t, 0, result)
	})

	t.Run("parameter not present", func(t *testing.T) {
		t.Parallel()
		params := url.Values{}
		result := parseIntParam(params, "limit")
		assert.Equal(t, 0, result)
	})

	t.Run("invalid integer value", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("limit=abc")
		result := parseIntParam(params, "limit")
		assert.Equal(t, 0, result)
	})

	t.Run("float value", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("limit=10.5")
		result := parseIntParam(params, "limit")
		assert.Equal(t, 0, result) // Can't parse float as int
	})

	t.Run("negative integer", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("offset=-10")
		result := parseIntParam(params, "offset")
		assert.Equal(t, -10, result)
	})

	t.Run("empty string value", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("limit=")
		result := parseIntParam(params, "limit")
		assert.Equal(t, 0, result)
	})

	t.Run("very large integer", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("limit=999999")
		result := parseIntParam(params, "limit")
		assert.Equal(t, 999999, result)
	})
}

func TestParseQueryParams(t *testing.T) {
	t.Parallel()

	t.Run("all parameters present", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("filter=age>18&fields=id,name&sort=-created_at&limit=10&offset=20")

		result := parseQueryParams(params)

		assert.Equal(t, "age>18", result.Filter)
		assert.Equal(t, []string{"id", "name"}, result.Fields)
		assert.Equal(t, []string{"-created_at"}, result.Sort)
		assert.Equal(t, 10, result.Limit)
		assert.Equal(t, 20, result.Offset)
	})

	t.Run("empty parameters", func(t *testing.T) {
		t.Parallel()
		params := url.Values{}

		result := parseQueryParams(params)

		assert.Empty(t, result.Filter)
		assert.Nil(t, result.Fields)
		assert.Nil(t, result.Sort)
		assert.Equal(t, 0, result.Limit)
		assert.Equal(t, 0, result.Offset)
	})

	t.Run("only filter parameter", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("filter=status='active'")

		result := parseQueryParams(params)

		assert.Equal(t, "status='active'", result.Filter)
		assert.Nil(t, result.Fields)
		assert.Nil(t, result.Sort)
		assert.Equal(t, 0, result.Limit)
		assert.Equal(t, 0, result.Offset)
	})

	t.Run("only fields parameter", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("fields=id,name,email")

		result := parseQueryParams(params)

		assert.Empty(t, result.Filter)
		assert.Equal(t, []string{"id", "name", "email"}, result.Fields)
		assert.Nil(t, result.Sort)
		assert.Equal(t, 0, result.Limit)
		assert.Equal(t, 0, result.Offset)
	})

	t.Run("only sort parameter", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("sort=-created_at,name")

		result := parseQueryParams(params)

		assert.Empty(t, result.Filter)
		assert.Nil(t, result.Fields)
		assert.Equal(t, []string{"-created_at", "name"}, result.Sort)
		assert.Equal(t, 0, result.Limit)
		assert.Equal(t, 0, result.Offset)
	})

	t.Run("only pagination parameters", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("limit=50&offset=100")

		result := parseQueryParams(params)

		assert.Empty(t, result.Filter)
		assert.Nil(t, result.Fields)
		assert.Nil(t, result.Sort)
		assert.Equal(t, 50, result.Limit)
		assert.Equal(t, 100, result.Offset)
	})

	t.Run("invalid limit is treated as zero", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("limit=invalid")

		result := parseQueryParams(params)

		assert.Equal(t, 0, result.Limit)
	})

	t.Run("invalid offset is treated as zero", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("offset=invalid")

		result := parseQueryParams(params)

		assert.Equal(t, 0, result.Offset)
	})

	t.Run("multiple values for same field uses first", func(t *testing.T) {
		t.Parallel()
		params := url.Values{}
		params.Add("filter", "age>18")
		params.Add("filter", "status='active'")

		result := parseQueryParams(params)

		// url.Values.Get() returns the first value
		assert.Equal(t, "age>18", result.Filter)
	})
}

func TestParse_WithValidation(t *testing.T) {
	t.Parallel()

	t.Run("parse then validate with allowed fields", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("filter=age>18&fields=id,name,age&sort=-created_at&limit=100")
		allowedFields := []string{"id", "name", "age", "created_at"}

		qb, err := Parse(params, "users")
		require.NoError(t, err)

		sql, args, err := qb.Validate(
			builder.WithAllowedFields(allowedFields),
		).ToSQL()

		require.NoError(t, err)
		assert.Contains(t, sql, "SELECT id, name, age FROM users")
		assert.Contains(t, sql, "WHERE age > ?")
		assert.Equal(t, []any{18}, args)
	})

	t.Run("parse then validate fails with disallowed field in filter", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("filter=password='secret'")
		allowedFields := []string{"id", "name", "email"}

		qb, err := Parse(params, "users")
		require.NoError(t, err)

		_, _, err = qb.Validate(
			builder.WithAllowedFields(allowedFields),
		).ToSQL()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "field 'password' is not allowed")
	})

	t.Run("parse then validate fails with disallowed field in fields", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("fields=id,name,password")
		allowedFields := []string{"id", "name", "email"}

		qb, err := Parse(params, "users")
		require.NoError(t, err)

		_, _, err = qb.Validate(
			builder.WithAllowedFields(allowedFields),
		).ToSQL()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "field 'password' is not allowed")
	})

	t.Run("parse then validate fails with disallowed field in sort", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("sort=-password")
		allowedFields := []string{"id", "name", "email"}

		qb, err := Parse(params, "users")
		require.NoError(t, err)

		_, _, err = qb.Validate(
			builder.WithAllowedFields(allowedFields),
		).ToSQL()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "field 'password' is not allowed")
	})

	t.Run("parse then validate with max limit", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("limit=50")

		qb, err := Parse(params, "users")
		require.NoError(t, err)

		sql, _, err := qb.Validate(
			builder.WithMaxLimit(100),
		).ToSQL()

		require.NoError(t, err)
		assert.Contains(t, sql, "LIMIT 50")
	})

	t.Run("parse then validate fails when limit exceeds max", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("limit=200")

		qb, err := Parse(params, "users")
		require.NoError(t, err)

		_, _, err = qb.Validate(
			builder.WithMaxLimit(100),
		).ToSQL()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "limit 200 exceeds maximum allowed limit of 100")
	})

	t.Run("parse then validate with max offset", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("offset=500")

		qb, err := Parse(params, "users")
		require.NoError(t, err)

		sql, _, err := qb.Validate(
			builder.WithMaxOffset(1000),
		).ToSQL()

		require.NoError(t, err)
		assert.Contains(t, sql, "OFFSET 500")
	})

	t.Run("parse then validate fails when offset exceeds max", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("offset=2000")

		qb, err := Parse(params, "users")
		require.NoError(t, err)

		_, _, err = qb.Validate(
			builder.WithMaxOffset(1000),
		).ToSQL()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "offset 2000 exceeds maximum allowed offset of 1000")
	})
}

func TestParse_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("URL encoded filter", func(t *testing.T) {
		t.Parallel()
		// URL encoded: "age>=18 && status='active'"
		params, _ := url.ParseQuery("filter=age%3E%3D18+%26%26+status%3D%27active%27")

		qb, err := Parse(params, "users")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "WHERE")
	})

	t.Run("very long fields list", func(t *testing.T) {
		t.Parallel()
		fields := "field1,field2,field3,field4,field5,field6,field7,field8,field9,field10"
		params, _ := url.ParseQuery("fields=" + fields)

		qb, err := Parse(params, "table")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "SELECT field1, field2")
		assert.Contains(t, sql, "field10 FROM table")
	})

	t.Run("single field in fields parameter", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("fields=id")

		qb, err := Parse(params, "users")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users", sql)
	})

	t.Run("descending sort without other parameters", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("sort=-created_at")

		qb, err := Parse(params, "posts")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "ORDER BY created_at DESC")
	})

	t.Run("ascending sort explicit", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("sort=name")

		qb, err := Parse(params, "users")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "ORDER BY name ASC")
	})

	t.Run("filter with NULL check", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("filter=deleted_at IS NULL")

		qb, err := Parse(params, "users")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "WHERE deleted_at IS NULL")
	})

	t.Run("filter with IN operator", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("filter=status IN ('active', 'pending')")

		qb, err := Parse(params, "orders")

		require.NoError(t, err)
		sql, args, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "WHERE status IN (?, ?)")
		assert.Equal(t, []any{"active", "pending"}, args)
	})

	t.Run("limit 1 for single record", func(t *testing.T) {
		t.Parallel()
		params, _ := url.ParseQuery("limit=1")

		qb, err := Parse(params, "users")

		require.NoError(t, err)
		sql, _, err := qb.ToSQL()
		require.NoError(t, err)
		assert.Contains(t, sql, "LIMIT 1")
	})
}

func TestParams_Struct(t *testing.T) {
	t.Parallel()

	t.Run("Params struct can be created and accessed", func(t *testing.T) {
		t.Parallel()
		params := &Params{
			Fields: []string{"id", "name"},
			Filter: "age>18",
			Sort:   []string{"-created_at"},
			Limit:  10,
			Offset: 20,
		}

		assert.Equal(t, []string{"id", "name"}, params.Fields)
		assert.Equal(t, "age>18", params.Filter)
		assert.Equal(t, []string{"-created_at"}, params.Sort)
		assert.Equal(t, 10, params.Limit)
		assert.Equal(t, 20, params.Offset)
	})

	t.Run("empty Params struct", func(t *testing.T) {
		t.Parallel()
		params := &Params{}

		assert.Nil(t, params.Fields)
		assert.Empty(t, params.Filter)
		assert.Nil(t, params.Sort)
		assert.Equal(t, 0, params.Limit)
		assert.Equal(t, 0, params.Offset)
	})
}
