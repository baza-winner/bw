package bwstring

import (
	"github.com/baza-winner/bwcore/bwtesting"
	"testing"
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
	bwtesting.BwRunTests(t, testsToRun, SmartQuote)
}
