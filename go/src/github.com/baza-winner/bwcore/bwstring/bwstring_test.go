package bwstring

import (
	"testing"

	"github.com/baza-winner/bwcore/bwtesting"
)

func TestSmartQuote(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"first": {
			In: []interface{}{
				[]string{"some", `thi"ng`, `go od`},
			},
			Out: []interface{}{
				`some "thi\"ng" "go od"`,
			},
		},
	}
	testsToRun := tests
	bwtesting.BwRunTests(t, SmartQuote, testsToRun)
}
