package syntax

import (
	"fmt"
	"golox/scanner"
)

type Environment struct {
	enclosing  *Environment
	exit       bool
	continuing bool
	values     map[string]interface{}
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing:  enclosing,
		exit:       false,
		continuing: false,
		values:     make(map[string]interface{}),
	}
}

func (env *Environment) assign(name *scanner.Token, value interface{}) {
	if _, ok := env.values[name.Lexeme]; ok {
		env.values[name.Lexeme] = value
		return
	}

	if env.enclosing != nil {
		env.enclosing.assign(name, value)
		return
	}

	throwRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
}

func (env *Environment) get(name *scanner.Token) interface{} {
	if value, ok := env.values[name.Lexeme]; ok {
		if value == nil {
			throwRuntimeError(name, fmt.Sprintf("Variable '%s' is uninitialized.", name.Lexeme))
			return nil
		}

		return value
	}

	if env.enclosing != nil {
		return env.enclosing.get(name)
	}

	throwRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
	return nil
}

func (env *Environment) define(name string, value interface{}) {
	env.values[name] = value
}
