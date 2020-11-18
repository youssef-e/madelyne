package interpreter

import (
	"fmt"
	"github.com/madelyne-io/madelyne/matcher/vm/ast"
)

var (
	ErrMissingFunc = fmt.Errorf("ErrMissingFunc")
)

type Interpreter struct {
	ast []*ast.NodeFunction
	fn  map[string]func(value interface{}, args []interface{}) error
}

func New(astFunc []*ast.NodeFunction, functions map[string]func(value interface{}, args []interface{}) error) *Interpreter {
	return &Interpreter{
		ast: astFunc,
		fn:  functions,
	}
}

func (i *Interpreter) runFunction(value interface{}, nf *ast.NodeFunction) error {
	fn, ok := i.fn[nf.Token().Literal]
	if !ok {
		return fmt.Errorf("%w : %s", ErrMissingFunc, nf.Token().Literal)
	}
	args := []interface{}{}

	for _, a := range nf.Arguments {
		switch a.(type) {
		case *ast.NodeValue:
			args = append(args, a.Token().Literal)
		case *ast.NodeFunction:
			err := i.runFunction(value, a.(*ast.NodeFunction))
			args = append(args, err)
		}
	}
	return fn(value, args)
}

func (i *Interpreter) Run(value interface{}) error {
	for _, nf := range i.ast {
		err := i.runFunction(value, nf)
		if err != nil {
			return err
		}
	}
	return nil
}
