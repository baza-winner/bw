/*
Предоставялет функции для тестирования.

Смотри defvalid_test.go в качестве образца использования
*/
package bwtesting

import (
	"fmt"
	// "log"
	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	// "github.com/baza-winner/bwcore/bwstring"
	"reflect"
	"runtime"
	"strings"
	"testing"
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
		testPrefix := "<ansiCmd>" + testFunc + "::<ansiOutline>test" + "<ansiCmd>.<ansiSecondaryLiteral>" + testName + "<ansiCmd>"
		t.Logf(ansi.Ansi(`Header`, "Running test case <ansiPrimaryLiteral>%s"), testName)
		if len(test.In) != numIn {
			bwerror.Panic(testPrefix+".Out<ansi>: ожидается <ansiPrimaryLiteral>%d<ansi> %s вместо <ansiSecondaryLiteral>%d", numIn, getPluralWord(numIn, "параметр", "", "а", "ов"), len(test.In))
		}
		if len(test.Out) != numOut {
			bwerror.Panic(testPrefix+".Out<ansi>: ожидается <ansiPrimaryLiteral>%d<ansi> %s вместо <ansiSecondaryLiteral>%d", numOut, getPluralWord(numOut, "параметр", "", "а", "ов"), len(test.Out))
		}
		in := []reflect.Value{}
		for i := 0; i < numIn; i++ {
			inItem := reflect.ValueOf(test.In[i])
			if inDef[i].Kind() != reflect.Interface && inItem.Kind() != inDef[i].Kind() {
				bwerror.Panic(
					testPrefix+".In[%d]<ansi>: ожидается <ansiPrimaryLiteral>%s<ansi> вместо <ansiPrimaryLiteral>%s<ansi> (<ansiSecondaryLiteral>%#v<ansi>)",
					i,
					inDef[i].Kind(),
					inItem.Kind(),
					test.In[i],
				)
			}
			if i < numIn - 1 || !fType.IsVariadic() {
				in = append(in, inItem)
			} else {
				for j := 0; j < inItem.Len(); j++ {
					in = append(in, inItem.Index(j))
				}
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
					vType.NumIn() == 1 &&
					vType.In(0).Kind() == reflect.Struct &&
					vType.In(0).Name() == "TestStruct" &&
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
							"    => <ansiErr>err<ansi>: '%v'\n" +
							", <ansiOK>want err<ansi>: '%v'\n" +
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
						if reflect.TypeOf(tstErr).Kind() == reflect.Struct {
							if sf, ok := reflect.TypeOf(tstErr).FieldByName("Where"); ok && sf.Type.Kind() == reflect.String {
								fmtString += "\n" +
									"<ansiHeader>errWhere<ansi>: <ansiCmd>%s"
								fmtArgs = append(fmtArgs, reflect.ValueOf(tstErr).FieldByName("Where").Interface())
							}
						}
						t.Errorf(ansi.Ansi("", fmtString), fmtArgs...)
					}
				}
			} else if !reflect.DeepEqual(out[i].Interface(), test.Out[i]) {
				tstResult := out[i].Interface()
				etaResult := test.Out[i]
				fmtString := testTitle + "    <ansiErr>=><ansi> %#v\n, <ansiOK>want<ansi> %#v"
				fmtArgs := []interface{}{tstResult, etaResult}

				fmtString += "\n" +
					"<ansiErr>tst(json)<ansi>: %s"
				if jsonable, ok := tstResult.(bwjson.Jsonable); ok {
					fmtArgs = append(fmtArgs, bwjson.PrettyJsonOf(jsonable))
				} else {
					fmtArgs = append(fmtArgs, bwjson.PrettyJson(tstResult))
				}

				fmtString += "\n" +
					"<ansiOK>eta(json)<ansi>: %s"
				if jsonable, ok := etaResult.(bwjson.Jsonable); ok {
					fmtArgs = append(fmtArgs, bwjson.PrettyJsonOf(jsonable))
				} else {
					fmtArgs = append(fmtArgs, bwjson.PrettyJson(etaResult))
				}

				t.Errorf(ansi.Ansi("", fmtString), fmtArgs...)
			}
		}
	}
}

func getPluralWord(count int, word string, word1 string, word2_4 string, _word5more ...string) (result string) {
	var word5more string
	if _word5more != nil { word5more = _word5more[0] }
	if len(word5more) == 0 { word5more = word2_4 }
	result = word5more
	decimal := count / 10 % 10
	if decimal != 1 {
		unit := count % 10
		if unit == 1 {
			result = word1
		} else if 2 <= unit && unit <= 4  {
			result = word2_4
		}
	}
	return word + result
}
// type WhenDiffErrFunc func(t *testing.T, bt BT, tstErr, etaErr error)

