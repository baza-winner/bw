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
				bw.MaxInt64:  int(bw.MaxInt64),
				bw.MaxUint32: int(bw.MaxUint32),
			} {
				tests[bw.Spew.Sprintf("%#v", in)] = bwtesting.Case{
					In:  []interface{}{in},
					Out: []interface{}{out},
				}
			}
			for in, out := range map[interface{}]string{
				bw.MaxUint: "\x1b[96;1m(uint)18446744073709551615\x1b[0m is not \x1b[97;1mInt\x1b[0m",
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
				bw.MaxInt64:  uint(bw.MaxInt64),
				bw.MaxUint64: uint(bw.MaxUint64),
			} {
				tests[bw.Spew.Sprintf("%#v", in)] = bwtesting.Case{
					In:  []interface{}{in},
					Out: []interface{}{out},
				}
			}
			for in, out := range map[interface{}]string{
				int8(-1): "\x1b[96;1m(int8)-1\x1b[0m is not \x1b[97;1mUint\x1b[0m",
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
