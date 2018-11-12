// +build trace

package val

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/bwdebug"
	"github.com/baza-winner/bwcore/bwjson"
)

func init() {
	trace = func(r rune, primary PrimaryState, secondary SecondaryState, stack []stackItem) {
		bwdebug.Print(
			"r", string(r),
			"primary", primary.String(),
			"secondary", secondary.String(),
			"stack", bwjson.Pretty(stack),
		)
	}
}

func (v stackItem) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{}
	result["it"] = v.it.String()
	result["s"] = v.s
	result["delimiter"] = string(v.delimiter)
	result["result"] = v.result
	return json.Marshal(result)
}
