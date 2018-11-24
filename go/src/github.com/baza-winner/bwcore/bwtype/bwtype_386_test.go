package bwtype_test

import (
	"testing"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwtype"
)

func TestMustIntPlatformSpecific(t *testing.T) {
	bwtesting.BwRunTests(t,
		bwtype.MustInt,
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for in, out := range map[interface{}]int{
				int64(bw.MaxInt):  bw.MaxInt,
				uint32(bw.MaxInt): bw.MaxInt,
			} {
				tests[bw.Spew.Sprintf("%#v", in)] = bwtesting.Case{
					In:  []interface{}{in},
					Out: []interface{}{out},
				}
			}
			for in, out := range map[interface{}]string{
				bw.MaxInt64:  "\x1b[96;1m(int64)9223372036854775807\x1b[0m is not \x1b[97;1mInt\x1b[0m",
				bw.MinInt64:  "\x1b[96;1m(int64)-9223372036854775808\x1b[0m is not \x1b[97;1mInt\x1b[0m",
				bw.MaxUint32: "\x1b[96;1m(uint32)4294967295\x1b[0m is not \x1b[97;1mInt\x1b[0m",
				bw.MaxUint:   "\x1b[96;1m(uint)4294967295\x1b[0m is not \x1b[97;1mInt\x1b[0m",
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

func TestMustUintPlatformSpecific(t *testing.T) {
	bwtesting.BwRunTests(t,
		bwtype.MustUint,
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{}
			for in, out := range map[interface{}]uint{
				int64(bw.MaxUint):  bw.MaxUint,
				uint64(bw.MaxUint): bw.MaxUint,
			} {
				tests[bw.Spew.Sprintf("%#v", in)] = bwtesting.Case{
					In:  []interface{}{in},
					Out: []interface{}{out},
				}
			}
			for in, out := range map[interface{}]string{
				int64(-1):    "\x1b[96;1m(int64)-1\x1b[0m is not \x1b[97;1mUint\x1b[0m",
				bw.MaxInt64:  "\x1b[96;1m(int64)9223372036854775807\x1b[0m is not \x1b[97;1mUint\x1b[0m",
				bw.MaxUint64: "\x1b[96;1m(uint64)18446744073709551615\x1b[0m is not \x1b[97;1mUint\x1b[0m",
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
