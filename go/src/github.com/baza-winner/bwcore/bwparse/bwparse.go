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

type ParseValKind uint8

const (
	ParseValNil ParseValKind = iota
	ParseValBool
	ParseValNumber
	ParseValInt
	ParseValUint
	ParseValRange
	ParseValString
	ParseValId
	ParseValMap
	ParseValArray
	ParseValArrayOfString
	ParseValPath
)

//go:generate bwsetter -type ParseValKind -test
//go:generate stringer -type ParseValKind -trimprefix ParseVal

const (
	_ParseValKindSetTestItemA = ParseValBool
	_ParseValKindSetTestItemB = ParseValInt
)

// ============================================================================

type I interface {
	Curr() *PosInfo
	Forward(count uint)
	UnexpectedA(a UnexpectedA) error
	LookAhead(ofs uint) *PosInfo
}

// ============================================================================

type IdFunc func(p I, s string, start *PosInfo) (result interface{}, ok bool, err error)

type ValidateMapKeyFunc func(p I, m map[string]interface{}, key string, start *PosInfo) (err error)
type ParseMapElemFunc func(p I, m map[string]interface{}, key string) (ok bool, err error)
type MapEndFunc func(p I, m map[string]interface{}) (err error)

type ParseArrayElemFunc func(p I, vals []interface{}) (ok bool, err error)
type ArrayEndFunc func(p I, vals []interface{}) (err error)

type ValidateNumberFunc func(p I, n bwtype.Number, start *PosInfo) (err error)

type ValidateStringFunc func(p I, s string, start *PosInfo) (err error)

type ValidateArrayOfStringElemFunc func(p I, ss []string, s string, start *PosInfo) (err error)
type ArrayOfStringEndFunc func(p I, ss []string) (err error)

// ============================================================================

type Opt struct {
	ExcludeKinds bool
	KindSet      ParseValKindSet

	IdVals            map[string]interface{}
	OnId              IdFunc
	NonNegativeNumber func(opt ...bwtype.RangeLimitKind) bool

	IdNil   bwset.String
	IdFalse bwset.String
	IdTrue  bwset.String

	// NonNegativeRangeMin bool

	// OnMapBegin       MapEndFunc
	OnValidateMapKey ValidateMapKeyFunc
	OnParseMapElem   ParseMapElemFunc
	OnMapEnd         MapEndFunc

	// OnArrayBegin     ArrayEndFunc
	OnParseArrayElem ParseArrayElemFunc
	OnArrayEnd       ArrayEndFunc

	OnValidateString            ValidateStringFunc
	OnValidateArrayOfStringElem ValidateArrayOfStringElemFunc
	OnArrayOfStringEnd          ArrayOfStringEndFunc

	OnValidateNumber ValidateNumberFunc
}

// ============================================================================

