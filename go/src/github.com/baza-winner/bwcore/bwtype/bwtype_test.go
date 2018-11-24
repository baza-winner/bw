package bwtype_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwtype"
)

func TestMustInt(t *testing.T) {
	bwtesting.BwRunTests(t,
		bwtype.MustInt,
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for in, out := range map[interface{}]int{
				bw.MaxInt8:          int(bw.MaxInt8),
				bw.MaxInt16:         int(bw.MaxInt16),
				bw.MaxInt32:         int(bw.MaxInt32),
				int64(bw.MaxInt32):  int(bw.MaxInt32),
				bw.MaxInt:           int(bw.MaxInt),
				bw.MaxUint8:         int(bw.MaxUint8),
				bw.MaxUint16:        int(bw.MaxUint16),
				uint32(bw.MaxInt32): int(bw.MaxInt32),
				uint64(bw.MaxInt):   bw.MaxInt,
				uint(bw.MaxInt):     bw.MaxInt,
				float32(2.0):        2,
				float64(2.0):        2,
			} {
				tests[bw.Spew.Sprintf("%#v", in)] = bwtesting.Case{
					In:  []interface{}{in},
					Out: []interface{}{out},
				}
			}
			for in, out := range map[interface{}]string{
				bw.MaxUint64:  "\x1b[96;1m(uint64)18446744073709551615\x1b[0m is not \x1b[97;1mInt\x1b[0m",
				float32(2.71): "\x1b[96;1m(float32)2.71\x1b[0m is not \x1b[97;1mInt\x1b[0m",
				float64(2.71): "\x1b[96;1m(float64)2.71\x1b[0m is not \x1b[97;1mInt\x1b[0m",
			} {
				tests[bw.Spew.Sprintf("%#v", in)] = bwtesting.Case{
					In:    []interface{}{in},
					Panic: out,
				}
			}
			return tests
		}(),
	)
}

func TestMustUint(t *testing.T) {
	bwtesting.BwRunTests(t,
		bwtype.MustUint,
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for in, out := range map[interface{}]uint{
				bw.MaxInt8:         uint(bw.MaxInt8),
				bw.MaxInt16:        uint(bw.MaxInt16),
				bw.MaxInt32:        uint(bw.MaxInt32),
				int64(bw.MaxInt32): uint(bw.MaxInt32),
				bw.MaxInt:          uint(bw.MaxInt),
				bw.MaxUint8:        uint(bw.MaxUint8),
				bw.MaxUint16:       uint(bw.MaxUint16),
				bw.MaxUint32:       uint(bw.MaxUint32),
				uint64(bw.MaxUint): bw.MaxUint,
				bw.MaxUint:         bw.MaxUint,
				float32(2.0):       uint(2),
				float64(2.0):       uint(2),
			} {
				tests[bw.Spew.Sprintf("%#v", in)] = bwtesting.Case{
					In:  []interface{}{in},
					Out: []interface{}{out},
				}
			}
			for in, out := range map[interface{}]string{
				int8(-1):      "\x1b[96;1m(int8)-1\x1b[0m is not \x1b[97;1mUint\x1b[0m",
				int16(-1):     "\x1b[96;1m(int16)-1\x1b[0m is not \x1b[97;1mUint\x1b[0m",
				int32(-1):     "\x1b[96;1m(int32)-1\x1b[0m is not \x1b[97;1mUint\x1b[0m",
				int(-1):       "\x1b[96;1m(int)-1\x1b[0m is not \x1b[97;1mUint\x1b[0m",
				float32(-1):   "\x1b[96;1m(float32)-1\x1b[0m is not \x1b[97;1mUint\x1b[0m",
				float32(3.14): "\x1b[96;1m(float32)3.14\x1b[0m is not \x1b[97;1mUint\x1b[0m",
				float64(-1):   "\x1b[96;1m(float64)-1\x1b[0m is not \x1b[97;1mUint\x1b[0m",
				float64(3.14): "\x1b[96;1m(float64)3.14\x1b[0m is not \x1b[97;1mUint\x1b[0m",
			} {
				tests[bw.Spew.Sprintf("%#v", in)] = bwtesting.Case{
					In:    []interface{}{in},
					Panic: out,
				}
			}
			return tests
		}(),
	)
}

func TestMustFloat64(t *testing.T) {
	bwtesting.BwRunTests(t,
		bwtype.MustFloat64,
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for in, out := range map[interface{}]float64{
				bw.MaxInt8:         float64(bw.MaxInt8),
				bw.MaxInt16:        float64(bw.MaxInt16),
				bw.MaxInt32:        float64(bw.MaxInt32),
				int64(bw.MaxInt64): float64(bw.MaxInt64),
				bw.MaxInt:          float64(bw.MaxInt),
				bw.MaxUint8:        float64(bw.MaxUint8),
				bw.MaxUint16:       float64(bw.MaxUint16),
				bw.MaxUint32:       float64(bw.MaxUint32),
				bw.MaxUint64:       float64(bw.MaxUint64),
				bw.MaxUint:         float64(bw.MaxUint),
				float32(0):         float64(0),
				float64(0):         float64(0),
			} {
				tests[bw.Spew.Sprintf("%#v", in)] = bwtesting.Case{
					In:  []interface{}{in},
					Out: []interface{}{out},
				}
			}
			for in, out := range map[interface{}]string{
				true: "\x1b[96;1m(bool)true\x1b[0m is not \x1b[97;1mFloat64\x1b[0m",
			} {
				tests[bw.Spew.Sprintf("%#v", in)] = bwtesting.Case{
					In:    []interface{}{in},
					Panic: out,
				}
			}
			return tests
		}(),
	)
}

