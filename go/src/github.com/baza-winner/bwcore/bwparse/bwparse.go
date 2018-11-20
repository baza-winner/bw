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

const Initial uint = 0

func (p *P) Forward(count uint) {
	if p.Curr.Pos < 0 || count > 0 && !p.Curr.IsEOF {
		if count <= 1 {
			p.forward()
		} else {
			for ; count > 0; count-- {
				p.forward()
			}
		}
	}
}

func (p *P) CheckNotEOF() (err error) {
	if p.Curr.IsEOF {
		err = p.Unexpected(p.Curr)
	}
	return
}

func (p *P) LookAhead(ofs uint) (ps PosInfo) {
	ps = p.Curr
	if ofs > 0 {
		idx := len(p.Next) - int(ofs)
		if idx >= 0 {
			ps = p.Next[idx]
		} else {
			if len(p.Next) > 0 {
				ps = p.Next[0]
			}
			var lookahead, newNext []PosInfo
			for i := idx; i < 0 && !ps.IsEOF; i++ {
				p.pullRune(&ps)
				lookahead = append(lookahead, ps)
			}
			for i := -idx - 1; i >= 0; i-- {
				newNext = append(newNext, lookahead[i])
			}
			p.Next = append(newNext, p.Next...)
		}
	}
	return
}

func (p *P) Unexpected(ps PosInfo, optFmt ...bw.I) error {
	var msg string
	if ps.Pos < p.Curr.Pos {
		if len(optFmt) > 0 {
			msg = bw.Spew.Sprintf(optFmt[0].FmtString(), optFmt[0].FmtArgs()...)
		} else {
			msg = fmt.Sprintf(ansiUnexpectedWord, p.Curr.Prefix[ps.Pos-p.Curr.PrefixStart:])
		}
	} else if !p.Curr.IsEOF {
		msg = fmt.Sprintf(ansiUnexpectedChar, ps.Rune, ps.Rune)
	} else {
		msg = ansiUnexpectedEOF
	}
	return bwerr.From(msg + p.suffix(ps))
}

// ============================================================================

const (
	TillNonEOF bool = false
	TillEOF    bool = true
)

func (p *P) SkipSpace(tillEOF bool) (err error) {
	p.Forward(Initial)
	for !p.Curr.IsEOF && unicode.IsSpace(p.Curr.Rune) {
		p.Forward(1)
	}
	if !tillEOF && p.Curr.IsEOF || tillEOF && !p.Curr.IsEOF {
		err = p.Unexpected(p.Curr)
	}
	return
}

// ============================================================================

