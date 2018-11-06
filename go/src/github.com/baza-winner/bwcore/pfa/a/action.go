package a

import (
	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"

	// "github.com/baza-winner/bwcore/pfa"
	"github.com/baza-winner/bwcore/pfa/b"
	"github.com/baza-winner/bwcore/pfa/core"
	"github.com/baza-winner/bwcore/pfa/d"
	"github.com/baza-winner/bwcore/pfa/formatted"
)

// ============================================================================

type PushRune struct{}

func (v PushRune) Execute(pfa *core.PfaStruct) (err error) {
	pfa.Proxy.PushRune()
	if pfa.TraceLevel > core.TraceNone {
		pfa.TraceAction("<ansiFunc>PushRune<ansi>: %s", runeVarPath)
	}
	return
}

func (v PushRune) GetAction() core.ProcessorAction {
	return v
}

// ============================================================================

type PullRune struct{}

var runeVarPath core.VarPath = core.MustVarPathFrom("rune")
var stackLenVarPath core.VarPath = core.MustVarPathFrom("stackLen")

func (v PullRune) Execute(pfa *core.PfaStruct) (err error) {
	pfa.Proxy.PullRune()
	if pfa.TraceLevel > core.TraceNone {
		pfa.TraceAction("<ansiFunc>PullRune<ansi>: %s", runeVarPath)
	}
	return
}

func (v PullRune) GetAction() core.ProcessorAction {
	return v
}

// ============================================================================

// type PushItem struct{}

// func (v PushItem) Execute(pfa *core.PfaStruct) {
// 	pfa.PushStackItem()
// 	if pfa.TraceLevel > core.TraceNone {
// 		pfa.TraceAction("<ansiFunc>PushItem<ansi>: %s", stackLenVarPath)
// 	}
// }

// func (v PushItem) GetAction() core.ProcessorAction {
// 	return v
// }

// // ============================================================================

// type PopItem struct{}

// func (v PopItem) Execute(pfa *core.PfaStruct) {
// 	pfa.PopStackItem()
// 	if pfa.TraceLevel > core.TraceNone {
// 		pfa.TraceAction("<ansiFunc>PopItem<ansi>: %s", stackLenVarPath)
// 	}
// 	// pfa.traceAction("<ansiFunc>PopItem<ansi>: <ansiPath>stackLen<ansi>(<ansiVal>%d<ansi>)", len(pfa.Stack))
// }

// func (v PopItem) GetAction() core.ProcessorAction {
// 	return v
// }

// ============================================================================

// type Panic struct {
// 	FmtString string
// 	FmtArgs   []interface{}
// }

// func (v Panic) GetAction() core.ProcessorAction {
// 	fmtArgs := []core.ValProvider{}
// 	for _, i := range v.FmtArgs {
// 		switch t := i.(type) {
// 		case bool, rune, string, int, int8, int16 /* int32, */, int64, uint, uint8, uint16, uint32, uint64:
// 			fmtArgs = append(fmtArgs, common.JustVal{t})
// 		case val.Var:
// 			fmtArgs = append(fmtArgs, common.VarVal{core.MustVarPathFrom(t.VarPathStr)})
// 		default:
// 			bwerr.Panic("%#v", i)
// 		}
// 	}
// 	return panicAction{v.FmtString, fmtArgs}
// }

// type panicAction struct {
// 	fmtString string
// 	fmtArgs   []core.ValProvider
// }

// func (v panicAction) Execute(pfa *core.PfaStruct) {
// 	var s string
// 	for _, i := range v.fmtArgs {
// 		s += ", " + string(i.GetSource(pfa))
// 	}
// 	if pfa.TraceLevel > core.TraceNone {
// 		pfa.TraceAction("<ansiFunc>Panic{%s%s}", v.fmtString, s)
// 	}
// 	fmtArgs := []interface{}{}
// 	for _, i := range v.fmtArgs {
// 		fmtArgs = append(fmtArgs, i.GetVal(pfa))
// 	}
// 	pfa.Panic(bw.StructFrom(v.fmtString, fmtArgs))
// }

// ============================================================================

type SetVar struct {
	VarPathStr string
	Val        interface{}
}

func (v SetVar) GetAction() core.ProcessorAction {
	return _setVarBy{
		core.MustVarPathFrom(v.VarPathStr),
		d.ValFrom(v.Val),
		// val.MustValProviderFrom(v.VarValue),
		b.By{},
		noAppend,
	}
}

// ============================================================================

