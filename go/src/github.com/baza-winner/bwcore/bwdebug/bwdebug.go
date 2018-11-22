package bwdebug

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/baza-winner/bwcore/ansi"
	_ "github.com/baza-winner/bwcore/ansi/tags"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwerr/where"
	"github.com/baza-winner/bwcore/bwjson"
)

var (
	ansiDebugVarValue              string
	ansiDebugVarValueAsVerboseHash string
	ansiDebugVarValueAsString      string
	ansiDebugMark                  string
	ansiDebugVarName               string
	ansiExpectsVarVal              string
	ansiMustBeString               string
	ansiMustBeNonEmptyString       string
	ansiMustBeNonEmptyVarName      string
)

func init() {
	ansi.MustAddTag("ansiDebugMark",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorYellow, Bright: true}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
	ansiDebugVarValue = ansi.String("<ansiVal>%v<ansi>")
	ansiDebugVarValueAsVerboseHash = ansi.String("<ansiVal>%#v<ansi>")
	ansiDebugVarValueAsString = ansi.String("<ansiVal>%s<ansi>")
	ansiDebugMark = ansi.String("<ansiDebugMark>%s<ansi>, ")
	ansiDebugVarName = ansi.String("<ansiVar>%s<ansi>: ")
	ansiExpectsVarVal = ansi.String("expects val for <ansiVar>%s")
	ansiMustBeString = "<ansiVar>args<ansiPath>.%d<ansi> (<ansiVal>%#v<ansi>) must be <ansiType>string"
	ansiMustBeNonEmptyString = "<ansiVar>args<ansiPath>.%d<ansi> must be <ansiType>non empty string"
}

func Print(args ...interface{}) {
	if s, err := stringToPrint(1, args...); err != nil {
		panic(err)
	} else {
		fmt.Println(s)
	}
}

var asJSONSuffix = regexp.MustCompile(":json$")
var asStringSuffix = regexp.MustCompile(":s$")
var asVerboseSuffix = regexp.MustCompile(":v$")
var asVerboseHashSuffix = regexp.MustCompile(":#v$")

func stringToPrint(depth uint, args ...interface{}) (result string, err error) {
	markPrefix := ""
	fmtString := ""
	fmtArgs := []interface{}{}
	expectsVal := false
	lastVar := ""
	i := 0
	type valueFormat uint8
	const (
		vfVerbose valueFormat = iota
		vfVerboseHash
		vfString
		vfJSON
	)
	vfDefault := vfVerbose
	var vf valueFormat
	// var fmtDebugVarValue string
	for _, arg := range args {
		i++
		if expectsVal == true {
			// fmt.Printf("vf: %d\n", vf)
			switch vf {
			case vfString:
				fmtString += ansiDebugVarValueAsString
				switch t := arg.(type) {
				case rune:
					fmtArgs = append(fmtArgs, string(t))
				case fmt.Stringer:
					// fmt.Printf("typeof(arg): %T\n", arg)
					// fmt.Printf("t.String(): %s\n", t.String())
					fmtArgs = append(fmtArgs, t.String())
				default:
					// fmt.Printf("typeof(arg): %T\n", arg)
					fmtArgs = append(fmtArgs, arg)
				}
			case vfJSON:
				fmtString += ansiDebugVarValueAsString
				fmtArgs = append(fmtArgs, bwjson.Pretty(arg))
			case vfVerboseHash:
				fmtString += ansiDebugVarValueAsVerboseHash
				fmtArgs = append(fmtArgs, arg)
			default:
				fmtString += ansiDebugVarValue
				fmtArgs = append(fmtArgs, arg)
			}
			expectsVal = false
		} else if valueOf := reflect.ValueOf(arg); valueOf.Kind() != reflect.String {
			err = bwerr.FromA(bwerr.A{Depth: 1, Fmt: ansiMustBeString, Args: bw.Args(i, arg)})
			return
		} else if s := valueOf.String(); len(s) == 0 {
			err = bwerr.FromA(bwerr.A{Depth: 1, Fmt: ansiMustBeNonEmptyString, Args: bw.Args(i)})
			return
		} else if s[0:1] == "!" {
			markPrefix += fmt.Sprintf(ansiDebugMark, s)
		} else {
			if len(fmtArgs) > 0 {
				fmtString += ", "
			}
			if asJSONSuffix.MatchString(s) {
				s = s[:len(s)-5]
				vf = vfJSON
			} else if asStringSuffix.MatchString(s) {
				s = s[:len(s)-2]
				vf = vfString
			} else if asVerboseSuffix.MatchString(s) {
				s = s[:len(s)-2]
				vf = vfVerbose
			} else if asVerboseHashSuffix.MatchString(s) {
				s = s[:len(s)-3]
				vf = vfVerboseHash
			} else {
				vf = vfDefault
			}
			if len(s) == 0 {
				vfDefault = vf
				// err = bwerr.FromA(bwerr.A{Depth: 1, Fmt: ansiMustBeNonEmptyVarName, Args: bw.Args(i)})
				// return
			} else {
				fmtString += fmt.Sprintf(ansiDebugVarName, s)
				lastVar = s
				expectsVal = true
			}
		}
	}
	if expectsVal {
		err = bwerr.FromA(bwerr.A{Depth: 1, Fmt: ansiExpectsVarVal, Args: bw.Args(lastVar)})
		return
	}
	result = markPrefix +
		where.MustFrom(1+depth).String() +
		": " +
		ansi.String(bw.Spew.Sprintf(fmtString, fmtArgs...))
	return
}
