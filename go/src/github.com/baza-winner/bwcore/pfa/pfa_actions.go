package pfa

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwint"
	"github.com/baza-winner/bwcore/bwset"
)

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

type UnicodeCategory uint8

const (
	IsUnicodeSpace UnicodeCategory = iota
	IsUnicodeLetter
	IsUnicodeDigit
	IsUnicodeOpenBraces
	IsUnicodePunct
	IsUnicodeSymbol
)

func (v UnicodeCategory) HasRune(r rune) (result bool) {
	switch v {
	case IsUnicodeSpace:
		result = unicode.IsSpace(r)
	case IsUnicodeLetter:
		result = unicode.IsLetter(r) || r == '_'
	case IsUnicodeDigit:
		result = unicode.IsDigit(r)
	case IsUnicodeOpenBraces:
		result = r == '(' || r == '{' || r == '[' || r == '<'
	case IsUnicodePunct:
		result = unicode.IsPunct(r)
	case IsUnicodeSymbol:
		result = unicode.IsSymbol(r)
	default:
		bwerror.Panic("UnicodeCategory: %s", v)
	}
	return
}

func (v UnicodeCategory) Len() int {
	return 1
}

type hasRuneSlice []hasRune

type IsEOF struct{}

func (v IsEOF) HasRune(r rune) (result bool) {
	bwerror.Panic("Unreachable")
	return
}

func (v IsEOF) Len() int {
	return 1
}

