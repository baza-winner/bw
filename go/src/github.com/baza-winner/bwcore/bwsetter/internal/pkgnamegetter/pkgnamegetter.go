package pkgnamegetter

import (
	"unicode"

	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwrune"
)

type parsePrimaryState uint8

const (
	ppsBelow parsePrimaryState = iota
	ppsSeekPackage
	ppsSeekComment
	ppsSeekEndOfLine
	ppsSeekEndOfMultilineComment
	ppsSeekPackageName
	ppsDone
	ppsAbove
)

type parseSecondaryState uint8

const (
	pssBelow parseSecondaryState = iota
	pssNone
	pssSeekAsterisk
	pssSeekSlash
	pssSeekNonSpace
	pssSeekSpace
	pssSeekEnd
	pssAbove
)

type parseState struct {
	primary   parsePrimaryState
	secondary parseSecondaryState
}

//go:generate stringer -type=parsePrimaryState,parseSecondaryState

func GetPackageName(fileSpec string) (packageName string, err error) {
	var p bwrune.Provider
	p, err = bwrune.FromFile(fileSpec)
	if err != nil {
		err = bwerr.FromA(bwerr.E{Error: err})
	} else {
		defer p.Close()
		var word string
		var wordLine, wordCol, wordPos int
		state := parseState{ppsSeekPackage, pssSeekNonSpace}
		for {
			var currRunePtr *rune
			var currRune rune
			var isEOF bool
			currRunePtr, err = p.PullRune()
			if currRunePtr == nil {
				isEOF = true
			} else {
				currRune = *currRunePtr
			}
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
							wordLine, wordCol, wordPos = p.Line(), p.Col(), p.Pos()
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
						err = bwerr.From(
							"unexpected word <ansiVal>%s<ansi> at line <ansiPath>%d<ansi>, col <ansiPath>%d<ansi> (pos <ansiPath>%d<ansi>) in file <ansiPath>%s",
							word, wordLine, wordCol, wordPos, fileSpec,
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
					// 	bwerr.Panic("no handler for %s.%s", state.primary, state.secondary)
				}
				if isUnexpectedRune {
					if isEOF {
						err = bwerr.From(
							"unexpected end of file at line <ansiPath>%d<ansi>, col <ansiPath>%d<ansi> (pos <ansiPath>%d<ansi>) in file <ansiPath>%s",
							p.Line(), p.Col(), p.Pos(), fileSpec,
						)
					} else {
						err = bwerr.From(
							"unexpected rune <ansiVal>%q<ansi> at line <ansiPath>%d<ansi>, col <ansiPath>%d<ansi> (pos <ansiPath>%d<ansi>) in file <ansiPath>%s",
							currRune, p.Line(), p.Col(), p.Pos(), fileSpec,
						)
					}
				}
			}
			if isEOF || err != nil || (state == parseState{ppsDone, pssNone}) {
				break
			}
		}
	}
	return
}
