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
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/bwstr"
)

// ============================================================================

type PosInfo struct {
	isEOF       bool
	rune        rune
	pos         int
	line        uint
	col         uint
	prefix      string
	prefixStart int
}

func (p PosInfo) IsEOF() bool {
	return p.isEOF
}

func (p PosInfo) Rune() rune {
	return p.rune
}

// ============================================================================

type OnMapKeyFunc func(p *P, m map[string]interface{}, key string) (ok bool, err error)
type OnMapEndFunc func(p *P, m map[string]interface{}) (err error)

// type OnArrayIdxFunc func(p *P, vals []interface{}, idx int) (ok bool, err error)
type OnArrayEndFunc func(p *P, arr interface{}) (err error)

type P struct {
	prov          bwrune.Provider
	curr          PosInfo
	next          []PosInfo
	preLineCount  uint
	postLineCount uint
	idVals        map[string]interface{}
	onMapKey      OnMapKeyFunc
	onMapEnd      OnMapEndFunc
	// onArrayIdx    OnArrayIdxFunc
	onArrayEnd OnArrayEndFunc
}

func From(p bwrune.Provider, opt ...map[string]interface{}) (result *P) {
	result = &P{
		prov:          p,
		curr:          PosInfo{pos: -1, line: 1},
		next:          []PosInfo{},
		preLineCount:  3,
		postLineCount: 3,
	}
	if len(opt) > 0 {
		m := opt[0]
		if m != nil {
			keys := bwset.String{}
			if i, ok := optKeyUint(m, "preLineCount", &keys); ok {
				result.preLineCount = i
			}
			if i, ok := optKeyUint(m, "postLineCount", &keys); ok {
				result.postLineCount = i
			}
			if m, ok := optKeyMap(m, "idVals", &keys); ok {
				result.idVals = m
			}
			if unexpectedKeys := bwmap.MustUnexpectedKeys(m, keys); len(unexpectedKeys) > 0 {
				bwerr.Panic(ansiOptHasUnexpectedKeys, bwjson.Pretty(m), unexpectedKeys)
			}
		}
	}
	return
}

const Initial uint = 0

func (p *P) Curr() PosInfo {
	return p.curr
}

