package bwparse_test

import (
	"fmt"
	"testing"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwparse"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwtype"
)

func TestUnexpected(t *testing.T) {
	testUnexpectedHelper := func(s string, ofs uint) *bwparse.P {
		p := bwparse.From(bwrune.FromString(s))
		p.Forward(ofs)
		return p
	}
	bwtesting.BwRunTests(t,
		"UnexpectedA",
		func() map[string]bwtesting.Case {
			p := testUnexpectedHelper("some", 2)
			pi := p.LookAhead(3)
			tests := map[string]bwtesting.Case{}
			tests["panic"] = bwtesting.Case{
				V: p,
				// In:    []interface{}{bwparse.PosInfo{Pos: 4}},
				In:    []interface{}{bwparse.UnexpectedA{PosInfo: pi}},
				Panic: "\x1b[38;5;201;1mps.pos\x1b[0m (\x1b[96;1m4\x1b[0m) > \x1b[38;5;201;1mp.curr.pos\x1b[0m (\x1b[96;1m1\x1b[0m)\x1b[0m",
			}
			p = testUnexpectedHelper("{\n key wrong \n} ", 0)
			p.Forward(7)
			pi = p.Curr()
			p.Forward(5)
			tests["normal"] = bwtesting.Case{
				V: p,
				// In:  []interface{}{pi},
				In:  []interface{}{bwparse.UnexpectedA{PosInfo: pi}},
				Out: []interface{}{"unexpected \x1b[91;1m\"wrong\"\x1b[0m at line \x1b[38;5;252;1m2\x1b[0m, col \x1b[38;5;252;1m6\x1b[0m (pos \x1b[38;5;252;1m7\x1b[0m)\x1b[0m:\n\x1b[32m{\n key \x1b[91mwrong\x1b[0m \n} \n"},
			}
			return tests
		}(),
	)
}

func TestLookAhead(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(p *bwparse.P, i int) rune {
			return p.LookAhead(uint(i)).Rune()
		},
		func() map[string]bwtesting.Case {
			s := "s\no\nm\ne\nt\nhing"
			p := bwparse.From(bwrune.FromString(s))
			p.Forward(0)
			tests := map[string]bwtesting.Case{}
			for i, r := range s {
				tests[fmt.Sprintf("%d", i)] = bwtesting.Case{
					In:  []interface{}{p, i},
					Out: []interface{}{r},
				}
			}
			return tests
		}(),
	)
}

