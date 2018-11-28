package bwparse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwdebug"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/bwstr"
	"github.com/baza-winner/bwcore/bwtype"
)

// ============================================================================

var (
	ansiOptKeyIsNotOfType    string
	ansiOptHasUnexpectedKeys string

	ansiOK  string
	ansiErr string

	ansiPos             string
	ansiLineCol         string
	ansiGetSuffixAssert string
	ansiUnexpectedEOF   string
	ansiUnexpectedChar  string
	ansiUnexpectedWord  string
	ansiOutOfRange      string
)

func init() {
	ansiOptKeyIsNotOfType = ansi.String("<ansiVar>opt.%s<ansi> (<ansiVal>%#v<ansi>) is not <ansiType>%s")
	ansiOptHasUnexpectedKeys = ansi.String("<ansiVar>opt<ansi> (<ansiVal>%s<ansi>) has unexpected keys <ansiVal>%s")

	ansiOK = ansi.CSIFromSGRCodes(ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen, Bright: false})).String()
	ansiErr = ansi.CSIFromSGRCodes(ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorRed, Bright: true})).String()

	ansiPos = ansi.String(" at pos <ansiPath>%d<ansi>")
	ansiLineCol = ansi.String(" at line <ansiPath>%d<ansi>, col <ansiPath>%d<ansi> (pos <ansiPath>%d<ansi>)")
	ansiGetSuffixAssert = ansi.String("<ansiVar>ps.pos<ansi> (<ansiVal>%d<ansi>) > <ansiVar>p.curr.pos<ansi> (<ansiVal>%d<ansi>)")
	ansiUnexpectedEOF = ansi.String("unexpected end of string")
	ansiUnexpectedChar = ansi.String("unexpected char <ansiVal>%q<ansiReset> (<ansiVar>charCode<ansi>: <ansiVal>%d<ansi>)")
	ansiUnexpectedWord = ansi.String("unexpected <ansiErr>`%s`<ansi>")
}

func optKeyUint(opt map[string]interface{}, key string, keys *bwset.String) (result uint, ok bool) {
	var val interface{}
	keys.Add(key)
	if val, ok = opt[key]; ok && val != nil {
		if result, ok = bwtype.Uint(val); !ok {
			bwerr.Panic(ansiOptKeyIsNotOfType, key, val, "Uint")
		}
	}
	return
}

// ============================================================================

func getOpt(optOpt []Opt) (result Opt) {
	if len(optOpt) > 0 {
		result = optOpt[0]
	}
	return
}

func getPathOpt(optOpt []PathOpt) (result PathOpt) {
	if len(optOpt) > 0 {
		result = optOpt[0]
	}
	return
}

func isOneOfId(p I, ss []string) (needForward uint, ok bool) {
	p.Forward(Initial)
	for _, s := range ss {
		if ok = CanSkipRunes(p, []rune(s)...); ok {
			u := uint(len(s))
			r := p.LookAhead(u).rune
			if ok = !(IsLetter(r) || IsDigit(r)); ok {
				needForward = u
				// p.Forward(u)
				return
			}
		}
	}
	return
}

// ============================================================================

func (p *P) pullRune(ps PosInfo) PosInfo {
	runePtr := bwrune.MustPull(p.prov)
	if !ps.isEOF {
		if ps.pos >= 0 {
			ps.prefix += string(ps.rune)
		}
		if runePtr != nil {
			if ps.rune != '\n' {
				ps.col++
			} else {
				ps.line++
				ps.col = 1
				if int(ps.line) > int(p.preLineCount)+1 {
					i := strings.Index(ps.prefix, "\n")
					ps.prefix = ps.prefix[i+1:]
					ps.prefixStart += i + 1
				}
			}
		}
		if runePtr == nil {
			ps.rune, ps.isEOF = '\000', true
		} else {
			ps.rune, ps.isEOF = *runePtr, false
		}
		ps.pos++
	}
	return ps
}

func (p *P) suffix(start Start) (suffix string) {
	if start.ps.pos > p.curr.pos {
		bwerr.Panic(ansiGetSuffixAssert, start.ps.pos, p.curr.pos)
	}

	postLineCount := int(p.postLineCount)

	var separator string
	if p.curr.line > 1 {
		suffix += fmt.Sprintf(ansiLineCol, start.ps.line, start.ps.col, start.ps.pos)
		separator = "\n"
	} else {
		suffix += fmt.Sprintf(ansiPos, start.ps.pos)
		separator = " "
	}
	suffix += ":" + separator + ansiOK + start.ps.prefix

	var needPostLines, noNeedNewline bool
	if start.ps.pos < p.curr.pos {
		suffix += ansiErr + start.suffix + ansi.Reset()
		needPostLines = true
	} else if !p.curr.isEOF {
		suffix += ansiErr + string(p.curr.rune) + ansi.Reset()
		p.Forward(1)
		needPostLines = true
	}
	noNeedNewline = p.curr.rune == '\n'

	for needPostLines && !p.curr.isEOF && postLineCount >= 0 {
		suffix += string(p.curr.rune)
		if noNeedNewline = p.curr.rune == '\n'; noNeedNewline {
			postLineCount -= 1
		}
		p.Forward(1)
	}

	if !noNeedNewline {
		suffix += string('\n')
	}
	return
}

