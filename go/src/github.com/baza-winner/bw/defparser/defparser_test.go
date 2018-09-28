package defparser

import (
	"encoding/json"
	// "fmt"
	"github.com/baza-winner/bw/ansi"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := map[string]struct {
		source string
		result interface{}
		err    error
	}{
		"zero number": {
			source: `0`,
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
		"double quoted string with escapes": {
			source: `"so\"me\nthing"`,
			result: "so\"me\nthing",
			err:    nil,
		},
		"single quoted string": {
			source: `'some'`,
			result: "some",
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
		"false": {
			source: `false'`,
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
		"qw/": {
			source: `[ qw/one two tree/ ]`,
			result: []interface{}{"one", "two", "tree"},
			err:    nil,
		},
		// "qw|": {
		// 	source: `[ qw|one two tree| ]`,
		// 	result: []interface{}{"one", "two", "tree"},
		// 	err:    nil,
		// },
		// "qw#": {
		// 	source: `[ qw#one two tree# ]`,
		// 	result: []interface{}{"one", "two", "tree"},
		// 	err:    nil,
		// },
		// "qw[": {
		// 	source: `[ qw[one two tree] ]`,
		// 	result: []interface{}{"one", "two", "tree"},
		// 	err:    nil,
		// },
		// "qw<": {
		// 	source: `[ qw<one two tree> ]`,
		// 	result: []interface{}{"one", "two", "tree"},
		// 	err:    nil,
		// },
		// "qw(": {
		// 	source: `[ qw(one two tree) ]`,
		// 	result: []interface{}{"one", "two", "tree"},
		// 	err:    nil,
		// },
		// "qw{": {
		// 	source: `[ qw{one two tree} ]`,
		// 	result: []interface{}{"one", "two", "tree"},
		// 	err:    nil,
		// },
		"map": {
			source: `{
				some => 0,
				thing: true
				'go\'od' "str\ning"
				"go\"od" nil
			}`,
			result: map[string]interface{}{
				"some":   0,
				"thing":  true,
				"go'od":  "str\ning",
				"go\"od": nil,
			},
			err: nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		result, err := Parse(test.source)
		if err != nil {
			if err != test.err {
				// t.Fatalf("Parse(%s) => err %v, want: %v", test.source, err, test.err)
				t.Errorf(ansi.Ansi("", "Parse(%s)\n    => err:<ansiErr> %v<ansi>\n, want err:<ansiOK>%v"), test.source, err, test.err)
			}
		} else if !reflect.DeepEqual(result, test.result) {
			// fmt.Printf("%v %v\n", result, test.result)
			// fmt.Println(reflect.TypeOf(result), reflect.TypeOf(test.result))
			tstJson, _ := json.MarshalIndent(result, ``, `  `)
			etaJson, _ := json.MarshalIndent(test.result, ``, `  `)
			// t.Fatalf("Parse(%s) => %s, want: %s", test.source, tstJson, etaJson)
			t.Errorf(ansi.Ansi("", "Parse(%s)\n    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), test.source, tstJson, etaJson)
		}
	}

}
