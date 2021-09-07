package lox

type Stmt interface {
	Accept(p Visitor) interface{}
}

type DeclarationStmt struct {
}

type FunctionStmt struct {
	name   Token
	params []Token
	body   []Stmt
}

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

type PrintStmt struct {
	expression Expr
}

type ReturnStmt struct {
	keyword Token
	value   Expr
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

type ExprStmt struct {
	expression Expr
}

type VariableStmt struct {
	name        Token
	initializer Expr
}

type BlockStmt struct {
	statements []Stmt
}

type ClassStmt struct {
	name       Token
	superclass VariableStmt
	methods    []FunctionStmt
}

func NewIfStmt(condition Expr, thenBranch Stmt, elseBranch Stmt) Stmt {
	return &IfStmt{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}
}

func NewWhileStmt(condition Expr, statement Stmt) Stmt {
	return &WhileStmt{
		condition: condition,
		body:      statement,
	}
}

func NewBlockStmt(statements []Stmt) Stmt {
	return &BlockStmt{statements: statements}
}

func NewPrintStmt(expression Expr) Stmt {
	return &PrintStmt{expression: expression}
}

func NewVariableStmt(name Token, initializer Expr) Stmt {
	return &VariableStmt{name: name, initializer: initializer}
}

func NewExprStmt(expr Expr) Stmt {
	return &ExprStmt{expression: expr}
}

func NewFunctionStmt(name Token, params []Token, body []Stmt) Stmt {
	return &FunctionStmt{
		name:   name,
		params: params,
		body:   body,
	}
}

func NewReturnStmt(keyword Token, value Expr) Stmt {
	return &ReturnStmt{
		keyword: keyword,
		value:   value,
	}
}

func (p *PrintStmt) Accept(v Visitor) interface{} {
	v.VisitPrintStmt(p)
	return nil
}

func (s *VariableStmt) Accept(v Visitor) interface{} {
	v.VisitVariableStmt(s)
	return nil
}

func (e *ExprStmt) Accept(v Visitor) interface{} {
	v.VisitExprStmt(e)
	return nil
}

func (b *BlockStmt) Accept(v Visitor) interface{} {
	return v.VisitBlockStmt(b)
}

func (i *IfStmt) Accept(p Visitor) interface{} {
	return p.VisitIfStmt(i)
}

func (w *WhileStmt) Accept(p Visitor) interface{} {
	return p.VisitWhileStmt(w)
}

func (f *FunctionStmt) Accept(p Visitor) interface{} {
	return p.VisitFunctionStmt(f)
}

func (r *ReturnStmt) Accept(p Visitor) interface{} {
	return p.VisitReturnStmt(r)
}
