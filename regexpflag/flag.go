// Package regexpflag implements a convenience wrapper for defining flags that
// contain *regexp.Regexp values.
//
// Example:
//  import (
//    "flag"
//
//    "bitbucket.org/creachadair/goflags/regexpflag"
//  )
//
//  var match regexpflag.Value
//  var skip = regexpflag.MustCompile(`^\s+`)
//
//  func init() {
//    flag.Var(&match, "match", "Regular expression to match")
//    flag.Var(&skip, "skip", "Regular expression to skip")
//  }
//
package regexpflag

import "regexp"

const empty = "ø"

// A Value represents a regular expression.  The methods of the embedded
// *regexp.Regexp are available directly.  A pointer to a Value satisfies the
// flag.Value and flag.Getter interfaces.
type Value struct{ *regexp.Regexp }

// MustCompile parses s using the standard regexp.MustCompile function, and
// returns a Value containing the resulting *regexp.Regexp value.
func MustCompile(expr string) Value { return Value{regexp.MustCompile(expr)} }

// String satisfies part of the flag.Value interface.
func (v Value) String() string {
	if v.Regexp == nil {
		return empty
	}
	return v.Regexp.String()
}

// Set satisfies part of the flag.Value interface.
func (v *Value) Set(s string) error {
	r, err := regexp.Compile(s)
	if err != nil {
		return err
	}
	v.Regexp = r
	return nil
}

// Get satisfies the flag.Getter interface.
// The concrete value has type *regexp.Regexp.
func (v *Value) Get() interface{} { return v.Regexp }
