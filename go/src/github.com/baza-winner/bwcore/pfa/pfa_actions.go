package pfa

import (
	"fmt"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/pfa/core"
	"github.com/baza-winner/bwcore/pfa/formatted"
)

// ============================== hasRune =====================================

// type EOF struct{}

// func (v EOF) String() string {
// 	return "EOF"
// }

// ============================================================================

type VarIs struct {
	VarPathStr string
	VarValue   interface{}
}

// ======================== core.ProcessorAction ====================================

type ProccessorActionProvider interface {
	GetAction() core.ProcessorAction
}

type PushRune struct{}

func (v PushRune) Execute(pfa *core.PfaStruct) {
	pfa.p.PushRune()
	pfa.traceAction("<ansiGreen>PushRune<ansi>: %s", runeVarPath)
}

func (v PushRune) GetAction() core.ProcessorAction {
	return v
}

type PullRune struct{}

var runeVarPath core.VarPath = MustVarPathFrom("rune")
var stackLenVarPath core.VarPath = MustVarPathFrom("stackLen")

func (v PullRune) Execute(pfa *core.PfaStruct) {
	pfa.p.PullRune()
	pfa.traceAction("<ansiGreen>PullRune<ansi>: %s", runeVarPath)
}

func (v PullRune) GetAction() core.ProcessorAction {
	return v
}

type PushItem struct{}

func (v PushItem) Execute(pfa *core.PfaStruct) {
	pfa.pushStackItem()
	// pfa.traceAction("<ansiGreen>PushItem<ansi>: <ansiCmd>stackLen<ansi>(<ansiPrimary>%d<ansi>)", len(pfa.Stack))
	pfa.traceAction("<ansiGreen>PushItem<ansi>: %s", stackLenVarPath)
}

func (v PushItem) GetAction() core.ProcessorAction {
	return v
}

type SubRules struct {
	Def Rules
}

func (v SubRules) GetAction() core.ProcessorAction {
	return v
}

func (v SubRules) Execute(pfa *core.PfaStruct) {
	pfa.traceIncLevel()
	pfa.processRules(v.Def)
	pfa.traceDecLevel()
}

type PopItem struct{}

func (v PopItem) Execute(pfa *core.PfaStruct) {
	pfa.popStackItem()
	pfa.traceAction("<ansiGreen>PopItem<ansi>: %s", stackLenVarPath)
	// pfa.traceAction("<ansiGreen>PopItem<ansi>: <ansiCmd>stackLen<ansi>(<ansiPrimary>%d<ansi>)", len(pfa.Stack))
}

func (v PopItem) GetAction() core.ProcessorAction {
	return v
}

type Var struct {
	VarPathStr string
}

// func (v UnexpectedItem) String() string {
// 	return "UnexpectedItem"
// }

type Panic struct {
	FmtString string
	FmtArgs   []interface{}
}

func (v Panic) GetAction() core.ProcessorAction {
	fmtArgs := []core.ValProvider{}
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
	fmtArgs   []core.ValProvider
}

func (v panicAction) Execute(pfa *core.PfaStruct) {
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
		fmtArgs = append(fmtArgs, i.GetVal(pfa))
	}
	pfa.panic(fmtArgs...)
}

// func (v Panic) GetVal(pfa *core.PfaStruct) interface{} {

// 	return unexpectedItemAction{MustVarPathFrom(v.VarPathStr)}
// 	return v
// }

// func (v Panic) GetSource(pfa *core.PfaStruct) formatted.String {
// 	return pfa.traceVal(v)
// }

// func (v Panic) String() string {
// 	return "Panic"
// }

type justVal struct {
	val interface{}
}

func (v justVal) conforms(pfa *core.PfaStruct, val interface{}, varPath core.VarPath) (result bool) {
	result = val == v.val
	pfa.traceCondition(varPath, val, result)
	return
}

func (v justVal) GetVal(pfa *core.PfaStruct) interface{} {
	return v.val
}

func (v justVal) GetSource(pfa *core.PfaStruct) formatted.String {
	return pfa.traceVal(v.val)
}

type varVal struct {
	varPath core.VarPath
}

func (v varVal) conforms(pfa *core.PfaStruct, val interface{}, varPath core.VarPath) (result bool) {
	varValue := pfa.getVarValue(v.varPath)
	if pfa.err == nil {
		result = varValue.val == val
		pfa.traceCondition(varPath, v.varPath, result)
	}
	return
}

func (v varVal) GetVal(pfa *core.PfaStruct) (result interface{}) {
	varValue := pfa.getVarValue(v.varPath)
	if pfa.err == nil {
		result = varValue.val
	}
	return
}

func (v varVal) GetSource(pfa *core.PfaStruct) formatted.String {
	// return fmt.Sprintf(ansi.Ansi("", "%s"), pfa.traceVarPath(v.varPath))
	return pfa.traceVal(v.varPath)
}

type SetVar struct {
	VarPathStr string
	VarValue   interface{}
}

