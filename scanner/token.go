package scanner

import (
	"fmt"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func NewToken(t TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		Type:    t,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (token *Token) String() string {
	return fmt.Sprintf("%d %s %v", int(token.Type), token.Lexeme, token.Literal)
}
