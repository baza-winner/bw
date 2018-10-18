package runeprovider

import (
	"bufio"
	"io"
	"os"
	"unicode/utf8"

	"github.com/baza-winner/bwcore/bwerror"
)

type RuneProvider interface {
	PullRune() (*rune, error)
	Close()
	Line() int
	Col() int
	Pos() int
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
