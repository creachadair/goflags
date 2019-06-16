# goflags

http://godoc.org/github.com/creachadair/goflags

The packages in this repository define extensions to the standard Go "flag"
package for additional types beyond the standard set, by providing
implementations of the [`flag.Value`](http://golang.org/pkg/flag#Value) and
[`flag.Getter`](http://golang.org/pkg/flag#Getter) interfaces.

## Subpackages

### sizeflag

Defines a flag that accepts human-readable integer values, with shorthand for
multiples of powers of 1024. For example, "1K" for 1024, "2.1m" for 2202009,
and so forth.

### enumflag

Defines a flag that accepts a single string value taken from a fixed set of
values chosen when the flag is defined. Values are compared without regard to
case.

### regexpflag

Defines a flag that accepts a valid
[`*regexp.Regexp`](http://golang.org/pkg/regexp#Regexp)
value compiled from a string literal.

### timeflag

Defines a flag that accepts a [`time.Time`](http://golang.org/pkg/time#Time)
value parsed via a [standard format string](http://golang.org/pkg/time#Parse).
