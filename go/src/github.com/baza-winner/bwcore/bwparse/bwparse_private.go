package bwparse

import (
	"fmt"
	"strings"
	"unicode"

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
	ansiGetSuffixAssert = ansi.String("<ansiVar>ps.Pos<ansi> (<ansiVal>%d<ansi>) > <ansiVar>p.Curr.Pos<ansi> (<ansiVal>%d<ansi>)")
	ansiUnexpectedEOF = ansi.String("unexpected end of string")
	ansiUnexpectedChar = ansi.String("unexpected char <ansiVal>%q<ansiReset> (<ansiVar>charCode<ansi>: <ansiVal>%d<ansi>)")
	ansiUnexpectedWord = ansi.String("unexpected <ansiErr>%q<ansi>")
	ansiOutOfRange = ansi.String("<ansiVal>%d<ansi> is out of range <ansiVal>%d..%d")
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

func optKeyMap(opt map[string]interface{}, key string, keys *bwset.String) (result map[string]interface{}, ok bool) {
	var val interface{}
	keys.Add(key)
	if val, ok = opt[key]; ok && val != nil {
		if result, ok = val.(map[string]interface{}); !ok {
			bwerr.Panic(ansiOptKeyIsNotOfType, key, val, "map[string]interface{}")
		}
	}
	return
}

var ()

func init() {
}

// ============================================================================

func (p *P) pullRune(ps *PosInfo) {
	runePtr := bwrune.MustPull(p.Prov)
	if !ps.IsEOF {
		if ps.Pos >= 0 {
			ps.Prefix += string(ps.Rune)
			// bwdebug.Print("ps.Prefix", ps.Prefix, "ps.Rune:s", ps.Rune)
		}
		if runePtr != nil {
			if ps.Rune != '\n' {
				ps.Col++
			} else {
				ps.Line++
				ps.Col = 1
				if int(ps.Line) > int(p.preLineCount)+1 {
					i := strings.Index(ps.Prefix, "\n")
					// bwdebug.Print("ps.Prefix", ps.Prefix, "i", i)
					ps.Prefix = ps.Prefix[i+1:]
					ps.PrefixStart += i + 1
					// bwdebug.Print("ps.Line", ps.Line, "p.preLineCount", p.preLineCount, "ps.Prefix", ps.Prefix, "ps.PrefixStart", ps.PrefixStart)
				}
			}
		}
		if runePtr == nil {
			ps.Rune, ps.IsEOF = '\000', true
		} else {
			ps.Rune, ps.IsEOF = *runePtr, false
		}
		ps.Pos++
	}
}

func (p *P) suffix(ps PosInfo) (suffix string) {
	if ps.Pos > p.Curr.Pos {
		bwerr.Panic(ansiGetSuffixAssert, ps.Pos, p.Curr.Pos)
	}

	preLineCount, postLineCount := int(p.preLineCount), int(p.postLineCount)
	if p.Curr.IsEOF {
		preLineCount += postLineCount
	}

	var separator string
	if p.Curr.Line > 1 {
		suffix += fmt.Sprintf(ansiLineCol, ps.Line, ps.Col, ps.Pos)
		separator = "\n"
	} else {
		suffix += fmt.Sprintf(ansiPos, ps.Pos)
		separator = " "
	}
	suffix += ":" + separator + ansiOK + p.Curr.Prefix[0:ps.Pos-p.Curr.PrefixStart]

	var needPostLines, noNeedNewline bool
	if ps.Pos < p.Curr.Pos {
		noNeedNewline = p.Curr.Prefix[len(p.Curr.Prefix)-1] == '\n'
		suffix += ansiErr + p.Curr.Prefix[ps.Pos-p.Curr.PrefixStart:] + ansi.Reset()
		needPostLines = true
	} else if !p.Curr.IsEOF {
		noNeedNewline = p.Curr.Rune == '\n'
		suffix += ansiErr + string(p.Curr.Rune) + ansi.Reset()
		p.Forward(1)
		needPostLines = true
	}

	// bwdebug.Print("needPostLines", needPostLines, "p.Curr.IsEOF", p.Curr.IsEOF, "postLineCount", postLineCount)
	for needPostLines && !p.Curr.IsEOF && postLineCount >= 0 {
		suffix += string(p.Curr.Rune)
		// bwdebug.Print("p.Curr.Rune:s", p.Curr.Rune, "postLineCount", postLineCount)
		if noNeedNewline = p.Curr.Rune == '\n'; noNeedNewline {
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
	if len(p.Next) == 0 {
		p.pullRune(&p.Curr)
	} else {
		last := len(p.Next) - 1
		p.Curr, p.Next = p.Next[last], p.Next[:last]
	}
}

// ============================================================================

type on interface {
	IsOn()
}

type onInt struct {
	f func(idx int, start PosInfo) (err error)
}

func (onInt) IsOn() {}

type onNumber struct {
	f func(val interface{}, start PosInfo) (err error)
}

func (onNumber) IsOn() {}

type onId struct {
	f func(s string, start PosInfo) (err error)
}

func (onId) IsOn() {}

type onString struct {
	f func(s string, start PosInfo) (err error)
}

func (onString) IsOn() {}

type onSubPath struct {
	f func(path bw.ValPath, start PosInfo) (err error)
	a PathA
}

func (onSubPath) IsOn() {}

type onPath struct {
	f func(path bw.ValPath, start PosInfo) (err error)
	a PathA
}

func (onPath) IsOn() {}

type onArray struct {
	f func(vals []interface{}, start PosInfo) (err error)
}

func (onArray) IsOn() {}

type onArrayOfString struct {
	f func(ss []string, start PosInfo) (err error)
}

func (onArrayOfString) IsOn() {}

type onMap struct {
	f func(m map[string]interface{}, start PosInfo) (err error)
}

func (onMap) IsOn() {}

// ============================================================================

func (p *P) processOn(processors ...on) (ok bool, err error) {
	var (
		start PosInfo
		idx   int
		s     string
		path  bw.ValPath
		val   interface{}
		vals  []interface{}
		ss    []string
		m     map[string]interface{}
	)
	for _, processor := range processors {
		switch t := processor.(type) {
		case onInt:
			idx, start, ok, err = p.Int()
		case onNumber:
			val, start, ok, err = p.Number()
		case onString:
			s, start, ok, err = p.String()
		case onId:
			s, start, ok, err = p.Id()
		case onSubPath:
			path, start, ok, err = p.subPath(t.a)
		case onPath:
			path, start, ok, err = p.Path(t.a)
		case onArray:
			vals, start, ok, err = p.Array()
		case onArrayOfString:
			ss, start, ok, err = p.ArrayOfString()
		case onMap:
			m, start, ok, err = p.Map()
		}
		if err != nil {
			return
		}
		if ok {
			switch t := processor.(type) {
			case onInt:
				err = t.f(idx, start)
			case onNumber:
				err = t.f(val, start)
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
			}
			return
		}
	}
	return
}

// ============================================================================

func (p *P) start() PosInfo {
	p.Forward(Initial)
	return p.Curr
}

func (p *P) parseDelimitedOptionalCommaSeparated(openDelimiter, closeDelimiter rune, f func() error) (start PosInfo, ok bool, err error) {
	start = p.start()
	if ok = p.skipRunes(openDelimiter); ok {
	LOOP:
		for err == nil {
			if err = p.SkipSpace(TillNonEOF); err == nil {
			NEXT:
				if p.skipRunes(closeDelimiter) {
					break LOOP
				}
				if err = f(); err == nil {
					if err = p.SkipSpace(TillNonEOF); err == nil {
						if !p.skipRunes(',') {
							goto NEXT
						}
					}
				}
			}
		}
	}
	return
}

func (p *P) parseVal() (result interface{}, err error) {
	var ok bool
	if result, _, ok, err = p.Val(); err == nil && !ok {
		err = p.Unexpected(p.Curr)
	}
	return
}

func (p *P) looksLikeNumber() (s string, start PosInfo, ok bool, err error) {
	start = p.start()
	var (
		r         rune
		needDigit bool
	)
	if r = p.Curr.Rune; r == '-' || r == '+' {
		needDigit = true
	} else if !isDigit(r) {
		return
	}
	ok = true
	p.Forward(1)
	s = string(r)
	if needDigit {
		if r = p.Curr.Rune; !isDigit(r) {
			err = p.Unexpected(p.Curr)
		} else {
			p.Forward(1)
			s += string(r)
		}
	}
	return
}

// func parseInt(s string) (result int, err error) {
// 	if _int64, err := strconv.ParseInt(s, 10, 64); err != nil {
// 		return 0, err
// 	} else {
// 		// } else if int64(bw.MinInt) <= _int64 && _int64 <= int64(bw.MaxInt) {
// 		return int(_int64), nil
// 		// } else {
// 		// 	return 0, bwerr.From(ansiOutOfRange, _int64, bw.MinInt, bw.MaxInt)
// 	}
// }

func (p *P) addDigit(r rune, s string) (string, bool) {
	if isDigit(r) {
		s += string(r)
	} else if r != '_' {
		return s, false
	}
	p.Forward(1)
	return s, true
}

func (p *P) canSkipRunes(rr ...rune) bool {
	for i, r := range rr {
		if pi := p.LookAhead(uint(i)); pi.IsEOF || pi.Rune != r {
			return false
		}
	}
	return true
}

func (p *P) skipRunes(rr ...rune) (ok bool) {
	if ok = p.canSkipRunes(rr...); ok {
		p.Forward(uint(len(rr)))
	}
	return
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isPunctOrSymbol(r rune) bool {
	return unicode.IsPunct(r) || unicode.IsSymbol(r)
}

// ============================================================================

func (p *P) subPath(a PathA) (result bw.ValPath, start PosInfo, ok bool, err error) {
	start = p.start()
	if ok = p.skipRunes('('); ok {
		if err = p.SkipSpace(TillNonEOF); err == nil {
			if result, err = p.PathContent(PathA{isSubPath: true, Bases: a.Bases}); err == nil {
				if err = p.SkipSpace(TillNonEOF); err == nil {
					if p.Curr.Rune == ')' {
						p.Forward(1)
					} else {
						err = p.Unexpected(p.Curr)
					}
				}
			}
		}
	}
	return
}

// ============================================================================
