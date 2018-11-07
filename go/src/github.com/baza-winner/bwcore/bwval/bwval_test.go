package bwval_test

import (
	"fmt"
	"testing"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwval"
)

func TestMustPathVal(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"nil": {
			In:  []interface{}{bwval.FromVal(nil), bwval.PathFrom("some")},
			Out: []interface{}{nil},
		},
		"self": {
			In:  []interface{}{bwval.FromVal(1), bwval.PathFrom(".")},
			Out: []interface{}{1},
		},
		"by key": {
			In: []interface{}{
				bwval.FromVal(map[string]interface{}{"some": "thing"}),
				bwval.PathFrom("some"),
			},
			Out: []interface{}{"thing"},
		},
		"by idx (1)": {
			In: []interface{}{
				bwval.FromVal([]interface{}{"some", "thing"}),
				bwval.PathFrom("1"),
			},
			Out: []interface{}{"thing"},
		},
		"by idx (-1)": {
			In: []interface{}{
				bwval.FromVal([]interface{}{"some", "thing"}),
				bwval.PathFrom("-1"),
			},
			Out: []interface{}{"thing"},
		},
		"by idx (len)": {
			In: []interface{}{
				bwval.FromVal([]interface{}{"some", "thing"}),
				bwval.PathFrom("2"),
			},
			Out: []interface{}{nil},
		},
		"by idx of nil": {
			In: []interface{}{
				bwval.FromVal(nil),
				bwval.PathFrom("2"),
			},
			Out: []interface{}{nil},
		},
		"by # of nil": {
			In: []interface{}{
				bwval.FromVal(nil),
				bwval.PathFrom("1.#"),
			},
			Out: []interface{}{0},
		},
		"by # of Array": {
			In: []interface{}{
				bwval.FromVal([]interface{}{"a", "b"}),
				bwval.PathFrom("#"),
			},
			Out: []interface{}{2},
		},
		"by # of Map": {
			In: []interface{}{
				bwval.FromVal(
					[]interface{}{
						"a",
						map[string]interface{}{"c": "d", "e": "f", "i": "g"},
					},
				),
				bwval.PathFrom("1.#"),
			},
			Out: []interface{}{3},
		},
		"by path (idx)": {
			In: []interface{}{
				bwval.FromVal(map[string]interface{}{"some": []interface{}{"good", "thing"}, "idx": 1}),
				bwval.PathFrom("some.{idx}"),
			},
			Out: []interface{}{"thing"},
		},
		"by path (key)": {
			In: []interface{}{
				bwval.FromVal(map[string]interface{}{"some": []interface{}{"good", "thing"}, "key": "some"}),
				bwval.PathFrom("{key}.1"),
			},
			Out: []interface{}{"thing"},
		},
		"some.{$idx}": {
			In: []interface{}{
				bwval.FromVal(map[string]interface{}{"some": []interface{}{"good", "thing"}}),
				bwval.PathFrom("some.{$idx}"),
				map[string]interface{}{"idx": 1},
			},
			Out: []interface{}{"thing"},
		},
		"$idx": {
			In: []interface{}{
				bwval.FromVal(map[string]interface{}{"some": []interface{}{"good", "thing"}}),
				bwval.PathFrom("$idx"),
			},
			Out: []interface{}{nil},
		},
		"err: is not Map": {
			In: []interface{}{
				bwval.FromVal(1),
				bwval.PathFrom("some.{$idx}"),
			},
			Out:   []interface{}{nil},
			Panic: "\x1b[96;1m(int)1\x1b[0m::\x1b[38;5;252;1msome\x1b[0m (\x1b[96;1m(int)1\x1b[0m) is not \x1b[97;1mMap\x1b[0m",
		},
		"err: is not Array": {
			In: []interface{}{
				bwval.FromVal(1),
				bwval.PathFrom("{$idx}"),
				map[string]interface{}{"idx": 1},
			},
			Out:   []interface{}{nil},
			Panic: "\x1b[96;1m(int)1\x1b[0m::\x1b[38;5;252;1m{$idx}\x1b[0m (\x1b[96;1m(int)1\x1b[0m) is not \x1b[97;1mArray\x1b[0m",
		},
		"err: nor Map, neither Array": {
			In: []interface{}{
				bwval.FromVal(1),
				bwval.PathFrom("#"),
			},
			Out:   []interface{}{nil},
			Panic: "\x1b[96;1m(int)1\x1b[0m::\x1b[38;5;252;1m.#\x1b[0m (\x1b[96;1m(int)1\x1b[0m) is nor \x1b[97;1mMap\x1b[0m, neither \x1b[97;1mArray\x1b[0m",
		},
		"err: nor Int, neither String": {
			In: []interface{}{
				bwval.FromVal(1),
				bwval.PathFrom("{$idx}"),
				map[string]interface{}{"idx": nil},
			},
			Out:   []interface{}{nil},
			Panic: "\x1b[96;1m(int)1\x1b[0m::\x1b[38;5;252;1m.{$idx}\x1b[0m (\x1b[96;1m(interface {})<nil>\x1b[0m) is nor \x1b[97;1mInt\x1b[0m, neither \x1b[97;1mString\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	bwtesting.BwRunTests(t, bwval.MustPathVal, tests)
}

func TestKind(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"Nil": {
			In: []interface{}{nil},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} { return test.In[0] },
				bwval.ValNil,
			},
		},
		"Bool": {
			In: []interface{}{true},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} { return test.In[0] },
				bwval.ValBool,
			},
		},
		"String": {
			In: []interface{}{"some"},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} { return test.In[0] },
				bwval.ValString,
			},
		},
		"Map": {
			In: []interface{}{map[string]interface{}{}},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} { return test.In[0] },
				bwval.ValMap,
			},
		},
		"Array": {
			In: []interface{}{[]interface{}{}},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} { return test.In[0] },
				bwval.ValArray,
			},
		},
		"Int(int8)": {
			In: []interface{}{bw.MaxInt8},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(int16)": {
			In: []interface{}{bw.MaxInt16},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(int32)": {
			In: []interface{}{bw.MaxInt32},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(int64(bw.MaxInt32))": {
			In: []interface{}{int64(bw.MaxInt32)},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(bw.MaxInt64)": {
			In: []interface{}{int64(bw.MaxInt64)},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} { return test.In[0] },
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
				func(test bwtesting.Case) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(uint8)": {
			In: []interface{}{bw.MaxUint8},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(uint16)": {
			In: []interface{}{bw.MaxUint16},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(uint32)": {
			In: []interface{}{bw.MaxUint32},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(uint64(bw.MaxInt32))": {
			In: []interface{}{uint64(bw.MaxUint32)},
			Out: []interface{}{
				func(test bwtesting.Case) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
		"Int(bw.MaxUint64)": {
			In: []interface{}{bw.MaxUint64},
			Out: []interface{}{
				func(test bwtesting.Case) (result interface{}) {
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
				func(test bwtesting.Case) (result interface{}) {
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
				func(test bwtesting.Case) interface{} { return test.In[0] },
				bwval.ValInt,
			},
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "Bool")
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
			map[string]bwtesting.Case{
				v: {
					In:  []interface{}{k},
					Out: []interface{}{fmt.Sprintf("%q", v)},
				},
			},
		)
	}
}

