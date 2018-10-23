package pathparse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
)

func init() {
	pfaPrimaryStateMethodsCheck()
	pfaItemFinishMethodsCheck()
	pfaErrorValidatorsCheck()
}

// ============================================================================

type parseStackItemType uint16

const (
	parseStackItemBelow parseStackItemType = iota
	// parseStackItemOpenIndex
	// parseStackItemArrayIndex
	parseStackItemNumber
	parseStackItemKey
	parseStackItemString
	// parseStackItemMap
	// parseStackItemArray
	// parseStackItemQw
	// parseStackItemQwItem
	parseStackItemAbove
)

//go:generate stringer -type=parseStackItemType

// ============================================================================

type parseStackItem struct {
	itemType parseStackItemType
	start    runePtrStruct
	// itemArray  []interface{}
	// itemMap    map[string]interface{}
	delimiter *rune
	// currentKey string
	itemString string
	value      interface{}
}

func (stackItem *parseStackItem) DataForJSON() interface{} {
	result := map[string]interface{}{}
	result["itemType"] = stackItem.itemType.String()
	result["start"] = stackItem.start.DataForJSON()
	if stackItem.itemType == parseStackItemNumber && stackItem.delimiter != nil {
		result["delimiter"] = stackItem.delimiter
	}
	result["itemString"] = stackItem.itemString
	return result
}

func (stackItem *parseStackItem) String() (result string) {
	return bwjson.PrettyJsonOf(stackItem)
}

// ============================================================================

type parseStack []parseStackItem

func (stack *parseStack) DataForJSON() interface{} {
	result := []interface{}{}
	for _, item := range *stack {
		result = append(result, item.DataForJSON())
	}
	return result
}

func (stack *parseStack) String() (result string) {
	return bwjson.PrettyJsonOf(stack)
}

// ============================================================================

type parsePrimaryState uint16

const (
	parsePrimaryStateBelow parsePrimaryState = iota
	expectEOF
	expectPathSegment
	// expectRocket
	expectMapKey
	// expectWord
	expectDigit
	expectContentOf
	expectEscapedContentOf
	// expectSpaceOrMapKey
	// expectSpaceOrQwItemOrDelimiter
	// expectEndOfQwItem
	parsePrimaryStateAbove
)

//go:generate stringer -type=parsePrimaryState

type parseSecondaryState uint16

const (
	noSecondaryState parseSecondaryState = iota
	orPathSegmentDelimiter

	// orMapKeySeparator
	// orArrayItemSeparator

	orUnderscore
	// orCloseBracket

	doubleQuoted
	singleQuoted

	// orMapValueSeparator
)

//go:generate stringer -type=parseSecondaryState

type parseTertiaryState uint16

const (
	noTertiaryState parseTertiaryState = iota
	// stringToken
	// keyToken
	// orUnderscore
)

//go:generate stringer -type=parseTertiaryState

type parseState struct {
	primary   parsePrimaryState
	secondary parseSecondaryState
	tertiary  parseTertiaryState
}

func (state *parseState) setPrimary(primary parsePrimaryState) {
	state.setSecondary(primary, noSecondaryState)
}

func (state *parseState) setSecondary(primary parsePrimaryState, secondary parseSecondaryState) {
	state.setTertiary(primary, secondary, noTertiaryState)
}

func (state *parseState) setTertiary(primary parsePrimaryState, secondary parseSecondaryState, tertiary parseTertiaryState) {
	state.primary = primary
	state.secondary = secondary
	state.tertiary = tertiary
}

func (state parseState) String() string {
	if state.tertiary != noTertiaryState {
		return fmt.Sprintf(`%s.%s.%s`, state.primary, state.secondary, state.tertiary)
	} else if state.secondary != noSecondaryState {
		return fmt.Sprintf(`%s.%s`, state.primary, state.secondary)
	} else {
		return state.primary.String()
	}
}

// ============================================================================

type runePtrStruct struct {
	runePtr     *rune
	pos         int
	line        uint
	col         uint
	prefix      string
	prefixStart int
}

