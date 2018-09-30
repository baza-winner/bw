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

type pfaStruct struct {
	stack ParseStack
	state ParseState
}

func (self ParseStack) String() (result string) {
	array4Json := []map[string]interface{}{}
	for _, item := range self {
		array4Json = append(array4Json, item.map4Json())
	}
	bytes, _ := json.MarshalIndent(array4Json, ``, `  `)
	result = string(bytes[:]) // https://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string/18615786#18615786
	return
}

type unexpectedCharError struct{}

func (e unexpectedCharError) Error() string {
	return `unexpected char`
}

func Parse(source string) (result interface{}, err error) {
	var pos int
	var char rune
	pfa := pfaStruct{stack: ParseStack{}, state: expectSpaceOrValue}
	for pos, char = range source {
		wasState := pfa.state
		if err = pfa.processCharAtPos(pos, char); err != nil {
			log.Printf(`%v`, err)
			if _, ok := err.(unexpectedCharError); ok {
				return nil, core.Error("unexpected <ansiOutline>char <ansiPrimaryLiteral>%+q<ansi> (code <ansiSecondaryLiteral>%v<ansi>) at <ansiOutline>pos <ansiSecondaryLiteral>%d<ansi> while <ansiOutline>state <ansiSecondaryLiteral>%s", char, char, pos, wasState)
			} else {
				log.Printf(`%v is NOT unexpectedCharError`, err)
			}
			return nil, err
		}
	}

	switch len(pfa.stack) {
	case 0:
		return nil, nil
	case 1:
		switch pfa.state {
		case expectDigitOrUnderscore, expectDigitOrUnderscoreOrDot, expectWord:
			if err = pfa.finishTopStackItem(nil); err != nil {
				return nil, err
			}
		}
		switch pfa.state {
		case tokenFinished:
			result = pfa.stack[0].value
		default:
			return nil, core.Error("unexpected <ansiOutline>state<ansi> <ansiPrimaryLiteral>%s<ansi> while at end of source", pfa.state)
			return
		}
	default:
		return nil, core.Error("<ansiOutline>stack<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expects to have one item at end of source and while <ansiOutline>state<ansi> <ansiSecondaryLiteral>%s", pfa.stack, pfa.state)
	}

	return
}

