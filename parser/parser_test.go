package parser_test

import(
	"github.com/deahtstroke/toml-ast/parser"
	"github.com/deahtstroke/toml-ast/scanner"
	"testing"
)

func TestParser_Key(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		tokens []scanner.Token
		want   parser.Node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.tokens)
			got := p.Key()
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Key() = %v, want %v", got, tt.want)
			}
		})
	}
}

