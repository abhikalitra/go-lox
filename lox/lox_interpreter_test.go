package lox

import (
	"fmt"
	"testing"
)

func TestInterpreter_Print(t *testing.T) {
	expr := "print \"hello\";"
	lexer := NewScanner()
	lexer.Eval(expr)
	parser := NewParser(lexer.Tokens)
	ast := parser.Parse()
	//printer := NewAstPrinter()
	//printer.Print(ast)
	interpreter := NewInterpreter()
	interpreter.Interpret(ast)
	fmt.Println()
}

func TestInterpreter_Multiple(t *testing.T) {
	prog := `var a = 1;
	var b = 2;
	print a + b;`
	lexer := NewScanner()
	lexer.Eval(prog)
	parser := NewParser(lexer.Tokens)
	ast := parser.Parse()
	interpreter := NewInterpreter()
	interpreter.Interpret(ast)
	fmt.Println()
}
