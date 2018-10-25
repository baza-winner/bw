package b

import (
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/pfa/core"
)

type ValTransformer interface {
	TransformValue(pfa *core.PfaStruct, i interface{}) interface{}
	String() string
}

type By []ValTransformer

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
		pfa.ErrVal = ParseNumberFailed{"niether String, nor Rune"}
	}
	if pfa.Err != nil {
		return
	}

	result, err := core.ParseNumber(s)
	if err != nil {
		pfa.ErrVal = ParseNumberFailed{err.Error()}
	}
	return
}

func (v ParseNumber) String() string {
	return "ParseNumber"
}

type ParseNumberFailed struct{ S string }

func (v ParseNumberFailed) Error(pfa *core.PfaStruct) error {
	bwerror.Unreachable()
	return nil
}

// ============================================================================

type Append struct{}

func (v Append) TransformValue(pfa *core.PfaStruct, i interface{}) (result interface{}) {
	bwerror.Unreachable()
	return
}

func (v Append) String() string {
	return "Append"
}

// ============================================================================

type AppendSlice struct{}

func (v AppendSlice) TransformValue(pfa *core.PfaStruct, i interface{}) (result interface{}) {
	bwerror.Unreachable()
	// bwerror.Unreachable()
	return
}

func (v AppendSlice) String() string {
	return "AppendSlice"
}

// ============================================================================
