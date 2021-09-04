package lox

import (
	"testing"
)

func TestParser(t *testing.T) {
	expr := "45.65"
	lexer := NewScanner()
	lexer.Eval(expr)
	parser := NewParser(lexer.Tokens)
	ast := parser.Parse()
	printer := NewAstPrinter()
	printer.Print(ast)
}
