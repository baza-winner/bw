package defparser

import (
	"encoding/json"
	"fmt"
	"github.com/baza-winner/bw/ansi"
	"github.com/baza-winner/bw/core"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type ParseState uint16

const (
	expectSpaceOrValue ParseState = iota
	expectSpaceOrArrayItem
	expectDigit
	expectDigitOrUnderscoreOrDot
	expectDigitOrUnderscore
	expectDoubleQuotedStringContent
	expectSingleQuotedStringContent
	expectDoubleQuotedStringEscapedContent
	expectSingleQuotedStringEscapedContent
	expectWord
	expectKey
	tokenFinished
	unexpectedChar
	expectArrayItemSeparatorOrSpaceOrArrayValue
	expectMapKeySeparatorOrSpace
	expectMapKeySeparatorOrSpaceOrMapValue
	expectSpaceOrMapValue
	expectRocket
	expectMapValueSeparatorOrSpaceOrMapValue
	expectSpaceOrMapKey
	expectMapValueSeparatorOrSpaceOrMapKey
	expectDoubleQuotedKeyContent
	expectSingleQuotedKeyContent
	expectDoubleQuotedKeyEscapedContent
	expectSingleQuotedKeyEscapedContent
	// expectQwStartDelimiter
	expectSpaceOrQwItemOrDelimiter
	expectEndOfQwItem
	expectArrayItemSeparatorOrSpace
)

//go:generate stringer -type=ParseState

type ParseStackItemType uint16

const (
	parseStackItemArray ParseStackItemType = iota
	parseStackItemQw
	parseStackItemQwItem
	parseStackItemMap
	parseStackItemNumber
	parseStackItemString
	parseStackItemWord
	parseStackItemKey
)

//go:generate stringer -type=ParseStackItemType

type ParseStackItem struct {
	itemType   ParseStackItemType
	pos        int
	itemArray  []interface{}
	itemMap    map[string]interface{}
	delimiter  rune
	currentKey string
	itemString string
	value      interface{}
}

func (self ParseStackItem) map4Json() (result map[string]interface{}) {
	result = map[string]interface{}{}
	result["itemType"] = self.itemType.String()
	result["pos"] = self.pos
	switch self.itemType {
	case parseStackItemArray:
		result["itemArray"] = self.itemArray
		result["value"] = self.value
	case parseStackItemQw:
		result["delimiter"] = string(self.delimiter)
		result["itemArray"] = self.itemArray
		result["value"] = self.value
	case parseStackItemQwItem:
		result["delimiter"] = string(self.delimiter)
		result["itemString"] = self.itemString
	case parseStackItemMap:
		result["itemMap"] = self.itemMap
		result["value"] = self.value
	case parseStackItemNumber, parseStackItemString, parseStackItemWord, parseStackItemKey:
		result["itemString"] = self.itemString
		result["value"] = self.value
	}
	return
}

func (self ParseStackItem) String() (result string) {
	bytes, _ := json.MarshalIndent(self.map4Json(), ``, `  `)
	result = string(bytes[:]) // https://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string/18615786#18615786
	return
}

func getPosTitle(pos int) (posTitle string) {
	if pos < 0 {
		posTitle = "end of source"
	} else {
		posTitle = fmt.Sprintf("<ansiOutline>pos <ansiSecondaryLiteral>%d<ansi>", pos)
	}
	return
}

func getTopStackItem(stack []ParseStackItem, itemType ParseStackItemType, pos int, state ParseState) (stackItem *ParseStackItem) {
	if !(len(stack) >= 1 && stack[len(stack)-1].itemType == itemType) {
		log.Panicf(ansi.Ansi(`Err`, "<ansiOutline>stack<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expects to have top item of type <ansiPrimaryLiteral>%s<ansi> while at "+getPosTitle(pos)+" and <ansiOutline>state <ansiSecondaryLiteral>%s"), stack, itemType, state)
	}
	stackItem = &stack[len(stack)-1]
	return
}

type ParseStack []ParseStackItem

func (self ParseStack) String() (result string) {
	array4Json := []map[string]interface{}{}
	for _, item := range self {
		array4Json = append(array4Json, item.map4Json())
	}
	bytes, _ := json.MarshalIndent(array4Json, ``, `  `)
	result = string(bytes[:]) // https://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string/18615786#18615786
	return
}

func Parse(source string) (result interface{}, err error) {
	state := expectSpaceOrValue
	stack := ParseStack{}
	var stackItem *ParseStackItem
	var pos int
	var char rune
	var wasState ParseState
	for pos, char = range source {

		needFinishTopStackItem := false
		wasState = state
		switch state {
		case expectSpaceOrQwItemOrDelimiter:
			stackItem = getTopStackItem(stack, parseStackItemQw, pos, state)
			switch {
			case unicode.IsSpace(char):
			case char == stackItem.delimiter:
				needFinishTopStackItem = true
			default:
				stack = append(stack, ParseStackItem{itemType: parseStackItemQwItem, pos: pos, itemString: string(char), delimiter: stackItem.delimiter})
				state = expectEndOfQwItem
			}
		case expectEndOfQwItem:
			stackItem = getTopStackItem(stack, parseStackItemQwItem, pos, state)
			switch {
			case unicode.IsSpace(char) || char == stackItem.delimiter:
				needFinishTopStackItem = true
			default:
				stackItem.itemString += string(char)
			}
		case expectArrayItemSeparatorOrSpace:
			switch {
			case unicode.IsSpace(char):
				state = expectArrayItemSeparatorOrSpaceOrArrayValue
			case char == ',':
				state = expectSpaceOrArrayItem
			default:
				state = unexpectedChar
			}
		case expectSpaceOrMapKey, expectMapValueSeparatorOrSpaceOrMapKey:
			switch {
			case unicode.IsSpace(char):
			case unicode.IsLetter(char):
				stack = append(stack, ParseStackItem{itemType: parseStackItemKey, pos: pos, itemString: string(char)})
				state = expectKey
			case char == '"':
				stack = append(stack, ParseStackItem{itemType: parseStackItemKey, pos: pos, itemString: ``})
				state = expectDoubleQuotedKeyContent
			case char == '\'':
				stack = append(stack, ParseStackItem{itemType: parseStackItemKey, pos: pos, itemString: ``})
				state = expectSingleQuotedKeyContent
			case char == ',' && state == expectMapValueSeparatorOrSpaceOrMapKey:
				state = expectSpaceOrMapKey
			case char == '}':
				stackItem = getTopStackItem(stack, parseStackItemMap, pos, state)
				needFinishTopStackItem = true
			default:
				state = unexpectedChar
			}
		case
			expectMapKeySeparatorOrSpaceOrMapValue,
			expectSpaceOrMapValue,
			expectSpaceOrValue,
			expectSpaceOrArrayItem,
			expectArrayItemSeparatorOrSpaceOrArrayValue:
			switch {
			case state == expectMapKeySeparatorOrSpaceOrMapValue && char == '=':
				state = expectRocket
			case state == expectMapKeySeparatorOrSpaceOrMapValue && char == ':':
				state = expectSpaceOrMapValue
			case unicode.IsSpace(char):
			case char == '{':
				stack = append(stack, ParseStackItem{itemType: parseStackItemMap, pos: pos, itemMap: map[string]interface{}{}})
				state = expectSpaceOrMapKey
			case char == '[':
				stack = append(stack, ParseStackItem{itemType: parseStackItemArray, pos: pos, itemArray: []interface{}{}})
				state = expectSpaceOrArrayItem
			case char == ']' && (state == expectSpaceOrArrayItem || state == expectArrayItemSeparatorOrSpaceOrArrayValue):
				stackItem = getTopStackItem(stack, parseStackItemArray, pos, state)
				needFinishTopStackItem = true
			case char == '-' || char == '+' || unicode.IsNumber(char):
				stack = append(stack, ParseStackItem{itemType: parseStackItemNumber, pos: pos, itemString: string(char)})
				if unicode.IsNumber(char) {
					state = expectDigitOrUnderscoreOrDot
				} else {
					state = expectDigit
				}
			case char == '"' || char == '\'':
				stack = append(stack, ParseStackItem{itemType: parseStackItemString, pos: pos + 1, itemString: ``})
				if char == '"' {
					state = expectDoubleQuotedStringContent
				} else {
					state = expectSingleQuotedStringContent
				}
			case unicode.IsLetter(char):
				stack = append(stack, ParseStackItem{itemType: parseStackItemWord, pos: pos, itemString: string(char)})
				state = expectWord
			case char == ',' && state == expectArrayItemSeparatorOrSpaceOrArrayValue:
				state = expectSpaceOrArrayItem
			default:
				state = unexpectedChar
			}

		case expectWord:
			stackItem = getTopStackItem(stack, parseStackItemWord, pos, state)
			if unicode.IsLetter(char) {
				stackItem.itemString = stackItem.itemString + string(char)
			} else {
				needFinishTopStackItem = true
			}

		case expectKey:
			stackItem = getTopStackItem(stack, parseStackItemKey, pos, state)
			if unicode.IsLetter(char) {
				stackItem.itemString = stackItem.itemString + string(char)
			} else {
				needFinishTopStackItem = true
			}

		case expectDigit, expectDigitOrUnderscore, expectDigitOrUnderscoreOrDot:
			stackItem = getTopStackItem(stack, parseStackItemNumber, pos, state)
			if state == expectDigitOrUnderscoreOrDot {
				switch {
				case char == '.':
					state = expectDigitOrUnderscore
					fallthrough
				case unicode.IsNumber(char) || char == '_':
					stackItem.itemString = stackItem.itemString + string(char)
				default:
					needFinishTopStackItem = true
				}
			} else {
				if unicode.IsDigit(char) || (state == expectDigitOrUnderscore) && char == '_' {
					stackItem.itemString = stackItem.itemString + string(char)
					if state == expectDigit {
						state = expectDigitOrUnderscoreOrDot
					}
				} else {
					state = unexpectedChar
				}
			}

		case
			expectDoubleQuotedKeyContent,
			expectSingleQuotedKeyContent,
			expectDoubleQuotedKeyEscapedContent,
			expectSingleQuotedKeyEscapedContent,
			expectDoubleQuotedStringContent,
			expectDoubleQuotedStringEscapedContent,
			expectSingleQuotedStringContent,
			expectSingleQuotedStringEscapedContent:
			switch state {
			case expectDoubleQuotedStringContent,
				expectDoubleQuotedStringEscapedContent,
				expectSingleQuotedStringContent,
				expectSingleQuotedStringEscapedContent:
				stackItem = getTopStackItem(stack, parseStackItemString, pos, state)
			default:
				stackItem = getTopStackItem(stack, parseStackItemKey, pos, state)
			}
			switch state {
			case
				expectDoubleQuotedKeyContent,
				expectSingleQuotedKeyContent,
				expectDoubleQuotedStringContent,
				expectSingleQuotedStringContent:
				if (state == expectDoubleQuotedStringContent || state == expectDoubleQuotedKeyContent) && char == '"' ||
					(state == expectSingleQuotedStringContent || state == expectSingleQuotedKeyContent) && char == '\'' {
					needFinishTopStackItem = true
				} else if char == '\\' {
					switch state {
					case expectDoubleQuotedStringContent:
						state = expectDoubleQuotedStringEscapedContent
					case expectSingleQuotedStringContent:
						state = expectSingleQuotedStringEscapedContent
					case expectDoubleQuotedKeyContent:
						state = expectDoubleQuotedKeyEscapedContent
					case expectSingleQuotedKeyContent:
						state = expectSingleQuotedKeyEscapedContent
					}
				} else {
					stackItem.itemString = stackItem.itemString + string(char)
				}
			case
				expectDoubleQuotedKeyEscapedContent,
				expectSingleQuotedKeyEscapedContent,
				expectDoubleQuotedStringEscapedContent,
				expectSingleQuotedStringEscapedContent:
				var actualVal string
				switch char {
				case '"':
					actualVal = "\""
				case '\'':
					actualVal = "'"
				case '\\':
					actualVal = "\\"
				default:
					switch state {
					case expectDoubleQuotedStringEscapedContent, expectDoubleQuotedKeyEscapedContent:
						switch char {
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
					state = unexpectedChar
				}
				stackItem.itemString = stackItem.itemString + actualVal
				switch state {
				case expectDoubleQuotedStringEscapedContent:
					state = expectDoubleQuotedStringContent
				case expectSingleQuotedStringEscapedContent:
					state = expectSingleQuotedStringContent
				case expectDoubleQuotedKeyEscapedContent:
					state = expectDoubleQuotedKeyContent
				case expectSingleQuotedKeyEscapedContent:
					state = expectSingleQuotedKeyContent
				}
			}

		case expectMapKeySeparatorOrSpace:
			switch {
			case unicode.IsSpace(char):
				state = expectMapKeySeparatorOrSpaceOrMapValue
			case char == ':':
				state = expectSpaceOrMapValue
			case char == '=':
				state = expectRocket
			default:
				state = unexpectedChar
			}

		case expectRocket:
			if char == '>' {
				state = expectSpaceOrMapValue
			} else {
				state = unexpectedChar
			}

		default:
			log.Panicf(ansi.Ansi(`Err`, "unexpected <ansiOutline>state <ansiSecondaryLiteral>%s<ansi> <ansiOutline>pos <ansiSecondaryLiteral>%d <ansiOutline>stack <ansiSecondaryLiteral>%s"), state, pos, stack)
		}

		if needFinishTopStackItem {
			log.Printf(ansi.Ansi(`Err`, `needFinishTopStackItem: <ansiOutline>stack <ansiSecondaryLiteral>%s`), stack)
			if stack, state, err = finishTopStackItem(stack, state, &char); err != nil {
				return nil, err
			}
		}

		if state == unexpectedChar {
			break
		}
	}
	if state == unexpectedChar {
		return nil, core.Error("unexpected <ansiOutline>char <ansiPrimaryLiteral>%+q<ansi> (code <ansiSecondaryLiteral>%v<ansi>) at <ansiOutline>pos <ansiSecondaryLiteral>%d<ansi> while <ansiOutline>state <ansiSecondaryLiteral>%s", char, char, pos, wasState)
	}

	switch len(stack) {
	case 0:
		return nil, nil
	case 1:
		switch state {
		case expectDigitOrUnderscore, expectDigitOrUnderscoreOrDot, expectWord:
			if stack, state, err = finishTopStackItem(stack, state, nil); err != nil {
				return nil, err
			}
		}
		switch state {
		case tokenFinished:
			result = stack[0].value
		default:
			return nil, core.Error("unexpected <ansiOutline>state<ansi> <ansiPrimaryLiteral>%s<ansi> while at end of source", state)
			return
		}
	default:
		return nil, core.Error("<ansiOutline>stack<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expects to have one item at end of source and while <ansiOutline>state<ansi> <ansiSecondaryLiteral>%s", stack, state)
	}

	return
}

func processCharAtPos(stack *ParseStack, pos int, char rune, state ParseState) (newState ParseState, err error) {
	return
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

var underscoreRegexp = regexp.MustCompile("[_]+")

func finishTopStackItem(stack ParseStack, state ParseState, charPtr *rune) (newStack ParseStack, newState ParseState, err error) {
	if len(stack) < 1 {
		log.Panic(`stack should have at least one item`)
	}
	newStack = stack
	newState = tokenFinished
	stackItem := &stack[len(stack)-1]
	switch stackItem.itemType {

	case parseStackItemQwItem:
		if len(stack) < 2 {
			log.Panicf("len(stack) < 2")
		}
		stackSubItem := stack[len(stack)-1]
		newStack = stack[:len(stack)-1]
		stackItem = getTopStackItem(newStack, parseStackItemQw, stackSubItem.pos, state)
		stackItem.itemArray = append(stackItem.itemArray, stackSubItem.itemString)
		if charPtr == nil {
			log.Panicf("charPtr == nil")
		}
		if unicode.IsSpace(*charPtr) {
			newState = expectSpaceOrQwItemOrDelimiter
		} else {
			stack = newStack
			if len(stack) < 2 {
				log.Panicf("len(stack) < 2")
			}
			stackSubItem := stack[len(stack)-1]
			newStack = stack[:len(stack)-1]
			stackItem = getTopStackItem(newStack, parseStackItemArray, stackSubItem.pos, state)
			stackItem.itemArray = append(stackItem.itemArray, stackSubItem.itemArray...)
			newState = expectArrayItemSeparatorOrSpace
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
			if len(stack) >= 2 && stack[len(stack)-2].itemType == parseStackItemArray && charPtr != nil {
				newState = expectSpaceOrQwItemOrDelimiter
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
						newState = unexpectedChar
					}
				}
				if newState == expectSpaceOrQwItemOrDelimiter {
					stackItem.itemType = parseStackItemQw
					stackItem.itemArray = []interface{}{}
				}
			} else {
				err = core.Error("unexpected word <ansiPrimaryLiteral>%s<ansi> in non array context at "+getPosTitle(stackItem.pos), stackItem.itemString)
			}
			return
		default:
			err = core.Error("unexpected word <ansiPrimaryLiteral>%s<ansi> at "+getPosTitle(stackItem.pos)+" <ansiOutline>stack <ansiSecondaryLiteral>%s", stackItem.itemString, stack)
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

	stack = newStack
	if len(stack) > 1 && charPtr != nil {
		var stackSubItem ParseStackItem
		stackSubItem, newStack = stack[len(stack)-1], stack[:len(stack)-1] // https://github.com/golang/go/wiki/SliceTricks
		stackItem = &newStack[len(newStack)-1]
		switch stackItem.itemType {
		case parseStackItemArray:
			stackItem.itemArray = append(stackItem.itemArray, stackSubItem.value)
			switch stackSubItem.itemType {
			case parseStackItemNumber, parseStackItemWord:
				switch {
				case unicode.IsSpace(*charPtr):
					newState = expectArrayItemSeparatorOrSpaceOrArrayValue
				case *charPtr == ',':
					newState = expectSpaceOrArrayItem
				default:
					newState = unexpectedChar
				}
			default:
				newState = expectArrayItemSeparatorOrSpaceOrArrayValue
			}

		case parseStackItemMap:
			switch stackSubItem.itemType {
			case parseStackItemKey:
				stackItem.currentKey = stackSubItem.itemString
				newState = expectMapKeySeparatorOrSpace
			default:
				stackItem.itemMap[stackItem.currentKey] = stackSubItem.value
				switch stackSubItem.itemType {
				case parseStackItemNumber, parseStackItemWord:
					switch {
					case unicode.IsSpace(*charPtr):
						newState = expectMapValueSeparatorOrSpaceOrMapKey
					case *charPtr == ',':
						newState = expectSpaceOrMapKey
					default:
						newState = unexpectedChar
					}
				default:
					newState = expectMapValueSeparatorOrSpaceOrMapKey
				}
			}
		default:
			err = core.Error("<ansiOutline>stackItem <ansiSecondaryLiteral>%s<ansi> can not have subitem <ansiSecondaryLiteral>%s<ansiOutline>stack<ansiSecondaryLiteral>%s", stackItem, stackSubItem, stack)
		}
	} else {
		charPtrTitle := ``
		if charPtr != nil {
			charPtrTitle = fmt.Sprintf(` <ansiOutline>*charPtr <ansiSecondaryLiteral>%v`, *charPtr)
		}
		log.Printf(ansi.Ansi(`Err`, `finishTopStackItem: <ansiOutline>stack <ansiSecondaryLiteral>%s <ansiOutline>charPtr <ansiSecondaryLiteral>%v`+charPtrTitle), stack, charPtr)
	}

	return
}