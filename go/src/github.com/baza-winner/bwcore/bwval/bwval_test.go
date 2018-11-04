package bwval

import (
	"testing"

	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
)

func TestVarPathFrom(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"some.thing": {
			In: []interface{}{"some.thing"},
			Out: []interface{}{
				VarPath{[]varPathItem{
					{Type: varPathItemKey, Key: "some"},
					{Type: varPathItemKey, Key: "thing"},
				}},
				nil,
			},
		},
	}

	// testsToRun := tests
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "zero number", "int number with underscore")
	// bwmap.CropMap(tests, "UnexpectedItem")
	// bwmap.CropMap(tests, "qw/Bool String Int Number Map Array ArrayOf/")
	bwtesting.BwRunTests(t, VarPathFrom, tests)
}
