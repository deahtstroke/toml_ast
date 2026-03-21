package parser

import (
	"strings"

	"github.com/deahtstroke/toml-ast/scanner"
)

type Node interface {
	TokenLiteral() string
}

type TableNode struct {
	Segment string
	Tokens []scanner.Token
}



type KeyNode struct {
	Segments []string
	Tokens []scanner.Token
}

func (n *KeyNode) TokenLiteral() string{
	var str strings.Builder
	for _, t := range n.Tokens {
		str.WriteString(t.Lexeme + " ")
	}
	return str.String()
}

type KeyValueNode struct {
	Key string
	Value any
	Tokens []scanner.Token
}

func (n *KeyValueNode) TokenLiteral() string {
	var str strings.Builder
	for _, t := range n.Tokens {
		str.WriteString(t.Lexeme + " ")
	}
	return str.String()
}

type StringNode struct {
	Value string
	Token scanner.Token
}

func (n *StringNode) TokenLiteral() string {
	return n.Token.Lexeme
}

type IntegerNode struct {
	Value int64
	Token scanner.Token
}

func (n *IntegerNode) TokenLiteral() string {
	return n.Token.Lexeme
}

type FloatNode struct {
	Value float64
	Token scanner.Token
}

func (n *FloatNode) TokenLiteral() string {
	return n.Token.Lexeme
}
