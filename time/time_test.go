package time

import (
	"bytes"
	"flag"
	"testing"
	"time"
)

func TestFlagBits(t *testing.T) {
	var ktime Value // default layout
	ptime := Value{
		Layout: "2006-01-02 15:04",
		Time:   time.Now(), // Default value
	}

	var buf bytes.Buffer
	fs := flag.NewFlagSet("time", flag.PanicOnError)
	fs.Var(&ktime, "ktime", ktime.Help("When to wake up"))
	fs.Var(&ptime, "ptime", ptime.Help("When the work is due"))

	fs.SetOutput(&buf)
	fs.PrintDefaults()
	t.Logf("Value flag set:\n%s", buf.String())
	buf.Reset()

	if err := fs.Parse([]string{"-ktime", "4:25AM", "-ptime", "2010-10-04 11:22"}); err != nil {
		t.Fatalf("Argument parsing failed: %v", err)
	}

	if got, want := ktime.Time.String(), "0000-01-01 04:25:00 +0000 UTC"; got != want {
		t.Errorf("Value for -ktime: got %q, want %q", got, want)
	}
	if got, want := ptime.Time.String(), "2010-10-04 11:22:00 +0000 UTC"; got != want {
		t.Errorf("Value for -ptime: got %q want %q", got, want)
	}
}