// func WhenDiffErrFuncDefault(t *testing.T, bt BT, tstErr, etaErr error) {
// 	fmtString := bt.GetTitle()+
// 	"    => err: <ansiErr>'%v'<ansi>\n"+
// 	", want err: <ansiOK>'%v'\n"+
// 	"tstErr(go): <ansiErr>%#v<ansi>\n"+
// 	"etaErr(go): <ansiOK>%#v<ansi>"
// 	fmtArgs := []interface{}{ tstErr, etaErr, tstErr, etaErr }
// 	// if jsonable, ok := tstErr.(bwjson.Jsonable); ok {
// 	// 	fmtString+="\n"+
// 	// 	"tstErr(json): <ansiErr>%s<ansi>"
// 	// 	fmtArgs=append(fmtArgs, bwjson.PrettyJsonOf(jsonable))
// 	// }
// 	if jsonable, ok := etaErr.(bwjson.Jsonable); ok {
// 		fmtString+="\n"+
// 		"etaErr(json): <ansiOK>%s<ansi>"
// 		fmtArgs=append(fmtArgs, bwjson.PrettyJsonOf(jsonable))
// 	}
// 	if reflect.TypeOf(tstErr).Kind() == reflect.Struct {
// 		if sf, ok := reflect.TypeOf(tstErr).FieldByName("Where"); ok && sf.Type.Kind() == reflect.String {
// 			fmtString+="\n"+
// 			"errWhere: <ansiCmd>%s"
// 			fmtArgs=append(fmtArgs, reflect.ValueOf(tstErr).FieldByName("Where").Interface())
// 		}
// 	}
// 	t.Errorf(ansi.Ansi("", fmtString), fmtArgs...)
// }

// type WhenDiffResultFunc func(t *testing.T, bt BT, tstResult, etaResult interface{})

// func WhenDiffResultFuncDefault(t *testing.T, bt BT, tstResult, etaResult interface{}) {
// 	t.Errorf(
// 		ansi.Ansi("", bt.GetTitle()+"    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"),
// 		bwjson.PrettyJson(bt.GetResultDataForJson(tstResult)),
// 		bwjson.PrettyJson(bt.GetResultDataForJson(etaResult)),
// 	)
// 	// t.Errorf(ansi.Ansi("", bt.GetTitle()+"    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), bwjson.PrettyJson(tstResult), bwjson.PrettyJson(etaResult))
// 	fmt.Printf("tstResult: %#v\netaResult: %#v\n", tstResult, etaResult)
// }

// type GetEtaErr func (testIntf interface{}) (err error)

// type BT interface {
// 	IsDiffResult(tstResult, etaResult interface{}) bool
// 	GetTitle() string
// 	GetTstResultErr() (interface{}, error)
// 	GetResultDataForJson(result interface{}) interface{}
// }

// func GetResultDataForJsonDefault(result bwjson.Jsonable) interface{} {
// 	return result.GetDataForJson()
// }

// func IsDiffResultDefault(tstResult, etaResult interface{}) bool {
// 	return !reflect.DeepEqual(tstResult, etaResult)
// }

// func BtCheckErr(t *testing.T, bt BT, tstErr error, opts ...map[BtOptFuncType]interface{}) bool {
// 	if reflect.TypeOf(bt).Kind() != reflect.Struct {
// 		bwerror.Panic("reflect.TypeOf(bt).Kind(): %s\n", reflect.TypeOf(bt).Kind())
// 	}
// 	if _, ok := reflect.TypeOf(bt).FieldByName("Err"); !ok {
// 		bwerror.Panic("reflect.TypeOf(bt).FieldByName(%q) not found\n", "Err")
// 	}
// 	etaIntf := reflect.ValueOf(bt).FieldByName("Err").Interface()
// 	var etaErr error
// 	var ok bool
// 	if etaIntf != nil {
// 		if etaErr, ok = etaIntf.(error); !ok {
// 			if getEtaErr, ok := etaIntf.(func (interface{}) error); !ok {
// 				bwerror.Panic("<ansiOutline>bt.Err<ansi> (<ansiSecondaryLiteral>%#v<ansi>) expected to be either <ansiPrimaryLiteral>error<ansi> or <ansiPrimaryLiteral>bwerror.GetEtaErr", etaIntf)
// 			} else {
// 				etaErr = GetEtaErr(getEtaErr)(bt)
// 			}
// 		}
// 	}
// 	if tstErr != etaErr {
// 		if tstErr == nil || etaErr == nil || tstErr.Error() != etaErr.Error() {
// 			if opts != nil {
// 				if i, ok := opts[0][WhenDiffErr]; ok {
// 					if f, ok := i.(func(t *testing.T, bt BT, tstErr, etaErr error)); !ok {
// 						bwerror.Panic("opt[%s] (%#v) is not of <ansiPrimaryLiteral>WhenDiffErrFunc", WhenDiffErr, i)
// 					} else {
// 						f(t, bt, tstErr, etaErr)
// 						return false
// 					}
// 				}
// 			}
// 			WhenDiffErrFuncDefault(t, bt, tstErr, etaErr)
// 		}
// 		return false
// 	}
// 	return true
// }

