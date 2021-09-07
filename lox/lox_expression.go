package lox

type Expr interface {
	Accept(p Visitor) interface{}
}

type BinaryExpr struct {
	left     Expr
	operator Token
	right    Expr
}

type AssignExpr struct {
	name  Token
	value Expr
}

type CallExpr struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

type GetExpr struct {
	object Expr
	name   Token
}

type GroupExpr struct {
	expression Expr
}

type LiteralExpr struct {
	value interface{}
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

func NewCallExpr(callee Expr, paren Token, arguments []Expr) Expr {
	return &CallExpr{
		callee:    callee,
		paren:     paren,
		arguments: arguments,
	}
}

func NewGetExpr(object Expr, name Token) Expr {
	return &GetExpr{
		object: object,
		name:   name,
	}
}

func NewGroupExpr(expression Expr) Expr {
	return &GroupExpr{expression: expression}
}

func NewBinaryExpr(left Expr, operator Token, right Expr) Expr {
	return &BinaryExpr{
		left:     left,
		operator: operator,
		right:    right,
	}
}

func NewAssignExpr(name Token, value Expr) Expr {
	return &AssignExpr{
		name:  name,
		value: value,
	}
}

func NewUnaryExpr(operator Token, right Expr) Expr {
	return &UnaryExpr{
		operator: operator,
		right:    right,
	}
}

func NewLiteralExpr(value interface{}) Expr {
	return &LiteralExpr{value: value}
}

func NewVariableExpr(name Token) Expr {
	return &VariableExpr{name: name}
}

func NewLogicalExpr(left Expr, operator Token, right Expr) Expr {
	return &LogicalExpr{
		left:     left,
		operator: operator,
		right:    right,
	}
}

func NewSetExpr(object Expr, name Token, value Expr) Expr {
	return &SetExpr{
		object: object,
		name:   name,
		value:  value,
	}
}

func NewThisExpr(keyword Token) Expr {
	return &ThisExpr{keyword: keyword}
}

func (e *UnaryExpr) Accept(p Visitor) interface{} {
	return p.VisitUnaryExpr(e)
}

func (e *VariableExpr) Accept(p Visitor) interface{} {
	return p.VisitVariableExpr(e)
}

func (e *BinaryExpr) Accept(p Visitor) interface{} {
	return p.VisitBinaryExpr(e)
}

func (b *AssignExpr) Accept(v Visitor) interface{} {
	return v.VisitAssignExpr(b)
}

func (e *LiteralExpr) Accept(p Visitor) interface{} {
	return p.VisitLiteralExpr(e)
}

func (b *GroupExpr) Accept(p Visitor) interface{} {
	return p.VisitGroupExpr(b)
}

func (l *LogicalExpr) Accept(p Visitor) interface{} {
	return p.VisitLogicalExpr(l)
}

func (b *GetExpr) Accept(p Visitor) interface{} {
	return p.VisitGetExpr(b)
}

func (s *SetExpr) Accept(p Visitor) interface{} {
	return p.VisitSetExpr(s)
}

func (t *ThisExpr) Accept(p Visitor) interface{} {
	return p.VisitThisExpr(t)
}

func (b *CallExpr) Accept(p Visitor) interface{} {
	panic("not implemented")
}
