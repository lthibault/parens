package parser

import (
	"fmt"
	"io"
)

// ParseModule parses till the EOF and returns all s-exprs as a single ModuleExpr.
// This should be used to build an entire module from a file or string etc.
func ParseModule(name string, sc io.RuneScanner) (Expr, error) {
	me := ModuleExpr{}
	me.Name = name

	var expr Expr
	var err error
	for {
		expr, err = Parse(sc)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		me.Exprs = append(me.Exprs, expr)
	}

	return &me, nil
}

// Parse consumes runes from the reader until a single s-expression is extracted.
// Returns EOF if the reaches end-of-file before an s-exp is found. Returns any
// other errors from reader. This should be used when a continuous parse-eval from
// a stream is necessary (e.g. TCP socket).
func Parse(sc io.RuneScanner) (Expr, error) {
	var expr Expr
	var err error
	for {
		expr, err = buildExpr(sc)
		if err != nil {
			return nil, err
		}

		if expr != nil {
			return expr, nil
		}
	}
}

func buildExpr(rd io.RuneScanner) (Expr, error) {
	ru, _, err := rd.ReadRune()
	if err != nil {
		return nil, err
	}

	switch ru {
	case '"':
		rd.UnreadRune()
		return buildStrExpr(rd)
	case '(':
		rd.UnreadRune()
		return buildListExpr(rd)
	case '[':
		rd.UnreadRune()
		return buildVectorExpr(rd)
	case '\'':
		rd.UnreadRune()
		return buildQuoteExpr(rd)
	case ':':
		rd.UnreadRune()
		return buildKeywordExpr(rd)
	case ' ', '\t', '\n':
		return nil, nil
	case ';':
		rd.UnreadRune()
		return buildCommentExpr(rd)
	case ')', ']':
		return nil, io.EOF
	default:
		rd.UnreadRune()
		return buildSymbolOrNumberExpr(rd)
	}

}

// Expr represents an evaluatable expression.
type Expr interface {
	Eval(env Scope) (interface{}, error)
}

// Scope is responsible for managing bindings.
type Scope interface {
	Get(name string) (interface{}, error)
	Doc(name string) string
	Bind(name string, v interface{}, doc ...string) error
	Root() Scope
}

func assertPrefix(rd io.RuneScanner, prefix rune) error {
	ru, _, err := rd.ReadRune()
	if err != nil {
		return err
	}

	if ru != prefix {
		return fmt.Errorf("expected '%c' at the beginning, found '%c'", prefix, ru)
	}

	return nil
}

func isSepratingChar(ru rune) bool {
	return oneOf(ru, ' ', '\t', '\n', '\r', '(', ')', '[', ']', '{', '}', '"', '\'')
}

func oneOf(ru rune, set ...rune) bool {
	for _, rs := range set {
		if ru == rs {
			return true
		}
	}
	return false
}
