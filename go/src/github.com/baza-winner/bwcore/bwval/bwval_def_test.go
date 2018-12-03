package bwval_test

import (
	"testing"

	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwparse"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwtype"
	"github.com/baza-winner/bwcore/bwval"
)

// func TestDefMarshalJSON(t *testing.T) {
// 	bwtesting.BwRunTests(t,
// 		bwjson.Pretty,
// 		func() map[string]bwtesting.Case {
// 			// tests := map[string]bwtesting.Case{}
// 			tests := map[string]bwtesting.Case{
// 				"": {
// 					In: []interface{}{
// 						bwval.Def{
// 							Types:      bwval.ValKindSetFrom(bwval.ValInt),
// 							IsOptional: true,
// 							Enum:       bwset.StringFrom("valueA", "valueB"),
// 							// Range:      bwval.IntRange{MinPtr: bwval.PtrToInt(-1), MaxPtr: bwval.PtrToInt(1)},
// 							Range: bwtype.MustRangeFrom(bwtype.A{Min: -1, Max: 1}),
// 							Keys: map[string]bwval.Def{
// 								"boolKey": {Types: bwval.ValKindSetFrom(bwval.ValBool)},
// 							},
// 							Elem: &bwval.Def{
// 								Types: bwval.ValKindSetFrom(bwval.ValBool),
// 							},
// 							ArrayElem: &bwval.Def{
// 								Types: bwval.ValKindSetFrom(bwval.ValBool),
// 							},
// 							Default: "default value",
// 						},
// 					},
// 					Out: []interface{}{
// 						"{\n  \"ArrayElem\": {\n    \"IsOptional\": false,\n    \"Types\": [\n      \"Bool\"\n    ]\n  },\n  \"Default\": \"default value\",\n  \"Elem\": {\n    \"IsOptional\": false,\n    \"Types\": [\n      \"Bool\"\n    ]\n  },\n  \"Enum\": [\n    \"valueA\",\n    \"valueB\"\n  ],\n  \"IsOptional\": true,\n  \"Range\": \"-1..1\",\n  \"Types\": [\n    \"Int\"\n  ],\n  \"keys\": {\n    \"boolKey\": {\n      \"IsOptional\": false,\n      \"Types\": [\n        \"Bool\"\n      ]\n    }\n  }\n}",
// 					},
// 				},
// 			}
// 			return tests
// 		}(),
// 	)
// }

