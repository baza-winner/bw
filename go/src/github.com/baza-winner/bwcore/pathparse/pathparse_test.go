package pathparse

import (
	"testing"

	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
)

func TestParse(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"0": {
			In: []interface{}{`0`},
			Out: []interface{}{
				[]interface{}{0},
				nil,
			},
		},
		"[0]": {
			In: []interface{}{`0`},
			Out: []interface{}{
				[]interface{}{0},
				nil,
			},
		},
		"#0": {
			In: []interface{}{`0`},
			Out: []interface{}{
				[]interface{}{0},
				nil,
			},
		},
		"'so me'": {
			In: []interface{}{`'so me'`},
			Out: []interface{}{
				[]interface{}{"so me"},
				nil,
			},
		},
		`"so\nme"`: {
			In: []interface{}{`"so\nme"`},
			Out: []interface{}{
				[]interface{}{"so\nme"}, nil,
			},
		},
		`complex`: {
			In: []interface{}{`[131070].keys.'some \'\"\\thing'.#8_589_934_591."\a\b\f\n\r\t\v".510.31`},
			Out: []interface{}{
				[]interface{}{int32(131070), "keys", "some '\"\\thing", int64(8589934591), "\a\b\f\n\r\t\v", int16(510), int8(31)},
				nil,
			},
		},
	}

	testsToRun := tests
	bwmap.CropMap(testsToRun)
	// bwmap.CropMap(testsToRun, "_expectSpaceOrQwItemOrDelimiter && fa.curr.runePtr == EOF")
	// bwmap.CropMap(testsToRun, "qw ")
	bwtesting.BwRunTests(t, Parse, testsToRun)
}
