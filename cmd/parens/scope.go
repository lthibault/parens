package main

import (
	"fmt"

	"github.com/spy16/parens"
	"github.com/spy16/parens/stdlib"
)

const version = "1.0.0"

const help = `
Welcome to Parens!

Type (exit) or Ctrl+C to exit the REPL.

Use (dump-scope) to see the list of symbols available in
the current scope.

Use (doc <symbol>) to get help about symbols in scope.

See "cmd/parens/main.go" in the github repository for
more information.

https://github.com/spy16/parens
`

func makeGlobalScope() parens.Scope {

	// user-defined values can be exposed too and their methods
	// can be accessed.
	st := &sampleType{val: "initial"}

	return stdlib.RegisterAll(parens.NewScope(nil)).
		Bind("parens-version", version).
		Bind("?", func() string { return help }).
		Bind("sample", st)
}

type sampleType struct {
	val string
}

func (st *sampleType) SetVal(s string) string {
	st.val = s
	return s
}

func (st sampleType) String() string {
	return fmt.Sprintf("sampleType[val=%s]", st.val)
}
