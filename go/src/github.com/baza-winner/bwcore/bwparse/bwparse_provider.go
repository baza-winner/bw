package bwparse

import (
	"fmt"
	"strings"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwrune"
)

// ============================================================================

type PosStruct struct {
	IsEOF       bool
	Rune        rune
	Pos         int
	Line        uint
	Col         uint
	Prefix      string
	PrefixStart int
}

type PosStructs []PosStruct

// ============================================================================

type Provider struct {
	Prov          bwrune.Provider
	Curr          PosStruct
	Next          PosStructs
	preLineCount  uint
	postLineCount uint
}

func ProviderFrom(p bwrune.Provider) (result *Provider) {
	result = &Provider{
		Prov:          p,
		Curr:          PosStruct{Pos: -1, Line: 1},
		Next:          PosStructs{},
		preLineCount:  3,
		postLineCount: 3,
	}
	return
}

const Initial bool = false

func (p *Provider) Forward(nonInitial bool) {
	if p.Curr.Pos < 0 || nonInitial && !p.Curr.IsEOF {
		if len(p.Next) == 0 {
			p.pullRune(&p.Curr)
		} else {
			p.Curr = p.Next[len(p.Next)-1]
			p.Next = p.Next[:len(p.Next)-1]
		}
	}
}

func (p *Provider) CheckNotEOF() (err error) {
	if p.Curr.IsEOF {
		err = p.Unexpected(p.Curr)
	}
	return
}

func (p *Provider) PosStruct(ofs uint) (ps PosStruct, err error) {
	if ofs > 0 {
		if len(p.Next) >= int(ofs) {
			ps = p.Next[len(p.Next)-int(ofs)]
		} else {
			if len(p.Next) > 0 {
				ps = p.Next[0]
			} else {
				ps = p.Curr
			}
			lookahead := PosStructs{}
			for i := int(ofs) - len(p.Next); i > 0; i-- {
				p.pullRune(&ps)
				lookahead = append(lookahead, ps)
				if ps.IsEOF {
					break
				}
			}
			newNext := PosStructs{}
			for i := len(lookahead) - 1; i >= 0; i-- {
				newNext = append(newNext, lookahead[i])
			}
			p.Next = append(newNext, p.Next...)
		}
	} else {
		ps = p.Curr
	}
	return
}

func (p *Provider) Unexpected(ps PosStruct, optFmt ...bw.I) (result error) {
	var msg string
	if ps.IsEOF {
		msg = ansiUnexpectedEOF
	} else if len(optFmt) == 0 {
		r := ps.Rune
		msg = fmt.Sprintf(ansiUnexpectedChar, r, r)
	} else {
		msg = bw.Spew.Sprintf(optFmt[0].FmtString(), optFmt[0].FmtArgs()...)
	}
	result = bwerr.From(msg + p.suffix(ps))
	return
}

// ============================================================================

var (
	ansiPosStructFailed string
	ansiOK              string
	ansiErr             string
	ansiPos             string
	ansiLineCol         string
	ansiGetSuffixAssert string
	ansiUnexpectedEOF   string
	ansiUnexpectedChar  string
)

func init() {
	ansiOK = ansi.CSIFromSGRCodes(ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen, Bright: false})).String()
	ansiErr = ansi.CSIFromSGRCodes(ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorRed, Bright: true})).String()

	ansiPosStructFailed = "<ansiPath>bwparse.Provider.<ansiFunc>PosStruct<ansi>(%d) failed, must <ansiFunc>.SetMaxBackwardCount<ansi>(<ansiVal>%d<ansi>)"
	ansiPos = ansi.String(" at pos <ansiPath>%d<ansi>")
	ansiLineCol = ansi.String(" at line <ansiPath>%d<ansi>, col <ansiPath>%d<ansi> (pos <ansiPath>%d<ansi>)")
	ansiGetSuffixAssert = ansi.String("<ansiVar>ps.Pos<ansi> (<ansiVal>%d<ansi>) > <ansiVar>p.Curr.Pos<ansi> (<ansiVal>%d<ansi>)")
	ansiUnexpectedEOF = ansi.String("unexpected end of string")
	ansiUnexpectedChar = ansi.String("unexpected char <ansiVal>%q<ansiReset> (<ansiVar>charCode<ansi>: <ansiVal>%d<ansi>)")
}

func (p *Provider) pullRune(ps *PosStruct) {
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

func (p *Provider) suffix(ps PosStruct) (suffix string) {
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
