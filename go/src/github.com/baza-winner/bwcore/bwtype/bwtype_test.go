package bwtype_test

import (
	"fmt"
	"testing"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwtype"
)

func TestInt(t *testing.T) {
	tests := map[string]bwtesting.Case{}
	for k, v := range map[string]struct {
		in  interface{}
		out int
	}{
		"int8":   {in: bw.MaxInt8, out: int(bw.MaxInt8)},
		"int16":  {in: bw.MaxInt16, out: int(bw.MaxInt16)},
		"int32":  {in: bw.MaxInt32, out: int(bw.MaxInt32)},
		"int64":  {in: int64(bw.MaxInt32), out: int(bw.MaxInt32)},
		"int":    {in: bw.MaxInt, out: int(bw.MaxInt)},
		"uint8":  {in: bw.MaxUint8, out: int(bw.MaxUint8)},
		"uint16": {in: bw.MaxUint16, out: int(bw.MaxUint16)},
		"uint32": {in: uint32(bw.MaxInt32), out: int(bw.MaxInt32)},
		"uint64": {in: uint64(bw.MaxInt), out: bw.MaxInt},
		"uint":   {in: uint(bw.MaxInt), out: bw.MaxInt},
	} {
		tests[k+", ok"] = bwtesting.Case{
			In:  []interface{}{v.in},
			Out: []interface{}{v.out, true},
		}
	}
	for k, v := range map[string]struct {
		in interface{}
	}{
		"uint":   {in: uint(bw.MaxUint)},
		"uint64": {in: bw.MaxUint64},
	} {
		tests[k+", !ok"] = bwtesting.Case{
			In:  []interface{}{v.in},
			Out: []interface{}{0, false},
		}
	}

	bwtesting.BwRunTests(t, bwtype.Int, tests,
		nil,
	)
}

func TestFloat64(t *testing.T) {
	tests := map[string]bwtesting.Case{}
	for k, v := range map[string]struct {
		in  interface{}
		out float64
	}{
		"int8":    {in: bw.MaxInt8, out: float64(bw.MaxInt8)},
		"int16":   {in: bw.MaxInt16, out: float64(bw.MaxInt16)},
		"int32":   {in: bw.MaxInt32, out: float64(bw.MaxInt32)},
		"int64":   {in: int64(bw.MaxInt64), out: float64(bw.MaxInt64)},
		"int":     {in: bw.MaxInt, out: float64(bw.MaxInt)},
		"uint8":   {in: bw.MaxUint8, out: float64(bw.MaxUint8)},
		"uint16":  {in: bw.MaxUint16, out: float64(bw.MaxUint16)},
		"uint32":  {in: bw.MaxUint32, out: float64(bw.MaxUint32)},
		"uint64":  {in: bw.MaxUint64, out: float64(bw.MaxUint64)},
		"uint":    {in: bw.MaxUint, out: float64(bw.MaxUint)},
		"float32": {in: float32(0), out: float64(0)},
		"float64": {in: float64(0), out: float64(0)},
	} {
		tests[k+", ok"] = bwtesting.Case{
			In:  []interface{}{v.in},
			Out: []interface{}{v.out, true},
		}
	}
	for k, v := range map[string]struct {
		in interface{}
	}{
		"bool": {in: true},
	} {
		tests[k+", !ok"] = bwtesting.Case{
			In:  []interface{}{v.in},
			Out: []interface{}{float64(0), false},
		}
	}

	bwtesting.BwRunTests(t, bwtype.Float64, tests,
		nil,
	)
}