func (p *P) forward() {
	if !p.curr.isEOF {
		for _, start := range p.starts {
			start.suffix += string(p.curr.rune)
		}
	}
	if len(p.next) == 0 {
		newCurr := p.pullRune(*p.curr)
		p.curr = &newCurr
	} else {
		last := len(p.next) - 1
		p.curr, p.next = p.next[last], p.next[:last]
	}
}

// ============================================================================

type on interface {
	IsOn()
}

type onInt struct {
	f   func(i int, start *Start) (err error)
	opt Opt
}

func (onInt) IsOn() {}

type onUint struct {
	f   func(u uint, start *Start) (err error)
	opt Opt
}

func (onUint) IsOn() {}

type onNumber struct {
	f   func(n bwtype.Number, start *Start) (err error)
	opt Opt
}

func (onNumber) IsOn() {}

type onRange struct {
	f   func(rng bwtype.Range, start *Start) (err error)
	opt Opt
}

func (onRange) IsOn() {}

type onId struct {
	f   func(s string, start *Start) (err error)
	opt Opt
}

func (onId) IsOn() {}

type onString struct {
	f   func(s string, start *Start) (err error)
	opt Opt
}

func (onString) IsOn() {}

type onSubPath struct {
	f   func(path bw.ValPath, start *Start) (err error)
	opt PathOpt
}

func (onSubPath) IsOn() {}

type onPath struct {
	f   func(path bw.ValPath, start *Start) (err error)
	opt PathOpt
}

func (onPath) IsOn() {}

type onArray struct {
	f   func(vals []interface{}, start *Start) (err error)
	opt Opt
}

func (onArray) IsOn() {}

type onArrayOfString struct {
	f   func(ss []string, start *Start) (err error)
	opt Opt
}

func (onArrayOfString) IsOn() {}

type onMap struct {
	f   func(m map[string]interface{}, start *Start) (err error)
	opt Opt
}

func (onMap) IsOn() {}

type onNil struct {
	f   func(start *Start) (err error)
	opt Opt
}

func (onNil) IsOn() {}

type onBool struct {
	f   func(b bool, start *Start) (err error)
	opt Opt
}

func (onBool) IsOn() {}

// ============================================================================

func processOn(p I, processors ...on) (status Status) {
	var (
		i    int
		u    uint
		n    bwtype.Number
		s    string
		path bw.ValPath
		vals []interface{}
		ss   []string
		m    map[string]interface{}
		b    bool
		rng  bwtype.Range
	)
	for _, processor := range processors {
		switch t := processor.(type) {
		case onInt:
			i, status = Int(p, t.opt)
		case onUint:
			u, status = Uint(p, t.opt)
		case onNumber:
			n, status = Number(p, t.opt)
		case onRange:
			rng, status = Range(p, t.opt)
		case onString:
			s, status = String(p, t.opt)
		case onId:
			s, status = Id(p, t.opt)
		case onSubPath:
			path, status = subPath(p, t.opt)
		case onPath:
			path, status = Path(p, t.opt)
		case onArray:
			vals, status = Array(p, t.opt)
		case onArrayOfString:
			ss, status = ArrayOfString(p, t.opt)
		case onMap:
			m, status = Map(p, t.opt)
		case onNil:
			status = Nil(p, t.opt)
		case onBool:
			b, status = Bool(p, t.opt)
		}
		if status.Err != nil {
			return
		}
		if status.OK {
			switch t := processor.(type) {
			case onInt:
				status.Err = t.f(i, status.Start)
			case onUint:
				status.Err = t.f(u, status.Start)
			case onNumber:
				status.Err = t.f(n, status.Start)
			case onRange:
				status.Err = t.f(rng, status.Start)
			case onString:
				status.Err = t.f(s, status.Start)
			case onId:
				status.Err = t.f(s, status.Start)
			case onSubPath:
				status.Err = t.f(path, status.Start)
			case onPath:
				status.Err = t.f(path, status.Start)
			case onArray:
				status.Err = t.f(vals, status.Start)
			case onArrayOfString:
				status.Err = t.f(ss, status.Start)
			case onMap:
				status.Err = t.f(m, status.Start)
			case onNil:
				status.Err = t.f(status.Start)
			case onBool:
				status.Err = t.f(b, status.Start)
			}
			return
		}
	}
	return
}

