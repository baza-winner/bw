package defparser

import (
	"encoding/json"
)

type parseStack []parseStackItem

func (stack parseStack) String() (result string) {
	array4Json := []map[string]interface{}{}
	for _, item := range stack {
		array4Json = append(array4Json, item.map4Json())
	}
	bytes, _ := json.MarshalIndent(array4Json, ``, `  `)
	result = string(bytes[:]) // https://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string/18615786#18615786
	return
}
