package syntax

import (
	"fmt"
	"golox/references"
	"golox/scanner"
)

type LoxInstance struct {
	class  *LoxClass
	fields map[string]interface{}
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		class:  class,
		fields: class.fields,
	}
}

func (instance *LoxInstance) call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return nil
}

func (instance *LoxInstance) arity() int {
	return 0
}

func (instance *LoxInstance) name() string {
	return fmt.Sprintf("%s instance", instance.class.name())
}

func (instance *LoxInstance) callableType() references.FunctionType {
	return references.Klass
}

func (instance *LoxInstance) getMethod(name *scanner.Token) interface{} {
	if method := instance.class.findMethod(name.Lexeme); method != nil && !method.isStatic {
		return method.bind(instance)
	}

	throwRuntimeError(name, fmt.Sprintf("Undefined method '%s'.", name.Lexeme))
	return nil
}

func (instance *LoxInstance) getField(name *scanner.Token) interface{} {
	if val, ok := instance.fields[name.Lexeme]; ok {
		return val
	}

	throwRuntimeError(name, fmt.Sprintf("Undefined field '%s'.", name.Lexeme))
	return nil
}

func (instance *LoxInstance) set(name *scanner.Token, value interface{}) {
	instance.fields[name.Lexeme] = value
}
