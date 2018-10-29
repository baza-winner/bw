package e

import (
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwfmt"
	"github.com/baza-winner/bwcore/pfa/core"
	"github.com/baza-winner/bwcore/runeprovider"
)

// ============================================================================

type Unexpected struct{ VarPathStr string }

func (v Unexpected) GetAction() core.ProcessorAction {
	return UnexpectedAction{core.MustVarPathFrom(v.VarPathStr)}
}

type UnexpectedAction struct{ varPath core.VarPath }

func (v UnexpectedAction) Execute(pfa *core.PfaStruct) {
	var err error
	if len(v.varPath) == 0 {
		err = bwerror.Error("unexpected")
		// return pfa.Error()
	} else {
		varValue := pfa.VarValue(v.varPath)
		if pfa.Err != nil {
			pfa.PanicErr(pfa.Err)
		} else if ps, ok := varValue.Val.(runeprovider.PosStruct); !ok {
			pfa.Panic(bwfmt.StructFrom("%#v", varValue.Val))
		} else if ps.Pos < pfa.Proxy.Curr.Pos {
			item := pfa.Proxy.Curr.Prefix[ps.Pos-pfa.Proxy.Curr.PrefixStart:]
			err = pfa.Proxy.Unexpected(ps, bwfmt.StructFrom("unexpected \"<ansiPrimary>%s<ansi>\"", item))
		} else {
			err = pfa.Proxy.Unexpected(ps)
		}
	}
	pfa.SetUnexpectedError(err)
	if pfa.TraceLevel > core.TraceNone {
		pfa.TraceAction("<ansiGreen>Unexpected %s", v.varPath)
	}
}

// func (v UnexpectedAction) Error(pfa *core.PfaStruct) (err error) {
// 	if len(v.varPath) == 0 {
// 		err = bwerror.Error("unexpected")
// 		// return pfa.Error()
// 	} else {
// 		varValue := pfa.VarValue(v.varPath)
// 		if pfa.Err != nil {
// 			pfa.PanicErr(pfa.Err)
// 		} else if ps, ok := varValue.Val.(runeprovider.PosStruct); !ok {
// 			pfa.Panic(bwfmt.StructFrom("%#v", varValue.Val))
// 		} else if ps.Pos < pfa.Proxy.Curr.Pos {
// 			item := pfa.Proxy.Curr.Prefix[ps.Pos-pfa.Proxy.Curr.PrefixStart:]
// 			err = pfa.Proxy.Unexpected(ps, bwfmt.StructFrom("unexpected \"<ansiPrimary>%s<ansi>\"", item))
// 		} else {
// 			err = pfa.Proxy.Unexpected(ps)
// 		}
// 	}
// 	return
// }

// // ============================================================================

// type UnexpectedRune struct{}

// func (v UnexpectedRune) GetAction() core.ProcessorAction {
// 	return v
// }

// func (v UnexpectedRune) Execute(pfa *core.PfaStruct) {
// 	pfa.Err = v
// 	if pfa.TraceLevel > core.TraceNone {
// 		pfa.TraceAction("<ansiGreen>UnexpectedRune")
// 	}
// }

// func (v UnexpectedRune) Error(pfa *core.PfaStruct) error {
// 	return pfa.Proxy.UnexpectedRuneError()
// }

// // ============================================================================

// type UnexpectedItem struct{}

// func (v UnexpectedItem) GetAction() core.ProcessorAction {
// 	return v
// }

// func (v UnexpectedItem) Execute(pfa *core.PfaStruct) {
// 	pfa.Err = v
// 	if pfa.TraceLevel > core.TraceNone {
// 		pfa.TraceAction("<ansiGreen>UnexpectedItem")
// 	}
// }

// func (v UnexpectedItem) Error(pfa *core.PfaStruct) error {
// 	stackItem := pfa.GetTopStackItem()
// 	start := stackItem.Start
// 	item := pfa.Proxy.Curr.Prefix[start.Pos-pfa.Proxy.Curr.PrefixStart:]
// 	return pfa.Proxy.ItemError(start, "unexpected \"<ansiPrimary>%s<ansi>\"", item)
// }

// // ============================================================================
