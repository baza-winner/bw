package b

import (
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/pfa/core"
)

type ValTransformer interface {
	TransformValue(pfa *core.PfaStruct, val interface{}) (interface{}, error)
	String() string
}

type By []ValTransformer

// ============================================================================

// type TransformError struct {
// 	content *TransformErrorContent
// 	Where   string
// }

// type TransformErrorContent struct {
// 	pfa       *PfaStruct
// 	reason    string
// 	fmtString string
// 	fmtArgs   []interface{}
// }

// func (v *TransformError) Prepare(pfa *core.PfaStruct, fmtString string, fmtArgs ...interface{}) {
// 	v.content.fmtString = t.content.fmtString + " (" + t.content.reason + ")"
// 	v.content.fmtArgs = fmtArgs
// 	// t.content.s = fmt.Sprintf(t.content.fmt+" (%s)", source, t.content.reason)
// }

// func (v TransformError) PfaError(pfa *core.PfaStruct) error {
// 	// return pfa.Proxy.ItemError(pfa.GetTopStackItem().Start, v.content.s)
// 	return pfa.Proxy.Unexpected(pfa.GetTopStackItem().Start, bwfmt.StructFrom(v.content.s))
// }

// func (v PfaError) Error() (result string) {
// 	switch v.content.state {
// 	case PecsNeedPrepare:
// 		bwerror.Panic("%#v", v.content.state)
// 	case PecsPrepared:
// 		result = bwerror.Error(v.content.fmtString, v.content.fmtArgs)
// 	}
// 	return
// }

// ============================================================================

type PairBrace struct{}

var pairs = map[rune]rune{
	'<': '>',
	'{': '}',
	'(': ')',
	'[': ']',
}

const fmtPairBrace = "failed to get PairBrace for %s"

func (v PairBrace) TransformValue(pfa *core.PfaStruct, val interface{}) (result interface{}, err error) {
	var r rune
	switch t := val.(type) {
	case rune:
		r = t
	default:
		err = pfa.TransformError(fmtPairBrace, "nor Rune")
		return
	}

	var ok bool
	if result, ok = pairs[r]; !ok {
		result = r
	}

	return
}

func (v PairBrace) String() string {
	return "PairBrace"
}

// ============================================================================

type Escape struct{}

var escape = map[rune]rune{
	'a': '\a',
	'b': '\b',
	'f': '\f',
	'n': '\n',
	'r': '\r',
	't': '\t',
	'v': '\v',
}

const fmtEscape = "failed to get Escape rune for %s"

func (v Escape) TransformValue(pfa *core.PfaStruct, val interface{}) (result interface{}, err error) {
	var r rune
	switch t := val.(type) {
	case rune:
		r = t
	default:
		err = pfa.TransformError(fmtEscape, "nor Rune")
		return
	}

	var ok bool
	if result, ok = escape[r]; !ok {
		err = pfa.TransformError(fmtEscape, "no escape rune")
		return
	}

	return
}

func (v Escape) String() string {
	return "Escape"
}

// ============================================================================

type ParseNumber struct{}

const fmtParseNumber = "failed to ParseNumber from %s"

func (v ParseNumber) TransformValue(pfa *core.PfaStruct, val interface{}) (result interface{}, err error) {
	var s string
	switch t := val.(type) {
	case string:
		s = t
	case rune:
		s = string(t)
	default:
		pfa.TransformError(fmtParseNumber, "niether String, nor Rune")
		return
	}

	result, err = core.ParseNumber(s)
	if err != nil {
		err = pfa.TransformError(fmtParseNumber, err.Error())
	}
	return
}

func (v ParseNumber) String() string {
	return "ParseNumber"
}

// ============================================================================

type ButLast struct{}

const fmtButLast = "failed to get ButLast from %s"

func (v ButLast) TransformValue(pfa *core.PfaStruct, val interface{}) (result interface{}, err error) {
	var vals []interface{}
	switch t := val.(type) {
	case []interface{}:
		vals = t
	// case rune:
	// s = string(t)
	default:
		err = pfa.TransformError(fmtButLast, "nor Array")
		return
	}

	if len(vals) > 0 {
		result = vals[:len(vals)-1]
	} else {
		err = pfa.TransformError(fmtButLast, "Array is empty")
		return
	}

	return
}

func (v ButLast) String() string {
	return "ButLast"
}

// ============================================================================

type Append struct{}

func (v Append) TransformValue(pfa *core.PfaStruct, val interface{}) (result interface{}, err error) {
	bwerror.Unreachable()
	return
}

func (v Append) String() string {
	return "Append"
}

// ============================================================================

type AppendSlice struct{}

func (v AppendSlice) TransformValue(pfa *core.PfaStruct, val interface{}) (result interface{}, err error) {
	bwerror.Unreachable()
	// bwerror.Unreachable()
	return
}

func (v AppendSlice) String() string {
	return "AppendSlice"
}

// ============================================================================
