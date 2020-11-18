package ast

import (
	"github.com/madelyne-io/madelyne/matcher/vm/token"
	"testing"
)

func TestProgram(t *testing.T) {
	program := []Node{
		&NodeFunction{
			Function: token.Token{Type: token.IDENTIFIER, Literal: "oneOf"},
			Arguments: []Node{
				&NodeFunction{
					Function: token.Token{Type: token.IDENTIFIER, Literal: "contains"},
					Arguments: []Node{
						&NodeValue{
							Value: token.Token{Type: token.NUMBER, Literal: "-23.6"},
						},
					},
				},
				&NodeFunction{
					Function: token.Token{Type: token.IDENTIFIER, Literal: "maxLength"},
					Arguments: []Node{
						&NodeValue{
							Value: token.Token{Type: token.NUMBER, Literal: "12"},
						},
					},
				},
				&NodeFunction{
					Function: token.Token{Type: token.IDENTIFIER, Literal: "startsWith"},
					Arguments: []Node{
						&NodeValue{
							Value: token.Token{Type: token.STRING, Literal: "abc"},
						},
					},
				},
			},
		},
		&NodeFunction{
			Function:  token.Token{Type: token.IDENTIFIER, Literal: "IsEmail"},
			Arguments: []Node{},
		},
	}

	expectedDump := `Function oneOf
	Function contains
		Value -23.6
	Function maxLength
		Value 12
	Function startsWith
		Value abc
Function IsEmail
`

	dump := ""
	for _, p := range program {
		dump = dump + p.Dump("")
	}

	findDifferenceInString(t, dump, expectedDump)
}

func findDifferenceInString(t *testing.T, left string, right string) {
	if len(left) != len(right) {
		t.Fatalf("not same length  got\n '%s' want\n '%s'", left, right)
	}

	diff := -1

	for i := 0; i < len(left); i++ {
		if left[i] != right[i] {
			diff = i
		}
	}

	if diff >= 0 {
		t.Errorf("program.Dump() wrong at %d. \ngot=%q\nexp=%q\n", diff, left, right)
		t.Errorf("left was '%s' right was '%s'", string(left[diff-1:diff+1]), string(right[diff-1:diff+1]))
	}
}
