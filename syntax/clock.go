package syntax

import (
	"golox/references"
	"time"
)

type Clock struct{}

func NewClock() LoxCallable {
	return &Clock{}
}

func (clock *Clock) arity() int {
	return 0
}

func (clock *Clock) call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return time.Now().UnixNano() / int64(time.Millisecond) / int64(time.Second)
}

func (clock *Clock) callableType() references.FunctionType {
	return references.Function
}

func (clock *Clock) String() string {
	return "<native fn>"
}

func (clock *Clock) name() string {
	return "clock"
}