// ============================================================================

func parseArrayOfString(p I, opt Opt, isEmbeded bool) (result []string, status Status) {
	var needForward uint
	if status.OK = p.Curr().rune == '<'; !status.OK {
		if status.OK = CanSkipRunes(p, 'q', 'w') && IsPunctOrSymbol(p.LookAhead(2).rune); !status.OK {
			return
		}
		needForward = 2
	}
	status.Start = p.Start()
	defer func() { p.Stop(status.Start) }()
	p.Forward(needForward)

	delimiter := p.Curr().rune
	if r, b := bw.Braces()[delimiter]; b {
		delimiter = r
	}
	p.Forward(1)
	result = []string{}
	on := On{p, status.Start, &opt}
	if !isEmbeded {
		base := opt.path
		on.Opt.path = append(base, bw.ValPathItem{Type: bw.ValPathItemIdx})
		defer func() { on.Opt.path = base }()
	}
	parseItem := func(r rune) {
		on.Start = p.Start()
		defer func() { p.Stop(on.Start) }()
		var s string
		for status.Err == nil && !(unicode.IsSpace(r) || r == delimiter) {
			s += string(r)
			p.Forward(1)
			if status.Err = CheckNotEOF(p); status.Err == nil {
				r = p.Curr().rune
			}
		}
		if status.Err == nil {
			if !isEmbeded && opt.OnValidateArrayOfStringElem != nil {
				status.Err = opt.OnValidateArrayOfStringElem(on, result, s)
			} else if opt.OnValidateString != nil {
				status.Err = opt.OnValidateString(on, s)
			}
			if status.Err == nil {
				result = append(result, s)
			}
		}
		on.Opt.path[len(on.Opt.path)-1].Idx++
	}
	for status.Err == nil {
		if _, status.Err = SkipSpace(p, TillNonEOF); status.Err == nil {
			r := p.Curr().rune
			if r == delimiter {
				p.Forward(1)
				break
			}
			parseItem(r)
		}
	}
	return
}

// ============================================================================

func parsePath(p I, opt PathOpt) (result bw.ValPath, status Status) {
	var justParsed pathResult
	curr := p.Curr()
	if justParsed, status.OK = curr.justParsed.(pathResult); status.OK {
		result = justParsed.path
		status.Start = justParsed.start
		// p.Forward(curr.justForward)
		p.Forward(uint(len(justParsed.start.suffix)))
		return
	} else if status.OK = curr.rune == '$'; status.OK {
		status.Start = p.Start()
		result, status.Err = PathContent(p, opt)
	} else if status.OK = CanSkipRunes(p, '{', '{'); status.OK {
		status.Start = p.Start()
		p.Forward(2)
		if _, status.Err = SkipSpace(p, TillNonEOF); status.Err == nil {
			if result, status.Err = PathContent(p, opt); status.Err == nil {
				if _, status.Err = SkipSpace(p, TillNonEOF); status.Err == nil {
					if !SkipRunes(p, '}', '}') {
						status.Err = Unexpected(p)
					}
				}
			}
		}
	}
	return
}

// ============================================================================

type numberResult struct {
	n     bwtype.Number
	start *Start
}

type pathResult struct {
	path  bw.ValPath
	start *Start
}

// ============================================================================

