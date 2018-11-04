package runeprovider

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
)

type RuneProvider interface {
	PullRune() (*rune, error)
	Close()
	Line() int
	Col() int
	Pos() int
	IsEOF() bool
}

type PosStruct struct {
	IsEOF       bool
	RunePtr     *rune
	Pos         int
	Line        uint
	Col         uint
	Prefix      string
	PrefixStart int
}

// func (v PosStruct) copyPtr() *PosStruct {
// 	return &PosStruct{v.IsEOF, v.RunePtr, v.Pos, v.Line, v.Col, v.Prefix, v.PrefixStart}
// }

func (v PosStruct) DataForJSON() interface{} {
	result := map[string]interface{}{}
	if v.RunePtr == nil {
		result["rune"] = "EOF"
	} else {
		result["rune"] = string(*(v.RunePtr))
	}
	result["line"] = v.Line
	result["col"] = v.Col
	result["pos"] = v.Pos
	// result["prefix"] = v.Prefix
	// result["prefixStart"] = v.PrefixStart
	return result
}

type Proxy struct {
	Curr               PosStruct
	Prev               RunePtrStructs
	Next               RunePtrStructs
	Prov               RuneProvider
	preLineCount       int
	postLineCount      int
	maxBehindRuneCount int
}

func ProxyFrom(p RuneProvider) *Proxy {
	return &Proxy{
		PosStruct{Pos: -1, Line: 1},
		RunePtrStructs{},
		RunePtrStructs{},
		p,
		3,
		3,
		2,
	}
}

type RunePtrStructs []PosStruct

func (v RunePtrStructs) DataForJSON() interface{} {
	result := []interface{}{}
	for _, i := range v {
		result = append(result, i.DataForJSON())
	}
	return result
}

func (p *Proxy) DataForJSON() interface{} {
	result := map[string]interface{}{}
	result["curr"] = p.Curr.DataForJSON()
	// if len(p.Prev) > 0 {
	// 	result["prev"] = p.Prev.DataForJSON()
	// }
	// if len(p.Next) > 0 {
	// 	result["next"] = p.Next.DataForJSON()
	// }
	return result
}

func (p *Proxy) PullRune() {
	if p.Curr.Pos < 0 || p.Curr.RunePtr != nil {
		p.Prev = append(p.Prev, p.Curr)
		if len(p.Prev) > p.maxBehindRuneCount {
			p.Prev = p.Prev[len(p.Prev)-p.maxBehindRuneCount:]
		}
		if len(p.Next) > 0 {
			p.Curr = p.Next[len(p.Next)-1]
			p.Next = p.Next[:len(p.Next)-1]
		} else {
			p.pullRune(&p.Curr)
		}
	}
}

func (p *Proxy) pullRune(ps *PosStruct) {
	runePtr, err := p.Prov.PullRune()
	if err != nil {
		bwerr.PanicA(bwerr.E{Error: err})
	}
	ps.Pos++
	if runePtr != nil && ps.RunePtr != nil {
		if *(ps.RunePtr) != '\n' {
			ps.Col++
		} else {
			ps.Line++
			ps.Col = 1
			if int(ps.Line) > p.preLineCount {
				i := strings.Index(ps.Prefix, "\n")
				ps.Prefix = ps.Prefix[i+1:]
				ps.PrefixStart += i + 1
			}
		}
	}
	ps.RunePtr = runePtr
	ps.IsEOF = runePtr == nil
	if !ps.IsEOF {
		ps.Prefix += string(*runePtr)
	}
}

func (p *Proxy) PushRune() {
	if len(p.Prev) == 0 {
		bwerr.Panic("len(p.Prev) == 0")
	} else {
		p.Next = append(p.Next, p.Curr)
		p.Curr = p.Prev[len(p.Prev)-1]
		p.Prev = p.Prev[:len(p.Prev)-1]
	}
}

func (p *Proxy) Rune(optOfs ...int) (result rune, isEOF bool) {
	ps := p.PosStruct(optOfs...)
	if ps.RunePtr == nil {
		result = '\000'
		isEOF = true
	} else {
		result = *ps.RunePtr
		isEOF = false
	}
	return
}

