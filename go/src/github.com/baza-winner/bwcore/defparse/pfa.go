package defparse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
)

var unexpectedEof, unexpectedRune []interface{}

var pfaPrimaryStateDef map[parsePrimaryState]*stateDef
var finishItemStateDef, postProcessStateDef *stateDef

func init() {
	// pfaPrimaryStateMethodsCheck()
	pfaItemFinishMethodsCheck()
	pfaErrorValidatorsCheck()

	unexpectedEof = []interface{}{eofRune{}, setError{unexpectedRuneError}}
	unexpectedRune = []interface{}{setError{unexpectedRuneError}}

	pfaPrimaryStateDef = map[parsePrimaryState]*stateDef{
		expectEOF: createStateDef(
			[]interface{}{eofRune{}, setPrimary{expectEOF}},
			[]interface{}{unicodeSpace},
			unexpectedRune,
		),
		expectRocket: createStateDef(
			[]interface{}{'>', setPrimary{expectValueOrSpace}},
			unexpectedRune,
		),
		expectWord: createStateDef(
			[]interface{}{unicodeLetter, unicodeDigit,
				appendCurrRune{},
			},
			[]interface{}{
				pushRune{},
				needFinish{},
			},
		),
		expectSpaceOrQwItemOrDelimiter: createStateDef(
			unexpectedEof,
			[]interface{}{unicodeSpace},
			[]interface{}{delimiterRune{},
				needFinish{},
			},
			[]interface{}{
				pushItem{itemType: parseStackItemQwItem, itemString: fromCurrRune{}, delimiter: fromParentItem{}},
				setPrimary{expectEndOfQwItem},
			},
		),
		expectEndOfQwItem: createStateDef(
			unexpectedEof,
			[]interface{}{unicodeSpace, delimiterRune{},
				pushRune{},
				needFinish{},
			},
			[]interface{}{
				appendCurrRune{},
			},
		),
		expectContentOf: createStateDef(
			unexpectedEof,
			[]interface{}{delimiterRune{},
				needFinish{},
			},
			[]interface{}{'\\',
				changePrimary{expectEscapedContentOf},
			},
			[]interface{}{
				appendCurrRune{},
			},
		),
		expectDigit: createStateDef(
			[]interface{}{unicodeDigit, noSecondaryState,
				appendCurrRune{},
				changeSecondary{orUnderscoreOrDot},
			},
			[]interface{}{'.', orUnderscoreOrDot,
				appendCurrRune{},
				changeSecondary{orUnderscore},
			},
			[]interface{}{'_', unicodeDigit, orUnderscoreOrDot, orUnderscore,
				appendCurrRune{},
			},
			[]interface{}{noSecondaryState,
				setError{unexpectedRuneError},
			},
			[]interface{}{
				pushRune{},
				needFinish{},
			},
		),
		expectSpaceOrMapKey: createStateDef(
			[]interface{}{unicodeSpace},
			[]interface{}{unicodeLetter,
				pushItem{itemType: parseStackItemKey, itemString: fromCurrRune{}},
				setPrimary{expectWord},
			},
			[]interface{}{'"', '\'',
				pushItem{itemType: parseStackItemKey, delimiter: fromCurrRune{}},
				setSecondary{expectContentOf, keyToken},
			},
			[]interface{}{',', orMapValueSeparator,
				setPrimary{expectSpaceOrMapKey},
			},
			[]interface{}{delimiterRune{}, parseStackItemMap,
				needFinish{},
			},
			unexpectedRune,
		),
		expectEscapedContentOf: createStateDef(
			[]interface{}{'"', '\'', '\\',
				appendCurrRune{},
				changePrimary{expectContentOf},
			},
			[]interface{}{delimiterIs{'"'},
				processStateDef{createStateDef(
					[]interface{}{'a', appendRune{'\a'}},
					[]interface{}{'b', appendRune{'\b'}},
					[]interface{}{'f', appendRune{'\f'}},
					[]interface{}{'n', appendRune{'\n'}},
					[]interface{}{'r', appendRune{'\r'}},
					[]interface{}{'t', appendRune{'\t'}},
					[]interface{}{'v', appendRune{'\v'}},
					unexpectedRune,
				)},
				changePrimary{expectContentOf},
			},
			unexpectedRune,
		),
		expectValueOrSpace: createStateDef(
			[]interface{}{eofRune{}, stackLenIs{0},
				setPrimary{expectEOF},
			},
			unexpectedEof,
			[]interface{}{'=', orMapKeySeparator,
				setPrimary{expectRocket},
			},
			[]interface{}{':', orMapKeySeparator,
				setPrimary{expectValueOrSpace},
			},
			[]interface{}{',', orArrayItemSeparator,
				setPrimary{expectValueOrSpace},
			},
			[]interface{}{unicodeSpace},
			[]interface{}{'{',
				pushItem{itemType: parseStackItemMap, delimiter: pairForCurrRune{}},
				setPrimary{expectSpaceOrMapKey},
			},
			[]interface{}{'<',
				pushItem{itemType: parseStackItemQw, delimiter: pairForCurrRune{}},
				setPrimary{expectSpaceOrQwItemOrDelimiter},
			},
			[]interface{}{'[',
				pushItem{itemType: parseStackItemArray, delimiter: pairForCurrRune{}},
				setPrimary{expectValueOrSpace},
			},
			[]interface{}{parseStackItemArray, delimiterRune{},
				needFinish{},
			},
			[]interface{}{'-', '+',
				pushItem{itemType: parseStackItemNumber, itemString: fromCurrRune{}},
				setPrimary{expectDigit},
			},
			[]interface{}{unicodeDigit,
				pushItem{itemType: parseStackItemNumber, itemString: fromCurrRune{}},
				setSecondary{expectDigit, orUnderscoreOrDot},
			},
			[]interface{}{'"', '\'',
				pullRune{},
				pushItem{itemType: parseStackItemString, delimiter: fromCurrRune{}},
				pushRune{},
				setSecondary{expectContentOf, stringToken},
			},
			[]interface{}{unicodeLetter,
				pushItem{itemType: parseStackItemWord, itemString: fromCurrRune{}},
				setPrimary{expectWord},
			},
			unexpectedRune,
		),
	}

	expect := parsePrimaryState_below_ + 1
	for expect < parsePrimaryState_above_ {
		if _, ok := pfaPrimaryStateDef[expect]; !ok {
			bwerror.Panic("not defined <ansiOutline>pfaPrimaryStateDef<ansi>[<ansiPrimaryLiteral>%s<ansi>]", expect)
		}
		expect += 1
	}

	finishItemStateDef = createStateDef(
		[]interface{}{parseStackItemString, parseStackItemQwItem,
			setTopItemValueAsString{},
		},
		[]interface{}{parseStackItemMap,
			setTopItemValueAsMap{},
		},
		[]interface{}{parseStackItemArray, parseStackItemQw,
			setTopItemValueAsArray{},
		},
		[]interface{}{parseStackItemNumber,
			setTopItemValueAsNumber{},
		},
		[]interface{}{parseStackItemWord,
			processStateDef{createStateDef(
				[]interface{}{"true",
					setTopItemValueAsBool{true},
				},
				[]interface{}{"false",
					setTopItemValueAsBool{false},
				},
				[]interface{}{"nil", "null"},
				[]interface{}{"Bool", "String", "Int", "Number", "Map", "Array", "ArrayOf",
					setTopItemValueAsString{},
				},
				[]interface{}{"qw",
					pullRune{},
					processStateDef{createStateDef(
						[]interface{}{unicodeOpenBraces, unicodePunct, unicodeSymbol,
							setPrimary{expectSpaceOrQwItemOrDelimiter},
							setTopItemDelimiter{pairForCurrRune{}},
							setTopItemType{parseStackItemQw},
						},
						[]interface{}{
							setError{unexpectedRuneError},
						},
					)},
					setSkipPostProcess{},
				},
				[]interface{}{
					setError{unknownWordError},
				},
			)},
		},
	)

	postProcessStateDef = createStateDef(
		[]interface{}{stackLenIs{0}},
		[]interface{}{stackLenIs{1},
			setSecondary{expectEOF, orSpace},
		},
		[]interface{}{
			popSubItem{},
			processStateDef{createStateDef(
				[]interface{}{parseStackItemQw,
					appendItemArray{fromSubItemValue{}},
					setPrimary{expectSpaceOrQwItemOrDelimiter},
				},
				[]interface{}{parseStackItemArray,
					processStateDef{createStateDef(
						[]interface{}{subItem{parseStackItemQw},
							appendItemArray{fromSubItemArray{}},
						},
						[]interface{}{
							appendItemArray{fromSubItemValue{}},
						},
					)},
					setSecondary{expectValueOrSpace, orArrayItemSeparator},
				},
				[]interface{}{parseStackItemMap,
					processStateDef{createStateDef(
						[]interface{}{subItem{parseStackItemKey},
							setTopItemStringFromSubItem{},
							setSecondary{expectValueOrSpace, orMapKeySeparator},
						},
						[]interface{}{
							setTopItemMapKeyValueFromSubItem{},
							setSecondary{expectSpaceOrMapKey, orMapValueSeparator},
						},
					)},
				},
				[]interface{}{
					unreachable{},
				},
			)},
		},
	)
}

