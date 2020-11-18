package parser

import (
	"fmt"
	"github.com/madelyne-io/madelyne/matcher/vm/ast"
	"github.com/madelyne-io/madelyne/matcher/vm/lexer"
	"github.com/madelyne-io/madelyne/matcher/vm/token"
)

var (
	ErrBadToken     = fmt.Errorf("Bad token type encounter")
	ErrIllegalToken = fmt.Errorf("Illegal token encournter")
)

type Parser struct {
	lexer   *lexer.Lexer
	current token.Token
	next    token.Token
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:   lexer,
		current: token.Token{},
		next:    token.Token{},
	}

	parser.moveToNextToken()
	parser.moveToNextToken()

	return parser
}

func (parser *Parser) moveToNextToken() {
	parser.current = parser.next
	parser.next = parser.lexer.GetNextToken()
}

func (parser *Parser) moveToNextTokenIfTypeIs(t token.TokenType) error {
	if !parser.isNextTokenA(t) {
		return fmt.Errorf("1 %w : expected %s, got %s", ErrBadToken, t, parser.next.Type)
	}
	parser.moveToNextToken()
	return nil
}

func (parser *Parser) isCurrentTokenA(t token.TokenType) bool {
	return parser.current.Type == t
}

func (parser *Parser) isNextTokenA(t token.TokenType) bool {
	return parser.next.Type == t
}

func (parser *Parser) Parse() ([]*ast.NodeFunction, error) {
	program := []*ast.NodeFunction{}
	for parser.current.Type != token.EOF {
		if parser.current.Type == token.ILLEGAL {
			return nil, fmt.Errorf("2 %w : %s", ErrIllegalToken, parser.current.Literal)
		}
		if !parser.isCurrentTokenA(token.DOT) {
			return nil, fmt.Errorf("3 %w : expected %s, got %s", ErrBadToken, token.DOT, parser.current.Type)
		}
		parser.moveToNextToken()
		fn, err := parser.parseFunction()
		if err != nil {
			return nil, err
		}
		if fn != nil {
			program = append(program, fn)
		}
		parser.moveToNextToken()
	}
	return program, nil
}

func (parser *Parser) parseFunction() (*ast.NodeFunction, error) {
	if !parser.isCurrentTokenA(token.IDENTIFIER) {
		return nil, fmt.Errorf("4 %w : expected %s, got %s", ErrBadToken, token.IDENTIFIER, parser.current.Type)
	}

	fn := &ast.NodeFunction{
		Function: parser.current,
	}
	if err := parser.moveToNextTokenIfTypeIs(token.LEFT_PARENTHESIS); err != nil {
		return nil, err
	}

	for parser.next.Type != token.RIGHT_PARENTHESIS {
		switch parser.next.Type {
		case token.IDENTIFIER:
			{
				parser.moveToNextToken()
				f, err := parser.parseFunction()
				if err != nil {
					return nil, err
				}
				fn.Arguments = append(fn.Arguments, f)
			}
		case token.NUMBER:
			{
				parser.moveToNextToken()
				v := &ast.NodeValue{parser.current}
				fn.Arguments = append(fn.Arguments, v)
			}
		case token.STRING:
			{
				parser.moveToNextToken()
				v := &ast.NodeValue{parser.current}
				fn.Arguments = append(fn.Arguments, v)
			}
		default:
			return nil, fmt.Errorf("5 %w : expected %s,%s or %s, got %s",
				ErrBadToken, token.IDENTIFIER, token.NUMBER, token.STRING, parser.next.Type)
		}

		switch parser.next.Type {
		case token.COMMA:
			parser.moveToNextToken()
		case token.RIGHT_PARENTHESIS:
		default:
			return nil, fmt.Errorf("6 %w : expected %s or %s, got %s",
				ErrBadToken, token.COMMA, token.RIGHT_PARENTHESIS, parser.next.Type)
		}
	}
	parser.moveToNextToken()
	return fn, nil
}
