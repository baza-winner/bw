package bwparse_test

import (
	"fmt"
	"testing"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwparse"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwtesting"
)

func TestUnexpected(t *testing.T) {
	testUnexpectedHelper := func(s string, ofs uint) *bwparse.P {
		p := bwparse.From(bwrune.ProviderFromString(s))
		p.Forward(ofs)
		return p
	}
	bwtesting.BwRunTests(t,
		"Unexpected",
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			tests["panic"] = bwtesting.Case{
				V:     testUnexpectedHelper("some", 2),
				In:    []interface{}{bwparse.PosInfo{Pos: 4}},
				Panic: "\x1b[38;5;201;1mps.Pos\x1b[0m (\x1b[96;1m4\x1b[0m) > \x1b[38;5;201;1mp.Curr.Pos\x1b[0m (\x1b[96;1m1\x1b[0m)\x1b[0m",
			}
			p := testUnexpectedHelper("{\n key wrong \n} ", 0)
			p.Forward(7)
			pi := p.Curr
			p.Forward(5)
			tests["normal"] = bwtesting.Case{
				V:   p,
				In:  []interface{}{pi},
				Out: []interface{}{"unexpected \x1b[91;1m\"wrong\"\x1b[0m at line \x1b[38;5;252;1m2\x1b[0m, col \x1b[38;5;252;1m6\x1b[0m (pos \x1b[38;5;252;1m7\x1b[0m)\x1b[0m:\n\x1b[32m{\n key \x1b[91mwrong\x1b[0m \n} \n"},
			}
			return tests
		}(),
	)
}

func TestLookAhead(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(p *bwparse.P, i int) rune {
			return p.LookAhead(uint(i)).Rune
		},
		func() map[string]bwtesting.Case {
			s := "s\no\nm\ne\nt\nhing"
			p := bwparse.From(bwrune.ProviderFromString(s))
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
				pco := bwparse.PathA{}
				if len(optBases) > 0 {
					pco.Bases = optBases[0]
				}
				p := bwparse.From(bwrune.ProviderFromString(s))
				if result, err = p.PathContent(pco); err == nil {
					err = end(p, true)
				}
				return
			}(s, optBases...); err != nil {
				bwerr.PanicA(bwerr.Err(err))
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
				p := bwparse.From(bwrune.ProviderFromString(s))
				var ok bool
				if result, _, ok, err = p.Int(); err == nil {
					err = end(p, ok)
				}
				return
			}(s); err != nil {
				bwerr.PanicA(bwerr.Err(err))
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
				"+1_000_000_000_000_000_000_000_000": "strconv.ParseInt: parsing \"+1000000000000000000000000\": value out of range at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\x1b[91m+1_000_000_000_000_000_000_000_000\x1b[0m\n",
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
		err = p.Unexpected(p.Curr)
	} else {
		err = p.SkipSpace(bwparse.TillEOF)
	}
	return
}

