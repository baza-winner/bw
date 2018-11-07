// Предоставялет функции для тестирования.
package bwtesting

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwerr/where"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/kylelemons/godebug/pretty"
	// "log"
)

type Case struct {
	V     interface{}
	In    []interface{}
	Out   []interface{}
	Panic interface{}
}

var (
	ansiSecondArgToBeFunc              string
	ansiTestPrefix                     string
	ansiTestHeading                    string
	ansiExpectsCountParams             string
	ansiExpectsCountParamsWithVariadic string
	ansiExpectsFuncParams              string
	ansiExpectsParamType               string
	ansiExpectsOneReturnValue          string
	ansiExpectsTypeOfReturnValue       string
	ansiPath                           string
	ansiTestTitleFunc                  string
	ansiTestTitleOpenBrace             string
	ansiTestTitleSep                   string
	ansiTestTitleVal                   string
	ansiTestTitleCloseBrace            string
	ansiErr                            string
	ansiVal                            string
	ansiDiffBegin                      string
	ansiDiffEnd                        string
	errorType                          reflect.Type
)

func init() {
	errorType = reflect.TypeOf((*error)(nil)).Elem()
	ansiSecondArgToBeFunc = ansi.String("BwRunTests: second arg (<ansiVal>%#v<ansi>) to be <ansiType>func")
	ansiTestHeading = ansi.StringA(ansi.A{
		Default: []ansi.SGRCode{
			ansi.SGRCodeOfColor256(ansi.Color256{Code: 248}),
			ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
		},
		S: "Running test case <ansiVal>%q",
	})
	testPrefixFmt := "<ansiFunc>%s<ansiVar>.tests<ansiPath>.%q"
	ansiExpectsCountParams = ansi.String(testPrefixFmt + ".%s<ansi>: ожидается <ansiVal>%d<ansi> %s вместо <ansiVal>%d")
	ansiExpectsCountParamsWithVariadic = ansi.String(testPrefixFmt + ".%s<ansi>: ожидается не менее <ansiVal>%d<ansi> %s вместо <ansiVal>%d")
	ansiExpectsFuncParams = ansi.String(testPrefixFmt + ".%s.%d<ansi>: ожидаются следующие параметры функции: (<ansiVar>testName <ansiType>string<ansi>) или (<ansiVar>test <ansiType>bwtesting.Case<ansi>) или (<ansiVar>testName <ansiType>string<ansi>, <ansiVar>test <ansiType>bwtesting.Case<ansi>) или (<ansiVar>test <ansiType>bwtesting.Case<ansi>, <ansiVar>testName <ansiType>string<ansi>)")
	ansiExpectsParamType = ansi.String(testPrefixFmt + ".%s.%d<ansi>: ожидается <ansiType>%s<ansi> вместо <ansiType>%s<ansi> (<ansiVal>%#v<ansi>)")
	ansiExpectsOneReturnValue = ansi.String(testPrefixFmt + ".%s.%d<ansi>: ожидается <ansiVal>1<ansi> возвращаемое значение вместо <ansiVal>%d")
	ansiExpectsTypeOfReturnValue = ansi.String(testPrefixFmt + ".%s.%d<ansi>: в качесте возвращаемого значения ожидается <ansiType>%s<ansi> вместо <ansiType>%s<ansi>")
	// ansiSeparator = ":\n"
	ansiPath = ansi.String("<ansiPath>.%d<ansi>")
	ansiTestTitleFunc = ansi.String("<ansiFunc>%s")
	ansiTestTitleOpenBrace = ansi.String("(")
	ansiTestTitleSep = ansi.String(",")
	ansiTestTitleVal = ansi.String("<ansiVal>%#v")
	ansiTestTitleCloseBrace = ansi.String(")")
	ansiErr = ansi.String(
		"<ansiErr>tst err<ansi>: '%s'" +
			"\n<ansiOK>eta err<ansi>: '%s'" +
			"\n------------------------\n" +
			"\n<ansiErr>tst(q)<ansi>: %q" +
			"\n<ansiOK>eta(q)<ansi>: %q" +
			"\n------------------------" +
			"\n<ansiErr>tst(json)<ansi>: %s" +
			"\n<ansiOK>eta(json)<ansi>: %s" +
			"\n------------------------",
	)
	ansiVal = ansi.String(
		"<ansiErr>tst(v)<ansi>: %#v" +
			"\n<ansiOK>eta(v)<ansi>: %#v" +
			"\n------------------------" +
			"\n<ansiErr>tst(json)<ansi>: %s" +
			"\n<ansiOK>eta(json)<ansi>: %s" +
			"\n------------------------",
	)
	ansiDiffBegin = "------ BEGIN DIFF ------\n"
	ansiDiffEnd = "\n------- END DIFF -------\n"
}

