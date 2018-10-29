package pfa

import (
	"github.com/baza-winner/bwcore/bwfmt"
	"github.com/baza-winner/bwcore/pfa/core"
	"github.com/baza-winner/bwcore/pfa/r"
	"github.com/baza-winner/bwcore/runeprovider"
)

func init() {
}

// ============================================================================

// ============================================================================

type SubRules struct {
	R r.Rules
}

func (v SubRules) GetAction() core.ProcessorAction {
	return v
}

func (v SubRules) Execute(pfa *core.PfaStruct) {
	if pfa.TraceLevel > core.TraceNone {
		pfa.TraceIncLevel()
	}
	v.R.Process(pfa)
	if pfa.TraceLevel > core.TraceNone {
		pfa.TraceDecLevel()
	}
}

// ============================================================================

func Run(p runeprovider.RuneProvider, rules r.Rules, optTraceLevel ...core.TraceLevel) (result interface{}, err error) {
	traceLevel := core.TraceNone
	if optTraceLevel != nil {
		traceLevel = optTraceLevel[0]
	}
	pfa := core.PfaFrom(p, traceLevel)
	for {
		rules.Process(pfa)
		if pfa.Err != nil || pfa.Proxy.Curr.IsEOF {
			break
		}
	}
	if pfa.Err != nil {
		err = pfa.Err
	} else {
		if len(pfa.Stack) > 1 {
			pfa.Panic(bwfmt.StructFrom("len(pfa.Stack) > 1"))
		} else if len(pfa.Stack) > 0 {
			result = pfa.GetTopStackItem().Vars["result"]
		}
	}
	return
}

// func runePtr(r rune) *rune {
// 	return &r
// }

// ============================================================================
