package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFilter_ASTStructure(t *testing.T) {
	t.Parallel()

	t.Run("simple equality - validate complete AST", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("age=18")

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Expression)
		require.Len(t, result.Expression.And, 1)
		require.Len(t, result.Expression.And[0].Comparison, 1)

		comparison := result.Expression.And[0].Comparison[0]
		assert.Equal(t, "age", comparison.Left.Field)
		assert.NotNil(t, comparison.Op)
		assert.True(t, comparison.Op.Equal)
		assert.NotNil(t, comparison.Right)
		assert.NotNil(t, comparison.Right.Int)
		assert.Equal(t, 18, *comparison.Right.Int)
	})

	t.Run("integer value", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("count=42")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.NotNil(t, comparison.Right.Int)
		assert.Equal(t, 42, *comparison.Right.Int)
		assert.Nil(t, comparison.Right.Number)
		assert.Nil(t, comparison.Right.String)
	})

	t.Run("float value", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("price=19.99")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.NotNil(t, comparison.Right.Number)
		assert.Equal(t, 19.99, *comparison.Right.Number)
		assert.Nil(t, comparison.Right.Int)
		assert.Nil(t, comparison.Right.String)
	})

	t.Run("negative integer", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("balance<-100")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.NotNil(t, comparison.Right.Int)
		assert.Equal(t, -100, *comparison.Right.Int)
	})

	t.Run("negative float", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("temperature>-3.5")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.NotNil(t, comparison.Right.Number)
		assert.Equal(t, -3.5, *comparison.Right.Number)
	})

	t.Run("string value with single quotes", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("name='john'")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.NotNil(t, comparison.Right.String)
		assert.Equal(t, "'john'", *comparison.Right.String)
		assert.Nil(t, comparison.Right.Int)
		assert.Nil(t, comparison.Right.Number)
	})

	t.Run("string value with double quotes", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter(`email="test@example.com"`)

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.NotNil(t, comparison.Right.String)
		assert.Equal(t, `"test@example.com"`, *comparison.Right.String)
	})

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("name=''")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.NotNil(t, comparison.Right.String)
		assert.Equal(t, "''", *comparison.Right.String)
	})

	t.Run("boolean true", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("active=true")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.NotNil(t, comparison.Right.Boolean)
		assert.True(t, comparison.Right.Boolean.True)
		assert.False(t, comparison.Right.Boolean.False)
	})

	t.Run("boolean false", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("deleted=false")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.NotNil(t, comparison.Right.Boolean)
		assert.True(t, comparison.Right.Boolean.False)
		assert.False(t, comparison.Right.Boolean.True)
	})

	t.Run("boolean TRUE uppercase", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("active=TRUE")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.NotNil(t, comparison.Right.Boolean)
		assert.True(t, comparison.Right.Boolean.True)
	})

	t.Run("boolean FALSE uppercase", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("deleted=FALSE")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.NotNil(t, comparison.Right.Boolean)
		assert.True(t, comparison.Right.Boolean.False)
	})
}

