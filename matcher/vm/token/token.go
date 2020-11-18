package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENTIFIER = "IDENTIFIER"
	NUMBER     = "NUMBER"
	STRING     = "STRING"

	COMMA = ","
	DOT   = "."

	LEFT_PARENTHESIS  = "("
	RIGHT_PARENTHESIS = ")"
)
