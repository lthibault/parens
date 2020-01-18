package stdlib

import "github.com/spy16/parens"

// RegisterAll registers different built-in functions into the
// given scope.
func RegisterAll(scope parens.Scope) parens.Scope {
	return doUntilErr(scope,
		RegisterCore,
		RegisterMath,
		RegisterIO,
		RegisterSystem,
	)
}

// RegisterSystem binds system functions into the scope.
func RegisterSystem(scope parens.Scope) parens.Scope {
	return registerList(scope, system)
}

// RegisterIO binds input/output functions into the scope.
func RegisterIO(scope parens.Scope) parens.Scope {
	return registerList(scope, io)
}

// RegisterCore binds all the core macros and functions into
// the scope.
func RegisterCore(scope parens.Scope) parens.Scope {
	scope.Bind("eval", Eval(scope),
		"Executes given LISP string in the current scope",
		"Usage: (eval <form>)",
	)

	scope.Bind("load", LoadFile(scope),
		"Reads and executes the file in the current scope",
		"Example: (load \"sample.lisp\")",
	)

	return registerList(scope, core)
}

// RegisterMath binds basic math operators into the scope.
func RegisterMath(scope parens.Scope) parens.Scope {
	return registerList(scope, math)
}

func doUntilErr(scope parens.Scope, fns ...func(scope parens.Scope) parens.Scope) parens.Scope {
	for _, fn := range fns {
		scope = fn(scope)
	}

	return scope
}

func registerList(scope parens.Scope, entries []mapEntry) parens.Scope {
	for _, entry := range entries {
		scope = scope.Bind(entry.name, entry.val, entry.doc...)
	}

	return scope
}

func entry(name string, val interface{}, doc ...string) mapEntry {
	return mapEntry{
		name: name,
		val:  val,
		doc:  doc,
	}
}

type mapEntry struct {
	name string
	val  interface{}
	doc  []string
}