func TestPath(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(s string, optBases ...[]bw.ValPath) (result bw.ValPath) {
			var err error
			if result, err = func(s string, optBases ...[]bw.ValPath) (result bw.ValPath, err error) {
				defer func() {
					if err != nil {
						result = nil
					}
				}()
				opt := bwparse.PathOpt{}
				if len(optBases) > 0 {
					opt.Bases = optBases[0]
				}
				p := bwparse.From(bwrune.FromString(s))
				if result, err = bwparse.PathContent(p, opt); err == nil {
					err = end(p, true)
				}
				return
			}(s, optBases...); err != nil {
				bwerr.PanicErr(err)
			}
			return result
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for k, v := range map[string]bw.ValPath{
				".": bw.ValPath{},
				"some.thing": bw.ValPath{
					{Type: bw.ValPathItemKey, Key: "some"},
					{Type: bw.ValPathItemKey, Key: "thing"},
				},
				"some.1": bw.ValPath{
					{Type: bw.ValPathItemKey, Key: "some"},
					{Type: bw.ValPathItemIdx, Idx: 1},
				},
				"some.#": bw.ValPath{
					{Type: bw.ValPathItemKey, Key: "some"},
					{Type: bw.ValPathItemHash},
				},
				"(some.thing).good": bw.ValPath{
					{Type: bw.ValPathItemPath,
						Path: bw.ValPath{
							{Type: bw.ValPathItemKey, Key: "some"},
							{Type: bw.ValPathItemKey, Key: "thing"},
						},
					},
					{Type: bw.ValPathItemKey, Key: "good"},
				},
				"$some.thing.(good)": bw.ValPath{
					{Type: bw.ValPathItemVar, Key: "some"},
					{Type: bw.ValPathItemKey, Key: "thing"},
					{Type: bw.ValPathItemPath,
						Path: bw.ValPath{
							{Type: bw.ValPathItemKey, Key: "good"},
						},
					},
				},
				"1.some": bw.ValPath{
					{Type: bw.ValPathItemIdx, Idx: 1},
					{Type: bw.ValPathItemKey, Key: "some"},
				},
				"-1.some": bw.ValPath{
					{Type: bw.ValPathItemIdx, Idx: -1},
					{Type: bw.ValPathItemKey, Key: "some"},
				},
				"2?": bw.ValPath{
					{Type: bw.ValPathItemIdx, Idx: 2, IsOptional: true},
				},
				"some.2?": bw.ValPath{
					{Type: bw.ValPathItemKey, Key: "some"},
					{Type: bw.ValPathItemIdx, Idx: 2, IsOptional: true},
				},
			} {
				tests[k] = bwtesting.Case{
					In:  []interface{}{func(testName string) string { return testName }},
					Out: []interface{}{v},
				}
			}
			for k, v := range map[string]string{
				"":          "unexpected end of string at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\n",
				"1.":        "unexpected end of string at pos \x1b[38;5;252;1m2\x1b[0m: \x1b[32m1.\n",
				"1.@":       "unexpected char \u001b[96;1m'@'\u001b[0m (\u001b[38;5;201;1mcharCode\u001b[0m: \u001b[96;1m64\u001b[0m)\u001b[0m at pos \u001b[38;5;252;1m2\u001b[0m: \u001b[32m1.\u001b[91m@\u001b[0m\n",
				"-a":        "unexpected char \u001b[96;1m'a'\u001b[0m (\u001b[38;5;201;1mcharCode\u001b[0m: \u001b[96;1m97\u001b[0m)\u001b[0m at pos \u001b[38;5;252;1m1\u001b[0m: \u001b[32m-\u001b[91ma\u001b[0m\n",
				"1a":        "unexpected char \u001b[96;1m'a'\u001b[0m (\u001b[38;5;201;1mcharCode\u001b[0m: \u001b[96;1m97\u001b[0m)\u001b[0m at pos \u001b[38;5;252;1m1\u001b[0m: \u001b[32m1\u001b[91ma\u001b[0m\n",
				"12.#.4":    "unexpected char \x1b[96;1m'.'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m46\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m4\x1b[0m: \x1b[32m12.#\x1b[91m.\x1b[0m4\n",
				"12.(4":     "unexpected end of string at pos \u001b[38;5;252;1m5\u001b[0m: \u001b[32m12.(4\n",
				"12.$a":     "unexpected char \u001b[96;1m'$'\u001b[0m (\u001b[38;5;201;1mcharCode\u001b[0m: \u001b[96;1m36\u001b[0m)\u001b[0m at pos \u001b[38;5;252;1m3\u001b[0m: \u001b[32m12.\u001b[91m$\u001b[0ma\n",
				"$1.some":   "unexpected base path idx \x1b[96;1m1\x1b[0m (len(bases): \x1b[96;1m0)\x1b[0m at pos \x1b[38;5;252;1m1\x1b[0m: \x1b[32m$\x1b[91m1\x1b[0m.some\n",
				"some.(2?)": "unexpected char \x1b[96;1m'?'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m63\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m7\x1b[0m: \x1b[32msome.(2\x1b[91m?\x1b[0m)\n",
			} {
				tests[k] = bwtesting.Case{
					In:    []interface{}{func(testName string) string { return testName }},
					Panic: v,
				}
			}
			tests["$0.some"] = bwtesting.Case{
				In: []interface{}{
					func(testName string) string { return testName },
					[]bw.ValPath{{bw.ValPathItem{Type: bw.ValPathItemKey, Key: "thing"}}},
				},
				Out: []interface{}{
					bw.ValPath{
						{Type: bw.ValPathItemKey, Key: "thing"},
						{Type: bw.ValPathItemKey, Key: "some"},
					},
				},
			}
			tests[".some"] = bwtesting.Case{
				In: []interface{}{
					func(testName string) string { return testName },
					[]bw.ValPath{{bw.ValPathItem{Type: bw.ValPathItemKey, Key: "thing"}}},
				},
				Out: []interface{}{
					bw.ValPath{
						{Type: bw.ValPathItemKey, Key: "thing"},
						{Type: bw.ValPathItemKey, Key: "some"},
					},
				},
			}
			return tests
		}(),
	)
}

