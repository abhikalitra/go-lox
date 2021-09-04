package lox

type LoxEnvironment struct {
	values map[string]interface{}
	parent *LoxEnvironment
}

func NewLoxEnvironment() *LoxEnvironment {
	env := &LoxEnvironment{}
	env.values = make(map[string]interface{})
	return env
}

func NewLoxEnvironmentWithParent(parent *LoxEnvironment) *LoxEnvironment {
	env := &LoxEnvironment{}
	env.values = make(map[string]interface{})
	env.parent = parent
	return env
}

func (e *LoxEnvironment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *LoxEnvironment) Get(name string) (interface{}, bool) {
	value, ok := e.values[name]

	if !ok && e.parent != nil {
		return e.parent.Get(name)
	}
	return value, ok
}

func (e *LoxEnvironment) Assign(name string, value interface{}) {
	value, ok := e.values[name]

	if !ok && e.parent != nil {
		e.parent.Assign(name, value)
	} else {
		e.values[name] = value
	}

}
