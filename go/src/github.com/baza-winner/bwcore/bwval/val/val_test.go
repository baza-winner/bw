package val_test

import (
	"testing"

	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwval/val"
)

func TestMustParse(t *testing.T) {
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
		`qw `:                  "unexpected char \x1b[96;1m' '\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m32\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m2\x1b[0m: \x1b[32mqw\x1b[91m \x1b[0m\n",
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
	// bwmap.CropMap(tests, "\"some\" \"thing\"")
	bwtesting.BwRunTests(t, val.MustParse, tests)
}
