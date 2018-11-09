package bwval_test

import (
	"fmt"
	"testing"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwval"
	"github.com/baza-winner/bwcore/defvalid/deftype"
)

func TestMustSetPathVal(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"keyA.keyB": {
			In: []interface{}{
				"something",
				bwval.FromVal(map[string]interface{}{
					"keyA": map[string]interface{}{},
				}),
				func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
			},
			Out: []interface{}{
				bwval.FromVal(map[string]interface{}{
					"keyA": map[string]interface{}{
						"keyB": "something",
					},
				}),
				nil,
			},
		},
		"2.1": {
			In: []interface{}{
				"good",
				bwval.FromVal([]interface{}{
					"string",
					273,
					[]interface{}{"some", "thing"},
				}),
				func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
				// bwval.PathFrom("2.1"),
			},
			Out: []interface{}{
				bwval.FromVal([]interface{}{
					"string",
					273,
					[]interface{}{"some", "good"},
				}),
				nil,
			},
		},
		"2.{$idx}": {
			In: []interface{}{
				"good",
				bwval.FromVal([]interface{}{
					"string",
					273,
					[]interface{}{"some", "thing"},
				}),
				func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
				// bwval.PathFrom("2.{$idx}"),
				map[string]interface{}{
					"idx": 1,
				},
			},
			Out: []interface{}{
				bwval.FromVal([]interface{}{
					"string",
					273,
					[]interface{}{"some", "good"},
				}),
				map[string]interface{}{
					"idx": 1,
				},
			},
		},
		"2.{0}": {
			In: []interface{}{
				"good",
				bwval.FromVal([]interface{}{
					1,
					"string",
					[]interface{}{"some", "thing"},
				}),
				func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
			},
			Out: []interface{}{
				bwval.FromVal([]interface{}{
					1,
					"string",
					[]interface{}{"some", "good"},
				}),
				nil,
			},
		},
		"2.{0.idx}": {
			In: []interface{}{
				"good",
				bwval.FromVal([]interface{}{
					map[string]interface{}{"idx": 1},
					"string",
					[]interface{}{"some", "thing"},
				}),
				func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
			},
			Out: []interface{}{
				bwval.FromVal([]interface{}{
					map[string]interface{}{"idx": 1},
					"string",
					[]interface{}{"some", "good"},
				}),
				nil,
			},
		},
		".": {
			In: []interface{}{
				"good",
				bwval.FromVal(nil),
				func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
			},
			Out: []interface{}{
				bwval.FromVal("good"),
				nil,
			},
		},
		"2.#": {
			In: []interface{}{
				"good",
				bwval.FromVal([]interface{}{
					map[string]interface{}{"idx": 1},
					"string",
					[]interface{}{"some", "thing"},
				}),
				func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
			},
			Out: []interface{}{
				bwval.FromVal(nil),
				nil,
			},
			Panic: "Failed to set \x1b[38;5;252;1m2.#\x1b[0m of \x1b[96;1m[\n  {\n    \"idx\": 1\n  },\n  \"string\",\n  [\n    \"some\",\n    \"thing\"\n  ]\n]\x1b[0m: \x1b[38;5;252;1m2.#\x1b[0m is \x1b[91;1mreadonly path\x1b[0m\x1b[0m",
		},
		"1.nonMapKey.some": {
			In: []interface{}{
				"good",
				bwval.FromVal([]interface{}{
					map[string]interface{}{"idx": 1},
					"string",
					[]interface{}{"some", "thing"},
				}),
				func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
			},
			Out: []interface{}{
				bwval.FromVal(nil),
				nil,
			},
			Panic: "Failed to set \x1b[38;5;252;1m1.nonMapKey.some\x1b[0m of \x1b[96;1m[\n  {\n    \"idx\": 1\n  },\n  \"string\",\n  [\n    \"some\",\n    \"thing\"\n  ]\n]\x1b[0m: \x1b[38;5;252;1m1.nonMapKey\x1b[0m (\x1b[96;1m\"string\"\x1b[0m)\x1b[0m is not \x1b[97;1mMap\x1b[0m\x1b[0m",
		},
		"(nil).some": {
			In: []interface{}{
				"good",
				bwval.FromVal(nil),
				bwval.PathFrom("some"),
			},
			Out: []interface{}{
				bwval.FromVal(nil),
				nil,
			},
			Panic: "Failed to set \x1b[38;5;252;1msome\x1b[0m of \x1b[96;1mnull\x1b[0m: \x1b[38;5;252;1m.\x1b[0m is \x1b[91;1mnil\x1b[0m\x1b[0m",
		},
		"err: nor Int, neither String": {
			In: []interface{}{
				"good",
				bwval.FromVal(map[string]interface{}{"some": 1}),
				bwval.PathFrom("some.{$idx}"),
				map[string]interface{}{"idx": nil},
			},
			Out: []interface{}{
				bwval.FromVal(nil),
				nil,
			},
			Panic: "Failed to set \x1b[38;5;252;1msome.{$idx}\x1b[0m of \x1b[96;1m{\n  \"some\": 1\n}\x1b[0m with \x1b[38;5;201;1mvars\x1b[0m \x1b[96;1m{\n  \"idx\": null\n}\x1b[0m: \x1b[38;5;252;1m$idx\x1b[0m (\x1b[96;1mnull\x1b[0m)\x1b[0m is none of \x1b[97;1mInt\x1b[0m or \x1b[97;1mString\x1b[0m\x1b[0m",
		},
		"$arr.{some}": {
			In: []interface{}{
				"good",
				bwval.FromVal(map[string]interface{}{"some": 1}),
				bwval.PathFrom("$arr.{some}"),
				map[string]interface{}{"arr": []interface{}{"some", "thing"}},
			},
			Out: []interface{}{
				bwval.FromVal(map[string]interface{}{"some": 1}),
				map[string]interface{}{"arr": []interface{}{"some", "good"}},
			},
		},
		"$arr": {
			In: []interface{}{
				"good",
				bwval.FromVal(map[string]interface{}{"some": 1}),
				bwval.PathFrom("$arr.{some}"),
				map[string]interface{}{"arr": []interface{}{"some", "thing"}},
			},
			Out: []interface{}{
				bwval.FromVal(map[string]interface{}{"some": 1}),
				map[string]interface{}{"arr": []interface{}{"some", "good"}},
			},
		},
		"$arr (vars is nil)": {
			In: []interface{}{
				"good",
				bwval.FromVal(map[string]interface{}{"some": 1}),
				bwval.PathFrom("$arr.{some}"),
				// map[string]interface{}{"arr": []interface{}{"some", "thing"}},
			},
			Out: []interface{}{
				bwval.FromVal(nil),
				nil,
			},
			Panic: "Failed to set \x1b[38;5;252;1m$arr.{some}\x1b[0m of \x1b[96;1m{\n  \"some\": 1\n}\x1b[0m with \x1b[38;5;201;1mvars\x1b[0m \x1b[96;1mnull\x1b[0m: \x1b[38;5;201;1mvars\x1b[0m is \x1b[91;1mnil\x1b[0m\x1b[0m",
		},
		"valAtPathIsNil": {
			In: []interface{}{
				"good",
				bwval.FromVal(map[string]interface{}{"some": []interface{}{0}}),
				bwval.PathFrom("some.1.key"),
			},
			Out: []interface{}{
				bwval.FromVal(nil),
				nil,
			},
			Panic: "Failed to set \x1b[38;5;252;1msome.1.key\x1b[0m of \x1b[96;1m{\n  \"some\": [\n    0\n  ]\n}\x1b[0m: \x1b[38;5;252;1msome.1\x1b[0m is \x1b[91;1mnil\x1b[0m\x1b[0m",
		},
		"ansiValAtPathHasNotEnoughRange": {
			In: []interface{}{
				"good",
				bwval.FromVal(map[string]interface{}{"some": []interface{}{0}}),
				bwval.PathFrom("some.1"),
			},
			Out: []interface{}{
				bwval.FromVal(nil),
				nil,
			},
			Panic: "Failed to set \x1b[38;5;252;1msome.1\x1b[0m of \x1b[96;1m{\n  \"some\": [\n    0\n  ]\n}\x1b[0m: \x1b[38;5;252;1msome\x1b[0m (\x1b[96;1m[\n  0\n]\x1b[0m)\x1b[0m has not enough length (\x1b[96;1m1\x1b[0m) for idx (\x1b[96;1m1)\x1b[0m\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "$arr.{some}")
	bwtesting.BwRunTests(t, MustSetPathValWrapper, tests)
}

func MustSetPathValWrapper(val interface{}, v bw.Val, path bw.ValPath, optVars ...map[string]interface{}) (bw.Val, map[string]interface{}) {
	bwval.MustSetPathVal(val, v, path, optVars...)
	var vars map[string]interface{}
	if len(optVars) > 0 {
		vars = optVars[0]
	}
	return v, vars
}

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
				bwval.PathFrom("some.{$key}"),
				map[string]interface{}{"key": "thing"},
			},
			Out:   []interface{}{nil},
			Panic: "Failed to get \x1b[38;5;252;1msome.{$key}\x1b[0m of \x1b[96;1m1\x1b[0m with \x1b[38;5;201;1mvars\x1b[0m \x1b[96;1m{\n  \"key\": \"thing\"\n}\x1b[0m: \x1b[38;5;252;1msome\x1b[0m (\x1b[96;1m1\x1b[0m)\x1b[0m is not \x1b[97;1mMap\x1b[0m\x1b[0m",
		},
		"err: is not Array": {
			In: []interface{}{
				bwval.FromVal(1),
				bwval.PathFrom("{$idx}"),
				map[string]interface{}{"idx": 1},
			},
			Out:   []interface{}{nil},
			Panic: "Failed to get \x1b[38;5;252;1m{$idx}\x1b[0m of \x1b[96;1m1\x1b[0m with \x1b[38;5;201;1mvars\x1b[0m \x1b[96;1m{\n  \"idx\": 1\n}\x1b[0m: \x1b[38;5;252;1m{$idx}\x1b[0m (\x1b[96;1m1\x1b[0m)\x1b[0m is not \x1b[97;1mArray\x1b[0m\x1b[0m",
		},
		"err: nor Map, neither Array": {
			In: []interface{}{
				bwval.FromVal(1),
				bwval.PathFrom("#"),
			},
			Out:   []interface{}{nil},
			Panic: "Failed to get \x1b[38;5;252;1m#\x1b[0m of \x1b[96;1m1\x1b[0m: \x1b[38;5;252;1m.\x1b[0m (\x1b[96;1m1\x1b[0m)\x1b[0m is none of \x1b[97;1mMap\x1b[0m or \x1b[97;1mArray\x1b[0m\x1b[0m",
		},
		"err: nor Int, neither String": {
			In: []interface{}{
				bwval.FromVal(map[string]interface{}{"some": 1}),
				bwval.PathFrom("some.{$idx}"),
				map[string]interface{}{"idx": nil},
			},
			Out:   []interface{}{nil},
			Panic: "Failed to get \x1b[38;5;252;1msome.{$idx}\x1b[0m of \x1b[96;1m{\n  \"some\": 1\n}\x1b[0m with \x1b[38;5;201;1mvars\x1b[0m \x1b[96;1m{\n  \"idx\": null\n}\x1b[0m: \x1b[38;5;252;1m$idx\x1b[0m (\x1b[96;1mnull\x1b[0m)\x1b[0m is none of \x1b[97;1mInt\x1b[0m or \x1b[97;1mString\x1b[0m\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "err: nor Map, neither Array")
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

func TestDefFrom(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"nil": {
			In:    []interface{}{func(testName string) interface{} { return bwval.MustPathVal(bwval.From(testName), bwval.PathFrom(".")) }},
			Out:   []interface{}{bwval.Def{}},
			Panic: "\x1b[38;5;252;1m$def\x1b[0m is \x1b[91;1mnil\x1b[0m",
		},
		"true": {
			In:    []interface{}{func(testName string) interface{} { return bwval.MustPathVal(bwval.From(testName), bwval.PathFrom(".")) }},
			Out:   []interface{}{bwval.Def{}},
			Panic: "\x1b[38;5;252;1m$def\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is none of \x1b[97;1mString\x1b[0m, \x1b[97;1mArray\x1b[0m or \x1b[97;1mMap\x1b[0m",
		},
		"[ Bool true ]": {
			In:    []interface{}{func(testName string) interface{} { return bwval.MustPathVal(bwval.From(testName), bwval.PathFrom(".")) }},
			Out:   []interface{}{bwval.Def{}},
			Panic: "\x1b[38;5;252;1m$def.1\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mString\x1b[0m",
		},
		"{type true}": {
			In:    []interface{}{func(testName string) interface{} { return bwval.MustPathVal(bwval.From(testName), bwval.PathFrom(".")) }},
			Out:   []interface{}{bwval.Def{}},
			Panic: "\x1b[38;5;252;1m$def.type\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is none of \x1b[97;1mString\x1b[0m or \x1b[97;1mArray\x1b[0m",
		},
		"Bool": {
			In:  []interface{}{func(testName string) interface{} { return bwval.MustPathVal(bwval.From(testName), bwval.PathFrom(".")) }},
			Out: []interface{}{bwval.Def{Types: deftype.From(deftype.Bool)}},
		},
		`{ type [ Int "bool" ] }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.MustPathVal(bwval.From(testName), bwval.PathFrom(".")) }},
			Out:   []interface{}{bwval.Def{Types: deftype.From(deftype.Bool)}},
			Panic: "\x1b[38;5;252;1m$def.type.1\x1b[0m (\x1b[96;1m\"bool\"\x1b[0m)\x1b[0m is \x1b[91;1mnon supported\x1b[0m value\x1b[0m",
		},
		`{ type String enum <a b c> }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.MustPathVal(bwval.From(testName), bwval.PathFrom(".")) }},
			Out: []interface{}{bwval.Def{
				Types: deftype.From(deftype.String),
				Enum:  bwset.StringFrom("a", "b", "c"),
			}},
			// Panic: "\x1b[38;5;252;1m$def.type.1\x1b[0m (\x1b[96;1m\"bool\"\x1b[0m)\x1b[0m is \x1b[91;1mnon supported\x1b[0m value\x1b[0m",
		},
		`{ type String enum [ "a" true ] }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.MustPathVal(bwval.From(testName), bwval.PathFrom(".")) }},
			Out: []interface{}{bwval.Def{
				Types: deftype.From(deftype.String),
				Enum:  bwset.StringFrom("a", "b", "c"),
			}},
			Panic: "\x1b[38;5;252;1m$def.enum.1\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mString\x1b[0m",
		},
		`{ type Map keys { a Bool } }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.MustPathVal(bwval.From(testName), bwval.PathFrom(".")) }},
			Out: []interface{}{bwval.Def{
				Types: deftype.From(deftype.Map),
				Keys:  map[string]bwval.Def{"a": {Types: deftype.From(deftype.Bool)}},
			}},
		},
		`{ type Array arrayElem Bool }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.MustPathVal(bwval.From(testName), bwval.PathFrom(".")) }},
			Out: []interface{}{bwval.Def{
				Types:     deftype.From(deftype.Array),
				ArrayElem: &bwval.Def{Types: deftype.From(deftype.Bool)},
			}},
		},
		`{ type Int min 1 max 2 }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.MustPathVal(bwval.From(testName), bwval.PathFrom(".")) }},
			Out: []interface{}{bwval.Def{
				Types: deftype.From(deftype.Int),
				Range: bwval.IntRange{bwval.PtrToInt(1), bwval.PtrToInt(2)},
				// ArrayElem: &bwval.Def{Types: deftype.From(deftype.Bool)},
			}},
		},
		`{ type Number min 1 max 2 }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.MustPathVal(bwval.From(testName), bwval.PathFrom(".")) }},
			Out: []interface{}{bwval.Def{
				Types: deftype.From(deftype.Number),
				Range: bwval.NumberRange{bwval.PtrToNumber(1), bwval.PtrToNumber(2)},
				// ArrayElem: &bwval.Def{Types: deftype.From(deftype.Bool)},
			}},
		},
		`{ type Number min 1 max 2 default 3 }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.MustPathVal(bwval.From(testName), bwval.PathFrom(".")) }},
			Out: []interface{}{bwval.Def{
				Types: deftype.From(deftype.Number),
				Range: bwval.NumberRange{bwval.PtrToNumber(1), bwval.PtrToNumber(2)},
			}},
			Panic: "\x1b[38;5;252;1m$def.default\x1b[0m (\x1b[96;1m3\x1b[0m)\x1b[0m is \x1b[91;1mout of range\x1b[0m \x1b[96;1m1..2\x1b[0m",
		},
		`{ type String default "some" }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.MustPathVal(bwval.From(testName), bwval.PathFrom(".")) }},
			Out: []interface{}{bwval.Def{
				Types:      deftype.From(deftype.String),
				Default:    "some",
				IsOptional: true,
				// ArrayElem: &bwval.Def{Types: deftype.From(deftype.Bool)},
			}},
		},
		// "[ Bool true ]": {
		// 	In:    []interface{}{bwval.From("[ Bool true ]")},
		// 	Out:   []interface{}{bwval.Def{}},
		// 	Panic: "",
		// },
		// "Bool": {
		// 	In: []interface{}{
		// 		false,
		// 	},
		// 	Out: []interface{}{
		// 		false,
		// 		// nil,
		// 	},
		// },
		// "non Bool": {
		// 	In: []interface{}{
		// 		"some",
		// 	},
		// 	Out: []interface{}{
		// 		false,
		// 	},
		// 	Panic: "\x1b[96;1m(string)some\x1b[0m is not \x1b[97;1mBool\x1b[0m",
		// },
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "{ type [ Int \"bool\" ] }")
	bwtesting.BwRunTests(t, bwval.DefFrom, tests)
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
