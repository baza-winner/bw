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
	"github.com/baza-winner/bwcore/bwint"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/pfa/formatted"
	"github.com/baza-winner/bwcore/runeprovider"
	"github.com/jimlawless/whereami"
)

type ProcessorAction interface {
	Execute(pfa *PfaStruct)
}

// ============================================================================

type ValProvider interface {
	GetVal(pfa *PfaStruct) interface{}
	GetSource(pfa *PfaStruct) formatted.String
}

// ============================================================================

type TraceLevel uint8

const (
	TraceNone TraceLevel = iota
	TraceBrief
	TraceAll
)

type VarValue struct {
	val interface{}
	pfa *PfaStruct
}

// ============================================================================

func (v VarValue) GetVal(varPath VarPath) (result VarValue) {
	// fmt.Printf("GetVal: %s\n", varPath.formatted.String(nil))
	if v.pfa.Err != nil || len(varPath) == 0 {
		result = v
	} else {
		result = VarValue{nil, v.pfa}
		v.helper(varPath, nil,
			func(s []interface{}, idx int, varVal interface{}) {
				if 0 <= idx && idx < len(s) {
					result.val = s[idx]
				}
				return
			},
			func(m map[string]interface{}, key string, varVal interface{}) {
				result.val = m[key]
				return
			},
		)
		if result.pfa.Err == nil && len(varPath) > 1 {
			result = result.GetVal(varPath[1:])
		}
	}

	return
}

func (v VarValue) helper(
	varPath VarPath,
	varVal interface{},
	onSlice func(s []interface{}, idx int, varVal interface{}),
	onMap func(m map[string]interface{}, key string, varVal interface{}),
) {
	if v.val == nil {
		return
	}
	isIdx, idx, key, err := varPath[0].GetIdxKey(v.pfa)
	// fmt.Printf("helper: %s,isIdx: %s, idx: %s, key: %s, err: %s \n", varPath.formatted.String(nil), isIdx, idx, key, err)
	if err != nil {
		v.pfa.Err = err
	} else if isIdx {
		if s, ok := v.val.([]interface{}); !ok {
			v.pfa.ErrVal = helperFailed{formatted.StringFrom("%s is not <ansiOutline>Array", varPath.formattedString())}
			// v.pfa.Panic()
		} else {
			onSlice(s, idx, varVal)
		}
	} else {
		if m, ok := v.val.(map[string]interface{}); !ok {
			v.pfa.ErrVal = helperFailed{formatted.StringFrom("%s is not <ansiOutline>Map", varPath.formattedString())}
			// v.pfa.Err = bwerror.Error("<ansiPrimary>%#v<ansi> is not <ansiOutline>Map<ansi>", v)
		} else {
			onMap(m, key, varVal)
		}
	}
}

type helperFailed struct{ s formatted.String }

// type getValFailed struct{ s formatted.String }
type setValFailed struct{ s formatted.String }

func (v VarValue) SetVal(varPath VarPath, varVal interface{}) {
	if len(varPath) == 0 {
		v.pfa.Panic("varPath: %#v", varPath)
	} else {
		target := VarValue{nil, v.pfa}
		v.helper(varPath, varVal,
			func(s []interface{}, idx int, varVal interface{}) {
				if 0 > idx || idx >= len(s) {
					v.pfa.ErrVal = setValFailed{formatted.StringFrom("%d is out of range [%d, %d] of %s", idx, 0, len(s)-1, v.pfa.traceVal(varPath))}
					// v.pfa.Err = bwerror.Error("%d is out of range [%d, %d]", idx, 0, len(s)-1)
				} else {
					if len(varPath) == 1 {
						s[idx] = varVal
					} else {
						target.val = s[idx]
					}
				}
			},
			func(m map[string]interface{}, key string, varVal interface{}) {
				if len(varPath) == 1 {
					m[key] = varVal
				} else {
					if kv, ok := m[key]; !ok {
						v.pfa.Err = bwerror.Error("Map (#%v) has no key <ansiPrimary>%s<ansi>", m, key)
					} else {
						target.val = kv
					}
				}
			},
		)
		if target.pfa.Err == nil && len(varPath) > 1 {
			target.SetVal(varPath[1:], varVal)
		}
	}
}

