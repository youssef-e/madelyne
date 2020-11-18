package vm

import (
	"errors"
	"fmt"
	"github.com/madelyne-io/madelyne/matcher/vm/parser"
	"testing"
)

var (
	errNotWrongArgNumber        = fmt.Errorf("Not : wants one param got")
	errNotNilarg                = fmt.Errorf("Not : provided argument is nil")
	errNotNeitherNilNorErrorArg = fmt.Errorf("Not argument should be nil or an error got")
)

func NotFunc(value interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("%w %d", errNotWrongArgNumber, len(args))
	}
	if args[0] == nil {
		return errNotNilarg
	}
	_, ok := args[0].(error)
	if !ok {
		return fmt.Errorf("%w %#v", errNotNeitherNilNorErrorArg, args[0])
	}
	return nil
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
		functions  map[string]func(value interface{}, args []interface{}) error
		program    string
		buildError error
		runError   error
	}{
		{
			program: ".oneOf(contains(-23.6), maxLength(12), startsWith('abc')).IsEmail()",
			functions: map[string]func(value interface{}, args []interface{}) error{
				"oneOf":      NopFunc,
				"contains":   NopFunc,
				"maxLength":  NopFunc,
				"startsWith": NopFunc,
				"IsEmail":    NopFunc,
			},
			buildError: nil,
			runError:   nil,
		},
		{
			program: ".oneOf(contains(-23.6), maxLength(12), startsWith('abc')).IsEmail()",
			functions: map[string]func(value interface{}, args []interface{}) error{
				"oneOf":      NopFunc,
				"contains":   NopFunc,
				"maxLength":  NopFunc,
				"startsWith": NopFunc,
				"IsEmail":    FailedFunc(customError),
			},
			buildError: nil,
			runError:   customError,
		},
		{
			program:    ".oneOf(",
			functions:  map[string]func(value interface{}, args []interface{}) error{},
			buildError: parser.ErrBadToken,
			runError:   ErrFailedBuilding,
		},
		{
			program: ".Not()",
			functions: map[string]func(value interface{}, args []interface{}) error{
				"Not": NotFunc,
			},
			buildError: nil,
			runError:   errNotWrongArgNumber,
		},
		{
			program: ".Not(Nop())",
			functions: map[string]func(value interface{}, args []interface{}) error{
				"Not": NotFunc,
				"Nop": NopFunc,
			},
			buildError: nil,
			runError:   errNotNilarg,
		},
		{
			program: ".Not(Failed())",
			functions: map[string]func(value interface{}, args []interface{}) error{
				"Not":    NotFunc,
				"Failed": FailedFunc(customError),
			},
			buildError: nil,
			runError:   nil,
		},
	}

	for i, tt := range tests {

		match, err := BuildProgramMatcher(tt.program, tt.functions)
		if !errors.Is(err, tt.buildError) {
			t.Fatalf("%d : exp %v got %v", i, tt.buildError, err)
		}
		if match == nil {
			t.Fatalf("failed nil returned")
		}
		err = match(nil)
		if !errors.Is(err, tt.runError) {
			t.Fatalf("%d : exp %v got %v", i, tt.runError, err)
		}
	}
}
