package interpreter

import (
	"errors"
	"github.com/madelyne-io/madelyne/matcher/vm/lexer"
	"github.com/madelyne-io/madelyne/matcher/vm/parser"
	"testing"
)

func CheckArg(count int, t *testing.T) func(value interface{}, args []interface{}) error {
	return func(value interface{}, args []interface{}) error {
		if len(args) != count {
			t.Fatalf("failed not enough arg want %d got %d", count, len(args))
		}
		for i, a := range args {
			if a == nil {
				continue
			}
			_, ok := args[0].(error)
			if ok {
				continue
			}
			t.Fatalf("arg %d should be nil or error got %#v", i, a)
		}
		return nil
	}
}

func NopFunc(value interface{}, args []interface{}) error {
	return nil
}

func FailedFunc(ret error) func(value interface{}, args []interface{}) error {
	return func(value interface{}, args []interface{}) error {
		return ret
	}
}

func TestParseProgram1(t *testing.T) {
	customError := errors.New("customError")
	tests := []struct {
		functions map[string]func(value interface{}, args []interface{}) error
		program   string
		expected  error
	}{
		{
			program: ".oneOf(contains(-23.6), maxLength(12), startsWith('abc')).IsEmail()",
			functions: map[string]func(value interface{}, args []interface{}) error{
				"oneOf":      CheckArg(3, t),
				"contains":   NopFunc,
				"maxLength":  NopFunc,
				"startsWith": NopFunc,
				"IsEmail":    NopFunc,
			},
			expected: nil,
		},
		{
			program: ".oneOf(contains(-23.6), maxLength(12), startsWith('abc')).IsEmail()",
			functions: map[string]func(value interface{}, args []interface{}) error{
				"oneOf":      CheckArg(3, t),
				"contains":   NopFunc,
				"maxLength":  NopFunc,
				"startsWith": NopFunc,
				"IsEmail":    FailedFunc(customError),
			},
			expected: customError,
		},
		{
			program:   ".oneOf(contains(-23.6), maxLength(12), startsWith('abc')).IsEmail()",
			functions: map[string]func(value interface{}, args []interface{}) error{},
			expected:  ErrMissingFunc,
		},
	}

	for i, tt := range tests {

		l := lexer.New(tt.program)
		p := parser.New(l)
		ast, err := p.Parse()
		if err != nil {
			t.Fatalf("%d failed %v", i, err)
		}
		int := New(ast, tt.functions)

		err = int.Run(nil)
		if !errors.Is(err, tt.expected) {
			t.Fatalf("%d : exp %v got %v", i, tt.expected, err)
		}
	}
}
