package bwval_test

import (
	"fmt"
	"testing"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwval"
)

func TestPathVal(t *testing.T) {
	for _, test := range []struct {
		name string
		val  interface{}
		path bw.ValPath
		vars map[string]interface{}
		eta  interface{}
		err  error
	}{
		{
			name: "nil",
			val:  nil,
			path: bwval.MustPathFrom("some"),
			eta:  nil,
		},
		{
			name: "self",
			val:  1,
			path: bwval.MustPathFrom(""),
			eta:  1,
		},
		{
			name: "by key",
			val:  map[string]interface{}{"some": "thing"},
			path: bwval.MustPathFrom("some"),
			eta:  "thing",
		},
		{
			name: "by idx (1)",
			val:  []interface{}{"some", "thing"},
			path: bwval.MustPathFrom("1"),
			eta:  "thing",
		},
		{
			name: "by idx (-1)",
			val:  []interface{}{"some", "thing"},
			path: bwval.MustPathFrom("-1"),
			eta:  "thing",
		},
		{
			name: "by idx (len)",
			val:  []interface{}{"some", "thing"},
			path: bwval.MustPathFrom("2"),
			eta:  nil,
		},
		{
			name: "by idx of nil",
			val:  nil,
			path: bwval.MustPathFrom("2"),
			eta:  nil,
		},
		{
			name: "by # of nil",
			val:  nil,
			path: bwval.MustPathFrom("1.#"),
			eta:  0,
		},
		{
			name: "by # of Array",
			val:  []interface{}{"a", "b"},
			path: bwval.MustPathFrom("#"),
			eta:  2,
		},
		{
			name: "by # of Map",
			val: []interface{}{
				"a",
				map[string]interface{}{"c": "d", "e": "f", "i": "g"},
			},
			path: bwval.MustPathFrom("1.#"),
			eta:  3,
		},
		{
			name: "by # of nil",
			val:  nil,
			path: bwval.MustPathFrom("1.#"),
			eta:  0,
		},
		{
			name: "by path (idx)",
			val:  map[string]interface{}{"some": []interface{}{"good", "thing"}, "idx": 1},
			path: bwval.MustPathFrom("some.{idx}"),
			eta:  "thing",
		},
		{
			name: "by path (key)",
			val:  map[string]interface{}{"some": []interface{}{"good", "thing"}, "key": "some"},
			path: bwval.MustPathFrom("{key}.1"),
			eta:  "thing",
		},
		{
			name: "by var",
			val:  map[string]interface{}{"some": []interface{}{"good", "thing"}},
			vars: map[string]interface{}{"idx": 1},
			path: bwval.MustPathFrom("some.{$idx}"),
			eta:  "thing",
		},
		{
			name: "by var",
			val:  map[string]interface{}{"some": []interface{}{"good", "thing"}},
			vars: nil,
			path: bwval.MustPathFrom("$idx"),
			eta:  nil,
		},
		{
			name: "err: is not Map",
			val:  1,
			path: bwval.MustPathFrom("some.{$idx}"),
			err:  bwerr.Error{S: "\x1b[96;1m(int)1\x1b[0m::\x1b[38;5;252;1msome\x1b[0m (\x1b[96;1m(int)1\x1b[0m) is not \x1b[97;1mMap\x1b[0m"},
		},
		{
			name: "err: is not Array",
			val:  1,
			path: bwval.MustPathFrom("{$idx}"),
			vars: map[string]interface{}{"idx": 1},
			err:  bwerr.Error{S: "\x1b[96;1m(int)1\x1b[0m::\x1b[38;5;252;1m{$idx}\x1b[0m (\x1b[96;1m(int)1\x1b[0m) is not \x1b[97;1mArray\x1b[0m"},
		},
		{
			name: "err: nor Map, neither Array",
			val:  1,
			path: bwval.MustPathFrom("#"),
			err:  bwerr.Error{S: "\x1b[96;1m(int)1\x1b[0m::\x1b[38;5;252;1m.#\x1b[0m (\x1b[96;1m(int)1\x1b[0m) is nor \x1b[97;1mMap\x1b[0m, neither \x1b[97;1mArray\x1b[0m"},
		},
		{
			name: "err: nor Int, neither String",
			val:  1,
			path: bwval.MustPathFrom("{$idx}"),
			vars: map[string]interface{}{"idx": nil},
			err:  bwerr.Error{S: "\x1b[96;1m(int)1\x1b[0m::\x1b[38;5;252;1m.{$idx}\x1b[0m (\x1b[96;1m(interface {})<nil>\x1b[0m) is nor \x1b[97;1mInt\x1b[0m, neither \x1b[97;1mString\x1b[0m"},
		},
	} {
		crop := false
		// crop = true
		if crop && test.name != "by # of Map" {
			continue
		}
		bwtesting.BwRunTests(t,
			bwval.FromVal(test.val).PathVal,
			map[string]bwtesting.TestCaseStruct{
				test.name: {
					In: []interface{}{
						test.path,
						test.vars,
					},
					Out: []interface{}{
						test.eta,
						test.err,
					},
				},
			},
		)
	}
}

