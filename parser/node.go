package parser

import (
	"strings"

	"github.com/deahtstroke/toml-ast/scanner"
)

type Node interface {
	TokenLiteral() string
}

type Document struct {
	Nodes []Node
}

type TableNode struct {
	Key    Node
	Tokens []scanner.Token
}

func (n *TableNode) TokenLiteral() string {
	return "[" + n.Key.TokenLiteral() + "]"
}

type KeyNode struct {
	Segments []string
	Tokens   []scanner.Token
}

func (n *KeyNode) TokenLiteral() string {
	return strings.Join(n.Segments, ".")
}

type KeyValueNode struct {
	Key    *KeyNode
	Value  Node
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

type BooleanNode struct {
	Value bool
	Token scanner.Token
}

func (n *BooleanNode) TokenLiteral() string {
	if n.Value {
		return "true"
	} else {
		return "false"
	}
}
