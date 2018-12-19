package parser

import (
	"fmt"
	"io"
)

// QuoteExpr implements the quote-literal form.
type QuoteExpr struct {
	Expr Expr
}

// Eval returns the expression itself without evaluating it.
func (qe QuoteExpr) Eval(scope Scope) (interface{}, error) {
	return qe.Expr, nil
}

// UnquoteEval unquotes and evaluates the underlying expression.
func (qe QuoteExpr) UnquoteEval(scope Scope) (interface{}, error) {
	return qe.Expr.Eval(scope)
}

func (qe QuoteExpr) String() string {
	return fmt.Sprintf("'%s", qe.Expr)
}

func buildQuoteExpr(rd io.RuneScanner) (Expr, error) {
	if err := assertPrefix(rd, '\''); err != nil {
		return nil, err
	}

	expr, err := buildExpr(rd)
	if err != nil {
		return nil, err
	}

	return QuoteExpr{Expr: expr}, nil
}
