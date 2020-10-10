package syntax

type LoxCallable interface {
	call(interpreter *Interpreter, arguments []interface{}) interface{}
	arity() int
	name() string
}
