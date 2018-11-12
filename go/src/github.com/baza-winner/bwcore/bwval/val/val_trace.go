// +build trace

package val

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/bwdebug"
	"github.com/baza-winner/bwcore/bwjson"
)

func init() {
	trace = func(r rune, primary PrimaryState, secondary SecondaryState, stack []StackItem, boolVarName string, boolVarVal bool) {
		bwdebug.Print(
			"r", string(r),
			"primary", primary.String(),
			"secondary", secondary.String(),
			boolVarName, boolVarVal,
			"stack", bwjson.Pretty(stack),
		)
	}
}

func (v StackItem) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{}
	result["Kind"] = v.Kind.String()
	result["S"] = v.S
	result["Delimiter"] = string(v.Delimiter)
	result["Result"] = v.Result
	return json.Marshal(result)
}
