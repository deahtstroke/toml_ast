package parser

import (
	"math"
	"reflect"
	"testing"

	"github.com/deahtstroke/toml-ast/scanner"
)

func Test_TomlTables(t *testing.T) {
	tests := map[string]struct {
		tokens          []scanner.Token
		expectedLiteral string
		expectedNodes   int
		shouldErr       bool
		err             error
	}{
		"Table with basic string key": {
			tokens: []scanner.Token{
				{
					Type:    scanner.LEFT_BRACKET,
					Lexeme:  "[",
					Literal: string("["),
				},
				{
					Type:    scanner.BARE_KEY,
					Lexeme:  "HelloWorld",
					Literal: string("HelloWorld"),
				},
				{
					Type:    scanner.RIGHT_BRACKET,
					Lexeme:  "]",
					Literal: string("]"),
				},
				{
					Type: scanner.EOF,
				},
			},
			expectedLiteral: "[HelloWorld]",
			expectedNodes:   1,
			shouldErr:       false,
		},
		"Table with basic dotted string key": {
			tokens: []scanner.Token{
				{
					Type:    scanner.LEFT_BRACKET,
					Lexeme:  "[",
					Literal: string("["),
				},
				{
					Type:    scanner.BASIC_STRING,
					Lexeme:  "\"Hello.World\"",
					Literal: string("\"Hello.World\""),
				},
				{
					Type:    scanner.RIGHT_BRACKET,
					Lexeme:  "]",
					Literal: string("]"),
				},
				{
					Type: scanner.EOF,
				},
			},
			expectedLiteral: "[\"Hello.World\"]",
			expectedNodes:   1,
			shouldErr:       false,
		},
		"Table with bare dotted key": {
			tokens: []scanner.Token{
				{
					Type:    scanner.LEFT_BRACKET,
					Lexeme:  "[",
					Literal: string("["),
				},
				{
					Type:    scanner.BARE_KEY,
					Lexeme:  "hello",
					Literal: string("hello"),
				},
				{
					Type:    scanner.DOT,
					Lexeme:  ".",
					Literal: string("."),
				},
				{
					Type:    scanner.BARE_KEY,
					Lexeme:  "world",
					Literal: string("world"),
				},
				{
					Type:    scanner.RIGHT_BRACKET,
					Lexeme:  "]",
					Literal: string("]"),
				},
				{
					Type: scanner.EOF,
				},
			},
			expectedLiteral: "[hello.world]",
			expectedNodes:   1,
			shouldErr:       false,
		},
		"Table with bare dotted key and basic string": {
			tokens: []scanner.Token{
				{
					Type:    scanner.LEFT_BRACKET,
					Lexeme:  "[",
					Literal: string("["),
				},
				{
					Type:    scanner.BASIC_STRING,
					Lexeme:  "\"hello.world\"",
					Literal: string("\"hello.world\""),
				},
				{
					Type:    scanner.DOT,
					Lexeme:  ".",
					Literal: string("."),
				},
				{
					Type:    scanner.BARE_KEY,
					Lexeme:  "bar",
					Literal: string("bar"),
				},
				{
					Type:    scanner.RIGHT_BRACKET,
					Lexeme:  "]",
					Literal: string("]"),
				},
				{
					Type: scanner.EOF,
				},
			},
			expectedLiteral: "[\"hello.world\".bar]",
			expectedNodes:   1,
			shouldErr:       false,
		},
	}

	for test, params := range tests {
		t.Run(test, func(t *testing.T) {
			parser := NewParser(params.tokens)
			doc := parser.Parse()

			length := len(doc.Nodes)
			if length != params.expectedNodes {
				t.Errorf("Incorrect length of nodes for root document node: expected: %d, got: %d", params.expectedNodes, length)
			}

			tokenLiteral := doc.Nodes[0].TokenLiteral()
			if tokenLiteral != params.expectedLiteral {
				t.Errorf("Incorrect token literal. Expected: %s. Got: %s", params.expectedLiteral, tokenLiteral)
			}
		})
	}
}

func Test_Table(t *testing.T) {
	tokens := []scanner.Token{
		{
			Type:    scanner.BASIC_STRING,
			Lexeme:  "HelloWorld",
			Literal: "HelloWorld",
			Line:    0,
		},
		{
			Type:    scanner.RIGHT_BRACKET,
			Lexeme:  "[",
			Literal: "[",
			Line:    0,
		},
	}
	p := NewParser(tokens)
	tableNode := p.Table()

	if tableNode == nil {
		t.Errorf("Not expecting node to be nil")
	}

	keyNode, _ := tableNode.Key.(*KeyNode)
	if keyNode.Segments[0] != "HelloWorld" {
		t.Errorf("Wrong key value. Expecting: HelloWorld. Got: %s", tableNode.Key.TokenLiteral())
	}
}

