package defparse

import (
	"fmt"
	"testing"

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
            "go\"od\\" nil,
          }`,
			},
			Out: []interface{}{
				map[string]interface{}{
					"some":     []interface{}{0, 100000, 5000.5, -3.14},
					"thing":    true,
					"go'od":    "str\ning",
					"go\"od\\": nil,
				},
				nil},
		},
		"failedToGetNumberError": {
			In: []interface{}{`{ someBigNumber: 1_000_000_000_000_000_000_000 }`},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"failed to get number from string <ansiPrimaryLiteral>1_000_000_000_000_000_000_000 at pos <ansiCmd>17<ansi>: <ansiDarkGreen>{ someBigNumber: <ansiLightRed>1_000_000_000_000_000_000_000<ansiReset> }\n",
				),
			},
		},
		"unexpectedWordError": {
			In: []interface{}{`[ Bool Something String ]`},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unknown word <ansiPrimaryLiteral>Something at pos <ansiCmd>7<ansi>: <ansiDarkGreen>[ Bool <ansiLightRed>Something<ansiReset> String ]\n",
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
					"unknown word <ansiPrimaryLiteral>def at line <ansiCmd>4<ansi>, col <ansiCmd>3<ansi> (pos <ansiCmd>25<ansi>):\n<ansiDarkGreen>[\n  qw/one two three/\n  <ansiLightRed>def<ansiReset>\n  qw/\n    four\n",
				),
			},
		},
		"unexpectedRuneError": {
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
					"unexpected char <ansiPrimaryLiteral>'*'<ansiReset> (charCode: 42, pfa.state: expectValueOrSpace.orArrayItemSeparator) at line <ansiCmd>5<ansi>, col <ansiCmd>3<ansi> (pos <ansiCmd>20<ansi>):\n<ansiDarkGreen>  1000,\n  true\n  <ansiLightRed>*<ansiReset>\n  'value',\n]\n",
				),
			},
		},
		"unexpectedRuneError(EOF)": {
			In: []interface{}{` [1000, true 'value', `},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected end of string (pfa.state: expectValueOrSpace) at pos <ansiCmd>22<ansi>: <ansiDarkGreen> [1000, true 'value', \n",
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
					"unknown word <ansiPrimaryLiteral>type at line <ansiCmd>2<ansi>, col <ansiCmd>9<ansi> (pos <ansiCmd>9<ansi>):\n<ansiDarkGreen>\n        <ansiLightRed>type<ansiReset>: 'map',\n        keys: {\n          v: {\n",
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
		"qw/Bool String Int Number Map Array ArrayOf/": {
			In: []interface{}{`qw/Bool String Int Number Map Array ArrayOf/`},
			Out: []interface{}{
				[]interface{}{"Bool", "String", "Int", "Number", "Map", "Array", "ArrayOf"},
				nil,
			},
		},
		"< bool string int number map array arrayof >": {
			In: []interface{}{`< bool string int number map array arrayof >`},
			Out: []interface{}{
				[]interface{}{"bool", "string", "int", "number", "map", "array", "arrayof"},
				nil,
			},
		},
		"[ 'bool' <string int number map array> 'arrayof' ]": {
			In: []interface{}{`[ 'bool' <string int number map array> 'arrayof' ]`},
			Out: []interface{}{
				[]interface{}{"bool", "string", "int", "number", "map", "array", "arrayof"},
				nil,
			},
		},
		"_expectEOF && non space fa.curr.runePtr": {
			In: []interface{}{`Map Bool`},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected char <ansiPrimaryLiteral>'B'<ansiReset> (charCode: 66, pfa.state: expectEOF.orSpace) at pos <ansiCmd>4<ansi>: <ansiDarkGreen>Map <ansiLightRed>B<ansiReset>ool\n",
				),
			},
		},
		"_expectRocket && fa.curr.runePtr != >": {
			In: []interface{}{`{ key =Bool`},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected char <ansiPrimaryLiteral>'B'<ansiReset> (charCode: 66, pfa.state: expectRocket) at pos <ansiCmd>7<ansi>: <ansiDarkGreen>{ key =<ansiLightRed>B<ansiReset>ool\n",
				),
			},
		},
		"_expectSpaceOrMapKey && fa.curr.runePtr == EOF": {
			In: []interface{}{`{ `},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected end of string (pfa.state: expectSpaceOrMapKey) at pos <ansiCmd>2<ansi>: <ansiDarkGreen>{ \n",
				),
			},
		},
		"_expectSpaceOrMapKey && fa.curr.runePtr == unexpected char": {
			In: []interface{}{`{ ,`},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected char <ansiPrimaryLiteral>','<ansiReset> (charCode: 44, pfa.state: expectSpaceOrMapKey) at pos <ansiCmd>2<ansi>: <ansiDarkGreen>{ <ansiLightRed>,<ansiReset>\n",
				),
			},
		},
		"_expectMapKey && fa.curr.runePtr == EOF": {
			In: []interface{}{`{ key`},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected end of string (pfa.state: %s) at pos <ansiCmd>%d<ansi>: <ansiDarkGreen>%s\n",
					"expectValueOrSpace.orMapKeySeparator", 5, "{ key",
				),
			},
		},
		"_expectEndOfQwItem && fa.curr.runePtr == EOF": {
			In: []interface{}{`<some`},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected end of string (pfa.state: expectEndOfQwItem) at pos <ansiCmd>5<ansi>: <ansiDarkGreen><some\n",
				),
			},
		},
		"expectDigit && fa.curr.runePtr == EOF": {
			In: []interface{}{`-`},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected end of string (pfa.state: %s) at pos <ansiCmd>%d<ansi>: <ansiDarkGreen>%s\n",
					expectDigit, 1, "-",
				),
			},
		},
		"expectContentOf && fa.curr.runePtr == EOF": {
			In: []interface{}{`"`},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected end of string (pfa.state: %s) at pos <ansiCmd>%d<ansi>: <ansiDarkGreen>%s\n",
					"expectContentOf.stringToken", 1, "\"",
				),
			},
		},
		"_expectEscapedContentOf && fa.curr.runePtr == EOF": {
			In: []interface{}{`"\`},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected end of string (pfa.state: %s) at pos <ansiCmd>%d<ansi>: <ansiDarkGreen>%s\n",
					"expectEscapedContentOf.stringToken", 2, "\"\\",
				),
			},
		},
		"_expectEscapedContentOf && fa.curr.runePtr == unexpected char": {
			In: []interface{}{`"\j`},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected char <ansiPrimaryLiteral>%q<ansiReset> (charCode: %[1]d, pfa.state: %s) at pos <ansiCmd>%d<ansi>: <ansiDarkGreen>%s<ansiLightRed>%s<ansiReset>\n",
					'j', "expectEscapedContentOf.stringToken", 2, "\"\\", "j",
				),
			},
		},
		"_expectSpaceOrQwItemOrDelimiter && fa.curr.runePtr == EOF": {
			In: []interface{}{`<some `},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected end of string (pfa.state: %s) at pos <ansiCmd>%d<ansi>: <ansiDarkGreen>%s\n",
					expectSpaceOrQwItemOrDelimiter, 6, "<some ",
				),
			},
		},
		"_parseStackItemNumber int64Val": {
			In: []interface{}{`8_589_934_591`},
			Out: []interface{}{
				int64(8589934591),
				nil,
			},
		},
		"qw && fa.curr.runePtr == EOF": {
			In: []interface{}{`qw`},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected end of string (pfa.state: %s) at pos <ansiCmd>%d<ansi>: <ansiDarkGreen>%s\n",
					expectWord, 2, "qw",
				),
			},
		},
		"qw ": {
			In: []interface{}{`qw `},
			Out: []interface{}{
				nil,
				bwerror.Error(
					"unexpected char <ansiPrimaryLiteral>%q<ansiReset> (charCode: %[1]d, pfa.state: %s) at pos <ansiCmd>%d<ansi>: <ansiDarkGreen>%s<ansiLightRed>%s<ansiReset>\n",
					' ', expectWord, 2, "qw", " ",
				),
			},
		},
		// "key starts with underscore": {
		// 	In: []interface{}{` { _key: some } `},
		// 	Out: []interface{}{
		// 		nil,
		// 		bwerror.Error(
		// 			"unexpected char <ansiPrimaryLiteral>%q<ansiReset> (charCode: %[1]d, pfa.state: %s) at pos <ansiCmd>%d<ansi>: <ansiDarkGreen>%s<ansiLightRed>%s<ansiReset>\n",
		// 			' ', expectSpaceOrQwItemOrDelimiter, 2, "qw", " ",
		// 		),
		// 	},
		// },
	}

	testsToRun := tests
	bwmap.CropMap(testsToRun)
	// bwmap.CropMap(testsToRun, "qw ")
	// bwmap.CropMap(testsToRun, "qw && fa.curr.runePtr == EOF")
	bwtesting.BwRunTests(t, Parse, testsToRun)
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