func (v SetVar) GetAction() core.ProcessorAction {
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

func (v SetVarBy) GetAction() core.ProcessorAction {
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
	varPath     core.VarPath
	valProvider core.ValProvider
	by          By
	at          appendType
}

func (v _setVarBy) Execute(pfa *core.PfaStruct) {
	val := v.valProvider.GetVal(pfa)
	if pfa.err != nil {
		return
	}
	for _, b := range v.by {
		val = b.TransformValue(pfa, val)
		if pfa.err != nil {
			return
		}
		if pfa.ErrVal != nil {
			break
		}
	}
	op := ""
	var isNotAppendee bool
	var expectedToBeAppendable formatted.String
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
								expectedToBeAppendable = formatted.StringFrom("<ansiSecondary>String<ansi> or <ansiSecondary>Rune")
							}
						case rune:
							switch aval := val.(type) {
							case string:
								val = string(oval) + aval
							case rune:
								val = string(oval) + string(aval)
							default:
								expectedToBeAppendable = formatted.StringFrom("<ansiSecondary>String<ansi> or <ansiSecondary>Rune")
							}
						case []interface{}:
							val = append(oval, val)
						default:
							isNotAppendee = true
						}
					}
				} else if aval, ok := val.([]interface{}); !ok {
					expectedToBeAppendable = formatted.StringFrom("<ansiPrimary>Array<ansi>")
				} else {
					if orig.val == nil {
						val = aval
					} else if oval, ok := orig.val.([]interface{}); !ok {
						isNotAppendee = true
					} else {
						val = append(oval, aval...)
					}
				}
			}
		}
		if pfa.err != nil {
			return
		}
	}
	if isNotAppendee {
		pfa.panic("can not append to %s", string(pfa.traceVal(v.varPath)))
	}
	var source formatted.String
	var target formatted.String
	if pfa.TraceLevel > TraceNone || len(expectedToBeAppendable) > 0 || pfa.ErrVal != nil {
		source = v.valProvider.GetSource(pfa)
		target = pfa.traceVal(v.varPath)
		for _, b := range v.by {
			if pfa.TraceLevel > TraceNone {
				source = source.concat(formatted.StringFrom(" <ansiGreen>| %s<ansi> ", b.String()))
			}
		}
		if pfa.ErrVal != nil {
			switch t := pfa.ErrVal.(type) {
			case parseNumberFailed:
				pfa.ErrVal = failedToTransformToNumber{fmt.Sprintf("failed to transform %s to number (%s)", source, t.s)}
			case getValFailed:
			default:
				pfa.panic("%#v", pfa.ErrVal)
			}
		} else if len(expectedToBeAppendable) > 0 {
			pfa.panic("%s expected to be %s", string(source), string(expectedToBeAppendable))
		}
	}
	pfa.setVarVal(v.varPath, val)
	if pfa.TraceLevel > TraceNone {
		pfa.traceAction("%s %s %s: %s", source, formatted.String(ansi.Ansi("Green", op)), target, v.varPath)
	}
}

type parseNumberFailed struct{ s string }

type failedToTransformToNumber struct{ s string }

type ValTransformer interface {
	TransformValue(pfa *core.PfaStruct, i interface{}) interface{}
	String() string
}

type By []ValTransformer

type ValChecker interface {
	conforms(pfa *core.PfaStruct, val interface{}, varPath core.VarPath) bool
}

func valProviderFrom(i interface{}) (result core.ValProvider, err error) {
	switch t := i.(type) {
	case Var:
		var varPath core.VarPath
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

func MustValProviderFrom(i interface{}) (result core.ValProvider) {
	var err error
	if result, err = valProviderFrom(i); err != nil {
		bwerror.PanicErr(err)
	}
	return
}

type Debug struct{ Message string }

func (v *Debug) Execute(pfa *core.PfaStruct) {
	fmt.Printf("%s: %s\n", v.Message, bwjson.PrettyJsonOf(pfa))
}

// ============================================================================

type ParseNumber struct{}

func (v ParseNumber) TransformValue(pfa *core.PfaStruct, i interface{}) (result interface{}) {
	var s string
	switch t := i.(type) {
	case string:
		s = t
	case rune:
		s = string(t)
	default:
		pfa.ErrVal = parseNumberFailed{"niether String, nor Rune"}
		// pfa.err = bwerror.Error("niether String, nor Rune")
		// pfa.ErrVal = parseNumberFailed{}
	}
	if pfa.err != nil {
		return
	}

	result, err := _parseNumber(s)
	if err != nil {
		pfa.ErrVal = parseNumberFailed{err.Error()}
	}
	return
}

func (v ParseNumber) String() string {
	return "ParseNumber"
}

type Append struct{}

func (v Append) TransformValue(pfa *core.PfaStruct, i interface{}) (result interface{}) {
	bwerror.Panic("unreachable")
	return
}

func (v Append) String() string {
	return "Append"
}

type AppendSlice struct{}

func (v AppendSlice) TransformValue(pfa *core.PfaStruct, i interface{}) (result interface{}) {
	bwerror.Panic("unreachable")
	// bwerror.Unreachable()
	return
}

func (v AppendSlice) String() string {
	return "AppendSlice"
}

// ============================================================================

// ============================================================================
