// Package sizeflag provides flag.Value implementations that support a
// convenient human-readable notation for integer sizes.
//
// The grammar of size strings is as follows:
//
//   size   = number unit [size]
//          | digits
//   number = digits ['.' digits]
//   unit   = 'k' | 'm' | 'g' | 't' | 'p' | 'e'
//   digits = [0-9]+
//
// For example: 25, 3K, 6.5g, 1.1T.
// Whitespace surrounding or separating size terms is ignored.
//
// The units are case-insensitive, and represent the following quantities:
//        Base10           Base2
//    k = 1000   = 10^3    1024    = 2^10
//    m = 1000*k = 10^6    1024*ki = 2^20
//    g = 1000*m = 10^9    1024*mi = 2^30
//    t = 1000*g = 10^12   1024*gi = 2^40
//    p = 1000*t = 10^15   1024*ti = 2^50
//    e = 1000*p = 10^18   1024*pi = 2^60
//
// A number without a tag is interpreted as a number of units, e.g., 25.
//
// If multiple sizes are concatenated, the resulting size is the sum of the
// terms, e.g., 2k15 represents 2k + 15 or 2048 + 15 = 2063 units.
//
// Fractional values are rounded toward -∞, e.g., 2.3k = 2355.
//
// Each size term is separately rounded in this way, so that
// 1.7M0.3K = 1782579 + 307 = 1782886.
package sizeflag

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// A Value2 represents a flaggable integer value scaled by powers of 2.
// A *Value2 satisfies the flag.Getter interface.
type Value2 int

// A Value10 represents a flaggable integer value scaled by powers of 10.
// A *Value10 satisfies the flag.Getter interface.
type Value10 int

// Int returns the value of the flag as an int.
func (v Value2) Int() int { return int(v) }

// Int return sthe value of the flag as an int.
func (v Value10) Int() int { return int(v) }

// String renders the current value of the flag as a string.
func (v Value2) String() string { return unparse(int(v), 1024, mult2) }

// String renders the current value of the flag as a string.
func (v Value10) String() string { return unparse(int(v), 1000, mult10) }

// Get retrieves the current value of the flag with concrete type int.
func (v Value2) Get() interface{} { return int(v) }

// Get retrieves the current value of the flag with concrete type int.
func (v Value10) Get() interface{} { return int(v) }

// Set sets the value of the flag from the specified string.
func (v *Value2) Set(s string) error {
	z, err := parse(s, units2)
	if err == nil {
		*v = Value2(z)
	}
	return err
}

// Set sets the value of the flag from the specified string.
func (v *Value10) Set(s string) error {
	z, err := parse(s, units10)
	if err == nil {
		*v = Value10(z)
	}
	return err
}

// Base2 returns a *Value2 initialized by v.
//
// If v has type *int, the parsed value will be stored in *v, and the default
// flag value will be taken from *v.
//
// If v == nil the default flag value is 0 and a fresh location is allocated
// and returned to receive the parsed value.
//
// If v has type int, the default flag value will be v, and a fresh location is
// allocated and returned to receive the parsed value.
//
// Any other value will cause Base2 to panic.
func Base2(v interface{}) *Value2 {
	switch t := v.(type) {
	case nil:
		return new(Value2)
	case *Value2:
		return t
	case int:
		return (*Value2)(&t)
	case *int:
		return (*Value2)(t)
	default:
		panic("invalid flag initializer")
	}
}

// Base10 returns a *Value10 initialized by v.
//
// If v has type *int, the parsed value will be stored in *v, and the default
// flag value will be taken from *v.
//
// If v == nil the default flag value is 0 and a fresh location is allocated
// and returned to receive the parsed value.
//
// If v has type int, the default flag value will be v, and a fresh location is
// allocated and returned to receive the parsed value.
//
// Any other value will cause Base10 to panic.
func Base10(v interface{}) *Value10 {
	switch t := v.(type) {
	case nil:
		return new(Value10)
	case *Value10:
		return t
	case int:
		return (*Value10)(&t)
	case *int:
		return (*Value10)(t)
	default:
		panic("invalid flag initializer")
	}
}

var sizeRE = regexp.MustCompile(`^(?i)([0-9]+(?:\.[0-9]+)?)([bkmgtp])`)

const (
	kd = 1000
	md = kd * kd
	gd = md * kd
	td = gd * kd
	pd = td * kd
	ed = pd * kd

	ki = 1024
	mi = ki * ki
	gi = mi * ki
	ti = gi * ki
	pi = ti * ki
	ei = pi * ki
)

var (
	units2  = map[string]float64{"k": ki, "m": mi, "g": gi, "t": ti, "p": pi, "e": ei}
	units10 = map[string]float64{"k": kd, "m": md, "g": gd, "t": td, "p": pd, "e": ed}
	mult2   = []int{ei, pi, ti, gi, mi, ki}              // descending order
	mult10  = []int{ed, pd, td, gd, md, kd}              // descending order
	labels  = []string{"", "E", "P", "T", "G", "M", "K"} // descending order

	// N.B. labels[0] is a sentinel.
)

// parse parses a human-readable string defining a number of units in the given
// base, and returns the number of units so defined.
func parse(s string, unit map[string]float64) (int, error) {
	var size int
	var ok bool
	for {
		s = strings.TrimSpace(s)
		m := sizeRE.FindStringSubmatch(s)
		if m == nil {
			break
		}
		v, err := strconv.ParseFloat(m[1], 64)
		if err != nil {
			return 0, fmt.Errorf("sizeflag: invalid size %q", m[0])
			// Should not be structurally possible, though.
		}
		if mul, ok := unit[strings.ToLower(m[2])]; ok {
			v *= mul
		} else {
			return size, fmt.Errorf("sizeflag: invalid unit %q", m[2])
		}
		size += int(v)
		s = s[len(m[0]):]
		ok = true
	}
	if s = strings.TrimSpace(s); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("sizeflag: invalid size %q", s)
		}
		size += v
	} else if !ok {
		return 0, fmt.Errorf("sizeflag: invalid size %q", s)
	}
	return size, nil
}

// unparse renders a non-negative int into a human-readable string, reversing
// the grammar understood by parse, so that the resulting values round-trip.
// Specificaly, if
//
//    n, err := parse(s, unitsN)
//
// and err == nil, then
//
//    p, err := parse(unparse(n, multN))
//
// yields err == nil and p == n.
func unparse(v int, pow int, mult []int) string {
	type term struct {
		n int
		u string
	}
	var terms []term
	add := func(n int, u, prev string, v int) int {
		// If the remaining value is zero and there is a previous term one place
		// higher, lower the previous term by one place and combine them.
		// For example, 1G+1M = 1025M with pow == 1024.

		if p := len(terms) - 1; p >= 0 && v == 0 && terms[p].u == prev {
			terms[p].n = terms[p].n*pow + n
			terms[p].u = u
		} else {
			terms = append(terms, term{n, u})
		}
		return v
	}

	z := int(v)
	for i, div := range mult {
		if n := int(z / div); n > 0 {
			z = add(n, labels[i+1], labels[i], z%div)
		}
	}
	if len(terms) == 0 || z > 0 {
		add(z, "", "K", 0)
	}

	parts := make([]string, len(terms))
	for i, t := range terms {
		parts[i] = fmt.Sprintf("%d%s", t.n, t.u)
	}
	return strings.Join(parts, " ")
}