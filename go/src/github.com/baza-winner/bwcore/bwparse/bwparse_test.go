package bwparse_test

import (
	"testing"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwparse"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwtesting"
)

func TestPath(t *testing.T) {
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
		"{some.thing}.good": bw.ValPath{
			{Type: bw.ValPathItemPath,
				Path: bw.ValPath{
					{Type: bw.ValPathItemKey, Key: "some"},
					{Type: bw.ValPathItemKey, Key: "thing"},
				},
			},
			{Type: bw.ValPathItemKey, Key: "good"},
		},
		"$some.thing.{good}": bw.ValPath{
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
	} {
		tests[k] = bwtesting.Case{
			In:  []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{v},
		}
	}
	for k, v := range map[string]string{
		"":        "unexpected end of string at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\n",
		"1.":      "unexpected end of string at pos \x1b[38;5;252;1m2\x1b[0m: \x1b[32m1.\n",
		"1.@":     "unexpected char \u001b[96;1m'@'\u001b[0m (\u001b[38;5;201;1mcharCode\u001b[0m: \u001b[96;1m64\u001b[0m)\u001b[0m at pos \u001b[38;5;252;1m2\u001b[0m: \u001b[32m1.\u001b[91m@\u001b[0m\n",
		"-a":      "unexpected char \u001b[96;1m'a'\u001b[0m (\u001b[38;5;201;1mcharCode\u001b[0m: \u001b[96;1m97\u001b[0m)\u001b[0m at pos \u001b[38;5;252;1m1\u001b[0m: \u001b[32m-\u001b[91ma\u001b[0m\n",
		"1a":      "unexpected char \u001b[96;1m'a'\u001b[0m (\u001b[38;5;201;1mcharCode\u001b[0m: \u001b[96;1m97\u001b[0m)\u001b[0m at pos \u001b[38;5;252;1m1\u001b[0m: \u001b[32m1\u001b[91ma\u001b[0m\n",
		"12.#.4":  "unexpected char \x1b[96;1m'.'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m46\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m4\x1b[0m: \x1b[32m12.#\x1b[91m.\x1b[0m4\n",
		"12.{4":   "unexpected end of string at pos \u001b[38;5;252;1m5\u001b[0m: \u001b[32m12.{4\n",
		"12.$a":   "unexpected char \u001b[96;1m'$'\u001b[0m (\u001b[38;5;201;1mcharCode\u001b[0m: \u001b[96;1m36\u001b[0m)\u001b[0m at pos \u001b[38;5;252;1m3\u001b[0m: \u001b[32m12.\u001b[91m$\u001b[0ma\n",
		"$1.some": "unexpected base path idx \x1b[96;1m1\x1b[0m (len(bases): \x1b[96;1m0)\x1b[0m at pos \x1b[38;5;252;1m1\x1b[0m: \x1b[32m$\x1b[91m1\x1b[0m.some\n",
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
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, ".some")
	bwtesting.BwRunTests(t, mustPath, tests)
}

func mustPath(s string, optBases ...[]bw.ValPath) (result bw.ValPath) {
	var err error
	if result, err = func(s string, optBases ...[]bw.ValPath) (result bw.ValPath, err error) {
		defer func() {
			if err != nil {
				result = nil
			}
		}()
		var (
			// isEOF bool
			// r  rune
			ok bool
		)
		p := bwparse.ProviderFrom(bwrune.ProviderFromString(s))

		if err = p.Forward(bwparse.NonEOF); err != nil {
			return
		}
		// p.MustPullRune(bwparse.NonEOF)
		// if r, isEOF, err = p.Rune(); err != nil || isEOF {
		// 	if err == nil {
		// 		err = p.Unexpected(p.Curr)
		// 	}
		// 	return
		// }
		if result, _, ok, err = p.Path(optBases...); err != nil || !ok {
			if err == nil {
				err = p.Unexpected(p.Curr)
			}
			return
		}
		// bwdebug.Print("*p.Curr.RunePtr", string(*p.Curr.RunePtr))
		if err = p.SkipOptionalSpaceTillEOF(); err != nil {
			return
		}
		return
	}(s, optBases...); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return result
}

func TestInt(t *testing.T) {
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
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "1a")
	bwtesting.BwRunTests(t, mustInt, tests)
}