func TestMustNumberFrom(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(val interface{}) (i int, u uint, f float64) {
			n := bwtype.MustNumberFrom(val)
			i, _ = bwtype.Int(n.Val())
			u, _ = bwtype.Uint(n.Val())
			f, _ = bwtype.Float64(n.Val())
			return
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for in, out := range map[interface{}]struct {
				i int
				u uint
				f float64
			}{
				int(1):     {i: 1, u: 1, f: float64(1)},
				bw.MaxUint: {i: 0, u: bw.MaxUint, f: float64(bw.MaxUint)},
				-1:         {i: -1, u: 0, f: float64(-1)},
				3.0:        {i: 3, u: 3, f: float64(3.0)},
				-3.0:       {i: -3, u: 0, f: float64(-3.0)},
				3.14:       {i: 0, u: 0, f: float64(3.14)},
			} {
				tests[bw.Spew.Sprintf("%#v", in)] = bwtesting.Case{
					In:  []interface{}{in},
					Out: []interface{}{out.i, out.u, out.f},
				}
			}
			for in, out := range map[interface{}]string{
				true: "\x1b[96;1m(bool)true\x1b[0m can not be a \x1b[97;1mNumber\x1b[0m",
			} {
				tests[bw.Spew.Sprintf("%#v", in)] = bwtesting.Case{
					In:    []interface{}{in},
					Panic: out,
				}
			}
			return tests
		}(),
	)
}

func TestNumberIsEqualTo(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(a, b bwtype.Number) bool {
			return a.IsEqualTo(b)
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			type key struct{ a, b bwtype.Number }
			for in, out := range map[key]bool{
				key{a: bwtype.MustNumberFrom(1), b: bwtype.MustNumberFrom(1)}:       true,
				key{a: bwtype.MustNumberFrom(1), b: bwtype.MustNumberFrom(2)}:       false,
				key{a: bwtype.MustNumberFrom(1), b: bwtype.MustNumberFrom(-1)}:      false,
				key{a: bwtype.MustNumberFrom(1), b: bwtype.MustNumberFrom(3.14)}:    false,
				key{a: bwtype.MustNumberFrom(-1), b: bwtype.MustNumberFrom(1)}:      false,
				key{a: bwtype.MustNumberFrom(-1), b: bwtype.MustNumberFrom(-1)}:     true,
				key{a: bwtype.MustNumberFrom(-1), b: bwtype.MustNumberFrom(-2)}:     false,
				key{a: bwtype.MustNumberFrom(-1), b: bwtype.MustNumberFrom(3.14)}:   false,
				key{a: bwtype.MustNumberFrom(3.14), b: bwtype.MustNumberFrom(3.14)}: true,
				key{a: bwtype.MustNumberFrom(3), b: bwtype.MustNumberFrom(3.00)}:    true,
				key{a: bwtype.MustNumberFrom(-3), b: bwtype.MustNumberFrom(-3.00)}:  true,
			} {
				tests[bw.Spew.Sprintf("%s.isEqualTo(%s) => %v", in.a, in.b, out)] = bwtesting.Case{
					In:  []interface{}{in.a, in.b},
					Out: []interface{}{out},
				}
			}
			return tests
		}(),
	)
}

func TestNumberIsLessThan(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(a, b bwtype.Number) bool {
			return a.IsLessThan(b)
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			type key struct{ a, b bwtype.Number }
			for in, out := range map[key]bool{
				key{a: bwtype.MustNumberFrom(1), b: bwtype.MustNumberFrom(2)}:       true,
				key{a: bwtype.MustNumberFrom(1), b: bwtype.MustNumberFrom(1)}:       false,
				key{a: bwtype.MustNumberFrom(1), b: bwtype.MustNumberFrom(-1)}:      false,
				key{a: bwtype.MustNumberFrom(1), b: bwtype.MustNumberFrom(3.14)}:    true,
				key{a: bwtype.MustNumberFrom(-1), b: bwtype.MustNumberFrom(1)}:      true,
				key{a: bwtype.MustNumberFrom(-1), b: bwtype.MustNumberFrom(-1)}:     false,
				key{a: bwtype.MustNumberFrom(-1), b: bwtype.MustNumberFrom(-2)}:     false,
				key{a: bwtype.MustNumberFrom(-1), b: bwtype.MustNumberFrom(3.14)}:   true,
				key{a: bwtype.MustNumberFrom(3.14), b: bwtype.MustNumberFrom(3.14)}: false,
				key{a: bwtype.MustNumberFrom(3), b: bwtype.MustNumberFrom(3.00)}:    false,
				key{a: bwtype.MustNumberFrom(-3), b: bwtype.MustNumberFrom(-3.00)}:  false,
			} {
				tests[bw.Spew.Sprintf("%s.isLessThan(%s) => %v", in.a, in.b, out)] = bwtesting.Case{
					In:  []interface{}{in.a, in.b},
					Out: []interface{}{out},
				}
			}
			return tests
		}(),
	)
}