func (v VarValue) Rune() (result rune, err error) {
	if v.pfa != nil && v.pfa.Err != nil {
		err = v.pfa.Err
	} else {
		var ok bool
		if result, ok = v.val.(rune); !ok {
			err = bwerror.Error("%#v is not rune", v.val)
		}
	}
	return
}

func (v VarValue) Int() (result int, err error) {
	if v.pfa != nil && v.pfa.Err != nil {
		err = v.pfa.Err
	} else {
		vValue := reflect.ValueOf(v.val)
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
		if result, ok = v.val.(string); !ok {
			err = bwerror.Error("<ansiPrimary>%#v<ansi> is not of type <ansiSecondary>string", v)
		}
	}
	return
}

// ============================================================================

type VarPathItem struct{ val interface{} }

func (v VarPathItem) GetIdxKey(pfa *PfaStruct) (isIdx bool, idx int, key string, err error) {
	varValue := VarValue{v.val, pfa}
	if varPath, ok := v.val.(VarPath); ok {
		if pfa == nil {
			err = bwerror.Error("VarPath requires pfa")
		} else {
			varValue = pfa.getVarValue(varPath)
			err = pfa.Err
		}
	}
	if err == nil && (pfa == nil || pfa.Err == nil) {
		var err error
		if idx, err = varValue.Int(); err == nil {
			isIdx = true
		} else if key, err = varValue.String(); err != nil {
			err = bwerror.Error("%s is nor int, neither string", varValue.val)
		}
	}
	if pfa != nil && pfa.Err != nil {
		err = pfa.Err
	}
	return
}

// ============================================================================

type VarPath []VarPathItem

