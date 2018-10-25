/*
Предоставляет функции для генерации ошибок.
*/
package bwerror

import (
	"fmt"
	"log"
	"os"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/davecgh/go-spew/spew"
	"github.com/jimlawless/whereami"
)

var Spew spew.ConfigState

func init() {
	Spew = spew.ConfigState{SortKeys: true}
}

const (
	ansiErr     = `Reset`
	errPrefix   = `<ansiErr>ERR:<ansi> `
	wherePrefix = ` <ansiCmd>at `
)

func ExitWithError(exitCode int, fmtString string, fmtArgs ...interface{}) {
	msg := ansi.Ansi(ansiErr, fmt.Sprintf(errPrefix+fmtString, fmtArgs...))
	if byte(msg[len(msg)-1]) != '\n' {
		msg += string('\n')
	}
	fmt.Print(msg)
	os.Exit(exitCode)
}

type WhereError struct {
	Where  string
	errStr string
}

func (v WhereError) Error() string {
	return v.errStr
}

func Error(msgFmt string, args ...interface{}) error {
	return WhereError{
		whereami.WhereAmI(2),
		Spew.Sprintf(ansi.Ansi(ansiErr, errPrefix+msgFmt), args...),
	}
}

func Errord(depth uint, msgFmt string, args ...interface{}) error {
	return WhereError{
		whereami.WhereAmI(int(depth) + 2),
		Spew.Sprintf(ansi.Ansi(ansiErr, errPrefix+msgFmt), args...),
	}
}

func ErrorErr(err error, optDepth ...uint) error {
	var depth uint
	if optDepth != nil {
		depth = optDepth[0]
	}
	return WhereError{
		whereami.WhereAmI(int(depth) + 2),
		Spew.Sprintf(ansi.Ansi(ansiErr, errPrefix+err.Error())),
	}
}

func PanicErr(err error, optDepth ...uint) {
	var depth uint
	if optDepth != nil {
		depth = optDepth[0]
	}
	log.Panic(err.Error() + ansi.Ansi("", wherePrefix+whereami.WhereAmI(int(depth)+2)))
}

func Panic(msgFmt string, args ...interface{}) {
	log.Panicf(ansi.Ansi(ansiErr, errPrefix+msgFmt+wherePrefix+whereami.WhereAmI(2)), args...)
}

func Panicd(depth uint, msgFmt string, args ...interface{}) {
	log.Panicf(ansi.Ansi(ansiErr, errPrefix+msgFmt+wherePrefix+whereami.WhereAmI(int(depth)+2)), args...)
}

func Noop(args ...interface{}) {
}

// func Unreachable(args ...string) {
// 	var suffix string
// 	if args != nil {
// 		suffix = " <ansi>" + strings.Join(args, " ")
// 	}
// 	Panicd(1, "<ansiErr>UNREACHABLE"+suffix)
// }

// func TODO() {
// 	Panicd(1, "<ansiErr>UNREACHABLE")
// }
