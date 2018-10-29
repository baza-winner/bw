package core

import (
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/pfa/formatted"
	"github.com/jimlawless/whereami"
)

// ============================================================================

type PfaError struct {
	pfa     *PfaStruct
	content *PfaErrorContent
	Where   string
}

type PfaErrorContentState uint8

const (
	PecsNeedPrepare PfaErrorContentState = iota
	PecsPrepared
)

type PfaErrorContent struct {
	state     PfaErrorContentState
	reason    string
	fmtString string
	errStr    string
}

func (v PfaError) DataForJSON() interface{} {
	result := map[string]interface{}{}
	switch v.content.state {
	case PecsNeedPrepare:
		result["reason"] = v.content.reason
		if len(v.content.fmtString) > 0 {
			result["fmtString"] = v.content.fmtString
		}
	case PecsPrepared:
		result["pfa"] = v.pfa.DataForJSON()
		result["err"] = v.content.errStr
		result["Where"] = v.Where
	}
	return result
}

func (pfa *PfaStruct) SetError(fmtString string, fmtArgs ...interface{}) {
	pfa.Err = PfaError{
		pfa: pfa,
		content: &PfaErrorContent{
			reason: string(formatted.StringFrom(fmtString, fmtArgs)),
		},
		Where: whereami.WhereAmI(3),
	}
}

func (pfa *PfaStruct) SetTransformError(fmtString, reason string) {
	pfa.Err = PfaError{
		pfa:     pfa,
		content: &PfaErrorContent{fmtString: fmtString, reason: reason},
		Where:   whereami.WhereAmI(3),
	}
}

func (pfa *PfaStruct) SetUnexpectedError(err error) {
	pfa.Err = PfaError{
		pfa:     pfa,
		content: &PfaErrorContent{errStr: err.Error(), state: PecsPrepared},
		Where:   whereami.WhereAmI(3),
	}
}

func (v *PfaError) PrepareErr(fmtString string, fmtArgs ...interface{}) {
	if v.content.state == PecsPrepared {
		bwerror.Panic("Already prepared %s ", bwjson.PrettyJsonOf(v))
	} else {
		v.content.fmtString = fmtString + ": " + v.content.reason
		v.content.errStr = bwerror.Error(v.content.fmtString, fmtArgs...).Error()
		v.content.state = PecsPrepared
	}
}

func (v *PfaError) SetErr(errStr string) {
	if v.content.state == PecsPrepared {
		bwerror.Panic("Already prepared %s ", bwjson.PrettyJsonOf(v))
	} else {
		v.content.errStr = errStr
		v.content.state = PecsPrepared
	}
}

func (v PfaError) Error() (result string) {
	switch v.content.state {
	case PecsNeedPrepare:
		bwerror.Panic("NeedPrepare %s ", bwjson.PrettyJsonOf(v))
	case PecsPrepared:
		result = v.content.errStr
	}
	return
}

func (v PfaError) State() PfaErrorContentState {
	return v.content.state
}

// ============================================================================
