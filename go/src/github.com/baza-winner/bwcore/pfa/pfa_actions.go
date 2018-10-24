package pfa

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/baza-winner/bwcore/ansi"
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

func (v UnicodeCategory) conforms(pfa *pfaStruct, val interface{}, varPath VarPath) (result bool) {
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
	pfa.traceCondition(varPath, v, result)
	return
}

type EOF struct{}

func (v EOF) String() string {
	return "EOF"
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
	pfa.traceAction("<ansiGreen>PushRune<ansi>: %s", runeVarPath)
}

func (v PushRune) GetAction() processorAction {
	return v
}

type PullRune struct{}

var runeVarPath VarPath = MustVarPathFrom("rune")
var stackLenVarPath VarPath = MustVarPathFrom("stackLen")

func (v PullRune) execute(pfa *pfaStruct) {
	pfa.p.PullRune()
	pfa.traceAction("<ansiGreen>PullRune<ansi>: %s", runeVarPath)
}

func (v PullRune) GetAction() processorAction {
	return v
}

type PushItem struct{}

func (v PushItem) execute(pfa *pfaStruct) {
	pfa.pushStackItem()
	// pfa.traceAction("<ansiGreen>PushItem<ansi>: <ansiCmd>stackLen<ansi>(<ansiPrimary>%d<ansi>)", len(pfa.stack))
	pfa.traceAction("<ansiGreen>PushItem<ansi>: %s", stackLenVarPath)
}

func (v PushItem) GetAction() processorAction {
	return v
}

type SubRules struct {
	Def Rules
}

func (v SubRules) GetAction() processorAction {
	return v
}

func (v SubRules) execute(pfa *pfaStruct) {
	pfa.traceIncLevel()
	pfa.processRules(v.Def)
	pfa.traceDecLevel()
}

type PopItem struct{}

func (v PopItem) execute(pfa *pfaStruct) {
	pfa.popStackItem()
	pfa.traceAction("<ansiGreen>PopItem<ansi>: %s", stackLenVarPath)
	// pfa.traceAction("<ansiGreen>PopItem<ansi>: <ansiCmd>stackLen<ansi>(<ansiPrimary>%d<ansi>)", len(pfa.stack))
}

func (v PopItem) GetAction() processorAction {
	return v
}

type Var struct {
	VarPathStr string
}

type ValProvider interface {
	getVal(pfa *pfaStruct) interface{}
	getSource(pfa *pfaStruct) formattedString
}

type Map struct{}

func (v Map) getVal(pfa *pfaStruct) interface{} {
	return map[string]interface{}{}
}

func (v Map) getSource(pfa *pfaStruct) formattedString {
	return pfa.traceVal(v)
}

func (v Map) String() string {
	return "Map"
}

type Array struct{}

func (v Array) getVal(pfa *pfaStruct) interface{} {
	return []interface{}{}
}

func (v Array) getSource(pfa *pfaStruct) formattedString {
	return pfa.traceVal(v)
}

func (v Array) String() string {
	return "Array"
}

type UnexpectedRune struct{}

func (v UnexpectedRune) getVal(pfa *pfaStruct) interface{} {
	return v
}

func (v UnexpectedRune) getSource(pfa *pfaStruct) formattedString {
	return pfa.traceVal(v)
}

func (v UnexpectedRune) String() string {
	return "UnexpectedRune"
}

type UnknownWord struct{}

func (v UnknownWord) getVal(pfa *pfaStruct) interface{} {
	return v
}

func (v UnknownWord) getSource(pfa *pfaStruct) formattedString {
	return pfa.traceVal(v)
}

func (v UnknownWord) String() string {
	return "UnknownWord"
}

type Panic struct{}

func (v Panic) getVal(pfa *pfaStruct) interface{} {
	return v
}

func (v Panic) getSource(pfa *pfaStruct) formattedString {
	return pfa.traceVal(v)
}

func (v Panic) String() string {
	return "Panic"
}

type justVal struct {
	val interface{}
}

func (v justVal) conforms(pfa *pfaStruct, val interface{}, varPath VarPath) (result bool) {
	result = val == v.val
	pfa.traceCondition(varPath, val, result)
	return
}

