package defparse

import (
	"encoding/json"
)

type parseStackItem struct {
	itemType   parseStackItemType
	pos        int
	itemArray  []interface{}
	itemMap    map[string]interface{}
	delimiter  rune
	currentKey string
	itemString string
	value      interface{}
}

func (self *parseStackItem) getDataForJson() interface{} {
	result := map[string]interface{}{}
	result["itemType"] = self.itemType.String()
	result["pos"] = self.pos
	switch self.itemType {
	case parseStackItemArray:
		result["itemArray"] = self.itemArray
		result["value"] = self.value
	case parseStackItemQw:
		result["delimiter"] = string(self.delimiter)
		result["itemArray"] = self.itemArray
		result["value"] = self.value
	case parseStackItemQwItem:
		result["delimiter"] = string(self.delimiter)
		result["itemString"] = self.itemString
	case parseStackItemMap:
		result["itemMap"] = self.itemMap
		result["value"] = self.value
	case parseStackItemNumber, parseStackItemString, parseStackItemWord, parseStackItemKey:
		result["itemString"] = self.itemString
		result["value"] = self.value
	}
	return result
}

func (stackItem *parseStackItem) String() (result string) {
	bytes, _ := json.MarshalIndent(stackItem.getDataForJson(), ``, `  `)
	result = string(bytes[:]) // https://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string/18615786#18615786
	return
}
