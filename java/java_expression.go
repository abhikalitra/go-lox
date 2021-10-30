package java

import "fmt"

type Expr interface {
	String() string
}

type VariableExpr struct {
	name Token
}

type LiteralExpr struct {
	value interface{}
}

type DotExpr struct {
	first Token
	next  Token
}

type AssignExpr struct {
	name  Token
	value Expr
}

func NewLiteralExpr(value interface{}) Expr {
	return &LiteralExpr{value: value}
}

func NewVariableExpr(name Token) Expr {
	return &VariableExpr{name: name}
}

func NewAssignExpr(name Token, value Expr) Expr {
	return &AssignExpr{
		name:  name,
		value: value,
	}
}

func NewDotExpr(first Token, next Token) Expr {
	return &DotExpr{first: first, next: next}
}

func (l LiteralExpr) String() string {
	return fmt.Sprintf("%v", l.value)
}

func (v VariableExpr) String() string {
	return v.name.Lexeme
}

func (a AssignExpr) String() string {
	return fmt.Sprintf("%s = %s", a.name.Lexeme, a.value.String())
}

func (d DotExpr) String() string {
	return fmt.Sprintf("%s.%s", d.first.Lexeme, d.next.Lexeme)
}
