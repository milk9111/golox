package syntax

import "golox/scanner"

type Expr interface{
	accept(visitor ExprVisitor) interface{}
	String() string}

type ExprVisitor interface {
	visitAssignExpr(expr *Assign) interface{}
	visitBinaryExpr(expr *Binary) interface{}
	visitCallExpr(expr *Call) interface{}
	visitGetExpr(expr *Get) interface{}
	visitSetExpr(expr *Set) interface{}
	visitThisExpr(expr *This) interface{}
	visitGroupingExpr(expr *Grouping) interface{}
	visitLiteralExpr(expr *Literal) interface{}
	visitLogicalExpr(expr *Logical) interface{}
	visitUnaryExpr(expr *Unary) interface{}
	visitVariableExpr(expr *Variable) interface{}
}

type Assign struct {
	name *scanner.Token
	value Expr
}

func NewAssign(name *scanner.Token, value Expr) Expr {
	return &Assign{
		name: name,
		value: value,
	}
}

func (assign *Assign) accept(visitor ExprVisitor) interface{} {
	return visitor.visitAssignExpr(assign)
}

func (assign *Assign) String() string {
	return "Assign"}


type Binary struct {
	left Expr
	operator *scanner.Token
	right Expr
}

func NewBinary(left Expr, operator *scanner.Token, right Expr) Expr {
	return &Binary{
		left: left,
		operator: operator,
		right: right,
	}
}

func (binary *Binary) accept(visitor ExprVisitor) interface{} {
	return visitor.visitBinaryExpr(binary)
}

func (binary *Binary) String() string {
	return "Binary"}


type Call struct {
	callee Expr
	paren *scanner.Token
	arguments []Expr
}

func NewCall(callee Expr, paren *scanner.Token, arguments []Expr) Expr {
	return &Call{
		callee: callee,
		paren: paren,
		arguments: arguments,
	}
}

func (call *Call) accept(visitor ExprVisitor) interface{} {
	return visitor.visitCallExpr(call)
}

func (call *Call) String() string {
	return "Call"}


type Get struct {
	object Expr
	name *scanner.Token
}

func NewGet(object Expr, name *scanner.Token) Expr {
	return &Get{
		object: object,
		name: name,
	}
}

func (get *Get) accept(visitor ExprVisitor) interface{} {
	return visitor.visitGetExpr(get)
}

func (get *Get) String() string {
	return "Get"}


type Set struct {
	object Expr
	name *scanner.Token
	value Expr
}

func NewSet(object Expr, name *scanner.Token, value Expr) Expr {
	return &Set{
		object: object,
		name: name,
		value: value,
	}
}

func (set *Set) accept(visitor ExprVisitor) interface{} {
	return visitor.visitSetExpr(set)
}

func (set *Set) String() string {
	return "Set"}


type This struct {
	keyword *scanner.Token
}

func NewThis(keyword *scanner.Token) Expr {
	return &This{
		keyword: keyword,
	}
}

func (this *This) accept(visitor ExprVisitor) interface{} {
	return visitor.visitThisExpr(this)
}

func (this *This) String() string {
	return "This"}


type Grouping struct {
	expression Expr
}

func NewGrouping(expression Expr) Expr {
	return &Grouping{
		expression: expression,
	}
}

func (grouping *Grouping) accept(visitor ExprVisitor) interface{} {
	return visitor.visitGroupingExpr(grouping)
}

func (grouping *Grouping) String() string {
	return "Grouping"}


type Literal struct {
	value interface{}
}

func NewLiteral(value interface{}) Expr {
	return &Literal{
		value: value,
	}
}

func (literal *Literal) accept(visitor ExprVisitor) interface{} {
	return visitor.visitLiteralExpr(literal)
}

func (literal *Literal) String() string {
	return "Literal"}


type Logical struct {
	left Expr
	operator *scanner.Token
	right Expr
}

func NewLogical(left Expr, operator *scanner.Token, right Expr) Expr {
	return &Logical{
		left: left,
		operator: operator,
		right: right,
	}
}

func (logical *Logical) accept(visitor ExprVisitor) interface{} {
	return visitor.visitLogicalExpr(logical)
}

func (logical *Logical) String() string {
	return "Logical"}


type Unary struct {
	operator *scanner.Token
	right Expr
}

func NewUnary(operator *scanner.Token, right Expr) Expr {
	return &Unary{
		operator: operator,
		right: right,
	}
}

func (unary *Unary) accept(visitor ExprVisitor) interface{} {
	return visitor.visitUnaryExpr(unary)
}

func (unary *Unary) String() string {
	return "Unary"}


type Variable struct {
	name *scanner.Token
}

func NewVariable(name *scanner.Token) Expr {
	return &Variable{
		name: name,
	}
}

func (variable *Variable) accept(visitor ExprVisitor) interface{} {
	return visitor.visitVariableExpr(variable)
}

func (variable *Variable) String() string {
	return "Variable"}