func TestParseDef(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(s string) (result bwval.Def) {
			p := bwparse.From(bwrune.FromString(s))
			var st bwparse.Status
			if result, st = bwval.ParseDef(p); st.Err == nil {
				_, st.Err = bwparse.SkipSpace(p, bwparse.TillEOF)
			}
			if st.Err != nil {
				bwerr.PanicErr(st.Err)
			}
			return
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for _, v := range []struct {
				in  string
				out bwval.Def
			}{
				{in: "Int", out: bwval.Def{Types: bwtype.ValKindSetFrom(bwtype.ValInt)}},
				{in: `"String"`, out: bwval.Def{Types: bwtype.ValKindSetFrom(bwtype.ValString)}},
				{in: `[ Int "String"]`, out: bwval.Def{Types: bwtype.ValKindSetFrom(bwtype.ValString, bwtype.ValInt)}},
				{in: `<Int String >`, out: bwval.Def{Types: bwtype.ValKindSetFrom(bwtype.ValString, bwtype.ValInt)}},
				{in: "{ type Int }", out: bwval.Def{Types: bwtype.ValKindSetFrom(bwtype.ValInt)}},
				{in: `{ type "String" }`, out: bwval.Def{Types: bwtype.ValKindSetFrom(bwtype.ValString)}},
				{in: `{ type [ Int "String"] }`, out: bwval.Def{Types: bwtype.ValKindSetFrom(bwtype.ValString, bwtype.ValInt)}},
				{in: `{ type <Int String > }`, out: bwval.Def{Types: bwtype.ValKindSetFrom(bwtype.ValString, bwtype.ValInt)}},
				{in: `{ type String enum <some thing> }`,
					out: bwval.Def{
						Types: bwtype.ValKindSetFrom(bwtype.ValString),
						Enum:  bwset.StringFrom("some", "thing"),
					},
				},
				{in: `{ type String enum ["some" <thing>] }`,
					out: bwval.Def{
						Types: bwtype.ValKindSetFrom(bwtype.ValString),
						Enum:  bwset.StringFrom("some", "thing"),
					},
				},
				{in: `{ type String enum "some" }`,
					out: bwval.Def{
						Types: bwtype.ValKindSetFrom(bwtype.ValString),
						Enum:  bwset.StringFrom("some"),
					},
				},
				{in: `{ type Int range 0.. }`,
					out: bwval.Def{
						Types: bwtype.ValKindSetFrom(bwtype.ValInt),
						Range: bwtype.MustRangeFrom(bwtype.A{Min: 0}),
					},
				},
				{in: `{ type Number range 2.71..3.14 }`,
					out: bwval.Def{
						Types: bwtype.ValKindSetFrom(bwtype.ValNumber),
						Range: bwtype.MustRangeFrom(bwtype.A{Min: 2.71, Max: 3.14}),
					},
				},
				{in: `<ArrayOf String>`,
					out: bwval.Def{
						Types:     bwtype.ValKindSetFrom(bwtype.ValString),
						IsArrayOf: true,
					},
				},
				{in: `{ type Map keys { some Bool thing Int } }`,
					out: bwval.Def{
						Types: bwtype.ValKindSetFrom(bwtype.ValMap),
						Keys: map[string]bwval.Def{
							"some": {
								Types: bwtype.ValKindSetFrom(bwtype.ValBool),
							},
							"thing": {
								Types: bwtype.ValKindSetFrom(bwtype.ValInt),
							},
						},
					},
				},
				{in: `{ type Map elem Int }`,
					out: bwval.Def{
						Types: bwtype.ValKindSetFrom(bwtype.ValMap),
						Elem: &bwval.Def{
							Types: bwtype.ValKindSetFrom(bwtype.ValInt),
						},
					},
				},
				{in: `{ type Array arrayElem Int }`,
					out: bwval.Def{
						Types: bwtype.ValKindSetFrom(bwtype.ValArray),
						ArrayElem: &bwval.Def{
							Types: bwtype.ValKindSetFrom(bwtype.ValInt),
						},
					},
				},
				{in: `{ type Int isOptional true  }`,
					out: bwval.Def{
						Types:      bwtype.ValKindSetFrom(bwtype.ValInt),
						IsOptional: true,
					},
				},
				{in: `{ type <ArrayOf Int> default [1 2 3]  }`,
					out: bwval.Def{
						Types:      bwtype.ValKindSetFrom(bwtype.ValInt),
						IsArrayOf:  true,
						Default:    []interface{}{1, 2, 3},
						IsOptional: true,
					},
				},
				{in: `{
					type Map
					keys {
						some {
							type Int
							default 273
						}
					}
					default {}
				}`,
					out: bwval.Def{
						Types:      bwtype.ValKindSetFrom(bwtype.ValMap),
						Default:    map[string]interface{}{"some": 273},
						IsOptional: true,
						Keys: map[string]bwval.Def{
							"some": {
								Types:      bwtype.ValKindSetFrom(bwtype.ValInt),
								Default:    273,
								IsOptional: true,
							},
						},
					},
				},
				{in: `{
					type [ArrayOf Map]
					keys {
						some {
							type Int
							default 273
						}
					}
					default [ {} ]
				}`,
					out: bwval.Def{
						Types:      bwtype.ValKindSetFrom(bwtype.ValMap),
						IsArrayOf:  true,
						Default:    []interface{}{map[string]interface{}{"some": 273}},
						IsOptional: true,
						Keys: map[string]bwval.Def{
							"some": {
								Types:      bwtype.ValKindSetFrom(bwtype.ValInt),
								Default:    273,
								IsOptional: true,
							},
						},
					},
				},
				{in: `{
					type [ArrayOf Map]
					keys {
						some {
							type Int
							default 273
						}
					}
					default {}
				}`,
					out: bwval.Def{
						Types:      bwtype.ValKindSetFrom(bwtype.ValMap),
						IsArrayOf:  true,
						Default:    []interface{}{map[string]interface{}{"some": 273}},
						IsOptional: true,
						Keys: map[string]bwval.Def{
							"some": {
								Types:      bwtype.ValKindSetFrom(bwtype.ValInt),
								Default:    273,
								IsOptional: true,
							},
						},
					},
				},
				{in: `{
					type [Array]
					elem {
						type Map
						keys {
							some {
								type Int
								default 273
							}
						}
					}
					default [ { some 300} {some 400} ]
				}`,
					out: bwval.Def{
						Types:      bwtype.ValKindSetFrom(bwtype.ValArray),
						IsOptional: true,
						Elem: &bwval.Def{
							Types: bwtype.ValKindSetFrom(bwtype.ValMap),
							Keys: map[string]bwval.Def{
								"some": {
									Types:      bwtype.ValKindSetFrom(bwtype.ValInt),
									Default:    273,
									IsOptional: true,
								},
							},
						},
						Default: []interface{}{
							map[string]interface{}{"some": 300},
							map[string]interface{}{"some": 400},
						},
					},
				},
				{in: `{
					type [ArrayOf Map]
					default { some: 273 }
				}`,
					out: bwval.Def{
						Types:      bwtype.ValKindSetFrom(bwtype.ValMap),
						IsArrayOf:  true,
						Default:    []interface{}{map[string]interface{}{"some": 273}},
						IsOptional: true,
					},
				},
				{in: `{
					type [Array]
					default [ { some: 273 } ]
				}`,
					out: bwval.Def{
						Types:      bwtype.ValKindSetFrom(bwtype.ValArray),
						Default:    []interface{}{map[string]interface{}{"some": 273}},
						IsOptional: true,
					},
				},
				{in: `{
					type [Array]
					arrayElem {
						type Map
						keys {
							some {
								type Int
								isOptional true
							}
							thing {
								type String
								default "not bad"
							}
						}
					}
					default [ { some: 273 thing "good" } {} ]
				}`,
					out: bwval.Def{
						Types: bwtype.ValKindSetFrom(bwtype.ValArray),
						ArrayElem: &bwval.Def{
							Types: bwtype.ValKindSetFrom(bwtype.ValMap),
							Keys: map[string]bwval.Def{
								"some": bwval.Def{
									Types:      bwtype.ValKindSetFrom(bwtype.ValInt),
									IsOptional: true,
								},
								"thing": bwval.Def{
									Types:      bwtype.ValKindSetFrom(bwtype.ValString),
									Default:    "not bad",
									IsOptional: true,
								},
							},
						},
						Default: []interface{}{
							map[string]interface{}{"some": 273, "thing": "good"},
							map[string]interface{}{"thing": "not bad"},
						},
						IsOptional: true,
					},
				},
			} {
				tests[v.in] = bwtesting.Case{
					In:  []interface{}{v.in},
					Out: []interface{}{v.out},
				}
			}
			for _, v := range []struct {
				in  string
				out string
			}{
				{in: "Uknown",
					out: "unexpected `\x1b[91;1mUknown\x1b[0m`\x1b[0m at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\x1b[91mUknown\x1b[0m\n",
				},
				{in: `"Uknown"`,
					out: "unexpected `\x1b[91;1m\"Uknown\"\x1b[0m`\x1b[0m at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\x1b[91m\"Uknown\"\x1b[0m\n",
				},
				{in: "<Int Number>",
					out: "\x1b[91;1mNumber\x1b[0m can not be mixed with \x1b[96;1mInt\x1b[0m at pos \x1b[38;5;252;1m5\x1b[0m: \x1b[32m<Int \x1b[91mNumber\x1b[0m>\n",
				},
				{in: "<Number Int >",
					out: "\x1b[91;1mInt\x1b[0m can not be mixed with \x1b[96;1mNumber\x1b[0m at pos \x1b[38;5;252;1m8\x1b[0m: \x1b[32m<Number \x1b[91mInt\x1b[0m >\n",
				},
				{in: "<>",
					out: "expects non empty \x1b[97;1mArray\x1b[0m instead of unexpected `\x1b[91;1m<>\x1b[0m`\x1b[0m at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\x1b[91m<>\x1b[0m\n\x1b[0m",
				},
				{in: "[]",
					out: "expects non empty \x1b[97;1mArray\x1b[0m instead of unexpected `\x1b[91;1m[]\x1b[0m`\x1b[0m at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\x1b[91m[]\x1b[0m\n\x1b[0m",
				},
				{in: "{ type <> }",
					out: "expects non empty \x1b[97;1mArray\x1b[0m instead of unexpected `\x1b[91;1m<>\x1b[0m`\x1b[0m at pos \x1b[38;5;252;1m7\x1b[0m: \x1b[32m{ type \x1b[91m<>\x1b[0m }\n\x1b[0m",
				},
				{in: "{ type [] }",
					out: "expects non empty \x1b[97;1mArray\x1b[0m instead of unexpected `\x1b[91;1m[]\x1b[0m`\x1b[0m at pos \x1b[38;5;252;1m7\x1b[0m: \x1b[32m{ type \x1b[91m[]\x1b[0m }\n\x1b[0m",
				},
				{in: "{ Type <Number Int > }",
					out: "unexpected key `\x1b[91;1mType\x1b[0m`\x1b[0m at pos \x1b[38;5;252;1m2\x1b[0m: \x1b[32m{ \x1b[91mType\x1b[0m <Number Int > }\n",
				},
				{in: `{ "Type" <Number Int > }`,
					out: "unexpected key `\x1b[91;1m\"Type\"\x1b[0m`\x1b[0m at pos \x1b[38;5;252;1m2\x1b[0m: \x1b[32m{ \x1b[91m\"Type\"\x1b[0m <Number Int > }\n",
				},
				{in: `{ "type" Int type Bool }`,
					out: "duplicate key \x1b[91;1mtype\x1b[0m at pos \x1b[38;5;252;1m13\x1b[0m: \x1b[32m{ \"type\" Int \x1b[91mtype\x1b[0m Bool }\n",
				},
				{in: `{ "type" Int enum <some thing> }`,
					out: "unexpected key `\x1b[91;1menum\x1b[0m`\x1b[0m at pos \x1b[38;5;252;1m13\x1b[0m: \x1b[32m{ \"type\" Int \x1b[91menum\x1b[0m <some thing> }\n",
				},
				{in: `{ enum "some" type Int }`,
					out: "key \x1b[38;5;201;1menum\x1b[0m is specified, so value of key \x1b[38;5;201;1mtype\x1b[0m expects to have \x1b[96;1mString\x1b[0m at pos \x1b[38;5;252;1m19\x1b[0m: \x1b[32m{ enum \"some\" type \x1b[91mInt\x1b[0m }\n\x1b[0m",
				},
				{in: `{ type Int range 15 }`,
					out: "expects \x1b[97;1mRange\x1b[0m instead of unexpected char \x1b[96;1m'1'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m49\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m17\x1b[0m: \x1b[32m{ type Int range \x1b[91m1\x1b[0m5 }\n\x1b[0m",
				},
				{in: `{ type String enum [ 1 ] }`,
					out: "expects \x1b[97;1mString\x1b[0m instead of unexpected char \x1b[96;1m'1'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m49\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m21\x1b[0m: \x1b[32m{ type String enum [ \x1b[91m1\x1b[0m ] }\n\x1b[0m",
				},
				{in: `{ range 0.. type <String Array> }`,
					out: "key \x1b[38;5;201;1mrange\x1b[0m is specified, so value of key \x1b[38;5;201;1mtype\x1b[0m expects to have \x1b[96;1mInt\x1b[0m or \x1b[96;1mNumber\x1b[0m at pos \x1b[38;5;252;1m17\x1b[0m: \x1b[32m{ range 0.. type \x1b[91m<String Array>\x1b[0m }\n\x1b[0m",
				},
				{in: `ArrayOf`,
					out: "\x1b[97;1mArrayOf\x1b[0m must be followed by another \x1b[38;5;201;1mType\x1b[0m at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\x1b[91mArrayOf\x1b[0m\n\x1b[0m",
				},
				{in: `<ArrayOf>`,
					out: "\x1b[97;1mArrayOf\x1b[0m must be followed by another \x1b[38;5;201;1mType\x1b[0m at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\x1b[91m<ArrayOf>\x1b[0m\n\x1b[0m",
				},
				{in: `{ type <ArrayOf> }`,
					out: "\x1b[97;1mArrayOf\x1b[0m must be followed by another \x1b[38;5;201;1mType\x1b[0m at pos \x1b[38;5;252;1m7\x1b[0m: \x1b[32m{ type \x1b[91m<ArrayOf>\x1b[0m }\n\x1b[0m",
				},
				{in: `<ArrayOf Array>`,
					out: "\x1b[91;1mArray\x1b[0m can not be mixed with \x1b[96;1mArrayOf\x1b[0m at pos \x1b[38;5;252;1m9\x1b[0m: \x1b[32m<ArrayOf \x1b[91mArray\x1b[0m>\n",
				},
				{in: `[ Array <ArrayOf>]`,
					out: "\x1b[91;1mArrayOf\x1b[0m can not be mixed with \x1b[96;1mArray\x1b[0m at pos \x1b[38;5;252;1m9\x1b[0m: \x1b[32m[ Array <\x1b[91mArrayOf\x1b[0m>]\n",
				},
				{in: `{ type <ArrayOf Array> }`,
					out: "\x1b[91;1mArray\x1b[0m can not be mixed with \x1b[96;1mArrayOf\x1b[0m at pos \x1b[38;5;252;1m16\x1b[0m: \x1b[32m{ type <ArrayOf \x1b[91mArray\x1b[0m> }\n",
				},
				{in: `{ type Map keys { some Bool thing 1 } }`,
					out: "expects \x1b[97;1mDef\x1b[0m instead of unexpected char \x1b[96;1m'1'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m49\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m34\x1b[0m: \x1b[32m{ type Map keys { some Bool thing \x1b[91m1\x1b[0m } }\n\x1b[0m",
				},
				{in: `{ type Map elem 1 }`,
					out: "expects \x1b[97;1mDef\x1b[0m instead of unexpected char \x1b[96;1m'1'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m49\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m16\x1b[0m: \x1b[32m{ type Map elem \x1b[91m1\x1b[0m }\n\x1b[0m",
				},
				{in: `{ type Map keys Int }`,
					out: "expects \x1b[97;1mMap\x1b[0m instead of unexpected char \x1b[96;1m'I'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m73\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m16\x1b[0m: \x1b[32m{ type Map keys \x1b[91mI\x1b[0mnt }\n\x1b[0m",
				},
				{in: `{ type Int isOptional false  default 1 }`,
					out: "while \x1b[38;5;201;1misOptional \x1b[96;1mfalse\x1b[0m, key \x1b[38;5;201;1mdefault\x1b[0m can not be specified at pos \x1b[38;5;252;1m29\x1b[0m: \x1b[32m{ type Int isOptional false  \x1b[91mdefault\x1b[0m 1 }\n\x1b[0m",
				},
				{in: `{ keys { some Int } type Array }`,
					out: "key \x1b[38;5;201;1mkeys\x1b[0m is specified, so value of key \x1b[38;5;201;1mtype\x1b[0m expects to have \x1b[96;1mMap\x1b[0m at pos \x1b[38;5;252;1m25\x1b[0m: \x1b[32m{ keys { some Int } type \x1b[91mArray\x1b[0m }\n\x1b[0m",
				},
				{in: `{ arrayElem Int type <Int Map> }`,
					out: "key \x1b[38;5;201;1marrayElem\x1b[0m is specified, so value of key \x1b[38;5;201;1mtype\x1b[0m expects to have \x1b[96;1mArray\x1b[0m at pos \x1b[38;5;252;1m21\x1b[0m: \x1b[32m{ arrayElem Int type \x1b[91m<Int Map>\x1b[0m }\n\x1b[0m",
				},
				{in: `{ arrayElem Int type <Int Map> }`,
					out: "key \x1b[38;5;201;1marrayElem\x1b[0m is specified, so value of key \x1b[38;5;201;1mtype\x1b[0m expects to have \x1b[96;1mArray\x1b[0m at pos \x1b[38;5;252;1m21\x1b[0m: \x1b[32m{ arrayElem Int type \x1b[91m<Int Map>\x1b[0m }\n\x1b[0m",
				},
				{in: `{ elem Int type String }`,
					out: "key \x1b[38;5;201;1melem\x1b[0m is specified, so value of key \x1b[38;5;201;1mtype\x1b[0m expects to have \x1b[96;1mMap\x1b[0m or \x1b[96;1mArray\x1b[0m at pos \x1b[38;5;252;1m16\x1b[0m: \x1b[32m{ elem Int type \x1b[91mString\x1b[0m }\n\x1b[0m",
				},
				{in: `{ arrayElem Bool elem Int type Array }`,
					out: "key \x1b[38;5;201;1melem\x1b[0m is specified, so value of key \x1b[38;5;201;1mtype\x1b[0m expects to have \x1b[96;1mMap\x1b[0m at pos \x1b[38;5;252;1m31\x1b[0m: \x1b[32m{ arrayElem Bool elem Int type \x1b[91mArray\x1b[0m }\n\x1b[0m",
				},
				{in: `{
					type [ArrayOf Map]
					keys {
						some {
							type Int
							default 273
						}
					}
					default [ {some true} ]
				}`,
					out: "expects \x1b[97;1mInt\x1b[0m instead of unexpected char \x1b[96;1m't'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m116\x1b[0m)\x1b[0m at line \x1b[38;5;252;1m9\x1b[0m, col \x1b[38;5;252;1m22\x1b[0m (pos \x1b[38;5;252;1m122\x1b[0m)\x1b[0m:\n\x1b[32m\t\t\t\t\t\t\tdefault 273\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t\tdefault [ {some \x1b[91mt\x1b[0mrue} ]\n\t\t\t\t}\n",
				},
				{in: `{
					type [ArrayOf Map]
					keys {
						some {
							type Int
							default 273
						}
					}
					default [ {thing true} ]
				}`,
					out: "unexpected key `\x1b[91;1mthing\x1b[0m`\x1b[0m at line \x1b[38;5;252;1m9\x1b[0m, col \x1b[38;5;252;1m17\x1b[0m (pos \x1b[38;5;252;1m117\x1b[0m)\x1b[0m:\n\x1b[32m\t\t\t\t\t\t\tdefault 273\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t\tdefault [ {\x1b[91mthing\x1b[0m true} ]\n\t\t\t\t}\n",
				},
				{in: `{
					type [ArrayOf Map]
					keys {
						some {
							type Int
							range 200..
							default 273
						}
					}
					default [ {some 1} ]
				}`,
					out: "\x1b[38;5;252;1m$Def.default.0.some\x1b[0m (\x1b[96;1m1\x1b[0m)\x1b[0m is \x1b[91;1mout of range\x1b[0m \x1b[96;1m200..\x1b[0m at line \x1b[38;5;252;1m10\x1b[0m, col \x1b[38;5;252;1m22\x1b[0m (pos \x1b[38;5;252;1m141\x1b[0m)\x1b[0m:\n\x1b[32m\t\t\t\t\t\t\tdefault 273\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t\tdefault [ {some \x1b[91m1\x1b[0m} ]\n\t\t\t\t}\n",
				},
				{in: `{
					type [ArrayOf Map]
					elem {
							type Int
							range 200..
							default 273
					}
					default [ {some 1} ]
				}`,
					out: "\x1b[38;5;252;1m$Def.default.0.some\x1b[0m (\x1b[96;1m1\x1b[0m)\x1b[0m is \x1b[91;1mout of range\x1b[0m \x1b[96;1m200..\x1b[0m at line \x1b[38;5;252;1m8\x1b[0m, col \x1b[38;5;252;1m22\x1b[0m (pos \x1b[38;5;252;1m120\x1b[0m)\x1b[0m:\n\x1b[32m\t\t\t\t\t\t\trange 200..\n\t\t\t\t\t\t\tdefault 273\n\t\t\t\t\t}\n\t\t\t\t\tdefault [ {some \x1b[91m1\x1b[0m} ]\n\t\t\t\t}\n",
				},
				{in: `{
					default [ {some 1} ]
					type [ArrayOf Map]
					elem {
							default 273
							type Int
							range 300..
					}
				}`,
					out: "key \x1b[38;5;201;1mtype\x1b[0m must be specified first at line \x1b[38;5;252;1m2\x1b[0m, col \x1b[38;5;252;1m6\x1b[0m (pos \x1b[38;5;252;1m7\x1b[0m)\x1b[0m:\n\x1b[32m{\n\t\t\t\t\t\x1b[91mdefault\x1b[0m [ {some 1} ]\n\t\t\t\t\ttype [ArrayOf Map]\n\t\t\t\t\telem {\n\t\t\t\t\t\t\tdefault 273\n\x1b[0m",
				},
				{in: `{
					type [ArrayOf Map]
					default [ {some 1} ]
					elem {
							default 273
							type Int
							range 300..
					}
				}`,
					out: "key \x1b[38;5;201;1mdefault\x1b[0m must be LAST specified, but found key \x1b[91;1melem\x1b[0m after it at line \x1b[38;5;252;1m4\x1b[0m, col \x1b[38;5;252;1m6\x1b[0m (pos \x1b[38;5;252;1m57\x1b[0m)\x1b[0m:\n\x1b[32m{\n\t\t\t\t\ttype [ArrayOf Map]\n\t\t\t\t\tdefault [ {some 1} ]\n\t\t\t\t\t\x1b[91melem\x1b[0m {\n\t\t\t\t\t\t\tdefault 273\n\t\t\t\t\t\t\ttype Int\n\t\t\t\t\t\t\trange 300..\n\x1b[0m",
				},
				{in: `{
					type [Array]
					elem {
						type Map
						keys {
							some {
								type Int
								range 1..
								default 273
							}
						}
					}
					default [ {thing 1} ]
				}`,
					out: "unexpected key `\x1b[91;1mthing\x1b[0m`\x1b[0m at line \x1b[38;5;252;1m13\x1b[0m, col \x1b[38;5;252;1m17\x1b[0m (pos \x1b[38;5;252;1m169\x1b[0m)\x1b[0m:\n\x1b[32m\t\t\t\t\t\t\t}\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t\tdefault [ {\x1b[91mthing\x1b[0m 1} ]\n\t\t\t\t}\n",
				},
				{in: `{
					type [Array]
					arrayElem {
						type Map
						keys {
							some {
								type Int
								range 1..
								default 273
							}
						}
					}
					default [ {thing 1} ]
				}`,
					out: "unexpected key `\x1b[91;1mthing\x1b[0m`\x1b[0m at line \x1b[38;5;252;1m13\x1b[0m, col \x1b[38;5;252;1m17\x1b[0m (pos \x1b[38;5;252;1m174\x1b[0m)\x1b[0m:\n\x1b[32m\t\t\t\t\t\t\t}\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t\tdefault [ {\x1b[91mthing\x1b[0m 1} ]\n\t\t\t\t}\n",
				},

				// {in: `{
				// 	type [ArrayOf Map]
				// 	keys {
				// 		some {
				// 			type Int
				// 			default 273
				// 		}
				// 	}
				// 	default [ { some 300} {thing 400} ]
				// }`,
				// 	out: " ",
				// },

				// {in: `{
				// 	default [ {some 1} ]
				// 	type [ArrayOf Map]
				// 	elem {
				// 			type Int
				// 			range 200..
				// 			default 273
				// 	}
				// }`,
				// 	out: " ",
				// 	// out: "\x1b[38;5;252;1m$Def.default.some\x1b[0m (\x1b[96;1m1\x1b[0m)\x1b[0m is \x1b[91;1mout of range\x1b[0m \x1b[96;1m200..\x1b[0m at line \x1b[38;5;252;1m8\x1b[0m, col \x1b[38;5;252;1m22\x1b[0m (pos \x1b[38;5;252;1m120\x1b[0m)\x1b[0m:\n\x1b[32m\t\t\t\t\t\t\trange 200..\n\t\t\t\t\t\t\tdefault 273\n\t\t\t\t\t}\n\t\t\t\t\tdefault [ {some \x1b[91m1\x1b[0m} ]\n\t\t\t\t}\n",
				// },
				// {in: `{
				// 	type [ArrayOf Map]
				// 	elem {
				// 		type Int
				// 		range 200..
				// 		default 273
				// 	}
				// 	default [ {some 1} ]
				// }`,
				// 	out: "\x1b[38;5;252;1m.\x1b[0m (\x1b[96;1m1\x1b[0m)\x1b[0m is \x1b[91;1mout of range\x1b[0m \x1b[96;1m200..\x1b[0m at line \x1b[38;5;252;1m10\x1b[0m, col \x1b[38;5;252;1m22\x1b[0m (pos \x1b[38;5;252;1m141\x1b[0m)\x1b[0m:\n\x1b[32m\t\t\t\t\t\t\tdefault 273\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t\tdefault [ {some \x1b[91m1\x1b[0m} ]\n\t\t\t\t}\n",
				// },
			} {
				tests[v.in] = bwtesting.Case{
					In:    []interface{}{v.in},
					Panic: v.out,
				}
			}
			return tests
		}(),
		// `{
		// 			type [Array]
		// 			elem {
		// 				type Map
		// 				keys {
		// 					some {
		// 						type Int
		// 						range 1..
		// 						default 273
		// 					}
		// 				}
		// 			}
		// 			default [ {thing 1} ]
		// 		}`,
		// `{
		// 			type [ArrayOf Map]
		// 			keys {
		// 				some {
		// 					type Int
		// 					default 273
		// 				}
		// 			}
		// 			default [ {} ]
		// 		}`,
		// `{
		// 			default [ {some 1} ]
		// 			type [ArrayOf Map]
		// 			elem {
		// 					default 273
		// 					type Int
		// 					range 300..
		// 			}
		// 		}`,
		// `{
		// 			default [ {some 1} ]
		// 			type [ArrayOf Map]
		// 			elem {
		// 					type Int
		// 					range 200..
		// 					default 273
		// 			}
		// 		}`,
		// `{
		// 			type [ArrayOf Map]
		// 			default { some: 273 }
		// 		}`,

		// `{ type Map keys { some Bool thing 1 } }`,
		// "<>",
		// `{ arrayElem Bool elem Int type Array }`,
	)
}

