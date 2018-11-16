package bwparse

import (
	"encoding/json"
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

func (v PosStruct) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{}
	if v.IsEOF {
		result["rune"] = "EOF"
	} else {
		result["rune"] = string(v.Rune)
	}
	result["line"] = v.Line
	result["col"] = v.Col
	result["pos"] = v.Pos
	// result["prefix"] = v.Prefix
	// result["prefixStart"] = v.PrefixStart
	return json.Marshal(result)
}

func (v PosStructs) MarshalJSON() ([]byte, error) {
	result := []interface{}{}
	for _, i := range v {
		result = append(result, i)
	}
	return json.Marshal(result)
}

// ============================================================================

type Provider struct {
	Curr          PosStruct
	Prev          PosStructs
	Next          PosStructs
	Prov          bwrune.Provider
	preLineCount  uint
	postLineCount uint
	maxBackward   uint
}

func ProviderFrom(p bwrune.Provider, optMaxBackward ...uint) (result *Provider) {
	result = &Provider{
		Curr:          PosStruct{Pos: -1, Line: 1},
		Prev:          PosStructs{},
		Next:          PosStructs{},
		Prov:          p,
		preLineCount:  3,
		postLineCount: 3,
		maxBackward:   1,
	}
	if len(optMaxBackward) > 0 {
		result.maxBackward = optMaxBackward[0]
	}
	return
}

// ============================================================================

func (p *Provider) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{}
	result["curr"] = p.Curr
	// if len(p.Prev) > 0 {
	//  result["prev"] = p.Prev.DataForJSON()
	// }
	// if len(p.Next) > 0 {
	//  result["next"] = p.Next.DataForJSON()
	// }
	return json.Marshal(result)
}

const NonEOF = true

func (p *Provider) SetMaxBackwardCount(newValue uint) (prev uint) {
	prev = p.maxBackward
	p.maxBackward = newValue
	if len(p.Prev) > int(p.maxBackward) {
		p.Prev = p.Prev[len(p.Prev)-int(p.maxBackward):]
	}
	return
}

func (p *Provider) Forward() {
	if p.Curr.Pos < 0 || !p.Curr.IsEOF {
		if p.maxBackward > 0 {
			if p.maxBackward == 1 && len(p.Prev) == 1 {
				p.Prev[0] = p.Curr
			} else {
				if len(p.Prev) >= int(p.maxBackward) {
					p.Prev = p.Prev[len(p.Prev)-int(p.maxBackward)+1:]
				}
				p.Prev = append(p.Prev, p.Curr)
			}
		}
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

func (p *Provider) pullRune(ps *PosStruct) {
	var runePtr *rune
	var err error
	if runePtr, err = p.Prov.PullRune(); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	ps.Pos++
	if runePtr != nil && !ps.IsEOF {
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
	ps.IsEOF = runePtr == nil
	if runePtr == nil {
		ps.Rune = '\000'
	} else {
		ps.Rune = *runePtr
		ps.Prefix += string(*runePtr)
	}
	return
}

var ansiBackwardFailed string

func init() {
	ansiBackwardFailed = "<ansiPath>bwparse.Provider.<ansiFunc>Backward<ansi>() failed, must <ansiFunc>.SetMaxBackwardCount<ansi>(<ansiVal>%d<ansi>)"
}

func (p *Provider) Backward() {
	if len(p.Prev) == 0 {
		bwerr.Panic(ansiBackwardFailed, p.maxBackward+1)
	} else {
		p.Next = append(p.Next, p.Curr)
		p.Curr = p.Prev[len(p.Prev)-1]
		p.Prev = p.Prev[:len(p.Prev)-1]
	}
	return
}

var ansiPosStructFailed string

func init() {
	ansiPosStructFailed = "<ansiPath>bwparse.Provider.<ansiFunc>PosStruct<ansi>(%d) failed, must <ansiFunc>.SetMaxBackwardCount<ansi>(<ansiVal>%d<ansi>)"
}

func (p *Provider) PosStruct(optOfs ...int) (ps PosStruct, err error) {
	var ofs int
	if optOfs != nil {
		ofs = optOfs[0]
	}
	if ofs < 0 {
		if len(p.Prev) >= -ofs {
			ps = p.Prev[len(p.Prev)+ofs]
		} else {
			err = bwerr.From(ansiPosStructFailed, ofs, -ofs+1)
			return
		}
	} else if ofs > 0 {
		if len(p.Next) >= ofs {
			ps = p.Next[len(p.Next)-ofs]
		} else {
			if len(p.Next) > 0 {
				ps = p.Next[0]
			} else {
				ps = p.Curr
			}
			lookahead := PosStructs{}
			for i := ofs - len(p.Next); i > 0; i-- {
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

var (
	ansiOK              string
	ansiErr             string
	ansiPos             string
	ansiLineCol         string
	ansiGetSuffixAssert string
)

func init() {
	ansiOK = ansi.CSIFromSGRCodes(ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen, Bright: false})).String()
	ansiErr = ansi.CSIFromSGRCodes(ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorRed, Bright: true})).String()
	ansiPos = ansi.String(" at pos <ansiPath>%d<ansi>")
	ansiLineCol = ansi.String(" at line <ansiPath>%d<ansi>, col <ansiPath>%d<ansi> (pos <ansiPath>%d<ansi>)")
	ansiGetSuffixAssert = ansi.String("<ansiVar>ps.Pos<ansi> (<ansiVal>%d<ansi>) > <ansiVar>p.Curr.Pos<ansi> (<ansiVal>%d<ansi>)")
}

func (p *Provider) GetSuffix(ps PosStruct) (suffix string) {
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
	if !p.Curr.IsEOF {
		suffix += ansiErr
		suffix += p.Curr.Prefix[ps.Pos-p.Curr.PrefixStart:]
		suffix += ansi.Reset()
		for !p.Curr.IsEOF && postLineCount > 0 {
			p.Forward()
			if !p.Curr.IsEOF {
				suffix += string(p.Curr.Rune)
				if p.Curr.Rune == '\n' {
					postLineCount -= 1
				}
			}
		}
	}
	if byte(suffix[len(suffix)-1]) != '\n' {
		suffix += string('\n')
	}
	return suffix
}

var (
	ansiUnexpectedEOF  string
	ansiUnexpectedChar string
)

func init() {
	ansiUnexpectedEOF = ansi.String("unexpected end of string")
	ansiUnexpectedChar = ansi.String("unexpected char <ansiVal>%q<ansiReset> (<ansiVar>charCode<ansi>: <ansiVal>%d<ansi>)")
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
	result = bwerr.From(msg + p.GetSuffix(ps))
	return
}

// ============================================================================
