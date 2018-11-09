package reflect

import (
	"fmt"
	"reflect"
)

// String represents string values.
type String string

// Interface returns the string value as empty-interface.
func (str String) Interface() interface{} {
	return string(str)
}

// Type returns the type information for string.
func (str String) Type() reflect.Type {
	return reflect.TypeOf(str)
}

// Repr returns lisp representation of string.
func (str String) Repr() string {
	return fmt.Sprintf("\"%s\"", string(str))
}
