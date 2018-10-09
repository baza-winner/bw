package defparse

import (
	"fmt"
	"testing"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
)

func TestParse(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"zero number": {
			In:  []interface{}{`0`},
			Out: []interface{}{0, nil},
		},
		"zero number space surrounded": {
			In:  []interface{}{` 0 `},
			Out: []interface{}{0, nil},
		},
		"int number": {
			In:  []interface{}{` 100 `},
			Out: []interface{}{100, nil},
		},
		"int number with underscore": {
			In:  []interface{}{`100_000`},
			Out: []interface{}{100000, nil},
		},
		"int number with plus sign": {
			In:  []interface{}{`+100_000`},
			Out: []interface{}{100000, nil},
		},
		"int number with minus sign": {
			In:  []interface{}{`-100_000`},
			Out: []interface{}{-100000, nil},
		},
		"float number": {
			In:  []interface{}{`1.0`},
			Out: []interface{}{1.0, nil},
		},
		"double quoted string": {
			In:  []interface{}{`"some"`},
			Out: []interface{}{"some", nil},
		},
		"double quoted string with newline inside": {
			In: []interface{}{`"so
me"`},
			Out: []interface{}{"so\nme", nil},
		},
		"double quoted string space surrounded": {
			In:  []interface{}{` "some" `},
			Out: []interface{}{"some", nil},
		},
		"double quoted string with escapes": {
			In:  []interface{}{`"so\"me\n\a\b\f\r\t\vthing"`},
			Out: []interface{}{"so\"me\n\a\b\f\r\t\vthing", nil},
		},
		"single quoted string": {
			In:  []interface{}{`'some'`},
			Out: []interface{}{"some", nil},
		},
		"single quoted string with newline inside": {
			In: []interface{}{`'so
me'`},
			Out: []interface{}{"so\nme", nil},
		},
		"single quoted string with escapes": {
			In:  []interface{}{`'so\'me'`},
			Out: []interface{}{"so'me", nil},
		},
		"true": {
			In:  []interface{}{`'true'`},
			Out: []interface{}{"true", nil},
		},
		"true space surrounded": {
			In:  []interface{}{` true `},
			Out: []interface{}{true, nil},
		},
		"false": {
			In:  []interface{}{`false`},
			Out: []interface{}{false, nil},
		},
		"empty": {
			In:  []interface{}{``},
			Out: []interface{}{nil, nil},
		},
		"array": {
			In:  []interface{}{`[ 0 'so\'me', "so\"me" ]`},
			Out: []interface{}{[]interface{}{0, "so'me", "so\"me"}, nil},
		},
		"qw": {
			In: []interface{}{`[
          qw/ one two tree /
          qw|one two tree|
          qw#one two tree#
          qw[one two tree]
          qw<one two tree>
          qw(one two tree)
          qw{one two tree}
        ]`,
			},
			Out: []interface{}{
				[]interface{}{
					"one", "two", "tree",
					"one", "two", "tree",
					"one", "two", "tree",
					"one", "two", "tree",
					"one", "two", "tree",
					"one", "two", "tree",
					"one", "two", "tree",
				},
				nil},
		},
		"map": {
			In: []interface{}{`{
            some => [ 0, 100_000, 5_000.5, -3.14 ],
            thing: true
            'go\'od' "str\ning"
            "go\"od" nil,
          }`,
			},
			Out: []interface{}{
				map[string]interface{}{
					"some":   []interface{}{0, 100000, 5000.5, -3.14},
					"thing":  true,
					"go'od":  "str\ning",
					"go\"od": nil,
				},
				nil},
		},
		"failedToGetNumberError": {
			In: []interface{}{`{ someBigNumber: 1_000_000_000_000_000_000_000 }`},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"failed to get number from string <ansiPrimaryLiteral>1_000_000_000_000_000_000_000" +
						ansi.Ansi("Reset", " at pos <ansiCmd>17<ansi>: <ansiDarkGreen>{ someBigNumber: <ansiLightRed>1_000_000_000_000_000_000_000<ansiReset> }\n"),
				),
			},
		},
		"unexpectedWordError": {
			In: []interface{}{` qw/abc def/  `},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected word <ansiPrimaryLiteral>qw" +
						ansi.Ansi("Reset", " at pos <ansiCmd>1<ansi>: <ansiDarkGreen> <ansiLightRed>qw<ansiReset>/abc def/  \n"),
				),
			},
		},

		"unknownWordError": {
			In: []interface{}{`
[
  qw/one two three/
  def
  qw/
    four
    five
    six
 /
]`,
			},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unknown word <ansiPrimaryLiteral>def" +
						ansi.Ansi("Reset", " at line <ansiCmd>4<ansi>, col <ansiCmd>3<ansi> (pos <ansiCmd>25<ansi>):\n<ansiDarkGreen>[\n  qw/one two three/\n  <ansiLightRed>def<ansiReset>\n  qw/\n    four\n"),
				),
			},
		},
		"unexpectedCharError": {
			In: []interface{}{`
[
  1000,
  true
  *
  'value',
]`,
			},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected char <ansiPrimaryLiteral>'*'<ansiReset> (charCode: 42, pfa.state: expectValueOrSpace.orArrayItemSeparator)" +
						ansi.Ansi("Reset", " at line <ansiCmd>5<ansi>, col <ansiCmd>3<ansi> (pos <ansiCmd>20<ansi>):\n<ansiDarkGreen>  1000,\n  true\n  <ansiLightRed>*<ansiReset>\n  'value',\n]\n"),
				),
			},
		},
		"unexpectedCharError(EOF)": {
			In: []interface{}{` [1000, true 'value', `},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected end of string (pfa.state: expectValueOrSpace)" +
						ansi.Ansi("Reset", " at pos <ansiCmd>22<ansi>: <ansiDarkGreen> [1000, true 'value', \n"),
				),
			},
		},
		"non map": {

			In: []interface{}{`
        type: 'map',
        keys: {
          v: {
            type: enum
            enum: qw/all Err ok none/
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
			},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unknown word <ansiPrimaryLiteral>type" +
						ansi.Ansi("Reset", " at line <ansiCmd>2<ansi>, col <ansiCmd>9<ansi> (pos <ansiCmd>9<ansi>):\n<ansiDarkGreen>\n        <ansiLightRed>type<ansiReset>: 'map',\n        keys: {\n          v: {\n"),
				),
			},
		},
		"{some: true}": {
			In: []interface{}{`{some: true}`},
			Out: []interface{}{
				map[string]interface{}{
					"some": true,
				},
				nil},
		},
		"[true]": {
			In:  []interface{}{`[true]`},
			Out: []interface{}{[]interface{}{true}, nil},
		},
		"[qw/one two three/]": {
			In:  []interface{}{`[qw/one two three/]`},
			Out: []interface{}{[]interface{}{"one", "two", "three"}, nil},
		},
		"[ Bool String Int Number Map Array ArrayOf ]": {
			In: []interface{}{`[ Bool String Int Number Map Array ArrayOf ]`},
			Out: []interface{}{
				[]interface{}{"Bool", "String", "Int", "Number", "Map", "Array", "ArrayOf"},
				nil,
			},
		},
		"qw/ Bool String Int Number Map Array ArrayOf /": {
			In: []interface{}{`[ Bool String Int Number Map Array ArrayOf ]`},
			Out: []interface{}{
				[]interface{}{"Bool", "String", "Int", "Number", "Map", "Array", "ArrayOf"},
				nil,
			},
		},
	}

	testsToRun := tests
	bwmap.CropMap(testsToRun)
	// bwmap.CropMap(testsToRun, "[qw/one two three/]")
	bwtesting.BwRunTests(t, testsToRun, Parse)
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
	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
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
	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
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
	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
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
	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
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
	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
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
	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
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
	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
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

func TestParseMap(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"map": {
			In: []interface{}{`{ some: "thing" }`},
			Out: []interface{}{
				map[string]interface{}{"some": "thing"},
				nil,
			},
		},
		"non map": {
			In: []interface{}{`
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
      `},
			Out: []interface{}{
				map[string]interface{}(nil),
				bwerror.Error("unknown word <ansiPrimaryLiteral>type" +
					ansi.Ansi("Reset", " at line <ansiCmd>2<ansi>, col <ansiCmd>9<ansi> (pos <ansiCmd>9<ansi>):\n<ansiDarkGreen>\n        <ansiLightRed>type<ansiReset>: 'map',\n        keys: {\n          v: {\n"),
				),
			},
		},
	}
	testsToRun := tests
	bwmap.CropMap(testsToRun)
	// bwmap.CropMap(testsToRun, "[qw/one two three/]")
	bwtesting.BwRunTests(t, testsToRun, ParseMap)
}