func BwRunTests(t *testing.T, testee interface{}, tests map[string]Case) {

	w := where.MustFrom(1)
	testFunc := w.FuncName()

	var (
		testeeType    reflect.Type
		testeeValue   reflect.Value
		testeeFunc    string
		ansiTestTitle string
		inDef         []reflect.Type
		outDef        []reflect.Type
		numIn         int
		numOut        int
		inVals        []interface{}
		fmtString     string
		fmtArgs       []interface{}
	)

	checkTesInQt := func(test Case, testName string) {
		if !testeeType.IsVariadic() {
			if len(test.In) != numIn {
				bwerr.Panic(
					ansiExpectsCountParams,
					testFunc, testName, "In",
					numIn, bw.PluralWord(numIn, "параметр", "", "а", "ов"), len(test.In),
				)
			}
		} else {
			if len(test.In) < numIn-1 {
				bwerr.Panic(
					ansiExpectsCountParamsWithVariadic,
					testFunc, testName, "In",
					numIn-1, bw.PluralWord(numIn, "параметр", "а", "ов"), len(test.In),
				)
			}
		}
	}

	prepareInDefAndAnsiTestTitle := func(qt int) {
		inDef = []reflect.Type{}
		for i := 0; i < numIn; i++ {
			inType := testeeType.In(i)
			inDef = append(inDef, inType)
		}

		ansiTestTitle = ansiTestTitleFunc + ansiTestTitleOpenBrace
		for i := 0; i < qt; i++ {
			if i > 0 {
				ansiTestTitle += ansiTestTitleSep
			}
			ansiTestTitle += ansiTestTitleVal
		}
		ansiTestTitle += ansiTestTitleCloseBrace
	}

	initFmt := func(suffixProvider func() string) {
		fmtArgs = append(bw.Args(testeeFunc), inVals...)
		fmtString = ansiTestTitle + suffixProvider() + ":\n"
		// bwdebug.Print("fmtArgs", fmtArgs, "fmtString", fmtString)
	}

	if _, ok := testee.(string); !ok {
		testeeType = reflect.TypeOf(testee)
		if testeeType.Kind() != reflect.Func {
			bwerr.Panic("reflect.TypeOf(testee).Kind(): " + testeeType.Kind().String() + "\n")
		}
		testeeValue = reflect.ValueOf(testee)
		testeeFunc = testFunc[4:]

		numIn = testeeType.NumIn()
		if !testeeType.IsVariadic() {
			prepareInDefAndAnsiTestTitle(numIn)
		}
		outDef = []reflect.Type{}
		numOut = testeeType.NumOut()
		for i := 0; i < numOut; i++ {
			outDef = append(outDef, testeeType.Out(i))
		}
	}

	for testName, test := range tests {
		t.Logf(ansiTestHeading, testName)

		if s, ok := testee.(string); ok {
			testeeValue = reflect.ValueOf(test.V).MethodByName(s)
			testeeType = testeeValue.Type()
			testeeFunc = fmt.Sprintf("%T.%s", test.V, s)
			ansiTestTitle = ansiTestTitleFunc + ansiTestTitleOpenBrace
			numIn = testeeType.NumIn()
			checkTesInQt(test, testName)
			prepareInDefAndAnsiTestTitle(len(test.In))

			ansiTestTitle += ansiTestTitleCloseBrace
			outDef = []reflect.Type{}
			numOut = testeeType.NumOut()
			for i := 0; i < numOut; i++ {
				outDef = append(outDef, testeeType.Out(i))
			}
		} else if testeeType.IsVariadic() {
			checkTesInQt(test, testName)
			prepareInDefAndAnsiTestTitle(len(test.In))
		} else {
			checkTesInQt(test, testName)
		}

		inValues := []reflect.Value{}
		inVals = []interface{}{}
		for i := 0; i < len(test.In); i++ {
			inValue := getInOutValue(test.In, "In", inDef, i, testFunc, testName, test, testeeType.IsVariadic())
			inValues = append(inValues, inValue)
			inVals = append(inVals, inValue.Interface())
		}

		if len(test.Out) != numOut {
			bwerr.Panic(
				ansiExpectsCountParams,
				testFunc, testName, "Out",
				numOut, bw.PluralWord(numOut, "параметр", "", "а", "ов"), len(test.Out),
			)
		}
		outEtaVals := []interface{}{}
		for i := 0; i < numOut; i++ {
			outValue := getInOutValue(test.Out, "Out", outDef, i, testFunc, testName, test)
			outEtaVals = append(outEtaVals, outValue.Interface())
		}

		var panicVal interface{}
		var outValues []reflect.Value
		func() {
			if test.Panic != nil {
				defer func() { panicVal = recover() }()
			}
			outValues = testeeValue.Call(inValues)
		}()

		if panicVal != nil {
			initFmt(func() string { return ".Panic" })
			if cmpErrs(panicVal, test.Panic, &fmtString, &fmtArgs) {
				t.Error(bw.Spew.Sprintf(fmtString, fmtArgs...))
			}
		} else {
			for i := 0; i < numOut; i++ {
				initFmt(func() (result string) {
					if numOut > 1 {
						result = ansiPath
						fmtArgs = append(fmtArgs, i)
					}
					return
				})

				cmpFunc := cmpVals
				if outDef[i].Implements(reflect.TypeOf((*error)(nil)).Elem()) {
					cmpFunc = cmpErrs
				}

				if cmpFunc(outValues[i].Interface(), outEtaVals[i], &fmtString, &fmtArgs) {
					t.Error(bw.Spew.Sprintf(fmtString, fmtArgs...))
				}
			}
		}
	}
}

