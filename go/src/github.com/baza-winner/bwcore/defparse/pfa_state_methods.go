package defparse

import (
	"unicode"

	"github.com/baza-winner/bwcore/bwerror"
)

type pfaPrimaryStateMethod func(*pfaStruct) (bool, error)

var pfaPrimaryStateMethods = map[parsePrimaryState]pfaPrimaryStateMethod{
	expectEOF:                      _expectEOF,
	expectValueOrSpace:             _expectValueOrSpace,
	expectRocket:                   _expectRocket,
	expectSpaceOrMapKey:            _expectSpaceOrMapKey,
	expectWord:                     _expectWord,
	expectDigit:                    _expectDigit,
	expectContentOf:                _expectContentOf,
	expectEscapedContentOf:         _expectEscapedContentOf,
	expectSpaceOrQwItemOrDelimiter: _expectSpaceOrQwItemOrDelimiter,
	expectEndOfQwItem:              _expectEndOfQwItem,
}

func pfaPrimaryStateMethodsCheck() {
	expect := parsePrimaryState_below_ + 1
	for expect < parsePrimaryState_above_ {
		if _, ok := pfaPrimaryStateMethods[expect]; !ok {
			bwerror.Panic("not defined <ansiOutline>pfaPrimaryStateMethods<ansi>[<ansiPrimaryLiteral>%s<ansi>]", expect)
		}
		expect += 1
	}
}

// type ruleProcessor func(*pfaStruct, rune) (bool, bool)

// func asEofProcessor(arg interface{}) (result eofProcessor, ok bool) {
// 	if result, ok = arg.(func(*pfaStruct) (bool, bool)); !ok {
// 		result, ok = arg.(eofProcessor)
// 	}
// 	return
// }
// func asRuleProcessor(arg interface{}) (result ruleProcessor, ok bool) {
// 	if result, ok = arg.(func(*pfaStruct, rune) (bool, bool)); !ok {
// 		result, ok = arg.(ruleProcessor)
// 	}
// 	return
// }

type ruleKind uint8

const (
	ruleNormal ruleKind = iota
	ruleDefault
	ruleEof
)

//go:generate stringer -type=ruleKind

type rule struct {
	kind        ruleKind
	runeChecker hasRuneSlice
	conditions  ruleConditions
	// parseSecondaryStates parseSecondaryStateSet
	// parseStackItemTypes  parseStackItemTypeSet
	// delimiterChecker     *delimiterIs
	// stackLen             *stackLenIs
	processorActions []processorAction
}

// type ruleSlice []rule
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
	return v.delimiter == pfa.getTopStackItem().delimiter
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
	execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune)
}

// func asProcessorAction(f interface{}) (result processorAction, ok bool) {
// 	if result, ok = f.(processorAction); !ok {
// 		result, ok = f.(processorAction)
// 	}
// 	return
// }

type unexpectedRune struct{}

func (v unexpectedRune) execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune) {
	*isUnexpectedRune = true
}

type needFinish struct{}

func (v needFinish) execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune) {
	*needFinishTopStackItem = true
}

type appendRuneToItemString struct{}

func (v appendRuneToItemString) execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune) {
	pfa.getTopStackItem().itemString += string(currRune)
}

type pushRune struct{}

func (v pushRune) execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune) {
	pfa.pushRune()
}

type pullRune struct{}

func (v pullRune) execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune) {
	pfa.pullRune()
}

type changePrimary struct {
	primary parsePrimaryState
}

func (v changePrimary) execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune) {
	pfa.state.primary = v.primary
}

type setPrimary struct {
	primary parsePrimaryState
}

func (v setPrimary) execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune) {
	pfa.state.setPrimary(v.primary)
}

type changeSecondary struct {
	secondary parseSecondaryState
}

func (v changeSecondary) execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune) {
	pfa.state.secondary = v.secondary
}

type setSecondary struct {
	primary   parsePrimaryState
	secondary parseSecondaryState
}

