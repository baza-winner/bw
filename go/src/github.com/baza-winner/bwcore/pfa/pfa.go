package pfa

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwint"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/runeprovider"
	"github.com/jimlawless/whereami"
)

//go:generate stringer -type=UnicodeCategory

func init() {
}

// ============================================================================

type parseStackItem struct {
	start runeprovider.RunePtrStruct
	vars  map[string]interface{}
}

func (stackItem *parseStackItem) DataForJSON() interface{} {
	result := map[string]interface{}{}
	result["start"] = stackItem.start.DataForJSON()
	result["vars"] = stackItem.vars
	return result
}

func (stackItem *parseStackItem) String() (result string) {
	return bwjson.PrettyJsonOf(stackItem)
}

// ============================================================================

type parseStack []parseStackItem

func (stack *parseStack) DataForJSON() interface{} {
	result := []interface{}{}
	for _, item := range *stack {
		result = append(result, item.DataForJSON())
	}
	return result
}

func (stack *parseStack) String() (result string) {
	return bwjson.PrettyJsonOf(stack)
}

// ============================================================================

type pfaStruct struct {
	stack parseStack
	p     *runeprovider.Proxy
	// isUnexpectedRune bool
	errVal interface{}
	err    error
	vars   map[string]interface{}
	// trace support
	traceLevel      TraceLevel
	traceConditions []string
	ruleLevel       int
}

func pfaFrom(p runeprovider.RuneProvider, traceLevel TraceLevel) *pfaStruct {
	return &pfaStruct{
		stack:      parseStack{},
		p:          runeprovider.ProxyFrom(p),
		vars:       map[string]interface{}{},
		traceLevel: traceLevel,
	}
}

func (pfa pfaStruct) DataForJSON() interface{} {
	result := map[string]interface{}{}
	result["stack"] = pfa.stack.DataForJSON()
	result["p"] = pfa.p.DataForJSON()
	if len(pfa.vars) > 0 {
		result["vars"] = pfa.vars
	}
	return result
}

func (pfa pfaStruct) String() string {
	return bwjson.PrettyJsonOf(pfa)
}

type TraceLevel uint8

const (
	TraceNone TraceLevel = iota
	TraceBrief
	TraceAll
)

func Run(p runeprovider.RuneProvider, rules Rules, optTraceLevel ...TraceLevel) (result interface{}, err error) {
	traceLevel := TraceNone
	if optTraceLevel != nil {
		traceLevel = optTraceLevel[0]
	}
	pfa := pfaFrom(p, traceLevel)
	for {
		pfa.processRules(rules)
		if pfa.err != nil || pfa.p.Curr.IsEOF {
			break
		}
	}
	if pfa.err != nil {
		err = pfa.err
	} else {
		if len(pfa.stack) > 1 {
			pfa.panic("len(pfa.stack) > 1")
		} else if len(pfa.stack) > 0 {
			result = pfa.getTopStackItem().vars["result"]
		}
	}
	return
}

func runePtr(r rune) *rune {
	return &r
}