func (pfa *pfaStruct) processCharAtPos(pos int, char rune) (err error) {
	var stackItem *ParseStackItem
	needFinishTopStackItem := false
	switch pfa.state {
	case expectSpaceOrQwItemOrDelimiter:
		stackItem = getTopStackItem(pfa.stack, parseStackItemQw, pos, pfa.state)
		switch {
		case unicode.IsSpace(char):
		case char == stackItem.delimiter:
			needFinishTopStackItem = true
		default:
			pfa.stack = append(pfa.stack, ParseStackItem{itemType: parseStackItemQwItem, pos: pos, itemString: string(char), delimiter: stackItem.delimiter})
			pfa.state = expectEndOfQwItem
		}
	case expectEndOfQwItem:
		stackItem = getTopStackItem(pfa.stack, parseStackItemQwItem, pos, pfa.state)
		switch {
		case unicode.IsSpace(char) || char == stackItem.delimiter:
			needFinishTopStackItem = true
		default:
			stackItem.itemString += string(char)
		}
	case expectArrayItemSeparatorOrSpace:
		switch {
		case unicode.IsSpace(char):
			pfa.state = expectArrayItemSeparatorOrSpaceOrArrayValue
		case char == ',':
			pfa.state = expectSpaceOrArrayItem
		default:
			err = unexpectedCharError{}
			return
		}
	case expectSpaceOrMapKey, expectMapValueSeparatorOrSpaceOrMapKey:
		switch {
		case unicode.IsSpace(char):
		case unicode.IsLetter(char):
			pfa.stack = append(pfa.stack, ParseStackItem{itemType: parseStackItemKey, pos: pos, itemString: string(char)})
			pfa.state = expectKey
		case char == '"':
			pfa.stack = append(pfa.stack, ParseStackItem{itemType: parseStackItemKey, pos: pos, itemString: ``})
			pfa.state = expectDoubleQuotedKeyContent
		case char == '\'':
			pfa.stack = append(pfa.stack, ParseStackItem{itemType: parseStackItemKey, pos: pos, itemString: ``})
			pfa.state = expectSingleQuotedKeyContent
		case char == ',' && pfa.state == expectMapValueSeparatorOrSpaceOrMapKey:
			pfa.state = expectSpaceOrMapKey
		case char == '}':
			stackItem = getTopStackItem(pfa.stack, parseStackItemMap, pos, pfa.state)
			needFinishTopStackItem = true
		default:
			err = unexpectedCharError{}
			return
		}
	case
		expectMapKeySeparatorOrSpaceOrMapValue,
		expectSpaceOrMapValue,
		expectSpaceOrValue,
		expectSpaceOrArrayItem,
		expectArrayItemSeparatorOrSpaceOrArrayValue:
		switch {
		case pfa.state == expectMapKeySeparatorOrSpaceOrMapValue && char == '=':
			pfa.state = expectRocket
		case pfa.state == expectMapKeySeparatorOrSpaceOrMapValue && char == ':':
			pfa.state = expectSpaceOrMapValue
		case unicode.IsSpace(char):
		case char == '{':
			pfa.stack = append(pfa.stack, ParseStackItem{itemType: parseStackItemMap, pos: pos, itemMap: map[string]interface{}{}})
			pfa.state = expectSpaceOrMapKey
		case char == '[':
			pfa.stack = append(pfa.stack, ParseStackItem{itemType: parseStackItemArray, pos: pos, itemArray: []interface{}{}})
			pfa.state = expectSpaceOrArrayItem
		case char == ']' && (pfa.state == expectSpaceOrArrayItem || pfa.state == expectArrayItemSeparatorOrSpaceOrArrayValue):
			stackItem = getTopStackItem(pfa.stack, parseStackItemArray, pos, pfa.state)
			needFinishTopStackItem = true
		case char == '-' || char == '+' || unicode.IsNumber(char):
			pfa.stack = append(pfa.stack, ParseStackItem{itemType: parseStackItemNumber, pos: pos, itemString: string(char)})
			if unicode.IsNumber(char) {
				pfa.state = expectDigitOrUnderscoreOrDot
			} else {
				pfa.state = expectDigit
			}
		case char == '"' || char == '\'':
			pfa.stack = append(pfa.stack, ParseStackItem{itemType: parseStackItemString, pos: pos + 1, itemString: ``})
			if char == '"' {
				pfa.state = expectDoubleQuotedStringContent
			} else {
				pfa.state = expectSingleQuotedStringContent
			}
		case unicode.IsLetter(char):
			pfa.stack = append(pfa.stack, ParseStackItem{itemType: parseStackItemWord, pos: pos, itemString: string(char)})
			pfa.state = expectWord
		case char == ',' && pfa.state == expectArrayItemSeparatorOrSpaceOrArrayValue:
			pfa.state = expectSpaceOrArrayItem
		default:
			err = unexpectedCharError{}
			return
		}

	case expectWord:
		stackItem = getTopStackItem(pfa.stack, parseStackItemWord, pos, pfa.state)
		if unicode.IsLetter(char) {
			stackItem.itemString = stackItem.itemString + string(char)
		} else {
			needFinishTopStackItem = true
		}

	case expectKey:
		stackItem = getTopStackItem(pfa.stack, parseStackItemKey, pos, pfa.state)
		if unicode.IsLetter(char) {
			stackItem.itemString = stackItem.itemString + string(char)
		} else {
			needFinishTopStackItem = true
		}

	case expectDigit, expectDigitOrUnderscore, expectDigitOrUnderscoreOrDot:
		stackItem = getTopStackItem(pfa.stack, parseStackItemNumber, pos, pfa.state)
		if pfa.state == expectDigitOrUnderscoreOrDot {
			switch {
			case char == '.':
				pfa.state = expectDigitOrUnderscore
				fallthrough
			case unicode.IsNumber(char) || char == '_':
				stackItem.itemString = stackItem.itemString + string(char)
			default:
				needFinishTopStackItem = true
			}
		} else {
			if unicode.IsDigit(char) || (pfa.state == expectDigitOrUnderscore) && char == '_' {
				stackItem.itemString = stackItem.itemString + string(char)
				if pfa.state == expectDigit {
					pfa.state = expectDigitOrUnderscoreOrDot
				}
			} else {
				err = unexpectedCharError{}
				return
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
		switch pfa.state {
		case expectDoubleQuotedStringContent,
			expectDoubleQuotedStringEscapedContent,
			expectSingleQuotedStringContent,
			expectSingleQuotedStringEscapedContent:
			stackItem = getTopStackItem(pfa.stack, parseStackItemString, pos, pfa.state)
		default:
			stackItem = getTopStackItem(pfa.stack, parseStackItemKey, pos, pfa.state)
		}
		switch pfa.state {
		case
			expectDoubleQuotedKeyContent,
			expectSingleQuotedKeyContent,
			expectDoubleQuotedStringContent,
			expectSingleQuotedStringContent:
			if (pfa.state == expectDoubleQuotedStringContent || pfa.state == expectDoubleQuotedKeyContent) && char == '"' ||
				(pfa.state == expectSingleQuotedStringContent || pfa.state == expectSingleQuotedKeyContent) && char == '\'' {
				needFinishTopStackItem = true
			} else if char == '\\' {
				switch pfa.state {
				case expectDoubleQuotedStringContent:
					pfa.state = expectDoubleQuotedStringEscapedContent
				case expectSingleQuotedStringContent:
					pfa.state = expectSingleQuotedStringEscapedContent
				case expectDoubleQuotedKeyContent:
					pfa.state = expectDoubleQuotedKeyEscapedContent
				case expectSingleQuotedKeyContent:
					pfa.state = expectSingleQuotedKeyEscapedContent
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
				switch pfa.state {
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
				err = unexpectedCharError{}
				return
			}
			stackItem.itemString = stackItem.itemString + actualVal
			switch pfa.state {
			case expectDoubleQuotedStringEscapedContent:
				pfa.state = expectDoubleQuotedStringContent
			case expectSingleQuotedStringEscapedContent:
				pfa.state = expectSingleQuotedStringContent
			case expectDoubleQuotedKeyEscapedContent:
				pfa.state = expectDoubleQuotedKeyContent
			case expectSingleQuotedKeyEscapedContent:
				pfa.state = expectSingleQuotedKeyContent
			}
		}

	case expectMapKeySeparatorOrSpace:
		switch {
		case unicode.IsSpace(char):
			pfa.state = expectMapKeySeparatorOrSpaceOrMapValue
		case char == ':':
			pfa.state = expectSpaceOrMapValue
		case char == '=':
			pfa.state = expectRocket
		default:
			err = unexpectedCharError{}
			return
		}

	case expectRocket:
		if char == '>' {
			pfa.state = expectSpaceOrMapValue
		} else {
			err = unexpectedCharError{}
			return
		}

	default:
		log.Panicf(ansi.Ansi(`Err`, "unexpected <ansiOutline>pfa.state <ansiSecondaryLiteral>%s<ansi> <ansiOutline>pos <ansiSecondaryLiteral>%d <ansiOutline>pfa.stack <ansiSecondaryLiteral>%s"), pfa.state, pos, pfa.stack)
	}

	if needFinishTopStackItem {
		if err = pfa.finishTopStackItem(&char); err != nil {
			return
		}
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
		log.Panic(`pfa.stack should have at least one item`)
	}
	newStack := pfa.stack
	newState := tokenFinished
	stackItem := &pfa.stack[len(pfa.stack)-1]
	switch stackItem.itemType {

	case parseStackItemQwItem:
		if len(pfa.stack) < 2 {
			log.Panicf("len(pfa.stack) < 2")
		}
		stackSubItem := pfa.stack[len(pfa.stack)-1]
		newStack = pfa.stack[:len(pfa.stack)-1]
		stackItem = getTopStackItem(newStack, parseStackItemQw, stackSubItem.pos, pfa.state)
		stackItem.itemArray = append(stackItem.itemArray, stackSubItem.itemString)
		if charPtr == nil {
			log.Panicf("charPtr == nil")
		}
		if unicode.IsSpace(*charPtr) {
			newState = expectSpaceOrQwItemOrDelimiter
		} else {
			pfa.stack = newStack
			if len(pfa.stack) < 2 {
				log.Panicf("len(pfa.stack) < 2")
			}
			stackSubItem := pfa.stack[len(pfa.stack)-1]
			newStack = pfa.stack[:len(pfa.stack)-1]
			stackItem = getTopStackItem(newStack, parseStackItemArray, stackSubItem.pos, pfa.state)
			stackItem.itemArray = append(stackItem.itemArray, stackSubItem.itemArray...)
			newState = expectArrayItemSeparatorOrSpace
		}
		pfa.state = newState
		pfa.stack = newStack
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
						err = unexpectedCharError{}
						return
					}
				}
				if newState == expectSpaceOrQwItemOrDelimiter {
					stackItem.itemType = parseStackItemQw
					stackItem.itemArray = []interface{}{}
				}
			} else {
				err = core.Error("unexpected word <ansiPrimaryLiteral>%s<ansi> in non array context at "+getPosTitle(stackItem.pos), stackItem.itemString)
			}
			pfa.state = newState
			pfa.stack = newStack
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

	pfa.stack = newStack
	if len(pfa.stack) > 1 && charPtr != nil {
		var stackSubItem ParseStackItem
		stackSubItem, newStack = pfa.stack[len(pfa.stack)-1], pfa.stack[:len(pfa.stack)-1] // https://github.com/golang/go/wiki/SliceTricks
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
					err = unexpectedCharError{}
					return
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
						err = unexpectedCharError{}
						return
					}
				default:
					newState = expectMapValueSeparatorOrSpaceOrMapKey
				}
			}
		default:
			err = core.Error("<ansiOutline>stackItem <ansiSecondaryLiteral>%s<ansi> can not have subitem <ansiSecondaryLiteral>%s<ansiOutline>pfa.stack<ansiSecondaryLiteral>%s", stackItem, stackSubItem, pfa.stack)
		}
	}
	pfa.state = newState
	pfa.stack = newStack

	return
}

