/*
Предоставялет функции для тестирования.

Смотри defvalid_test.go в качестве образца использования
*/
package bwtesting

import (
	"fmt"
	// "log"
	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwerror"
	"reflect"
	"testing"
)

type BtOptFuncType uint16

const (
	BtOptFuncType_below_ BtOptFuncType = iota
	WhenDiffErr
	WhenDiffResult
	BtOptFuncType_above_
)

//go:generate stringer -type=BtOptFuncType

type WhenDiffErrFunc func(t *testing.T, bt BT, tstErr, etaErr error)

func WhenDiffErrFuncDefault(t *testing.T, bt BT, tstErr, etaErr error) {
	t.Errorf(ansi.Ansi("", bt.GetTitle()+"    => err: <ansiErr>'%v'<ansi>\n, want err: <ansiOK>'%v'"), tstErr, etaErr)
	fmt.Printf("tstErr: %+q\netaErr: %+q\n", tstErr, etaErr)
	if whereErr, ok := tstErr.(bwerror.WhereError); ok {
		fmt.Printf("errWhere: %s\n", whereErr.WhereError())
	}
}

type WhenDiffResultFunc func(t *testing.T, bt BT, tstResult, etaResult interface{})

func WhenDiffResultFuncDefault(t *testing.T, bt BT, tstResult, etaResult interface{}) {
	t.Errorf(
		ansi.Ansi("", bt.GetTitle()+"    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"),
		bwjson.PrettyJson(bt.GetResultDataForJson(tstResult)),
		bwjson.PrettyJson(bt.GetResultDataForJson(etaResult)),
	)
	// t.Errorf(ansi.Ansi("", bt.GetTitle()+"    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), bwjson.PrettyJson(tstResult), bwjson.PrettyJson(etaResult))
	fmt.Printf("tstResult: %#v\netaResult: %#v\n", tstResult, etaResult)
}

type GetEtaErr func (testIntf interface{}) (err error)

type BT interface {
	GetEtaErr() interface{}
	GetEtaResult() interface{}
	IsDiffResult(tstResult, etaResult interface{}) bool
	GetTitle() string
	GetTstResultErr() (interface{}, error)
	GetResultDataForJson(result interface{}) interface{}
}

func GetResultDataForJsonDefault(result bwjson.Jsonable) interface{} {
	return result.GetDataForJson()
}

func IsDiffResultDefault(tstResult, etaResult interface{}) bool {
	return !reflect.DeepEqual(tstResult, etaResult)
}

func BtCheckErr(t *testing.T, bt BT, tstErr error, opts ...map[BtOptFuncType]interface{}) bool {
	eta := bt.GetEtaErr()
	var etaErr error
	var ok bool
	if eta != nil {
		if etaErr, ok = eta.(error); !ok {
			if getEtaErr, ok := eta.(func (interface{}) error); !ok {
				bwerror.Panic("<ansiSecondaryLiteral>%#v<ansi> expected to be either <ansiPrimaryLiteral>error<ansi> or <ansiPrimaryLiteral>bwerror.GetEtaErr", eta)
			} else {
				etaErr = GetEtaErr(getEtaErr)(bt)
			}
		}
	}
	if tstErr != etaErr {
		if tstErr == nil || etaErr == nil || tstErr.Error() != etaErr.Error() {
			if opts != nil {
				if i, ok := opts[0][WhenDiffErr]; ok {
					if f, ok := i.(func(t *testing.T, bt BT, tstErr, etaErr error)); !ok {
						bwerror.Panic("opt[%s] (%#v) is not of <ansiPrimaryLiteral>WhenDiffErrFunc", WhenDiffErr, i)
					} else {
						f(t, bt, tstErr, etaErr)
						return false
					}
				}
			}
			WhenDiffErrFuncDefault(t, bt, tstErr, etaErr)
		}
		return false
	}
	return true
}

