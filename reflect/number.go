package reflect

import (
	"fmt"
	"reflect"
)

// Number represents numeric data.
type Number float64

// Interface returns number as empty-interface.
func (num Number) Interface() interface{} {
	return float64(num)
}

// Type returns underlying type information of number.
func (num Number) Type() reflect.Type {
	return reflect.TypeOf(num)
}

// Repr returns LISP representation of number.
func (num Number) Repr() string {
	return fmt.Sprintf("%f", float64(num))
}