func (v setSecondary) execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune) {
	pfa.state.setSecondary(v.primary, v.secondary)
}

type addRune struct {
	r rune
}

func (v addRune) execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune) {
	pfa.getTopStackItem().itemString += string(v.r)
}

type pushStackItemQwItem struct{}

func (v pushStackItemQwItem) execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune) {
	stackItem := pfa.getTopStackItem()
	// pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemQwItem, start: pfa.curr, itemString: string(currRune), delimiter: stackItem.delimiter})
	pfa.pushStackItem(parseStackItemQwItem, string(currRune), stackItem.delimiter)
}

type pushStackItemQwWithDelimiter struct{ delimiter rune }

func (v pushStackItemQwWithDelimiter) execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune) {
	// pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemQw, start: pfa.curr, itemArray: []interface{}{}, delimiter: v.delimiter})
	pfa.pushStackItem(parseStackItemQw, "", v.delimiter)
}

type pushStackItemArray struct{}

func (v pushStackItemArray) execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune) {
	// pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemArray, start: pfa.curr, itemArray: []interface{}{}, delimiter: ']'})
	pfa.pushStackItem(parseStackItemArray, "", ']')
}

type pushStackItem struct{ itemType parseStackItemType }

func (v pushStackItem) execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune) {
	pfa.pushStackItem(v.itemType, string(currRune), '\000')
}

type pushStackItemWithCurrRuneAsDelimter struct{ itemType parseStackItemType }

func (v pushStackItemWithCurrRuneAsDelimter) execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune) {
	pfa.pushStackItem(v.itemType, "", currRune)
}

type pushStackItemMap struct{}