func VarPathFrom(s string) (result VarPath, err error) {
	p := runeprovider.ProxyFrom(runeprovider.FromString(s))
	Stack := []VarPath{VarPath{}}
	state := "begin"
	var item string
	for {
		p.PullRune()
		currRune, isEOF := p.Rune()
		if err == nil {
			isUnexpectedRune := false
			switch state {
			case "begin":
				if isEOF {
					if len(Stack) == 1 && len(Stack[0]) == 0 {
						state = "done"
					} else {
						isUnexpectedRune = true
					}
				} else if unicode.IsDigit(currRune) {
					item = string(currRune)
					state = "idx"
				} else if currRune == '-' || currRune == '+' {
					item = string(currRune)
					state = "digit"
				} else if unicode.IsLetter(currRune) || currRune == '_' {
					item = string(currRune)
					state = "key"
				} else if currRune == '{' {
					Stack = append(Stack, VarPath{})
					state = "begin"
				} else {
					isUnexpectedRune = true
				}
			case "digit":
				if unicode.IsDigit(currRune) {
					item += string(currRune)
					state = "idx"
				} else {
					isUnexpectedRune = true
				}
			case "end":
				if isEOF {
					if len(Stack) == 1 {
						state = "done"
					} else {
						isUnexpectedRune = true
					}
				} else if currRune == '.' {
					state = "begin"
				} else if currRune == '}' && len(Stack) > 0 {
					Stack[len(Stack)-2] = append(Stack[len(Stack)-2], VarPathItem{Stack[len(Stack)-1]})
					Stack = Stack[0 : len(Stack)-1]
				} else {
					isUnexpectedRune = true
				}
			case "idx":
				if unicode.IsDigit(currRune) {
					item += string(currRune)
				} else {
					var i interface{}
					if i, err = _parseNumber(item); err == nil {
						Stack[len(Stack)-1] = append(Stack[len(Stack)-1], VarPathItem{i})
					}
					p.PushRune()
					state = "end"
				}
			case "key":
				if unicode.IsLetter(currRune) || currRune == '_' || unicode.IsDigit(currRune) {
					item += string(currRune)
				} else {
					Stack[len(Stack)-1] = append(Stack[len(Stack)-1], VarPathItem{item})
					p.PushRune()
					state = "end"
				}
			default:
				bwerror.Panic("no handler for %s", state)
			}
			if isUnexpectedRune {
				err = p.UnexpectedRuneError(fmt.Sprintf("state = %s", state))
			}
		}
		if isEOF || err != nil || (state == "done") {
			break
		}
	}
	if err == nil {
		result = Stack[0]
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

func (v VarPath) formattedString(optPfa ...*PfaStruct) formatted.String {
	var pfa *PfaStruct
	if optPfa != nil {
		pfa = optPfa[0]
	}
	ss := []string{}
	for _, i := range v {
		switch t := i.val.(type) {
		case VarPath:
			if pfa == nil {
				ss = append(ss, fmt.Sprintf("{%s}", t.formattedString(nil)))
			} else {
				ss = append(ss, fmt.Sprintf("{%s(%s)}", t.formattedString(pfa), pfa.traceVal(pfa.getVarValue(t).val)))
			}
		case string:
			ss = append(ss, t)
		default:
			vv := VarValue{t, nil}
			if _int, err := vv.Int(); err == nil {
				ss = append(ss, strconv.FormatInt(int64(_int), 10))
			}
		}
	}
	return formatted.StringFrom("<ansiCmd>%s", strings.Join(ss, "."))
}

// ============================================================================

type PfaStruct struct {
	Stack           ParseStack
	Proxy           *runeprovider.Proxy
	ErrVal          interface{}
	Err             error
	vars            map[string]interface{}
	TraceLevel      TraceLevel
	traceConditions []string
	ruleLevel       int
}

func PfaFrom(p runeprovider.RuneProvider, TraceLevel TraceLevel) *PfaStruct {
	return &PfaStruct{
		Stack:      ParseStack{},
		Proxy:      runeprovider.ProxyFrom(p),
		vars:       map[string]interface{}{},
		TraceLevel: TraceLevel,
	}
}

func (pfa *PfaStruct) TraceAction(fmtString string, fmtArgs ...interface{}) {
	fmt.Printf(pfa.indent(pfa.ruleLevel+1)+ansi.Ansi("", fmtString+"\n"), pfa.fmtArgs(fmtArgs...)...)
}

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
		result = append(result, pfa.traceVal(arg))
	}
	return result
}

// type formatted.String string

func (pfa *PfaStruct) traceVal(val interface{}) (result formatted.String) {
	// if pfa.TraceLevel > TraceNone {
	switch t := val.(type) {
	case formatted.String:
		result = t
	case rune, string:
		result = formatted.StringFrom("<ansiPrimary>%q", val)
	// result = formatted.String(fmt.Sprintf(ansi.Ansi("", "<ansiPrimary>%q"), val))
	case formatted.FormattedString:
		result = t.FormattedString()
	// case Map, Array:
	// UnicodeCategory, EOF,
	// , UnexpectedRune, UnexpectedItem, Panic
	// result = formatted.StringFrom("<ansiOutline>%s", t)
	// result = formatted.String(fmt.Sprintf(ansi.Ansi("", "<ansiOutline>%s"), t))
	case VarPath:
		// var val interface{}
		var valStr formatted.String
		if t[0].val == "rune" {
			ofs := 0
			if len(t) > 1 {
				isIdx, idx, _, err := t[1].GetIdxKey(pfa)
				if err == nil && isIdx {
					ofs = idx
				}
			}
			if r, isEOF := pfa.Proxy.Rune(ofs); isEOF {
				// val = EOF{}
				valStr = formatted.StringFrom("<ansiOutline>EOF")
			} else {
				valStr = pfa.traceVal(r)
				// val = r
			}
		} else {
			// val = pfa.getVarValue(t).val
			valStr = pfa.traceVal(pfa.getVarValue(t).val)
		}
		result = formatted.String(fmt.Sprintf("%s(%s)", t.formattedString(pfa), valStr))
	case bwset.String, bwset.Rune, bwset.Int:
		value := reflect.ValueOf(t)
		keys := value.MapKeys()
		if len(keys) == 1 {
			result = formatted.StringFrom("<ansiPrimary>%s", traceValHelper(keys[0].Interface()))
		} else if len(keys) > 1 {
			ss := []string{}
			for _, k := range keys {
				ss = append(ss, traceValHelper(k.Interface()))
			}
			result = formatted.StringFrom("<<ansiSecondary>%s>", strings.Join(ss, " "))
		}
	default:
		result = formatted.StringFrom("<ansiPrimary>%#v", val)
	}
	// }
	return
}

