package bwparse

import (
	"fmt"
	"strings"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
)

// ============================================================================

var (
	ansiPosInfoFailed   string
	ansiOK              string
	ansiErr             string
	ansiPos             string
	ansiLineCol         string
	ansiGetSuffixAssert string
	ansiUnexpectedEOF   string
	ansiUnexpectedChar  string
	ansiUnexpectedWord  string
)

func init() {
	ansiOK = ansi.CSIFromSGRCodes(ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen, Bright: false})).String()
	ansiErr = ansi.CSIFromSGRCodes(ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorRed, Bright: true})).String()

	ansiPosInfoFailed = "<ansiPath>bwparse.Provider.<ansiFunc>PosInfo<ansi>(%d) failed, must <ansiFunc>.SetMaxBackwardCount<ansi>(<ansiVal>%d<ansi>)"
	ansiPos = ansi.String(" at pos <ansiPath>%d<ansi>")
	ansiLineCol = ansi.String(" at line <ansiPath>%d<ansi>, col <ansiPath>%d<ansi> (pos <ansiPath>%d<ansi>)")
	ansiGetSuffixAssert = ansi.String("<ansiVar>ps.Pos<ansi> (<ansiVal>%d<ansi>) > <ansiVar>p.Curr.Pos<ansi> (<ansiVal>%d<ansi>)")
	ansiUnexpectedEOF = ansi.String("unexpected end of string")
	ansiUnexpectedChar = ansi.String("unexpected char <ansiVal>%q<ansiReset> (<ansiVar>charCode<ansi>: <ansiVal>%d<ansi>)")
	ansiUnexpectedWord = ansi.String("unexpected <ansiErr>%q<ansi>")
}

func (p *P) pullRune(ps *PosInfo) {
	var runePtr *rune
	var err error
	if runePtr, err = p.Prov.PullRune(); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	if !ps.IsEOF {
		if ps.Pos >= 0 {
			ps.Prefix += string(ps.Rune)
		}
		if runePtr != nil {
			if ps.Rune != '\n' {
				ps.Col++
			} else {
				ps.Line++
				ps.Col = 1
				if int(ps.Line) > int(p.preLineCount) {
					i := strings.Index(ps.Prefix, "\n")
					ps.Prefix = ps.Prefix[i+1:]
					ps.PrefixStart += i + 1
				}
			}
		}
	}
	if runePtr == nil {
		ps.Rune = '\000'
		ps.IsEOF = true
	} else {
		ps.Rune = *runePtr
		ps.IsEOF = false
	}
	ps.Pos++
	return
}

func (p *P) suffix(ps PosInfo) (suffix string) {
	if ps.Pos > p.Curr.Pos {
		bwerr.Panic(ansiGetSuffixAssert, ps.Pos, p.Curr.Pos)
	}
	preLineCount := p.preLineCount
	postLineCount := p.postLineCount
	if p.Curr.IsEOF {
		preLineCount += postLineCount
	}

	separator := "\n"
	if p.Curr.Line <= 1 {
		suffix += fmt.Sprintf(ansiPos, ps.Pos)
		separator = " "
	} else {
		suffix += fmt.Sprintf(ansiLineCol, ps.Line, ps.Col, ps.Pos)
	}
	suffix += ":" + separator + ansiOK

	suffix += p.Curr.Prefix[0 : ps.Pos-p.Curr.PrefixStart]
	var needPostLines, noNeedNewline bool
	if ps.Pos < p.Curr.Pos {
		suffix += ansiErr
		suffix += p.Curr.Prefix[ps.Pos-p.Curr.PrefixStart:]
		noNeedNewline = suffix[len(suffix)-1] == '\n'
		suffix += ansi.Reset()
		needPostLines = true
	} else if !p.Curr.IsEOF {
		suffix += ansiErr
		suffix += string(p.Curr.Rune)
		noNeedNewline = p.Curr.Rune == '\n'
		suffix += ansi.Reset()
		p.Forward(true)
		needPostLines = true
	}

	if needPostLines {
		for !p.Curr.IsEOF && postLineCount > 0 {
			suffix += string(p.Curr.Rune)
			if p.Curr.Rune == '\n' {
				postLineCount -= 1
				noNeedNewline = true
			} else {
				noNeedNewline = false
			}
			p.Forward(true)
		}
	}
	_ = noNeedNewline

	if !noNeedNewline {
		suffix += string('\n')
	}

	return suffix
}

// ============================================================================

func (p *P) subPath(a PathA) (result bw.ValPath, start PosInfo, ok bool, err error) {
	if p.Curr.Rune == '(' {
		ok = true
		start = p.Curr
		p.Forward(true)
		if err = p.SkipOptionalSpace(); err != nil {
			return
		}
		a.isSubPath = true
		if result, err = p.PathContent(a); err != nil {
			return
		}
		if err = p.SkipOptionalSpace(); err != nil {
			return
		}
		if p.Curr.Rune != ')' {
			err = p.Unexpected(p.Curr)
			return
		}
		p.Forward(true)
	}
	return
}

// ============================================================================