// func finishTopStackItem(stack ParseStack, state ParseState, charPtr *rune) (newStack ParseStack, newState ParseState, err error) {
// 	if len(stack) < 1 {
// 		log.Panic(`stack should have at least one item`)
// 	}
// 	newStack = stack
// 	newState = tokenFinished
// 	stackItem := &stack[len(stack)-1]
// 	switch stackItem.itemType {

// 	case parseStackItemQwItem:
// 		if len(stack) < 2 {
// 			log.Panicf("len(stack) < 2")
// 		}
// 		stackSubItem := stack[len(stack)-1]
// 		newStack = stack[:len(stack)-1]
// 		stackItem = getTopStackItem(newStack, parseStackItemQw, stackSubItem.pos, state)
// 		stackItem.itemArray = append(stackItem.itemArray, stackSubItem.itemString)
// 		if charPtr == nil {
// 			log.Panicf("charPtr == nil")
// 		}
// 		if unicode.IsSpace(*charPtr) {
// 			newState = expectSpaceOrQwItemOrDelimiter
// 		} else {
// 			stack = newStack
// 			if len(stack) < 2 {
// 				log.Panicf("len(stack) < 2")
// 			}
// 			stackSubItem := stack[len(stack)-1]
// 			newStack = stack[:len(stack)-1]
// 			stackItem = getTopStackItem(newStack, parseStackItemArray, stackSubItem.pos, state)
// 			stackItem.itemArray = append(stackItem.itemArray, stackSubItem.itemArray...)
// 			newState = expectArrayItemSeparatorOrSpace
// 		}
// 		return

