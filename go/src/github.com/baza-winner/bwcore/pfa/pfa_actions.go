package pfa

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwint"
	"github.com/baza-winner/bwcore/bwjson"
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

type ProccessorActionProvider interface {
	GetAction() processorAction
}

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

func (v UnexpectedRune) GetAction() processorAction {
	return v
}
func (v UnexpectedRune) execute(pfa *pfaStruct) {
	// pfa.isUnexpectedRune = true
	pfa.errVal = v
	pfa.traceAction("<ansiGreen>UnexpectedRune")
	// word, pfa.err = pfa.getVarValue(v.varPath).String()
	// if pfa.err == nil {
	// 	pfa.err = pfa.p.WordError(
	// 		"unknown word <ansiPrimary>%s<ansi>",
	// 		word,
	// 		pfa.getTopStackItem().start,
	// 	)
	// 	pfa.traceAction("<ansiGreen>UnexpectedItem{%s}", v.varPath)
	// }
}

// func (v UnexpectedRune) getVal(pfa *pfaStruct) interface{} {
// 	return v
// }

// func (v UnexpectedRune) getSource(pfa *pfaStruct) formattedString {
// 	return pfa.traceVal(v)
// }

// func (v UnexpectedRune) String() string {
// 	return "UnexpectedRune"
// }

type UnexpectedItem struct {
	// VarPathStr string
}

func (v UnexpectedItem) GetAction() processorAction {
	return v
	// return unexpectedItemAction{MustVarPathFrom(v.VarPathStr)}
}

// type unexpectedItemAction struct{ varPath VarPath }

func (v UnexpectedItem) execute(pfa *pfaStruct) {
	// var word string
	// word, pfa.err = pfa.getVarValue(v.varPath).String()
	// if pfa.err == nil {
	// pfa.err = pfa.p.ItemError(
	// 	pfa.getTopStackItem().start,
	// 	"unexpected item",
	// )

	pfa.errVal = v
	pfa.traceAction("<ansiGreen>UnexpectedItem")
	// }
}

// func (v UnexpectedItem) String() string {
// 	return "UnexpectedItem"
// }

type Panic struct {
	FmtString string
	FmtArgs   []interface{}
}

func (v Panic) GetAction() processorAction {
	fmtArgs := []ValProvider{}
	for _, i := range v.FmtArgs {
		switch t := i.(type) {
		case bool, rune, string, int, int8, int16 /* int32, */, int64, uint, uint8, uint16, uint32, uint64:
			fmtArgs = append(fmtArgs, justVal{t})
		case Var:
			fmtArgs = append(fmtArgs, varVal{MustVarPathFrom(t.VarPathStr)})
		default:
			bwerror.Panic("%#v", i)
		}
	}
	return panicAction{v.FmtString, fmtArgs}
}

type panicAction struct {
	fmtString string
	fmtArgs   []ValProvider
}

func (v panicAction) execute(pfa *pfaStruct) {
	var s string
	for _, i := range v.fmtArgs {
		var val interface{}
		switch t := i.(type) {
		case justVal:
			val = t.val
		case varVal:
			val = t.varPath
		default:
			bwerror.Panic("%#v", i)
		}
		s += ", " + string(pfa.traceVal(val))
	}
	pfa.traceAction("<ansiGreen>Panic{%s%s}", v.fmtString, s)
	fmtArgs := []interface{}{v.fmtString}
	for _, i := range v.fmtArgs {
		fmtArgs = append(fmtArgs, i.getVal(pfa))
	}
	pfa.panic(fmtArgs...)
}

// func (v Panic) getVal(pfa *pfaStruct) interface{} {

// 	return unexpectedItemAction{MustVarPathFrom(v.VarPathStr)}
// 	return v
// }

// func (v Panic) getSource(pfa *pfaStruct) formattedString {
// 	return pfa.traceVal(v)
// }

