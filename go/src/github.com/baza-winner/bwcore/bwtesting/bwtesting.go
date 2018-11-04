// Предоставялет функции для тестирования.
package bwtesting

import (
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

type TestCaseStruct struct {
	In  []interface{}
	Out []interface{}
}

var (
	ansiSecondArgToBeFunc        string
	ansiTestPrefix               string
	ansiTestHeading              string
	ansiExpectsCountParams       string
	ansiExpectsParamType         string
	ansiExpectsOneReturnValue    string
	ansiExpectsTypeOfReturnValue string
	ansiPath                     string
	ansiTestTitleFunc            string
	ansiTestTitleOpenBrace       string
	ansiTestTitleSep             string
	ansiTestTitleVal             string
	ansiTestTitleCloseBrace      string
	ansiErr                      string
	ansiVal                      string
	ansiWhere                    string
)

func init() {
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
	ansiExpectsParamType = ansi.String(testPrefixFmt + ".%s.%d<ansi>: ожидается <ansiType>%s<ansi> вместо <ansiType>%s<ansi> (<ansiVal>%#v<ansi>)")
	ansiExpectsOneReturnValue = ansi.String(testPrefixFmt + ".%s.%d<ansi>: ожидается <ansiVal>1<ansi> возвращаемое значение вместо <ansiVal>%d")
	ansiExpectsTypeOfReturnValue = ansi.String(testPrefixFmt + ".%s.%d<ansi>: в качесте возвращаемого значения ожидается <ansiType>%s<ansi> вместо <ansiType>%s<ansi>")
	ansiPath = ansi.String("<ansiPath>.[%d]<ansi>:\n")
	ansiTestTitleFunc = ansi.String("<ansiFunc>%s")
	ansiTestTitleOpenBrace = ansi.String("(")
	ansiTestTitleSep = ansi.String(",")
	ansiTestTitleVal = ansi.String("<ansiVal>%#v")
	ansiTestTitleCloseBrace = ansi.String(")")
	ansiErr = ansi.String(
		"<ansiErr>tst err<ansi>: '%s'" +
			"\n<ansiOK>eta err<ansi>: '%s'" +
			"\n<ansiErr>tst(q)<ansi>: %q" +
			"\n<ansiOK>eta(q)<ansi>: %q" +
			"\n<ansiErr>tst(json)<ansi>: %s" +
			"\n<ansiOK>eta(json)<ansi>: %s",
	)
	ansiVal = ansi.String(
		"<ansiErr>tst(v)<ansi>: %#v" +
			"\n<ansiOK>eta(v)<ansi>: %#v" +
			"\n<ansiErr>tst(json)<ansi>: %s" +
			"\n<ansiOK>eta(json)<ansi>: %s",
	)
	ansiWhere = ansi.String("\n<ansiVar>Where<ansi>: <ansiPath>%s")
}

func BwRunTests(t *testing.T, f interface{}, tests map[string]TestCaseStruct) {
	if err := tryBwRunTests(t, f, tests, 1); err != nil {
		bwerr.PanicA(bwerr.E{Error: err})
	}
}

func tryBwRunTests(t *testing.T, f interface{}, tests map[string]TestCaseStruct, depth uint) error {
	fType := reflect.TypeOf(f)
	if fType.Kind() != reflect.Func {
		return bwerr.From("reflect.TypeOf(f).Kind(): " + fType.Kind().String() + "\n")
	}
	inDef := []reflect.Type{}
	numIn := fType.NumIn()

	w := where.MustFrom(depth + 1)
	testFunc := w.FuncName()
	testeeFunc := testFunc[4:]

	ansiTestTitle := ansiTestTitleFunc + ansiTestTitleOpenBrace
	// ansiTestTitle := ansi.Concat(ansiTestTitleFunc, ansiTestTitleOpenBrace) //"<ansiFunc>%s<ansi>("
	for i := 0; i < numIn; i++ {
		inType := fType.In(i)
		inDef = append(inDef, inType)
		if i > 0 {
			ansiTestTitle += ansiTestTitleSep
			// ansiTestTitle = ansi.Concat(ansiTestTitle, ansiTestTitleSep)
		}
		ansiTestTitle += ansiTestTitleVal
		// ansiTestTitle = ansi.Concat(ansiTestTitle, ansiTestTitleVal)
	}
	ansiTestTitle += ansiTestTitleCloseBrace
	// ansiTestTitle = ansi.Concat(ansiTestTitle, ansiTestTitleCloseBrace)
	outDef := []reflect.Type{}
	numOut := fType.NumOut()
	for i := 0; i < numOut; i++ {
		outType := fType.Out(i)
		outDef = append(outDef, outType)
	}
	fValue := reflect.ValueOf(f)

	for testName, test := range tests {
		t.Logf(ansiTestHeading, testName)
		if len(test.In) != numIn {
			return bwerr.FromA(bwerr.A{depth + 1,
				ansiExpectsCountParams,
				bw.Args(
					testFunc, testName, "In",
					numIn, bw.PluralWord(numIn, "параметр", "", "а", "ов"), len(test.In),
				),
			})
		}
		if len(test.Out) != numOut {
			return bwerr.FromA(bwerr.A{depth + 1,
				ansiExpectsCountParams,
				bw.Args(
					testFunc, testName, "Out",
					numOut, bw.PluralWord(numOut, "параметр", "", "а", "ов"), len(test.Out),
				),
			})
		}
		in := []reflect.Value{}
		for i := 0; i < numIn; i++ {
			var inItem reflect.Value
			if test.In[i] == nil {
				inItem = reflect.New(inDef[i]).Elem()
			} else {
				inItem = reflect.ValueOf(test.In[i])
			}
			if inDef[i].Kind() != reflect.Interface && inItem.Kind() != inDef[i].Kind() {
				return bwerr.FromA(bwerr.A{depth + 1,
					ansiExpectsParamType,
					bw.Args(
						testFunc, testName, "In",
						i,
						inDef[i].Kind(),
						inItem.Kind(),
						test.In[i],
					),
				})
			}
			if i == numIn-1 && fType.IsVariadic() {
				for j := 0; j < inItem.Len(); j++ {
					in = append(in, inItem.Index(j))
				}
			} else {
				in = append(in, inItem)
			}
		}
		out := fValue.Call(in)
		outEta := []interface{}{}
		for i := 0; i < numOut; i++ {
			v := test.Out[i]
			if v != nil {
				vType := reflect.TypeOf(v)
				if true &&
					vType.Kind() == reflect.Func &&
					true {
					if vType.NumOut() != 1 {
						return bwerr.FromA(bwerr.A{depth + 1,
							ansiExpectsOneReturnValue,
							bw.Args(
								testFunc, testName, "Out",
								vType.NumOut(),
							),
						})
					} else if vType.Out(0) != outDef[i] {
						return bwerr.FromA(bwerr.A{depth + 1,
							ansiExpectsTypeOfReturnValue,
							bw.Args(
								testFunc, testName, "Out",
								outDef[i],
								vType.Out(0),
							),
						})
					} else {
						vOut := reflect.ValueOf(v).Call([]reflect.Value{reflect.ValueOf(test)})
						outEta = append(outEta, vOut[0].Interface())
					}
					continue
				}
			}
			outEta = append(outEta, v)
		}

		for i := 0; i < numOut; i++ {
			fmtString := ansiTestTitle
			// testTitle := fmt.Sprintf(ansiTestTitle, test.In...)
			fmtArgs := append(bw.Args(testeeFunc), test.In...) //bw.Args(testeeFunc, test.In...)
			if numOut > 1 {
				fmtString += ansiPath
				fmtArgs = append(fmtArgs, i)
				// testTitle += fmt.Sprintf("<ansiPath>.[%d]<ansi>:\n", i)
			}
			var hasDiff bool
			// fmtString := testTitle
			// fmtArgs := []interface{}{}
			if outDef[i].Implements(reflect.TypeOf((*error)(nil)).Elem()) {
				etaErr, _ := outEta[i].(error)
				var tstErr error
				if out[i].IsNil() {
					tstErr = nil
				} else {
					tstErr, _ = out[i].Interface().(error)
				}
				tstErrStr := bwerr.FmtStringOf(tstErr)
				etaErrStr := bwerr.FmtStringOf(etaErr)
				if tstErrStr != etaErrStr {
					hasDiff = true
					fmtString += ansiErr
					fmtArgs = append(fmtArgs,
						tstErrStr,
						etaErrStr,
						tstErrStr,
						etaErrStr,
						bwjson.Pretty(tstErr),
						bwjson.Pretty(etaErr),
					)

					if tstErr != nil {
						tstErrType := reflect.TypeOf(tstErr)
						if tstErrType.Kind() == reflect.Struct {
							if sf, ok := tstErrType.FieldByName("Where"); ok && sf.Type.Kind() == reflect.String {
								fmtString += ansiWhere
								fmtArgs = append(fmtArgs, reflect.ValueOf(tstErr).FieldByName("Where").Interface())
							}
						}
					}
				}
			} else {
				tstResult := out[i].Interface()
				etaResult := outEta[i]
				if cmp := pretty.Compare(tstResult, etaResult); len(cmp) > 0 {
					hasDiff = true
					fmtString += colorizedCmp(cmp) + ansiVal
					fmtArgs = append(fmtArgs,
						tstResult,
						etaResult,
						bwjson.Pretty(tstResult),
						bwjson.Pretty(etaResult),
					)
				}
			}
			if hasDiff {
				t.Error(bw.Spew.Sprintf(fmtString, fmtArgs...))
			}
		}
	}
	return nil
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
	return strings.Join(result, "\n") + "\n"
}

// func getPluralWord(count int, word string, word1 string, word2_4 string, _word5more ...string) (result string) {
// 	var word5more string
// 	if _word5more != nil {
// 		word5more = _word5more[0]
// 	}
// 	if len(word5more) == 0 {
// 		word5more = word2_4
// 	}
// 	result = word5more
// 	decimal := count / 10 % 10
// 	if decimal != 1 {
// 		unit := count % 10
// 		if unit == 1 {
// 			result = word1
// 		} else if 2 <= unit && unit <= 4 {
// 			result = word2_4
// 		}
// 	}
// 	return word + result
// }