func TestParseFilter_Operators(t *testing.T) {
	t.Parallel()

	t.Run("equal operator (=)", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("age=18")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.Equal)
		assert.Equal(t, "=", comparison.Op.String())
	})

	t.Run("not equal operator (!=)", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("status!='inactive'")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.NotEqual)
		assert.Equal(t, "!=", comparison.Op.String())
	})

	t.Run("not equal operator (<>)", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("status<>'inactive'")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.NotEqual)
		assert.Equal(t, "!=", comparison.Op.String())
	})

	t.Run("greater than operator (>)", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("age>18")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.Greater)
		assert.Equal(t, ">", comparison.Op.String())
	})

	t.Run("less than operator (<)", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("price<100")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.Less)
		assert.Equal(t, "<", comparison.Op.String())
	})

	t.Run("greater or equal operator (>=)", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("rating>=4.5")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.GreaterOrEqual)
		assert.Equal(t, ">=", comparison.Op.String())
	})

	t.Run("less or equal operator (<=)", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("stock<=10")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.LessOrEqual)
		assert.Equal(t, "<=", comparison.Op.String())
	})

	t.Run("LIKE operator uppercase", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("name LIKE '%john%'")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.Like)
		assert.Equal(t, "LIKE", comparison.Op.String())
	})

	t.Run("LIKE operator lowercase", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("name like '%john%'")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.Like)
		assert.Equal(t, "LIKE", comparison.Op.String())
	})

	t.Run("ILIKE operator uppercase", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("email ILIKE '%@gmail.com'")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.ILike)
		assert.Equal(t, "ILIKE", comparison.Op.String())
	})

	t.Run("ILIKE operator lowercase", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("email ilike '%@gmail.com'")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.ILike)
		assert.Equal(t, "ILIKE", comparison.Op.String())
	})

	t.Run("NOT LIKE operator uppercase", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("name NOT LIKE '%test%'")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.NotLike)
		assert.Equal(t, "NOT LIKE", comparison.Op.String())
	})

	t.Run("NOT LIKE operator lowercase", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("name not like '%test%'")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.NotLike)
		assert.Equal(t, "NOT LIKE", comparison.Op.String())
	})

	t.Run("IN operator uppercase", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("status IN ('active', 'pending')")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.In)
		assert.Equal(t, "IN", comparison.Op.String())
	})

	t.Run("IN operator lowercase", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("status in ('active', 'pending')")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.In)
		assert.Equal(t, "IN", comparison.Op.String())
	})

	t.Run("NOT IN operator uppercase", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("role NOT IN ('admin', 'superadmin')")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.NotIn)
		assert.Equal(t, "NOT IN", comparison.Op.String())
	})

	t.Run("NOT IN operator lowercase", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("role not in ('admin', 'superadmin')")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.NotIn)
		assert.Equal(t, "NOT IN", comparison.Op.String())
	})

	t.Run("IS operator for NULL checks", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("deleted_at IS NULL")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.Is)
		assert.Equal(t, "IS", comparison.Op.String())
	})
}

func TestParseFilter_OperatorPrecedence(t *testing.T) {
	t.Parallel()

	t.Run("AND before OR - no parentheses", func(t *testing.T) {
		t.Parallel()
		// "a=1 && b=2 || c=3" should parse as "(a=1 && b=2) || c=3"
		result, err := ParseFilter("a=1 && b=2 || c=3")

		require.NoError(t, err)
		require.NotNil(t, result.Expression)

		// Should have 2 elements in the OR level
		assert.Len(t, result.Expression.And, 2)

		// First element should have 2 comparisons (a=1 && b=2)
		assert.Len(t, result.Expression.And[0].Comparison, 2)
		assert.Equal(t, "a", result.Expression.And[0].Comparison[0].Left.Field)
		assert.Equal(t, "b", result.Expression.And[0].Comparison[1].Left.Field)

		// Second element should have 1 comparison (c=3)
		assert.Len(t, result.Expression.And[1].Comparison, 1)
		assert.Equal(t, "c", result.Expression.And[1].Comparison[0].Left.Field)
	})

	t.Run("multiple AND operations", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("a=1 && b=2 && c=3")

		require.NoError(t, err)

		// Should have 1 OR element with 3 AND comparisons
		assert.Len(t, result.Expression.And, 1)
		assert.Len(t, result.Expression.And[0].Comparison, 3)
	})

	t.Run("multiple OR operations", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("a=1 || b=2 || c=3")

		require.NoError(t, err)

		// Should have 3 OR elements, each with 1 comparison
		assert.Len(t, result.Expression.And, 3)
		assert.Len(t, result.Expression.And[0].Comparison, 1)
		assert.Len(t, result.Expression.And[1].Comparison, 1)
		assert.Len(t, result.Expression.And[2].Comparison, 1)
	})
}

func TestParseFilter_GroupedExpressions(t *testing.T) {
	t.Parallel()

	t.Run("simple parenthesized expression", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("(age>18)")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]

		// Should create a SubExpr
		assert.NotNil(t, comparison.Left.SubExpr)
		assert.Empty(t, comparison.Left.Field)
	})

	t.Run("complex grouped expression", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("(age>18 && status='active') || role='admin'")

		require.NoError(t, err)
		assert.NotNil(t, result)

		// Should have 2 OR elements
		assert.Len(t, result.Expression.And, 2)

		// First should be a subexpression
		assert.NotNil(t, result.Expression.And[0].Comparison[0].Left.SubExpr)

		// Second should be a simple field comparison
		assert.Equal(t, "role", result.Expression.And[1].Comparison[0].Left.Field)
	})

	t.Run("nested parentheses", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("((age>18))")

		require.NoError(t, err)
		assert.NotNil(t, result)

		// Should create nested SubExpr
		comparison := result.Expression.And[0].Comparison[0]
		assert.NotNil(t, comparison.Left.SubExpr)
	})

	t.Run("multiple groups with OR", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("(age>=18 && country='US') || (age>=21 && country='UK')")

		require.NoError(t, err)
		assert.NotNil(t, result)

		// Should have 2 OR elements, both subexpressions
		assert.Len(t, result.Expression.And, 2)
		assert.NotNil(t, result.Expression.And[0].Comparison[0].Left.SubExpr)
		assert.NotNil(t, result.Expression.And[1].Comparison[0].Left.SubExpr)
	})
}

