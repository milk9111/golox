package scanner

import (
	"fmt"
	"golox/references"
)

type Token struct {
	Type    references.TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func NewToken(t references.TokenType, lexeme string, literal interface{}, line int) *Token {
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
