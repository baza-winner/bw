package bwparse

import (
	"fmt"
	"strings"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwset"
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
	ansiUnexpectedWord = ansi.String("unexpected <ansiErr>%q<ansi>")
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

// func optKeyMap(opt map[string]interface{}, key string, keys *bwset.String) (result map[string]interface{}, ok bool) {
// 	var val interface{}
// 	keys.Add(key)
// 	if val, ok = opt[key]; ok && val != nil {
// 		if result, ok = val.(map[string]interface{}); !ok {
// 			bwerr.Panic(ansiOptKeyIsNotOfType, key, val, "map[string]interface{}")
// 		}
// 	}
// 	return
// }

// func optKeyValidateMapKeyFunc(opt map[string]interface{}, key string, keys *bwset.String) (result ValidateMapKeyFunc, ok bool) {
// 	var val interface{}
// 	keys.Add(key)
// 	if val, ok = opt[key]; ok && val != nil {
// 		if result, ok = val.(ValidateMapKeyFunc); !ok {
// 			bwerr.Panic(ansiOptKeyIsNotOfType, key, val, "ValidateMapKeyFunc")
// 		}
// 	}
// 	return
// }

// func optKeyParseMapElemFunc(opt map[string]interface{}, key string, keys *bwset.String) (result ParseMapElemFunc, ok bool) {
// 	var val interface{}
// 	keys.Add(key)
// 	if val, ok = opt[key]; ok && val != nil {
// 		if result, ok = val.(ParseMapElemFunc); !ok {
// 			bwerr.Panic(ansiOptKeyIsNotOfType, key, val, "ParseMapElemFunc")
// 		}
// 	}
// 	return
// }

// func optKeyOnMapBeginEndFunc(opt map[string]interface{}, key string, keys *bwset.String) (result OnMapBeginEndFunc, ok bool) {
// 	var val interface{}
// 	keys.Add(key)
// 	if val, ok = opt[key]; ok && val != nil {
// 		if result, ok = val.(OnMapBeginEndFunc); !ok {
// 			bwerr.Panic(ansiOptKeyIsNotOfType, key, val, "OnMapBeginEndFunc")
// 		}
// 	}
// 	return
// }

// func optKeyParseArrayElemFunc(opt map[string]interface{}, key string, keys *bwset.String) (result ParseArrayElemFunc, ok bool) {
// 	var val interface{}
// 	keys.Add(key)
// 	if val, ok = opt[key]; ok && val != nil {
// 		if result, ok = val.(ParseArrayElemFunc); !ok {
// 			bwerr.Panic(ansiOptKeyIsNotOfType, key, val, "ParseArrayElemFunc")
// 		}
// 	}
// 	return
// }

// func optKeyOnArrayBeginEndFunc(opt map[string]interface{}, key string, keys *bwset.String) (result OnArrayBeginEndFunc, ok bool) {
// 	var val interface{}
// 	keys.Add(key)
// 	if val, ok = opt[key]; ok && val != nil {
// 		if result, ok = val.(OnArrayBeginEndFunc); !ok {
// 			bwerr.Panic(ansiOptKeyIsNotOfType, key, val, "OnArrayBeginEndFunc")
// 		}
// 	}
// 	return
// }

// func optKeyValidateArrayOfStringElemFunc(opt map[string]interface{}, key string, keys *bwset.String) (result ValidateArrayOfStringElemFunc, ok bool) {
// 	var val interface{}
// 	keys.Add(key)
// 	if val, ok = opt[key]; ok && val != nil {
// 		if result, ok = val.(ValidateArrayOfStringElemFunc); !ok {
// 			bwerr.Panic(ansiOptKeyIsNotOfType, key, val, "ValidateArrayOfStringElemFunc")
// 		}
// 	}
// 	return
// }

// ============================================================================

func getOpt(optOpt []Opt) Opt {
	opt := Opt{}
	if len(optOpt) > 0 {
		opt = optOpt[0]
	}
	return opt
}

func isOneOfId(p I, ss []string) (ok bool) {
	for _, s := range ss {
		if ok = CanSkipRunes(p, []rune(s)...); ok {
			u := uint(len(s))
			r := p.LookAhead(u).rune
			if ok = !(IsLetter(r) || IsDigit(r)); ok {
				p.Forward(u)
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

func (p *P) suffix(ps PosInfo) (suffix string) {
	if ps.pos > p.curr.pos {
		bwerr.Panic(ansiGetSuffixAssert, ps.pos, p.curr.pos)
	}

	preLineCount, postLineCount := int(p.preLineCount), int(p.postLineCount)
	if p.curr.isEOF {
		preLineCount += postLineCount
	}

	var separator string
	if p.curr.line > 1 {
		suffix += fmt.Sprintf(ansiLineCol, ps.line, ps.col, ps.pos)
		separator = "\n"
	} else {
		suffix += fmt.Sprintf(ansiPos, ps.pos)
		separator = " "
	}
	suffix += ":" + separator + ansiOK + p.curr.prefix[0:ps.pos-p.curr.prefixStart]

	var needPostLines, noNeedNewline bool
	if ps.pos < p.curr.pos {
		noNeedNewline = p.curr.prefix[len(p.curr.prefix)-1] == '\n'
		suffix += ansiErr + p.curr.prefix[ps.pos-p.curr.prefixStart:] + ansi.Reset()
		needPostLines = true
	} else if !p.curr.isEOF {
		noNeedNewline = p.curr.rune == '\n'
		suffix += ansiErr + string(p.curr.rune) + ansi.Reset()
		p.Forward(1)
		needPostLines = true
	}

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
	if len(p.next) == 0 {
		// ps := p.curr
		newCurr := p.pullRune(*p.curr)
		p.curr = &newCurr
		// bwdebug.Print("ps:#v", ps, "p.curr:#v", p.curr)
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
	f func(i int, start *PosInfo) (err error)
}

func (onInt) IsOn() {}

type onUint struct {
	f func(u uint, start *PosInfo) (err error)
}

func (onUint) IsOn() {}

type onNumber struct {
	f func(n bwtype.Number, start *PosInfo) (err error)
}

func (onNumber) IsOn() {}

type onRange struct {
	f func(rng bwtype.Range, start *PosInfo) (err error)
}

func (onRange) IsOn() {}

type onId struct {
	f func(s string, start *PosInfo) (err error)
}

func (onId) IsOn() {}

type onString struct {
	f func(s string, start *PosInfo) (err error)
}

func (onString) IsOn() {}

type onSubPath struct {
	f func(path bw.ValPath, start *PosInfo) (err error)
	a PathA
}

func (onSubPath) IsOn() {}

type onPath struct {
	f func(path bw.ValPath, start *PosInfo) (err error)
	a PathA
}

func (onPath) IsOn() {}

type onArray struct {
	f func(vals []interface{}, start *PosInfo) (err error)
}

func (onArray) IsOn() {}

type onArrayOfString struct {
	f func(ss []string, start *PosInfo) (err error)
}

func (onArrayOfString) IsOn() {}

type onMap struct {
	f func(m map[string]interface{}, start *PosInfo) (err error)
}

func (onMap) IsOn() {}

type onNil struct {
	f func(start *PosInfo) (err error)
}

func (onNil) IsOn() {}

type onBool struct {
	f func(b bool, start *PosInfo) (err error)
}

func (onBool) IsOn() {}

// ============================================================================

func processOn(p I, opt Opt, processors ...on) (ok bool, err error) {
	var (
		start *PosInfo
		i     int
		u     uint
		n     bwtype.Number
		s     string
		path  bw.ValPath
		// val   interface{}
		vals []interface{}
		ss   []string
		m    map[string]interface{}
		b    bool
		rng  bwtype.Range
	)
	for _, processor := range processors {
		switch t := processor.(type) {
		case onInt:
			i, start, ok, err = Int(p, opt)
		case onUint:
			u, start, ok, err = Uint(p, opt)
		case onNumber:
			n, start, ok, err = Number(p, opt)
		case onRange:
			rng, start, ok, err = Range(p, opt)
		case onString:
			s, start, ok, err = String(p, opt)
		case onId:
			s, start, ok, err = Id(p, opt)
		case onSubPath:
			path, start, ok, err = subPath(p, t.a, opt)
		case onPath:
			path, start, ok, err = Path(p, t.a, opt)
		case onArray:
			vals, start, ok, err = Array(p, opt)
		case onArrayOfString:
			ss, start, ok, err = ArrayOfString(p, opt)
		case onMap:
			m, start, ok, err = Map(p, opt)
		case onNil:
			start, ok = Nil(p, opt)
		case onBool:
			b, start, ok = Bool(p, opt)
		}
		if err != nil {
			return
		}
		if ok {
			switch t := processor.(type) {
			case onInt:
				err = t.f(i, start)
			case onUint:
				err = t.f(u, start)
			case onNumber:
				err = t.f(n, start)
			case onRange:
				err = t.f(rng, start)
			case onString:
				err = t.f(s, start)
			case onId:
				err = t.f(s, start)
			case onSubPath:
				err = t.f(path, start)
			case onPath:
				err = t.f(path, start)
			case onArray:
				err = t.f(vals, start)
			case onArrayOfString:
				err = t.f(ss, start)
			case onMap:
				err = t.f(m, start)
			case onNil:
				err = t.f(start)
			case onBool:
				err = t.f(b, start)
			}
			return
		}
	}
	return
}

// func (p *P) processOn(processors ...on) (ok bool, err error) {
// 	var (
// 		start *PosInfo
// 		idx   int
// 		s     string
// 		path  bw.ValPath
// 		val   interface{}
// 		vals  []interface{}
// 		ss    []string
// 		m     map[string]interface{}
// 		b     bool
// 		rng   bwtype.Range
// 	)
// 	for _, processor := range processors {
// 		switch t := processor.(type) {
// 		case onInt:
// 			idx, start, ok, err = p.Int()
// 		case onNumber:
// 			val, start, ok, err = p.Number()
// 		case onRange:
// 			rng, start, ok, err = p.Range()
// 		case onString:
// 			s, start, ok, err = p.String()
// 		case onId:
// 			s, start, ok, err = p.Id()
// 			// bwdebug.Print("start:#v", start)
// 		case onSubPath:
// 			path, start, ok, err = p.subPath(t.a)
// 		case onPath:
// 			path, start, ok, err = p.Path(t.a)
// 		case onArray:
// 			vals, start, ok, err = p.Array()
// 		case onArrayOfString:
// 			ss, start, ok, err = p.ArrayOfString()
// 		case onMap:
// 			m, start, ok, err = p.Map()
// 		case onNil:
// 			start, ok = p.Nil()
// 		case onBool:
// 			b, start, ok = p.Bool()
// 		}
// 		if err != nil {
// 			return
// 		}
// 		if ok {
// 			switch t := processor.(type) {
// 			case onInt:
// 				err = t.f(idx, start)
// 			case onNumber:
// 				err = t.f(val, start)
// 			case onRange:
// 				err = t.f(rng, start)
// 			case onString:
// 				err = t.f(s, start)
// 			case onId:
// 				err = t.f(s, start)
// 			case onSubPath:
// 				err = t.f(path, start)
// 			case onPath:
// 				err = t.f(path, start)
// 			case onArray:
// 				err = t.f(vals, start)
// 			case onArrayOfString:
// 				err = t.f(ss, start)
// 			case onMap:
// 				err = t.f(m, start)
// 			case onNil:
// 				err = t.f(start)
// 			case onBool:
// 				err = t.f(b, start)
// 			}
// 			return
// 		}
// 	}
// 	return
// }

// ============================================================================

// ============================================================================

// func (p *P) mapBegin(m map[string]interface{}) (err error) {
// 	if p.OnMapBegin != nil {
// 		err = p.OnMapBegin(p, m)
// 	}
// 	return
// }
// func (p *P) validateMapKey(m map[string]interface{}, key string, start *PosInfo) (err error) {
// 	if p.OnValidateMapKey != nil {
// 		ok, err = p.OnValidateMapKey(p, m, key, start)
// 	}
// 	return
// }

// func (p *P) parseMapElem(m map[string]interface{}, key string) (ok bool, err error) {
// 	if p.OnParseMapElem != nil {
// 		ok, err = p.OnParseMapElem(p, m, key)
// 	}
// 	return
// }

// func (p *P) mapEnd(m map[string]interface{}) (err error) {
// 	if p.OnMapEnd != nil {
// 		err = p.OnMapEnd(p, m)
// 	}
// 	return
// }

// func (p *P) arrayBegin(vals []interface{}) (err error) {
// 	if p.OnArrayBegin != nil {
// 		err = p.OnArrayBegin(p, vals)
// 	}
// 	return
// }

// func (p *P) parseArrayElem(vals []interface{}) (ok bool, err error) {
// 	if p.OnParseArrayElem != nil {
// 		ok, err = p.OnParseArrayElem(p, vals)
// 	}
// 	return
// }

// func (p *P) arrayEnd(vals []interface{}) (err error) {
// 	if p.OnArrayBegin != nil {
// 		err = p.OnArrayBegin(p, vals)
// 	}
// 	return
// }

// func (p *P) validateArrayOfStringElem(ss []string, s string, start *PosInfo) (err error) {
// 	if p.OnValidateArrayOfStringElem != nil {
// 		err = p.OnValidateArrayOfStringElem(p, ss, s, start)
// 	}
// 	return
// }

// ============================================================================

func getStart(p I) *PosInfo {
	p.Forward(Initial)
	return p.Curr()
}

func parseDelimitedOptionalCommaSeparated(p I, openDelimiter, closeDelimiter rune, f func() error) (start *PosInfo, ok bool, err error) {
	start = getStart(p)
	if ok = SkipRunes(p, openDelimiter); ok {
	LOOP:
		for err == nil {
			if err = SkipSpace(p, TillNonEOF); err == nil {
			NEXT:
				if SkipRunes(p, closeDelimiter) {
					break LOOP
				}
				if err = f(); err == nil {
					if err = SkipSpace(p, TillNonEOF); err == nil {
						if !SkipRunes(p, ',') {
							goto NEXT
						}
					}
				}
			}
		}
	}
	return
}

func looksLikeNumber(p I, nonNegative bool) (s string, start *PosInfo, isNegative bool, ok bool, err error) {
	start = getStart(p)
	var (
		r         rune
		needDigit bool
	)
	// bwdebug.Print("!HERE")
	r = p.Curr().rune
	if ok = r == '+'; ok {
		needDigit = true
	} else if ok = !nonNegative && r == '-'; ok {
		s = string(r)
		needDigit = true
		isNegative = true
	} else if ok = IsDigit(r); ok {
		s = string(r)
	} else {
		return
	}
	p.Forward(1)
	if needDigit {
		if r = p.Curr().rune; !IsDigit(r) {
			err = Unexpected(p)
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

func subPath(p I, a PathA, opt Opt) (result bw.ValPath, start *PosInfo, ok bool, err error) {
	start = getStart(p)
	if ok = SkipRunes(p, '('); ok {
		if err = SkipSpace(p, TillNonEOF); err == nil {
			if result, err = PathContent(p, PathA{isSubPath: true, Bases: a.Bases}, opt); err == nil {
				if err = SkipSpace(p, TillNonEOF); err == nil {
					if p.Curr().rune == ')' {
						p.Forward(1)
					} else {
						err = Unexpected(p)
					}
				}
			}
		}
	}
	return
}

// ============================================================================
