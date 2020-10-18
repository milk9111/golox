package scanner

import (
	"golox/loxerror"
	"golox/references"
	"strconv"
)

var keywords = map[string]references.TokenType{
	"and":      references.And,
	"new":      references.New,
	"class":    references.Class,
	"else":     references.Else,
	"false":    references.False,
	"for":      references.For,
	"fun":      references.Fun,
	"if":       references.If,
	"nil":      references.Nil,
	"or":       references.Or,
	"print":    references.Print,
	"return":   references.Return,
	"super":    references.Super,
	"this":     references.This,
	"true":     references.True,
	"var":      references.Var,
	"while":    references.While,
	"continue": references.Continue,
	"break":    references.Break,
}

type Scanner struct {
	Source  string
	Tokens  []*Token
	Start   int
	Current int
	Line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		Source:  source,
		Start:   0,
		Current: 0,
		Line:    1,
	}
}

func (scanner *Scanner) ScanTokens() []*Token {
	for !scanner.isAtEnd() {
		scanner.Start = scanner.Current
		scanner.scanToken()
	}

	scanner.Tokens = append(scanner.Tokens, NewToken(references.EOF, "", nil, scanner.Line))
	return scanner.Tokens
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.Current >= len(scanner.Source)
}

func (scanner *Scanner) scanToken() {
	c := scanner.advance()
	switch c {
	case '(':
		scanner.addToken(references.LeftParen)
		break
	case ')':
		scanner.addToken(references.RightParen)
		break
	case '{':
		scanner.addToken(references.LeftBrace)
		break
	case '}':
		scanner.addToken(references.RightBrace)
		break
	case ',':
		scanner.addToken(references.Comma)
		break
	case '.':
		scanner.addToken(references.Dot)
		break
	case '%':
		scanner.addToken(references.Modulo)
		break
	case '-':
		token := references.Minus
		if scanner.peek() == '-' {
			scanner.advance()
			token = references.DecrementOne
		}

		if scanner.peek() == '=' {
			scanner.advance()
			token = references.Decrement
		}

		scanner.addToken(token)
		break
	case '+':
		token := references.Plus
		if scanner.peek() == '+' {
			scanner.advance()
			token = references.IncrementOne
		}

		if scanner.peek() == '=' {
			scanner.advance()
			token = references.Increment
		}

		scanner.addToken(token)
		break
	case ';':
		scanner.addToken(references.Semicolon)
		break
	case '*':
		scanner.addToken(references.Star)
		break
	case '!':
		token := references.Bang
		if scanner.match('=') {
			token = references.BangEqual
		}
		scanner.addToken(token)
		break
	case '=':
		token := references.Equal
		if scanner.match('=') {
			token = references.EqualEqual
		}
		scanner.addToken(token)
		break
	case '<':
		token := references.Less
		if scanner.match('=') {
			token = references.LessEqual
		}
		scanner.addToken(token)
		break
	case '>':
		token := references.Greater
		if scanner.match('=') {
			token = references.GreaterEqual
		}
		scanner.addToken(token)
		break
	case '/':
		if scanner.match('/') {
			for scanner.peek() != '\n' && !scanner.isAtEnd() {
				scanner.advance()
			}
		} else if scanner.match('*') {
			for !scanner.isAtEnd() {
				if scanner.match('*') && scanner.match('/') {
					break
				}

				if c := scanner.advance(); c == '\n' {
					scanner.Line++
				}
			}
		} else {
			scanner.addToken(references.Slash)
		}
		break
	case ' ':
		break
	case '\r':
		break
	case '\t':
		break
	case '\n':
		scanner.Line++
		break
	case '"':
		scanner.parseString()
		break
	default:
		if isDigit(c) {
			scanner.number()
		} else if isAlpha(c) {
			scanner.identifier()
		} else {
			loxerror.Error(scanner.Line, "Unexpected character.")
		}

		break
	}
}

func (scanner *Scanner) advance() rune {
	scanner.Current++
	return rune(scanner.Source[scanner.Current-1])
}

func (scanner *Scanner) addToken(t references.TokenType) {
	scanner.addTokenLiteral(t, nil)
}

func (scanner *Scanner) addTokenLiteral(t references.TokenType, literal interface{}) {
	text := scanner.Source[scanner.Start:scanner.Current]
	scanner.Tokens = append(scanner.Tokens, NewToken(t, text, literal, scanner.Line))
}

func (scanner *Scanner) match(expected rune) bool {
	if scanner.isAtEnd() ||
		rune(scanner.Source[scanner.Current]) != expected {
		return false
	}

	scanner.Current++
	return true
}

func (scanner *Scanner) peek() rune {
	if scanner.isAtEnd() {
		return '\000'
	}

	return rune(scanner.Source[scanner.Current])
}

func (scanner *Scanner) parseString() {
	for scanner.peek() != '"' && !scanner.isAtEnd() {
		if scanner.peek() == '\n' {
			scanner.Line++
		}

		scanner.advance()
	}

	if scanner.isAtEnd() {
		loxerror.Error(scanner.Line, "Unterminated string.")
		return
	}

	scanner.advance()

	value := string(scanner.Source[scanner.Start+1 : scanner.Current-1])
	scanner.addTokenLiteral(references.String, value)
}

func (scanner *Scanner) number() {
	for isDigit(scanner.peek()) {
		scanner.advance()
	}

	if scanner.peek() == '.' && isDigit(scanner.peekNext()) {
		scanner.advance()

		for isDigit(scanner.peek()) {
			scanner.advance()
		}
	}

	number, err := strconv.ParseFloat(scanner.Source[scanner.Start:scanner.Current], 64)
	if err != nil {
		loxerror.Error(scanner.Line, "Invalid number.")
		return
	}

	scanner.addTokenLiteral(references.Number, number)
}

func (scanner *Scanner) peekNext() rune {
	if scanner.Current+1 >= len(scanner.Source) {
		return '\000'
	}

	return rune(scanner.Source[scanner.Current+1])
}

func (scanner *Scanner) identifier() {
	for isAlphaNumeric(scanner.peek()) {
		scanner.advance()
	}

	text := scanner.Source[scanner.Start:scanner.Current]

	t, ok := keywords[text]
	if !ok {
		t = references.Identifier
	}

	scanner.addToken(t)
}
