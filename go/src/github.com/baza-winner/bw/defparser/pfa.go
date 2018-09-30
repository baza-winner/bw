package defparser

import (
	"encoding/json"
	"fmt"
	"github.com/baza-winner/bw/core"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type pfaStruct struct {
	stack                  parseStack
	state                  parseState
	result                 interface{}
	needFinishTopStackItem bool
}

func (pfa *pfaStruct) getDataForJson() interface{} {
	result := map[string]interface{}{}
	result["stack"] = pfa.stack.getDataForJson()
	result["state"] = pfa.state.String()
	return result
}

func (pfa *pfaStruct) String() (result string) {
	bytes, _ := json.MarshalIndent(pfa.getDataForJson(), ``, `  `)
	result = string(bytes[:]) // https://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string/18615786#18615786
	return
}

type unexpectedCharError struct{}

func (e unexpectedCharError) Error() string {
	return `unexpected char`
}

func getPosTitle(pos int) (posTitle string) {
	if pos < 0 {
		posTitle = "end of source"
	} else {
		posTitle = fmt.Sprintf("<ansiOutline>pos <ansiSecondaryLiteral>%d<ansi>", pos)
	}
	return
}

func (pfa *pfaStruct) getTopStackItem(itemType parseStackItemType, pos int) (stackItem *parseStackItem) {
	if !(len(pfa.stack) >= 1 && pfa.stack[len(pfa.stack)-1].itemType == itemType) {
		core.Panic("<ansiOutline>stack<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expects to have top item of type <ansiPrimaryLiteral>%s<ansi> while at "+getPosTitle(pos)+" and <ansiOutline>state <ansiSecondaryLiteral>%s", pfa.stack, itemType, pfa.state)
	}
	stackItem = &pfa.stack[len(pfa.stack)-1]
	return
}

type pfaPrimaryStateMethod func(*pfaStruct, *rune, int) error

func _expectEOF(pfa *pfaStruct, charPtr *rune, pos int) error {
	if charPtr == nil {
		pfa.state.setPrimary(expectEOF)
	} else if pfa.state.secondary != orSpace || !unicode.IsSpace(*charPtr) {
		return unexpectedCharError{}
	}
	return nil
}

func _expectValueOrSpace(pfa *pfaStruct, charPtr *rune, pos int) error {
	if charPtr == nil && len(pfa.stack) == 0 {
		pfa.state.setPrimary(expectEOF)
	} else {
		switch {
		case *charPtr == '=' && pfa.state.secondary == orMapKeySeparator:
			pfa.state.setPrimary(expectRocket)
		case *charPtr == ':' && pfa.state.secondary == orMapKeySeparator:
			pfa.state.setPrimary(expectValueOrSpace)
		case *charPtr == ',' && pfa.state.secondary == orArrayItemSeparator:
			pfa.state.setPrimary(expectValueOrSpace)

		case unicode.IsSpace(*charPtr):
		case *charPtr == '{':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemMap, pos: pos, itemMap: map[string]interface{}{}})
			pfa.state.setPrimary(expectSpaceOrMapKey)
		case *charPtr == '[':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemArray, pos: pos, itemArray: []interface{}{}})
			pfa.state.setPrimary(expectValueOrSpace)

		case *charPtr == ']' && (len(pfa.stack) > 0 && pfa.stack[len(pfa.stack)-1].itemType == parseStackItemArray):
			_ = pfa.getTopStackItem(parseStackItemArray, pos)
			pfa.needFinishTopStackItem = true

		case *charPtr == '-' || *charPtr == '+' || unicode.IsDigit(*charPtr):
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemNumber, pos: pos, itemString: string(*charPtr)})
			if unicode.IsDigit(*charPtr) {
				pfa.state.setSecondary(expectDigit, orUnderscoreOrDot)
			} else {
				pfa.state.setPrimary(expectDigit)
			}
		case *charPtr == '"' || *charPtr == '\'':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemString, pos: pos + 1, itemString: ``})
			if *charPtr == '"' {
				pfa.state.setTertiary(expectContentOf, doubleQuoted, stringToken)
			} else {
				pfa.state.setTertiary(expectContentOf, singleQuoted, stringToken)
			}
		case unicode.IsLetter(*charPtr):
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemWord, pos: pos, itemString: string(*charPtr)})
			pfa.state.setPrimary(expectWord)
		default:
			return unexpectedCharError{}
		}
	}
	return nil
}