func (p *P) Forward(count uint) {
	if p.curr.pos < 0 || count > 0 && !p.curr.isEOF {
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
	if p.curr.isEOF {
		err = p.Unexpected()
	}
	return
}

func (p *P) LookAhead(ofs uint) (ps PosInfo) {
	ps = p.curr
	if ofs > 0 {
		idx := len(p.next) - int(ofs)
		if idx >= 0 {
			ps = p.next[idx]
		} else {
			if len(p.next) > 0 {
				ps = p.next[0]
			}
			var lookahead, newNext []PosInfo
			for i := idx; i < 0 && !ps.isEOF; i++ {
				p.pullRune(&ps)
				lookahead = append(lookahead, ps)
			}
			for i := -idx - 1; i >= 0; i-- {
				newNext = append(newNext, lookahead[i])
			}
			p.next = append(newNext, p.next...)
		}
	}
	return
}

type UnexpectedA struct {
	PosInfo PosInfo
	Fmt     bw.I
}

func (p *P) UnexpectedA(a UnexpectedA) error {
	var ps PosInfo
	var zeroPosInfo = PosInfo{}
	if a.PosInfo == zeroPosInfo {
		ps = p.curr
	} else {
		ps = a.PosInfo
	}
	var msg string
	if ps.pos < p.curr.pos {
		if a.Fmt != nil {
			msg = bw.Spew.Sprintf(a.Fmt.FmtString(), a.Fmt.FmtArgs()...)
		} else {
			msg = fmt.Sprintf(ansiUnexpectedWord, p.curr.prefix[ps.pos-p.curr.prefixStart:])
		}
	} else if !p.curr.isEOF {
		msg = fmt.Sprintf(ansiUnexpectedChar, ps.rune, ps.rune)
	} else {
		msg = ansiUnexpectedEOF
	}
	return bwerr.From(msg + p.suffix(ps))
}

func (p *P) Unexpected(optPosInfo ...PosInfo) error {
	var a UnexpectedA
	if len(optPosInfo) > 0 {
		a.PosInfo = optPosInfo[0]
	}
	return p.UnexpectedA(a)
}

// ============================================================================

const (
	TillNonEOF bool = false
	TillEOF    bool = true
)

func (p *P) SkipSpace(tillEOF bool) (err error) {
	p.Forward(Initial)
REDO:
	for !p.curr.isEOF && unicode.IsSpace(p.curr.rune) {
		p.Forward(1)
	}
	if p.curr.isEOF && !tillEOF {
		err = p.Unexpected()
		return
	}
	if p.canSkipRunes('/', '/') {
		p.Forward(2)
		for !p.curr.isEOF && p.curr.rune != '\n' {
			p.Forward(1)
		}
		if !p.curr.isEOF {
			p.Forward(1)
		}
		goto REDO
	} else if p.canSkipRunes('/', '*') {
		p.Forward(2)
		for !p.curr.isEOF && !p.canSkipRunes('*', '/') {
			p.Forward(1)
		}
		if !p.curr.isEOF {
			p.Forward(2)
		}
		goto REDO
	}
	if tillEOF && !p.curr.isEOF {
		err = p.Unexpected()
	}
	return
}

// ============================================================================

func (p *P) ArrayOfString() (result []string, start PosInfo, ok bool, err error) {
	start = p.start()
	if ok = p.curr.rune == '<'; !ok {
		if ok = p.canSkipRunes('q', 'w') && isPunctOrSymbol(p.LookAhead(2).rune); !ok {
			return
		}
		p.Forward(2)
	}
	delimiter := p.curr.rune
	if r, b := Braces[delimiter]; b {
		delimiter = r
	}
	p.Forward(1)
	result = []string{}
	for err == nil {
		if err = p.SkipSpace(TillNonEOF); err == nil {
			r := p.curr.rune
			if r == delimiter {
				p.Forward(1)
				break
			}
			var s string
			for err == nil && !(unicode.IsSpace(r) || r == delimiter) {
				s += string(r)
				p.Forward(1)
				if err = p.CheckNotEOF(); err == nil {
					r = p.curr.rune
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
	r := p.curr.rune
	if ok = isLetter(r); ok {
		for isLetter(r) || unicode.IsDigit(r) {
			result += string(r)
			p.Forward(1)
			r = p.curr.rune
		}
	}
	return
}

// ============================================================================

func (p *P) String() (result string, start PosInfo, ok bool, err error) {
	start = p.start()
	delimiter := p.curr.rune
	if ok = p.skipRunes('"') || p.skipRunes('\''); ok {
		expectEscapedContent := false
		b := true
		for err == nil {
			r := p.curr.rune
			if !expectEscapedContent {
				if p.curr.isEOF {
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
				err = p.Unexpected()
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
			s, b = p.addDigit(p.curr.rune, s)
		}
		if result, err = bwstr.ParseInt(s); err != nil {
			err = p.UnexpectedA(UnexpectedA{start, bwerr.Err(err)})
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
			if s, b = p.addDigit(p.curr.rune, s); !b {
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
			if i, err = bwstr.ParseInt(s); err == nil {
				result = i
			}
		}
		if err != nil {
			err = p.UnexpectedA(UnexpectedA{start, bwerr.Err(err)})
		}
	}
	return
}

var zeroAfterDotRegexp = regexp.MustCompile(`\.0+$`)

// ============================================================================

func (p *P) Array() (result []interface{}, start PosInfo, ok bool, err error) {
	start, ok, err = p.parseDelimitedOptionalCommaSeparated('[', ']', func() (err error) {
		if result == nil {
			result = map[string]interface{}{}
		}
		p.parentArray = result
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
	if p.onArrayEnd != nil {
		err = p.onArrayEnd(p, result)
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
			err = p.Unexpected()
		} else if err == nil {
			if err = p.SkipSpace(TillNonEOF); err == nil {
				if p.skipRunes(':') || p.skipRunes('=', '>') {
					err = p.SkipSpace(TillNonEOF)
				}
				if err == nil {
					if result == nil {
						result = map[string]interface{}{}
					}
					p.parentArray = nil
					var b bool
					if p.onMapKey != nil {
						b, err = p.onMapKey(p, result, key)
					}
					if err == nil && !b {
						result[key], err = p.parseVal()
					}
				}
			}
		}
		return
	})
	if ok && result == nil {
		result = map[string]interface{}{}
	}
	if p.onMapEnd != nil {
		err = p.onMapEnd(p, result)
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
			default:
				if val, ok := p.idVals[s]; ok {
					result = val
				} else {
					err = p.UnexpectedA(UnexpectedA{pi, bw.Fmt(ansiUnexpectedWord, s)})
				}
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
	if ok = p.curr.rune == '$'; ok {
		result, err = p.PathContent(a)
	} else if ok = p.skipRunes('{', '{'); ok {
		if err = p.SkipSpace(TillNonEOF); err == nil {
			if result, err = p.PathContent(a); err == nil {
				if err = p.SkipSpace(TillNonEOF); err == nil {
					if !p.skipRunes('}', '}') {
						err = p.Unexpected()
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
		if isEmptyResult && p.curr.rune == '.' {
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
						err = p.UnexpectedA(UnexpectedA{start, bw.Fmt(ansi.String("unexpected base path idx <ansiVal>%d<ansi> (len(bases): <ansiVal>%d)"), idx, l)})
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
			err = p.Unexpected()
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
