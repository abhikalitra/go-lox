package lox

type Visitor interface {
	VisitGroupExpr(e *GroupExpr) interface{}
	VisitBinaryExpr(e *BinaryExpr) interface{}
	VisitLiteralExpr(e *LiteralExpr) interface{}
	VisitUnaryExpr(e *UnaryExpr) interface{}
	VisitVariableExpr(e *VariableExpr) interface{}
	VisitVariableStmt(s *VariableStmt) interface{}
	VisitAssignExpr(b *AssignExpr) interface{}
	VisitBlockStmt(b *BlockStmt) interface{}
}
