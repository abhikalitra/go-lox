package lox

import (
	"log"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

//Parse ::program        → declaration* EOF
func (p *Parser) Parse() []Stmt {
	var statements []Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

//declaration    → varDecl | statement
func (p *Parser) declaration() Stmt {
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
	//TODO implement error handling and return error
	//synchronise it here if something goes wrong
}

//varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
func (p *Parser) varDeclaration() Stmt {

	name := p.consume(IDENTIFIER, "Expect variable name.")

	var expr Expr
	if p.match(EQUAL) {
		expr = p.expression()
	}
	p.consume(SEMICOLON, "Expected ; after value.")
	return NewVariableStmt(name, expr)
}

//statement      → exprStmt | printStmt | block
func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(LeftBrace) {
		return p.block()
	}

	return p.expressionStatement()
}

//block          → "{" declaration* "}" ;
func (p *Parser) block() Stmt {
	var statements []Stmt
	for !p.check(RightBrace) && !p.isAtEnd() {
		stmt := p.declaration()
		statements = append(statements, stmt)
	}
	p.consume(RightBrace, "Expect '}' after block.")
	return NewBlockStmt(statements)
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expected ; after value.")
	return NewPrintStmt(value)
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expected ; after value.")
	return NewExprStmt(expr)
}

//expression     → assignment ;
func (p *Parser) expression() Expr {
	return p.assignment()
}

//assignment     → IDENTIFIER "=" assignment | equality ;
func (p *Parser) assignment() Expr {
	expr := p.equality()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		valexpr, ok := expr.(*VariableExpr)

		if ok {
			return NewAssignExpr(valexpr.name, value)
		}

		p.error(equals, "Invalid assignment target.")
	}
	return expr
}

//equality       → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr
}

//comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(GREATER, GreaterEqual, LESS, LessEqual) {
		operator := p.previous()
		right := p.term()
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr
}

//term           → factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr
}

//factor         → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr
}

//unary          → ( "!" | "-" ) unary | primary ;
func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return NewUnaryExpr(operator, right)
	}

	return p.primary()
}

//primary        →  "true" | "false" | "nil" | NUMBER | STRING | "(" expression ") | IDENTIFIER" ;
func (p *Parser) primary() Expr {

	if p.match(TRUE) {
		return NewLiteralExpr(true)
	}

	if p.match(FALSE) {
		return NewLiteralExpr(false)
	}

	if p.match(NIL) {
		return NewLiteralExpr(nil)
	}

	if p.match(NUMBER, STRING) {
		return NewLiteralExpr(p.previous().Literal)
	}

	if p.match(IDENTIFIER) {
		return NewVariableExpr(p.previous())
	}

	if p.match(LeftParen) {
		expr := p.expression()
		p.consume(RightParen, "Expecting ) after expression")
		return NewGroupExpr(expr)
	}

	p.error(p.peek(), "Expected expression.")
	return nil
}

func (p *Parser) match(types ...TokenType) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(typ TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == typ
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType TokenType, msg string) Token {
	if p.check(tokenType) {
		return p.advance()
	}
	p.error(p.peek(), msg)
	return NewErrorToken()
}

func (p *Parser) error(token Token, msg string) {
	if token.TokenType == EOF {
		p.report(token.Line, " at end", msg)
	} else {
		p.report(token.Line, " at '"+token.Lexeme+"'", msg)
	}
}

func (p *Parser) report(line int, s string, msg string) {
	log.Fatalf("%d %s : %s", line, s, msg)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == SEMICOLON {
			return
		}
		switch p.peek().TokenType {
		case CLASS:
			return
		case FUN:
			return
		case VAR:
			return
		case FOR:
			return
		case IF:
			return
		case WHILE:
			return
		case PRINT:
			return
		case RETURN:
			return
		}
		p.advance()
	}
}
