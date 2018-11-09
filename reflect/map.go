package reflect

import (
	"fmt"
	"reflect"
	"strings"
)

// HashMap represents hash map or map data structure.
type HashMap struct {
	data map[Value]Value
}

// Interface returns the underlying value as empty-interface.
func (hm *HashMap) Interface() interface{} {
	return hm.data
}

// Type returns the type information for the underlying value.
func (hm *HashMap) Type() reflect.Type {
	return reflect.TypeOf(hm.data)
}

// Repr returns LISP representation of a hash-map.
func (hm *HashMap) Repr() string {
	kvPairs := []string{}

	for key, val := range hm.data {
		kvPairs = append(kvPairs, fmt.Sprintf("%s %s", key.Repr(), val.Repr()))
	}

	return fmt.Sprintf("{%s}", strings.Join(kvPairs, " "))
}
