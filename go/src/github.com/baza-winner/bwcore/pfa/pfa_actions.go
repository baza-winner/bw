package pfa

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwint"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/jimlawless/whereami"
)

// ============================== hasRune =====================================

type UnicodeCategory uint8

const (
	UnicodeSpace UnicodeCategory = iota
	UnicodeLetter
	UnicodeLower
	UnicodeUpper
	UnicodeDigit
	UnicodeOpenBraces
	UnicodePunct
	UnicodeSymbol
)

func (v UnicodeCategory) conforms(pfa *pfaStruct, val interface{}) (result bool) {
	if r, ok := val.(rune); ok {
		switch v {
		case UnicodeSpace:
			result = unicode.IsSpace(r)
		case UnicodeLetter:
			result = unicode.IsLetter(r) || r == '_'
		case UnicodeLower:
			result = unicode.IsLower(r)
		case UnicodeUpper:
			result = unicode.IsUpper(r)
		case UnicodeDigit:
			result = unicode.IsDigit(r)
		case UnicodeOpenBraces:
			result = r == '(' || r == '{' || r == '[' || r == '<'
		case UnicodePunct:
			result = unicode.IsPunct(r)
		case UnicodeSymbol:
			result = unicode.IsSymbol(r)
		default:
			bwerror.Panic("UnicodeCategory: %s", v)
		}
	}
	return
}

// func (v UnicodeCategory) Len() int {
// 	return 1
// }

type EOF struct{}

// func (v IsEOF) HasRune(pfa *pfaStruct) (result bool) {
// 	_, isEOF := pfa.p.Rune()
// 	return isEOF
// }

// func (v IsEOF) Len() int {
// 	return 1
// }

// ============================================================================

type VarIs struct {
	VarPathStr string
	VarValue   interface{}
}

// ======================== processorAction ====================================

type processorAction interface {
	execute(pfa *pfaStruct)
}

type PushRune struct{}

func (v PushRune) execute(pfa *pfaStruct) {
	pfa.p.PushRune()
}

type PullRune struct{}

func (v PullRune) execute(pfa *pfaStruct) {
	pfa.p.PullRune()
}

type PushItem struct{}

func (v PushItem) execute(pfa *pfaStruct) {
	pfa.pushStackItem()
}

type SubRules struct {
	Def Rules
}

func (v SubRules) execute(pfa *pfaStruct) {
	pfa.processRules(v.Def)
}

type PopItem struct{}

func (v PopItem) execute(pfa *pfaStruct) {
	pfa.popStackItem()
}

type Var struct {
	VarPathStr string
}

type justVal struct {
	val interface{}
}

func (v justVal) conforms(pfa *pfaStruct, val interface{}) bool {
	return val == v.val
}

func (v justVal) getVal(pfa *pfaStruct) interface{} {
	return v.val
}

type varVal struct {
	varPath VarPath
}

func (v varVal) conforms(pfa *pfaStruct, val interface{}) (result bool) {
	varValue := pfa.getVarValue(v.varPath)
	if pfa.err == nil {
		result = varValue.val == val
	}
	return
}

func (v varVal) getVal(pfa *pfaStruct) (result interface{}) {
	varValue := pfa.getVarValue(v.varPath)
	if pfa.err == nil {
		result = varValue.val
	}
	return
}

type ProccessorActionProvider interface {
	GetAction() processorAction
}

type SetVar struct {
	VarPathStr string
	VarValue   interface{}
}

func (v SetVar) GetAction() processorAction {
	return _setVar{
		MustVarPathFrom(v.VarPathStr),
		MustValProviderFrom(v.VarValue),
	}
}

type SetVarBy struct {
	VarPathStr   string
	VarValue     interface{}
	Transformers By
}

func (v SetVarBy) GetAction() processorAction {
	by := By{}
	var needAppend, appendSlice bool
	for _, b := range v.Transformers {
		if _, ok := b.(Append); ok {
			needAppend = true
		} else if _, ok := b.(AppendSlice); ok {
			needAppend = true
			appendSlice = true
		} else {
			by = append(by, b)
		}
	}
	return _setVarBy{
		MustVarPathFrom(v.VarPathStr),
		MustValProviderFrom(v.VarValue),
		by,
		needAppend,
		appendSlice,
	}
}

