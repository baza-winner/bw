package path

import (
	"encoding/json"
	"unicode"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwstr"
	"github.com/baza-winner/bwcore/runeprovider"
)

//go:generate stringer -type=varPathParseState

type varPathParseState uint8

const (
	vppsBegin varPathParseState = iota
	vppsEnd
	vppsDone
	vppsIdx
	vppsDigit
	vppsKey
	vppsVar
	vppsString
	vppsForceEnd
)

func (v varPathParseState) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func Parse(s string) (result bw.ValPath, err error) {
	p := runeprovider.ProxyFrom(runeprovider.FromString(s))
	stack := []bw.ValPath{bw.ValPath{}}
	state := vppsBegin
	var item string
	for {
		p.PullRune()
		var currRune rune
		var isEOF bool
		currRune, isEOF, err = p.Rune()
		if err != nil {
			return
		}
		isUnexpectedRune := false
		switch {
		case state == vppsBegin:
			if isEOF {
				if len(stack) == 1 && len(stack[0]) == 0 {
					state = vppsDone
				} else {
					isUnexpectedRune = true
				}
			} else if unicode.IsDigit(currRune) &&
				len(stack) > 0 &&
				len(stack[len(stack)-1]) > 0 {
				item = string(currRune)
				state = vppsIdx
			} else if (currRune == '-' || currRune == '+') &&
				len(stack) > 0 &&
				len(stack[len(stack)-1]) > 0 {
				item = string(currRune)
				state = vppsDigit
			} else if unicode.IsLetter(currRune) || currRune == '_' {
				item = string(currRune)
				state = vppsKey
			} else if currRune == '{' {
				stack = append(stack, bw.ValPath{})
				state = vppsBegin
			} else if currRune == '$' {
				state = vppsVar
			} else if currRune == '#' {
				stack[len(stack)-1] = append(
					stack[len(stack)-1],
					bw.ValPathItem{Type: bw.ValPathItemHash},
				)
				state = vppsForceEnd
			} else {
				isUnexpectedRune = true
			}
		case state == vppsDigit:
			if unicode.IsDigit(currRune) {
				item += string(currRune)
				state = vppsIdx
			} else {
				isUnexpectedRune = true
			}
		case state == vppsForceEnd:
			if isEOF && len(stack) == 1 {
				state = vppsDone
			} else {
				isUnexpectedRune = true
			}
		case state == vppsEnd:
			if isEOF {
				if len(stack) == 1 {
					state = vppsDone
				} else {
					isUnexpectedRune = true
				}
			} else if currRune == '.' {
				state = vppsBegin
			} else if currRune == '}' && len(stack) > 0 {
				stack[len(stack)-2] = append(
					stack[len(stack)-2],
					bw.ValPathItem{Type: bw.ValPathItemPath, Path: stack[len(stack)-1]},
				)
				stack = stack[0 : len(stack)-1]
			} else {
				isUnexpectedRune = true
			}
		case state == vppsIdx:
			if unicode.IsDigit(currRune) {
				item += string(currRune)
			} else {
				if i, err := bwstr.ParseInt(item); err == nil {
					stack[len(stack)-1] = append(
						stack[len(stack)-1],
						bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: i},
					)
				}
				p.PushRune()
				state = vppsEnd
			}
		case state == vppsKey || state == vppsVar:
			if unicode.IsLetter(currRune) || currRune == '_' || unicode.IsDigit(currRune) {
				item += string(currRune)
			} else {
				tp := bw.ValPathItemKey
				if state == vppsVar {
					tp = bw.ValPathItemVar
				}
				stack[len(stack)-1] = append(
					stack[len(stack)-1],
					bw.ValPathItem{Type: tp, Key: item},
				)
				p.PushRune()
				state = vppsEnd
			}
		default:
			bwerr.Panic("no handler for %s", state)
		}
		if isUnexpectedRune {
			err = p.Unexpected(p.Curr)
			return
		}
		if isEOF || state == vppsDone {
			break
		}
	}
	result = stack[0]
	return
}
