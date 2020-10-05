package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: generate_ast <output directory>")
		os.Exit(64)
	}

	defineAst(os.Args[1], "expression.go", "Expr", []string{
		"Assign : name *scanner.Token, value Expr",
		"Binary : left Expr, operator *scanner.Token, right Expr",
		"Grouping : expression Expr",
		"Literal : value interface{}",
		"Logical : left Expr, operator *scanner.Token, right Expr",
		"Unary : operator *scanner.Token, right Expr",
		"Variable : name *scanner.Token",
	})

	defineAst(os.Args[1], "statement.go", "Stmt", []string{
		"Block : statements []Stmt, isLoopIncrementer bool",
		"Expression : expression Expr",
		"IfCmd : condition Expr, thenBranch Stmt, elseBranch Stmt",
		"Print : expression Expr",
		"VarCmd : name *scanner.Token, initializer Expr",
		"WhileLoop : condition Expr, body Stmt",
		"BreakCmd : envDepth int",
		"ContinueCmd : envDepth int",
	})
}

func defineAst(outputDir string, filename string, baseName string, types []string) {
	path := fmt.Sprintf("%s/%s", outputDir, filename)

	visitorName := fmt.Sprintf("%sVisitor", baseName)

	sb := strings.Builder{}

	sb.WriteString("package syntax\n")
	sb.WriteString("\n")
	sb.WriteString("import \"golox/scanner\"\n")
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("type %s interface{\n", baseName))
	sb.WriteString(fmt.Sprintf("\taccept(visitor %s) interface{}\n", visitorName))
	sb.WriteString("\tString() string")
	sb.WriteString("}\n")
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("type %s interface {\n", visitorName))
	for _, t := range types {
		name := strings.TrimSpace(strings.Split(t, ":")[0])
		sb.WriteString(fmt.Sprintf("\tvisit%s%s(%s *%s) interface{}\n", name, baseName, strings.ToLower(baseName), name))
	}
	sb.WriteString("}\n")
	sb.WriteString("\n")

	for _, t := range types {
		parts := strings.Split(t, ":")
		structName := strings.TrimSpace(parts[0])
		fields := strings.TrimSpace(parts[1])
		defineType(&sb, baseName, structName, visitorName, fields)
	}

	err := ioutil.WriteFile(path, []byte(sb.String()), 0644)
	if err != nil {
		panic(err)
	}
}

func defineType(sb *strings.Builder, baseName string, structName string, visitorName string, fieldList string) {
	sb.WriteString(fmt.Sprintf("type %s struct {\n", structName))
	for _, f := range strings.Split(fieldList, ",") {
		sb.WriteString(fmt.Sprintf("\t%s\n", strings.TrimSpace(f)))
	}
	sb.WriteString("}\n")

	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("func New%s(%s) %s {\n", structName, fieldList, baseName))
	sb.WriteString(fmt.Sprintf("\treturn &%s{\n", structName))
	for _, f := range strings.Split(fieldList, ",") {
		name := strings.TrimSpace(strings.Split(strings.TrimSpace(f), " ")[0])
		sb.WriteString(fmt.Sprintf("\t\t%s: %s,\n", name, name))
	}
	sb.WriteString("\t}\n")
	sb.WriteString("}\n")
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("func (%s *%s) accept(visitor %s) interface{} {\n", strings.ToLower(structName), structName, visitorName))
	sb.WriteString(fmt.Sprintf("\treturn visitor.visit%s%s(%s)\n", structName, baseName, strings.ToLower(structName)))
	sb.WriteString("}\n")
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("func (%s *%s) String() string {\n", strings.ToLower(structName), structName))
	sb.WriteString(fmt.Sprintf("\treturn \"%s\"", structName))
	sb.WriteString("}\n")

	sb.WriteString("\n\n")
}
