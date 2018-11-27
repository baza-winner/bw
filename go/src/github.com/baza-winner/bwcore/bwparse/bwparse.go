package bwparse

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwdebug"
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

type Start struct {
	ps     *PosInfo
	suffix string
}

func (start Start) Suffix() string {
	return start.suffix
}

// ============================================================================

type ValKind uint8

const (
	ValNil ValKind = iota
	ValBool
	ValNumber
	ValInt
	ValUint
	ValRange
	ValString
	ValId
	ValMap
	ValArray
	ValPath
)

//go:generate bwsetter -type ValKind -test
//go:generate stringer -type ValKind -trimprefix Val

const (
	_ValKindSetTestItemA = ValBool
	_ValKindSetTestItemB = ValInt
)

// ============================================================================

type I interface {
	Curr() *PosInfo
	Forward(count uint)
	UnexpectedA(a UnexpectedA) error
	LookAhead(ofs uint) *PosInfo
	Start() *Start
	Stop(start *Start)
}

// ============================================================================

type On struct {
	P     I
	Start *Start
	Opt   *Opt
}

type IdFunc func(on On, s string) (val interface{}, ok bool, err error)

type ValidateMapKeyFunc func(on On, m map[string]interface{}, key string) (err error)
type ParseMapElemFunc func(on On, m map[string]interface{}, key string) (ok bool, err error)
type MapEndFunc func(on On, m map[string]interface{}) (err error)

type ParseArrayElemFunc func(on On, vals []interface{}) (outVals []interface{}, ok bool, err error)
type ArrayEndFunc func(on On, vals []interface{}) (err error)

type ValidateNumberFunc func(on On, n bwtype.Number) (err error)
type ValidateRangeFunc func(on On, rng bwtype.Range) (err error)

type ValidateStringFunc func(on On, s string) (err error)

type ValidateArrayOfStringElemFunc func(on On, ss []string, s string) (err error)
type ArrayOfStringEndFunc func(on On, ss []string) (err error)

// ============================================================================

type RangeLimitKind uint8

const (
	RangeLimitNone RangeLimitKind = iota
	RangeLimitMin
	RangeLimitMax
)

type Opt struct {
	ExcludeKinds bool
	KindSet      ValKindSet

	Base bw.ValPath

	path bw.ValPath

	IdVals            map[string]interface{}
	OnId              IdFunc
	NonNegativeNumber func(rangeLimitKind RangeLimitKind) bool

	IdNil   bwset.String
	IdFalse bwset.String
	IdTrue  bwset.String

	OnValidateMapKey ValidateMapKeyFunc
	OnParseMapElem   ParseMapElemFunc
	OnMapEnd         MapEndFunc

	OnParseArrayElem ParseArrayElemFunc
	OnArrayEnd       ArrayEndFunc

	OnValidateString            ValidateStringFunc
	OnValidateArrayOfStringElem ValidateArrayOfStringElemFunc
	OnArrayOfStringEnd          ArrayOfStringEndFunc

	OnValidateNumber ValidateNumberFunc
	OnValidateRange  ValidateRangeFunc
}

func (opt Opt) Path() bw.ValPath {
	return opt.path
}

// ============================================================================

type P struct {
	prov          bwrune.Provider
	curr          *PosInfo
	next          []*PosInfo
	preLineCount  uint
	postLineCount uint
	starts        map[int]*Start
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
			for i := len(lookahead) - 1; i >= 0; i-- {
				newNext = append(newNext, &lookahead[i])
			}
			p.next = append(newNext, p.next...)
			if len(p.next) > 0 {
				result = p.next[0]
			}
		}
	}
	return
}

func (p *P) Start() (result *Start) {
	p.Forward(Initial)
	var ok bool
	if result, ok = p.starts[p.curr.pos]; !ok {
		result = &Start{ps: p.curr}
		if p.starts == nil {
			p.starts = map[int]*Start{}
		}
		p.starts[p.curr.pos] = result
	}
	return
}

func (p *P) Stop(start *Start) {
	delete(p.starts, start.ps.pos)
}

type UnexpectedA struct {
	Start *Start
	Fmt   bw.I
}