func TestNumberString(t *testing.T) {
	bwtesting.BwRunTests(t,
		"String",
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			type key struct{ a, b bwtype.Number }
			for in, out := range map[interface{}]string{
				3.14:  "3.14",
				3.0:   "3",
				-1:    "-1",
				-2.71: "-2.71",
			} {
				n := bwtype.MustNumberFrom(in)
				tests[bw.Spew.Sprintf("%s.String() => %q", n, out)] = bwtesting.Case{
					V:   n,
					Out: []interface{}{out},
				}
			}
			return tests
		}(),
	)
}

func TestNumberMarshalJSON(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(n bwtype.Number) string {
			bytes, _ := json.Marshal(n)
			return string(bytes)
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			type key struct{ a, b bwtype.Number }
			for in, out := range map[interface{}]string{
				3.14:  "3.14",
				3.0:   "3",
				-1:    "-1",
				-2.71: "-2.71",
			} {
				n := bwtype.MustNumberFrom(in)
				tests[bw.Spew.Sprintf("%s.MarshalJSON() => %q", n, out)] = bwtesting.Case{
					In:  []interface{}{n},
					Out: []interface{}{out},
				}
			}
			return tests
		}(),
	)
}

func TestMustRangeLimitFrom(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(val interface{}) (isNil bool, numStr string, pathStr string) {
			rl := bwtype.MustRangeLimitFrom(val)
			if n, isNumber := rl.Number(); isNumber {
				numStr = n.String()
			}
			if path, isPath := rl.Path(); isPath {
				pathStr = path.String()
			}
			isNil = rl.Nil()
			return
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			type key struct{ a, b bwtype.Number }
			for _, v := range []struct {
				val     interface{}
				isNil   bool
				numStr  string
				pathStr string
			}{
				{val: nil, isNil: true},
				{val: 3.14, numStr: "3.14"},
				{val: bw.ValPath{}, pathStr: "."},
			} {
				tests[bw.Spew.Sprintf("%#v", v.val)] = bwtesting.Case{
					In:  []interface{}{v.val},
					Out: []interface{}{v.isNil, v.numStr, v.pathStr},
				}
				v.val = bwtype.MustRangeLimitFrom(v.val)
				tests[bw.Spew.Sprintf("%#v", v.val)] = bwtesting.Case{
					In:  []interface{}{v.val},
					Out: []interface{}{v.isNil, v.numStr, v.pathStr},
				}
			}
			for in, out := range map[interface{}]string{
				true: "\x1b[96;1m(bool)true\x1b[0m can not be a \x1b[97;1mRangeLimit\x1b[0m",
			} {
				tests[bw.Spew.Sprintf("%#v", in)] = bwtesting.Case{
					In:    []interface{}{in},
					Panic: out,
				}
			}
			return tests
		}(),
	)
}

func TestRangeLimitMustNumber(t *testing.T) {
	bwtesting.BwRunTests(t,
		"MustNumber",
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for _, v := range []struct {
				in  interface{}
				out bwtype.Number
			}{
				{in: 3.14, out: bwtype.MustNumberFrom(3.14)},
			} {
				tests[bw.Spew.Sprintf(ansi.String("<ansi>MustRangeLimitFrom(<ansiVal>%#v<ansi>).MustNumber() => <ansiVal>%#v"), v.in, v.out)] = bwtesting.Case{
					V:   bwtype.MustRangeLimitFrom(v.in),
					Out: []interface{}{v.out},
				}
			}
			for _, v := range []struct {
				in  interface{}
				out string
			}{
				{in: nil, out: "\x1b[96;1m(interface {})<nil>\x1b[0m is not \x1b[97;1mNumber\x1b[0m"},
				{in: bw.ValPath{}, out: "\x1b[96;1m(bw.ValPath).\x1b[0m is not \x1b[97;1mNumber\x1b[0m"},
			} {
				tests[bw.Spew.Sprintf(ansi.String("<ansi>MustRangeLimitFrom(<ansiVal>%#v<ansi>).MustNumber() => %s"), v.in, v.out)] = bwtesting.Case{
					V:     bwtype.MustRangeLimitFrom(v.in),
					Panic: v.out,
				}
			}
			return tests
		}(),
	)
}