func ExampleParse_7() {
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
        key5 [
          qw/one two three/
          qw{four five six}
          qw< seven eight nine >
        ]
        key6 qw( ten eleven twelve )
        key7 qw# thirteen fourteen #
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
	//         "three",
	//         "four",
	//         "five",
	//         "six",
	//         "seven",
	//         "eight",
	//         "nine"
	//       ],
	//       "key6": [
	//         "ten",
	//         "eleven",
	//         "twelve"
	//       ],
	//       "key7": [
	//         "thirteen",
	//         "fourteen"
	//       ]
	//     },
	//     "keyOfNil": null,
	//     "keyOfNull": null
	//   }
	// ]
}

func ExampleParse_8() {
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
        key5 [
          <one two three>
          'four' <five > 'six'
          < seven eight nine >
        ]
        key6 < ten eleven twelve>
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
	//         "three",
	//         "four",
	//         "five",
	//         "six",
	//         "seven",
	//         "eight",
	//         "nine"
	//       ],
	//       "key6": [
	//         "ten",
	//         "eleven",
	//         "twelve"
	//       ]
	//     },
	//     "keyOfNil": null,
	//     "keyOfNull": null
	//   }
	// ]
}

func ExampleParse_9() {
	result, err := Parse(`{
    type Map
    keys {
      v {
        type String
        enum <all err ok none>
        default 'none'
      }
      s {
        type String
        enum <none stderr stdout all>
        default 'all'
      }
    }
  }`)
	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
	// Output:
	// err: <nil>
	// result: {
	//   "keys": {
	//     "s": {
	//       "default": "all",
	//       "enum": [
	//         "none",
	//         "stderr",
	//         "stdout",
	//         "all"
	//       ],
	//       "type": "String"
	//     },
	//     "v": {
	//       "default": "none",
	//       "enum": [
	//         "all",
	//         "err",
	//         "ok",
	//         "none"
	//       ],
	//       "type": "String"
	//     }
	//   },
	//   "type": "Map"
	// }
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
				bwerror.Error("unknown word <ansiPrimaryLiteral>type at line <ansiCmd>2<ansi>, col <ansiCmd>9<ansi> (pos <ansiCmd>9<ansi>):\n<ansiDarkGreen>\n        <ansiLightRed>type<ansiReset>: 'map',\n        keys: {\n          v: {\n"),
			},
		},
	}
	testsToRun := tests
	bwmap.CropMap(testsToRun)
	// bwmap.CropMap(testsToRun, "[qw/one two three/]")
	bwtesting.BwRunTests(t, ParseMap, testsToRun)
}
