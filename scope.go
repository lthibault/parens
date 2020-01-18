package parens

import (
	"fmt"
	"strings"
)

// NewScope initializes a new scope with given parent scope. parent
// can be nil.
func NewScope(parent Scope) Scope {
	return defaultScope{
		parent: parent,
		vals:   newPersistentScope([]PersistentScopeItem{}),
	}
}

type defaultScope struct {
	parent Scope
	vals   *persistentScope
}

func (sc defaultScope) Root() Scope {
	if sc.parent == nil {
		return sc
	}

	return sc.parent.Root()
}

func (sc defaultScope) Bind(name string, v interface{}, doc ...string) Scope {
	return defaultScope{
		parent: sc.parent,
		vals: sc.vals.Store(name, scopeEntry{
			val: newValue(v),
			doc: strings.TrimSpace(strings.Join(doc, "\n")),
		}),
	}
}

func (sc defaultScope) Doc(name string) string {
	if entry, found := sc.entry(name); found {
		return entry.doc
	}

	if sc.parent != nil {
		if swd, ok := sc.parent.(scopeWithDoc); ok {
			return swd.Doc(name)
		}
	}

	return ""
}

func (sc defaultScope) Get(name string) (interface{}, error) {
	entry, found := sc.entry(name)
	if !found {
		if sc.parent != nil {
			return sc.parent.Get(name)
		}
		return nil, fmt.Errorf("name '%s' not found", name)
	}

	return entry.val.RVal.Interface(), nil
}

func (sc defaultScope) String() string {
	str := []string{}
	sc.vals.Range(func(name string, _ scopeEntry) bool {
		str = append(str, fmt.Sprintf("%s", name))
		return true
	})
	return strings.Join(str, "\n")
}

func (sc defaultScope) entry(name string) (entry scopeEntry, found bool) {
	var v interface{}
	if v, found = sc.vals.Load(name); found {
		entry = v.(scopeEntry)
	}

	return
}

type scopeEntry struct {
	val reflectVal
	doc string
}

type scopeWithDoc interface {
	Doc(name string) string
}
