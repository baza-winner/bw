package bwparse

import (
	"fmt"
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
}

func (p PosInfo) IsEOF() bool {
	return p.isEOF
}

func (p PosInfo) Rune() rune {
	return p.rune
}

// ============================================================================

func (start Start) Suffix() string {
	return start.suffix
}

// ============================================================================

type I interface {
	FileSpec() string
	Close() error
	Curr() *PosInfo
	Forward(count uint)
	Error(a E) error
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

type NonNegativeNumberFunc func(rangeLimitKind RangeLimitKind) bool
type IdFunc func(on On, s string) (val interface{}, ok bool, err error)

type ValidateMapKeyFunc func(on On, m map[string]interface{}, key string) (err error)
type ParseMapElemFunc func(on On, m map[string]interface{}, key string) (status Status)
type ValidateMapFunc func(on On, m map[string]interface{}) (err error)

type ParseArrayElemFunc func(on On, vals []interface{}) (outVals []interface{}, status Status)
type ValidateArrayFunc func(on On, vals []interface{}) (err error)

type ValidateNumberFunc func(on On, n bwtype.Number) (err error)
type ValidateRangeFunc func(on On, rng bwtype.Range) (err error)
type ValidatePathFunc func(on On, path bw.ValPath) (err error)

type ValidateStringFunc func(on On, s string) (err error)

type ValidateArrayOfStringElemFunc func(on On, ss []string, s string) (err error)

// ============================================================================

type RangeLimitKind uint8

const (
	RangeLimitNone RangeLimitKind = iota
	RangeLimitMin
	RangeLimitMax
)

type Opt struct {
	ExcludeKinds bool
	KindSet      bwtype.ValKindSet

	RangeLimitMinOwnKindSet bwtype.ValKindSet // empty means 'inherits'
	RangeLimitMaxOwnKindSet bwtype.ValKindSet // empty means 'inherits'

	Path bw.ValPath

	IdVals            map[string]interface{}
	OnId              IdFunc
	NonNegativeNumber NonNegativeNumberFunc

	IdNil   bwset.String
	IdFalse bwset.String
	IdTrue  bwset.String

	OnValidateMapKey ValidateMapKeyFunc
	OnParseMapElem   ParseMapElemFunc
	OnValidateMap    ValidateMapFunc

	OnParseArrayElem ParseArrayElemFunc
	OnValidateArray  ValidateArrayFunc

	OnValidateString            ValidateStringFunc
	OnValidateArrayOfStringElem ValidateArrayOfStringElemFunc
	// OnValidateArrayOfString     ValidateArrayOfStringFunc

	OnValidateNumber ValidateNumberFunc
	OnValidateRange  ValidateRangeFunc
	OnValidatePath   ValidatePathFunc

	ValFalse bool
}

// ============================================================================

type P struct {
	pp            bwrune.ProviderProvider
	prov          bwrune.Provider
	curr          *PosInfo
	next          []*PosInfo
	preLineCount  uint
	postLineCount uint
	starts        map[int]*Start
}

func From(pp bwrune.ProviderProvider, opt ...map[string]interface{}) (result *P) {
	result = &P{
		prov:          pp.Provider(),
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

func (p *P) FileSpec() string {
	return p.prov.FileSpec()
}

func (p *P) Close() error {
	return p.prov.Close()
}

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

type E struct {
	Start *Start
	Fmt   bw.I
}

func (p *P) Error(a E) error {
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
	return bwerr.From(msg + suffix(&proxy{p: p}, start, p.postLineCount))
}

// ============================================================================

func Unexpected(p I, optStart ...*Start) error {
	var a E
	if len(optStart) > 0 {
		a.Start = optStart[0]
	}
	return p.Error(a)
}

func Expects(p I, err error, what string) error {
	if err == nil {
		err = Unexpected(p)
	}
	return bwerr.Refine(err, "expects %s instead of {Error}", what)
}

func ExpectsSpace(p I) error {
	return Expects(p, nil, ansi.String(ansiVarSpace))
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

func SkipSpace(p I, tillEOF bool) (ok bool, err error) {
	p.Forward(Initial)
REDO:
	for !p.Curr().isEOF && unicode.IsSpace(p.Curr().rune) {
		ok = true
		p.Forward(1)
	}
	if p.Curr().isEOF && !tillEOF {
		err = Unexpected(p)
		return
	}
	if CanSkipRunes(p, '/', '/') {
		ok = true
		p.Forward(2)
		for !p.Curr().isEOF && p.Curr().rune != '\n' {
			p.Forward(1)
		}
		if !p.Curr().isEOF {
			p.Forward(1)
		}
		goto REDO
	} else if CanSkipRunes(p, '/', '*') {
		ok = true
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

type Status struct {
	Start *Start
	OK    bool
	Err   error
}

func (v Status) IsOK() bool {
	return v.OK && v.Err == nil
}

func (v Status) NoErr() bool {
	return v.Err == nil
}

func (v *Status) UnexpectedIfErr(p I) {
	if v.Err != nil {
		v.Err = p.Error(E{v.Start, bwerr.Err(v.Err)})
	}
}

// ============================================================================

func Nil(p I, optOpt ...Opt) (status Status) {
	opt := getOpt(optOpt)

	ss := []string{"nil"}
	if len(opt.IdNil) > 0 {
		ss = append(ss, opt.IdNil.ToSliceOfStrings()...)
	}

	var needForward uint
	if needForward, status.OK = isOneOfId(p, ss); status.OK {
		status.Start = p.Start()
		defer func() { p.Stop(status.Start) }()
		p.Forward(needForward)
	}
	return
}

// ============================================================================

func Bool(p I, optOpt ...Opt) (result bool, status Status) {
	opt := getOpt(optOpt)

	ss := []string{"true"}
	if len(opt.IdTrue) > 0 {
		ss = append(ss, opt.IdTrue.ToSliceOfStrings()...)
	}

	var needForward uint
	if needForward, status.OK = isOneOfId(p, ss); status.OK {
		result = true
	} else {
		ss = []string{"false"}
		if len(opt.IdFalse) > 0 {
			ss = append(ss, opt.IdFalse.ToSliceOfStrings()...)
		}
		needForward, status.OK = isOneOfId(p, ss)
	}
	if status.OK {
		status.Start = p.Start()
		defer func() { p.Stop(status.Start) }()
		p.Forward(needForward)
	}
	return
}

// ============================================================================

func Id(p I, optOpt ...Opt) (result string, status Status) {
	r := p.Curr().rune
	if status.OK = IsLetter(r); status.OK {
		status.Start = p.Start()
		defer func() { p.Stop(status.Start) }()
		for IsLetter(r) || unicode.IsDigit(r) {
			result += string(r)
			p.Forward(1)
			r = p.Curr().rune
		}
	}
	return
}

// ============================================================================

func String(p I, optOpt ...Opt) (result string, status Status) {
	delimiter := p.Curr().rune
	if status.OK = CanSkipRunes(p, '"') || CanSkipRunes(p, '\''); status.OK {
		status.Start = p.Start()
		defer func() { p.Stop(status.Start) }()
		p.Forward(1)
		expectEscapedContent := false
		b := true
		for status.NoErr() {
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
				status.Err = Unexpected(p)
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

func Int(p I, optOpt ...Opt) (result int, status Status) {
	opt := getOpt(optOpt)
	result, status = parseInt(p, opt, RangeLimitNone)
	if status.OK {
		p.Stop(status.Start)
	}
	return
}

// ============================================================================

func Uint(p I, optOpt ...Opt) (result uint, status Status) {
	var justParsed numberResult
	curr := p.Curr()
	if justParsed, status.OK = curr.justParsed.(numberResult); status.OK {
		if result, status.OK = bwtype.Uint(justParsed.n.Val()); status.OK {
			status.Start = justParsed.start
			p.Forward(uint(len(justParsed.start.suffix)))
			return
		}
	}
	var s string
	if s, _, status = looksLikeNumber(p, true); status.IsOK() {
		defer func() { p.Stop(status.Start) }()
		b := true
		for b {
			s, b = addDigit(p, s)
		}
		result, status.Err = bwstr.ParseUint(s)
		status.UnexpectedIfErr(p)
	}
	return
}

// ============================================================================

func Number(p I, optOpt ...Opt) (result bwtype.Number, status Status) {
	opt := getOpt(optOpt)
	result, status = parseNumber(p, opt, RangeLimitNone)
	if status.OK {
		p.Stop(status.Start)
	}
	return
}

// ============================================================================

// func ArrayOfString(p I, optOpt ...Opt) (result []string, status Status) {
// 	opt := getOpt(optOpt)
// 	return parseArrayOfString(p, opt, false)
// }

// ============================================================================

func Array(p I, optOpt ...Opt) (result []interface{}, status Status) {
	opt := getOpt(optOpt)
	if status = parseDelimitedOptionalCommaSeparated(p, '[', ']', opt, func(on On, base bw.ValPath) (err error) {
		if result == nil {
			result = []interface{}{}
			on.Opt.Path = append(base, bw.ValPathItem{Type: bw.ValPathItemIdx})
		}
		if err == nil {
			var vals []interface{}
			var st Status
			if vals, st = parseArrayOfString(p, *on.Opt, true); st.Err == nil {
				if st.OK {
					result = append(result, vals...)
				} else {
					if opt.OnParseArrayElem != nil {
						var newResult []interface{}
						if newResult, st = opt.OnParseArrayElem(on, result); st.IsOK() {
							result = newResult
						}
					}
					if st.Err == nil && !st.OK {
						var val interface{}
						if val, st = Val(p, opt); st.IsOK() {
							result = append(result, val)
						}
					}
				}
				on.Opt.Path[len(on.Opt.Path)-1].Idx = len(result)
			}
			err = st.Err
		}
		return
	}); status.IsOK() {
		if result == nil {
			result = []interface{}{}
		}
	}
	return
}

// ============================================================================

func Map(p I, optOpt ...Opt) (result map[string]interface{}, status Status) {
	opt := getOpt(optOpt)
	var path bw.ValPath
	if status = parseDelimitedOptionalCommaSeparated(p, '{', '}', opt, func(on On, base bw.ValPath) error {
		if result == nil {
			result = map[string]interface{}{}
			path = append(base, bw.ValPathItem{Type: bw.ValPathItemKey})
		}
		var (
			key string
		)
		onKey := func(s string, start *Start) (err error) {
			key = s
			if opt.OnValidateMapKey != nil {
				on.Opt.Path = base
				on.Start = start
				if _, b := result[key]; b {
					err = p.Error(E{
						Start: start,
						Fmt:   bw.Fmt(ansi.String("duplicate key <ansiErr>%s<ansi>"), start.Suffix()),
					})
				} else {
					err = opt.OnValidateMapKey(on, result, key)
				}
			}
			return
		}
		var st Status
		if st = processOn(p,
			onString{opt: opt, f: onKey},
			onId{opt: opt, f: onKey},
		); st.IsOK() {
			var isSpaceSkipped bool
			if isSpaceSkipped, st.Err = SkipSpace(p, TillNonEOF); st.Err == nil {

				var isSeparatorSkipped bool
				if SkipRunes(p, ':') {
					isSeparatorSkipped = true
				} else if SkipRunes(p, '=') {
					if SkipRunes(p, '>') {
						isSeparatorSkipped = true
					} else {
						return Expects(p, nil, fmt.Sprintf(ansiVal, string('>')))
					}
				}

				if st.Err == nil {
					if isSeparatorSkipped {
						_, st.Err = SkipSpace(p, TillNonEOF)
					}
					if !(isSpaceSkipped || isSeparatorSkipped) {
						st.Err = ExpectsSpace(p)
					} else {
						path[len(path)-1].Key = key
						on.Opt.Path = path
						on.Start = p.Start()
						defer func() { p.Stop(on.Start) }()

						st.OK = false
						if opt.OnParseMapElem != nil {
							st = opt.OnParseMapElem(on, result, key)
						}
						if st.Err == nil && !st.OK {
							result[key], st = Val(p, opt)
						}
					}
				}
			}
		}
		if !st.OK && st.Err == nil {
			st.Err = Unexpected(p)
		}
		return st.Err
	}); status.IsOK() {
		if result == nil {
			result = map[string]interface{}{}
		}
	}
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

func Path(p I, optOpt ...PathOpt) (result bw.ValPath, status Status) {
	opt := getPathOpt(optOpt)
	result, status = parsePath(p, opt)
	if status.OK {
		p.Stop(status.Start)
	}
	return
}

func MustPathFrom(s string, bases ...bw.ValPath) (result bw.ValPath) {
	var err error
	p := From(bwrune.S{s})
	if result, err = PathContent(p, PathOpt{Bases: bases}); err == nil {
		_, err = SkipSpace(p, TillEOF)
	}
	if err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func PathContent(p I, optOpt ...PathOpt) (bw.ValPath, error) {
	opt := getPathOpt(optOpt)
	p.Forward(Initial)
	var (
		vpi           bw.ValPathItem
		isEmptyResult bool
		st            Status
	)
	result := bw.ValPath{}
	for st.Err == nil {
		isEmptyResult = len(result) == 0
		if isEmptyResult && p.Curr().rune == '.' {
			p.Forward(1)
			if len(opt.Bases) > 0 {
				result = append(result, opt.Bases[0]...)
			} else {
				break
			}
		} else if st = processOn(p,
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
		); st.Err == nil {
			if st.OK {
				result = append(result, vpi)
			} else if SkipRunes(p, '#') {
				result = append(result, bw.ValPathItem{Type: bw.ValPathItemHash})
				break
			} else if isEmptyResult && SkipRunes(p, '$') {
				st = processOn(p,
					onInt{opt: opt.Opt, f: func(idx int, start *Start) (err error) {
						l := len(opt.Bases)
						if nidx, b := bw.NormalIdx(idx, l); b {
							result = append(result, opt.Bases[nidx]...)
						} else {
							err = p.Error(E{start, bw.Fmt(ansiUnexpectedBasePathIdx, idx, l)})
						}
						return
					}},
					onId{opt: opt.Opt, f: func(s string, start *Start) (err error) {
						result = append(result, bw.ValPathItem{Type: bw.ValPathItemVar, Key: s})
						return
					}},
				)
			}
			if st.Err == nil {
				if !st.OK {
					st.Err = Unexpected(p)
				} else {
					if !opt.isSubPath && SkipRunes(p, '?') {
						result[len(result)-1].IsOptional = true
					}
					if CanSkipRunes(p, '.', '.') || !SkipRunes(p, '.') {
						break
					}
				}
			}
		}
	}
	return result, st.Err
}

// ============================================================================

func Range(p I, optOpt ...Opt) (result bwtype.Range, status Status) {
	opt := getOpt(optOpt)

	var (
		min, rangeLimitVal interface{}
		isNumber           bool
		isPath             bool
		justParsedPath     bw.ValPath
	)

	hasKind := func(kind bwtype.ValKind, rlk RangeLimitKind) (result bool) {
		var ks bwtype.ValKindSet
		if rlk == RangeLimitMin {
			ks = opt.RangeLimitMinOwnKindSet
		} else {
			ks = opt.RangeLimitMaxOwnKindSet
		}
		if len(ks) != 0 {
			result = ks.Has(kind)
		} else if len(opt.KindSet) == 0 {
			result = true
		} else if !opt.ExcludeKinds {
			result = opt.KindSet.Has(kind)
		} else if opt.ExcludeKinds {
			result = !opt.KindSet.Has(kind)
		}
		return
	}

	onArgs := func(rlk RangeLimitKind) (onArgs []on) {
		if hasKind(bwtype.ValPath, rlk) {
			onArgs = append(onArgs, onPath{opt: PathOpt{Opt: opt}, f: func(path bw.ValPath, start *Start) (err error) {
				justParsedPath = path
				rangeLimitVal = path
				isPath = true
				return
			}})
		}

		if hasKind(bwtype.ValNumber, rlk) {
			onArgs = append(onArgs, onNumber{opt: opt, rlk: rlk, f: func(n bwtype.Number, start *Start) (err error) {
				rangeLimitVal = n
				isNumber = true
				return
			}})
		} else if hasKind(bwtype.ValInt, rlk) {
			onArgs = append(onArgs, onInt{opt: opt, rlk: rlk, f: func(i int, start *Start) (err error) {
				rangeLimitVal = i
				isNumber = true
				return
			}})
		} else if hasKind(bwtype.ValUint, rlk) {
			onArgs = append(onArgs, onUint{opt: opt, f: func(u uint, start *Start) (err error) {
				rangeLimitVal = u
				isNumber = true
				return
			}})
		}
		return
	}

	pp := &proxy{p: p}
	{
		if status = processOn(pp, onArgs(RangeLimitMin)...); status.Err != nil {
			if status.OK {
				pp.Stop(status.Start)
			}
			return
		}
		min = rangeLimitVal

		if status.OK = CanSkipRunes(pp, '.', '.'); !status.OK {
			if isNumber || isPath {
				pp.Stop(status.Start)
				ps := status.Start.ps
				if isNumber {
					ps.justParsed = numberResult{bwtype.MustNumberFrom(min), status.Start}
				} else if isPath {
					ps.justParsed = pathResult{justParsedPath, status.Start}
				}
			}
			status = Status{}
			return
		}
	}
	status.Start = p.Start()
	defer func() { p.Stop(status.Start) }()
	p.Forward(pp.ofs + 2)
	pp = nil

	rangeLimitVal = nil
	var st Status
	if st = processOn(p, onArgs(RangeLimitMax)...); st.OK {
		p.Stop(st.Start)
	}
	if st.Err != nil {
		status.Err = st.Err
	} else {
		result, status.Err = bwtype.RangeFrom(bwtype.A{Min: min, Max: rangeLimitVal})
	}

	return
}

// ============================================================================

func Val(p I, optOpt ...Opt) (result interface{}, status Status) {
	opt := getOpt(optOpt)
	var onArgs []on
	kinds := []bwtype.ValKind{}
	kindSetIsEmpty := len(opt.KindSet) == 0 //|| len(opt.KindSet) == 1 && opt.KindSet.has(ValUn)
	hasKind := func(kind bwtype.ValKind) (result bool) {
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
		return
	}
	if hasKind(bwtype.ValArray) {
		onArgs = append(onArgs, onArray{opt: opt, f: func(vals []interface{}, start *Start) (err error) {
			if opt.OnValidateArray != nil {
				err = opt.OnValidateArray(On{p, start, &opt}, vals)
			}
			if err == nil {
				result = vals
			}
			return
		}})
		onArgs = append(onArgs, onArrayOfString{opt: opt, f: func(vals []interface{}, start *Start) (err error) {
			if opt.OnValidateArray != nil {
				err = opt.OnValidateArray(On{p, start, &opt}, vals)
			}
			if err == nil {
				result = vals
			}
			return
		}})
	}
	if hasKind(bwtype.ValString) {
		onArgs = append(onArgs, onString{opt: opt, f: func(s string, start *Start) (err error) {
			if opt.OnValidateString != nil {
				err = opt.OnValidateString(On{p, start, &opt}, s)
			}
			if err == nil {
				result = s
			}
			return
		}})
	}
	if hasKind(bwtype.ValRange) {
		onArgs = append(onArgs, onRange{opt: opt, f: func(rng bwtype.Range, start *Start) (err error) {
			if opt.OnValidateRange != nil {
				err = opt.OnValidateRange(On{p, start, &opt}, rng)
			}
			if err == nil {
				result = rng
			}
			return
		}})
	}
	if hasKind(bwtype.ValPath) {
		onArgs = append(onArgs, onPath{opt: PathOpt{Opt: opt}, f: func(path bw.ValPath, start *Start) (err error) {
			if opt.OnValidatePath != nil {
				err = opt.OnValidatePath(On{p, start, &opt}, path)
			}
			if err == nil {
				result = path
			}
			return
		}})
	}
	if hasKind(bwtype.ValMap) {
		onArgs = append(onArgs, onMap{opt: opt, f: func(m map[string]interface{}, start *Start) (err error) {
			if opt.OnValidateMap != nil {
				err = opt.OnValidateMap(On{p, start, &opt}, m)
			}
			if err == nil {
				result = m
			}
			return
		}})
	}
	if hasKind(bwtype.ValNumber) {
		onArgs = append(onArgs, onNumber{opt: opt, f: func(n bwtype.Number, start *Start) (err error) {
			if opt.OnValidateNumber != nil {
				err = opt.OnValidateNumber(On{p, start, &opt}, n)
			}
			if err == nil {
				result = n
			}
			return
		}})
	} else if hasKind(bwtype.ValInt) {
		onArgs = append(onArgs, onInt{opt: opt, f: func(i int, start *Start) (err error) {
			if opt.OnValidateNumber != nil {
				err = opt.OnValidateNumber(On{p, start, &opt}, bwtype.MustNumberFrom(i))
			}
			if err == nil {
				result = i
			}
			return
		}})
	} else if hasKind(bwtype.ValUint) {
		onArgs = append(onArgs, onUint{opt: opt, f: func(u uint, start *Start) (err error) {
			if opt.OnValidateNumber != nil {
				err = opt.OnValidateNumber(On{p, start, &opt}, bwtype.MustNumberFrom(u))
			}
			if err == nil {
				result = u
			}
			return
		}})
	}
	if hasKind(bwtype.ValNil) {
		onArgs = append(onArgs, onNil{opt: opt, f: func(start *Start) (err error) { return }})
	}
	if hasKind(bwtype.ValBool) {
		onArgs = append(onArgs, onBool{opt: opt, f: func(b bool, start *Start) (err error) { result = b; return }})
	}
	if (len(opt.IdVals) > 0 || opt.OnId != nil) && hasKind(bwtype.ValId) {
		onArgs = append(onArgs,
			onId{opt: opt, f: func(s string, start *Start) (err error) {
				var b bool
				if result, b = opt.IdVals[s]; !b {
					if opt.OnId != nil {
						result, b, err = opt.OnId(On{p, start, &opt}, s)
					}
				}
				if !b && err == nil {
					err = p.Error(E{start, bw.Fmt(ansiUnexpectedWord, s)})
					if expects := getIdExpects(opt, ""); len(expects) > 0 {
						err = Expects(p, err, expects)
					}
				}
				return
			}},
		)
	}
	if status = processOn(p, onArgs...); !status.OK && !opt.ValFalse {
		var expects []string
		asType := func(kind bwtype.ValKind) (result string) {
			s := kind.String()
			switch kind {
			case bwtype.ValNumber, bwtype.ValInt:
				if opt.NonNegativeNumber != nil && opt.NonNegativeNumber(RangeLimitNone) {
					s = "NonNegative" + s
				}
			case bwtype.ValRange:
				if opt.NonNegativeNumber != nil && opt.NonNegativeNumber(RangeLimitMin) {
					s = s + "(Min: NonNegative)"
				}
			}
			result = fmt.Sprintf(ansiType, s)
			if kind == bwtype.ValId {
				if expects := getIdExpects(opt, "  "); len(expects) > 0 {
					result += "(" + expects + ")"
				}
			}
			return
		}
		addExpects := func(kind bwtype.ValKind) {
			expects = append(expects, asType(kind))
		}
		for _, kind := range kinds {
			addExpects(kind)
		}
		status.Err = Expects(p, status.Err,
			bwstr.SmartJoin(bwstr.A{
				Source: bwstr.SS{
					SS: expects,
				},
				MaxLen:              60,
				NoJoinerOnMutliline: true,
			}),
		)
	}
	return
}

// ============================================================================