type _setVarBy struct {
	varPath     VarPath
	valProvider ValProvider
	by          By
	needAppend  bool
	appendSlice bool
}

func (v _setVarBy) execute(pfa *pfaStruct) {
	val := v.valProvider.getVal(pfa)
	if pfa.err == nil {
		for _, b := range v.by {
			val = b.TransformValue(pfa, val)
			if pfa.err != nil {
				break
			}
		}
		if pfa.err == nil {
			if !v.needAppend {
				pfa.setVarVal(v.varPath, val)
			} else {
				if orig := pfa.getVarValue(v.varPath); pfa.err == nil {
					if s, ok := orig.val.(string); ok {
						if a, ok := val.(string); ok {
							val = s + a
						} else if r, ok := val.(rune); ok {
							val = s + string(r)
						} else {
							pfa.err = bwerror.Error("%#v expected to be string or rune", val)
						}
					} else if reflect.TypeOf(orig.val).Kind() == reflect.Slice {
						valueOfOrigVal := reflect.ValueOf(orig.val)
						valueOfVal := reflect.ValueOf(val)
						if !v.appendSlice {
							val = reflect.Append(valueOfOrigVal, valueOfVal).Interface()
						} else if reflect.TypeOf(val).Kind() == reflect.Slice {
							val = reflect.AppendSlice(valueOfOrigVal, valueOfVal).Interface()
						} else {
							pfa.err = bwerror.Error("%#v expected to be slice", val)
						}
					} else {
						pfa.err = bwerror.Error(
							"value (<ansiPrimary>%#v<ansi>) at <ansiCmd>%s<ansi> expected to be <ansiSecondary>string<ansi> or <ansiSecondary>slice<ansi> to be <ansiOutline>Appendable",
							val, v.varPath,
						)
					}
				}
				if pfa.err == nil {
					pfa.setVarVal(v.varPath, val)
				}
			}
		}
	}
}

type ValTransformer interface {
	TransformValue(pfa *pfaStruct, i interface{}) interface{}
}

type By []ValTransformer

type _setVar struct {
	varPath     VarPath
	valProvider ValProvider
}

type ValChecker interface {
	conforms(pfa *pfaStruct, val interface{}) bool
}

type ValProvider interface {
	getVal(pfa *pfaStruct) interface{}
}

func valProviderFrom(i interface{}) (result ValProvider, err error) {
	if v, ok := i.(Var); !ok {
		result = justVal{i}
	} else {
		var varPath VarPath
		varPath, err = VarPathFrom(v.VarPathStr)
		if err == nil {
			result = varVal{varPath}
		}
	}
	return
}

func MustValProviderFrom(i interface{}) (result ValProvider) {
	var err error
	if result, err = valProviderFrom(i); err != nil {
		bwerror.PanicErr(err)
	}
	return
}

func (v _setVar) execute(pfa *pfaStruct) {
	val := v.valProvider.getVal(pfa)
	if pfa.err == nil {
		pfa.setVarVal(v.varPath, val)
	}
}

type Debug struct{ Message string }

func (v *Debug) execute(pfa *pfaStruct) {
	fmt.Printf("%s: %s\n", v.Message, bwjson.PrettyJsonOf(pfa))
}

// ============================================================================

type ParseNumber struct{}

func (v ParseNumber) TransformValue(pfa *pfaStruct, i interface{}) (result interface{}) {
	if s, ok := i.(string); !ok {
		pfa.err = bwerror.Error("ParseNumber expects string for TransformValue, not %#v", i)
	} else {
		result, pfa.err = _parseNumber(s)
		if pfa.err != nil {
			stackItem := pfa.getTopStackItem()
			errStr := pfa.p.WordError("failed to get number from string <ansiPrimary>%s<ansi>", s, stackItem.start).Error()
			pfa.err = pfaError{pfa, "failedToGetNumber", errStr, whereami.WhereAmI(2)}
		}
	}
	return
}

type Append struct{}

func (v Append) TransformValue(pfa *pfaStruct, i interface{}) (result interface{}) {
	bwerror.Unreachable()
	return
}

type AppendSlice struct{}

func (v AppendSlice) TransformValue(pfa *pfaStruct, i interface{}) (result interface{}) {
	bwerror.Unreachable()
	return
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
