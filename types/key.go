package types

func (s *Scanner) key() {
	for !s.isAtEnd() && !s.isKeyEnd() {
		s.advance()
	}

	lexeme := s.Source[s.start:s.current]
	s.addTokenValue(BARE_KEY, string(lexeme))
}

func isKeyStart(b byte) bool {
	return isAlphanumeric(b) || b == '_' || b == '-'
}

func (s *Scanner) isKeyEnd() bool {
	return s.peek() == '=' || s.peek() == ' ' || s.peek() == '\t'
}
