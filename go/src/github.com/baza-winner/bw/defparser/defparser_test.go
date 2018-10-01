package defparser

import (
	"encoding/json"
	"fmt"
	"github.com/baza-winner/bw/ansi"
	"reflect"
	"testing"
)

type testParseStruct struct {
	source string
	result interface{}
	err    error
}

func TestParse(t *testing.T) {

	tests := map[string]testParseStruct{
		"zero number": {
			source: `0`,
			result: 0,
			err:    nil,
		},
		"zero number space surrounded": {
			source: ` 0 `,
			result: 0,
			err:    nil,
		},
		"int number": {
			source: `100`,
			result: 100,
			err:    nil,
		},
		"int number with underscore": {
			source: `100_000`,
			result: 100000,
			err:    nil,
		},
		"int number with plus sign": {
			source: `+100_000`,
			result: 100000,
			err:    nil,
		},
		"int number with minus sign": {
			source: `-100_000`,
			result: -100000,
			err:    nil,
		},
		"float number": {
			source: `1.0`,
			result: 1.0,
			err:    nil,
		},
		"double quoted string": {
			source: `"some"`,
			result: "some",
			err:    nil,
		},
		"double quoted string with newline inside": {
			source: `"so
me"`,
			result: "so\nme",
			err:    nil,
		},
		"double quoted string space surrounded": {
			source: ` "some" `,
			result: "some",
			err:    nil,
		},
		"double quoted string with escapes": {
			source: `"so\"me\n\a\b\f\r\t\vthing"`,
			result: "so\"me\n\a\b\f\r\t\vthing",
			err:    nil,
		},
		"single quoted string": {
			source: `'some'`,
			result: "some",
			err:    nil,
		},
		"single quoted string with newline inside": {
			source: `'so
me'`,
			result: "so\nme",
			err:    nil,
		},
		"single quoted string with escapes": {
			source: `'so\'me'`,
			result: "so'me",
			err:    nil,
		},
		"true": {
			source: `true'`,
			result: true,
			err:    nil,
		},
		"true space surrounded": {
			source: ` true `,
			result: true,
			err:    nil,
		},
		"false": {
			source: `false`,
			result: false,
			err:    nil,
		},
		"empty": {
			source: ``,
			result: nil,
			err:    nil,
		},
		"array": {
			source: `[ 0 'so\'me', "so\"me" ]`,
			result: []interface{}{0, "so'me", "so\"me"},
			err:    nil,
		},
		"qw": {
			source: `[
				qw/one two tree/
				qw|one two tree|
				qw#one two tree#
				qw[one two tree]
				qw<one two tree>
				qw(one two tree)
				qw{one two tree}
			]`,
			result: []interface{}{
				"one", "two", "tree",
				"one", "two", "tree",
				"one", "two", "tree",
				"one", "two", "tree",
				"one", "two", "tree",
				"one", "two", "tree",
				"one", "two", "tree",
			},
			err: nil,
		},
		"map": {
			source: `{
				some => [ 0, 100_000, 5_000.5, -3.14 ],
				thing: true
				'go\'od' "str\ning"
				"go\"od" nil,
			}`,
			result: map[string]interface{}{
				"some":   []interface{}{0, 100000, 5000.5, -3.14},
				"thing":  true,
				"go'od":  "str\ning",
				"go\"od": nil,
			},
			err: nil,
		},
		// "unexpected word": {
		// 	source: `[ qw/abc def/ some ]`,
		// 	result: nil,
		// 	err:    fmt.Errorf(ansi.Ansi(`Err`, `unexpected <ansiOutline>char <ansiPrimaryLiteral>'*'<ansi> (code <ansiSecondaryLiteral>42<ansi>) at <ansiOutline>pos <ansiSecondaryLiteral>0<ansi> while <ansiOutline>state <ansiSecondaryLiteral>expectValueOrSpace`)),
		// },

		"unexpected word": {
			source: `
 qw/abc def/  `,
			result: nil,
			err:    fmt.Errorf(ansi.Ansi(`Err`, "unexpected word <ansiPrimaryLiteral>qw<ansi> at line <ansiCmd>2<ansi>, col <ansiCmd>2<ansi> (pos <ansiCmd>2<ansi>):\n<ansiOK>\n <ansiErr>qw<ansiReset>/abc def/ ")),
		},

		"unknown word": {
			source: ` [ qw/one two three/ def qw/four five six/ ] `,
			result: nil,
			err:    fmt.Errorf(ansi.Ansi(`Err`, `unknown word <ansiPrimaryLiteral>def<ansi> at pos <ansiSecondaryLiteral>3`)),
		},
		"unexpected char": {
			source: `*`,
			result: nil,
			err:    fmt.Errorf(ansi.Ansi(`Err`, `unexpected <ansiOutline>char <ansiPrimaryLiteral>'*'<ansi> (code <ansiSecondaryLiteral>42<ansi>) at <ansiOutline>pos <ansiSecondaryLiteral>0<ansi> while <ansiOutline>state <ansiSecondaryLiteral>expectValueOrSpace`)),
		},
	}

	testsToRun := tests
	// testsToRun = map[string]testParseStruct{"empty": tests["empty"], "true": tests["true"], "float number": tests["float number"], "double quoted string": tests["double quoted string"], "array": tests["array"]}
	// testsToRun = map[string]testParseStruct{"empty": tests["empty"]}
	testsToRun = map[string]testParseStruct{"unknown word": tests["unknown word"]}
	for testName, test := range testsToRun {
		t.Logf(ansi.Ansi(`Header`, "Running test case <ansiPrimaryLiteral>%s"), testName)
		result, err := Parse(test.source)
		if err != test.err {
			if err == nil || test.err == nil || err.Error() != test.err.Error() {
				t.Errorf(ansi.Ansi("", "Parse(%s)\n    => err: <ansiErr>'%v'<ansi>\n, want err: <ansiOK>'%v'"), test.source, err, test.err)
			}
		} else if !reflect.DeepEqual(result, test.result) {
			tstJson, _ := json.MarshalIndent(result, ``, `  `)
			etaJson, _ := json.MarshalIndent(test.result, ``, `  `)
			t.Errorf(ansi.Ansi("", "Parse(%s)\n    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), test.source, tstJson, etaJson)
		}
	}

}
