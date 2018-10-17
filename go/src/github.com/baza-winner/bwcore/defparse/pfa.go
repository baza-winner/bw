package defparse

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwset"
)

var pfaStateDef *stateDef

func init() {
	pfaErrorValidatorsCheck()
	pfaStateDef = prepareStateDef()
}

// ============================================================================

// ============================================================================

type parseStackItemType uint16

func (v parseStackItemType) DataForJson() interface{} {
	return v.String()
}

// ============================================================================

type parseStackItem struct {
	itemType   parseStackItemType
	start      runePtrStruct
	itemArray  []interface{}
	itemMap    map[string]interface{}
	delimiter  *rune
	itemString string
	value      interface{}
}

func (stackItem *parseStackItem) DataForJson() interface{} {
	result := map[string]interface{}{}
	result["itemType"] = stackItem.itemType.String()
	result["start"] = stackItem.start.DataForJson()
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

func (stack *parseStack) DataForJson() interface{} {
	result := []interface{}{}
	for _, item := range *stack {
		result = append(result, item.DataForJson())
	}
	return result
}

func (stack *parseStack) String() (result string) {
	return bwjson.PrettyJsonOf(stack)
}

// ============================================================================

type parsePrimaryState uint16

func (v parsePrimaryState) DataForJson() interface{} {
	return v.String()
}

type parseSecondaryState uint16

func (v parseSecondaryState) DataForJson() interface{} {
	return v.String()
}

type parseState struct {
	primary   parsePrimaryState
	secondary parseSecondaryState
}

func (state *parseState) setPrimary(primary parsePrimaryState) {
	state.setSecondary(primary, noSecondaryState)
}

func (state *parseState) setSecondary(primary parsePrimaryState, secondary parseSecondaryState) {
	state.primary = primary
	state.secondary = secondary
}

func (state parseState) String() string {
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

func (v runePtrStruct) DataForJson() interface{} {
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
	stackSubItem  *parseStackItem
	prev          *runePtrStruct
	curr          runePtrStruct
	next          *runePtrStruct
	runeProvider  PfaRuneProvider
	preLineCount  int
	postLineCount int
	err           error
	errorType     pfaErrorType
	vars          map[string]interface{}
}

func (pfa pfaStruct) DataForJson() interface{} {
	result := map[string]interface{}{}
	result["stack"] = pfa.stack.DataForJson()
	result["state"] = pfa.state.String()
	if pfa.stackSubItem != nil {
		result["stackSubItem"] = pfa.stackSubItem.DataForJson()
	}
	result["pos"] = strconv.FormatInt(int64(pfa.curr.pos), 10)
	result["curr"] = pfa.curr.DataForJson()
	if pfa.prev != nil {
		result["prev"] = pfa.prev.DataForJson()
	}
	if pfa.next != nil {
		result["next"] = pfa.prev.DataForJson()
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
		vars:          map[string]interface{}{},
	}
	for {
		pfa.processStateDef(pfaStateDef)
		if pfa.err != nil || pfa.curr.runePtr == nil {
			break
		}
	}
	if pfa.err != nil {
		err = pfa.err
	} else {
		if pfa.state.primary != expectEOF {
			pfa.panic("pfa.state.primary != expectEOF")
		}
		if len(pfa.stack) > 1 {
			pfa.panic("len(pfa.stack) > 1")
		} else if len(pfa.stack) > 0 {
			result = pfa.getTopStackItem().value
		}
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

func (pfa *pfaStruct) mustCurrRune() (result rune) {
	var isEof bool
	if result, isEof = pfa.currRune(); isEof {
		pfa.panic("mustCurrRune: isEof")
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
		itemType:   itemType,
		start:      pfa.curr,
		itemArray:  []interface{}{},
		itemMap:    map[string]interface{}{},
		delimiter:  delimiter,
		itemString: itemString,
	})
}

// ============================================================================

type ruleKind uint8

const (
	ruleNormal ruleKind = iota
	ruleDefault
	ruleEof
)

type rule struct {
	kind             ruleKind
	runeChecker      hasRuneSlice
	conditions       ruleConditions
	processorActions []processorAction
}

type stateDef map[ruleKind][]rule

func createStateDef(args ...[]interface{}) *stateDef {
	result := stateDef{}
	for _, arg := range args {
		r := createRule(arg)
		if _, ok := result[r.kind]; !ok {
			result[r.kind] = []rule{}
		}
		result[r.kind] = append(result[r.kind], r)
	}
	return &result
}

func createRule(args []interface{}) rule {
	result := rule{
		ruleNormal,
		hasRuneSlice{},
		ruleConditions{},
		[]processorAction{},
	}
	parsePrimaryStateValues := parsePrimaryStateSet{}
	parseSecondaryStateValues := parseSecondaryStateSet{}
	itemTypeValues := parseStackItemTypeSet{}
	subItemTypeValues := parseStackItemTypeSet{}
	delimiterValues := bwset.RuneSet{}
	itemStringValues := bwset.StringSet{}
	stackLenValues := bwset.IntSet{}
	currRuneValues := bwset.RuneSet{}
	vars := map[string]bwset.InterfaceSet{}
	for _, arg := range args {
		if typedArg, ok := arg.(rune); ok {
			currRuneValues.Add(typedArg)
		} else if _, ok := arg.(eofRune); ok {
			if len(result.runeChecker) > 0 {
				bwerror.Panic("combined eofRune and non eof hasRune")
			}
			result.kind = ruleEof
		} else if typedArg, ok := arg.(hasRune); ok {
			if result.kind == ruleEof {
				bwerror.Panic("combined eofRune and non eof hasRune")
			}
			if typedArg.Len() > 0 {
				result.runeChecker = append(result.runeChecker, typedArg)
			}
		} else if typedArg, ok := arg.(parsePrimaryState); ok {
			parsePrimaryStateValues.Add(typedArg)
		} else if typedArg, ok := arg.(parseSecondaryState); ok {
			parseSecondaryStateValues.Add(typedArg)
			// } else if typedArg, ok := arg.(parseStackItemType); ok {
			// 	itemTypeValues.Add(typedArg)
		} else if typedArg, ok := arg.(topItem); ok {
			itemTypeValues.Add(typedArg.itemType)
		} else if typedArg, ok := arg.(subItem); ok {
			subItemTypeValues.Add(typedArg.itemType)
		} else if typedArg, ok := arg.(processorAction); ok {
			result.processorActions = append(result.processorActions, typedArg)
		} else if typedArg, ok := arg.(delimiterIs); ok {
			delimiterValues.Add(typedArg.delimiter)
		} else if typedArg, ok := arg.(varIs); ok {
			varValues := vars[typedArg.varName]
			if varValues == nil {
				varValues = bwset.InterfaceSet{}
				vars[typedArg.varName] = varValues
			}
			varValues.Add(typedArg.varValue)
		} else if typedArg, ok := arg.(stackLenIs); ok {
			stackLenValues.Add(typedArg.i)
		} else if typedArg, ok := arg.(string); ok {
			itemStringValues.Add(typedArg)
		} else {
			bwerror.Panic("unexpected %#v", arg)
		}
	}
	for varName, varValues := range vars {
		result.conditions = append(result.conditions, varChecker{varName, varValues})
	}
	if len(itemStringValues) > 0 {
		result.conditions = append(result.conditions, itemStringChecker{itemStringValues})
	}
	if len(parsePrimaryStateValues) > 0 {
		result.conditions = append(result.conditions, parsePrimaryStateChecker{parsePrimaryStateValues})
	}
	if len(parseSecondaryStateValues) > 0 {
		result.conditions = append(result.conditions, parseSecondaryStateChecker{parseSecondaryStateValues})
	}
	if len(itemTypeValues) > 0 {
		result.conditions = append(result.conditions, parseStackItemTypeChecker{topStackItemProvider{}, itemTypeValues})
	}
	if len(subItemTypeValues) > 0 {
		result.conditions = append(result.conditions, parseStackItemTypeChecker{stackSubItemProvider{}, subItemTypeValues})
	}
	if len(delimiterValues) > 0 {
		result.conditions = append(result.conditions, delimiterChecker{delimiterValues})
	}
	if len(stackLenValues) > 0 {
		result.conditions = append(result.conditions, stackLenChecker{stackLenValues})
	}
	if len(currRuneValues) > 0 {
		result.runeChecker = append(result.runeChecker, runeSet{currRuneValues})
	}
	if len(result.runeChecker) == 0 && result.kind == ruleNormal {
		result.kind = ruleDefault
	}
	return result
}

func (pfa *pfaStruct) processStateDef(def *stateDef) {
	var needBreak bool
	pfa.err = nil
	pfa.errorType = pfaError_below_
	currRune, isEof := pfa.currRune()
	if rules, ok := (*def)[ruleEof]; isEof && ok {
		for _, rule := range rules {
			needBreak = pfa.tryToProcessRule(rule, currRune)
			if needBreak {
				break
			}
		}
	} else {
		if rules, ok := (*def)[ruleNormal]; ok {
			for _, rule := range rules {
				if rule.runeChecker.HasRune(pfa, currRune) {
					needBreak = pfa.tryToProcessRule(rule, currRune)
					if needBreak {
						goto End
					}
				}
			}
		}
		if rules, ok := (*def)[ruleDefault]; ok {
			for _, rule := range rules {
				needBreak = pfa.tryToProcessRule(rule, currRune)
				if needBreak {
					goto End
				}
			}
		}
	End:
	}
	switch pfa.errorType {
	case unexpectedRuneError, failedToGetNumberError, unknownWordError:
		pfa.err = pfaErrorMake(pfa, pfa.errorType)
		pfa.errorType = pfaError_below_
	}
	return
}

func (pfa *pfaStruct) tryToProcessRule(r rule, currRune rune) (needBreak bool) {
	needBreak = false
	if r.conditions.conformsTo(pfa) {
		needBreak = true
		for _, pa := range r.processorActions {
			pa.execute(pfa)
			if pfa.err != nil || pfa.errorType > pfaError_below_ {
				break
			}
		}
	}
	return
}

// ========================= ruleCondition =====================================

type ruleConditions []ruleCondition

func (v ruleConditions) conformsTo(pfa *pfaStruct) (result bool) {
	result = true
	for _, i := range v {
		if !i.ConformsTo(pfa) {
			result = false
			break
		}
	}
	return
}

type ruleCondition interface {
	ConformsTo(pfa *pfaStruct) bool
}

type delimiterChecker struct {
	values bwset.RuneSet
}

func (v delimiterChecker) ConformsTo(pfa *pfaStruct) bool {
	stackItem := pfa.getTopStackItem()
	return stackItem.delimiter != nil && v.values.Has(*stackItem.delimiter)
}

type stackLenChecker struct {
	values bwset.IntSet
}

func (v stackLenChecker) ConformsTo(pfa *pfaStruct) bool {
	return v.values.Has(len(pfa.stack))
}

type varChecker struct {
	varName   string
	varValues bwset.InterfaceSet
}

func (v varChecker) ConformsTo(pfa *pfaStruct) bool {
	return v.varValues.Has(pfa.vars[v.varName])
}

type parsePrimaryStateChecker struct {
	values parsePrimaryStateSet
}

func (v parsePrimaryStateChecker) ConformsTo(pfa *pfaStruct) bool {
	return v.values.Has(pfa.state.primary)
}

type parseSecondaryStateChecker struct {
	values parseSecondaryStateSet
}

func (v parseStackItemTypeChecker) ConformsTo(pfa *pfaStruct) bool {
	stackItem := v.itemProvider.StackItem(pfa)
	return stackItem != nil && v.itemTypes.Has(stackItem.itemType)
}

type itemStringChecker struct {
	values bwset.StringSet
}

func (v itemStringChecker) ConformsTo(pfa *pfaStruct) bool {
	return len(pfa.stack) > 0 && v.values.Has(pfa.getTopStackItem().itemString)
}

type parseStackItemTypeChecker struct {
	itemProvider parseStackItemProvider
	itemTypes    parseStackItemTypeSet
}

func (v parseSecondaryStateChecker) ConformsTo(pfa *pfaStruct) bool {
	return v.values.Has(pfa.state.secondary)
}

type parseStackItemProvider interface {
	StackItem(pfa *pfaStruct) *parseStackItem
}

type topStackItemProvider struct{}

func (v topStackItemProvider) StackItem(pfa *pfaStruct) (result *parseStackItem) {
	if len(pfa.stack) > 0 {
		result = pfa.getTopStackItem()
	}
	return
}

type stackSubItemProvider struct{}

func (v stackSubItemProvider) StackItem(pfa *pfaStruct) *parseStackItem {
	return pfa.stackSubItem
}

// ============================== hasRune =====================================

type hasRune interface {
	HasRune(rune) bool
	Len() int
}

type runeSet struct {
	values bwset.RuneSet
}

func (v runeSet) HasRune(r rune) (result bool) {
	result = v.values.Has(r)
	return result
}

func (v runeSet) Len() int {
	return len(v.values)
}

type unicodeCategory uint8

const (
	unicodeSpace unicodeCategory = iota
	unicodeLetter
	unicodeDigit
	unicodeOpenBraces
	unicodePunct
	unicodeSymbol
)

func (v unicodeCategory) HasRune(r rune) (result bool) {
	switch v {
	case unicodeSpace:
		result = unicode.IsSpace(r)
	case unicodeLetter:
		result = unicode.IsLetter(r) || r == '_'
	case unicodeDigit:
		result = unicode.IsDigit(r)
	case unicodeOpenBraces:
		result = r == '(' || r == '{' || r == '[' || r == '<'
	case unicodePunct:
		result = unicode.IsPunct(r)
	case unicodeSymbol:
		result = unicode.IsSymbol(r)
	default:
		bwerror.Panic("unicodeCategory: %s", v)
	}
	return
}

func (v unicodeCategory) Len() int {
	return 1
}

type hasRuneSlice []hasRune

type eofRune struct{}

func (v eofRune) HasRune(r rune) (result bool) {
	bwerror.Panic("unreachable")
	return
}

func (v eofRune) Len() int {
	return 1
}

func (v hasRuneSlice) HasRune(pfa *pfaStruct, r rune) (result bool) {
	result = false
	for _, i := range v {
		if _, ok := i.(delimiterRune); ok {
			if len(pfa.stack) > 0 {
				stackItem := pfa.getTopStackItem()
				result = stackItem.delimiter != nil && *stackItem.delimiter == r
			}
		} else {
			result = i.HasRune(r)
		}
		if result {
			break
		}
	}
	return result
}

type delimiterRune struct{}

func (v delimiterRune) HasRune(r rune) (result bool) {
	bwerror.Panic("unreachable")
	return
}

func (v delimiterRune) Len() int {
	return 1
}

// ============================================================================

type delimiterIs struct {
	delimiter rune
}

type stackLenIs struct {
	i int
}

type varIs struct {
	varName  string
	varValue interface{}
}

type topItem struct {
	itemType parseStackItemType
}

type subItem struct {
	itemType parseStackItemType
}

// ============== delimiterProvider, itemStringProvider ========================

type delimiterProvider interface {
	Delimiter(pfa *pfaStruct) *rune
}

type itemStringProvider interface {
	ItemString(pfa *pfaStruct) string
}

type delim struct{ r rune }

func (v delim) Delimiter(pfa *pfaStruct, currRune rune) *rune {
	return &v.r
}

type fromParentItem struct{}

func (v fromParentItem) Delimiter(pfa *pfaStruct) *rune {
	return pfa.getTopStackItem().delimiter
}

type fromCurrRune struct{}

func (v fromCurrRune) ItemString(pfa *pfaStruct) string {
	return string(pfa.mustCurrRune())
}

func (v fromCurrRune) Delimiter(pfa *pfaStruct) *rune {
	currRune := pfa.mustCurrRune()
	return &currRune
}

type pairForCurrRune struct{}

func (v pairForCurrRune) Delimiter(pfa *pfaStruct) *rune {
	r := pfa.mustCurrRune()
	switch r {
	case '<':
		r = '>'
	case '[':
		r = ']'
	case '(':
		r = ')'
	case '{':
		r = '}'
	}
	return &r
}

// =============================================================================
