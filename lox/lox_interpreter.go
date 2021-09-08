package lox

import (
	"fmt"
	"log"
)

type Interpreter struct {
	env     *LoxEnvironment
	globals *LoxEnvironment
	locals  map[Expr]int
}

func NewInterpreter() *Interpreter {
	env := NewLoxEnvironment()
	env.Define("clock", &ClockFunction{})

	i := &Interpreter{globals: NewLoxEnvironment()}
	return i
}

func (i *Interpreter) Interpret(statements []Stmt) {

	resolver := NewResolver(i)
	resolver.Resolve(statements)
	for _, stmt := range statements {
		i.execute(stmt)
	}
}

func (i *Interpreter) VisitVariableStmt(s *VariableStmt) interface{} {
	var value interface{}
	if s.initializer != nil {
		value = i.evaluate(s.initializer)
	}
	i.globals.Define(s.name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitExprStmt(e *ExprStmt) interface{} {
	i.evaluate(e.expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(p *PrintStmt) interface{} {
	fmt.Printf("%v", i.evaluate(p.expression))
	return nil
}

func (i *Interpreter) VisitVariableExpr(e *VariableExpr) interface{} {
	return i.lookupVariable(e.name, e)
}

func (i *Interpreter) VisitGroupExpr(e *GroupExpr) interface{} {
	return i.evaluate(e.expression)
}

func (i *Interpreter) VisitBinaryExpr(e *BinaryExpr) interface{} {
	left := i.evaluate(e.left)
	right := i.evaluate(e.right)

	var leftf float64
	var rightf float64

	switch right.(type) {
	case float64:
		rightf = right.(float64)
	default:
		return "**ERR**"
	}

	switch left.(type) {
	case float64:
		leftf = left.(float64)
	default:
		return "**ERR**"
	}

	switch e.operator.TokenType {
	case MINUS:
		return leftf - rightf
	case SLASH:
		return leftf / rightf
	case STAR:
		return leftf * rightf
	case PLUS:
		//TODO string concat
		return leftf + rightf
	case GREATER:
		return leftf > rightf
	case GreaterEqual:
		return leftf >= rightf
	case LESS:
		return leftf < rightf
	case LessEqual:
		return leftf <= rightf
	case BangEqual:
		return !i.isEqual(left, right)
	case EqualEqual:
		return i.isEqual(left, right)
	}
	return "**ERR**"
}

func (i *Interpreter) VisitLiteralExpr(e *LiteralExpr) interface{} {
	return e.value
}

func (i *Interpreter) VisitUnaryExpr(e *UnaryExpr) interface{} {
	right := i.evaluate(e.right)

	switch e.operator.TokenType {
	case BANG:
		return !i.isTrue(right)
	case MINUS:
		switch right.(type) {
		case int:
			return -1 * right.(int)
		case float32:
			return -1 * right.(float32)
		case float64:
			return -1 * right.(float64)
		}
		return "**ERR**"
	}
	return "*ERR*"
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) isTrue(right interface{}) bool {
	if right == nil {
		return false
	}

	switch right.(type) {
	case bool:
		return !right.(bool)
	}
	return true
}

func (i *Interpreter) isEqual(left interface{}, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil && right != nil {
		return false
	}
	if left != nil && right == nil {
		return false
	}
	//TODO implement!
	return true
}

func (i *Interpreter) execute(stmt Stmt) {
	stmt.Accept(i)
}

//TODO cannot assign to unknown variable!!
func (i *Interpreter) VisitAssignExpr(e *AssignExpr) interface{} {

	value := i.evaluate(e.value)
	dist, ok := i.locals[e]
	if ok {
		i.env.AssignAt(dist, e.name, value)
	} else {
		i.globals.Assign(e.name.Lexeme, value)
	}

	return value
}

func (i *Interpreter) VisitBlockStmt(b *BlockStmt) interface{} {
	i.executeBlock(b.statements, NewLoxEnvironmentWithParent(i.globals))
	return nil
}

func (i *Interpreter) executeBlock(statements []Stmt, env *LoxEnvironment) {
	prev := i.globals
	i.globals = env
	for _, stmt := range statements {
		i.execute(stmt)
	}
	i.globals = prev
}

func (i *Interpreter) VisitIfStmt(ifstmt *IfStmt) interface{} {
	if i.isTrue(i.evaluate(ifstmt.condition)) {
		i.execute(ifstmt.thenBranch)
	} else if ifstmt.elseBranch != nil {
		i.execute(ifstmt.elseBranch)
	}
	return nil
}

func (i *Interpreter) VisitLogicalExpr(e *LogicalExpr) interface{} {
	left := i.evaluate(e.left)

	switch e.operator.TokenType {
	case OR:
		if i.isTrue(left) {
			return left
		}
	case AND:
		if !i.isTrue(left) {
			return left
		}
	}

	return i.evaluate(e.right)
}

func (i *Interpreter) VisitWhileStmt(w *WhileStmt) interface{} {
	for i.isTrue(i.evaluate(w.condition)) {
		i.execute(w.body)
	}
	return nil
}

func (i *Interpreter) VisitCallExpr(c *CallExpr) interface{} {
	callee := i.evaluate(c.callee)

	var arguments []interface{}
	for _, arg := range c.arguments {
		arguments = append(arguments, i.evaluate(arg))
	}

	//TODO call the function
	fn, ok := callee.(LoxCallable)

	if !ok {
		i.error("Can only call functions and classes.")
	}

	if len(arguments) != fn.Arity() {
		i.error(fmt.Sprintf("Expected %d but got %d arguments.", fn.Arity, len(arguments)))
	}

	fn.Call(i, arguments)
	return nil
}

func (i *Interpreter) VisitFunctionStmt(f *FunctionStmt) interface{} {
	fn := NewLoxFunction(f, i.globals, false)
	i.globals.Define(f.name.Lexeme, fn)
	return nil
}

func (i *Interpreter) error(msg string) {
	log.Fatal(msg)
}

func (i *Interpreter) VisitReturnStmt(r *ReturnStmt) interface{} {
	//var value interface{}
	//if r.value != nil {
	//	value = i.evaluate(r.value)
	//}
	//TODO we need to escape to call fn somehow
	//in java its done with throw, dont know how to implement it
	//cleanly in golang :(
	return nil
}

func (i *Interpreter) resolve(e Expr, depth int) {
	i.locals[e] = depth
}

func (i *Interpreter) lookupVariable(name Token, e Expr) interface{} {
	dist, ok := i.locals[e]
	if ok {
		varr, _ := i.env.GetAt(dist, name.Lexeme)
		return varr
	} else {
		varr, _ := i.globals.Get(name.Lexeme)
		return varr
	}
}

func (i *Interpreter) VisitClassStmt(c *ClassStmt) interface{} {

	if c.superclass != nil {
		super := i.evaluate(c.superclass)
		_, ok := super.(*LoxClass)
		if !ok {
			i.error(c.superclass.name.Lexeme + "<Superclass must be a class.")
		}
	}
	i.env.Define(c.name.Lexeme, nil)

	if c.superclass != nil {
		i.env = NewLoxEnvironmentWithParent(i.env)
		i.env.Define("super", c.superclass)
	}

	var methods map[string]LoxCallable
	for _, method := range c.methods {
		function := NewLoxFunction(method.(*FunctionStmt), i.env, method.(*FunctionStmt).name.Lexeme == "init")
		methods[method.(*FunctionStmt).name.Lexeme] = function
	}
	class := NewLoxClass(c.name.Lexeme, methods)

	if c.superclass != nil {
		i.env = i.env.parent
	}

	i.env.Assign(c.name.Lexeme, class)
	return nil
}

func (i *Interpreter) VisitGetExpr(g *GetExpr) interface{} {
	object := i.evaluate(g.object)
	instance, ok := object.(LoxInstance)

	if ok {
		instance.Get(g.name.Lexeme)
	} else {
		i.error(g.name.Lexeme + ":Only instances have properties.")
	}
	return nil
}

func (i *Interpreter) VisitSetExpr(s *SetExpr) interface{} {
	object := i.evaluate(s.object)
	instance, ok := object.(LoxInstance)

	if ok {
		value := i.evaluate(s.value)
		instance.Set(s.name.Lexeme, value)
	} else {
		i.error(s.name.Lexeme + ":Only instances have properties.")
	}
	return nil
}

func (i *Interpreter) VisitThisExpr(t *ThisExpr) interface{} {
	return i.lookupVariable(t.keyword, t)
}

func (i *Interpreter) VisitSuperExpr(s *SuperExpr) interface{} {
	distance := i.locals[s]
	superclass, _ := i.env.GetAt(distance, "super")
	instance, _ := i.env.GetAt(distance-1, "this")
	method := superclass.(*LoxClass).findMethod(s.method.Lexeme)

	if method == nil {
		i.error("Undefined property '" + s.method.Lexeme + "'.")
	}

	return method.(*LoxFunction).bind(instance.(*LoxInstance))
}
