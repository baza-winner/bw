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

func MustParse(s string, optBase ...bw.ValPath) (result bw.ValPath) {
	var err error
	if result, err = Parse(s, optBase...); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return result
}

func Parse(s string, optBase ...bw.ValPath) (result bw.ValPath, err error) {
	base := bw.ValPath{}
	if len(optBase) > 0 {
		base = optBase[0]
	}
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
		switch {
		case state == vppsBegin:
			if isEOF {
				// if len(stack) == 1 && len(stack[0]) == 0 {
				// state = vppsDone
				// } else {
				err = p.Unexpected(p.Curr)
				return
				// }
			} else if currRune == '.' &&
				len(stack) > 0 &&
				len(stack[len(stack)-1]) == 0 {
				if len(base) == 0 {
					state = vppsForceEnd
				} else if len(stack) == 1 && len(stack[0]) == 0 {
					stack[0] = append(stack[0], base...)
					// state = vppsForceEnd
				} else {
					err = p.Unexpected(p.Curr)
					return
				}
			} else if unicode.IsDigit(currRune) {
				item = string(currRune)
				state = vppsIdx
			} else if currRune == '-' || currRune == '+' {
				item = string(currRune)
				state = vppsDigit
			} else if unicode.IsLetter(currRune) || currRune == '_' {
				item = string(currRune)
				state = vppsKey
			} else if currRune == '{' {
				stack = append(stack, bw.ValPath{})
				state = vppsBegin
			} else if currRune == '$' &&
				len(stack) > 0 &&
				len(stack[len(stack)-1]) == 0 {
				item = ""
				state = vppsVar
			} else if currRune == '#' {
				stack[len(stack)-1] = append(
					stack[len(stack)-1],
					bw.ValPathItem{Type: bw.ValPathItemHash},
				)
				state = vppsForceEnd
			} else {
				err = p.Unexpected(p.Curr)
				return
			}
		case state == vppsDigit:
			if unicode.IsDigit(currRune) {
				item += string(currRune)
				state = vppsIdx
			} else {
				err = p.Unexpected(p.Curr)
				return
			}
		case state == vppsForceEnd:
			if isEOF && len(stack) == 1 {
				state = vppsDone
			} else {
				err = p.Unexpected(p.Curr)
				return
			}
		case state == vppsEnd:
			if isEOF {
				if len(stack) == 1 {
					state = vppsDone
				} else {
					err = p.Unexpected(p.Curr)
					return
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
				err = p.Unexpected(p.Curr)
				return
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
		// bwdebug.Print("state", state, "currRune", currRune, "stack", stack, "item", item)
		if state == vppsDone {
			break
		}
	}
	result = stack[0]
	return
}
