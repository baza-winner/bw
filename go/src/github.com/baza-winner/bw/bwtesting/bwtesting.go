package bwtesting

import (
	"fmt"
	"github.com/baza-winner/bw/ansi"
	"github.com/baza-winner/bw/bwerror"
	"github.com/baza-winner/bw/bwjson"
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
		t.Errorf(ansi.Ansi("", testTitle+"    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), bwjson.PrettyJson(tst), bwjson.PrettyJson(eta))
		fmt.Printf("tst: %+q\neta: %+q\n", tst, eta)
	}
}

func CheckTestErrResult(t *testing.T, tstErr, etaErr error, tstVal, etaVal interface{}, testTitle string) {
	if CompareErrors(t, tstErr, etaErr, testTitle) {
		DeepEqual(t, tstVal, etaVal, testTitle)
	}
}

func CropMap(m map[string]interface{}, crop interface{}) (err error) {
	var keysToRemove = []string{}
	if keyName, ok := crop.(string); ok {
		for k, _ := range m {
			if k != keyName {
				keysToRemove = append(keysToRemove, k)
			}
		}
	} else if keyNames, ok := crop.([]string); ok {
		keyNameMap := map[string]struct{}{}
		for _, k := range keyNames {
			keyNameMap[k] = struct{}{}
		}
		for k, _ := range m {
			if _, ok := keyNameMap[k]; !ok {
				keysToRemove = append(keysToRemove, k)
			}
		}
	} else if keyNameMap, ok := crop.(map[string]interface{}); ok {
		for k, _ := range m {
			if _, ok := keyNameMap[k]; !ok {
				keysToRemove = append(keysToRemove, k)
			}
		}
	} else {
		err = bwerror.Error("<ansiOutline>crop<ansi> (<ansiPrimaryLiteral>%+v<ansi>) neither <ansiSecondaryLiteral>string<ansi>, nor <ansiSecondaryLiteral>[]string<ansi>, nor <ansiSecondaryLiteral>map[string]interface", crop)
	}
	for _, k := range keysToRemove {
		delete(m, k)
	}
	return
}

func MustCropMap(m map[string]interface{}, crop interface{}) {
	if err := CropMap(m, crop); err != nil {
		bwerror.Panic(err.Error())
	}
}
