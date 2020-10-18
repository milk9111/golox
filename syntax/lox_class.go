package syntax

import (
	"golox/references"
)

type LoxClass struct {
	className string
	methods   map[string]*LoxFunction
}

func NewLoxClass(name string, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		className: name,
		methods:   methods,
	}
}

func (class *LoxClass) findMethod(name string) *LoxFunction {
	return class.methods[name]
}

func (class *LoxClass) call(interpreter *Interpreter, arguments []interface{}) interface{} {
	instance := NewLoxInstance(class)

	initializer := class.findMethod("init")
	if initializer != nil {
		initializer.bind(instance).call(interpreter, arguments)
	}

	return instance
}

func (class *LoxClass) arity() int {
	init := class.findMethod("init")
	if init == nil {
		return 0
	}

	return init.arity()
}

func (class *LoxClass) name() string {
	return class.className
}

func (class *LoxClass) callableType() references.FunctionType {
	return references.Klass
}
