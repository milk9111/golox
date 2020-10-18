package syntax

import (
	"fmt"
	"golox/scanner"
)

var level = -1

type Environment struct {
	enclosing  *Environment
	exit       bool
	continuing bool
	values     map[string]interface{}
	name       string
}

func NewEnvironment(enclosing *Environment) *Environment {
	level++
	return &Environment{
		enclosing:  enclosing,
		exit:       false,
		continuing: false,
		values:     make(map[string]interface{}),
		name:       fmt.Sprintf("env: %d", level),
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

func (env *Environment) print() {
	fmt.Printf("%s\n", env.name)
	for key, val := range env.values {
		fmt.Printf("\t%s: %v\n", key, val)
	}
}

func (env *Environment) assignAt(depth int, token *scanner.Token, value interface{}) {
	env.ancestor(depth).values[token.Lexeme] = value
}

func (env *Environment) getAt(depth int, name string) interface{} {
	val, ok := env.ancestor(depth).values[name]
	if !ok {
		return nil
	}

	return val
}

func (env *Environment) ancestor(depth int) *Environment {
	currEnv := env
	for i := 0; i < depth; i++ {
		currEnv = currEnv.enclosing
	}

	return currEnv
}
