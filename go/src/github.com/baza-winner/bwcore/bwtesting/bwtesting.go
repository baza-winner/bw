/*
Предоставялет функции для тестирования.
*/
package bwtesting

import (
	"fmt"
	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwjson"
	"reflect"
	"testing"
)

/*
CompareErrors - сравнивает тестовое и эталонное значение ошибки, и в случае их расхождения вызывает t.Errorf
*/
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

/*
DeepEqual - сравнивает тестовое и эталонное значение, и в случае их расхождения вызывает t.Errorf
*/
func DeepEqual(t *testing.T, tst, eta interface{}, testTitle string) {
	if !reflect.DeepEqual(tst, eta) { // https://stackoverflow.com/questions/18208394/testing-equivalence-of-maps-golang
		t.Errorf(ansi.Ansi("", testTitle+"    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), bwjson.PrettyJson(tst), bwjson.PrettyJson(eta))
		fmt.Printf("tst: %+q\neta: %+q\n", tst, eta)
	}
}

/*
DeepEqual - сначала сравнивает ошибки (CompareErrors), потом значения (DeepEqual)
*/
func CheckTestErrResult(t *testing.T, tstErr, etaErr error, tstVal, etaVal interface{}, testTitle string) {
	if CompareErrors(t, tstErr, etaErr, testTitle) {
		DeepEqual(t, tstVal, etaVal, testTitle)
	}
}