func TestMustNumberFrom(t *testing.T) {

	bwtesting.BwRunTests(t,
		func(val interface{}) int {
			return bwtype.MustNumberFrom(val).MustInt()
		}, map[string]bwtesting.Case{
			"nil": {
				In:    []interface{}{nil},
				Panic: "\x1b[96;1m(interface {})<nil>\x1b[0m is not \x1b[97;1mInt\x1b[0m",
			},
			"273": {
				In:  []interface{}{273},
				Out: []interface{}{273},
			},
			"3.14": {
				In:    []interface{}{3.14},
				Panic: "\x1b[96;1m(float64)3.14\x1b[0m is not \x1b[97;1mInt\x1b[0m",
			},
			"Number(-273)": {
				In:  []interface{}{bwtype.MustNumberFrom(-273)},
				Out: []interface{}{-273},
			},
			"true": {
				In:    []interface{}{true},
				Panic: "\x1b[96;1m(bool)true\x1b[0m can not be a \x1b[97;1mNumber\x1b[0m",
			},
		},
	)

	bwtesting.BwRunTests(t,
		func(val interface{}) bool {
			return bwtype.MustNumberFrom(val).IsInt()
		}, map[string]bwtesting.Case{
			"nil": {
				In:  []interface{}{nil},
				Out: []interface{}{false},
			},
			"273": {
				In:  []interface{}{273},
				Out: []interface{}{true},
			},
			"3.14": {
				In:  []interface{}{3.14},
				Out: []interface{}{false},
			},
			"Number(-273)": {
				In:  []interface{}{bwtype.MustNumberFrom(-273)},
				Out: []interface{}{true},
			},
			"Number(-2.71)": {
				In:  []interface{}{bwtype.MustNumberFrom(-2.71)},
				Out: []interface{}{false},
			},
		},
	)

	bwtesting.BwRunTests(t,
		func(val interface{}) float64 {
			return bwtype.MustNumberFrom(val).MustFloat64()
		}, map[string]bwtesting.Case{
			"nil": {
				In:    []interface{}{nil},
				Panic: "\x1b[96;1m(interface {})<nil>\x1b[0m is not \x1b[97;1mFloat64\x1b[0m",
			},
		},
	)
}

func TestMustRangeFrom(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(a bwtype.A) (string, bwtype.Number, bwtype.Number) {
			r := bwtype.MustRangeFrom(a)
			return r.String(), r.Min(), r.Max()
		}, map[string]bwtesting.Case{
			"..": {
				In: []interface{}{bwtype.A{}},
				Out: []interface{}{func(testName string) string { return testName },
					bwtype.Number{},
					bwtype.Number{},
				},
			},
			"2..": {
				In: []interface{}{bwtype.A{Min: 2}},
				Out: []interface{}{func(testName string) string { return testName },
					bwtype.MustNumberFrom(2),
					bwtype.Number{},
				},
			},
			"..-3.14": {
				In: []interface{}{bwtype.A{Max: -3.14}},
				Out: []interface{}{func(testName string) string { return testName },
					bwtype.Number{},
					bwtype.MustNumberFrom(-3.14),
				},
			},
			"2.71..273": {
				In: []interface{}{bwtype.A{Min: 2.71, Max: 273}},
				Out: []interface{}{func(testName string) string { return testName },
					bwtype.MustNumberFrom(2.71),
					bwtype.MustNumberFrom(273),
				},
			},
			"2.71 > -273": {
				In:    []interface{}{bwtype.A{Min: 2.71, Max: -273}},
				Panic: "\x1b[38;5;201;1ma.Max\x1b[0m (\x1b[96;1m-273\x1b[0m) must not be \x1b[91;1mless\x1b[0m then \x1b[38;5;201;1ma.Min\x1b[0m (\x1b[96;1m2.71\x1b[0m)\x1b[0m",
			},
			"Min: true": {
				In:    []interface{}{bwtype.A{Min: true}},
				Panic: "\x1b[38;5;201;1ma.Min\x1b[0m (\x1b[96;1m(bool)true\x1b[0m) can not be a \x1b[97;1mNumber\x1b[0m",
			},
			"Max: true": {
				In:    []interface{}{bwtype.A{Max: true}},
				Panic: "\x1b[38;5;201;1ma.Max\x1b[0m (\x1b[96;1m(bool)true\x1b[0m) can not be a \x1b[97;1mNumber\x1b[0m",
			},
		},
	)
}

func TestRangeContains(t *testing.T) {
	bwtesting.BwRunTests(t,
		"Contains",
		map[string]bwtesting.Case{
			"nil not in ..": {
				V:   bwtype.MustRangeFrom(bwtype.A{}),
				In:  []interface{}{bwtype.MustNumberFrom(nil)},
				Out: []interface{}{false},
			},
			"-273 in ..": {
				V:   bwtype.MustRangeFrom(bwtype.A{}),
				In:  []interface{}{bwtype.MustNumberFrom(-273)},
				Out: []interface{}{true},
			},
			"-273 in ..0": {
				V:   bwtype.MustRangeFrom(bwtype.A{Max: 0}),
				In:  []interface{}{bwtype.MustNumberFrom(-273)},
				Out: []interface{}{true},
			},
			"-273 not in 0..": {
				V:   bwtype.MustRangeFrom(bwtype.A{Min: 0}),
				In:  []interface{}{bwtype.MustNumberFrom(-273)},
				Out: []interface{}{false},
			},
			"2.71 in 0..3.14": {
				V:   bwtype.MustRangeFrom(bwtype.A{Min: 0, Max: 3.14}),
				In:  []interface{}{bwtype.MustNumberFrom(2.71)},
				Out: []interface{}{true},
			},
			"3.14 not in 0...2.71": {
				V:   bwtype.MustRangeFrom(bwtype.A{Min: 0, Max: 2.71}),
				In:  []interface{}{bwtype.MustNumberFrom(3.14)},
				Out: []interface{}{false},
			},
		},
	)
}

