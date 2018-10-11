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
	currRune, isEOF := pfa.currRune()
	switch {
	case isEOF:
		pfa.state.setPrimary(expectEOF)
	case pfa.state.secondary == orSpace && unicode.IsSpace(currRune):
	default:
		err = pfaErrorMake(pfa, unexpectedRuneError)
	}
	return
}

func _expectValueOrSpace(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	currRune, isEOF := pfa.currRune()
	switch {
	case isEOF:
		if len(pfa.stack) == 0 {
			pfa.state.setPrimary(expectEOF)
		} else {
			err = pfaErrorMake(pfa, unexpectedRuneError)
		}

	case currRune == '=' && pfa.state.secondary == orMapKeySeparator:
		pfa.state.setPrimary(expectRocket)

	case currRune == ':' && pfa.state.secondary == orMapKeySeparator:
		pfa.state.setPrimary(expectValueOrSpace)

	case currRune == ',' && pfa.state.secondary == orArrayItemSeparator:
		pfa.state.setPrimary(expectValueOrSpace)

	case unicode.IsSpace(currRune):

	case currRune == '{':
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemMap, start: pfa.curr, itemMap: map[string]interface{}{}, delimiter: '}'})
		pfa.state.setPrimary(expectSpaceOrMapKey)

	case currRune == '<':
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemQw, start: pfa.curr, itemArray: []interface{}{}, delimiter: '>'})
		pfa.state.setPrimary(expectSpaceOrQwItemOrDelimiter)

	case currRune == '[':
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemArray, start: pfa.curr, itemArray: []interface{}{}, delimiter: ']'})
		pfa.state.setPrimary(expectValueOrSpace)

	case pfa.isTopStackItemOfType(parseStackItemArray) && pfa.getTopStackItem().delimiter == currRune:
		needFinishTopStackItem = true

	case currRune == '-' || currRune == '+':
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemNumber, start: pfa.curr, itemString: string(currRune)})
		pfa.state.setPrimary(expectDigit)

	case unicode.IsDigit(currRune):
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemNumber, start: pfa.curr, itemString: string(currRune)})
		pfa.state.setSecondary(expectDigit, orUnderscoreOrDot)

	case currRune == '"' || currRune == '\'':
		pfa.pullRune()
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemString, start: pfa.curr, itemString: ``, delimiter: currRune})
		pfa.pushRune()
		pfa.state.setSecondary(expectContentOf, stringToken)

	case unicode.IsLetter(currRune):
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemWord, start: pfa.curr, itemString: string(currRune)})
		pfa.state.setPrimary(expectWord)

	default:
		err = pfaErrorMake(pfa, unexpectedRuneError)
	}
	return
}

func _expectRocket(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	currRune, _ := pfa.currRune()
	switch currRune {
	case '>':
		pfa.state.setPrimary(expectValueOrSpace)
	default:
		err = pfaErrorMake(pfa, unexpectedRuneError)
	}
	return
}

func _expectSpaceOrMapKey(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	currRune, _ := pfa.currRune()
	switch {
	case unicode.IsSpace(currRune):
	case unicode.IsLetter(currRune):
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, start: pfa.curr, itemString: string(currRune)})
		pfa.state.setPrimary(expectWord)
	case currRune == '"' || currRune == '\'':
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, start: pfa.curr, itemString: ``, delimiter: currRune})
		pfa.state.setSecondary(expectContentOf, keyToken)
	case currRune == ',' && pfa.state.secondary == orMapValueSeparator:
		pfa.state.setPrimary(expectSpaceOrMapKey)
	case pfa.isTopStackItemOfType(parseStackItemMap) && currRune == pfa.getTopStackItem().delimiter:
		needFinishTopStackItem = true
	default:
		err = pfaErrorMake(pfa, unexpectedRuneError)
	}
	return
}

func _expectWord(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItem()
	currRune, _ := pfa.currRune()
	switch {
	case unicode.IsLetter(currRune) || currRune == '_' || unicode.IsDigit(currRune):
		stackItem.itemString = stackItem.itemString + string(currRune)
	default:
		pfa.pushRune()
		needFinishTopStackItem = true
	}
	return
}