func TestMustMap(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"Map": {
			In: []interface{}{
				map[string]interface{}{},
			},
			Out: []interface{}{
				map[string]interface{}{},
			},
		},
		"non Map": {
			In: []interface{}{
				1,
			},
			Out: []interface{}{
				(map[string]interface{})(nil),
			},
			Panic: "\x1b[96;1m(int)1\x1b[0m is not \x1b[97;1mMap\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, bwval.MustMap, tests)
}

func TestMustArray(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"Array": {
			In: []interface{}{
				[]interface{}{},
			},
			Out: []interface{}{
				[]interface{}{},
			},
		},
		"non Array": {
			In: []interface{}{
				1,
			},
			Out: []interface{}{
				([]interface{})(nil),
			},
			Panic: "\x1b[96;1m(int)1\x1b[0m is not \x1b[97;1mArray\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, bwval.MustArray, tests)
}

func TestMustString(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"String": {
			In: []interface{}{
				"some",
			},
			Out: []interface{}{
				"some",
			},
		},
		"non String": {
			In: []interface{}{
				1,
			},
			Out: []interface{}{
				"",
			},
			Panic: "\x1b[96;1m(int)1\x1b[0m is not \x1b[97;1mString\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, bwval.MustString, tests)
}

func TestMustInt(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"Int": {
			In: []interface{}{
				273,
			},
			Out: []interface{}{
				273,
				// nil,
			},
		},
		"non Int": {
			In: []interface{}{
				"some",
			},
			Out: []interface{}{
				0,
			},
			Panic: "\x1b[96;1m(string)some\x1b[0m is not \x1b[97;1mInt\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, bwval.MustInt, tests)
}

func TestMustBool(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"Bool": {
			In: []interface{}{
				false,
			},
			Out: []interface{}{
				false,
				// nil,
			},
		},
		"non Bool": {
			In: []interface{}{
				"some",
			},
			Out: []interface{}{
				false,
			},
			Panic: "\x1b[96;1m(string)some\x1b[0m is not \x1b[97;1mBool\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, bwval.MustBool, tests)
}

// func TestPathFrom(t *testing.T) {
// 	tests := map[string]bwtesting.Case{
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
// 			map[string]bwtesting.Case{
// 				test.eta: {
// 					In:  []interface{}{},
// 					Out: []interface{}{test.eta},
// 				},
// 			},
// 		)
// 	}
// }
