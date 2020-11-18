package lexer

import (
	"github.com/madelyne-io/madelyne/matcher/vm/token"
)

type Lexer struct {
	source  string
	current int
	read    int
	ch      byte
}

func New(source string) *Lexer {
	l := &Lexer{
		source:  source,
		current: 0,
		read:    0,
		ch:      0,
	}
	l.readChar()
	return l
}

func (lexer *Lexer) GetNextToken() token.Token {
	for isWhiteSpace(lexer.ch) {
		lexer.readChar()
	}
	tok := token.Token{
		Type:    "",
		Literal: string(lexer.ch),
	}

	switch lexer.ch {
	case '(':
		tok.Type = token.LEFT_PARENTHESIS
		lexer.readChar()
	case ')':
		tok.Type = token.RIGHT_PARENTHESIS
		lexer.readChar()
	case ',':
		tok.Type = token.COMMA
		lexer.readChar()
	case '.':
		tok.Type = token.DOT
		lexer.readChar()
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
		lexer.readChar()
	case '\'':
		tok.Type = token.STRING
		lexer.readChar()
		tok.Literal = lexer.readString()
	default:
		switch {
		case isLetter(lexer.ch):
			tok.Type = token.IDENTIFIER
			tok.Literal = lexer.readIdentifier()
		case isMinus(lexer.ch):
			tok.Literal = lexer.readNumeric()
			tok.Type = token.NUMBER
		case isNumeric(lexer.ch):
			tok.Literal = lexer.readNumeric()
			tok.Type = token.NUMBER
		default:
			tok.Type = token.ILLEGAL
			lexer.readChar()
		}
	}

	return tok
}

func (lexer *Lexer) readChar() {
	lexer.ch = 0

	if lexer.read < len(lexer.source) {
		lexer.ch = lexer.source[lexer.read]
	}

	lexer.current = lexer.read
	lexer.read += 1
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func (lexer *Lexer) readIdentifier() string {
	start := lexer.current
	for isLetter(lexer.ch) {
		lexer.readChar()
	}
	return lexer.source[start:lexer.current]
}

func isNumeric(ch byte) bool {
	return (ch >= '0' && ch <= '9')
}

func isDot(ch byte) bool {
	return ch == '.'
}

func isMinus(ch byte) bool {
	return ch == '-'
}

func (lexer *Lexer) readNumeric() string {
	start := lexer.current
	if isMinus(lexer.ch) {
		lexer.readChar()
	}
	for isNumeric(lexer.ch) {
		lexer.readChar()
	}

	if isDot(lexer.ch) {
		lexer.readChar()
		for isNumeric(lexer.ch) {
			lexer.readChar()
		}
	}
	return lexer.source[start:lexer.current]
}

func isWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n'
}

func (lexer *Lexer) readString() string {
	start := lexer.current
	for lexer.ch != '\'' && lexer.ch != 0 {
		lexer.readChar()
	}
	defer lexer.readChar()
	return lexer.source[start:lexer.current]
}
