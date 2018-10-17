package reflag

import (
	"flag"
	"testing"
)

func TestFlagBits(t *testing.T) {
	const probe = `^\s+`

	var match Value
	var skip = MustCompile(probe)

	fs := flag.NewFlagSet("regexp", flag.PanicOnError)
	fs.Var(&match, "match", "Match expression")
	fs.Var(&skip, "skip", "Skip expression")

	if v := match.String(); v != Empty {
		t.Errorf("Initial value for -match: got %q, want %q", v, Empty)
	}
	if v := skip.String(); v != probe {
		t.Errorf("Initial value for -skip: got %q, want %q", v, probe)
	}

	if err := fs.Parse([]string{"-match", `foo`, "-skip", `bar`}); err != nil {
		t.Fatalf("Argument parsing failed: %v", err)
	}

	if v := match.String(); v != "foo" {
		t.Errorf("Value for -match: got %q, want %q", v, "foo")
	}
	if v := skip.String(); v != "bar" {
		t.Errorf("Value for -skip: got %q, want %q", v, "bar")
	}

	// The methods of the underlying matcher can be used from the flag value.

	const needle = "well fooey on you too!"

	if match.FindString(needle) == "" {
		t.Errorf("Missing match for %q in %q", match, needle)
	}
	if m := skip.FindString(needle); m != "" {
		t.Errorf("Unexpected match for %q in %q: %v", skip, needle, m)
	}
}
