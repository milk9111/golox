package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"lox/error"
	"lox/scanner"
	"os"
	"strings"
)

func main() {
	length := len(os.Args)
	fmt.Println(os.Args)
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
		panic(fmt.Errorf("%s does not exist", path))
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	run(string(data))

	if error.HadError() {
		os.Exit(65)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
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

	for _, token := range tokens {
		if token == nil {
			continue
		}

		fmt.Println(token.String())
	}
}
