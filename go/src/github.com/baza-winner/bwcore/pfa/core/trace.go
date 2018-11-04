package core

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/pfa/formatted"
)

func (pfa *PfaStruct) TraceVal(val interface{}) (result formatted.String) {
	switch t := val.(type) {
	case formatted.String:
		result = t
	case rune, string:
		result = formatted.StringFrom("<ansiVal>%q", val)
	case formatted.FormattedString:
		result = t.FormattedString()
	case VarPath:
		var valStr formatted.String
		if t[0].Type == VarPathItemKey && (t[0].Key == "rune" || t[0].Key == "runePos") {
			ofs := 0
			if len(t) > 1 {
				vt, idx, _, err := t[1].TypeIdxKey(pfa)
				if err == nil && vt == VarPathItemIdx {
					ofs = idx
				}
			}
			switch t[0].Key {
			case "rune":
				if r, isEOF := pfa.Proxy.Rune(ofs); isEOF {
					valStr = formatted.StringFrom("<ansiVar>EOF")
				} else {
					valStr = pfa.TraceVal(r)
				}
			case "runePos":
				ps := pfa.Proxy.PosStruct(ofs)
				valStr = pfa.TraceVal(ps.Pos)
			default:
				bwerr.Unreachable()
			}
		} else {
			// {
			// 	p := MustVarPathFrom("stack")
			// 	fmt.Printf("%s: %#v\n", p.FormattedString(), pfa.VarValue(p).Val)
			// }
			// {
			// 	p := MustVarPathFrom("stack.-1")
			// 	fmt.Printf("%s: %#v\n", p.FormattedString(), pfa.VarValue(p).Val)
			// }
			varValue, err := pfa.VarValue(t)
			if err != nil {
				bwerr.PanicA(bwerr.E{Error: err})
			}
			// val :=
			fmt.Printf("%s: %#v\n", t.FormattedString(), varValue.Val)
			// fmt.Printf("||| %#v\n, t: %#v", pfa.VarValue(t).Val, t)
			valStr = pfa.TraceVal(varValue.Val)
		}
		result = formatted.String(fmt.Sprintf("%s(%s)", t.FormattedString(pfa), valStr))
	case bwset.String, bwset.Rune, bwset.Int:
		value := reflect.ValueOf(t)
		keys := value.MapKeys()
		if len(keys) == 1 {
			result = formatted.StringFrom("<ansiVal>%s", traceValHelper(keys[0].Interface()))
		} else if len(keys) > 1 {
			ss := []string{}
			for _, k := range keys {
				ss = append(ss, traceValHelper(k.Interface()))
			}
			result = formatted.StringFrom("<<ansiVal>%s<ansi>>", strings.Join(ss, " "))
		}
	default:
		result = formatted.StringFrom("<ansiVal>%#v", val)
	}
	return
}

func (pfa *PfaStruct) TraceAction(fmtString string, fmtArgs ...interface{}) {
	fmt.Printf(pfa.indent(pfa.ruleLevel+1)+ansi.String(fmtString+"\n"), pfa.fmtArgs(fmtArgs...)...)
}

func (pfa *PfaStruct) TraceCondition(varPath VarPath, arg interface{}, result bool) {
	// if pfa.TraceLevel > core.TraceNone {
	fmtArgs := pfa.fmtArgs(varPath, arg)
	pfa.traceConditions = append(pfa.traceConditions, fmt.Sprintf("%s %s", fmtArgs...))
	// }
}

func (pfa *PfaStruct) TraceBeginConditions() {
	pfa.traceConditions = nil
}

func (pfa *PfaStruct) TraceFailedConditions() {
	pfa.traceConditionsHelper(" <ansiErr>Failed")
}

func (pfa *PfaStruct) TraceBeginActions() {
	pfa.traceConditionsHelper("")
}

func (pfa *PfaStruct) traceConditionsHelper(suffix string) {
	fmt.Printf(
		pfa.indent(pfa.ruleLevel) +
			strings.Join(pfa.traceConditions, " ") +
			ansi.StringA(ansi.A{
				Default: []ansi.SGRCode{
					ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorYellow}),
					ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
				},
				S: ":" + suffix + "\n",
			}),
	)
	pfa.traceConditions = nil
}

func (pfa *PfaStruct) TraceIncLevel() {
	// if pfa.TraceLevel > core.TraceNone {
	pfa.ruleLevel++
	// }
}

func (pfa *PfaStruct) TraceDecLevel() {
	// if pfa.TraceLevel > core.TraceNone {
	pfa.ruleLevel--
	// }
}

func traceValHelper(i interface{}) (s string) {
	switch t := i.(type) {
	case rune, string:
		s = fmt.Sprintf("%q", t)
	default:
		s = fmt.Sprintf("%d", t)
	}
	return
}

func (pfa *PfaStruct) fmtArgs(fmtArgs ...interface{}) []interface{} {
	result := []interface{}{}
	for _, arg := range fmtArgs {
		if f, ok := arg.(func(pfa *PfaStruct) interface{}); ok {
			arg = f(pfa)
		}
		result = append(result, pfa.TraceVal(arg))
	}
	return result
}
