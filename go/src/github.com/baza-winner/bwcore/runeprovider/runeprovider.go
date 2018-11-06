package runeprovider

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/baza-winner/bwcore/ansi"
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

func (v PosStruct) MarshalJSON() ([]byte, error) {
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
	return json.Marshal(result)
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

func (v RunePtrStructs) MarshalJSON() ([]byte, error) {
	result := []interface{}{}
	for _, i := range v {
		result = append(result, i)
	}
	return json.Marshal(result)
}

func (p *Proxy) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{}
	result["curr"] = p.Curr
	// if len(p.Prev) > 0 {
	// 	result["prev"] = p.Prev.DataForJSON()
	// }
	// if len(p.Next) > 0 {
	// 	result["next"] = p.Next.DataForJSON()
	// }
	return json.Marshal(result)
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

func (p *Proxy) pullRune(ps *PosStruct) (err error) {
	var runePtr *rune
	runePtr, err = p.Prov.PullRune()
	if err != nil {
		return
		// bwerr.PanicA(bwerr.Err(err))
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
	return
}

func (p *Proxy) PushRune() (err error) {
	if len(p.Prev) == 0 {
		err = bwerr.From("len(p.Prev) == 0")
		// bwerr.Panic("len(p.Prev) == 0")
	} else {
		p.Next = append(p.Next, p.Curr)
		p.Curr = p.Prev[len(p.Prev)-1]
		p.Prev = p.Prev[:len(p.Prev)-1]
	}
	return
}

func (p *Proxy) Rune(optOfs ...int) (result rune, isEOF bool, err error) {
	var ps PosStruct
	ps, err = p.PosStruct(optOfs...)
	if err != nil {
		return
	}
	if ps.RunePtr == nil {
		result = '\000'
		isEOF = true
	} else {
		result = *ps.RunePtr
		isEOF = false
	}
	return
}

func (p *Proxy) PosStruct(optOfs ...int) (ps PosStruct, err error) {
	var ofs int
	if optOfs != nil {
		ofs = optOfs[0]
	}
	if ofs < 0 {
		if len(p.Prev) >= -ofs {
			ps = p.Prev[len(p.Prev)+ofs]
		} else {
			err = bwerr.From("len(p.Prev) < %d", -ofs)
			return
			// bwerr.Panic("len(p.Prev) < %d", -ofs)
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

// func init() {
// 	ansi.MustAddTag("ansiDarkGreen",
// 		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen}),
// 		ansi.MustSGRCodeOfCmd(ansi.SGRCmdFaint),
// 		// ansi.SGRCodeOfColor256(ansi.Color256{Code: 201}),
// 		// ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen, Bright: true}),
// 	)
// 	ansi.MustAddTag("ansiLightRed",
// 		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorRed, Bright: true}),
// 		// ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen}),
// 		// ansi.MustSGRCodeOfCmd(ansi.SGRCmdFaint),
// 		// ansi.SGRCodeOfColor256(ansi.Color256{Code: 201}),
// 		// ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen, Bright: true}),
// 	)
// 	// ansi.MustAddTag("ansiDarkGreen",
// 	// ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen}),
// 	// ansi.MustSGRCodeOfCmd(ansi.SGRCmdFaint),
// }

var (
	ansiOK      string
	ansiErr     string
	ansiPos     string
	ansiLineCol string
)

func init() {
	ansiOK = ansi.CSIFromSGRCodes(ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen, Bright: false})).String()
	ansiErr = ansi.CSIFromSGRCodes(ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorRed, Bright: true})).String()
	ansiPos = ansi.String(" at pos <ansiPath>%d<ansi>")
	ansiLineCol = ansi.String(" at line <ansiPath>%d<ansi>, col <ansiPath>%d<ansi> (pos <ansiPath>%d<ansi>)")
}

func (p *Proxy) GetSuffix(ps PosStruct) (suffix string) {
	preLineCount := p.preLineCount
	postLineCount := p.postLineCount
	if p.Curr.RunePtr == nil {
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
	if p.Curr.RunePtr != nil {
		suffix += ansiErr
		suffix += p.Curr.Prefix[ps.Pos-p.Curr.PrefixStart:]
		suffix += ansi.Reset()
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

var (
	ansiUnexpectedEOF  string
	ansiUnexpectedChar string
)

func init() {
	ansiUnexpectedEOF = ansi.String("unexpected end of string")
	ansiUnexpectedChar = ansi.String("unexpected char <ansiVal>%q<ansiReset> (<ansiVar>charCode<ansi>: <ansiVal>%d<ansi>)")
}

func (p *Proxy) Unexpected(ps PosStruct) (result error) {
	var msg string
	if ps.RunePtr == nil {
		msg = ansiUnexpectedEOF
	} else {
		r := *ps.RunePtr
		msg = fmt.Sprintf(ansiUnexpectedChar, r, r)
	}
	result = bwerr.From(msg + p.GetSuffix(ps))
	return
}

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

var (
	ansiInvalidByte string
)

func init() {
	ansiInvalidByte = ansi.String("utf-8 encoding <ansiVal>%#v<ansi> is invalid at pos <ansiPath>%d<ansi>")
}

func (v *fileRuneProvider) PullRune() (result *rune, err error) {
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
				err = bwerr.From(ansiInvalidByte, v.bytePos, v.pos)
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
