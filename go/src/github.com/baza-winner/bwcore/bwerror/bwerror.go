/*
Предоставляет функции для генерации ошибок.
*/
package bwerror

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"

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
	var where string
	typeOf := reflect.TypeOf(err)
	if typeOf.Kind() == reflect.Struct {
		if sf, ok := typeOf.FieldByName("Where"); ok && sf.Type.Kind() == reflect.String {
			where = wherePrefix + reflect.ValueOf(err).FieldByName("Where").String()
		}
	}
	if len(where) == 0 {
		var depth uint
		if optDepth != nil {
			depth = optDepth[0]
		}
		where = wherePrefix + whereami.WhereAmI(int(depth)+2)
	}
	log.Panic(err.Error() + ansi.Ansi("", where))
}

func Panic(msgFmt string, args ...interface{}) {
	log.Panicf(ansi.Ansi(ansiErr, errPrefix+msgFmt+wherePrefix+whereami.WhereAmI(2)), args...)
}

func Panicd(depth uint, msgFmt string, args ...interface{}) {
	log.Panicf(ansi.Ansi(ansiErr, errPrefix+msgFmt+wherePrefix+whereami.WhereAmI(int(depth)+2)), args...)
}

func Unreachable() {
	Panicd(1, "<ansiErr>UNREACHABLE")
}

func TODO() {
	Panicd(1, "<ansiErr>TODO")
}

func Debug(args ...interface{}) {
	prefixFmt := ""
	prefixArgs := []interface{}{}
	fmtString := ""
	fmtArgs := []interface{}{}
	expectsVal := false
	lastVar := ""
	i := 0
	for _, arg := range args {
		i++
		if expectsVal == true {
			fmtString += "<ansiPrimary>%#v<ansi>"
			fmtArgs = append(fmtArgs, arg)
			expectsVal = false
		} else if valueOf := reflect.ValueOf(arg); valueOf.Kind() != reflect.String {
			Panicd(1, "expects string as arg #%d", i)
		} else if s := valueOf.String(); len(s) == 0 {
			Panicd(1, "expects non empty string as arg #%d", i)
		} else if s[0:1] == "!" {
			prefixFmt += "<ansiBold><ansiYellow>%s<ansi>, "
			prefixArgs = append(prefixArgs, arg)
		} else {
			if len(fmtArgs) > 0 {
				fmtString += ", "
			}
			fmtString += "<ansiOutline>%s<ansi>: "
			fmtArgs = append(fmtArgs, s)
			lastVar = s
			expectsVal = true
		}
	}
	if expectsVal {
		Panicd(1, "expects val for <ansiOutline>%s", lastVar)
	}

	if len(fmtString) == 0 || fmtString[len(fmtString)-1:] != "\n" {
		fmtString += "\n"
	}
	function, file, line, _ := runtime.Caller(1)
	fmtString = prefixFmt + "<ansiDarkGreen>%s<ansiDarkGray>@<ansiCmd>%s:%d<ansiBold><ansiYellow>:<ansi> " + fmtString
	fmtArgs = append([]interface{}{runtime.FuncForPC(function).Name(), chopPath(file), line}, fmtArgs...)
	fmtArgs = append(prefixArgs, fmtArgs...)
	fmt.Printf(ansi.Ansi("", fmtString), fmtArgs...)
}

// return the source filename after the last slash
func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i == -1 {
		return original
	} else {
		return original[i+1:]
	}
}
