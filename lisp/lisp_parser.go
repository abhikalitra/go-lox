package lisp

import "log"

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

func (p *Parser) Parse() []Expr {
	var expressions []Expr
	for !p.isAtEnd() {
		expressions = append(expressions, p.expression())
	}
	return expressions
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) expression() Expr {

	if p.match(LeftParen) {
		expr := p.expression()

		var list []Expr

		for {
			if !p.check(RightParen) {
				list = append(list, p.expression())
			} else {
				p.consume(RightParen, "Expecting )")
				return NewListExpr(expr, list)
			}
		}
	}

	return p.atom()
}

func (p *Parser) atom() Expr {

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

	//if p.match(IDENTIFIER) {
	//	return NewVariableExpr(p.previous())
	//}

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
		p.report(token.Line, " near '"+token.Lexeme+"'", msg)
	}
}

func (p *Parser) report(line int, s string, msg string) {
	log.Fatalf("at line %d %s : %s", line, s, msg)
}
