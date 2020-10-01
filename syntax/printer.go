package syntax

import (
	"fmt"
	"strings"
)

type AstPrettyPrinter struct{}

func NewAstPrettyPrinter() *AstPrettyPrinter {
	return &AstPrettyPrinter{}
}

func (printer *AstPrettyPrinter) Print(expr Expr) string {
	return expr.accept(printer).(string)
}

func (printer *AstPrettyPrinter) visitAssignExpr(expr *Assign) interface{} {
	return printer.parenthesize(expr.name.Lexeme, expr.value)
}

func (printer *AstPrettyPrinter) visitVartStmt(stmt *Vart) interface{} {
	return printer.parenthesize(stmt.name.Lexeme, stmt.initializer)
}

func (printer *AstPrettyPrinter) visitVariableExpr(expr *Variable) interface{} {
	return printer.parenthesize(expr.name.Lexeme)
}

func (printer *AstPrettyPrinter) visitBinaryExpr(expr *Binary) interface{} {
	return printer.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (printer *AstPrettyPrinter) visitGroupingExpr(expr *Grouping) interface{} {
	return printer.parenthesize("group", expr.expression)
}

func (printer *AstPrettyPrinter) visitLiteralExpr(expr *Literal) interface{} {
	if expr.value == nil {
		return "nil"
	}

	return fmt.Sprintf("%v", expr.value)
}

func (printer *AstPrettyPrinter) visitUnaryExpr(expr *Unary) interface{} {
	return printer.parenthesize(expr.operator.Lexeme, expr.right)
}

func (printer *AstPrettyPrinter) parenthesize(name string, exprs ...Expr) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("(%s", name))
	for _, expr := range exprs {
		sb.WriteString(" ")
		sb.WriteString(expr.accept(printer).(string))
	}
	sb.WriteString(")")

	return sb.String()
}
