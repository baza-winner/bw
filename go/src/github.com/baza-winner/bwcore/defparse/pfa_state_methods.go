package defparse

import (
	"unicode"

	"github.com/baza-winner/bwcore/bwerror"
)

type pfaPrimaryStateMethod func(*pfaStruct) (bool, error)

var pfaPrimaryStateMethods = map[parsePrimaryState]pfaPrimaryStateMethod{
	expectEOF:                      _expectEOF,
	expectValueOrSpace:             _expectValueOrSpace,
	expectRocket:                   _expectRocket,
	expectSpaceOrMapKey:            _expectSpaceOrMapKey,
	expectMapKey:                   _expectMapKey,
	expectWord:                     _expectWord,
	expectDigit:                    _expectDigit,
	expectContentOf:                _expectContentOf,
	expectEscapedContentOf:         _expectEscapedContentOf,
	expectSpaceOrQwItemOrDelimiter: _expectSpaceOrQwItemOrDelimiter,
	expectEndOfQwItem:              _expectEndOfQwItem,
}

func pfaPrimaryStateMethodsCheck() {
	expect := parsePrimaryState_below_ + 1
	for expect < parsePrimaryState_above_ {
		if _, ok := pfaPrimaryStateMethods[expect]; !ok {
			bwerror.Panic("not defined <ansiOutline>pfaPrimaryStateMethods<ansi>[<ansiPrimaryLiteral>%s<ansi>]", expect)
		}
		expect += 1
	}
}

func _expectEOF(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if pfa.curr.runePtr == nil {
		pfa.state.setPrimary(expectEOF)
	} else if pfa.state.secondary != orSpace || !unicode.IsSpace(*pfa.curr.runePtr) {
		err = pfaErrorMake(pfa, unexpectedRuneError)
	}
	return
}

func _expectValueOrSpace(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if pfa.curr.runePtr == nil && len(pfa.stack) == 0 {
		pfa.state.setPrimary(expectEOF)
	} else {
		switch {
		case pfa.curr.runePtr == nil:
			err = pfaErrorMake(pfa, unexpectedRuneError)

		case *pfa.curr.runePtr == '=' && pfa.state.secondary == orMapKeySeparator:
			pfa.state.setPrimary(expectRocket)

		case *pfa.curr.runePtr == ':' && pfa.state.secondary == orMapKeySeparator:
			pfa.state.setPrimary(expectValueOrSpace)

		case *pfa.curr.runePtr == ',' && pfa.state.secondary == orArrayItemSeparator:
			pfa.state.setPrimary(expectValueOrSpace)

		case unicode.IsSpace(*pfa.curr.runePtr):

		case *pfa.curr.runePtr == '{':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemMap, start: pfa.curr, itemMap: map[string]interface{}{}})
			pfa.state.setPrimary(expectSpaceOrMapKey)

		case *pfa.curr.runePtr == '<':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemQw, start: pfa.curr, itemArray: []interface{}{}, delimiter: '>'})
			pfa.state.setPrimary(expectSpaceOrQwItemOrDelimiter)

		case *pfa.curr.runePtr == '[':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemArray, start: pfa.curr, itemArray: []interface{}{}})
			pfa.state.setPrimary(expectValueOrSpace)

		case *pfa.curr.runePtr == ']' && pfa.isTopStackItemOfType(parseStackItemArray):
			needFinishTopStackItem = true

		case *pfa.curr.runePtr == '-' || *pfa.curr.runePtr == '+' || unicode.IsDigit(*pfa.curr.runePtr):
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemNumber, start: pfa.curr, itemString: string(*pfa.curr.runePtr)})
			if unicode.IsDigit(*pfa.curr.runePtr) {
				pfa.state.setSecondary(expectDigit, orUnderscoreOrDot)
			} else {
				pfa.state.setPrimary(expectDigit)
			}

		case *pfa.curr.runePtr == '"' || *pfa.curr.runePtr == '\'':
			pfa.pullRune()
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemString, start: pfa.curr, itemString: ``})
			pfa.pushRune()
			quoted := singleQuoted
			if *pfa.curr.runePtr == '"' {
				quoted = doubleQuoted
			}
			pfa.state.setTertiary(expectContentOf, quoted, stringToken)

		case unicode.IsLetter(*pfa.curr.runePtr):
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemWord, start: pfa.curr, itemString: string(*pfa.curr.runePtr)})
			pfa.state.setPrimary(expectWord)

		default:
			err = pfaErrorMake(pfa, unexpectedRuneError)
		}
	}
	return
}

func _expectRocket(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if pfa.curr.runePtr != nil && *pfa.curr.runePtr == '>' {
		pfa.state.setPrimary(expectValueOrSpace)
	} else {
		err = pfaErrorMake(pfa, unexpectedRuneError)
	}
	return
}

func _expectSpaceOrMapKey(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	switch {
	case pfa.curr.runePtr == nil:
		err = pfaErrorMake(pfa, unexpectedRuneError)
	case unicode.IsSpace(*pfa.curr.runePtr):
	case unicode.IsLetter(*pfa.curr.runePtr):
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, start: pfa.curr, itemString: string(*pfa.curr.runePtr)})
		pfa.state.setPrimary(expectMapKey)
	case *pfa.curr.runePtr == '"':
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, start: pfa.curr, itemString: ``})
		pfa.state.setTertiary(expectContentOf, doubleQuoted, keyToken)
	case *pfa.curr.runePtr == '\'':
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, start: pfa.curr, itemString: ``})
		pfa.state.setTertiary(expectContentOf, singleQuoted, keyToken)
	case *pfa.curr.runePtr == ',' && pfa.state.primary == expectSpaceOrMapKey && pfa.state.secondary == orMapValueSeparator:
		pfa.state.setPrimary(expectSpaceOrMapKey)
	case *pfa.curr.runePtr == '}':
		_ = pfa.getTopStackItemOfType(parseStackItemMap)
		needFinishTopStackItem = true
	default:
		err = pfaErrorMake(pfa, unexpectedRuneError)
	}
	return
}