type P struct {
	prov          bwrune.Provider
	curr          *PosInfo
	next          []*PosInfo
	preLineCount  uint
	postLineCount uint

	// IdVals map[string]interface{}

	// OnMapBegin       MapEndFunc
	// OnValidateMapKey ValidateMapKeyFunc
	// OnParseMapElem   ParseMapElemFunc
	// OnMapEnd         MapEndFunc

	// OnArrayBegin     ArrayEndFunc
	// OnParseArrayElem ParseArrayElemFunc
	// OnArrayEnd       ArrayEndFunc

	// OnValidateArrayOfStringElem ValidateArrayOfStringElemFunc
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
			// if f, ok := optKeyValidateMapKeyFunc(m, "OnValidateMapKey", &keys); ok {
			// 	result.OnValidateMapKey = f
			// }
			// if f, ok := optKeyParseMapElemFunc(m, "OnParseMapElem", &keys); ok {
			// 	result.OnParseMapElem = f
			// }
			// if f, ok := optKeyOnMapBeginEndFunc(m, "OnMapEnd", &keys); ok {
			// 	result.OnMapEnd = f
			// }

			// if f, ok := optKeyOnArrayBeginEndFunc(m, "OnArrayBegin", &keys); ok {
			// 	result.OnArrayBegin = f
			// }
			// if f, ok := optKeyParseArrayElemFunc(m, "OnParseArrayElem", &keys); ok {
			// 	result.OnParseArrayElem = f
			// }
			// if f, ok := optKeyOnArrayBeginEndFunc(m, "OnArrayEnd", &keys); ok {
			// 	result.OnArrayEnd = f
			// }

			// if f, ok := optKeyValidateArrayOfStringElemFunc(m, "OnValidateArrayOfStringElem", &keys); ok {
			// 	result.OnValidateArrayOfStringElem = f
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

// func (p *P) CheckNotEOF() (err error) {
// 	return CheckNotEOF(p)
// }

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

// ============================================================================

func Unexpected(p I, optPosInfo ...*PosInfo) error {
	var a UnexpectedA
	if len(optPosInfo) > 0 {
		a.PosInfo = optPosInfo[0]
	}
	return p.UnexpectedA(a)
}

// ============================================================================

func CheckNotEOF(p I) (err error) {
	if p.Curr().isEOF {
		err = Unexpected(p)
	}
	return
}

// ============================================================================

func CanSkipRunes(p I, rr ...rune) bool {
	for i, r := range rr {
		if pi := p.LookAhead(uint(i)); pi.isEOF || pi.rune != r {
			return false
		}
	}
	return true
}

func SkipRunes(p I, rr ...rune) (ok bool) {
	if ok = CanSkipRunes(p, rr...); ok {
		p.Forward(uint(len(rr)))
	}
	return
}

// ============================================================================

func IsDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func IsLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func IsPunctOrSymbol(r rune) bool {
	return unicode.IsPunct(r) || unicode.IsSymbol(r)
}

// ============================================================================

const (
	TillNonEOF bool = false
	TillEOF    bool = true
)

func SkipSpace(p I, tillEOF bool) (err error) {
	p.Forward(Initial)
REDO:
	for !p.Curr().isEOF && unicode.IsSpace(p.Curr().rune) {
		p.Forward(1)
	}
	if p.Curr().isEOF && !tillEOF {
		err = Unexpected(p)
		return
	}
	if CanSkipRunes(p, '/', '/') {
		p.Forward(2)
		for !p.Curr().isEOF && p.Curr().rune != '\n' {
			p.Forward(1)
		}
		if !p.Curr().isEOF {
			p.Forward(1)
		}
		goto REDO
	} else if CanSkipRunes(p, '/', '*') {
		p.Forward(2)
		for !p.Curr().isEOF && !CanSkipRunes(p, '*', '/') {
			p.Forward(1)
		}
		if !p.Curr().isEOF {
			p.Forward(2)
		}
		goto REDO
	}
	if tillEOF && !p.Curr().isEOF {
		err = Unexpected(p)
	}
	return
}

// ============================================================================

func Id(p I, optOpt ...Opt) (result string, start *PosInfo, ok bool, err error) {
	start = getStart(p)
	r := p.Curr().rune
	if ok = IsLetter(r); ok {
		for IsLetter(r) || unicode.IsDigit(r) {
			result += string(r)
			p.Forward(1)
			r = p.Curr().rune
		}
	}
	return
}

// ============================================================================

func String(p I, optOpt ...Opt) (result string, start *PosInfo, ok bool, err error) {
	start = getStart(p)
	delimiter := p.Curr().rune
	if ok = SkipRunes(p, '"') || SkipRunes(p, '\''); ok {
		expectEscapedContent := false
		b := true
		for err == nil {
			r := p.Curr().rune
			if !expectEscapedContent {
				if p.Curr().isEOF {
					b = false
				} else if SkipRunes(p, delimiter) {
					break
				} else if SkipRunes(p, '\\') {
					expectEscapedContent = true
					continue
				}
			} else if !(r == '"' || r == '\'' || r == '\\') {
				r, b = EscapeRunes[r]
				b = b && delimiter == '"'
			}
			if !b {
				err = Unexpected(p)
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

func Int(p I, optOpt ...Opt) (result int, start *PosInfo, ok bool, err error) {
	opt := getOpt(optOpt)
	var s string
	nonNegativeNumber := false
	if opt.NonNegativeNumber != nil {
		nonNegativeNumber = opt.NonNegativeNumber()
	}
	if s, start, _, ok, err = looksLikeNumber(p, nonNegativeNumber); err == nil && ok {
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

func Uint(p I, optOpt ...Opt) (result uint, start *PosInfo, ok bool, err error) {
	// opt := getOpt(optOpt)
	var s string
	if s, start, _, ok, err = looksLikeNumber(p, true); err == nil && ok {
		b := true
		for b {
			s, b = addDigit(p, s)
		}
		if result, err = bwstr.ParseUint(s); err != nil {
			err = p.UnexpectedA(UnexpectedA{start, bwerr.Err(err)})
		}
	}
	return
}

// ============================================================================

const dotRune = '.'

func Number(p I, optOpt ...Opt) (result bwtype.Number, start *PosInfo, ok bool, err error) {
	opt := getOpt(optOpt)
	var (
		s          string
		hasDot     bool
		b          bool
		isNegative bool
	)
	nonNegativeNumber := false
	if opt.NonNegativeNumber != nil {
		nonNegativeNumber = opt.NonNegativeNumber()
	}
	if s, start, isNegative, ok, err = looksLikeNumber(p, nonNegativeNumber); err == nil && ok {
		for {
			if s, b = addDigit(p, s); !b {
				if !hasDot && CanSkipRunes(p, dotRune) {
					// bwdebug.Print("!hi")
					pi := p.LookAhead(1)
					if IsDigit(pi.rune) {
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
		}
		if hasDot && !zeroAfterDotRegexp.MatchString(s) {
			var f float64
			if f, err = strconv.ParseFloat(s, 64); err == nil {
				result = bwtype.MustNumberFrom(f)
			}
		} else {
			if pos := strings.LastIndex(s, string(dotRune)); pos >= 0 {
				s = s[:pos]
			}
			if isNegative {
				var i int
				if i, err = bwstr.ParseInt(s); err == nil {
					result = bwtype.MustNumberFrom(i)
				}
			} else {
				var u uint
				if u, err = bwstr.ParseUint(s); err == nil {
					result = bwtype.MustNumberFrom(u)
				}
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

func ArrayOfString(p I, optOpt ...Opt) (result []string, start *PosInfo, ok bool, err error) {
	opt := getOpt(optOpt)
	start = getStart(p)
	if ok = p.Curr().rune == '<'; !ok {
		if ok = CanSkipRunes(p, 'q', 'w') && IsPunctOrSymbol(p.LookAhead(2).rune); !ok {
			return
		}
		p.Forward(2)
	}
	delimiter := p.Curr().rune
	if r, b := Braces[delimiter]; b {
		delimiter = r
	}
	p.Forward(1)
	result = []string{}
	for err == nil {
		if err = SkipSpace(p, TillNonEOF); err == nil {
			r := p.Curr().rune
			if r == delimiter {
				p.Forward(1)
				break
			}
			pi := p.Curr()
			var s string
			for err == nil && !(unicode.IsSpace(r) || r == delimiter) {
				s += string(r)
				p.Forward(1)
				if err = CheckNotEOF(p); err == nil {
					r = p.Curr().rune
				}
			}
			if err == nil {
				if opt.OnValidateArrayOfStringElem != nil {
					err = opt.OnValidateArrayOfStringElem(p, result, s, pi)
				} else if opt.OnValidateString != nil {
					err = opt.OnValidateString(p, s, pi)
				}
				if err == nil {
					result = append(result, s)
				}
			}
		}
	}
	if opt.OnArrayOfStringEnd != nil {
		err = opt.OnArrayOfStringEnd(p, result)
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

// func (p *P) Array(opt Opt) (result []interface{}, start *PosInfo, ok bool, err error) {
// 	return Array(p, opt)
// }

func Array(p I, optOpt ...Opt) (result []interface{}, start *PosInfo, ok bool, err error) {
	opt := getOpt(optOpt)
	if start, ok, err = parseDelimitedOptionalCommaSeparated(p, '[', ']', func() (err error) {
		if result == nil {
			result = []interface{}{}
			// if opt.OnArrayBegin != nil {
			// 	err = opt.OnArrayBegin(p, result)
			// }
		}
		if err == nil {
			var b bool
			if opt.OnParseArrayElem != nil {
				b, err = opt.OnParseArrayElem(p, result)
			}
			if err == nil && !b {
				var val interface{}
				if val, _, b, err = Val(p, opt); err == nil && !b {
					err = Unexpected(p)
				}
				if err == nil {
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
	}); ok {
		if result == nil {
			result = []interface{}{}
			// if opt.OnArrayBegin != nil {
			// 	err = opt.OnArrayBegin(p, result)
			// }
		}
		if opt.OnArrayEnd != nil {
			err = opt.OnArrayEnd(p, result)
		}
	}
	return
}

// func (p *P) Map() (result map[string]interface{}, start *PosInfo, ok bool, err error) {
// 	return parseMap(p)
// }

func Map(p I, optOpt ...Opt) (result map[string]interface{}, start *PosInfo, ok bool, err error) {
	opt := getOpt(optOpt)
	if start, ok, err = parseDelimitedOptionalCommaSeparated(p, '{', '}', func() (err error) {
		var (
			key string
			b   bool
		)
		onKey := func(s string, pi *PosInfo) (err error) {
			key = s
			if opt.OnValidateMapKey != nil {
				err = opt.OnValidateMapKey(p, result, key, pi)
			}
			return
		}
		if b, err = processOn(p, opt, onString{f: onKey}, onId{f: onKey}); !b {
			err = Unexpected(p)
		} else if err == nil {
			if err = SkipSpace(p, TillNonEOF); err == nil {
				if SkipRunes(p, ':') || SkipRunes(p, '=', '>') {
					err = SkipSpace(p, TillNonEOF)
				}
				if err == nil {
					if result == nil {
						result = map[string]interface{}{}
						// if opt.OnMapBegin != nil {
						// 	err = opt.OnMapBegin(p, result)
						// }
					}
					if err == nil {
						var b bool
						if opt.OnParseArrayElem != nil {
							b, err = opt.OnParseMapElem(p, result, key)
						}
						if err == nil && !b {
							if result[key], _, b, err = Val(p, opt); err == nil && !b {
								err = Unexpected(p)
							}
						}
					}
				}
			}
		}
		return
	}); ok {
		if result == nil {
			result = map[string]interface{}{}
			// if opt.OnMapBegin != nil {
			// 	err = opt.OnMapBegin(p, result)
			// }
		}
		if opt.OnMapEnd != nil {
			err = opt.OnMapEnd(p, result)
		}
	}
	return
}

func Nil(p I, optOpt ...Opt) (start *PosInfo, ok bool) {
	opt := getOpt(optOpt)
	start = getStart(p)
	p.Forward(Initial)

	ss := []string{"nil"}
	if len(opt.IdNil) > 0 {
		ss = append(ss, opt.IdNil.ToSliceOfStrings()...)
	}
	ok = isOneOfId(p, ss)
	return
}

func Bool(p I, optOpt ...Opt) (result bool, start *PosInfo, ok bool) {
	opt := getOpt(optOpt)
	start = getStart(p)
	p.Forward(Initial)

	ss := []string{"true"}
	if len(opt.IdTrue) > 0 {
		ss = append(ss, opt.IdTrue.ToSliceOfStrings()...)
	}
	if ok = isOneOfId(p, ss); ok {
		result = true
	} else {
		ss = []string{"false"}
		if len(opt.IdFalse) > 0 {
			ss = append(ss, opt.IdFalse.ToSliceOfStrings()...)
		}
		if ok = isOneOfId(p, ss); ok {
			return
		}
	}
	return
}

func Val(p I, optOpt ...Opt) (result interface{}, start *PosInfo, ok bool, err error) {
	opt := getOpt(optOpt)
	start = getStart(p)
	// onArgs := []on{}
	var onArgs []on
	kinds := []ParseValKind{}
	// addKind := func(kind ParseValKind) {
	// 	kinds = append(kinds, kind)
	// }
	kindSetIsEmpty := len(opt.KindSet) == 0
	hasKind := func(kind ParseValKind) (result bool) {
		kinds = append(kinds, kind)
		// addKind(kind)
		if kindSetIsEmpty {
			result = true
		} else if !opt.ExcludeKinds {
			result = opt.KindSet.Has(ParseValArray)
		} else if opt.ExcludeKinds {
			result = !opt.KindSet.Has(ParseValArray)
		}
		return
	}
	if hasKind(ParseValArray) {
		onArgs = append(onArgs, onArray{f: func(vals []interface{}, pi *PosInfo) (err error) { result = vals; return }})
	}
	if hasKind(ParseValString) {
		onArgs = append(onArgs, onString{f: func(s string, pi *PosInfo) (err error) {
			if opt.OnValidateString != nil {
				if err = opt.OnValidateString(p, s, start); err != nil {
					return
				}
			}
			result = s
			return
		}})
	}
	if hasKind(ParseValRange) {
		onArgs = append(onArgs, onRange{f: func(rng bwtype.Range, pi *PosInfo) (err error) { result = rng; return }})
	}

	if hasKind(ParseValNumber) {
		onArgs = append(onArgs, onNumber{f: func(n bwtype.Number, pi *PosInfo) (err error) {
			if opt.OnValidateNumber != nil {
				if err = opt.OnValidateNumber(p, n, start); err != nil {
					return
				}
			}
			result = n
			return
		}})
	} else if hasKind(ParseValInt) {
		onArgs = append(onArgs, onInt{f: func(i int, pi *PosInfo) (err error) {
			if opt.OnValidateNumber != nil {
				if err = opt.OnValidateNumber(p, bwtype.MustNumberFrom(i), start); err != nil {
					return
				}
			}
			result = i
			return
		}})
	} else if hasKind(ParseValUint) {
		onArgs = append(onArgs, onUint{f: func(u uint, pi *PosInfo) (err error) {
			if opt.OnValidateNumber != nil {
				if err = opt.OnValidateNumber(p, bwtype.MustNumberFrom(u), start); err != nil {
					return
				}
			}
			result = u
			return
		}})
	}

	if hasKind(ParseValPath) {
		onArgs = append(onArgs, onPath{f: func(path bw.ValPath, pi *PosInfo) (err error) { result = path; return }})
	}
	if hasKind(ParseValMap) {
		onArgs = append(onArgs, onMap{f: func(m map[string]interface{}, pi *PosInfo) (err error) { result = m; return }})
	}
	if hasKind(ParseValArrayOfString) {
		onArgs = append(onArgs, onArrayOfString{f: func(ss []string, pi *PosInfo) (err error) { result = ss; return }})
	}
	if hasKind(ParseValNil) {
		onArgs = append(onArgs, onNil{f: func(pi *PosInfo) (err error) { return }})
	}
	if hasKind(ParseValBool) {
		onArgs = append(onArgs, onBool{f: func(b bool, pi *PosInfo) (err error) { result = b; return }})
	}
	if len(opt.IdVals) > 0 || opt.OnId != nil {
		onArgs = append(onArgs,
			onId{f: func(s string, pi *PosInfo) (err error) {
				var b bool
				if result, b = opt.IdVals[s]; !b {
					if opt.OnId != nil {
						result, b, err = opt.OnId(p, s, pi)
					}
				}
				if !ok && err == nil {
					err = p.UnexpectedA(UnexpectedA{pi, bw.Fmt(ansiUnexpectedWord, s)})
				}
				return
			}},
		)
	}
	ok, err = processOn(p, opt, onArgs...)
	if !ok && err == nil {
		var expects []string
		addExpects := func(kind ParseValKind) {
			expects = append(expects, ansi.String("<ansiType>"+kind.String()))
		}
		for _, kind := range kinds {
			if len(opt.KindSet) == 0 || opt.KindSet.Has(ParseValArray) {
				addExpects(kind)
			}
		}
		if len(opt.IdVals) > 0 || opt.OnId != nil {
			bwerr.TODO()
		}
	}
	return
}

// ============================================================================

type proxy struct {
	p   I
	ofs uint
}

func (p *proxy) Curr() *PosInfo {
	result := p.p.LookAhead(p.ofs)
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

// ============================================================================

func Range(p I, opt Opt) (result bwtype.Range, start *PosInfo, ok bool, err error) {
	start = getStart(p)
	pp := &proxy{p: p}
	var min, max interface{}
	if min, _, ok, err = Number(pp, opt); err != nil {
		return
	} else if !ok {

	}
	if ok {
		start.justParsed = min
		start.justForward = pp.ofs
	}
	if ok = SkipRunes(pp, '.', '.'); !ok {
		return
	}
	p.Forward(pp.ofs)
	var b bool
	if max, _, b, err = Number(p, opt); err != nil {
		return
	} else if !b {

	}
	result = bwtype.MustRangeFrom(bwtype.A{Min: min, Max: max})

	// bwdebug.Print("!GERE")

	// p.Forward(Initial)
	// ok, err = p.processOn(
	// 	// onArray{f: func(vals []interface{}, pi *PosInfo) (err error) { result = vals; return }},
	// 	// onString{f: func(s string, pi *PosInfo) (err error) { result = s; return }},
	// 	onFloat64{f: func(val interface{}, pi *PosInfo) (err error) { result = val; return }},
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

// func (p *P) Path(a PathA) (result bw.ValPath, start *PosInfo, ok bool, err error) {
// 	return parsePath(p, a)
// 	// start = getStart(p)
// 	// if ok = p.curr.rune == '$'; ok {
// 	// 	result, err = p.PathContent(a)
// 	// } else if ok = SkipRunes(p, '{', '{'); ok {
// 	// 	if err = p.SkipSpace(TillNonEOF); err == nil {
// 	// 		if result, err = p.PathContent(a); err == nil {
// 	// 			if err = p.SkipSpace(TillNonEOF); err == nil {
// 	// 				if !SkipRunes(p, '}', '}') {
// 	// 					err = p.Unexpected()
// 	// 				}
// 	// 			}
// 	// 		}
// 	// 	}
// 	// }
// 	// return
// }

func Path(p I, a PathA, opt Opt) (result bw.ValPath, start *PosInfo, ok bool, err error) {
	start = getStart(p)
	if ok = p.Curr().rune == '$'; ok {
		result, err = PathContent(p, a, opt)
	} else if ok = SkipRunes(p, '{', '{'); ok {
		if err = SkipSpace(p, TillNonEOF); err == nil {
			if result, err = PathContent(p, a, opt); err == nil {
				if err = SkipSpace(p, TillNonEOF); err == nil {
					if !SkipRunes(p, '}', '}') {
						err = Unexpected(p)
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

func PathContent(p I, a PathA, optOpt ...Opt) (result bw.ValPath, err error) {
	opt := getOpt(optOpt)
	p.Forward(Initial)

	var (
		vpi              bw.ValPathItem
		b, isEmptyResult bool
	)

	result = bw.ValPath{}
	for err == nil {
		isEmptyResult = len(result) == 0
		b = true
		if isEmptyResult && p.Curr().rune == '.' {
			if len(a.Bases) > 0 {
				result = append(result, a.Bases[0]...)
			} else {
				p.Forward(1)
				break
			}
		} else if b, err = processOn(p, opt,
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
		} else if SkipRunes(p, '#') {
			result = append(result, bw.ValPathItem{Type: bw.ValPathItemHash})
			break
		} else if isEmptyResult && SkipRunes(p, '$') {
			b, err = processOn(p, opt,
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
			err = Unexpected(p)
		}
		if err == nil {
			if !a.isSubPath && SkipRunes(p, '?') {
				result[len(result)-1].IsOptional = true
			}
			if !SkipRunes(p, '.') {
				break
			}
		}
	}
	return
}

// func (p *P) PathContent(a PathA) (result bw.ValPath, err error) {
// 	return parsePathContent(p, a)
// 	// p.Forward(Initial)

// var (
// 	vpi              bw.ValPathItem
// 	b, isEmptyResult bool
// )

// result = bw.ValPath{}
// for err == nil {
// 	isEmptyResult = len(result) == 0
// 	b = true
// 	if isEmptyResult && p.curr.rune == '.' {
// 		if len(a.Bases) > 0 {
// 			result = append(result, a.Bases[0]...)
// 		} else {
// 			p.Forward(1)
// 			break
// 		}
// 	} else if b, err = p.processOn(
// 		onInt{f: func(idx int, start *PosInfo) (err error) {
// 			vpi = bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: idx}
// 			return
// 		}},
// 		onId{f: func(s string, start *PosInfo) (err error) {
// 			vpi = bw.ValPathItem{Type: bw.ValPathItemKey, Key: s}
// 			return
// 		}},
// 		onSubPath{a: a, f: func(path bw.ValPath, start *PosInfo) (err error) {
// 			vpi = bw.ValPathItem{Type: bw.ValPathItemPath, Path: path}
// 			return
// 		}},
// 	); b {
// 		result = append(result, vpi)
// 	} else if SkipRunes(p, '#') {
// 		result = append(result, bw.ValPathItem{Type: bw.ValPathItemHash})
// 		break
// 	} else if isEmptyResult && SkipRunes(p, '$') {
// 		b, err = p.processOn(
// 			onInt{f: func(idx int, start *PosInfo) (err error) {
// 				l := len(a.Bases)
// 				if nidx, b := bw.NormalIdx(idx, l); b {
// 					result = append(result, a.Bases[nidx]...)
// 				} else {
// 					err = p.UnexpectedA(UnexpectedA{start, bw.Fmt(ansi.String("Unexpected base path idx <ansiVal>%d<ansi> (len(bases): <ansiVal>%d)"), idx, l)})
// 				}
// 				return
// 			}},
// 			onId{f: func(s string, start *PosInfo) (err error) {
// 				result = append(result, bw.ValPathItem{Type: bw.ValPathItemVar, Key: s})
// 				return
// 			}},
// 		)
// 	} else {
// 		b = false
// 	}
// 	if err == nil && !b {
// 		err = p.Unexpected()
// 	}
// 	if err == nil {
// 		if !a.isSubPath && SkipRunes(p, '?') {
// 			result[len(result)-1].IsOptional = true
// 		}
// 		if !SkipRunes(p, '.') {
// 			break
// 		}
// 	}
// }
// return
// }

// ============================================================================