func (v justVal) getVal(pfa *pfaStruct) interface{} {
	return v.val
}

func (v justVal) getSource(pfa *pfaStruct) formattedString {
	return pfa.traceVal(v.val)
}

type varVal struct {
	varPath VarPath
}

func (v varVal) conforms(pfa *pfaStruct, val interface{}, varPath VarPath) (result bool) {
	varValue := pfa.getVarValue(v.varPath)
	if pfa.err == nil {
		result = varValue.val == val
		pfa.traceCondition(varPath, v.varPath, result)
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

func (v varVal) getSource(pfa *pfaStruct) formattedString {
	// return fmt.Sprintf(ansi.Ansi("", "%s"), pfa.traceVarPath(v.varPath))
	return pfa.traceVal(v.varPath)
}

type ProccessorActionProvider interface {
	GetAction() processorAction
}

type SetVar struct {
	VarPathStr string
	VarValue   interface{}
}

func (v SetVar) GetAction() processorAction {
	return _setVarBy{
		MustVarPathFrom(v.VarPathStr),
		MustValProviderFrom(v.VarValue),
		By{},
		noAppend,
	}
}

type SetVarBy struct {
	VarPathStr   string
	VarValue     interface{}
	Transformers By
}

func (v SetVarBy) GetAction() processorAction {
	by := By{}
	at := noAppend
	for _, b := range v.Transformers {
		if _, ok := b.(Append); ok {
			if at == noAppend {
				at = appendScalar
			}
		} else if _, ok := b.(AppendSlice); ok {
			at = appendSlice
		} else {
			by = append(by, b)
		}
	}
	return _setVarBy{
		MustVarPathFrom(v.VarPathStr),
		MustValProviderFrom(v.VarValue),
		by,
		at,
	}
}

type appendType uint8

const (
	noAppend appendType = iota
	appendScalar
	appendSlice
)

type _setVarBy struct {
	varPath     VarPath
	valProvider ValProvider
	by          By
	at          appendType
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
			var source formattedString
			var target formattedString
			op := ""
			if pfa.traceLevel > TraceNone {
				source = v.valProvider.getSource(pfa)
				target = pfa.traceVal(v.varPath)
			}
			if v.at == noAppend {
				pfa.setVarVal(v.varPath, val)
				op = ">"
			} else {
				op = ">>"
				if orig := pfa.getVarValue(v.varPath); pfa.err == nil {
					if orig.val == nil {
						if a, ok := val.(string); ok {
							val = a
						} else if r, ok := val.(rune); ok {
							val = string(r)
						} else if v.at == appendScalar {
							val = []interface{}{val}
						} else if reflect.TypeOf(val).Kind() != reflect.Slice {
							pfa.err = bwerror.Error("%#v expected to be slice", val)
						} else {
							op = ">>>"
						}

					} else if s, ok := orig.val.(string); ok {
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
						if v.at == appendScalar {
							val = reflect.Append(valueOfOrigVal, valueOfVal).Interface()
						} else if reflect.TypeOf(val).Kind() == reflect.Slice {
							val = reflect.AppendSlice(valueOfOrigVal, valueOfVal).Interface()
							op = ">>>"
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
			if pfa.err == nil {
				pfa.traceAction("%s %s %s: %s", source, formattedString(ansi.Ansi("Green", op)), target, v.varPath)
			}
		}
	}
}

type ValTransformer interface {
	TransformValue(pfa *pfaStruct, i interface{}) interface{}
}

type By []ValTransformer

type ValChecker interface {
	conforms(pfa *pfaStruct, val interface{}, varPath VarPath) bool
}

func valProviderFrom(i interface{}) (result ValProvider, err error) {
	switch t := i.(type) {
	case Var:
		var varPath VarPath
		varPath, err = VarPathFrom(t.VarPathStr)
		if err == nil {
			result = varVal{varPath}
		}
	case Array:
		result = Array{}
	case Map:
		result = Map{}
	default:
		result = justVal{i}
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
			err := pfa.p.WordError("failed to get number from string <ansiPrimary>%s<ansi>", s, stackItem.start)
			pfa.err = pfaError{pfa, "failedToGetNumber", err, whereami.WhereAmI(2)}
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
