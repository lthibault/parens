package parser

import (
	"fmt"
	"io"
)

// KeywordExpr represents a keyword literal.
type KeywordExpr struct {
	Keyword string
}

// Eval returns the keyword itself.
func (ke KeywordExpr) Eval(scope Scope) (interface{}, error) {
	return ke.Keyword, nil
}

func (ke KeywordExpr) String() string {
	return ke.Keyword
}

func buildKeywordExpr(rd io.RuneScanner) (Expr, error) {
	if err := assertPrefix(rd, ':'); err != nil {
		return nil, err
	}

	kw := []rune{}
	for {
		ru, _, err := rd.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if isSepratingChar(ru) {
			rd.UnreadRune()
			break
		}

		if oneOf(ru, '\\') {
			return nil, fmt.Errorf("unexpected character '%c'", ru)
		}

		kw = append(kw, ru)
	}

	return KeywordExpr{Keyword: string(kw)}, nil
}