// 	case parseStackItemNumber:
// 		source := underscoreRegexp.ReplaceAllLiteralString(stackItem.itemString, ``)
// 		if strings.Contains(source, `.`) {
// 			var float64Val float64
// 			if float64Val, err = strconv.ParseFloat(source, 64); err == nil {
// 				stackItem.value = float64Val
// 			}
// 		} else {
// 			var int64Val int64
// 			if int64Val, err = strconv.ParseInt(source, 10, 64); err == nil {
// 				if int64(MinInt) <= int64Val && int64Val <= int64(MaxInt) {
// 					stackItem.value = int(int64Val)
// 				} else {
// 					stackItem.value = int64Val
// 				}
// 			}
// 		}
// 		if err != nil {
// 			err = core.Error("failed to get number from string <ansiPrimaryLiteral>%s<ansi> at "+getPosTitle(stackItem.pos)+": %v", stackItem.itemString, err)
// 		}

// 	case parseStackItemString:
// 		stackItem.value = stackItem.itemString

// 	case parseStackItemWord:
// 		switch stackItem.itemString {
// 		case "true":
// 			stackItem.value = true
// 		case "false":
// 			stackItem.value = false
// 		case "nil":
// 			stackItem.value = nil
// 		case "qw":
// 			if len(stack) >= 2 && stack[len(stack)-2].itemType == parseStackItemArray && charPtr != nil {
// 				newState = expectSpaceOrQwItemOrDelimiter
// 				switch *charPtr {
// 				case '<':
// 					stackItem.delimiter = '>'
// 				case '[':
// 					stackItem.delimiter = ']'
// 				case '(':
// 					stackItem.delimiter = ')'
// 				case '{':
// 					stackItem.delimiter = '}'
// 				default:
// 					switch {
// 					case unicode.IsPunct(*charPtr) || unicode.IsSymbol(*charPtr):
// 						stackItem.delimiter = *charPtr
// 					default:
// 						err = unexpectedCharError{}
// 						return
// 					}
// 				}
// 				if newState == expectSpaceOrQwItemOrDelimiter {
// 					stackItem.itemType = parseStackItemQw
// 					stackItem.itemArray = []interface{}{}
// 				}
// 			} else {
// 				err = core.Error("unexpected word <ansiPrimaryLiteral>%s<ansi> in non array context at "+getPosTitle(stackItem.pos), stackItem.itemString)
// 			}
// 			return
// 		default:
// 			err = core.Error("unexpected word <ansiPrimaryLiteral>%s<ansi> at "+getPosTitle(stackItem.pos)+" <ansiOutline>stack <ansiSecondaryLiteral>%s", stackItem.itemString, stack)
// 		}