//go:generate setter -type=rune

// ============================================================================

type unicodeCategory uint8

const (
	unicodeSpace unicodeCategory = iota
	unicodeLetter
	unicodeDigit
	unicodeOpenBraces
	unicodePunct
	unicodeSymbol
)

//go:generate stringer -type=unicodeCategory

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
//go:generate setter -type=parseStackItemType

// ============================================================================

type parseStackItem struct {
	itemType  parseStackItemType
	start     runePtrStruct
	itemArray []interface{}
	itemMap   map[string]interface{}
	delimiter *rune
	// currentKey string
	itemString string
	value      interface{}
}

func (stackItem *parseStackItem) GetDataForJson() interface{} {
	result := map[string]interface{}{}
	result["itemType"] = stackItem.itemType.String()
	result["start"] = stackItem.start.GetDataForJson()
	if stackItem.delimiter != nil {
		result["delimiter"] = string(*stackItem.delimiter)
	}
	switch stackItem.itemType {
	case parseStackItemArray, parseStackItemQw:
		result["itemArray"] = stackItem.itemArray
		result["value"] = stackItem.value
	case parseStackItemQwItem:
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
//go:generate setter -type=parseSecondaryState

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
	stack                  parseStack
	state                  parseState
	stackSubItem           *parseStackItem
	prev                   *runePtrStruct
	curr                   runePtrStruct
	next                   *runePtrStruct
	runeProvider           PfaRuneProvider
	preLineCount           int
	postLineCount          int
	err                    error
	errorType              pfaErrorType
	needFinishTopStackItem bool
	skipPostProcess        bool
}

func (pfa pfaStruct) GetDataForJson() interface{} {
	result := map[string]interface{}{}
	result["stack"] = pfa.stack.GetDataForJson()
	result["state"] = pfa.state.String()
	if pfa.stackSubItem != nil {
		result["stackSubItem"] = pfa.stackSubItem.GetDataForJson()
	}
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

func pfaParse(runeProvider PfaRuneProvider) (result interface{}, err error) {
	pfa := pfaStruct{
		stack:         parseStack{},
		state:         parseState{primary: expectValueOrSpace},
		runeProvider:  runeProvider,
		curr:          runePtrStruct{pos: -1, line: 1},
		preLineCount:  3,
		postLineCount: 3,
	}
	// var err error
	for {
		pfa.pullRune()
		var def *stateDef
		var ok bool
		if def, ok = pfaPrimaryStateDef[pfa.state.primary]; !ok {
			bwerror.Panic("pfa.state.primary: %s", pfa.state.primary)
		}
		if pfa.processStateDef(def); pfa.err == nil && pfa.needFinishTopStackItem {
			pfa.finishTopStackItem()
		}
		if pfa.err != nil {
			break
		}
		if pfa.curr.runePtr == nil {
			if pfa.state.primary != expectEOF {
				pfa.panic("pfa.state.primary != expectEOF")
			}
			if len(pfa.stack) > 1 {
				pfa.panic("len(pfa.stack) > 1")
			}
			break
		}
	}
	if pfa.err != nil {
		err = pfa.err
	} else if len(pfa.stack) > 0 {
		result = pfa.getTopStackItem().value
	}
	return
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
	// _, isEof := pfa.currRune()
	// if isEof {
	// 	pfa.panic("HeRe")
	// }
	// fmt.Printf("pullRune: %s\n", bwjson.PrettyJsonOf(pfa))
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

func runePtr(r rune) *rune {
	return &r
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

func (pfa *pfaStruct) popStackItem() *parseStackItem {
	pfa.mustStackLen(1)
	stackItem := pfa.stack[len(pfa.stack)-1]
	pfa.stack = pfa.stack[:len(pfa.stack)-1]
	return &stackItem
}

func (pfa *pfaStruct) pushStackItem(
	itemType parseStackItemType,
	itemString string,
	delimiter *rune,
) {
	pfa.stack = append(pfa.stack, parseStackItem{
		itemType:  itemType,
		start:     pfa.curr,
		itemArray: []interface{}{},
		itemMap:   map[string]interface{}{},
		delimiter: delimiter,
		// currentKey: "",
		itemString: itemString,
		// value:      nil,
	})
}

func (pfa *pfaStruct) finishTopStackItem() {
	stackItem := pfa.getTopStackItem()
	if pfa.processStateDef(finishItemStateDef); pfa.err == nil && !pfa.skipPostProcess {

		if len(pfa.stack) == 1 {
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
					stackItem.itemString = stackSubItem.itemString
					pfa.state.setSecondary(expectValueOrSpace, orMapKeySeparator)
				default:
					stackItem.itemMap[stackItem.itemString] = stackSubItem.value
					pfa.state.setSecondary(expectSpaceOrMapKey, orMapValueSeparator)
				}

			default:
				pfa.panic("stackItem.itemType: %s", stackItem.itemType)
			}
		}
	}
}
