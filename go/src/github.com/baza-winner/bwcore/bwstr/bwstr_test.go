package bwstr_test

import (
	"testing"

	"github.com/baza-winner/bwcore/bwstr"
	"github.com/baza-winner/bwcore/bwtesting"
)

// ============================================================================

func TestSmartQuote(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"first": {
			In: []interface{}{
				[]string{"some", `thi"ng`, `go od`},
			},
			Out: []interface{}{
				`some "thi\"ng" "go od"`,
			},
		},
	}
	testsToRun := tests
	bwtesting.BwRunTests(t, bwstr.SmartQuote, testsToRun)
}

// ============================================================================

// func PluralWord(count int, word string, word1 string, word2_4 string, _word5more ...string) (result string) {
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

// // ============================================================================

// var underscoreRegexp = regexp.MustCompile("[_]+")

// func ParseInt(s string) (result int, err error) {
// 	var _int64 int64
// 	if _int64, err = strconv.ParseInt(underscoreRegexp.ReplaceAllLiteralString(s, ""), 10, 64); err == nil {
// 		if int64(MinInt) <= _int64 && _int64 <= int64(MaxInt) {
// 			result = int(_int64)
// 		} else {
// 			err = fmt.Errorf("%d is out of range [%d, %d]", _int64, MinInt, MaxInt)
// 		}
// 	}
// 	return
// }

// func ParseNumber(s string) (value interface{}, err error) {
// 	s = underscoreRegexp.ReplaceAllLiteralString(s, ``)
// 	if strings.Contains(s, `.`) {
// 		var _float64 float64
// 		if _float64, err = strconv.ParseFloat(s, 64); err == nil {
// 			value = _float64
// 		}
// 	} else {
// 		var _int64 int64
// 		if _int64, err = strconv.ParseInt(s, 10, 64); err == nil {
// 			if int64(MinInt8) <= _int64 && _int64 <= int64(MaxInt8) {
// 				value = int8(_int64)
// 			} else if int64(MinInt16) <= _int64 && _int64 <= int64(MaxInt16) {
// 				value = int16(_int64)
// 			} else if int64(MinInt32) <= _int64 && _int64 <= int64(MaxInt32) {
// 				value = int32(_int64)
// 			} else {
// 				value = _int64
// 			}
// 		}
// 	}
// 	return
// }

// // ============================================================================
