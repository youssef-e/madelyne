package lexer

import (
	"github.com/madelyne-io/madelyne/matcher/vm/token"
	"testing"
)

func TestLexing(t *testing.T) {
	tests := []struct {
		program  string
		expected []token.Token
	}{
		{
			program: ".oneOf(contains(-23.6), maxLength(12), startsWith('abc')).IsEmail()",
			expected: []token.Token{
				{Type: token.DOT, Literal: "."},
				{Type: token.IDENTIFIER, Literal: "oneOf"},
				{Type: token.LEFT_PARENTHESIS, Literal: "("},
				{Type: token.IDENTIFIER, Literal: "contains"},
				{Type: token.LEFT_PARENTHESIS, Literal: "("},
				{Type: token.NUMBER, Literal: "-23.6"},
				{Type: token.RIGHT_PARENTHESIS, Literal: ")"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.IDENTIFIER, Literal: "maxLength"},
				{Type: token.LEFT_PARENTHESIS, Literal: "("},
				{Type: token.NUMBER, Literal: "12"},
				{Type: token.RIGHT_PARENTHESIS, Literal: ")"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.IDENTIFIER, Literal: "startsWith"},
				{Type: token.LEFT_PARENTHESIS, Literal: "("},
				{Type: token.STRING, Literal: "abc"},
				{Type: token.RIGHT_PARENTHESIS, Literal: ")"},
				{Type: token.RIGHT_PARENTHESIS, Literal: ")"},
				{Type: token.DOT, Literal: "."},
				{Type: token.IDENTIFIER, Literal: "IsEmail"},
				{Type: token.LEFT_PARENTHESIS, Literal: "("},
				{Type: token.RIGHT_PARENTHESIS, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			program: "test{",
			expected: []token.Token{
				{Type: token.IDENTIFIER, Literal: "test"},
				{Type: token.ILLEGAL, Literal: "{"},
				{Type: token.EOF, Literal: ""},
			},
		},
	}

	for _, tt := range tests {

		lexer := New(tt.program)

		for i, expectedToken := range tt.expected {
			actualToken := lexer.GetNextToken()

			if actualToken.Type != expectedToken.Type {
				t.Fatalf("['%s':%d] - token type wrong. expected=%q, got=%q, literal=%q",
					tt.program, i, expectedToken.Type, actualToken.Type, actualToken.Literal)
			}
			if actualToken.Literal != expectedToken.Literal {
				t.Fatalf("['%s':%d] - token literal wrong. expected=%q, got=%q",
					tt.program, i, expectedToken.Literal, actualToken.Literal)
			}
		}
	}
}
