package defparser

import (
	"fmt"
	"github.com/baza-winner/bw/core"
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
)

//go:generate stringer -type=ParseState

type ParseStackItemType uint16

const (
	parseStackItemArray ParseStackItemType = iota
	parseStackItemQw
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
	currentKey string
	itemString string
	value      interface{}
}

func getPosTitle(pos int) (posTitle string) {
	if pos < 0 {
		posTitle = "end of source"
	} else {
		posTitle = fmt.Sprintf("<ansiOutline>pos <ansiSecondaryLiteral>%d<ansi>", pos)
	}
	return
}

func getTopStackItem(stack []ParseStackItem, itemType ParseStackItemType, pos int, state ParseState) (stackItem *ParseStackItem, err error) {
	if !(len(stack) >= 1 && stack[len(stack)-1].itemType == itemType) {
		return nil, core.Error("<ansiOutline>stack<ansi> (<ansiSecondaryLiteral>%v<ansi>) expects to have top item of type <ansiPrimaryLiteral>%s<ansi> while at "+getPosTitle(pos)+" and <ansiOutline>state <ansiSecondaryLiteral>%s", stack, itemType, state)
	}
	stackItem = &stack[len(stack)-1]
	return
}

func Parse(source string) (result interface{}, err error) {
	state := expectSpaceOrValue
	stack := []ParseStackItem{}
	var stackItem *ParseStackItem
	var pos int
	var char rune
	var wasState ParseState
	for pos, char = range source {
		wasState = state
		switch state {
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
				if stackItem, err = getTopStackItem(stack, parseStackItemMap, pos, state); err != nil {
					return nil, err
				} else if err = finishTopStackItem(&stack, &state); err != nil {
					return nil, err
				}
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
				if stackItem, err = getTopStackItem(stack, parseStackItemArray, pos, state); err != nil {
					return nil, err
				} else if err = finishTopStackItem(&stack, &state); err != nil {
					return nil, err
				}
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
			if stackItem, err = getTopStackItem(stack, parseStackItemWord, pos, state); err != nil {
				return nil, err
			}
			if unicode.IsLetter(char) {
				stackItem.itemString = stackItem.itemString + string(char)
			} else if err = finishTopStackItem(&stack, &state); err != nil {
				return nil, err
			}

		case expectKey:
			if stackItem, err = getTopStackItem(stack, parseStackItemKey, pos, state); err != nil {
				return nil, err
			}
			if unicode.IsLetter(char) {
				stackItem.itemString = stackItem.itemString + string(char)
			} else if err = finishTopStackItem(&stack, &state); err != nil {
				return nil, err
			}

		case expectDigit, expectDigitOrUnderscore, expectDigitOrUnderscoreOrDot:
			if stackItem, err = getTopStackItem(stack, parseStackItemNumber, pos, state); err != nil {
				return nil, err
			}
			if state == expectDigitOrUnderscoreOrDot {
				switch {
				case char == '.':
					state = expectDigitOrUnderscore
					fallthrough
				case unicode.IsNumber(char) || char == '_':
					stackItem.itemString = stackItem.itemString + string(char)
				default:
					if err = finishTopStackItem(&stack, &state); err != nil {
						return nil, err
					}
				}
			} else {
				if unicode.IsDigit(char) || (state == expectDigitOrUnderscore) && char == '_' {
					stackItem.itemString = stackItem.itemString + string(char)
					if state == expectDigit {
						state = expectDigitOrUnderscoreOrDot
					}
				} else {
					state = unexpectedChar
					// return nil, core.Error("unexpected <ansiOutline>char <ansiPrimaryLiteral>%+q<ansi> (code <ansiSecondaryLiteral>%v<ansi>) at <ansiOutline>pos <ansiSecondaryLiteral>%d while <ansiOutline>state <ansiSecondaryLiteral>%s", char, char, pos, state)
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
				if stackItem, err = getTopStackItem(stack, parseStackItemString, pos, state); err != nil {
					return nil, err
				}
			default:
				if stackItem, err = getTopStackItem(stack, parseStackItemKey, pos, state); err != nil {
					return nil, err
				}
			}
			switch state {
			case
				expectDoubleQuotedKeyContent,
				expectSingleQuotedKeyContent,
				expectDoubleQuotedStringContent,
				expectSingleQuotedStringContent:
				if (state == expectDoubleQuotedStringContent || state == expectDoubleQuotedKeyContent) && char == '"' ||
					(state == expectSingleQuotedStringContent || state == expectSingleQuotedKeyContent) && char == '\'' {
					if err = finishTopStackItem(&stack, &state); err != nil {
						return nil, err
					}
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
			return nil, core.Error("unexpected <ansiOutline>state<ansi> <ansiPrimaryLiteral>%s<ansi> while at pos <ansiSecondaryLiteral>%d", state, pos)
		}

		if state == tokenFinished && len(stack) > 1 {
			var stackSubItem ParseStackItem
			stackSubItem, stack = stack[len(stack)-1], stack[:len(stack)-1] // https://github.com/golang/go/wiki/SliceTricks
			stackItem = &stack[len(stack)-1]
			switch stackItem.itemType {
			case parseStackItemArray:
				stackItem.itemArray = append(stackItem.itemArray, stackSubItem.value)
				switch stackSubItem.itemType {
				case parseStackItemNumber, parseStackItemWord:
					switch {
					case unicode.IsSpace(char):
						state = expectArrayItemSeparatorOrSpaceOrArrayValue
					case char == ',':
						state = expectSpaceOrArrayItem
					default:
						state = unexpectedChar
					}
				default:
					state = expectArrayItemSeparatorOrSpaceOrArrayValue
				}
			case parseStackItemMap:
				switch stackSubItem.itemType {
				case parseStackItemKey:
					stackItem.currentKey = stackSubItem.itemString
					state = expectMapKeySeparatorOrSpace
				default:
					stackItem.itemMap[stackItem.currentKey] = stackSubItem.value
					switch stackSubItem.itemType {
					case parseStackItemNumber, parseStackItemWord:
						switch {
						case unicode.IsSpace(char):
							state = expectMapValueSeparatorOrSpaceOrMapKey
						case char == ',':
							state = expectSpaceOrMapKey
						default:
							state = unexpectedChar
						}
					default:
						state = expectMapValueSeparatorOrSpaceOrMapKey
					}
				}
			default:
				return nil, core.Error("unexpected <ansiOutline>stackItem <ansiSecondaryLiteral>%v<ansi> can not have subitem %v", stackItem, stackSubItem)
			}
		}
		if state == unexpectedChar {
			break
		}
	}
	if state == unexpectedChar {
		return nil, core.Error("unexpected <ansiOutline>char <ansiPrimaryLiteral>%+q<ansi> (code <ansiSecondaryLiteral>%v<ansi>) at <ansiOutline>pos <ansiSecondaryLiteral>%d while <ansiOutline>state <ansiSecondaryLiteral>%s", char, char, pos, wasState)
	}

	switch len(stack) {
	case 0:
		return nil, nil
	case 1:
		switch state {
		case expectDigitOrUnderscore, expectDigitOrUnderscoreOrDot, expectWord:
			if err = finishTopStackItem(&stack, &state); err != nil {
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
		return nil, core.Error("<ansiOutline>stack<ansi> (<ansiSecondaryLiteral>%v<ansi>) expects to have one item at end of source and while <ansiOutline>state<ansi> <ansiSecondaryLiteral>%s", stack, state)
	}

	return
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

var underscoreRegexp = regexp.MustCompile("[_]+")

func finishTopStackItem(stackItem *ParseStackItem, state *ParseState) (err error) {
	switch stackItem.itemType {

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
		// case "qw":
		default:
			err = core.Error("unexpected word <ansiPrimaryLiteral>%s<ansi> at "+getPosTitle(stackItem.pos), stackItem.itemString)
		}

	case parseStackItemArray:
		stackItem.value = stackItem.itemArray

	case parseStackItemMap:
		stackItem.value = stackItem.itemMap

	case parseStackItemKey:

	default:
		err = core.Error("can not finish item of type <ansiPrimaryLiteral>%s<ansi> at "+getPosTitle(stackItem.pos), stackItem.itemType)
	}
	return
}
