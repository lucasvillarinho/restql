package restql

// Filter represents the root of the filter expression tree.
type Filter struct {
	Expression *OrExpr `parser:"@@"`
}

// OrExpr represents an OR expression (lowest precedence).
type OrExpr struct {
	And []*AndExpr `parser:"@@ ( \"||\" @@ )*"`
}

// AndExpr represents an AND expression.
type AndExpr struct {
	Comparison []*Comparison `parser:"@@ ( \"&&\" @@ )*"`
}

// Comparison represents a comparison operation.
type Comparison struct {
	Left  *Primary   `parser:"@@"`
	Op    *Operator  `parser:"@@?"`
	Right *Value     `parser:"@@?"`
	Null  *NullCheck `parser:"@@?"`
}

// Primary represents a field or a parenthesized expression.
type Primary struct {
	Field   string  `parser:"@Ident |"`
	SubExpr *OrExpr `parser:"\"(\" @@ \")\""`
}

// Operator represents comparison operators.
type Operator struct {
	Equal          bool `parser:"@\"=\""`
	NotEqual       bool `parser:"| @(\"!=\" | \"<>\")"`
	GreaterOrEqual bool `parser:"| @\">=\""`
	LessOrEqual    bool `parser:"| @\"<=\""`
	Greater        bool `parser:"| @\">\""`
	Less           bool `parser:"| @\"<\""`
	Like           bool `parser:"| @(\"LIKE\" | \"like\")"`
	ILike          bool `parser:"| @(\"ILIKE\" | \"ilike\")"`
	NotLike        bool `parser:"| @(\"NOT\" \"LIKE\" | \"not\" \"like\")"`
	In             bool `parser:"| @(\"IN\" | \"in\")"`
	NotIn          bool `parser:"| @(\"NOT\" \"IN\" | \"not\" \"in\")"`
	Is             bool `parser:"| @(\"IS\" | \"is\")"`
}

// String returns the operator as a string.
func (o *Operator) String() string {
	switch {
	case o.Equal:
		return "="
	case o.NotEqual:
		return "!="
	case o.GreaterOrEqual:
		return ">="
	case o.LessOrEqual:
		return "<="
	case o.Greater:
		return ">"
	case o.Less:
		return "<"
	case o.Like:
		return "LIKE"
	case o.ILike:
		return "ILIKE"
	case o.NotLike:
		return "NOT LIKE"
	case o.In:
		return "IN"
	case o.NotIn:
		return "NOT IN"
	case o.Is:
		return "IS"
	default:
		return ""
	}
}

// NullCheck represents NULL checks (IS NULL, IS NOT NULL).
type NullCheck struct {
	IsNull    bool `parser:"@(\"NULL\" | \"null\")"`
	IsNotNull bool `parser:"| @(\"NOT\" \"NULL\" | \"not\" \"null\")"`
}

// Value represents a value in a comparison.
type Value struct {
	String  *string  `parser:"  @String"`
	Number  *float64 `parser:"| @Float"`
	Int     *int     `parser:"| @Int"`
	Boolean *Boolean `parser:"| @@"`
	Array   *Array   `parser:"| @@"`
}

// Boolean represents a boolean value.
type Boolean struct {
	True  bool `parser:"  @(\"true\" | \"TRUE\")"`
	False bool `parser:"| @(\"false\" | \"FALSE\")"`
}

// Value returns the boolean value.
func (b *Boolean) Value() bool {
	return b.True
}

// Array represents an array of values for IN/NOT IN operations.
type Array struct {
	Values []*Value `parser:"\"(\" @@ ( \",\" @@ )* \")\""`
}