// func (v Panic) String() string {
// 	return "Panic"
// }

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
	if pfa.err != nil {
		return
	}
	for _, b := range v.by {
		val = b.TransformValue(pfa, val)
		if pfa.err != nil {
			return
		}
		if pfa.errVal != nil {
			break
		}
	}
	op := ""
	var isNotAppendee bool
	var expectedToBeAppendable formattedString
	if pfa.err == nil {
		if v.at == noAppend {
			op = ">"
		} else {
			op = ">>"
			if orig := pfa.getVarValue(v.varPath); pfa.err == nil {
				if v.at == appendScalar {
					if orig.val == nil {
						switch aval := val.(type) {
						case string:
							val = aval
						case rune:
							val = string(aval)
						default:
							val = []interface{}{val}
						}
					} else {
						switch oval := orig.val.(type) {
						case string:
							switch aval := val.(type) {
							case string:
								val = oval + aval
							case rune:
								val = oval + string(aval)
							default:
								expectedToBeAppendable = formattedStringFrom("<ansiSecondary>String<ansi> or <ansiSecondary>Rune")
							}
						case rune:
							switch aval := val.(type) {
							case string:
								val = string(oval) + aval
							case rune:
								val = string(oval) + string(aval)
							default:
								expectedToBeAppendable = formattedStringFrom("<ansiSecondary>String<ansi> or <ansiSecondary>Rune")
							}
						case []interface{}:
							val = append(oval, val)
						default:
							isNotAppendee = true
						}
					}
				} else if aval, ok := val.([]interface{}); !ok {
					expectedToBeAppendable = formattedStringFrom("<ansiPrimary>Array<ansi>")
				} else {
					if orig.val == nil {
						val = aval
					} else if oval, ok := orig.val.([]interface{}); !ok {
						isNotAppendee = true
					} else {
						val = append(oval, aval...)
					}
				}

				// } else {

				// }

				// switch aval := val.(type) {

				// // if orig.val == nil {

				// }
				// if orig.val == nil {
				// 	// isNotAppendee = true
				// 	switch aval := val.(type) {
				// 	case string:
				// 		val = aval
				// 	case rune:
				// 		val = string(aval)
				// 	default:
				// 		val = []interface{}{val}
				// 		// expectedToBeAppendable = formattedStringFrom("<ansiSecondary>String<ansi> or <ansiSecondary>Rune")
				// 	}
				// } else {
				// 	switch oval := orig.val.(type) {
				// 	case string:
				// 		switch aval := val.(type) {
				// 		case string:
				// 			val = oval + aval
				// 		case rune:
				// 			val = oval + string(aval)
				// 		default:
				// 			expectedToBeAppendable = formattedStringFrom("<ansiSecondary>String<ansi> or <ansiSecondary>Rune")
				// 		}
				// 	case rune:
				// 		switch aval := val.(type) {
				// 		case string:
				// 			val = string(oval) + aval
				// 		case rune:
				// 			val = string(oval) + string(aval)
				// 		default:
				// 			expectedToBeAppendable = formattedStringFrom("<ansiSecondary>String<ansi> or <ansiSecondary>Rune")
				// 		}
				// 	case []interface{}:
				// 		if v.at == appendScalar {
				// 			val = append(oval, val)
				// 		} else {
				// 			switch aval := val.(type) {
				// 			case []interface{}:
				// 				val = append(oval, aval...)
				// 				op = ">>>"
				// 			default:
				// 				expectedToBeAppendable = formattedStringFrom("<ansiPrimary>Array<ansi>")
				// 			}
				// 		}
				// 	default:
				// 		isNotAppendee = true
				// 	}
				// }

			}
		}
		if pfa.err != nil {
			return
		}
	}
	if isNotAppendee {
		pfa.panic("can not append to %s", string(pfa.traceVal(v.varPath)))
	}
	var source formattedString
	var target formattedString
	if pfa.traceLevel > TraceNone || len(expectedToBeAppendable) > 0 || pfa.errVal != nil {
		source = v.valProvider.getSource(pfa)
		target = pfa.traceVal(v.varPath)
		for _, b := range v.by {
			if pfa.traceLevel > TraceNone {
				source = source.concat(formattedStringFrom(" <ansiGreen>| %s<ansi> ", b.String()))
			}
		}
		if pfa.errVal != nil {
			switch t := pfa.errVal.(type) {
			case parseNumberFailed:
				pfa.errVal = failedToTransformToNumber{fmt.Sprintf("failed to transform %s to number (%s)", source, t.s)}
			case getValFailed:
			default:
				pfa.panic("%#v", pfa.errVal)
			}
		} else if len(expectedToBeAppendable) > 0 {
			pfa.panic("%s expected to be %s", string(source), string(expectedToBeAppendable))
		}
	}
	pfa.setVarVal(v.varPath, val)
	if pfa.traceLevel > TraceNone {
		pfa.traceAction("%s %s %s: %s", source, formattedString(ansi.Ansi("Green", op)), target, v.varPath)
	}
	// if pfa.err != nil {
	// return
	// }
}

type parseNumberFailed struct{ s string }

type failedToTransformToNumber struct{ s string }

type ValTransformer interface {
	TransformValue(pfa *pfaStruct, i interface{}) interface{}
	String() string
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
	var s string
	switch t := i.(type) {
	case string:
		s = t
	case rune:
		s = string(t)
	default:
		pfa.errVal = parseNumberFailed{"niether String, nor Rune"}
		// pfa.err = bwerror.Error("niether String, nor Rune")
		// pfa.errVal = parseNumberFailed{}
	}
	if pfa.err != nil {
		return
	}

	result, err := _parseNumber(s)
	if err != nil {
		pfa.errVal = parseNumberFailed{err.Error()}
	}
	return
}

func (v ParseNumber) String() string {
	return "ParseNumber"
}

type Append struct{}

func (v Append) TransformValue(pfa *pfaStruct, i interface{}) (result interface{}) {
	bwerror.Panic("unreachable")
	return
}

func (v Append) String() string {
	return "Append"
}

type AppendSlice struct{}

func (v AppendSlice) TransformValue(pfa *pfaStruct, i interface{}) (result interface{}) {
	bwerror.Panic("unreachable")
	// bwerror.Unreachable()
	return
}

func (v AppendSlice) String() string {
	return "AppendSlice"
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
