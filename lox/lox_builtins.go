package lox

import "time"

type ClockFunction struct{}

func (c *ClockFunction) Arity() int {
	return 0
}

func (c *ClockFunction) Call(i *Interpreter, arguments ...interface{}) interface{} {
	return time.Now().Second()
}
