package bwval_test

import (
	"testing"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwtype"
	"github.com/baza-winner/bwcore/bwval"
)

func TestPathFrom(t *testing.T) {
	bwtesting.BwRunTests(t, bwval.PathFrom,
		map[string]bwtesting.Case{
			"some.key(.2)": {
				In: []interface{}{".2", []bw.ValPath{bwval.PathFrom("some.key")}},
				Out: []interface{}{
					bwval.PathFrom("some.key.2"),
				},
			},
			"some.key(.#.2)": {
				In:    []interface{}{".#.2", []bw.ValPath{bwval.PathFrom("some.key")}},
				Panic: "unexpected char \x1b[96;1m'.'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m46\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m2\x1b[0m: \x1b[32m.#\x1b[91m.\x1b[0m2\n",
			},
		},
	)
}

func TestMustArrayOfString(t *testing.T) {
	bwtesting.BwRunTests(t, bwval.MustArrayOfString,
		map[string]bwtesting.Case{
			`<abc def>`: {
				In: []interface{}{
					func(testName string) interface{} { return bwval.From(testName) },
				},
				Out: []interface{}{
					[]string{"abc", "def"},
				},
			},
			`["abc" "def"]`: {
				In: []interface{}{
					func(testName string) interface{} { return bwval.From(testName) },
				},
				Out: []interface{}{
					[]string{"abc", "def"},
				},
			},
			`["abc" true]`: {
				In: []interface{}{
					func(testName string) interface{} { return bwval.From(testName) },
				},
				Panic: "\x1b[96;1m([]interface {})[(string)abc (bool)true]\x1b[0m is not \x1b[97;1mArrayOfString\x1b[0m",
			},
		},
	)
}

func TestFrom(t *testing.T) {
	bwtesting.BwRunTests(t, bwval.From,
		map[string]bwtesting.Case{
			`{{$a}}`: {
				In: []interface{}{
					func(testName string) string { return testName },
					map[string]interface{}{
						"a": "valueC",
					},
				},
				Out: []interface{}{
					"valueC",
				},
			},
			`{ keyA: "valueA" keyB: [ "valueB" {{keyA}} ] keyC: {{$a}}}`: {
				In: []interface{}{
					func(testName string) string { return testName },
					map[string]interface{}{
						"a": "valueC",
					},
				},
				Out: []interface{}{
					map[string]interface{}{
						"keyA": "valueA",
						"keyB": []interface{}{"valueB", "valueA"},
						"keyC": "valueC",
					},
				},
			},
			`} `: {
				In: []interface{}{
					func(testName string) string { return testName },
				},
				Panic: "unexpected char \x1b[96;1m'}'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m125\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\x1b[91m}\x1b[0m \n",
			},
			`{ key $a.1 }`: {
				In: []interface{}{
					func(testName string) string { return testName },
				},
				Panic: "var \x1b[38;5;201;1ma\x1b[0m is not defined\x1b[0m",
			},
			`[ $a.1 ]`: {
				In: []interface{}{
					func(testName string) string { return testName },
				},
				Panic: "var \x1b[38;5;201;1ma\x1b[0m is not defined\x1b[0m",
			},
		},
	)
}

