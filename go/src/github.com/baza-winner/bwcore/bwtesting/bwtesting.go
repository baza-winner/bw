// Предоставялет функции для тестирования.
package bwtesting

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/kylelemons/godebug/pretty"
	// "log"
)

type TestCaseStruct struct {
	In  []interface{}
	Out []interface{}
}

func BwRunTests(t *testing.T, tests map[string]TestCaseStruct, f interface{}) {
	fType := reflect.TypeOf(f)
	if fType.Kind() != reflect.Func {
		bwerror.Panic("reflect.TypeOf(f).Kind(): %s\n", fType.Kind())
	}
	inDef := []reflect.Type{}
	numIn := fType.NumIn()

	function, _, _, _ := runtime.Caller(1)
	testFunc := runtime.FuncForPC(function).Name()
	pointPos := strings.LastIndex(testFunc, ".") + 1
	if pointPos > 0 {
		testFunc = testFunc[pointPos:]
	}
	testeeFunc := testFunc[4:]

	fmtTestTitle := testeeFunc + "("
	for i := 0; i < numIn; i++ {
		inType := fType.In(i)
		inDef = append(inDef, inType)
		if i > 0 {
			fmtTestTitle += ", "
		}
		fmtTestTitle += "%#v"
	}
	fmtTestTitle += ")"
	outDef := []reflect.Type{}
	numOut := fType.NumOut()
	for i := 0; i < numOut; i++ {
		outType := fType.Out(i)
		outDef = append(outDef, outType)
	}
	fValue := reflect.ValueOf(f)

	for testName, test := range tests {
		testPrefix := "<ansiCmd>" + testFunc + "::<ansiOutline>test" + "<ansiCmd>.<ansiSecondaryLiteral>[\"" + testName + "\"]<ansiCmd>"
		t.Logf(ansi.Ansi(`Header`, "Running test case <ansiPrimaryLiteral>%s"), testName)
		if len(test.In) != numIn {
			bwerror.Panic(testPrefix+".In<ansi>: ожидается <ansiPrimaryLiteral>%d<ansi> %s вместо <ansiSecondaryLiteral>%d", numIn, getPluralWord(numIn, "параметр", "", "а", "ов"), len(test.In))
		}
		if len(test.Out) != numOut {
			bwerror.Panic(testPrefix+".Out<ansi>: ожидается <ansiPrimaryLiteral>%d<ansi> %s вместо <ansiSecondaryLiteral>%d", numOut, getPluralWord(numOut, "параметр", "", "а", "ов"), len(test.Out))
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
				bwerror.Panic(
					testPrefix+".In[%d]<ansi>: ожидается <ansiPrimaryLiteral>%s<ansi> вместо <ansiPrimaryLiteral>%s<ansi> (<ansiSecondaryLiteral>%#v<ansi>)",
					i,
					inDef[i].Kind(),
					inItem.Kind(),
					test.In[i],
				)
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
						bwerror.Panic(
							testPrefix+".Out[%d]<ansi>: ожидается <ansiPrimaryLiteral>1<ansi> одно возвращаемое значение вместо <ansiSecondaryLiteral>%d",
							vType.NumOut(),
						)
					} else if vType.Out(0) != outDef[i] {
						bwerror.Panic(
							testPrefix+".Out[%d]<ansi>: в качесте возвращаемого значения ожидается <ansiPrimaryLiteral>%s<ansi> вместо <ansiPrimaryLiteral>%s<ansi>",
							outDef[i],
							vType.Out(0),
						)
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
			testTitle := fmt.Sprintf(fmtTestTitle, test.In...)
			if numOut > 1 {
				testTitle += fmt.Sprintf("<ansiCmd>.[%d]<ansi>:\n", i)
			}
			if outDef[i].Implements(reflect.TypeOf((*error)(nil)).Elem()) {
				etaErr, _ := outEta[i].(error)
				var tstErr error
				if out[i].IsNil() {
					tstErr = nil
				} else {
					tstErr, _ = out[i].Interface().(error)
				}
				if tstErr != etaErr {
					if tstErr == nil || etaErr == nil || tstErr.Error() != etaErr.Error() {
						fmtString := testTitle +
							"    => <ansiErr>err<ansi>: '%s'\n" +
							", <ansiOK>want err<ansi>: '%s'\n" +
							"<ansiErr>tstErr(q)<ansi>: %q\n" +
							"<ansiOK>etaErr(q)<ansi>: %q"
						fmtArgs := []interface{}{tstErr, etaErr, tstErr, etaErr}

						if jsonable, ok := tstErr.(bwjson.Jsonable); ok {
							fmtString += "\n" +
								"<ansiErr>tstErr(json)<ansi>: %s"
							fmtArgs = append(fmtArgs, bwjson.PrettyJsonOf(jsonable))
						}
						if jsonable, ok := etaErr.(bwjson.Jsonable); ok {
							fmtString += "\n" +
								"<ansiOK>etaErr(json)<ansi>: %s"
							fmtArgs = append(fmtArgs, bwjson.PrettyJsonOf(jsonable))
						}

						if tstErr != nil {
							tstErrType := reflect.TypeOf(tstErr)
							if tstErrType.Kind() == reflect.Struct {
								if sf, ok := tstErrType.FieldByName("Where"); ok && sf.Type.Kind() == reflect.String {
									fmtString += "\n" +
										"<ansiHeader>errWhere<ansi>: <ansiCmd>%s"
									fmtArgs = append(fmtArgs, reflect.ValueOf(tstErr).FieldByName("Where").Interface())
								}
							}
						}
						t.Error(bwerror.Spew.Sprintf(ansi.Ansi("", fmtString), fmtArgs...))
					}
				}
			} else {
				tstResult := out[i].Interface()
				etaResult := outEta[i]
				if cmp := pretty.Compare(tstResult, etaResult); len(cmp) > 0 {

					fmtString := testTitle
					fmtString += cmp
					fmtString += "\n <ansiErr>got<ansi>: %#v\n<ansiOK>want<ansi>: %#v\n"
					fmtArgs := []interface{}{tstResult, etaResult}

					// fmtString += "\n" +
					// 	"<ansiErr>tst(json)<ansi>: %s"
					// if jsonable, ok := tstResult.(bwjson.Jsonable); ok {
					// 	fmtArgs = append(fmtArgs, bwjson.PrettyJsonOf(jsonable))
					// } else {
					// 	fmtArgs = append(fmtArgs, bwjson.PrettyJson(tstResult))
					// }

					// fmtString += "\n" +
					// 	"<ansiOK>eta(json)<ansi>: %s"
					// if jsonable, ok := etaResult.(bwjson.Jsonable); ok {
					// 	fmtArgs = append(fmtArgs, bwjson.PrettyJsonOf(jsonable))
					// } else {
					// 	fmtArgs = append(fmtArgs, bwjson.PrettyJson(etaResult))
					// }

					t.Error(bwerror.Spew.Sprintf(ansi.Ansi("", fmtString), fmtArgs...))
				}
			}
		}
	}
}

func getPluralWord(count int, word string, word1 string, word2_4 string, _word5more ...string) (result string) {
	var word5more string
	if _word5more != nil {
		word5more = _word5more[0]
	}
	if len(word5more) == 0 {
		word5more = word2_4
	}
	result = word5more
	decimal := count / 10 % 10
	if decimal != 1 {
		unit := count % 10
		if unit == 1 {
			result = word1
		} else if 2 <= unit && unit <= 4 {
			result = word2_4
		}
	}
	return word + result
}
