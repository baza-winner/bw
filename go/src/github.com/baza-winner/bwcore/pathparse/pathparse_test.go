package pathparse

import (
	"testing"

	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
)

func TestParse(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"0": {
			In: []interface{}{`0`},
			Out: []interface{}{
				[]interface{}{0},
				nil,
			},
		},
		"[0]": {
			In: []interface{}{`0`},
			Out: []interface{}{
				[]interface{}{0},
				nil,
			},
		},
		"#0": {
			In: []interface{}{`0`},
			Out: []interface{}{
				[]interface{}{0},
				nil,
			},
		},
		"'so me'": {
			In: []interface{}{`'so me'`},
			Out: []interface{}{
				[]interface{}{"so me"},
				nil,
			},
		},
		`"so\nme"`: {
			In: []interface{}{`"so\nme"`},
			Out: []interface{}{
				[]interface{}{"so\nme"}, nil,
			},
		},
		`complex`: {
			In: []interface{}{`[131070].keys.'some \'\"\\thing'.#8_589_934_591."\a\b\f\n\r\t\v".510.31`},
			Out: []interface{}{
				[]interface{}{int32(131070), "keys", "some '\"\\thing", int64(8589934591), "\a\b\f\n\r\t\v", int16(510), int8(31)},
				nil,
			},
		},
		// "qw ": {
		// 	In: []interface{}{`qw `},
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
	// bwmap.CropMap(testsToRun, "_expectSpaceOrQwItemOrDelimiter && fa.curr.runePtr == EOF")
	// bwmap.CropMap(testsToRun, "qw ")
	bwtesting.BwRunTests(t, testsToRun, Parse)
}

// func ExampleParse_1() {
// 	result, err := Parse(`[
//    {
//      "keyOfStringValue": "stringValue",
//      "keyOfBoolValue": false,
//      "keyOfNumberValue": 12345000.678001
//    }, {
//       "keyOfNull": null,
//       "keyOfNil": nil,
//       "keyOfArrayValue": [ "stringValue", true, 876.54321 ],
//       "keyOfMapValue": {
//         "key1": "value1",
//         "key2": true,
//         "key3": -3.14,
//         "key4": nil,
//         "key5": [ "one", "two", "three" ]
//      }
//    }
//   ]`)
// 	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
// 	// Output:
// 	// err: <nil>
// 	// result: [
// 	//   {
// 	//     "keyOfBoolValue": false,
// 	//     "keyOfNumberValue": 12345000.678001,
// 	//     "keyOfStringValue": "stringValue"
// 	//   },
// 	//   {
// 	//     "keyOfArrayValue": [
// 	//       "stringValue",
// 	//       true,
// 	//       876.54321
// 	//     ],
// 	//     "keyOfMapValue": {
// 	//       "key1": "value1",
// 	//       "key2": true,
// 	//       "key3": -3.14,
// 	//       "key4": null,
// 	//       "key5": [
// 	//         "one",
// 	//         "two",
// 	//         "three"
// 	//       ]
// 	//     },
// 	//     "keyOfNil": null,
// 	//     "keyOfNull": null
// 	//   }
// 	// ]
// }

// func ExampleParse_2() {
// 	result, err := Parse(`[
//    {
//      "keyOfStringValue": "stringValue"
//      "keyOfBoolValue": false
//      "keyOfNumberValue": 12345000.678001
//    } {
//       "keyOfNull": null
//       "keyOfNil": nil
//       "keyOfArrayValue": [ "stringValue", true, 876.54321 ]
//       "keyOfMapValue": {
//         "key1": "value1"
//         "key2": true
//         "key3": -3.14
//         "key4": nil
//         "key5": [ "one", "two", "three" ]
//      }
//    }
//   ]`)
// 	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
// 	// Output:
// 	// err: <nil>
// 	// result: [
// 	//   {
// 	//     "keyOfBoolValue": false,
// 	//     "keyOfNumberValue": 12345000.678001,
// 	//     "keyOfStringValue": "stringValue"
// 	//   },
// 	//   {
// 	//     "keyOfArrayValue": [
// 	//       "stringValue",
// 	//       true,
// 	//       876.54321
// 	//     ],
// 	//     "keyOfMapValue": {
// 	//       "key1": "value1",
// 	//       "key2": true,
// 	//       "key3": -3.14,
// 	//       "key4": null,
// 	//       "key5": [
// 	//         "one",
// 	//         "two",
// 	//         "three"
// 	//       ]
// 	//     },
// 	//     "keyOfNil": null,
// 	//     "keyOfNull": null
// 	//   }
// 	// ]
// }