func TestMustSetPathVal(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(val interface{}, v bwval.Holder, path bw.ValPath, optVars ...map[string]interface{}) (bwval.Holder, map[string]interface{}) {
			bwval.MustSetPathVal(val, &v, path, optVars...)
			var vars map[string]interface{}
			if len(optVars) > 0 {
				vars = optVars[0]
			}
			return v, vars
		},
		map[string]bwtesting.Case{
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
			"2.($idx)": {
				In: []interface{}{
					"good",
					bwval.Holder{Val: []interface{}{
						"string",
						273,
						[]interface{}{"some", "thing"},
					}},
					func(testName string) bw.ValPath { return bwval.PathFrom(testName) },
					// bwval.PathFrom("2.($idx)"),
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
			"2.(0)": {
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
			"2.(0.idx)": {
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
				Panic: "Failed to set \x1b[38;5;252;1m1.nonMapKey.some\x1b[0m of \x1b[96;1m[\n  {\n    \"idx\": 1\n  },\n  \"string\",\n  [\n    \"some\",\n    \"thing\"\n  ]\n]\x1b[0m: \x1b[38;5;252;1m1.nonMapKey\x1b[0m (\x1b[96;1m\"string\"\x1b[0m)\x1b[0m is not \x1b[97;1mMap\x1b[0m\x1b[0m",
			},
			"(nil).some": {
				In: []interface{}{
					"good",
					bwval.Holder{},
					bwval.PathFrom("some"),
				},
				Panic: "Failed to set \x1b[38;5;252;1msome\x1b[0m of \x1b[96;1mnull\x1b[0m: \x1b[38;5;252;1m.\x1b[0m is \x1b[91;1m(interface {})<nil>\x1b[0m\x1b[0m",
			},
			"err: neither Int nor String": {
				In: []interface{}{
					"good",
					bwval.Holder{Val: map[string]interface{}{"some": 1}},
					bwval.PathFrom("some.($idx)"),
					map[string]interface{}{"idx": nil},
				},
				Panic: "Failed to set \x1b[38;5;252;1msome.($idx)\x1b[0m of \x1b[96;1m{\n  \"some\": 1\n}\x1b[0m with \x1b[38;5;201;1mvars\x1b[0m \x1b[96;1m{\n  \"idx\": null\n}\x1b[0m: \x1b[38;5;252;1m$idx\x1b[0m (\x1b[96;1mnull\x1b[0m)\x1b[0m neither \x1b[97;1mInt\x1b[0m nor \x1b[97;1mString\x1b[0m\x1b[0m",
			},
			"$arr.(some)": {
				In: []interface{}{
					"good",
					bwval.Holder{Val: map[string]interface{}{"some": 1}},
					bwval.PathFrom("$arr.(some)"),
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
					bwval.PathFrom("$arr.(some)"),
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
					bwval.PathFrom("$arr.(some)"),
					// map[string]interface{}{"arr": []interface{}{"some", "thing"}},
				},
				Panic: "Failed to set \x1b[38;5;252;1m$arr.(some)\x1b[0m of \x1b[96;1m{\n  \"some\": 1\n}\x1b[0m with \x1b[38;5;201;1mvars\x1b[0m \x1b[96;1mnull\x1b[0m: \x1b[38;5;201;1mvars\x1b[0m is \x1b[91;1mnil\x1b[0m\x1b[0m",
			},
			"err: some.1.key": {
				In: []interface{}{
					"good",
					bwval.Holder{Val: map[string]interface{}{"some": []interface{}{0}}},
					bwval.PathFrom("some.1.key"),
				},
				Panic: "Failed to set \x1b[38;5;252;1msome.1.key\x1b[0m of \x1b[96;1m{\n  \"some\": [\n    0\n  ]\n}\x1b[0m: \x1b[38;5;252;1msome.1\x1b[0m (\x1b[96;1m[\n  0\n]\x1b[0m)\x1b[0m has not enough length (\x1b[96;1m1\x1b[0m) for idx (\x1b[96;1m1)\x1b[0m\x1b[0m",
			},
			"ansiValAtPathHasNotEnoughRange": {
				In: []interface{}{
					"good",
					bwval.Holder{Val: map[string]interface{}{"some": []interface{}{0}}},
					bwval.PathFrom("some.1"),
				},
				Panic: "Failed to set \x1b[38;5;252;1msome.1\x1b[0m of \x1b[96;1m{\n  \"some\": [\n    0\n  ]\n}\x1b[0m: \x1b[38;5;252;1msome\x1b[0m (\x1b[96;1m[\n  0\n]\x1b[0m)\x1b[0m has not enough length (\x1b[96;1m1\x1b[0m) for idx (\x1b[96;1m1)\x1b[0m\x1b[0m",
			},
		},
	)
}

func TestMustPathVal(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(v bwval.Holder, path bw.ValPath, optVars ...map[string]interface{}) interface{} {
			return bwval.MustPathVal(&v, path, optVars...)
		},
		map[string]bwtesting.Case{
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

			"nil::some?": {
				In:  []interface{}{bwval.Holder{}, bwval.PathFrom("some?")},
				Out: []interface{}{nil},
			},
			"nil::2?": {
				In: []interface{}{
					bwval.Holder{},
					bwval.PathFrom("2?"),
				},
				Out: []interface{}{nil},
			},
			"<some thing>::2?": {
				In: []interface{}{
					bwval.Holder{Val: []interface{}{"some", "thing"}},
					bwval.PathFrom("2?"),
				},
				Out: []interface{}{nil},
			},
			"<good thing>::$idx?": {
				In: []interface{}{
					bwval.Holder{Val: map[string]interface{}{"some": []interface{}{"good", "thing"}}},
					bwval.PathFrom("$idx?"),
				},
				Out: []interface{}{nil},
			},

			"<some thing>::2": {
				In: []interface{}{
					bwval.Holder{Val: []interface{}{"some", "thing"}},
					bwval.PathFrom("2"),
				},
				Panic: "Failed to get \x1b[38;5;252;1m2\x1b[0m of \x1b[96;1m[\n  \"some\",\n  \"thing\"\n]\x1b[0m: \x1b[38;5;252;1m.\x1b[0m (\x1b[96;1m[\n  \"some\",\n  \"thing\"\n]\x1b[0m)\x1b[0m has not enough length (\x1b[96;1m2\x1b[0m) for idx (\x1b[96;1m2)\x1b[0m\x1b[0m",
			},
			"nil::2": {
				In: []interface{}{
					bwval.Holder{},
					bwval.PathFrom("2"),
				},
				// Out: []interface{}{nil},
				Panic: "Failed to get \x1b[38;5;252;1m2\x1b[0m of \x1b[96;1mnull\x1b[0m: \x1b[38;5;252;1m.\x1b[0m is \x1b[91;1m(interface {})<nil>\x1b[0m\x1b[0m",
			},
			"<some thing>::$idx": {
				In: []interface{}{
					bwval.Holder{Val: map[string]interface{}{"some": []interface{}{"good", "thing"}}},
					bwval.PathFrom("$idx"),
				},
				// Out: []interface{}{nil},
				Panic: "Failed to get \x1b[38;5;252;1m$idx\x1b[0m of \x1b[96;1m{\n  \"some\": [\n    \"good\",\n    \"thing\"\n  ]\n}\x1b[0m with \x1b[38;5;201;1mvars\x1b[0m \x1b[96;1mnull\x1b[0m: var \x1b[38;5;201;1midx\x1b[0m is not defined\x1b[0m\x1b[0m",
				// Panic: "Failed to get \x1b[38;5;252;1m2\x1b[0m of \x1b[96;1m[\n  \"some\",\n  \"thing\"\n]\x1b[0m: \x1b[38;5;252;1m.\x1b[0m (\x1b[96;1m[\n  \"some\",\n  \"thing\"\n]\x1b[0m)\x1b[0m has not enough length (\x1b[96;1m2\x1b[0m) for idx (\x1b[96;1m2)\x1b[0m\x1b[0m",
			},

			"by # of nil": {
				In: []interface{}{
					bwval.Holder{},
					bwval.PathFrom("1.#"),
				},
				Panic: "Failed to get \x1b[38;5;252;1m1.#\x1b[0m of \x1b[96;1mnull\x1b[0m: \x1b[38;5;252;1m.\x1b[0m is \x1b[91;1m(interface {})<nil>\x1b[0m\x1b[0m",
				// Out: []interface{}{0},
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
					bwval.PathFrom("some.(idx)"),
				},
				Out: []interface{}{"thing"},
			},
			"by path (key)": {
				In: []interface{}{
					bwval.Holder{Val: map[string]interface{}{"some": []interface{}{"good", "thing"}, "key": "some"}},
					bwval.PathFrom("(key).1"),
				},
				Out: []interface{}{"thing"},
			},
			"some.($idx)": {
				In: []interface{}{
					bwval.Holder{Val: map[string]interface{}{"some": []interface{}{"good", "thing"}}},
					bwval.PathFrom("some.($idx)"),
					map[string]interface{}{"idx": 1},
				},
				Out: []interface{}{"thing"},
			},
			"err: is not Map": {
				In: []interface{}{
					bwval.Holder{Val: 1},
					bwval.PathFrom("some.($key)"),
					map[string]interface{}{"key": "thing"},
				},
				Panic: "Failed to get \x1b[38;5;252;1msome.($key)\x1b[0m of \x1b[96;1m1\x1b[0m with \x1b[38;5;201;1mvars\x1b[0m \x1b[96;1m{\n  \"key\": \"thing\"\n}\x1b[0m: \x1b[38;5;252;1m.\x1b[0m (\x1b[96;1m1\x1b[0m)\x1b[0m is not \x1b[97;1mMap\x1b[0m\x1b[0m",
			},
			"err: is not Array": {
				In: []interface{}{
					bwval.Holder{Val: "some"},
					bwval.PathFrom("($idx)"),
					map[string]interface{}{"idx": 1},
				},
				Panic: "Failed to get \x1b[38;5;252;1m($idx)\x1b[0m of \x1b[96;1m\"some\"\x1b[0m with \x1b[38;5;201;1mvars\x1b[0m \x1b[96;1m{\n  \"idx\": 1\n}\x1b[0m: \x1b[38;5;252;1m.\x1b[0m (\x1b[96;1m\"some\"\x1b[0m)\x1b[0m is not \x1b[97;1mArray\x1b[0m\x1b[0m",
			},
			"err: neither Array nor Map": {
				In: []interface{}{
					bwval.Holder{Val: 1},
					bwval.PathFrom("#"),
				},
				Panic: "Failed to get \x1b[38;5;252;1m#\x1b[0m of \x1b[96;1m1\x1b[0m: \x1b[38;5;252;1m.\x1b[0m (\x1b[96;1m1\x1b[0m)\x1b[0m neither \x1b[97;1mArray\x1b[0m nor \x1b[97;1mMap\x1b[0m\x1b[0m",
			},
			"err: neither Int nor String": {
				In: []interface{}{
					bwval.Holder{Val: map[string]interface{}{"some": 1}},
					bwval.PathFrom("some.($idx)"),
					map[string]interface{}{"idx": nil},
				},
				Panic: "Failed to get \x1b[38;5;252;1msome.($idx)\x1b[0m of \x1b[96;1m{\n  \"some\": 1\n}\x1b[0m with \x1b[38;5;201;1mvars\x1b[0m \x1b[96;1m{\n  \"idx\": null\n}\x1b[0m: \x1b[38;5;252;1m$idx\x1b[0m (\x1b[96;1mnull\x1b[0m)\x1b[0m neither \x1b[97;1mInt\x1b[0m nor \x1b[97;1mString\x1b[0m\x1b[0m",
			},
		},
	)
}

