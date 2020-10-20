package syntax

import (
	"fmt"
	"golox/loxerror"
	"golox/references"
	"golox/scanner"
	"math"
	"strconv"
	"strings"
)

var globals = NewEnvironment(nil)
var locals = map[Expr]*int{}

type Interpreter struct {
	env  *Environment
	prev *Environment
}

func NewInterpreter() *Interpreter {
	globals.define("clock", NewClock())

	return &Interpreter{
		env:  globals,
		prev: nil,
	}
}

func (interpreter *Interpreter) Interpret(statements []Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if !loxerror.HadRuntimeError() {
				if err, ok := r.(error); ok {
					fmt.Println(err.Error())
				} else {
					fmt.Println("Runtime error occurred.")
				}
			}
		}
	}()

	for _, stmt := range statements {
		interpreter.execute(stmt)
	}
}

func (interpreter *Interpreter) execute(stmt Stmt) {
	stmt.accept(interpreter)
}

func (interpreter *Interpreter) resolve(expr Expr, depth *int) {
	locals[expr] = depth
}

func (interpreter *Interpreter) visitReturnCmdStmt(stmt *ReturnCmd) interface{} {
	var value interface{}
	if stmt.value != nil {
		value = interpreter.evaluate(stmt.value)
	}

	throwReturn(value)
	return nil
}

func (interpreter *Interpreter) visitFunctionStmt(stmt *Function) interface{} {
	function := NewLoxFunction(stmt, interpreter.env, false, false)
	globals.define(stmt.name.Lexeme, function)

	return nil
}

func (interpreter *Interpreter) visitContinueCmdStmt(continueCmd *ContinueCmd) interface{} {
	env := interpreter.env
	for depth := continueCmd.envDepth; depth > 0; depth-- {
		env.continuing = true
		env = env.enclosing
	}

	return nil
}

func (interpreter *Interpreter) visitBreakCmdStmt(breakCmd *BreakCmd) interface{} {
	env := interpreter.env
	for depth := breakCmd.envDepth; depth > 0; depth-- {
		env.exit = true
		env = env.enclosing
	}

	return nil
}

func (interpreter *Interpreter) visitWhileLoopStmt(whileLoop *WhileLoop) interface{} {
	for isTruthy(interpreter.evaluate(whileLoop.condition)) {
		interpreter.execute(whileLoop.body)

		if interpreter.env.continuing {
			interpreter.env = interpreter.prev
			continue
		}

		if interpreter.env.exit {
			interpreter.env = interpreter.prev
			break
		}
	}

	return nil
}

func (interpreter *Interpreter) visitLogicalExpr(expr *Logical) interface{} {
	left := interpreter.evaluate(expr.left)

	if expr.operator.Type == references.Or {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}

	return interpreter.evaluate(expr.right)
}

func (interpreter *Interpreter) visitIfCmdStmt(stmt *IfCmd) interface{} {
	if isTruthy(interpreter.evaluate(stmt.condition)) {
		interpreter.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		interpreter.execute(stmt.elseBranch)
	}

	return nil
}

func (interpreter *Interpreter) visitGetMethodExpr(expr *GetMethod) interface{} {
	object := interpreter.evaluate(expr.object)
	if val, ok := object.(*LoxInstance); ok {
		return val.getMethod(expr.name)
	}

	if val, ok := object.(*LoxClass); ok {
		return val.getStaticMethod(expr.name)
	}

	throwRuntimeError(expr.name, "Only instances have properties.")
	return nil
}

func (interpreter *Interpreter) visitGetFieldExpr(expr *GetField) interface{} {
	object := interpreter.evaluate(expr.object)
	if val, ok := object.(*LoxInstance); ok {
		return val.getField(expr.name)
	}

	throwRuntimeError(expr.name, "Only instances have properties.")
	return nil
}

func (interpreter *Interpreter) visitSetExpr(expr *Set) interface{} {
	object := interpreter.evaluate(expr.object)

	val, ok := object.(*LoxInstance)
	if !ok {
		throwRuntimeError(expr.name, "Only instances have fields.")
	}

	value := interpreter.evaluate(expr.value)
	val.set(expr.name, value)

	return value
}

func (interpreter *Interpreter) visitBlockStmt(stmt *Block) interface{} {
	interpreter.executeBlock(stmt.statements, NewEnvironment(interpreter.env), stmt)
	return nil
}

func (interpreter *Interpreter) executeBlock(statements []Stmt, env *Environment, block *Block) {
	previous := interpreter.env

	interpreter.env = env
	for _, statement := range statements {
		interpreter.execute(statement)
		if interpreter.env.continuing && block != nil && block.isLoopIncrementer {
			interpreter.env.continuing = false
			continue
		}

		if interpreter.env.exit || interpreter.env.continuing {
			interpreter.prev = previous
			return
		}
	}

	interpreter.env = previous
}

