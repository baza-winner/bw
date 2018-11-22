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
	bwtesting.BwRunTests(t,
		bwtype.Int,
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for k, v := range map[string]struct {
				in  interface{}
				out int
			}{
				"int8":         {in: bw.MaxInt8, out: int(bw.MaxInt8)},
				"int16":        {in: bw.MaxInt16, out: int(bw.MaxInt16)},
				"int32":        {in: bw.MaxInt32, out: int(bw.MaxInt32)},
				"int64":        {in: int64(bw.MaxInt32), out: int(bw.MaxInt32)},
				"int":          {in: bw.MaxInt, out: int(bw.MaxInt)},
				"uint8":        {in: bw.MaxUint8, out: int(bw.MaxUint8)},
				"uint16":       {in: bw.MaxUint16, out: int(bw.MaxUint16)},
				"uint32":       {in: uint32(bw.MaxInt32), out: int(bw.MaxInt32)},
				"uint64":       {in: uint64(bw.MaxInt), out: bw.MaxInt},
				"uint":         {in: uint(bw.MaxInt), out: bw.MaxInt},
				"float32(2.0)": {in: float32(2.0), out: 2},
				"float64(2.0)": {in: float64(2.0), out: 2},
			} {
				tests[k+", ok"] = bwtesting.Case{
					In:  []interface{}{v.in},
					Out: []interface{}{v.out, true},
				}
			}
			for k, v := range map[string]interface{}{
				"uint":          uint(bw.MaxUint),
				"uint64":        bw.MaxUint64,
				"float32(2.71)": float32(2.71),
				"float64(2.71)": float64(2.71),
			} {
				tests[k+", !ok"] = bwtesting.Case{
					In:  []interface{}{v},
					Out: []interface{}{0, false},
				}
			}
			return tests
		}(),
	)
}

func TestUint(t *testing.T) {
	bwtesting.BwRunTests(t,
		bwtype.Uint,
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for k, v := range map[string]struct {
				in  interface{}
				out uint
			}{
				"bw.MaxInt8":         {in: bw.MaxInt8, out: uint(bw.MaxInt8)},
				"bw.MaxInt16":        {in: bw.MaxInt16, out: uint(bw.MaxInt16)},
				"bw.MaxInt32":        {in: bw.MaxInt32, out: uint(bw.MaxInt32)},
				"int64(bw.MaxInt32)": {in: int64(bw.MaxInt32), out: uint(bw.MaxInt32)},
				"bw.MaxInt":          {in: bw.MaxInt, out: uint(bw.MaxInt)},
				"bw.MaxUint8":        {in: bw.MaxUint8, out: uint(bw.MaxUint8)},
				"bw.MaxUint16":       {in: bw.MaxUint16, out: uint(bw.MaxUint16)},
				"bw.MaxUint32":       {in: bw.MaxUint32, out: uint(bw.MaxUint32)},
				"uint64(bw.MaxUint)": {in: uint64(bw.MaxUint), out: bw.MaxUint},
				"bw.MaxUint":         {in: bw.MaxUint, out: bw.MaxUint},
				"float32(2.0)":       {in: float32(2.0), out: uint(2)},
				"float64(2.0)":       {in: float64(2.0), out: uint(2)},
			} {
				tests[k+", ok"] = bwtesting.Case{
					In:  []interface{}{v.in},
					Out: []interface{}{v.out, true},
				}
			}
			for k, v := range map[string]interface{}{
				"bw.MinInt8":    bw.MinInt8,
				"bw.MinInt16":   bw.MinInt16,
				"bw.MinInt32":   bw.MinInt32,
				"float32(2.71)": float32(2.71),
				"float64(2.71)": float64(2.71),
			} {
				tests[k+", !ok"] = bwtesting.Case{
					In:  []interface{}{v},
					Out: []interface{}{uint(0), false},
				}
			}
			return tests
		}(),
	)
}

func TestFloat64(t *testing.T) {
	bwtesting.BwRunTests(t,
		bwtype.Float64,
		func() map[string]bwtesting.Case {
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
			return tests
		}(),
	)
}

