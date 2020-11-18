package vm

import (
	"fmt"
	"github.com/madelyne-io/madelyne/matcher/vm/interpreter"
	"github.com/madelyne-io/madelyne/matcher/vm/lexer"
	"github.com/madelyne-io/madelyne/matcher/vm/parser"
)

var (
	ErrFailedBuilding = fmt.Errorf("ErrFailedBuilding")
)

func BuildProgramMatcher(program string, functions map[string]func(value interface{}, args []interface{}) error) (func(interface{}) error, error) {
	ast, err := parser.New(lexer.New(program)).Parse()
	if err != nil {
		return func(value interface{}) error {
			return ErrFailedBuilding
		}, err
	}
	int := interpreter.New(ast, functions)
	return func(value interface{}) error {
		return int.Run(value)
	}, nil
}
