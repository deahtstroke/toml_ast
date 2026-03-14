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
	if isNumberStart(currentChar) {
		s.number()
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

func (s *Scanner) number() {
	for !s.isAtEnd() && (isNumberStart(s.peek()) || s.isValidUnderscore()) {
		s.advance()
	}

	var isFloatingPoint bool
	if s.peek() == '.' && isDigit(s.peekNext()) {

		isFloatingPoint = true
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	lexeme := s.Source[s.start:s.current]

	// Cleanup any underscores
	cleaned := strings.ReplaceAll(string(lexeme), "_", "")

	if isFloatingPoint {
		floatVal, _ := strconv.ParseFloat(cleaned, 64)
		s.addTokenValue(FLOAT, floatVal)
	} else {
		intVal, _ := strconv.Atoi(cleaned)
		s.addTokenValue(INTEGER, intVal)
	}
}

// Valid underscore means that it should be proceded by another digit value
// otherwise is not valid
func (s *Scanner) isValidUnderscore() bool {
	return s.peek() == '_' && isNumberStart(s.peekNext())
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isNumberStart(b byte) bool {
	return isDigit(b) || b == '+' || b == '-'
}

func (s *Scanner) advance() byte {
	curr := s.Source[s.current]
	s.current++
	return curr
}

func (s *Scanner) str() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		panic("Unterminated string")
	}

	// trims the initial and final quotes
	s.advance()
	strValue := s.Source[s.start+1 : s.current-1]

	s.addTokenValue(STRING, string(strValue))
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
