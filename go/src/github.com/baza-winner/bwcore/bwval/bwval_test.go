package bwval_test

import (
	"testing"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwval"
)

func TestMustSetPathVal(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"keyA.keyB": {
			In: []interface{}{
				"something",
				bwval.Holder{Val: map[string]interface{}{
					"keyA": map[string]interface{}{},
				}},
				func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
			},
			Out: []interface{}{
				bwval.Holder{Val: map[string]interface{}{
					"keyA": map[string]interface{}{
						"keyB": "something",
					},
				}},
				nil,
			},
		},
		"2.1": {
			In: []interface{}{
				"good",
				bwval.Holder{Val: []interface{}{
					"string",
					273,
					[]interface{}{"some", "thing"},
				}},
				func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
				// bwval.PathFrom("2.1"),
			},
			Out: []interface{}{
				bwval.Holder{Val: []interface{}{
					"string",
					273,
					[]interface{}{"some", "good"},
				}},
				nil,
			},
		},
		"2.{$idx}": {
			In: []interface{}{
				"good",
				bwval.Holder{Val: []interface{}{
					"string",
					273,
					[]interface{}{"some", "thing"},
				}},
				func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
				// bwval.PathFrom("2.{$idx}"),
				map[string]interface{}{
					"idx": 1,
				},
			},
			Out: []interface{}{
				bwval.Holder{Val: []interface{}{
					"string",
					273,
					[]interface{}{"some", "good"},
				}},
				map[string]interface{}{
					"idx": 1,
				},
			},
		},
		"2.{0}": {
			In: []interface{}{
				"good",
				bwval.Holder{Val: []interface{}{
					1,
					"string",
					[]interface{}{"some", "thing"},
				}},
				func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
			},
			Out: []interface{}{
				bwval.Holder{Val: []interface{}{
					1,
					"string",
					[]interface{}{"some", "good"},
				}},
				nil,
			},
		},
		"2.{0.idx}": {
			In: []interface{}{
				"good",
				bwval.Holder{Val: []interface{}{
					map[string]interface{}{"idx": 1},
					"string",
					[]interface{}{"some", "thing"},
				}},
				func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
			},
			Out: []interface{}{
				bwval.Holder{Val: []interface{}{
					map[string]interface{}{"idx": 1},
					"string",
					[]interface{}{"some", "good"},
				}},
				nil,
			},
		},
		".": {
			In: []interface{}{
				"good",
				bwval.Holder{},
				func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
			},
			Out: []interface{}{
				bwval.Holder{Val: "good"},
				nil,
			},
		},
		"2.#": {
			In: []interface{}{
				"good",
				bwval.Holder{Val: []interface{}{
					map[string]interface{}{"idx": 1},
					"string",
					[]interface{}{"some", "thing"},
				}},
				func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
			},
			Out: []interface{}{
				bwval.Holder{},
				nil,
			},
			Panic: "Failed to set \x1b[38;5;252;1m2.#\x1b[0m of \x1b[96;1m[\n  {\n    \"idx\": 1\n  },\n  \"string\",\n  [\n    \"some\",\n    \"thing\"\n  ]\n]\x1b[0m: \x1b[38;5;252;1m2.#\x1b[0m is \x1b[91;1mreadonly path\x1b[0m\x1b[0m",
		},
		"1.nonMapKey.some": {
			In: []interface{}{
				"good",
				bwval.Holder{Val: []interface{}{
					map[string]interface{}{"idx": 1},
					"string",
					[]interface{}{"some", "thing"},
				}},
				func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
			},
			Out: []interface{}{
				bwval.Holder{},
				nil,
			},
			Panic: "Failed to set \x1b[38;5;252;1m1.nonMapKey.some\x1b[0m of \x1b[96;1m[\n  {\n    \"idx\": 1\n  },\n  \"string\",\n  [\n    \"some\",\n    \"thing\"\n  ]\n]\x1b[0m: \x1b[38;5;252;1m1.nonMapKey\x1b[0m (\x1b[96;1m\"string\"\x1b[0m)\x1b[0m is not \x1b[97;1mMap\x1b[0m\x1b[0m",
		},
		"(nil).some": {
			In: []interface{}{
				"good",
				bwval.Holder{},
				bwval.PathFrom("some"),
			},
			Out: []interface{}{
				bwval.Holder{},
				nil,
			},
			Panic: "Failed to set \x1b[38;5;252;1msome\x1b[0m of \x1b[96;1mnull\x1b[0m: \x1b[38;5;252;1m.\x1b[0m is \x1b[91;1m(interface {})<nil>\x1b[0m\x1b[0m",
		},
		"err: nor Int, neither String": {
			In: []interface{}{
				"good",
				bwval.Holder{Val: map[string]interface{}{"some": 1}},
				bwval.PathFrom("some.{$idx}"),
				map[string]interface{}{"idx": nil},
			},
			Out: []interface{}{
				bwval.Holder{},
				nil,
			},
			Panic: "Failed to set \x1b[38;5;252;1msome.{$idx}\x1b[0m of \x1b[96;1m{\n  \"some\": 1\n}\x1b[0m with \x1b[38;5;201;1mvars\x1b[0m \x1b[96;1m{\n  \"idx\": null\n}\x1b[0m: \x1b[38;5;252;1m$idx\x1b[0m (\x1b[96;1mnull\x1b[0m)\x1b[0m is none of \x1b[97;1mInt\x1b[0m or \x1b[97;1mString\x1b[0m\x1b[0m",
		},
		"$arr.{some}": {
			In: []interface{}{
				"good",
				bwval.Holder{Val: map[string]interface{}{"some": 1}},
				bwval.PathFrom("$arr.{some}"),
				map[string]interface{}{"arr": []interface{}{"some", "thing"}},
			},
			Out: []interface{}{
				bwval.Holder{Val: map[string]interface{}{"some": 1}},
				map[string]interface{}{"arr": []interface{}{"some", "good"}},
			},
		},
		"$arr": {
			In: []interface{}{
				"good",
				bwval.Holder{Val: map[string]interface{}{"some": 1}},
				bwval.PathFrom("$arr.{some}"),
				map[string]interface{}{"arr": []interface{}{"some", "thing"}},
			},
			Out: []interface{}{
				bwval.Holder{Val: map[string]interface{}{"some": 1}},
				map[string]interface{}{"arr": []interface{}{"some", "good"}},
			},
		},
		"$arr (vars is nil)": {
			In: []interface{}{
				"good",
				bwval.Holder{Val: map[string]interface{}{"some": 1}},
				bwval.PathFrom("$arr.{some}"),
				// map[string]interface{}{"arr": []interface{}{"some", "thing"}},
			},
			Out: []interface{}{
				bwval.Holder{},
				nil,
			},
			Panic: "Failed to set \x1b[38;5;252;1m$arr.{some}\x1b[0m of \x1b[96;1m{\n  \"some\": 1\n}\x1b[0m with \x1b[38;5;201;1mvars\x1b[0m \x1b[96;1mnull\x1b[0m: \x1b[38;5;201;1mvars\x1b[0m is \x1b[91;1mnil\x1b[0m\x1b[0m",
		},
		"valAtPathIsNil": {
			In: []interface{}{
				"good",
				bwval.Holder{Val: map[string]interface{}{"some": []interface{}{0}}},
				bwval.PathFrom("some.1.key"),
			},
			Out: []interface{}{
				bwval.Holder{},
				nil,
			},
			Panic: "Failed to set \x1b[38;5;252;1msome.1.key\x1b[0m of \x1b[96;1m{\n  \"some\": [\n    0\n  ]\n}\x1b[0m: \x1b[38;5;252;1msome.1\x1b[0m is \x1b[91;1m(interface {})<nil>\x1b[0m\x1b[0m",
		},
		"ansiValAtPathHasNotEnoughRange": {
			In: []interface{}{
				"good",
				bwval.Holder{Val: map[string]interface{}{"some": []interface{}{0}}},
				bwval.PathFrom("some.1"),
			},
			Out: []interface{}{
				bwval.Holder{},
				nil,
			},
			Panic: "Failed to set \x1b[38;5;252;1msome.1\x1b[0m of \x1b[96;1m{\n  \"some\": [\n    0\n  ]\n}\x1b[0m: \x1b[38;5;252;1msome\x1b[0m (\x1b[96;1m[\n  0\n]\x1b[0m)\x1b[0m has not enough length (\x1b[96;1m1\x1b[0m) for idx (\x1b[96;1m1)\x1b[0m\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "$arr.{some}")
	bwtesting.BwRunTests(t, MustSetPathValWrapper, tests)
}

func MustSetPathValWrapper(val interface{}, v bwval.Holder, path bw.ValPath, optVars ...map[string]interface{}) (bwval.Holder, map[string]interface{}) {
	bwval.MustSetPathVal(val, &v, path, optVars...)
	var vars map[string]interface{}
	if len(optVars) > 0 {
		vars = optVars[0]
	}
	return v, vars
}

func TestMustPathVal(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"nil": {
			In:  []interface{}{bwval.Holder{}, bwval.PathFrom("some")},
			Out: []interface{}{nil},
		},
		"self": {
			In:  []interface{}{bwval.Holder{Val: 1}, bwval.PathFrom(".")},
			Out: []interface{}{1},
		},
		"by key": {
			In: []interface{}{
				bwval.Holder{Val: map[string]interface{}{"some": "thing"}},
				bwval.PathFrom("some"),
			},
			Out: []interface{}{"thing"},
		},
		"by idx (1)": {
			In: []interface{}{
				bwval.Holder{Val: []interface{}{"some", "thing"}},
				bwval.PathFrom("1"),
			},
			Out: []interface{}{"thing"},
		},
		"by idx (-1)": {
			In: []interface{}{
				bwval.Holder{Val: []interface{}{"some", "thing"}},
				bwval.PathFrom("-1"),
			},
			Out: []interface{}{"thing"},
		},
		"by idx (len)": {
			In: []interface{}{
				bwval.Holder{Val: []interface{}{"some", "thing"}},
				bwval.PathFrom("2"),
			},
			Out: []interface{}{nil},
		},
		"by idx of nil": {
			In: []interface{}{
				bwval.Holder{},
				bwval.PathFrom("2"),
			},
			Out: []interface{}{nil},
		},
		"by # of nil": {
			In: []interface{}{
				bwval.Holder{},
				bwval.PathFrom("1.#"),
			},
			Out: []interface{}{0},
		},
		"by # of Array": {
			In: []interface{}{
				bwval.Holder{Val: []interface{}{"a", "b"}},
				bwval.PathFrom("#"),
			},
			Out: []interface{}{2},
		},
		"by # of Map": {
			In: []interface{}{
				bwval.Holder{Val: []interface{}{
					"a",
					map[string]interface{}{"c": "d", "e": "f", "i": "g"},
				}},
				bwval.PathFrom("1.#"),
			},
			Out: []interface{}{3},
		},
		"by path (idx)": {
			In: []interface{}{
				bwval.Holder{Val: map[string]interface{}{"some": []interface{}{"good", "thing"}, "idx": 1}},
				bwval.PathFrom("some.{idx}"),
			},
			Out: []interface{}{"thing"},
		},
		"by path (key)": {
			In: []interface{}{
				bwval.Holder{Val: map[string]interface{}{"some": []interface{}{"good", "thing"}, "key": "some"}},
				bwval.PathFrom("{key}.1"),
			},
			Out: []interface{}{"thing"},
		},
		"some.{$idx}": {
			In: []interface{}{
				bwval.Holder{Val: map[string]interface{}{"some": []interface{}{"good", "thing"}}},
				bwval.PathFrom("some.{$idx}"),
				map[string]interface{}{"idx": 1},
			},
			Out: []interface{}{"thing"},
		},
		"$idx": {
			In: []interface{}{
				bwval.Holder{Val: map[string]interface{}{"some": []interface{}{"good", "thing"}}},
				bwval.PathFrom("$idx"),
			},
			Out: []interface{}{nil},
		},
		"err: is not Map": {
			In: []interface{}{
				bwval.Holder{Val: 1},
				bwval.PathFrom("some.{$key}"),
				map[string]interface{}{"key": "thing"},
			},
			Out:   []interface{}{nil},
			Panic: "Failed to get \x1b[38;5;252;1msome.{$key}\x1b[0m of \x1b[96;1m1\x1b[0m with \x1b[38;5;201;1mvars\x1b[0m \x1b[96;1m{\n  \"key\": \"thing\"\n}\x1b[0m: \x1b[38;5;252;1msome\x1b[0m (\x1b[96;1m1\x1b[0m)\x1b[0m is not \x1b[97;1mMap\x1b[0m\x1b[0m",
		},
		"err: is not Array": {
			In: []interface{}{
				bwval.Holder{Val: 1},
				bwval.PathFrom("{$idx}"),
				map[string]interface{}{"idx": 1},
			},
			Out:   []interface{}{nil},
			Panic: "Failed to get \x1b[38;5;252;1m{$idx}\x1b[0m of \x1b[96;1m1\x1b[0m with \x1b[38;5;201;1mvars\x1b[0m \x1b[96;1m{\n  \"idx\": 1\n}\x1b[0m: \x1b[38;5;252;1m{$idx}\x1b[0m (\x1b[96;1m1\x1b[0m)\x1b[0m is not \x1b[97;1mArray\x1b[0m\x1b[0m",
		},
		"err: nor Map, neither Array": {
			In: []interface{}{
				bwval.Holder{Val: 1},
				bwval.PathFrom("#"),
			},
			Out:   []interface{}{nil},
			Panic: "Failed to get \x1b[38;5;252;1m#\x1b[0m of \x1b[96;1m1\x1b[0m: \x1b[38;5;252;1m.\x1b[0m (\x1b[96;1m1\x1b[0m)\x1b[0m is none of \x1b[97;1mMap\x1b[0m or \x1b[97;1mArray\x1b[0m\x1b[0m",
		},
		"err: nor Int, neither String": {
			In: []interface{}{
				bwval.Holder{Val: map[string]interface{}{"some": 1}},
				bwval.PathFrom("some.{$idx}"),
				map[string]interface{}{"idx": nil},
			},
			Out:   []interface{}{nil},
			Panic: "Failed to get \x1b[38;5;252;1msome.{$idx}\x1b[0m of \x1b[96;1m{\n  \"some\": 1\n}\x1b[0m with \x1b[38;5;201;1mvars\x1b[0m \x1b[96;1m{\n  \"idx\": null\n}\x1b[0m: \x1b[38;5;252;1m$idx\x1b[0m (\x1b[96;1mnull\x1b[0m)\x1b[0m is none of \x1b[97;1mInt\x1b[0m or \x1b[97;1mString\x1b[0m\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "by key")
	bwtesting.BwRunTests(t, MustPathValWrapper, tests)
}

func MustPathValWrapper(v bwval.Holder, path bw.ValPath, optVars ...map[string]interface{}) interface{} {
	return bwval.MustPathVal(&v, path, optVars...)
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

func TestMustNumber(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"Number(float32)": {
			In: []interface{}{
				float32(273),
			},
			Out: []interface{}{
				float64(273),
				// nil,
			},
		},
		"Number(float64)": {
			In: []interface{}{
				float64(273),
			},
			Out: []interface{}{
				float64(273),
				// nil,
			},
		},
		"non Number": {
			In: []interface{}{
				"some",
			},
			Out: []interface{}{
				float64(0),
			},
			Panic: "\x1b[96;1m(string)some\x1b[0m is not \x1b[97;1mNumber\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, bwval.MustNumber, tests)
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
