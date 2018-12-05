package bwval_test

import (
	"fmt"
	"testing"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwparse"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwtype"
	"github.com/baza-winner/bwcore/bwval"
)

func TestHolderMustKindSwitch(t *testing.T) {
	type A struct {
		kindCases   map[bwtype.ValKind]bwval.KindCase
		defaultCase bwval.KindCase
	}
	bwtesting.BwRunTests(t,
		func(v bwval.Holder, a A) interface{} {
			return v.MustKindSwitch(a.kindCases, a.defaultCase)
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{
				"default case": {
					In: []interface{}{
						bwval.Holder{Val: nil},
						A{
							map[bwtype.ValKind]bwval.KindCase{
								bwtype.ValInt: func(val interface{}, kind bwtype.ValKind) (interface{}, error) { return nil, nil },
							},
							func(val interface{}, kind bwtype.ValKind) (interface{}, error) { return "default", nil },
						},
					},
					Out: []interface{}{"default"},
				},
				"default case with error": {
					In: []interface{}{
						bwval.Holder{Val: nil},
						A{
							map[bwtype.ValKind]bwval.KindCase{
								bwtype.ValInt: func(val interface{}, kind bwtype.ValKind) (interface{}, error) { return nil, nil },
							},
							func(val interface{}, kind bwtype.ValKind) (interface{}, error) { return nil, bwerr.From("no default") },
						},
					},
					Panic: "no default",
				},
			}
			return tests
		}(),
	)
}

func TestHolderMustSetKeyVal(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(target bwval.Holder, key string, val interface{}) bwval.Holder {
			target.MustSetKeyVal(key, val)
			return target
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{
				"ok": {
					In: []interface{}{
						bwval.Holder{Val: map[string]interface{}{"some": "not bad"}, Pth: bwval.MustPath(bwval.PathStr{S: "1.some"})},
						"thing", "good",
					},
					Out: []interface{}{
						bwval.Holder{Val: map[string]interface{}{"some": "not bad", "thing": "good"}, Pth: bwval.MustPath(bwval.PathStr{S: "1.some"})},
					},
				},
				"panic": {
					In: []interface{}{
						bwval.Holder{Val: nil, Pth: bwval.MustPath(bwval.PathStr{S: "1.some"})},
						"thing", "good",
					},
					Panic: "\x1b[38;5;252;1m1.some\x1b[0m (\x1b[96;1mnull\x1b[0m)\x1b[0m is not \x1b[97;1mMap\x1b[0m",
				},
			}
			return tests
		}(),
	)
}

func TestHolderMustKeyVal(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(v bwval.Holder, key string, opt ...bool) interface{} {
			if len(opt) == 0 {
				return v.MustKeyVal(key)
			} else {
				return v.MustKeyVal(key, nil)
				// } else {
				// bwdebug.Print("opt:#v", opt, "opt[0] == nil", opt[0] == interface{}(nil), "opt[0]:#v", opt[0])
			}
			return nil
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{
				"ok": {
					In: []interface{}{
						bwval.Holder{Val: map[string]interface{}{"thing": "good"}},
						"thing",
					},
					Out: []interface{}{"good"},
				},
				"default nil": {
					In: []interface{}{
						bwval.Holder{Val: nil},
						"thing",
						true,
					},
					Out: []interface{}{nil},
				},
				"panic": {
					In: []interface{}{
						bwval.Holder{Val: nil, Pth: bwval.MustPath(bwval.PathStr{S: "1.some"})},
						"thing",
					},
					Panic: "\x1b[38;5;252;1m1.some\x1b[0m is \x1b[91;1m(interface {})<nil>\x1b[0m",
				},
			}
			return tests
		}(),
		// "default nil",
	)
}

func TestHolderMustKey(t *testing.T) {
	bwtesting.BwRunTests(t,
		"MustKey",
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{
				"panic": {
					V:     bwval.Holder{Val: nil, Pth: bwval.MustPath(bwval.PathStr{S: "1.some"})},
					In:    []interface{}{"thing"},
					Panic: "\x1b[38;5;252;1m1.some\x1b[0m is \x1b[91;1m(interface {})<nil>\x1b[0m",
				},
				"ok": {
					V:   bwval.Holder{Val: map[string]interface{}{"thing": "good"}, Pth: bwval.MustPath(bwval.PathStr{S: "1.some"})},
					In:  []interface{}{"thing"},
					Out: []interface{}{bwval.Holder{Val: "good", Pth: bwval.MustPath(bwval.PathStr{S: "1.some.thing"})}},
				},
			}
			return tests
		}(),
	)
}

func TestHolderMustIdxVal(t *testing.T) {
	bwtesting.BwRunTests(t,
		"MustIdxVal",
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{
				"ok": {
					V:   bwval.Holder{Val: []interface{}{"thing", "good"}},
					In:  []interface{}{1},
					Out: []interface{}{"good"},
				},
				"panic": {
					V:     bwval.Holder{Val: nil, Pth: bwval.MustPath(bwval.PathStr{S: "1.some"})},
					In:    []interface{}{1},
					Panic: "\x1b[38;5;252;1m1.some\x1b[0m is \x1b[91;1m(interface {})<nil>\x1b[0m",
				},
			}
			return tests
		}(),
	)
}

