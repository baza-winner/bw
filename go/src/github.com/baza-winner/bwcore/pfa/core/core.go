package core

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwfmt"
	"github.com/baza-winner/bwcore/bwint"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/pfa/formatted"
	"github.com/baza-winner/bwcore/runeprovider"
	"github.com/jimlawless/whereami"
)

// ============================================================================

// type PfaErrorProvider interface {
// 	PfaError(pfa *PfaStruct) error
// }

// ============================================================================

type ProcessorAction interface {
	Execute(pfa *PfaStruct)
}

// ============================================================================

type ProccessorActionProvider interface {
	GetAction() ProcessorAction
}

// ============================================================================

type ValProvider interface {
	GetVal(pfa *PfaStruct) interface{}
	GetSource(pfa *PfaStruct) formatted.String
}

// ============================================================================

type ValChecker interface {
	Conforms(pfa *PfaStruct, val interface{}, varPath VarPath) bool
}

// ============================================================================

type ValCheckerProvider interface {
	GetChecker() ValChecker
}

// ============================================================================

type TraceLevel uint8

const (
	TraceNone TraceLevel = iota
	TraceBrief
	TraceAll
)

// ============================================================================

type VarValue struct {
	Val interface{}
	pfa *PfaStruct
}

func VarValueFrom(val interface{}) VarValue {
	return VarValue{val, nil}
}

func (v VarValue) GetVal(varPath VarPath) (result VarValue) {
	// fmt.Printf("GetVal: %s\n", varPath.formatted.String(nil))
	if v.pfa.Err != nil || len(varPath) == 0 {
		result = v
	} else {
		result = VarValue{nil, v.pfa}
		v.helper(varPath, nil,
			func(vt valType, vals []interface{}, m map[string]interface{}, idx int, key string, VarVal interface{}) {
				switch vt {
				case valTypeSlice:
					minIdx := -len(vals)
					maxIdx := len(vals) - 1
					if minIdx <= idx && idx <= maxIdx {
						result.Val = vals[idx]
					}
				case valTypeMap:
					result.Val = m[key]
				}
			},
			func(vt valType, vals []interface{}, m map[string]interface{}) {
				switch vt {
				case valTypeSlice:
					result.Val = len(vals)
				case valTypeMap:
					result.Val = len(m)
				default:
					result.Val = 0
				}
			},
		)
		if result.pfa.Err == nil && len(varPath) > 1 {
			result = result.GetVal(varPath[1:])
		}
	}

	return
}

type valType uint8

const (
	valTypeNil valType = iota
	valTypeSlice
	valTypeMap
)

//go:generate stringer -type=valType

func (v VarValue) helper(
	varPath VarPath,
	VarVal interface{},
	onVal func(vt valType, vals []interface{}, m map[string]interface{}, idx int, key string, VarVal interface{}),
	onLen func(vt valType, vals []interface{}, m map[string]interface{}),
) {
	// varPathItem := varPath[0]
	// fmt.Printf("helper: %s,isIdx: %s, idx: %s, key: %s, err: %s \n", varPath.formatted.String(nil), isIdx, idx, key, err)
	// if err != nil {
	// 	v.pfa.Err = err
	// } else
	if varPath[0].Type == VarPathItemHash {
		if v.Val == nil {
			onLen(valTypeNil, nil, nil)
		} else {
			switch t := v.Val.(type) {
			case []interface{}:
				onLen(valTypeSlice, t, nil)
			case map[string]interface{}:
				onLen(valTypeSlice, nil, t)
			default:
				v.pfa.SetError("%s nor <ansiOutline>Array, neither <ansiOutline>Map", varPath.FormattedString())
				// v.pfa.Err = PfaError{formatted.StringFrom("%s nor <ansiOutline>Array, neither <ansiOutline>Map", varPath.FormattedString())}
			}
		}
	} else if v.Val == nil {
		onVal(valTypeNil, nil, nil, 0, "", VarVal)
	} else {
		vt, idx, key, err := varPath[0].TypeIdxKey(v.pfa)
		if err != nil {
			v.pfa.Err = err
		} else if vt == VarPathItemIdx {
			if vals, ok := v.Val.([]interface{}); !ok {
				// v.pfa.Err = PfaError{formatted.StringFrom("%s is not <ansiOutline>Array", varPath.FormattedString())}
				v.pfa.SetError("%s is not <ansiOutline>Array", varPath.FormattedString())
			} else {
				onVal(valTypeSlice, vals, nil, idx, "", VarVal)
			}
		} else if m, ok := v.Val.(map[string]interface{}); !ok {
			v.pfa.SetError("%s is not <ansiOutline>Map", varPath.FormattedString())
		} else {
			onVal(valTypeSlice, nil, m, 0, key, VarVal)
		}
	}
}

