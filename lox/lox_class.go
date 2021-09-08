package lox

type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]LoxCallable
}

type LoxInstance struct {
	class  *LoxClass
	fields map[string]interface{}
}

func NewLoxClass(name string, methods map[string]LoxCallable) *LoxClass {
	return &LoxClass{name: name, methods: methods}
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{class: class, fields: map[string]interface{}{}}
}

func (c *LoxClass) Arity() int {
	initializer := c.findMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.(*LoxFunction).Arity()
}

func (c *LoxClass) Call(i *Interpreter, arguments ...interface{}) interface{} {
	instance := NewLoxInstance(c)
	initializer := c.findMethod("init")
	if initializer != nil {
		initializer.(*LoxFunction).bind(instance).(*LoxFunction).Call(i, arguments)
	}
	return instance
}

func (c *LoxClass) findMethod(name string) interface{} {
	method, ok := c.methods[name]
	if ok {
		return method
	}

	if c.superclass != nil {
		return c.superclass.findMethod(name)
	}
	return nil
}

func (i *LoxInstance) Get(name string) interface{} {
	value, ok := i.fields[name]
	if ok {
		return value
	}

	method := i.class.findMethod(name)

	if method != nil {
		method.(*LoxFunction).bind(i)
		return method
	}

	//TODO throw error!
	return nil
}

func (i *LoxInstance) Set(name string, value interface{}) {
	i.fields[name] = value
}