func traceValHelper(i interface{}) (s string) {
	switch t := i.(type) {
	case rune, string:
		s = fmt.Sprintf("%q", t)
	default:
		s = fmt.Sprintf("%d", t)
	}
	return
}

func (pfa *PfaStruct) Panic(args ...interface{}) {
	fmtString := "<ansiOutline>pfa<ansi> <ansiSecondary>%s<ansi>"
	fmtArgs := []interface{}{pfa}
	if args == nil {
		bwerror.Panicd(1, fmtString, fmtArgs...)
	} else {
		switch t := args[0].(type) {
		case string:
			fmtString += " " + t
			// fmtArgs = append(fmtArgs, args[1:]...)
			if len(args) > 1 {
				fmtArgs = append(fmtArgs, args[1:]...)
			}
			bwerror.Panicd(1, fmtString, fmtArgs...)
		case error:
			bwerror.PanicErr(fmt.Errorf(t.Error()+"\n"+ansi.Ansi("", fmtString), fmtArgs), 1)
		default:
			bwerror.Panic("%#v", args)
		}
	}
}

func (pfa *PfaStruct) getVarValue(varPath VarPath) (result VarValue) {
	// fmt.Printf("getVarValue: %s\n", varPath.formatted.String(nil))
	result = VarValue{nil, pfa}
	// if pfa.ErrVal != nil {

	// pfa.Panic("%#v", pfa.ErrVal)
	// return
	// }
	pfa.getSetHelper(varPath, nil,
		func(stackItemVars VarValue, varVal interface{}) {
			if stackItemVars.val != nil {
				result = stackItemVars.GetVal(varPath[1:])
			}
			return
		},
		func(name string, ofs int) {
			currRune, _ := pfa.Proxy.Rune(ofs)
			result.val = currRune
		},
		func(name string) {
			result.val = len(pfa.Stack)
		},
		func(pfaVars VarValue, varVal interface{}) {
			result = pfaVars.GetVal(varPath)
			return
		},
	)
	if pfa.ErrVal != nil {
		switch t := pfa.ErrVal.(type) {
		case helperFailed:
			pfa.ErrVal = nil
			pfa.Err = pfa.error(bwerror.Error("failed to get %s: "+string(t.s), varPath.formattedString()))
			// pfa.Err = pfaError{
			// 	pfa,
			// 	bwerror.Error("failed to get %s: "+string(t.s), varPath.formatted.String(nil)),
			// 	whereami.WhereAmI(2),
			// }
			// pfa.Panic(pfa.Err)
		}
	}
	return
}

func (pfa *PfaStruct) error(err error) error {
	return pfaError{
		pfa,
		err,
		// bwerror.Error("failed to get %s: "+string(t.s), varPath.formatted.String(nil)),
		whereami.WhereAmI(3),
	}
}

