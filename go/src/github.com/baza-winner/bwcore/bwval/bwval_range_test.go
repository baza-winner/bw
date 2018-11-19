package bwval_test

import (
	"fmt"
	"testing"

	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwval"
)

func TestRangeKind(t *testing.T) {
	example := map[bwval.RangeKindValue]bwval.Range{
		bwval.RangeMinMax: bwval.Range{Min: bwval.NumberFromInt(-1), Max: bwval.NumberFromInt(2)},
		// bwval.RangeMinMax: bwval.Range{Min: bwval.NumberFromInt(-1), Max: bwval.NumberFromInt(2)},
		// bwval.RangeMin: bwval.Range{Min: bwval.NumberFromFloat64(-1)},
		bwval.RangeMin: bwval.Range{Min: bwval.NumberFromInt(-1)},
		// bwval.RangeMax: bwval.Range{Max: bwval.NumberFromInt(2)},
		bwval.RangeMax: bwval.Range{Max: bwval.NumberFromInt(2)},
		// bwval.RangeNo:  bwval.Range{},
		bwval.RangeNo: bwval.Range{},
	}
	tests := map[string]bwtesting.Case{}
	for k, v := range example {
		tests[k.String()] = bwtesting.Case{
			V: v,
			// In:  []interface{}{v},
			Out: []interface{}{k},
		}
	}

	bwtesting.BwRunTests(t, "Kind", tests,
		nil,
	)
}

func TestRangeString(t *testing.T) {
	example := map[bwval.RangeKindValue]struct {
		r bwval.Range
		s string
	}{
		bwval.RangeMinMax: {
			// r: bwval.Range{Min: bwval.NumberFromInt(-1), Max: bwval.NumberFromInt(2)},
			r: bwval.Range{Min: bwval.NumberFromInt(-1), Max: bwval.NumberFromInt(2)},
			s: "-1..2",
		},
		bwval.RangeMin: {
			r: bwval.Range{Min: bwval.NumberFromFloat64(-1)},
			s: "-1..",
		},
		bwval.RangeMax: {
			r: bwval.Range{Max: bwval.NumberFromInt(2)},
			s: "..2",
		},
		bwval.RangeNo: {
			r: bwval.Range{},
			s: "..",
		},
	}
	tests := map[string]bwtesting.Case{}
	for _, v := range example {
		tests[v.s] = bwtesting.Case{
			V: v.r,
			// In:  []interface{}{v.r},
			Out: []interface{}{v.s},
		}
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, "String", tests)
}

func TestNumberRangeMarshalJSON(t *testing.T) {
	example := map[bwval.RangeKindValue]struct {
		r bwval.Range
		s string
	}{
		bwval.RangeMinMax: {
			r: bwval.Range{Min: bwval.NumberFromFloat64(-1), Max: bwval.NumberFromFloat64(2)},
			s: "-1..2",
		},
		bwval.RangeMin: {
			r: bwval.Range{Min: bwval.NumberFromFloat64(-1)},
			s: "-1..",
		},
		bwval.RangeMax: {
			r: bwval.Range{Max: bwval.NumberFromFloat64(2)},
			s: "..2",
		},
		bwval.RangeNo: {
			r: bwval.Range{},
			s: "..",
		},
	}
	tests := map[string]bwtesting.Case{}
	for _, v := range example {
		tests[v.s] = bwtesting.Case{
			In:  []interface{}{v.r},
			Out: []interface{}{fmt.Sprintf("%q", v.s)},
		}
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, NumberRangePretty, tests)
}

func NumberRangePretty(r bwval.Range) string {
	return bwjson.Pretty(r)
}

func TestRangeContains(t *testing.T) {
	example := map[bwval.RangeKindValue]struct {
		// r   bwval.Range
		min, max bwval.Number
		in       int
		out      int
	}{
		bwval.RangeMinMax: {
			min: bwval.NumberFromInt(-1),
			max: bwval.NumberFromInt(2),
			in:  0,
			out: -2,
		},
		bwval.RangeMin: {
			min: bwval.NumberFromInt(-1),
			in:  0,
			out: -2,
		},
		bwval.RangeMax: {
			max: bwval.NumberFromInt(2),
			in:  0,
			out: 3,
		},
		bwval.RangeNo: {
			in:  0,
			out: 0,
		},
	}
	tests := map[string]bwtesting.Case{}
	for k, v := range example {
		r := bwval.Range{v.min, v.max}
		// intRange := bwval.Range{v.min, v.max}
		// numRange := bwval.Range{numPtrFromIntPtr(v.min), numPtrFromIntPtr(v.max)}
		s := r.String()
		tests[fmt.Sprintf("%d in %s", v.in, s)] = bwtesting.Case{
			V:   r,
			In:  []interface{}{v.in},
			Out: []interface{}{true},
		}
		tests[fmt.Sprintf("float64(%d) in %s", v.in, s)] = bwtesting.Case{
			V:   r,
			In:  []interface{}{float64(v.in)},
			Out: []interface{}{true},
		}
		if k != bwval.RangeNo {
			tests[fmt.Sprintf("%d out of %s", v.out, s)] = bwtesting.Case{
				V:   r,
				In:  []interface{}{v.out},
				Out: []interface{}{false},
			}
			tests[fmt.Sprintf("float64(%d) out of %s", v.out, s)] = bwtesting.Case{
				V:   r,
				In:  []interface{}{float64(v.out)},
				Out: []interface{}{false},
			}
		}
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, "Contains", tests)
}

// func numPtrFromIntPtr(intPtr *int) (result *float64) {
// 	if intPtr == nil {
// 		return
// 	}
// 	n := float64(*intPtr)
// 	result = &n
// 	return
// }
