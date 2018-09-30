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
	stack  parseStack
	state  parseState
	result interface{}
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

func (pfa *pfaStruct) processCharAtPos(pos int, charPtr *rune) (err error) {
	var stackItem *parseStackItem
	needFinishTopStackItem := false
	switch pfa.state {

	case expectEOF, expectEOFOrSpace:
		if charPtr == nil {
			pfa.state = expectEOF
		} else if !unicode.IsSpace(*charPtr) {
			return unexpectedCharError{}
		}

	case
		expectValueOrSpace,
		expectValueOrSpaceOrMapKeySeparator,
		expectValueOrSpaceOrArrayItemSeparator:
		if charPtr == nil && len(pfa.stack) == 0 {
			pfa.state = expectEOF
		} else {
			switch {
			case pfa.state == expectValueOrSpaceOrMapKeySeparator && *charPtr == '=':
				pfa.state = expectRocket
			case pfa.state == expectValueOrSpaceOrMapKeySeparator && *charPtr == ':':
				pfa.state = expectValueOrSpace
			case *charPtr == ',' && pfa.state == expectValueOrSpaceOrArrayItemSeparator:
				pfa.state = expectValueOrSpace

			case unicode.IsSpace(*charPtr):
			case *charPtr == '{':
				pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemMap, pos: pos, itemMap: map[string]interface{}{}})
				pfa.state = expectSpaceOrMapKey
			case *charPtr == '[':
				pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemArray, pos: pos, itemArray: []interface{}{}})
				pfa.state = expectValueOrSpace

			// case *charPtr == ']' && (pfa.state == expectSpaceOrArrayItem || pfa.state == expectArrayItemSeparatorOrSpaceOrArrayValue):
			case *charPtr == ']' && (len(pfa.stack) > 0 && pfa.stack[len(pfa.stack)-1].itemType == parseStackItemArray):
				stackItem = pfa.getTopStackItem(parseStackItemArray, pos)
				needFinishTopStackItem = true

			case *charPtr == '-' || *charPtr == '+' || unicode.IsDigit(*charPtr):
				pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemNumber, pos: pos, itemString: string(*charPtr)})
				if unicode.IsDigit(*charPtr) {
					pfa.state = expectDigitOrUnderscoreOrDot
				} else {
					pfa.state = expectDigit
				}
			case *charPtr == '"' || *charPtr == '\'':
				pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemString, pos: pos + 1, itemString: ``})
				if *charPtr == '"' {
					pfa.state = expectContentOfDoubleQuotedString
				} else {
					pfa.state = expectContentOfSingleQuotedString
				}
			case unicode.IsLetter(*charPtr):
				pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemWord, pos: pos, itemString: string(*charPtr)})
				pfa.state = expectWord
			default:
				return unexpectedCharError{}
			}
		}

	case expectArrayItemSeparatorOrSpace:
		switch {
		case charPtr == nil:
			return unexpectedCharError{}
		case unicode.IsSpace(*charPtr):
			pfa.state = expectValueOrSpaceOrArrayItemSeparator
		case *charPtr == ',':
			pfa.state = expectValueOrSpace
		default:
			return unexpectedCharError{}
		}

	case expectMapKeySeparatorOrSpace:
		switch {
		case charPtr == nil:
			return unexpectedCharError{}
		case unicode.IsSpace(*charPtr):
			pfa.state = expectValueOrSpaceOrMapKeySeparator
		case *charPtr == ':':
			pfa.state = expectValueOrSpace
		case *charPtr == '=':
			pfa.state = expectRocket
		default:
			return unexpectedCharError{}
		}

	case expectRocket:
		if charPtr != nil && *charPtr == '>' {
			pfa.state = expectValueOrSpace
		} else {
			return unexpectedCharError{}
		}

	case expectSpaceOrMapKey, expectSpaceOrMapKeyOrMapValueSeparator:
		switch {
		case charPtr == nil:
			return unexpectedCharError{}
		case unicode.IsSpace(*charPtr):
		case unicode.IsLetter(*charPtr):
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, pos: pos, itemString: string(*charPtr)})
			pfa.state = expectMapKey
		case *charPtr == '"':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, pos: pos, itemString: ``})
			pfa.state = expectContentOfDoubleQuotedKey
		case *charPtr == '\'':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, pos: pos, itemString: ``})
			pfa.state = expectContentOfSingleQuotedKey
		case *charPtr == ',' && pfa.state == expectSpaceOrMapKeyOrMapValueSeparator:
			pfa.state = expectSpaceOrMapKey
		case *charPtr == '}':
			stackItem = pfa.getTopStackItem(parseStackItemMap, pos)
			needFinishTopStackItem = true
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
			needFinishTopStackItem = true
		}

	case expectWord:
		stackItem = pfa.getTopStackItem(parseStackItemWord, pos)
		switch {
		case charPtr != nil && unicode.IsLetter(*charPtr):
			stackItem.itemString = stackItem.itemString + string(*charPtr)
		default:
			needFinishTopStackItem = true
		}

	case expectDigit, expectDigitOrUnderscore, expectDigitOrUnderscoreOrDot:
		stackItem = pfa.getTopStackItem(parseStackItemNumber, pos)
		if pfa.state == expectDigit {
			switch {
			case charPtr != nil && unicode.IsDigit(*charPtr):
				stackItem.itemString = stackItem.itemString + string(*charPtr)
				pfa.state = expectDigitOrUnderscoreOrDot
			default:
				return unexpectedCharError{}
			}
		} else {
			switch {
			case pfa.state == expectDigitOrUnderscoreOrDot && charPtr != nil && *charPtr == '.':
				pfa.state = expectDigitOrUnderscore
				fallthrough
			case charPtr != nil && (unicode.IsDigit(*charPtr) || *charPtr == '_'):
				stackItem.itemString = stackItem.itemString + string(*charPtr)
			default:
				needFinishTopStackItem = true
			}
		}

	case
		expectContentOfDoubleQuotedKey,
		expectContentOfSingleQuotedKey,
		expectEscapedContentOfDoubleQuotedKey,
		expectEscapedContentOfSingleQuotedKey,
		expectContentOfDoubleQuotedString,
		expectEscapedContentOfDoubleQuotedString,
		expectContentOfSingleQuotedString,
		expectEscapedContentOfSingleQuotedString:
		if charPtr == nil {
			return unexpectedCharError{}
		}
		switch pfa.state {
		case expectContentOfDoubleQuotedString,
			expectEscapedContentOfDoubleQuotedString,
			expectContentOfSingleQuotedString,
			expectEscapedContentOfSingleQuotedString:
			stackItem = pfa.getTopStackItem(parseStackItemString, pos)
		default:
			stackItem = pfa.getTopStackItem(parseStackItemKey, pos)
		}
		switch pfa.state {
		case
			expectContentOfDoubleQuotedKey,
			expectContentOfSingleQuotedKey,
			expectContentOfDoubleQuotedString,
			expectContentOfSingleQuotedString:
			if (pfa.state == expectContentOfDoubleQuotedString || pfa.state == expectContentOfDoubleQuotedKey) && *charPtr == '"' ||
				(pfa.state == expectContentOfSingleQuotedString || pfa.state == expectContentOfSingleQuotedKey) && *charPtr == '\'' {
				needFinishTopStackItem = true
			} else if *charPtr == '\\' {
				switch pfa.state {
				case expectContentOfDoubleQuotedString:
					pfa.state = expectEscapedContentOfDoubleQuotedString
				case expectContentOfSingleQuotedString:
					pfa.state = expectEscapedContentOfSingleQuotedString
				case expectContentOfDoubleQuotedKey:
					pfa.state = expectEscapedContentOfDoubleQuotedKey
				case expectContentOfSingleQuotedKey:
					pfa.state = expectEscapedContentOfSingleQuotedKey
				}
			} else {
				stackItem.itemString = stackItem.itemString + string(*charPtr)
			}
		case
			expectEscapedContentOfDoubleQuotedKey,
			expectEscapedContentOfSingleQuotedKey,
			expectEscapedContentOfDoubleQuotedString,
			expectEscapedContentOfSingleQuotedString:
			var actualVal string
			switch *charPtr {
			case '"':
				actualVal = "\""
			case '\'':
				actualVal = "'"
			case '\\':
				actualVal = "\\"
			default:
				switch pfa.state {
				case expectEscapedContentOfDoubleQuotedString, expectEscapedContentOfDoubleQuotedKey:
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
			switch pfa.state {
			case expectEscapedContentOfDoubleQuotedString:
				pfa.state = expectContentOfDoubleQuotedString
			case expectEscapedContentOfSingleQuotedString:
				pfa.state = expectContentOfSingleQuotedString
			case expectEscapedContentOfDoubleQuotedKey:
				pfa.state = expectContentOfDoubleQuotedKey
			case expectEscapedContentOfSingleQuotedKey:
				pfa.state = expectContentOfSingleQuotedKey
			}
		}

	case expectSpaceOrQwItemOrDelimiter:
		stackItem = pfa.getTopStackItem(parseStackItemQw, pos)
		switch {
		case charPtr == nil:
			return unexpectedCharError{}
		case unicode.IsSpace(*charPtr):
		case *charPtr == stackItem.delimiter:
			needFinishTopStackItem = true
		default:
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemQwItem, pos: pos, itemString: string(*charPtr), delimiter: stackItem.delimiter})
			pfa.state = expectEndOfQwItem
		}

	case expectEndOfQwItem:
		stackItem = pfa.getTopStackItem(parseStackItemQwItem, pos)
		switch {
		case charPtr == nil:
			return unexpectedCharError{}
		case unicode.IsSpace(*charPtr) || *charPtr == stackItem.delimiter:
			needFinishTopStackItem = true
		default:
			stackItem.itemString += string(*charPtr)
		}

	default:
		core.Panic("<ansiOutline>pfa<ansi> <ansiSecondaryLiteral>%s<ansi>", pfa)
	}

	if needFinishTopStackItem {
		if err = pfa.finishTopStackItem(charPtr); err != nil {
			return
		}
	}

	if charPtr == nil && !(pfa.state == expectEOF || pfa.state == expectEOFOrSpace) {
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
			pfa.state = expectSpaceOrQwItemOrDelimiter
		} else {
			if len(pfa.stack) < 2 {
				core.Panic("len(pfa.stack) < 2")
			}
			stackSubItem := pfa.stack[len(pfa.stack)-1]
			pfa.stack = pfa.stack[:len(pfa.stack)-1]
			stackItem = pfa.getTopStackItem(parseStackItemArray, stackSubItem.pos)
			stackItem.itemArray = append(stackItem.itemArray, stackSubItem.itemArray...)
			pfa.state = expectArrayItemSeparatorOrSpace
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
				pfa.state = expectSpaceOrQwItemOrDelimiter
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
				if pfa.state == expectSpaceOrQwItemOrDelimiter {
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
		pfa.state = expectEOFOrSpace
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
					pfa.state = expectValueOrSpaceOrArrayItemSeparator
				case *charPtr == ',':
					pfa.state = expectValueOrSpace
				default:
					return unexpectedCharError{}
				}
			default:
				pfa.state = expectValueOrSpaceOrArrayItemSeparator
			}

		case parseStackItemMap:
			switch stackSubItem.itemType {
			case parseStackItemKey:
				stackItem.currentKey = stackSubItem.itemString
				pfa.state = expectMapKeySeparatorOrSpace
			default:
				stackItem.itemMap[stackItem.currentKey] = stackSubItem.value
				switch stackSubItem.itemType {
				case parseStackItemNumber, parseStackItemWord:
					switch {
					case unicode.IsSpace(*charPtr):
						pfa.state = expectSpaceOrMapKeyOrMapValueSeparator
					case *charPtr == ',':
						pfa.state = expectSpaceOrMapKey
					default:
						return unexpectedCharError{}
					}
				default:
					pfa.state = expectSpaceOrMapKeyOrMapValueSeparator
				}
			}
		default:
			err = core.Error("<ansiOutline>stackItem <ansiSecondaryLiteral>%s<ansi> can not have subitem <ansiSecondaryLiteral>%s<ansiOutline>pfa.stack<ansiSecondaryLiteral>%s", stackItem, stackSubItem, pfa.stack)
		}
	}

	return
}
