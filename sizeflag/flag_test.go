package sizeflag

import (
	"bytes"
	"flag"
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"0", 0},
		{"1", 1},
		{"12345678", 12345678},

		// Fractional values work, and round toward ­∞
		{"0.5k", 512},
		{"2.3K", 2355},
		{"1.2m", 1258291},
		{"3.7P", 4165829655317709},
		{"0.25p", 1 * pUnits / 4},

		// Concatenated sizes are additive
		{"0K00m0", 0 + 0 + 0},
		{"1k1k1k", 3 * 1024},
		{"2K15", 2063},
		{"2k15", 2063},
		{"1.7M0.3K", 1782886},
		{"3m2.5k17", 3148305},

		// Order of units does not matter, and units can be repeated
		{"2k 100", 100 + 2*kUnits},
		{"2k100", 100 + 2*kUnits},
		{"1k 1k 100", 100 + 2*kUnits},

		// Whitespace around sizes is ignored
		{" 123 ", 123},
		{"1.4M 2K 17", 1470071},
		{"\n\t2.1k\n\r\t 11 \t", 2161},
	}
	for _, test := range tests {
		got, err := Parse(test.in)
		t.Logf("Parsed %q as %d", test.in, got)
		if err != nil {
			t.Errorf("Parse(%q): unexpected error: %v", test.in, err)
			continue
		}
		if got != test.want {
			t.Errorf("Parse(%q): got %d, want %d", test.in, got, test.want)
		}
	}
}

func TestParseErrors(t *testing.T) {
	tests := []string{
		"",         // Empty input
		"  ",       // Blank input
		"k",        // No value on unit
		".1",       // Missing leader
		"1.",       // Missing trailer
		"2.1",      // Unqualified fraction
		"17kk",     // Bare unit
		"-16.3mbb", // Negatives
		"1b2kqq",   // Garbage at end
		"6.7q",     // Invalid unit
		"0b1ce0",   // Invalid unit
	}
	for _, test := range tests {
		got, err := Parse(test)
		if err == nil {
			t.Errorf("ParseByteSize(%q): got %d, wanted error", test, got)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	tests := []int{
		0, 1, 2, 3, 7, 8, 11, 16, 100, 1023, 1024, 1025, 1597, 2584, 8192, 16384,
		1 * kUnits, 3 * (kUnits / 2), 2 * kUnits, 1027 * kUnits,
		mUnits, mUnits + 5, 1500 * mUnits,
		gUnits,
		tUnits, tUnits + 2*gUnits + 3*mUnits + 7,
		3 * (tUnits / 4), 1*pUnits + 32*kUnits,
	}
	for _, v := range tests {
		test := New(v)
		u := test.Unparse()
		got, err := Parse(u)
		if err != nil {
			t.Errorf("Invalid unparse result for %d: got %q", test, u)
			continue
		}

		if want := test.Int(); got != want {
			t.Errorf("Round trip failed: got %d, want %d [via %q]", got, want, u)
		} else {
			t.Logf("Unparsed %d as %q", want, u)
		}
	}
}

func TestFlagBits(t *testing.T) {
	const initial = 1024
	const sizeValue = 32256
	const countValue = 256
	const initString = "1K"

	size := New(1024)
	if s := size.String(); s != initString {
		t.Errorf("Size string: got %q, want %q", s, initString)
	}
	var byteCount int

	var buf bytes.Buffer
	fs := flag.NewFlagSet("size", flag.PanicOnError)
	fs.Var(size, "size", "The number of bytes to corrupt on disk")
	fs.Var(New(&byteCount), "count", "The number of bytes that matter")
	fs.SetOutput(&buf)
	fs.PrintDefaults()
	t.Logf("Size flag set:\n%s", buf.String())
	buf.Reset()

	if err := fs.Parse([]string{"-size", "31.5k", "-count", "0.25k"}); err != nil {
		t.Fatalf("Argument parsing failed: %v", err)
	}

	if v := size.Int(); v != sizeValue {
		t.Errorf("Size value after parsing: got %d, want %d", v, sizeValue)
	}
	if byteCount != countValue {
		t.Errorf("Count value after parsing: got %d, want %d", byteCount, countValue)
	}

	if err := size.Set("bogus"); err == nil {
		t.Errorf("Set bogus value: got %v, want error", err)
	} else {
		t.Logf("Set bogus value gave expected error: %v", err)
	}
}

func ExampleNew() {
	size := New(nil)
	flag.Var(size, "size", "The size of the thing")

	dim := 1024
	flag.Var(New(&dim), "dim", "The dimension of the thing")

	mass := New(2000)
	flag.Var(mass, "mass", "The mass of the thing")

	fmt.Printf("size %d %q\n", size.Int(), flag.Lookup("size").Value.String())
	fmt.Printf("dim %d %q\n", dim, flag.Lookup("dim").Value.String())
	fmt.Printf("mass %d %q\n", mass.Int(), flag.Lookup("mass").Value.String())
	// Output:
	// size 0 "0"
	// dim 1024 "1K"
	// mass 2000 "2000"
}
