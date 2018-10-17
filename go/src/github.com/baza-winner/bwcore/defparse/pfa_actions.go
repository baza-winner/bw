package defparse

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/baza-winner/bwcore/bwint"
)

// ======================== processorAction ====================================

type processorAction interface {
	execute(pfa *pfaStruct)
}

type setError struct{ errorType pfaErrorType }

func (v setError) execute(pfa *pfaStruct) {
	pfa.errorType = v.errorType
}

type pushRune struct{}

func (v pushRune) execute(pfa *pfaStruct) {
	pfa.pushRune()
}

type pullRune struct{}

func (v pullRune) execute(pfa *pfaStruct) {
	pfa.pullRune()

}

type changePrimary struct {
	primary parsePrimaryState
}

func (v changePrimary) execute(pfa *pfaStruct) {
	pfa.state.primary = v.primary
}

type setPrimary struct {
	primary parsePrimaryState
}

func (v setPrimary) execute(pfa *pfaStruct) {
	pfa.state.setPrimary(v.primary)
}

type changeSecondary struct {
	secondary parseSecondaryState
}

func (v changeSecondary) execute(pfa *pfaStruct) {
	pfa.state.secondary = v.secondary
}

type setSecondary struct {
	primary   parsePrimaryState
	secondary parseSecondaryState
}

func (v setSecondary) execute(pfa *pfaStruct) {
	pfa.state.setSecondary(v.primary, v.secondary)
}

type appendCurrRune struct{}

func (v appendCurrRune) execute(pfa *pfaStruct) {
	pfa.getTopStackItem().itemString += string(pfa.mustCurrRune())
}

type appendRune struct {
	r rune
}

func (v appendRune) execute(pfa *pfaStruct) {
	pfa.getTopStackItem().itemString += string(v.r)
}

type setTopItemDelimiter struct{ d delimiterProvider }

func (v setTopItemDelimiter) execute(pfa *pfaStruct) {
	pfa.getTopStackItem().delimiter = v.d.Delimiter(pfa)
}

type setTopItemType struct{ itemType parseStackItemType }

func (v setTopItemType) execute(pfa *pfaStruct) {
	pfa.getTopStackItem().itemType = v.itemType
}

type pushItem struct {
	itemType   parseStackItemType
	itemString itemStringProvider
	delimiter  delimiterProvider
}

func (v pushItem) execute(pfa *pfaStruct) {
	var itemString string
	if v.itemString != nil {
		itemString = v.itemString.ItemString(pfa)
	}
	var delimiter *rune
	if v.delimiter != nil {
		delimiter = v.delimiter.Delimiter(pfa)
	}
	pfa.pushStackItem(
		v.itemType,
		itemString,
		delimiter,
	)
}

type setTopItemValueAsString struct{}

func (v setTopItemValueAsString) execute(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemString
}

type setTopItemValueAsMap struct{}

func (v setTopItemValueAsMap) execute(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemMap
}

type setTopItemValueAsArray struct{}

func (v setTopItemValueAsArray) execute(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	stackItem.value = stackItem.itemArray
}

type setTopItemValueAsNumber struct{}

func (v setTopItemValueAsNumber) execute(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	var err error
	if stackItem.value, err = _parseNumber(stackItem.itemString); err != nil {
		pfa.errorType = failedToGetNumberError
	}
}

type setTopItemValueAsBool struct{ b bool }

func (v setTopItemValueAsBool) execute(pfa *pfaStruct) {
	pfa.getTopStackItem().value = v.b
}

type processStateDef struct{ def *stateDef }

func (v processStateDef) execute(pfa *pfaStruct) {
	pfa.processStateDef(v.def)
	if pfa.errorType == unexpectedRuneError {
		pfa.panic("SOME")
	}
}

type popSubItem struct{}

func (v popSubItem) execute(pfa *pfaStruct) {
	pfa.stackSubItem = pfa.popStackItem()
}

type appendItemArray struct{ p arrayProvider }

func (v appendItemArray) execute(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	stackItem.itemArray = append(stackItem.itemArray, v.p.GetArray(pfa)...)
}

type setTopItemStringFromSubItem struct{}

func (v setTopItemStringFromSubItem) execute(pfa *pfaStruct) {
	pfa.getTopStackItem().itemString = pfa.stackSubItem.itemString
}

type setTopItemMapKeyValueFromSubItem struct{}

func (v setTopItemMapKeyValueFromSubItem) execute(pfa *pfaStruct) {
	stackItem := pfa.getTopStackItem()
	stackItem.itemMap[stackItem.itemString] = pfa.stackSubItem.value
}

type unreachable struct{}

func (v unreachable) execute(pfa *pfaStruct) {
	pfa.panic("unreachabe")
}

type setVar struct {
	varName  string
	varValue interface{}
}

func (v setVar) execute(pfa *pfaStruct) {
	pfa.vars[v.varName] = v.varValue
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
