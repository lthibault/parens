package reflect

import "reflect"

// Value is returned by Expr implementations when Evaluated.
type Value interface {
	Interface() interface{}
	Type() reflect.Type
	Repr() string
}
