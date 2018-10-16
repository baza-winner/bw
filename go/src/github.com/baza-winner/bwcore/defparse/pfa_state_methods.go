package defparse

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwint"
	"github.com/baza-winner/bwcore/bwset"
)

type ruleKind uint8

const (
	ruleNormal ruleKind = iota
	ruleDefault
	ruleEof
)

//go:generate stringer -type=ruleKind

type rule struct {
	kind             ruleKind
	runeChecker      hasRuneSlice
	conditions       ruleConditions
	processorActions []processorAction
}

type stateDef map[ruleKind][]rule

func createStateDef(args ...[]interface{}) *stateDef {
	result := stateDef{} //[]rule{}, []rule{}, []rule{}}
	for _, arg := range args {
		r := createRule(arg)
		if _, ok := result[r.kind]; !ok {
			result[r.kind] = []rule{}
		}
		result[r.kind] = append(result[r.kind], r)
	}
	return &result
}

type itemStringSet struct{ bwset.StringSet }

type parseStackItemTypeChecker struct {
	itemProvider parseStackItemProvider
	itemTypes    parseStackItemTypeSet
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

type subItem struct{ itemType parseStackItemType }

func createRule(args []interface{}) rule {
	result := rule{
		ruleNormal,
		hasRuneSlice{},
		ruleConditions{},
		[]processorAction{},
	}
	parseSecondaryStateChecker := parseSecondaryStateSet{}
	itemChecker := parseStackItemTypeChecker{topStackItemProvider{}, parseStackItemTypeSet{}}
	subItemChecker := parseStackItemTypeChecker{stackSubItemProvider{}, parseStackItemTypeSet{}}
	delimiterChecker := (*delimiterIs)(nil)
	itemStringChecker := itemStringSet{bwset.StringSet{}}
	stackLenChecker := (*stackLenIs)(nil)
	currRuneChecker := runeSet{}
	for _, arg := range args {
		if typedArg, ok := arg.(rune); ok {
			currRuneChecker.Add(typedArg)
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
		} else if typedArg, ok := arg.(parseSecondaryState); ok {
			parseSecondaryStateChecker.Add(typedArg)
		} else if typedArg, ok := arg.(parseStackItemType); ok {
			itemChecker.itemTypes.Add(typedArg)
		} else if typedArg, ok := arg.(subItem); ok {
			subItemChecker.itemTypes.Add(typedArg.itemType)
		} else if typedArg, ok := arg.(processorAction); ok {
			result.processorActions = append(result.processorActions, typedArg)
		} else if typedArg, ok := arg.(delimiterIs); ok {
			if delimiterChecker == nil {
				delimiterChecker = &typedArg
			} else {
				bwerror.Panic("result.delimiterChecker already set")
			}
		} else if typedArg, ok := arg.(stackLenIs); ok {
			if stackLenChecker == nil {
				stackLenChecker = &typedArg
			} else {
				bwerror.Panic("result.stackLenChecker already set")
			}
		} else if typedArg, ok := arg.(string); ok {
			itemStringChecker.Add(typedArg)
		} else {
			bwerror.Panic("unexpected %#v", arg)
		}
	}
	if len(itemStringChecker.StringSet) > 0 {
		result.conditions = append(result.conditions, itemStringChecker)
	}
	if len(parseSecondaryStateChecker) > 0 {
		result.conditions = append(result.conditions, parseSecondaryStateChecker)
	}
	if len(itemChecker.itemTypes) > 0 {
		result.conditions = append(result.conditions, itemChecker)
	}
	if len(subItemChecker.itemTypes) > 0 {
		result.conditions = append(result.conditions, subItemChecker)
	}
	if delimiterChecker != nil {
		result.conditions = append(result.conditions, *delimiterChecker)
	}
	if stackLenChecker != nil {
		result.conditions = append(result.conditions, *stackLenChecker)
	}
	if currRuneChecker.Len() > 0 {
		result.runeChecker = append(result.runeChecker, currRuneChecker)
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
	pfa.needFinishTopStackItem = false
	pfa.skipPostProcess = false
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
	// if pfa.errorType == unexpectedRuneError {
	// fmt.Printf("pfa: %s\n", bwjson.PrettyJsonOf(pfa))
	// pfa.panic("THERE")
	// }
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
			pa.execute(pfa, currRune)
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

type delimiterIs struct {
	delimiter rune
}

func (v delimiterIs) ConformsTo(pfa *pfaStruct) bool {
	stackItem := pfa.getTopStackItem()
	return stackItem.delimiter != nil && v.delimiter == *stackItem.delimiter
}

type stackLenIs struct {
	i int
}

func (v stackLenIs) ConformsTo(pfa *pfaStruct) bool {
	return v.i == len(pfa.stack)
}

func (v parseSecondaryStateSet) ConformsTo(pfa *pfaStruct) bool {
	return v.Has(pfa.state.secondary)
}

func (v parseStackItemTypeChecker) ConformsTo(pfa *pfaStruct) bool {
	stackItem := v.itemProvider.StackItem(pfa)
	return stackItem != nil && v.itemTypes.Has(stackItem.itemType)
}

// func (v parseStackItemTypeSet) ConformsTo(pfa *pfaStruct) bool {
// 	return len(pfa.stack) > 0 && v.Has(pfa.stackSubItem.itemType)
// }

func (v itemStringSet) ConformsTo(pfa *pfaStruct) bool {
	return len(pfa.stack) > 0 && v.Has(pfa.getTopStackItem().itemString)
}

// ============== delimiterProvider, itemStringProvider ========================

type delimiterProvider interface {
	Delimiter(pfa *pfaStruct, currRune rune) *rune
}

type itemStringProvider interface {
	ItemString(pfa *pfaStruct, currRune rune) string
}

type delim struct{ r rune }

func (v delim) Delimiter(pfa *pfaStruct, currRune rune) *rune {
	return &v.r
}

type fromParentItem struct{}

func (v fromParentItem) Delimiter(pfa *pfaStruct, currRune rune) *rune {
	return pfa.getTopStackItem().delimiter
}

type fromCurrRune struct{}

func (v fromCurrRune) ItemString(pfa *pfaStruct, currRune rune) string {
	return string(currRune)
}

func (v fromCurrRune) Delimiter(pfa *pfaStruct, currRune rune) *rune {
	return &currRune
}

type pairForCurrRune struct{}

func (v pairForCurrRune) Delimiter(pfa *pfaStruct, currRune rune) *rune {
	var r rune
	switch currRune {
	case '<':
		r = '>'
	case '[':
		r = ']'
	case '(':
		r = ')'
	case '{':
		r = '}'
	default:
		r = currRune
	}
	return &r
}

// ============================== hasRune =====================================

type hasRune interface {
	HasRune(rune) bool
	Len() int
}

func (v runeSet) HasRune(r rune) (result bool) {
	result = v.Has(r)
	return result
}

func (v runeSet) Len() int {
	return len(v)
}

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

// ======================== processorAction ====================================

type processorAction interface {
	execute(pfa *pfaStruct, currRune rune)
}

type setError struct{ errorType pfaErrorType }

func (v setError) execute(pfa *pfaStruct, currRune rune) {
	pfa.errorType = v.errorType
}

type needFinish struct{}

func (v needFinish) execute(pfa *pfaStruct, currRune rune) {
	pfa.needFinishTopStackItem = true
}

type pushRune struct{}

func (v pushRune) execute(pfa *pfaStruct, currRune rune) {
	pfa.pushRune()
}

type pullRune struct{}

func (v pullRune) execute(pfa *pfaStruct, currRune rune) {
	pfa.pullRune()

}

type changePrimary struct {
	primary parsePrimaryState
}

func (v changePrimary) execute(pfa *pfaStruct, currRune rune) {
	pfa.state.primary = v.primary
}

type setPrimary struct {
	primary parsePrimaryState
}

func (v setPrimary) execute(pfa *pfaStruct, currRune rune) {
	pfa.state.setPrimary(v.primary)
}

type changeSecondary struct {
	secondary parseSecondaryState
}

func (v changeSecondary) execute(pfa *pfaStruct, currRune rune) {
	pfa.state.secondary = v.secondary
}

type setSecondary struct {
	primary   parsePrimaryState
	secondary parseSecondaryState
}

func (v setSecondary) execute(pfa *pfaStruct, currRune rune) {
	pfa.state.setSecondary(v.primary, v.secondary)
}

type appendCurrRune struct{}

func (v appendCurrRune) execute(pfa *pfaStruct, currRune rune) {
	pfa.getTopStackItem().itemString += string(currRune)
}

type appendRune struct {
	r rune
}

func (v appendRune) execute(pfa *pfaStruct, currRune rune) {
	pfa.getTopStackItem().itemString += string(v.r)
}

type setTopItemDelimiter struct{ d delimiterProvider }

func (v setTopItemDelimiter) execute(pfa *pfaStruct, currRune rune) {
	pfa.getTopStackItem().delimiter = v.d.Delimiter(pfa, currRune)
}

type setTopItemType struct{ itemType parseStackItemType }

func (v setTopItemType) execute(pfa *pfaStruct, currRune rune) {
	pfa.getTopStackItem().itemType = v.itemType
	// pfa.getTopStackItem().stackItem.delimiter = delimiterProvider.Delimiter(pfa, currRune)
}

type pushItem struct {
	itemType   parseStackItemType
	itemString itemStringProvider
	delimiter  delimiterProvider
}

func (v pushItem) execute(pfa *pfaStruct, currRune rune) {
	var itemString string
	if v.itemString != nil {
		itemString = v.itemString.ItemString(pfa, currRune)
	}
	var delimiter *rune
	if v.delimiter != nil {
		delimiter = v.delimiter.Delimiter(pfa, currRune)
	}
	pfa.pushStackItem(
		v.itemType,
		itemString,
		delimiter,
	)
}

type setTopItemValueAsString struct{}

func (v setTopItemValueAsString) execute(pfa *pfaStruct, currRune rune) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemString
}

type setTopItemValueAsMap struct{}

func (v setTopItemValueAsMap) execute(pfa *pfaStruct, currRune rune) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemMap
}

type setTopItemValueAsArray struct{}

func (v setTopItemValueAsArray) execute(pfa *pfaStruct, currRune rune) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemArray
}

