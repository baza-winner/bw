package pathparse

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
)

type pfaItemFinishMethod func(*pfaStruct) (bool, error)

var pfaItemFinishMethods = map[parseStackItemType]pfaItemFinishMethod{
	parseStackItemNumber: _parseStackItemNumber,
	parseStackItemKey:    _parseStackItemKey,
	parseStackItemString: _parseStackItemString,
	// parseStackItemMap:    _parseStackItemMap,
	// parseStackItemArray:  _parseStackItemArray,
	// parseStackItemQw:     _parseStackItemQw,
	// parseStackItemQwItem: _parseStackItemQwItem,
	// parseStackItemNumber: _parseStackItemNumber,
	// parseStackItemWord:   _parseStackItemWord,
}

func pfaItemFinishMethodsCheck() {
	itemType := parseStackItemBelow + 1
	for itemType < parseStackItemAbove {
		if _, ok := pfaItemFinishMethods[itemType]; !ok {
			bwerr.Panic("not defined <ansiVar>pfaItemFinishMethods<ansi>[<ansiVal>%s<ansi>]", itemType)
		}
		itemType += 1
	}
}

// func _parseStackItemNumber(pfa *pfaStruct) (skipPostProcess bool, err error) {

// 	return false, nil
// }

func _parseStackItemKey(pfa *pfaStruct) (skipPostProcess bool, err error) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemString
	// pfa.pushRune()
	return false, nil
	// return false, nil
}

func _parseStackItemString(pfa *pfaStruct) (skipPostProcess bool, err error) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemString
	return false, nil
	return false, nil
}

// func _parsestackitemkey(pfa *pfaStruct) (skipPostProcess bool, err error) {
// 	return false, nil
// }

// func _parsestackitemkey(pfa *pfaStruct) (skipPostProcess bool, err error) {
// 	return false, nil
// }

// func _parseStackItemString(pfa *pfaStruct) (skipPostProcess bool, err error) {
// 	stackItem := pfa.getTopStackItem()
// 	stackItem.value = stackItem.itemString
// 	return false, nil
// }

// func _parseStackItemMap(pfa *pfaStruct) (skipPostProcess bool, err error) {
// 	stackItem := pfa.getTopStackItem()
// 	stackItem.value = stackItem.itemMap
// 	return false, nil
// }

// func _parseStackItemArray(pfa *pfaStruct) (skipPostProcess bool, err error) {
// 	stackItem := pfa.getTopStackItem()
// 	stackItem.value = stackItem.itemArray
// 	return false, nil
// }

// func _parseStackItemQw(pfa *pfaStruct) (skipPostProcess bool, err error) {
// 	stackItem := pfa.getTopStackItemOfType(parseStackItemQw)
// 	stackItem.value = stackItem.itemArray
// 	return false, nil
// }

// func _parseStackItemQwItem(pfa *pfaStruct) (skipPostProcess bool, err error) {
// 	stackItem := pfa.getTopStackItemOfType(parseStackItemQwItem)
// 	stackItem.value = stackItem.itemString
// 	return false, nil
// }

var underscoreRegexp = regexp.MustCompile("[_]+")

func _parseStackItemNumber(pfa *pfaStruct) (skipPostProcess bool, err error) {
	stackItem := pfa.getTopStackItem()
	source := underscoreRegexp.ReplaceAllLiteralString(stackItem.itemString, ``)
	if strings.Contains(source, `.`) {
		var float64Val float64
		if float64Val, err = strconv.ParseFloat(source, 64); err == nil {
			stackItem.value = float64Val
		}
	} else {
		var int64Val int64
		if int64Val, err = strconv.ParseInt(source, 10, 64); err == nil {
			if int64(bw.MinInt8) <= int64Val && int64Val <= int64(bw.MaxInt8) {
				stackItem.value = int8(int64Val)
			} else if int64(bw.MinInt16) <= int64Val && int64Val <= int64(bw.MaxInt16) {
				stackItem.value = int16(int64Val)
			} else if int64(bw.MinInt32) <= int64Val && int64Val <= int64(bw.MaxInt32) {
				stackItem.value = int32(int64Val)
			} else {
				stackItem.value = int64Val
			}
		}
	}
	if err != nil {
		err = pfaErrorMake(pfa, failedToGetNumberError)
	}
	return false, err
}

// func _parseStackItemWord(pfa *pfaStruct) (skipPostProcess bool, err error) {
// 	stackItem := pfa.getTopStackItem()
// 	switch stackItem.itemString {
// 	case "true":
// 		stackItem.value = true
// 	case "false":
// 		stackItem.value = false
// 	case "nil", "null":
// 		stackItem.value = nil
// 	case "Bool", "String", "Int", "Number", "Map", "Array", "ArrayOf":
// 		stackItem.value = stackItem.itemString
// 	case "qw":
// 		pfa.pullRune()
// 		if pfa.curr.runePtr == nil {
// 			err = pfaErrorMake(pfa, unexpectedRuneError)
// 		} else {
// 			pfa.state.setPrimary(expectSpaceOrQwItemOrDelimiter)
// 			switch *pfa.curr.runePtr {
// 			case '<':
// 				stackItem.delimiter = '>'
// 			case '[':
// 				stackItem.delimiter = ']'
// 			case '(':
// 				stackItem.delimiter = ')'
// 			case '{':
// 				stackItem.delimiter = '}'
// 			default:
// 				switch {
// 				case unicode.IsPunct(*pfa.curr.runePtr) || unicode.IsSymbol(*pfa.curr.runePtr):
// 					stackItem.delimiter = *pfa.curr.runePtr
// 				default:
// 					err = pfaErrorMake(pfa, unexpectedRuneError)
// 				}
// 			}
// 			stackItem.itemType = parseStackItemQw
// 			stackItem.itemArray = []interface{}{}
// 		}
// 		return true, err

// 	default:
// 		err = pfaErrorMake(pfa, unknownWordError)
// 	}
// 	return false, err
// }
