package bwtype_test

import (
	"testing"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwtype"
)

func TestIntSpecific(t *testing.T) {
	tests := map[string]bwtesting.Case{}
	for k, v := range map[string]struct {
		in  interface{}
		out int
	}{
		"uint64": {in: bw.MaxInt64, out: bw.MaxInt},
	} {
		tests[k+", ok"] = bwtesting.Case{
			In:  []interface{}{v.in},
			Out: []interface{}{v.out, true},
		}
	}

	bwtesting.BwRunTests(t, bwtype.Int, tests,
		nil,
	)
}