func TestParseFilter_NullChecks(t *testing.T) {
	t.Parallel()

	t.Run("IS NULL uppercase", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("deleted_at IS NULL")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]

		assert.Equal(t, "deleted_at", comparison.Left.Field)
		assert.True(t, comparison.Op.Is)
		assert.NotNil(t, comparison.Null)
		assert.True(t, comparison.Null.IsNull)
		assert.False(t, comparison.Null.IsNotNull)
	})

	t.Run("IS NULL lowercase", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("deleted_at is null")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]

		assert.True(t, comparison.Op.Is)
		assert.True(t, comparison.Null.IsNull)
	})

	t.Run("IS NOT NULL uppercase", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("email IS NOT NULL")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]

		assert.Equal(t, "email", comparison.Left.Field)
		assert.True(t, comparison.Op.Is)
		assert.NotNil(t, comparison.Null)
		assert.True(t, comparison.Null.IsNotNull)
		assert.False(t, comparison.Null.IsNull)
	})

	t.Run("IS NOT NULL lowercase", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("email is not null")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]

		assert.True(t, comparison.Op.Is)
		assert.True(t, comparison.Null.IsNotNull)
	})

	t.Run("NULL check in complex expression", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("deleted_at IS NULL && status='active'")

		require.NoError(t, err)

		// Should have 2 comparisons in AND
		assert.Len(t, result.Expression.And[0].Comparison, 2)

		// First is NULL check
		assert.True(t, result.Expression.And[0].Comparison[0].Null.IsNull)

		// Second is regular comparison
		assert.Equal(t, "status", result.Expression.And[0].Comparison[1].Left.Field)
	})
}

func TestParseFilter_Arrays(t *testing.T) {
	t.Parallel()

	t.Run("IN with multiple string values", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("status IN ('active', 'pending', 'approved')")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]

		assert.Equal(t, "status", comparison.Left.Field)
		assert.True(t, comparison.Op.In)
		assert.NotNil(t, comparison.Right.Array)
		assert.Len(t, comparison.Right.Array.Values, 3)

		assert.Equal(t, "'active'", *comparison.Right.Array.Values[0].String)
		assert.Equal(t, "'pending'", *comparison.Right.Array.Values[1].String)
		assert.Equal(t, "'approved'", *comparison.Right.Array.Values[2].String)
	})

	t.Run("IN with single value", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("status IN ('active')")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]

		assert.NotNil(t, comparison.Right.Array)
		assert.Len(t, comparison.Right.Array.Values, 1)
		assert.Equal(t, "'active'", *comparison.Right.Array.Values[0].String)
	})

	t.Run("IN with integer values", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("id IN (1, 2, 3, 5, 8)")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]

		assert.NotNil(t, comparison.Right.Array)
		assert.Len(t, comparison.Right.Array.Values, 5)

		assert.Equal(t, 1, *comparison.Right.Array.Values[0].Int)
		assert.Equal(t, 2, *comparison.Right.Array.Values[1].Int)
		assert.Equal(t, 8, *comparison.Right.Array.Values[4].Int)
	})

	t.Run("NOT IN with string values", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("role NOT IN ('admin', 'superadmin')")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]

		assert.True(t, comparison.Op.NotIn)
		assert.NotNil(t, comparison.Right.Array)
		assert.Len(t, comparison.Right.Array.Values, 2)
	})

	t.Run("IN with mixed numeric types", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("value IN (1, 2.5, 3)")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]

		assert.Len(t, comparison.Right.Array.Values, 3)
		assert.NotNil(t, comparison.Right.Array.Values[0].Int)
		assert.NotNil(t, comparison.Right.Array.Values[1].Number)
		assert.NotNil(t, comparison.Right.Array.Values[2].Int)
	})
}