func Test_Value(t *testing.T) {
	tests := map[string]struct {
		tokens      []scanner.Token
		expNodeType any
		expValue    any
	}{
		"negative integer": {
			tokens: []scanner.Token{
				{
					Type:    scanner.MINUS,
					Lexeme:  "-",
					Literal: "-",
					Line:    1,
				},
				{
					Type:    scanner.INTEGER,
					Lexeme:  "1234",
					Literal: int64(1234),
					Line:    1,
				},
			},
			expNodeType: &IntegerNode{},
			expValue:    int64(-1234),
		},
		"positive integer": {
			tokens: []scanner.Token{
				{
					Type:    scanner.PLUS,
					Lexeme:  "+",
					Literal: "+",
					Line:    1,
				},
				{
					Type:    scanner.INTEGER,
					Lexeme:  "12341",
					Literal: int64(12341),
					Line:    1,
				},
			},
			expNodeType: &IntegerNode{},
			expValue:    int64(12341),
		},
		"unsigned integer": {
			tokens: []scanner.Token{
				{
					Type:    scanner.INTEGER,
					Lexeme:  "12341",
					Literal: int64(12341),
					Line:    1,
				},
			},
			expNodeType: &IntegerNode{},
			expValue:    int64(12341),
		},
		"negative floating point": {
			tokens: []scanner.Token{
				{
					Type:    scanner.MINUS,
					Lexeme:  "-",
					Literal: "-",
					Line:    0,
				},
				{
					Type:    scanner.FLOAT,
					Lexeme:  "3.12451",
					Literal: float64(3.12451),
					Line:    0,
				},
			},
			expNodeType: &FloatNode{},
			expValue:    float64(-3.12451),
		},
		"positive floating point": {
			tokens: []scanner.Token{
				{
					Type:    scanner.PLUS,
					Lexeme:  "+",
					Literal: "+",
					Line:    0,
				},
				{
					Type:    scanner.FLOAT,
					Lexeme:  "3.12451",
					Literal: float64(3.12451),
					Line:    0,
				},
			},
			expNodeType: &FloatNode{},
			expValue:    float64(3.12451),
		},
		"unsigned floating point": {
			tokens: []scanner.Token{
				{
					Type:    scanner.FLOAT,
					Lexeme:  "3.12451",
					Literal: float64(3.12451),
					Line:    0,
				},
			},
			expNodeType: &FloatNode{},
			expValue:    float64(3.12451),
		},
		"negative infinity": {
			tokens: []scanner.Token{
				{
					Type:    scanner.MINUS,
					Lexeme:  "-",
					Literal: "-",
					Line:    0,
				},
				{
					Type:    scanner.INF,
					Lexeme:  "inf",
					Literal: nil,
					Line:    0,
				},
			},
			expNodeType: &IntegerNode{},
			expValue:    -int64(math.MaxInt64),
		},
		"positive infinity": {
			tokens: []scanner.Token{
				{
					Type:    scanner.PLUS,
					Lexeme:  "+",
					Literal: "+",
					Line:    0,
				},
				{
					Type:    scanner.INF,
					Lexeme:  "inf",
					Literal: nil,
					Line:    0,
				},
			},
			expNodeType: &IntegerNode{},
			expValue:    int64(math.MaxInt64),
		},
		"unsigned infinity": {
			tokens: []scanner.Token{
				{
					Type:    scanner.INF,
					Lexeme:  "inf",
					Literal: nil,
					Line:    0,
				},
			},
			expNodeType: &IntegerNode{},
			expValue:    int64(math.MaxInt64),
		},
		"false": {
			tokens: []scanner.Token{
				{
					Type:   scanner.FALSE,
					Lexeme: "false",
					Line:   0,
				},
			},
			expNodeType: &BooleanNode{},
			expValue:    false,
		},
		"true": {
			tokens: []scanner.Token{
				{
					Type:   scanner.TRUE,
					Lexeme: "true",
					Line:   0,
				},
			},
			expNodeType: &BooleanNode{},
			expValue:    true,
		},
		"basic string": {
			tokens: []scanner.Token{
				{
					Type:    scanner.BASIC_STRING,
					Lexeme:  "hello world!",
					Literal: "hello world!",
					Line:    0,
				},
			},
			expNodeType: &StringNode{},
			expValue:    "hello world!",
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			parser := NewParser(tt.tokens)
			node := parser.value()

			gotType := reflect.TypeOf(node)
			expType := reflect.TypeOf(tt.expNodeType)

			if gotType != expType {
				t.Fatalf("Expected node type %s, got %s", expType, gotType)
			}

			gotValue := reflect.ValueOf(node).Elem().FieldByName("Value").Interface()
			if gotValue != tt.expValue {
				t.Fatalf("expected value %v (%T), got %v (%T)", tt.expValue, tt.expValue, gotValue, gotValue)
			}
		})
	}
}