func TestMustRangeLimitFrom(t *testing.T) {

	bwtesting.BwRunTests(t,
		func(val interface{}) int {
			return bwtype.MustRangeLimitFrom(val).MustInt()
		},
		map[string]bwtesting.Case{
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
			"RangeLimit(-273)": {
				In:  []interface{}{bwtype.MustRangeLimitFrom(-273)},
				Out: []interface{}{-273},
			},
			"true": {
				In:    []interface{}{true},
				Panic: "\x1b[96;1m(bool)true\x1b[0m can not be a \x1b[97;1mRangeLimit\x1b[0m",
			},
		},
	)

	bwtesting.BwRunTests(t,
		func(val interface{}) bwtype.RangeLimitKind {
			return bwtype.MustRangeLimitFrom(val).Kind()
		},
		map[string]bwtesting.Case{
			"nil": {
				In:  []interface{}{nil},
				Out: []interface{}{bwtype.RangeLimitNil},
			},
			"273": {
				In:  []interface{}{273},
				Out: []interface{}{bwtype.RangeLimitInt},
			},
			"3.14": {
				In:  []interface{}{3.14},
				Out: []interface{}{bwtype.RangeLimitFloat64},
			},
			"path": {
				In:  []interface{}{bw.ValPath{}},
				Out: []interface{}{bwtype.RangeLimitPath},
			},
			"RangeLimit(nil)": {
				In:  []interface{}{bwtype.MustRangeLimitFrom(nil)},
				Out: []interface{}{bwtype.RangeLimitNil},
			},
			"RangeLimit(-273)": {
				In:  []interface{}{bwtype.MustRangeLimitFrom(-273)},
				Out: []interface{}{bwtype.RangeLimitInt},
			},
			"RangeLimit(-2.71)": {
				In:  []interface{}{bwtype.MustRangeLimitFrom(-2.71)},
				Out: []interface{}{bwtype.RangeLimitFloat64},
			},
			"RangeLimit(path)": {
				In:  []interface{}{bwtype.MustRangeLimitFrom(bw.ValPath{})},
				Out: []interface{}{bwtype.RangeLimitPath},
			},
		},
	)

	bwtesting.BwRunTests(t,
		"String",
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for k, v := range map[bwtype.RangeLimitKind]string{
				bwtype.RangeLimitNil:      "Nil",
				bwtype.RangeLimitInt:      "Int",
				bwtype.RangeLimitFloat64:  "Float64",
				bwtype.RangeLimitPath:     "Path",
				bwtype.RangeLimitPath + 1: "RangeLimitKind(4)",
			} {
				tests[v] = bwtesting.Case{
					V:   k,
					Out: []interface{}{v},
				}
			}
			return tests
		}(),
	)

	bwtesting.BwRunTests(t,
		func(val interface{}) float64 {
			return bwtype.MustRangeLimitFrom(val).MustFloat64()
		},
		map[string]bwtesting.Case{
			"nil": {
				In:    []interface{}{nil},
				Panic: "\x1b[96;1m(interface {})<nil>\x1b[0m is not \x1b[97;1mFloat64\x1b[0m",
			},
		},
	)

	bwtesting.BwRunTests(t,
		func(val interface{}) float64 {
			return bwtype.MustRangeLimitFrom(val).MustFloat64()
		},
		map[string]bwtesting.Case{
			"3.14": {
				In:  []interface{}{3.14},
				Out: []interface{}{3.14},
			},
		},
	)

	bwtesting.BwRunTests(t,
		func(val interface{}) bw.ValPath {
			return bwtype.MustRangeLimitFrom(val).MustPath()
		},
		map[string]bwtesting.Case{
			"nil": {
				In:    []interface{}{nil},
				Panic: "\x1b[96;1m(interface {})<nil>\x1b[0m is not \x1b[97;1mbw.ValPath\x1b[0m",
			},
		},
	)

}

func TestMustRangeFrom(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(a bwtype.A) (string, bwtype.RangeLimit, bwtype.RangeLimit) {
			r := bwtype.MustRangeFrom(a)
			return r.String(), r.Min(), r.Max()
		},
		map[string]bwtesting.Case{
			"..": {
				In: []interface{}{bwtype.A{}},
				Out: []interface{}{func(testName string) string { return testName },
					bwtype.RangeLimit{},
					bwtype.RangeLimit{},
				},
			},
			"2..": {
				In: []interface{}{bwtype.A{Min: 2}},
				Out: []interface{}{func(testName string) string { return testName },
					bwtype.MustRangeLimitFrom(2),
					bwtype.RangeLimit{},
				},
			},
			"..-3.14": {
				In: []interface{}{bwtype.A{Max: -3.14}},
				Out: []interface{}{func(testName string) string { return testName },
					bwtype.RangeLimit{},
					bwtype.MustRangeLimitFrom(-3.14),
				},
			},
			"2.71..273": {
				In: []interface{}{bwtype.A{Min: 2.71, Max: 273}},
				Out: []interface{}{func(testName string) string { return testName },
					bwtype.MustRangeLimitFrom(2.71),
					bwtype.MustRangeLimitFrom(273),
				},
			},
			"{{.}}..{{some}}": {
				In: []interface{}{bwtype.A{Min: bw.ValPath{}, Max: bw.ValPath{bw.ValPathItem{Type: bw.ValPathItemKey, Key: "some"}}}},
				Out: []interface{}{func(testName string) string { return testName },
					bwtype.MustRangeLimitFrom(bw.ValPath{}),
					bwtype.MustRangeLimitFrom(bw.ValPath{bw.ValPathItem{Type: bw.ValPathItemKey, Key: "some"}}),
				},
			},
			"$some.thing..{{good}}": {
				In: []interface{}{bwtype.A{
					Min: bw.ValPath{
						bw.ValPathItem{Type: bw.ValPathItemVar, Key: "some"},
						bw.ValPathItem{Type: bw.ValPathItemKey, Key: "thing"},
					},
					Max: bw.ValPath{bw.ValPathItem{Type: bw.ValPathItemKey, Key: "good"}},
				}},
				Out: []interface{}{func(testName string) string { return testName },
					bwtype.MustRangeLimitFrom(bw.ValPath{
						bw.ValPathItem{Type: bw.ValPathItemVar, Key: "some"},
						bw.ValPathItem{Type: bw.ValPathItemKey, Key: "thing"},
					}),
					bwtype.MustRangeLimitFrom(bw.ValPath{bw.ValPathItem{Type: bw.ValPathItemKey, Key: "good"}}),
				},
			},
			"2.71 > -273": {
				In:    []interface{}{bwtype.A{Min: 2.71, Max: -273}},
				Panic: "\x1b[38;5;201;1ma.Max\x1b[0m (\x1b[96;1m-273\x1b[0m) must not be \x1b[91;1mless\x1b[0m then \x1b[38;5;201;1ma.Min\x1b[0m (\x1b[96;1m2.71\x1b[0m)\x1b[0m",
			},
			"Min: true": {
				In:    []interface{}{bwtype.A{Min: true}},
				Panic: "\x1b[38;5;201;1ma.Min\x1b[0m (\x1b[96;1m(bool)true\x1b[0m) can not be a \x1b[97;1mRangeLimit\x1b[0m",
			},
			"Max: true": {
				In:    []interface{}{bwtype.A{Max: true}},
				Panic: "\x1b[38;5;201;1ma.Max\x1b[0m (\x1b[96;1m(bool)true\x1b[0m) can not be a \x1b[97;1mRangeLimit\x1b[0m",
			},
		},
		// "{{.}}..{{some}}",
	)
}