func TestMustMap(t *testing.T) {
	bwtesting.BwRunTests(t, bwval.MustMap,
		map[string]bwtesting.Case{
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
				Panic: "\x1b[96;1m(int)1\x1b[0m is not \x1b[97;1mMap\x1b[0m",
			},
		},
	)
}

func TestMustArray(t *testing.T) {
	bwtesting.BwRunTests(t, bwval.MustArray,
		map[string]bwtesting.Case{
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
				Panic: "\x1b[96;1m(int)1\x1b[0m is not \x1b[97;1mArray\x1b[0m",
			},
		},
	)
}

func TestMustString(t *testing.T) {
	bwtesting.BwRunTests(t, bwval.MustString,
		map[string]bwtesting.Case{
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
				Panic: "\x1b[96;1m(int)1\x1b[0m is not \x1b[97;1mString\x1b[0m",
			},
		},
	)
}

func TestMustInt(t *testing.T) {
	bwtesting.BwRunTests(t, bwval.MustInt,
		map[string]bwtesting.Case{
			"273": {
				In: []interface{}{
					273,
				},
				Out: []interface{}{
					273,
				},
			},
			"Number(-273)": {
				In: []interface{}{
					bwtype.MustNumberFrom(-273),
				},
				Out: []interface{}{
					-273,
				},
			},
			"non Int": {
				In: []interface{}{
					"some",
				},
				Panic: "\x1b[96;1m(string)some\x1b[0m is not \x1b[97;1mInt\x1b[0m",
			},
		},
	)
}

