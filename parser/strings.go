package parser

import (
	"io"
)

// StringExpr represents single and double quoted strings.
type StringExpr struct {
	Value string
}

// Eval returns unquoted version of the STRING token.
func (se StringExpr) Eval(_ Scope) (interface{}, error) {
	return se.Value, nil
}

func buildStrExpr(rd io.RuneScanner) (Expr, error) {
	if err := assertPrefix(rd, '"'); err != nil {
		return nil, err
	}

	val := []rune{}

	for {
		ru, _, err := rd.ReadRune()
		if err != nil {
			return nil, err
		}

		lastI := len(val) - 1
		if len(val) >= 1 && val[lastI] == '\\' {
			var esc byte
			switch ru {
			case 'n':
				esc = '\n'
			case 't':
				esc = '\t'
			case 'r':
				esc = '\r'
			case '"':
				esc = '"'
			}

			if esc != 0 {
				val[lastI] = rune(esc)
				continue
			}

		}
		if oneOf(ru, 't', 'n', '"', 'r') {
		}
		if ru == '"' {
			break
		}

		val = append(val, ru)
	}

	return StringExpr{Value: string(val)}, nil
}
