/*
Предоставляет функции для генерации ошибок.
*/
package bwerror

import (
	"errors"
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

func Error(msgFmt string, args ...interface{}) error {
	return errors.New(Spew.Sprintf(ansi.Ansi(ansiErr, errPrefix+msgFmt), args...))
}

func ErrorErr(err error) error {
	return Error(err.Error())
}

func PanicErr(err error) {
	log.Panic(err.Error() + wherePrefix + whereami.WhereAmI(2))
}

func Panic(msgFmt string, args ...interface{}) {
	log.Panicf(ansi.Ansi(ansiErr, errPrefix+msgFmt+wherePrefix+whereami.WhereAmI(2)), args...)
}

func Panicd(depth uint, msgFmt string, args ...interface{}) {
	log.Panicf(ansi.Ansi(ansiErr, errPrefix+msgFmt+wherePrefix+whereami.WhereAmI(int(depth)+2)), args...)
}

func Noop(args ...interface{}) {
}
