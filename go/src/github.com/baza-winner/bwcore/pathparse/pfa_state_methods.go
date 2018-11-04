package pathparse

import (
	"unicode"

	"github.com/baza-winner/bwcore/bwerr"
)

type pfaPrimaryStateMethod func(*pfaStruct) (bool, error)

var pfaPrimaryStateMethods = map[parsePrimaryState]pfaPrimaryStateMethod{
	expectEOF:              _expectEOF,
	expectPathSegment:      _expectPathSegment,
	expectMapKey:           _expectMapKey,
	expectDigit:            _expectDigit,
	expectContentOf:        _expectContentOf,
	expectEscapedContentOf: _expectEscapedContentOf,
}

func pfaPrimaryStateMethodsCheck() {
	expect := parsePrimaryStateBelow + 1
	for expect < parsePrimaryStateAbove {
		if _, ok := pfaPrimaryStateMethods[expect]; !ok {
			bwerr.Panic("not defined <ansiVar>pfaPrimaryStateMethods<ansi>[<ansiVal>%s<ansi>]", expect)
		}
		expect += 1
	}
}

func _expectEOF(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if pfa.curr.runePtr == nil {
		pfa.state.setPrimary(expectEOF)
	} else if pfa.state.secondary == orPathSegmentDelimiter && *pfa.curr.runePtr == '.' {
		pfa.state.setPrimary(expectPathSegment)
	} else {
		err = pfaErrorMake(pfa, unexpectedRuneError)
	}
	return
}

func _expectPathSegment(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	switch {
	case pfa.curr.runePtr == nil:
		err = pfaErrorMake(pfa, unexpectedRuneError)

	case *pfa.curr.runePtr == '#':
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemNumber, start: pfa.curr})
		pfa.state.setPrimary(expectDigit)

	case *pfa.curr.runePtr == '[':
		delimiter := ']'
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemNumber, start: pfa.curr, delimiter: &delimiter})
		pfa.state.setPrimary(expectDigit)

	case *pfa.curr.runePtr == '"' || *pfa.curr.runePtr == '\'':
		pfa.pullRune()
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemString, start: pfa.curr, itemString: ``})
		pfa.pushRune()
		quoted := singleQuoted
		if *pfa.curr.runePtr == '"' {
			quoted = doubleQuoted
		}
		pfa.state.setSecondary(expectContentOf, quoted)

	case unicode.IsLetter(*pfa.curr.runePtr) || *pfa.curr.runePtr == '_':
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, start: pfa.curr, itemString: string(*pfa.curr.runePtr)})
		pfa.state.setPrimary(expectMapKey)

	case unicode.IsDigit(*pfa.curr.runePtr):
		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemNumber, start: pfa.curr, itemString: string(*pfa.curr.runePtr)})
		pfa.state.setSecondary(expectDigit, orUnderscore)

	default:
		err = pfaErrorMake(pfa, unexpectedRuneError)
	}
	return
}

func _expectMapKey(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemKey)
	currRune := pfa.currRune()
	switch {
	case unicode.IsLetter(currRune) || currRune == '_' || unicode.IsDigit(currRune):
		stackItem.itemString = stackItem.itemString + string(currRune)
	default:
		pfa.pushRune()
		needFinishTopStackItem = true
	}
	return
}

func _expectDigit(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	stackItem := pfa.getTopStackItemOfType(parseStackItemNumber)
	currRune := pfa.currRune()
	if pfa.state.secondary == noSecondaryState {
		switch {
		case unicode.IsDigit(currRune):
			stackItem.itemString = stackItem.itemString + string(currRune)
			pfa.state.secondary = orUnderscore
		default:
			err = pfaErrorMake(pfa, unexpectedRuneError)
		}
	} else {
		switch {
		case unicode.IsDigit(currRune) || currRune == '_':
			stackItem.itemString = stackItem.itemString + string(*pfa.curr.runePtr)
		default:
			if stackItem.delimiter == nil {
				pfa.pushRune()
			} else if currRune != *stackItem.delimiter {
				return false, pfaErrorMake(pfa, unexpectedRuneError)
			}
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
