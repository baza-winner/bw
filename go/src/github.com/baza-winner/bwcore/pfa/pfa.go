package pfa

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/runeprovider"
)

//go:generate stringer -type=UnicodeCategory,ErrorType,ruleKind

func init() {
	pfaErrorValidatorsCheck()
}

// ============================================================================

type parseStackItem struct {
	itemType   string
	start      runePtrStruct
	itemArray  []interface{}
	itemMap    map[string]interface{}
	delimiter  *rune
	itemString string
	value      interface{}
}

func (stackItem *parseStackItem) DataForJSON() interface{} {
	result := map[string]interface{}{}
	result["itemType"] = stackItem.itemType
	result["start"] = stackItem.start.DataForJSON()
	if stackItem.delimiter != nil {
		result["delimiter"] = string(*stackItem.delimiter)
	}
	switch stackItem.itemType {
	case "array", "qw":
		result["itemArray"] = stackItem.itemArray
		result["value"] = stackItem.value
	case "qwItem":
		result["itemString"] = stackItem.itemString
	case "map":
		result["itemMap"] = stackItem.itemMap
		result["value"] = stackItem.value
	case "number", "string", "word", "key":
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

// State - состояние pfa
type State struct {
	Primary   string
	Secondary string
}

// SetPrimary устанавливает Primary значение и сбрасывает Secondary
func (state *State) SetPrimary(Primary string) {
	state.SetSecondary(Primary, "")
}

// SetSecondary устанавливает как Primary так и Secondary значение
func (state *State) SetSecondary(Primary string, Secondary string) {
	state.Primary = Primary
	state.Secondary = Secondary
}

func (state State) String() (result string) {
	if state.Secondary != "" {
		result = fmt.Sprintf(`%s.%s`, state.Primary, state.Secondary)
	} else {
		result = state.Primary
	}
	return
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
	state         State
	stackSubItem  *parseStackItem
	prev          *runePtrStruct
	curr          runePtrStruct
	next          *runePtrStruct
	runeProvider  runeprovider.RuneProvider
	preLineCount  int
	postLineCount int
	err           error
	errorType     ErrorType
	vars          map[string]interface{}
}

func (pfa pfaStruct) DataForJSON() interface{} {
	result := map[string]interface{}{}
	result["stack"] = pfa.stack.DataForJSON()
	result["state"] = pfa.state.String()
	if pfa.stackSubItem != nil {
		result["stackSubItem"] = pfa.stackSubItem.DataForJSON()
	}
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

func Run(runeProvider runeprovider.RuneProvider, pfaStateDef *LogicDef, initialState State) (result interface{}, err error) {
	pfa := pfaStruct{
		stack:         parseStack{},
		state:         initialState,
		runeProvider:  runeProvider,
		curr:          runePtrStruct{pos: -1, line: 1},
		preLineCount:  3,
		postLineCount: 3,
		vars:          map[string]interface{}{},
	}
	for {
		pfa.processLogic(pfaStateDef)
		if pfa.err != nil || pfa.curr.runePtr == nil {
			break
		}
	}
	if pfa.err != nil {
		err = pfa.err
	} else {
		if pfa.state.Primary != "expectEOF" {
			pfa.panic("pfa.state.Primary != expectEOF")
		}
		if len(pfa.stack) > 1 {
			pfa.panic("len(pfa.stack) > 1")
		} else if len(pfa.stack) > 0 {
			result = pfa.getTopStackItem().value
		}
	}
	return
}

func (pfa *pfaStruct) PullRune() {
	if pfa.curr.pos < 0 || pfa.curr.runePtr != nil {
		pfa.prev = pfa.curr.copyPtr()
		if pfa.next != nil {
			pfa.curr = *(pfa.next)
			pfa.next = nil
		} else {
			runePtr, err := pfa.runeProvider.PullRune()
			if err != nil {
				bwerror.PanicErr(err)
			}
			// runePtr := pfa.runeProvider.PullRune()
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

func (pfa *pfaStruct) PushRune() {
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
	var IsEOF bool
	if result, IsEOF = pfa.currRune(); IsEOF {
		pfa.panic("mustCurrRune: IsEOF")
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

// func (pfa *pfaStruct) isTopStackItemOfType(itemType string) bool {
// 	return pfa.ifStackLen(1) && pfa.getTopStackItem().itemType == itemType
// }

func (pfa *pfaStruct) getTopStackItemOfType(itemType string) (stackItem *parseStackItem) {
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
	itemType string,
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

type LogicDef map[ruleKind][]rule

func CreateLogicDef(args ...[]interface{}) *LogicDef {
	result := LogicDef{}
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
	primaryStateValues := bwset.StringSet{}
	secondaryStateValues := bwset.StringSet{}
	itemTypeValues := bwset.StringSet{}
	subItemTypeValues := bwset.StringSet{}
	delimiterValues := bwset.RuneSet{}
	topItemStringValues := bwset.StringSet{}
	stackLenValues := bwset.IntSet{}
	currRuneValues := bwset.RuneSet{}
	vars := map[string]bwset.InterfaceSet{}
	for _, arg := range args {
		if typedArg, ok := arg.(rune); ok {
			currRuneValues.Add(typedArg)
		} else if _, ok := arg.(IsEOF); ok {
			if len(result.runeChecker) > 0 {
				bwerror.Panic("combined IsEOF and non eof hasRune")
			}
			result.kind = ruleEof
		} else if typedArg, ok := arg.(hasRune); ok {
			if result.kind == ruleEof {
				bwerror.Panic("combined IsEOF and non eof hasRune")
			}
			if typedArg.Len() > 0 {
				result.runeChecker = append(result.runeChecker, typedArg)
			}
		} else if typedArg, ok := arg.(PrimaryIs); ok {
			primaryStateValues.Add(typedArg.Value)
		} else if typedArg, ok := arg.(SecondaryIs); ok {
			secondaryStateValues.Add(typedArg.Value)
		} else if typedArg, ok := arg.(TopItemIs); ok {
			itemTypeValues.Add(typedArg.Value)
		} else if typedArg, ok := arg.(SubItemIs); ok {
			subItemTypeValues.Add(typedArg.Value)
		} else if typedArg, ok := arg.(processorAction); ok {
			result.processorActions = append(result.processorActions, typedArg)
		} else if typedArg, ok := arg.(DelimiterIs); ok {
			delimiterValues.Add(typedArg.Value)
		} else if typedArg, ok := arg.(VarIs); ok {
			varValues := vars[typedArg.VarName]
			if varValues == nil {
				varValues = bwset.InterfaceSet{}
				vars[typedArg.VarName] = varValues
			}
			varValues.Add(typedArg.VarValue)
		} else if typedArg, ok := arg.(StackLenIs); ok {
			stackLenValues.Add(typedArg.Value)
		} else if typedArg, ok := arg.(TopItemStringIs); ok {
			topItemStringValues.Add(typedArg.Value)
		} else if typedArg, ok := arg.(TopItemStringIsOneOf); ok {
			topItemStringValues.AddSet(typedArg.Value)
		} else {
			bwerror.Panic("unexpected %#v", arg)
		}
	}
	for varName, varValues := range vars {
		result.conditions = append(result.conditions, varChecker{varName, varValues})
	}
	if len(topItemStringValues) > 0 {
		result.conditions = append(result.conditions, itemStringChecker{topItemStringValues})
	}
	if len(primaryStateValues) > 0 {
		result.conditions = append(result.conditions, parsePrimaryStateChecker{primaryStateValues})
	}
	if len(secondaryStateValues) > 0 {
		result.conditions = append(result.conditions, parseSecondaryStateChecker{secondaryStateValues})
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

func (pfa *pfaStruct) processLogic(def *LogicDef) {
	var needBreak bool
	pfa.err = nil
	pfa.errorType = pfaErrorBelow
	currRune, IsEOF := pfa.currRune()
	if rules, ok := (*def)[ruleEof]; IsEOF && ok {
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
	case UnexpectedRune, FailedToGetNumber, UnknownWord:
		pfa.err = pfaErrorMake(pfa, pfa.errorType)
		pfa.errorType = pfaErrorBelow
	}
	return
}

func (pfa *pfaStruct) tryToProcessRule(r rule, currRune rune) (needBreak bool) {
	needBreak = false
	if r.conditions.conformsTo(pfa) {
		needBreak = true
		for _, pa := range r.processorActions {
			pa.execute(pfa)
			if pfa.err != nil || pfa.errorType > pfaErrorBelow {
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
	values bwset.StringSet
}

func (v parsePrimaryStateChecker) ConformsTo(pfa *pfaStruct) bool {
	return v.values.Has(pfa.state.Primary)
}

type parseSecondaryStateChecker struct {
	values bwset.StringSet
}

func (v parseSecondaryStateChecker) ConformsTo(pfa *pfaStruct) bool {
	return v.values.Has(pfa.state.Secondary)
}

type itemStringChecker struct {
	values bwset.StringSet
}

func (v itemStringChecker) ConformsTo(pfa *pfaStruct) bool {
	return len(pfa.stack) > 0 && v.values.Has(pfa.getTopStackItem().itemString)
}

type parseStackItemTypeChecker struct {
	itemProvider parseStackItemProvider
	itemTypes    bwset.StringSet
}

func (v parseStackItemTypeChecker) ConformsTo(pfa *pfaStruct) bool {
	stackItem := v.itemProvider.StackItem(pfa)
	return stackItem != nil && v.itemTypes.Has(stackItem.itemType)
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

// =============================================================================
