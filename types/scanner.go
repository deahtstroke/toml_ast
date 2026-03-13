package types

import (
	"strconv"
	"strings"
)

type Scanner struct {
	Source []byte
	Tokens []Token

	// Internal reading state
	current int
	line    int
	start   int
}

func (s *Scanner) ScanTokens() {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.Tokens = append(s.Tokens, Token{Type: EOF, Lexeme: "", Line: int64(s.line)})
}

func (s *Scanner) scanToken() {
	currentChar := s.advance()
	if isDigit(currentChar) {
		s.integer()
		return
	}

	switch currentChar {
	case '#':
		s.comment()
	case '"':
		s.str()
	case '=':
		s.addToken(EQUAL)
	case '\n':
		s.line++
	case '\t':
	case ' ':
	case '\r':
		break
	default:
	}
}

func (s *Scanner) integer() {
	for !s.isAtEnd() && (isDigit(s.peek()) || s.isValidUnderscore()) {
		s.advance()
	}

	val := s.Source[s.start:s.current]

	cleaned := strings.ReplaceAll(string(val), "_", "")
	intVal, _ := strconv.Atoi(cleaned)
	s.addTokenValue(INTEGER, intVal)
}

// Valid underscore means that it should be proceded by another digit value
// otherwise is not valid
func (s *Scanner) isValidUnderscore() bool {
	return s.peek() == '_' && isDigit(s.peekNext())
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func (s *Scanner) advance() byte {
	curr := s.Source[s.current]
	s.current++
	return curr
}

func (s *Scanner) str() {
	for s.peek() != '"' && s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		panic("Unterminated string")
	}

	s.advance()
	// trims the initial and final quotes
	strValue := s.Source[s.start+1 : s.current-1]

	s.addTokenValue(STRING, strValue)
}

func (s *Scanner) comment() {
	for s.peek() != '\n' && !s.isAtEnd() {
		s.advance()
	}
	commentValue := s.Source[s.start:s.current]

	// make up for finding a newline character
	s.line++
	s.addTokenValue(COMMENT, commentValue)
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}

	return s.Source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.Source) {
		return 0
	}
	return s.Source[s.current+1]
}

func (s *Scanner) addToken(t TokenType) {
	s.addTokenValue(t, nil)
}

func (s *Scanner) addTokenValue(t TokenType, value any) {
	text := string(s.Source[s.start:s.current])
	s.Tokens = append(s.Tokens, Token{
		Type:    t,
		Lexeme:  text,
		Literal: value,
		Line:    int64(s.line),
	})
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}