func TestMustFloat64(t *testing.T) {
	bwtesting.BwRunTests(t, bwval.MustFloat64,
		map[string]bwtesting.Case{
			"Float64(Int)": {
				In: []interface{}{
					int(273),
				},
				Out: []interface{}{
					float64(273),
				},
			},
			"Float64(Float64)": {
				In: []interface{}{
					float64(273),
				},
				Out: []interface{}{
					float64(273),
				},
			},
			"Float64(Number)": {
				In: []interface{}{
					bwtype.MustNumberFrom(273),
				},
				Out: []interface{}{
					float64(273),
				},
			},
			"non Float64": {
				In: []interface{}{
					"some",
				},
				Panic: "\x1b[96;1m(string)some\x1b[0m is not \x1b[97;1mFloat64\x1b[0m",
			},
		},
	)
}

func TestMustNumber(t *testing.T) {
	bwtesting.BwRunTests(t, bwval.MustNumber,
		map[string]bwtesting.Case{
			"Number(int)": {
				In: []interface{}{
					int(273),
				},
				Out: []interface{}{
					bwval.MustNumberFrom(273),
				},
			},
			"Number(float64)": {
				In: []interface{}{
					float64(273),
				},
				Out: []interface{}{
					bwval.MustNumberFrom(float64(273)),
				},
			},
			"Number(Number)": {
				In: []interface{}{
					bwval.NumberFromFloat64(273),
				},
				Out: []interface{}{
					bwval.NumberFromFloat64(273),
				},
			},
			"non Number": {
				In: []interface{}{
					"some",
				},
				Panic: "\x1b[96;1m(string)some\x1b[0m is not \x1b[97;1mNumber\x1b[0m",
			},
		},
	)
}

func TestMustBool(t *testing.T) {
	bwtesting.BwRunTests(t, bwval.MustBool,
		map[string]bwtesting.Case{
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
				Panic: "\x1b[96;1m(string)some\x1b[0m is not \x1b[97;1mBool\x1b[0m",
			},
		},
	)
}