func (v pushStackItemMap) execute(pfa *pfaStruct, isUnexpectedRune *bool, needFinishTopStackItem *bool, currRune rune) {
	pfa.pushStackItem(parseStackItemMap, "", '}')
	// pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemMap, start: pfa.curr, itemMap: map[string]interface{}{}, delimiter: '}'})
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

func (pfa *pfaStruct) processStateDef(def *stateDef) (needFinishTopStackItem bool, err error) {
	var isUnexpectedRune, needBreak bool
	currRune, isEof := pfa.currRune()
	// bwerror.Spew.Printf("def: %#v\n", def)
	if rules, ok := (*def)[ruleEof]; isEof && ok {
		// bwerror.Spew.Printf("ruleEof: %#v\n", rules)
		for _, rule := range rules {
			needFinishTopStackItem, isUnexpectedRune, needBreak = pfa.tryToProcessRule(rule, currRune)
			if needBreak {
				break
			}
		}
	} else {
		if rules, ok := (*def)[ruleNormal]; ok {
			// bwerror.Spew.Printf("ruleNormal: %#v\n", rules)
			for _, rule := range rules {
				if rule.runeChecker.HasRune(pfa, currRune) {
					needFinishTopStackItem, isUnexpectedRune, needBreak = pfa.tryToProcessRule(rule, currRune)
					if needBreak {
						goto End
					}
				}
			}
		}
		if rules, ok := (*def)[ruleDefault]; ok {
			// bwerror.Spew.Printf("ruleDefault: %#v\n", rules)
			for _, rule := range rules {
				needFinishTopStackItem, isUnexpectedRune, needBreak = pfa.tryToProcessRule(rule, currRune)
				if needBreak {
					goto End
				}
			}
		}
	End:
	}
	if isUnexpectedRune {
		err = pfaErrorMake(pfa, unexpectedRuneError)
	}
	return
}

func (pfa *pfaStruct) tryToProcessRule(r rule, currRune rune) (needFinishTopStackItem bool, isUnexpectedRune bool, needBreak bool) {
	needBreak = false
	if r.conditions.conformsTo(pfa) {
		// bwerror.Spew.Printf("r.conditions: %#v, pfa: %#v\n", r.conditions, pfa)
		for _, pa := range r.processorActions {
			pa.execute(pfa, &isUnexpectedRune, &needFinishTopStackItem, currRune)
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

// const delimiter delimiterRune = delimiterRune{}

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
			result = len(pfa.stack) > 0 && pfa.getTopStackItem().delimiter == r
		} else {
			result = i.HasRune(r)
		}
		if result {
			break
		}
	}
	return result
}

func _expectEOF(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if true {
		needFinishTopStackItem, err = pfa.processStateDef(createStateDef(
			[]interface{}{eofRune{}, setPrimary{expectEOF}},
			[]interface{}{unicodeSpace},
			[]interface{}{unexpectedRune{}},
		))
	} else {
		currRune, isEOF := pfa.currRune()
		switch {
		case isEOF:
			pfa.state.setPrimary(expectEOF)
		case pfa.state.secondary == orSpace && unicode.IsSpace(currRune):
		default:
			err = pfaErrorMake(pfa, unexpectedRuneError)
		}
	}
	return
}

func _expectRocket(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if true {
		needFinishTopStackItem, err = pfa.processStateDef(createStateDef(
			[]interface{}{'>', setPrimary{expectValueOrSpace}},
			[]interface{}{unexpectedRune{}},
		))
	} else {
		currRune, _ := pfa.currRune()
		switch currRune {
		case '>':
			pfa.state.setPrimary(expectValueOrSpace)
		default:
			err = pfaErrorMake(pfa, unexpectedRuneError)
		}
	}
	return
}

func _expectWord(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if true {
		needFinishTopStackItem, err = pfa.processStateDef(createStateDef(
			[]interface{}{unicodeLetter, unicodeDigit,
				appendRuneToItemString{},
			},
			[]interface{}{
				pushRune{},
				needFinish{},
			},
		))
	} else {

		currRune, _ := pfa.currRune()
		switch {
		case unicode.IsLetter(currRune) || currRune == '_' || unicode.IsDigit(currRune):
			stackItem := pfa.getTopStackItem()
			stackItem.itemString = stackItem.itemString + string(currRune)
		default:
			pfa.pushRune()
			needFinishTopStackItem = true
		}
	}
	return
}

func _expectSpaceOrQwItemOrDelimiter(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if true {
		needFinishTopStackItem, err = pfa.processStateDef(createStateDef(
			unexpectedEof,
			[]interface{}{unicodeSpace},
			[]interface{}{delimiterRune{},
				needFinish{},
			},
			[]interface{}{
				pushStackItemQwItem{},
				setPrimary{expectEndOfQwItem},
			},
		))

	} else {

		stackItem := pfa.getTopStackItem()
		currRune, isEOF := pfa.currRune()
		switch {
		case isEOF:
			err = pfaErrorMake(pfa, unexpectedRuneError)
		case unicode.IsSpace(currRune):
		case currRune == stackItem.delimiter:
			needFinishTopStackItem = true
		default:
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemQwItem, start: pfa.curr, itemString: string(*pfa.curr.runePtr), delimiter: stackItem.delimiter})
			pfa.state.setPrimary(expectEndOfQwItem)
		}
	}
	return
}

func _expectEndOfQwItem(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if true {
		needFinishTopStackItem, err = pfa.processStateDef(createStateDef(
			unexpectedEof,
			[]interface{}{unicodeSpace, delimiterRune{},
				pushRune{},
				needFinish{},
			},
			[]interface{}{
				appendRuneToItemString{},
			},
		))
	} else {

		stackItem := pfa.getTopStackItemOfType(parseStackItemQwItem)
		currRune, isEOF := pfa.currRune()
		switch {
		case isEOF:
			err = pfaErrorMake(pfa, unexpectedRuneError)
		case unicode.IsSpace(currRune) || currRune == stackItem.delimiter:
			pfa.pushRune()
			needFinishTopStackItem = true
		default:
			stackItem.itemString += string(currRune)
		}
	}
	return
}

