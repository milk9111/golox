package error

import (
	"fmt"
)

var hadError = false

func Error(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where string, message string) {
	fmt.Printf("%s\n", fmt.Errorf("[line %d] Error%s: %s", line, where, message).Error())
	hadError = true
}

func HadError() bool {
	return hadError
}
