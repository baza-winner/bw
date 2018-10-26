package b

import (
	"fmt"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/pfa/core"
	"github.com/baza-winner/bwcore/pfa/formatted"
)

type ValTransformer interface {
	TransformValue(pfa *core.PfaStruct, i interface{}) interface{}
	String() string
}

type By []ValTransformer

// ============================================================================

type PairBrace struct{}

var pairs = map[rune]rune{
	'<': '>',
	'{': '}',
	'(': ')',
	'[': ']',
}

const fmtPairBrace = "failed to get pair for %s"

func (v PairBrace) TransformValue(pfa *core.PfaStruct, i interface{}) (result interface{}) {
	var r rune
	switch t := i.(type) {
	case rune:
		r = t
	default:
		pfa.ErrVal = FailedToTransformFrom(fmtPairBrace, "nor Rune")
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

const fmtEscape = "failed to escape rune for %s"

func (v Escape) TransformValue(pfa *core.PfaStruct, i interface{}) (result interface{}) {
	var r rune
	switch t := i.(type) {
	case rune:
		r = t
	default:
		pfa.ErrVal = FailedToTransformFrom(fmtEscape, "nor Rune")
		return
	}

	var ok bool
	if result, ok = escape[r]; !ok {
		pfa.ErrVal = FailedToTransformFrom(fmtEscape, "no escape rune")
	}

	return
}

func (v Escape) String() string {
	return "Escape"
}

// ============================================================================

type ParseNumber struct{}

const fmtParseNumber = "failed to transform %s to number"

func (v ParseNumber) TransformValue(pfa *core.PfaStruct, i interface{}) (result interface{}) {
	var s string
	switch t := i.(type) {
	case string:
		s = t
	case rune:
		s = string(t)
	default:
		pfa.ErrVal = FailedToTransformFrom(fmtParseNumber, "niether String, nor Rune")
		return
	}

	result, err := core.ParseNumber(s)
	if err != nil {
		pfa.ErrVal = FailedToTransformFrom(fmtParseNumber, err.Error())
	}
	return
}

func (v ParseNumber) String() string {
	return "ParseNumber"
}

// ============================================================================

type FailedToTransform struct {
	content *FailedToTransformContent
}

type FailedToTransformContent struct {
	fmt    string
	reason string
	s      string
}

func FailedToTransformFrom(fmt, reason string) FailedToTransform {
	return FailedToTransform{&FailedToTransformContent{fmt: fmt, reason: reason}}
}

func (t *FailedToTransform) Prepare(source formatted.String) {
	t.content.s = fmt.Sprintf(t.content.fmt+" (%s)", source, t.content.reason)
}

func (v FailedToTransform) Error(pfa *core.PfaStruct) error {
	return pfa.Proxy.ItemError(pfa.GetTopStackItem().Start, v.content.s)
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
