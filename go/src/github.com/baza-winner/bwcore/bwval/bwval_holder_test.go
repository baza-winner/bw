package bwval_test

import (
	"fmt"
	"testing"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwparse"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwtype"
	"github.com/baza-winner/bwcore/bwval"
)

func TestHolderMustPath(t *testing.T) {
	bwtesting.BwRunTests(t, "MustPath",
		map[string]bwtesting.Case{
			"ok": {
				V:   bwval.Holder{Val: []interface{}{nil, map[string]interface{}{"some": "thing"}}},
				In:  []interface{}{bwval.PathStr{S: "1.some"}.MustPath()},
				Out: []interface{}{bwval.Holder{Val: "thing", Pth: bwval.PathStr{S: "1.some"}.MustPath()}},
			},
			"panic": {
				V:     bwval.Holder{Val: []interface{}{map[string]interface{}{"some": "thing"}}},
				In:    []interface{}{bwval.PathStr{S: "1.some"}.MustPath()},
				Panic: "\x1b[38;5;252;1m.\x1b[0m (\x1b[96;1m[\n  {\n    \"some\": \"thing\"\n  }\n]\x1b[0m)\x1b[0m has not enough length (\x1b[96;1m1\x1b[0m) for idx (\x1b[96;1m1)\x1b[0m",
			},
		},
	)
}

func TestHolderMarshalJSON(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(v bwval.Holder) string {
			return bwjson.Pretty(v)
		},
		map[string]bwtesting.Case{
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
					bwval.Holder{Val: "string", Pth: bwval.PathStr{S: "key"}.MustPath()},
				},
				Out: []interface{}{
					"{\n  \"path\": \"key\",\n  \"val\": \"string\"\n}",
				},
			},
		},
	)
}

func TestHolderMustBool(t *testing.T) {
	bwtesting.BwRunTests(t,
		"MustBool",
		map[string]bwtesting.Case{
			"true": {
				V:   bwval.Holder{Val: true},
				Out: []interface{}{true},
			},
			"false": {
				V:   bwval.Holder{Val: false},
				Out: []interface{}{false},
			},
			"non Bool": {
				V:     bwval.Holder{Val: "s", Pth: bwval.PathStr{S: "some.1.boolKey"}.MustPath()},
				Panic: "\x1b[38;5;252;1msome.1.boolKey\x1b[0m (\x1b[96;1m\"s\"\x1b[0m)\x1b[0m is not \x1b[97;1mBool\x1b[0m",
			},
		},
	)
}