var pfaPrimaryStateMethods = map[parsePrimaryState]pfaPrimaryStateMethod{
	expectEOF:          _expectEOF,
	expectValueOrSpace: _expectValueOrSpace,
}

func (pfa *pfaStruct) processCharAtPos(pos int, charPtr *rune) (err error) {
	var stackItem *parseStackItem
	pfa.needFinishTopStackItem = false
	switch pfa.state.primary {

	case expectEOF:
		err = pfaPrimaryStateMethods[pfa.state.primary](pfa, charPtr, pos)
		// if charPtr == nil {
		// 	pfa.state.setPrimary(expectEOF)
		// } else if pfa.state.secondary != orSpace || !unicode.IsSpace(*charPtr) {
		// 	return unexpectedCharError{}
		// }

	case expectValueOrSpace:
		err = pfaPrimaryStateMethods[pfa.state.primary](pfa, charPtr, pos)
		// if charPtr == nil && len(pfa.stack) == 0 {
		// 	pfa.state.setPrimary(expectEOF)
		// } else {
		// 	switch {
		// 	case *charPtr == '=' && pfa.state.secondary == orMapKeySeparator:
		// 		pfa.state.setPrimary(expectRocket)
		// 	case *charPtr == ':' && pfa.state.secondary == orMapKeySeparator:
		// 		pfa.state.setPrimary(expectValueOrSpace)
		// 	case *charPtr == ',' && pfa.state.secondary == orArrayItemSeparator:
		// 		pfa.state.setPrimary(expectValueOrSpace)

		// 	case unicode.IsSpace(*charPtr):
		// 	case *charPtr == '{':
		// 		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemMap, pos: pos, itemMap: map[string]interface{}{}})
		// 		pfa.state.setPrimary(expectSpaceOrMapKey)
		// 	case *charPtr == '[':
		// 		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemArray, pos: pos, itemArray: []interface{}{}})
		// 		pfa.state.setPrimary(expectValueOrSpace)

		// 	case *charPtr == ']' && (len(pfa.stack) > 0 && pfa.stack[len(pfa.stack)-1].itemType == parseStackItemArray):
		// 		stackItem = pfa.getTopStackItem(parseStackItemArray, pos)
		// 		pfa.needFinishTopStackItem = true

		// 	case *charPtr == '-' || *charPtr == '+' || unicode.IsDigit(*charPtr):
		// 		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemNumber, pos: pos, itemString: string(*charPtr)})
		// 		if unicode.IsDigit(*charPtr) {
		// 			pfa.state.setSecondary(expectDigit, orUnderscoreOrDot)
		// 		} else {
		// 			pfa.state.setPrimary(expectDigit)
		// 		}
		// 	case *charPtr == '"' || *charPtr == '\'':
		// 		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemString, pos: pos + 1, itemString: ``})
		// 		if *charPtr == '"' {
		// 			pfa.state.setTertiary(expectContentOf, doubleQuoted, stringToken)
		// 		} else {
		// 			pfa.state.setTertiary(expectContentOf, singleQuoted, stringToken)
		// 		}
		// 	case unicode.IsLetter(*charPtr):
		// 		pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemWord, pos: pos, itemString: string(*charPtr)})
		// 		pfa.state.setPrimary(expectWord)
		// 	default:
		// 		return unexpectedCharError{}
		// 	}
		// }

	case expectArrayItemSeparatorOrSpace:
		switch {
		case charPtr == nil:
			return unexpectedCharError{}
		case unicode.IsSpace(*charPtr):
			pfa.state.setSecondary(expectValueOrSpace, orArrayItemSeparator)
		case *charPtr == ',':
			pfa.state.setPrimary(expectValueOrSpace)
		default:
			return unexpectedCharError{}
		}

	case expectMapKeySeparatorOrSpace:
		switch {
		case charPtr == nil:
			return unexpectedCharError{}
		case unicode.IsSpace(*charPtr):
			pfa.state.setSecondary(expectValueOrSpace, orMapKeySeparator)
		case *charPtr == ':':
			pfa.state.setPrimary(expectValueOrSpace)
		case *charPtr == '=':
			pfa.state.setPrimary(expectRocket)
		default:
			return unexpectedCharError{}
		}

	case expectRocket:
		if charPtr != nil && *charPtr == '>' {
			pfa.state.setPrimary(expectValueOrSpace)
		} else {
			return unexpectedCharError{}
		}

	case expectSpaceOrMapKey:
		switch {
		case charPtr == nil:
			return unexpectedCharError{}
		case unicode.IsSpace(*charPtr):
		case unicode.IsLetter(*charPtr):
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, pos: pos, itemString: string(*charPtr)})
			pfa.state.setPrimary(expectMapKey)
		case *charPtr == '"':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, pos: pos, itemString: ``})
			pfa.state.setTertiary(expectContentOf, doubleQuoted, keyToken)
		case *charPtr == '\'':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, pos: pos, itemString: ``})
			pfa.state.setTertiary(expectContentOf, singleQuoted, keyToken)
		case *charPtr == ',' && pfa.state.primary == expectSpaceOrMapKey && pfa.state.secondary == orMapValueSeparator:
			pfa.state.setPrimary(expectSpaceOrMapKey)
		case *charPtr == '}':
			stackItem = pfa.getTopStackItem(parseStackItemMap, pos)
			pfa.needFinishTopStackItem = true
		default:
			return unexpectedCharError{}
		}

	case expectMapKey:
		stackItem = pfa.getTopStackItem(parseStackItemKey, pos)
		switch {
		case charPtr == nil:
			return unexpectedCharError{}
		case unicode.IsLetter(*charPtr):
			stackItem.itemString = stackItem.itemString + string(*charPtr)
		default:
			pfa.needFinishTopStackItem = true
		}

	case expectWord:
		stackItem = pfa.getTopStackItem(parseStackItemWord, pos)
		switch {
		case charPtr != nil && unicode.IsLetter(*charPtr):
			stackItem.itemString = stackItem.itemString + string(*charPtr)
		default:
			pfa.needFinishTopStackItem = true
		}

	case expectDigit:
		stackItem = pfa.getTopStackItem(parseStackItemNumber, pos)
		if pfa.state.secondary == noSecondaryState {
			switch {
			case charPtr != nil && unicode.IsDigit(*charPtr):
				stackItem.itemString = stackItem.itemString + string(*charPtr)
				pfa.state.secondary = orUnderscoreOrDot
			default:
				return unexpectedCharError{}
			}
		} else {
			switch {
			case pfa.state.secondary == orUnderscoreOrDot && charPtr != nil && *charPtr == '.':
				pfa.state.secondary = orUnderscore
				fallthrough
			case charPtr != nil && (unicode.IsDigit(*charPtr) || *charPtr == '_'):
				stackItem.itemString = stackItem.itemString + string(*charPtr)
			default:
				pfa.needFinishTopStackItem = true
			}
		}

	case expectContentOf:
		if charPtr == nil {
			return unexpectedCharError{}
		}
		itemType := parseStackItemString
		if pfa.state.tertiary == keyToken {
			itemType = parseStackItemKey
		}
		stackItem = pfa.getTopStackItem(itemType, pos)
		if pfa.state.secondary == doubleQuoted && *charPtr == '"' ||
			pfa.state.secondary == singleQuoted && *charPtr == '\'' {
			pfa.needFinishTopStackItem = true
		} else if *charPtr == '\\' {
			pfa.state.primary = expectEscapedContentOf
		} else {
			stackItem.itemString = stackItem.itemString + string(*charPtr)
		}

	case expectEscapedContentOf:
		if charPtr == nil {
			return unexpectedCharError{}
		}
		itemType := parseStackItemString
		if pfa.state.tertiary == keyToken {
			itemType = parseStackItemKey
		}
		stackItem = pfa.getTopStackItem(itemType, pos)
		var actualVal string
		switch *charPtr {
		case '"':
			actualVal = "\""
		case '\'':
			actualVal = "'"
		case '\\':
			actualVal = "\\"
		default:
			if pfa.state.secondary == doubleQuoted {
				switch *charPtr {
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
			return unexpectedCharError{}
		}
		stackItem.itemString = stackItem.itemString + actualVal
		pfa.state.primary = expectContentOf

	case expectSpaceOrQwItemOrDelimiter:
		stackItem = pfa.getTopStackItem(parseStackItemQw, pos)
		switch {
		case charPtr == nil:
			return unexpectedCharError{}
		case unicode.IsSpace(*charPtr):
		case *charPtr == stackItem.delimiter:
			pfa.needFinishTopStackItem = true
		default:
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemQwItem, pos: pos, itemString: string(*charPtr), delimiter: stackItem.delimiter})
			pfa.state.setPrimary(expectEndOfQwItem)
		}

	case expectEndOfQwItem:
		stackItem = pfa.getTopStackItem(parseStackItemQwItem, pos)
		switch {
		case charPtr == nil:
			return unexpectedCharError{}
		case unicode.IsSpace(*charPtr) || *charPtr == stackItem.delimiter:
			pfa.needFinishTopStackItem = true
		default:
			stackItem.itemString += string(*charPtr)
		}

	default:
		core.Panic("<ansiOutline>pfa<ansi> <ansiSecondaryLiteral>%s<ansi>", pfa)
	}

	if pfa.needFinishTopStackItem {
		if err = pfa.finishTopStackItem(charPtr); err != nil {
			return
		}
	}

	if charPtr == nil && pfa.state.primary != expectEOF {
		core.Panic("<ansiOutline>pfa<ansi> <ansiSecondaryLiteral>%s<ansi>", pfa)
	}

	return
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