func TestInt(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(s string) (result interface{}) {
			var err error
			if result, err = func(s string) (result int, err error) {
				defer func() {
					if err != nil {
						result = 0
					}
				}()
				p := bwparse.From(bwrune.FromString(s))
				var ok bool
				if result, _, ok, err = bwparse.Int(p); err == nil {
					err = end(p, ok)
				}
				return
			}(s); err != nil {
				bwerr.PanicErr(err)
			}
			return result
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for k, v := range map[string]int{
				"0":                                  0,
				"-273":                               -273,
				"+1_000_000":                         1000000,
				"+1_000_000_000_000_000_000_000_000": 1000000,
			} {
				tests[k] = bwtesting.Case{
					In:  []interface{}{func(testName string) string { return testName }},
					Out: []interface{}{v},
				}
			}
			for k, v := range map[string]string{
				"+1_000_000_000_000_000_000_000_000": "strconv.ParseInt: parsing \"1000000000000000000000000\": value out of range at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\x1b[91m+1_000_000_000_000_000_000_000_000\x1b[0m\n",
			} {
				tests[k] = bwtesting.Case{
					In:    []interface{}{func(testName string) string { return testName }},
					Panic: v,
				}
			}
			return tests
		}(),
	)
}

func end(p *bwparse.P, ok bool) (err error) {
	if !ok {
		err = bwparse.Unexpected(p)
	} else {
		err = bwparse.SkipSpace(p, bwparse.TillEOF)
	}
	return
}

// func TestParseRange(t *testing.T) {
// 	bwtesting.BwRunTests(t,
// 		func(s string) (result interface{}) {
// 		},
// 		func() map[string]bwtesting.Case {

// 			tests := map[string]bwtesting.Case{}
// 			return tests
// 		},
// 	)
// }

