package defparse

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
		"failedToGetNumberError": {
			source: `{ someBigNumber: 1_000_000_000_000_000_000_000 }`,
			result: nil,
			err:    fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiReset>failed to get number from string <ansiPrimaryLiteral>1_000_000_000_000_000_000_000"+ansi.Ansi("Reset", " at pos <ansiCmd>17<ansi>: <ansiDarkGreen>{ someBigNumber: <ansiLightRed>1_000_000_000_000_000_000_000<ansiReset> }\n"))),
		},
		"unexpectedWordError": {
			source: ` qw/abc def/  `,
			result: nil,
			err:    fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiReset>unexpected word <ansiPrimaryLiteral>qw"+ansi.Ansi("Reset", " at pos <ansiCmd>1<ansi>: <ansiDarkGreen> <ansiLightRed>qw<ansiReset>/abc def/  \n"))),
		},

		"unknownWordError": {
			source: `
[
  qw/one two three/
  def
  qw/
    four
    five
    six
 /
]`,
			result: nil,
			err:    fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiReset>unknown word <ansiPrimaryLiteral>def"+ansi.Ansi("Reset", " at line <ansiCmd>4<ansi>, col <ansiCmd>3<ansi> (pos <ansiCmd>25<ansi>):\n<ansiDarkGreen>[\n  qw/one two three/\n  <ansiLightRed>def<ansiReset>\n  qw/\n    four\n"))),
		},
		"unexpectedCharError": {
			source: `
[
  1000,
  true
  *
  'value',
]`,
			result: nil,
			err:    fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiReset>unexpected char <ansiPrimaryLiteral>'*'<ansiReset> (charCode: 42, pfa.state: expectValueOrSpace.orArrayItemSeparator)"+ansi.Ansi("Reset", " at line <ansiCmd>5<ansi>, col <ansiCmd>3<ansi> (pos <ansiCmd>20<ansi>):\n<ansiDarkGreen>  1000,\n  true\n  <ansiLightRed>*<ansiReset>\n  'value',\n]\n"))),
		},
		"unexpectedCharError(EOF)": {
			source: ` [1000, true 'value', `,
			result: nil,
			err:    fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiReset>unexpected end of string (pfa.state: expectValueOrSpace)"+ansi.Ansi("Reset", " at pos <ansiCmd>22<ansi>: <ansiDarkGreen> [1000, true 'value', \n"))),
		},
	}

	testsToRun := tests
	// testsToRun = map[string]testParseStruct{"empty": tests["empty"], "true": tests["true"], "float number": tests["float number"], "double quoted string": tests["double quoted string"], "array": tests["array"]}
	// testsToRun = map[string]testParseStruct{"empty": tests["empty"]}
	// testsToRun = map[string]testParseStruct{"unexpectedCharError(EOF)": tests["unexpectedCharError(EOF)"]}
	for testName, test := range testsToRun {
		t.Logf(ansi.Ansi(`Header`, "Running test case <ansiPrimaryLiteral>%s"), testName)
		result, err := Parse(test.source)
		if err != test.err {
			if err == nil || test.err == nil || err.Error() != test.err.Error() {
				t.Errorf(ansi.Ansi("", "Parse(%s)\n    => err: <ansiErr>'%v'<ansi>\n, want err: <ansiOK>'%v'"), test.source, err, test.err)
				fmt.Printf("eta: %q\ntst: %q\n", err, test.err)
			}
		} else if !reflect.DeepEqual(result, test.result) {
			tstJson, _ := json.MarshalIndent(result, ``, `  `)
			etaJson, _ := json.MarshalIndent(test.result, ``, `  `)
			t.Errorf(ansi.Ansi("", "Parse(%s)\n    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), test.source, tstJson, etaJson)
		}
	}

}

func ExampleParse_1() {
	result, err := Parse(`[
   {
     "keyOfStringValue": "stringValue",
     "keyOfBoolValue": false,
     "keyOfNumberValue": 12345000.678001
   }, {
      "keyOfNull": null,
      "keyOfNil": nil,
      "keyOfArrayValue": [ "stringValue", true, 876.54321 ],
      "keyOfMapValue": {
        "key1": "value1",
        "key2": true,
        "key3": -3.14,
        "key4": nil,
        "key5": [ "one", "two", "three" ]
     }
   }
  ]`)
	resultJson, _ := json.MarshalIndent(result, ``, `  `)
	fmt.Printf("err: %v\nresult: %s", err, resultJson)
	// Output:
	// err: <nil>
	// result: [
	//   {
	//     "keyOfBoolValue": false,
	//     "keyOfNumberValue": 12345000.678001,
	//     "keyOfStringValue": "stringValue"
	//   },
	//   {
	//     "keyOfArrayValue": [
	//       "stringValue",
	//       true,
	//       876.54321
	//     ],
	//     "keyOfMapValue": {
	//       "key1": "value1",
	//       "key2": true,
	//       "key3": -3.14,
	//       "key4": null,
	//       "key5": [
	//         "one",
	//         "two",
	//         "three"
	//       ]
	//     },
	//     "keyOfNil": null,
	//     "keyOfNull": null
	//   }
	// ]
}

func ExampleParse_2() {
	result, err := Parse(`[
   {
     "keyOfStringValue": "stringValue"
     "keyOfBoolValue": false
     "keyOfNumberValue": 12345000.678001
   } {
      "keyOfNull": null
      "keyOfNil": nil
      "keyOfArrayValue": [ "stringValue", true, 876.54321 ]
      "keyOfMapValue": {
        "key1": "value1"
        "key2": true
        "key3": -3.14
        "key4": nil
        "key5": [ "one", "two", "three" ]
     }
   }
  ]`)
	resultJson, _ := json.MarshalIndent(result, ``, `  `)
	fmt.Printf("err: %v\nresult: %s", err, resultJson)
	// Output:
	// err: <nil>
	// result: [
	//   {
	//     "keyOfBoolValue": false,
	//     "keyOfNumberValue": 12345000.678001,
	//     "keyOfStringValue": "stringValue"
	//   },
	//   {
	//     "keyOfArrayValue": [
	//       "stringValue",
	//       true,
	//       876.54321
	//     ],
	//     "keyOfMapValue": {
	//       "key1": "value1",
	//       "key2": true,
	//       "key3": -3.14,
	//       "key4": null,
	//       "key5": [
	//         "one",
	//         "two",
	//         "three"
	//       ]
	//     },
	//     "keyOfNil": null,
	//     "keyOfNull": null
	//   }
	// ]
}

func ExampleParse_3() {
	result, err := Parse(`[
   {
     keyOfStringValue: "stringValue"
     keyOfBoolValue: false
     keyOfNumberValue: 12345000.678001
   } {
      keyOfNull: null
      keyOfNil: nil
      keyOfArrayValue: [ "stringValue", true, 876.54321 ]
      keyOfMapValue: {
        key1: "value1"
        key2: true
        key3: -3.14
        key4: nil
        key5: [ "one", "two", "three" ]
     }
   }
  ]`)
	resultJson, _ := json.MarshalIndent(result, ``, `  `)
	fmt.Printf("err: %v\nresult: %s", err, resultJson)
	// Output:
	// err: <nil>
	// result: [
	//   {
	//     "keyOfBoolValue": false,
	//     "keyOfNumberValue": 12345000.678001,
	//     "keyOfStringValue": "stringValue"
	//   },
	//   {
	//     "keyOfArrayValue": [
	//       "stringValue",
	//       true,
	//       876.54321
	//     ],
	//     "keyOfMapValue": {
	//       "key1": "value1",
	//       "key2": true,
	//       "key3": -3.14,
	//       "key4": null,
	//       "key5": [
	//         "one",
	//         "two",
	//         "three"
	//       ]
	//     },
	//     "keyOfNil": null,
	//     "keyOfNull": null
	//   }
	// ]
}

func ExampleParse_4() {
	result, err := Parse(`[
   {
     keyOfStringValue => "stringValue"
     keyOfBoolValue => false
     keyOfNumberValue => 12345000.678001
   } {
      keyOfNull => null
      keyOfNil => nil
      keyOfArrayValue => [ "stringValue", true, 876.54321 ]
      keyOfMapValue => {
        key1 => "value1"
        key2 => true
        key3 => -3.14
        key4 => nil
        key5 => [ "one", "two", "three" ]
     }
   }
  ]`)
	resultJson, _ := json.MarshalIndent(result, ``, `  `)
	fmt.Printf("err: %v\nresult: %s", err, resultJson)
	// Output:
	// err: <nil>
	// result: [
	//   {
	//     "keyOfBoolValue": false,
	//     "keyOfNumberValue": 12345000.678001,
	//     "keyOfStringValue": "stringValue"
	//   },
	//   {
	//     "keyOfArrayValue": [
	//       "stringValue",
	//       true,
	//       876.54321
	//     ],
	//     "keyOfMapValue": {
	//       "key1": "value1",
	//       "key2": true,
	//       "key3": -3.14,
	//       "key4": null,
	//       "key5": [
	//         "one",
	//         "two",
	//         "three"
	//       ]
	//     },
	//     "keyOfNil": null,
	//     "keyOfNull": null
	//   }
	// ]
}

func ExampleParse_5() {
	result, err := Parse(`[
   {
     keyOfStringValue "stringValue"
     keyOfBoolValue false
     keyOfNumberValue 12345000.678001
   } {
      keyOfNull null
      keyOfNil nil
      keyOfArrayValue [ "stringValue" true 876.54321 ]
      keyOfMapValue {
        key1 "value1"
        key2 true
        key3 -3.14
        key4 nil
        key5 [ "one" "two" "three" ]
     }
   }
  ]`)
	resultJson, _ := json.MarshalIndent(result, ``, `  `)
	fmt.Printf("err: %v\nresult: %s", err, resultJson)
	// Output:
	// err: <nil>
	// result: [
	//   {
	//     "keyOfBoolValue": false,
	//     "keyOfNumberValue": 12345000.678001,
	//     "keyOfStringValue": "stringValue"
	//   },
	//   {
	//     "keyOfArrayValue": [
	//       "stringValue",
	//       true,
	//       876.54321
	//     ],
	//     "keyOfMapValue": {
	//       "key1": "value1",
	//       "key2": true,
	//       "key3": -3.14,
	//       "key4": null,
	//       "key5": [
	//         "one",
	//         "two",
	//         "three"
	//       ]
	//     },
	//     "keyOfNil": null,
	//     "keyOfNull": null
	//   }
	// ]
}

func ExampleParse_6() {
	result, err := Parse(`[
   {
     keyOfStringValue "stringValue"
     keyOfBoolValue false
     keyOfNumberValue 12345000.678001
   } {
      keyOfNull null
      keyOfNil nil
      keyOfArrayValue [ "stringValue" true 876.54321 ]
      keyOfMapValue {
        key1 "value1"
        key2 true
        key3 -3.14
        key4 nil
        key5 [ qw/one two three/ ]
     }
   }
  ]`)
	resultJson, _ := json.MarshalIndent(result, ``, `  `)
	fmt.Printf("err: %v\nresult: %s", err, resultJson)
	// Output:
	// err: <nil>
	// result: [
	//   {
	//     "keyOfBoolValue": false,
	//     "keyOfNumberValue": 12345000.678001,
	//     "keyOfStringValue": "stringValue"
	//   },
	//   {
	//     "keyOfArrayValue": [
	//       "stringValue",
	//       true,
	//       876.54321
	//     ],
	//     "keyOfMapValue": {
	//       "key1": "value1",
	//       "key2": true,
	//       "key3": -3.14,
	//       "key4": null,
	//       "key5": [
	//         "one",
	//         "two",
	//         "three"
	//       ]
	//     },
	//     "keyOfNil": null,
	//     "keyOfNull": null
	//   }
	// ]
}

func ExampleParse_7() {
	result, err := Parse(`[
   {
     keyOfStringValue "stringValue"
     keyOfBoolValue false
     keyOfNumberValue 12_345_000.678_001
   } {
      keyOfNull null
      keyOfNil nil
      keyOfArrayValue [ "stringValue" true 876.543_21 ]
      keyOfMapValue {
        key1 "value1"
        key2 true
        key3 -3.14
        key4 nil
        key5 [ qw/one two three/ ]
     }
   }
  ]`)
	resultJson, _ := json.MarshalIndent(result, ``, `  `)
	fmt.Printf("err: %v\nresult: %s", err, resultJson)
	// Output:
	// err: <nil>
	// result: [
	//   {
	//     "keyOfBoolValue": false,
	//     "keyOfNumberValue": 12345000.678001,
	//     "keyOfStringValue": "stringValue"
	//   },
	//   {
	//     "keyOfArrayValue": [
	//       "stringValue",
	//       true,
	//       876.54321
	//     ],
	//     "keyOfMapValue": {
	//       "key1": "value1",
	//       "key2": true,
	//       "key3": -3.14,
	//       "key4": null,
	//       "key5": [
	//         "one",
	//         "two",
	//         "three"
	//       ]
	//     },
	//     "keyOfNil": null,
	//     "keyOfNull": null
	//   }
	// ]
}