func TestRangeLimitMustPath(t *testing.T) {
	bwtesting.BwRunTests(t,
		"MustPath",
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for _, v := range []struct {
				in  interface{}
				out bw.ValPath
			}{
				{in: bw.ValPath{}, out: bw.ValPath{}},
			} {
				tests[bw.Spew.Sprintf(ansi.String("<ansi>MustRangeLimitFrom(<ansiVal>%#v<ansi>).MustPath() => <ansiVal>%#v"), v.in, v.out)] = bwtesting.Case{
					V:   bwtype.MustRangeLimitFrom(v.in),
					Out: []interface{}{v.out},
				}
			}
			for _, v := range []struct {
				in  interface{}
				out string
			}{
				{in: nil, out: "\x1b[96;1m(interface {})<nil>\x1b[0m is not \x1b[97;1mbw.ValPath\x1b[0m"},
				{in: 3.14, out: "\x1b[96;1m(bwtype.Number)3.14\x1b[0m is not \x1b[97;1mbw.ValPath\x1b[0m"},
			} {
				tests[bw.Spew.Sprintf(ansi.String("<ansi>MustRangeLimitFrom(<ansiVal>%#v<ansi>).MustPath() => %s"), v.in, v.out)] = bwtesting.Case{
					V:     bwtype.MustRangeLimitFrom(v.in),
					Panic: v.out,
				}
			}
			return tests
		}(),
	)
}

func TestRangeLimitString(t *testing.T) {
	bwtesting.BwRunTests(t,
		"String",
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for _, v := range []struct {
				in  interface{}
				out string
			}{
				{in: nil, out: ""},
				{in: 3.14, out: "3.14"},
				{in: bw.ValPath{
					bw.ValPathItem{Type: bw.ValPathItemKey, Key: "some"},
					bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: 3},
				}, out: "{{some.3}}"},
				{in: bw.ValPath{
					bw.ValPathItem{Type: bw.ValPathItemVar, Key: "some"},
					bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: 3},
				}, out: "$some.3"},
			} {
				tests[bw.Spew.Sprintf(ansi.String("%#v<ansi> => <ansiVal>%q"), v.in, v.out)] = bwtesting.Case{
					V:   bwtype.MustRangeLimitFrom(v.in),
					Out: []interface{}{v.out},
				}
			}
			return tests
		}(),
	)
}

func TestMustRangeFrom(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(a bwtype.A) (string, bwtype.RangeLimit, bwtype.RangeLimit) {
			rl := bwtype.MustRangeFrom(a)
			return rl.String(), rl.Min(), rl.Max()
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for _, v := range []struct {
				in       bwtype.A
				out      string
				min, max bwtype.RangeLimit
			}{
				{in: bwtype.A{}, out: ".."},
				{in: bwtype.A{Min: 3.14}, out: "3.14..", min: bwtype.MustRangeLimitFrom(3.14)},
				{in: bwtype.A{Max: 3.14}, out: "..3.14", max: bwtype.MustRangeLimitFrom(3.14)},
				{in: bwtype.A{Min: 2.71, Max: 3.14}, out: "2.71..3.14", min: bwtype.MustRangeLimitFrom(2.71), max: bwtype.MustRangeLimitFrom(3.14)},
				{
					in: bwtype.A{
						Min: bw.ValPath{bw.ValPathItem{Type: bw.ValPathItemVar, Key: "var"}},
						Max: bw.ValPath{bw.ValPathItem{Type: bw.ValPathItemKey, Key: "some"}},
					},
					out: "$var..{{some}}",
					min: bwtype.MustRangeLimitFrom(bw.ValPath{bw.ValPathItem{Type: bw.ValPathItemVar, Key: "var"}}),
					max: bwtype.MustRangeLimitFrom(bw.ValPath{bw.ValPathItem{Type: bw.ValPathItemKey, Key: "some"}}),
				},
			} {
				tests[bw.Spew.Sprintf(ansi.String("<ansi>MustRangeFrom(Min: <ansiVal>%#v<ansi>, Max: <ansiVal>%#v<ansi>) => <ansiVal>%s"), v.in.Min, v.in.Max, v.out)] = bwtesting.Case{
					In:  []interface{}{v.in},
					Out: []interface{}{v.out, v.min, v.max},
				}
			}
			for _, v := range []struct {
				in  bwtype.A
				out string
			}{
				{in: bwtype.A{Max: 2.71, Min: 3.14},
					out: "\x1b[38;5;201;1ma.Max\x1b[0m (\x1b[96;1m2.71\x1b[0m) must not be \x1b[91;1mless\x1b[0m then \x1b[38;5;201;1ma.Min\x1b[0m (\x1b[96;1m3.14\x1b[0m)\x1b[0m",
				},
				{in: bwtype.A{Max: true},
					out: "\x1b[38;5;201;1ma.Max\x1b[0m (\x1b[96;1m(bool)true\x1b[0m) can not be a \x1b[97;1mRangeLimit\x1b[0m",
				},
				{in: bwtype.A{Min: true},
					out: "\x1b[38;5;201;1ma.Min\x1b[0m (\x1b[96;1m(bool)true\x1b[0m) can not be a \x1b[97;1mRangeLimit\x1b[0m",
				},
			} {
				tests[bw.Spew.Sprintf(ansi.String("<ansi>MustRangeFrom(Min: <ansiVal>%#v<ansi>, Max: <ansiVal>%#v<ansi>) => <ansiVal>%s"), v.in.Min, v.in.Max, v.out)] = bwtesting.Case{
					In:    []interface{}{v.in},
					Panic: v.out,
				}
			}
			return tests
		}(),
	)
}

