package lox

type Expr interface {
	Accept(p Visitor) interface{}
}

type BinaryExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func NewBinaryExpr(left Expr, operator Token, right Expr) Expr {
	return BinaryExpr{
		left:     left,
		operator: operator,
		right:    right,
	}
}

func (e BinaryExpr) Accept(p Visitor) interface{} {
	return p.VisitBinaryExpr(&e)
}

type AssignExpr struct {
	name  Token
	value Expr
}

func NewAssignExpr(name Token, value Expr) Expr {
	return &AssignExpr{
		name:  name,
		value: value,
	}
}

func (b *AssignExpr) Accept(v Visitor) interface{} {
	return v.VisitAssignExpr(b)
}

type CallExpr struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

func NewCallExpr(callee Expr, paren Token, arguments []Expr) Expr {
	return CallExpr{
		callee:    callee,
		paren:     paren,
		arguments: arguments,
	}
}

func (b CallExpr) Accept(p Visitor) interface{} {
	return "TODO"
}

type GetExpr struct {
	object Expr
	name   Token
}

func NewGetExpr(object Expr, name Token) Expr {
	return GetExpr{
		object: object,
		name:   name,
	}
}

func (b GetExpr) Accept(p Visitor) interface{} {
	return "TODO"
}

type GroupExpr struct {
	expression Expr
}

func NewGroupExpr(expression Expr) Expr {
	return GroupExpr{expression: expression}
}

func (b GroupExpr) Accept(p Visitor) interface{} {
	return p.VisitGroupExpr(&b)
}

type LiteralExpr struct {
	value interface{}
}

func NewLiteralExpr(value interface{}) Expr {
	return LiteralExpr{value: value}
}

func (e LiteralExpr) Accept(p Visitor) interface{} {
	return p.VisitLiteralExpr(&e)
}

type VariableExpr struct {
	name Token
}

type LogicalExpr struct {
	left     Expr
	operator Token
	right    Expr
}

type SetExpr struct {
	object Expr
	name   Token
	value  Expr
}

type SuperExpr struct {
	keyword Token
	method  Token
}

type ThisExpr struct {
	keyword Token
}

type UnaryExpr struct {
	operator Token
	right    Expr
}

func NewUnaryExpr(operator Token, right Expr) Expr {
	return &UnaryExpr{
		operator: operator,
		right:    right,
	}
}

func NewVariableExpr(name Token) Expr {
	return &VariableExpr{name: name}
}

func (e *UnaryExpr) Accept(p Visitor) interface{} {
	return p.VisitUnaryExpr(e)
}

func (e *VariableExpr) Accept(p Visitor) interface{} {
	return p.VisitVariableExpr(e)
}