func (pfa *PfaStruct) getSetHelper(
	varPath VarPath,
	varVal interface{},
	onStackItemVar func(stackItemVars VarValue, varVal interface{}),
	onRune func(name string, ofs int),
	onStackLen func(name string),
	onPfaVar func(pfaVars VarValue, varVal interface{}),
) {
	if len(varPath) == 0 {
		pfa.Err = bwerror.Error("varPath is empty")
	} else {
		isIdx, idx, key, err := varPath[0].GetIdxKey(pfa)
		if err != nil {
			pfa.Err = err
		} else if isIdx {
			stackItemVars := VarValue{nil, pfa}
			if len(pfa.Stack) == 0 {
				// pfa.Err = bwerror.Error("Stack is empty")
			} else if 0 > idx || idx >= len(pfa.Stack) {
				// pfa.Err = bwerror.Error("%d is out of range [%d, %d]", idx, 0, len(pfa.Stack)-1)
			} else if len(varPath) == 1 {
				pfa.Err = bwerror.Error("%#v requires var name", varPath)
			} else {
				stackItemVars.val = pfa.GetTopStackItem(uint(idx)).Vars
			}
			if pfa.Err == nil {
				onStackItemVar(stackItemVars, varVal)
			}
		} else {
			if key == "rune" || key == "stackLen" || key == "error" {
				switch key {
				case "rune":
					var ofs int
					if len(varPath) > 2 {
						pfa.Err = bwerror.Error("%#v requires no additional VarPathItem", varPath)
					} else if len(varPath) > 1 {
						isIdx, idx, _, err := varPath[1].GetIdxKey(pfa)
						if err != nil {
							pfa.Err = err
						} else {
							if !isIdx {
								pfa.Err = bwerror.Error("%#v expects idx after rune", varPath)
							} else {
								ofs = idx
							}
						}
					}
					if pfa.Err == nil {
						onRune(key, ofs)
					}
				case "stackLen":
					if len(varPath) > 1 {
						pfa.Err = bwerror.Error("%#v requires no additional VarPathItem", varPath)
					} else {
						onStackLen(key)
					}
				}
			} else {
				onPfaVar(VarValue{pfa.vars, pfa}, varVal)
			}
		}
	}
}

func (pfa *PfaStruct) setVarVal(varPath VarPath, varVal interface{}) {
	pfa.getSetHelper(varPath, varVal,
		func(stackItemVars VarValue, varVal interface{}) {
			if stackItemVars.val == nil {
				if len(pfa.Stack) == 0 {
					pfa.Err = bwerror.Error("Stack is empty")
				} else {
					_, idx, _, _ := varPath[0].GetIdxKey(pfa)
					if 0 > idx || idx >= len(pfa.Stack) {
						pfa.Err = bwerror.Error("%d is out of range [%d, %d]", idx, 0, len(pfa.Stack)-1)
					}
				}
			} else {
				stackItemVars.SetVal(varPath[1:], varVal)
			}
		},
		func(name string, idx int) {
			pfa.Err = bwerror.Error("<ansiOutline>%s<ansi> is read only", name)
		},
		func(name string) {
			pfa.Err = bwerror.Error("<ansiOutline>%s<ansi> is read only", name)
		},
		// func(key string) {
		//  if key == "error" {
		//    pfaVars.SetVal(varPath, varVal)
		//  } else {
		//    pfa.Err = bwerror.Error("<ansiOutline>%s<ansi> is read only", key)
		//  }
		// },
		func(pfaVars VarValue, varVal interface{}) {
			pfaVars.SetVal(varPath, varVal)
		},
	)
	if pfa.ErrVal != nil {
		switch t := pfa.ErrVal.(type) {
		case helperFailed:
			pfa.ErrVal = nil
			pfa.Err = pfa.error(bwerror.Error("failed to set %s: "+string(t.s), varPath.formattedString(nil)))
			// pfa.Err = pfaError{
			// 	pfa,
			// 	bwerror.Error("failed to set %s: "+string(t.s), varPath.formatted.String(nil)),
			// 	whereami.WhereAmI(2),
			// }
			// pfa.Panic(pfa.Err)
		}
	}
}

func (pfa *PfaStruct) ifStackLen(minLen int) bool {
	return len(pfa.Stack) >= minLen
}

