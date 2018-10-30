package core

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwfmt"
	"github.com/baza-winner/bwcore/pfa/formatted"
	"github.com/baza-winner/bwcore/runeprovider"
)

// ============================================================================

type VarPathItem struct {
	Type VarPathItemType
	Idx  int
	Key  string
	Path VarPath
	// Val interface{}
}

type varPathHash struct{}

type VarPathItemType uint8

const (
	VarPathItemHash VarPathItemType = iota
	VarPathItemIdx
	VarPathItemKey
	VarPathItemPath
)

func (v VarPathItem) TypeIdxKey(pfa *PfaStruct) (itemType VarPathItemType, idx int, key string, err error) {
	itemType = v.Type
	switch v.Type {
	case VarPathItemHash:
	case VarPathItemIdx:
		idx = v.Idx
	case VarPathItemKey:
		key = v.Key
	case VarPathItemPath:
		if pfa == nil {
			err = bwerror.Error("VarPath requires pfa")
		} else if varValue := pfa.VarValue(v.Path); pfa.Err != nil {
			err = pfa.Err
		} else if idx, err = varValue.Int(); err == nil {
			itemType = VarPathItemIdx
		} else if key, err = varValue.String(); err == nil {
			itemType = VarPathItemKey
		} else {
			err = bwerror.Error("%s is nor Idx, neither Key", pfa.TraceVal(v.Path))
		}
	}
	return
}

// ============================================================================

type VarPath []VarPathItem

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

//go:generate stringer -type=varPathParseState

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
					if len(stack) == 1 && len(stack[0]) == 0 {
						state = vppsDone
					} else {
						isUnexpectedRune = true
					}
				} else if unicode.IsDigit(currRune) && len(stack[len(stack)-1]) > 0 {
					item = string(currRune)
					state = vppsIdx
				} else if (currRune == '-' || currRune == '+') && len(stack[len(stack)-1]) > 0 {
					item = string(currRune)
					state = vppsDigit
				} else if unicode.IsLetter(currRune) || currRune == '_' {
					item = string(currRune)
					state = vppsKey
				} else if currRune == '{' {
					stack = append(stack, VarPath{})
					state = vppsBegin
				} else if currRune == '#' {
					stack[len(stack)-1] = append(stack[len(stack)-1], VarPathItem{Type: VarPathItemHash})
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
					stack[len(stack)-2] = append(stack[len(stack)-2], VarPathItem{Type: VarPathItemPath, Path: stack[len(stack)-1]})
					stack = stack[0 : len(stack)-1]
				} else {
					isUnexpectedRune = true
				}
			case state == vppsIdx:
				if unicode.IsDigit(currRune) {
					item += string(currRune)
				} else {
					if i, err := ParseInt(item); err == nil {
						stack[len(stack)-1] = append(stack[len(stack)-1], VarPathItem{Type: VarPathItemIdx, Idx: i})
					}
					p.PushRune()
					state = vppsEnd
				}
			case state == vppsKey:
				if unicode.IsLetter(currRune) || currRune == '_' || unicode.IsDigit(currRune) {
					item += string(currRune)
				} else {
					stack[len(stack)-1] = append(stack[len(stack)-1], VarPathItem{Type: VarPathItemKey, Key: item})
					p.PushRune()
					state = vppsEnd
				}
			default:
				bwerror.Panic("no handler for %s", state)
			}
			if isUnexpectedRune {
				// err = p.UnexpectedRuneError(fmt.Sprintf("state = %s", state))
				err = p.Unexpected(p.Curr, bwfmt.StructFrom("state = %s", state))
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

func MustVarPathFrom(s string) (result VarPath) {
	var err error
	if result, err = VarPathFrom(s); err != nil {
		bwerror.PanicErr(err)
	}
	return
}

func (v VarPath) FormattedString(optPfa ...*PfaStruct) formatted.String {
	var pfa *PfaStruct
	if optPfa != nil {
		pfa = optPfa[0]
	}
	ss := []string{}
	for _, vpi := range v {
		switch vpi.Type {
		case VarPathItemPath:
			if pfa == nil {
				ss = append(ss, fmt.Sprintf("{%s}", vpi.Path.FormattedString(nil)))
			} else {
				ss = append(ss, fmt.Sprintf("{%s(%s)}", vpi.Path.FormattedString(pfa), pfa.TraceVal(vpi.Path)))
			}
		case VarPathItemKey:
			ss = append(ss, vpi.Key)
		case VarPathItemIdx:
			ss = append(ss, strconv.FormatInt(int64(vpi.Idx), 10))
		case VarPathItemHash:
			ss = append(ss, "#")
		}
	}
	return formatted.StringFrom("<ansiCmd>%s", strings.Join(ss, "."))
}

// ============================================================================
