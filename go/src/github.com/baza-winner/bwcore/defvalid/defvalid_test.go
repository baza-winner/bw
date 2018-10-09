package defvalid

import (
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/defparse"
	"github.com/baza-winner/bwcore/defvalid/deftype"
	"testing"
)

// ============================================================================

func TestCompileDef(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"def: nil": {
			In: []interface{}{nil},
			Out: []interface{}{
				(*Def)(nil),
				bwerror.Error("<ansiOutline>def<ansiCmd><ansi> (<ansiSecondaryLiteral>null<ansi>) has non supported value"),
			},
		},
		"def: invalid type": {
			In: []interface{}{defparse.MustParse(`false`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerror.Error("<ansiOutline>def<ansiCmd><ansi> (<ansiSecondaryLiteral>false<ansi>) is not of type <ansiPrimaryLiteral>[]string<ansi>, or <ansiPrimaryLiteral>map[string]<ansi>, or <ansiPrimaryLiteral>string"),
			},
		},
		"def: simple valid": {
			In: []interface{}{defparse.MustParse(`"bool"`)},
			Out: []interface{}{
				&Def{tp: deftype.FromArgs(deftype.Bool)},
				nil,
			},
		},
		"def: invalid deftypeItem": {
			In: []interface{}{defparse.MustParse(`[ qw/ bool int some / ]`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerror.Error("<ansiOutline>def<ansiCmd>.#2<ansi> (<ansiSecondaryLiteral>\"some\"<ansi>) has non supported value"),
			},
		},
		"def: enum": {
			In: []interface{}{defparse.MustParse(`{ type: "string", enum: [qw/one two three/]}`)},
			Out: []interface{}{
				&Def{tp: deftype.FromArgs(deftype.String), enum: bwset.StringsFromArgs("one", "two", "three")},
				nil,
			},
		},
		"def: map with keys": {
			In: []interface{}{defparse.MustParse(`{ type: "map", keys: { keyBool: ['bool'] }}`)},
			Out: []interface{}{
				&Def{
					tp: deftype.FromArgs(deftype.Map),
					keys: map[string]Def{
						"keyBool": Def{tp: deftype.FromArgs(deftype.Bool)},
					}},
				nil,
			},
		},
		"def: unexpected keys": {
			In: []interface{}{defparse.MustParse(`{ type: "map", kyes: { keyBool: ['bool'] }}`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.TestCaseStruct) error {
					return bwerror.Error(
						"<ansiOutline>def<ansiCmd><ansi> (<ansiSecondaryLiteral>" +
							bwjson.PrettyJson(test.In[0]) +
							"<ansi>) has unexpected key <ansiPrimaryLiteral>kyes",
					)
				},
			},
		},
		"def: array with arrayElem": {
			In: []interface{}{defparse.MustParse(`{ type: "array", arrayElem: 'int' }`)},
			Out: []interface{}{
				&Def{
					tp:        deftype.FromArgs(deftype.Array),
					arrayElem: &Def{tp: deftype.FromArgs(deftype.Int)},
				},
				nil,
			},
		},
		"def: array with arrayElem and Elem": {
			In: []interface{}{defparse.MustParse(`{ type: "array", arrayElem: 'int', elem: 'bool' }`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.TestCaseStruct) error {
					return bwerror.Error(
						"<ansiOutline>def<ansiCmd><ansi> (<ansiSecondaryLiteral>" +
							bwjson.PrettyJson(test.In[0]) +
							"<ansi>) has unexpected key <ansiPrimaryLiteral>elem",
					)
				},
			},
		},
		"def: minInt, maxInt": {
			In: []interface{}{defparse.MustParse(`{ type: "int", minInt: -6, maxInt: 10 }`)},
			Out: []interface{}{
				&Def{
					tp:     deftype.FromArgs(deftype.Int),
					minInt: ptrToInt64(-6),
					maxInt: ptrToInt64(10),
				},
				nil,
			},
		},
		"def: minInt > maxInt": {
			In: []interface{}{defparse.MustParse(`{ type: "int", minInt: 6, maxInt: -10 }`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.TestCaseStruct) error {
					return bwerror.Error(
						"<ansiOutline>def<ansiCmd><ansi> (<ansiSecondaryLiteral>" +
							bwjson.PrettyJson(test.In[0]) +
							"<ansi>) has conflicting keys: <ansiSecondaryLiteral>" +
							bwjson.PrettyJson(defparse.MustParse("{ minInt: 6, maxInt: -10 }")),
					)
				},
			},
		},
		"def: minNumber, maxNumber": {
			In: []interface{}{defparse.MustParse(`{ type: "number", minNumber: -6, maxNumber: 10 }`)},
			Out: []interface{}{
				&Def{
					tp:        deftype.FromArgs(deftype.Number),
					minNumber: ptrToFloat64(float64(-6)),
					maxNumber: ptrToFloat64(float64(10)),
				},
				nil,
			},
		},
		"def: default": {
			In: []interface{}{defparse.MustParse(`{ type: "bool", default: true }`)},
			Out: []interface{}{
				&Def{
					tp:         deftype.FromArgs(deftype.Bool),
					dflt:       true,
					isOptional: true,
				},
				nil,
			},
		},
		"def: default 'string' for bool": {
			In: []interface{}{defparse.MustParse(`{ type: "bool", default: "string" }`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerror.Error("<ansiOutline>def<ansiCmd>.default<ansi> (<ansiSecondaryLiteral>\"string\"<ansi>) is not of type <ansiPrimaryLiteral>Bool"),
			},
		},
		"def: isOptional": {
			In: []interface{}{defparse.MustParse(`{ type: "bool", isOptional: true }`)},
			Out: []interface{}{
				&Def{
					tp:         deftype.FromArgs(deftype.Bool),
					isOptional: true,
				},
				nil,
			},
		},
		"def: isOptional = false conflicts with dflt": {
			In: []interface{}{defparse.MustParse(`{ type: "bool", isOptional: false, default: false }`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.TestCaseStruct) error {
					return bwerror.Error(
						"<ansiOutline>def<ansiCmd><ansi> (<ansiSecondaryLiteral>" +
							bwjson.PrettyJson(test.In[0]) +
							"<ansi>) has conflicting keys: <ansiSecondaryLiteral>" +
							bwjson.PrettyJson(defparse.MustParse("{ isOptional: false, default: false }")),
					)
				},
			},
		},
		"def: arrayOf without follower": {
			In: []interface{}{defparse.MustParse(`{ type:  "arrayOf"  }`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerror.Error("<ansiOutline>def<ansiCmd>.type<ansi> (<ansiSecondaryLiteral>\"arrayOf\"<ansi>) must be followed by some type"),
			},
		},
		"simple.def: arrayOf without follower": {
			In: []interface{}{defparse.MustParse(`"arrayOf"`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerror.Error("<ansiOutline>def<ansiCmd><ansi> (<ansiSecondaryLiteral>\"arrayOf\"<ansi>) must be followed by some type"),
			},
		},
		"simple.def: arrayOf can not be combined with array": {
			In: []interface{}{defparse.MustParse(`["arrayOf", "array", "string"]`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.TestCaseStruct) error {
					return bwerror.Error(
						"<ansiOutline>def<ansiCmd><ansi> (<ansiSecondaryLiteral>" +
							bwjson.PrettyJson(test.In[0]) +
							"<ansi>) following values can not be combined: <ansiSecondaryLiteral>" +
							bwjson.PrettyJson(defparse.MustParse("[ 'array', 'arrayOf' ]")),
					)
				},
			},
		},
	}
	testsToRun := tests
	bwmap.CropMap(testsToRun)
	// bwmap.CropMap(testsToRun, "def: arrayOf without follower")
	bwtesting.BwRunTests(t, testsToRun, CompileDef)
}

// ============================================================================

func TestValidateVal(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"val: nil, simple.def: bool": {
			In: []interface{}{
				"val",
				nil,
				MustCompileDef(defparse.MustParse("'bool'")),
			},
			Out: []interface{}{
				nil,
				bwerror.Error("<ansiOutline>val<ansiCmd><ansi> (<ansiSecondaryLiteral>null<ansi>) is not of type <ansiPrimaryLiteral>Bool"),
			},
		},
		"val: nil, def: bool.isOptional": {
			In: []interface{}{
				"val",
				nil,
				MustCompileDef(defparse.MustParse("{ type: 'bool', isOptional: true }")),
			},
			Out: []interface{}{
				nil,
				nil,
			},
		},
		"val: nil, def: bool.default=true": {
			In: []interface{}{
				"val",
				nil,
				MustCompileDef(defparse.MustParse("{ type: 'bool', default: true }")),
			},
			Out: []interface{}{
				true,
				nil,
			},
		},

		// ==============================

		"val: valid, def: string.enum": {
			In: []interface{}{
				"val",
				"one",
				MustCompileDef(defparse.MustParse("{ type: 'string', enum: [qw/one two three/] }")),
			},
			Out: []interface{}{
				"one",
				nil,
			},
		},
		"val: invalid, def: string.enum": {
			In: []interface{}{
				"val",
				"One",
				MustCompileDef(defparse.MustParse("{ type: 'string', enum: [qw/one two three/] }")),
			},
			Out: []interface{}{
				nil,
				bwerror.Error("<ansiOutline>val<ansiCmd><ansi> (<ansiSecondaryLiteral>\"One\"<ansi>) has non supported value"),
			},
		},

		// ==============================

		"val: invalid, def: int": {
			In: []interface{}{
				"val",
				"1",
				MustCompileDef(defparse.MustParse("{ type: 'int'  }")),
			},
			Out: []interface{}{
				nil,
				bwerror.Error("<ansiOutline>val<ansiCmd><ansi> (<ansiSecondaryLiteral>\"1\"<ansi>) is not of type <ansiPrimaryLiteral>Int"),
			},
		},
		"val: valid, def: int": {
			In: []interface{}{
				"val",
				1,
				MustCompileDef(defparse.MustParse("{ type: 'int' }")),
			},
			Out: []interface{}{
				1,
				nil,
			},
		},
		"val: valid, def: int.min": {
			In: []interface{}{
				"val",
				1,
				MustCompileDef(defparse.MustParse("{ type: 'int', minInt: 0 }")),
			},
			Out: []interface{}{
				1,
				nil,
			},
		},
		"val: invalid, def: int.min": {
			In: []interface{}{
				"val",
				1,
				MustCompileDef(defparse.MustParse("{ type: 'int', minInt: 2 }")),
			},
			Out: []interface{}{
				nil,
				bwerror.Error("<ansiOutline>val<ansiCmd><ansi> (<ansiSecondaryLiteral>1<ansi>) is less then <ansiOutline>minLimit <ansiPrimaryLiteral>2"),
			},
		},
		"val: invalid, def: int.min.max": {
			In: []interface{}{
				"val",
				1,
				MustCompileDef(defparse.MustParse("{ type: 'int', minInt: 2, maxInt: 3 }")),
			},
			Out: []interface{}{
				nil,
				bwerror.Error("<ansiOutline>val<ansiCmd><ansi> (<ansiSecondaryLiteral>1<ansi>) is out of <ansiOutline>range <ansiSecondaryLiteral>[2, 3]"),
			},
		},
		"val: invalid, def: int.max": {
			In: []interface{}{
				"val",
				1,
				MustCompileDef(defparse.MustParse("{ type: 'int', maxInt: 0 }")),
			},
			Out: []interface{}{
				nil,
				bwerror.Error("<ansiOutline>val<ansiCmd><ansi> (<ansiSecondaryLiteral>1<ansi>) is greater then <ansiOutline>maxLimit <ansiPrimaryLiteral>0"),
			},
		},

		// ==============================

		"val: invalid, def: number": {
			In: []interface{}{
				"val",
				"3.14",
				MustCompileDef(defparse.MustParse("{ type: 'number'  }")),
			},
			Out: []interface{}{
				nil,
				bwerror.Error("<ansiOutline>val<ansiCmd><ansi> (<ansiSecondaryLiteral>\"3.14\"<ansi>) is not of type <ansiPrimaryLiteral>Number"),
			},
		},
		"val: valid, def: number": {
			In: []interface{}{
				"val",
				3.14,
				MustCompileDef(defparse.MustParse("{ type: 'number' }")),
			},
			Out: []interface{}{
				3.14,
				nil,
			},
		},
		"val: valid, def: number.min": {
			In: []interface{}{
				"val",
				3.14,
				MustCompileDef(defparse.MustParse("{ type: 'number', minNumber: 2.71 }")),
			},
			Out: []interface{}{
				3.14,
				nil,
			},
		},
		"val: invalid, def: number.min": {
			In: []interface{}{
				"val",
				2.71,
				MustCompileDef(defparse.MustParse("{ type: 'number', minNumber: 3.14 }")),
			},
			Out: []interface{}{
				nil,
				bwerror.Error("<ansiOutline>val<ansiCmd><ansi> (<ansiSecondaryLiteral>2.71<ansi>) is less then <ansiOutline>minLimit <ansiPrimaryLiteral>3.14"),
			},
		},
		"val: invalid, def: number.min.max": {
			In: []interface{}{
				"val",
				2.71,
				MustCompileDef(defparse.MustParse("{ type: 'number', minNumber: 3.14, maxNumber: 273 }")),
			},
			Out: []interface{}{
				nil,
				bwerror.Error("<ansiOutline>val<ansiCmd><ansi> (<ansiSecondaryLiteral>2.71<ansi>) is out of <ansiOutline>range <ansiSecondaryLiteral>[3.14, 273]"),
			},
		},
		"val: invalid, def: number.max": {
			In: []interface{}{
				"val",
				3.14,
				MustCompileDef(defparse.MustParse("{ type: 'number', maxNumber: 2.71 }")),
			},
			Out: []interface{}{
				nil,
				bwerror.Error("<ansiOutline>val<ansiCmd><ansi> (<ansiSecondaryLiteral>3.14<ansi>) is greater then <ansiOutline>maxLimit <ansiPrimaryLiteral>2.71"),
			},
		},

		// ==============================

		"val: nil, simple.def: map": {
			In: []interface{}{
				"val",
				nil,
				MustCompileDef(defparse.MustParse("'map'")),
			},
			Out: []interface{}{
				map[string]interface{}{},
				nil,
			},
		},

		"val: valid, def: map": {
			In: []interface{}{
				"val",
				map[string]interface{}{
					"boolKey":   true,
					"intKey":    273,
					"numberKey": 3.14,
					"stringKey": "something",
				},
				MustCompileDef(defparse.MustParse(`{
					type 'map'
					keys {
						boolKey 'bool'
						intKey 'int'
						numberKey 'number'
						stringKey 'string'
					}
				}`)),
			},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} {
					return test.In[1]
				},
				nil,
			},
		},
		"val: unexpected keys, def: map": {
			In: []interface{}{
				"val",
				map[string]interface{}{
					"boolKey":   true,
					"intKey":    273,
					"numberKey": 3.14,
					"stringKey": "something",
				},
				MustCompileDef(defparse.MustParse(`{
					type 'map'
					keys {
						boolKey 'bool'
						intKey 'int'
					}
				}`)),
			},
			Out: []interface{}{
				nil,
				func(test bwtesting.TestCaseStruct) error {
					return bwerror.Error(
						"<ansiOutline>val<ansiCmd><ansi> (<ansiSecondaryLiteral>" +
							bwjson.PrettyJson(test.In[1]) +
							"<ansi>) has unexpected keys <ansiSecondaryLiteral>" +
							bwjson.PrettyJson(defparse.MustParse(`[qw/numberKey stringKey/]`)),
					)
				},
			},
		},
		"val: valid, def: map.keys.elem": {
			In: []interface{}{
				"val",
				map[string]interface{}{
					"boolKey":    true,
					"intKey":     273,
					"numberKey1": 3.14,
					"numberKey2": 2.71,
				},
				MustCompileDef(defparse.MustParse(`{
					type 'map'
					keys {
						boolKey 'bool'
						intKey 'int'
					}
					elem 'number'
				}`)),
			},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} {
					return test.In[1]
				},
				nil,
			},
		},
		"val: valid, def: map.elem": {
			In: []interface{}{
				"val",
				map[string]interface{}{
					"numberKey1": 3.14,
					"numberKey2": 2.71,
				},
				MustCompileDef(defparse.MustParse(`{
					type 'map'
					elem 'number'
				}`)),
			},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} {
					return test.In[1]
				},
				nil,
			},
		},

		// ==============================

		"val: valid, def: array": {
			In: []interface{}{
				"val",
				[]int{1, 2, 3},
				MustCompileDef(defparse.MustParse(`{
					type 'array'
				}`)),
			},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} {
					return test.In[1]
				},
				nil,
			},
		},

		"val: valid, def: array.arrayElem": {
			In: []interface{}{
				"val",
				[]int{1, 2, 3},
				MustCompileDef(defparse.MustParse(`{
					type 'array'
					arrayElem 'number'
				}`)),
			},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} {
					return test.In[1]
				},
				nil,
			},
		},

		"val: valid, def: array.elem": {
			In: []interface{}{
				"val",
				defparse.MustParse(`[1 2 3]`),
				// []int{1, 2, 3},
				MustCompileDef(defparse.MustParse(`{
					type 'array'
					arrayElem 'int'
				}`)),
			},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} {
					return test.In[1]
				},
				nil,
			},
		},

		// "ExecCmd opt": {
		// 	what: "ExecCmd.opt",
		// 	val:  nil,
		// 	def: defparse.MustParse(`{
		//      type: 'map',
		//      keys: {
		//        v: {
		//          type: 'enum'
		//          enum: [ qw/all err ok none/ ]
		//          default: 'none'
		//        }
		//        s: {
		//          type: 'enum'
		//          enum: [ qw/none stderr stdout all/ ]
		//          default: 'all'
		//        }
		//        exitOnError: {
		//          type: 'bool'
		//          default: false
		//        }
		//      }
		//    }`),
		// 	result: defparse.MustParse(`{
		// 		v: 'none'
		// 		s: 'all'
		// 		exitOnError: false
		// 	}`),
		// },
	}
	testsToRun := tests
	bwmap.CropMap(testsToRun)
	bwmap.CropMap(testsToRun, "val: valid, def: array.arrayElem")
	bwtesting.BwRunTests(t, testsToRun, ValidateVal)
}
