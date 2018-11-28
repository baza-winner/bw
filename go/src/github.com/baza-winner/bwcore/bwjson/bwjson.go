// Package bwjson содержит bw-дополнение для package encoding/json
package bwjson

import (
	"encoding/json"
	"fmt"
)

// Pretty - wrapper для json.MarshalIndent
func Pretty(v interface{}) (result string) {
	if bytes, err := json.MarshalIndent(v, "", "  "); err != nil {
		panic(fmt.Sprintf("Pretty: failed to encode to json value %#v", v))
	} else {
		result = string(bytes[:])
	}
	return
}

func S(v interface{}) (result string) {
	if bytes, err := json.Marshal(v); err != nil {
		panic(fmt.Sprintf("S: failed to encode to json value %#v", v))
	} else {
		result = string(bytes[:])
	}
	return
}
