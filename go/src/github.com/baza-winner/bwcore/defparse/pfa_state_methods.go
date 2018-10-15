package defparse

import (
	"unicode"

	"github.com/baza-winner/bwcore/bwerror"
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

type ruleCondition interface {
	ConformsTo(pfa *pfaStruct) bool
}

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

func (v parseStackItemTypeSet) ConformsTo(pfa *pfaStruct) bool {
	return len(pfa.stack) > 0 && v.Has(pfa.getTopStackItem().itemType)
}

type processorAction interface {
	execute(pfa *pfaStruct, currRune rune)
}

type unexpectedRune struct{}

func (v unexpectedRune) execute(pfa *pfaStruct, currRune rune) {
	pfa.isUnexpectedRune = true
}

type needFinish struct{}

func (v needFinish) execute(pfa *pfaStruct, currRune rune) {
	pfa.needFinishTopStackItem = true
}

type appendRuneToItemString struct{}

func (v appendRuneToItemString) execute(pfa *pfaStruct, currRune rune) {
	pfa.getTopStackItem().itemString += string(currRune)
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

type addRune struct {
	r rune
}

func (v addRune) execute(pfa *pfaStruct, currRune rune) {
	pfa.getTopStackItem().itemString += string(v.r)
}

type delimiterProvider interface {
	Delimiter(pfa *pfaStruct, currRune rune) *rune
}

type delim struct{ r rune }

func (v delim) Delimiter(pfa *pfaStruct, currRune rune) *rune {
	return &v.r
}

type fromParentItem struct{}

func (v fromParentItem) Delimiter(pfa *pfaStruct, currRune rune) *rune {
	return pfa.getTopStackItem().delimiter
}

type itemStringProvider interface {
	ItemString(pfa *pfaStruct, currRune rune) string
}

type fromCurrRune struct{}

func (v fromCurrRune) ItemString(pfa *pfaStruct, currRune rune) string {
	return string(currRune)
}

func (v fromCurrRune) Delimiter(pfa *pfaStruct, currRune rune) *rune {
	return runePtr(currRune)
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

func createRule(args []interface{}) rule {
	result := rule{
		ruleNormal,
		hasRuneSlice{},
		ruleConditions{},
		[]processorAction{},
	}
	parseSecondaryStates := parseSecondaryStateSet{}
	parseStackItemTypes := parseStackItemTypeSet{}
	delimiterChecker := (*delimiterIs)(nil)
	stackLen := (*stackLenIs)(nil)
	runes := runeSet{}
	for _, arg := range args {
		if typedArg, ok := arg.(rune); ok {
			runes.Add(typedArg)
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
			parseSecondaryStates.Add(typedArg)
		} else if typedArg, ok := arg.(parseStackItemType); ok {
			parseStackItemTypes.Add(typedArg)
		} else if typedArg, ok := arg.(processorAction); ok {
			result.processorActions = append(result.processorActions, typedArg)
		} else if typedArg, ok := arg.(delimiterIs); ok {
			if delimiterChecker == nil {
				delimiterChecker = &typedArg
			} else {
				bwerror.Panic("result.delimiterChecker already set")
			}
		} else if typedArg, ok := arg.(stackLenIs); ok {
			if stackLen == nil {
				stackLen = &typedArg
			} else {
				bwerror.Panic("result.stackLen already set")
			}
		} else {
			bwerror.Panic("unexpected %#v", arg)
		}
	}
	if len(parseSecondaryStates) > 0 {
		result.conditions = append(result.conditions, parseSecondaryStates)
	}
	if len(parseStackItemTypes) > 0 {
		result.conditions = append(result.conditions, parseStackItemTypes)
	}
	if delimiterChecker != nil {
		result.conditions = append(result.conditions, *delimiterChecker)
	}
	if stackLen != nil {
		result.conditions = append(result.conditions, *stackLen)
	}
	if runes.Len() > 0 {
		result.runeChecker = append(result.runeChecker, runes)
	}
	if len(result.runeChecker) == 0 && result.kind == ruleNormal {
		result.kind = ruleDefault
	}
	return result
}

func (pfa *pfaStruct) processStateDef(def *stateDef) (err error) {
	var needBreak bool
	pfa.isUnexpectedRune = false
	pfa.needFinishTopStackItem = false
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
	if pfa.isUnexpectedRune {
		err = pfaErrorMake(pfa, unexpectedRuneError)
	}
	return
}

func (pfa *pfaStruct) tryToProcessRule(r rule, currRune rune) (needBreak bool) {
	needBreak = false
	if r.conditions.conformsTo(pfa) {
		for _, pa := range r.processorActions {
			pa.execute(pfa, currRune)
		}
		needBreak = true
	}
	return
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
	default:
		bwerror.Panic("unicodeCategory: %s", v)
	}
	return
}

func (v unicodeCategory) Len() int {
	return 1
}

type delimiterRune struct{}

func (v delimiterRune) HasRune(r rune) (result bool) {
	bwerror.Panic("unreachable")
	return
}

func (v delimiterRune) Len() int {
	return 1
}

type hasRune interface {
	HasRune(rune) bool
	Len() int
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
