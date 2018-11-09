package reflect

import (
	"fmt"
	"reflect"
	"strings"
)

// NewVector creates Vector data type with given values.
func NewVector(vals []Value) *Vector {
	return &Vector{
		vals: vals,
	}
}

// Vector represents the vector/array data structure.
type Vector struct {
	vals []Value
}

// Interface returns underlying value as empty-interface.
func (vec *Vector) Interface() interface{} {
	return vec.vals
}

// Type returns type information of underlying value.
func (vec *Vector) Type() reflect.Type {
	return reflect.TypeOf(vec.vals)
}

// Repr returns LISP representation of the vector.
func (vec *Vector) Repr() string {
	reprs := []string{}
	for _, val := range vec.vals {
		reprs = append(reprs, val.Repr())
	}
	return fmt.Sprintf("[%s]", strings.Join(reprs, " "))
}
