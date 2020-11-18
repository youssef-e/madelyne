package ast

import (
	"bytes"
	"github.com/madelyne-io/madelyne/matcher/vm/token"
)

const (
	incrementIndentation = "	"
)

type Node interface {
	Token() token.Token
	Dump(indent string) string
}

type NodeValue struct {
	Value token.Token
}

func (nv *NodeValue) Token() token.Token { return nv.Value }
func (nv *NodeValue) Dump(indent string) string {
	return indent + "Value " + nv.Token().Literal + "\n"
}

type NodeFunction struct {
	Function  token.Token
	Arguments []Node
}

func (nf *NodeFunction) Token() token.Token { return nf.Function }
func (nf *NodeFunction) Dump(indent string) string {
	var out bytes.Buffer
	out.WriteString(indent + "Function " + nf.Token().Literal + "\n")
	for _, a := range nf.Arguments {
		out.WriteString(a.Dump(indent + incrementIndentation))
	}
	return out.String()
}