type setTopItemValueAsNumber struct{}

func (v setTopItemValueAsNumber) execute(pfa *pfaStruct, currRune rune) {
	stackItem := pfa.getTopStackItem()
	var err error
	if stackItem.value, err = _parseNumber(stackItem.itemString); err != nil {
		pfa.errorType = failedToGetNumberError
	}
}

type setTopItemValueAsBool struct{ b bool }

func (v setTopItemValueAsBool) execute(pfa *pfaStruct, currRune rune) {
	pfa.getTopStackItem().value = v.b
}

type setSkipPostProcess struct{ b bool }

func (v setSkipPostProcess) execute(pfa *pfaStruct, currRune rune) {
	pfa.skipPostProcess = true
}

type processStateDef struct{ def *stateDef }

func (v processStateDef) execute(pfa *pfaStruct, currRune rune) {
	pfa.processStateDef(v.def)
	if pfa.errorType == unexpectedRuneError {
		pfa.panic("SOME")
	}
}

type popSubItem struct{}

func (v popSubItem) execute(pfa *pfaStruct, currRune rune) {
	pfa.stackSubItem = pfa.popStackItem()
}

type appendItemArray struct{ p arrayProvider }

func (v appendItemArray) execute(pfa *pfaStruct, currRune rune) {
	stackItem := pfa.getTopStackItem()
	stackItem.itemArray = append(stackItem.itemArray, v.p.GetArray(pfa)...)
}