func (pfa *pfaStruct) panic(args ...interface{}) {
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

func (pfa *pfaStruct) ifStackLen(minLen int) bool {
	return len(pfa.stack) >= minLen
}

func (pfa *pfaStruct) mustStackLen(minLen int) {
	if !pfa.ifStackLen(minLen) {
		pfa.panic("<ansiOutline>minLen <ansiSecondary>%d", minLen)
	}
}

func (pfa *pfaStruct) getTopStackItem(optDeep ...uint) *parseStackItem {
	ofs := -1
	if optDeep != nil {
		ofs = ofs - int(optDeep[0])
	}
	pfa.mustStackLen(-ofs)
	return &pfa.stack[len(pfa.stack)+ofs]
}

func (pfa *pfaStruct) popStackItem() {
	pfa.mustStackLen(1)
	pfa.stack = pfa.stack[:len(pfa.stack)-1]
}

func (pfa *pfaStruct) pushStackItem() {
	pfa.stack = append(pfa.stack, parseStackItem{
		start: pfa.p.Curr,
		vars:  map[string]interface{}{},
	})
}

type VarValue struct {
	val interface{}
	pfa *pfaStruct
}

func (v VarValue) GetVal(varPath VarPath) (result VarValue) {
	// fmt.Printf("GetVal: %s\n", varPath.formattedString(nil))
	if v.pfa.err != nil || len(varPath) == 0 {
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
		if result.pfa.err == nil && len(varPath) > 1 {
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
	// fmt.Printf("helper: %s,isIdx: %s, idx: %s, key: %s, err: %s \n", varPath.formattedString(nil), isIdx, idx, key, err)
	if err != nil {
		v.pfa.err = err
	} else if isIdx {
		if s, ok := v.val.([]interface{}); !ok {
			v.pfa.errVal = helperFailed{formattedStringFrom("%s is not <ansiOutline>Array", varPath.formattedString())}
			// v.pfa.panic()
		} else {
			onSlice(s, idx, varVal)
		}
	} else {
		if m, ok := v.val.(map[string]interface{}); !ok {
			v.pfa.errVal = helperFailed{formattedStringFrom("%s is not <ansiOutline>Map", varPath.formattedString())}
			// v.pfa.err = bwerror.Error("<ansiPrimary>%#v<ansi> is not <ansiOutline>Map<ansi>", v)
		} else {
			onMap(m, key, varVal)
		}
	}
}

type helperFailed struct{ s formattedString }
type getValFailed struct{ s formattedString }
type setValFailed struct{ s formattedString }

func (v VarValue) SetVal(varPath VarPath, varVal interface{}) {
	if len(varPath) == 0 {
		v.pfa.panic("varPath: %#v", varPath)
	} else {
		target := VarValue{nil, v.pfa}
		v.helper(varPath, varVal,
			func(s []interface{}, idx int, varVal interface{}) {
				if 0 > idx || idx >= len(s) {
					v.pfa.errVal = setValFailed{
						formattedStringFrom("%d is out of range [%d, %d] of %s", idx, 0, len(s)-1, v.pfa.traceVal(varPath)),
					}
					// v.pfa.err = bwerror.Error("%d is out of range [%d, %d]", idx, 0, len(s)-1)
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
						v.pfa.err = bwerror.Error("Map (#%v) has no key <ansiPrimary>%s<ansi>", m, key)
					} else {
						target.val = kv
					}
				}
			},
		)
		if target.pfa.err == nil && len(varPath) > 1 {
			target.SetVal(varPath[1:], varVal)
		}
	}
}

func (v VarValue) Rune() (result rune, err error) {
	if v.pfa != nil && v.pfa.err != nil {
		err = v.pfa.err
	} else {
		var ok bool
		if result, ok = v.val.(rune); !ok {
			err = bwerror.Error("%#v is not rune", v.val)
		}
	}
	return
}

func (v VarValue) Int() (result int, err error) {
	if v.pfa != nil && v.pfa.err != nil {
		err = v.pfa.err
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
	if v.pfa != nil && v.pfa.err != nil {
		err = v.pfa.err
	} else {
		var ok bool
		if result, ok = v.val.(string); !ok {
			err = bwerror.Error("<ansiPrimary>%#v<ansi> is not of type <ansiSecondary>string", v)
		}
	}
	return
}

func (pfa *pfaStruct) getVarValue(varPath VarPath) (result VarValue) {
	// fmt.Printf("getVarValue: %s\n", varPath.formattedString(nil))
	result = VarValue{nil, pfa}
	// if pfa.errVal != nil {

	// pfa.panic("%#v", pfa.errVal)
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
			currRune, _ := pfa.p.Rune(ofs)
			result.val = currRune
		},
		func(name string) {
			result.val = len(pfa.stack)
		},
		func(pfaVars VarValue, varVal interface{}) {
			result = pfaVars.GetVal(varPath)
			return
		},
	)
	if pfa.errVal != nil {
		switch t := pfa.errVal.(type) {
		case helperFailed:
			pfa.errVal = nil
			pfa.err = pfaError{
				pfa,
				bwerror.Error("failed to get %s: "+string(t.s), varPath.formattedString(nil)),
				whereami.WhereAmI(2),
			}
			pfa.panic(pfa.err)
		}
	}
	return
}

func (pfa *pfaStruct) getSetHelper(
	varPath VarPath,
	varVal interface{},
	onStackItemVar func(stackItemVars VarValue, varVal interface{}),
	onRune func(name string, ofs int),
	onStackLen func(name string),
	onPfaVar func(pfaVars VarValue, varVal interface{}),
) {
	if len(varPath) == 0 {
		pfa.err = bwerror.Error("varPath is empty")
	} else {
		isIdx, idx, key, err := varPath[0].GetIdxKey(pfa)
		if err != nil {
			pfa.err = err
		} else if isIdx {
			stackItemVars := VarValue{nil, pfa}
			if len(pfa.stack) == 0 {
				// pfa.err = bwerror.Error("stack is empty")
			} else if 0 > idx || idx >= len(pfa.stack) {
				// pfa.err = bwerror.Error("%d is out of range [%d, %d]", idx, 0, len(pfa.stack)-1)
			} else if len(varPath) == 1 {
				pfa.err = bwerror.Error("%#v requires var name", varPath)
			} else {
				stackItemVars.val = pfa.getTopStackItem(uint(idx)).vars
			}
			if pfa.err == nil {
				onStackItemVar(stackItemVars, varVal)
			}
		} else {
			if key == "rune" || key == "stackLen" || key == "error" {
				switch key {
				case "rune":
					var ofs int
					if len(varPath) > 2 {
						pfa.err = bwerror.Error("%#v requires no additional VarPathItem", varPath)
					} else if len(varPath) > 1 {
						isIdx, idx, _, err := varPath[1].GetIdxKey(pfa)
						if err != nil {
							pfa.err = err
						} else {
							if !isIdx {
								pfa.err = bwerror.Error("%#v expects idx after rune", varPath)
							} else {
								ofs = idx
							}
						}
					}
					if pfa.err == nil {
						onRune(key, ofs)
					}
				case "stackLen":
					if len(varPath) > 1 {
						pfa.err = bwerror.Error("%#v requires no additional VarPathItem", varPath)
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

func (pfa *pfaStruct) setVarVal(varPath VarPath, varVal interface{}) {
	pfa.getSetHelper(varPath, varVal,
		func(stackItemVars VarValue, varVal interface{}) {
			if stackItemVars.val == nil {
				if len(pfa.stack) == 0 {
					pfa.err = bwerror.Error("stack is empty")
				} else {
					_, idx, _, _ := varPath[0].GetIdxKey(pfa)
					if 0 > idx || idx >= len(pfa.stack) {
						pfa.err = bwerror.Error("%d is out of range [%d, %d]", idx, 0, len(pfa.stack)-1)
					}
				}
			} else {
				stackItemVars.SetVal(varPath[1:], varVal)
			}
		},
		func(name string, idx int) {
			pfa.err = bwerror.Error("<ansiOutline>%s<ansi> is read only", name)
		},
		func(name string) {
			pfa.err = bwerror.Error("<ansiOutline>%s<ansi> is read only", name)
		},
		// func(key string) {
		// 	if key == "error" {
		// 		pfaVars.SetVal(varPath, varVal)
		// 	} else {
		// 		pfa.err = bwerror.Error("<ansiOutline>%s<ansi> is read only", key)
		// 	}
		// },
		func(pfaVars VarValue, varVal interface{}) {
			pfaVars.SetVal(varPath, varVal)
		},
	)
	if pfa.errVal != nil {
		switch t := pfa.errVal.(type) {
		case helperFailed:
			pfa.errVal = nil
			pfa.err = pfaError{
				pfa,
				bwerror.Error("failed to set %s: "+string(t.s), varPath.formattedString(nil)),
				whereami.WhereAmI(2),
			}
			// pfa.panic(pfa.err)
		}
	}
}

// ============================================================================

type VarPathItem struct{ val interface{} }

func (v VarPathItem) GetIdxKey(pfa *pfaStruct) (isIdx bool, idx int, key string, err error) {
	varValue := VarValue{v.val, pfa}
	if varPath, ok := v.val.(VarPath); ok {
		if pfa == nil {
			err = bwerror.Error("VarPath requires pfa")
		} else {
			varValue = pfa.getVarValue(varPath)
			err = pfa.err
		}
	}
	if err == nil && (pfa == nil || pfa.err == nil) {
		var err error
		if idx, err = varValue.Int(); err == nil {
			isIdx = true
		} else if key, err = varValue.String(); err != nil {
			err = bwerror.Error("%s is nor int, neither string", varValue.val)
		}
	}
	if pfa != nil && pfa.err != nil {
		err = pfa.err
	}
	return
}

type VarPath []VarPathItem

func VarPathFrom(s string) (result VarPath, err error) {
	p := runeprovider.ProxyFrom(runeprovider.FromString(s))
	stack := []VarPath{VarPath{}}
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
					if len(stack) == 1 && len(stack[0]) == 0 {
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
					stack = append(stack, VarPath{})
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
					if len(stack) == 1 {
						state = "done"
					} else {
						isUnexpectedRune = true
					}
				} else if currRune == '.' {
					state = "begin"
				} else if currRune == '}' && len(stack) > 0 {
					stack[len(stack)-2] = append(stack[len(stack)-2], VarPathItem{stack[len(stack)-1]})
					stack = stack[0 : len(stack)-1]
				} else {
					isUnexpectedRune = true
				}
			case "idx":
				if unicode.IsDigit(currRune) {
					item += string(currRune)
				} else {
					var i interface{}
					if i, err = _parseNumber(item); err == nil {
						stack[len(stack)-1] = append(stack[len(stack)-1], VarPathItem{i})
					}
					p.PushRune()
					state = "end"
				}
			case "key":
				if unicode.IsLetter(currRune) || currRune == '_' || unicode.IsDigit(currRune) {
					item += string(currRune)
				} else {
					stack[len(stack)-1] = append(stack[len(stack)-1], VarPathItem{item})
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

func (v VarPath) formattedString(optPfa ...*pfaStruct) formattedString {
	var pfa *pfaStruct
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
	return formattedStringFrom("<ansiCmd>%s", strings.Join(ss, "."))
}

// ============================================================================

type rule struct {
	conditions       ruleConditions
	processorActions []processorAction
}

type Rules []rule

func CreateRules(args ...[]interface{}) Rules {
	result := Rules{}
	for _, arg := range args {
		result = append(result, createRule(arg))
	}
	return result
}

func getVarIs(varIsMap map[string]*_varIs, varPathStr string) *_varIs {
	varIs := varIsMap[varPathStr]
	if varIs == nil {
		varIs = &_varIs{varPath: MustVarPathFrom(varPathStr)}
		varIsMap[varPathStr] = varIs
	}
	return varIs
}

func (v *_varIs) AddRune(r rune) {
	if v.runeSet == nil {
		v.runeSet = bwset.Rune{}
	}
	v.runeSet.Add(r)
}

func (v *_varIs) AddInt(i int) {
	if v.intSet == nil {
		v.intSet = bwset.Int{}
	}
	v.intSet.Add(i)
}

func (v *_varIs) AddStr(s string) {
	if v.strSet == nil {
		v.strSet = bwset.String{}
	}
	v.strSet.Add(s)
}

func (v *_varIs) AddValChecker(vc ValChecker) {
	if v.valCheckers == nil {
		v.valCheckers = []ValChecker{}
	}
	v.valCheckers = append(v.valCheckers, vc)
}

func createRule(args []interface{}) rule {
	result := rule{
		ruleConditions{},
		[]processorAction{},
	}
	varIsMap := map[string]*_varIs{}
	for _, arg := range args {
		switch typedArg := arg.(type) {
		case rune:
			getVarIs(varIsMap, "rune").AddRune(typedArg)
		case EOF:
			getVarIs(varIsMap, "rune").isNil = true
		case UnicodeCategory:
			getVarIs(varIsMap, "rune").AddValChecker(typedArg)
		case ProccessorActionProvider:
			result.processorActions = append(result.processorActions,
				typedArg.GetAction(),
			)
		case VarIs:
			if len(typedArg.VarPathStr) == 0 {
				bwerror.Panic("len(typedArg.VarPathStr) == 0, typedArg: %#v", typedArg)
			}
			varPath := MustVarPathFrom(typedArg.VarPathStr)
			switch varPath[0].val {
			case "rune":
				if len(varPath) > 2 {
					bwerror.Panic("len(varPath) > 2, varPath: %s", typedArg.VarPathStr)
				} else if len(varPath) > 1 {
					isIdx, _, _, err := varPath[1].GetIdxKey(nil)
					if err != nil {
						bwerror.PanicErr(err)
					} else {
						if !isIdx {
							bwerror.Panic("varPath[1] expects to be idx, varPath: %#v, %s", varPath[1], reflect.TypeOf(varPath[1]).Kind())
						} else if varPath[1].val == 0 {
							typedArg.VarPathStr = "rune"
						}
					}
				}
				varIs := getVarIs(varIsMap, typedArg.VarPathStr)
				switch typedArg.VarValue.(type) {
				case rune:
					r, _ := typedArg.VarValue.(rune)
					varIs.AddRune(r)
				case EOF:
					varIs.isNil = true
				case ValChecker:
					vc, _ := typedArg.VarValue.(ValChecker)
					varIs.AddValChecker(vc)
				case Var:
					vp, _ := typedArg.VarValue.(Var)
					if varPath, err := VarPathFrom(vp.VarPathStr); err == nil {
						varIs.AddValChecker(varVal{varPath})
					} else {
						bwerror.PanicErr(err)
					}
				default:
					bwerror.Panic("arg: %#v", arg)
				}
			case "stackLen":
				if len(varPath) > 1 {
					bwerror.Panic("len(varPath) > 2, varPath: %s", typedArg.VarPathStr)
				}
				varIs := getVarIs(varIsMap, typedArg.VarPathStr)
				switch typedArg.VarValue.(type) {
				case int:
					i, _ := typedArg.VarValue.(int)
					varIs.AddInt(i)
				default:
					bwerror.Panic("typedArg.VarValue: %#v", typedArg.VarValue)
				}
			default:
				varIs := getVarIs(varIsMap, typedArg.VarPathStr)
				if typedArg.VarValue == nil {
					varIs.isNil = true
				} else {
					switch t := typedArg.VarValue.(type) {
					case rune:
						varIs.AddRune(t)
					case EOF:
						bwerror.Panic("EOF is appliable only for rune, varPath: %s", typedArg.VarPathStr)
					case int:
						varIs.AddInt(t)
					case string:
						varIs.AddStr(t)
					case int8, int16 /*int32, */, int64:
						_int64 := reflect.ValueOf(typedArg.VarValue).Int()
						if int64(bwint.MinInt) <= _int64 && _int64 <= int64(bwint.MaxInt) {
							varIs.AddInt(int(_int64))
						} else {
							varIs.AddValChecker(justVal{typedArg.VarValue})
						}
					case bool:
						varIs.AddValChecker(justVal{typedArg.VarValue})
					case uint, uint8, uint16, uint32, uint64:
						_uint64 := reflect.ValueOf(typedArg.VarValue).Uint()
						if _uint64 <= uint64(bwint.MaxInt) {
							varIs.AddInt(int(_uint64))
						} else {
							varIs.AddValChecker(justVal{typedArg.VarValue})
						}
					case ValChecker:
						varIs.AddValChecker(t)
					default:
						bwerror.Panic("typedArg.VarValue: %#v", typedArg.VarValue)
					}
				}
			}
		default:
			bwerror.Panic("unexpected %#v", arg)
		}
	}
	for _, v := range varIsMap {
		result.conditions = append(result.conditions, v)
	}
	// bwerror.Spew.Printf("result: %#v\n", result)
	return result
}

func (pfa *pfaStruct) processRules(rules Rules) {
	pfa.err = nil
	pfa.errVal = nil
rules:
	for _, rule := range rules {
		pfa.traceBeginConditions()
		if !rule.conditions.conformsTo(pfa) {
			pfa.traceFailedConditions()
		} else {
			pfa.traceBeginActions()
			if pfa.err != nil {
				bwerror.Panic("r.conditions: %#v, pfa.err: %#v", rule.conditions, pfa.err)
			}
			for _, pa := range rule.processorActions {
				pa.execute(pfa)
				if pfa.err != nil {
					break
				} else if pfa.vars["error"] != nil {
					break
				}
			}
			break rules
		}
	}
	if pfa.err == nil {
		if pfa.errVal != nil {
			var err error
			switch t := pfa.errVal.(type) {
			case UnexpectedRune:
				err = pfa.p.UnexpectedRuneError()
			case UnexpectedItem:
				stackItem := pfa.getTopStackItem()
				start := stackItem.start
				item := pfa.p.Curr.Prefix[start.Pos-pfa.p.Curr.PrefixStart:]
				err = pfa.p.ItemError(start, "unexpected \"<ansiPrimary>%s<ansi>\"", item)
			case failedToTransformToNumber:
				stackItem := pfa.getTopStackItem()
				start := stackItem.start
				err = pfa.p.ItemError(start, t.s)
			default:
				pfa.panic("pfa.errVal: %#v", pfa.errVal)
			}
			pfa.err = pfaError{
				pfa,
				err,
				whereami.WhereAmI(2),
			}
		}
	}
	return
}

func (pfa *pfaStruct) traceCondition(varPath VarPath, arg interface{}, result bool) {
	if pfa.traceLevel > TraceNone {
		fmtArgs := pfa.fmtArgs(varPath, arg)
		pfa.traceConditions = append(pfa.traceConditions, fmt.Sprintf("%s %s", fmtArgs...))
	}
}

func (pfa *pfaStruct) traceAction(fmtString string, fmtArgs ...interface{}) {
	if pfa.traceLevel > TraceNone {
		fmt.Printf(pfa.indent(pfa.ruleLevel+1)+ansi.Ansi("", fmtString+"\n"), pfa.fmtArgs(fmtArgs...)...)
	}
}

func (pfa *pfaStruct) traceBeginConditions() {
	if pfa.traceLevel > TraceNone {
		pfa.traceConditions = nil
	}
}

func (pfa *pfaStruct) traceFailedConditions() {
	if pfa.traceLevel >= TraceAll {
		pfa.traceConditionsHelper(" <ansiErr>Failed")
	}
}

func (pfa *pfaStruct) traceBeginActions() {
	pfa.traceConditionsHelper("")
}

func (pfa *pfaStruct) traceConditionsHelper(suffix string) {
	if pfa.traceLevel > TraceNone {
		fmt.Printf(pfa.indent(pfa.ruleLevel) + strings.Join(pfa.traceConditions, " ") + ansi.Ansi("", "<ansiYellow><ansiBold>:"+suffix+"\n"))
		pfa.traceConditions = nil
	}
}

func (pfa *pfaStruct) indent(indentLevel int) string {
	indentAtom := "  "
	indent := ""
	for i := 0; i <= indentLevel; i++ {
		indent += indentAtom
	}
	return indent
}

func (pfa *pfaStruct) traceIncLevel() {
	if pfa.traceLevel > TraceNone {
		pfa.ruleLevel++
	}
}

func (pfa *pfaStruct) traceDecLevel() {
	if pfa.traceLevel > TraceNone {
		pfa.ruleLevel--
	}
}

func (pfa *pfaStruct) fmtArgs(fmtArgs ...interface{}) []interface{} {
	result := []interface{}{}
	for _, arg := range fmtArgs {
		if f, ok := arg.(func(pfa *pfaStruct) interface{}); ok {
			arg = f(pfa)
		}
		result = append(result, pfa.traceVal(arg))
	}
	return result
}

type formattedString string

func (pfa *pfaStruct) traceVal(val interface{}) (result formattedString) {
	// if pfa.traceLevel > TraceNone {
	switch t := val.(type) {
	case formattedString:
		result = t
	case rune, string:
		result = formattedStringFrom("<ansiPrimary>%q", val)
		// result = formattedString(fmt.Sprintf(ansi.Ansi("", "<ansiPrimary>%q"), val))
	case UnicodeCategory, EOF, Map, Array, UnexpectedRune, UnexpectedItem, Panic:
		result = formattedStringFrom("<ansiOutline>%s", t)
		// result = formattedString(fmt.Sprintf(ansi.Ansi("", "<ansiOutline>%s"), t))
	case VarPath:
		var val interface{}
		if t[0].val == "rune" {
			ofs := 0
			if len(t) > 1 {
				isIdx, idx, _, err := t[1].GetIdxKey(pfa)
				if err == nil && isIdx {
					ofs = idx
				}
			}
			if r, isEOF := pfa.p.Rune(ofs); isEOF {
				val = EOF{}
			} else {
				val = r
			}
		} else {
			val = pfa.getVarValue(t).val
		}
		result = formattedString(fmt.Sprintf("%s(%s)", t.formattedString(pfa), pfa.traceVal(val)))
	case bwset.String, bwset.Rune, bwset.Int:
		value := reflect.ValueOf(t)
		keys := value.MapKeys()
		if len(keys) == 1 {
			result = formattedStringFrom("<ansiPrimary>%s", traceValHelper(keys[0].Interface()))
		} else if len(keys) > 1 {
			ss := []string{}
			for _, k := range keys {
				ss = append(ss, traceValHelper(k.Interface()))
			}
			result = formattedStringFrom("<<ansiSecondary>%s>", strings.Join(ss, " "))
		}
	default:
		result = formattedStringFrom("<ansiPrimary>%#v", val)
	}
	// }
	return
}

func formattedStringFrom(fmtString string, fmtArgs ...interface{}) formattedString {
	return formattedString(fmt.Sprintf(ansi.Ansi("", fmtString), fmtArgs...))
}

func (v formattedString) concat(s formattedString) formattedString {
	return formattedString(string(v) + string(s))
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

// ========================= ruleCondition =====================================

type ruleConditions []ruleCondition

func (v ruleConditions) conformsTo(pfa *pfaStruct) (result bool) {
	result = true
	for _, i := range v {
		if !i.ConformsTo(pfa) {
			result = false
			break
		}
	}
	return
}

type ruleCondition interface {
	ConformsTo(pfa *pfaStruct) bool
}

type _varIs struct {
	varPath     VarPath
	valCheckers []ValChecker
	isNil       bool
	runeSet     bwset.Rune
	intSet      bwset.Int
	strSet      bwset.String
}

func (v *_varIs) ConformsTo(pfa *pfaStruct) (result bool) {
	if v.varPath[0].val == "rune" {
		var ofs int
		if len(v.varPath) > 1 {
			ofs, _ = VarValue{v.varPath[1].val, nil}.Int()
		}
		r, isEOF := pfa.p.Rune(ofs)
		if v.isNil {
			result = isEOF
			pfa.traceCondition(v.varPath, EOF{}, result)
		}
		if !result && !isEOF {
			if v.runeSet != nil {
				result = v.runeSet.Has(r)
				pfa.traceCondition(v.varPath, v.runeSet, result)
			}
			if !result {
				for _, p := range v.valCheckers {
					if p.conforms(pfa, r, v.varPath) {
						result = true
						break
					}
				}
			}
		}
	} else if v.varPath[0].val == "stackLen" {
		i := len(pfa.stack)
		if v.intSet != nil {
			result = v.intSet.Has(i)
			pfa.traceCondition(v.varPath, v.intSet, result)
		}
		if !result {
			for _, p := range v.valCheckers {
				if p.conforms(pfa, i, v.varPath) {
					result = true
					break
				}
			}
		}
	} else {
		varValue := pfa.getVarValue(v.varPath)
		if varValue.val == nil {
			result = v.isNil
			pfa.traceCondition(v.varPath, nil, result)
		} else if r, err := varValue.Rune(); err == nil {
			if v.runeSet != nil {
				result = v.runeSet.Has(r)
				pfa.traceCondition(v.varPath, v.runeSet, result)
			}
		} else if i, err := varValue.Int(); err == nil {
			if v.intSet != nil {
				result = v.intSet.Has(i)
				pfa.traceCondition(v.varPath, v.intSet, result)
			}
		} else if s, err := varValue.String(); err == nil {
			if v.strSet != nil {
				result = v.strSet.Has(s)
				pfa.traceCondition(v.varPath, v.strSet, result)
			}
		}
		if !result && pfa.err == nil {
			for _, p := range v.valCheckers {
				if p.conforms(pfa, varValue.val, v.varPath) {
					result = true
					break
				} else if pfa.err != nil {
					break
				}
			}
		}
	}
	return
}

// =============================================================================
