package parser

import (
	"io"
	"regexp"
)

// TODO:
// - Support for hex (0x) and binary (0x) numbers
// - Support for scientific notation (1.3e10)
// - Clear differentiation between symbol and number
func buildSymbolOrNumberExpr(rd io.RuneScanner) (Expr, error) {
	seq := []rune{}
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

		seq = append(seq, ru)
	}

	s := string(seq)
	if numberRegex.MatchString(s) {
		return NumberExpr{
			NumStr: s,
		}, nil
	}

	return SymbolExpr{
		Symbol: s,
	}, nil
}

var numberRegex = regexp.MustCompile("^(\\+|-)?\\d+(\\.\\d+)?$")
