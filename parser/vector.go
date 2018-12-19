package parser

import (
	"fmt"
	"io"
	"strings"
)

// VectorExpr represents a vector form.
type VectorExpr struct {
	List []Expr
}

// Eval creates a golang slice.
func (ve VectorExpr) Eval(scope Scope) (interface{}, error) {
	lst := []interface{}{}

	for _, expr := range ve.List {
		val, err := expr.Eval(scope)
		if err != nil {
			return nil, err
		}
		lst = append(lst, val)
	}

	return lst, nil
}

func (ve VectorExpr) String() string {
	strs := []string{}
	for _, expr := range ve.List {
		strs = append(strs, fmt.Sprint(expr))
	}

	return fmt.Sprintf("[%s]", strings.Join(strs, " "))
}

func buildVectorExpr(rd io.RuneScanner) (Expr, error) {
	if err := assertPrefix(rd, '['); err != nil {
		return nil, err
	}

	vals := []Expr{}
	for {
		ru, _, err := rd.ReadRune()
		if err != nil {
			return nil, err
		}

		if ru == ']' {
			break
		}
		rd.UnreadRune()

		expr, err := buildExpr(rd)
		if err != nil {
			return nil, err
		}

		if expr != nil {
			vals = append(vals, expr)
		}
	}

	return VectorExpr{List: vals}, nil
}
