package bwtesting

import (
	"fmt"
	"github.com/baza-winner/bw/ansi"
	"github.com/baza-winner/bw/core"
	"reflect"
	"testing"
)

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

func DeepEqual(t *testing.T, tst, eta interface{}, testTitle string) {
	if !reflect.DeepEqual(tst, eta) { // https://stackoverflow.com/questions/18208394/testing-equivalence-of-maps-golang
		t.Errorf(ansi.Ansi("", testTitle+"    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), core.PrettyJson(tst), core.PrettyJson(eta))
    fmt.Printf("tst: %+q\neta: %+q\n", tst, eta)
	}
}

func CheckTestErrResult(t *testing.T, tstErr, etaErr error, tstVal, etaVal interface{}, testTitle string) {
  if CompareErrors(t, tstErr, etaErr, testTitle) {
    DeepEqual(t, tstVal, etaVal, testTitle)
  }
}

