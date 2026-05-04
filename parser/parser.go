package parser

import (
	"math"
	"slices"

	"github.com/deahtstroke/toml-ast/scanner"
)

type Parser struct {
	Tokens  []scanner.Token
	current int
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{
		Tokens: tokens,
	}
}

func (p *Parser) Parse() *Document {
	documentNode := Document{}
	for !p.isAtEnd() {
		node := p.parseEntry()
		if node != nil {
			documentNode.Nodes = append(documentNode.Nodes, node)
		}
	}
	return &documentNode
}

func (p *Parser) parseEntry() Node {
	switch {
	case p.Match(scanner.LEFT_BRACKET):
		return p.Table()
	case p.Match(scanner.BARE_KEY, scanner.BASIC_STRING):
		return p.KeyValue()
	default:
		p.advance()
		return nil
	}
}

func (p *Parser) KeyValue() *KeyValueNode {
	key := p.Key()

	if !p.Match(scanner.EQUAL) {
		return nil
	}

	value := p.value()

	return &KeyValueNode{
		Key:   key,
		Value: value,
	}
}

func (p *Parser) value() Node {
	if p.Match(scanner.MINUS, scanner.PLUS) {
		operator := p.previous().Type
		switch {
		case p.Match(scanner.FLOAT):
			return createFloatNode(p, operator)
		case p.Match(scanner.INTEGER):
			return createIntNode(p, operator)
		case p.Match(scanner.INF):
			return createInfinityNode(p, operator)
		default:
			return nil
		}
	}

	switch {
	case p.Match(scanner.FLOAT):
		return createFloatNode(p, 0)
	case p.Match(scanner.INTEGER):
		return createIntNode(p, 0)
	case p.Match(scanner.FALSE):
		return createBoolNode(p, scanner.FALSE)
	case p.Match(scanner.TRUE):
		return createBoolNode(p, scanner.TRUE)
	case p.Match(scanner.INF):
		return createInfinityNode(p, 0)
	case p.Match(scanner.BASIC_STRING, scanner.MULTILINE_BASIC_STRING):
		return createStringNode(p)
	default:
	}

	return nil
}

func createStringNode(p *Parser) Node {
	val, ok := p.previous().Literal.(string)
	if !ok {
		return nil
	}

	return &StringNode{
		Value: val,
		Token: p.previous(),
	}
}

func createBoolNode(p *Parser, b scanner.TokenType) Node {
	return &BooleanNode{
		Value: b == scanner.TRUE,
		Token: p.previous(),
	}
}

func createInfinityNode(p *Parser, operator scanner.TokenType) Node {
	val := math.MaxInt64
	if operator == scanner.MINUS {
		val = -val
	}

	return &IntegerNode{
		Value: int64(val),
		Token: p.previous(),
	}
}

func createIntNode(p *Parser, operator scanner.TokenType) Node {
	val, ok := p.previous().Literal.(int64)
	if !ok {
		return nil
	}

	if operator == scanner.MINUS {
		val = -val
	}

	return &IntegerNode{
		Value: val,
		Token: p.previous(),
	}
}

func createFloatNode(p *Parser, operator scanner.TokenType) Node {
	val, ok := p.previous().Literal.(float64)
	if !ok {
		return nil
	}

	if operator == scanner.MINUS {
		val = -val
	}

	return &FloatNode{
		Value: val,
		Token: p.previous(),
	}
}

// Parse a TOML table which follows the grammar rule:
// table -> LEFT_BRACKET  RIGHT_BRACKET
func (p *Parser) Table() *TableNode {
	if !p.Match(scanner.BARE_KEY, scanner.BASIC_STRING) {
		return nil
	}
	key := p.Key()

	if !p.Match(scanner.RIGHT_BRACKET) {
		return nil
	}

	return &TableNode{
		Key: key,
	}
}

// Parse a TOML key which follows the grammar rule:
// key -> (BARE_KEY | STRING) (DOT (BARE_KEY | STRING))*
func (p *Parser) Key() *KeyNode {
	curr := p.previous()
	node := &KeyNode{
		Segments: []string{curr.Literal.(string)},
		Tokens:   []scanner.Token{curr},
	}

	for p.Match(scanner.DOT) {
		if !p.Match(scanner.BASIC_STRING, scanner.BARE_KEY) {
			return nil
		}

		segment := p.previous()
		node.Segments = append(node.Segments, segment.Literal.(string))
		node.Tokens = append(node.Tokens, segment)
	}

	return node
}

func (p *Parser) Match(types ...scanner.TokenType) bool {
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
	}

	return false
}

func (p *Parser) check(token scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == token
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) peek() scanner.Token {
	return p.Tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scanner.EOF
}

func (p *Parser) previous() scanner.Token {
	return p.Tokens[p.current-1]
}
