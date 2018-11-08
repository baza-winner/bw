// Package bwerr содержит реализацию bw-типа Error и утилиты для работы с ним.
package bwerr

import (
	"regexp"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr/where"

	_ "github.com/baza-winner/bwcore/ansi/tags"
)

func Panic(fmtString string, fmtArgs ...interface{}) {
	PanicA(A{1, fmtString, fmtArgs})
}

func IncDepth(a bw.I, incOpt ...uint) (result bw.I) {
	inc := uint(1)
	switch t := a.(type) {
	case A:
		t.Depth += inc
		result = t
	case E:
		t.Depth += inc
		result = t
	case bw.A:
		result = A{inc, t.Fmt, t.Args}
	default:
		panic(bw.Spew.Sprintf("%#v", a))
	}
	return
}

func PanicA(a bw.I) {
	panic(FromA(IncDepth(a)))
}

func Unreachable() {
	PanicA(A{Depth: 1, Fmt: ansiUnreachable})
}

func TODO() {
	PanicA(A{Depth: 1, Fmt: ansiTODO})
}

// ============================================================================

type Error struct {
	S  string
	WW where.WW
}

// Error for error implemention
func (v Error) Error() string {
	var suffix string
	for _, w := range v.WW {
		suffix += "\n    " + w.String()
	}
	return ansiErrPrefix + v.S + ansiAt + suffix
}

var findRefineRegexp = regexp.MustCompile("{Error}")

func Refine(err error, fmtString string, fmtArgs ...interface{}) Error {
	var t Error
	var ok bool
	if t, ok = err.(Error); !ok {
		t = FromA(Err(err))
	}
	return t.RefineA(bw.A{fmtString, fmtArgs})
}

func (v Error) RefineA(a bw.I) (result Error) {
	result = v
	fmtString := findRefineRegexp.ReplaceAllStringFunc(a.FmtString(), func(s string) (result string) { return v.S })
	result.S = ansi.String(bw.Spew.Sprintf(fmtString, a.FmtArgs()...))
	return
}

// ============================================================================

type A struct {
	Depth uint
	Fmt   string
	Args  []interface{}
}

// FmtString for bw.I implementation
func (v A) FmtString() string { return v.Fmt }

// FmtArgs for bw.I implementation
func (v A) FmtArgs() []interface{} { return v.Args }

// ============================================================================

type E struct {
	Depth uint
	Error error
}

func Err(err error, optDepth ...uint) E {
	var depth uint
	if optDepth != nil {
		depth = optDepth[0]
	}
	return E{depth, err}
}

func FmtStringOf(err error) (result string) {
	if err != nil {
		if t, ok := err.(Error); ok {
			result = t.S
		} else {
			result = err.Error()
		}
	}
	return
}

// FmtString for bw.I implementation
func (v E) FmtString() string { return FmtStringOf(v.Error) }

// FmtArgs for bw.I implementation
func (v E) FmtArgs() []interface{} { return nil }

// From constructs Error
func From(fmtString string, fmtArgs ...interface{}) Error {
	return FromA(IncDepth(bw.A{fmtString, fmtArgs}))
}

func FromA(a bw.I) Error {
	var depth uint
	var fmtString string
	var fmtArgs []interface{}
	// var noNeedErrPrefix bool
	switch t := a.(type) {
	case E:
		if e, ok := t.Error.(Error); ok {
			return e
		}
		depth = t.Depth
	case A:
		depth = t.Depth
	}
	fmtString = a.FmtString()
	fmtArgs = a.FmtArgs()
	return Error{
		ansi.String(bw.Spew.Sprintf(fmtString, fmtArgs...)),
		where.WWFrom(depth + 1),
	}
}

// ============================================================================

var ansiWherePrefix string
var ansiErrPrefix string
var ansiUnreachable string
var ansiTODO string
var ansiAt string

func init() {
	ansiUnreachable = ansi.String("<ansiErr>UNREACHABLE")
	ansiTODO = ansi.String("<ansiErr>TODO")
	ansiErrPrefix = ansi.String("<ansiErr>ERR: ")
	ansiAt = ansi.String("\n  at ")
}

var newlineAtTheEnd, _ = regexp.Compile(`\n\s*$`)

// ============================================================================
