package defparse

import (
	"github.com/baza-winner/bwcore/bwjson"
)

type parseStack []parseStackItem

func (stack *parseStack) GetDataForJson() interface{} {
	result := []interface{}{}
	for _, item := range *stack {
		result = append(result, item.GetDataForJson())
	}
	return result
}

func (stack *parseStack) String() (result string) {
	return bwjson.PrettyJsonOf(stack)
}
