package bwstring

import (
	"fmt"
	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwtesting"
	"testing"
)

type testSmartQuoteStruct struct {
	ss     []string
	result string
}

func TestSmartQuote(t *testing.T) {
	tests := map[string]testSmartQuoteStruct{
		"first": {
			ss:     []string{"some", `thi"ng`, `go od`},
			result: `some "thi\"ng" "go od"`,
		},
	}
	testsToRun := tests
	for testName, test := range testsToRun {
		t.Logf(ansi.Ansi(`Header`, "Running test case <ansiPrimaryLiteral>%s"), testName)
		result := SmartQuote(test.ss...)
		testTitle := fmt.Sprintf("SmartQuote(%v)\n", test.ss)
		bwtesting.DeepEqual(t, result, test.result, testTitle)
	}
}
