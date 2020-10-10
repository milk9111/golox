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

	if parser.match(references.Fun) {
		return parser.function("function")
	}

	if parser.match(references.Var) {
		return parser.varDeclaration()
	}

	return parser.statement()
}

func (parser *AstParser) function(kind string) Stmt {
	name := parser.consume(references.Identifier, fmt.Sprintf("Expect %s name.", kind))
	parser.consume(references.LeftParen, fmt.Sprintf("Expect '(' after %s name", kind))

	var params []*scanner.Token
	if !parser.check(references.RightParen) {
		for ok := true; ok; ok = parser.match(references.Comma) {
			if len(params) > 255 {
				throwError(parser.peek(), "Can't have more than 255 parameters.")
			}

			params = append(params, parser.consume(references.Identifier, "Expect parameter name."))
		}
	}

	parser.consume(references.RightParen, "Expect ')' after parameters.")
	parser.consume(references.LeftBrace, fmt.Sprintf("Expect '{' before %s body.", kind))
	body := parser.block()

	return NewFunction(name, params, body)
}

func (parser *AstParser) varDeclaration() Stmt {
	name := parser.consume(references.Identifier, "Expect variable name.")

	var initializer Expr
	if parser.match(references.Equal) {
		initializer = parser.expression()
	}

	parser.consume(references.Semicolon, "Expect ';' after variable declaration.")
	return NewVarCmd(name, initializer)
}

func (parser *AstParser) statement() Stmt {
	if parser.match(references.For) {
		return parser.forStatement()
	}

	if parser.match(references.If) {
		return parser.ifStatement()
	}

	if parser.match(references.Print) {
		return parser.printStatement()
	}

	if parser.match(references.Return) {
		return parser.returnStatement()
	}

	if parser.match(references.While) {
		return parser.whileStatement()
	}

	if parser.match(references.LeftBrace) {
		return NewBlock(parser.block(), false)
	}

	if parser.match(references.Break) {
		return parser.breakStatement()
	}

	if parser.match(references.Continue) {
		return parser.continueStatement()
	}

	return parser.expressionStatement()
}

func (parser *AstParser) returnStatement() Stmt {
	keyword := parser.previous()

	var value Expr
	if !parser.check(references.Semicolon) {
		value = parser.expression()
	}

	parser.consume(references.Semicolon, "Expect ';' after return value.")
	return NewReturnCmd(keyword, value)
}

func (parser *AstParser) continueStatement() Stmt {
	depth, found := parser.calculateDepth()

	if !found {
		throwError(parser.previous(), "Expect 'continue' in a loop.")
	}

	parser.consume(references.Semicolon, "Expect ';' after continue.")
	return NewContinueCmd(depth)
}

func (parser *AstParser) breakStatement() Stmt {
	depth, found := parser.calculateDepth()

	if !found {
		throwError(parser.previous(), "Expect 'break' in a loop.")
	}

	parser.consume(references.Semicolon, "Expect ';' after break.")
	return NewBreakCmd(depth)
}

func (parser *AstParser) forStatement() Stmt {
	parser.consume(references.LeftParen, "Expect '(' after for.")

	var initializer Stmt
	if parser.match(references.Semicolon) {
		initializer = nil
	} else if parser.match(references.Var) {
		initializer = parser.varDeclaration()
	} else {
		initializer = parser.expressionStatement()
	}

	var conditional Expr
	if !parser.check(references.Semicolon) {
		conditional = parser.expression()
	}
	parser.consume(references.Semicolon, "Expect ';' after for loop condition.")

	var increment Expr
	if !parser.check(references.RightParen) {
		increment = parser.expression()
	}
	parser.consume(references.RightParen, "Expect ')' after for loop clauses.")

	body := parser.statement()

	if increment != nil {
		body = NewBlock([]Stmt{body, NewExpression(increment)}, true)
	}

	if conditional == nil {
		conditional = NewLiteral(true)
	}
	body = NewWhileLoop(conditional, body)

	if initializer != nil {
		body = NewBlock([]Stmt{initializer, body}, false)
	}

	return body
}

func (parser *AstParser) whileStatement() Stmt {
	parser.consume(references.LeftParen, "Expect '(' after while.")
	condition := parser.expression()
	parser.consume(references.RightParen, "Expect ')' after while condition.")

	body := parser.statement()

	return NewWhileLoop(condition, body)
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

	return NewIfCmd(condition, thenStatement, elseStatement)
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
	expr := parser.or()

	// TODO - Add in ++ and -- here
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

func (parser *AstParser) or() Expr {
	expr := parser.and()

	for parser.match(references.Or) {
		operator := parser.previous()
		right := parser.and()
		expr = NewLogical(expr, operator, right)
	}

	return expr
}

func (parser *AstParser) and() Expr {
	expr := parser.equality()

	for parser.match(references.And) {
		operator := parser.previous()
		right := parser.equality()
		expr = NewLogical(expr, operator, right)
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

	return parser.call()
}

func (parser *AstParser) call() Expr {
	expr := parser.primary()

	for {
		if parser.match(references.LeftParen) {
			expr = parser.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (parser *AstParser) finishCall(callee Expr) Expr {
	var arguments []Expr
	if !parser.check(references.RightParen) {
		for ok := true; ok; ok = parser.match(references.Comma) {
			if len(arguments) > 255 {
				throwError(parser.peek(), "Can't have more than 255 arguments.")
			}
			arguments = append(arguments, parser.expression())
		}
	}

	paren := parser.consume(references.RightParen, "Expect ')' after arguments.")
	return NewCall(callee, paren, arguments)
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

func (parser *AstParser) previousIndex(index int) *scanner.Token {
	if index < 0 || index >= len(parser.Tokens) || parser.Tokens[index].Type == references.EOF {
		return nil
	}

	return parser.Tokens[index]
}

func (parser *AstParser) calculateDepth() (int, bool) {
	curr := parser.Current - 1
	depth := 0
	leftBraces := 0
	rightBraces := 0
	prev := parser.previousIndex(curr)
	found := false
	for prev != nil {
		if prev.Type == references.RightBrace {
			rightBraces++
		}

		if prev.Type == references.LeftBrace {
			leftBraces++
			if leftBraces > rightBraces {
				depth++
			}
		}

		if leftBraces > rightBraces && (prev.Type == references.For || prev.Type == references.While) {
			// if type is for loop then look for the initializer and increment the depth
			if prev.Type == references.For {
				forCurr := curr + 1
				forPrev := parser.previousIndex(forCurr)
				forFound := false
				for forPrev != nil && forPrev.Type != references.Semicolon && forPrev.Type != references.RightParen {
					if forPrev.Type == references.Var {
						forFound = true
						break
					}

					forCurr++
					forPrev = parser.previousIndex(forCurr)
				}

				// increment depth because an initializer means a second enclosing block
				if forFound {
					depth++
				}
			}

			found = true
			break
		}

		curr--
		prev = parser.previousIndex(curr)
	}

	return depth, found
}

func throwError(token *scanner.Token, message string) {
	loxerror.TokenError(token.Type, token.Line, token.Lexeme, message)

	panic(fmt.Errorf(message))
}

func throwRuntimeError(token *scanner.Token, message string) {
	loxerror.TokenRuntimeError(token.Type, token.Line, token.Lexeme, message, true)

	panic(fmt.Errorf(message))
}

func throwReturn(obj interface{}) {
	panic(obj)
}