func _expectContentOf(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if true {
		needFinishTopStackItem, err = pfa.processStateDef(createStateDef(
			unexpectedEof,
			[]interface{}{delimiterRune{},
				needFinish{},
			},
			[]interface{}{'\\',
				changePrimary{expectEscapedContentOf},
			},
			[]interface{}{
				appendRuneToItemString{},
			},
		))
	} else {

		stackItem := pfa.getTopStackItem()
		currRune, isEOF := pfa.currRune()
		switch {
		case isEOF:
			err = pfaErrorMake(pfa, unexpectedRuneError)
		case currRune == stackItem.delimiter:
			needFinishTopStackItem = true
		case currRune == '\\':
			pfa.state.primary = expectEscapedContentOf
		default:
			stackItem.itemString = stackItem.itemString + string(currRune)
		}
	}
	return
}

func _expectDigit(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if true {
		needFinishTopStackItem, err = pfa.processStateDef(createStateDef(
			[]interface{}{unicodeDigit, noSecondaryState,
				appendRuneToItemString{},
				changeSecondary{orUnderscoreOrDot},
			},
			[]interface{}{'.', orUnderscoreOrDot,
				appendRuneToItemString{},
				changeSecondary{orUnderscore},
			},
			[]interface{}{'_', unicodeDigit, orUnderscoreOrDot, orUnderscore,
				appendRuneToItemString{},
			},
			[]interface{}{noSecondaryState,
				unexpectedRune{},
			},
			[]interface{}{
				pushRune{},
				needFinish{},
			},
		))
	} else {
		stackItem := pfa.getTopStackItemOfType(parseStackItemNumber)
		currRune, _ := pfa.currRune()
		switch pfa.state.secondary {
		case noSecondaryState:
			switch {
			case unicode.IsDigit(currRune):
				stackItem.itemString = stackItem.itemString + string(currRune)
				pfa.state.secondary = orUnderscoreOrDot
			default:
				err = pfaErrorMake(pfa, unexpectedRuneError)
			}
		case orUnderscoreOrDot:
			switch {
			case currRune == '.':
				pfa.state.secondary = orUnderscore
				stackItem.itemString = stackItem.itemString + string(currRune)
			case unicode.IsDigit(currRune) || currRune == '_':
				stackItem.itemString = stackItem.itemString + string(currRune)
			default:
				pfa.pushRune()
				needFinishTopStackItem = true
			}
		case orUnderscore:
			switch {
			case unicode.IsDigit(currRune) || currRune == '_':
				stackItem.itemString = stackItem.itemString + string(currRune)
			default:
				pfa.pushRune()
				needFinishTopStackItem = true
			}
		}
	}
	return
}

func _expectSpaceOrMapKey(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if true {
		needFinishTopStackItem, err = pfa.processStateDef(createStateDef(
			[]interface{}{unicodeSpace},
			[]interface{}{unicodeLetter,
				pushStackItem{parseStackItemKey},
				// pushStackItemItemKey{false},
				setPrimary{expectWord},
			},
			[]interface{}{'"', '\'',
				pushStackItemWithCurrRuneAsDelimter{parseStackItemKey},
				// pushStackItemItemKey{true},
				setSecondary{expectContentOf, keyToken},
			},
			[]interface{}{',', orMapValueSeparator,
				setPrimary{expectSpaceOrMapKey},
			},
			[]interface{}{delimiterRune{}, parseStackItemMap,
				needFinish{},
			},
			[]interface{}{
				unexpectedRune{},
			},
		))
	} else {
		currRune, _ := pfa.currRune()
		switch {
		case unicode.IsSpace(currRune):
		case unicode.IsLetter(currRune):
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, start: pfa.curr, itemString: string(currRune)})
			pfa.state.setPrimary(expectWord)
		case currRune == '"' || currRune == '\'':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemKey, start: pfa.curr, itemString: ``, delimiter: currRune})
			pfa.state.setSecondary(expectContentOf, keyToken)
		case currRune == ',' && pfa.state.secondary == orMapValueSeparator:
			pfa.state.setPrimary(expectSpaceOrMapKey)
		case pfa.isTopStackItemOfType(parseStackItemMap) && currRune == pfa.getTopStackItem().delimiter:
			needFinishTopStackItem = true
		default:
			err = pfaErrorMake(pfa, unexpectedRuneError)
		}
	}
	return
}

