package lox

import (
	"fmt"
	"testing"
)

func TestPrintAst(t *testing.T) {

	e := BinaryExpr{
		left: &UnaryExpr{
			operator: Token{
				TokenType: MINUS,
				Lexeme:    "-",
				Literal:   nil,
				Line:      1,
			},
			right: LiteralExpr{value: 123},
		},
		operator: Token{
			TokenType: STAR,
			Lexeme:    "*",
			Literal:   nil,
			Line:      1,
		},
		right: GroupExpr{
			expression: LiteralExpr{value: 45.67},
		},
	}

	printer := NewAstPrinter()
	want := "(* (- 123) (group 45.67))"
	got := printer.PrintExpr(e)
	fmt.Println(got)
	if want != got {
		t.Fatalf("AST: %s != %s", want, got)
	}

	//(* (- 123) (group 45.67))

	/*
		Expr expression = new Expr.Binary(
			new Expr.Unary(
			new Token(TokenType.MINUS, "-", null, 1),
			new Expr.Literal(123)),
		new Token(TokenType.STAR, "*", null, 1),
			new Expr.Grouping(
			new Expr.Literal(45.67)));

		System.out.println(new AstPrinter().print(expression));

	*/
}
