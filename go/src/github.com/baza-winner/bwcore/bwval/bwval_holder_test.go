package bwval_test

import (
	"testing"

	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwval"
)

func TestHolderMarshalJSON(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"just val": {
			In: []interface{}{
				bwval.Holder{},
			},
			Out: []interface{}{
				"null",
			},
		},
		"val, path": {
			In: []interface{}{
				bwval.Holder{Val: "string", Pth: bwval.PathFrom("key")},
			},
			Out: []interface{}{
				"{\n  \"path\": \"key\",\n  \"val\": \"string\"\n}",
			},
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "Bool")
	bwtesting.BwRunTests(t, HolderPretty, tests)
}

func HolderPretty(v bwval.Holder) string {
	return bwjson.Pretty(v)
}

func TestHolderMustBool(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"Bool": {
			V:   bwval.Holder{Val: true, Pth: bwval.PathFrom("some.1.boolKey")},
			In:  []interface{}{},
			Out: []interface{}{true},
		},
		"non Bool": {
			V:     bwval.Holder{Val: "s", Pth: bwval.PathFrom("some.1.boolKey")},
			In:    []interface{}{},
			Panic: "\x1b[38;5;252;1msome.1.boolKey\x1b[0m (\x1b[96;1m\"s\"\x1b[0m)\x1b[0m is not \x1b[97;1mBool\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, "MustBool", tests)
}

func TestHolderMustString(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"String": {
			V:   bwval.Holder{Val: "value", Pth: bwval.PathFrom("some.1.key")},
			In:  []interface{}{},
			Out: []interface{}{"value"},
		},
		"non String": {
			V:     bwval.Holder{Val: true, Pth: bwval.PathFrom("some.1.key")},
			In:    []interface{}{},
			Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mString\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, "MustString", tests)
}

func TestHolderMustInt(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"Int": {
			V:   bwval.Holder{Val: 273, Pth: bwval.PathFrom("some.1.key")},
			In:  []interface{}{},
			Out: []interface{}{273},
		},
		"non Int": {
			V:     bwval.Holder{Val: true, Pth: bwval.PathFrom("some.1.key")},
			In:    []interface{}{},
			Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mInt\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, "MustInt", tests)
}

func TestHolderMustNumber(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"Number": {
			V:   bwval.Holder{Val: 273, Pth: bwval.PathFrom("some.1.key")},
			In:  []interface{}{},
			Out: []interface{}{float64(273)},
		},
		"non Number": {
			V:     bwval.Holder{Val: true, Pth: bwval.PathFrom("some.1.key")},
			In:    []interface{}{},
			Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mNumber\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, "MustNumber", tests)
}

func TestHolderMustArray(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"Array": {
			V:   bwval.Holder{Val: []interface{}{0, 1}, Pth: bwval.PathFrom("some.1.key")},
			In:  []interface{}{},
			Out: []interface{}{[]interface{}{0, 1}},
		},
		"non Array": {
			V:     bwval.Holder{Val: true, Pth: bwval.PathFrom("some.1.key")},
			In:    []interface{}{},
			Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mArray\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, "MustArray", tests)
}

func TestHolderMustArrayOfString(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"ArrayOfString": {
			V:   bwval.Holder{Val: []interface{}{"a", "b"}, Pth: bwval.PathFrom("some.1.key")},
			In:  []interface{}{},
			Out: []interface{}{[]string{"a", "b"}},
		},
		"non ArrayOfString": {
			V:     bwval.Holder{Val: true, Pth: bwval.PathFrom("some.1.key")},
			In:    []interface{}{},
			Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mArray\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, "MustArrayOfString", tests)
}

func TestHolderMustMap(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"Map": {
			V:   bwval.Holder{Val: map[string]interface{}{"a": 1}, Pth: bwval.PathFrom("some.1.key")},
			In:  []interface{}{},
			Out: []interface{}{map[string]interface{}{"a": 1}},
		},
		"non Map": {
			V:     bwval.Holder{Val: true, Pth: bwval.PathFrom("some.1.key")},
			In:    []interface{}{},
			Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mMap\x1b[0m",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, "MustMap", tests)
}

func TestHolderValidVal(t *testing.T) {
	tests := map[string]bwtesting.Case{}
	for k, v := range map[string]interface{}{
		"{ val: true, def: Bool }":                                            true,
		"{ val: nil def: { type Int default 273} }":                           273,
		"{ val: nil def: { type Int isOptional true} }":                       nil,
		"{ val: nil def: { type Map keys { some {type Int default 273} } } }": map[string]interface{}{"some": 273},
	} {
		h := bwval.HolderFrom(k)
		tests[k] = bwtesting.Case{
			V:   bwval.Holder{Val: h.MustKeyVal("val")},
			In:  []interface{}{bwval.DefFrom(h.MustKeyVal("def"))},
			Out: []interface{}{v},
		}
	}
	for k, v := range map[string]string{
		"{ val: 0, def: Bool }": "\x1b[96;1m0\x1b[0m::\x1b[38;5;252;1m.\x1b[0m (\x1b[96;1m0\x1b[0m)\x1b[0m is not \x1b[97;1mBool\x1b[0m\x1b[0m",

		"{ val: [true 0], def: [ArrayOf Bool] }": "\x1b[96;1m[\n  true,\n  0\n]\x1b[0m::\x1b[38;5;252;1m1\x1b[0m (\x1b[96;1m0\x1b[0m)\x1b[0m is not \x1b[97;1mBool\x1b[0m\x1b[0m",
		"{ val: nil def: Array}":                 "\x1b[96;1mnull\x1b[0m::\x1b[38;5;252;1m.\x1b[0m (\x1b[96;1mnull\x1b[0m)\x1b[0m is not \x1b[97;1mArray\x1b[0m\x1b[0m",
		"{ val: { key: <some thing> } def: { type Map elem { type <ArrayOf String> enum <some good>}} }": "\x1b[96;1m{\n  \"key\": [\n    \"some\",\n    \"thing\"\n  ]\n}\x1b[0m::\x1b[38;5;252;1mkey.1\x1b[0m: expected one of \x1b[96;1m[\n  \"good\",\n  \"some\"\n]\x1b[0m instead of \x1b[91;1m\"thing\"\x1b[0m\x1b[0m",
	} {
		h := bwval.HolderFrom(k)
		tests[k] = bwtesting.Case{
			V:     bwval.Holder{Val: h.MustKeyVal("val")},
			In:    []interface{}{bwval.DefFrom(h.MustKeyVal("def"))},
			Panic: v,
		}
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "{ val: { key: <some thing> } def: { type Map elem { type <ArrayOf String> enum <some good>}} }")
	// for _, test := range tests {
	// 	bwdebug.Print("test", bwjson.Pretty(test))
	// }
	bwtesting.BwRunTests(t, "MustValidVal", tests)
}
