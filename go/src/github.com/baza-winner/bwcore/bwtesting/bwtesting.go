/*
Предоставялет функции для тестирования.

Смотри defvalid_test.go в качестве образца использования
*/
package bwtesting

import (
	"fmt"
	"log"
	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwerror"
	"reflect"
	"testing"
)

func Debug(args ...interface{}) {
	for _, arg := range args {
		log.Printf("%s ", bwjson.PrettyJson(arg))
	}
	log.Println()
}

type GetEtaErr func (testIntf interface{}) (err error)

type BT interface {
	GetEtaErr() interface{}
	GetEtaResult() interface{}
	GetTitle() string
	GetTstResultErr() (interface{}, error)
}

func BtCheckErr(t *testing.T, bt BT, tstErr error) bool {
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
			t.Errorf(ansi.Ansi("", bt.GetTitle()+"    => err: <ansiErr>'%v'<ansi>\n, want err: <ansiOK>'%v'"), tstErr, etaErr)
			fmt.Printf("tstErr: %+q\netaErr: %+q\n", tstErr, etaErr)
		}
		return false
	}
	return true
}

func BtCheckResult(t *testing.T, bt BT, tstResult interface{}) {
	eta := bt.GetEtaResult()
	if !reflect.DeepEqual(tstResult, eta) { // https://stackoverflow.com/questions/18208394/testing-equivalence-of-maps-golang
		t.Errorf(ansi.Ansi("", bt.GetTitle()+"    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), bwjson.PrettyJson(tstResult), bwjson.PrettyJson(eta))
		fmt.Printf("tstResult: %+q\netaResult: %+q\n", tstResult, eta)
	}
}

func BtCheckErrResult(t *testing.T, bt BT, tstErr error, tstResult interface{}) {
	if BtCheckErr(t, bt, tstErr) {
		BtCheckResult(t, bt, tstResult)
	}
}

func BtRunTest(t *testing.T, testName string, bt BT) {
	t.Logf(ansi.Ansi(`Header`, "Running test case <ansiPrimaryLiteral>%s"), testName)
	result, err := bt.GetTstResultErr()
	BtCheckErrResult(t, bt, err, result)
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
