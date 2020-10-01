package main

import (
	"fmt"
	"golox/scanner"
	"golox/syntax"
)

func main() {
	expression := syntax.NewBinary(
		syntax.NewUnary(
			scanner.NewToken(scanner.Minus, "-", nil, 1),
			syntax.NewLiteral(123),
		),
		scanner.NewToken(scanner.Star, "*", nil, 1),
		syntax.NewGrouping(
			syntax.NewLiteral(45.67),
		),
	)

	fmt.Println(syntax.NewAstPrettyPrinter().Print(expression))
}