// func ExampleParse_3() {
// 	result, err := Parse(`[
//    {
//      keyOfStringValue: "stringValue"
//      keyOfBoolValue: false
//      keyOfNumberValue: 12345000.678001
//    } {
//       keyOfNull: null
//       keyOfNil: nil
//       keyOfArrayValue: [ "stringValue", true, 876.54321 ]
//       keyOfMapValue: {
//         key1: "value1"
//         key2: true
//         key3: -3.14
//         key4: nil
//         key5: [ "one", "two", "three" ]
//      }
//    }
//   ]`)
// 	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
// 	// Output:
// 	// err: <nil>
// 	// result: [
// 	//   {
// 	//     "keyOfBoolValue": false,
// 	//     "keyOfNumberValue": 12345000.678001,
// 	//     "keyOfStringValue": "stringValue"
// 	//   },
// 	//   {
// 	//     "keyOfArrayValue": [
// 	//       "stringValue",
// 	//       true,
// 	//       876.54321
// 	//     ],
// 	//     "keyOfMapValue": {
// 	//       "key1": "value1",
// 	//       "key2": true,
// 	//       "key3": -3.14,
// 	//       "key4": null,
// 	//       "key5": [
// 	//         "one",
// 	//         "two",
// 	//         "three"
// 	//       ]
// 	//     },
// 	//     "keyOfNil": null,
// 	//     "keyOfNull": null
// 	//   }
// 	// ]
// }

// func ExampleParse_4() {
// 	result, err := Parse(`[
//    {
//      keyOfStringValue => "stringValue"
//      keyOfBoolValue => false
//      keyOfNumberValue => 12345000.678001
//    } {
//       keyOfNull => null
//       keyOfNil => nil
//       keyOfArrayValue => [ "stringValue", true, 876.54321 ]
//       keyOfMapValue => {
//         key1 => "value1"
//         key2 => true
//         key3 => -3.14
//         key4 => nil
//         key5 => [ "one", "two", "three" ]
//      }
//    }
//   ]`)
// 	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
// 	// Output:
// 	// err: <nil>
// 	// result: [
// 	//   {
// 	//     "keyOfBoolValue": false,
// 	//     "keyOfNumberValue": 12345000.678001,
// 	//     "keyOfStringValue": "stringValue"
// 	//   },
// 	//   {
// 	//     "keyOfArrayValue": [
// 	//       "stringValue",
// 	//       true,
// 	//       876.54321
// 	//     ],
// 	//     "keyOfMapValue": {
// 	//       "key1": "value1",
// 	//       "key2": true,
// 	//       "key3": -3.14,
// 	//       "key4": null,
// 	//       "key5": [
// 	//         "one",
// 	//         "two",
// 	//         "three"
// 	//       ]
// 	//     },
// 	//     "keyOfNil": null,
// 	//     "keyOfNull": null
// 	//   }
// 	// ]
// }

// func ExampleParse_5() {
// 	result, err := Parse(`[
//    {
//      keyOfStringValue "stringValue"
//      keyOfBoolValue false
//      keyOfNumberValue 12345000.678001
//    } {
//       keyOfNull null
//       keyOfNil nil
//       keyOfArrayValue [ "stringValue" true 876.54321 ]
//       keyOfMapValue {
//         key1 "value1"
//         key2 true
//         key3 -3.14
//         key4 nil
//         key5 [ "one" "two" "three" ]
//      }
//    }
//   ]`)
// 	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
// 	// Output:
// 	// err: <nil>
// 	// result: [
// 	//   {
// 	//     "keyOfBoolValue": false,
// 	//     "keyOfNumberValue": 12345000.678001,
// 	//     "keyOfStringValue": "stringValue"
// 	//   },
// 	//   {
// 	//     "keyOfArrayValue": [
// 	//       "stringValue",
// 	//       true,
// 	//       876.54321
// 	//     ],
// 	//     "keyOfMapValue": {
// 	//       "key1": "value1",
// 	//       "key2": true,
// 	//       "key3": -3.14,
// 	//       "key4": null,
// 	//       "key5": [
// 	//         "one",
// 	//         "two",
// 	//         "three"
// 	//       ]
// 	//     },
// 	//     "keyOfNil": null,
// 	//     "keyOfNull": null
// 	//   }
// 	// ]
// }

