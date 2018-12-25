package enumflag_test

import (
	"flag"
	"fmt"
	"log"

	"bitbucket.org/creachadair/goflags/enumflag"
)

func Example() {
	fs := flag.NewFlagSet("example", flag.ContinueOnError)

	var feature = enumflag.New("auto", "on", "off")
	fs.Var(feature, "feature", feature.Help("Enable the new behaviour"))

	if err := fs.Parse([]string{"-feature", "off"}); err != nil {
		log.Fatalf("Parse: %v", err)
	}

	fmt.Printf("Chose index %d (%q)\n", feature.Index(), feature.Key())
	// Output:
	// Chose index 2 ("off")
}
