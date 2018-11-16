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
	"github.com/baza-winner/bwcore/bwrune"
)

// ============================================================================

type PosInfo struct {
	IsEOF       bool
	Rune        rune
	Pos         int
	Line        uint
	Col         uint
	Prefix      string
	PrefixStart int
}

// ============================================================================

type P struct {
	Prov          bwrune.Provider
	Curr          PosInfo
	Next          []PosInfo
	preLineCount  uint
	postLineCount uint
}

func From(p bwrune.Provider) (result *P) {
	result = &P{
		Prov:          p,
		Curr:          PosInfo{Pos: -1, Line: 1},
		Next:          []PosInfo{},
		preLineCount:  3,
		postLineCount: 3,
	}
	return
}

const Initial bool = false

func (p *P) Forward(nonInitial bool) {
	if p.Curr.Pos < 0 || nonInitial && !p.Curr.IsEOF {
		if len(p.Next) == 0 {
			p.pullRune(&p.Curr)
		} else {
			p.Curr = p.Next[len(p.Next)-1]
			p.Next = p.Next[:len(p.Next)-1]
		}
	}
}

func (p *P) CheckNotEOF() (err error) {
	if p.Curr.IsEOF {
		err = p.Unexpected(p.Curr)
	}
	return
}

func (p *P) PosInfo(ofs uint) (ps PosInfo, err error) {
	if ofs > 0 {
		if len(p.Next) >= int(ofs) {
			ps = p.Next[len(p.Next)-int(ofs)]
		} else {
			if len(p.Next) > 0 {
				ps = p.Next[0]
			} else {
				ps = p.Curr
			}
			lookahead := []PosInfo{}
			for i := int(ofs) - len(p.Next); i > 0; i-- {
				p.pullRune(&ps)
				lookahead = append(lookahead, ps)
				if ps.IsEOF {
					break
				}
			}
			newNext := []PosInfo{}
			for i := len(lookahead) - 1; i >= 0; i-- {
				newNext = append(newNext, lookahead[i])
			}
			p.Next = append(newNext, p.Next...)
		}
	} else {
		ps = p.Curr
	}
	return
}

func (p *P) Unexpected(ps PosInfo, optFmt ...bw.I) (result error) {
	var msg string
	if ps.IsEOF {
		msg = ansiUnexpectedEOF
	} else if len(optFmt) == 0 {
		r := ps.Rune
		msg = fmt.Sprintf(ansiUnexpectedChar, r, r)
	} else {
		msg = bw.Spew.Sprintf(optFmt[0].FmtString(), optFmt[0].FmtArgs()...)
	}
	result = bwerr.From(msg + p.suffix(ps))
	return
}

// ============================================================================

const (
	TillNonEOF bool = false
	TillEOF    bool = true
)

func (p *P) SkipSpace(tillEOF bool) (err error) {
	p.Forward(Initial)

	for {
		if p.Curr.IsEOF || !unicode.IsSpace(p.Curr.Rune) {
			break
		}
		p.Forward(true)
	}
	if !tillEOF && p.Curr.IsEOF || tillEOF && !p.Curr.IsEOF {
		err = p.Unexpected(p.Curr)
		return
	}
	return
}

// ============================================================================

