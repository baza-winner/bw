package val_test

import (
	"testing"

	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwval/val"
)

func TestMustParse(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"": {
			In: []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{
				nil,
			},
			Panic: "unexpected end of string at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\n",
		},
		"nil": {
			In: []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{
				nil,
			},
		},
		"true": {
			In: []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{
				true,
			},
		},
		"false": {
			In: []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{
				false,
			},
		},
		"0": {
			In: []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{
				0,
			},
		},
		"-1_000_000": {
			In: []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{
				-1000000,
			},
		},
		"+2.0": {
			In: []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{
				2.0,
			},
		},
		"[0, 1]": {
			In: []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{
				[]interface{}{0, 1},
			},
		},
		`"a"`: {
			In: []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{
				"a",
			},
		},
		`<a b c>`: {
			In: []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{
				[]interface{}{"a", "b", "c"},
			},
		},
		`[<a b c>]`: {
			In: []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{
				[]interface{}{"a", "b", "c"},
			},
		},
		`["x" <a b c> "z"]`: {
			In: []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{
				[]interface{}{"x", "a", "b", "c", "z"},
			},
		},
		`{ key "value" bool true }`: {
			In: []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{
				map[string]interface{}{
					"key":  "value",
					"bool": true,
				},
			},
		},
		`{ key => "\"value\n", 'bool': true keyword Bool}`: {
			In: []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{
				map[string]interface{}{
					"key":     "\"value\n",
					"bool":    true,
					"keyword": "Bool",
				},
			},
		},
		`[ qw/a b c/ qw{ d e f} qw(g i j ) qw<h k l> qw[ m n ogo ]]`: {
			In: []interface{}{func(testName string) string { return testName }},
			Out: []interface{}{
				[]interface{}{"a", "b", "c", "d", "e", "f", "g", "i", "j", "h", "k", "l", "m", "n", "ogo"},
			},
		},
		`"some" "thing"`: {
			In:    []interface{}{func(testName string) string { return testName }},
			Out:   []interface{}{nil},
			Panic: "unexpected char \x1b[96;1m'\"'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m34\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m7\x1b[0m: \x1b[32m\"some\" \x1b[91m\"\x1b[0mthing\"\n",
		},
		`{ some = > "thing" }`: {
			In:    []interface{}{func(testName string) string { return testName }},
			Out:   []interface{}{nil},
			Panic: "unexpected char \x1b[96;1m' '\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m32\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m8\x1b[0m: \x1b[32m{ some =\x1b[91m \x1b[0m> \"thing\" }\n",
		},
		`qw/ one two three`: {
			In:    []interface{}{func(testName string) string { return testName }},
			Out:   []interface{}{nil},
			Panic: "unexpected end of string at pos \x1b[38;5;252;1m17\x1b[0m: \x1b[32mqw/ one two three\n",
		},
		`qw/ one two three `: {
			In:    []interface{}{func(testName string) string { return testName }},
			Out:   []interface{}{nil},
			Panic: "unexpected end of string at pos \x1b[38;5;252;1m18\x1b[0m: \x1b[32mqw/ one two three \n",
		},
		`"one two three `: {
			In:    []interface{}{func(testName string) string { return testName }},
			Out:   []interface{}{nil},
			Panic: "unexpected end of string at pos \x1b[38;5;252;1m15\x1b[0m: \x1b[32m\"one two three \n",
		},
		`-`: {
			In:    []interface{}{func(testName string) string { return testName }},
			Out:   []interface{}{nil},
			Panic: "unexpected end of string at pos \x1b[38;5;252;1m1\x1b[0m: \x1b[32m-\n",
		},
		`"\z"`: {
			In:    []interface{}{func(testName string) string { return testName }},
			Out:   []interface{}{nil},
			Panic: "unexpected char \x1b[96;1m'z'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m122\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m2\x1b[0m: \x1b[32m\"\\\x1b[91mz\x1b[0m\"\n",
		},
		`{key:`: {
			In:    []interface{}{func(testName string) string { return testName }},
			Out:   []interface{}{nil},
			Panic: "unexpected end of string at pos \x1b[38;5;252;1m5\x1b[0m: \x1b[32m{key:\n",
		},
		`}`: {
			In:    []interface{}{func(testName string) string { return testName }},
			Out:   []interface{}{nil},
			Panic: "unexpected char \x1b[96;1m'}'\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m125\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m0\x1b[0m: \x1b[32m\x1b[91m}\x1b[0m\n",
		},
		`qw `: {
			In:    []interface{}{func(testName string) string { return testName }},
			Out:   []interface{}{nil},
			Panic: "unexpected char \x1b[96;1m' '\x1b[0m (\x1b[38;5;201;1mcharCode\x1b[0m: \x1b[96;1m32\x1b[0m)\x1b[0m at pos \x1b[38;5;252;1m2\x1b[0m: \x1b[32mqw\x1b[91m \x1b[0m\n",
		},
		`{ key: 1_000_000_000_000_000_000_000_000 }`: {
			In:    []interface{}{func(testName string) string { return testName }},
			Out:   []interface{}{nil},
			Panic: "strconv.ParseInt: parsing \"1000000000000000000000000\": value out of range at pos \x1b[38;5;252;1m7\x1b[0m: \x1b[32m{ key: \x1b[91m1_000_000_000_000_000_000_000_000\x1b[0m }\n",
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "[ qw/a b c/ qw{ d e f} qw(g i j ) qw<h k l> qw[ m n o ]]")
	bwtesting.BwRunTests(t, val.MustParse, tests)
}
