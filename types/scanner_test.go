package types

import (
	"testing"
)

func Test_CommentNode(t *testing.T) {
	str := `# Some comment here
# Some comment there`

	scanner := Scanner{
		Source:  []byte(str),
		line:    0,
		current: 0,
		start:   0,
	}
	scanner.ScanTokens()

	if len(scanner.Tokens)-1 != 2 { // Minus the EOF token
		t.Fatalf("Expecting two comment tokens parsed: %v", len(scanner.Tokens))
	}
}

func Test_IntegerNode(t *testing.T) {
	tests := map[string]struct {
		source    string
		text      string
		literal   any
		shouldErr bool
	}{
		"whole integer": {
			source:  `12345`,
			text:    "12345",
			literal: int(12345),
		},
		"integer with underscores": {
			source:  `12_345`,
			text:    "12_345",
			literal: int(12_345),
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			s := Scanner{
				Source:  []byte(tt.source),
				start:   0,
				line:    0,
				current: 0,
			}

			s.ScanTokens()

			if s.Tokens[0].Type != INTEGER {
				t.Fatalf("Incorrect amount of tokens: %d", len(s.Tokens)-1)
			}

			if s.Tokens[0].Literal != tt.literal {
				t.Fatalf("Incorrect literal value for token: Expected: %v. Got: %v", tt.literal, s.Tokens[0].Literal)
			}

			if s.Tokens[0].Lexeme != tt.text {
				t.Fatalf("Incorrect text value for token: Expected: %s. Got: %s", tt.text, s.Tokens[0].Lexeme)
			}
		})
	}
}
