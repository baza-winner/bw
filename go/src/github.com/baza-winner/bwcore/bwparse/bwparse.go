package bwparse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/runeprovider"
)

// ============================================================================

func ArrayOfString(p *runeprovider.Proxy, r rune) (result []string, start runeprovider.PosStruct, ok bool, err error) {
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
		isEOF     bool
	)

	// bwdebug.Print("r", string(r))
	if r == '<' {
		delimiter = '>'
	} else {
		// if r, isEOF, err = p.Rune(); err != nil || isEOF || r != 'q' {
		if r != 'q' {
			return
		}
		if r, isEOF, err = p.Rune(1); err != nil || isEOF || r != 'w' {
			return
		}
		if r, isEOF, err = p.Rune(2); err != nil || isEOF {
			return
		}
		if r2, b = Braces[r]; !(b || unicode.IsPunct(r) || unicode.IsSymbol(r)) {
			return
		}
		if b {
			delimiter = r2
		} else {
			delimiter = r
		}
		p.PullRune()
		p.PullRune()
	}
	// p.PullRune()
	// bwdebug.Print("!!!", "r", string(r))
	start = p.Curr
	ok = true
	result = []string{}
	state = expectSpaceOrQwItemOrDelimiter

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
		// bwdebug.Print("r", string(r), "result", result)
	}
	// bwdebug.Print("r", string(r), "result", result)
	return
}

var Braces = map[rune]rune{
	'(': ')',
	'{': '}',
	'<': '>',
	'[': ']',
}

// ============================================================================