func (p *P) ArrayOfString() (result []string, start PosInfo, ok bool, err error) {
	start = p.start()
	if ok = p.Curr.Rune == '<'; !ok {
		if ok = p.canSkipRunes('q', 'w') && isPunctOrSymbol(p.LookAhead(2).Rune); !ok {
			return
		}
		p.Forward(2)
	}
	delimiter := p.Curr.Rune
	if r, b := Braces[delimiter]; b {
		delimiter = r
	}
	p.Forward(1)
	result = []string{}
	for err == nil {
		if err = p.SkipSpace(TillNonEOF); err == nil {
			r := p.Curr.Rune
			if r == delimiter {
				p.Forward(1)
				break
			}
			var s string
			for err == nil && !(unicode.IsSpace(r) || r == delimiter) {
				s += string(r)
				p.Forward(1)
				if err = p.CheckNotEOF(); err == nil {
					r = p.Curr.Rune
				}
			}
			result = append(result, s)
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

func (p *P) Id() (result string, start PosInfo, ok bool, err error) {
	start = p.start()
	r := p.Curr.Rune
	if ok = isLetter(r); ok {
		for isLetter(r) || unicode.IsDigit(r) {
			result += string(r)
			p.Forward(1)
			r = p.Curr.Rune
		}
	}
	return
}

// ============================================================================

func (p *P) String() (result string, start PosInfo, ok bool, err error) {
	start = p.start()
	delimiter := p.Curr.Rune
	if ok = p.skipRunes('"') || p.skipRunes('\''); ok {
		expectEscapedContent := false
		b := true
		for err == nil {
			r := p.Curr.Rune
			if !expectEscapedContent {
				if p.Curr.IsEOF {
					b = false
				} else if p.skipRunes(delimiter) {
					break
				} else if p.skipRunes('\\') {
					expectEscapedContent = true
					continue
				}
			} else if !(r == '"' || r == '\'' || r == '\\') {
				r, b = EscapeRunes[r]
				b = b && delimiter == '"'
			}
			if !b {
				err = p.Unexpected(p.Curr)
			} else {
				result += string(r)
				p.Forward(1)
			}
			expectEscapedContent = false
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
	var s string
	if s, start, ok, err = p.looksLikeNumber(); err == nil && ok {
		b := true
		for b {
			s, b = p.addDigit(p.Curr.Rune, s)
		}
		if result, err = parseInt(s); err != nil {
			err = p.Unexpected(start, bwerr.Err(err))
		}
	}
	return
}

// ============================================================================

const dotRune = '.'

func (p *P) Number() (result interface{}, start PosInfo, ok bool, err error) {
	var (
		s      string
		hasDot bool
		b      bool
	)
	if s, start, ok, err = p.looksLikeNumber(); err == nil && ok {
		for {
			if s, b = p.addDigit(p.Curr.Rune, s); !b {
				if hasDot || !p.skipRunes(dotRune) {
					break
				} else {
					s += string(dotRune)
					hasDot = true
				}
			}
		}
		if hasDot && !zeroAfterDotRegexp.MatchString(s) {
			var f float64
			if f, err = strconv.ParseFloat(s, 64); err == nil {
				result = f
			}
		} else {
			if pos := strings.LastIndex(s, string(dotRune)); pos >= 0 {
				s = s[:pos]
			}
			var i int
			if i, err = parseInt(s); err == nil {
				result = i
			}
		}
		if err != nil {
			err = p.Unexpected(start, bwerr.Err(err))
		}
	}
	return
}

var zeroAfterDotRegexp = regexp.MustCompile(`\.0+$`)

// ============================================================================

func (p *P) Array() (result []interface{}, start PosInfo, ok bool, err error) {
	start, ok, err = p.parseDelimitedOptionalCommaSeparated('[', ']', func() (err error) {
		var val interface{}
		if val, err = p.parseVal(); err == nil {
			switch t := val.(type) {
			case []string:
				for _, s := range t {
					result = append(result, s)
				}
			default:
				result = append(result, val)
			}
		}
		return
	})
	if ok && result == nil {
		result = []interface{}{}
	}
	return
}

func (p *P) Map() (result map[string]interface{}, start PosInfo, ok bool, err error) {
	start, ok, err = p.parseDelimitedOptionalCommaSeparated('{', '}', func() (err error) {
		var (
			key string
			b   bool
		)
		if b, err = p.processOn(
			onString{f: func(s string, pi PosInfo) (err error) { key = s; return }},
			onId{f: func(s string, pi PosInfo) (err error) { key = s; return }},
		); !b {
			err = p.Unexpected(p.Curr)
		} else if err == nil {
			if err = p.SkipSpace(TillNonEOF); err == nil {
				if p.skipRunes(':') || p.skipRunes('=', '>') {
					err = p.SkipSpace(TillNonEOF)
				}
				if err == nil {
					if result == nil {
						result = map[string]interface{}{}
					}
					result[key], err = p.parseVal()
				}
			}
		}
		return
	})
	if ok && result == nil {
		result = map[string]interface{}{}
	}
	return
}

func (p *P) Val() (result interface{}, start PosInfo, ok bool, err error) {
	start = p.start()
	ok, err = p.processOn(
		onArray{f: func(vals []interface{}, pi PosInfo) (err error) { result = vals; return }},
		onString{f: func(s string, pi PosInfo) (err error) { result = s; return }},
		onNumber{f: func(val interface{}, pi PosInfo) (err error) { result = val; return }},
		onPath{f: func(path bw.ValPath, pi PosInfo) (err error) { result = path; return }},
		onMap{f: func(m map[string]interface{}, pi PosInfo) (err error) { result = m; return }},
		onArrayOfString{f: func(ss []string, pi PosInfo) (err error) { result = ss; return }},
		onId{f: func(s string, pi PosInfo) (err error) {
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
				err = p.Unexpected(pi, bw.Fmt(ansiUnexpectedWord, s))
				return
			}
			return
		}},
	)
	return
}

// ============================================================================

func (p *P) Path(a PathA) (result bw.ValPath, start PosInfo, ok bool, err error) {
	start = p.start()
	if ok = p.Curr.Rune == '$'; ok {
		result, err = p.PathContent(a)
	} else if ok = p.skipRunes('{', '{'); ok {
		if err = p.SkipSpace(TillNonEOF); err == nil {
			if result, err = p.PathContent(a); err == nil {
				if err = p.SkipSpace(TillNonEOF); err == nil {
					if !p.skipRunes('}', '}') {
						err = p.Unexpected(p.Curr)
					}
				}
			}
		}
	}
	return
}

type PathA struct {
	Bases     []bw.ValPath
	isSubPath bool
}

func (p *P) PathContent(a PathA) (result bw.ValPath, err error) {
	p.Forward(Initial)

	var (
		vpi              bw.ValPathItem
		b, isEmptyResult bool
	)

	result = bw.ValPath{}
	for err == nil {
		isEmptyResult = len(result) == 0
		b = true
		if isEmptyResult && p.Curr.Rune == '.' {
			if len(a.Bases) > 0 {
				result = append(result, a.Bases[0]...)
			} else {
				p.Forward(1)
				break
			}
		} else if b, err = p.processOn(
			onInt{f: func(idx int, start PosInfo) (err error) {
				vpi = bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: idx}
				return
			}},
			onId{f: func(s string, start PosInfo) (err error) {
				vpi = bw.ValPathItem{Type: bw.ValPathItemKey, Key: s}
				return
			}},
			onSubPath{a: a, f: func(path bw.ValPath, start PosInfo) (err error) {
				vpi = bw.ValPathItem{Type: bw.ValPathItemPath, Path: path}
				return
			}},
		); b {
			result = append(result, vpi)
		} else if p.skipRunes('#') {
			result = append(result, bw.ValPathItem{Type: bw.ValPathItemHash})
			break
		} else if isEmptyResult && p.skipRunes('$') {
			b, err = p.processOn(
				onInt{f: func(idx int, start PosInfo) (err error) {
					l := len(a.Bases)
					if nidx, b := bw.NormalIdx(idx, l); b {
						result = append(result, a.Bases[nidx]...)
					} else {
						err = p.Unexpected(start, bw.Fmt(ansi.String("unexpected base path idx <ansiVal>%d<ansi> (len(bases): <ansiVal>%d)"), idx, l))
					}
					return
				}},
				onId{f: func(s string, start PosInfo) (err error) {
					result = append(result, bw.ValPathItem{Type: bw.ValPathItemVar, Key: s})
					return
				}},
			)
		} else {
			b = false
		}
		if err == nil && !b {
			err = p.Unexpected(p.Curr)
		}
		if err == nil {
			if !a.isSubPath && p.skipRunes('?') {
				result[len(result)-1].IsOptional = true
			}
			if !p.skipRunes('.') {
				break
			}
		}
	}
	return
}

// ============================================================================