func (p *P) UnexpectedA(a UnexpectedA) error {
	// var ps PosInfo
	var start Start
	if a.Start == nil {
		start = Start{ps: p.curr}
	} else {
		start = *a.Start
	}
	var msg string
	if start.ps.pos < p.curr.pos {
		if a.Fmt != nil {
			msg = bw.Spew.Sprintf(a.Fmt.FmtString(), a.Fmt.FmtArgs()...)
		} else {
			msg = fmt.Sprintf(ansiUnexpectedWord, start.suffix)
		}
	} else if !p.curr.isEOF {
		msg = fmt.Sprintf(ansiUnexpectedChar, p.curr.rune, p.curr.rune)
	} else {
		msg = ansiUnexpectedEOF
	}
	return bwerr.From(msg + p.suffix(start))
}

// ============================================================================

func Unexpected(p I, optStart ...*Start) error {
	var a UnexpectedA
	if len(optStart) > 0 {
		a.Start = optStart[0]
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

func Id(p I, optOpt ...Opt) (result string, start *Start, ok bool, err error) {
	start = p.Start()
	defer func() { p.Stop(start) }()
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

func String(p I, optOpt ...Opt) (result string, start *Start, ok bool, err error) {
	start = p.Start()
	defer func() { p.Stop(start) }()
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

func Int(p I, optOpt ...Opt) (result int, start *Start, ok bool, err error) {
	opt := getOpt(optOpt)
	nonNegativeNumber := false
	if opt.NonNegativeNumber != nil {
		nonNegativeNumber = opt.NonNegativeNumber(RangeLimitNone)
	}
	start = p.Start()
	defer func() { p.Stop(start) }()
	var justParsed numberResult
	// bwdebug.Print("start.ps.justParsed", start.ps.justParsed)
	if justParsed, ok = start.ps.justParsed.(numberResult); ok {
		if result, ok = bwtype.Int(justParsed.n.Val()); ok {
			if ok = !nonNegativeNumber || result >= 0; ok {
				p.Forward(start.ps.justForward)
			}
			return
		}
	}
	var s string
	if s, _, ok, err = looksLikeNumber(p, nonNegativeNumber); err == nil && ok {
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

func Uint(p I, optOpt ...Opt) (result uint, start *Start, ok bool, err error) {
	start = p.Start()
	defer func() { p.Stop(start) }()
	var justParsed numberResult
	if justParsed, ok = start.ps.justParsed.(numberResult); ok {
		if result, ok = bwtype.Uint(justParsed.n.Val()); ok {
			p.Forward(start.ps.justForward)
			return
		}
	}
	var s string
	if s, _, ok, err = looksLikeNumber(p, true); err == nil && ok {
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

func Number(p I, optOpt ...Opt) (result bwtype.Number, start *Start, ok bool, err error) {
	opt := getOpt(optOpt)
	return parseNumber(p, opt, RangeLimitNone)
}

// ============================================================================

func ArrayOfString(p I, optOpt ...Opt) (result []string, start *Start, ok bool, err error) {
	opt := getOpt(optOpt)
	return parseArrayOfString(p, opt, false)
}

// ============================================================================

func Array(p I, optOpt ...Opt) (result []interface{}, start *Start, ok bool, err error) {
	opt := getOpt(optOpt)
	start = p.Start()
	defer func() { p.Stop(start) }()
	on := On{p, start, &opt}
	base := opt.path
	on.Opt.path = append(base, bw.ValPathItem{Type: bw.ValPathItemIdx})
	if ok, err = parseDelimitedOptionalCommaSeparated(p, '[', ']', func() (err error) {
		if result == nil {
			result = []interface{}{}
		}
		if err == nil {
			var b bool
			var ss []string
			if ss, _, b, err = parseArrayOfString(p, opt, true); err == nil {
				if b {
					for _, s := range ss {
						result = append(result, s)
					}
				} else {
					if opt.OnParseArrayElem != nil {
						var newResult []interface{}
						if newResult, b, err = opt.OnParseArrayElem(on, result); b && err == nil {
							result = newResult
						}
					}
					if err == nil && !b {
						var val interface{}
						if val, _, err = Val(p, opt); err == nil {
							result = append(result, val)
						}
					}
				}
				// bwdebug.Print("len(result)", len(result))
				on.Opt.path[len(on.Opt.path)-1].Idx = len(result)
			}
		}
		return
	}); ok {
		if result == nil {
			result = []interface{}{}
		}
		on.Opt.path = base
		if opt.OnArrayEnd != nil {
			err = opt.OnArrayEnd(on, result)
		}
	}
	return
}

func Map(p I, optOpt ...Opt) (result map[string]interface{}, start *Start, ok bool, err error) {
	opt := getOpt(optOpt)
	start = p.Start()
	defer func() { p.Stop(start) }()
	on := On{p, start, &opt}
	base := opt.path
	if ok, err = parseDelimitedOptionalCommaSeparated(p, '{', '}', func() (err error) {
		if result == nil {
			result = map[string]interface{}{}
		}
		var (
			key string
			b   bool
		)
		path := append(base, bw.ValPathItem{Type: bw.ValPathItemKey})
		onKey := func(s string, start *Start) (err error) {
			key = s
			if opt.OnValidateMapKey != nil {
				on.Opt.path = path[:len(path)-1]
				err = opt.OnValidateMapKey(on, result, key)
			}
			return
		}
		if b, err = processOn(p,
			onString{opt: opt, f: onKey},
			onId{opt: opt, f: onKey},
		); !b {
			err = Unexpected(p)
		} else if err == nil {
			if err = SkipSpace(p, TillNonEOF); err == nil {
				if SkipRunes(p, ':') || SkipRunes(p, '=', '>') {
					err = SkipSpace(p, TillNonEOF)
				}
				if err == nil {
					if err == nil {
						var b bool
						path[len(path)-1].Key = key
						on.Opt.path = path
						if opt.OnParseArrayElem != nil {
							b, err = opt.OnParseMapElem(on, result, key)
						}
						// bwdebug.Print("result:json", result)
						if err == nil && !b {
							result[key], _, err = Val(p, opt)
						}
					}
				}
			}
		}
		return
	}); ok {
		if result == nil {
			result = map[string]interface{}{}
		}
		on.Opt.path = base
		if opt.OnMapEnd != nil {
			err = opt.OnMapEnd(on, result)
		}
	}
	return
}

func Nil(p I, optOpt ...Opt) (start *Start, ok bool) {
	opt := getOpt(optOpt)
	start = p.Start()
	defer func() { p.Stop(start) }()
	p.Forward(Initial)

	ss := []string{"nil"}
	if len(opt.IdNil) > 0 {
		ss = append(ss, opt.IdNil.ToSliceOfStrings()...)
	}
	ok = isOneOfId(p, ss)
	return
}

func Bool(p I, optOpt ...Opt) (result bool, start *Start, ok bool) {
	opt := getOpt(optOpt)
	start = p.Start()
	defer func() { p.Stop(start) }()
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

func Val(p I, optOpt ...Opt) (result interface{}, start *Start, err error) {
	opt := getOpt(optOpt)
	start = p.Start()
	defer func() { p.Stop(start) }()
	var onArgs []on
	kinds := []ValKind{}
	kindSetIsEmpty := len(opt.KindSet) == 0
	hasKind := func(kind ValKind) (result bool) {
		if kindSetIsEmpty {
			result = true
		} else if !opt.ExcludeKinds {
			result = opt.KindSet.Has(kind)
		} else if opt.ExcludeKinds {
			result = !opt.KindSet.Has(kind)
		}
		if result {
			kinds = append(kinds, kind)
		}
		// bwdebug.Print("result", result, "kind", kind, "kinds", kinds)
		return
	}
	if hasKind(ValArray) {
		onArgs = append(onArgs, onArray{opt: opt, f: func(vals []interface{}, start *Start) (err error) { result = vals; return }})
		onArgs = append(onArgs, onArrayOfString{opt: opt, f: func(ss []string, start *Start) (err error) { result = ss; return }})
	}
	if hasKind(ValString) {
		onArgs = append(onArgs, onString{opt: opt, f: func(s string, start *Start) (err error) {
			if opt.OnValidateString != nil {
				if err = opt.OnValidateString(On{p, start, &opt}, s); err != nil {
					return
				}
			}
			result = s
			return
		}})
	}
	if hasKind(ValRange) {
		onArgs = append(onArgs, onRange{opt: opt, f: func(rng bwtype.Range, start *Start) (err error) {
			if opt.OnValidateRange != nil {
				if err = opt.OnValidateRange(On{p, start, &opt}, rng); err != nil {
					return
				}
			}
			result = rng
			return
		}})
	}

	if hasKind(ValNumber) {
		onArgs = append(onArgs, onNumber{opt: opt, f: func(n bwtype.Number, start *Start) (err error) {
			if opt.OnValidateNumber != nil {
				if err = opt.OnValidateNumber(On{p, start, &opt}, n); err != nil {
					return
				}
			}
			result = n
			return
		}})
	} else if hasKind(ValInt) {
		onArgs = append(onArgs, onInt{opt: opt, f: func(i int, start *Start) (err error) {
			if opt.OnValidateNumber != nil {
				if err = opt.OnValidateNumber(On{p, start, &opt}, bwtype.MustNumberFrom(i)); err != nil {
					return
				}
			}
			result = i
			return
		}})
	} else if hasKind(ValUint) {
		onArgs = append(onArgs, onUint{opt: opt, f: func(u uint, start *Start) (err error) {
			if opt.OnValidateNumber != nil {
				if err = opt.OnValidateNumber(On{p, start, &opt}, bwtype.MustNumberFrom(u)); err != nil {
					return
				}
			}
			result = u
			return
		}})
	}
	if hasKind(ValPath) {
		onArgs = append(onArgs, onPath{opt: PathOpt{Opt: opt}, f: func(path bw.ValPath, start *Start) (err error) { result = path; return }})
	}
	if hasKind(ValMap) {
		onArgs = append(onArgs, onMap{opt: opt, f: func(m map[string]interface{}, start *Start) (err error) { result = m; return }})
	}
	if hasKind(ValNil) {
		onArgs = append(onArgs, onNil{opt: opt, f: func(start *Start) (err error) { return }})
	}
	if hasKind(ValBool) {
		onArgs = append(onArgs, onBool{opt: opt, f: func(b bool, start *Start) (err error) { result = b; return }})
	}
	if len(opt.IdVals) > 0 || opt.OnId != nil {
		onArgs = append(onArgs,
			onId{opt: opt, f: func(s string, start *Start) (err error) {
				var b bool
				if result, b = opt.IdVals[s]; !b {
					if opt.OnId != nil {
						result, b, err = opt.OnId(On{p, start, &opt}, s)
					}
				}
				if !b && err == nil {
					err = p.UnexpectedA(UnexpectedA{start, bw.Fmt(ansiUnexpectedWord, s)})
				}
				return
			}},
		)
	}
	var ok bool
	if ok, err = processOn(p, onArgs...); !ok && err == nil {
		var expects []string
		addExpects := func(kind ValKind) {
			s := kind.String()
			switch kind {
			case ValNumber, ValInt:
				if opt.NonNegativeNumber != nil && opt.NonNegativeNumber(RangeLimitNone) {
					s = "NonNegative" + s
				}
			case ValRange:
				if opt.NonNegativeNumber != nil && opt.NonNegativeNumber(RangeLimitMin) {
					s = s + "(Min: NonNegative)"
				}
			}
			expects = append(expects, ansi.String("<ansiType>"+s))
		}
		for _, kind := range kinds {
			addExpects(kind)
		}
		if len(opt.IdVals) > 0 || opt.OnId != nil {
			if len(opt.IdVals) == 0 {
				addExpects(ValId)
			} else {
				sset := bwset.String{}
				for s := range opt.IdVals {
					sset.Add(s)
				}
				s := ansi.String("<ansiType>%s<ansi>")
				if len(opt.IdVals) == 1 {
					s += fmt.Sprintf(ansi.String(" (<ansiVal>%s<ansi>")+ValId.String(), sset.ToSliceOfStrings()[0])
				} else {
					bytes, _ := json.Marshal(sset.ToSliceOfStrings())
					s += fmt.Sprintf(ansi.String(" (one of <ansiVal>%s<ansi>")+ValId.String(), string(bytes))
				}
				if opt.OnId != nil {
					s += ansi.String(" or <ansiVar>custom<ansi>)")
				}
				expects = append(expects, s)
			}
		}

		var what string
		if len(expects) <= 2 {
			what = strings.Join(expects, " or ")
		} else {
			what = "one of ["
			for _, s := range expects {
				what += "\n  " + s
			}
			what += "\n]"
		}
		err = bwerr.Refine(Unexpected(p), "expects %s instead of {Error}", what)
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

func (p *proxy) Start() *Start {
	return p.p.Start()
}

func (p *proxy) Stop(start *Start) {
	p.p.Stop(start)
}

// ============================================================================

type numberResult struct {
	n bwtype.Number
}

type pathResult struct {
	path bw.ValPath
}

// ============================================================================

func Range(p I, optOpt ...Opt) (result bwtype.Range, start *Start, ok bool, err error) {
	opt := getOpt(optOpt)
	start = p.Start()
	defer func() { p.Stop(start) }()
	pp := &proxy{p: p}
	var (
		min, max interface{}
		n        bwtype.Number
		isNumber bool
		isPath   bool
		path     bw.ValPath
	)

	if n, _, ok, err = parseNumber(pp, opt, RangeLimitMin); err != nil {
		return
	} else if ok {
		min = n
		isNumber = true
	} else if path, _, ok, err = Path(pp, PathOpt{Opt: opt}); err != nil {
		return
	} else if ok {
		min = path
		isPath = true
	}
	if ok = SkipRunes(pp, '.', '.'); !ok {
		if isNumber {
			start.ps.justParsed = numberResult{n}
			start.ps.justForward = pp.ofs
			bwdebug.Print("start.justParsed", start.ps.justParsed)
		} else if isPath {
			start.ps.justParsed = pathResult{path}
			start.ps.justForward = pp.ofs
		}
		return
	}
	p.Forward(pp.ofs)
	var b bool
	if max, _, b, err = parseNumber(p, opt, RangeLimitMax); err != nil {
		return
	} else if !b {
		if max, _, b, err = Path(p, PathOpt{Opt: opt}); err != nil {
			return
		} else if !b {
			max = nil
		}
	}
	result = bwtype.MustRangeFrom(bwtype.A{Min: min, Max: max})

	return
}

// ============================================================================

type PathA struct {
	Bases     []bw.ValPath
	isSubPath bool
}

type PathOpt struct {
	Opt
	Bases     []bw.ValPath
	isSubPath bool
}

func Path(p I, optOpt ...PathOpt) (result bw.ValPath, start *Start, ok bool, err error) {
	opt := getPathOpt(optOpt)
	start = p.Start()
	defer func() { p.Stop(start) }()
	var justParsed pathResult
	if justParsed, ok = start.ps.justParsed.(pathResult); ok {
		result = justParsed.path
		p.Forward(start.ps.justForward)
		return
	} else if ok = p.Curr().rune == '$'; ok {
		result, err = PathContent(p, opt)
	} else if ok = SkipRunes(p, '{', '{'); ok {
		if err = SkipSpace(p, TillNonEOF); err == nil {
			if result, err = PathContent(p, opt); err == nil {
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

func PathContent(p I, optOpt ...PathOpt) (result bw.ValPath, err error) {
	opt := getPathOpt(optOpt)
	p.Forward(Initial)
	var (
		vpi           bw.ValPathItem
		b             bool
		isEmptyResult bool
	)
	result = bw.ValPath{}
	for err == nil {
		isEmptyResult = len(result) == 0
		b = true
		if isEmptyResult && p.Curr().rune == '.' {
			if len(opt.Bases) > 0 {
				result = append(result, opt.Bases[0]...)
			} else {
				p.Forward(1)
				break
			}
		} else if b, err = processOn(p,
			onInt{opt: opt.Opt, f: func(idx int, start *Start) (err error) {
				vpi = bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: idx}
				return
			}},
			onId{opt: opt.Opt, f: func(s string, start *Start) (err error) {
				vpi = bw.ValPathItem{Type: bw.ValPathItemKey, Key: s}
				return
			}},
			onSubPath{opt: opt, f: func(path bw.ValPath, start *Start) (err error) {
				vpi = bw.ValPathItem{Type: bw.ValPathItemPath, Path: path}
				return
			}},
		); b {
			result = append(result, vpi)
		} else if SkipRunes(p, '#') {
			result = append(result, bw.ValPathItem{Type: bw.ValPathItemHash})
			break
		} else if isEmptyResult && SkipRunes(p, '$') {
			b, err = processOn(p,
				onInt{opt: opt.Opt, f: func(idx int, start *Start) (err error) {
					l := len(opt.Bases)
					if nidx, b := bw.NormalIdx(idx, l); b {
						result = append(result, opt.Bases[nidx]...)
					} else {
						err = p.UnexpectedA(UnexpectedA{start, bw.Fmt(ansi.String("unexpected base path idx <ansiVal>%d<ansi> (len(bases): <ansiVal>%d)"), idx, l)})
					}
					return
				}},
				onId{opt: opt.Opt, f: func(s string, start *Start) (err error) {
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
			if !opt.isSubPath && SkipRunes(p, '?') {
				result[len(result)-1].IsOptional = true
			}
			if CanSkipRunes(p, '.', '.') || !SkipRunes(p, '.') {
				break
			}
		}
	}
	return
}

// ============================================================================
