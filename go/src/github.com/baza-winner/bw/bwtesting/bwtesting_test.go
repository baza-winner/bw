package bwtesting

import (
	"fmt"
	"github.com/baza-winner/bw/ansi"
	"github.com/baza-winner/bw/bwerror"
	// "github.com/baza-winner/bw/bwtesting"
	"testing"
)

type testCropMapStruct struct {
	m      map[string]interface{}
	crop   interface{}
	result map[string]interface{}
	err    error
}

func TestCropMap(t *testing.T) {
	tests := map[string]testCropMapStruct{
		"string": {
			m: map[string]interface{}{
				"some": "thing",
				"good": "is not bad",
			},
			crop: `some`,
			result: map[string]interface{}{
				"some": "thing",
			},
		},
		"[]string": {
			m: map[string]interface{}{
				"A": 1,
				"B": 2,
				"C": 3,
				"D": 4,
			},
			crop: []string{"B", "C"},
			result: map[string]interface{}{
				"B": 2,
				"C": 3,
			},
		},
		"map[string]interface{}": {
			m: map[string]interface{}{
				"A": 1,
				"B": 2,
				"C": 3,
				"D": 4,
			},
			crop: map[string]interface{}{
				"A": struct{}{},
				"D": struct{}{},
			},
			result: map[string]interface{}{
				"A": 1,
				"D": 4,
			},
		},
		"error": {
			m: map[string]interface{}{
				"A": 1,
				"B": 2,
				"C": 3,
				"D": 4,
			},
			crop: true,
			err:  bwerror.Error("<ansiOutline>crop<ansi> (<ansiPrimaryLiteral>%+v<ansi>) neither <ansiSecondaryLiteral>string<ansi>, nor <ansiSecondaryLiteral>[]string<ansi>, nor <ansiSecondaryLiteral>map[string]interface", true),
		},
	}
	testsToRun := tests
	for testName, test := range testsToRun {
		t.Logf(ansi.Ansi(`Header`, "Running test case <ansiPrimaryLiteral>%s"), testName)
		result := map[string]interface{}{}
		deepCopyJSON(test.m, result)
		err := CropMap(result, test.crop)
		testTitle := fmt.Sprintf("CropMap(%+v, %+v)\n", test.m, test.crop)
		CheckTestErrResult(t, err, test.err, result, test.result, testTitle)
	}
}

// https://stackoverflow.com/questions/51459083/deep-copying-maps-in-golang/51684750#51684750
func deepCopyJSON(src map[string]interface{}, dest map[string]interface{}) {
	for key, value := range src {
		switch src[key].(type) {
		case map[string]interface{}:
			dest[key] = map[string]interface{}{}
			deepCopyJSON(src[key].(map[string]interface{}), dest[key].(map[string]interface{}))
		default:
			dest[key] = value
		}
	}
}