func TestVal(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(s string, opt bwparse.Opt) (result interface{}) {
			var err error
			if result, err = func(s string) (result interface{}, err error) {
				defer func() {
					if err != nil {
						result = nil
					}
				}()
				p := bwparse.From(bwrune.FromString(s))

				var ok bool
				if result, _, ok, err = bwparse.Val(p, opt); err == nil {
					err = end(p, ok)
				}
				// 	bwparse.Opt{
				// 	OnId: func(p bwparse.I, s string, start *bwparse.PosInfo) (result interface{}, ok bool, err error) {
				// 		defIds := bwset.StringFrom("Bool", "String", "Int", "Number", "Array", "ArrayOf")
				// 		if ok = defIds.Has(s); ok {
				// 			result = s
				// 		}
				// 		return
				// 	},
				// }
				return
			}(s); err != nil {
				bwerr.PanicErr(err)
			}
			return result
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for _, v := range []struct {
				in  string
				opt bwparse.Opt
				out interface{}
			}{
				{in: "nil", out: nil},
				{in: "true", out: true},
				{in: "false", out: false},
				{in: "0", out: bwtype.MustNumberFrom(0)},

				{in: "0..1", out: bwtype.MustRangeFrom(bwtype.A{Min: 0, Max: 1})},
				{in: "0.5..1", out: bwtype.MustRangeFrom(bwtype.A{Min: 0.5, Max: 1})},
				{in: "..3.14", out: bwtype.MustRangeFrom(bwtype.A{Max: 3.14})},
				{in: "..", out: bwtype.MustRangeFrom(bwtype.A{})},
				{in: "$idx.3..{{some.thing}}", out: bwtype.MustRangeFrom(bwtype.A{
					Min: bw.ValPath{
						bw.ValPathItem{Type: bw.ValPathItemVar, Key: "idx"},
						bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: 3},
					},
					Max: bw.ValPath{
						bw.ValPathItem{Type: bw.ValPathItemKey, Key: "some"},
						bw.ValPathItem{Type: bw.ValPathItemKey, Key: "thing"},
					},
				})},

				{in: "-1_000_000", out: bwtype.MustNumberFrom(-1000000)},
				{in: "+3.14", out: bwtype.MustNumberFrom(3.14)},
				{in: "+2.0", out: bwtype.MustNumberFrom(2)},
				{in: "[0, 1]", out: []interface{}{bwtype.MustNumberFrom(0), bwtype.MustNumberFrom(1)}},
				{in: `"a"`, out: "a"},
				{in: `<a b c>`, out: []string{"a", "b", "c"}},
				{in: `[<a b c>]`, out: []interface{}{"a", "b", "c"}},
				{in: `["x" <a b c> "z"]`, out: []interface{}{"x", "a", "b", "c", "z"}},
				{in: `{ key "value" bool true }`, out: map[string]interface{}{
					"key":  "value",
					"bool": true,
				}},
				{in: `{ key => "\"value\n", 'bool': true keyword Bool}`,
					opt: bwparse.Opt{IdVals: map[string]interface{}{"Bool": "Bool"}},
					out: map[string]interface{}{
						"key":     "\"value\n",
						"bool":    true,
						"keyword": "Bool",
					}},
				{in: `[ qw/a b c/ qw{ d e f} qw(g i j ) qw<h k l> qw[ m n ogo ]]`, out: []interface{}{"a", "b", "c", "d", "e", "f", "g", "i", "j", "h", "k", "l", "m", "n", "ogo"}},
				{in: `{{$a}}`, out: bw.ValPath{{Type: bw.ValPathItemVar, Key: "a"}}},
				{in: `{ some {{ $a }} }`, out: map[string]interface{}{
					"some": bw.ValPath{{Type: bw.ValPathItemVar, Key: "a"}},
				}},
				{in: `{ some $a.thing }`, out: map[string]interface{}{
					"some": bw.ValPath{
						{Type: bw.ValPathItemVar, Key: "a"},
						{Type: bw.ValPathItemKey, Key: "thing"},
					},
				}},
				{in: `{ some: {} }`, out: map[string]interface{}{
					"some": map[string]interface{}{},
				}},
				{in: `{ some: [] }`, out: map[string]interface{}{
					"some": []interface{}{},
				}},
				{in: `{ some: /* comment */ [] }`, out: map[string]interface{}{
					"some": []interface{}{},
				}},

				{in: `{ some: // comment
					[] }`, out: map[string]interface{}{
					"some": []interface{}{},
				}},
				{in: `{ some: <> }`, out: map[string]interface{}{
					"some": []string{},
				}},
			} {
				// bwdebug.Print("v.in", v.in)
				tests[v.in] = bwtesting.Case{
					In:  []interface{}{v.in, v.opt},
					Out: []interface{}{v.out},
				}
			}
			for _, v := range []struct {
				in  string
				opt bwparse.Opt
				out string
			}{
				// "":               "unexpected end of string at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\n",
				{in: `"some" "thing"`,
					out: "unexpected char \x1b[96;1m'\"'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m34\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m7\x1b[0m: \x1b[32m\"some\" \x1b[91m\"\x1b[0mthing\"\n",
				},
				// `{ some = > "thing" }`: "expects \x1b[97;1mArray\x1b[0m or \x1b[97;1mString\x1b[0m or \x1b[97;1mRange\x1b[0m or \x1b[97;1mNumber\x1b[0m or \x1b[97;1mPath\x1b[0m or \x1b[97;1mMap\x1b[0m or \x1b[97;1mArrayOfString\x1b[0m or \x1b[97;1mNil\x1b[0m or \x1b[97;1mBool\x1b[0m or \x1b[97;1mId\x1b[0m instead of unexpected char \x1b[96;1m'='\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m61\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m7\x1b[0m: \x1b[32m{ some \x1b[91m=\x1b[0m > \"thing\" }\n",
				{in: `qw/ one two three`,
					out: "unexpected end of string at pos \x1b[38;5;252;1m17\x1b[0m: \x1b[32mqw/ one two three\n",
				},
				{in: `qw/ one two three `,
					out: "unexpected end of string at pos \x1b[38;5;252;1m18\x1b[0m: \x1b[32mqw/ one two three \n",
				},
				{in: `"one two three `,
					out: "unexpected end of string at pos \x1b[38;5;252;1m15\x1b[0m: \x1b[32m\"one two three \n",
				},
				{in: `-`,
					out: "unexpected end of string at pos \x1b[38;5;252;1m1\x1b[0m: \x1b[32m-\n",
				},
				{in: `"\z"`,
					out: "unexpected char \x1b[96;1m'z'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m122\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m2\x1b[0m: \x1b[32m\"\\\x1b[91mz\x1b[0m\"\n",
				},
				{in: `{key:`,
					out: "unexpected end of string at pos \x1b[38;5;252;1m5\x1b[0m: \x1b[32m{key:\n",
				},
				{in: `qw `,
					out: "expects \x1b[97;1mArray\x1b[0m or \x1b[97;1mString\x1b[0m or \x1b[97;1mRange\x1b[0m or \x1b[97;1mNumber\x1b[0m or \x1b[97;1mPath\x1b[0m or \x1b[97;1mMap\x1b[0m or \x1b[97;1mArrayOfString\x1b[0m or \x1b[97;1mNil\x1b[0m or \x1b[97;1mBool\x1b[0m instead of unexpected char \x1b[96;1m'q'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m113\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\x1b[91mq\x1b[0mw \n",
				},
				{in: `{ key: 1_000_000_000_000_000_000_000_000 }`,
					out: "strconv.ParseUint: parsing \"1000000000000000000000000\": value out of range at pos \x1b[38;5;252;1m7\x1b[0m: \x1b[32m{ key: \x1b[91m1_000_000_000_000_000_000_000_000\x1b[0m }\n",
				},
				{in: `{ type Number keyA valA keyB valB }`,
					opt: bwparse.Opt{IdVals: map[string]interface{}{"Number": "Number"}},
					out: "unexpected \x1b[91;1m\"valA\"\x1b[0m at pos \x1b[38;5;252;1m19\x1b[0m: \x1b[32m{ type Number keyA \x1b[91mvalA\x1b[0m keyB valB }\n"},
				{in: "{ val: nil def: Array",
					opt: bwparse.Opt{IdVals: map[string]interface{}{"Array": "Array"}},
					out: "unexpected end of string at pos \x1b[38;5;252;1m21\x1b[0m: \x1b[32m{ val: nil def: Array\n"},
				{in: `{ some { { $a }} }`,
					out: "unexpected char \x1b[96;1m'{'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m123\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m9\x1b[0m: \x1b[32m{ some { \x1b[91m{\x1b[0m $a }} }\n",
				},
				{in: `{ some {{ $a } } }`,
					out: "unexpected char \x1b[96;1m'}'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m125\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m13\x1b[0m: \x1b[32m{ some {{ $a \x1b[91m}\x1b[0m } }\n",
				},
			} {
				tests[v.in] = bwtesting.Case{
					In:    []interface{}{v.in, v.opt},
					Panic: v.out,
				}
			}
			return tests
		}(),
		// "+2.0",
		// "$idx.3..{{some.thing}}",
		// "..",
		// `{ key => "\"value\n", 'bool': true keyword Bool}`,
		// "0..1",
		// "0.5..1",
		// "..3.14",
		// "..",
		// "..3.14",
		// "0..1",
	)
}

