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
	"github.com/baza-winner/bwcore/bwset"
	"github.com/jimlawless/whereami"
)

// ============================== hasRune =====================================

type hasRune interface {
	HasRune(pfa *pfaStruct) bool
	Len() int
}

type runeSet struct {
	values bwset.RuneSet
}

func (v runeSet) HasRune(pfa *pfaStruct) (result bool) {
	r, isEOF := pfa.p.Rune()
	result = !isEOF && v.values.Has(r)
	return result
}

func (v runeSet) Len() int {
	return len(v.values)
}

type currRuneVarPaths struct {
	values []VarPath
}

func (v currRuneVarPaths) HasRune(pfa *pfaStruct) (result bool) {
	currRune, isEOF := pfa.p.Rune()
	if !isEOF {
		for _, k := range v.values {
			if r, err := pfa.getVarValue(k).AsRune(); err == nil {
				if r == currRune {
					result = true
					break
				}
			}
		}
	}
	return
}

func (v currRuneVarPaths) Len() int {
	return len(v.values)
}

type UnicodeCategory uint8

const (
	IsUnicodeSpace UnicodeCategory = iota
	IsUnicodeLetter
	IsUnicodeLower
	IsUnicodeUpper
	IsUnicodeDigit
	IsUnicodeOpenBraces
	IsUnicodePunct
	IsUnicodeSymbol
)

func (v UnicodeCategory) HasRune(pfa *pfaStruct) (result bool) {
	r, isEOF := pfa.p.Rune()
	if !isEOF {
		switch v {
		case IsUnicodeSpace:
			result = unicode.IsSpace(r)
		case IsUnicodeLetter:
			result = unicode.IsLetter(r) || r == '_'
		case IsUnicodeLower:
			result = unicode.IsLower(r)
		case IsUnicodeUpper:
			result = unicode.IsUpper(r)
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
	}
	return
}

func (v UnicodeCategory) Len() int {
	return 1
}

type IsEOF struct{}

func (v IsEOF) HasRune(pfa *pfaStruct) (result bool) {
	_, isEOF := pfa.p.Rune()
	return isEOF
}

func (v IsEOF) Len() int {
	return 1
}

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

func (v justVal) GetVal(pfa *pfaStruct) (val interface{}, err error) {
	val = v.val
	return
}

type varVal struct {
	varPath VarPath
}

func (v varVal) GetVal(pfa *pfaStruct) (val interface{}, err error) {
	varValue := pfa.getVarValue(v.varPath)
	if varValue.Err == nil {
		val = varValue.Val
	} else {
		err = varValue.Err
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
	val, err := v.valProvider.GetVal(pfa)
	if err == nil {
		for _, b := range v.by {
			val = b.TransformValue(pfa, val)
			if pfa.err != nil {
				break
			}
		}
		if err == nil {
			if !v.needAppend {
				err = pfa.setVarVal(v.varPath, val)
			} else {
				if orig := pfa.getVarValue(v.varPath); orig.Err != nil {
					err = orig.Err
				} else if s, ok := orig.Val.(string); ok {
					if a, ok := val.(string); ok {
						val = s + a
					} else if r, ok := val.(rune); ok {
						val = s + string(r)
					} else {
						err = bwerror.Error("%#v expected to be string or rune", val)
					}
				} else if reflect.TypeOf(orig.Val).Kind() == reflect.Slice {
					valueOfOrigVal := reflect.ValueOf(orig.Val)
					valueOfVal := reflect.ValueOf(val)
					if !v.appendSlice {
						val = reflect.Append(valueOfOrigVal, valueOfVal).Interface()
					} else if reflect.TypeOf(val).Kind() == reflect.Slice {
						val = reflect.AppendSlice(valueOfOrigVal, valueOfVal).Interface()
					} else {
						err = bwerror.Error("%#v expected to be slice", val)
					}

					// if reflect.TypeOf(val).Kind() == reflect.Slice {
					// 	val = reflect.AppendSlice(valueOfOrigVal, valueOfVal).Interface()
					// } else {
					// 	val = reflect.Append(valueOfOrigVal, valueOfVal).Interface()
					// }
				} else {
					err = bwerror.Error(
						"value (<ansiPrimary>%#v<ansi>) at <ansiCmd>%s<ansi> expected to be <ansiSecondary>string<ansi> or <ansiSecondary>slice<ansi> to be <ansiOutline>Appendable",
						val, v.varPath,
					)
				}
				if err == nil {
					err = pfa.setVarVal(v.varPath, val)
				}
			}
		}
	}
	pfa.err = err
}

type ValTransformer interface {
	TransformValue(pfa *pfaStruct, i interface{}) interface{}
}

type By []ValTransformer

type _setVar struct {
	varPath     VarPath
	valProvider ValProvider
}

type ValProvider interface {
	GetVal(pfa *pfaStruct) (val interface{}, err error)
}

func ValProviderFrom(i interface{}) (result ValProvider, err error) {
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
	if result, err = ValProviderFrom(i); err != nil {
		bwerror.PanicErr(err)
	}
	return
}

func (v _setVar) execute(pfa *pfaStruct) {
	val, err := v.valProvider.GetVal(pfa)
	if err == nil {
		err = pfa.setVarVal(v.varPath, val)
	}
	pfa.err = err
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
