/*
Предоставляет функцию PrettyJson - wrapper для json.MarshalIndent.
*/
package bwjson

import (
	"encoding/json"
	"github.com/baza-winner/bwcore/bwerror"
)

/*
PrettyJson - wrapper для json.MarshalIndent
*/
func PrettyJson(v interface{}) (result string) {
	if bytes, err := json.MarshalIndent(v, "", "  "); err != nil {
		bwerror.Panic("failed to encode to json value %+v", v)
	} else {
		result = string(bytes[:])
	}
	return
}