type PfaError struct {
	pfa     *PfaStruct
	content *PfaErrorContent
	Where   string
}

type PfaErrorContentState uint8

const (
	PecsNeedPrepare PfaErrorContentState = iota
	PecsPrepared
)

type PfaErrorContent struct {
	state     PfaErrorContentState
	reason    string
	fmtString string
	errStr    string
}

func (v PfaError) DataForJSON() interface{} {
	result := map[string]interface{}{}
	switch v.content.state {
	case PecsNeedPrepare:
		result["reason"] = v.content.reason
		if len(v.content.fmtString) > 0 {
			result["fmtString"] = v.content.fmtString
		}
	case PecsPrepared:
		result["pfa"] = v.pfa.DataForJSON()
		result["err"] = v.content.errStr
		result["Where"] = v.Where
	}
	return result
}

func (pfa *PfaStruct) SetError(fmtString string, fmtArgs ...interface{}) {
	pfa.Err = PfaError{
		pfa: pfa,
		content: &PfaErrorContent{
			reason: string(formatted.StringFrom(fmtString, fmtArgs)),
		},
		Where: whereami.WhereAmI(3),
	}
}

func (pfa *PfaStruct) SetTransformError(fmtString, reason string) {
	pfa.Err = PfaError{
		pfa:     pfa,
		content: &PfaErrorContent{fmtString: fmtString, reason: reason},
		Where:   whereami.WhereAmI(3),
	}
	// pfa.Err = PfaError{
	// 	pfa: pfa,
	// 	content: &PfaErrorContent{
	// 		reason: string(formatted.StringFrom(fmtString, fmtArgs)),
	// 	},
	// 	Where: whereami.WhereAmI(3),
	// }
}

func (pfa *PfaStruct) SetUnexpectedError(err error) {
	pfa.Err = PfaError{
		pfa:     pfa,
		content: &PfaErrorContent{errStr: err.Error(), state: PecsPrepared},
		Where:   whereami.WhereAmI(3),
	}
}

// func PfaErrorFrom(pfa *PfaStruct, fmtString string, fmtArgs ...interface{}) PfaError {
// 	return PfaError{
// 		pfa: pfa,
// 		content: &PfaErrorContent{
// 			reason: string(formatted.StringFrom(fmtString, fmtArgs)),
// 		},
// 		Where: whereami.WhereAmI(3),
// 	}
// }

// func TransformErrorFrom(pfa *PfaStruct, fmtString, reason string) PfaError {
// 	return PfaError{
// 		pfa:     pfa,
// 		content: &PfaErrorContent{fmtString: fmtString, reason: reason},
// 		Where:   whereami.WhereAmI(3),
// 	}
// }

// func UnexpetedErrorFrom(pfa *PfaStruct, err error) PfaError {
// 	return PfaError{
// 		pfa:     pfa,
// 		content: &PfaErrorContent{errStr: err.Error(), state: PecsPrepared},
// 		Where:   whereami.WhereAmI(3),
// 	}
// }

func (v *PfaError) PrepareErr(fmtString string, fmtArgs ...interface{}) {
	if v.content.state == PecsPrepared {
		bwerror.Panic("Already prepared %s ", bwjson.PrettyJsonOf(v))
	} else {
		v.content.fmtString = fmtString + ": " + v.content.reason
		// v.content.fmtArgs = fmtArgs
		v.content.errStr = bwerror.Error(v.content.fmtString, fmtArgs...).Error()
		v.content.state = PecsPrepared
	}
	// v.content.err = bwerror.Error(fmtString+": "+string(v.content.reason), fmtArgs...)
}

