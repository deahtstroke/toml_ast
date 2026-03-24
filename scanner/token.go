package scanner

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
	EQUAL
	HASHTAG

	QUOTE

	BASIC_STRING
	MULTILINE_BASIC_STRING

	LITERAL_STRING
	MULTILINE_LITERAL_STRING

	FLOAT
	INTEGER

	BARE_KEY

	// Reserved keywords
	FALSE
	TRUE
	INF
	NAN

	EOF
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int64
}
