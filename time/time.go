// Package time defines a flag.Value implementation that parses time.Time
// values through a format string.
//
// Example:
//   import (
//     "flag"
//
//     "bitbucket.org/creachadair/goflags/time"
//   )
//
//   var dueDate = time.Value{
//     Layout: "2006/01/02",
//     Time:   time.Now().Add(24*time.Hour),
//   }
//   func init() {
//     flag.Var(&dueDate, "due_date", dueDate.Help("When the work is due"))
//   }
//
package time

import (
	"fmt"
	"time"
)

// Value implements the flag.Value interface to a time.Time value.
type Value struct {
	// The layout of the string to parse, as accepted by time.Parse.
	// Defaults to time.Kitchen.
	Layout string

	// The time value parsed from the flag.
	Time time.Time
}

// String satisfies part of the flag.Value interface.
func (v *Value) String() string { return format(v.Time, v.Layout) }

// Help concatenates a human-readable string summarizing the format of t to h,
// for use in generating a documentation string.
func (v *Value) Help(h string) string {
	if v.Layout == "" {
		return fmt.Sprintf("%s (e.g., %q)", h, time.Kitchen)
	}
	return fmt.Sprintf("%s (e.g., %q)", h, v.Layout)
}

// Set satisfies part of the flag.Value interface.
func (v *Value) Set(s string) error {
	var err error
	v.Time, err = parse(s, v.Layout)
	return err
}

// Get satisfies the flag.Getter interface.
// The concrete value has type time.Time.
func (v *Value) Get() interface{} { return v.Time }

func parse(s string, format string) (time.Time, error) {
	if format == "" {
		format = time.Kitchen
	}
	return time.Parse(format, s)
}

func format(t time.Time, format string) string {
	if format == "" {
		return fmt.Sprintf("%q", t.Format(time.Kitchen))
	}
	return fmt.Sprintf("%q", t.Format(format))
}
