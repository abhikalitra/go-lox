package lisp

import "fmt"

type Atom interface {
}

type Expr interface {
	PrintExpr()
}

type LiteralExpr struct {
	value interface{}
}

type ListExpr struct {
	expr Expr
	list []Expr
}

func NewLiteralExpr(value interface{}) Expr {
	return &LiteralExpr{value: value}
}

func NewListExpr(expr Expr, list []Expr) Expr {
	return &ListExpr{expr: expr, list: list}
}

func (l *LiteralExpr) PrintExpr() {
	fmt.Printf("%v", l.value)
}

func (l *ListExpr) PrintExpr() {
	fmt.Print("(")
	l.expr.PrintExpr()
	fmt.Print(" ")
	for _, expr := range l.list {
		expr.PrintExpr()
		fmt.Print(" ")
	}
	fmt.Println(")")
}
