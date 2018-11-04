// Package bwjson Предоставляет функцию PrettyJson - wrapper для json.MarshalIndent.
package bwjson

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/bwerr"
)

// type Jsonable interface {
// 	DataForJSON() interface{}
// }

// Pretty - wrapper для json.MarshalIndent
func Pretty(v interface{}) (result string) {
	if bytes, err := json.MarshalIndent(v, "", "  "); err != nil {
		bwerr.Panic("Pretty: failed to encode to json value %#v", v)
	} else {
		result = string(bytes[:])
	}
	return
}
