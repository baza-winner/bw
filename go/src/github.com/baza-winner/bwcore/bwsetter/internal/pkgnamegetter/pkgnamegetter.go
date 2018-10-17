package pkgnamegetter

import (
	"bufio"
	"io"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/baza-winner/bwcore/bwerror"
)

// func GetPackageName(packageDir string) (packageName string, err error) {

// }

type RuneProvider struct {
	fileSpec string
	data     *os.File
	buf      []byte
	reader   *bufio.Reader
	pos      int
	line     int
	col      int
	bytePos  int
	isEof    bool
}

const (
	chunksize int = 1024
)

func InitRuneProvider(fileSpec string) (result *RuneProvider, err error) {
	result = &RuneProvider{fileSpec: fileSpec, pos: -1, bytePos: -1, line: 1}
	result.data, err = os.Open(fileSpec)
	if err == nil {
		result.reader = bufio.NewReader(result.data)
	}
	return
}

func (v *RuneProvider) Close() {
	v.data.Close()
}

func (v *RuneProvider) PullRune() (currRune rune, isEof bool, err error) {
	var size int
	if len(v.buf) < utf8.UTFMax && !v.isEof {
		chunk := make([]byte, chunksize)
		var count int
		count, err = v.reader.Read(chunk)
		if err == io.EOF {
			v.isEof = true
			err = nil
		} else if err == nil {
			v.buf = append(v.buf, chunk[:count]...)
		}
	}
	if err == nil {
		if len(v.buf) == 0 {
			isEof = true
		} else {
			currRune, size = utf8.DecodeRune(v.buf)
			if currRune == utf8.RuneError {
				err = bwerror.Error("utf-8 encoding is invalid at pos %d (byte #%d)", v.pos, v.bytePos)
			} else {
				v.buf = v.buf[size:]
				v.pos += 1
				v.bytePos += size
				if currRune == '\n' {
					v.line += 1
					v.col = 0
				} else {
					v.col += 1
				}
			}
		}
	}
	return
}

func getFirstLine(fileSpec string) (result string, err error) {
	var p *RuneProvider
	p, err = InitRuneProvider(fileSpec)
	if err != nil {
		err = bwerror.ErrorErr(err)
	} else {
		defer p.Close()
		for {
			var currRune rune
			var isEof bool
			currRune, isEof, err = p.PullRune()
			if isEof || err != nil || currRune == '\n' {
				break
			} else {
				result += string(currRune)
			}
		}
	}
	return
}

type parsePrimaryState uint8

const (
	pps_below_ parsePrimaryState = iota
	ppsSeekPackage
	ppsSeekComment
	ppsSeekEndOfLine
	ppsSeekEndOfMultilineComment
	ppsSeekPackageName
	ppsDone
	pps_above_
)

type parseSecondaryState uint8

const (
	pss_below_ parseSecondaryState = iota
	pssNone
	pssSeekAsterisk
	pssSeekSlash
	pssSeekNonSpace
	pssSeekSpace
	pssSeekEnd
	pss_above_
)

type parseState struct {
	primary   parsePrimaryState
	secondary parseSecondaryState
}

//go:generate stringer -type=parsePrimaryState,parseSecondaryState

func GetPackageName(fileSpec string) (packageName string, err error) {
	var p *RuneProvider
	p, err = InitRuneProvider(fileSpec)
	if err != nil {
		err = bwerror.ErrorErr(err)
	} else {
		defer p.Close()
		var word string
		var word_line, word_col, word_pos int
		state := parseState{ppsSeekPackage, pssSeekNonSpace}
		for {
			var currRune rune
			var isEof bool
			currRune, isEof, err = p.PullRune()
			if err == nil {
				isUnexpectedRune := false
				switch state {
				case parseState{ppsSeekPackage, pssSeekNonSpace}:
					if !unicode.IsSpace(currRune) {
						switch currRune {
						case '/':
							state = parseState{ppsSeekComment, pssNone}
						case 'p':
							state = parseState{ppsSeekPackage, pssSeekSpace}
							word = string(currRune)
							word_line, word_col, word_pos = p.line, p.col, p.pos
						default:
							isUnexpectedRune = true
						}
					}
				case parseState{ppsSeekPackage, pssSeekSpace}:
					if !unicode.IsSpace(currRune) {
						word += string(currRune)
					} else if word == "package" {
						state = parseState{ppsSeekPackageName, pssSeekNonSpace}
					} else {
						err = bwerror.Error(
							"unexpected word <ansiPrimaryLiteral>%s<ansi> at line <ansiCmd>%d<ansi>, col <ansiCmd>%d<ansi> (pos <ansiCmd>%d<ansi>) in file <ansiCmd>%s",
							word, word_line, word_col, word_pos, fileSpec,
						)
					}
				case parseState{ppsSeekPackageName, pssSeekNonSpace}:
					if !unicode.IsSpace(currRune) {
						switch {
						case currRune == '/':
							state = parseState{ppsSeekComment, pssSeekAsterisk}
						case unicode.IsLetter(currRune):
							state = parseState{ppsSeekPackageName, pssSeekEnd}
							word = string(currRune)
						default:
							isUnexpectedRune = true
						}
					}
				case parseState{ppsSeekPackageName, pssSeekEnd}:
					if unicode.IsLetter(currRune) || currRune == '_' || unicode.IsDigit(currRune) {
						word += string(currRune)
					} else {
						packageName = word
						state = parseState{ppsDone, pssNone}
					}
				case parseState{ppsSeekComment, pssSeekAsterisk}:
					if currRune == '*' {
						state = parseState{ppsSeekEndOfMultilineComment, pssSeekAsterisk}
					} else {
						isUnexpectedRune = true
					}
				case parseState{ppsSeekComment, pssNone}:
					switch currRune {
					case '/':
						state = parseState{ppsSeekEndOfLine, pssNone}
					case '*':
						state = parseState{ppsSeekEndOfMultilineComment, pssSeekAsterisk}
					default:
						isUnexpectedRune = true
					}
				case parseState{ppsSeekEndOfLine, pssNone}:
					if currRune == '\n' {
						state = parseState{ppsSeekPackage, pssSeekNonSpace}
					}
				case parseState{ppsSeekEndOfMultilineComment, pssSeekAsterisk}:
					if currRune == '*' {
						state = parseState{ppsSeekEndOfMultilineComment, pssSeekSlash}
					}
				case parseState{ppsSeekEndOfMultilineComment, pssSeekSlash}:
					if currRune == '/' {
						if word == "package" {
							state = parseState{ppsSeekPackageName, pssSeekNonSpace}
						} else {
							state = parseState{ppsSeekPackage, pssSeekNonSpace}
						}
					}
					// default:
					// 	bwerror.Panic("no handler for %s.%s", state.primary, state.secondary)
				}
				if isUnexpectedRune {
					if isEof {
						err = bwerror.Error(
							"unexpected end of file at line <ansiCmd>%d<ansi>, col <ansiCmd>%d<ansi> (pos <ansiCmd>%d<ansi>) in file <ansiCmd>%s",
							p.line, p.col, p.pos, fileSpec,
						)
					} else {
						err = bwerror.Error(
							"unexpected rune <ansiPrimaryLiteral>%q<ansi> at line <ansiCmd>%d<ansi>, col <ansiCmd>%d<ansi> (pos <ansiCmd>%d<ansi>) in file <ansiCmd>%s",
							currRune, p.line, p.col, p.pos, fileSpec,
						)
					}
				}
			}
			if isEof || err != nil || (state == parseState{ppsDone, pssNone}) {
				break
			}
		}
	}
	return
}
