package defparse

import (
	"github.com/baza-winner/bwcore/bwjson"
)

type parseStackItem struct {
	itemType   parseStackItemType
	start      runePtrStruct
	itemArray  []interface{}
	itemMap    map[string]interface{}
	delimiter  rune
	currentKey string
	itemString string
	value      interface{}
}

func (stackItem *parseStackItem) GetDataForJson() interface{} {
	result := map[string]interface{}{}
	result["itemType"] = stackItem.itemType.String()
	result["start"] = stackItem.start.GetDataForJson()
	switch stackItem.itemType {
	case parseStackItemArray:
		result["itemArray"] = stackItem.itemArray
		result["value"] = stackItem.value
	case parseStackItemQw:
		result["delimiter"] = string(stackItem.delimiter)
		result["itemArray"] = stackItem.itemArray
		result["value"] = stackItem.value
	case parseStackItemQwItem:
		result["delimiter"] = string(stackItem.delimiter)
		result["itemString"] = stackItem.itemString
	case parseStackItemMap:
		result["itemMap"] = stackItem.itemMap
		result["value"] = stackItem.value
	case parseStackItemNumber, parseStackItemString, parseStackItemWord, parseStackItemKey:
		result["itemString"] = stackItem.itemString
		result["value"] = stackItem.value
	}
	return result
}

func (stackItem *parseStackItem) String() (result string) {
	return bwjson.PrettyJsonOf(stackItem)
}