func (v runePtrStruct) copyPtr() *runePtrStruct {
	return &runePtrStruct{v.runePtr, v.pos, v.line, v.col, v.prefix, v.prefixStart}
}

func (v runePtrStruct) DataForJSON() interface{} {
	result := map[string]interface{}{}
	if v.runePtr == nil {
		result["rune"] = "EOF"
	} else {
		result["rune"] = string(*(v.runePtr))
	}
	result["line"] = v.line
	result["col"] = v.col
	result["pos"] = v.pos
	return result
}

// ============================================================================

type pfaStruct struct {
	stack         parseStack
	state         parseState
	result        []interface{}
	prev          *runePtrStruct
	curr          runePtrStruct
	next          *runePtrStruct
	runeProvider  pfaRuneProvider
	preLineCount  int
	postLineCount int
}

func (pfa pfaStruct) DataForJSON() interface{} {
	result := map[string]interface{}{}
	result["stack"] = pfa.stack.DataForJSON()
	result["state"] = pfa.state.String()
	result["result"] = pfa.result
	result["pos"] = strconv.FormatInt(int64(pfa.curr.pos), 10)
	result["curr"] = pfa.curr.DataForJSON()
	if pfa.prev != nil {
		result["prev"] = pfa.prev.DataForJSON()
	}
	if pfa.next != nil {
		result["next"] = pfa.prev.DataForJSON()
	}
	return result
}

func (pfa pfaStruct) String() string {
	return bwjson.PrettyJsonOf(pfa)
}

type pfaRuneProvider interface {
	PullRune() *rune
}

func pfaParse(runeProvider pfaRuneProvider) ([]interface{}, error) {
	pfa := pfaStruct{
		stack:         parseStack{},
		state:         parseState{primary: expectPathSegment},
		result:        []interface{}{},
		runeProvider:  runeProvider,
		curr:          runePtrStruct{pos: -1, line: 1},
		preLineCount:  3,
		postLineCount: 3,
	}
	var err error
	for {
		pfa.pullRune()
		var needFinishTopStackItem bool
		if needFinishTopStackItem, err = pfaPrimaryStateMethods[pfa.state.primary](&pfa); err == nil && needFinishTopStackItem {
			err = pfa.finishTopStackItem()
		}
		if err != nil {
			break
		}
		if pfa.curr.runePtr == nil {
			if pfa.state.primary != expectEOF {
				pfa.panic("pfa.state.primary != expectEOF")
			}
			break
		}
	}
	if err != nil {
		return nil, err
	} else {
		return pfa.result, nil
	}
}

func (pfa *pfaStruct) pullRune() {
	if pfa.curr.pos < 0 || pfa.curr.runePtr != nil {
		pfa.prev = pfa.curr.copyPtr()
		if pfa.next != nil {
			pfa.curr = *(pfa.next)
			pfa.next = nil
		} else {
			runePtr := pfa.runeProvider.PullRune()
			pos := pfa.prev.pos + 1
			line := pfa.prev.line
			col := pfa.prev.col
			prefix := pfa.prev.prefix
			prefixStart := pfa.prev.prefixStart
			if runePtr != nil && pfa.prev.runePtr != nil {
				if *(pfa.prev.runePtr) != '\n' {
					col += 1
				} else {
					line += 1
					col = 1
					if int(line) > pfa.preLineCount {
						i := strings.Index(prefix, "\n")
						prefix = prefix[i+1:]
						prefixStart += i + 1
					}
				}
			}
			if runePtr != nil {
				prefix += string(*runePtr)
			}
			pfa.curr = runePtrStruct{runePtr, pos, line, col, prefix, prefixStart}
		}
	}
}

func (pfa *pfaStruct) pushRune() {
	if pfa.prev == nil {
		pfa.panic("pfa.prev == nil")
	} else {
		pfa.next = pfa.curr.copyPtr()
		pfa.curr = *(pfa.prev)
	}
}