func TestHolderMustIdx(t *testing.T) {
	bwtesting.BwRunTests(t,
		"MustIdx",
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{
				"ok": {
					V:   bwval.Holder{Val: []interface{}{"thing", "good"}, Pth: bwval.MustPath(bwval.PathStr{S: "1.some"})},
					In:  []interface{}{1},
					Out: []interface{}{bwval.Holder{Val: "good", Pth: bwval.MustPath(bwval.PathStr{S: "1.some.1"})}},
				},
				"panic": {
					V:     bwval.Holder{Val: nil, Pth: bwval.MustPath(bwval.PathStr{S: "1.some"})},
					In:    []interface{}{1},
					Panic: "\x1b[38;5;252;1m1.some\x1b[0m is \x1b[91;1m(interface {})<nil>\x1b[0m",
				},
			}
			return tests
		}(),
	)
}

func TestHolderMustPath(t *testing.T) {
	bwtesting.BwRunTests(t, "MustPath",
		map[string]bwtesting.Case{
			"ok": {
				V:   bwval.Holder{Val: []interface{}{nil, map[string]interface{}{"some": "thing"}}},
				In:  []interface{}{bwval.MustPath(bwval.PathStr{S: "1.some"})},
				Out: []interface{}{bwval.Holder{Val: "thing", Pth: bwval.MustPath(bwval.PathStr{S: "1.some"})}},
			},
			"panic": {
				V:     bwval.Holder{Val: []interface{}{map[string]interface{}{"some": "thing"}}},
				In:    []interface{}{bwval.MustPath(bwval.PathStr{S: "1.some"})},
				Panic: "\x1b[38;5;252;1m.\x1b[0m (\x1b[96;1m[\n  {\n    \"some\": \"thing\"\n  }\n]\x1b[0m)\x1b[0m has not enough length (\x1b[96;1m1\x1b[0m) for idx (\x1b[96;1m1)\x1b[0m",
			},
			"invalid path": {
				V:     bwval.Holder{Val: []interface{}{nil, map[string]interface{}{"some": "thing"}}},
				In:    []interface{}{bwval.PathStr{S: ".some"}},
				Panic: "invalid path: unexpected char \x1b[96;1m's'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m115\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m1\x1b[0m: \x1b[32m.\x1b[91ms\x1b[0mome\n",
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
					bwval.Holder{Val: "string", Pth: bwval.MustPath(bwval.PathStr{S: "key"})},
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
			"non Bool, but true by default": {
				V:   bwval.Holder{Val: "s"},
				In:  []interface{}{true},
				Out: []interface{}{true},
			},
			"non Bool, but false by default": {
				V:   bwval.Holder{Val: "s"},
				In:  []interface{}{false},
				Out: []interface{}{false},
			},
			"non Bool, no default": {
				V:     bwval.Holder{Val: "s", Pth: bwval.MustPath(bwval.PathStr{S: "some.1.boolKey"})},
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
			"non String, but has default": {
				V:   bwval.Holder{Val: true},
				In:  []interface{}{"good"},
				Out: []interface{}{"good"},
			},
			"non String": {
				V:     bwval.Holder{Val: true, Pth: bwval.MustPath(bwval.PathStr{S: "some.1.key"})},
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
				V:     bwval.Holder{Val: bwtype.MustNumberFrom(bw.MaxUint), Pth: bwval.MustPath(bwval.PathStr{S: "some.1.key"})},
				Panic: fmt.Sprintf("\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1m%d\x1b[0m)\x1b[0m is not \x1b[97;1mInt\x1b[0m", bw.MaxUint),
			},
			"non Int: bwtype.MustNumberFrom(bw.MaxUint), but has default": {
				V:   bwval.Holder{Val: bwtype.MustNumberFrom(bw.MaxUint)},
				In:  []interface{}{273},
				Out: []interface{}{273},
			},
			"non Int: true": {
				V:     bwval.Holder{Val: true, Pth: bwval.MustPath(bwval.PathStr{S: "some.1.key"})},
				Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mInt\x1b[0m",
			},
			"non Int: true, but has default": {
				V:   bwval.Holder{Val: true},
				In:  []interface{}{273},
				Out: []interface{}{273},
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
				V:     bwval.Holder{Val: bwtype.MustNumberFrom(float64(-273)), Pth: bwval.MustPath(bwval.PathStr{S: "some.1.key"})},
				Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1m-273\x1b[0m)\x1b[0m is not \x1b[97;1mUint\x1b[0m",
			},
			"non Uint: true": {
				V:     bwval.Holder{Val: true, Pth: bwval.MustPath(bwval.PathStr{S: "some.1.key"})},
				Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mUint\x1b[0m",
			},
			"non Uint: bwtype.MustNumberFrom(float64(-273)), but has default": {
				V:   bwval.Holder{Val: bwtype.MustNumberFrom(float64(-273))},
				In:  []interface{}{uint(273)},
				Out: []interface{}{uint(273)},
			},
			"non Uint: true, but has default": {
				V:   bwval.Holder{Val: true},
				In:  []interface{}{uint(273)},
				Out: []interface{}{uint(273)},
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
				V:     bwval.Holder{Val: true, Pth: bwval.MustPath(bwval.PathStr{S: "some.1.key"})},
				Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mFloat64\x1b[0m",
			},
			"non Float64, but has default": {
				V:   bwval.Holder{Val: true},
				In:  []interface{}{float64(273)},
				Out: []interface{}{float64(273)},
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
				V:     bwval.Holder{Val: true, Pth: bwval.MustPath(bwval.PathStr{S: "some.1.key"})},
				Panic: "\x1b[38;5;252;1msome.1.key\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mArray\x1b[0m",
			},
			"non Array, but has default": {
				V:   bwval.Holder{Val: true},
				In:  []interface{}{[]interface{}{}},
				Out: []interface{}{[]interface{}{}},
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
				V:     bwval.Holder{Val: true, Pth: bwval.MustPath(bwval.PathStr{S: "some.1.key"})},
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
				p := bwparse.From(bwrune.S{source})
				if _, err := bwparse.SkipSpace(p, bwparse.TillNonEOF); err != nil {
					bwerr.PanicErr(err)
				}
				bwparse.Map(p, bwparse.Opt{
					OnValidateMapKey: func(on bwparse.On, m map[string]interface{}, key string) (err error) {
						if !bwset.StringFrom("val", "def").Has(key) {
							err = on.P.Error(bwparse.E{
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
				`
					{
					 	val: nil
					 	def:
					 	  {
								type Map
								keys {
									v {
										type String
										enum <all err ok none>
										default "none"
									}
									s {
										type String
										enum <none stderr stdout all>
										default "all"
									}
									exitOnError {
										type Bool
										default false
									}
								}
							}
					}
				`: map[string]interface{}{"v": "none", "s": "all", "exitOnError": false},
			} {
				test := testProto(k)
				test.Out = []interface{}{v}
				tests[k] = test
			}

			for k, v := range map[string]string{
				"{ val: 0, def: Bool }":                  "\x1b[96;1m0\x1b[0m::\x1b[38;5;252;1m.\x1b[0m (\x1b[96;1m0\x1b[0m)\x1b[0m is not \x1b[97;1mBool\x1b[0m\x1b[0m",
				"{ val: [true 0], def: [ArrayOf Bool] }": "\x1b[96;1m[\n  true,\n  0\n]\x1b[0m::\x1b[38;5;252;1m1\x1b[0m (\x1b[96;1m0\x1b[0m)\x1b[0m is not \x1b[97;1mBool\x1b[0m\x1b[0m",
				"{ val: nil def: Array}":                 "\x1b[96;1mnull\x1b[0m::\x1b[38;5;252;1m.\x1b[0m (\x1b[96;1mnull\x1b[0m)\x1b[0m is not \x1b[97;1mArray\x1b[0m\x1b[0m",
				"{ val: { key: <some thing> } def: { type Map elem { type <ArrayOf String> enum <some good>}} }": "\x1b[96;1m{\n  \"key\": [\n    \"some\",\n    \"thing\"\n  ]\n}\x1b[0m::\x1b[38;5;252;1mkey.1\x1b[0m: expected one of \x1b[96;1m[\n  \"good\",\n  \"some\"\n]\x1b[0m instead of \x1b[91;1m\"thing\"\x1b[0m\x1b[0m",
				"{ val: { some: 0 thing: 1 } def: { type Map keys { some Int } } }":                              "\x1b[96;1m{\n  \"some\": 0,\n  \"thing\": 1\n}\x1b[0m::\x1b[38;5;252;1m.\x1b[0m (\x1b[96;1m{\n  \"some\": 0,\n  \"thing\": 1\n}\x1b[0m)\x1b[0m has unexpected key \x1b[96;1m\"thing\"\x1b[0m\x1b[0m",
				"{ val: { some: 0 thing: 1 good: 2 } def: { type Map keys { some Int } } }":                      "\x1b[96;1m{\n  \"good\": 2,\n  \"some\": 0,\n  \"thing\": 1\n}\x1b[0m::\x1b[38;5;252;1m.\x1b[0m (\x1b[96;1m{\n  \"good\": 2,\n  \"some\": 0,\n  \"thing\": 1\n}\x1b[0m)\x1b[0m has unexpected keys: [\x1b[96;1m\"good\"\x1b[0m, \x1b[96;1m\"thing\"\x1b[0m]\x1b[0m",
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
		// `
		// 			{
		// 			 	val: nil
		// 			 	def:
		// 			 	  {
		// 						type Map
		// 						keys {
		// 							v {
		// 								type String
		// 								enum <all err ok none>
		// 								default "none"
		// 							}
		// 							s {
		// 								type String
		// 								enum <none stderr stdout all>
		// 								default "all"
		// 							}
		// 							exitOnError {
		// 								type Bool
		// 								default false
		// 							}
		// 						}
		// 					}
		// 			}
		// 		`,
		// "{ val: nil def: { type Map keys { some {type Int default 273} } } }",
	)
}