func (p *Proxy) PosStruct(optOfs ...int) (ps PosStruct) {
	var ofs int
	if optOfs != nil {
		ofs = optOfs[0]
	}
	if ofs < 0 {
		if len(p.Prev) >= -ofs {
			ps = p.Prev[len(p.Prev)+ofs]
		} else {
			bwerr.Panic("len(p.Prev) < %d", -ofs)
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
			lookahead := RunePtrStructs{}
			for i := ofs - len(p.Next); i > 0; i-- {
				p.pullRune(&ps)
				lookahead = append(lookahead, ps)
				if ps.IsEOF {
					break
				}
			}
			newNext := RunePtrStructs{}
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

func init() {
	ansi.MustAddTag("ansiDarkGreen",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdFaint),
		// ansi.SGRCodeOfColor256(ansi.Color256{Code: 201}),
		// ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen, Bright: true}),
	)
	ansi.MustAddTag("ansiLightRed",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorRed, Bright: true}),
		// ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen}),
		// ansi.MustSGRCodeOfCmd(ansi.SGRCmdFaint),
		// ansi.SGRCodeOfColor256(ansi.Color256{Code: 201}),
		// ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen, Bright: true}),
	)
	// ansi.MustAddTag("ansiDarkGreen",
	// ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen}),
	// ansi.MustSGRCodeOfCmd(ansi.SGRCmdFaint),
}