func (pfa *pfaStruct) currRune() (result rune) {
	if pfa.curr.runePtr == nil {
		result = '\000'
	} else {
		result = *pfa.curr.runePtr
	}
	return
}

func (pfa *pfaStruct) panic(args ...interface{}) {
	fmtString := "<ansiOutline>pfa<ansi> <ansiSecondary>%s<ansi>"
	if args != nil {
		fmtString += " " + args[0].(string)
	}
	fmtArgs := []interface{}{pfa}
	if len(args) > 1 {
		fmtArgs = append(fmtArgs, args[1:])
	}
	bwerror.Panicd(1, fmtString, fmtArgs...)
}

func (pfa *pfaStruct) ifStackLen(minLen int) bool {
	return len(pfa.stack) >= minLen
}

func (pfa *pfaStruct) mustStackLen(minLen int) {
	if !pfa.ifStackLen(minLen) {
		pfa.panic("<ansiOutline>minLen <ansiSecondary>%d", minLen)
	}
}

// func (pfa *pfaStruct) isTopStackItemOfType(itemType parseStackItemType) bool {
// 	return pfa.ifStackLen(1) && pfa.getTopStackItem().itemType == itemType
// }

func (pfa *pfaStruct) getTopStackItemOfType(itemType parseStackItemType) (stackItem *parseStackItem) {
	stackItem = pfa.getTopStackItem()
	if stackItem.itemType != itemType {
		pfa.panic("<ansiOutline>itemType<ansiSecondary>%s", itemType)
	}
	return
}

func (pfa *pfaStruct) getTopStackItem() *parseStackItem {
	pfa.mustStackLen(1)
	return &pfa.stack[len(pfa.stack)-1]
}

func (pfa *pfaStruct) popStackItem() (stackItem parseStackItem) {
	pfa.mustStackLen(1)
	stackItem = pfa.stack[len(pfa.stack)-1]
	pfa.stack = pfa.stack[:len(pfa.stack)-1]
	return
}

func (pfa *pfaStruct) finishTopStackItem() (err error) {
	stackItem := pfa.getTopStackItem()
	var skipPostProcess bool
	if skipPostProcess, err = pfaItemFinishMethods[stackItem.itemType](pfa); err == nil && !skipPostProcess {
		// pfa.panic()

		// if len(pfa.stack) == 1 {
		// 	pfa.result = stackItem.value
		// 	pfa.state.setSecondary(expectEOF, orSpace)
		// } else if len(pfa.stack) > 1 {
		// 	stackSubItem := pfa.popStackItem()
		_ = pfa.popStackItem()
		pfa.result = append(pfa.result, stackItem.value)
		pfa.state.setSecondary(expectEOF, orPathSegmentDelimiter)
		// switch stackItem.itemType {
		// case parseStackItemNumber:

		// 	stackItem.itemArray = append(stackItem.itemArray, stackSubItem.value)
		// 	pfa.state.setPrimary(expectSpaceOrQwItemOrDelimiter)

		// case parseStackItemArray:
		// 	if stackSubItem.itemType == parseStackItemQw {
		// 		stackItem.itemArray = append(stackItem.itemArray, stackSubItem.itemArray...)
		// 	} else {
		// 		stackItem.itemArray = append(stackItem.itemArray, stackSubItem.value)
		// 	}
		// 	pfa.state.setSecondary(expectPathSegment, orArrayItemSeparator)

		// case parseStackItemMap:
		// 	switch stackSubItem.itemType {
		// 	case parseStackItemKey:
		// 		stackItem.currentKey = stackSubItem.itemString
		// 		pfa.state.setSecondary(expectPathSegment, orMapKeySeparator)
		// 	default:
		// 		stackItem.itemMap[stackItem.currentKey] = stackSubItem.value
		// 		pfa.state.setSecondary(expectSpaceOrMapKey, orMapValueSeparator)
		// 	}
		// default:
		// 	pfa.panic()
		// }
		// }
	}
	return
}
