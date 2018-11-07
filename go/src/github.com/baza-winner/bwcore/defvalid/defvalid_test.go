package defvalid

import (
	"testing"

	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/defparse"
	"github.com/baza-winner/bwcore/defvalid/deftype"
)

// ============================================================================

func TestCompileDef(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"def: nil": {
			In: []interface{}{nil},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.Case) error {
					return bwerr.From(
						"<ansiVar>def<ansiPath><ansi> (<ansiVal>%#v<ansi>) has non supported value",
						test.In[0],
					)
				},
			},
		},
		"simple.def: invalid type": {
			In: []interface{}{defparse.MustParse(`false`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerr.From("<ansiVar>def<ansiPath><ansi> (<ansiVal>%#v<ansi>) is not of type <ansiVal>[]string<ansi>, or <ansiVal>map[string]<ansi>, or <ansiVal>string", false),
			},
		},
		"def: no .type": {
			In: []interface{}{defparse.MustParse(`{ }`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.Case) error {
					return bwerr.From(
						"<ansiVar>def<ansiPath><ansi> (<ansiVal>%#v<ansi>) has no key <ansiVal>type",
						test.In[0],
					)
				},
			},
		},
		"def: invalid type": {
			In: []interface{}{defparse.MustParse(`{ type false }`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerr.From("<ansiVar>def<ansiPath>.type<ansi> (<ansiVal>%#v<ansi>) is not of type <ansiVal>[]string<ansi>, or <ansiVal>string", false),
			},
		},
		"def: simple valid": {
			In: []interface{}{defparse.MustParse(`"Bool"`)},
			Out: []interface{}{
				&Def{tp: deftype.From(deftype.Bool)},
				nil,
			},
		},
		"def: invalid deftypeItem": {
			In: []interface{}{defparse.MustParse(`[ qw/ Bool Int some / ]`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerr.From("<ansiVar>def<ansiPath>.#2<ansi> (<ansiVal>%#v<ansi>) has non supported value", "some"),
			},
		},
		"def: enum": {
			In: []interface{}{defparse.MustParse(`{ type: "String", enum: [qw/one two three/]}`)},
			Out: []interface{}{
				&Def{tp: deftype.From(deftype.String), enum: bwset.StringFrom("one", "two", "three")},
				nil,
			},
		},
		"def: invalid enum": {
			In: []interface{}{defparse.MustParse(`{ type: "String", enum: [qw/one two three/ true]}`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerr.From("<ansiVar>def<ansiPath>.enum<ansi> (<ansiVal>%#v<ansi>) is not of type <ansiVal>[]string", defparse.MustParse(`[qw/one two three/ true]`)),
			},
		},
		"def: map with keys": {
			In: []interface{}{defparse.MustParse(`{ type: "Map", keys: { keyBool: ['Bool'] }}`)},
			Out: []interface{}{
				&Def{
					tp: deftype.From(deftype.Map),
					keys: map[string]Def{
						"keyBool": Def{tp: deftype.From(deftype.Bool)},
					}},
				nil,
			},
		},
		"def: map with invalid keys": {
			In: []interface{}{
				map[string]interface{}{
					"type": "Map",
					"keys": map[int]interface{}{
						0: nil,
					},
				},
			},
			Out: []interface{}{
				(*Def)(nil),
				// nil,
				bwerr.From(
					"<ansiVar>def<ansiPath>.keys<ansi> (<ansiVal>%#v<ansi>) is not of type <ansiVal>map[string]",
					map[int]interface{}{0: nil},
				),
			},
		},
		"def: map with invalid Def in keys": {
			In: []interface{}{
				defparse.MustParse(`{ type: "Map", keys: { keyBool: { type 'Boolean' } }}`),
			},
			Out: []interface{}{
				(*Def)(nil),
				bwerr.From(
					"<ansiVar>def<ansiPath>.keys.keyBool.type<ansi> (<ansiVal>%#v<ansi>) has non supported value",
					"Boolean",
				),
			},
		},
		"def: map with valid Def in keys": {
			In: []interface{}{
				defparse.MustParse(`{ type: "Map", keys: { keyBool: { type 'Bool' } }}`),
			},
			Out: []interface{}{
				&Def{
					tp:         deftype.From(deftype.Map),
					isOptional: false,
					keys: map[string]Def{
						"keyBool": Def{tp: deftype.From(deftype.Bool), isOptional: false},
					},
				},
				nil,
			},
		},
		"def: map with invalid Def in elem": {
			In: []interface{}{
				defparse.MustParse(`{ type: "Map", keys: { keyBool: { type 'Bool' } }, elem: 'Boolean'}`),
			},
			Out: []interface{}{
				(*Def)(nil),
				bwerr.From(
					"<ansiVar>def<ansiPath>.elem<ansi> (<ansiVal>%#v<ansi>) has non supported value",
					"Boolean",
				),
			},
		},
		"def: unexpected keys": {
			In: []interface{}{defparse.MustParse(`{ type: "Map", kyes: { keyBool: ['Bool'] }, some: 'thing'}`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.Case) error {
					return bwerr.From(
						"<ansiVar>def<ansiPath><ansi> (<ansiVal>%#v<ansi>) has unexpected keys <ansiVal>%s",
						test.In[0],
						bwjson.Pretty(defparse.MustParse(`[qw/kyes some/]`)),
					)
				},
			},
		},
		"def: array with arrayElem": {
			In: []interface{}{defparse.MustParse(`{ type: "Array", arrayElem: 'Int' }`)},
			Out: []interface{}{
				&Def{
					tp:        deftype.From(deftype.Array),
					arrayElem: &Def{tp: deftype.From(deftype.Int)},
				},
				nil,
			},
		},
		"def: array with invalidDef in arrayElem": {
			In: []interface{}{defparse.MustParse(`{ type: "Array", arrayElem: 'Integer' }`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerr.From(
					"<ansiVar>def<ansiPath>.arrayElem<ansi> (<ansiVal>%#v<ansi>) has non supported value",
					"Integer",
				),
			},
		},
		"def: array with arrayElem and Elem": {
			In: []interface{}{defparse.MustParse(`{ type: "Array", arrayElem: 'Int', elem: 'bool' }`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.Case) error {
					return bwerr.From(
						"<ansiVar>def<ansiPath><ansi> (<ansiVal>%#v<ansi>) has unexpected key <ansiVal>elem",
						test.In[0],
					)
				},
			},
		},
		"def: invalid minInt": {
			In: []interface{}{defparse.MustParse(`{ type: "Int", minInt: true }`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerr.From(
					"<ansiVar>def<ansiPath>.minInt<ansi> (<ansiVal>%#v<ansi>) is not of type <ansiVal>int64",
					true,
				),
			},
		},
		"def: invalid maxInt": {
			In: []interface{}{defparse.MustParse(`{ type: "Int", maxInt: true }`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerr.From(
					"<ansiVar>def<ansiPath>.maxInt<ansi> (<ansiVal>%#v<ansi>) is not of type <ansiVal>int64",
					true,
				),
			},
		},
		"def: minInt, maxInt": {
			In: []interface{}{defparse.MustParse(`{ type: "Int", minInt: -6, maxInt: 10 }`)},
			Out: []interface{}{
				&Def{
					tp:     deftype.From(deftype.Int),
					minInt: ptrToInt64(-6),
					maxInt: ptrToInt64(10),
				},
				nil,
			},
		},
		"def: minInt > maxInt": {
			In: []interface{}{defparse.MustParse(`{ type: "Int", minInt: 6, maxInt: -10 }`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.Case) error {
					return bwerr.From(
						"<ansiVar>def<ansiPath><ansi> (<ansiVal>%#v<ansi>) has conflicting keys: <ansiVal>%s",
						test.In[0], bwjson.Pretty(defparse.MustParse("{ minInt: 6, maxInt: -10 }")),
					)
				},
			},
		},
		"def: invalid minNumber": {
			In: []interface{}{defparse.MustParse(`{ type: "Number", minNumber: true }`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerr.From(
					"<ansiVar>def<ansiPath>.minNumber<ansi> (<ansiVal>%#v<ansi>) is not of type <ansiVal>float64",
					true,
				),
			},
		},
		"def: invalid maxNumber": {
			In: []interface{}{defparse.MustParse(`{ type: "Number", maxNumber: true }`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerr.From(
					"<ansiVar>def<ansiPath>.maxNumber<ansi> (<ansiVal>%#v<ansi>) is not of type <ansiVal>float64",
					true,
				),
			},
		},
		"def: minNumber, maxNumber": {
			In: []interface{}{defparse.MustParse(`{ type: "Number", minNumber: -6, maxNumber: 10 }`)},
			Out: []interface{}{
				&Def{
					tp:        deftype.From(deftype.Number),
					minNumber: ptrToFloat64(float64(-6)),
					maxNumber: ptrToFloat64(float64(10)),
				},
				nil,
			},
		},
		"def: minNumber > maxNumber": {
			In: []interface{}{defparse.MustParse(`{ type: "Number", minNumber: 3.14, maxNumber: -2.71 }`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.Case) error {
					return bwerr.From(
						"<ansiVar>def<ansiPath><ansi> (<ansiVal>%#v<ansi>) has conflicting keys: <ansiVal>%s",
						test.In[0], bwjson.Pretty(defparse.MustParse("{ minNumber: 3.14, maxNumber: -2.71 }")),
					)
				},
			},
		},
		"def: invalid isOptional": {
			In: []interface{}{defparse.MustParse(`{ type: "Bool", isOptional: 0 }`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerr.From(
					"<ansiVar>def<ansiPath>.isOptional<ansi> (<ansiVal>%#v<ansi>) is not of type <ansiVal>bool",
					int8(0),
				),
			},
		},
		"def: default": {
			In: []interface{}{defparse.MustParse(`{ type: "Bool", default: true }`)},
			Out: []interface{}{
				&Def{
					tp:         deftype.From(deftype.Bool),
					dflt:       true,
					isOptional: true,
				},
				nil,
			},
		},
		"def: default 'string' for bool": {
			In: []interface{}{defparse.MustParse(`{ type: "Bool", default: "string" }`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerr.From("<ansiVar>def<ansiPath>.default<ansi> (<ansiVal>%#v<ansi>) is not of type <ansiVal>Bool", "string"),
			},
		},
		"def: isOptional": {
			In: []interface{}{defparse.MustParse(`{ type: "Bool", isOptional: true }`)},
			Out: []interface{}{
				&Def{
					tp:         deftype.From(deftype.Bool),
					isOptional: true,
				},
				nil,
			},
		},
		"def: isOptional = false conflicts with dflt": {
			In: []interface{}{defparse.MustParse(`{ type: "Bool", isOptional: false, default: false }`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.Case) error {
					return bwerr.From(
						"<ansiVar>def<ansiPath><ansi> (<ansiVal>%#v<ansi>) has conflicting keys: <ansiVal>%s",
						test.In[0], bwjson.Pretty(defparse.MustParse("{ isOptional: false, default: false }")),
					)
				},
			},
		},
		"def: ArrayOf without follower": {
			In: []interface{}{defparse.MustParse(`{ type: "ArrayOf"  }`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerr.From("<ansiVar>def<ansiPath>.type<ansi> (<ansiVal>%#v<ansi>) must be followed by some type", "ArrayOf"),
			},
		},
		"simple.def: ArrayOf without follower": {
			In: []interface{}{defparse.MustParse(`"ArrayOf"`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerr.From("<ansiVar>def<ansiPath><ansi> (<ansiVal>%#v<ansi>) must be followed by some type", "ArrayOf"),
			},
		},
		"simple.def: ArrayOf can not be combined with array": {
			In: []interface{}{defparse.MustParse(`["ArrayOf", "Array", "String"]`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.Case) error {
					return bwerr.From(
						"<ansiVar>def<ansiPath><ansi> (<ansiVal>%#v<ansi>) following values can not be combined: <ansiVal>%s",
						test.In[0], bwjson.Pretty(defparse.MustParse("[ 'Array', 'ArrayOf' ]")),
					)
				},
			},
		},
		"simple.def: ArrayOf with non Array default": {
			In: []interface{}{defparse.MustParse(`{
				type ["ArrayOf", "Int"]
				default 3
			}`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerr.From(
					"<ansiVar>def<ansiPath>.default<ansi> (<ansiVal>%#v<ansi>) is not of type <ansiVal>Array", int8(3),
				),
			},
		},
		"simple.def: ArrayOf with valid default": {
			In: []interface{}{defparse.MustParse(`{
				type ["ArrayOf", "Int"]
				default [3]
			}`)},
			Out: []interface{}{
				&Def{
					tp:         deftype.From(deftype.ArrayOf, deftype.Int),
					isOptional: true,
					dflt:       []interface{}{3},
				},
				nil,
			},
		},
	}
	testsToRun := tests
	bwmap.CropMap(testsToRun)
	// bwmap.CropMap(testsToRun, "simple.def: ArrayOf with valid default")
	bwtesting.BwRunTests(t, CompileDef, testsToRun)
}

// ============================================================================

func TestValidateVal(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"val: nil, simple.def: bool": {
			In: []interface{}{
				"val",
				nil,
				MustCompileDef(defparse.MustParse("'Bool'")),
			},
			Out: []interface{}{
				nil,
				bwerr.From("<ansiVar>val<ansiPath><ansi> (<ansiVal>%#v<ansi>) is not of type <ansiVal>Bool", nil),
			},
		},
		"val: nil, def: bool.isOptional": {
			In: []interface{}{
				"val",
				nil,
				MustCompileDef(defparse.MustParse("{ type: 'Bool', isOptional: true }")),
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
				MustCompileDef(defparse.MustParse("{ type: 'Bool', default: true }")),
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
				MustCompileDef(defparse.MustParse("{ type: 'String', enum: [qw/one two three/] }")),
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
				MustCompileDef(defparse.MustParse("{ type: 'String', enum: [qw/one two three/] }")),
			},
			Out: []interface{}{
				nil,
				bwerr.From("<ansiVar>val<ansiPath><ansi> (<ansiVal>%#v<ansi>) has non supported value", "One"),
			},
		},

		// ==============================

		"val: invalid, def: int": {
			In: []interface{}{
				"val",
				"1",
				MustCompileDef(defparse.MustParse("{ type: 'Int'  }")),
			},
			Out: []interface{}{
				nil,
				bwerr.From("<ansiVar>val<ansiPath><ansi> (<ansiVal>%#v<ansi>) is not of type <ansiVal>Int", "1"),
			},
		},
		"val: valid, def: int": {
			In: []interface{}{
				"val",
				1,
				MustCompileDef(defparse.MustParse("{ type: 'Int' }")),
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
				MustCompileDef(defparse.MustParse("{ type: 'Int', minInt: 0 }")),
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
				MustCompileDef(defparse.MustParse("{ type: 'Int', minInt: 2 }")),
			},
			Out: []interface{}{
				nil,
				bwerr.From("<ansiVar>val<ansiPath><ansi> (<ansiVal>%#v<ansi>) is less then <ansiVar>minLimit <ansiVal>2", 1),
			},
		},
		"val: invalid, def: int.min.max": {
			In: []interface{}{
				"val",
				1,
				MustCompileDef(defparse.MustParse("{ type: 'Int', minInt: 2, maxInt: 3 }")),
			},
			Out: []interface{}{
				nil,
				bwerr.From("<ansiVar>val<ansiPath><ansi> (<ansiVal>%#v<ansi>) is out of <ansiVar>range <ansiVal>[2, 3]", 1),
			},
		},
		"val: invalid, def: int.max": {
			In: []interface{}{
				"val",
				1,
				MustCompileDef(defparse.MustParse("{ type: 'Int', maxInt: 0 }")),
			},
			Out: []interface{}{
				nil,
				bwerr.From("<ansiVar>val<ansiPath><ansi> (<ansiVal>%#v<ansi>) is greater then <ansiVar>maxLimit <ansiVal>0", 1),
			},
		},

		// ==============================

		"val: invalid, def: Number": {
			In: []interface{}{
				"val",
				"3.14",
				MustCompileDef(defparse.MustParse("{ type: 'Number'  }")),
			},
			Out: []interface{}{
				nil,
				bwerr.From("<ansiVar>val<ansiPath><ansi> (<ansiVal>%#v<ansi>) is not of type <ansiVal>Number", "3.14"),
			},
		},
		"val: valid, def: Number": {
			In: []interface{}{
				"val",
				3.14,
				MustCompileDef(defparse.MustParse("{ type: 'Number' }")),
			},
			Out: []interface{}{
				3.14,
				nil,
			},
		},
		"val: valid, def: Number.min": {
			In: []interface{}{
				"val",
				3.14,
				MustCompileDef(defparse.MustParse("{ type: 'Number', minNumber: 2.71 }")),
			},
			Out: []interface{}{
				3.14,
				nil,
			},
		},
		"val: invalid, def: Number.min": {
			In: []interface{}{
				"val",
				2.71,
				MustCompileDef(defparse.MustParse("{ type: 'Number', minNumber: 3.14 }")),
			},
			Out: []interface{}{
				nil,
				bwerr.From("<ansiVar>val<ansiPath><ansi> (<ansiVal>%#v<ansi>) is less then <ansiVar>minLimit <ansiVal>3.14", 2.71),
			},
		},
		"val: invalid, def: Number.min.max": {
			In: []interface{}{
				"val",
				2.71,
				MustCompileDef(defparse.MustParse("{ type: 'Number', minNumber: 3.14, maxNumber: 273 }")),
			},
			Out: []interface{}{
				nil,
				bwerr.From("<ansiVar>val<ansiPath><ansi> (<ansiVal>%#v<ansi>) is out of <ansiVar>range <ansiVal>[3.14, 273]", 2.71),
			},
		},
		"val: invalid, def: Number.max": {
			In: []interface{}{
				"val",
				3.14,
				MustCompileDef(defparse.MustParse("{ type: 'Number', maxNumber: 2.71 }")),
			},
			Out: []interface{}{
				nil,
				bwerr.From("<ansiVar>val<ansiPath><ansi> (<ansiVal>%#v<ansi>) is greater then <ansiVar>maxLimit <ansiVal>2.71", 3.14),
			},
		},

		// ==============================

		"val: nil, simple.def: map": {
			In: []interface{}{
				"val",
				nil,
				MustCompileDef(defparse.MustParse("'Map'")),
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
					type 'Map'
					keys {
						boolKey 'Bool'
						intKey 'Int'
						numberKey 'Number'
						stringKey 'String'
					}
				}`)),
			},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} {
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
					type 'Map'
					keys {
						boolKey 'Bool'
						intKey 'Int'
					}
				}`)),
			},
			Out: []interface{}{
				nil,
				func(test bwtesting.Case) error {
					return bwerr.From(
						"<ansiVar>val<ansiPath><ansi> (<ansiVal>%#v<ansi>) has unexpected keys <ansiVal>%s",
						test.In[1], bwjson.Pretty(defparse.MustParse(`[qw/numberKey stringKey/]`)),
					)
				},
			},
		},
		"val: invalid key value, def: map": {
			In: []interface{}{
				"val",
				map[string]interface{}{
					"boolKey": 0,
				},
				MustCompileDef(defparse.MustParse(`{
					type 'Map'
					keys {
						boolKey 'Bool'
					}
				}`)),
			},
			Out: []interface{}{
				nil,
				bwerr.From(
					"<ansiVar>val<ansiPath>.boolKey<ansi> (<ansiVal>%#v<ansi>) is not of type <ansiVal>Bool",
					0,
				),
			},
		},
		"val: invalid elem value, def: map": {
			In: []interface{}{
				"val",
				map[string]interface{}{
					"boolKey": true,
					"intKey":  false,
				},
				MustCompileDef(defparse.MustParse(`{
					type 'Map'
					keys { boolKey: 'Bool' }
					elem 'Int'
				}`)),
			},
			Out: []interface{}{
				nil,
				bwerr.From(
					"<ansiVar>val<ansiPath>.intKey<ansi> (<ansiVal>%#v<ansi>) is not of type <ansiVal>Int",
					false,
				),
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
					type 'Map'
					keys {
						boolKey 'Bool'
						intKey 'Int'
					}
					elem 'Number'
				}`)),
			},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} {
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
					type 'Map'
					elem 'Number'
				}`)),
			},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} {
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
					type 'Array'
				}`)),
			},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} {
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
					type 'Array'
					arrayElem 'Number'
				}`)),
			},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} {
					return test.In[1]
				},
				nil,
			},
		},

		"val: valid, def: array.elem": {
			In: []interface{}{
				"val",
				defparse.MustParse(`[1 2 3]`),
				MustCompileDef(defparse.MustParse(`{
					type 'Array'
					arrayElem 'Int'
				}`)),
			},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} {
					return test.In[1]
				},
				nil,
			},
		},

		"val: scalar, def: ArrayOf,int": {
			In: []interface{}{
				"val",
				1,
				MustCompileDef(defparse.MustParse(`{
					type [ 'ArrayOf' 'Int' ]
				}`)),
			},
			Out: []interface{}{
				[]interface{}{1},
				nil,
			},
		},
		"val: array, def: ArrayOf,int": {
			In: []interface{}{
				"val",
				[]int{1, 2},
				MustCompileDef(defparse.MustParse(`{
					type [ 'ArrayOf' 'Int' ]
				}`)),
			},
			Out: []interface{}{
				[]int{1, 2},
				nil,
			},
		},
		"val: array, def: ArrayOf,int.default": {
			In: []interface{}{
				"val",
				[]int{1, 2},
				MustCompileDef(defparse.MustParse(`{
					type [ 'ArrayOf' 'Int' ]
					default [3]
				}`)),
			},
			Out: []interface{}{
				[]int{1, 2},
				nil,
			},
		},
		"val: nil, def: ArrayOf,int.default": {
			In: []interface{}{
				"val",
				nil,
				MustCompileDef(defparse.MustParse(`{
					type [ 'ArrayOf' 'Int' ]
					default [3]
				}`)),
			},
			Out: []interface{}{
				[]interface{}{3},
				nil,
			},
		},
		"ExecCmd opt": {
			In: []interface{}{
				"ExecCmd.opt",
				nil,
				MustCompileDef(defparse.MustParse(`{
		     type: 'Map',
		     keys: {
		       v: {
		         type: 'String'
		         enum: [ qw/all err ok none/ ]
		         default: 'none'
		       }
		       s: {
		         type: 'String'
		         enum: [ qw/none stderr stdout all/ ]
		         default: 'all'
		       }
		       exitOnError: {
		         type: 'Bool'
		         default: false
		       }
		     }
				}`)),
			},
			Out: []interface{}{
				defparse.MustParse(`{
					v: 'none'
					s: 'all'
					exitOnError: false
				}`),
				nil,
			},
		},
	}
	testsToRun := tests
	bwmap.CropMap(testsToRun)
	// bwmap.CropMap(testsToRun, "val: invalid, def: int.min.max")
	bwtesting.BwRunTests(t, ValidateVal, testsToRun)
}