func (pfa *PfaStruct) mustStackLen(minLen int) {
	if !pfa.ifStackLen(minLen) {
		pfa.Panic("<ansiOutline>minLen <ansiSecondary>%d", minLen)
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

func (pfa *PfaStruct) popStackItem() {
	pfa.mustStackLen(1)
	pfa.Stack = pfa.Stack[:len(pfa.Stack)-1]
}

func (pfa *PfaStruct) pushStackItem() {
	pfa.Stack = append(pfa.Stack, ParseStackItem{
		Start: pfa.Proxy.Curr,
		Vars:  map[string]interface{}{},
	})
}

func (pfa PfaStruct) DataForJSON() interface{} {
	result := map[string]interface{}{}
	result["Stack"] = pfa.Stack.DataForJSON()
	result["p"] = pfa.Proxy.DataForJSON()
	if len(pfa.vars) > 0 {
		result["vars"] = pfa.vars
	}
	return result
}

func (pfa PfaStruct) String() string {
	return bwjson.PrettyJsonOf(pfa)
}

func (pfa *PfaStruct) traceCondition(varPath VarPath, arg interface{}, result bool) {
	if pfa.TraceLevel > TraceNone {
		fmtArgs := pfa.fmtArgs(varPath, arg)
		pfa.traceConditions = append(pfa.traceConditions, fmt.Sprintf("%s %s", fmtArgs...))
	}
}

func (pfa *PfaStruct) TraceBeginConditions() {
	pfa.traceConditions = nil
}

func (pfa *PfaStruct) TraceFailedConditions() {
	pfa.traceConditionsHelper(" <ansiErr>Failed")
}

func (pfa *PfaStruct) TraceBeginActions() {
	pfa.traceConditionsHelper("")
}

func (pfa *PfaStruct) traceConditionsHelper(suffix string) {
	fmt.Printf(pfa.indent(pfa.ruleLevel) + strings.Join(pfa.traceConditions, " ") + ansi.Ansi("", "<ansiYellow><ansiBold>:"+suffix+"\n"))
	pfa.traceConditions = nil
}

func (pfa *PfaStruct) traceIncLevel() {
	if pfa.TraceLevel > TraceNone {
		pfa.ruleLevel++
	}
}

func (pfa *PfaStruct) traceDecLevel() {
	if pfa.TraceLevel > TraceNone {
		pfa.ruleLevel--
	}
}

// ============================================================================

type pfaError struct {
	pfa *PfaStruct
	// ErrVal interface{}
	err   error
	Where string
}

func (err pfaError) Error() string {
	return err.err.Error()
}

func (v pfaError) DataForJSON() interface{} {
	result := map[string]interface{}{}
	result["pfa"] = v.pfa.DataForJSON()
	// result["ErrVal"] = v.ErrVal
	result["err"] = v.err
	result["Where"] = v.Where
	return result
}

// ============================================================================

var underscoreRegexp = regexp.MustCompile("[_]+")

func _parseNumber(source string) (value interface{}, err error) {
	source = underscoreRegexp.ReplaceAllLiteralString(source, ``)
	if strings.Contains(source, `.`) {
		var float64Val float64
		if float64Val, err = strconv.ParseFloat(source, 64); err == nil {
			value = float64Val
		}
	} else {
		var int64Val int64
		if int64Val, err = strconv.ParseInt(source, 10, 64); err == nil {
			if int64(bwint.MinInt8) <= int64Val && int64Val <= int64(bwint.MaxInt8) {
				value = int8(int64Val)
			} else if int64(bwint.MinInt16) <= int64Val && int64Val <= int64(bwint.MaxInt16) {
				value = int16(int64Val)
			} else if int64(bwint.MinInt32) <= int64Val && int64Val <= int64(bwint.MaxInt32) {
				value = int32(int64Val)
			} else {
				value = int64Val
			}
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
	Start runeprovider.RunePtrStruct
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
