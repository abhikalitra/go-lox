package lox

type Visitor interface {
	VisitGroupExpr(e *GroupExpr) interface{}
	VisitBinaryExpr(e *BinaryExpr) interface{}
	VisitLogicalExpr(l *LogicalExpr) interface{}
	VisitLiteralExpr(e *LiteralExpr) interface{}
	VisitUnaryExpr(e *UnaryExpr) interface{}
	VisitVariableExpr(e *VariableExpr) interface{}
	VisitAssignExpr(b *AssignExpr) interface{}
	VisitCallExpr(c *CallExpr) interface{}
	VisitGetExpr(g *GetExpr) interface{}
	VisitThisExpr(t *ThisExpr) interface{}
	VisitIfStmt(i *IfStmt) interface{}
	VisitBlockStmt(b *BlockStmt) interface{}
	VisitVariableStmt(s *VariableStmt) interface{}
	VisitWhileStmt(i *WhileStmt) interface{}
	VisitFunctionStmt(f *FunctionStmt) interface{}
	VisitReturnStmt(r *ReturnStmt) interface{}
	VisitExprStmt(e *ExprStmt) interface{}
	VisitPrintStmt(p *PrintStmt) interface{}
	VisitClassStmt(c *ClassStmt) interface{}
	VisitSetExpr(s *SetExpr) interface{}
	VisitSuperExpr(e *SuperExpr) interface{}
}
