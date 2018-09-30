package defparser

import (
	"encoding/json"
)

type parseStack []parseStackItem

func (stack *parseStack) getDataForJson() interface{} {
	result := []interface{}{}
	for _, item := range *stack {
		result = append(result, item.getDataForJson())
	}
	return result
}

func (stack *parseStack) String() (result string) {
	bytes, _ := json.MarshalIndent(stack.getDataForJson(), ``, `  `)
	result = string(bytes[:]) // https://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string/18615786#18615786
	return
}
