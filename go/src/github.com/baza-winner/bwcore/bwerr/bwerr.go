// Package bwerr предоставляет функции для генерации ошибок.
package bwerr

import (
	"log"
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

// func depthOf(a bw.I) (result uint) {
// 	switch t := a.(type) {
// 	case A:
// 		result = t.Depth
// 	case E:
// 		result = t.Depth
// 	}
// 	return result
// }

// func ModifyBy(a bw.I, b bw.I, prepend bool) (result bw.I) {
// 	var fmtString string
// 	var fmtArgs []interface{}
// 	switch t := a.(type) {
// 	case A, bw.A:
// 		switch t2 := b.(type) {
// 		case A, bw.A:
// 			if prepend {
// 				fmtString = b.FmtString() + a.FmtString()
// 				fmtArgs = append(b.FmtArgs(), a.FmtArgs())
// 			} else {
// 				fmtString = a.FmtString() + b.FmtString()
// 				fmtArgs = append(a.FmtArgs(), b.FmtArgs())
// 			}
// 		case E:
// 			var errStr string
// 			if e, ok := t2.Error.(Error); ok {
// 				errStr = e.Ansi
// 			} else {
// 				errStr = t2.Error.Error()
// 			}
// 			fmtArgs = a.FmtArgs()
// 			if prepend {
// 				fmtString = errStr + a.FmtString()
// 			} else {
// 				fmtString = a.FmtString() + errStr
// 			}
// 		default:
// 			panic(bw.Spew.Sprintf("%#v", a))
// 		}
// 	case E:
// 		var errStr string
// 		if e, ok := t.Error.(Error); ok {
// 			errStr = e.Ansi
// 		} else {
// 			errStr = t.Error.Error()
// 		}
// 		switch t2 := b.(type) {
// 		case A, bw.A:
// 			if prepend {
// 				fmtString = b.FmtString() + errStr
// 			} else {
// 				fmtString = errStr + b.FmtString()
// 			}
// 		case E:
// 			var errStr2 string
// 			if e2, ok := t2.Error.(Error); ok {
// 				errStr2 = e2.Ansi
// 			} else {
// 				errStr2 = t2.Error.Error()
// 			}
// 			if prepend {
// 				fmtString = errStr2 + errStr
// 			} else {
// 				fmtString = errStr + errStr2
// 			}
// 		default:
// 			panic(bw.Spew.Sprintf("%#v", a))
// 		}
// 	default:
// 		panic(bw.Spew.Sprintf("%#v", a))
// 	}
// 	result = A{depthOf(a), fmtString, fmtArgs}
// 	return
// }

func PanicA(a bw.I) {
	log.Panic(FromA(IncDepth(a)).Error())
}

func Unreachable() {
	PanicA(A{Depth: 1, Fmt: ansiUnreachable})
}

func TODO() {
	PanicA(A{Depth: 1, Fmt: ansiTODO})
}

// func Exit(exitCode int, a bw.I) {
// func Exit(exitCode int, a bw.I) {
// 	err := FromA(a)
// 	msg := err.Ansi
// 	if !newlineAtTheEnd.MatchString(ansi.ChopReset(msg)) {
// 		msg += string('\n')
// 	}
// 	fmt.Print(msg)
// 	os.Exit(exitCode)
// }

// ============================================================================

type Error struct {
	Ansi  string
	Where where.W
}

// func (v Error) EqualToTest(err error) (result bool) {
// 	result = v.Ansi == FmtStringOf(err)
// 	return
// }

// Error for error implemention
func (v Error) Error() string {
	return ansiErrPrefix + v.Ansi + " at " + v.Where.String()
	// ,
	// return ansi.Concat(
	// 	ansiErrPrefix,
	// 	v.Ansi,
	// 	" at ",
	// 	v.Where.String(),
	// )
}

var findRefineRegexp = regexp.MustCompile("{Error}")

func (v Error) Refine(fmtString string, fmtArgs ...interface{}) Error {
	return v.RefineA(bw.A{fmtString, fmtArgs})
}

func (v Error) RefineA(a bw.I) (result Error) {
	result = v
	result.Ansi = ansi.String(bw.Spew.Sprintf(
		findRefineRegexp.ReplaceAllString(a.FmtString(), v.Ansi),
		a.FmtArgs()...,
	))
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
			result = t.Ansi
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
		where.MustFrom(depth + 1),
	}
}

// func (v Error) Refine() {

// }

// ============================================================================

var ansiWherePrefix string
var ansiErrPrefix string
var ansiUnreachable string
var ansiTODO string

func init() {
	ansiUnreachable = ansi.String("<ansiErr>UNREACHABLE")
	ansiTODO = ansi.String("<ansiErr>TODO")
	ansiErrPrefix = ansi.String("<ansiErr>ERR: ")
}

var newlineAtTheEnd, _ = regexp.Compile(`\n\s*$`)

// ============================================================================