// func TestDefFrom(t *testing.T) {
// 	bwtesting.BwRunTests(t,
// 		bwval.DefFrom,
// 		func() map[string]bwtesting.Case {
// 			tests := map[string]bwtesting.Case{}
// 			for k, v := range map[string]bwval.Def{
// 				"Bool": {
// 					Types: bwval.ValKindSetFrom(bwval.ValBool),
// 				},
// 				`{ type String enum <a b c> }`: {
// 					Types: bwval.ValKindSetFrom(bwval.ValString),
// 					Enum:  bwset.StringFrom("a", "b", "c")},
// 				`{ type Map keys { a Bool } }`: {
// 					Types: bwval.ValKindSetFrom(bwval.ValMap),
// 					Keys:  map[string]bwval.Def{"a": {Types: bwval.ValKindSetFrom(bwval.ValBool)}},
// 				},
// 				`{ type Array arrayElem Bool }`: {
// 					Types:     bwval.ValKindSetFrom(bwval.ValArray),
// 					ArrayElem: &bwval.Def{Types: bwval.ValKindSetFrom(bwval.ValBool)},
// 				},
// 				`{ type Int min 1 max 2 }`: {
// 					Types: bwval.ValKindSetFrom(bwval.ValInt),
// 					Range: bwtype.MustRangeFrom(bwtype.A{Min: 1, Max: 2}),
// 				},
// 				`{ type Float64 min 1 max 2 }`: {
// 					Types: bwval.ValKindSetFrom(bwval.ValFloat64),
// 					Range: bwtype.MustRangeFrom(bwtype.A{Min: 1, Max: 2}),
// 				},
// 				`{ type String default "some" }`: {
// 					Types:      bwval.ValKindSetFrom(bwval.ValString),
// 					Default:    "some",
// 					IsOptional: true,
// 				},
// 				`{ type [ArrayOf Int] default [1 2 3] }`: {
// 					Types:      bwval.ValKindSetFrom(bwval.ValArrayOf, bwval.ValInt),
// 					Default:    []interface{}{1, 2, 3},
// 					IsOptional: true,
// 				},
// 				`{ type Int min 2.0 }`: {
// 					Types:      bwval.ValKindSetFrom(bwval.ValInt),
// 					IsOptional: false,
// 					Range:      bwtype.MustRangeFrom(bwtype.A{Min: 2}),
// 				},
// 			} {
// 				tests[k] = bwtesting.Case{
// 					In:  []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
// 					Out: []interface{}{v},
// 				}
// 			}