func TestParseFilter_ErrorCases(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		filter string
	}{
		{"unclosed parenthesis", "(age>18"},
		{"invalid operator", "age ~~ 18"},
		{"double operators", "age >> 18"},
		{"missing field name", ">18"},
		{"empty IN array", "status IN ()"},
		{"unclosed IN array", "status IN ('active'"},
		{"missing comma in array", "status IN ('active' 'pending')"},
		{"invalid AND syntax", "age>18 & status='active'"},
		{"invalid OR syntax", "age>18 | status='active'"},
		{"standalone operator", "&&"},
		{"trailing operator", "age>18 &&"},
		{"leading operator", "&& age>18"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := ParseFilter(tc.filter)

			require.Error(t, err, "expected error for filter: %s", tc.filter)
			assert.Nil(t, result)
			if err != nil {
				assert.Contains(t, err.Error(), "invalid filter syntax")
			}
		})
	}
}

func TestOperator_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		operator Operator
		expected string
	}{
		{"equal", Operator{Equal: true}, "="},
		{"not equal", Operator{NotEqual: true}, "!="},
		{"greater", Operator{Greater: true}, ">"},
		{"less", Operator{Less: true}, "<"},
		{"greater or equal", Operator{GreaterOrEqual: true}, ">="},
		{"less or equal", Operator{LessOrEqual: true}, "<="},
		{"like", Operator{Like: true}, "LIKE"},
		{"ilike", Operator{ILike: true}, "ILIKE"},
		{"not like", Operator{NotLike: true}, "NOT LIKE"},
		{"in", Operator{In: true}, "IN"},
		{"not in", Operator{NotIn: true}, "NOT IN"},
		{"is", Operator{Is: true}, "IS"},
		{"empty operator", Operator{}, ""},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, tc.operator.String())
		})
	}
}

func TestBoolean_Value(t *testing.T) {
	t.Parallel()

	t.Run("true value", func(t *testing.T) {
		t.Parallel()
		b := &Boolean{True: true}
		assert.True(t, b.Value())
	})

	t.Run("false value", func(t *testing.T) {
		t.Parallel()
		b := &Boolean{False: true}
		assert.False(t, b.Value())
	})

	t.Run("uninitialized boolean defaults to false", func(t *testing.T) {
		t.Parallel()
		b := &Boolean{}
		assert.False(t, b.Value())
	})
}

func TestParseFilter_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("field name with underscores", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("user_id=123")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.Equal(t, "user_id", comparison.Left.Field)
	})

	t.Run("field name with numbers", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("field123=456")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.Equal(t, "field123", comparison.Left.Field)
	})

	t.Run("multiple spaces between tokens", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("age   >   18")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.Equal(t, "age", comparison.Left.Field)
		assert.True(t, comparison.Op.Greater)
		assert.Equal(t, 18, *comparison.Right.Int)
	})

	t.Run("tabs and spaces mixed", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("age\t>\t18")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.Equal(t, "age", comparison.Left.Field)
	})

	t.Run("very long field name", func(t *testing.T) {
		t.Parallel()
		longField := "this_is_a_very_long_field_name_with_many_underscores_and_characters"
		result, err := ParseFilter(longField + "=1")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.Equal(t, longField, comparison.Left.Field)
	})

	t.Run("zero values", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("count=0")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.Equal(t, 0, *comparison.Right.Int)
	})

	t.Run("zero float", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("balance=0.0")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.Equal(t, 0.0, *comparison.Right.Number)
	})

	t.Run("LIKE with wildcard patterns", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("name LIKE '%John_Doe%'")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.True(t, comparison.Op.Like)
		assert.Equal(t, "'%John_Doe%'", *comparison.Right.String)
	})

	t.Run("complex real-world filter", func(t *testing.T) {
		t.Parallel()
		filter := "(age>=18 && age<=65) && (status='active' || status='pending') && deleted_at IS NULL"
		result, err := ParseFilter(filter)

		require.NoError(t, err)
		assert.NotNil(t, result)
		// Just verify it parses without error - structure validation would be complex
	})

	t.Run("string with special characters", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("email='user@example.com'")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.Equal(t, "'user@example.com'", *comparison.Right.String)
	})

	t.Run("string with spaces", func(t *testing.T) {
		t.Parallel()
		result, err := ParseFilter("name='John Doe'")

		require.NoError(t, err)
		comparison := result.Expression.And[0].Comparison[0]
		assert.Equal(t, "'John Doe'", *comparison.Right.String)
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