func (interpreter *Interpreter) visitAssignExpr(expr *Assign) interface{} {
	value := interpreter.evaluate(expr.value)

	distance, ok := locals[expr]
	if !ok {
		interpreter.env.assign(expr.name, value)
		return value
	}

	if distance != nil {
		interpreter.env.assignAt(*distance, expr.name, value)
	} else {
		globals.assign(expr.name, value)
	}

	return value
}

func (interpreter *Interpreter) visitVarCmdStmt(stmt *VarCmd) interface{} {
	var value interface{}
	if stmt.initializer != nil {
		value = interpreter.evaluate(stmt.initializer)
	}

	interpreter.env.define(stmt.name.Lexeme, value)
	return nil
}

func (interpreter *Interpreter) visitClassStmt(stmt *Class) interface{} {
	var superclass *LoxClass
	if stmt.superclass != nil {
		s := interpreter.evaluate(stmt.superclass)
		ok := false
		if superclass, ok = s.(*LoxClass); !ok {
			throwRuntimeError(stmt.superclass.name, "Superclass must be a class.")
		}
	}

	interpreter.env.define(stmt.name.Lexeme, nil)

	if stmt.superclass != nil {
		interpreter.env = NewEnvironment(interpreter.env)
		interpreter.env.define("super", interpreter.env.get(stmt.superclass.name))
	}

	methods := make(map[string]*LoxFunction)
	for _, method := range stmt.methods {
		methods[method.name.Lexeme] = NewLoxFunction(method, interpreter.env, method.name.Lexeme == "init" && !method.isStatic, method.isStatic)
	}

	fields := make(map[string]interface{})
	for _, field := range stmt.fields {
		var value interface{}
		if field.initializer != nil {
			value = interpreter.evaluate(field.initializer)
		}

		fields[field.name.Lexeme] = value
	}

	class := NewLoxClass(stmt.name.Lexeme, superclass, methods, fields)

	if stmt.superclass != nil {
		interpreter.env = interpreter.env.enclosing
	}

	interpreter.env.assign(stmt.name, class)
	return nil
}

func (interpreter *Interpreter) visitVariableExpr(expr *Variable) interface{} {
	return interpreter.lookupVariable(expr.name, expr)
}

func (interpreter *Interpreter) lookupVariable(name *scanner.Token, expr Expr) interface{} {
	distance, ok := locals[expr]
	if !ok {
		return interpreter.env.get(name)
	}

	if distance != nil {
		return interpreter.env.getAt(*distance, name.Lexeme)
	}

	return globals.get(name)
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
	case references.Modulo:
		checkNumberOperand(expr.operator, left, right)
		return math.Mod(left.(float64), right.(float64))
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

func (interpreter *Interpreter) visitSuperExpr(expr *Super) interface{} {
	distance := locals[expr]
	superclass := interpreter.env.getAt(*distance, "super").(*LoxClass)
	object := interpreter.env.getAt(*distance-1, "this").(*LoxInstance)

	method := superclass.findMethod(expr.method.Lexeme)
	if method == nil {
		throwRuntimeError(expr.method, fmt.Sprintf("Undefined property '%s'.", expr.method.Lexeme))
	}

	return method.bind(object)
}

func (interpreter *Interpreter) visitThisExpr(expr *This) interface{} {
	return interpreter.lookupVariable(expr.keyword, expr)
}

func (interpreter *Interpreter) visitCallExpr(expr *Call) interface{} {
	callee := interpreter.evaluate(expr.callee)
	if v, ok := callee.(*LoxFunction); ok && v == nil {
		throwRuntimeError(expr.paren, "Could not find function or method.")
	}

	var arguments []interface{}
	for _, arg := range expr.arguments {
		arguments = append(arguments, interpreter.evaluate(arg))
	}

	if _, ok := callee.(LoxCallable); !ok {
		throwRuntimeError(expr.paren, fmt.Sprintf("Can only call functions and classes but tried to call '%v'.", callee))
	}

	function := callee.(LoxCallable)
	if len(arguments) != function.arity() {
		throwRuntimeError(expr.paren, fmt.Sprintf("Expected %d arguments but got %d for %s '%s'.", function.arity(), len(arguments), strings.ToLower(references.GetFunctionTypeName(function.callableType())), function.name()))
	}

	return function.call(interpreter, arguments)
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

	if val, ok := obj.(LoxCallable); ok {
		return val.name()
	}

	return fmt.Sprintf("%v", obj)
}