func _expectEscapedContentOf(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	// needFinishTopStackItem, err = pfa.processState(
	// 	nil,
	// 	[]rule{
	// 		rule{
	// 			[]hasRune{runeSetFrom('\'')},
	// 			parseSecondaryStateSet{},
	// 			parseStackItemTypeSet{},
	// 			func(currRune rune) (needFinishTopStackItem bool, isUnexpectedRune bool) {
	// 				if actualVal, ok := singleQuotedEscapedContent[currRune]; ok {
	// 					stackItem.itemString = stackItem.itemString + actualVal
	// 					pfa.state.primary = expectContentOf
	// 				} else {
	// 					isUnexpectedRune = true
	// 				}
	// 				return
	// 			},
	// 		},
	// 	},
	// 	func(currRune rune) (needFinishTopStackItem bool, isUnexpectedRune bool) {
	// 		isUnexpectedRune = true
	// 		return
	// 	},
	// )

	if true {
		needFinishTopStackItem, err = pfa.processStateDef(createStateDef(
			[]interface{}{'"',
				addRune{'"'},
				changePrimary{expectContentOf},
			},
			[]interface{}{'\'',
				addRune{'\''},
				changePrimary{expectContentOf},
			},
			[]interface{}{'\\',
				addRune{'\\'},
				changePrimary{expectContentOf},
			},
			[]interface{}{'a', delimiterIs{'"'},
				addRune{'\a'},
				changePrimary{expectContentOf},
			},
			[]interface{}{'b', delimiterIs{'"'},
				addRune{'\b'},
				changePrimary{expectContentOf},
			},
			[]interface{}{'f', delimiterIs{'"'},
				addRune{'\f'},
				changePrimary{expectContentOf},
			},
			[]interface{}{'n', delimiterIs{'"'},
				addRune{'\n'},
				changePrimary{expectContentOf},
			},
			[]interface{}{'r', delimiterIs{'"'},
				addRune{'\r'},
				changePrimary{expectContentOf},
			},
			[]interface{}{'t', delimiterIs{'"'},
				addRune{'\t'},
				changePrimary{expectContentOf},
			},
			[]interface{}{'v', delimiterIs{'"'},
				addRune{'\v'},
				changePrimary{expectContentOf},
			},
			[]interface{}{
				unexpectedRune{},
			},
		))
	} else {

		stackItem := pfa.getTopStackItem()
		currRune, _ := pfa.currRune()
		if stackItem.delimiter == '\'' {
			if actualVal, ok := singleQuotedEscapedContent[currRune]; ok {
				stackItem.itemString = stackItem.itemString + actualVal
				pfa.state.primary = expectContentOf
			} else {
				err = pfaErrorMake(pfa, unexpectedRuneError)
			}
		} else {
			if actualVal, ok := doubleQuotedEscapedContent[currRune]; ok {
				stackItem.itemString = stackItem.itemString + actualVal
				pfa.state.primary = expectContentOf
			} else {
				err = pfaErrorMake(pfa, unexpectedRuneError)
			}
		}
	}
	return
}