func (v *PfaError) SetErr(errStr string) {
	if v.content.state == PecsPrepared {
		bwerror.Panic("Already prepared %s ", bwjson.PrettyJsonOf(v))
	} else {
		v.content.errStr = errStr
		v.content.state = PecsPrepared
	}
}

func (v PfaError) Error() (result string) {
	switch v.content.state {
	case PecsNeedPrepare:
		bwerror.Panic("NeedPrepare %s ", bwjson.PrettyJsonOf(v))
	case PecsPrepared:
		result = v.content.errStr
		// result = bwerror.Error(v.content.FmtString, v.content.fmtArgs).Error()
	}
	return
}

func (v PfaError) State() PfaErrorContentState {
	return v.content.state
}

// func (v PfaError) Error(pfa *PfaStruct) error {
// 	bwerror.Unreachable()
// 	return nil
// }

func (v VarValue) SetVal(varPath VarPath, VarVal interface{}) {
	if len(varPath) == 0 {
		v.pfa.Panic(bwfmt.StructFrom("varPath: %#v", varPath))
	} else {
		target := VarValue{nil, v.pfa}
		v.helper(varPath, VarVal,
			func(vt valType, vals []interface{}, m map[string]interface{}, idx int, key string, VarVal interface{}) {
				switch vt {
				case valTypeSlice:
					if len(vals) == 0 {
						// v.pfa.Err = PfaError{formatted.StringFrom("path does not exist (no elem with idx <ansiPrimary>%d<ansi> at empty Array)", idx)}
						v.pfa.SetError("path does not exist (no elem with idx <ansiPrimary>%d<ansi> at empty Array)", idx)
					} else {
						minIdx := -len(vals)
						maxIdx := len(vals) - 1
						if !(minIdx <= idx && idx <= maxIdx) {
							// v.pfa.Err = PfaError{formatted.StringFrom("path does not exist (<ansiPrimary>%d<ansi> is out of range <ansiSecondary>[%d, %d]<ansi>)", idx, minIdx, maxIdx)}
							v.pfa.SetError("path does not exist (<ansiPrimary>%d<ansi> is out of range <ansiSecondary>[%d, %d]<ansi>)", idx, minIdx, maxIdx)
						} else {
							if idx < 0 {
								idx = len(vals) + idx
							}
							if len(varPath) == 1 {
								vals[idx] = VarVal
							} else {
								target.Val = vals[idx]
							}
						}
					}
				case valTypeMap:
					if len(varPath) == 1 {
						m[key] = VarVal
					} else if kv, ok := m[key]; !ok {
						// v.pfa.Err = PfaError{formatted.StringFrom("path does not exist (no key <ansiPrimary>%s)", key)}
						v.pfa.SetError("path does not exist (no key <ansiPrimary>%s)", key)
					} else {
						target.Val = kv
					}
				case valTypeNil:
					// v.pfa.Err = PfaError{formatted.StringFrom("can not set to nil value")}
					v.pfa.SetError("can not set to nil value")
				}
			},
			func(vt valType, vals []interface{}, m map[string]interface{}) {
				// v.pfa.Err = PfaError{formatted.StringFrom("<ansiOutline>hash path<ansi> is readonly")}
				v.pfa.SetError("<ansiOutline>hash path<ansi> is readonly")
			},
		)
		if target.pfa.Err == nil && len(varPath) > 1 {
			target.SetVal(varPath[1:], VarVal)
		}
	}
}

func (v VarValue) Rune() (result rune, err error) {
	if v.pfa != nil && v.pfa.Err != nil {
		err = v.pfa.Err
	} else {
		var ok bool
		if result, ok = v.Val.(rune); !ok {
			err = bwerror.Error("%#v is not rune", v.Val)
		}
	}
	return
}

