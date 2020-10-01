package syntax

import (
	"fmt"
	"golox/loxerror"
	"golox/references"
	"golox/scanner"
)

type AstParser struct {
	Tokens  []*scanner.Token
	Current int
}

func NewAstParser(tokens []*scanner.Token) *AstParser {
	return &AstParser{
		Tokens:  tokens,
		Current: 0,
	}
}

func (parser *AstParser) Parse() []Stmt {
	var statements []Stmt
	for !parser.isAtEnd() {
		statements = append(statements, parser.declaration())
	}

	return statements
}

func (parser *AstParser) declaration() Stmt {
	defer func() {
		if r := recover(); r != nil {
			parser.synchronize()
		}
	}()

	if parser.match(references.Var) {
		return parser.varDeclaration()
	}

	return parser.statement()
}

func (parser *AstParser) varDeclaration() Stmt {
	name := parser.consume(references.Identifier, "Expect variable name.")

	var initializer Expr
	if parser.match(references.Equal) {
		initializer = parser.expression()
	}

	parser.consume(references.Semicolon, "Expect ';' after variable declaration.")
	return NewVart(name, initializer)
}

func (parser *AstParser) statement() Stmt {
	if parser.match(references.If) {
		return parser.ifStatement()
	}

	if parser.match(references.Print) {
		return parser.printStatement()
	}

	if parser.match(references.LeftBrace) {
		return NewBlock(parser.block())
	}

	return parser.expressionStatement()
}

func (parser *AstParser) ifStatement() Stmt {
	parser.consume(references.LeftParen, "Expect '(' after if.")
	condition := parser.expression()
	parser.consume(references.RightParen, "Expect ')' after if condition.")

	thenStatement := parser.statement()
	var elseStatement Stmt
	if parser.match(references.Else) {
		elseStatement = parser.statement()
	}

	return NewIft(condition, thenStatement, elseStatement)
}

func (parser *AstParser) block() []Stmt {
	var statements []Stmt
	for !parser.check(references.RightBrace) && !parser.isAtEnd() {
		statements = append(statements, parser.declaration())
	}

	parser.consume(references.RightBrace, "Expect '}' after block.")
	return statements
}

func (parser *AstParser) printStatement() Stmt {
	value := parser.expression()
	parser.consume(references.Semicolon, "Expect ';' after value.")

	return NewPrint(value)
}

func (parser *AstParser) expressionStatement() Stmt {
	value := parser.expression()
	parser.consume(references.Semicolon, "Expect ';' after value.")

	return NewExpression(value)
}

func (parser *AstParser) expression() Expr {
	return parser.assignment()
}

func (parser *AstParser) assignment() Expr {
	expr := parser.equality()

	if parser.match(references.Equal) {
		equals := parser.previous()
		value := parser.assignment()

		if v, ok := expr.(*Variable); ok {
			return NewAssign(v.name, value)
		}

		throwError(equals, "Invalid assignment target.")
	}

	return expr
}

func (parser *AstParser) equality() Expr {
	expr := parser.comparison()

	for parser.match(references.BangEqual, references.EqualEqual) {
		operator := parser.previous()
		right := parser.comparison()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (parser *AstParser) comparison() Expr {
	expr := parser.addition()

	for parser.match(references.Greater, references.GreaterEqual, references.Less, references.LessEqual) {
		operator := parser.previous()
		right := parser.addition()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (parser *AstParser) addition() Expr {
	expr := parser.multiplication()

	for parser.match(references.Minus, references.Plus) {
		operator := parser.previous()
		right := parser.multiplication()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (parser *AstParser) multiplication() Expr {
	expr := parser.unary()

	for parser.match(references.Slash, references.Star) {
		operator := parser.previous()
		right := parser.unary()

		val := parser.previous().Literal
		if val != nil {
			if f, ok := val.(float64); operator.Type == references.Slash && ok && f == 0 {
				throwError(operator, "Cannot divide by zero.")
			}
		}

		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (parser *AstParser) unary() Expr {
	if parser.match(references.Bang, references.Minus) {
		operator := parser.previous()
		right := parser.unary()
		return NewUnary(operator, right)
	}

	return parser.primary()
}

func (parser *AstParser) primary() Expr {
	if parser.match(references.False) {
		return NewLiteral(false)
	}

	if parser.match(references.True) {
		return NewLiteral(true)
	}

	if parser.match(references.Nil) {
		return NewLiteral(nil)
	}

	if parser.match(references.Number, references.String) {
		return NewLiteral(parser.previous().Literal)
	}

	if parser.match(references.Identifier) {
		return NewVariable(parser.previous())
	}

	if parser.match(references.LeftParen) {
		expr := parser.expression()
		parser.consume(references.RightParen, "Expected ')' after expression.")
		return NewGrouping(expr)
	}

	throwError(parser.peek(), "Expect expression.")
	return nil
}

func (parser *AstParser) consume(tokenType references.TokenType, message string) *scanner.Token {
	if parser.check(tokenType) {
		return parser.advance()
	}

	throwError(parser.peek(), message)
	return nil
}

func (parser *AstParser) synchronize() {
	parser.advance()

	for !parser.isAtEnd() {
		if parser.previous().Type == references.Semicolon {
			return
		}

		switch parser.peek().Type {
		case references.Class:
			return
		case references.Fun:
			return
		case references.Var:
			return
		case references.For:
			return
		case references.If:
			return
		case references.While:
			return
		case references.Print:
			return
		case references.Return:
			return
		}

		parser.advance()
	}
}

func (parser *AstParser) match(types ...references.TokenType) bool {
	for _, t := range types {
		if parser.check(t) {
			parser.advance()
			return true
		}
	}

	return false
}

func (parser *AstParser) check(t references.TokenType) bool {
	if parser.isAtEnd() {
		return false
	}

	return parser.peek().Type == t
}

func (parser *AstParser) advance() *scanner.Token {
	if !parser.isAtEnd() {
		parser.Current++
	}

	return parser.previous()
}

func (parser *AstParser) isAtEnd() bool {
	return parser.peek().Type == references.EOF
}

func (parser *AstParser) peek() *scanner.Token {
	return parser.Tokens[parser.Current]
}

func (parser *AstParser) previous() *scanner.Token {
	return parser.Tokens[parser.Current-1]
}

func throwError(token *scanner.Token, message string) {
	loxerror.TokenError(token.Type, token.Line, token.Lexeme, message)

	panic(fmt.Errorf(message))
}

func throwRuntimeError(token *scanner.Token, message string) {
	loxerror.TokenRuntimeError(token.Type, token.Line, token.Lexeme, message, true)

	panic(fmt.Errorf(message))
}