func TestKind(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"Nil": {
			In: []interface{}{nil},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} { return test.In[0] },
				bwval.ValNil,
			},
		},
		"Nool": {
			In: []interface{}{true},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} { return test.In[0] },
				bwval.ValBool,
			},
		},
		"String": {
			In: []interface{}{"some"},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} { return test.In[0] },
				bwval.ValString,
			},
		},
		"Map": {
			In: []interface{}{map[string]interface{}{}},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} { return test.In[0] },
				bwval.ValMap,
			},
		},
		"Array": {
			In: []interface{}{[]interface{}{}},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} { return test.In[0] },
				bwval.ValArray,
			},
		},
		"Int(int8)": {
			In: []interface{}{bw.MaxInt8},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(int16)": {
			In: []interface{}{bw.MaxInt16},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(int32)": {
			In: []interface{}{bw.MaxInt32},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(int64(bw.MaxInt32))": {
			In: []interface{}{int64(bw.MaxInt32)},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(bw.MaxInt64)": {
			In: []interface{}{int64(bw.MaxInt64)},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} { return test.In[0] },
				func() (result bwval.ValKind) {
					if bw.MaxInt64 > int64(bw.MaxInt) {
						result = bwval.ValUnknown
					} else {
						result = bwval.ValInt
					}
					return
				},
			},
		},
		"Int(int)": {
			In: []interface{}{bw.MaxInt},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(uint8)": {
			In: []interface{}{bw.MaxUint8},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(uint16)": {
			In: []interface{}{bw.MaxUint16},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(uint32)": {
			In: []interface{}{bw.MaxUint32},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(uint64(bw.MaxInt32))": {
			In: []interface{}{uint64(bw.MaxUint32)},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(bw.MaxUint64)": {
			In: []interface{}{bw.MaxUint64},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) (result interface{}) {
					if bw.MaxUint64 <= uint64(bw.MaxInt) {
						result = test.In[0]
					}
					return
				},
				func() (result bwval.ValKind) {
					if bw.MaxUint64 > uint64(bw.MaxInt) {
						result = bwval.ValUnknown
					} else {
						result = bwval.ValInt
					}
					return
				},
			},
		},
		"Int(bw.MaxUint)": {
			In: []interface{}{bw.MaxUint},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) (result interface{}) {
					if bw.MaxUint <= uint(bw.MaxInt) {
						result = test.In[0]
					}
					return
				},
				func() (result bwval.ValKind) {
					if bw.MaxUint > uint(bw.MaxInt) {
						result = bwval.ValUnknown
					} else {
						result = bwval.ValInt
					}
					return
				},
			},
		},
		"Int(uint(0))": {
			In: []interface{}{uint(0)},
			Out: []interface{}{
				func(test bwtesting.TestCaseStruct) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, bwval.Kind, tests)
}

func ValKindPretty(vk bwval.ValKind) string {
	return bwjson.Pretty(vk)
}

func TestValKindMarshalJSON(t *testing.T) {
	var m = map[bwval.ValKind]string{
		bwval.ValUnknown: "ValUnknown",
		bwval.ValNil:     "ValNil",
		bwval.ValString:  "ValString",
		bwval.ValInt:     "ValInt",
		bwval.ValMap:     "ValMap",
		bwval.ValArray:   "ValArray",
	}
	for k, v := range m {
		bwtesting.BwRunTests(t,
			ValKindPretty,
			map[string]bwtesting.TestCaseStruct{
				v: {
					In:  []interface{}{k},
					Out: []interface{}{fmt.Sprintf("%q", v)},
				},
			},
		)
	}
}

