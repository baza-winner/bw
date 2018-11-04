package bwos

import (
	"fmt"
	"os"
	"testing"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
)

func ExampleShortenFileSpec() {
	fmt.Printf(`%q`, ShortenFileSpec(os.Getenv(`HOME`)+`/bw`))
	// Output: "~/bw"
}

func ExampleShortenFileSpec_2() {
	fmt.Printf(`%q`, ShortenFileSpec(`/lib/bw`))
	// Output: "/lib/bw"
}

func TestExitMsg(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"with newline": {
			In: []interface{}{
				bw.A{Fmt: "some exit msg\n"},
			},
			Out: []interface{}{
				ansi.String("some exit msg\n"),
			},
		},
		"without newline": {
			In: []interface{}{
				bw.A{Fmt: "some exit msg"},
			},
			Out: []interface{}{
				ansi.String("some exit msg") + "\n",
			},
		},
	}
	bwmap.CropMap(tests)
	bwtesting.BwRunTests(t, exitMsg, tests)
}
