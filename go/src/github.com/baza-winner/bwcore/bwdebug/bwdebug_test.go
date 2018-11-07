package bwdebug

import (
	"fmt"
	"testing"

	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
)

func TestPrint(t *testing.T) {
	varA := struct {
		S string
		N int
	}{"string value", 273}
	tests := map[string]bwtesting.Case{
		"zero number": {
			// in: []interface{}{"!HERE", "varA", varA},
			// out: Out{
			// 	result: fmt.Sprintf(
			// 		"\x1b[93;1m!HERE\x1b[0m, \x1b[32;1mgithub.com/baza-winner/bwcore/bwdebug.TestPrint.func1\x1b[38;5;243m@\x1b[97;1mbwdebug_test.go:%d: \x1b[38;5;201;1mvarA\x1b[0m: \x1b[0m\x1b[96;1m(struct { S string; N int }){S:(string)string value N:(int)273}\x1b[0m\x1b[0m\x1b[0m",
			// 		25,
			// 	),
			// 	err: bwerr.Error{},
			// },
			In: []interface{}{
				uint(5),
				"!HERE", "varA", varA,
			},
			Out: []interface{}{
				// 0,
				fmt.Sprintf(
					"\x1b[93;1m!HERE\x1b[0m, \x1b[0m\x1b[32;1mgithub.com/baza-winner/bwcore/bwdebug.TestPrint\x1b[38;5;243m@\x1b[97;1mbwdebug_test.go:%d\x1b[0m: \x1b[38;5;201;1mvarA\x1b[0m: \x1b[0m\x1b[96;1m(struct { S string; N int }){S:(string)string value N:(int)273}\x1b[0m",
					42,
				),
				nil,
			},
		},
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, ansiString, tests)

	// type Out struct {
	// 	result string
	// 	err    bwerr.Error
	// }
	// varA := struct {
	// 	S string
	// 	N int
	// }{"string value", 273}
	// type Test struct {
	// 	in  []interface{}
	// 	out Out
	// }
	// runTest := func(test Test) {
	// 	got, err := ansiString(0, test.in...)
	// 	tstErrStr := bwerr.FmtStringOf(err)
	// 	if test.out.err.Ansi != tstErrStr {
	// 		t.Errorf(ansi.String(fmt.Sprintf(
	// 			"From(%#v)\n  err: <ansiErr>%q<ansi>\n want: <ansiOK>%q<ansi>\ntst: %s\neta: %s",
	// 			test.in, tstErrStr, test.out.err.Ansi, tstErrStr, test.out.err.Ansi,
	// 		)))
	// 	}
	// 	if got != test.out.result {
	// 		t.Errorf(ansi.String(fmt.Sprintf(
	// 			"From(%#v)\n    => <ansiErr>%q<ansi>\n, want <ansiOK>%q<ansi>\ntst: %s\neta: %s",
	// 			test.in, got, test.out.result, got, test.out.result,
	// 		)))
	// 	}
	// }
	// tests := []Test{
	// 	{
	// 		in: []interface{}{"!HERE", "varA", varA},
	// 		out: Out{
	// 			result: fmt.Sprintf(
	// 				"\x1b[93;1m!HERE\x1b[0m, \x1b[32;1mgithub.com/baza-winner/bwcore/bwdebug.TestPrint.func1\x1b[38;5;243m@\x1b[97;1mbwdebug_test.go:%d: \x1b[38;5;201;1mvarA\x1b[0m: \x1b[0m\x1b[96;1m(struct { S string; N int }){S:(string)string value N:(int)273}\x1b[0m\x1b[0m\x1b[0m",
	// 				25,
	// 			),
	// 			err: bwerr.Error{},
	// 		},
	// 	},
	// 	{
	// 		in: []interface{}{"!HERE", varA},
	// 		out: Out{
	// 			result: "",
	// 			err:    bwerr.Error{Ansi: "\x1b[38;5;201;1margs\x1b[38;5;252;1m.2\x1b[0m (\x1b[96;1m(struct { S string; N int }){S:(string)string value N:(int)273}\x1b[0m) must be \x1b[97;1mstring\x1b[0m"},
	// 		},
	// 	},
	// 	{
	// 		in: []interface{}{"!HERE", "", varA},
	// 		out: Out{
	// 			result: "",
	// 			err:    bwerr.Error{Ansi: "\x1b[38;5;201;1margs\x1b[38;5;252;1m.2\x1b[0m must be \x1b[97;1mnon empty string\x1b[0m"},
	// 		},
	// 	},
	// 	{
	// 		in: []interface{}{"!HERE", "varA", varA, "varB"},
	// 		out: Out{
	// 			result: "",
	// 			err:    bwerr.Error{Ansi: "expects val for \x1b[38;5;201;1mvarB\x1b[0m\x1b[0m"},
	// 		},
	// 	},
	// }
	// for _, test := range tests {
	// 	runTest(test)
	// }
}
