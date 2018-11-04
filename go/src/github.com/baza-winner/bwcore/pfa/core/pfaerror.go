package core

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/bwerr"
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

func (v PfaError) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{}
	switch v.content.state {
	case PecsNeedPrepare:
		result["reason"] = v.content.reason
		if len(v.content.fmtString) > 0 {
			result["fmtString"] = v.content.fmtString
		}
	case PecsPrepared:
		result["pfa"] = v.pfa
		result["err"] = v.content.errStr
		result["Where"] = v.Where
	}
	return json.Marshal(result)
}

func (pfa *PfaStruct) Error(fmtString string, fmtArgs ...interface{}) error {
	return PfaError{
		pfa: pfa,
		content: &PfaErrorContent{
			reason: string(formatted.StringFrom(fmtString, fmtArgs...)),
		},
		Where: whereami.WhereAmI(3),
	}
}

func (pfa *PfaStruct) TransformError(fmtString, reason string) error {
	return PfaError{
		pfa:     pfa,
		content: &PfaErrorContent{fmtString: fmtString, reason: reason},
		Where:   whereami.WhereAmI(3),
	}
}

func (pfa *PfaStruct) UnexpectedError(err error) error {
	return PfaError{
		pfa:     pfa,
		content: &PfaErrorContent{errStr: err.Error(), state: PecsPrepared},
		Where:   whereami.WhereAmI(3),
	}
}

func (v *PfaError) PrepareErr(fmtString string, fmtArgs ...interface{}) {
	if v.content.state == PecsPrepared {
		bwerr.Panic("Already prepared %s ", bwjson.Pretty(v))
	} else {
		v.content.fmtString = fmtString + ": " + v.content.reason
		v.content.errStr = bwerr.From(v.content.fmtString, fmtArgs...).Error()
		v.content.state = PecsPrepared
	}
}

func (v *PfaError) SetErr(errStr string) {
	if v.content.state == PecsPrepared {
		bwerr.Panic("Already prepared %s ", bwjson.Pretty(v))
	} else {
		v.content.errStr = errStr
		v.content.state = PecsPrepared
	}
}

func (v PfaError) Error() (result string) {
	switch v.content.state {
	case PecsNeedPrepare:
		bwerr.Panic("NeedPrepare %s ", bwjson.Pretty(v))
	case PecsPrepared:
		result = v.content.errStr
	}
	return
}

func (v PfaError) State() PfaErrorContentState {
	return v.content.state
}

// ============================================================================
