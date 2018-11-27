package bwparse

import (
	"encoding/json"
	"fmt"
	"unicode"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwdebug"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwerr/where"
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
	// justForward uint
}

func (p PosInfo) IsEOF() bool {
	return p.isEOF
}

func (p PosInfo) Rune() rune {
	return p.rune
}

type Start struct {
	ps      *PosInfo
	suffix  string
	stopped where.WW
}

func (start Start) Suffix() string {
	return start.suffix
}

func (p PosInfo) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{}
	result["isEOF"] = p.isEOF
	result["rune"] = string(p.rune)
	result["pos"] = p.pos
	result["line"] = p.line
	result["col"] = p.col
	result["prefix"] = p.prefix
	result["prefixStart"] = p.prefixStart
	if p.justParsed != nil {
		result["justParsed"] = p.justParsed
		// result["justForward"] = p.justForward
	}
	return json.Marshal(result)
}

func (start Start) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{}
	result["ps"] = *start.ps
	result["suffix"] = start.suffix
	if len(start.stopped) > 0 {
		result["stopped"] = start.stopped
	}
	return json.Marshal(result)
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
type ParseMapElemFunc func(on On, m map[string]interface{}, key string) (status Status)
type ValidateMapFunc func(on On, m map[string]interface{}) (err error)

type ParseArrayElemFunc func(on On, vals []interface{}) (outVals []interface{}, status Status)
type ValidateArrayFunc func(on On, vals []interface{}) (err error)

type ValidateNumberFunc func(on On, n bwtype.Number) (err error)
type ValidateRangeFunc func(on On, rng bwtype.Range) (err error)
type ValidatePathFunc func(on On, path bw.ValPath) (err error)

type ValidateStringFunc func(on On, s string) (err error)

type ValidateArrayOfStringElemFunc func(on On, ss []string, s string) (err error)
type ValidateArrayOfStringFunc func(on On, ss []string) (err error)

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
	OnValidateMap    ValidateMapFunc

	OnParseArrayElem ParseArrayElemFunc
	OnValidateArray  ValidateArrayFunc

	OnValidateString            ValidateStringFunc
	OnValidateArrayOfStringElem ValidateArrayOfStringElemFunc
	OnValidateArrayOfString     ValidateArrayOfStringFunc

	OnValidateNumber ValidateNumberFunc
	OnValidateRange  ValidateRangeFunc
	OnValidatePath   ValidatePathFunc
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
	if len(start.stopped) > 0 {
		return
	}
	start.stopped = where.WWFrom(1)
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

func (v Status) UnexpectedIfErr(p I) {
	if v.Err != nil {
		v.Err = p.UnexpectedA(UnexpectedA{v.Start, bwerr.Err(v.Err)})
	}
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
	nonNegativeNumber := false
	if opt.NonNegativeNumber != nil {
		nonNegativeNumber = opt.NonNegativeNumber(RangeLimitNone)
	}
	var justParsed numberResult
	curr := p.Curr()
	if justParsed, status.OK = curr.justParsed.(numberResult); status.OK {
		if result, status.OK = bwtype.Int(justParsed.n.Val()); status.OK {
			if status.OK = !nonNegativeNumber || result >= 0; status.OK {
				status.Start = justParsed.start
				// p.Forward(curr.justForward)
				p.Forward(uint(len(justParsed.start.suffix)))
			}
			return
		}
	}
	var s string
	if s, _, status = looksLikeNumber(p, nonNegativeNumber); status.IsOK() {
		defer func() { p.Stop(status.Start) }()
		b := true
		for b {
			s, b = addDigit(p, s)
		}
		result, status.Err = bwstr.ParseInt(s)
		status.UnexpectedIfErr(p)
	}
	return
}

// ============================================================================