func _expectMapKey(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemKey)
	switch {
	case pfa.curr.runePtr == nil:
		err = pfaErrorMake(pfa, unexpectedRuneError)
	case unicode.IsLetter(*pfa.curr.runePtr) || *pfa.curr.runePtr == '_' || unicode.IsDigit(*pfa.curr.runePtr):
		stackItem.itemString = stackItem.itemString + string(*pfa.curr.runePtr)
	default:
		pfa.pushRune()
		needFinishTopStackItem = true
	}
	return
	// stackItem := pfa.getTopStackItemOfType(parseStackItemKey)
	// switch {
	// case pfa.curr.runePtr == nil:
	// 	err = pfaErrorMake(pfa, unexpectedRuneError)
	// case !(unicode.IsSpace(*pfa.curr.runePtr) || *pfa.curr.runePtr == ':' || *pfa.curr.runePtr == '='):
	// 	stackItem.itemString = stackItem.itemString + string(*pfa.curr.runePtr)
	// default:
	// 	needFinishTopStackItem = true
	// }
	// return
}

func _expectWord(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemWord)
	switch {
	case pfa.curr.runePtr != nil && unicode.IsLetter(*pfa.curr.runePtr):
		stackItem.itemString = stackItem.itemString + string(*pfa.curr.runePtr)
	default:
		pfa.pushRune()
		needFinishTopStackItem = true
	}
	return
}

func _expectEndOfQwItem(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemQwItem)
	switch {
	case pfa.curr.runePtr == nil:
		err = pfaErrorMake(pfa, unexpectedRuneError)
	case unicode.IsSpace(*pfa.curr.runePtr) || *pfa.curr.runePtr == stackItem.delimiter:
		pfa.pushRune()
		needFinishTopStackItem = true
	default:
		stackItem.itemString += string(*pfa.curr.runePtr)
	}
	return
}

func _expectDigit(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemNumber)
	if pfa.state.secondary == noSecondaryState {
		switch {
		case pfa.curr.runePtr != nil && unicode.IsDigit(*pfa.curr.runePtr):
			stackItem.itemString = stackItem.itemString + string(*pfa.curr.runePtr)
			pfa.state.secondary = orUnderscoreOrDot
		default:
			err = pfaErrorMake(pfa, unexpectedRuneError)
		}
	} else {
		switch {
		case pfa.state.secondary == orUnderscoreOrDot && pfa.curr.runePtr != nil && *pfa.curr.runePtr == '.':
			pfa.state.secondary = orUnderscore
			fallthrough
		case pfa.curr.runePtr != nil && (unicode.IsDigit(*pfa.curr.runePtr) || *pfa.curr.runePtr == '_'):
			stackItem.itemString = stackItem.itemString + string(*pfa.curr.runePtr)
		default:
			pfa.pushRune()
			needFinishTopStackItem = true
		}
	}
	return
}

func _expectContentOf(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if pfa.curr.runePtr == nil {
		err = pfaErrorMake(pfa, unexpectedRuneError)
	} else {
		itemType := parseStackItemString
		if pfa.state.tertiary == keyToken {
			itemType = parseStackItemKey
		}
		stackItem := pfa.getTopStackItemOfType(itemType)
		if pfa.state.secondary == doubleQuoted && *pfa.curr.runePtr == '"' ||
			pfa.state.secondary == singleQuoted && *pfa.curr.runePtr == '\'' {
			needFinishTopStackItem = true
		} else if *pfa.curr.runePtr == '\\' {
			pfa.state.primary = expectEscapedContentOf
		} else {
			stackItem.itemString = stackItem.itemString + string(*pfa.curr.runePtr)
		}
	}
	return
}

func _expectEscapedContentOf(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if pfa.curr.runePtr == nil {
		err = pfaErrorMake(pfa, unexpectedRuneError)
	} else {
		itemType := parseStackItemString
		if pfa.state.tertiary == keyToken {
			itemType = parseStackItemKey
		}
		stackItem := pfa.getTopStackItemOfType(itemType)
		var actualVal string
		switch *pfa.curr.runePtr {
		case '"':
			actualVal = "\""
		case '\'':
			actualVal = "'"
		case '\\':
			actualVal = "\\"
		default:
			if pfa.state.secondary == doubleQuoted {
				switch *pfa.curr.runePtr {
				case 'a':
					actualVal = "\a"
				case 'b':
					actualVal = "\b"
				case 'f':
					actualVal = "\f"
				case 'n':
					actualVal = "\n"
				case 'r':
					actualVal = "\r"
				case 't':
					actualVal = "\t"
				case 'v':
					actualVal = "\v"
				}
			}
		}
		if len(actualVal) == 0 {
			err = pfaErrorMake(pfa, unexpectedRuneError)
		} else {
			stackItem.itemString = stackItem.itemString + actualVal
			pfa.state.primary = expectContentOf
		}
	}
	return
}

func _expectSpaceOrQwItemOrDelimiter(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemQw)
	switch {
	case pfa.curr.runePtr == nil:
		err = pfaErrorMake(pfa, unexpectedRuneError)
	case unicode.IsSpace(*pfa.curr.runePtr):
	case *pfa.curr.runePtr == stackItem.delimiter:
		needFinishTopStackItem = true
	default:
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemQwItem, start: pfa.curr, itemString: string(*pfa.curr.runePtr), delimiter: stackItem.delimiter})
		pfa.state.setPrimary(expectEndOfQwItem)
	}
	return
}