func TestHolderMustString(t *testing.T) {
	bwtesting.BwRunTests(t, "MustString",
		map[string]bwtesting.Case{
			"String": {
				V:   bwval.Holder{Val: "value"},
				Out: []interface{}{"value"},
			},
			"non String": {
				V:     bwval.Holder{Val: true, Pth: bwval.PathStr{S: "some.1.key"}.MustPath()},
				In:    []interface{}{},
				Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mString\x1b[0m",
			},
		},
	)
}

func TestHolderMustInt(t *testing.T) {
	bwtesting.BwRunTests(t,
		"MustInt",
		map[string]bwtesting.Case{
			"-273": {
				V:   bwval.Holder{Val: -273},
				Out: []interface{}{-273},
			},
			"273": {
				V:   bwval.Holder{Val: 273},
				Out: []interface{}{273},
			},
			"float64(273)": {
				V:   bwval.Holder{Val: float64(273)},
				Out: []interface{}{273},
			},
			"bwtype.MustNumberFrom(float64(273))": {
				V:   bwval.Holder{Val: bwtype.MustNumberFrom(float64(273))},
				Out: []interface{}{273},
			},
			"non Int: bwtype.MustNumberFrom(bw.MaxUint)": {
				V:     bwval.Holder{Val: bwtype.MustNumberFrom(bw.MaxUint), Pth: bwval.PathStr{S: "some.1.key"}.MustPath()},
				Panic: fmt.Sprintf("\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1m%d\x1b[0m)\x1b[0m is not \x1b[97;1mInt\x1b[0m", bw.MaxUint),
			},
			"non Int: true": {
				V:     bwval.Holder{Val: true, Pth: bwval.PathStr{S: "some.1.key"}.MustPath()},
				Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mInt\x1b[0m",
			},
		},
	)
}

func TestHolderMustUint(t *testing.T) {
	bwtesting.BwRunTests(t,
		"MustUint",
		map[string]bwtesting.Case{
			"273": {
				V:   bwval.Holder{Val: 273},
				Out: []interface{}{uint(273)},
			},
			"float64(273)": {
				V:   bwval.Holder{Val: float64(273)},
				Out: []interface{}{uint(273)},
			},
			"bwtype.MustNumberFrom(float64(273))": {
				V:   bwval.Holder{Val: bwtype.MustNumberFrom(float64(273))},
				Out: []interface{}{uint(273)},
			},
			"non Uint: bwtype.MustNumberFrom(float64(-273))": {
				V:     bwval.Holder{Val: bwtype.MustNumberFrom(float64(-273)), Pth: bwval.PathStr{S: "some.1.key"}.MustPath()},
				Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1m-273\x1b[0m)\x1b[0m is not \x1b[97;1mUint\x1b[0m",
			},
			"non Uint: true": {
				V:     bwval.Holder{Val: true, Pth: bwval.PathStr{S: "some.1.key"}.MustPath()},
				Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mUint\x1b[0m",
			},
		},
	)
}

func TestHolderMustFloat64(t *testing.T) {
	bwtesting.BwRunTests(t, "MustFloat64",
		map[string]bwtesting.Case{
			"273": {
				V:   bwval.Holder{Val: 273},
				Out: []interface{}{float64(273)},
			},
			"uint(273)": {
				V:   bwval.Holder{Val: uint(273)},
				Out: []interface{}{float64(273)},
			},
			"float32(273)": {
				V:   bwval.Holder{Val: float32(273)},
				Out: []interface{}{float64(273)},
			},
			"float64(273)": {
				V:   bwval.Holder{Val: float64(273)},
				Out: []interface{}{float64(273)},
			},
			"non Float64": {
				V:     bwval.Holder{Val: true, Pth: bwval.PathStr{S: "some.1.key"}.MustPath()},
				Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mFloat64\x1b[0m",
			},
		},
	)
}

func TestHolderMustArray(t *testing.T) {
	bwtesting.BwRunTests(t, "MustArray",
		map[string]bwtesting.Case{
			"[0 1]": {
				V:   bwval.Holder{Val: []interface{}{0, 1}},
				Out: []interface{}{[]interface{}{0, 1}},
			},
			"<some thing>": {
				V:   bwval.Holder{Val: []string{"some", "thing"}},
				Out: []interface{}{[]interface{}{"some", "thing"}},
			},
			"non Array": {
				V:     bwval.Holder{Val: true, Pth: bwval.PathStr{S: "some.1.key"}.MustPath()},
				Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mArray\x1b[0m",
			},
		},
	)
}

func TestHolderMustMap(t *testing.T) {
	bwtesting.BwRunTests(t, "MustMap",
		map[string]bwtesting.Case{
			"map[string]interface{}": {
				V:   bwval.Holder{Val: map[string]interface{}{"a": 1}},
				Out: []interface{}{map[string]interface{}{"a": 1}},
			},
			"map[string]string": {
				V:   bwval.Holder{Val: map[string]string{"a": "some"}},
				Out: []interface{}{map[string]interface{}{"a": "some"}},
			},
			"non Map": {
				V:     bwval.Holder{Val: true, Pth: bwval.PathStr{S: "some.1.key"}.MustPath()},
				Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mMap\x1b[0m",
			},
		},
	)
}

func TestHolderValidVal(t *testing.T) {
	bwtesting.BwRunTests(t,
		"MustValidVal",
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			testProto := func(source string) bwtesting.Case {
				test := bwtesting.Case{}
				bwparse.Map(bwparse.From(bwrune.FromString(source)), bwparse.Opt{
					OnValidateMapKey: func(on bwparse.On, m map[string]interface{}, key string) (err error) {
						if !bwset.StringFrom("val", "def").Has(key) {
							err = on.P.Error(bwparse.A{
								Start: on.Start,
								Fmt:   bw.Fmt(ansi.String("unexpected key `<ansiErr>%s<ansi>`"), on.Start.Suffix()),
							})
						}
						return
					},
					OnParseMapElem: func(on bwparse.On, m map[string]interface{}, key string) (status bwparse.Status) {
						switch on.Opt.Path.String() {
						case "val":
							var val interface{}
							if val, status = bwparse.Val(on.P); status.IsOK() {
								test.V = bwval.Holder{Val: val}
							}
						case "def":
							var def bwval.Def
							if def, status = bwval.ParseDef(on.P); status.IsOK() {
								test.In = []interface{}{def}
							}
						}
						m[key] = nil
						return
					},
				})
				return test
			}
			for k, v := range map[string]interface{}{
				"{ val: true, def: Bool }":                                              true,
				"{ val: nil def: { type Int default 273} }":                             273,
				"{ val: nil def: { type Int isOptional true} }":                         nil,
				"{ val: nil def: { type Map keys { some {type Int default 273} } } }":   map[string]interface{}{"some": 273},
				"{ val: <some thing> def Array }":                                       []interface{}{"some", "thing"},
				"{ val: <some thing> def: { type [ArrayOf String] enum <some thing>} }": []interface{}{"some", "thing"},
				`{
          val {
            some 273
            thing 3.14
          }
          def {
            type Map
            elem Number
          }
        }`: map[string]interface{}{"some": 273, "thing": 3.14},
			} {
				test := testProto(k)
				test.Out = []interface{}{v}
				tests[k] = test
			}

			for k, v := range map[string]string{
				"{ val: 0, def: Bool }": "\x1b[96;1m0\x1b[0m::\x1b[38;5;252;1m.\x1b[0m (\x1b[96;1m0\x1b[0m)\x1b[0m is not \x1b[97;1mBool\x1b[0m\x1b[0m",

				"{ val: [true 0], def: [ArrayOf Bool] }": "\x1b[96;1m[\n  true,\n  0\n]\x1b[0m::\x1b[38;5;252;1m1\x1b[0m (\x1b[96;1m0\x1b[0m)\x1b[0m is not \x1b[97;1mBool\x1b[0m\x1b[0m",
				"{ val: nil def: Array}":                 "\x1b[96;1mnull\x1b[0m::\x1b[38;5;252;1m.\x1b[0m (\x1b[96;1mnull\x1b[0m)\x1b[0m is not \x1b[97;1mArray\x1b[0m\x1b[0m",
				"{ val: { key: <some thing> } def: { type Map elem { type <ArrayOf String> enum <some good>}} }": "\x1b[96;1m{\n  \"key\": [\n    \"some\",\n    \"thing\"\n  ]\n}\x1b[0m::\x1b[38;5;252;1mkey.1\x1b[0m: expected one of \x1b[96;1m[\n  \"good\",\n  \"some\"\n]\x1b[0m instead of \x1b[91;1m\"thing\"\x1b[0m\x1b[0m",
				"{ val: { some: 0 thing: 1 } def: { type Map keys { some Int } } }":                              "\x1b[96;1m{\n  \"some\": 0,\n  \"thing\": 1\n}\x1b[0m::\x1b[38;5;252;1m.\x1b[0m (\x1b[96;1m{\n  \"some\": 0,\n  \"thing\": 1\n}\x1b[0m)\x1b[0m has unexpected key \x1b[96;1m\"thing\"\x1b[0m\x1b[0m",
				"{ val: { some: 0 thing: 1 } def: { type Map keys { some Int } elem Bool } }":                    "\x1b[96;1m{\n  \"some\": 0,\n  \"thing\": 1\n}\x1b[0m::\x1b[38;5;252;1mthing\x1b[0m (\x1b[96;1m1\x1b[0m)\x1b[0m is not \x1b[97;1mBool\x1b[0m\x1b[0m",
				"{ val: { some: 0 } def: { type Map keys { some Bool } } }":                                      "\x1b[96;1m{\n  \"some\": 0\n}\x1b[0m::\x1b[38;5;252;1msome\x1b[0m (\x1b[96;1m0\x1b[0m)\x1b[0m is not \x1b[97;1mBool\x1b[0m\x1b[0m",
				"{ val: [0, true] def: { type Array elem Int } }":                                                "\x1b[96;1m[\n  0,\n  true\n]\x1b[0m::\x1b[38;5;252;1m1\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mInt\x1b[0m\x1b[0m",
			} {
				test := testProto(k)
				test.Panic = v
				tests[k] = test
			}
			return tests
		}(),
	)
}
