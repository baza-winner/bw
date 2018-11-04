package bwval

import (
	"encoding/json"
	"unicode"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwstr"
	"github.com/baza-winner/bwcore/runeprovider"
	// "github.com/baza-winner/bwcore/formatted"
)

// type VarPath interface {
// FormattedString() (result formatted.String)
// }

type varPathItemType uint8

const (
	varPathItemHash varPathItemType = iota
	varPathItemIdx
	varPathItemKey
	varPathItemPath
)

func (v varPathItemType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

type varPathItem struct {
	Type varPathItemType
	Idx  int
	Key  string
	Path VarPath
}

type VarPath struct {
	items []varPathItem
}

type varPathParseState uint8

const (
	vppsBegin varPathParseState = iota
	vppsEnd
	vppsDone
	vppsIdx
	vppsDigit
	vppsKey
	vppsForceEnd
)

//go:generate stringer -type=varPathParseState,varPathItemType

func (v varPathParseState) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func VarPathFrom(s string) (result VarPath, err error) {
	p := runeprovider.ProxyFrom(runeprovider.FromString(s))
	stack := []VarPath{VarPath{}}
	state := vppsBegin
	var item string
	for {
		p.PullRune()
		currRune, isEOF := p.Rune()
		if err == nil {
			isUnexpectedRune := false
			switch {
			case state == vppsBegin:
				if isEOF {
					if len(stack) == 1 && len(stack[0].items) == 0 {
						state = vppsDone
					} else {
						isUnexpectedRune = true
					}
				} else if unicode.IsDigit(currRune) && len(stack[len(stack)-1].items) > 0 {
					item = string(currRune)
					state = vppsIdx
				} else if (currRune == '-' || currRune == '+') && len(stack[len(stack)-1].items) > 0 {
					item = string(currRune)
					state = vppsDigit
				} else if unicode.IsLetter(currRune) || currRune == '_' {
					item = string(currRune)
					state = vppsKey
				} else if currRune == '{' {
					stack = append(stack, VarPath{})
					state = vppsBegin
				} else if currRune == '#' {
					stack[len(stack)-1].items = append(stack[len(stack)-1].items, varPathItem{Type: varPathItemHash})
					// stack = append(stack, VarPath{})
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
					stack[len(stack)-2].items = append(stack[len(stack)-2].items, varPathItem{Type: varPathItemPath, Path: stack[len(stack)-1]})
					stack = stack[0 : len(stack)-1]
				} else {
					isUnexpectedRune = true
				}
			case state == vppsIdx:
				if unicode.IsDigit(currRune) {
					item += string(currRune)
				} else {
					if i, err := bwstr.ParseInt(item); err == nil {
						stack[len(stack)-1].items = append(stack[len(stack)-1].items, varPathItem{Type: varPathItemIdx, Idx: i})
					}
					p.PushRune()
					state = vppsEnd
				}
			case state == vppsKey:
				if unicode.IsLetter(currRune) || currRune == '_' || unicode.IsDigit(currRune) {
					item += string(currRune)
				} else {
					stack[len(stack)-1].items = append(stack[len(stack)-1].items, varPathItem{Type: varPathItemKey, Key: item})
					p.PushRune()
					state = vppsEnd
				}
			default:
				bwerr.Panic("no handler for %s", state)
			}
			if isUnexpectedRune {
				// err = p.UnexpectedRuneError(fmt.Sprintf("state = %s", state))
				err = p.Unexpected(p.Curr, bw.Fmt("state = %s", state))
				// err = p.UnexpectedRuneError(p.Curr, "state = %s", state)
			}
		}
		if isEOF || err != nil || (state == vppsDone) {
			break
		}
	}
	if err == nil {
		result = stack[0]
	}
	return
}

// func MustVarPath(varPath VarPath, err error) (result VarPath) {
// 	bwerr.TODO()
// 	return
// }

// type QualVarPath interface {
// 	// TypeIdxKey(i int) (itemType varPathItemType, idx int, key string, err error)
// 	// FormattedString(optVal ...Val) (result formatted.String)
// }

// type Opts struct {
// 	Vars   map[string]Val
// 	Consts map[string]Val
// }

// func QualVarPathFrom(varPath VarPath, opts ...Opts) (result QualVarPath, err error) {
// 	bwerr.TODO()
// 	return
// }

// type Val interface {
// 	PathVal(path QualVarPath) (result Val, err error)
// 	SetPathVal(path QualVarPath, val []interface{}) (err error)
// 	// Val
// }
