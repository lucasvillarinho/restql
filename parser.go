package restql

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	// filterLexer defines the lexer for filter expressions.
	filterLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "whitespace", Pattern: `\s+`},
		{Name: "Float", Pattern: `[-+]?\d+\.\d+`},
		{Name: "Int", Pattern: `[-+]?\d+`},
		{Name: "String", Pattern: `'[^']*'|"[^"]*"`},
		{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
		{Name: "Operators", Pattern: `>=|<=|!=|<>|&&|\|\||=|>|<`},
		{Name: "Punct", Pattern: `[(),]`},
	})

	// filterParser is the global parser instance.
	filterParser = participle.MustBuild[Filter](
		participle.Lexer(filterLexer),
		participle.Elide("whitespace"),
		participle.UseLookahead(2),
	)
)

// ParseFilter parses a filter string into an AST.
func ParseFilter(filter string) (*Filter, error) {
	if filter == "" {
		return nil, nil
	}

	ast, err := filterParser.ParseString("", filter)
	if err != nil {
		return nil, fmt.Errorf("invalid filter syntax: %s (filter: %s)", err.Error(), filter)
	}

	return ast, nil
}
