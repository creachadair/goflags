// Package sizeflag provides flag.Value implementation that supports a
// convenient human-readable notation for integer sizes scaled by powers of
// 1024.
package sizeflag

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// A Value represents a flaggable integer value.  A pointer to a Value
// satisfies the flag.Value interface taking values accepted by Parse.
type Value int

// Int returns the value of the flag as an int.
func (v Value) Int() int { return int(v) }

// String implements part of the flag.Value interface.
func (v Value) String() string { return v.Unparse() }

// Get implements part of the flag.Getter interface.  The value returned has
// concrete type int.
func (v Value) Get() interface{} { return int(v) }

// Set implements part of the flag.Value interface.
func (v *Value) Set(s string) error {
	z, err := Parse(s)
	if err == nil {
		*v = Value(z)
	}
	return err
}

// New returns a *Value that satisfies the flag.Getter interface.  If v has
// type *int, it is converted and used as the flag location, with initial value
// *v.  If v == nil or v has type int, a fresh location is allocated; in the
// former case, the initial value is 0; in the latter it is v.  Any other value
// will cause New to panic.
func New(v interface{}) *Value {
	switch t := v.(type) {
	case nil:
		return new(Value)
	case int:
		return (*Value)(&t)
	case *int:
		return (*Value)(t)
	default:
		panic("invalid flag initializer")
	}
}

var sizeRE = regexp.MustCompile(`^(?i)([0-9]+(?:\.[0-9]+)?)[bkmgtp]`)

const (
	nUnits = 1
	scale  = 1024
	kUnits = nUnits * scale
	mUnits = kUnits * scale
	gUnits = mUnits * scale
	tUnits = gUnits * scale
	pUnits = tUnits * scale
)

// Parse parses a human-readable string defining a number of units, and returns
// the number of units so defined.
//
// Grammar:
//   size   = number unit [size]
//          | digits
//   number = digits ['.' digits]
//   unit   = 'k' | 'm' | 'g' | 't' | 'p'
//   digits = [0-9]+
//
// For example: 25, 3K, 6.5g, 1.1T.
// Whitespace surrounding or separating size terms is ignored.
//
// The units are case-insensitive, and represent the following quantities:
//    k  -- 1024 units
//    m  -- 1024 kUnits -- 1024*1024 (1,048,576) units
//    g  -- 1024 mUnits -- 1024*1024*1024 (1,073,741,824) units
//    t  -- 1024 gUnits -- 1024*1024*1024*1024 (1,099,511,627,776) units
//    p  -- 1024 tUnits -- 1024*1024*1024*1024*1024 (1,125,899,906,842,624) units
//
// A number without a tag is interpreted as a number of units.  If multiple
// sizes are concatenated, the resulting size is the sum of the terms, e.g.,
// 2k15 represents 2k + 15 or 2048 + 15 = 2063 units.  Fractional values are
// rounded toward -∞, e.g., 2.3k = 2355.  Each size term is separately rounded
// in this way, e.g. 1.7M0.3K = 1782579 + 307 = 1782886.
//
func Parse(s string) (int, error) {
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
			return 0, fmt.Errorf("intsize: invalid size %q", m[0])
			// Should not be structurally possible, though.
		}
		switch u := strings.ToLower(m[0][len(m[0])-1:]); u {
		case "k":
			v *= kUnits
		case "m":
			v *= mUnits
		case "g":
			v *= gUnits
		case "t":
			v *= tUnits
		case "p":
			v *= pUnits
		default:
			return size, fmt.Errorf("intsize: invalid unit %q", u)
		}
		size += int(v)
		s = s[len(m[0]):]
		ok = true
	}
	if s = strings.TrimSpace(s); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("intsize: invalid size %q", s)
		}
		size += v
	} else if !ok {
		return 0, fmt.Errorf("intsize: invalid size %q", s)
	}
	return size, nil
}

// Unparse renders a Value into a human-readable string following the same grammar
// understood by Parse, so that Parse(Value(n).Unparse()) == n.
func (v Value) Unparse() string {
	type term struct {
		n int
		u string
	}
	var terms []term
	add := func(n int, u, prev string, v int) int {
		// If the remaining value is zero and there is a previous term one place
		// higher, lower the previous term by one place and combine them.
		// For example, 1G+1M = 1025M.

		if p := len(terms) - 1; p >= 0 && v == 0 && terms[p].u == prev {
			terms[p].n = terms[p].n*1024 + n
			terms[p].u = u
		} else {
			terms = append(terms, term{n, u})
		}
		return v
	}

	z := int(v)
	if n := int(z / pUnits); n > 0 {
		z = add(n, "P", "", z%pUnits)
	}
	if n := int(z / tUnits); n > 0 {
		z = add(n, "T", "P", z%tUnits)
	}
	if n := int(z / gUnits); n > 0 {
		z = add(n, "G", "T", z%gUnits)
	}
	if n := int(z / mUnits); n > 0 {
		z = add(n, "M", "G", z%mUnits)
	}
	if n := int(z / kUnits); n > 0 {
		z = add(n, "K", "M", z%kUnits)
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
