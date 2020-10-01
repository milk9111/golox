package syntax

import (
	"fmt"
	"golox/references"
	"golox/scanner"
	"strconv"
)

type Interpreter struct {
	env *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: NewEnvironment(nil),
	}
}

func (interpreter *Interpreter) Interpret(statements []Stmt) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	for _, stmt := range statements {
		interpreter.execute(stmt)
	}
}

func (interpreter *Interpreter) execute(stmt Stmt) {
	stmt.accept(interpreter)
}

func (interpreter *Interpreter) visitIftStmt(stmt *Ift) interface{} {
	if isTruthy(interpreter.evaluate(stmt.condition)) {
		interpreter.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		interpreter.execute(stmt.elseBranch)
	}

	return nil
}

func (interpreter *Interpreter) visitBlockStmt(stmt *Block) interface{} {
	interpreter.executeBlock(stmt.statements, NewEnvironment(interpreter.env))
	return nil
}

func (interpreter *Interpreter) executeBlock(statements []Stmt, env *Environment) {
	previous := interpreter.env

	interpreter.env = env
	for _, statement := range statements {
		interpreter.execute(statement)
	}

	interpreter.env = previous
}

func (interpreter *Interpreter) visitAssignExpr(expr *Assign) interface{} {
	value := interpreter.evaluate(expr.value)
	interpreter.env.assign(expr.name, value)
	return value
}

func (interpreter *Interpreter) visitVartStmt(stmt *Vart) interface{} {
	var value interface{}
	if stmt.initializer != nil {
		value = interpreter.evaluate(stmt.initializer)
	}

	interpreter.env.define(stmt.name.Lexeme, value)
	return nil
}

func (interpreter *Interpreter) visitVariableExpr(expr *Variable) interface{} {
	return interpreter.env.get(expr.name)
}

func (interpreter *Interpreter) visitExpressionStmt(stmt *Expression) interface{} {
	interpreter.evaluate(stmt.expression)
	return nil
}

func (interpreter *Interpreter) visitPrintStmt(stmt *Print) interface{} {
	value := interpreter.evaluate(stmt.expression)
	fmt.Println(stringify(value))
	return nil
}

func (interpreter *Interpreter) visitLiteralExpr(expr *Literal) interface{} {
	return expr.value
}

func (interpreter *Interpreter) visitGroupingExpr(expr *Grouping) interface{} {
	return interpreter.evaluate(expr.expression)
}

func (interpreter *Interpreter) visitUnaryExpr(expr *Unary) interface{} {
	right := interpreter.evaluate(expr.right)

	switch expr.operator.Type {
	case references.Bang:
		return !isTruthy(right)
	case references.Minus:
		checkNumberOperand(expr.operator, right)
		return -(right.(float64))
	}

	return nil
}

func (interpreter *Interpreter) visitBinaryExpr(expr *Binary) interface{} {
	left := interpreter.evaluate(expr.left)
	right := interpreter.evaluate(expr.right)

	switch expr.operator.Type {
	case references.Greater:
		checkNumberOperand(expr.operator, left, right)
		return left.(float64) > right.(float64)
	case references.GreaterEqual:
		checkNumberOperand(expr.operator, left, right)
		return left.(float64) >= right.(float64)
	case references.Less:
		checkNumberOperand(expr.operator, left, right)
		return left.(float64) < right.(float64)
	case references.LessEqual:
		checkNumberOperand(expr.operator, left, right)
		return left.(float64) <= right.(float64)
	case references.BangEqual:
		return !isEqual(left, right)
	case references.EqualEqual:
		return isEqual(left, right)
	case references.Minus:
		checkNumberOperand(expr.operator, left, right)
		return left.(float64) - right.(float64)
	case references.Slash:
		checkNumberOperand(expr.operator, left, right)
		if right.(float64) == 0 {
			throwRuntimeError(expr.operator, "Cannot divide by zero.")
		}
		return left.(float64) / right.(float64)
	case references.Star:
		checkNumberOperand(expr.operator, left, right)
		return left.(float64) * right.(float64)
	case references.Plus:
		lFl, lOk := left.(float64)
		rFl, rOk := right.(float64)
		if lOk && rOk {
			return lFl + rFl
		}

		_, lOk = left.(string)
		_, rOk = right.(string)
		if lOk || rOk {
			return fmt.Sprintf("%v%v", left, right)
		}

		throwRuntimeError(expr.operator, "Operands must be two numbers or two strings.")
	}

	return nil
}

func (interpreter *Interpreter) evaluate(expr Expr) interface{} {
	return expr.accept(interpreter)
}

func checkNumberOperand(operator *scanner.Token, operands ...interface{}) {

	good := true
	for _, val := range operands {
		if _, ok := val.(float64); !ok {
			good = false
			break
		}
	}

	if !good {
		s := ""
		if len(operands) > 1 {
			s = "s"
		}
		throwRuntimeError(operator, fmt.Sprintf("Operand%s must be a number.", s))
	}

}

func isTruthy(obj interface{}) bool {
	if obj == nil {
		return false
	}

	if b, ok := obj.(bool); ok {
		return b
	}

	return true
}

func isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil {
		return false
	}

	return a == b
}

func stringify(obj interface{}) string {
	if obj == nil {
		return "nil"
	}

	if f, ok := obj.(float64); ok {
		return strconv.FormatFloat(f, 'f', -1, 64)
	}

	return fmt.Sprintf("%v", obj)
}
