package lox

import (
	"log"
)

type FunctionType int
type ClassType int

const (
	NONE FunctionType = iota
	FUNCTION
	INITIALIZER
	METHOD
)

const (
	CNONE ClassType = iota
	CCLASS
	SUBCLASS
)

type Scope struct {
	s map[string]bool
}

func NewScope() *Scope {
	return &Scope{s: make(map[string]bool)}
}

func (s *Scope) put(key string, value bool) {
	s.s[key] = value
}

func (s *Scope) get(key string) (bool, bool) {
	val, ok := s.s[key]
	return val, ok
}

func (s *Scope) containsKey(key string) bool {
	_, ok := s.s[key]
	return ok
}

type LoxResolver struct {
	i               *Interpreter
	scopes          []*Scope
	currentFunction FunctionType
	currentClass    ClassType
}

func NewResolver(i *Interpreter) *LoxResolver {
	return &LoxResolver{i: i, currentFunction: NONE, currentClass: CNONE}
}

func (l *LoxResolver) VisitBlockStmt(b *BlockStmt) interface{} {
	l.beginScope()
	l.resolveStatements(b.statements)
	l.endScope()
	return nil
}

func (l *LoxResolver) VisitVariableStmt(s *VariableStmt) interface{} {
	l.declare(s.name)
	if s.initializer != nil {
		l.resolveExpr(s.initializer)
	}
	l.define(s.name)
	return nil
}

func (l *LoxResolver) VisitVariableExpr(e *VariableExpr) interface{} {

	if len(l.scopes) > 0 {
		n := len(l.scopes) - 1
		scope := l.scopes[n]
		varr, ok := scope.s[e.name.Lexeme]
		if !varr || !ok {
			l.error(e.name, "Can't read local variable in its own initializer.")
		}
	}
	l.resolveLocal(e, e.name)
	return nil
}

func (l *LoxResolver) VisitAssignExpr(b *AssignExpr) interface{} {
	l.resolveExpr(b.value)
	l.resolveLocal(b, b.name)
	return nil
}

func (l *LoxResolver) VisitFunctionStmt(f *FunctionStmt) interface{} {
	l.declare(f.name)
	l.define(f.name)
	l.resolveFunction(f, FUNCTION)
	return nil
}

func (l *LoxResolver) VisitExprStmt(e *ExprStmt) interface{} {
	l.resolveExpr(e.expression)
	return nil
}

func (l *LoxResolver) VisitIfStmt(i *IfStmt) interface{} {
	l.resolveExpr(i.condition)
	l.resolveStatement(i.thenBranch)
	if i.elseBranch != nil {
		l.resolveStatement(i.elseBranch)
	}
	return nil
}

func (l *LoxResolver) VisitPrintStmt(p *PrintStmt) interface{} {
	l.resolveExpr(p.expression)
	return nil
}

func (l *LoxResolver) VisitReturnStmt(r *ReturnStmt) interface{} {

	if l.currentFunction == NONE {
		l.error(r.keyword, "Can't return from top-level code.")
	}

	if r.value != nil {
		if l.currentFunction == INITIALIZER {
			l.error(r.keyword, "Can't return a value from an initializer.")
		}
		l.resolveExpr(r.value)
	}
	return nil
}

func (l *LoxResolver) VisitWhileStmt(i *WhileStmt) interface{} {
	l.resolveExpr(i.condition)
	l.resolveStatement(i.body)
	return nil
}

func (l *LoxResolver) VisitBinaryExpr(e *BinaryExpr) interface{} {
	l.resolveExpr(e.left)
	l.resolveExpr(e.right)
	return nil
}

func (l *LoxResolver) VisitCallExpr(c *CallExpr) interface{} {
	l.resolveExpr(c.callee)
	for _, expr := range c.arguments {
		l.resolveExpr(expr)
	}
	return nil
}

func (l *LoxResolver) VisitGroupExpr(e *GroupExpr) interface{} {
	l.resolveExpr(e.expression)
	return nil
}

func (l *LoxResolver) VisitLogicalExpr(le *LogicalExpr) interface{} {
	l.resolveExpr(le.left)
	l.resolveExpr(le.right)
	return nil
}

func (l *LoxResolver) VisitLiteralExpr(e *LiteralExpr) interface{} {
	return nil
}