func TestRangeContains(t *testing.T) {
	bwtesting.BwRunTests(t,
		"Contains",
		map[string]bwtesting.Case{
			"nil not in ..": {
				V:   bwtype.MustRangeFrom(bwtype.A{}),
				In:  []interface{}{bwtype.MustRangeLimitFrom(nil)},
				Out: []interface{}{false},
			},
			"-273 in ..": {
				V:   bwtype.MustRangeFrom(bwtype.A{}),
				In:  []interface{}{bwtype.MustRangeLimitFrom(-273)},
				Out: []interface{}{true},
			},
			"-273 in ..0": {
				V:   bwtype.MustRangeFrom(bwtype.A{Max: 0}),
				In:  []interface{}{bwtype.MustRangeLimitFrom(-273)},
				Out: []interface{}{true},
			},
			"-273 not in 0..": {
				V:   bwtype.MustRangeFrom(bwtype.A{Min: 0}),
				In:  []interface{}{bwtype.MustRangeLimitFrom(-273)},
				Out: []interface{}{false},
			},
			"2.71 in 0..3.14": {
				V:   bwtype.MustRangeFrom(bwtype.A{Min: 0, Max: 3.14}),
				In:  []interface{}{bwtype.MustRangeLimitFrom(2.71)},
				Out: []interface{}{true},
			},
			"3.14 not in 0...2.71": {
				V:   bwtype.MustRangeFrom(bwtype.A{Min: 0, Max: 2.71}),
				In:  []interface{}{bwtype.MustRangeLimitFrom(3.14)},
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
				V:   bwtype.MustRangeLimitFrom(-273),
				In:  []interface{}{bwtype.MustRangeLimitFrom(bwtype.MustRangeLimitFrom(-273))},
				Out: []interface{}{true},
			},
			"-273 == -273.0": {
				V:   bwtype.MustRangeLimitFrom(-273),
				In:  []interface{}{bwtype.MustRangeLimitFrom(bwtype.MustRangeLimitFrom(-273.0))},
				Out: []interface{}{true},
			},
			"3.14 != 2.71": {
				V:   bwtype.MustRangeLimitFrom(3.14),
				In:  []interface{}{bwtype.MustRangeLimitFrom(2.71)},
				Out: []interface{}{false},
			},
		},
	)
}

func TestRangeKind(t *testing.T) {
	bwtesting.BwRunTests(t, "Kind",
		func() map[string]bwtesting.Case {
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
			return tests
		}(),
	)
	bwtesting.BwRunTests(t,
		"String",
		map[string]bwtesting.Case{
			"(bwtype.RangeMinMax + 1).String()": {
				V:   bwtype.RangeMinMax + 1,
				Out: []interface{}{"RangeKindValue(4)"},
			},
		},
	)
}

func TestRangeString(t *testing.T) {
	bwtesting.BwRunTests(t,
		"String",
		func() map[string]bwtesting.Case {
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
			return tests
		}(),
	)
}

func TestRangeLimitRangeMarshalJSON(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(r bwtype.Range) string {
			return bwjson.Pretty(r)
		},
		func() map[string]bwtesting.Case {
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
					In:  []interface{}{v.r},
					Out: []interface{}{fmt.Sprintf("%q", v.s)},
				}
			}
			return tests
		}(),
		nil,
	)
}