func (p *Proxy) GetSuffix(ps PosStruct) (suffix string) {
	preLineCount := p.preLineCount
	postLineCount := p.postLineCount
	if p.Curr.RunePtr == nil {
		preLineCount += postLineCount
	}

	separator := "\n"
	if p.Curr.Line <= 1 {
		suffix += fmt.Sprintf(" at pos <ansiPath>%d<ansi>", ps.Pos)
		separator = " "
	} else {
		suffix += fmt.Sprintf(" at line <ansiPath>%d<ansi>, col <ansiPath>%d<ansi> (pos <ansiPath>%d<ansi>)", ps.Line, ps.Col, ps.Pos)
	}
	suffix += ":" + separator + "<ansiDarkGreen>"

	// if p.Curr.Pos == ps.Pos {
	// suffix += p.Curr.Prefix[0 : ps.Pos-p.Curr.PrefixStart]
	// }
	suffix += p.Curr.Prefix[0 : ps.Pos-p.Curr.PrefixStart]
	if p.Curr.RunePtr != nil {
		suffix += "<ansiLightRed>"
		suffix += p.Curr.Prefix[ps.Pos-p.Curr.PrefixStart:]
		// suffix += redString
		suffix += "<ansiReset>"
		for p.Curr.RunePtr != nil && postLineCount > 0 {
			p.PullRune()
			if p.Curr.RunePtr != nil {
				suffix += string(*p.Curr.RunePtr)
				if *p.Curr.RunePtr == '\n' {
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

// func (p *Proxy) UnexpectedRuneError(infix ...string) error {
// 	var fmtString string
// 	fmtArgs := []interface{}{}
// 	if p.Curr.RunePtr == nil {
// 		suffix := p.GetSuffix(p.Curr)
// 		fmtString = "unexpected end of string"
// 		if infix != nil {
// 			fmtString += "(" + strings.Join(infix, " ") + ")"
// 		}
// 		fmtString += suffix
// 	} else {
// 		rune := *p.Curr.RunePtr
// 		suffix := p.GetSuffix(p.Curr)
// 		fmtString = "unexpected char <ansiVal>%q<ansiReset> (charCode: %v"
// 		if infix != nil {
// 			fmtString += ", " + strings.Join(infix, " ")
// 		}
// 		fmtString += ")" + suffix
// 		fmtArgs = []interface{}{rune, rune}
// 	}
// 	return bwerr.Errord(1, fmtString, fmtArgs...)
// }

// type FmtStruct struct {
// 	FmtString string
// 	FmtArgs   []interface{}
// }

// func FmtStructFrom(fmtString string, fmtArgs ...interface{}) FmtStruct {
// 	return FmtStruct{fmtString, fmtArgs}
// }

// func (p *Proxy) Unexpected(ps PosStruct, fmtString string, fmtArgs ...interface{}) error {
func (p *Proxy) Unexpected(ps PosStruct, optFmtStruct ...bw.A) (err error) {
	var fmtString string
	fmtArgs := []interface{}{}
	if p.Curr.RunePtr == nil {
		suffix := p.GetSuffix(p.Curr)
		// if fmtString
		fmtString = "unexpected end of string"
		if optFmtStruct != nil {
			fmtString += "(" + optFmtStruct[0].Fmt + ")"
			fmtArgs = append(fmtArgs, optFmtStruct[0].Args...)
		}
		fmtString += suffix
	} else if ps.Pos == p.Curr.Pos {
		r := *p.Curr.RunePtr
		suffix := p.GetSuffix(p.Curr)
		fmtString = "unexpected char <ansiVal>%q<ansiReset> (charCode: %v"
		fmtArgs = []interface{}{r, r}
		if optFmtStruct != nil {
			fmtString += ", " + optFmtStruct[0].Fmt
			fmtArgs = append(fmtArgs, optFmtStruct[0].Args...)
		}
		fmtString += ")" + suffix
	} else if ps.Pos < p.Curr.Pos {
		if optFmtStruct != nil {
			fmtString += optFmtStruct[0].Fmt
			fmtArgs = append(fmtArgs, optFmtStruct[0].Args...)
		}
		fmtString += p.GetSuffix(ps)
	} else {
		bwerr.Panic("ps.Pos: %#v, p.Curr.Pos: %#v", ps.Pos, p.Curr.Pos)
	}
	// bwerr.Panic("%#v, fmtArgs: %#v", fmtString, fmtArgs)
	err = bwerr.FromA(bwerr.A{1, fmtString, fmtArgs})
	// bwerr.Panic("%#v", err)
	return
}

// func (p *Proxy) ItemError(start PosStruct, fmtString string, fmtArgs ...interface{}) error {
// 	fmtString += p.GetSuffix(start)
// 	return bwerr.Errord(1, fmtString, fmtArgs...)
// }

// func (p *Proxy) unknownWordError(start int, word string) {
// 	// suffix := pfa.p.GetSuffix(start, word)
// 	return p.wordError("unknown word <ansiVal>%s<ansi>", word, start)
// }

func FromString(source string) RuneProvider {
	result := stringRuneProvider{pos: -1, src: []rune(source)}
	return &result
}

func FromFile(fileSpec string) (result RuneProvider, err error) {
	p := &fileRuneProvider{fileSpec: fileSpec, pos: -1, bytePos: -1, line: 1}
	p.data, err = os.Open(fileSpec)
	if err == nil {
		p.reader = bufio.NewReader(p.data)
		result = p
	}
	return
}

type stringRuneProvider struct {
	src  []rune
	line int
	col  int
	pos  int
}

func (v *stringRuneProvider) PullRune() (result *rune, err error) {
	v.pos++
	if v.pos < len(v.src) {
		currRune := v.src[v.pos]
		result = &currRune
		if currRune == '\n' {
			v.line++
			v.col = 0
		} else {
			v.col++
		}
	}
	return
}

func (v *stringRuneProvider) Close() {
	v.src = nil
}

func (v *stringRuneProvider) Pos() int {
	return v.pos
}

func (v *stringRuneProvider) Line() int {
	return v.line
}

func (v *stringRuneProvider) Col() int {
	return v.col
}

func (v *stringRuneProvider) IsEOF() bool {
	return v.pos >= len(v.src)
}

const chunksize int = 1024

type fileRuneProvider struct {
	fileSpec string
	data     *os.File
	buf      []byte
	reader   *bufio.Reader
	pos      int
	line     int
	col      int
	bytePos  int
	isEOF    bool
}

func (v *fileRuneProvider) PullRune() (result *rune, err error) {
	// var size int
	if len(v.buf) < utf8.UTFMax && !v.isEOF {
		chunk := make([]byte, chunksize)
		var count int
		count, err = v.reader.Read(chunk)
		if err == io.EOF {
			v.isEOF = true
			err = nil
		} else if err == nil {
			v.buf = append(v.buf, chunk[:count]...)
		}
	}
	if err == nil {
		if len(v.buf) != 0 {
			currRune, size := utf8.DecodeRune(v.buf)
			if currRune == utf8.RuneError {
				err = bwerr.From("utf-8 encoding is invalid at pos %d (byte #%d)", v.pos, v.bytePos)
			} else {
				result = &currRune
				v.buf = v.buf[size:]
				v.pos++
				v.bytePos += size
				if currRune == '\n' {
					v.line++
					v.col = 0
				} else {
					v.col++
				}
			}
		}
	}
	return
}

func (v *fileRuneProvider) Close() {
	v.data.Close()
}

func (v *fileRuneProvider) Pos() int {
	return v.pos
}

func (v *fileRuneProvider) Line() int {
	return v.line
}

func (v *fileRuneProvider) Col() int {
	return v.col
}

func (v *fileRuneProvider) IsEOF() bool {
	return v.isEOF && len(v.buf) == 0
}
