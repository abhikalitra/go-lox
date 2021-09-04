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

//statement      → exprStmt | forStmt | ifStmt | printStmt | whileStmt | block
func (p *Parser) statement() Stmt {

	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(LeftBrace) {
		return p.block()
	}

	return p.expressionStatement()
}

//whileStmt      → "while" "(" expression ")" statement ;
func (p *Parser) whileStatement() Stmt {
	p.consume(LeftParen, "Expected '(' after 'while'.)")
	condition := p.expression()
	p.consume(RightParen, "Expected ')' after 'while' condition.)")
	return NewWhileStmt(condition, p.statement())
}

//forStmt        → "for" "(" ( varDecl | exprStmt | ";" )
//                 expression? ";"
//                 expression? ")" statement ;
func (p *Parser) forStatement() Stmt {
	p.consume(LeftParen, "Expected '(' after 'for'.")
	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expr
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after loop condition.")

	var increment Expr
	if !p.check(RightParen) {
		increment = p.expression()
	}
	p.consume(RightParen, "Expect ')' after for clauses.")
	body := p.statement()

	//desugar to while loop
	if increment != nil {
		body = NewBlockStmt([]Stmt{body, NewExprStmt(increment)})
	}

	if condition == nil {
		condition = NewLiteralExpr(true)
	}

	body = NewWhileStmt(condition, body)

	if initializer != nil {
		body = NewBlockStmt([]Stmt{initializer, body})
	}

	return body
}

//ifStmt         → "if" "(" expression ")" statement ( "else" statement )? ;
func (p *Parser) ifStatement() Stmt {
	p.consume(LeftParen, "Expected '(' after 'if'.")
	condition := p.expression()
	p.consume(RightParen, "Expected ')' after if condition.")
	thenBranch := p.statement()
	var elseBranch Stmt

	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return NewIfStmt(condition, thenBranch, elseBranch)
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

//assignment     → IDENTIFIER "=" assignment | logic_or ;
func (p *Parser) assignment() Expr {

	expr := p.or()

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

//logic_or       → logic_and ( "or" logic_and )* ;
func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = NewLogicalExpr(expr, operator, right)
	}
	return expr
}

//logic_and      → equality ( "and" equality )* ;
func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = NewLogicalExpr(expr, operator, right)
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
