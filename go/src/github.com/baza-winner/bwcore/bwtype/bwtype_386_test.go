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
		in interface{}
	}{
		"uint32": {in: bw.MaxUint32},
		"int64":  {in: bw.MaxInt64},
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
