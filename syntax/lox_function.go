package syntax

import "fmt"

type LoxFunction struct {
	declaration *Function
	closure     *Environment
}

func NewLoxFunction(declaration *Function, closure *Environment) LoxCallable {
	return &LoxFunction{
		declaration: declaration,
		closure:     closure,
	}
}

func (fun *LoxFunction) name() string {
	return fun.declaration.name.Lexeme
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
				resp = r
			}
		}()

		interpreter.executeBlock(fun.declaration.body, env, nil)
	}()

	interpreter.env = previous
	return resp
}

func (fun *LoxFunction) arity() int {
	return len(fun.declaration.params)
}

func (fun *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", fun.declaration.name.Lexeme)
}