func Uint(p I, optOpt ...Opt) (result uint, status Status) {
	var justParsed numberResult
	curr := p.Curr()
	if justParsed, status.OK = curr.justParsed.(numberResult); status.OK {
		if result, status.OK = bwtype.Uint(justParsed.n.Val()); status.OK {
			// p.Forward(curr.justForward)
			status.Start = justParsed.start
			// p.Forward(curr.justForward)
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

func ArrayOfString(p I, optOpt ...Opt) (result []string, status Status) {
	opt := getOpt(optOpt)
	return parseArrayOfString(p, opt, false)
}

// ============================================================================

func Array(p I, optOpt ...Opt) (result []interface{}, status Status) {
	opt := getOpt(optOpt)
	if status = parseDelimitedOptionalCommaSeparated(p, '[', ']', opt, func(on On, base bw.ValPath) (err error) {
		if result == nil {
			result = []interface{}{}
			on.Opt.path = append(base, bw.ValPathItem{Type: bw.ValPathItemIdx})
		}
		if err == nil {
			// var b bool
			var ss []string
			var st Status
			if ss, st = parseArrayOfString(p, opt, true); st.Err == nil {
				if st.OK {
					for _, s := range ss {
						result = append(result, s)
					}
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
				on.Opt.path[len(on.Opt.path)-1].Idx = len(result)
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

func Map(p I, optOpt ...Opt) (result map[string]interface{}, status Status) {
	opt := getOpt(optOpt)
	var path bw.ValPath
	if status = parseDelimitedOptionalCommaSeparated(p, '{', '}', opt, func(on On, base bw.ValPath) (err error) {
		if result == nil {
			result = map[string]interface{}{}
			path = append(base, bw.ValPathItem{Type: bw.ValPathItemKey})
		}
		var (
			key string
			// b   bool
		)
		onKey := func(s string, start *Start) (err error) {
			key = s
			if opt.OnValidateMapKey != nil {
				on.Opt.path = base
				on.Start = start
				err = opt.OnValidateMapKey(on, result, key)
			}
			return
		}
		var st Status
		if st = processOn(p,
			onString{opt: opt, f: onKey},
			onId{opt: opt, f: onKey},
		); !st.OK {
			err = Unexpected(p)
		} else if err != nil {
			err = st.Err
		} else {
			var isSpaceSkipped bool
			if isSpaceSkipped, err = SkipSpace(p, TillNonEOF); err == nil {
				if SkipRunes(p, ':') || SkipRunes(p, '=', '>') {
					isSpaceSkipped = true
					_, err = SkipSpace(p, TillNonEOF)
				}
				if err == nil && !isSpaceSkipped {
					err = Unexpected(p)
				}
				if err == nil {
					path[len(path)-1].Key = key
					on.Opt.path = path
					on.Start = p.Start()
					defer func() { p.Stop(on.Start) }()
					var st Status
					if opt.OnParseArrayElem != nil {
						st = opt.OnParseMapElem(on, result, key)
					}
					if st.Err == nil && !st.OK {
						result[key], st = Val(p, opt)
					}
					bwdebug.Print("key", key, "result", result)
					err = st.Err
				}
			}
		}
		return
	}); status.IsOK() {
		if result == nil {
			result = map[string]interface{}{}
		}
	}
	return
}

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

func Val(p I, optOpt ...Opt) (result interface{}, status Status) {
	opt := getOpt(optOpt)
	// start = p.Start()
	// defer func() { p.Stop(start) }()
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
		return
	}
	if hasKind(ValArray) {
		onArgs = append(onArgs, onArray{opt: opt, f: func(vals []interface{}, start *Start) (err error) {
			if opt.OnValidateArray != nil {
				if err = opt.OnValidateArray(On{p, start, &opt}, vals); err != nil {
					return
				}
			}
			result = vals
			return
		}})
		onArgs = append(onArgs, onArrayOfString{opt: opt, f: func(ss []string, start *Start) (err error) {
			if opt.OnValidateArrayOfString != nil {
				if err = opt.OnValidateArrayOfString(On{p, start, &opt}, ss); err != nil {
					return
				}
			}
			result = ss
			return
		}})
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
	if hasKind(ValPath) {
		onArgs = append(onArgs, onPath{opt: PathOpt{Opt: opt}, f: func(path bw.ValPath, start *Start) (err error) {
			if opt.OnValidatePath != nil {
				if err = opt.OnValidatePath(On{p, start, &opt}, path); err != nil {
					return
				}
			}
			result = path
			return
		}})
	}
	if hasKind(ValMap) {
		onArgs = append(onArgs, onMap{opt: opt, f: func(m map[string]interface{}, start *Start) (err error) {
			if opt.OnValidateMap != nil {
				if err = opt.OnValidateMap(On{p, start, &opt}, m); err != nil {
					return
				}
			}
			result = m
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
					if expects := getIdExpects(opt, ""); len(expects) > 0 {
						err = bwerr.Refine(err, "expects %s instead of {Error}", expects)
					}
				}
				return
			}},
		)
	}
	if status = processOn(p, onArgs...); !status.OK {
		var expects []string
		asType := func(kind ValKind) string {
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
			return ansi.String("<ansiType>" + s)
		}
		addExpects := func(kind ValKind) {
			expects = append(expects, asType(kind))
		}
		for _, kind := range kinds {
			addExpects(kind)
		}
		if len(opt.IdVals) > 0 || opt.OnId != nil {
			addExpects(ValId)
			s := asType(ValId)
			if expects := getIdExpects(opt, "  "); len(expects) > 0 {
				s += "(" + expects + ")"
			}
			expects = append(expects, s)
		}
		status.Err = bwerr.Refine(Unexpected(p), "expects %s instead of {Error}", bwstr.SmartJoin(bwstr.A{
			Source: bwstr.SS{
				SS: expects,
			},
			MaxLen:              80,
			NoJoinerOnMutliline: true,
		}))
	}
	return
}

// ============================================================================

type proxy struct {
	p      I
	ofs    uint
	starts map[int]*Start
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

func Range(p I, optOpt ...Opt) (result bwtype.Range, status Status) {
	opt := getOpt(optOpt)

	var (
		min, max interface{}
		n        bwtype.Number
		isNumber bool
		isPath   bool
		path     bw.ValPath
	)

	pp := &proxy{p: p}
	if n, status = parseNumber(pp, opt, RangeLimitMin); status.Err != nil {
		return
	} else if status.OK {
		min = n
		isNumber = true
	} else if path, status = Path(pp, PathOpt{Opt: opt}); status.Err != nil {
		return
	} else if status.OK {
		min = path
		isPath = true
	}
	if status.OK = CanSkipRunes(pp, '.', '.'); !status.OK {
		if isNumber || isPath {
			p.Stop(status.Start)
			ps := status.Start.ps
			if isNumber {
				ps.justParsed = numberResult{n, status.Start}
				// ps.justForward = pp.ofs
			} else if isPath {
				ps.justParsed = pathResult{path, status.Start}
				// ps.justForward = pp.ofs
			}
		}
		status = Status{}
		return
	} else if status.Start == nil {
		status.Start = p.Start()
	}
	defer func() { p.Stop(status.Start) }()

	p.Forward(pp.ofs)
	var st Status
	if max, st = parseNumber(p, opt, RangeLimitMax); st.Err != nil {
		status.Err = st.Err
		return
	} else if st.OK {
		p.Stop(st.Start)
	} else if max, st = parsePath(p, PathOpt{Opt: opt}); st.Err != nil {
		status.Err = st.Err
		return
	} else if st.OK {
		p.Stop(st.Start)
	} else {
		max = nil
	}

	result, status.Err = bwtype.RangeFrom(bwtype.A{Min: min, Max: max})

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

func PathContent(p I, optOpt ...PathOpt) (result bw.ValPath, err error) {
	opt := getPathOpt(optOpt)
	p.Forward(Initial)
	var (
		vpi bw.ValPathItem
		// b             bool
		isEmptyResult bool
	)
	result = bw.ValPath{}
	var st Status
	for st.Err == nil {
		isEmptyResult = len(result) == 0
		if isEmptyResult && p.Curr().rune == '.' {
			if len(opt.Bases) > 0 {
				result = append(result, opt.Bases[0]...)
			} else {
				p.Forward(1)
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
		); st.Err != nil {
			err = st.Err
		} else if st.OK {
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
			st.OK = false
		}
		if st.Err == nil && !st.OK {
			st.Err = Unexpected(p)
		}
		if st.Err == nil {
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