func TestFrom2(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(s string) (result interface{}) {
			var err error
			if result, err = func(s string) (result interface{}, err error) {
				defer func() {
					if err != nil {
						result = nil
					}
				}()
				p := bwparse.From(bwrune.FromString(s))
				var ok bool
				if result, _, ok, err = bwparse.Val(p, bwparse.Opt{IdVals: map[string]interface{}{"Int": "Int"}}); err == nil {
					err = end(p, ok)
				}
				return
			}(s); err != nil {
				bwerr.PanicErr(err)
			}
			return result
		},
		map[string]bwtesting.Case{
			`{ type Float64 }`: {
				In:    []interface{}{func(testName string) string { return testName }},
				Panic: "unexpected \x1b[91;1m\"Float64\"\x1b[0m at pos \x1b[38;5;252;1m7\x1b[0m: \x1b[32m{ type \x1b[91mFloat64\x1b[0m }\n",
			},
		},
	)
}

func TestFrom(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(s string, opt ...map[string]interface{}) {
			_ = bwparse.From(bwrune.FromString(s), opt...)
		},
		map[string]bwtesting.Case{
			`preLineCount non uint`: {
				In:    []interface{}{"", map[string]interface{}{"preLineCount": true}},
				Panic: "\x1b[38;5;201;1mopt.preLineCount\x1b[0m (\x1b[96;1m(bool)true\x1b[0m) is not \x1b[97;1mUint\x1b[0m",
			},
			`postLineCount non uint`: {
				In:    []interface{}{"", map[string]interface{}{"postLineCount": true}},
				Panic: "\x1b[38;5;201;1mopt.postLineCount\x1b[0m (\x1b[96;1m(bool)true\x1b[0m) is not \x1b[97;1mUint\x1b[0m",
			},
			`unexpected keys`: {
				In:    []interface{}{"", map[string]interface{}{"idvals": true}},
				Panic: "\x1b[38;5;201;1mopt\x1b[0m (\x1b[96;1m{\n  \"idvals\": true\n}\x1b[0m) has unexpected keys \x1b[96;1m[\"idvals\"]\x1b[0m",
			},
		},
	)
}

