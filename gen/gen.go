package gen

import (
	"go/ast"

	"github.com/spy16/parens/parser"
)

// Generator provides functions for analyzing and generating Go code
// from LISP code.
type Generator struct {
}

// Generate parses the given lisp string and creates a Go AST.
func (gen *Generator) Generate(name, src string) (*ast.File, error) {
	expr, err := parser.Parse(name, src)
	if err != nil {
		return nil, err
	}

	af := &ast.File{}
	af.Name = &ast.Ident{Name: name}

	if err := emit(expr, af); err != nil {
		return nil, err
	}
	return af, nil
}

func emit(expr parser.Expr, af *ast.File) error {
	switch ex := expr.(type) {
	case parser.ModuleExpr:
		for _, em := range ex.Exprs {
			if err := emit(em, af); err != nil {
				return err
			}
		}
	case parser.ListExpr:
		sym, err := ex.Symbol()
		if err != nil {
			return err
		}
		decl := &ast.FuncDecl{}
		decl.Name = &ast.Ident{
			Name: sym,
		}
		af.Decls = append(af.Decls, decl)
	}

	return nil
}
