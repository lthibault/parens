package parser

import (
	"io"
	"strings"
)

type CommentExpr struct {
	comment string
}

func (ce CommentExpr) Eval(_ Scope) (interface{}, error) {
	return ce.comment, nil
}

func buildCommentExpr(rd io.RuneScanner) (Expr, error) {
	if err := assertPrefix(rd, ';'); err != nil {
		return nil, err
	}

	comment := []rune{}
	for {
		ru, _, err := rd.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if ru == '\n' {
			break
		}

		comment = append(comment, ru)
	}

	return CommentExpr{
		comment: strings.TrimSpace(string(comment)),
	}, nil
}
