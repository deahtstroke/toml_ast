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
	EQUAL
	HASHTAG

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
