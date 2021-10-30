package java

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

//Parse ::program        → declaration* EOF
func (p *Parser) Parse() []Stmt {
	var statements []Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

//declaration    → package Decl | import Decl | classDecl
func (p *Parser) declaration() Stmt {

	if p.match(PACKAGE) {
		return p.packageDeclaration()
	}

	if p.match(IMPORT) {
		return p.importDeclaration()
	}

	return p.classDeclaration()
	//TODO implement error handling and return error
	//synchronise it here if something goes wrong
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
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
		p.report(token.Line, " at '"+token.Lexeme+"'", msg)
	}
}

func (p *Parser) report(line int, s string, msg string) {
	log.Fatalf("%d %s : %s", line, s, msg)
}

//packageDecl      → "package" IDENTIFIER ( "." IDENTIFIER )*;
func (p *Parser) packageDeclaration() Stmt {

	var packagename []Token
	name := p.consume(IDENTIFIER, "Expected package name name")
	packagename = append(packagename, name)

	for {
		if p.match(DOT) {
			name = p.consume(IDENTIFIER, "??")
			packagename = append(packagename, name)
		}
		if p.match(SEMICOLON) {
			break
		}
		if p.match(EOF) {
			p.error(name, "parse error, cannot locate ; after package declaration")
		}
	}
	//p.consume(SEMICOLON, "Expected ; after value.")
	return NewDeclarationStmt(PackageDecl, packagename)
}

//importDecl      → "package" IDENTIFIER ( "." IDENTIFIER )*;
func (p *Parser) importDeclaration() Stmt {
	var packagename []Token
	name := p.consume(IDENTIFIER, "Expected import name")
	packagename = append(packagename, name)

	for {
		if p.match(DOT) {
			name = p.consume(IDENTIFIER, "??")
			packagename = append(packagename, name)
		}
		if p.match(SEMICOLON) {
			break
		}
		if p.match(EOF) {
			p.error(name, "parse error, cannot locate ; after import declaration")
		}
	}
	//p.consume(SEMICOLON, "Expected ; after value.")
	return NewDeclarationStmt(ImportDecl, packagename)
}

//classDecl      → " (annotation)* class" IDENTIFIER ( "extends" IDENTIFIER )? "{" fields* "}" ;
func (p *Parser) classDeclaration() Stmt {

	var annotations []Stmt
	for {
		if p.match(ATRATE) {
			annotations = append(annotations, p.annotation())
		} else {
			break
		}
	}

	if p.match(PUBLIC, PRIVATE) {

	}

	p.consume(CLASS, "Expected class")
	name := p.consume(IDENTIFIER, "Expected class name")
	var superclass Expr

	if p.match(EXTENDS) {
		p.consume(IDENTIFIER, "Expect superclass name.")
		superclass = NewVariableExpr(p.previous())
	}

	p.consume(LeftBrace, "Expected '{' before class body.")
	var fields []Stmt
	for !p.check(RightBrace) && !p.isAtEnd() {
		fields = append(fields, p.field())
	}

	p.consume(RightBrace, "Expect '}' after class body.")
	return NewClassStmt(name, superclass, annotations, fields)
}

//field      → "@" Annotation ("public" | "private")? type var;
func (p *Parser) field() Stmt {
	var annotation Stmt
	if p.match(ATRATE) {
		annotation = p.annotation()
	}

	private := false
	if p.match(PRIVATE) {
		private = true
	}
	typeName := p.consume(IDENTIFIER, "Expected type name.")
	varName := p.consume(IDENTIFIER, "Expected var name.")
	p.consume(SEMICOLON, "Expected ; after field decl")
	return NewFieldStmt(annotation, private, typeName, varName)
}

//annotation      → "Type ( ("name" = "value")* )"
func (p *Parser) annotation() Stmt {
	typeName := p.consume(IDENTIFIER, "Expected annotation type")
	var fields []Expr

	if p.match(LeftParen) {
		for {
			if p.match(RightParen) || p.match(EOF) {
				break
			}
			name := p.consume(IDENTIFIER, "Expecting identifier")
			var value Expr
			var field Expr
			if p.match(EQUAL) {
				value = p.primary()
				field = NewAssignExpr(name, value)
			}

			if p.match(DOT) {
				def := NewToken(DEFAULT, "default", nil, name.Line)
				value = NewDotExpr(name, p.consume(IDENTIFIER, "Expecting identifier"))
				field = NewAssignExpr(def, value)
			}

			//field := NewAssignExpr(name, value)
			fields = append(fields, field)

			if p.match(COMMA) {
				continue
			} else {
				p.consume(RightParen, "Expecting )")
				break
			}
		}
	}
	return NewAnnotationStmt(typeName, fields)
}

//primary        →  "true" | "false" | "nil" | NUMBER | STRING | "this" | IDENTIFIER | "(" expression ") | "super" "." IDENTIFIER ;;
func (p *Parser) primary() Expr {

	if p.match(TRUE) {
		return NewLiteralExpr(true)
	}

	if p.match(FALSE) {
		return NewLiteralExpr(false)
	}

	if p.match(NULL) {
		return NewLiteralExpr(nil)
	}

	if p.match(NUMBER, STRING) {
		return NewLiteralExpr(p.previous().Literal)
	}

	if p.match(IDENTIFIER) {
		return NewVariableExpr(p.previous())
	}

	p.error(p.peek(), "Expected expression.")
	return nil
}