// func BtCheckResult(t *testing.T, bt BT, tstResult interface{}, opts ...map[BtOptFuncType]interface{}) {
// 	// etaResult := bt.GetEtaResult()
// 	if reflect.TypeOf(bt).Kind() != reflect.Struct {
// 		bwerror.Panic("reflect.TypeOf(bt).Kind(): %s\n", reflect.TypeOf(bt).Kind())
// 	}
// 	if _, ok := reflect.TypeOf(bt).FieldByName("Result"); !ok {
// 		bwerror.Panic("reflect.TypeOf(bt).FieldByName(%q) not found\n", "Result")
// 	}
// 	etaResult := reflect.ValueOf(bt).FieldByName("Result").Interface()
// 	// if !reflect.DeepEqual(tstResult, etaResult) { // https://stackoverflow.com/questions/18208394/testing-equivalence-of-maps-golang
// 	if bt.IsDiffResult(tstResult, etaResult) { // https://stackoverflow.com/questions/18208394/testing-equivalence-of-maps-golang
// 		if opts != nil {
// 			if i, ok := opts[0][WhenDiffResult]; ok {
// 				if f, ok := i.(*func(t *testing.T, bt BT, tstResult, etaResult interface{})); !ok {
// 					bwerror.Panic("opt[%s] (%#v) is not of <ansiPrimaryLiteral>WhenDiffResultFunc", WhenDiffResult, i)
// 				} else {
// 					(*f)(t, bt, tstResult, etaResult)
// 					return
// 				}
// 			}
// 		}
// 		WhenDiffResultFuncDefault(t, bt, tstResult, etaResult)
// 	}

// 	// if !bt.CompareResult(tstResult, etaResult) { // https://stackoverflow.com/questions/18208394/testing-equivalence-of-maps-golang
// 	// 	// t.Errorf(ansi.Ansi("", bt.GetTitle()+"    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), bwjson.PrettyJson(tstResult), bwjson.PrettyJson(etaResult))
// 	// 	t.Errorf(ansi.Ansi("", bt.GetTitle()+"    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), bt.GetResultDataForJson(tstResult), bt.GetResultDataForJson(etaResult))
// 	// 	// fmt.Printf("tstResult: %#v\netaResult: %#v\n", tstResult, etaResult)
// 	// }
// }

// func BtCheckErrResult(t *testing.T, bt BT, tstErr error, tstResult interface{}, opts ...map[BtOptFuncType]interface{}) {
// 	if BtCheckErr(t, bt, tstErr, opts...) {
// 		BtCheckResult(t, bt, tstResult, opts...)
// 	}
// }

// func BtRunTest(t *testing.T, testName string, bt BT, opts ...map[BtOptFuncType]interface{}) {
// 	t.Logf(ansi.Ansi(`Header`, "Running test case <ansiPrimaryLiteral>%s"), testName)
// 	result, err := bt.GetTstResultErr()
// 	BtCheckErrResult(t, bt, err, result, opts...)
// }

// CompareErrors - сравнивает тестовое и эталонное значение ошибки, и в случае их расхождения вызывает t.Errorf
// Deprecated: test should implement BT instead
// func CompareErrors(t *testing.T, tst, eta error, testTitle string) bool {
// 	if tst != eta {
// 		if tst == nil || eta == nil || tst.Error() != eta.Error() {
// 			t.Errorf(ansi.Ansi("", testTitle+"    => err: <ansiErr>'%v'<ansi>\n, want err: <ansiOK>'%v'"), tst, eta)
// 			fmt.Printf("tst: %+q\neta: %+q\n", tst, eta)
// 		}
// 		return false
// 	}
// 	return true
// }

// // DeepEqual - сравнивает тестовое и эталонное значение, и в случае их расхождения вызывает t.Errorf
// // Deprecated: test should implement BT instead
// func DeepEqual(t *testing.T, tst, eta interface{}, testTitle string) {
// 	if !reflect.DeepEqual(tst, eta) { // https://stackoverflow.com/questions/18208394/testing-equivalence-of-maps-golang
// 		t.Errorf(ansi.Ansi("", testTitle+"    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), bwjson.PrettyJson(tst), bwjson.PrettyJson(eta))
// 		fmt.Printf("tst: %+q\neta: %+q\n", tst, eta)
// 	}
// }

// // CheckTestErrResult - сначала сравнивает ошибки (CompareErrors), потом значения (DeepEqual)
// // Deprecated: test should implement BT instead
// func CheckTestErrResult(t *testing.T, tstErr, etaErr error, tstVal, etaVal interface{}, testTitle string) {
// 	if CompareErrors(t, tstErr, etaErr, testTitle) {
// 		DeepEqual(t, tstVal, etaVal, testTitle)
// 	}
// }