func parseNumber(p I, opt Opt, rangeLimitKind RangeLimitKind) (result bwtype.Number, status Status) {
	var (
		s          string
		hasDot     bool
		b          bool
		isNegative bool
		justParsed numberResult
	)
	nonNegativeNumber := false
	if opt.NonNegativeNumber != nil {
		nonNegativeNumber = opt.NonNegativeNumber(rangeLimitKind)
	}
	curr := p.Curr()
	if justParsed, status.OK = curr.justParsed.(numberResult); status.OK {
		if nonNegativeNumber {
			if _, status.OK = bwtype.Uint(justParsed.n.Val()); !status.OK {
				status.Err = Unexpected(p)
				return
			}
		}
		status.Start = justParsed.start
		result = justParsed.n
		// p.Forward(curr.justForward)
		p.Forward(uint(len(justParsed.start.suffix)))
		bwdebug.Print("justParsed.start:json", justParsed.start)
	} else if s, isNegative, status = looksLikeNumber(p, nonNegativeNumber); status.IsOK() {
		for {
			if s, b = addDigit(p, s); !b {
				if !hasDot && CanSkipRunes(p, dotRune) {
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
		// p.Stop(status.Start)

		if hasDot && !zeroAfterDotRegexp.MatchString(s) {
			var f float64
			if f, status.Err = strconv.ParseFloat(s, 64); status.Err == nil {
				result = bwtype.MustNumberFrom(f)
			}
		} else {
			if pos := strings.LastIndex(s, string(dotRune)); pos >= 0 {
				s = s[:pos]
			}
			if isNegative {
				var i int
				if i, status.Err = bwstr.ParseInt(s); status.Err == nil {
					result = bwtype.MustNumberFrom(i)
				}
			} else {
				var u uint
				if u, status.Err = bwstr.ParseUint(s); status.Err == nil {
					result = bwtype.MustNumberFrom(u)
				}
			}
		}
		status.UnexpectedIfErr(p)
	}
	return
}

const dotRune = '.'

var zeroAfterDotRegexp = regexp.MustCompile(`\.0+$`)

// ============================================================================

func getIdExpects(opt Opt, indent string) (expects string) {
	if len(opt.IdVals) > 0 {
		sset := bwset.String{}
		for s := range opt.IdVals {
			sset.Add(s)
		}
		var suffix string
		if opt.OnId != nil {
			suffix = ansi.String(" or <ansiVar>custom<ansi>")
		}
		expects = bwstr.SmartJoin(bwstr.A{
			Source: bwstr.SS{
				SS: sset.ToSliceOfStrings(),
				Preformat: func(s string) string {
					return fmt.Sprintf(ansi.String("<ansiVal>%s"), s)
				},
			},
			MaxLen:              uint(80 - len(suffix)),
			NoJoinerOnMutliline: true,
			InitialIndent:       indent,
		}) + suffix
	}
	return
}

// one of [Int, Number, String] or custom
// ============================================================================

func parseDelimitedOptionalCommaSeparated(p I, openDelimiter, closeDelimiter rune, opt Opt, fn func(on On, base bw.ValPath) error) (status Status) {
	p.Forward(Initial)
	if status.OK = CanSkipRunes(p, openDelimiter); status.OK {
		status.Start = p.Start()
		defer func() { p.Stop(status.Start) }()
		p.Forward(1)
		base := opt.path
		on := On{p, status.Start, &opt}
		defer func() { on.Opt.path = base }()
	LOOP:
		for status.Err == nil {
			if _, status.Err = SkipSpace(p, TillNonEOF); status.Err == nil {
			NEXT:
				if SkipRunes(p, closeDelimiter) {
					break LOOP
				}
				if status.Err = fn(on, base); status.Err == nil {
					var isSpaceSkipped bool
					if isSpaceSkipped, status.Err = SkipSpace(p, TillNonEOF); status.Err == nil {
						if !SkipRunes(p, ',') {
							if !isSpaceSkipped {
								status.Err = Unexpected(p)
							} else {
								goto NEXT
							}
						}
					}
				}
			}
		}
	}
	bwdebug.Print("!parseDelimitedOptionalCommaSeparated", "status", status)
	return
}

func looksLikeNumber(p I, nonNegative bool) (s string, isNegative bool, status Status) {
	var (
		r         rune
		needDigit bool
	)
	p.Forward(Initial)
	r = p.Curr().rune
	if status.OK = r == '+'; status.OK {
		needDigit = true
	} else if status.OK = !nonNegative && r == '-'; status.OK {
		s = string(r)
		needDigit = true
		isNegative = true
	} else if status.OK = IsDigit(r); status.OK {
		s = string(r)
	} else {
		return
	}
	status.Start = p.Start()
	p.Forward(1)
	if needDigit {
		if r = p.Curr().rune; !IsDigit(r) {
			status.Err = Unexpected(p)
		} else {
			p.Forward(1)
			s += string(r)
		}
	}
	return
}

func addDigit(p I, s string) (string, bool) {
	r := p.Curr().rune
	if IsDigit(r) {
		s += string(r)
	} else if r != '_' {
		return s, false
	}
	p.Forward(1)
	return s, true
}

// ============================================================================

func subPath(p I, opt PathOpt) (result bw.ValPath, status Status) {
	if status.OK = CanSkipRunes(p, '('); status.OK {
		status.Start = p.Start()
		defer func() { p.Stop(status.Start) }()
		p.Forward(1)
		if _, status.Err = SkipSpace(p, TillNonEOF); status.Err == nil {
			subOpt := opt
			subOpt.isSubPath = true
			subOpt.Bases = opt.Bases
			if result, status.Err = PathContent(p, subOpt); status.Err == nil {
				if _, status.Err = SkipSpace(p, TillNonEOF); status.Err == nil {
					if p.Curr().rune == ')' {
						p.Forward(1)
					} else {
						status.Err = Unexpected(p)
					}
				}
			}
		}
	}
	return
}

// ============================================================================
