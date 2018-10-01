package defparser

import (
	// "fmt"
	"unicode"
)

type pfaPrimaryStateMethod func(*pfaStruct) (bool, error)

var pfaPrimaryStateMethods = map[parsePrimaryState]pfaPrimaryStateMethod{
	expectEOF:                       _expectEOF,
	expectValueOrSpace:              _expectValueOrSpace,
	expectArrayItemSeparatorOrSpace: _expectArrayItemSeparatorOrSpace,
	expectMapKeySeparatorOrSpace:    _expectMapKeySeparatorOrSpace,
	expectRocket:                    _expectRocket,
	expectSpaceOrMapKey:             _expectSpaceOrMapKey,
	expectMapKey:                    _expectMapKey,
	expectWord:                      _expectWord,
	expectDigit:                     _expectDigit,
	expectContentOf:                 _expectContentOf,
	expectEscapedContentOf:          _expectEscapedContentOf,
	expectSpaceOrQwItemOrDelimiter:  _expectSpaceOrQwItemOrDelimiter,
	expectEndOfQwItem:               _expectEndOfQwItem,
}

func _expectEOF(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if pfa.charPtr == nil {
		pfa.state.setPrimary(expectEOF)
	} else if pfa.state.secondary != orSpace || !unicode.IsSpace(*pfa.charPtr) {
		err = unexpectedCharError{}
	}
	return
}

func _expectValueOrSpace(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if pfa.charPtr == nil && len(pfa.stack) == 0 {
		pfa.state.setPrimary(expectEOF)
	} else {
		switch {
		case pfa.charPtr == nil:
			err = unexpectedCharError{}

		case *pfa.charPtr == '=' && pfa.state.secondary == orMapKeySeparator:
			pfa.state.setPrimary(expectRocket)

		case *pfa.charPtr == ':' && pfa.state.secondary == orMapKeySeparator:
			pfa.state.setPrimary(expectValueOrSpace)

		case *pfa.charPtr == ',' && pfa.state.secondary == orArrayItemSeparator:
			pfa.state.setPrimary(expectValueOrSpace)

		case unicode.IsSpace(*pfa.charPtr):

		case *pfa.charPtr == '{':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemMap, pos: pfa.pos, itemMap: map[string]interface{}{}})
			pfa.state.setPrimary(expectSpaceOrMapKey)

		case *pfa.charPtr == '[':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemArray, pos: pfa.pos, itemArray: []interface{}{}})
			pfa.state.setPrimary(expectValueOrSpace)

		case *pfa.charPtr == ']' && pfa.isTopStackItemOfType(parseStackItemArray):
			_ = pfa.getTopStackItemOfType(parseStackItemArray)
			needFinishTopStackItem = true

		case *pfa.charPtr == '-' || *pfa.charPtr == '+' || unicode.IsDigit(*pfa.charPtr):
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemNumber, pos: pfa.pos, itemString: string(*pfa.charPtr)})
			if unicode.IsDigit(*pfa.charPtr) {
				pfa.state.setSecondary(expectDigit, orUnderscoreOrDot)
			} else {
				pfa.state.setPrimary(expectDigit)
			}

		case *pfa.charPtr == '"' || *pfa.charPtr == '\'':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemString, pos: pfa.pos + 1, itemString: ``})
			quoted := singleQuoted
			if *pfa.charPtr == '"' {
				quoted = doubleQuoted
			}
			pfa.state.setTertiary(expectContentOf, quoted, stringToken)

		case unicode.IsLetter(*pfa.charPtr):
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemWord, pos: pfa.pos, itemString: string(*pfa.charPtr)})
			pfa.state.setPrimary(expectWord)

		default:
			err = unexpectedCharError{}
		}
	}
	return
}

func _expectArrayItemSeparatorOrSpace(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	switch {
	case pfa.charPtr == nil:
		err = unexpectedCharError{}
	case unicode.IsSpace(*pfa.charPtr):
		pfa.state.setSecondary(expectValueOrSpace, orArrayItemSeparator)
	case *pfa.charPtr == ',':
		pfa.state.setPrimary(expectValueOrSpace)
	default:
		err = unexpectedCharError{}
	}
	return
}

func _expectMapKeySeparatorOrSpace(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	switch {
	case pfa.charPtr == nil:
		err = unexpectedCharError{}
	case unicode.IsSpace(*pfa.charPtr):
		pfa.state.setSecondary(expectValueOrSpace, orMapKeySeparator)
	case *pfa.charPtr == ':':
		pfa.state.setPrimary(expectValueOrSpace)
	case *pfa.charPtr == '=':
		pfa.state.setPrimary(expectRocket)
	default:
		err = unexpectedCharError{}
	}
	return
}

func _expectRocket(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if pfa.charPtr != nil && *pfa.charPtr == '>' {
		pfa.state.setPrimary(expectValueOrSpace)
	} else {
		err = unexpectedCharError{}
	}
	return
}

