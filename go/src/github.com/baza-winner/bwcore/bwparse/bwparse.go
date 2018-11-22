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
	"github.com/baza-winner/bwcore/bwtype"
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
	justParsed  interface{}
	justForward uint
}

func (p PosInfo) IsEOF() bool {
	return p.isEOF
}

func (p PosInfo) Rune() rune {
	return p.rune
}

// ============================================================================

type ValidateMapKeyFunc func(p *P, m map[string]interface{}, key string, start *PosInfo) (err error)
type ParseMapElemFunc func(p *P, m map[string]interface{}, key string) (ok bool, err error)
type OnMapBeginEndFunc func(p *P, m map[string]interface{}) (err error)

type ParseArrayElemFunc func(p *P, vals []interface{}) (ok bool, err error)
type OnArrayBeginEndFunc func(p *P, vals []interface{}) (err error)
type ValidateArrayOfStringElemFunc func(p *P, ss []string, s string, start *PosInfo) (err error)

type P struct {
	prov          bwrune.Provider
	curr          *PosInfo
	next          []*PosInfo
	preLineCount  uint
	postLineCount uint

	IdVals                    map[string]interface{}
	OnMapBegin                OnMapBeginEndFunc
	ValidateMapKey            ValidateMapKeyFunc
	ParseMapElem              ParseMapElemFunc
	OnMapEnd                  OnMapBeginEndFunc
	OnArrayBegin              OnArrayBeginEndFunc
	ParseArrayElem            ParseArrayElemFunc
	OnArrayEnd                OnArrayBeginEndFunc
	ValidateArrayOfStringElem ValidateArrayOfStringElemFunc
}

func From(p bwrune.Provider, opt ...map[string]interface{}) (result *P) {
	result = &P{
		prov:          p,
		curr:          &PosInfo{pos: -1, line: 1},
		next:          []*PosInfo{},
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
			// if m, ok := optKeyMap(m, "IdVals", &keys); ok {
			// 	result.IdVals = m
			// }

			// if f, ok := optKeyOnMapBeginEndFunc(m, "OnMapBegin", &keys); ok {
			// 	result.OnMapBegin = f
			// }
			// if f, ok := optKeyValidateMapKeyFunc(m, "ValidateMapKey", &keys); ok {
			// 	result.ValidateMapKey = f
			// }
			// if f, ok := optKeyParseMapElemFunc(m, "ParseMapElem", &keys); ok {
			// 	result.ParseMapElem = f
			// }
			// if f, ok := optKeyOnMapBeginEndFunc(m, "OnMapEnd", &keys); ok {
			// 	result.OnMapEnd = f
			// }

			// if f, ok := optKeyOnArrayBeginEndFunc(m, "OnArrayBegin", &keys); ok {
			// 	result.OnArrayBegin = f
			// }
			// if f, ok := optKeyParseArrayElemFunc(m, "ParseArrayElem", &keys); ok {
			// 	result.ParseArrayElem = f
			// }
			// if f, ok := optKeyOnArrayBeginEndFunc(m, "OnArrayEnd", &keys); ok {
			// 	result.OnArrayEnd = f
			// }

			// if f, ok := optKeyValidateArrayOfStringElemFunc(m, "ValidateArrayOfStringElem", &keys); ok {
			// 	result.ValidateArrayOfStringElem = f
			// }

			if unexpectedKeys := bwmap.MustUnexpectedKeys(m, keys); len(unexpectedKeys) > 0 {
				bwerr.Panic(ansiOptHasUnexpectedKeys, bwjson.Pretty(m), unexpectedKeys)
			}
		}
	}
	return
}

const Initial uint = 0