func (v hasRuneSlice) HasRune(pfa *pfaStruct, r rune) (result bool) {
	result = false
	for _, i := range v {
		if _, ok := i.(IsDelimiterRune); ok {
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

type IsDelimiterRune struct{}

func (v IsDelimiterRune) HasRune(r rune) (result bool) {
	bwerror.Panic("Unreachable")
	return
}

func (v IsDelimiterRune) Len() int {
	return 1
}

// ============================================================================

type DelimiterIs struct {
	Value rune
}

type StackLenIs struct {
	Value int
}

type VarIs struct {
	VarName  string
	VarValue interface{}
}

type TopItemIs struct {
	Value string
}

type SubItemIs struct {
	Value string
}

type TopItemStringIs struct {
	Value string
}

type TopItemStringIsOneOf struct {
	Value bwset.StringSet
}

type PrimaryIs struct {
	Value string
}

type SecondaryIs struct {
	Value string
}

// ============== DelimiterProvider, ItemStringProvider ========================

type DelimiterProvider interface {
	Delimiter(pfa *pfaStruct) *rune
}

type ItemStringProvider interface {
	ItemString(pfa *pfaStruct) string
}

type Delim struct{ r rune }

func (v Delim) Delimiter(pfa *pfaStruct, currRune rune) *rune {
	return &v.r
}

type FromParentItem struct{}

func (v FromParentItem) Delimiter(pfa *pfaStruct) *rune {
	return pfa.getTopStackItem().delimiter
}

type FromCurrRune struct{}

func (v FromCurrRune) ItemString(pfa *pfaStruct) string {
	return string(pfa.mustCurrRune())
}

func (v FromCurrRune) Delimiter(pfa *pfaStruct) *rune {
	currRune := pfa.mustCurrRune()
	return &currRune
}

type PairForCurrRune struct{}

func (v PairForCurrRune) Delimiter(pfa *pfaStruct) *rune {
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

// ======================== processorAction ====================================

type processorAction interface {
	execute(pfa *pfaStruct)
}

type SetError struct {
	Value ErrorType
}

func (v SetError) execute(pfa *pfaStruct) {
	pfa.errorType = v.Value
}

type PushRune struct{}

func (v PushRune) execute(pfa *pfaStruct) {
	pfa.PushRune()
}

type PullRune struct{}

func (v PullRune) execute(pfa *pfaStruct) {
	pfa.PullRune()

}

type ChangePrimary struct {
	Value string
}

func (v ChangePrimary) execute(pfa *pfaStruct) {
	pfa.state.Primary = v.Value
}

type SetPrimary struct {
	Value string
}

func (v SetPrimary) execute(pfa *pfaStruct) {
	pfa.state.SetPrimary(v.Value)
}

type ChangeSecondary struct {
	Value string
}

func (v ChangeSecondary) execute(pfa *pfaStruct) {
	pfa.state.Secondary = v.Value
}

type SetSecondary struct {
	Primary   string
	Secondary string
}

func (v SetSecondary) execute(pfa *pfaStruct) {
	pfa.state.SetSecondary(v.Primary, v.Secondary)
}

type AppendCurrRune struct{}

func (v AppendCurrRune) execute(pfa *pfaStruct) {
	pfa.getTopStackItem().itemString += string(pfa.mustCurrRune())
}

type AppendRune struct {
	Value rune
}

func (v AppendRune) execute(pfa *pfaStruct) {
	pfa.getTopStackItem().itemString += string(v.Value)
}

type SetTopItemDelimiter struct {
	P DelimiterProvider
}

func (v SetTopItemDelimiter) execute(pfa *pfaStruct) {
	pfa.getTopStackItem().delimiter = v.P.Delimiter(pfa)
}

type SetTopItemType struct {
	Value string
}

func (v SetTopItemType) execute(pfa *pfaStruct) {
	pfa.getTopStackItem().itemType = v.Value
}

type PushItem struct {
	ItemType   string
	ItemString ItemStringProvider
	Delimiter  DelimiterProvider
}

func (v PushItem) execute(pfa *pfaStruct) {
	var itemString string
	if v.ItemString != nil {
		itemString = v.ItemString.ItemString(pfa)
	}
	var delimiter *rune
	if v.Delimiter != nil {
		delimiter = v.Delimiter.Delimiter(pfa)
	}
	pfa.pushStackItem(
		v.ItemType,
		itemString,
		delimiter,
	)
}

type SetTopItemValueAsString struct{}

func (v SetTopItemValueAsString) execute(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemString
}

type SetTopItemValueAsMap struct{}

func (v SetTopItemValueAsMap) execute(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemMap
}

type SetTopItemValueAsArray struct{}

func (v SetTopItemValueAsArray) execute(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemArray
}

type SetTopItemValueAsNumber struct{}

func (v SetTopItemValueAsNumber) execute(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	var err error
	if stackItem.value, err = _parseNumber(stackItem.itemString); err != nil {
		pfa.errorType = FailedToGetNumber
	}
}

type SetTopItemValueAsBool struct {
	Value bool
}

func (v SetTopItemValueAsBool) execute(pfa *pfaStruct) {
	pfa.getTopStackItem().value = v.Value
}

type SubLogic struct {
	Def *LogicDef
}

func (v SubLogic) execute(pfa *pfaStruct) {
	pfa.processLogic(v.Def)
}

type PopSubItem struct{}

func (v PopSubItem) execute(pfa *pfaStruct) {
	pfa.stackSubItem = pfa.popStackItem()
}

type AppendItemArray struct {
	P ArrayProvider
}

func (v AppendItemArray) execute(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	stackItem.itemArray = append(stackItem.itemArray, v.P.GetArray(pfa)...)
}

type SetTopItemStringFromSubItem struct{}

func (v SetTopItemStringFromSubItem) execute(pfa *pfaStruct) {
	pfa.getTopStackItem().itemString = pfa.stackSubItem.itemString
}

type SetTopItemMapKeyValueFromSubItem struct{}

func (v SetTopItemMapKeyValueFromSubItem) execute(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	stackItem.itemMap[stackItem.itemString] = pfa.stackSubItem.value
}

type Unreachable struct{}

func (v Unreachable) execute(pfa *pfaStruct) {
	pfa.panic("unreachabe")
}

type SetVar struct {
	VarName  string
	VarValue interface{}
}

func (v SetVar) execute(pfa *pfaStruct) {
	pfa.vars[v.VarName] = v.VarValue
}

// ============================================================================

type ArrayProvider interface {
	GetArray(pfa *pfaStruct) []interface{}
}

type FromSubItemValue struct{}

func (v FromSubItemValue) GetArray(pfa *pfaStruct) []interface{} {
	return []interface{}{pfa.stackSubItem.value}
}

type FromSubItemArray struct{}

func (v FromSubItemArray) GetArray(pfa *pfaStruct) []interface{} {
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