func TestMap(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"Map": {
			In: []interface{}{
				map[string]interface{}{},
			},
			Out: []interface{}{
				map[string]interface{}{},
				true,
			},
		},
		"non Map": {
			In: []interface{}{
				1,
			},
			Out: []interface{}{
				(map[string]interface{})(nil),
				false,
			},
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, bwval.Map, tests)
}

func TestMustMap(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"Map": {
			In: []interface{}{
				map[string]interface{}{},
			},
			Out: []interface{}{
				map[string]interface{}{},
				nil,
			},
		},
		"non Map": {
			In: []interface{}{
				1,
			},
			Out: []interface{}{
				(map[string]interface{})(nil),
				"",
			},
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, MustMapWrapper, tests)
}

func MustMapWrapper(val interface{}) (result map[string]interface{}, panicrecover interface{}) {
	defer func() {
		panicrecover = recover()
		if t, ok := panicrecover.(bwerr.Error); ok {
			panicrecover = t.S
		}
	}()
	result = bwval.MustMap(val)
	return
}

func TestArray(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"Array": {
			In: []interface{}{
				[]interface{}{},
			},
			Out: []interface{}{
				[]interface{}{},
				true,
			},
		},
		"non Array": {
			In: []interface{}{
				1,
			},
			Out: []interface{}{
				([]interface{})(nil),
				false,
			},
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, bwval.Array, tests)
}

func TestMustArray(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"Array": {
			In: []interface{}{
				[]interface{}{},
			},
			Out: []interface{}{
				[]interface{}{},
				nil,
			},
		},
		"non Array": {
			In: []interface{}{
				1,
			},
			Out: []interface{}{
				([]interface{})(nil),
				"\x1b[91;1mERR: \x1b[0m\x1b[96;1m(int)1\x1b[0m is not \x1b[97;1mArray\x1b[0m at \x1b[32;1mgithub.com/baza-winner/bwcore/bwval.MustArray\x1b[38;5;243m@\x1b[97;1mbwval.go:119\x1b[0m",
			},
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, MustArrayWrapper, tests)
}

func MustArrayWrapper(val interface{}) (result []interface{}, panicrecover interface{}) {
	defer func() {
		panicrecover = recover()
	}()
	result = bwval.MustArray(val)
	return
}

func TestString(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"String": {
			In: []interface{}{
				"some",
			},
			Out: []interface{}{
				"some",
				true,
			},
		},
		"non String": {
			In: []interface{}{
				1,
			},
			Out: []interface{}{
				"",
				false,
			},
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, bwval.String, tests)
}

func TestInt(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"Int": {
			In: []interface{}{
				1,
			},
			Out: []interface{}{
				1,
				true,
			},
		},
		"non Int": {
			In: []interface{}{
				"some",
			},
			Out: []interface{}{
				0,
				false,
			},
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, bwval.Int, tests)
}

func TestBool(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"Bool": {
			In: []interface{}{
				false,
			},
			Out: []interface{}{
				false,
				true,
			},
		},
		"non Bool": {
			In: []interface{}{
				"some",
			},
			Out: []interface{}{
				false,
				false,
			},
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, bwval.Bool, tests)
}

func TestMustBool(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"Bool": {
			In: []interface{}{
				false,
			},
			Out: []interface{}{
				false,
				nil,
			},
		},
		"non Bool": {
			In: []interface{}{
				"some",
			},
			Out: []interface{}{
				false,
				"\x1b[91;1mERR: \x1b[0m\x1b[96;1m(string)some\x1b[0m is not \x1b[97;1mBool\x1b[0m at \x1b[32;1mgithub.com/baza-winner/bwcore/bwval.MustBool\x1b[38;5;243m@\x1b[97;1mbwval.go:59\x1b[0m",
			},
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, MustBoolWrapper, tests)
}

func MustBoolWrapper(val interface{}) (result bool, panicrecover interface{}) {
	defer func() {
		panicrecover = recover()
	}()
	result = bwval.MustBool(val)
	return
}

