package bwtesting

import (
	"fmt"
	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/defparse"
	"testing"
)


type testParseMapStruct struct {
	source string
	result map[string]interface{}
	err    error
}

func ExampleCheckTestErrResult() {
	t := &testing.T{}
	tests := map[string]testParseMapStruct{
		"map": {
			source: `{ some: "thing" }`,
			result: map[string]interface{}{
				"some": "thing",
			},
			err: nil,
		},
		"non map": {
			source: `
        type: 'map',
        keys: {
          v: {
            type: enum
            enum: qw/all err ok none/
            default: none
          }
          s: {
            type: enum
            enum: qw/none stderr stdout all/
            default: all
          }
          exitOnError: {
            type: bool
            default: false
          }
        }
      `,
			result: nil,
			err:    fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiReset>unknown word <ansiPrimaryLiteral>type"+ansi.Ansi("Reset", " at line <ansiCmd>2<ansi>, col <ansiCmd>9<ansi> (pos <ansiCmd>9<ansi>):\n<ansiDarkGreen>\n        <ansiLightRed>type<ansiReset>: 'map',\n        keys: {\n          v: {\n"))),
		},
	}
	testsToRun := tests
	for testName, test := range testsToRun {
		t.Logf(ansi.Ansi(`Header`, "Running test case <ansiPrimaryLiteral>%s"), testName)
		result, err := defparse.ParseMap(test.source)
		testTitle := fmt.Sprintf("ParseMap(%s)\n", test.source)
		CheckTestErrResult(t, err, test.err, result, test.result, testTitle)
	}
}
