package syntax

import "golox/scanner"

type Stmt interface{
	accept(visitor StmtVisitor) interface{}
}

type StmtVisitor interface {
	visitBlockStmt(stmt *Block) interface{}
	visitExpressionStmt(stmt *Expression) interface{}
	visitIftStmt(stmt *Ift) interface{}
	visitPrintStmt(stmt *Print) interface{}
	visitVartStmt(stmt *Vart) interface{}
}

type Block struct {
	statements []Stmt
}

func NewBlock(statements []Stmt) Stmt {
	return &Block{
		statements: statements,
	}
}

func (block *Block) accept(visitor StmtVisitor) interface{} {
	return visitor.visitBlockStmt(block)
}


type Expression struct {
	expression Expr
}

func NewExpression(expression Expr) Stmt {
	return &Expression{
		expression: expression,
	}
}

func (expression *Expression) accept(visitor StmtVisitor) interface{} {
	return visitor.visitExpressionStmt(expression)
}


type Ift struct {
	condition Expr
	thenBranch Stmt
	elseBranch Stmt
}

func NewIft(condition Expr, thenBranch Stmt, elseBranch Stmt) Stmt {
	return &Ift{
		condition: condition,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}
}

func (ift *Ift) accept(visitor StmtVisitor) interface{} {
	return visitor.visitIftStmt(ift)
}


type Print struct {
	expression Expr
}

func NewPrint(expression Expr) Stmt {
	return &Print{
		expression: expression,
	}
}

func (print *Print) accept(visitor StmtVisitor) interface{} {
	return visitor.visitPrintStmt(print)
}


type Vart struct {
	name *scanner.Token
	initializer Expr
}

func NewVart(name *scanner.Token, initializer Expr) Stmt {
	return &Vart{
		name: name,
		initializer: initializer,
	}
}

func (vart *Vart) accept(visitor StmtVisitor) interface{} {
	return visitor.visitVartStmt(vart)
}


