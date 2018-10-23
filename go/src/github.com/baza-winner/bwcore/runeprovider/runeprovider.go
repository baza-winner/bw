package runeprovider

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/baza-winner/bwcore/bwerror"
)

type RuneProvider interface {
	PullRune() (*rune, error)
	Close()
	Line() int
	Col() int
	Pos() int
	IsEOF() bool
}

type RunePtrStruct struct {
	IsEOF       bool
	RunePtr     *rune
	Pos         int
	Line        uint
	Col         uint
	Prefix      string
	PrefixStart int
}

func (v RunePtrStruct) copyPtr() *RunePtrStruct {
	return &RunePtrStruct{v.IsEOF, v.RunePtr, v.Pos, v.Line, v.Col, v.Prefix, v.PrefixStart}
}

func (v RunePtrStruct) DataForJSON() interface{} {
	result := map[string]interface{}{}
	if v.RunePtr == nil {
		result["rune"] = "EOF"
	} else {
		result["rune"] = string(*(v.RunePtr))
	}
	result["line"] = v.Line
	result["col"] = v.Col
	result["pos"] = v.Pos
	result["prefix"] = v.Prefix
	result["prefixStart"] = v.PrefixStart
	return result
}

type Proxy struct {
	Prev          *RunePtrStruct
	Curr          RunePtrStruct
	Next          *RunePtrStruct
	Prov          RuneProvider
	preLineCount  int
	postLineCount int
}

func ProxyFrom(p RuneProvider) *Proxy {
	return &Proxy{
		Prov:          p,
		Curr:          RunePtrStruct{Pos: -1, Line: 1},
		preLineCount:  3,
		postLineCount: 3,
	}
}

func (p *Proxy) DataForJSON() interface{} {
	result := map[string]interface{}{}
	result["curr"] = p.Curr.DataForJSON()
	if p.Prev != nil {
		result["prev"] = p.Prev.DataForJSON()
	}
	if p.Next != nil {
		result["next"] = p.Next.DataForJSON()
	}
	return result
}

func (p *Proxy) PullRune() {
	if p.Curr.Pos < 0 || p.Curr.RunePtr != nil {
		p.Prev = p.Curr.copyPtr()
		if p.Next != nil {
			p.Curr = *(p.Next)
			p.Next = nil
		} else {
			runePtr, err := p.Prov.PullRune()
			if err != nil {
				bwerror.PanicErr(err)
			}
			pos := p.Prev.Pos + 1
			line := p.Prev.Line
			col := p.Prev.Col
			prefix := p.Prev.Prefix
			prefixStart := p.Prev.PrefixStart
			if runePtr != nil && p.Prev.RunePtr != nil {
				if *(p.Prev.RunePtr) != '\n' {
					col += 1
				} else {
					line += 1
					col = 1
					if int(line) > p.preLineCount {
						i := strings.Index(prefix, "\n")
						prefix = prefix[i+1:]
						prefixStart += i + 1
					}
				}
			}
			isEOF := runePtr == nil
			if !isEOF {
				prefix += string(*runePtr)
			}
			p.Curr = RunePtrStruct{isEOF, runePtr, pos, line, col, prefix, prefixStart}
		}
	}
}

func (p *Proxy) PushRune() {
	if p.Prev == nil {
		bwerror.Panic("p.Prev == nil")
	} else {
		p.Next = p.Curr.copyPtr()
		p.Curr = *(p.Prev)
	}
}

func (p *Proxy) Rune() (result rune, isEOF bool) {
	if p.Curr.RunePtr == nil {
		result = '\000'
		isEOF = true
	} else {
		result = *p.Curr.RunePtr
		isEOF = false
	}
	return
}

func (p *Proxy) GetSuffix(start RunePtrStruct, redString string) (suffix string) {
	preLineCount := p.preLineCount
	postLineCount := p.postLineCount
	if p.Curr.RunePtr == nil {
		preLineCount += postLineCount
	}

	separator := "\n"
	if p.Curr.Line <= 1 {
		suffix += fmt.Sprintf(" at pos <ansiCmd>%d<ansi>", start.Pos)
		separator = " "
	} else {
		suffix += fmt.Sprintf(" at line <ansiCmd>%d<ansi>, col <ansiCmd>%d<ansi> (pos <ansiCmd>%d<ansi>)", start.Line, start.Col, start.Pos)
	}
	suffix += ":" + separator + "<ansiDarkGreen>"

	suffix += p.Curr.Prefix[0 : start.Pos-p.Curr.PrefixStart]
	if p.Curr.RunePtr != nil {
		suffix += "<ansiLightRed>"
		suffix += redString
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

func (p *Proxy) UnexpectedRuneError(infix ...string) error {
	var fmtString string
	fmtArgs := []interface{}{}
	if p.Curr.RunePtr == nil {
		suffix := p.GetSuffix(p.Curr, "")
		fmtString = "unexpected end of string"
		if infix != nil {
			fmtString += "(" + strings.Join(infix, " ") + ")"
		}
		fmtString += suffix
	} else {
		rune := *p.Curr.RunePtr
		suffix := p.GetSuffix(p.Curr, string(rune))
		fmtString = "unexpected char <ansiPrimary>%q<ansiReset> (charCode: %v"
		if infix != nil {
			fmtString += ", " + strings.Join(infix, " ")
		}
		fmtString += ")" + suffix
		fmtArgs = []interface{}{rune, rune}
	}
	return bwerror.Error(fmtString, fmtArgs...)
}

func (p *Proxy) WordError(fmtString string, word string, start RunePtrStruct) error {
	suffix := p.GetSuffix(start, word)
	return bwerror.Error(fmtString+suffix, word)
}

// func (p *Proxy) unknownWordError(start int, word string) {
// 	// suffix := pfa.p.GetSuffix(start, word)
// 	return p.wordError("unknown word <ansiPrimary>%s<ansi>", word, start)
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
				err = bwerror.Error("utf-8 encoding is invalid at pos %d (byte #%d)", v.pos, v.bytePos)
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