func TestLineCount(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(opt ...map[string]interface{}) (result interface{}) {
			s := `{
				some "thing"
				type Float64
				another "key"
			}`
			p := bwparse.From(bwrune.FromString(s), opt...)
			var err error
			if result, _, _, err = bwparse.Val(p, bwparse.Opt{IdVals: map[string]interface{}{"Int": "Int"}}); err != nil {
				bwerr.PanicErr(err)
			}
			return
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for _, v := range []struct {
				preLineCount  uint
				postLineCount uint
				s             string
			}{
				{0, 0,
					"unexpected \x1b[91;1m\"Float64\"\x1b[0m at line \x1b[38;5;252;1m3\x1b[0m, col \x1b[38;5;252;1m10\x1b[0m (pos \x1b[38;5;252;1m28\x1b[0m)\x1b[0m:\n\x1b[32m\t\t\t\ttype \x1b[91mFloat64\x1b[0m\n",
				},
				{0, 1,
					"unexpected \x1b[91;1m\"Float64\"\x1b[0m at line \x1b[38;5;252;1m3\x1b[0m, col \x1b[38;5;252;1m10\x1b[0m (pos \x1b[38;5;252;1m28\x1b[0m)\x1b[0m:\n\x1b[32m\t\t\t\ttype \x1b[91mFloat64\x1b[0m\n\t\t\t\tanother \"key\"\n",
				},
				{0, 2,
					"unexpected \x1b[91;1m\"Float64\"\x1b[0m at line \x1b[38;5;252;1m3\x1b[0m, col \x1b[38;5;252;1m10\x1b[0m (pos \x1b[38;5;252;1m28\x1b[0m)\x1b[0m:\n\x1b[32m\t\t\t\ttype \x1b[91mFloat64\x1b[0m\n\t\t\t\tanother \"key\"\n\t\t\t}\n",
				},
				{0, 3,
					"unexpected \x1b[91;1m\"Float64\"\x1b[0m at line \x1b[38;5;252;1m3\x1b[0m, col \x1b[38;5;252;1m10\x1b[0m (pos \x1b[38;5;252;1m28\x1b[0m)\x1b[0m:\n\x1b[32m\t\t\t\ttype \x1b[91mFloat64\x1b[0m\n\t\t\t\tanother \"key\"\n\t\t\t}\n",
				},
				{1, 0,
					"unexpected \x1b[91;1m\"Float64\"\x1b[0m at line \x1b[38;5;252;1m3\x1b[0m, col \x1b[38;5;252;1m10\x1b[0m (pos \x1b[38;5;252;1m28\x1b[0m)\x1b[0m:\n\x1b[32m\t\t\t\tsome \"thing\"\n\t\t\t\ttype \x1b[91mFloat64\x1b[0m\n",
				},
				{1, 1,
					"unexpected \x1b[91;1m\"Float64\"\x1b[0m at line \x1b[38;5;252;1m3\x1b[0m, col \x1b[38;5;252;1m10\x1b[0m (pos \x1b[38;5;252;1m28\x1b[0m)\x1b[0m:\n\x1b[32m\t\t\t\tsome \"thing\"\n\t\t\t\ttype \x1b[91mFloat64\x1b[0m\n\t\t\t\tanother \"key\"\n",
				},
				{1, 2,
					"unexpected \x1b[91;1m\"Float64\"\x1b[0m at line \x1b[38;5;252;1m3\x1b[0m, col \x1b[38;5;252;1m10\x1b[0m (pos \x1b[38;5;252;1m28\x1b[0m)\x1b[0m:\n\x1b[32m\t\t\t\tsome \"thing\"\n\t\t\t\ttype \x1b[91mFloat64\x1b[0m\n\t\t\t\tanother \"key\"\n\t\t\t}\n",
				},
				{1, 3,
					"unexpected \x1b[91;1m\"Float64\"\x1b[0m at line \x1b[38;5;252;1m3\x1b[0m, col \x1b[38;5;252;1m10\x1b[0m (pos \x1b[38;5;252;1m28\x1b[0m)\x1b[0m:\n\x1b[32m\t\t\t\tsome \"thing\"\n\t\t\t\ttype \x1b[91mFloat64\x1b[0m\n\t\t\t\tanother \"key\"\n\t\t\t}\n",
				},
				{2, 0,
					"unexpected \x1b[91;1m\"Float64\"\x1b[0m at line \x1b[38;5;252;1m3\x1b[0m, col \x1b[38;5;252;1m10\x1b[0m (pos \x1b[38;5;252;1m28\x1b[0m)\x1b[0m:\n\x1b[32m{\n\t\t\t\tsome \"thing\"\n\t\t\t\ttype \x1b[91mFloat64\x1b[0m\n",
				},
				{2, 1,
					"unexpected \x1b[91;1m\"Float64\"\x1b[0m at line \x1b[38;5;252;1m3\x1b[0m, col \x1b[38;5;252;1m10\x1b[0m (pos \x1b[38;5;252;1m28\x1b[0m)\x1b[0m:\n\x1b[32m{\n\t\t\t\tsome \"thing\"\n\t\t\t\ttype \x1b[91mFloat64\x1b[0m\n\t\t\t\tanother \"key\"\n",
				},
				{2, 2,
					"unexpected \x1b[91;1m\"Float64\"\x1b[0m at line \x1b[38;5;252;1m3\x1b[0m, col \x1b[38;5;252;1m10\x1b[0m (pos \x1b[38;5;252;1m28\x1b[0m)\x1b[0m:\n\x1b[32m{\n\t\t\t\tsome \"thing\"\n\t\t\t\ttype \x1b[91mFloat64\x1b[0m\n\t\t\t\tanother \"key\"\n\t\t\t}\n",
				},
				{2, 3,
					"unexpected \x1b[91;1m\"Float64\"\x1b[0m at line \x1b[38;5;252;1m3\x1b[0m, col \x1b[38;5;252;1m10\x1b[0m (pos \x1b[38;5;252;1m28\x1b[0m)\x1b[0m:\n\x1b[32m{\n\t\t\t\tsome \"thing\"\n\t\t\t\ttype \x1b[91mFloat64\x1b[0m\n\t\t\t\tanother \"key\"\n\t\t\t}\n",
				},
				{3, 0,
					"unexpected \x1b[91;1m\"Float64\"\x1b[0m at line \x1b[38;5;252;1m3\x1b[0m, col \x1b[38;5;252;1m10\x1b[0m (pos \x1b[38;5;252;1m28\x1b[0m)\x1b[0m:\n\x1b[32m{\n\t\t\t\tsome \"thing\"\n\t\t\t\ttype \x1b[91mFloat64\x1b[0m\n",
				},
				{3, 1,
					"unexpected \x1b[91;1m\"Float64\"\x1b[0m at line \x1b[38;5;252;1m3\x1b[0m, col \x1b[38;5;252;1m10\x1b[0m (pos \x1b[38;5;252;1m28\x1b[0m)\x1b[0m:\n\x1b[32m{\n\t\t\t\tsome \"thing\"\n\t\t\t\ttype \x1b[91mFloat64\x1b[0m\n\t\t\t\tanother \"key\"\n",
				},
				{3, 2,
					"unexpected \x1b[91;1m\"Float64\"\x1b[0m at line \x1b[38;5;252;1m3\x1b[0m, col \x1b[38;5;252;1m10\x1b[0m (pos \x1b[38;5;252;1m28\x1b[0m)\x1b[0m:\n\x1b[32m{\n\t\t\t\tsome \"thing\"\n\t\t\t\ttype \x1b[91mFloat64\x1b[0m\n\t\t\t\tanother \"key\"\n\t\t\t}\n",
				},
				{3, 3,
					"unexpected \x1b[91;1m\"Float64\"\x1b[0m at line \x1b[38;5;252;1m3\x1b[0m, col \x1b[38;5;252;1m10\x1b[0m (pos \x1b[38;5;252;1m28\x1b[0m)\x1b[0m:\n\x1b[32m{\n\t\t\t\tsome \"thing\"\n\t\t\t\ttype \x1b[91mFloat64\x1b[0m\n\t\t\t\tanother \"key\"\n\t\t\t}\n",
				},
			} {
				tests[fmt.Sprintf(`"preLineCount": %d, "postLineCount": %d`, v.preLineCount, v.postLineCount)] = bwtesting.Case{
					In: []interface{}{
						map[string]interface{}{"preLineCount": v.preLineCount, "postLineCount": v.postLineCount},
					},
					Panic: v.s,
				}
			}
			return tests
		}(),
	)
}