// func ExampleParse_6() {
// 	result, err := Parse(`[
//    {
//      keyOfStringValue "stringValue"
//      keyOfBoolValue false
//      keyOfNumberValue 12_345_000.678_001
//    } {
//       keyOfNull null
//       keyOfNil nil
//       keyOfArrayValue [ "stringValue" true 876.543_21 ]
//       keyOfMapValue {
//         key1 "value1"
//         key2 true
//         key3 -3.14
//         key4 nil
//         key5 [ qw/one two three/ ]
//      }
//    }
//   ]`)
// 	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
// 	// Output:
// 	// err: <nil>
// 	// result: [
// 	//   {
// 	//     "keyOfBoolValue": false,
// 	//     "keyOfNumberValue": 12345000.678001,
// 	//     "keyOfStringValue": "stringValue"
// 	//   },
// 	//   {
// 	//     "keyOfArrayValue": [
// 	//       "stringValue",
// 	//       true,
// 	//       876.54321
// 	//     ],
// 	//     "keyOfMapValue": {
// 	//       "key1": "value1",
// 	//       "key2": true,
// 	//       "key3": -3.14,
// 	//       "key4": null,
// 	//       "key5": [
// 	//         "one",
// 	//         "two",
// 	//         "three"
// 	//       ]
// 	//     },
// 	//     "keyOfNil": null,
// 	//     "keyOfNull": null
// 	//   }
// 	// ]
// }

// func ExampleParse_7() {
// 	result, err := Parse(`[
//    {
//      keyOfStringValue "stringValue"
//      keyOfBoolValue false
//      keyOfNumberValue 12345000.678001
//    } {
//       keyOfNull null
//       keyOfNil nil
//       keyOfArrayValue [ "stringValue" true 876.54321 ]
//       keyOfMapValue {
//         key1 "value1"
//         key2 true
//         key3 -3.14
//         key4 nil
//         key5 [
//           qw/one two three/
//           qw{four five six}
//           qw< seven eight nine >
//         ]
//         key6 qw( ten eleven twelve )
//         key7 qw# thirteen fourteen #
//      }
//    }
//   ]`)
// 	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
// 	// Output:
// 	// err: <nil>
// 	// result: [
// 	//   {
// 	//     "keyOfBoolValue": false,
// 	//     "keyOfNumberValue": 12345000.678001,
// 	//     "keyOfStringValue": "stringValue"
// 	//   },
// 	//   {
// 	//     "keyOfArrayValue": [
// 	//       "stringValue",
// 	//       true,
// 	//       876.54321
// 	//     ],
// 	//     "keyOfMapValue": {
// 	//       "key1": "value1",
// 	//       "key2": true,
// 	//       "key3": -3.14,
// 	//       "key4": null,
// 	//       "key5": [
// 	//         "one",
// 	//         "two",
// 	//         "three",
// 	//         "four",
// 	//         "five",
// 	//         "six",
// 	//         "seven",
// 	//         "eight",
// 	//         "nine"
// 	//       ],
// 	//       "key6": [
// 	//         "ten",
// 	//         "eleven",
// 	//         "twelve"
// 	//       ],
// 	//       "key7": [
// 	//         "thirteen",
// 	//         "fourteen"
// 	//       ]
// 	//     },
// 	//     "keyOfNil": null,
// 	//     "keyOfNull": null
// 	//   }
// 	// ]
// }