func _expectSpaceOrMapKey(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	switch {
	case pfa.charPtr == nil:
		err = unexpectedCharError{}
	case unicode.IsSpace(*pfa.charPtr):
	case unicode.IsLetter(*pfa.charPtr):
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, pos: pfa.pos, itemString: string(*pfa.charPtr)})
		pfa.state.setPrimary(expectMapKey)
	case *pfa.charPtr == '"':
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, pos: pfa.pos, itemString: ``})
		pfa.state.setTertiary(expectContentOf, doubleQuoted, keyToken)
	case *pfa.charPtr == '\'':
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, pos: pfa.pos, itemString: ``})
		pfa.state.setTertiary(expectContentOf, singleQuoted, keyToken)
	case *pfa.charPtr == ',' && pfa.state.primary == expectSpaceOrMapKey && pfa.state.secondary == orMapValueSeparator:
		pfa.state.setPrimary(expectSpaceOrMapKey)
	case *pfa.charPtr == '}':
		_ = pfa.getTopStackItemOfType(parseStackItemMap)
		needFinishTopStackItem = true
	default:
		err = unexpectedCharError{}
	}
	return
}

func _expectMapKey(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemKey)
	switch {
	case pfa.charPtr == nil:
		err = unexpectedCharError{}
	case unicode.IsLetter(*pfa.charPtr):
		stackItem.itemString = stackItem.itemString + string(*pfa.charPtr)
	default:
		needFinishTopStackItem = true
	}
	return
}

func _expectWord(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemWord)
	switch {
	case pfa.charPtr != nil && unicode.IsLetter(*pfa.charPtr):
		stackItem.itemString = stackItem.itemString + string(*pfa.charPtr)
	default:
		needFinishTopStackItem = true
	}
	return
}

func _expectDigit(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemNumber)
	if pfa.state.secondary == noSecondaryState {
		switch {
		case pfa.charPtr != nil && unicode.IsDigit(*pfa.charPtr):
			stackItem.itemString = stackItem.itemString + string(*pfa.charPtr)
			pfa.state.secondary = orUnderscoreOrDot
		default:
			err = unexpectedCharError{}
		}
	} else {
		switch {
		case pfa.state.secondary == orUnderscoreOrDot && pfa.charPtr != nil && *pfa.charPtr == '.':
			pfa.state.secondary = orUnderscore
			fallthrough
		case pfa.charPtr != nil && (unicode.IsDigit(*pfa.charPtr) || *pfa.charPtr == '_'):
			stackItem.itemString = stackItem.itemString + string(*pfa.charPtr)
		default:
			needFinishTopStackItem = true
		}
	}
	return
}

func _expectContentOf(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if pfa.charPtr == nil {
		err = unexpectedCharError{}
	} else {
		itemType := parseStackItemString
		if pfa.state.tertiary == keyToken {
			itemType = parseStackItemKey
		}
		stackItem := pfa.getTopStackItemOfType(itemType)
		if pfa.state.secondary == doubleQuoted && *pfa.charPtr == '"' ||
			pfa.state.secondary == singleQuoted && *pfa.charPtr == '\'' {
			needFinishTopStackItem = true
		} else if *pfa.charPtr == '\\' {
			pfa.state.primary = expectEscapedContentOf
		} else {
			stackItem.itemString = stackItem.itemString + string(*pfa.charPtr)
		}
	}
	return
}

func _expectEscapedContentOf(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if pfa.charPtr == nil {
		err = unexpectedCharError{}
	} else {
		itemType := parseStackItemString
		if pfa.state.tertiary == keyToken {
			itemType = parseStackItemKey
		}
		stackItem := pfa.getTopStackItemOfType(itemType)
		var actualVal string
		switch *pfa.charPtr {
		case '"':
			actualVal = "\""
		case '\'':
			actualVal = "'"
		case '\\':
			actualVal = "\\"
		default:
			if pfa.state.secondary == doubleQuoted {
				switch *pfa.charPtr {
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
			err = unexpectedCharError{}
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
	case pfa.charPtr == nil:
		err = unexpectedCharError{}
	case unicode.IsSpace(*pfa.charPtr):
	case *pfa.charPtr == stackItem.delimiter:
		needFinishTopStackItem = true
	default:
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemQwItem, pos: pfa.pos, itemString: string(*pfa.charPtr), delimiter: stackItem.delimiter})
		pfa.state.setPrimary(expectEndOfQwItem)
	}
	return
}

func _expectEndOfQwItem(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemQwItem)
	switch {
	case pfa.charPtr == nil:
		err = unexpectedCharError{}
	case unicode.IsSpace(*pfa.charPtr) || *pfa.charPtr == stackItem.delimiter:
		needFinishTopStackItem = true
	default:
		stackItem.itemString += string(*pfa.charPtr)
	}
	return
}