func getInOutValue(vals []interface{}, path string, def []reflect.Type, i int, testFunc, testName string, test Case, optIsVariadic ...bool) (result reflect.Value) {
	val := vals[i]
	if val == nil {
		result = reflect.New(def[i]).Elem()
	} else {
		result = reflect.ValueOf(val)
		if result.Kind() == reflect.Func {
			if valType := reflect.TypeOf(val); valType.NumOut() != 1 {
				bwerr.Panic(
					ansiExpectsOneReturnValue,
					testFunc, testName, path,
					valType.NumOut(),
				)
			} else if valType.Out(0) != def[i] {
				bwerr.Panic(
					ansiExpectsTypeOfReturnValue,
					testFunc, testName, path, i,
					def[i],
					valType.Out(0),
				)
			} else {
				switch valType.NumIn() {
				case 0:
					result = reflect.ValueOf(val).Call([]reflect.Value{})[0]
				case 1:
					if valType.In(0).Kind() == reflect.String {
						result = reflect.ValueOf(val).Call([]reflect.Value{
							reflect.ValueOf(testName),
						})[0]
					} else if valType.In(0).Name() == "Case" {
						result = reflect.ValueOf(val).Call([]reflect.Value{
							reflect.ValueOf(test),
						})[0]
					} else {
						bwerr.Panic(
							ansiExpectsFuncParams,
							testFunc, testName, path, i,
						)
					}
				case 2:
					if valType.In(0).Kind() == reflect.String && valType.In(1).Name() == "Case" {
						result = reflect.ValueOf(val).Call([]reflect.Value{
							reflect.ValueOf(testName),
							reflect.ValueOf(test),
						})[0]
					} else if valType.In(1).Kind() == reflect.String && valType.In(0).Name() == "Case" {
						result = reflect.ValueOf(val).Call([]reflect.Value{
							reflect.ValueOf(test),
							reflect.ValueOf(testName),
						})[0]
					} else {
						bwerr.Panic(
							ansiExpectsFuncParams,
							testFunc, testName, path, i,
						)
					}
				default:
					bwerr.Panic(
						ansiExpectsFuncParams,
						testFunc, testName, path, i,
					)
				}
			}
		}
	}
	j := i
	if j >= len(def) {
		j = len(def) - 1
	}
	if def[j].Kind() != reflect.Interface {
		if i >= len(def)-1 && len(optIsVariadic) > 0 && optIsVariadic[0] {
			if def[j].Elem().Kind() != reflect.Interface && result.Kind() != def[j].Elem().Kind() {
				// bwdebug.Print("def", def, "j", j, "i", "i")
				bwerr.Panic(
					ansiExpectsParamType,
					testFunc, testName, path, i,
					def[j].Elem().Kind(),
					result.Kind(),
					val,
				)
			}
		} else if result.Kind() != def[i].Kind() {
			bwerr.Panic(
				ansiExpectsParamType,
				testFunc, testName, path, i,
				def[j].Kind(),
				result.Kind(),
				val,
			)
		}
	} else if def[j].Implements(errorType) {
		if result.Type().Kind() != reflect.String && !result.Type().Implements(errorType) {
			bwerr.Panic(
				ansiExpectsParamType,
				testFunc, testName, path, i,
				"error или string",
				result.Kind(),
				val,
			)
		}
	}
	return
}