func (p *P) Curr() *PosInfo {
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

func (p *P) LookAhead(ofs uint) (result *PosInfo) {
	result = p.curr
	if ofs > 0 {
		idx := len(p.next) - int(ofs)
		if idx >= 0 {
			result = p.next[idx]
		} else {
			var ps PosInfo
			if len(p.next) > 0 {
				ps = *p.next[0]
			} else {
				ps = *p.curr
			}
			var lookahead []PosInfo
			for i := idx; i < 0 && !ps.isEOF; i++ {
				ps = p.pullRune(ps)
				lookahead = append(lookahead, ps)
			}
			var newNext []*PosInfo
			for i := -idx - 1; i >= 0; i-- {
				newNext = append(newNext, &lookahead[i])
			}
			p.next = append(newNext, p.next...)
			result = p.next[0]
		}
	}
	return
}

type UnexpectedA struct {
	PosInfo *PosInfo
	Fmt     bw.I
}

func (p *P) UnexpectedA(a UnexpectedA) error {
	var ps PosInfo
	// var zeroPosInfo = PosInfo{}
	if a.PosInfo == nil {
		ps = *p.curr
	} else {
		ps = *a.PosInfo
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

func (p *P) Unexpected(optPosInfo ...*PosInfo) error {
	return unexpected(p, optPosInfo...)
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
	if canSkipRunes(p, '/', '/') {
		p.Forward(2)
		for !p.curr.isEOF && p.curr.rune != '\n' {
			p.Forward(1)
		}
		if !p.curr.isEOF {
			p.Forward(1)
		}
		goto REDO
	} else if canSkipRunes(p, '/', '*') {
		p.Forward(2)
		for !p.curr.isEOF && !canSkipRunes(p, '*', '/') {
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

func (p *P) Id() (result string, start *PosInfo, ok bool, err error) {
	start = getStart(p)
	// bwdebug.Print("start:#v", *start)
	r := p.curr.rune
	if ok = isLetter(r); ok {
		for isLetter(r) || unicode.IsDigit(r) {
			result += string(r)
			p.Forward(1)
			r = p.curr.rune
		}
	}
	// bwdebug.Print("start:#v", *start)
	return
}

// ============================================================================

func (p *P) String() (result string, start *PosInfo, ok bool, err error) {
	start = getStart(p)
	delimiter := p.curr.rune
	if ok = skipRunes(p, '"') || skipRunes(p, '\''); ok {
		expectEscapedContent := false
		b := true
		for err == nil {
			r := p.curr.rune
			if !expectEscapedContent {
				if p.curr.isEOF {
					b = false
				} else if skipRunes(p, delimiter) {
					break
				} else if skipRunes(p, '\\') {
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

func (p *P) Int() (result int, start *PosInfo, ok bool, err error) {
	var s string
	if s, start, ok, err = looksLikeNumber(p); err == nil && ok {
		b := true
		for b {
			s, b = addDigit(p, s)
		}
		if result, err = bwstr.ParseInt(s); err != nil {
			err = p.UnexpectedA(UnexpectedA{start, bwerr.Err(err)})
		}
	}
	return
}

// ============================================================================

func (p *P) ArrayOfString() (result []string, start *PosInfo, ok bool, err error) {
	start = getStart(p)
	if ok = p.curr.rune == '<'; !ok {
		if ok = canSkipRunes(p, 'q', 'w') && isPunctOrSymbol(p.LookAhead(2).rune); !ok {
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
			pi := p.Curr()
			var s string
			for err == nil && !(unicode.IsSpace(r) || r == delimiter) {
				s += string(r)
				p.Forward(1)
				if err = p.CheckNotEOF(); err == nil {
					r = p.curr.rune
				}
			}
			if err == nil {
				if p.ValidateArrayOfStringElem != nil {
					err = p.ValidateArrayOfStringElem(p, result, s, pi)
				}
				if err == nil {
					result = append(result, s)
				}
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

func (p *P) Array() (result []interface{}, start *PosInfo, ok bool, err error) {
	start, ok, err = p.parseDelimitedOptionalCommaSeparated('[', ']', func() (err error) {
		if result == nil {
			result = []interface{}{}
			if p.OnArrayBegin != nil {
				err = p.OnArrayBegin(p, result)
			}
		}
		if err == nil {
			var b bool
			if p.ParseArrayElem != nil {
				b, err = p.ParseArrayElem(p, result)
			}
			if err == nil && !b {
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
			}
		}
		return
	})
	if ok {
		if result == nil {
			result = []interface{}{}
		} else if p.OnArrayEnd != nil {
			err = p.OnArrayEnd(p, result)
		}
	}
	return
}

func (p *P) Map() (result map[string]interface{}, start *PosInfo, ok bool, err error) {
	start, ok, err = p.parseDelimitedOptionalCommaSeparated('{', '}', func() (err error) {
		var (
			key string
			b   bool
		)
		onKey := func(s string, pi *PosInfo) (err error) {
			key = s
			if p.ValidateMapKey != nil {
				if err = p.ValidateMapKey(p, result, key, pi); err != nil {
					return
				}
			}
			return
		}
		if b, err = p.processOn(onString{f: onKey}, onId{f: onKey}); !b {
			err = p.Unexpected()
		} else if err == nil {
			if err = p.SkipSpace(TillNonEOF); err == nil {
				if skipRunes(p, ':') || skipRunes(p, '=', '>') {
					err = p.SkipSpace(TillNonEOF)
				}
				if err == nil {
					if result == nil {
						result = map[string]interface{}{}
						if p.OnMapBegin != nil {
							err = p.OnMapBegin(p, result)
						}
					}
					if err == nil {
						var b bool
						if p.ParseMapElem != nil {
							b, err = p.ParseMapElem(p, result, key)
						}
						if err == nil && !b {
							result[key], err = p.parseVal()
						}
					}
				}
			}
		}
		return
	})
	if ok {
		if result == nil {
			result = map[string]interface{}{}
		} else if p.OnMapEnd != nil {
			err = p.OnMapEnd(p, result)
		}
	}
	return
}

func (p *P) Nil() (start *PosInfo, ok bool) {
	start = getStart(p)
	p.Forward(Initial)
	if ok = canSkipRunes(p, 'n', 'i', 'l'); ok {
		p.Forward(3)
	} else if ok = canSkipRunes(p, 'n', 'u', 'l', 'l'); ok {
		p.Forward(4)
	}
	return
}

func (p *P) Bool() (result bool, start *PosInfo, ok bool) {
	start = getStart(p)
	p.Forward(Initial)
	if ok = canSkipRunes(p, 't', 'r', 'u', 'e'); ok {
		result = true
		p.Forward(4)
	} else if ok = canSkipRunes(p, 'f', 'a', 'l', 's', 'e'); ok {
		p.Forward(5)
	}
	return
}

func (p *P) Val() (result interface{}, start *PosInfo, ok bool, err error) {
	start = getStart(p)
	ok, err = p.processOn(
		onArray{f: func(vals []interface{}, pi *PosInfo) (err error) { result = vals; return }},
		onString{f: func(s string, pi *PosInfo) (err error) { result = s; return }},
		onRange{f: func(rng bwtype.Range, pi *PosInfo) (err error) { result = rng; return }},
		onNumber{f: func(val interface{}, pi *PosInfo) (err error) { result = val; return }},
		onPath{f: func(path bw.ValPath, pi *PosInfo) (err error) { result = path; return }},
		onMap{f: func(m map[string]interface{}, pi *PosInfo) (err error) { result = m; return }},
		onArrayOfString{f: func(ss []string, pi *PosInfo) (err error) { result = ss; return }},
		onNil{f: func(pi *PosInfo) (err error) { return }},
		onBool{f: func(b bool, pi *PosInfo) (err error) { result = b; return }},
		onId{f: func(s string, pi *PosInfo) (err error) {
			if val, ok := p.IdVals[s]; ok {
				result = val
			} else {
				// bwdebug.Print("pi:#v", pi)
				err = p.UnexpectedA(UnexpectedA{pi, bw.Fmt(ansiUnexpectedWord, s)})
			}
			return
		}},
	)
	return
}

// ============================================================================

const dotRune = '.'

func (p *P) Number() (result interface{}, start *PosInfo, ok bool, err error) {
	return parseNumber(p)
}

type intf interface {
	Curr() *PosInfo
	Forward(count uint)
	UnexpectedA(a UnexpectedA) error
	LookAhead(ofs uint) *PosInfo
}

func parseNumber(p intf) (result interface{}, start *PosInfo, ok bool, err error) {
	var (
		s      string
		hasDot bool
		b      bool
	)
	if s, start, ok, err = looksLikeNumber(p); err == nil && ok {
		for {
			if s, b = addDigit(p, s); !b {
				if !hasDot && canSkipRunes(p, dotRune) {
					// bwdebug.Print("!hi")
					pi := p.LookAhead(1)
					if isDigit(pi.rune) {
						p.Forward(1)
						s += string(dotRune)
						hasDot = true
					} else {
						break
					}
				} else {
					break
				}
			}
			// bwdebug.Print("p.Curr().rune:s", p.Curr().rune)
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

// type RangeA struct {
// }

type proxy struct {
	p   *P
	ofs uint
}

func (p *proxy) Curr() *PosInfo {
	result := p.p.LookAhead(p.ofs)
	// bwdebug.Print("proxy.Curr().rune:s", result.rune, "proxy.ofs", p.ofs)
	return result
}

func (p *proxy) Forward(count uint) {
	if count == 0 {
		p.p.Forward(0)
	} else {
		p.ofs += count
	}
}

func (p *proxy) LookAhead(ofs uint) *PosInfo {
	return p.p.LookAhead(p.ofs + ofs)
}

func (p *proxy) UnexpectedA(a UnexpectedA) error {
	p.p.Forward(p.ofs)
	return p.p.UnexpectedA(a)
}

func (p *P) Range() (result bwtype.Range, start *PosInfo, ok bool, err error) {
	start = getStart(p)
	pp := &proxy{p: p}
	var min, max interface{}
	// bwdebug.Print("!GERE")
	if min, _, ok, err = parseNumber(pp); !ok || err != nil {
		return
	}
	// bwdebug.Print("!GERE")
	start.justParsed = min
	start.justForward = pp.ofs
	if ok = canSkipRunes(pp, '.', '.'); !ok {
		return
	}
	p.Forward(pp.ofs + 2)
	if max, _, ok, err = parseNumber(p); err != nil {
		return
	}
	result = bwtype.MustRangeFrom(bwtype.A{Min: min, Max: max})

	// bwdebug.Print("!GERE")

	// p.Forward(Initial)
	// ok, err = p.processOn(
	// 	// onArray{f: func(vals []interface{}, pi *PosInfo) (err error) { result = vals; return }},
	// 	// onString{f: func(s string, pi *PosInfo) (err error) { result = s; return }},
	// 	onNumber{f: func(val interface{}, pi *PosInfo) (err error) { result = val; return }},
	// 	onPath{f: func(path bw.ValPath, pi *PosInfo) (err error) { result = path; return }},
	// 	// onMap{f: func(m map[string]interface{}, pi *PosInfo) (err error) { result = m; return }},
	// 	// onArrayOfString{f: func(ss []string, pi *PosInfo) (err error) { result = ss; return }},
	// 	// onNil{f: func(pi *PosInfo) (err error) { return }},
	// 	// onBool{f: func(b bool, pi *PosInfo) (err error) { result = b; return }},
	// 	// onId{f: func(s string, pi *PosInfo) (err error) {
	// 	// 	if val, ok := p.IdVals[s]; ok {
	// 	// 		result = val
	// 	// 	} else {
	// 	// 		err = p.UnexpectedA(UnexpectedA{pi, bw.Fmt(ansiUnexpectedWord, s)})
	// 	// 	}
	// 	// 	return
	// 	// }},
	// )
	return
}

// ============================================================================

func (p *P) Path(a PathA) (result bw.ValPath, start *PosInfo, ok bool, err error) {
	start = getStart(p)
	if ok = p.curr.rune == '$'; ok {
		result, err = p.PathContent(a)
	} else if ok = skipRunes(p, '{', '{'); ok {
		if err = p.SkipSpace(TillNonEOF); err == nil {
			if result, err = p.PathContent(a); err == nil {
				if err = p.SkipSpace(TillNonEOF); err == nil {
					if !skipRunes(p, '}', '}') {
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
			onInt{f: func(idx int, start *PosInfo) (err error) {
				vpi = bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: idx}
				return
			}},
			onId{f: func(s string, start *PosInfo) (err error) {
				vpi = bw.ValPathItem{Type: bw.ValPathItemKey, Key: s}
				return
			}},
			onSubPath{a: a, f: func(path bw.ValPath, start *PosInfo) (err error) {
				vpi = bw.ValPathItem{Type: bw.ValPathItemPath, Path: path}
				return
			}},
		); b {
			result = append(result, vpi)
		} else if skipRunes(p, '#') {
			result = append(result, bw.ValPathItem{Type: bw.ValPathItemHash})
			break
		} else if isEmptyResult && skipRunes(p, '$') {
			b, err = p.processOn(
				onInt{f: func(idx int, start *PosInfo) (err error) {
					l := len(a.Bases)
					if nidx, b := bw.NormalIdx(idx, l); b {
						result = append(result, a.Bases[nidx]...)
					} else {
						err = p.UnexpectedA(UnexpectedA{start, bw.Fmt(ansi.String("unexpected base path idx <ansiVal>%d<ansi> (len(bases): <ansiVal>%d)"), idx, l)})
					}
					return
				}},
				onId{f: func(s string, start *PosInfo) (err error) {
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
			if !a.isSubPath && skipRunes(p, '?') {
				result[len(result)-1].IsOptional = true
			}
			if !skipRunes(p, '.') {
				break
			}
		}
	}
	return
}

// ============================================================================
