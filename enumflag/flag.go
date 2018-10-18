// Package enumflag defines a flag.Value implementation that accepts one of a
// specified collection of string keys.  Values are compared without respect to
// case, so that "foo" and "Foo" are accepted as equivalent to "FOO".
//
// Example:
//   import (
//     "flag"
//
//     "bitbucket.org/creachadair/goflags/enum"
//   )
//
//   // The first enumerated value is the default.
//   var color = enum.Flag("", "red", "orange", "yellow", "green", "blue")
//   func init() {
//     flag.Var(color, "color", color.Help("What color to paint the bikeshed"))
//   }
//
package enumflag

import (
	"fmt"
	"sort"
	"strings"
)

// A Value represents an enumeration of string values.  A pointer to a Value
// satisfies the flag.Value interface. Use the Key method to recover the
// currently-selected value of the enumeration.
type Value struct {
	keys  []string
	index int // The selected index in the enumeration
}

// Help concatenates a human-readable string summarizing the legal values of v
// to h, for use in generating a documentation string.
func (v Value) Help(h string) string {
	return fmt.Sprintf("%s (%s)", h, strings.Join(v.keys, "|"))
}

// New returns a *Value for the specified enumerators, where defaultKey is the
// default value and otherKeys are additional options.
func New(defaultKey string, otherKeys ...string) *Value {
	v := &Value{keys: append(otherKeys, defaultKey)}
	sort.Strings(v.keys)
	for i, key := range v.keys {
		if key == defaultKey {
			v.index = i
			break
		}
	}
	return v
}

// Key returns the currently-selected key in the enumeration.  The original
// spelling of the selected value is returned, as given to the Flag
// constructor, not the value as parsed.
func (v Value) Key() string {
	if len(v.keys) == 0 {
		return "" // BUG: https://github.com/golang/go/issues/16694
	}
	return v.keys[v.index]
}

// Get satisfies the flag.Getter interface.
// The concrete value is the the string of the current key.
func (v Value) Get() interface{} { return v.Key() }

// Index returns the currently-selected index in the enumeration.
func (v Value) Index() int { return v.index }

// String satisfies part of the flag.Value interface.
func (v Value) String() string { return fmt.Sprintf("%q", v.Key()) }

// Set satisfies part of the flag.Value interface.
func (v *Value) Set(s string) error {
	for i, key := range v.keys {
		if strings.EqualFold(s, key) {
			v.index = i
			return nil
		}
	}
	return fmt.Errorf("expected one of (%s)", strings.Join(v.keys, "|"))
}