func TestRangeContains(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(rng bwtype.Range, val interface{}) bool {
			return rng.Contains(val)
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for _, v := range []struct {
				in  bwtype.A
				val interface{}
				out bool
			}{
				{in: bwtype.A{}, val: nil, out: false},
				{in: bwtype.A{}, val: 0, out: true},
				{in: bwtype.A{Min: 3.14}, val: 2.71, out: false},
				{in: bwtype.A{Min: 3.14}, val: 4, out: true},
				{in: bwtype.A{Max: 3.14}, val: 3.14, out: true},
				{in: bwtype.A{Max: 3.14}, val: 4, out: false},
				{in: bwtype.A{Min: 2.71, Max: 3.14}, val: 0, out: false},
				{in: bwtype.A{Min: 2.71, Max: 3.14}, val: 3, out: true},
				{in: bwtype.A{Min: 2.71, Max: 3.14}, val: 4, out: false},

				{in: bwtype.A{Min: 3.14, Max: bw.ValPath{}}, val: 4, out: false},
				{in: bwtype.A{Min: 3.14, Max: bw.ValPath{}}, val: 3.14, out: true},
				{in: bwtype.A{Min: bw.ValPath{}, Max: 3.14}, val: 3.14, out: true},
				{in: bwtype.A{Min: bw.ValPath{}, Max: 3.14}, val: 3, out: false},
			} {
				rng := bwtype.MustRangeFrom(v.in)
				tests[bw.Spew.Sprintf(ansi.String("<ansi>[<ansiVal>%s<ansi>].Contains(<ansiVal>%#v<ansi>) => <ansiVal>%v"), rng, v.val, v.out)] = bwtesting.Case{
					In:  []interface{}{rng, v.val},
					Out: []interface{}{v.out},
				}
			}
			return tests
		}(),
	)
}

func TestRangeMarshalJSON(t *testing.T) {
	bwtesting.BwRunTests(t,
		func(rng bwtype.Range) string {
			bytes, _ := rng.MarshalJSON()
			return string(bytes)
		},
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for _, v := range []struct {
				in  bwtype.A
				out string
			}{
				{in: bwtype.A{}, out: ".."},
				{in: bwtype.A{Min: 3.14}, out: "3.14.."},
				{in: bwtype.A{Max: 3.14}, out: "..3.14"},
				{in: bwtype.A{Min: 2.71, Max: 3.14}, out: "2.71..3.14"},
			} {
				rng := bwtype.MustRangeFrom(v.in)
				tests[bw.Spew.Sprintf(ansi.String("<ansi>[<ansiVal>%s<ansi>].MarshalJSON() => <ansiVal>%q"), rng, v.out)] = bwtesting.Case{
					In:  []interface{}{rng},
					Out: []interface{}{fmt.Sprintf("%q", v.out)},
				}
			}
			return tests
		}(),
	)

}

// func TestMustIntFrom(t *testing.T) {
// 	bwtesting.BwRunTests(t,
// 		func(val interface{}) int {
// 			return bwtype.MustInt(bwtype.MustNumberFrom(val).Val())
// 		},
// 		map[string]bwtesting.Case{
// 			// "nil": {
// 			// 	In:    []interface{}{nil},
// 			// 	Panic: "\x1b[96;1m(interface {})<nil>\x1b[0m is not \x1b[97;1mInt\x1b[0m",
// 			// },
// 			"273": {
// 				In:  []interface{}{273},
// 				Out: []interface{}{273},
// 			},
// 			"3.14": {
// 				In:    []interface{}{3.14},
// 				Panic: "\x1b[96;1m(float64)3.14\x1b[0m is not \x1b[97;1mInt\x1b[0m",
// 			},
// 			"Number(-273)": {
// 				In:  []interface{}{bwtype.MustNumberFrom(-273)},
// 				Out: []interface{}{-273},
// 			},
// 			"RangeLimit(-273)": {
// 				In:  []interface{}{bwtype.MustRangeLimitFrom(-273)},
// 				Out: []interface{}{-273},
// 			},
// 			"true": {
// 				In:    []interface{}{true},
// 				Panic: "\x1b[96;1m(bool)true\x1b[0m can not be a \x1b[97;1mNumber\x1b[0m",
// 			},
// 		},
// 	)
// }

// func TestNumberMustFloat64From(t *testing.T) {
// }

// func TestMustRangeLimitFrom(t *testing.T) {