func _expectEndOfQwItem(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemQwItem)
	currRune, isEOF := pfa.currRune()
	switch {
	case isEOF:
		err = pfaErrorMake(pfa, unexpectedRuneError)
	case unicode.IsSpace(currRune) || currRune == stackItem.delimiter:
		pfa.pushRune()
		needFinishTopStackItem = true
	default:
		stackItem.itemString += string(currRune)
	}
	return
}

func _expectDigit(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemNumber)
	currRune, _ := pfa.currRune()
	switch pfa.state.secondary {
	case noSecondaryState:
		switch {
		case unicode.IsDigit(currRune):
			stackItem.itemString = stackItem.itemString + string(currRune)
			pfa.state.secondary = orUnderscoreOrDot
		default:
			err = pfaErrorMake(pfa, unexpectedRuneError)
		}
	case orUnderscoreOrDot:
		switch {
		case currRune == '.':
			pfa.state.secondary = orUnderscore
			stackItem.itemString = stackItem.itemString + string(currRune)
		case unicode.IsDigit(currRune) || currRune == '_':
			stackItem.itemString = stackItem.itemString + string(currRune)
		default:
			pfa.pushRune()
			needFinishTopStackItem = true
		}
	case orUnderscore:
		switch {
		case unicode.IsDigit(currRune) || currRune == '_':
			stackItem.itemString = stackItem.itemString + string(currRune)
		default:
			pfa.pushRune()
			needFinishTopStackItem = true
		}
	}
	return
}

func _expectContentOf(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItem()
	currRune, isEOF := pfa.currRune()
	switch {
	case isEOF:
		err = pfaErrorMake(pfa, unexpectedRuneError)
	case currRune == stackItem.delimiter:
		needFinishTopStackItem = true
	case currRune == '\\':
		pfa.state.primary = expectEscapedContentOf
	default:
		stackItem.itemString = stackItem.itemString + string(currRune)
	}
	return
}

var singleQuotedEscapedContent = map[rune]string{
	'"':  "\"",
	'\'': "'",
	'\\': "\\",
}

var doubleQuotedEscapedContent = map[rune]string{
	'"':  "\"",
	'\'': "'",
	'\\': "\\",
	'a':  "\a",
	'b':  "\b",
	'f':  "\f",
	'n':  "\n",
	'r':  "\r",
	't':  "\t",
	'v':  "\v",
}

// func getRunesOfEscapedContent() (result []rune) {
// 	result = []rune{}
// 	for r, _ := range escapedContent {
// 		result = append(result, r)
// 	}
// 	return
// }

func _expectEscapedContentOf(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItem()
	currRune, _ := pfa.currRune()
	if stackItem.delimiter == '\'' {
		if actualVal, ok := singleQuotedEscapedContent[currRune]; ok {
			stackItem.itemString = stackItem.itemString + actualVal
			pfa.state.primary = expectContentOf
		} else {
			err = pfaErrorMake(pfa, unexpectedRuneError)
		}
	} else {
		if actualVal, ok := doubleQuotedEscapedContent[currRune]; ok {
			stackItem.itemString = stackItem.itemString + actualVal
			pfa.state.primary = expectContentOf
		} else {
			err = pfaErrorMake(pfa, unexpectedRuneError)
		}
	}
	return
}

func _expectSpaceOrQwItemOrDelimiter(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItem()
	currRune, isEOF := pfa.currRune()
	switch {
	case isEOF:
		err = pfaErrorMake(pfa, unexpectedRuneError)
	case unicode.IsSpace(currRune):
	case currRune == stackItem.delimiter:
		needFinishTopStackItem = true
	default:
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemQwItem, start: pfa.curr, itemString: string(*pfa.curr.runePtr), delimiter: stackItem.delimiter})
		pfa.state.setPrimary(expectEndOfQwItem)
	}
	return
}