type SetVarBy struct {
	VarPathStr   string
	Val          interface{}
	Transformers b.By
}

func (v SetVarBy) GetAction() core.ProcessorAction {
	by := b.By{}
	at := noAppend
	for _, i := range v.Transformers {
		switch i.(type) {
		case b.Append:
			if at == noAppend {
				at = appendScalar
			}
		case b.AppendSlice:
			at = appendSlice
		default:
			by = append(by, i)
		}
	}
	return _setVarBy{
		core.MustVarPathFrom(v.VarPathStr),
		d.ValFrom(v.Val),
		// val.MustValProviderFrom(v.VarValue),
		by,
		at,
	}
}

// ============================================================================

type appendType uint8

const (
	noAppend appendType = iota
	appendScalar
	appendSlice
)

type _setVarBy struct {
	varPath  core.VarPath
	valToSet d.Val
	by       b.By
	at       appendType
}

func (v _setVarBy) Execute(pfa *core.PfaStruct) (err error) {
	var val interface{}
	val, err = v.valToSet.GetVal(pfa)
	if err != nil {
		return
	}
	for _, b := range v.by {
		val, err = b.TransformValue(pfa, val)
		if err != nil {
			return
		}
	}
	op := ""
	var isNotAppendee bool
	var expectedToBeAppendable formatted.String
	if pfa.Err == nil {
		if v.at == noAppend {
			op = ">"
		} else {
			op = ">>"
			var orig core.VarValue
			if orig, err = pfa.VarValue(v.varPath); err != nil {
				return
			}
			if v.at == appendScalar {
				// pfa.Panic(bw.StructFrom("%#v", val))
				if orig.Val == nil {
					switch aval := val.(type) {
					case string:
						val = aval
					case rune:
						val = string(aval)
					default:
						val = []interface{}{val}
					}
				} else {
					switch oval := orig.Val.(type) {
					case string:
						switch aval := val.(type) {
						case string:
							val = oval + aval
						case rune:
							val = oval + string(aval)
						default:
							expectedToBeAppendable = formatted.StringFrom("<ansiVal>String<ansi> or <ansiVal>Rune")
						}
					case rune:
						switch aval := val.(type) {
						case string:
							val = string(oval) + aval
						case rune:
							val = string(oval) + string(aval)
						default:
							expectedToBeAppendable = formatted.StringFrom("<ansiVal>String<ansi> or <ansiVal>Rune")
						}
					case []interface{}:
						val = append(oval, val)
					default:
						isNotAppendee = true
					}
				}
			} else if aval, ok := val.([]interface{}); !ok {
				expectedToBeAppendable = formatted.StringFrom("<ansiVal>Array<ansi>")
			} else {
				if orig.Val == nil {
					val = aval
				} else if oval, ok := orig.Val.([]interface{}); !ok {
					isNotAppendee = true
				} else {
					val = append(oval, aval...)
				}
			}
		}
		if pfa.Err != nil {
			return
		}
	}
	if isNotAppendee {
		pfa.PanicA(bw.Fmt("can not append to %s", string(pfa.TraceVal(v.varPath))))
	}
	var source formatted.String
	var target formatted.String
	if pfa.TraceLevel > core.TraceNone || len(expectedToBeAppendable) > 0 || pfa.Err != nil {
		// source = v.valProvider.GetSource(pfa)
		source = v.valToSet.GetSource(pfa)
		target = pfa.TraceVal(v.varPath)
		for _, b := range v.by {
			if pfa.TraceLevel > core.TraceNone {
				source = source.Concat(formatted.StringFrom(" <ansiFunc>| %s<ansi> ", b.String()))
			}
		}
		if pfa.Err != nil {
			if t, ok := pfa.Err.(core.PfaError); !ok {
				pfa.PanicA(bw.Fmt("%#v", pfa.Err))
			} else if t.State() != core.PecsPrepared {
				t.PrepareErr(string(source))
			}
		} else if len(expectedToBeAppendable) > 0 {
			pfa.PanicA(bw.Fmt("%s expected to be %s", string(source), string(expectedToBeAppendable)))
		}
	}
	pfa.SetVarVal(v.varPath, val)
	if pfa.TraceLevel > core.TraceNone {
		pfa.TraceAction("%s %s %s: %s", source,
			formatted.String(ansi.StringA(ansi.A{
				Default: []ansi.SGRCode{ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen})},
				S:       op,
			})), target, v.varPath)
		// if op == ">>" {
		// 	pfa.Panic()
		// }
	}
	return
}