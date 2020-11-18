package parser

import (
	"errors"
	"github.com/madelyne-io/madelyne/matcher/vm/lexer"
	"testing"
)

func TestParseSucced(t *testing.T) {
	tests := []struct {
		program  string
		expected string
	}{
		{
			program: ".oneOf(contains(-23.6), maxLength(12), startsWith('abc')).IsEmail()",
			expected: `Function oneOf
	Function contains
		Value -23.6
	Function maxLength
		Value 12
	Function startsWith
		Value abc
Function IsEmail
`,
		},
	}

	for i, tt := range tests {

		l := lexer.New(tt.program)
		parser := New(l)
		nodes, err := parser.Parse()
		if err != nil {
			t.Fatalf("%d %v", i, err)
		}

		dump := ""
		for _, n := range nodes {
			dump = dump + n.Dump("")
		}

		if dump != tt.expected {
			t.Fatalf("%d failed got\n '%s' want\n '%s'", i, dump, tt.expected)
		}
	}
}

func TestParseError(t *testing.T) {
	tests := []struct {
		program  string
		expected error
	}{
		{
			program:  ".oneOf(",
			expected: ErrBadToken,
		},
		{
			program:  "&",
			expected: ErrIllegalToken,
		},
		{
			program:  "a",
			expected: ErrBadToken,
		},
		{
			program:  ".(",
			expected: ErrBadToken,
		},
		{
			program:  ".a.",
			expected: ErrBadToken,
		},
		{
			program:  ".a(a().)",
			expected: ErrBadToken,
		},
		{
			program:  ".a(a)",
			expected: ErrBadToken,
		},
	}

	for i, tt := range tests {
		l := lexer.New(tt.program)
		parser := New(l)
		_, err := parser.Parse()
		if !errors.Is(err, tt.expected) {
			t.Fatalf("%d : exp %v got %v", i, tt.expected, err)
		}

	}
}
