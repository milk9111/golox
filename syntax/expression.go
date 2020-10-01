package syntax

import "golox/scanner"

type Expr interface{
	accept(visitor ExprVisitor) interface{}
}

type ExprVisitor interface {
	visitAssignExpr(expr *Assign) interface{}
	visitBinaryExpr(expr *Binary) interface{}
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


