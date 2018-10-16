package defparse

import (
	"github.com/baza-winner/bwcore/bwerror"
)

type pfaItemFinishMethod func(*pfaStruct)

var pfaItemFinishMethods = map[parseStackItemType]pfaItemFinishMethod{
	parseStackItemKey:    _parseStackItemKey,
	parseStackItemString: _parseStackItemString,
	parseStackItemMap:    _parseStackItemMap,
	parseStackItemArray:  _parseStackItemArray,
	parseStackItemQw:     _parseStackItemQw,
	parseStackItemQwItem: _parseStackItemQwItem,
	parseStackItemNumber: _parseStackItemNumber,
	parseStackItemWord:   _parseStackItemWord,
}

func pfaItemFinishMethodsCheck() {
	itemType := parseStackItem_below_ + 1
	for itemType < parseStackItem_above_ {
		if _, ok := pfaItemFinishMethods[itemType]; !ok {
			bwerror.Panic("not defined <ansiOutline>pfaItemFinishMethods<ansi>[<ansiPrimaryLiteral>%s<ansi>]", itemType)
		}
		itemType += 1
	}
}

func _parseStackItemKey(pfa *pfaStruct) {
	return
}

func _parseStackItemString(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemString
	return
}

func _parseStackItemMap(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemMap
	return
}

func _parseStackItemArray(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemArray
	return
}

func _parseStackItemQw(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemQw)
	stackItem.value = stackItem.itemArray
	return
}

func _parseStackItemQwItem(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemQwItem)
	stackItem.value = stackItem.itemString
	return
}

func _parseStackItemNumber(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	if stackItem.value, pfa.err = _parseNumber(stackItem.itemString); pfa.err != nil {
		pfa.err = pfaErrorMake(pfa, failedToGetNumberError)
	}
	return
}

func _parseStackItemWord(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	switch stackItem.itemString {
	case "true":
		stackItem.value = true
	case "false":
		stackItem.value = false
	case "nil", "null":
	case "Bool", "String", "Int", "Number", "Map", "Array", "ArrayOf":
		stackItem.value = stackItem.itemString
	case "qw":
		pfa.pullRune()
		pfa.processStateDef(createStateDef(
			[]interface{}{unicodeOpenBraces, unicodePunct, unicodeSymbol,
				setPrimary{expectSpaceOrQwItemOrDelimiter},
				setTopItemDelimiter{pairForCurrRune{}},
				setTopItemType{parseStackItemQw},
			},
			[]interface{}{
				setError{unexpectedRuneError},
			},
		))
		pfa.skipPostProcess = true
	default:
		pfa.err = pfaErrorMake(pfa, unknownWordError)
	}
	return
}
