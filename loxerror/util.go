package loxerror

import (
	"fmt"
	"golox/references"
)

var hadError = false
var hadRuntimeError = false

func Error(line int, message string) {
	Report(line, "", message, false)
}

func TokenError(t references.TokenType, line int, lexeme string, message string) {
	TokenRuntimeError(t, line, lexeme, message, false)
}

func TokenRuntimeError(t references.TokenType, line int, lexeme string, message string, isRuntimeError bool) {
	if t == references.EOF {
		Report(line, " at the end", message, isRuntimeError)
	} else {
		Report(line, fmt.Sprintf(" at '%s'", lexeme), message, isRuntimeError)
	}
}

func Report(line int, where string, message string, isRuntimeError bool) {
	fmt.Printf("%s\n", fmt.Errorf("[line %d] Error%s: %s", line, where, message).Error())
	hadError = !isRuntimeError
	hadRuntimeError = isRuntimeError
}

func HadError() bool {
	return hadError
}

func HadRuntimeError() bool {
	return hadRuntimeError
}