func Id(p *runeprovider.Proxy, r rune) (result string, start runeprovider.PosStruct, ok bool, err error) {
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

// // ============================================================================

// func ParseVarName(p *runeprovider.Proxy, r rune) (result string, start runeprovider.PosStruct, ok bool, err error) {
// 	if r == '$' {
// 		// result = string(r)
// 		start = p.Curr
// 		ok = true
// 	} else {
// 		ok = false
// 		return
// 	}
// LOOP:
// 	for {
// 		if r, _, err = p.PullRuneOrEOF(); err != nil {
// 			return
// 		}
// 		if unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r) {
// 			result += string(r)
// 		} else {
// 			_ = p.PushRune()
// 			break LOOP
// 		}
// 	}
// 	return
// }

// ============================================================================

func String(p *runeprovider.Proxy, r rune) (result string, start runeprovider.PosStruct, ok bool, err error) {
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

func ParseInt(p *runeprovider.Proxy, r rune) (result int, start runeprovider.PosStruct, ok bool, err error) {
	var s string

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
		default:
			err = p.Unexpected(p.Curr)
			return
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		start = p.Curr
		s = string(r)
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
			p.PushRune()
			break LOOP
		}
	}

	s = underscoreRegexp.ReplaceAllLiteralString(s, "")
	var _int64 int64
	if _int64, err = strconv.ParseInt(underscoreRegexp.ReplaceAllLiteralString(s, ""), 10, 64); err == nil {
		if int64(bw.MinInt) <= _int64 && _int64 <= int64(bw.MaxInt) {
			result = int(_int64)
		} else {
			err = fmt.Errorf("%d is out of range [%d, %d]", _int64, bw.MinInt, bw.MaxInt)
		}
	}
	if err != nil {
		err = p.Unexpected(start, bwerr.Err(err))
	}

	return
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

func ParseSpace(p *runeprovider.Proxy, r rune) (result rune, isEOF bool, start runeprovider.PosStruct, ok bool, err error) {
	if unicode.IsSpace(r) {
		start = p.Curr
		ok = true
	} else {
		ok = false
		return
	}
LOOP:
	for {
		result = r
		if r, isEOF, err = p.PullRuneOrEOF(); err != nil {
			return
		} else if isEOF || !unicode.IsSpace(r) {
			p.PushRune()
			break LOOP
		}
	}
	return
}

// ============================================================================

func SkipOptionalSpaceTillEOF(p *runeprovider.Proxy, r rune) (err error) {
	var isEOF, ok bool
	if r, isEOF, err = p.PullRuneOrEOF(); err != nil || isEOF {
		return
	}

	// bwdebug.Print("r", string(r))
	if _, isEOF, _, ok, err = ParseSpace(p, r); err != nil || ok && isEOF {
		return
	} else {
		// bwdebug.Print("p.Next[0]", bwjson.Pretty(p.Next[0]))
		_ = p.PullRune()
		err = p.Unexpected(p.Curr)
		return
	}
	return
}

func SkipOptionalSpace(p *runeprovider.Proxy) (r rune, err error) {
	if r, err = p.PullNonEOFRune(); err != nil {
		return
	}
	if unicode.IsSpace(r) {
	LOOP:
		for {
			// result = r
			if r, err = p.PullNonEOFRune(); err != nil {
				return
			} else if !unicode.IsSpace(r) {
				break LOOP
			}
			// else if isEOF || !unicode.IsSpace(r) {
			// 	p.PushRune()
			// 	break LOOP
			// }
		}
	}
	// var isEOF, ok bool
	// if result, isEOF, _, ok, err = ParseSpace(p, r); err != nil || ok && isEOF {
	// 	if err == nil {
	// 		_ = p.PullRune()
	// 		err = p.Unexpected(p.Curr)
	// 	}
	// 	return
	// } else if !ok {
	// 	// bwdebug.Print("r", string(r))
	// 	result = r
	// } else {
	// 	p.PullRuneOrEOF()

	// }
	// bwdebug.Print("r", string(r))
	return
}

func ParseArray(p *runeprovider.Proxy, r rune) (result []interface{}, start runeprovider.PosStruct, ok bool, err error) {
	if r != '[' {
		ok = false
		return
	}

	start = p.Curr
	result = []interface{}{}
	ok = true
	if r, err = SkipOptionalSpace(p); err != nil {
		return
	}
LOOP:
	for {
		if r == ']' {
			break LOOP
		}

		var val interface{}
		if val, err = ParseVal(p, r); err != nil {
			return
		}
		if ss, b := val.([]string); !b {
			result = append(result, val)
		} else {
			for _, s := range ss {
				result = append(result, s)
			}
		}
		if r, err = SkipOptionalSpace(p); err != nil {
			return
		}
		if r == ',' {
			if r, err = SkipOptionalSpace(p); err != nil {
				return
			}
		}
	}

	return
}

func Map(p *runeprovider.Proxy, r rune) (result map[string]interface{}, start runeprovider.PosStruct, ok bool, err error) {
	if r != '{' {
		ok = false
		return
	}

	start = p.Curr
	result = map[string]interface{}{}
	ok = true
	if r, err = SkipOptionalSpace(p); err != nil {
		return
	}
LOOP:
	for {
		if r == '}' {
			break LOOP
		}
		var (
			key string
			b   bool
		)

		if key, _, b, err = String(p, r); err != nil || b {
			if err != nil {
				return
			}
		} else if key, _, b, err = Id(p, r); err != nil || b {
			if err != nil {
				return
			}
		} else {
			err = p.Unexpected(p.Curr)
			return
		}
		// bwdebug.Print("key", key)

		if r, err = SkipOptionalSpace(p); err != nil {
			return
		}
		// bwdebug.Print("r", string(r))

		if r == ':' {
			if r, err = SkipOptionalSpace(p); err != nil {
				return
			}
		} else if r == '=' {
			if r, err = p.PullNonEOFRune(); err != nil {
				return
			}
			if r != '>' {
				err = p.Unexpected(p.Curr)
				return
			}
			if r, err = SkipOptionalSpace(p); err != nil {
				return
			}
		}

		var val interface{}
		if val, err = ParseVal(p, r); err != nil {
			return
		}
		// bwdebug.Print("val", val)

		result[key] = val

		if r, err = SkipOptionalSpace(p); err != nil {
			return
		}
		if r == ',' {
			if r, err = SkipOptionalSpace(p); err != nil {
				return
			}
		}
		// bwdebug.Print("r", string(r))

	}

	return
}

func ParseVal(p *runeprovider.Proxy, r rune) (result interface{}, err error) {

	// type PrimaryState uint8

	// const (
	// 	Begin PrimaryState = iota
	// 	// ExpectSpaceOrQwItemOrDelimiter
	// 	// ExpectSpaceOrMapKey
	// 	// End
	// )

	// type SecondaryState uint8

	// const (
	// 	None SecondaryState = iota
	// 	// orArrayItemSeparator
	// 	// orMapKeySeparator
	// 	// orMapValueSeparator
	// )

	type ItemKind uint8

	const (
		ItemString ItemKind = iota
		ItemQw
		ItemQwItem
		ItemNumber
		ItemWord
		ItemKey
		ItemMap
		ItemArray
	)

	type StackItem struct {
		PosStruct runeprovider.PosStruct
		Kind      ItemKind
		S         string
		Result    interface{}
		Delimiter rune
	}

	var (
		needFinish bool
		// skipPostProcess bool
		ok    bool
		s     string
		start runeprovider.PosStruct
		val   interface{}
		vals  []interface{}
		ss    []string
		m     map[string]interface{}
		stack []StackItem
		// primary PrimaryState
		// secondary SecondaryState
	)
LOOP:
	// for primary != End {
	for {
		r, _, _ = p.Rune()
		needFinish = false
		switch {
		// case unicode.IsSpace(r):

		case r == '{':
			if m, start, _, err = Map(p, r); err != nil {
				return
			}
			stack = append(stack, StackItem{
				PosStruct: p.Curr,
				Kind:      ItemMap,
				Result:    m,
				Delimiter: '}',
			})
			needFinish = true
		case r == '<':
			if ss, start, _, err = ArrayOfString(p, r); err != nil {
				return
			}
			stack = append(stack, StackItem{
				PosStruct: start,
				Kind:      ItemQw,
				Result:    ss,
				Delimiter: '>',
			})
			needFinish = true
		case r == '[':
			if vals, start, _, err = ParseArray(p, r); err != nil {
				return
			}
			stack = append(stack, StackItem{
				PosStruct: start,
				Kind:      ItemArray,
				Result:    vals,
				Delimiter: ']',
			})
			needFinish = true
		case len(stack) > 0 && stack[len(stack)-1].Kind == ItemArray && r == stack[len(stack)-1].Delimiter:
			needFinish = true
		case r == '-' || r == '+' || unicode.IsDigit(r):
			if val, start, _, err = ParseNumber(p, r); err != nil {
				return
			}
			stack = append(stack, StackItem{
				PosStruct: start,
				Kind:      ItemNumber,
				Result:    val,
			})
			needFinish = true

		case r == '"' || r == '\'':
			if s, start, _, err = String(p, r); err != nil {
				return
			}
			stack = append(stack, StackItem{
				PosStruct: start,
				Kind:      ItemString,
				Delimiter: r,
				Result:    s,
			})
			needFinish = true
		case unicode.IsLetter(r) || r == '_':
			if s, start, _, err = Id(p, r); err != nil {
				return
			}
			needFinish = true

			stack = append(stack, StackItem{
				PosStruct: start,
				Kind:      ItemWord,
				S:         s,
			})
		default:
			err = p.Unexpected(p.Curr)
			return
		}

		_ = needFinish
		// if needFinish {
		switch stack[len(stack)-1].Kind {
		case ItemWord:
			switch stack[len(stack)-1].S {
			case "true":
				stack[len(stack)-1].Result = true
			case "false":
				stack[len(stack)-1].Result = false
			case "nil", "null":
				stack[len(stack)-1].Result = nil
			case "Bool", "String", "Int", "Number", "Map", "Array", "ArrayOf":
				stack[len(stack)-1].Result = stack[len(stack)-1].S
			case "qw":
				if r, err = p.PullNonEOFRune(); err != nil {
					return
				}
				var b bool
				if _, b = Braces[r]; !(b || unicode.IsPunct(r) || unicode.IsSymbol(r)) {
					err = p.Unexpected(p.Curr)
					return
				}
				p.PushRune()
				p.PushRune()
				p.PushRune()
				r = *p.Curr.RunePtr
				if ss, start, ok, err = ArrayOfString(p, r); err != nil {
					return
				}
				if !ok {
					err = p.Unexpected(p.Curr)
					return
				}
				stack[len(stack)-1].Kind = ItemQw
				stack[len(stack)-1].Result = ss
			default:
				err = p.Unexpected(stack[len(stack)-1].PosStruct, bw.Fmt(ansi.String("unexpected <ansiErr>%q<ansi>"), stack[len(stack)-1].S))
				return
			}
		}

		break LOOP
		// } else {
		// 	bwdebug.Print("!HERE")
		// }
	}
	result = stack[0].Result
	return
}

// ============================================================================

func Path(p *runeprovider.Proxy, r rune, optBases ...[]bw.ValPath) (result bw.ValPath, start runeprovider.PosStruct, ok bool, err error) {

	ok = true
	start = p.Curr
	defer func() {
		if err != nil {
			ok = false
		}
	}()

	var bases []bw.ValPath
	if len(optBases) > 0 {
		bases = optBases[0]
	}

LOOP:
	for {
		var (
			idx int
			s   string
			b   bool
			sp  bw.ValPath
			ps  runeprovider.PosStruct
		)
		if r == '.' &&
			len(result) == 0 {
			if len(bases) == 0 {
				break LOOP
			} else if len(result) == 0 {
				result = append(result, bases[0]...)
				p.PushRune()

			} else {
				err = p.Unexpected(p.Curr)
				return
			}
		} else if idx, _, b, err = ParseInt(p, r); b || err != nil {
			if err != nil {
				return
			}

			result = append(
				result,
				bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: idx},
			)

		} else if s, _, b, err = Id(p, r); b || err != nil {

			if err != nil {
				return
			}
			result = append(
				result,
				bw.ValPathItem{Type: bw.ValPathItemKey, Key: s},
			)

		} else if sp, _, b, err = subPath(p, r); b || err != nil {
			if err != nil {
				return
			}
			result = append(
				result,
				bw.ValPathItem{Type: bw.ValPathItemPath, Path: sp},
			)

		} else if r == '#' {
			result = append(
				result,
				bw.ValPathItem{Type: bw.ValPathItemHash},
			)
			break LOOP
		} else {
			if len(result) == 0 {
				if r == '$' {
					if r, err = p.PullNonEOFRune(); err != nil {
						return
					}
					if s, _, b, err = Id(p, r); err != nil || b {
						if err != nil {
							return
						}
						result = append(
							result,
							bw.ValPathItem{Type: bw.ValPathItemVar, Key: s},
						)
						goto CONTINUE
					} else if idx, ps, b, err = ParseInt(p, r); err != nil || b {
						if err != nil {
							return
						}
						var nidx int
						l := len(bases)
						if nidx, b = bw.NormalIdx(idx, len(bases)); !b {
							err = p.Unexpected(ps, bw.Fmt(ansi.String("unexpected base path idx <ansiVal>%d<ansi> (len(bases): <ansiVal>%d)"), idx, l))
							return
						}
						result = append(result, bases[nidx]...)
						goto CONTINUE
					}
				}
			}
			if len(result) == 0 {
				b = false
				return
			}
			err = p.Unexpected(p.Curr)
			return
		}
	CONTINUE:

		if r, _, err = p.PullRuneOrEOF(); err != nil {
			return
		}

		if r != '.' {
			p.PushRune()
			break LOOP
		}

		if r, _, err = p.PullRuneOrEOF(); err != nil {
			return
		}
	}
	return
}

func subPath(p *runeprovider.Proxy, r rune, optBases ...[]bw.ValPath) (result bw.ValPath, start runeprovider.PosStruct, ok bool, err error) {
	b := true
	defer func() {
		if !b && err == nil {
			err = p.Unexpected(p.Curr)
		}
	}()
	if r == '{' {
		start = p.Curr
		r, err = p.PullNonEOFRune()
		if result, _, b, err = Path(p, r, optBases...); err != nil || !b {
			return
		}
		if r, err = p.PullNonEOFRune(); err != nil || r != '}' {
			b = false
			return
		}
		ok = true
	}
	return
}

// ============================================================================
