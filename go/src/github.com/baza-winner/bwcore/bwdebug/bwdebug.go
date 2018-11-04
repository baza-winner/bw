package bwdebug

import (
	"fmt"
	"reflect"

	"github.com/baza-winner/bwcore/ansi"
	_ "github.com/baza-winner/bwcore/ansi/tags"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwerr/where"
)

var (
	ansiDebugVarValue        string
	ansiDebugMark            string
	ansiDebugVarName         string
	ansiExpectsVarVal        string
	ansiMustBeString         string
	ansiMustBeNonEmptyString string
)

func init() {
	ansi.MustAddTag("ansiDebugMark",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorYellow, Bright: true}),
		// ansi.SGRCodeOfColor256(ansi.Color256{Code: 201}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
	ansiDebugVarValue = ansi.String("<ansiVal>%#v<ansi>")
	ansiDebugMark = ansi.String("<ansiDebugMark>%s<ansi>, ")
	ansiDebugVarName = ansi.String("<ansiVar>%s<ansi>: ")
	ansiExpectsVarVal = ansi.String("expects val for <ansiVar>%s")
	ansiMustBeString = "<ansiVar>args<ansiPath>.%d<ansi> (<ansiVal>%#v<ansi>) must be <ansiType>string"
	ansiMustBeNonEmptyString = "<ansiVar>args<ansiPath>.%d<ansi> must be <ansiType>non empty string"
}

func Print(args ...interface{}) {
	if s, err := ansiString(1, args...); err != nil {
		panic(err)
	} else {
		fmt.Println(s)
	}
}

func ansiString(depth uint, args ...interface{}) (result string, err error) {
	markPrefix := ""
	fmtString := ""
	fmtArgs := []interface{}{}
	expectsVal := false
	lastVar := ""
	i := 0
	for _, arg := range args {
		i++
		if expectsVal == true {
			fmtString += ansiDebugVarValue
			fmtArgs = append(fmtArgs, arg)
			expectsVal = false
		} else if valueOf := reflect.ValueOf(arg); valueOf.Kind() != reflect.String {
			err = bwerr.FromA(bwerr.A{1, ansiMustBeString, bw.Args(i, arg)})
			return
		} else if s := valueOf.String(); len(s) == 0 {
			err = bwerr.FromA(bwerr.A{1, ansiMustBeNonEmptyString, bw.Args(i)})
			return
		} else if s[0:1] == "!" {
			markPrefix += fmt.Sprintf(ansiDebugMark, s)
		} else {
			if len(fmtArgs) > 0 {
				fmtString += ", "
			}
			fmtString += fmt.Sprintf(ansiDebugVarName, s)
			lastVar = s
			expectsVal = true
		}
	}
	if expectsVal {
		err = bwerr.FromA(bwerr.A{1, ansiExpectsVarVal, bw.Args(lastVar)})
		return
	}
	result = markPrefix +
		where.MustFrom(1+depth).String() +
		": " +
		ansi.String(bw.Spew.Sprintf(fmtString, fmtArgs...))
	// result = ansi.Concat(
	// 	markPrefix,
	// 	where.MustFrom(1+depth).String(),
	// 	": ",
	// 	ansi.String(bw.Spew.Sprintf(fmtString, fmtArgs...)),
	// )
	return
}
