package lox

import "fmt"

type Interpreter struct {
	env *LoxEnvironment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{env: NewLoxEnvironment()}
}

func (i *Interpreter) Interpret(statements []Stmt) {
	for _, stmt := range statements {
		i.execute(stmt)
	}
}

func (i *Interpreter) VisitVariableStmt(s *VariableStmt) interface{} {
	var value interface{}
	if s.initializer != nil {
		value = i.evaluate(s.initializer)
	}
	i.env.Define(s.name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitExprStmt(e *ExprStmt) interface{} {
	i.evaluate(e.expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(e *PrintStmt) interface{} {
	result := i.evaluate(e.expression)
	fmt.Printf("%v", result)
	return nil
}

func (i *Interpreter) VisitVariableExpr(e *VariableExpr) interface{} {
	value, _ := i.env.Get(e.name.Lexeme)
	return value
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
	i.env.Assign(e.name.Lexeme, value)
	return value
}

func (i *Interpreter) VisitBlockStmt(b *BlockStmt) interface{} {
	i.executeBlock(b.statements, NewLoxEnvironmentWithParent(i.env))
	return nil
}

func (i *Interpreter) executeBlock(statements []Stmt, env *LoxEnvironment) {
	prev := i.env
	i.env = env
	for _, stmt := range statements {
		i.execute(stmt)
	}
	i.env = prev
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