var underscoreRegexp = regexp.MustCompile("[_]+")

func (pfa *pfaStruct) finishTopStackItem(charPtr *rune) (err error) {
	if len(pfa.stack) < 1 {
		core.Panic("<ansiOutline>pfa<ansi> <ansiSecondaryLiteral>%s<ansi>", pfa)
	}
	stackItem := &pfa.stack[len(pfa.stack)-1]
	switch stackItem.itemType {

	case parseStackItemQwItem:
		if len(pfa.stack) < 2 {
			core.Panic("len(pfa.stack) < 2")
		}
		stackSubItem := pfa.stack[len(pfa.stack)-1]
		pfa.stack = pfa.stack[:len(pfa.stack)-1]
		stackItem = pfa.getTopStackItem(parseStackItemQw, stackSubItem.pos)
		stackItem.itemArray = append(stackItem.itemArray, stackSubItem.itemString)
		if charPtr == nil {
			core.Panic("charPtr == nil")
		}
		if unicode.IsSpace(*charPtr) {
			pfa.state.setPrimary(expectSpaceOrQwItemOrDelimiter)
		} else {
			if len(pfa.stack) < 2 {
				core.Panic("len(pfa.stack) < 2")
			}
			stackSubItem := pfa.stack[len(pfa.stack)-1]
			pfa.stack = pfa.stack[:len(pfa.stack)-1]
			stackItem = pfa.getTopStackItem(parseStackItemArray, stackSubItem.pos)
			stackItem.itemArray = append(stackItem.itemArray, stackSubItem.itemArray...)
			pfa.state.setPrimary(expectArrayItemSeparatorOrSpace)
		}
		return

	case parseStackItemNumber:
		source := underscoreRegexp.ReplaceAllLiteralString(stackItem.itemString, ``)
		if strings.Contains(source, `.`) {
			var float64Val float64
			if float64Val, err = strconv.ParseFloat(source, 64); err == nil {
				stackItem.value = float64Val
			}
		} else {
			var int64Val int64
			if int64Val, err = strconv.ParseInt(source, 10, 64); err == nil {
				if int64(MinInt) <= int64Val && int64Val <= int64(MaxInt) {
					stackItem.value = int(int64Val)
				} else {
					stackItem.value = int64Val
				}
			}
		}
		if err != nil {
			err = core.Error("failed to get number from string <ansiPrimaryLiteral>%s<ansi> at "+getPosTitle(stackItem.pos)+": %v", stackItem.itemString, err)
		}

	case parseStackItemString:
		stackItem.value = stackItem.itemString

	case parseStackItemWord:
		switch stackItem.itemString {
		case "true":
			stackItem.value = true
		case "false":
			stackItem.value = false
		case "nil":
			stackItem.value = nil
		case "qw":
			if len(pfa.stack) >= 2 && pfa.stack[len(pfa.stack)-2].itemType == parseStackItemArray && charPtr != nil {
				pfa.state.setPrimary(expectSpaceOrQwItemOrDelimiter)
				switch *charPtr {
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
					case unicode.IsPunct(*charPtr) || unicode.IsSymbol(*charPtr):
						stackItem.delimiter = *charPtr
					default:
						return unexpectedCharError{}
					}
				}
				if pfa.state.primary == expectSpaceOrQwItemOrDelimiter {
					stackItem.itemType = parseStackItemQw
					stackItem.itemArray = []interface{}{}
				}
			} else {
				err = core.Error("unexpected word <ansiPrimaryLiteral>%s<ansi> in non array context at "+getPosTitle(stackItem.pos), stackItem.itemString)
			}

			return
		default:
			err = core.Error("unexpected word <ansiPrimaryLiteral>%s<ansi> at "+getPosTitle(stackItem.pos)+" <ansiOutline>pfa.stack <ansiSecondaryLiteral>%s", stackItem.itemString, pfa.stack)
		}

	case parseStackItemArray:
		stackItem.value = stackItem.itemArray

	case parseStackItemMap:
		stackItem.value = stackItem.itemMap

	case parseStackItemKey:

	default:
		err = core.Error("can not finish item of type <ansiPrimaryLiteral>%s<ansi> at "+getPosTitle(stackItem.pos), stackItem.itemType)
	}
	if err != nil {
		return
	}

	if len(pfa.stack) == 1 {
		pfa.result = stackItem.value
		pfa.state.setSecondary(expectEOF, orSpace)
		return
	} else if len(pfa.stack) > 1 && charPtr != nil {
		var stackSubItem parseStackItem
		stackSubItem, pfa.stack = pfa.stack[len(pfa.stack)-1], pfa.stack[:len(pfa.stack)-1] // https://github.com/golang/go/wiki/SliceTricks
		stackItem = &pfa.stack[len(pfa.stack)-1]
		switch stackItem.itemType {
		case parseStackItemArray:
			stackItem.itemArray = append(stackItem.itemArray, stackSubItem.value)
			switch stackSubItem.itemType {
			case parseStackItemNumber, parseStackItemWord:
				switch {
				case unicode.IsSpace(*charPtr):
					pfa.state.setSecondary(expectValueOrSpace, orArrayItemSeparator)
				case *charPtr == ',':
					pfa.state.setPrimary(expectValueOrSpace)
				default:
					return unexpectedCharError{}
				}
			default:
				pfa.state.setSecondary(expectValueOrSpace, orArrayItemSeparator)
			}

		case parseStackItemMap:
			switch stackSubItem.itemType {
			case parseStackItemKey:
				stackItem.currentKey = stackSubItem.itemString
				pfa.state.setPrimary(expectMapKeySeparatorOrSpace)
			default:
				stackItem.itemMap[stackItem.currentKey] = stackSubItem.value
				switch stackSubItem.itemType {
				case parseStackItemNumber, parseStackItemWord:
					switch {
					case unicode.IsSpace(*charPtr):
						pfa.state.setSecondary(expectSpaceOrMapKey, orMapValueSeparator)
					case *charPtr == ',':
						pfa.state.setPrimary(expectSpaceOrMapKey)
					default:
						return unexpectedCharError{}
					}
				default:
					pfa.state.setSecondary(expectSpaceOrMapKey, orMapValueSeparator)
				}
			}
		default:
			err = core.Error("<ansiOutline>stackItem <ansiSecondaryLiteral>%s<ansi> can not have subitem <ansiSecondaryLiteral>%s<ansiOutline>pfa.stack<ansiSecondaryLiteral>%s", stackItem, stackSubItem, pfa.stack)
		}
	}

	return
}
