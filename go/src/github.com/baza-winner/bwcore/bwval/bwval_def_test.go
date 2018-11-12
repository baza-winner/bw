package bwval_test

import (
	"testing"

	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwval"
	"github.com/baza-winner/bwcore/bwval/deftype"
)

func TestDefMarshalJSON(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"": {
			In: []interface{}{
				bwval.Def{
					Types:      deftype.From(deftype.Int),
					IsOptional: true,
					Enum:       bwset.StringFrom("valueA", "valueB"),
					Range:      bwval.IntRange{MinPtr: bwval.PtrToInt(-1), MaxPtr: bwval.PtrToInt(1)},
					Keys: map[string]bwval.Def{
						"boolKey": {Types: deftype.From(deftype.Bool)},
					},
					Elem: &bwval.Def{
						Types: deftype.From(deftype.Bool),
					},
					ArrayElem: &bwval.Def{
						Types: deftype.From(deftype.Bool),
					},
					Default: "default value",
				},
			},
			Out: []interface{}{
				"{\n  \"ArrayElem\": {\n    \"IsOptional\": false,\n    \"Types\": [\n      \"Bool\"\n    ]\n  },\n  \"Default\": \"default value\",\n  \"Elem\": {\n    \"IsOptional\": false,\n    \"Types\": [\n      \"Bool\"\n    ]\n  },\n  \"Enum\": [\n    \"valueA\",\n    \"valueB\"\n  ],\n  \"IsOptional\": true,\n  \"Range\": \"-1..1\",\n  \"Types\": [\n    \"Int\"\n  ],\n  \"keys\": {\n    \"boolKey\": {\n      \"IsOptional\": false,\n      \"Types\": [\n        \"Bool\"\n      ]\n    }\n  }\n}",
			},
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "Bool")
	bwtesting.BwRunTests(t, DefPretty, tests)
}

func DefPretty(v bwval.Def) string {
	return bwjson.Pretty(v)
}

func TestDefFrom(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"nil": {
			// In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out:   []interface{}{bwval.Def{}},
			Panic: "\x1b[38;5;252;1m$def\x1b[0m is \x1b[91;1m(interface {})<nil>\x1b[0m",
		},
		"true": {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out:   []interface{}{bwval.Def{}},
			Panic: "\x1b[38;5;252;1m$def\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is none of \x1b[97;1mString\x1b[0m, \x1b[97;1mArray\x1b[0m or \x1b[97;1mMap\x1b[0m",
		},
		"[ Bool true ]": {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out:   []interface{}{bwval.Def{}},
			Panic: "\x1b[38;5;252;1m$def.1\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mString\x1b[0m",
		},
		"{type true}": {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out:   []interface{}{bwval.Def{}},
			Panic: "\x1b[38;5;252;1m$def.type\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is none of \x1b[97;1mString\x1b[0m or \x1b[97;1mArray\x1b[0m",
		},
		"Bool": {
			In:  []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{Types: deftype.From(deftype.Bool)}},
		},
		`{ type [ Int "bool" ] }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out:   []interface{}{bwval.Def{Types: deftype.From(deftype.Bool)}},
			Panic: "\x1b[38;5;252;1m$def.type.1\x1b[0m (\x1b[96;1m\"bool\"\x1b[0m)\x1b[0m is \x1b[91;1mnon supported\x1b[0m value\x1b[0m",
		},
		`{ type String enum <a b c> }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types: deftype.From(deftype.String),
				Enum:  bwset.StringFrom("a", "b", "c"),
			}},
			// Panic: "\x1b[38;5;252;1m$def.type.1\x1b[0m (\x1b[96;1m\"bool\"\x1b[0m)\x1b[0m is \x1b[91;1mnon supported\x1b[0m value\x1b[0m",
		},
		`{ type String enum [ "a" true ] }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types: deftype.From(deftype.String),
				Enum:  bwset.StringFrom("a", "b", "c"),
			}},
			Panic: "\x1b[38;5;252;1m$def.enum.1\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mString\x1b[0m",
		},
		`{ type Map keys { a Bool } }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types: deftype.From(deftype.Map),
				Keys:  map[string]bwval.Def{"a": {Types: deftype.From(deftype.Bool)}},
			}},
		},
		`{ type Array arrayElem Bool }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types:     deftype.From(deftype.Array),
				ArrayElem: &bwval.Def{Types: deftype.From(deftype.Bool)},
			}},
		},
		`{ type Int min 1 max 2 }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types: deftype.From(deftype.Int),
				Range: bwval.IntRange{bwval.PtrToInt(1), bwval.PtrToInt(2)},
				// ArrayElem: &bwval.Def{Types: deftype.From(deftype.Bool)},
			}},
		},
		`{ type Number min 1 max 2 }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types: deftype.From(deftype.Number),
				Range: bwval.NumberRange{bwval.PtrToNumber(1), bwval.PtrToNumber(2)},
				// ArrayElem: &bwval.Def{Types: deftype.From(deftype.Bool)},
			}},
		},
		`{ type Number min 1 max 2 default 3 }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types:      deftype.From(deftype.Number),
				Range:      bwval.NumberRange{bwval.PtrToNumber(1), bwval.PtrToNumber(2)},
				IsOptional: true,
				Default:    3,
			}},
			Panic: "\x1b[38;5;252;1m$def.default\x1b[0m (\x1b[96;1m3\x1b[0m)\x1b[0m is \x1b[91;1mout of range\x1b[0m \x1b[96;1m1..2\x1b[0m",
		},
		`{ type String default "some" }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types:      deftype.From(deftype.String),
				Default:    "some",
				IsOptional: true,
			}},
		},
		`{ type [ArrayOf Int] default [1 2 3] }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types:      deftype.From(deftype.ArrayOf, deftype.Int),
				Default:    []interface{}{1, 2, 3},
				IsOptional: true,
			}},
		},
		`{ enum <Bool Int> }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types:      deftype.From(deftype.ArrayOf, deftype.Int),
				Default:    []interface{}{1, 2, 3},
				IsOptional: true,
			}},
			Panic: "\x1b[38;5;252;1m$def\x1b[0m (\x1b[96;1m{\n  \"enum\": [\n    \"Bool\",\n    \"Int\"\n  ]\n}\x1b[0m)\x1b[0m has no key \x1b[96;1mtype\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "{ type [ Int \"bool\" ] }")
	bwtesting.BwRunTests(t, bwval.DefFrom, tests)
}
