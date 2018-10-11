package defparse

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
	parseStackItem_below_ parseStackItemType = iota
	parseStackItemKey
	parseStackItemString
	parseStackItemMap
	parseStackItemArray
	parseStackItemQw
	parseStackItemQwItem
	parseStackItemNumber
	parseStackItemWord
	parseStackItem_above_
)

//go:generate stringer -type=parseStackItemType

// ============================================================================

type parseStackItem struct {
	itemType   parseStackItemType
	start      runePtrStruct
	itemArray  []interface{}
	itemMap    map[string]interface{}
	delimiter  rune
	currentKey string
	itemString string
	value      interface{}
}

func (stackItem *parseStackItem) GetDataForJson() interface{} {
	result := map[string]interface{}{}
	result["itemType"] = stackItem.itemType.String()
	result["start"] = stackItem.start.GetDataForJson()
	switch stackItem.itemType {
	case parseStackItemArray:
		result["itemArray"] = stackItem.itemArray
		result["value"] = stackItem.value
	case parseStackItemQw:
		result["delimiter"] = string(stackItem.delimiter)
		result["itemArray"] = stackItem.itemArray
		result["value"] = stackItem.value
	case parseStackItemQwItem:
		result["delimiter"] = string(stackItem.delimiter)
		result["itemString"] = stackItem.itemString
	case parseStackItemMap:
		result["itemMap"] = stackItem.itemMap
		result["value"] = stackItem.value
	case parseStackItemNumber, parseStackItemString, parseStackItemWord, parseStackItemKey:
		result["itemString"] = stackItem.itemString
		result["value"] = stackItem.value
	}
	return result
}

func (stackItem *parseStackItem) String() (result string) {
	return bwjson.PrettyJsonOf(stackItem)
}

// ============================================================================

type parseStack []parseStackItem

func (stack *parseStack) GetDataForJson() interface{} {
	result := []interface{}{}
	for _, item := range *stack {
		result = append(result, item.GetDataForJson())
	}
	return result
}

func (stack *parseStack) String() (result string) {
	return bwjson.PrettyJsonOf(stack)
}

// ============================================================================

type parsePrimaryState uint16

const (
	parsePrimaryState_below_ parsePrimaryState = iota
	expectEOF
	expectValueOrSpace
	expectRocket
	expectWord
	expectDigit
	expectContentOf
	expectEscapedContentOf
	expectSpaceOrMapKey
	expectSpaceOrQwItemOrDelimiter
	expectEndOfQwItem
	parsePrimaryState_above_
)

//go:generate stringer -type=parsePrimaryState

type parseSecondaryState uint16

const (
	anySecondaryState parseSecondaryState = iota
	noSecondaryState
	orSpace

	orMapKeySeparator
	orArrayItemSeparator

	orUnderscoreOrDot
	orUnderscore

	stringToken
	keyToken
	// doubleQuoted
	// singleQuoted

	orMapValueSeparator
)

//go:generate stringer -type=parseSecondaryState

// type parseTertiaryState uint16

// const (
// 	noTertiaryState parseTertiaryState = iota
// 	stringToken
// 	keyToken
// )

// go:generate stringer -type=parseTertiaryState

type parseState struct {
	primary   parsePrimaryState
	secondary parseSecondaryState
	// tertiary  parseTertiaryState
}

func (state *parseState) setPrimary(primary parsePrimaryState) {
	state.setSecondary(primary, noSecondaryState)
}

func (state *parseState) setSecondary(primary parsePrimaryState, secondary parseSecondaryState) {
	state.primary = primary
	state.secondary = secondary
	// state.setTertiary(primary, secondary, noTertiaryState)
}

// func (state *parseState) setTertiary(primary parsePrimaryState, secondary parseSecondaryState, tertiary parseTertiaryState) {
// 	state.primary = primary
// 	state.secondary = secondary
// 	state.tertiary = tertiary
// }

