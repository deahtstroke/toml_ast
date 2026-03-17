package scanner

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

	// Detects bare-keys only
	if isKeyStart(currentChar) {
		s.key()
		return
	}

	switch currentChar {
	case '#':
		s.comment()
	case '"':
		if s.isMultlineStart() {
			s.multilineBasicString()
		} else {
			s.basicString()
		}
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

// Valid underscore means that it should be proceded by another digit value
// otherwise is not valid
func (s *Scanner) isValidUnderscore() bool {
	return s.peek() == '_' && isNumberStart(s.peekNext())
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isAlphanumeric(b byte) bool {
	return isDigit(b) || (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func (s *Scanner) advance() byte {
	curr := s.Source[s.current]
	s.current++
	return curr
}

// Looks at the value of the source at the current index
// without consuming it
//
// Alias for peekAt0
func (s *Scanner) peek() byte {
	return s.peekAt(0)
}

// Looks at the value of the source at the current index + 1
// without consuming it
//
// Alias for peekAt 1
func (s *Scanner) peekNext() byte {
	return s.peekAt(1)
}

// Looks at the value of the source at the current index + an
// arbitrary offset value without consuming it
func (s *Scanner) peekAt(offset int) byte {
	if s.current+offset >= len(s.Source) {
		return 0
	}

	return s.Source[s.current+offset]
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

// Checks to see if the current pointer is off bounds from the
// length of the source byte array
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}