type setTopItemStringFromSubItem struct{}

func (v setTopItemStringFromSubItem) execute(pfa *pfaStruct, currRune rune) {
	pfa.getTopStackItem().itemString = pfa.stackSubItem.itemString
}

type setTopItemMapKeyValueFromSubItem struct{}

func (v setTopItemMapKeyValueFromSubItem) execute(pfa *pfaStruct, currRune rune) {
	stackItem := pfa.getTopStackItem()
	stackItem.itemMap[stackItem.itemString] = pfa.stackSubItem.value
}

type unreachable struct{}

func (v unreachable) execute(pfa *pfaStruct, currRune rune) {
	pfa.panic("unreachabe")
}

// ============================================================================

type arrayProvider interface {
	GetArray(pfa *pfaStruct) []interface{}
}

type fromSubItemValue struct{}

func (v fromSubItemValue) GetArray(pfa *pfaStruct) []interface{} {
	return []interface{}{pfa.stackSubItem.value}
}

type fromSubItemArray struct{}

func (v fromSubItemArray) GetArray(pfa *pfaStruct) []interface{} {
	return pfa.stackSubItem.itemArray
}

// ============================================================================

var underscoreRegexp = regexp.MustCompile("[_]+")

func _parseNumber(source string) (value interface{}, err error) {
	source = underscoreRegexp.ReplaceAllLiteralString(source, ``)
	if strings.Contains(source, `.`) {
		var float64Val float64
		if float64Val, err = strconv.ParseFloat(source, 64); err == nil {
			value = float64Val
		}
	} else {
		var int64Val int64
		if int64Val, err = strconv.ParseInt(source, 10, 64); err == nil {
			if int64(bwint.MinInt8) <= int64Val && int64Val <= int64(bwint.MaxInt8) {
				value = int8(int64Val)
			} else if int64(bwint.MinInt16) <= int64Val && int64Val <= int64(bwint.MaxInt16) {
				value = int16(int64Val)
			} else if int64(bwint.MinInt32) <= int64Val && int64Val <= int64(bwint.MaxInt32) {
				value = int32(int64Val)
			} else {
				value = int64Val
			}
		}
	}
	return
}

// ============================================================================
