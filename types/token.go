package types

type TokenType uint32

const (
	_ TokenType = iota
	COMMENT
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SLASH
	STAR
	ASSIGN

	QUOTE

	STRING
	FLOAT
	INTEGER
	KEY

	NEWLINE

	EOF
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int64
}

type Scanner struct {
	Source []byte
	Tokens []Token

	// Internal reading state
	current int
	line    int
	start   int
}

func (s *Scanner) ScanTokens() {
	for s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.Tokens = append(s.Tokens, Token{Type: EOF, Lexeme: "", Line: int64(s.line)})
}

func (s *Scanner) scanToken() {
	// Advance to the next token
	currentChar := s.Source[s.current]
	s.current++

	switch currentChar {
	case '[':
		for s.peek() != '\n' && !s.isAtEnd() {

		}
		s.addToken(LEFT_BRACE)
	case ']':
		s.addToken(RIGHT_BRACE)
	case '\n':
		s.line++
	case '"':
		s.str()
	default:
	}
}

func (s *Scanner) advance() byte {
	s.current++
	return s.Source[s.current]
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

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}

	return s.Source[s.current]
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
