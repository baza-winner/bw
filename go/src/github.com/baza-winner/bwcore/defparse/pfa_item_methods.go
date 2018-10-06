package defparse

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
  "github.com/baza-winner/bwcore/bwerror"
)

type pfaItemFinishMethod func(*pfaStruct) (bool, error)

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

func _parseStackItemKey(pfa *pfaStruct) (skipPostProcess bool, err error) {
	return false, nil
}

func _parseStackItemString(pfa *pfaStruct) (skipPostProcess bool, err error) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemString
	return false, nil
}

func _parseStackItemMap(pfa *pfaStruct) (skipPostProcess bool, err error) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemMap
	return false, nil
}

func _parseStackItemArray(pfa *pfaStruct) (skipPostProcess bool, err error) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemArray
	return false, nil
}

func _parseStackItemQw(pfa *pfaStruct) (skipPostProcess bool, err error) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemQw)
	if pfa.charPtr == nil {
		pfa.panic()
	}
	if unicode.IsSpace(*pfa.charPtr) {
		pfa.state.setPrimary(expectSpaceOrQwItemOrDelimiter)
	} else {
		if len(pfa.stack) < 2 {
			pfa.panic()
		}
		stackSubItem := pfa.stack[len(pfa.stack)-1]
		pfa.stack = pfa.stack[:len(pfa.stack)-1]
		stackItem = pfa.getTopStackItemOfType(parseStackItemArray)
		stackItem.itemArray = append(stackItem.itemArray, stackSubItem.itemArray...)
		pfa.state.setPrimary(expectArrayItemSeparatorOrSpace)
	}
	return true, nil
}

func _parseStackItemQwItem(pfa *pfaStruct) (skipPostProcess bool, err error) {
	if len(pfa.stack) < 2 {
		pfa.panic()
	}
	stackSubItem := pfa.popStackItem()

	stackItem := pfa.getTopStackItemOfType(parseStackItemQw)
	stackItem.itemArray = append(stackItem.itemArray, stackSubItem.itemString)
	return _parseStackItemQw(pfa)

	return true, nil
}

const (
	maxUint = ^uint(0)
	maxInt  = int(maxUint >> 1)
	minInt  = -maxInt - 1
)

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
			if int64(minInt) <= int64Val && int64Val <= int64(maxInt) {
				stackItem.value = int(int64Val)
			} else {
				stackItem.value = int64Val
			}
		}
	}
	if err != nil {
		err = failedToGetNumberError{}
	}
	return false, err
}

func _parseStackItemWord(pfa *pfaStruct) (skipPostProcess bool, err error) {
	stackItem := pfa.getTopStackItem()
	switch stackItem.itemString {
	case "true":
		stackItem.value = true
	case "false":
		stackItem.value = false
	case "nil", "null":
		stackItem.value = nil
	case "qw":
		if len(pfa.stack) < 2 || pfa.getTopStackItem(-2).itemType != parseStackItemArray {
			err = unexpectedWordError{}
		} else {
			if pfa.charPtr == nil {
				err = unexpectedCharError{}
			} else {
				pfa.state.setPrimary(expectSpaceOrQwItemOrDelimiter)
				switch *pfa.charPtr {
				case '<':
					stackItem.delimiter = '>'
				case '[':
					stackItem.delimiter = ']'
				case '(':
					stackItem.delimiter = ')'
				case '{':
					stackItem.delimiter = '}'
				default:
					switch {
					case unicode.IsPunct(*pfa.charPtr) || unicode.IsSymbol(*pfa.charPtr):
						stackItem.delimiter = *pfa.charPtr
					default:
						err = unexpectedCharError{}
					}
				}
				if pfa.state.primary == expectSpaceOrQwItemOrDelimiter {
					stackItem.itemType = parseStackItemQw
					stackItem.itemArray = []interface{}{}
				}
			}
		}
		return true, err

	default:
		err = unknownWordError{}
	}
	return false, err
}
