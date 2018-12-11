package bwrune

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerr"
)

// ============================================================================

type Provider interface {
	FileSpec() string
	Close() error
	PullRune() (*rune, error)
	Line() int
	Col() int
	Pos() int
	IsEOF() bool
}

// ============================================================================

type ProviderProvider interface {
	Provider() Provider
}

// ============================================================================

type S struct {
	S string
}

func (v S) Provider() Provider {
	return FromString(v.S)
}

// ============================================================================

type F struct {
	S string
}

func (v F) Provider() Provider {
	return MustFromFile(v.S)
}

// ============================================================================

func FromString(source string) Provider {
	result := stringProvider{pos: -1, src: []rune(source)}
	return &result
}

func FromFile(fileSpec string) (result Provider, err error) {
	if fileSpec, err = filepath.Abs(fileSpec); err != nil {
		return
	}
	p := &fileProvider{fileSpec: fileSpec, pos: -1, bytePos: -1, line: 1}
	p.data, err = os.Open(fileSpec)
	if err == nil {
		p.reader = bufio.NewReader(p.data)
		result = p
	}
	return
}

func MustFromFile(fileSpec string) (result Provider) {
	var err error
	if result, err = FromFile(fileSpec); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func MustPull(p Provider) (result *rune) {
	var err error
	if result, err = p.PullRune(); err != nil {
		bwerr.PanicErr(err)
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

func (v *stringProvider) FileSpec() string {
	return ""
}

func (v *stringProvider) Close() error {
	v.src = nil
	return nil
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

func (v *fileProvider) FileSpec() string {
	return v.fileSpec
}

func (v *fileProvider) Close() error {
	return v.data.Close()
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
