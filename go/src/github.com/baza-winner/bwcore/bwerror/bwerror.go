/*
Предоставляет функции для генерации ошибок.
*/
package bwerror

import (
	"errors"
	"fmt"
	"github.com/baza-winner/bwcore/ansi"
	"github.com/jimlawless/whereami"
	"log"
	"os"
)

func ExitWithError(exitCode int, fmtString string, fmtArgs ...interface{}) {
	log.Print(ansi.Ansi(`Err`, fmt.Sprintf(fmtString, fmtArgs...)))
	os.Exit(exitCode)
}

func Error(msgFmt string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(ansi.Ansi(`Err`, "ERR: "+msgFmt), args...))
}

func PanicErr(err error) {
	log.Panic(err.Error() + ansi.Ansi("Cmd", " at "+whereami.WhereAmI(2)))
}

func Panic(msgFmt string, args ...interface{}) {
	log.Panicf(ansi.Ansi(`Err`, "ERR: "+msgFmt+` <ansiCmd>at `+whereami.WhereAmI(2)), args...)
}

func Panicd(depth uint, msgFmt string, args ...interface{}) {
	log.Panicf(ansi.Ansi(`Err`, "ERR: "+msgFmt+` <ansiCmd>at `+whereami.WhereAmI(int(depth)+2)), args...)
}

func Noop(args ...interface{}) {
}