func BtCheckResult(t *testing.T, bt BT, tstResult interface{}, opts ...map[BtOptFuncType]interface{}) {
	etaResult := bt.GetEtaResult()
	// if !reflect.DeepEqual(tstResult, etaResult) { // https://stackoverflow.com/questions/18208394/testing-equivalence-of-maps-golang
	if bt.IsDiffResult(tstResult, etaResult) { // https://stackoverflow.com/questions/18208394/testing-equivalence-of-maps-golang
		if opts != nil {
			if i, ok := opts[0][WhenDiffResult]; ok {
				if f, ok := i.(*func(t *testing.T, bt BT, tstResult, etaResult interface{})); !ok {
					bwerror.Panic("opt[%s] (%#v) is not of <ansiPrimaryLiteral>WhenDiffResultFunc", WhenDiffResult, i)
				} else {
					(*f)(t, bt, tstResult, etaResult)
					return
				}
			}
		}
		WhenDiffResultFuncDefault(t, bt, tstResult, etaResult)
	}

	// if !bt.CompareResult(tstResult, etaResult) { // https://stackoverflow.com/questions/18208394/testing-equivalence-of-maps-golang
	// 	// t.Errorf(ansi.Ansi("", bt.GetTitle()+"    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), bwjson.PrettyJson(tstResult), bwjson.PrettyJson(etaResult))
	// 	t.Errorf(ansi.Ansi("", bt.GetTitle()+"    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), bt.GetResultDataForJson(tstResult), bt.GetResultDataForJson(etaResult))
	// 	// fmt.Printf("tstResult: %#v\netaResult: %#v\n", tstResult, etaResult)
	// }
}

func BtCheckErrResult(t *testing.T, bt BT, tstErr error, tstResult interface{}, opts ...map[BtOptFuncType]interface{}) {
	if BtCheckErr(t, bt, tstErr, opts...) {
		BtCheckResult(t, bt, tstResult, opts...)
	}
}

func BtRunTest(t *testing.T, testName string, bt BT, opts ...map[BtOptFuncType]interface{}) {
	t.Logf(ansi.Ansi(`Header`, "Running test case <ansiPrimaryLiteral>%s"), testName)
	result, err := bt.GetTstResultErr()
	BtCheckErrResult(t, bt, err, result, opts...)
}

// CompareErrors - сравнивает тестовое и эталонное значение ошибки, и в случае их расхождения вызывает t.Errorf
// Deprecated: test should implement BT instead
func CompareErrors(t *testing.T, tst, eta error, testTitle string) bool {
	if tst != eta {
		if tst == nil || eta == nil || tst.Error() != eta.Error() {
			t.Errorf(ansi.Ansi("", testTitle+"    => err: <ansiErr>'%v'<ansi>\n, want err: <ansiOK>'%v'"), tst, eta)
			fmt.Printf("tst: %+q\neta: %+q\n", tst, eta)
		}
		return false
	}
	return true
}

// DeepEqual - сравнивает тестовое и эталонное значение, и в случае их расхождения вызывает t.Errorf
// Deprecated: test should implement BT instead
func DeepEqual(t *testing.T, tst, eta interface{}, testTitle string) {
	if !reflect.DeepEqual(tst, eta) { // https://stackoverflow.com/questions/18208394/testing-equivalence-of-maps-golang
		t.Errorf(ansi.Ansi("", testTitle+"    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), bwjson.PrettyJson(tst), bwjson.PrettyJson(eta))
		fmt.Printf("tst: %+q\neta: %+q\n", tst, eta)
	}
}

// CheckTestErrResult - сначала сравнивает ошибки (CompareErrors), потом значения (DeepEqual)
// Deprecated: test should implement BT instead
func CheckTestErrResult(t *testing.T, tstErr, etaErr error, tstVal, etaVal interface{}, testTitle string) {
	if CompareErrors(t, tstErr, etaErr, testTitle) {
		DeepEqual(t, tstVal, etaVal, testTitle)
	}
}