func TestIsEqualTo(t *testing.T) {
	bwtesting.BwRunTests(t,
		"IsEqualTo",
		map[string]bwtesting.Case{
			"-273 == -273": {
				V:   bwtype.MustNumberFrom(-273),
				In:  []interface{}{bwtype.MustNumberFrom(bwtype.MustNumberFrom(-273))},
				Out: []interface{}{true},
			},
			"-273 == -273.0": {
				V:   bwtype.MustNumberFrom(-273),
				In:  []interface{}{bwtype.MustNumberFrom(bwtype.MustNumberFrom(-273.0))},
				Out: []interface{}{true},
			},
			"3.14 != 2.71": {
				V:   bwtype.MustNumberFrom(3.14),
				In:  []interface{}{bwtype.MustNumberFrom(2.71)},
				Out: []interface{}{false},
			},
		},
	)
}

func TestRangeKind(t *testing.T) {
	tests := map[string]bwtesting.Case{}
	for k, v := range map[bwtype.RangeKindValue]bwtype.Range{
		bwtype.RangeMinMax: bwtype.MustRangeFrom(bwtype.A{Min: -1, Max: 2}),
		bwtype.RangeMin:    bwtype.MustRangeFrom(bwtype.A{Min: -1}),
		bwtype.RangeMax:    bwtype.MustRangeFrom(bwtype.A{Max: 2}),
		bwtype.RangeNo:     bwtype.Range{},
	} {
		tests[k.String()] = bwtesting.Case{
			V:   v,
			Out: []interface{}{k},
		}
	}

	bwtesting.BwRunTests(t, "Kind", tests,
		nil,
	)
	bwtesting.BwRunTests(t, "String", map[string]bwtesting.Case{
		"(bwtype.RangeMinMax + 1).String()": {
			V:   bwtype.RangeMinMax + 1,
			Out: []interface{}{"RangeKindValue(4)"},
		},
	})
}

func TestRangeString(t *testing.T) {
	tests := map[string]bwtesting.Case{}
	for _, v := range map[bwtype.RangeKindValue]struct {
		r bwtype.Range
		s string
	}{
		bwtype.RangeMinMax: {
			r: bwtype.MustRangeFrom(bwtype.A{Min: -1, Max: 2}),
			s: "-1..2",
		},
		bwtype.RangeMin: {
			r: bwtype.MustRangeFrom(bwtype.A{Min: -1}),
			s: "-1..",
		},
		bwtype.RangeMax: {
			r: bwtype.MustRangeFrom(bwtype.A{Max: 2}),
			s: "..2",
		},
		bwtype.RangeNo: {
			r: bwtype.Range{},
			s: "..",
		},
	} {
		tests[v.s] = bwtesting.Case{
			V:   v.r,
			Out: []interface{}{v.s},
		}
	}

	bwtesting.BwRunTests(t, "String", tests,
		nil,
	)
}

func TestNumberRangeMarshalJSON(t *testing.T) {
	example := map[bwtype.RangeKindValue]struct {
		r bwtype.Range
		s string
	}{
		bwtype.RangeMinMax: {
			r: bwtype.MustRangeFrom(bwtype.A{Min: -1, Max: 2}),
			s: "-1..2",
		},
		bwtype.RangeMin: {
			r: bwtype.MustRangeFrom(bwtype.A{Min: -1}),
			s: "-1..",
		},
		bwtype.RangeMax: {
			r: bwtype.MustRangeFrom(bwtype.A{Max: 2}),
			s: "..2",
		},
		bwtype.RangeNo: {
			r: bwtype.Range{},
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

	bwtesting.BwRunTests(t,
		func(r bwtype.Range) string {
			return bwjson.Pretty(r)
		}, tests,
		nil,
	)
}
