package parser

import (
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

func (p *Parser) KeyValue() *KeyValueNode {
	key := p.Key()

	if !p.Match(scanner.EQUAL) {
		return nil
	}

	// value := p.value()

	return &KeyValueNode{
		Key:   key,
		// Value: value,
	}
}


// Parse a TOML table which follows the grammar rule:
// table -> LEFT_BRACKET  RIGHT_BRACKET
func (p *Parser) Table() *TableNode {
	if !p.Match(scanner.BARE_KEY, scanner.BASIC_STRING) {
		return nil
	}
	key := p.Key()

	if !p.Match(scanner.RIGHT_BRACE) {
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
