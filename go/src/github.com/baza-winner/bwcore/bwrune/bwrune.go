package bwrune

import (
	"bufio"
	"io"
	"os"
	"unicode/utf8"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerr"
)

// ============================================================================

type Provider interface {
	PullRune() (*rune, error)
	Close()
	Line() int
	Col() int
	Pos() int
	IsEOF() bool
}

// ============================================================================

func ProviderFromString(source string) Provider {
	result := stringProvider{pos: -1, src: []rune(source)}
	return &result
}

func ProviderFromFile(fileSpec string) (result Provider, err error) {
	p := &fileProvider{fileSpec: fileSpec, pos: -1, bytePos: -1, line: 1}
	p.data, err = os.Open(fileSpec)
	if err == nil {
		p.reader = bufio.NewReader(p.data)
		result = p
	}
	return
}

// ============================================================================

type stringProvider struct {
	src  []rune
	line int
	col  int
	pos  int
}

func (v *stringProvider) PullRune() (result *rune, err error) {
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

func (v *stringProvider) Close() {
	v.src = nil
}

func (v *stringProvider) Pos() int {
	return v.pos
}

func (v *stringProvider) Line() int {
	return v.line
}

func (v *stringProvider) Col() int {
	return v.col
}

func (v *stringProvider) IsEOF() bool {
	return v.pos >= len(v.src)
}

const chunksize int = 1024

type fileProvider struct {
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
	ansiInvalidByte = ansi.String("utf-8 encoding <ansiVal>%#v<ansi> is invalid at pos <ansiPath>%d<ansi> of file <ansiPath>%s")
}

func (v *fileProvider) PullRune() (result *rune, err error) {
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
				err = bwerr.From(ansiInvalidByte, v.bytePos, v.pos, v.fileSpec)
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

func (v *fileProvider) Close() {
	v.data.Close()
}

func (v *fileProvider) Pos() int {
	return v.pos
}

func (v *fileProvider) Line() int {
	return v.line
}

func (v *fileProvider) Col() int {
	return v.col
}

func (v *fileProvider) IsEOF() bool {
	return v.isEOF && len(v.buf) == 0
}

// ============================================================================