func (l *LoxResolver) VisitUnaryExpr(e *UnaryExpr) interface{} {
	l.resolveExpr(e.right)
	return nil
}

func (l *LoxResolver) resolveStatements(statements []Stmt) {
	for _, stmt := range statements {
		l.resolveStatement(stmt)
	}
}

func (l *LoxResolver) resolveStatement(stmt Stmt) {
	stmt.Accept(l)
}

func (l *LoxResolver) resolveExpr(expr Expr) {
	expr.Accept(l)
}

func (l *LoxResolver) beginScope() {
	l.scopes = append(l.scopes, NewScope())
}

func (l *LoxResolver) endScope() {
	n := len(l.scopes) - 1
	l.scopes[n] = nil
	l.scopes = l.scopes[:n]
}

func (l *LoxResolver) declare(name Token) {
	if len(l.scopes) == 0 {
		return
	}

	//TOP
	n := len(l.scopes) - 1
	scope := l.scopes[n]

	//if already defined
	if scope.containsKey(name.Lexeme) {
		l.error(name, "Already a variable with this name in this scope.")
	}
	scope.put(name.Lexeme, false)
}

func (l *LoxResolver) define(name Token) {
	if len(l.scopes) == 0 {
		return
	}
	n := len(l.scopes) - 1
	scope := l.scopes[n]
	scope.put(name.Lexeme, true)
}

func (l *LoxResolver) resolveLocal(e Expr, name Token) {
	for i := len(l.scopes) - 1; i >= 0; i-- {
		_, ok := l.scopes[i].s[name.Lexeme]
		if ok {
			l.i.resolve(e, len(l.scopes)-1-i)
		}
	}
}

func (l *LoxResolver) error(name Token, msg string) {
	log.Fatal(name, msg)
}

func (l *LoxResolver) resolveFunction(f *FunctionStmt, ftype FunctionType) {
	enclosingType := l.currentFunction
	l.currentFunction = ftype
	l.beginScope()
	for _, param := range f.params {
		l.declare(param)
		l.define(param)
	}
	l.resolveStatements(f.body)
	l.endScope()
	l.currentFunction = enclosingType
}

func (l *LoxResolver) Resolve(statements []Stmt) {
	l.resolveStatements(statements)
}

func (l *LoxResolver) VisitClassStmt(c *ClassStmt) interface{} {
	enclosingClass := l.currentClass
	l.currentClass = CCLASS
	l.declare(c.name)
	l.define(c.name)

	if c.superclass != nil && c.name.Lexeme == c.superclass.name.Lexeme {
		l.error(c.superclass.name,
			"A class can't inherit from itself.")
	}

	if c.superclass != nil {
		l.currentClass = SUBCLASS
		l.beginScope()
		n := len(l.scopes) - 1
		scope := l.scopes[n]
		scope.put("super", true)
		l.resolveExpr(c.superclass)
	}

	l.beginScope()
	n := len(l.scopes) - 1
	scope := l.scopes[n]
	scope.put("this", true)

	for _, method := range c.methods {
		declaration := METHOD
		if method.(*FunctionStmt).name.Lexeme == "init" {
			declaration = INITIALIZER
		}
		l.resolveFunction(method.(*FunctionStmt), declaration)
	}

	l.endScope()
	if c.superclass != nil {
		l.endScope()
	}
	l.currentClass = enclosingClass
	return nil
}

func (l *LoxResolver) VisitGetExpr(g *GetExpr) interface{} {
	l.resolveExpr(g.object)
	return nil
}

func (l *LoxResolver) VisitSetExpr(s *SetExpr) interface{} {
	l.resolveExpr(s.object)
	l.resolveExpr(s.value)
	return nil
}

func (l *LoxResolver) VisitThisExpr(e *ThisExpr) interface{} {
	if l.currentClass == CNONE {
		l.error(e.keyword, "Can't use 'this' outside of a class.")
	}
	l.resolveLocal(e, e.keyword)
	return nil
}

func (l *LoxResolver) VisitSuperExpr(e *SuperExpr) interface{} {

	if l.currentClass == CNONE {
		l.error(e.keyword, "Can't use 'super' outside of a class.")
	}
	if l.currentClass != SUBCLASS {
		l.error(e.keyword, "Can't use 'super' in a class with no superclass.")
	}
	l.resolveLocal(e, e.keyword)
	return nil
}