func (v VarValue) Int() (result int, err error) {
	if v.pfa != nil && v.pfa.Err != nil {
		err = v.pfa.Err
	} else {
		vValue := reflect.ValueOf(v.Val)
		switch vValue.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			_int64 := vValue.Int()
			if int64(bwint.MinInt) <= _int64 && _int64 <= int64(bwint.MaxInt) {
				result = int(_int64)
			} else {
				err = bwerror.Error("%d is out of range [%d, %d]", _int64, bwint.MinInt, bwint.MaxInt)
			}
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			_uint64 := vValue.Uint()
			if _uint64 <= uint64(bwint.MaxInt) {
				result = int(_uint64)
			} else {
				err = bwerror.Error("%d is more than %d", _uint64, bwint.MaxInt)
			}
		default:
			err = bwerror.Error("<ansiPrimary>%#v<ansi> is not of type <ansiSecondary>int", v)
		}
	}
	return
}

func (v VarValue) String() (result string, err error) {
	if v.pfa != nil && v.pfa.Err != nil {
		err = v.pfa.Err
	} else {
		var ok bool
		if result, ok = v.Val.(string); !ok {
			err = bwerror.Error("<ansiPrimary>%#v<ansi> is not of type <ansiSecondary>string", v)
		}
	}
	return
}

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
				} else if unicode.IsDigit(currRune) && len(stack) > 0 {
					item = string(currRune)
					state = vppsIdx
				} else if (currRune == '-' || currRune == '+') && len(stack) > 0 {
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

type PfaStruct struct {
	Stack ParseStack
	Proxy *runeprovider.Proxy
	// Err          ErrorProvider
	Err             error
	Vars            map[string]interface{}
	TraceLevel      TraceLevel
	traceConditions []string
	ruleLevel       int
}

func PfaFrom(p runeprovider.RuneProvider, TraceLevel TraceLevel) *PfaStruct {
	return &PfaStruct{
		Stack:      ParseStack{},
		Proxy:      runeprovider.ProxyFrom(p),
		Vars:       map[string]interface{}{},
		TraceLevel: TraceLevel,
	}
}

// func (pfa *PfaStruct) Value(val interface{}) VarValue {
// 	return VarValue{val, pfa}
// }

func (pfa *PfaStruct) indent(indentLevel int) string {
	indentAtom := "  "
	indent := ""
	for i := 0; i <= indentLevel; i++ {
		indent += indentAtom
	}
	return indent
}

func (pfa *PfaStruct) fmtArgs(fmtArgs ...interface{}) []interface{} {
	result := []interface{}{}
	for _, arg := range fmtArgs {
		if f, ok := arg.(func(pfa *PfaStruct) interface{}); ok {
			arg = f(pfa)
		}
		result = append(result, pfa.TraceVal(arg))
	}
	return result
}

// type formatted.String string

func (pfa *PfaStruct) Panic(optFmtStruct ...bwfmt.Struct) {
	fmtString := "<ansiOutline>pfa<ansi> <ansiSecondary>%s<ansi>"
	fmtArgs := []interface{}{pfa}
	if optFmtStruct == nil {
		bwerror.Panicd(1, fmtString, fmtArgs...)
	} else {
		fmtString += " " + optFmtStruct[0].FmtString
		fmtArgs = append(fmtArgs, optFmtStruct[0].FmtArgs...)
	}
}

func (pfa *PfaStruct) PanicErr(err error) {
	fmtString := "<ansiOutline>pfa<ansi> <ansiSecondary>%s<ansi>"
	fmtArgs := []interface{}{pfa}
	bwerror.PanicErr(fmt.Errorf(err.Error()+"\n"+ansi.Ansi("", fmtString), fmtArgs), 1)
}

func (pfa *PfaStruct) VarValue(varPath VarPath) (result VarValue) {
	// fmt.Printf("VarValue: %s\n", varPath.formatted.String(nil))
	result = VarValue{nil, pfa}
	if len(varPath) > 0 {
		pfa.getSetHelper(varPath, nil,
			func(name string, ofs int) {
				switch name {
				case "rune":
					currRune, _ := pfa.Proxy.Rune(ofs)
					result.Val = currRune
				case "runePos":
					ps := pfa.Proxy.PosStruct(ofs)
					result.Val = ps
				}
			},
			func(pfaVars VarValue, VarVal interface{}) {
				result = pfaVars.GetVal(varPath)
				return
			},
		)
		if pfa.Err != nil {
			switch t := pfa.Err.(type) {
			case PfaError:
				// pfa.Err = nil
				if t.State() == PecsNeedPrepare {
					t.PrepareErr("failed to get %s", varPath.FormattedString())
				}
				// pfa.Err.Prepare
				// pfa.Err = pfa.Error(bwerror.Error("failed to get %s: "+string(t.s), varPath.FormattedString()))
			}
		}
	}
	return
}

// func (pfa *PfaStruct) Error(err error) error {
// 	return pfaError{
// 		pfa,
// 		err,
// 		whereami.WhereAmI(3),
// 	}
// }

func (pfa *PfaStruct) getSetHelper(
	varPath VarPath,
	VarVal interface{},
	onSpecial func(name string, ofs int),
	onPfaVar func(pfaVars VarValue, VarVal interface{}),
) {
	if len(varPath) == 0 {
		pfa.Err = bwerror.Error("varPath is empty")
	} else {
		vt, _, key, err := varPath[0].TypeIdxKey(pfa)
		if err != nil {
			pfa.Err = err
		} else if vt != VarPathItemKey {
			pfa.SetError("path must start with key")
		} else if key == "rune" || key == "runePos" {
			var ofs int
			if len(varPath) > 2 {
				// pfa.Err = PfaError{formatted.StringFrom("<ansiPrimary>%s<ansi> path may have at most 2 items", key)}
				pfa.SetError("<ansiPrimary>%s<ansi> path may have at most 2 items", key)
			} else if len(varPath) > 1 {
				vt, idx, key, err := varPath[1].TypeIdxKey(pfa)
				if err != nil {
					pfa.Err = err
				} else if vt != VarPathItemIdx {
					// pfa.Err = PfaError{formatted.StringFrom("<ansiPrimary>%s<ansi> path expects <ansiOutline>idx<ansi> as second item", key)}
					pfa.SetError("<ansiPrimary>%s<ansi> path expects <ansiOutline>idx<ansi> as second item", key)
				} else {
					ofs = idx
				}
			}
			if pfa.Err == nil {
				onSpecial(key, ofs)
			}
		} else {
			onPfaVar(VarValue{pfa.Vars, pfa}, VarVal)
		}
	}
}

func (pfa *PfaStruct) SetVarVal(varPath VarPath, VarVal interface{}) {
	if len(varPath) == 0 {
		pfa.SetError("varPath is empty")
	} else {
		pfa.getSetHelper(varPath, VarVal,
			func(name string, idx int) {
				pfa.SetError("<ansiOutline>%s<ansi> is read only", name)
				// pfa.Err = bwerror.Error("<ansiOutline>%s<ansi> is read only", name)
			},
			func(pfaVars VarValue, VarVal interface{}) {
				pfaVars.SetVal(varPath, VarVal)
			},
		)
	}
	if pfa.Err != nil {
		switch t := pfa.Err.(type) {
		case PfaError:
			// pfa.Err = nil
			if t.State() == PecsNeedPrepare {
				t.PrepareErr("failed to set %s", varPath.FormattedString())
			}
			// pfa.Err = pfa.Error(bwerror.Error("failed to set %s: "+string(t.s), varPath.FormattedString(nil)))
		}
	}
}

func (pfa *PfaStruct) ifStackLen(minLen int) bool {
	return len(pfa.Stack) >= minLen
}

func (pfa *PfaStruct) mustStackLen(minLen int) {
	if !pfa.ifStackLen(minLen) {
		pfa.Panic(bwfmt.StructFrom("<ansiOutline>minLen <ansiSecondary>%d", minLen))
	}
}

func (pfa *PfaStruct) GetTopStackItem(optDeep ...uint) *ParseStackItem {
	ofs := -1
	if optDeep != nil {
		ofs = ofs - int(optDeep[0])
	}
	pfa.mustStackLen(-ofs)
	return &pfa.Stack[len(pfa.Stack)+ofs]
}

func (pfa *PfaStruct) PopStackItem() {
	pfa.mustStackLen(1)
	pfa.Stack = pfa.Stack[:len(pfa.Stack)-1]
}

func (pfa *PfaStruct) PushStackItem() {
	pfa.Stack = append(pfa.Stack, ParseStackItem{
		Start: pfa.Proxy.Curr,
		Vars:  map[string]interface{}{},
	})
}

func (pfa PfaStruct) DataForJSON() interface{} {
	result := map[string]interface{}{}
	result["Stack"] = pfa.Stack.DataForJSON()
	result["Proxy"] = pfa.Proxy.DataForJSON()
	if len(pfa.Vars) > 0 {
		result["Vars"] = pfa.Vars
	}
	return result
}

func (pfa PfaStruct) String() string {
	return bwjson.PrettyJsonOf(pfa)
}

// ============================================================================

// type pfaError struct {
// 	pfa *PfaStruct
// 	// Err interface{}
// 	err   error
// 	Where string
// }

// func (err pfaError) Error() string {
// 	return err.err.Error()
// }

// func (v pfaError) DataForJSON() interface{} {
// 	result := map[string]interface{}{}
// 	result["pfa"] = v.pfa.DataForJSON()
// 	// result["Err"] = v.Err
// 	result["err"] = v.err
// 	result["Where"] = v.Where
// 	return result
// }

// ============================================================================

var underscoreRegexp = regexp.MustCompile("[_]+")

func ParseNumber(source string) (value interface{}, err error) {
	source = underscoreRegexp.ReplaceAllLiteralString(source, ``)
	if strings.Contains(source, `.`) {
		var _float64 float64
		if _float64, err = strconv.ParseFloat(source, 64); err == nil {
			value = _float64
		}
	} else {
		var _int64 int64
		if _int64, err = strconv.ParseInt(source, 10, 64); err == nil {
			if int64(bwint.MinInt8) <= _int64 && _int64 <= int64(bwint.MaxInt8) {
				value = int8(_int64)
			} else if int64(bwint.MinInt16) <= _int64 && _int64 <= int64(bwint.MaxInt16) {
				value = int16(_int64)
			} else if int64(bwint.MinInt32) <= _int64 && _int64 <= int64(bwint.MaxInt32) {
				value = int32(_int64)
			} else {
				value = _int64
			}
		}
	}
	return
}

func ParseInt(source string) (value int, err error) {
	source = underscoreRegexp.ReplaceAllLiteralString(source, ``)
	var _int64 int64
	if _int64, err = strconv.ParseInt(source, 10, 64); err == nil {
		if int64(bwint.MinInt) <= _int64 && _int64 <= int64(bwint.MaxInt) {
			value = int(_int64)
		} else {
			err = bwerror.Error("<ansiPrimary>%d<ansi> is out of range <ansiSecondary>[%d, %d]", _int64, bwint.MinInt, bwint.MaxInt)
		}
	}
	return
}

// ============================================================================

type ParseStack []ParseStackItem

func (Stack *ParseStack) DataForJSON() interface{} {
	result := []interface{}{}
	for _, item := range *Stack {
		result = append(result, item.DataForJSON())
	}
	return result
}

func (Stack *ParseStack) String() (result string) {
	return bwjson.PrettyJsonOf(Stack)
}

// ============================================================================

type ParseStackItem struct {
	Start runeprovider.PosStruct
	Vars  map[string]interface{}
}

func (stackItem *ParseStackItem) DataForJSON() interface{} {
	result := map[string]interface{}{}
	result["Start"] = stackItem.Start.DataForJSON()
	result["Vars"] = stackItem.Vars
	return result
}

func (stackItem *ParseStackItem) String() (result string) {
	return bwjson.PrettyJsonOf(stackItem)
}