func mustInt(s string, optVars ...map[string]interface{}) (result interface{}) {
	var err error
	if result, err = func(s string, optVars ...map[string]interface{}) (result int, err error) {
		defer func() {
			if err != nil {
				result = 0
			}
		}()
		var (
			// isEOF bool
			// r  rune
			ok bool
		)
		p := bwparse.ProviderFrom(bwrune.ProviderFromString(s))

		if err = p.Forward(bwparse.NonEOF); err != nil {
			return
		}
		// if _, err = p.PullNonEOFRune(); err != nil {
		// 	return
		// }

		// if r, isEOF, err = p.Rune(); err != nil || isEOF {
		// 	if err == nil {
		// 		err = p.Unexpected(p.Curr)
		// 	}
		// 	return
		// }
		if result, _, ok, err = p.Int(); err != nil {
			return
		} else if !ok {
			err = p.Unexpected(p.Curr)
			return
		}
		if err = p.SkipOptionalSpaceTillEOF(); err != nil {
			return
		}
		// if r, isEOF, err = p.PullRuneOrEOF(); err != nil || isEOF {
		// 	return
		// }
		// if _, ok, err = bwparse.ParseSpace(p, r); err != nil {
		// 	return
		// } else if !ok {
		// 	err = p.Unexpected(p.Curr)
		// 	return
		// }
		return
	}(s, optVars...); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return result
}
func TestVal(t *testing.T) {
	tests := map[string]bwtesting.Case{}
	for k, v := range map[string]interface{}{
		"nil":               nil,
		"true":              true,
		"false":             false,
		"0":                 0,
		"-1_000_000":        -1000000,
		"+2.0":              2.0,
		"[0, 1]":            []interface{}{0, 1},
		`"a"`:               "a",
		`<a b c>`:           []interface{}{"a", "b", "c"},
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
	} {
		tests[k] = bwtesting.Case{
			In:  []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{v},
		}
	}
	for k, v := range map[string]string{
		"":                     "unexpected end of string at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\n",
		`"some" "thing"`:       "unexpected char \x1b[96;1m'\"'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m34\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m7\x1b[0m: \x1b[32m\"some\" \x1b[91m\"\x1b[0mthing\"\n",
		`{ some = > "thing" }`: "unexpected char \x1b[96;1m' '\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m32\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m8\x1b[0m: \x1b[32m{ some =\x1b[91m \x1b[0m> \"thing\" }\n",
		`qw/ one two three`:    "unexpected end of string at pos \x1b[38;5;252;1m17\x1b[0m: \x1b[32mqw/ one two three\n",
		`qw/ one two three `:   "unexpected end of string at pos \x1b[38;5;252;1m18\x1b[0m: \x1b[32mqw/ one two three \n",
		`"one two three `:      "unexpected end of string at pos \x1b[38;5;252;1m15\x1b[0m: \x1b[32m\"one two three \n",
		`-`:                    "unexpected end of string at pos \x1b[38;5;252;1m1\x1b[0m: \x1b[32m-\n",
		`"\z"`:                 "unexpected char \x1b[96;1m'z'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m122\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m2\x1b[0m: \x1b[32m\"\\\x1b[91mz\x1b[0m\"\n",
		`{key:`:                "unexpected end of string at pos \x1b[38;5;252;1m5\x1b[0m: \x1b[32m{key:\n",
		`}`:                    "unexpected char \x1b[96;1m'}'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m125\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\x1b[91m}\x1b[0m\n",
		`qw `:                  "unexpected \x1b[91;1m\"qw\"\x1b[0m at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\x1b[91mqw\x1b[0m \n",
		`{ key: 1_000_000_000_000_000_000_000_000 }`: "strconv.ParseInt: parsing \"1000000000000000000000000\": value out of range at pos \x1b[38;5;252;1m7\x1b[0m: \x1b[32m{ key: \x1b[91m1_000_000_000_000_000_000_000_000\x1b[0m }\n",
		`{ type Number keyA valA keyB valB }`:        "unexpected \x1b[91;1m\"valA\"\x1b[0m at pos \x1b[38;5;252;1m19\x1b[0m: \x1b[32m{ type Number keyA \x1b[91mvalA\x1b[0m keyB valB }\n",
		"{ val: nil def: Array":                      "unexpected end of string at pos \x1b[38;5;252;1m21\x1b[0m: \x1b[32m{ val: nil def: Array\n",
	} {
		tests[k] = bwtesting.Case{
			In:    []interface{}{func(testName string) string { return testName }},
			Panic: v,
		}
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "{ key => \"\\\"value\\n\", 'bool': true keyword Bool}")
	bwtesting.BwRunTests(t, mustVal, tests)
}

func mustVal(s string, optVars ...map[string]interface{}) (result interface{}) {
	var err error
	if result, err = func(s string, optVars ...map[string]interface{}) (result interface{}, err error) {
		defer func() {
			if err != nil {
				result = nil
			}
		}()
		var (
			// isEOF bool
			// r  rune
			ok bool
		)
		p := bwparse.ProviderFrom(bwrune.ProviderFromString(s))

		if err = p.Forward(bwparse.NonEOF); err != nil {
			return
		}
		// if r, isEOF, err = p.Rune(); err != nil || isEOF {
		// 	if err == nil {
		// 		err = p.Unexpected(p.Curr)
		// 	}
		// 	return
		// }
		if result, _, ok, err = p.Val(); err != nil || !ok {
			if err == nil {
				err = p.Unexpected(p.Curr)
			}
			return
		}
		if err = p.SkipOptionalSpaceTillEOF(); err != nil {
			return
		}
		return
	}(s, optVars...); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return result
}
