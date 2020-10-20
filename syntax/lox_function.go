package syntax

import (
	"fmt"
	"golox/references"
)

type LoxFunction struct {
	declaration   *Function
	closure       *Environment
	isInitializer bool
	isStatic      bool
}

func NewLoxFunction(declaration *Function, closure *Environment, isInit bool, isStatic bool) *LoxFunction {
	return &LoxFunction{
		declaration:   declaration,
		closure:       closure,
		isInitializer: isInit,
		isStatic:      isStatic,
	}
}

func (fun *LoxFunction) name() string {
	return fun.declaration.name.Lexeme
}

func (fun *LoxFunction) bind(instance *LoxInstance) *LoxFunction {
	env := NewEnvironment(fun.closure)
	env.define("this", instance)
	return NewLoxFunction(fun.declaration, env, fun.isInitializer, fun.isStatic)
}

func (fun *LoxFunction) call(interpreter *Interpreter, arguments []interface{}) interface{} {
	env := NewEnvironment(fun.closure)
	for i := 0; i < len(fun.declaration.params); i++ {
		env.define(fun.declaration.params[i].Lexeme, arguments[i])
	}

	var resp interface{}

	previous := interpreter.env
	interpreter.env = env
	func() {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					throwRuntimeError(fun.declaration.name, err.Error())
				}

				if fun.isInitializer {
					resp = fun.closure.getAt(0, "this")
				} else {
					resp = r
				}
			}
		}()

		interpreter.executeBlock(fun.declaration.body, env, nil)
	}()

	interpreter.env = previous

	if fun.isInitializer {
		resp = fun.closure.getAt(0, "this")
	}

	return resp
}

func (fun *LoxFunction) arity() int {
	return len(fun.declaration.params)
}

func (fun *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", fun.declaration.name.Lexeme)
}

func (fun *LoxFunction) callableType() references.FunctionType {
	return references.Function
}
