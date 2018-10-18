package sizeflag

import (
	"bytes"
	"flag"
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	tests := []struct {
		in            string
		want2, want10 int
	}{
		{"0", 0, 0},
		{"1", 1, 1},
		{"12345678", 12345678, 12345678},

		// Fractional values work, and round toward ­∞
		{"0.5k", 512, 500},
		{"2.3K", 2355, 2300},
		{"1.2m", 1258291, 1200000},
		{"3.7P", 4165829655317709, 3700000000000000},
		{"0.25p", 1 * pi / 4, 1 * pd / 4},

		// Concatenated sizes are additive
		{"0K00m0", 0 + 0 + 0, 0},
		{"1k1k1k", 3 * 1024, 3 * 1000},
		{"2K15", 2063, 2015},
		{"2k15", 2063, 2015},
		{"1.7M0.3K", 1782886, 1700300},
		{"3m2.5k17", 3148305, 3002517},

		// Order of units does not matter, and units can be repeated
		{"2k 100", 100 + 2*ki, 2100},
		{"2k100", 100 + 2*ki, 2100},
		{"1k 1k 100", 100 + 2*ki, 2100},

		// Whitespace around sizes is ignored
		{" 123 ", 123, 123},
		{"1.4M 2K 17", 1470071, 1402017},
		{"\n\t2.1k\n\r\t 11 \t", 2161, 2111},
	}
	for _, test := range tests {
		t.Logf("Testing parse of %q...", test.in)
		t.Run("Base2", func(t *testing.T) {
			v := Base2(0)
			if err := v.Set(test.in); err != nil {
				t.Errorf("Set(%q) failed: %v", test.in, err)
			} else if got := int(*v); got != test.want2 {
				t.Errorf("Set(%q): got %d, want %d", test.in, got, test.want2)
			}
		})
		t.Run("Base10", func(t *testing.T) {
			v := Base10(0)
			if err := v.Set(test.in); err != nil {
				t.Errorf("Set(%q) failed: %v", test.in, err)
			} else if got := int(*v); got != test.want10 {
				t.Errorf("Set(%q): got %d, want %d", test.in, got, test.want10)
			}
		})
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
		// These should fail regardless of the units involved.
		for _, val := range []flag.Value{Base2(nil), Base10(nil)} {
			if err := val.Set(test); err == nil {
				t.Errorf("Parsing %q: got flag %#v, wanted error", test, val)
			}
		}
	}
}

func TestRoundTrip(t *testing.T) {
	tests := []int{
		0, 1, 2, 3, 7, 8, 11, 16, 100, 1023, 1024, 1025, 1597, 2584, 8192, 16384,
		1 * ki, 3 * (ki / 2), 2 * ki, 1027 * ki,
		mi, mi + 5, 1500 * mi,
		gi,
		ti, ti + 2*gi + 3*mi + 7,
		3 * (ti / 4), 1*pi + 32*ki,
	}
	for _, test := range tests {
		// These values should round-trip regardless of interpretation.
		for _, v := range []flag.Getter{Base2(test), Base10(test)} {
			u := v.String()
			t.Logf("Unparsed %v as %q", test, u)
			if err := v.Set(u); err != nil {
				t.Errorf("[%v].Set(%q): unexpected error: %v", v, u, err)
			}
			if got := v.Get().(int); got != test {
				t.Errorf("Round trip for %d failed: string %q reported %d", test, u, got)
			}
		}
	}
}

func TestFlagBits(t *testing.T) {
	const initial = 1024
	const sizeValue = 32256
	const countValue = 250
	const initString = "1K"

	size := Base2(1024)
	if s := size.String(); s != initString {
		t.Errorf("Initial -size string: got %q, want %q", s, initString)
	}
	var byteCount int

	var buf bytes.Buffer
	fs := flag.NewFlagSet("size", flag.PanicOnError)
	fs.Var(size, "size", "The number of bytes to corrupt on disk")
	fs.Var(Base10(&byteCount), "count", "The number of bytes that matter")
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

func ExampleBase2() {
	size := Base2(11)
	flag.Var(size, "size", "The size of the thing")

	var dim = 1024
	flag.Var(Base2(&dim), "dim", "The dimension of the thing")

	mass := Base2(2000)
	flag.Var(mass, "mass", "The mass of the thing")

	fmt.Printf("size %d %q\n", size.Int(), flag.Lookup("size").Value.String())
	fmt.Printf("dim %d %q\n", dim, flag.Lookup("dim").Value.String())
	fmt.Printf("mass %d %q\n", mass.Int(), flag.Lookup("mass").Value.String())
	// Output:
	// size 11 "11"
	// dim 1024 "1K"
	// mass 2000 "2000"
}
