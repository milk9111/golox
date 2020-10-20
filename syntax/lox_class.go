package syntax

import (
	"golox/references"
	"golox/scanner"
)

type LoxClass struct {
	className  string
	superclass *LoxClass
	methods    map[string]*LoxFunction
	fields     map[string]interface{}
}

func NewLoxClass(name string, superclass *LoxClass, methods map[string]*LoxFunction, fields map[string]interface{}) *LoxClass {
	return &LoxClass{
		className:  name,
		superclass: superclass,
		methods:    methods,
		fields:     fields,
	}
}

func (class *LoxClass) findMethod(name string) *LoxFunction {
	if method, ok := class.methods[name]; ok {
		return method
	}

	if class.superclass != nil {
		method := class.superclass.findMethod(name)
		if method == nil || method.isStatic {
			return nil
		}

		return method
	}

	return nil
}

func (class *LoxClass) call(interpreter *Interpreter, arguments []interface{}) interface{} {
	instance := NewLoxInstance(class)

	initializer := class.findMethod("init")
	if initializer != nil {
		initializer.bind(instance).call(interpreter, arguments)
	}

	return instance
}

func (class *LoxClass) getStaticMethod(name *scanner.Token) *LoxFunction {
	method := class.findMethod(name.Lexeme)
	if method == nil || !method.isStatic {
		return nil
	}

	return method
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
