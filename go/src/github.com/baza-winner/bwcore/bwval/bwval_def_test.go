package bwval_test

import (
	"testing"

	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwval"
)

func TestDefMarshalJSON(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"": {
			In: []interface{}{
				bwval.Def{
					Types:      bwval.ValKindSetFrom(bwval.ValInt),
					IsOptional: true,
					Enum:       bwset.StringFrom("valueA", "valueB"),
					Range:      bwval.IntRange{MinPtr: bwval.PtrToInt(-1), MaxPtr: bwval.PtrToInt(1)},
					Keys: map[string]bwval.Def{
						"boolKey": {Types: bwval.ValKindSetFrom(bwval.ValBool)},
					},
					Elem: &bwval.Def{
						Types: bwval.ValKindSetFrom(bwval.ValBool),
					},
					ArrayElem: &bwval.Def{
						Types: bwval.ValKindSetFrom(bwval.ValBool),
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
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def\x1b[0m is \x1b[91;1m(interface {})<nil>\x1b[0m",
		},
		"true": {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m neither \x1b[97;1mString\x1b[0m nor \x1b[97;1mArray\x1b[0m nor \x1b[97;1mMap\x1b[0m",
		},
		"[ Bool true ]": {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.1\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mString\x1b[0m",
		},
		"{type true}": {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.type\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m neither \x1b[97;1mString\x1b[0m nor \x1b[97;1mArray\x1b[0m",
		},
		"Bool": {
			In:  []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{Types: bwval.ValKindSetFrom(bwval.ValBool)}},
		},
		`{ type [ Int "bool" ] }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.type.1\x1b[0m (\x1b[96;1m\"bool\"\x1b[0m)\x1b[0m is \x1b[91;1mnon supported\x1b[0m value\x1b[0m",
		},
		`{ type String enum <a b c> }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types: bwval.ValKindSetFrom(bwval.ValString),
				Enum:  bwset.StringFrom("a", "b", "c"),
			}},
			// Panic: "\x1b[38;5;252;1m$def.type.1\x1b[0m (\x1b[96;1m\"bool\"\x1b[0m)\x1b[0m is \x1b[91;1mnon supported\x1b[0m value\x1b[0m",
		},
		`{ type String enum [ "a" true ] }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.enum.1\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mString\x1b[0m",
		},
		`{ type Map keys { a Bool } }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types: bwval.ValKindSetFrom(bwval.ValMap),
				Keys:  map[string]bwval.Def{"a": {Types: bwval.ValKindSetFrom(bwval.ValBool)}},
			}},
		},
		`{ type Array arrayElem Bool }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types:     bwval.ValKindSetFrom(bwval.ValArray),
				ArrayElem: &bwval.Def{Types: bwval.ValKindSetFrom(bwval.ValBool)},
			}},
		},
		`{ type Int min 1 max 2 }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types: bwval.ValKindSetFrom(bwval.ValInt),
				Range: bwval.IntRange{bwval.PtrToInt(1), bwval.PtrToInt(2)},
				// ArrayElem: &bwval.Def{Types: bwval.ValKindSetFrom(bwval.ValBool)},
			}},
		},
		`{ type Number min 1 max 2 }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types: bwval.ValKindSetFrom(bwval.ValNumber),
				Range: bwval.NumberRange{bwval.PtrToNumber(1), bwval.PtrToNumber(2)},
				// ArrayElem: &bwval.Def{Types: bwval.ValKindSetFrom(bwval.ValBool)},
			}},
		},
		`{ type Number min 1 max 2 default 3 }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.default\x1b[0m (\x1b[96;1m3\x1b[0m)\x1b[0m is \x1b[91;1mout of range\x1b[0m \x1b[96;1m1..2\x1b[0m",
		},
		`{ type String default "some" }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types:      bwval.ValKindSetFrom(bwval.ValString),
				Default:    "some",
				IsOptional: true,
			}},
		},
		`{ type [ArrayOf Int] default [1 2 3] }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types:      bwval.ValKindSetFrom(bwval.ValArrayOf, bwval.ValInt),
				Default:    []interface{}{1, 2, 3},
				IsOptional: true,
			}},
		},
		`{ enum <Bool Int> }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def\x1b[0m (\x1b[96;1m{\n  \"enum\": [\n    \"Bool\",\n    \"Int\"\n  ]\n}\x1b[0m)\x1b[0m has no key \x1b[96;1mtype\x1b[0m",
		},
		`{ type Map keys true }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.keys\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mMap\x1b[0m",
		},
		`{ type Map keys { some true } }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.keys.some\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m neither \x1b[97;1mString\x1b[0m nor \x1b[97;1mArray\x1b[0m nor \x1b[97;1mMap\x1b[0m",
		},
		`{ type Array arrayElem true }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.arrayElem\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m neither \x1b[97;1mString\x1b[0m nor \x1b[97;1mArray\x1b[0m nor \x1b[97;1mMap\x1b[0m",
		},
		`{ type Array elem true }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.elem\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m neither \x1b[97;1mString\x1b[0m nor \x1b[97;1mArray\x1b[0m nor \x1b[97;1mMap\x1b[0m",
		},
		`{ type Int min 2.1 }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.min\x1b[0m (\x1b[96;1m2.1\x1b[0m)\x1b[0m is not \x1b[97;1mInt\x1b[0m",
		},
		`{ type Int min 2.0 }`: {
			In: []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Out: []interface{}{bwval.Def{
				Types:      bwval.ValKindSetFrom(bwval.ValInt),
				IsOptional: false,
				Range:      bwval.IntRange{MinPtr: bwval.PtrToInt(2)},
			}},
		},
		`{ type Int max 2.1 }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.max\x1b[0m (\x1b[96;1m2.1\x1b[0m)\x1b[0m is not \x1b[97;1mInt\x1b[0m",
		},
		`{ type Int min 3 max 2 }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def\x1b[0m (\x1b[96;1m{\n  \"max\": 2,\n  \"min\": 3,\n  \"type\": \"Int\"\n}\x1b[0m)\x1b[0m: \x1b[38;5;252;1m.max\x1b[0m (\x1b[96;1m2\x1b[0m) must not be \x1b[91;1mless\x1b[0m then \x1b[38;5;252;1m.min\x1b[0m (\x1b[96;1m3\x1b[0m)\x1b[0m",
		},
		`{ type Number min true }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.min\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mNumber\x1b[0m",
		},
		`{ type Number max true }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.max\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mNumber\x1b[0m",
		},
		`{ type Number min 3.14 max 2.71 }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def\x1b[0m (\x1b[96;1m{\n  \"max\": 2.71,\n  \"min\": 3.14,\n  \"type\": \"Number\"\n}\x1b[0m)\x1b[0m: \x1b[38;5;252;1m.max\x1b[0m (\x1b[96;1m2.71\x1b[0m) must not be \x1b[91;1mless\x1b[0m then \x1b[38;5;252;1m.min\x1b[0m (\x1b[96;1m3.14\x1b[0m)\x1b[0m",
		},
		`{ type Number isOptional 3 }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.isOptional\x1b[0m (\x1b[96;1m3\x1b[0m)\x1b[0m is not \x1b[97;1mBool\x1b[0m",
		},
		`{ type Number isOptional false default 3.14 }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def\x1b[0m (\x1b[96;1m{\n  \"default\": 3.14,\n  \"isOptional\": false,\n  \"type\": \"Number\"\n}\x1b[0m)\x1b[0m: having \x1b[38;5;252;1m.default\x1b[0m can not have \x1b[38;5;252;1m.isOptional\x1b[0m \x1b[96;1mtrue\x1b[0m",
		},
		`{ type Number keyA "valA" keyB "valB" }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def\x1b[0m (\x1b[96;1m{\n  \"keyA\": \"valA\",\n  \"keyB\": \"valB\",\n  \"type\": \"Number\"\n}\x1b[0m)\x1b[0m has unexpected keys: \x1b[96;1m[\n  \"keyA\",\n  \"keyB\"\n]\x1b[0m",
		},
		`{ type Number keyA "valA" }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def\x1b[0m (\x1b[96;1m{\n  \"keyA\": \"valA\",\n  \"type\": \"Number\"\n}\x1b[0m)\x1b[0m has unexpected key \x1b[96;1m\"keyA\"\x1b[0m",
		},
		`{ type <ArrayOf> }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.type\x1b[0m (\x1b[96;1m[\n  \"ArrayOf\"\n]\x1b[0m)\x1b[0m: \x1b[96;1mArrayOf\x1b[0m must be followed by some type, can not be \x1b[91;1mused alone\x1b[0m",
		},
		`{ type <ArrayOf Array> }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.type\x1b[0m (\x1b[96;1m[\n  \"ArrayOf\",\n  \"Array\"\n]\x1b[0m)\x1b[0m: values \x1b[96;1m\"ArrayOf\"\x1b[0m and \x1b[96;1m\"Array\"\x1b[0m are \x1b[91;1mmutually exclusive\x1b[0m, can not be \x1b[91;1mused both at once\x1b[0m",
		},
		`{ type <Int Number> }`: {
			In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
			Panic: "\x1b[38;5;252;1m$def.type\x1b[0m (\x1b[96;1m[\n  \"Int\",\n  \"Number\"\n]\x1b[0m)\x1b[0m: values \x1b[96;1m\"Int\"\x1b[0m and \x1b[96;1m\"Number\"\x1b[0m are \x1b[91;1mmutually exclusive\x1b[0m, can not be \x1b[91;1mused both at once\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "{ type Int min 2.0 }")
	bwtesting.BwRunTests(t, bwval.DefFrom, tests)
}
