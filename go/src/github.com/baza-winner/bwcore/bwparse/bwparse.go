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
)

// ============================================================================

func (p *Provider) ArrayOfString() (result []string, start PosStruct, ok bool, err error) {
	type State bool
	const (
		expectSpaceOrQwItemOrDelimiter State = true
		expectEndOfQwItem              State = false
	)
	var (
		state     State
		b         bool
		r, r2     rune
		delimiter rune
		s         string
		ps        PosStruct
	)

	start = p.Curr
	if p.Curr.Rune == '<' {
		delimiter = '>'
	} else {
		if p.Curr.Rune != 'q' {
			return
		}
		if ps, err = p.PosStruct(1); err != nil || ps.IsEOF || ps.Rune != 'w' {
			return
		}
		if ps, err = p.PosStruct(2); err != nil || ps.IsEOF {
			return
		}
		r = ps.Rune
		if r2, b = Braces[r]; !(b || unicode.IsPunct(r) || unicode.IsSymbol(r)) {
			return
		}
		if b {
			delimiter = r2
		} else {
			delimiter = r
		}
		p.Forward()
		p.Forward()
	}
	ok = true
	result = []string{}
	state = expectSpaceOrQwItemOrDelimiter

LOOP:
	for {
		if err = p.Forward(NonEOF); err != nil {
			return
		}
		r = p.Curr.Rune
		if state {
			if r == delimiter {
				break LOOP
			} else if !unicode.IsSpace(r) {
				s = string(r)
				state = expectEndOfQwItem
			}
		} else {
			if r == delimiter || unicode.IsSpace(r) {
				p.Backward()
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

func (p *Provider) Id() (result string, start PosStruct, ok bool, err error) {
	r := p.Curr.Rune
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
		if err = p.Forward(); err != nil {
			return
		}
		r = p.Curr.Rune
		if unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r) {
			result += string(r)
		} else {
			p.Backward()
			break LOOP
		}
	}
	return
}

// ============================================================================

func (p *Provider) String() (result string, start PosStruct, ok bool, err error) {
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
	r := p.Curr.Rune

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
		if err = p.Forward(NonEOF); err != nil {
			return
		}
		r = p.Curr.Rune
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

func (p *Provider) Int() (result int, start PosStruct, ok bool, err error) {
	var s string
	r := p.Curr.Rune

	switch r {
	case '-', '+':
		start = p.Curr
		s = string(r)
		ok = true
		start = p.Curr
		if err = p.Forward(NonEOF); err != nil {
			return
		}
		r = p.Curr.Rune
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
		if err = p.Forward(); err != nil {
			return
		}
		r = p.Curr.Rune
		switch r {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '_':
			s += string(r)
		default:
			p.Backward()
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

func (p *Provider) Number() (result interface{}, start PosStruct, ok bool, err error) {
	type State bool
	const (
		expectDigitOrUnderscore      State = true
		expectDigitOrUnderscoreOrDot State = false
	)
	var (
		s     string
		state State
	)
	r := p.Curr.Rune

	switch r {
	case '-', '+':
		start = p.Curr
		s = string(r)
		ok = true
		start = p.Curr
		if err = p.Forward(NonEOF); err != nil {
			return
		}
		r = p.Curr.Rune
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
		if err = p.Forward(); err != nil {
			return
		}
		r = p.Curr.Rune
		switch r {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '_':
			s += string(r)
		default:
			if state == expectDigitOrUnderscoreOrDot && r == '.' {
				s += string(r)
				state = expectDigitOrUnderscore
			} else {
				p.Backward()
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

func (p *Provider) SkipOptionalSpaceTillEOF() (err error) {
	for {
		if err = p.Forward(); err != nil || p.Curr.IsEOF {
			return
		} else if !unicode.IsSpace(p.Curr.Rune) {
			err = p.Unexpected(p.Curr)
			return
		}
	}
}

func (p *Provider) SkipOptionalSpace() (err error) {
	if err = p.Forward(NonEOF); err != nil {
		return
	}
	if unicode.IsSpace(p.Curr.Rune) {
	LOOP:
		for {
			if err = p.Forward(NonEOF); err != nil {
				return
			} else if !unicode.IsSpace(p.Curr.Rune) {
				break LOOP
			}
		}
	}
	return
}

// ============================================================================

func (p *Provider) Array() (result []interface{}, start PosStruct, ok bool, err error) {
	if p.Curr.Rune != '[' {
		ok = false
		return
	}

	start = p.Curr
	result = []interface{}{}
	ok = true
	if err = p.SkipOptionalSpace(); err != nil {
		return
	}
LOOP:
	for {
		if p.Curr.Rune == ']' {
			break LOOP
		}

		var val interface{}
		if val, _, ok, err = p.Val(); err != nil || !ok {
			if err == nil {
				err = p.Unexpected(p.Curr)
			}
			return
		}
		if ss, b := val.([]string); !b {
			result = append(result, val)
		} else {
			for _, s := range ss {
				result = append(result, s)
			}
		}
		if err = p.SkipOptionalSpace(); err != nil {
			return
		}
		if p.Curr.Rune == ',' {
			if err = p.SkipOptionalSpace(); err != nil {
				return
			}
		}
	}

	return
}

func (p *Provider) Map() (result map[string]interface{}, start PosStruct, ok bool, err error) {
	if p.Curr.Rune != '{' {
		ok = false
		return
	}

	start = p.Curr
	result = map[string]interface{}{}
	ok = true
	if err = p.SkipOptionalSpace(); err != nil {
		return
	}

LOOP:
	for {
		if p.Curr.Rune == '}' {
			break LOOP
		}
		var (
			key string
			b   bool
		)

		if key, _, b, err = p.String(); err != nil || b {
			if err != nil {
				return
			}
		} else if key, _, b, err = p.Id(); err != nil || b {
			if err != nil {
				return
			}
		} else {
			err = p.Unexpected(p.Curr)
			return
		}

		if err = p.SkipOptionalSpace(); err != nil {
			return
		}

		if p.Curr.Rune == ':' {
			if err = p.SkipOptionalSpace(); err != nil {
				return
			}
		} else if p.Curr.Rune == '=' {
			if err = p.Forward(NonEOF); err != nil {
				return
			}
			if p.Curr.Rune != '>' {
				err = p.Unexpected(p.Curr)
				return
			}
			if err = p.SkipOptionalSpace(); err != nil {
				return
			}
		}

		var val interface{}
		if val, _, ok, err = p.Val(); err != nil || !ok {
			if err == nil {
				err = p.Unexpected(p.Curr)
			}
			return
		}

		result[key] = val

		if err = p.SkipOptionalSpace(); err != nil {
			return
		}
		if p.Curr.Rune == ',' {
			if err = p.SkipOptionalSpace(); err != nil {
				return
			}
		}

	}

	return
}

func (p *Provider) Val() (result interface{}, start PosStruct, ok bool, err error) {
	var (
		s    string
		ps   PosStruct
		val  interface{}
		vals []interface{}
		ss   []string
		m    map[string]interface{}
		b    bool
		path bw.ValPath
	)

	ok = true
	start = p.Curr

	if vals, _, b, err = p.Array(); err != nil || b {
		if err != nil {
			return
		}
		result = vals
	} else if s, _, b, err = p.String(); err != nil || b {
		if err != nil {
			return
		}
		result = s
	} else if val, _, b, err = p.Number(); err != nil || b {
		if err != nil {
			return
		}
		result = val
	} else if path, _, b, err = p.Path(); err != nil || b {
		if err != nil {
			return
		}
		result = path
	} else if m, _, b, err = p.Map(); err != nil || b {
		if err != nil {
			return
		}
		result = m
	} else if ss, _, b, err = p.ArrayOfString(); err != nil || b {
		if err != nil {
			return
		}
		result = ss
	} else if s, ps, b, err = p.Id(); err != nil || b {
		if err != nil {
			return
		}
		switch s {
		case "true":
			result = true
		case "false":
			result = false
		case "nil", "null":
			result = nil
		case "Bool", "String", "Int", "Number", "Map", "Array", "ArrayOf":
			result = s
		default:
			err = p.Unexpected(ps, bw.Fmt(ansiUnexpectedWord, s))
			return
		}
	} else {
		ok = false
		return
	}

	return
}

var ansiUnexpectedWord string

func init() {
	ansiUnexpectedWord = ansi.String("unexpected <ansiErr>%q<ansi>")
}

// ============================================================================

func (p *Provider) Path(optBases ...[]bw.ValPath) (result bw.ValPath, start PosStruct, ok bool, err error) {
	start = p.Curr
	var ps PosStruct
	if p.Curr.Rune != '{' {
		ok = false
		return
	}
	if ps, err = p.PosStruct(1); err != nil || ps.IsEOF || ps.Rune != '{' {
		return
	}
	p.Forward()
	ok = true

	if err = p.SkipOptionalSpace(); err != nil {
		return
	}
	if result, err = p.PathContent(NoAutoForward, optBases...); err != nil {
		return
	}
	if err = p.SkipOptionalSpace(); err != nil {
		return
	}
	if p.Curr.Rune != '}' {
		err = p.Unexpected(p.Curr)
		return
	}
	if err = p.Forward(); err != nil {
		return
	}
	if p.Curr.Rune != '}' {
		err = p.Unexpected(p.Curr)
		return
	}

	return
}

func (p *Provider) subPath(optBases ...[]bw.ValPath) (result bw.ValPath, start PosStruct, ok bool, err error) {
	if p.Curr.Rune == '(' {
		ok = true
		start = p.Curr
		if err = p.SkipOptionalSpace(); err != nil {
			return
		}
		if result, err = p.PathContent(NoAutoForward, optBases...); err != nil {
			return
		}
		if err = p.SkipOptionalSpace(); err != nil {
			return
		}
		if p.Curr.Rune != ')' {
			err = p.Unexpected(p.Curr)
			return
		}
	}
	return
}

const (
	NoAutoForward bool = true
	AutoForward   bool = false
)

func (p *Provider) PathContent(noAutoForward bool, optBases ...[]bw.ValPath) (result bw.ValPath, err error) {

	if !noAutoForward {
		if err = p.Forward(NonEOF); err != nil {
			return
		}
	}

	var (
		idx   int
		s     string
		b     bool
		sp    bw.ValPath
		ps    PosStruct
		bases []bw.ValPath
	)
	if len(optBases) > 0 {
		bases = optBases[0]
	}

LOOP:
	for {
		if p.Curr.Rune == '.' &&
			len(result) == 0 {
			if len(bases) == 0 {
				break LOOP
			} else if len(result) == 0 {
				result = append(result, bases[0]...)
				p.Backward()

			} else {
				err = p.Unexpected(p.Curr)
				return
			}
		} else if idx, _, b, err = p.Int(); b || err != nil {
			if err != nil {
				return
			}

			result = append(
				result,
				bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: idx},
			)

		} else if s, _, b, err = p.Id(); b || err != nil {

			if err != nil {
				return
			}
			result = append(
				result,
				bw.ValPathItem{Type: bw.ValPathItemKey, Key: s},
			)

		} else if sp, _, b, err = p.subPath(); b || err != nil {
			if err != nil {
				return
			}
			result = append(
				result,
				bw.ValPathItem{Type: bw.ValPathItemPath, Path: sp},
			)

		} else if p.Curr.Rune == '#' {
			result = append(
				result,
				bw.ValPathItem{Type: bw.ValPathItemHash},
			)
			break LOOP
		} else {
			if len(result) == 0 {
				if p.Curr.Rune == '$' {
					if err = p.Forward(NonEOF); err != nil {
						return
					}
					if s, _, b, err = p.Id(); err != nil || b {
						if err != nil {
							return
						}
						result = append(
							result,
							bw.ValPathItem{Type: bw.ValPathItemVar, Key: s},
						)
						goto CONTINUE
					} else if idx, ps, b, err = p.Int(); err != nil || b {
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

		if err = p.Forward(); err != nil {
			return
		}

		if p.Curr.Rune != '.' {
			p.Backward()
			break LOOP
		}

		if err = p.Forward(); err != nil {
			return
		}
	}
	return
}

// ============================================================================