func (state parseState) String() string {
	// if state.tertiary != noTertiaryState {
	// 	return fmt.Sprintf(`%s.%s.%s`, state.primary, state.secondary, state.tertiary)
	// } else
	if state.secondary != noSecondaryState {
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

func (v runePtrStruct) GetDataForJson() interface{} {
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
	result        interface{}
	prev          *runePtrStruct
	curr          runePtrStruct
	next          *runePtrStruct
	runeProvider  PfaRuneProvider
	preLineCount  int
	postLineCount int
}

func (pfa pfaStruct) GetDataForJson() interface{} {
	result := map[string]interface{}{}
	result["stack"] = pfa.stack.GetDataForJson()
	result["state"] = pfa.state.String()
	result["result"] = pfa.result
	result["pos"] = strconv.FormatInt(int64(pfa.curr.pos), 10)
	result["curr"] = pfa.curr.GetDataForJson()
	if pfa.prev != nil {
		result["prev"] = pfa.prev.GetDataForJson()
	}
	if pfa.next != nil {
		result["next"] = pfa.prev.GetDataForJson()
	}
	return result
}

func (pfa pfaStruct) String() string {
	return bwjson.PrettyJsonOf(pfa)
}

type PfaRuneProvider interface {
	PullRune() *rune
}

func pfaParse(runeProvider PfaRuneProvider) (interface{}, error) {
	pfa := pfaStruct{
		stack:         parseStack{},
		state:         parseState{primary: expectValueOrSpace},
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

func (pfa *pfaStruct) currRune() (result rune, isEOF bool) {
	if pfa.curr.runePtr == nil {
		result = '\000'
		isEOF = true
	} else {
		result = *pfa.curr.runePtr
		isEOF = false
	}
	return
}

func (pfa *pfaStruct) panic(args ...interface{}) {
	fmtString := "<ansiOutline>pfa<ansi> <ansiSecondaryLiteral>%s<ansi>"
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
		pfa.panic("<ansiOutline>minLen <ansiSecondaryLiteral>%d", minLen)
	}
}

func (pfa *pfaStruct) isTopStackItemOfType(itemType parseStackItemType) bool {
	return pfa.ifStackLen(1) && pfa.getTopStackItem().itemType == itemType
}

func (pfa *pfaStruct) getTopStackItemOfType(itemType parseStackItemType) (stackItem *parseStackItem) {
	stackItem = pfa.getTopStackItem()
	if stackItem.itemType != itemType {
		pfa.panic("<ansiOutline>itemType<ansiSecondaryLiteral>%s", itemType)
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
		if len(pfa.stack) == 1 {
			pfa.result = stackItem.value
			pfa.state.setSecondary(expectEOF, orSpace)
		} else if len(pfa.stack) > 1 {
			stackSubItem := pfa.popStackItem()
			stackItem = pfa.getTopStackItem()
			switch stackItem.itemType {
			case parseStackItemQw:
				stackItem.itemArray = append(stackItem.itemArray, stackSubItem.value)
				pfa.state.setPrimary(expectSpaceOrQwItemOrDelimiter)

			case parseStackItemArray:
				if stackSubItem.itemType == parseStackItemQw {
					stackItem.itemArray = append(stackItem.itemArray, stackSubItem.itemArray...)
				} else {
					stackItem.itemArray = append(stackItem.itemArray, stackSubItem.value)
				}
				pfa.state.setSecondary(expectValueOrSpace, orArrayItemSeparator)

			case parseStackItemMap:
				switch stackSubItem.itemType {
				case parseStackItemKey:
					stackItem.currentKey = stackSubItem.itemString
					pfa.state.setSecondary(expectValueOrSpace, orMapKeySeparator)
				default:
					stackItem.itemMap[stackItem.currentKey] = stackSubItem.value
					pfa.state.setSecondary(expectSpaceOrMapKey, orMapValueSeparator)
				}
			default:
				pfa.panic()
			}
		}
	}
	return
}