func (p *P) ArrayOfString() (result []string, start PosInfo, ok bool, err error) {
	var (
		b         bool
		r, r2     rune
		delimiter rune
		s         string
		ps        PosInfo
	)

	p.Forward(Initial)
	start = p.Curr
	r = p.Curr.Rune
	if r == '<' {
		delimiter = '>'
	} else {
		if r != 'q' {
			return
		}
		if ps, err = p.PosInfo(1); err != nil || ps.IsEOF || ps.Rune != 'w' {
			return
		}
		if ps, err = p.PosInfo(2); err != nil || ps.IsEOF {
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
		p.Forward(true)
		p.Forward(true)
	}
	ok = true
	result = []string{}

	p.Forward(true)
LOOP:
	for {
		if err = p.SkipSpace(TillNonEOF); err != nil {
			return
		}

		r = p.Curr.Rune
		if r == delimiter {
			p.Forward(true)
			break LOOP
		}

		s = string(r)

		for {
			p.Forward(true)
			if err = p.CheckNotEOF(); err != nil {
				return
			}
			r = p.Curr.Rune
			if unicode.IsSpace(r) || r == delimiter {
				break
			}
			s += string(r)
		}
		result = append(result, s)
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

func (p *P) Id() (result string, start PosInfo, ok bool, err error) {
	p.Forward(Initial)
	r := p.Curr.Rune
	if !(unicode.IsLetter(r) || r == '_') {
		return
	}
	ok = true
	result = string(r)
	start = p.Curr

LOOP:
	for {
		p.Forward(true)
		r = p.Curr.Rune
		if unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r) {
			result += string(r)
		} else {
			break LOOP
		}
	}
	return
}

// ============================================================================

func (p *P) String() (result string, start PosInfo, ok bool, err error) {
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

	p.Forward(Initial)
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
		p.Forward(true)
		if err = p.CheckNotEOF(); err != nil {
			return
		}
		r = p.Curr.Rune
		if state {
			switch r {
			case delimiter:
				p.Forward(true)
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

func (p *P) Int() (result int, start PosInfo, ok bool, err error) {
	var (
		s string
		r rune
	)

	if s, start, ok, err = p.looksLikeNumber(); err != nil || !ok {
		return
	}

LOOP:
	for {
		r = p.Curr.Rune
		if '0' <= r && r <= '9' || r == '_' {
			s += string(r)
		} else {
			break LOOP
		}
		p.Forward(true)
	}

	if result, err = parseInt(underscoreRegexp.ReplaceAllLiteralString(s, "")); err != nil {
		err = p.Unexpected(start, bwerr.Err(err))
	}

	return
}

// ============================================================================

func (p *P) Number() (result interface{}, start PosInfo, ok bool, err error) {
	type State bool
	const (
		expectDigitOrUnderscore      State = false
		expectDigitOrUnderscoreOrDot State = true
	)
	var (
		s     string
		state State
		r     rune
	)

	if s, start, ok, err = p.looksLikeNumber(); err != nil || !ok {
		return
	}

	state = expectDigitOrUnderscoreOrDot

LOOP:
	for {
		r = p.Curr.Rune
		if '0' <= r && r <= '9' || r == '_' {
			s += string(r)
		} else if state && r == '.' {
			s += string(r)
			state = expectDigitOrUnderscore
		} else {
			break LOOP
		}
		p.Forward(true)
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
		result, err = parseInt(s)
	}
	if err != nil {
		err = p.Unexpected(start, bwerr.Err(err))
	}

	return
}

var underscoreRegexp = regexp.MustCompile("[_]+")

var zeroAfterDotRegexp = regexp.MustCompile(`\.0+$`)

// ============================================================================

func (p *P) Array() (result []interface{}, start PosInfo, ok bool, err error) {
	p.Forward(Initial)
	if p.Curr.Rune != '[' {
		return
	}
	ok = true
	start = p.Curr
	result = []interface{}{}

	var b bool
LOOP:
	for {
		p.Forward(true)
		if err = p.SkipSpace(TillNonEOF); err != nil {
			return
		}
	NEXT:
		if p.Curr.Rune == ']' {
			p.Forward(true)
			break LOOP
		}

		var val interface{}
		if val, err = p.parseVal(); err != nil {
			return
		}
		if ss, b := val.([]string); !b {
			result = append(result, val)
		} else {
			for _, s := range ss {
				result = append(result, s)
			}
		}

		if b, err = p.skipComma(); err != nil {
			return
		} else if !b {
			goto NEXT
		}
	}

	return
}

func (p *P) Map() (result map[string]interface{}, start PosInfo, ok bool, err error) {
	p.Forward(Initial)
	if p.Curr.Rune != '{' {
		return
	}
	ok = true
	start = p.Curr
	result = map[string]interface{}{}

	var (
		key string
		b   bool
	)
LOOP:
	for {
		p.Forward(true)
		if err = p.SkipSpace(TillNonEOF); err != nil {
			return
		}
	NEXT:
		if p.Curr.Rune == '}' {
			p.Forward(true)
			break LOOP
		}

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

		if err = p.SkipSpace(TillNonEOF); err != nil {
			return
		}

		var needForwardAndSkipSpace bool
		if p.Curr.Rune == ':' {
			needForwardAndSkipSpace = true
		} else if p.Curr.Rune == '=' {
			p.Forward(true)
			if p.Curr.Rune != '>' {
				err = p.Unexpected(p.Curr)
				return
			}
			needForwardAndSkipSpace = true
		}
		if needForwardAndSkipSpace {
			p.Forward(true)
			if err = p.SkipSpace(TillNonEOF); err != nil {
				return
			}
		}

		if result[key], err = p.parseVal(); err != nil {
			return
		}

		if b, err = p.skipComma(); err != nil {
			return
		} else if !b {
			goto NEXT
		}
	}

	return
}

func (p *P) Val() (result interface{}, start PosInfo, ok bool, err error) {
	var (
		s string
		b bool
	)

	ok = true

	p.Forward(Initial)

	if result, _, b, err = p.Array(); b {
	} else if result, start, b, err = p.String(); b {
	} else if result, start, b, err = p.Number(); b {
	} else if result, start, b, err = p.Path(PathA{}); b {
	} else if result, start, b, err = p.Map(); b {
	} else if result, start, b, err = p.ArrayOfString(); b {
	} else if s, start, b, err = p.Id(); b {
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
			err = p.Unexpected(start, bw.Fmt(ansiUnexpectedWord, s))
			return
		}
	} else {
		ok = false
		return
	}

	return
}

// ============================================================================

func (p *P) Path(a PathA) (result bw.ValPath, start PosInfo, ok bool, err error) {

	p.Forward(Initial)
	start = p.Curr
	if p.Curr.Rune != '{' {
		return
	}
	var ps PosInfo
	if ps, err = p.PosInfo(1); err != nil || ps.IsEOF || ps.Rune != '{' {
		return
	}
	ok = true
	p.Forward(true)
	p.Forward(true)

	if err = p.SkipSpace(TillNonEOF); err != nil {
		return
	}
	if result, err = p.PathContent(a); err != nil {
		return
	}
	if err = p.SkipSpace(TillNonEOF); err != nil {
		return
	}
	for i := 2; i > 0; i-- {
		if p.Curr.Rune != '}' {
			err = p.Unexpected(p.Curr)
			return
		}
		p.Forward(true)
	}

	return
}

const (
	NoAutoForward bool = true
	AutoForward   bool = false
)

type PathA struct {
	Bases     []bw.ValPath
	isSubPath bool
}

func (p *P) PathContent(a PathA) (result bw.ValPath, err error) {

	p.Forward(Initial)
	if err = p.CheckNotEOF(); err != nil {
		return
	}

	var (
		idx int
		s   string
		b   bool
		sp  bw.ValPath
		ps  PosInfo
	)

LOOP:
	for {
		if p.Curr.Rune == '.' && len(result) == 0 {
			if len(a.Bases) == 0 {
				p.Forward(true)
				break LOOP
			} else if len(result) == 0 {
				result = append(result, a.Bases[0]...)
			} else {
				err = p.Unexpected(p.Curr)
				return
			}
		} else if idx, _, b, err = p.Int(); b {
			if err != nil {
				return
			}
			result = append(
				result,
				bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: idx},
			)
		} else if s, _, b, err = p.Id(); b {
			if err != nil {
				return
			}
			result = append(
				result,
				bw.ValPathItem{Type: bw.ValPathItemKey, Key: s},
			)
		} else if sp, _, b, err = p.subPath(a); b {
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
			p.Forward(true)
			break LOOP
		} else {
			if len(result) == 0 {
				if p.Curr.Rune == '$' {
					p.Forward(true)
					if err = p.CheckNotEOF(); err != nil {
						return
					}
					if b, err = p.processOn([]on{
						onInt{f: func(idx int, start PosInfo) {
							result = append(
								result,
								bw.ValPathItem{Type: bw.ValPathItemVar, Key: s},
							)
						}},
						onId{f: func(s string, start PosInfo) {
							var nidx int
							l := len(a.Bases)
							if nidx, b = bw.NormalIdx(idx, l); !b {
								err = p.Unexpected(ps, bw.Fmt(ansi.String("unexpected base path idx <ansiVal>%d<ansi> (len(bases): <ansiVal>%d)"), idx, l))
								return
							}
							result = append(result, a.Bases[nidx]...)
						}},
					}); err != nil {
						return
					}
					if b {
						goto CONTINUE
					}
					// var gotoCONTINUE bool
					// if s, _, b, err = p.Id(); b {
					// 	if err != nil {
					// 		return
					// 	}
					// 	result = append(
					// 		result,
					// 		bw.ValPathItem{Type: bw.ValPathItemVar, Key: s},
					// 	)
					// 	gotoCONTINUE = true
					// } else if idx, ps, b, err = p.Int(); b {
					// 	if err != nil {
					// 		return
					// 	}
					// 	var nidx int
					// 	l := len(a.Bases)
					// 	if nidx, b = bw.NormalIdx(idx, l); !b {
					// 		err = p.Unexpected(ps, bw.Fmt(ansi.String("unexpected base path idx <ansiVal>%d<ansi> (len(bases): <ansiVal>%d)"), idx, l))
					// 		return
					// 	}
					// 	result = append(result, a.Bases[nidx]...)
					// 	gotoCONTINUE = true
					// }
					// if gotoCONTINUE {
					// 	goto CONTINUE
					// }
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

		if !a.isSubPath && p.Curr.Rune == '?' {
			result[len(result)-1].IsOptional = true
			p.Forward(true)
		}

		if p.Curr.Rune != '.' {
			break LOOP
		}

		p.Forward(true)
	}
	return
}

// ============================================================================