// 	case parseStackItemArray:
// 		stackItem.value = stackItem.itemArray

// 	case parseStackItemMap:
// 		stackItem.value = stackItem.itemMap

// 	case parseStackItemKey:

// 	default:
// 		err = core.Error("can not finish item of type <ansiPrimaryLiteral>%s<ansi> at "+getPosTitle(stackItem.pos), stackItem.itemType)
// 	}
// 	if err != nil {
// 		return
// 	}

// 	stack = newStack
// 	if len(stack) > 1 && charPtr != nil {
// 		var stackSubItem ParseStackItem
// 		stackSubItem, newStack = stack[len(stack)-1], stack[:len(stack)-1] // https://github.com/golang/go/wiki/SliceTricks
// 		stackItem = &newStack[len(newStack)-1]
// 		switch stackItem.itemType {
// 		case parseStackItemArray:
// 			stackItem.itemArray = append(stackItem.itemArray, stackSubItem.value)
// 			switch stackSubItem.itemType {
// 			case parseStackItemNumber, parseStackItemWord:
// 				switch {
// 				case unicode.IsSpace(*charPtr):
// 					newState = expectArrayItemSeparatorOrSpaceOrArrayValue
// 				case *charPtr == ',':
// 					newState = expectSpaceOrArrayItem
// 				default:
// 					err = unexpectedCharError{}
// 					return
// 				}
// 			default:
// 				newState = expectArrayItemSeparatorOrSpaceOrArrayValue
// 			}

// 		case parseStackItemMap:
// 			switch stackSubItem.itemType {
// 			case parseStackItemKey:
// 				stackItem.currentKey = stackSubItem.itemString
// 				newState = expectMapKeySeparatorOrSpace
// 			default:
// 				stackItem.itemMap[stackItem.currentKey] = stackSubItem.value
// 				switch stackSubItem.itemType {
// 				case parseStackItemNumber, parseStackItemWord:
// 					switch {
// 					case unicode.IsSpace(*charPtr):
// 						newState = expectMapValueSeparatorOrSpaceOrMapKey
// 					case *charPtr == ',':
// 						newState = expectSpaceOrMapKey
// 					default:
// 						err = unexpectedCharError{}
// 						return
// 					}
// 				default:
// 					newState = expectMapValueSeparatorOrSpaceOrMapKey
// 				}
// 			}
// 		default:
// 			err = core.Error("<ansiOutline>stackItem <ansiSecondaryLiteral>%s<ansi> can not have subitem <ansiSecondaryLiteral>%s<ansiOutline>stack<ansiSecondaryLiteral>%s", stackItem, stackSubItem, stack)
// 		}
// 	}

// 	return
// }
