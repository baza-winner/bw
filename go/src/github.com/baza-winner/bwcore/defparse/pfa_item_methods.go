package defparse

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwint"
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
	stackItem.value = stackItem.itemArray
	return false, nil
}

func _parseStackItemQwItem(pfa *pfaStruct) (skipPostProcess bool, err error) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemQwItem)
	stackItem.value = stackItem.itemString
	return false, nil
}

var underscoreRegexp = regexp.MustCompile("[_]+")

func _parseStackItemNumber(pfa *pfaStruct) (skipPostProcess bool, err error) {
	stackItem := pfa.getTopStackItem()
	source := underscoreRegexp.ReplaceAllLiteralString(stackItem.itemString, ``)
	if stackItem.value, err = parseNumber(source); err != nil {
		err = pfaErrorMake(pfa, failedToGetNumberError)
	}
	return false, err
}

func parseNumber(source string) (value interface{}, err error) {
	if strings.Contains(source, `.`) {
		var float64Val float64
		if float64Val, err = strconv.ParseFloat(source, 64); err == nil {
			value = float64Val
		}
	} else {
		var int64Val int64
		if int64Val, err = strconv.ParseInt(source, 10, 64); err == nil {
			if int64(bwint.MinInt8) <= int64Val && int64Val <= int64(bwint.MaxInt8) {
				value = int8(int64Val)
			} else if int64(bwint.MinInt16) <= int64Val && int64Val <= int64(bwint.MaxInt16) {
				value = int16(int64Val)
			} else if int64(bwint.MinInt32) <= int64Val && int64Val <= int64(bwint.MaxInt32) {
				value = int32(int64Val)
			} else {
				value = int64Val
			}
		}
	}
	return
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
	case "Bool", "String", "Int", "Number", "Map", "Array", "ArrayOf":
		stackItem.value = stackItem.itemString
	case "qw":
		pfa.pullRune()
		currRune, isEof := pfa.currRune()
		if isEof {
			err = pfaErrorMake(pfa, unexpectedRuneError)
		} else {
			pfa.state.setPrimary(expectSpaceOrQwItemOrDelimiter)
			switch currRune {
			case '<':
				stackItem.delimiter = runePtr('>')
			case '[':
				stackItem.delimiter = runePtr(']')
			case '(':
				stackItem.delimiter = runePtr(')')
			case '{':
				stackItem.delimiter = runePtr('}')
			default:
				switch {
				case unicode.IsPunct(currRune) || unicode.IsSymbol(currRune):
					stackItem.delimiter = runePtr(currRune)
				default:
					err = pfaErrorMake(pfa, unexpectedRuneError)
				}
			}
			stackItem.itemType = parseStackItemQw
		}
		return true, err

	default:
		err = pfaErrorMake(pfa, unknownWordError)
	}
	return false, err
}
