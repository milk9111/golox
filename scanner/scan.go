package scanner

import (
	"lox/error"
	"strconv"
)

var keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"for":    For,
	"fun":    Fun,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
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

	scanner.Tokens = append(scanner.Tokens, NewToken(EOF, "", nil, scanner.Line))
	return scanner.Tokens
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.Current >= len(scanner.Source)
}

func (scanner *Scanner) scanToken() {
	c := scanner.advance()
	switch c {
	case '(':
		scanner.addToken(LeftParen)
		break
	case ')':
		scanner.addToken(RightParen)
		break
	case '{':
		scanner.addToken(LeftBrace)
		break
	case '}':
		scanner.addToken(RightBrace)
		break
	case ',':
		scanner.addToken(Comma)
		break
	case '.':
		scanner.addToken(Dot)
		break
	case '-':
		scanner.addToken(Minus)
		break
	case '+':
		scanner.addToken(Plus)
		break
	case ';':
		scanner.addToken(Semicolon)
		break
	case '*':
		scanner.addToken(Star)
		break
	case '!':
		token := Bang
		if scanner.match('=') {
			token = BangEqual
		}
		scanner.addToken(token)
		break
	case '=':
		token := Equal
		if scanner.match('=') {
			token = EqualEqual
		}
		scanner.addToken(token)
		break
	case '<':
		token := Less
		if scanner.match('=') {
			token = LessEqual
		}
		scanner.addToken(token)
		break
	case '>':
		token := Greater
		if scanner.match('=') {
			token = GreaterEqual
		}
		scanner.addToken(token)
		break
	case '/':
		if scanner.match('/') {
			for scanner.peek() != '\n' && !scanner.isAtEnd() {
				scanner.advance()
			}
		} else {
			scanner.addToken(Slash)
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
			error.Error(scanner.Line, "Unexpected character.")
		}

		break
	}
}

func (scanner *Scanner) advance() rune {
	scanner.Current++
	return rune(scanner.Source[scanner.Current-1])
}

func (scanner *Scanner) addToken(t TokenType) {
	scanner.addTokenLiteral(t, nil)
}

func (scanner *Scanner) addTokenLiteral(t TokenType, literal interface{}) {
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
		error.Error(scanner.Line, "Unterminated string.")
		return
	}

	scanner.advance()

	value := string(scanner.Source[scanner.Start+1 : scanner.Current-1])
	scanner.addTokenLiteral(String, value)
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
		error.Error(scanner.Line, "Invalid number.")
		return
	}

	scanner.addTokenLiteral(Number, number)
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
		t = Identifier
	}

	scanner.addToken(t)
}
