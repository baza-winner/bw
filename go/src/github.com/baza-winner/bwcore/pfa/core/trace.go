package core

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/pfa/formatted"
)

func (pfa *PfaStruct) TraceVal(val interface{}) (result formatted.String) {
	// if pfa.TraceLevel > TraceNone {
	switch t := val.(type) {
	case formatted.String:
		result = t
	case rune, string:
		result = formatted.StringFrom("<ansiPrimary>%q", val)
	// result = formatted.String(fmt.Sprintf(ansi.Ansi("", "<ansiPrimary>%q"), val))
	case formatted.FormattedString:
		result = t.FormattedString()
	// case Map, Array:
	// UnicodeCategory, EOF,
	// , UnexpectedRune, UnexpectedItem, Panic
	// result = formatted.StringFrom("<ansiOutline>%s", t)
	// result = formatted.String(fmt.Sprintf(ansi.Ansi("", "<ansiOutline>%s"), t))
	case VarPath:
		// var val interface{}
		var valStr formatted.String
		if t[0].Val == "rune" {
			ofs := 0
			if len(t) > 1 {
				isIdx, idx, _, err := t[1].GetIdxKey(pfa)
				if err == nil && isIdx {
					ofs = idx
				}
			}
			if r, isEOF := pfa.Proxy.Rune(ofs); isEOF {
				// val = EOF{}
				valStr = formatted.StringFrom("<ansiOutline>EOF")
			} else {
				valStr = pfa.TraceVal(r)
				// val = r
			}
		} else {
			// val = pfa.VarValue(t).val
			valStr = pfa.TraceVal(pfa.VarValue(t).Val)
		}
		result = formatted.String(fmt.Sprintf("%s(%s)", t.formattedString(pfa), valStr))
	case bwset.String, bwset.Rune, bwset.Int:
		value := reflect.ValueOf(t)
		keys := value.MapKeys()
		if len(keys) == 1 {
			result = formatted.StringFrom("<ansiPrimary>%s", traceValHelper(keys[0].Interface()))
		} else if len(keys) > 1 {
			ss := []string{}
			for _, k := range keys {
				ss = append(ss, traceValHelper(k.Interface()))
			}
			result = formatted.StringFrom("<<ansiSecondary>%s>", strings.Join(ss, " "))
		}
	default:
		result = formatted.StringFrom("<ansiPrimary>%#v", val)
	}
	// }
	return
}

func (pfa *PfaStruct) TraceAction(fmtString string, fmtArgs ...interface{}) {
	fmt.Printf(pfa.indent(pfa.ruleLevel+1)+ansi.Ansi("", fmtString+"\n"), pfa.fmtArgs(fmtArgs...)...)
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
	fmt.Printf(pfa.indent(pfa.ruleLevel) + strings.Join(pfa.traceConditions, " ") + ansi.Ansi("", "<ansiYellow><ansiBold>:"+suffix+"\n"))
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
