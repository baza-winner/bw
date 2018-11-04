// Предоставялет функции для работы со строками.
package bwstr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/baza-winner/bwcore/bw"
)

// ============================================================================

func SmartQuote(ss ...string) (result string) {
	result = ``
	for i, s := range ss {
		if i > 0 {
			result += ` `
		}
		if strings.ContainsAny(s, ` "`) {
			result += fmt.Sprintf(`%q`, s)
		} else {
			result += s
		}
	}
	return
}

// ============================================================================

var underscoreRegexp = regexp.MustCompile("[_]+")

func ParseInt(s string) (result int, err error) {
	var _int64 int64
	if _int64, err = strconv.ParseInt(underscoreRegexp.ReplaceAllLiteralString(s, ""), 10, 64); err == nil {
		if int64(bw.MinInt) <= _int64 && _int64 <= int64(bw.MaxInt) {
			result = int(_int64)
		} else {
			err = fmt.Errorf("%d is out of range [%d, %d]", _int64, bw.MinInt, bw.MaxInt)
		}
	}
	return
}

func ParseNumber(s string) (value interface{}, err error) {
	s = underscoreRegexp.ReplaceAllLiteralString(s, ``)
	if strings.Contains(s, `.`) {
		var _float64 float64
		if _float64, err = strconv.ParseFloat(s, 64); err == nil {
			value = _float64
		}
	} else {
		var _int64 int64
		if _int64, err = strconv.ParseInt(s, 10, 64); err == nil {
			if int64(bw.MinInt8) <= _int64 && _int64 <= int64(bw.MaxInt8) {
				value = int8(_int64)
			} else if int64(bw.MinInt16) <= _int64 && _int64 <= int64(bw.MaxInt16) {
				value = int16(_int64)
			} else if int64(bw.MinInt32) <= _int64 && _int64 <= int64(bw.MaxInt32) {
				value = int32(_int64)
			} else {
				value = _int64
			}
		}
	}
	return
}

// ============================================================================