// func TestPathFrom(t *testing.T) {
// 	tests := map[string]bwtesting.TestCaseStruct{
// 		"some.thing": {
// 			In: []interface{}{
// 				func(testName string) string { return testName },
// 			},
// 			Out: []interface{}{
// 				Path{[]pathItem{
// 					{Type: pathItemKey, Key: "some"},
// 					{Type: pathItemKey, Key: "thing"},
// 				}},
// 				nil,
// 			},
// 		},
// 		"some.1": {
// 			In: []interface{}{
// 				func(testName string) string { return testName },
// 			},
// 			Out: []interface{}{
// 				Path{[]pathItem{
// 					{Type: pathItemKey, Key: "some"},
// 					{Type: pathItemIdx, Idx: 1},
// 				}},
// 				nil,
// 			},
// 		},
// 		"some.#": {
// 			In: []interface{}{
// 				func(testName string) string { return testName },
// 			},
// 			Out: []interface{}{
// 				Path{[]pathItem{
// 					{Type: pathItemKey, Key: "some"},
// 					{Type: pathItemHash},
// 				}},
// 				nil,
// 			},
// 		},
// 		"{some.thing}.good": {
// 			In: []interface{}{
// 				func(testName string) string { return testName },
// 			},
// 			Out: []interface{}{
// 				Path{[]pathItem{
// 					{Type: pathItemPath,
// 						Path: Path{[]pathItem{
// 							{Type: pathItemKey, Key: "some"},
// 							{Type: pathItemKey, Key: "thing"},
// 						}},
// 					},
// 					{Type: pathItemKey, Key: "good"},
// 				}},
// 				nil,
// 			},
// 		},
// 		"{$some.thing}.good": {
// 			In: []interface{}{
// 				func(testName string) string { return testName },
// 			},
// 			Out: []interface{}{
// 				Path{[]pathItem{
// 					{Type: pathItemPath,
// 						Path: Path{[]pathItem{
// 							{Type: pathItemVar, Key: "some"},
// 							{Type: pathItemKey, Key: "thing"},
// 						}},
// 					},
// 					{Type: pathItemKey, Key: "good"},
// 				}},
// 				nil,
// 			},
// 		},
// 	}

// 	bwmap.CropMap(tests)
// 	// bwmap.CropMap(tests, "UnexpectedItem")
// 	bwtesting.BwRunTests(t, PathFrom, tests)
// }

// func TestPathString(t *testing.T) {
// 	for _, test := range []struct {
// 		eta string
// 		v   Path
// 	}{
// 		{
// 			"some.thing",
// 			Path{[]pathItem{
// 				{Type: pathItemKey, Key: "some"},
// 				{Type: pathItemKey, Key: "thing"},
// 			}},
// 		},
// 		{
// 			"some.1",
// 			Path{[]pathItem{
// 				{Type: pathItemKey, Key: "some"},
// 				{Type: pathItemIdx, Idx: 1},
// 			}},
// 		},
// 		{
// 			"some.#",
// 			Path{[]pathItem{
// 				{Type: pathItemKey, Key: "some"},
// 				{Type: pathItemHash},
// 			}},
// 		},
// 		{
// 			"{some.thing}.good",
// 			Path{[]pathItem{
// 				{Type: pathItemPath,
// 					Path: Path{[]pathItem{
// 						{Type: pathItemKey, Key: "some"},
// 						{Type: pathItemKey, Key: "thing"},
// 					}},
// 				},
// 				{Type: pathItemKey, Key: "good"},
// 			}},
// 		},
// 		{
// 			"{$some.thing}.good",
// 			Path{[]pathItem{
// 				{Type: pathItemPath,
// 					Path: Path{[]pathItem{
// 						{Type: pathItemVar, Key: "some"},
// 						{Type: pathItemKey, Key: "thing"},
// 					}},
// 				},
// 				{Type: pathItemKey, Key: "good"},
// 			}},
// 		},
// 	} {
// 		bwtesting.BwRunTests(t,
// 			test.v.String,
// 			map[string]bwtesting.TestCaseStruct{
// 				test.eta: {
// 					In:  []interface{}{},
// 					Out: []interface{}{test.eta},
// 				},
// 			},
// 		)
// 	}
// }
