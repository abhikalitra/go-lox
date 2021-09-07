package lox

type LoxCallable interface {
	Arity() int
	Call(i *Interpreter, arguments ...interface{}) interface{}
}

type LoxFunction struct {
	isInitializer bool
	declaration   *FunctionStmt
	closure       *LoxEnvironment
}

func NewLoxFunction(declaration *FunctionStmt, closure *LoxEnvironment, isInitializer bool) LoxCallable {
	return &LoxFunction{declaration: declaration, closure: closure, isInitializer: isInitializer}
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
	if fn.isInitializer {
		result, _ := fn.closure.GetAt(0, "this")
		return result
	}
	return nil
}

func (fn *LoxFunction) bind(i *LoxInstance) LoxCallable {
	env := NewLoxEnvironmentWithParent(fn.closure)
	env.Define("this", i)
	return NewLoxFunction(fn.declaration, env, fn.isInitializer)
}
