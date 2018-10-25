package e

import "github.com/baza-winner/bwcore/pfa/core"

// ============================================================================

type UnexpectedRune struct{}

func (v UnexpectedRune) GetAction() core.ProcessorAction {
	return v
}

func (v UnexpectedRune) Execute(pfa *core.PfaStruct) {
	pfa.ErrVal = v
	if pfa.TraceLevel > core.TraceNone {
		pfa.TraceAction("<ansiGreen>UnexpectedRune")
	}
}

func (v UnexpectedRune) Error(pfa *core.PfaStruct) error {
	return pfa.Proxy.UnexpectedRuneError()
}

// ============================================================================

type UnexpectedItem struct{}

func (v UnexpectedItem) GetAction() core.ProcessorAction {
	return v
}

func (v UnexpectedItem) Execute(pfa *core.PfaStruct) {
	pfa.ErrVal = v
	if pfa.TraceLevel > core.TraceNone {
		pfa.TraceAction("<ansiGreen>UnexpectedItem")
	}
}

func (v UnexpectedItem) Error(pfa *core.PfaStruct) error {
	stackItem := pfa.GetTopStackItem()
	start := stackItem.Start
	item := pfa.Proxy.Curr.Prefix[start.Pos-pfa.Proxy.Curr.PrefixStart:]
	return pfa.Proxy.ItemError(start, "unexpected \"<ansiPrimary>%s<ansi>\"", item)
}

// ============================================================================
