package java

import (
	"log"
)

type DeclType string

const (
	PackageDecl DeclType = "PACKAGE"
	ImportDecl  DeclType = "IMPORT"
)

type Stmt interface {
	PrintStmt()
}

type DeclarationStmt struct {
	declType DeclType
	name     []Token
}

type ClassStmt struct {
	name        Token
	superclass  *VariableExpr
	annotations []Stmt
	fields      []Stmt
}

type FieldStmt struct {
	annotation Stmt
	private    bool
	typeName   Token
	varName    Token
}

type AnnotationStmt struct {
	name   Token
	fields []Expr
}

func NewDeclarationStmt(declType DeclType, name []Token) Stmt {
	return &DeclarationStmt{declType: declType, name: name}
}

func NewFieldStmt(annotation Stmt, private bool, typeName Token, varName Token) Stmt {
	return &FieldStmt{
		annotation: annotation,
		private:    private,
		typeName:   typeName,
		varName:    varName,
	}
}

func NewClassStmt(name Token, superclass Expr, annotations []Stmt, fields []Stmt) Stmt {
	return &ClassStmt{
		name:        name,
		superclass:  superclass.(*VariableExpr),
		fields:      fields,
		annotations: annotations,
	}
}

func NewAnnotationStmt(name Token, fields []Expr) Stmt {
	return &AnnotationStmt{
		name:   name,
		fields: fields,
	}
}

func (d *DeclarationStmt) PrintStmt() {
	name := ""
	for _, n := range d.name {
		name += n.Lexeme + "."
	}
	name = name[:len(name)-1]
	log.Printf("%s %s;", d.declType, name)
}

func (f *FieldStmt) PrintStmt() {
	f.annotation.PrintStmt()
	pvt := ""
	if f.private {
		pvt = "PRIVATE"
	}
	log.Printf("%s %s %s", pvt, f.typeName.Lexeme, f.varName.Lexeme)
}

func (c *ClassStmt) PrintStmt() {
	for _, a := range c.annotations {
		a.PrintStmt()
	}
	super := ""
	if c.superclass != nil {
		super = "extends " + c.superclass.String()
	}

	log.Printf("CLASS %s %s { ", c.name.Lexeme, super)
	for _, f := range c.fields {
		f.PrintStmt()
	}
	log.Println("}")
}

func (a *AnnotationStmt) PrintStmt() {
	fields := "("
	for _, f := range a.fields {
		fields += f.String()
	}
	fields += ")"
	log.Printf("@%s%s", a.name.Lexeme, fields)
}