// 			for k, v := range map[string]string{
// 				"nil":                                           "\x1b[38;5;252;1m$Def\x1b[0m is \x1b[91;1m(interface {})<nil>\x1b[0m",
// 				"true":                                          "\x1b[38;5;252;1m$Def\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m neither \x1b[97;1mString\x1b[0m nor \x1b[97;1mArray\x1b[0m nor \x1b[97;1mMap\x1b[0m",
// 				"[ Bool true ]":                                 "\x1b[38;5;252;1m$Def.1\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mString\x1b[0m",
// 				"{type true}":                                   "\x1b[38;5;252;1m$Def.type\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m neither \x1b[97;1mString\x1b[0m nor \x1b[97;1mArray\x1b[0m",
// 				`{ type [ Int "bool" ] }`:                       "\x1b[38;5;252;1m$Def.type.1\x1b[0m (\x1b[96;1m\"bool\"\x1b[0m)\x1b[0m is \x1b[91;1mnon supported\x1b[0m value\x1b[0m",
// 				`{ type String enum [ "a" true ] }`:             "\x1b[38;5;252;1m$Def.enum.1\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mString\x1b[0m",
// 				`{ type Number min 1 max 2 default 3 }`:         "\x1b[38;5;252;1m$Def.default\x1b[0m (\x1b[96;1m3\x1b[0m)\x1b[0m is \x1b[91;1mout of range\x1b[0m \x1b[96;1m1..2\x1b[0m",
// 				`{ enum <Bool Int> }`:                           "\x1b[38;5;252;1m$Def\x1b[0m (\x1b[96;1m{\n  \"enum\": [\n    \"Bool\",\n    \"Int\"\n  ]\n}\x1b[0m)\x1b[0m has no key \x1b[96;1mtype\x1b[0m",
// 				`{ type Map keys true }`:                        "\x1b[38;5;252;1m$Def.keys\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mMap\x1b[0m",
// 				`{ type Map keys { some true } }`:               "\x1b[38;5;252;1m$Def.keys.some\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m neither \x1b[97;1mString\x1b[0m nor \x1b[97;1mArray\x1b[0m nor \x1b[97;1mMap\x1b[0m",
// 				`{ type Array arrayElem true }`:                 "\x1b[38;5;252;1m$Def.arrayElem\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m neither \x1b[97;1mString\x1b[0m nor \x1b[97;1mArray\x1b[0m nor \x1b[97;1mMap\x1b[0m",
// 				`{ type Array elem true }`:                      "\x1b[38;5;252;1m$Def.elem\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m neither \x1b[97;1mString\x1b[0m nor \x1b[97;1mArray\x1b[0m nor \x1b[97;1mMap\x1b[0m",
// 				`{ type Int min 2.1 }`:                          "\x1b[38;5;252;1m$Def.min\x1b[0m (\x1b[96;1m2.1\x1b[0m)\x1b[0m is not \x1b[97;1mInt\x1b[0m",
// 				`{ type Int max 2.1 }`:                          "\x1b[38;5;252;1m$Def.max\x1b[0m (\x1b[96;1m2.1\x1b[0m)\x1b[0m is not \x1b[97;1mInt\x1b[0m",
// 				`{ type Int min 3 max 2 }`:                      "\x1b[38;5;252;1m$Def\x1b[0m (\x1b[96;1m{\n  \"max\": 2,\n  \"min\": 3,\n  \"type\": \"Int\"\n}\x1b[0m)\x1b[0m: \x1b[38;5;252;1m.max\x1b[0m (\x1b[96;1m2\x1b[0m) must not be \x1b[91;1mless\x1b[0m then \x1b[38;5;252;1m.min\x1b[0m (\x1b[96;1m3\x1b[0m)\x1b[0m",
// 				`{ type Number min true }`:                      "\x1b[38;5;252;1m$Def.min\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mNumber\x1b[0m",
// 				`{ type Number max true }`:                      "\x1b[38;5;252;1m$Def.max\x1b[0m (\x1b[96;1mtrue\x1b[0m)\x1b[0m is not \x1b[97;1mNumber\x1b[0m",
// 				`{ type Number min 3.14 max 2.71 }`:             "\x1b[38;5;252;1m$Def\x1b[0m (\x1b[96;1m{\n  \"max\": 2.71,\n  \"min\": 3.14,\n  \"type\": \"Number\"\n}\x1b[0m)\x1b[0m: \x1b[38;5;252;1m.max\x1b[0m (\x1b[96;1m2.71\x1b[0m) must not be \x1b[91;1mless\x1b[0m then \x1b[38;5;252;1m.min\x1b[0m (\x1b[96;1m3.14\x1b[0m)\x1b[0m",
// 				`{ type Number isOptional 3 }`:                  "\x1b[38;5;252;1m$Def.isOptional\x1b[0m (\x1b[96;1m3\x1b[0m)\x1b[0m is not \x1b[97;1mBool\x1b[0m",
// 				`{ type Number isOptional false default 3.14 }`: "\x1b[38;5;252;1m$Def\x1b[0m (\x1b[96;1m{\n  \"default\": 3.14,\n  \"isOptional\": false,\n  \"type\": \"Number\"\n}\x1b[0m)\x1b[0m: having \x1b[38;5;252;1m.default\x1b[0m can not have \x1b[38;5;252;1m.isOptional\x1b[0m \x1b[96;1mtrue\x1b[0m",
// 				`{ type Number keyA "valA" keyB "valB" }`:       "\x1b[38;5;252;1m$Def\x1b[0m (\x1b[96;1m{\n  \"keyA\": \"valA\",\n  \"keyB\": \"valB\",\n  \"type\": \"Number\"\n}\x1b[0m)\x1b[0m has unexpected keys: \x1b[96;1m[\n  \"keyA\",\n  \"keyB\"\n]\x1b[0m",
// 				`{ type Number keyA "valA" }`:                   "\x1b[38;5;252;1m$Def\x1b[0m (\x1b[96;1m{\n  \"keyA\": \"valA\",\n  \"type\": \"Number\"\n}\x1b[0m)\x1b[0m has unexpected key \x1b[96;1m\"keyA\"\x1b[0m",
// 				`{ type <ArrayOf> }`:                            "\x1b[38;5;252;1m$Def.type\x1b[0m (\x1b[96;1m[\n  \"ArrayOf\"\n]\x1b[0m)\x1b[0m: \x1b[96;1mArrayOf\x1b[0m must be followed by some type, can not be \x1b[91;1mused alone\x1b[0m",
// 				`{ type <ArrayOf Array> }`:                      "\x1b[38;5;252;1m$Def.type\x1b[0m (\x1b[96;1m[\n  \"ArrayOf\",\n  \"Array\"\n]\x1b[0m)\x1b[0m: values \x1b[96;1m\"ArrayOf\"\x1b[0m and \x1b[96;1m\"Array\"\x1b[0m are \x1b[91;1mmutually exclusive\x1b[0m, can not be \x1b[91;1mused both at once\x1b[0m",
// 				`{ type <Int Number> }`:                         "\x1b[38;5;252;1m$Def.type\x1b[0m (\x1b[96;1m[\n  \"Int\",\n  \"Number\"\n]\x1b[0m)\x1b[0m: values \x1b[96;1m\"Int\"\x1b[0m and \x1b[96;1m\"Number\"\x1b[0m are \x1b[91;1mmutually exclusive\x1b[0m, can not be \x1b[91;1mused both at once\x1b[0m",
// 			} {
// 				tests[k] = bwtesting.Case{
// 					In:    []interface{}{func(testName string) interface{} { return bwval.From(testName) }},
// 					Panic: v,
// 				}
// 			}
// 			return tests
// 		}(),
// 		// tests,
// 		// "{ type Int min 2.1 }",
// 	)
// }
