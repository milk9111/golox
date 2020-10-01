package main

import (
	"bufio"
	"fmt"
	"golox/loxerror"
	"golox/scanner"
	"golox/syntax"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var interpreter = syntax.NewInterpreter()

func main() {
	length := len(os.Args)
	if length > 2 {
		fmt.Printf("Usage: golox [script]")
		os.Exit(64)
	} else if length == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	if !fileExists(path) {
		fmt.Printf("%s does not exist\n", path)
		os.Exit(64)
	}

	if _, filename := filepath.Split(path); !strings.HasSuffix(filename, ".lox") {
		fmt.Println("Not a Lox file")
		os.Exit(64)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(64)
	}

	run(string(data))

	if loxerror.HadError() {
		os.Exit(65)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(64)
		}

		line = strings.Replace(line, "\n", "", -1)
		if line == "" {
			break
		}

		run(line)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func run(source string) {
	scanner := scanner.NewScanner(source)
	tokens := scanner.ScanTokens()

	parser := syntax.NewAstParser(tokens)
	statements := parser.Parse()

	if loxerror.HadError() {
		os.Exit(65)
	}

	interpreter.Interpret(statements)

	if loxerror.HadRuntimeError() {
		os.Exit(70)
	}
}