// 	// bwtesting.BwRunTests(t,
// 	// 	func(val interface{}) bwtype.RangeLimitKind {
// 	// 		return bwtype.MustRangeLimitFrom(val).Kind()
// 	// 	},
// 	// 	map[string]bwtesting.Case{
// 	// 		"nil": {
// 	// 			In:  []interface{}{nil},
// 	// 			Out: []interface{}{bwtype.RangeLimitNil},
// 	// 		},
// 	// 		"273": {
// 	// 			In:  []interface{}{273},
// 	// 			Out: []interface{}{bwtype.RangeLimitInt},
// 	// 		},
// 	// 		"3.14": {
// 	// 			In:  []interface{}{3.14},
// 	// 			Out: []interface{}{bwtype.RangeLimitFloat64},
// 	// 		},
// 	// 		"path": {
// 	// 			In:  []interface{}{bw.ValPath{}},
// 	// 			Out: []interface{}{bwtype.RangeLimitPath},
// 	// 		},
// 	// 		"RangeLimit(nil)": {
// 	// 			In:  []interface{}{bwtype.MustRangeLimitFrom(nil)},
// 	// 			Out: []interface{}{bwtype.RangeLimitNil},
// 	// 		},
// 	// 		"RangeLimit(-273)": {
// 	// 			In:  []interface{}{bwtype.MustRangeLimitFrom(-273)},
// 	// 			Out: []interface{}{bwtype.RangeLimitInt},
// 	// 		},
// 	// 		"RangeLimit(-2.71)": {
// 	// 			In:  []interface{}{bwtype.MustRangeLimitFrom(-2.71)},
// 	// 			Out: []interface{}{bwtype.RangeLimitFloat64},
// 	// 		},
// 	// 		"RangeLimit(path)": {
// 	// 			In:  []interface{}{bwtype.MustRangeLimitFrom(bw.ValPath{})},
// 	// 			Out: []interface{}{bwtype.RangeLimitPath},
// 	// 		},
// 	// 	},
// 	// )

// 	// bwtesting.BwRunTests(t,
// 	// 	"String",
// 	// 	func() map[string]bwtesting.Case {
// 	// 		tests := map[string]bwtesting.Case{}
// 	// 		for k, v := range map[bwtype.RangeLimitKind]string{
// 	// 			bwtype.RangeLimitNil:      "Nil",
// 	// 			bwtype.RangeLimitInt:      "Int",
// 	// 			bwtype.RangeLimitFloat64:  "Float64",
// 	// 			bwtype.RangeLimitPath:     "Path",
// 	// 			bwtype.RangeLimitPath + 1: "RangeLimitKind(4)",
// 	// 		} {
// 	// 			tests[v] = bwtesting.Case{
// 	// 				V:   k,
// 	// 				Out: []interface{}{v},
// 	// 			}
// 	// 		}
// 	// 		return tests
// 	// 	}(),
// 	// )

// 	// bwtesting.BwRunTests(t,
// 	// 	func(val interface{}) float64 {
// 	// 		return bwtype.MustRangeLimitFrom(val).MustFloat64()
// 	// 	},
// 	// 	map[string]bwtesting.Case{
// 	// 		"nil": {
// 	// 			In:    []interface{}{nil},
// 	// 			Panic: "\x1b[96;1m(interface {})<nil>\x1b[0m is not \x1b[97;1mFloat64\x1b[0m",
// 	// 		},
// 	// 	},
// 	// )

// 	// bwtesting.BwRunTests(t,
// 	// 	func(val interface{}) float64 {
// 	// 		return bwtype.MustRangeLimitFrom(val).MustFloat64()
// 	// 	},
// 	// 	map[string]bwtesting.Case{
// 	// 		"3.14": {
// 	// 			In:  []interface{}{3.14},
// 	// 			Out: []interface{}{3.14},
// 	// 		},
// 	// 	},
// 	// )

// 	bwtesting.BwRunTests(t,
// 		func(val interface{}) bw.ValPath {
// 			return bwtype.MustRangeLimitFrom(val).MustPath()
// 		},
// 		map[string]bwtesting.Case{
// 			"nil": {
// 				In:    []interface{}{nil},
// 				Panic: "\x1b[96;1m(interface {})<nil>\x1b[0m is not \x1b[97;1mbw.ValPath\x1b[0m",
// 			},
// 		},
// 	)

// }