func TestVal(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(s string, optVars ...map[string]interface{}) (result interface{}) {
			var err error
			if result, err = func(s string) (result interface{}, err error) {
				defer func() {
					if err != nil {
						result = nil
					}
				}()
				p := bwparse.From(bwrune.ProviderFromString(s), map[string]interface{}{
					"idVals": map[string]interface{}{
						"Bool":    "Bool",
						"String":  "String",
						"Int":     "Int",
						"Float64": "Float64",
						"Array":   "Array",
						"ArrayOf": "ArrayOf",
					}})
				var ok bool
				if result, _, ok, err = p.Val(); err == nil {
					err = end(p, ok)
				}
				return
			}(s); err != nil {
				bwerr.PanicA(bwerr.Err(err))
			}
			return result
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for k, v := range map[string]interface{}{
				"nil":               nil,
				"true":              true,
				"false":             false,
				"0":                 0,
				"-1_000_000":        -1000000,
				"+3.14":             3.14,
				"+2.0":              2,
				"[0, 1]":            []interface{}{0, 1},
				`"a"`:               "a",
				`<a b c>`:           []string{"a", "b", "c"},
				`[<a b c>]`:         []interface{}{"a", "b", "c"},
				`["x" <a b c> "z"]`: []interface{}{"x", "a", "b", "c", "z"},
				`{ key "value" bool true }`: map[string]interface{}{
					"key":  "value",
					"bool": true,
				},
				`{ key => "\"value\n", 'bool': true keyword Bool}`: map[string]interface{}{
					"key":     "\"value\n",
					"bool":    true,
					"keyword": "Bool",
				},
				`[ qw/a b c/ qw{ d e f} qw(g i j ) qw<h k l> qw[ m n ogo ]]`: []interface{}{"a", "b", "c", "d", "e", "f", "g", "i", "j", "h", "k", "l", "m", "n", "ogo"},
				`{{$a}}`: bw.ValPath{{Type: bw.ValPathItemVar, Key: "a"}},
				`{ some {{ $a }} }`: map[string]interface{}{
					"some": bw.ValPath{{Type: bw.ValPathItemVar, Key: "a"}},
				},
				`{ some $a.thing }`: map[string]interface{}{
					"some": bw.ValPath{
						{Type: bw.ValPathItemVar, Key: "a"},
						{Type: bw.ValPathItemKey, Key: "thing"},
					},
				},
				`{ some: {} }`: map[string]interface{}{
					"some": map[string]interface{}{},
				},
				`{ some: [] }`: map[string]interface{}{
					"some": []interface{}{},
				},
				`{ some: /* comment */ [] }`: map[string]interface{}{
					"some": []interface{}{},
				},

				`{ some: // comment
					[] }`: map[string]interface{}{
					"some": []interface{}{},
				},
				`{ some: <> }`: map[string]interface{}{
					"some": []string{},
				},
			} {
				tests[k] = bwtesting.Case{
					In:  []interface{}{func(testName string) string { return testName }},
					Out: []interface{}{v},
				}
			}
			for k, v := range map[string]string{
				"":                     "unexpected end of string at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\n",
				`"some" "thing"`:       "unexpected char \x1b[96;1m'\"'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m34\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m7\x1b[0m: \x1b[32m\"some\" \x1b[91m\"\x1b[0mthing\"\n",
				`{ some = > "thing" }`: "unexpected char \x1b[96;1m'='\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m61\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m7\x1b[0m: \x1b[32m{ some \x1b[91m=\x1b[0m > \"thing\" }\n",
				`qw/ one two three`:    "unexpected end of string at pos \x1b[38;5;252;1m17\x1b[0m: \x1b[32mqw/ one two three\n",
				`qw/ one two three `:   "unexpected end of string at pos \x1b[38;5;252;1m18\x1b[0m: \x1b[32mqw/ one two three \n",
				`"one two three `:      "unexpected end of string at pos \x1b[38;5;252;1m15\x1b[0m: \x1b[32m\"one two three \n",
				`-`:                    "unexpected end of string at pos \x1b[38;5;252;1m1\x1b[0m: \x1b[32m-\n",
				`"\z"`:                 "unexpected char \x1b[96;1m'z'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m122\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m2\x1b[0m: \x1b[32m\"\\\x1b[91mz\x1b[0m\"\n",
				`{key:`:                "unexpected end of string at pos \x1b[38;5;252;1m5\x1b[0m: \x1b[32m{key:\n",
				`}`:                    "unexpected char \x1b[96;1m'}'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m125\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\x1b[91m}\x1b[0m\n",
				`qw `:                  "unexpected \x1b[91;1m\"qw\"\x1b[0m at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\x1b[91mqw\x1b[0m \n",
				`{ key: 1_000_000_000_000_000_000_000_000 }`: "strconv.ParseInt: parsing \"1000000000000000000000000\": value out of range at pos \x1b[38;5;252;1m7\x1b[0m: \x1b[32m{ key: \x1b[91m1_000_000_000_000_000_000_000_000\x1b[0m }\n",
				`{ type Float64 keyA valA keyB valB }`:       "unexpected \x1b[91;1m\"valA\"\x1b[0m at pos \x1b[38;5;252;1m20\x1b[0m: \x1b[32m{ type Float64 keyA \x1b[91mvalA\x1b[0m keyB valB }\n",
				"{ val: nil def: Array":                      "unexpected end of string at pos \x1b[38;5;252;1m21\x1b[0m: \x1b[32m{ val: nil def: Array\n",
				`{ some { { $a }} }`:                         "unexpected char \x1b[96;1m'{'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m123\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m9\x1b[0m: \x1b[32m{ some { \x1b[91m{\x1b[0m $a }} }\n",
				`{ some {{ $a } } }`:                         "unexpected char \x1b[96;1m'}'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m125\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m13\x1b[0m: \x1b[32m{ some {{ $a \x1b[91m}\x1b[0m } }\n",
			} {
				tests[k] = bwtesting.Case{
					In:    []interface{}{func(testName string) string { return testName }},
					Panic: v,
				}
			}
			return tests
		}(),
	)

	bwtesting.BwRunTests(t,
		func(s string) (result interface{}) {
			var err error
			if result, err = func(s string) (result interface{}, err error) {
				defer func() {
					if err != nil {
						result = nil
					}
				}()
				p := bwparse.From(bwrune.ProviderFromString(s))
				var ok bool
				if result, _, ok, err = p.Val(); err == nil {
					err = end(p, ok)
				}
				return
			}(s); err != nil {
				bwerr.PanicA(bwerr.Err(err))
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
			_ = bwparse.From(bwrune.ProviderFromString(s), opt...)
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
			`idVals non map[string]interface{}`: {
				In:    []interface{}{"", map[string]interface{}{"idVals": true}},
				Panic: "\x1b[38;5;201;1mopt.idVals\x1b[0m (\x1b[96;1m(bool)true\x1b[0m) is not \x1b[97;1mmap[string]interface{}\x1b[0m",
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
			p := bwparse.From(bwrune.ProviderFromString(s), opt...)
			var err error
			if result, _, _, err = p.Val(); err != nil {
				bwerr.PanicA(bwerr.Err(err))
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
