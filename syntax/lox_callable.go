package syntax

import "golox/references"

type LoxCallable interface {
	call(interpreter *Interpreter, arguments []interface{}) interface{}
	arity() int
	name() string
	callableType() references.FunctionType
}
