package lox

type Visitor interface {
	VisitGroupExpr(e *GroupExpr) interface{}
	VisitBinaryExpr(e *BinaryExpr) interface{}
	VisitLogicalExpr(l *LogicalExpr) interface{}
	VisitLiteralExpr(e *LiteralExpr) interface{}
	VisitUnaryExpr(e *UnaryExpr) interface{}
	VisitVariableExpr(e *VariableExpr) interface{}
	VisitAssignExpr(b *AssignExpr) interface{}
	VisitIfStmt(i *IfStmt) interface{}
	VisitBlockStmt(b *BlockStmt) interface{}
	VisitVariableStmt(s *VariableStmt) interface{}
	VisitWhileStmt(i *WhileStmt) interface{}
}