func _expectValueOrSpace(pfa *pfaStruct) (needFinishTopStackItem bool, err error) {
	if true {
		needFinishTopStackItem, err = pfa.processStateDef(createStateDef(
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
				pushStackItemMap{},
				setPrimary{expectSpaceOrMapKey},
			},
			[]interface{}{'<',
				pushStackItemQwWithDelimiter{'>'},
				setPrimary{expectSpaceOrQwItemOrDelimiter},
			},
			[]interface{}{'[',
				pushStackItemArray{},
				setPrimary{expectValueOrSpace},
			},
			[]interface{}{parseStackItemArray, delimiterRune{},
				needFinish{},
			},
			[]interface{}{'-', '+',
				pushStackItem{parseStackItemNumber},
				// pushStackItemNumber{},
				setPrimary{expectDigit},
			},
			[]interface{}{unicodeDigit,
				pushStackItem{parseStackItemNumber},
				// pushStackItemNumber{},
				setSecondary{expectDigit, orUnderscoreOrDot},
			},
			[]interface{}{'"', '\'',
				pullRune{},
				pushStackItemWithCurrRuneAsDelimter{parseStackItemString},
				// pushStackItemString{},
				pushRune{},
				setSecondary{expectContentOf, stringToken},
			},
			[]interface{}{unicodeLetter,
				pushStackItem{parseStackItemWord},
				// pushStackItemWord{},
				setPrimary{expectWord},
			},
			[]interface{}{
				unexpectedRune{},
			},
		))
	} else {
		currRune, isEOF := pfa.currRune()
		switch {
		case isEOF:
			if len(pfa.stack) == 0 {
				pfa.state.setPrimary(expectEOF)
			} else {
				err = pfaErrorMake(pfa, unexpectedRuneError)
			}

		case currRune == '=' && pfa.state.secondary == orMapKeySeparator:
			pfa.state.setPrimary(expectRocket)

		case currRune == ':' && pfa.state.secondary == orMapKeySeparator:
			pfa.state.setPrimary(expectValueOrSpace)

		case currRune == ',' && pfa.state.secondary == orArrayItemSeparator:
			pfa.state.setPrimary(expectValueOrSpace)

		case unicode.IsSpace(currRune):

		case currRune == '{':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemMap, start: pfa.curr, itemMap: map[string]interface{}{}, delimiter: '}'})
			pfa.state.setPrimary(expectSpaceOrMapKey)

		case currRune == '<':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemQw, start: pfa.curr, itemArray: []interface{}{}, delimiter: '>'})
			pfa.state.setPrimary(expectSpaceOrQwItemOrDelimiter)

		case currRune == '[':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemArray, start: pfa.curr, itemArray: []interface{}{}, delimiter: ']'})
			pfa.state.setPrimary(expectValueOrSpace)

		case pfa.isTopStackItemOfType(parseStackItemArray) && pfa.getTopStackItem().delimiter == currRune:
			needFinishTopStackItem = true

		case currRune == '-' || currRune == '+':
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemNumber, start: pfa.curr, itemString: string(currRune)})
			pfa.state.setPrimary(expectDigit)

		case unicode.IsDigit(currRune):
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemNumber, start: pfa.curr, itemString: string(currRune)})
			pfa.state.setSecondary(expectDigit, orUnderscoreOrDot)

		case currRune == '"' || currRune == '\'':
			pfa.pullRune()
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemString, start: pfa.curr, itemString: ``, delimiter: currRune})
			pfa.pushRune()
			pfa.state.setSecondary(expectContentOf, stringToken)

		case unicode.IsLetter(currRune):
			pfa.stack = append(pfa.stack, parseStackItem{itemType: parseStackItemWord, start: pfa.curr, itemString: string(currRune)})
			pfa.state.setPrimary(expectWord)

		default:
			err = pfaErrorMake(pfa, unexpectedRuneError)
		}
	}
	return
}

var singleQuotedEscapedContent = map[rune]string{
	'"':  "\"",
	'\'': "'",
	'\\': "\\",
}

var doubleQuotedEscapedContent = map[rune]string{
	'"':  "\"",
	'\'': "'",
	'\\': "\\",
	'a':  "\a",
	'b':  "\b",
	'f':  "\f",
	'n':  "\n",
	'r':  "\r",
	't':  "\t",
	'v':  "\v",
}

// func getRunesOfEscapedContent() (result []rune) {
// 	result = []rune{}
// 	for r, _ := range escapedContent {
// 		result = append(result, r)
// 	}
// 	return
// }
