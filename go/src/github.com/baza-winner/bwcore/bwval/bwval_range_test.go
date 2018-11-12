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
		bwval.RangeMinMax: bwval.IntRange{MinPtr: bwval.PtrToInt(-1), MaxPtr: bwval.PtrToInt(2)},
		bwval.RangeMin:    bwval.NumberRange{MinPtr: bwval.PtrToNumber(-1)},
		bwval.RangeMax:    bwval.IntRange{MaxPtr: bwval.PtrToInt(2)},
		bwval.RangeNo:     bwval.NumberRange{},
	}
	tests := map[string]bwtesting.Case{}
	for k, v := range example {
		tests[k.String()] = bwtesting.Case{
			In:  []interface{}{v},
			Out: []interface{}{k},
		}
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, bwval.RangeKind, tests)
}

func TestRangeString(t *testing.T) {
	example := map[bwval.RangeKindValue]struct {
		r bwval.Range
		s string
	}{
		bwval.RangeMinMax: {
			r: bwval.IntRange{MinPtr: bwval.PtrToInt(-1), MaxPtr: bwval.PtrToInt(2)},
			s: "-1..2",
		},
		bwval.RangeMin: {
			r: bwval.NumberRange{MinPtr: bwval.PtrToNumber(-1)},
			s: "-1..",
		},
		bwval.RangeMax: {
			r: bwval.IntRange{MaxPtr: bwval.PtrToInt(2)},
			s: "..2",
		},
		bwval.RangeNo: {
			r: bwval.NumberRange{},
			s: "..",
		},
	}
	tests := map[string]bwtesting.Case{}
	for _, v := range example {
		tests[v.s] = bwtesting.Case{
			In:  []interface{}{v.r},
			Out: []interface{}{v.s},
		}
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, bwval.RangeString, tests)
}

func TestNumberRangeMarshalJSON(t *testing.T) {
	example := map[bwval.RangeKindValue]struct {
		r bwval.NumberRange
		s string
	}{
		bwval.RangeMinMax: {
			r: bwval.NumberRange{MinPtr: bwval.PtrToNumber(-1), MaxPtr: bwval.PtrToNumber(2)},
			s: "-1..2",
		},
		bwval.RangeMin: {
			r: bwval.NumberRange{MinPtr: bwval.PtrToNumber(-1)},
			s: "-1..",
		},
		bwval.RangeMax: {
			r: bwval.NumberRange{MaxPtr: bwval.PtrToNumber(2)},
			s: "..2",
		},
		bwval.RangeNo: {
			r: bwval.NumberRange{},
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

func NumberRangePretty(r bwval.NumberRange) string {
	return bwjson.Pretty(r)
}

func TestRangeContains(t *testing.T) {
	example := map[bwval.RangeKindValue]struct {
		// r   bwval.IntRange
		min *int
		max *int
		in  int
		out int
	}{
		bwval.RangeMinMax: {
			min: bwval.PtrToInt(-1),
			max: bwval.PtrToInt(2),
			in:  0,
			out: -2,
		},
		bwval.RangeMin: {
			min: bwval.PtrToInt(-1),
			in:  0,
			out: -2,
		},
		bwval.RangeMax: {
			max: bwval.PtrToInt(2),
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
		intRange := bwval.IntRange{v.min, v.max}
		numRange := bwval.NumberRange{numPtrFromIntPtr(v.min), numPtrFromIntPtr(v.max)}
		s := bwval.RangeString(intRange)
		tests[fmt.Sprintf("%d in %s", v.in, s)] = bwtesting.Case{
			In:  []interface{}{intRange, v.in},
			Out: []interface{}{true},
		}
		tests[fmt.Sprintf("float64(%d) in %s", v.in, s)] = bwtesting.Case{
			In:  []interface{}{numRange, float64(v.in)},
			Out: []interface{}{true},
		}
		if k != bwval.RangeNo {
			tests[fmt.Sprintf("%d out of %s", v.out, s)] = bwtesting.Case{
				In:  []interface{}{intRange, v.out},
				Out: []interface{}{false},
			}
			tests[fmt.Sprintf("float64(%d) out of %s", v.out, s)] = bwtesting.Case{
				In:  []interface{}{numRange, float64(v.out)},
				Out: []interface{}{false},
			}
		}
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, bwval.RangeContains, tests)
}

func numPtrFromIntPtr(intPtr *int) (result *float64) {
	if intPtr == nil {
		return
	}
	n := float64(*intPtr)
	result = &n
	return
}
