package e

import (
	"github.com/baza-winner/bwcore/pfa/core"
)

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