// func TestMustRangeFrom(t *testing.T) {
// 	bwtesting.BwRunTests(t,
// 		func(a bwtype.A) (string, bwtype.RangeLimit, bwtype.RangeLimit) {
// 			r := bwtype.MustRangeFrom(a)
// 			return r.String(), r.Min(), r.Max()
// 		},
// 		map[string]bwtesting.Case{
// 			"..": {
// 				In: []interface{}{bwtype.A{}},
// 				Out: []interface{}{func(testName string) string { return testName },
// 					bwtype.RangeLimit{},
// 					bwtype.RangeLimit{},
// 				},
// 			},
// 			"2..": {
// 				In: []interface{}{bwtype.A{Min: 2}},
// 				Out: []interface{}{func(testName string) string { return testName },
// 					bwtype.MustRangeLimitFrom(2),
// 					bwtype.RangeLimit{},
// 				},
// 			},
// 			"..-3.14": {
// 				In: []interface{}{bwtype.A{Max: -3.14}},
// 				Out: []interface{}{func(testName string) string { return testName },
// 					bwtype.RangeLimit{},
// 					bwtype.MustRangeLimitFrom(-3.14),
// 				},
// 			},
// 			"2.71..273": {
// 				In: []interface{}{bwtype.A{Min: 2.71, Max: 273}},
// 				Out: []interface{}{func(testName string) string { return testName },
// 					bwtype.MustRangeLimitFrom(2.71),
// 					bwtype.MustRangeLimitFrom(273),
// 				},
// 			},
// 			"{{.}}..{{some}}": {
// 				In: []interface{}{bwtype.A{Min: bw.ValPath{}, Max: bw.ValPath{bw.ValPathItem{Type: bw.ValPathItemKey, Key: "some"}}}},
// 				Out: []interface{}{func(testName string) string { return testName },
// 					bwtype.MustRangeLimitFrom(bw.ValPath{}),
// 					bwtype.MustRangeLimitFrom(bw.ValPath{bw.ValPathItem{Type: bw.ValPathItemKey, Key: "some"}}),
// 				},
// 			},
// 			"$some.thing..{{good}}": {
// 				In: []interface{}{bwtype.A{
// 					Min: bw.ValPath{
// 						bw.ValPathItem{Type: bw.ValPathItemVar, Key: "some"},
// 						bw.ValPathItem{Type: bw.ValPathItemKey, Key: "thing"},
// 					},
// 					Max: bw.ValPath{bw.ValPathItem{Type: bw.ValPathItemKey, Key: "good"}},
// 				}},
// 				Out: []interface{}{func(testName string) string { return testName },
// 					bwtype.MustRangeLimitFrom(bw.ValPath{
// 						bw.ValPathItem{Type: bw.ValPathItemVar, Key: "some"},
// 						bw.ValPathItem{Type: bw.ValPathItemKey, Key: "thing"},
// 					}),
// 					bwtype.MustRangeLimitFrom(bw.ValPath{bw.ValPathItem{Type: bw.ValPathItemKey, Key: "good"}}),
// 				},
// 			},
// 			"2.71 > -273": {
// 				In:    []interface{}{bwtype.A{Min: 2.71, Max: -273}},
// 				Panic: "\x1b[38;5;201;1ma.Max\x1b[0m (\x1b[96;1m-273\x1b[0m) must not be \x1b[91;1mless\x1b[0m then \x1b[38;5;201;1ma.Min\x1b[0m (\x1b[96;1m2.71\x1b[0m)\x1b[0m",
// 			},
// 			"Min: true": {
// 				In:    []interface{}{bwtype.A{Min: true}},
// 				Panic: "\x1b[38;5;201;1ma.Min\x1b[0m (\x1b[96;1m(bool)true\x1b[0m) can not be a \x1b[97;1mRangeLimit\x1b[0m",
// 			},
// 			"Max: true": {
// 				In:    []interface{}{bwtype.A{Max: true}},
// 				Panic: "\x1b[38;5;201;1ma.Max\x1b[0m (\x1b[96;1m(bool)true\x1b[0m) can not be a \x1b[97;1mRangeLimit\x1b[0m",
// 			},
// 		},
// 		// "{{.}}..{{some}}",
// 	)
// }

// func TestRangeContains(t *testing.T) {
// 	bwtesting.BwRunTests(t,
// 		"Contains",
// 		map[string]bwtesting.Case{
// 			"nil not in ..": {
// 				V:   bwtype.MustRangeFrom(bwtype.A{}),
// 				In:  []interface{}{bwtype.MustRangeLimitFrom(nil)},
// 				Out: []interface{}{false},
// 			},
// 			"-273 in ..": {
// 				V:   bwtype.MustRangeFrom(bwtype.A{}),
// 				In:  []interface{}{bwtype.MustRangeLimitFrom(-273)},
// 				Out: []interface{}{true},
// 			},
// 			"-273 in ..0": {
// 				V:   bwtype.MustRangeFrom(bwtype.A{Max: 0}),
// 				In:  []interface{}{bwtype.MustRangeLimitFrom(-273)},
// 				Out: []interface{}{true},
// 			},
// 			"-273 not in 0..": {
// 				V:   bwtype.MustRangeFrom(bwtype.A{Min: 0}),
// 				In:  []interface{}{bwtype.MustRangeLimitFrom(-273)},
// 				Out: []interface{}{false},
// 			},
// 			"2.71 in 0..3.14": {
// 				V:   bwtype.MustRangeFrom(bwtype.A{Min: 0, Max: 3.14}),
// 				In:  []interface{}{bwtype.MustRangeLimitFrom(2.71)},
// 				Out: []interface{}{true},
// 			},
// 			"3.14 not in 0...2.71": {
// 				V:   bwtype.MustRangeFrom(bwtype.A{Min: 0, Max: 2.71}),
// 				In:  []interface{}{bwtype.MustRangeLimitFrom(3.14)},
// 				Out: []interface{}{false},
// 			},
// 		},
// 	)
// }