func cmpErrs(tstResult, etaResult interface{}, fmtString *string, fmtArgs *[]interface{}) (hasDiff bool) {
	tstErrStr := getErrStr(tstResult)
	etaErrStr := getErrStr(etaResult)
	if tstErrStr != etaErrStr {
		hasDiff = true
		*fmtString += ansiErr
		*fmtArgs = append(*fmtArgs,
			tstErrStr,
			etaErrStr,
			tstErrStr,
			etaErrStr,
			bwjson.Pretty(tstResult),
			bwjson.Pretty(etaResult),
		)
	}
	return
}

func getErrStr(val interface{}) (result string) {
	switch t := val.(type) {
	case string:
		result = t
	case error:
		result = bwerr.FmtStringOf(t)
	}
	return
}

func cmpVals(tstResult, etaResult interface{}, fmtString *string, fmtArgs *[]interface{}) (hasDiff bool) {
	if cmp := pretty.Compare(tstResult, etaResult); len(cmp) > 0 {
		hasDiff = true
		*fmtString += ansiDiffBegin + colorizedCmp(cmp) + ansiDiffEnd + ansiVal
		*fmtArgs = append(*fmtArgs,
			tstResult,
			etaResult,
			bwjson.Pretty(tstResult),
			bwjson.Pretty(etaResult),
		)
	}
	return
}

func colorizedCmp(s string) string {
	ss := strings.Split(s, "\n")
	result := make([]string, 0, len(ss))
	ansiReset := ansi.Reset()
	ansiPlus := ansi.CSIFromSGRCodes(ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen})).String()
	ansiMinus := ansi.CSIFromSGRCodes(ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorRed})).String()
	for _, s := range ss {
		if len(s) > 0 {
			r := s[0]
			if r == '+' {
				s = ansiPlus + s + ansiReset
			} else if r == '-' {
				s = ansiMinus + s + ansiReset
			}
		}
		result = append(result, s)
	}
	return strings.Join(result, "\n")
}
