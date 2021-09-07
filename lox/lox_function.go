package lox

type LoxVariable struct {
}

type LoxCallable interface {
	Arity() int
	Call(i *Interpreter, arguments ...interface{}) interface{}
}

type LoxFunction struct {
	declaration *FunctionStmt
	closure     *LoxEnvironment
}

func NewLoxFunction(declaration *FunctionStmt, closure *LoxEnvironment) LoxCallable {
	return &LoxFunction{declaration: declaration, closure: closure}
}

func (fn *LoxFunction) Arity() int {
	return len(fn.declaration.params)
}

func (fn *LoxFunction) Call(i *Interpreter, arguments ...interface{}) interface{} {
	fnenv := NewLoxEnvironmentWithParent(fn.closure)
	for idx, param := range fn.declaration.params {
		fnenv.Define(param.Lexeme, arguments[idx])
	}
	i.executeBlock(fn.declaration.body, fnenv)
	return nil
}