// func TestIsEqualTo(t *testing.T) {
// 	bwtesting.BwRunTests(t,
// 		"IsEqualTo",
// 		map[string]bwtesting.Case{
// 			"-273 == -273": {
// 				V:   bwtype.MustNumberFrom(-273),
// 				In:  []interface{}{bwtype.MustNumberFrom(bwtype.MustNumberFrom(-273))},
// 				Out: []interface{}{true},
// 			},
// 			"-273 == -273.0": {
// 				V:   bwtype.MustNumberFrom(-273),
// 				In:  []interface{}{bwtype.MustNumberFrom(bwtype.MustNumberFrom(-273.0))},
// 				Out: []interface{}{true},
// 			},
// 			"3.14 != 2.71": {
// 				V:   bwtype.MustNumberFrom(3.14),
// 				In:  []interface{}{bwtype.MustNumberFrom(2.71)},
// 				Out: []interface{}{false},
// 			},
// 		},
// 	)
// }

// func TestRangeKind(t *testing.T) {
// 	bwtesting.BwRunTests(t, "Kind",
// 		func() map[string]bwtesting.Case {
// 			tests := map[string]bwtesting.Case{}
// 			for k, v := range map[bwtype.RangeKindValue]bwtype.Range{
// 				bwtype.RangeMinMax: bwtype.MustRangeFrom(bwtype.A{Min: -1, Max: 2}),
// 				bwtype.RangeMin:    bwtype.MustRangeFrom(bwtype.A{Min: -1}),
// 				bwtype.RangeMax:    bwtype.MustRangeFrom(bwtype.A{Max: 2}),
// 				bwtype.RangeNo:     bwtype.Range{},
// 			} {
// 				tests[k.String()] = bwtesting.Case{
// 					V:   v,
// 					Out: []interface{}{k},
// 				}
// 			}
// 			return tests
// 		}(),
// 	)
// 	bwtesting.BwRunTests(t,
// 		"String",
// 		map[string]bwtesting.Case{
// 			"(bwtype.RangeMinMax + 1).String()": {
// 				V:   bwtype.RangeMinMax + 1,
// 				Out: []interface{}{"RangeKindValue(4)"},
// 			},
// 		},
// 	)
// }

// func TestRangeString(t *testing.T) {
// 	bwtesting.BwRunTests(t,
// 		"String",
// 		func() map[string]bwtesting.Case {
// 			tests := map[string]bwtesting.Case{}
// 			for _, v := range map[bwtype.RangeKindValue]struct {
// 				r bwtype.Range
// 				s string
// 			}{
// 				bwtype.RangeMinMax: {
// 					r: bwtype.MustRangeFrom(bwtype.A{Min: -1, Max: 2}),
// 					s: "-1..2",
// 				},
// 				bwtype.RangeMin: {
// 					r: bwtype.MustRangeFrom(bwtype.A{Min: -1}),
// 					s: "-1..",
// 				},
// 				bwtype.RangeMax: {
// 					r: bwtype.MustRangeFrom(bwtype.A{Max: 2}),
// 					s: "..2",
// 				},
// 				bwtype.RangeNo: {
// 					r: bwtype.Range{},
// 					s: "..",
// 				},
// 			} {
// 				tests[v.s] = bwtesting.Case{
// 					V:   v.r,
// 					Out: []interface{}{v.s},
// 				}
// 			}
// 			return tests
// 		}(),
// 	)
// }

// func TestRangeLimitRangeMarshalJSON(t *testing.T) {
// 	bwtesting.BwRunTests(t,
// 		func(r bwtype.Range) string {
// 			return bwjson.Pretty(r)
// 		},
// 		func() map[string]bwtesting.Case {
// 			tests := map[string]bwtesting.Case{}
// 			for _, v := range map[bwtype.RangeKindValue]struct {
// 				r bwtype.Range
// 				s string
// 			}{
// 				bwtype.RangeMinMax: {
// 					r: bwtype.MustRangeFrom(bwtype.A{Min: -1, Max: 2}),
// 					s: "-1..2",
// 				},
// 				bwtype.RangeMin: {
// 					r: bwtype.MustRangeFrom(bwtype.A{Min: -1}),
// 					s: "-1..",
// 				},
// 				bwtype.RangeMax: {
// 					r: bwtype.MustRangeFrom(bwtype.A{Max: 2}),
// 					s: "..2",
// 				},
// 				bwtype.RangeNo: {
// 					r: bwtype.Range{},
// 					s: "..",
// 				},
// 			} {
// 				tests[v.s] = bwtesting.Case{
// 					In:  []interface{}{v.r},
// 					Out: []interface{}{fmt.Sprintf("%q", v.s)},
// 				}
// 			}
// 			return tests
// 		}(),
// 		// "-1..2",
// 	)
// }