// func ExampleParse_8() {
// 	result, err := Parse(`[
//    {
//      keyOfStringValue "stringValue"
//      keyOfBoolValue false
//      keyOfNumberValue 12345000.678001
//    } {
//       keyOfNull null
//       keyOfNil nil
//       keyOfArrayValue [ "stringValue" true 876.54321 ]
//       keyOfMapValue {
//         key1 "value1"
//         key2 true
//         key3 -3.14
//         key4 nil
//         key5 [
//           <one two three>
//           'four' <five > 'six'
//           < seven eight nine >
//         ]
//         key6 < ten eleven twelve>
//      }
//    }
//   ]`)
// 	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
// 	// Output:
// 	// err: <nil>
// 	// result: [
// 	//   {
// 	//     "keyOfBoolValue": false,
// 	//     "keyOfNumberValue": 12345000.678001,
// 	//     "keyOfStringValue": "stringValue"
// 	//   },
// 	//   {
// 	//     "keyOfArrayValue": [
// 	//       "stringValue",
// 	//       true,
// 	//       876.54321
// 	//     ],
// 	//     "keyOfMapValue": {
// 	//       "key1": "value1",
// 	//       "key2": true,
// 	//       "key3": -3.14,
// 	//       "key4": null,
// 	//       "key5": [
// 	//         "one",
// 	//         "two",
// 	//         "three",
// 	//         "four",
// 	//         "five",
// 	//         "six",
// 	//         "seven",
// 	//         "eight",
// 	//         "nine"
// 	//       ],
// 	//       "key6": [
// 	//         "ten",
// 	//         "eleven",
// 	//         "twelve"
// 	//       ]
// 	//     },
// 	//     "keyOfNil": null,
// 	//     "keyOfNull": null
// 	//   }
// 	// ]
// }

// func ExampleParse_9() {
// 	result, err := Parse(`{
//     type Map
//     keys {
//       v {
//         type String
//         enum <all err ok none>
//         default 'none'
//       }
//       s {
//         type String
//         enum <none stderr stdout all>
//         default 'all'
//       }
//     }
//   }`)
// 	fmt.Printf("err: %v\nresult: %s", err, bwjson.PrettyJson(result))
// 	// Output:
// 	// err: <nil>
// 	// result: {
// 	//   "keys": {
// 	//     "s": {
// 	//       "default": "all",
// 	//       "enum": [
// 	//         "none",
// 	//         "stderr",
// 	//         "stdout",
// 	//         "all"
// 	//       ],
// 	//       "type": "String"
// 	//     },
// 	//     "v": {
// 	//       "default": "none",
// 	//       "enum": [
// 	//         "all",
// 	//         "err",
// 	//         "ok",
// 	//         "none"
// 	//       ],
// 	//       "type": "String"
// 	//     }
// 	//   },
// 	//   "type": "Map"
// 	// }
// }

// func TestParseMap(t *testing.T) {
// 	tests := map[string]bwtesting.TestCaseStruct{
// 		"map": {
// 			In: []interface{}{`{ some: "thing" }`},
// 			Out: []interface{}{
// 				map[string]interface{}{"some": "thing"},
// 				nil,
// 			},
// 		},
// 		"non map": {
// 			In: []interface{}{`
//         type: 'map',
//         keys: {
//           v: {
//             type: enum
//             enum: qw/all err ok none/
//             default: none
//           }
//           s: {
//             type: enum
//             enum: qw/none stderr stdout all/
//             default: all
//           }
//           exitOnError: {
//             type: bool
//             default: false
//           }
//         }
//       `},
// 			Out: []interface{}{
// 				map[string]interface{}(nil),
// 				bwerror.Error("unknown word <ansiPrimaryLiteral>type at line <ansiCmd>2<ansi>, col <ansiCmd>9<ansi> (pos <ansiCmd>9<ansi>):\n<ansiDarkGreen>\n        <ansiLightRed>type<ansiReset>: 'map',\n        keys: {\n          v: {\n"),
// 			},
// 		},
// 	}
// 	testsToRun := tests
// 	bwmap.CropMap(testsToRun)
// 	// bwmap.CropMap(testsToRun, "[qw/one two three/]")
// 	bwtesting.BwRunTests(t, testsToRun, ParseMap)
// }
