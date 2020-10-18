package syntax

import "golox/scanner"

type Stmt interface{
	accept(visitor StmtVisitor) interface{}
	String() string}

type StmtVisitor interface {
	visitBlockStmt(stmt *Block) interface{}
	visitExpressionStmt(stmt *Expression) interface{}
	visitFunctionStmt(stmt *Function) interface{}
	visitIfCmdStmt(stmt *IfCmd) interface{}
	visitPrintStmt(stmt *Print) interface{}
	visitReturnCmdStmt(stmt *ReturnCmd) interface{}
	visitVarCmdStmt(stmt *VarCmd) interface{}
	visitWhileLoopStmt(stmt *WhileLoop) interface{}
	visitBreakCmdStmt(stmt *BreakCmd) interface{}
	visitContinueCmdStmt(stmt *ContinueCmd) interface{}
	visitClassStmt(stmt *Class) interface{}
}

type Block struct {
	statements []Stmt
	isLoopIncrementer bool
}

func NewBlock(statements []Stmt, isLoopIncrementer bool) Stmt {
	return &Block{
		statements: statements,
		isLoopIncrementer: isLoopIncrementer,
	}
}

func (block *Block) accept(visitor StmtVisitor) interface{} {
	return visitor.visitBlockStmt(block)
}

func (block *Block) String() string {
	return "Block"}


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

func (expression *Expression) String() string {
	return "Expression"}


type Function struct {
	name *scanner.Token
	params []*scanner.Token
	body []Stmt
}

func NewFunction(name *scanner.Token, params []*scanner.Token, body []Stmt) Stmt {
	return &Function{
		name: name,
		params: params,
		body: body,
	}
}

func (function *Function) accept(visitor StmtVisitor) interface{} {
	return visitor.visitFunctionStmt(function)
}

func (function *Function) String() string {
	return "Function"}


type IfCmd struct {
	condition Expr
	thenBranch Stmt
	elseBranch Stmt
}

func NewIfCmd(condition Expr, thenBranch Stmt, elseBranch Stmt) Stmt {
	return &IfCmd{
		condition: condition,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}
}

func (ifcmd *IfCmd) accept(visitor StmtVisitor) interface{} {
	return visitor.visitIfCmdStmt(ifcmd)
}

func (ifcmd *IfCmd) String() string {
	return "IfCmd"}


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

func (print *Print) String() string {
	return "Print"}


type ReturnCmd struct {
	keyword *scanner.Token
	value Expr
}

func NewReturnCmd(keyword *scanner.Token, value Expr) Stmt {
	return &ReturnCmd{
		keyword: keyword,
		value: value,
	}
}

func (returncmd *ReturnCmd) accept(visitor StmtVisitor) interface{} {
	return visitor.visitReturnCmdStmt(returncmd)
}

func (returncmd *ReturnCmd) String() string {
	return "ReturnCmd"}


type VarCmd struct {
	name *scanner.Token
	initializer Expr
}

func NewVarCmd(name *scanner.Token, initializer Expr) Stmt {
	return &VarCmd{
		name: name,
		initializer: initializer,
	}
}

func (varcmd *VarCmd) accept(visitor StmtVisitor) interface{} {
	return visitor.visitVarCmdStmt(varcmd)
}

func (varcmd *VarCmd) String() string {
	return "VarCmd"}


type WhileLoop struct {
	condition Expr
	body Stmt
}

func NewWhileLoop(condition Expr, body Stmt) Stmt {
	return &WhileLoop{
		condition: condition,
		body: body,
	}
}

func (whileloop *WhileLoop) accept(visitor StmtVisitor) interface{} {
	return visitor.visitWhileLoopStmt(whileloop)
}

func (whileloop *WhileLoop) String() string {
	return "WhileLoop"}


type BreakCmd struct {
	keyword *scanner.Token
	envDepth int
}

func NewBreakCmd(keyword *scanner.Token, envDepth int) Stmt {
	return &BreakCmd{
		keyword: keyword,
		envDepth: envDepth,
	}
}

func (breakcmd *BreakCmd) accept(visitor StmtVisitor) interface{} {
	return visitor.visitBreakCmdStmt(breakcmd)
}

func (breakcmd *BreakCmd) String() string {
	return "BreakCmd"}


type ContinueCmd struct {
	keyword *scanner.Token
	envDepth int
}

func NewContinueCmd(keyword *scanner.Token, envDepth int) Stmt {
	return &ContinueCmd{
		keyword: keyword,
		envDepth: envDepth,
	}
}

func (continuecmd *ContinueCmd) accept(visitor StmtVisitor) interface{} {
	return visitor.visitContinueCmdStmt(continuecmd)
}

func (continuecmd *ContinueCmd) String() string {
	return "ContinueCmd"}


type Class struct {
	name *scanner.Token
	methods []*Function
}

func NewClass(name *scanner.Token, methods []*Function) Stmt {
	return &Class{
		name: name,
		methods: methods,
	}
}

func (class *Class) accept(visitor StmtVisitor) interface{} {
	return visitor.visitClassStmt(class)
}

func (class *Class) String() string {
	return "Class"}


