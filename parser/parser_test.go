package parser

import (
	"math"
	"reflect"
	"testing"

	"github.com/deahtstroke/toml-ast/scanner"
)

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
			expValue: "hello world!",
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
