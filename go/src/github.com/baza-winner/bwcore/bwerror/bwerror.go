/*
Предоставляет функции для генерации ошибок.
*/
package bwerror

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwfmt"
	"github.com/jimlawless/whereami"
)

// const (
// 	wherePrefix = " <ansiCmd>at "
// 	errPrefix   = "<ansiErr>ERR:<ansi> "
// )

var ansiWherePrefix string
var ansiErrPrefix string
var ansiUnreachable string
var ansiTODO string

func init() {
}

func init() {
	ansi.MustAddTag("ansiErr",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorRed, Bright: true}),
	)
	ansi.MustAddTag("ansiDebugVarValue",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorCyan, Bright: true}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
	ansi.MustAddTag("ansiDebugVarName",
		ansi.SGRCodeOfColor256(ansi.Color256{Code: 201}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
	ansi.MustAddTag("ansiDebugMark",
		ansi.SGRCodeOfColor256(ansi.Color256{Code: 201}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
	ansi.MustAddTag("ansiDebugAt",
		ansi.SGRCodeOfColor256(ansi.Color256{Code: 243}),
	)
	ansiUnreachable = ansi.String(ansi.A{S: "<ansiErr>UNREACHABLE"})
	ansiTODO = ansi.String(ansi.A{S: "<ansiErr>TODO"})
	ansiWherePrefix = ansi.String(ansi.A{S: " <ansiCmd>at "})
	ansiErrPrefix = ansi.String(ansi.A{S: "<ansiErr>ERR:<ansi> "})
}

var newlineAtTheEnd, _ = regexp.Compile(`\n\s*$`)

func Exit(exitCode int, a bwfmt.I) {
	err := From(a)
	msg := err.ansiString
	if !newlineAtTheEnd.MatchString(ansi.ChopReset(err.ansiString)) {
		msg += string('\n')
	}
	fmt.Print(msg)
	os.Exit(exitCode)
}

type Error struct {
	ansiString string
	Where      string
}

func (v Error) Error() string {
	return v.ansiString
}

type A struct {
	Depth uint
	Fmt   string
	Args  []interface{}
}

func (v A) FmtString() string {
	return v.Fmt
}

func (v A) FmtArgs() []interface{} {
	return v.Args
}

type E struct {
	Depth uint
	Error error
}

func (v E) FmtString() string {
	return v.Error.Error()
}

func (v E) FmtArgs() []interface{} {
	return nil
}

func From(a bwfmt.I) Error {
	var depth uint
	var fmtString string
	var fmtArgs []interface{}
	var noNeedErrPrefix bool
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
	ansiPrefix := ""
	if !noNeedErrPrefix {
		ansiPrefix = ansiErrPrefix
	}
	return Error{
		ansi.Concat(
			ansiPrefix,
			ansi.String(ansi.A{S: bw.Spew.Sprintf(fmtString, fmtArgs...)}),
		),
		whereami.WhereAmI(int(depth) + 2),
	}
}

func Panic(a bwfmt.I) {
	err := From(a)
	log.Panic(ansi.Concat(
		err.ansiString,
		ansiWherePrefix,
		err.Where,
	))
}

func Unreachable() {
	Panic(A{Depth: 1, Fmt: ansiUnreachable})
}

func TODO() {
	Panic(A{Depth: 1, Fmt: ansiTODO})
}

var ansiDebugVarValue string
var ansiDebugMark string

func init() {
	ansiDebugVarValue = ansi.String(ansi.A{S: "<ansiDebugVarValue>%#v<ansi>"})
}

func Debug(args ...interface{}) {
	markPrefix := ""
	// prefixArgs := []interface{}{}
	fmtString := ""
	fmtArgs := []interface{}{}
	expectsVal := false
	lastVar := ""
	i := 0
	for _, arg := range args {
		i++
		if expectsVal == true {
			fmtString += "<ansiDebugVarValue>%#v<ansi>"
			fmtArgs = append(fmtArgs, arg)
			expectsVal = false
		} else if valueOf := reflect.ValueOf(arg); valueOf.Kind() != reflect.String {
			Panic(A{1, "expects string as arg #%d", bw.Args(i)})
		} else if s := valueOf.String(); len(s) == 0 {
			Panic(A{1, "expects non empty string as arg #%d", bw.Args(i)})
		} else if s[0:1] == "!" {
			markPrefix += "<ansiDebugMark>" + s + "<ansi>, "
			// prefixArgs = append(prefixArgs, arg)
		} else {
			if len(fmtArgs) > 0 {
				fmtString += ", "
			}
			fmtString += "<ansiDebugVarName>" + s + "<ansi>: "
			// fmtArgs = append(fmtArgs, s)
			lastVar = s
			expectsVal = true
		}
	}
	if expectsVal {
		Panic(A{1, "expects val for <ansiDebugVarName>%s", bw.Args(lastVar)})
	}
	if len(fmtString) == 0 || fmtString[len(fmtString)-1:] != "\n" {
		fmtString += "\n"
	}
	function, file, line, _ := runtime.Caller(1)
	fmtString = prefixFmt + "<ansiDebugFuncName>%s<ansiDebugAt>@<ansiDebugFile>%s:%d<ansiDebugMark>:<ansi> " + fmtString
	// fmtArgs = append([]interface{}{runtime.FuncForPC(function).Name(), chopPath(file), line}, fmtArgs...)
	// fmtArgs = append(prefixArgs, fmtArgs...)

	// bw.Spew.Sprintf(fmtString, fmtArgs...)
	fmt.Print(
		ansi.Concat(
			ansi.String(markPrefix),

		ansi.String(ansi.A{S: bw.Spew.Sprintf(fmtString, fmtArgs...)})
			)
		)
}

type Where struct {
	funcName string
	fileSpec string
	line     int
}

func WhereFrom(depth uint) (result Where, ok bool) {
	function, file, line, ok := runtime.Caller(int(depth) + 1)
	if !ok {
		return
	}
	result = Where{
		runtime.FuncForPC(function).Name(),
		file,
		line,
	}
	return
}

func (v Where) FuncName() string {
	return v.funcName
}

func (v Where) FileSpec() string {
	return v.fileSpec
}

func (v Where) FileName() string {
	return chopPath(v.fileSpec)
}

func (v Where) Line() int {
	return v.line
}

var ansiWhere string

func init() {
	ansi.AddTag("ansiWhereFuncName", ansi.SGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen}))
	ansi.AddTag("ansiWhereAt", ansi.SGRCodeOfColor256(ansi.Color256{Code: 243}))
	ansi.AddTag("ansiWhereFile", ansi.SGRCodeOfColor256(ansi.Color8{Color: ansi.SGRColorWhite, Bright: true}))
	ansiWhere = "<ansiWhereFuncName>%s<ansiWhereAt>@<ansiWhereFile>%s:%d"
}

func (v Where) String() string {
	return bw.Spew.Sprintf(ansiWhere, v.FuncName(), v.FileName(), v.Line())
	// return v.line
	// ansi.Str
	// fmtString = "<ansiWhereFuncName>%s<ansiWhereAt>@<ansiWhereFile>%s:%d<ansiWhereColon>" + fmtString

	// fmtArgs = append([]interface{}{runtime.FuncForPC(function).Name(), chopPath(file), line}, fmtArgs...)
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
