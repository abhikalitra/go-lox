package lox

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (p *AstPrinter) Print(statements []Stmt) {
	for _, stmt := range statements {
		p.print(stmt)
	}
}

func (p *AstPrinter) PrintExpr(e Expr) string {
	return fmt.Sprintf("%v", e.Accept(p))
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) interface{} {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		result := expr.Accept(p)
		builder.WriteString(fmt.Sprintf("%v", result))
	}
	builder.WriteString(")")
	return builder.String()
}

func (p *AstPrinter) VisitGroupExpr(e *GroupExpr) interface{} {
	return p.parenthesize("group", e.expression)
}

func (p *AstPrinter) VisitBinaryExpr(e *BinaryExpr) interface{} {
	return p.parenthesize(e.operator.Lexeme, e.left, e.right)
}

func (p *AstPrinter) VisitUnaryExpr(e *UnaryExpr) interface{} {
	return p.parenthesize(e.operator.Lexeme, e.right)
}

func (p *AstPrinter) VisitLiteralExpr(e *LiteralExpr) interface{} {
	if e.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", e.value)
}

func (p *AstPrinter) print(stmt Stmt) {
	stmt.Accept(p)
}

func (p *AstPrinter) VisitVariableStmt(s *VariableStmt) interface{} {
	return p.parenthesize(s.name.Lexeme, s.initializer)
}

func (p *AstPrinter) VisitVariableExpr(e *VariableExpr) interface{} {
	return p.parenthesize(e.name.Lexeme)
}
