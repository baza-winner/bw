package defparser

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := map[string]struct {
		source string
		result map[string]interface{}
		err    error
	}{
		"first": {
			source: ``,
			result: nil,
			err:    nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		result, err := Parse(test.source)
		if err != nil {
			if err != test.err {
				t.Fatalf("Parse(%s) => err %v, want: %v", test.source, err, test.err)
			}
		} else if !reflect.DeepEqual(result, test.result) {
			tstJson, _ := json.MarshalIndent(result, ``, `  `)
			etaJson, _ := json.MarshalIndent(test.result, ``, `  `)
			t.Fatalf("Parse(%s) => %s, want: %s", test.source, tstJson, etaJson)
		}
	}

}
