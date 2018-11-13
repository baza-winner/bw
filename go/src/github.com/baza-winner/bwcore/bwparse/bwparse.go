package bwparse

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/runeprovider"
)

// ============================================================================

func ParseQw(p *runeprovider.Proxy, r rune) (result []interface{}, start runeprovider.PosStruct, ok bool, err error) {
	type State bool
	const (
		expectSpaceOrQwItemOrDelimiter State = true
		expectEndOfQwItem              State = false
	)
	var (
		state     State
		b         bool
		r2        rune
		delimiter rune
		s         string
	)

	if r2, b = Braces[r]; b || unicode.IsPunct(r) || unicode.IsSymbol(r) {
		start = p.Curr
		ok = true
		result = []interface{}{}
		state = expectSpaceOrQwItemOrDelimiter
		if b {
			delimiter = r2
		} else {
			delimiter = r
		}
	} else {
		ok = false
		return
	}

LOOP:
	for {
		if r, err = p.PullNonEOFRune(); err != nil {
			return
		}
		if state {
			if r == delimiter {
				break LOOP
			} else if !unicode.IsSpace(r) {
				s = string(r)
				state = expectEndOfQwItem
			}
		} else {
			if r == delimiter || unicode.IsSpace(r) {
				_ = p.PushRune()
				result = append(result, s)
				state = expectSpaceOrQwItemOrDelimiter
			} else {
				s += string(r)
			}
		}
	}
	return
}

var Braces = map[rune]rune{
	'(': ')',
	'{': '}',
	'<': '>',
	'[': ']',
}

// ============================================================================

func ParseWord(p *runeprovider.Proxy, r rune) (result string, start runeprovider.PosStruct, ok bool, err error) {
	if unicode.IsLetter(r) || r == '_' {
		result = string(r)
		start = p.Curr
		ok = true
	} else {
		ok = false
		return
	}
LOOP:
	for {
		if r, _, err = p.PullRuneOrEOF(); err != nil {
			return
		}
		if unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r) {
			result += string(r)
		} else {
			_ = p.PushRune()
			break LOOP
		}
	}
	return
}

// ============================================================================

func ParseString(p *runeprovider.Proxy, r rune) (result string, start runeprovider.PosStruct, ok bool, err error) {
	type State bool
	const (
		expectContent        State = true
		expectEscapedContent State = false
	)

	var (
		delimiter rune
		r2        rune
		b         bool
		state     State
	)

	if r == '"' || r == '\'' {
		delimiter = r
		state = expectContent
		start = p.Curr
		ok = true
	} else {
		ok = false
		return
	}

LOOP:
	for {
		if r, err = p.PullNonEOFRune(); err != nil {
			return
		}
		if state {
			switch r {
			case delimiter:
				break LOOP
			case '\\':
				state = expectEscapedContent
			default:
				result += string(r)
			}
		} else {
			switch r {
			case '"', '\'', '\\':
				result += string(r)
			default:
				if delimiter == '"' {
					if r2, b = EscapeRunes[r]; b {
						result += string(r2)
					} else {
						err = p.Unexpected(p.Curr)
						return
					}
				}
			}
			state = expectContent
		}
	}
	return
}

var EscapeRunes = map[rune]rune{
	'a': '\a',
	'b': '\b',
	'f': '\f',
	'n': '\n',
	'r': '\r',
	't': '\t',
	'v': '\v',
}

// ============================================================================

func ParseNumber(p *runeprovider.Proxy, r rune) (result interface{}, start runeprovider.PosStruct, ok bool, err error) {
	type State bool
	const (
		expectDigitOrUnderscore      State = true
		expectDigitOrUnderscoreOrDot State = false
	)
	var (
		s     string
		state State
	)

	switch r {
	case '-', '+':
		start = p.Curr
		s = string(r)
		ok = true
		start = p.Curr
		if r, err = p.PullNonEOFRune(); err != nil {
			return
		}
		switch r {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			s += string(r)
			state = expectDigitOrUnderscoreOrDot
		default:
			err = p.Unexpected(p.Curr)
			return
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		start = p.Curr
		s = string(r)
		state = expectDigitOrUnderscoreOrDot
		ok = true
		start = p.Curr
	default:
		ok = false
		return
	}

LOOP:
	for {
		if r, _, err = p.PullRuneOrEOF(); err != nil {
			return
		}
		switch r {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '_':
			s += string(r)
		default:
			if state == expectDigitOrUnderscoreOrDot && r == '.' {
				s += string(r)
				state = expectDigitOrUnderscore
			} else {
				p.PushRune()
				break LOOP
			}
		}
	}

	s = underscoreRegexp.ReplaceAllLiteralString(s, "")
	if strings.Contains(s, ".") && !zeroAfterDotRegexp.MatchString(s) {
		var _float64 float64
		if _float64, err = strconv.ParseFloat(s, 64); err == nil {
			result = _float64
		}
	} else {
		if pos := strings.LastIndex(s, "."); pos >= 0 {
			s = s[:pos]
		}
		var _int64 int64
		if _int64, err = strconv.ParseInt(s, 10, 64); err == nil {
			if int64(bw.MinInt8) <= _int64 && _int64 <= int64(bw.MaxInt8) {
				result = int8(_int64)
			} else if int64(bw.MinInt16) <= _int64 && _int64 <= int64(bw.MaxInt16) {
				result = int16(_int64)
			} else if int64(bw.MinInt32) <= _int64 && _int64 <= int64(bw.MaxInt32) {
				result = int32(_int64)
			} else {
				result = _int64
			}
		}
	}
	if err != nil {
		err = p.Unexpected(start, bwerr.Err(err))
	}

	return
}

var underscoreRegexp = regexp.MustCompile("[_]+")

var zeroAfterDotRegexp = regexp.MustCompile(`\.0+$`)

// ============================================================================

func ParseSpace(p *runeprovider.Proxy, r rune) (start runeprovider.PosStruct, ok bool, err error) {
	var isEOF bool
	if unicode.IsSpace(r) {
		start = p.Curr
		ok = true
	} else {
		ok = false
		return
	}
LOOP:
	for {
		if r, isEOF, err = p.PullRuneOrEOF(); err != nil {
			return
		} else if isEOF {
			break LOOP
		} else if !unicode.IsSpace(r) {
			err = p.Unexpected(p.Curr)
			return
		}
	}
	return
}

// ============================================================================

// ============================================================================
