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
	// ansiSeparator                string
	ansiPath                string
	ansiTestTitleFunc       string
	ansiTestTitleOpenBrace  string
	ansiTestTitleSep        string
	ansiTestTitleVal        string
	ansiTestTitleCloseBrace string
	ansiErr                 string
	ansiVal                 string
	ansiDiffBegin           string
	ansiDiffEnd             string
	// ansiSeparator                string
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
	// ansiSeparator = ":\n"
	ansiPath = ansi.String("<ansiPath>.[%d]<ansi>")
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
	// ansiSeparator = "\n------------------------\n"
}

func BwRunTests(t *testing.T, f interface{}, tests map[string]TestCaseStruct) {
	if err := tryBwRunTests(t, f, tests, 1); err != nil {
		bwerr.PanicA(bwerr.Err(err))
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
	for i := 0; i < numIn; i++ {
		inType := fType.In(i)
		inDef = append(inDef, inType)
		if i > 0 {
			ansiTestTitle += ansiTestTitleSep
		}
		ansiTestTitle += ansiTestTitleVal
	}
	ansiTestTitle += ansiTestTitleCloseBrace
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
			v := test.In[i]
			if v == nil {
				inItem = reflect.New(inDef[i]).Elem()
			} else {
				inItem = reflect.ValueOf(v)
				// vType := reflect.TypeOf(v)
				if inItem.Kind() == reflect.Func {
					if vType := reflect.TypeOf(v); vType.NumOut() != 1 {
						return bwerr.FromA(bwerr.A{depth + 1,
							ansiExpectsOneReturnValue,
							bw.Args(
								testFunc, testName, "In",
								vType.NumOut(),
							),
						})
					} else if vType.Out(0) != inDef[i] {
						return bwerr.FromA(bwerr.A{depth + 1,
							ansiExpectsTypeOfReturnValue,
							bw.Args(
								testFunc, testName, "In", i,
								outDef[i],
								vType.Out(0),
							),
						})
					} else if vType.NumIn() == 0 {
						inItem = reflect.ValueOf(v).Call([]reflect.Value{})[0]
						// outEta = append(outEta, vOut[0].Interface())
					} else if vType.NumIn() == 1 {
						if vType.In(0).Kind() == reflect.String {
							inItem = reflect.ValueOf(v).Call([]reflect.Value{
								reflect.ValueOf(testName),
							})[0]
						} else {
							bwerr.TODO()
						}
					} else {
						bwerr.TODO()
					}
					// outEta = append(outEta, vOut[0].Interface())
					// }
					// continue
				}

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
				if vType := reflect.TypeOf(v); vType.Kind() == reflect.Func {
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
								testFunc, testName, "Out", i,
								outDef[i],
								vType.Out(0),
							),
						})
					} else if vType.NumIn() == 0 {
						outItem := reflect.ValueOf(v).Call([]reflect.Value{})[0]
						outEta = append(outEta, outItem.Interface())
					} else if vType.NumIn() == 1 {
						if vType.Out(0).Kind() == reflect.String {
							outItem := reflect.ValueOf(v).Call([]reflect.Value{
								reflect.ValueOf(testName),
							})[0]
							outEta = append(outEta, outItem.Interface())
						} else if vType.In(0).Name() == "TestCaseStruct" {
							outItem := reflect.ValueOf(v).Call([]reflect.Value{
								reflect.ValueOf(test),
							})[0]
							outEta = append(outEta, outItem.Interface())
						} else {
							bwerr.Panic("vType.In(0).Name(): %s", vType.In(0).Name())
						}
					} else {
						bwerr.TODO()
						outItem := reflect.ValueOf(v).Call([]reflect.Value{
							reflect.ValueOf(testName),
							reflect.ValueOf(test),
						})[0]
						outEta = append(outEta, outItem.Interface())
					}
					continue
				}

			}
			outEta = append(outEta, v)
		}

		for i := 0; i < numOut; i++ {
			fmtString := ansiTestTitle
			fmtArgs := append(bw.Args(testeeFunc), test.In...)
			if numOut > 1 {
				fmtString += ansiPath
				fmtArgs = append(fmtArgs, i)
			}
			fmtString += ":\n"
			var hasDiff bool
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
				}
			} else {
				tstResult := out[i].Interface()
				etaResult := outEta[i]
				if cmp := pretty.Compare(tstResult, etaResult); len(cmp) > 0 {
					hasDiff = true
					fmtString += ansiDiffBegin + colorizedCmp(cmp) + ansiDiffEnd + ansiVal
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
	return strings.Join(result, "\n") //+ "\n"
}
