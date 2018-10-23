package pfa

import (
	"fmt"
	"reflect"
	"unicode"

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
	err   error
	vars  map[string]interface{}
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

func Run(p runeprovider.RuneProvider, logicDef Rules) (result interface{}, err error) {
	pfa := pfaStruct{
		stack: parseStack{},
		p:     runeprovider.ProxyFrom(p),
		vars:  map[string]interface{}{},
	}
	for {
		pfa.processRules(logicDef)
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
	if args != nil {
		fmtString += " " + args[0].(string)
	}
	fmtArgs := []interface{}{pfa}
	if len(args) > 1 {
		fmtArgs = append(fmtArgs, args[1:])
	}
	bwerror.Panicd(1, fmtString, fmtArgs...)
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
	if v.pfa.err != nil || len(varPath) == 0 {
		result = v
	} else {
		result = VarValue{nil, v.pfa}
		v.helper(varPath, nil,
			func(vIndex reflect.Value, varVal interface{}) {
				result.val = vIndex.Interface()
				return
			},
			func(vValue reflect.Value, key string, varVal interface{}) {
				keyValue := reflect.ValueOf(key)
				valueOfKey := vValue.MapIndex(keyValue)
				zeroValue := reflect.Value{}
				if valueOfKey == zeroValue {
					v.pfa.err = bwerror.Error("no key %s", key)
				} else {
					result.val = valueOfKey.Interface()
				}
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
	onSlice func(vIndex reflect.Value, varVal interface{}),
	onMap func(vValue reflect.Value, key string, varVal interface{}),
) {
	vType := reflect.TypeOf(v.val)
	vValue := reflect.ValueOf(v.val)
	isIdx, idx, key, err := varPath[0].GetIdxKey(v.pfa)
	if err != nil {
		v.pfa.err = err
	} else {
		if isIdx {
			if vType.Kind() != reflect.Slice {
				v.pfa.err = bwerror.Error("%#v is not Slice", v)
			} else if 0 > idx || idx >= vValue.Len() {
				v.pfa.err = bwerror.Error("%d is out of range [%d, %d]", idx, 0, vValue.Len()-1)
			} else {
				vIndex := vValue.Index(idx)
				onSlice(vIndex, varVal)
			}
		} else if vType.Kind() != reflect.Map || vType.Key().Kind() != reflect.String {
			bwerror.Panic("%#v is not map[string]", v)
			v.pfa.err = bwerror.Error("%#v is not map[string]", v)
		} else {
			onMap(vValue, key, varVal)
		}
	}
}

func (v VarValue) SetVal(varPath VarPath, varVal interface{}) {
	if len(varPath) == 0 {
		v.pfa.err = bwerror.Error("varPath is empty")
	} else {
		target := VarValue{nil, v.pfa}
		v.helper(varPath, varVal,
			func(vIndex reflect.Value, varVal interface{}) {
				if len(varPath) == 1 {
					vIndex.Set(reflect.ValueOf(varVal)) // https://stackoverflow.com/questions/18115785/set-slice-index-using-reflect-in-go
				} else {
					target.val = vIndex.Interface()
				}
			},
			func(vValue reflect.Value, key string, varVal interface{}) {
				keyValue := reflect.ValueOf(key)
				if len(varPath) == 1 {
					vValue.SetMapIndex(keyValue, reflect.ValueOf(varVal))
				} else {
					target.val = vValue.MapIndex(keyValue).Interface()
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
	result = VarValue{nil, pfa}
	pfa.getSetHelper(varPath, nil,
		func(stackItemVars VarValue, varVal interface{}) {
			result = stackItemVars.GetVal(varPath[1:])
			return
		},
		func(key string) {
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
					currRune, _ := pfa.p.Rune(ofs)
					result = VarValue{currRune, pfa}
				}
			case "stackLen":
				if len(varPath) > 1 {
					pfa.err = bwerror.Error("%#v requires no additional VarPathItem", varPath)
				} else {
					result.val = len(pfa.stack)
				}
			}
			return
		},
		func(pfaVars VarValue, varVal interface{}) {
			result = pfaVars.GetVal(varPath)
			return
		},
	)
	return
}

func (pfa *pfaStruct) getSetHelper(
	varPath VarPath,
	varVal interface{},
	onStackItemVar func(stackItemVars VarValue, varVal interface{}),
	onSpecial func(key string),
	onPfaVar func(pfaVars VarValue, varVal interface{}),
) {
	if len(varPath) == 0 {
		pfa.err = bwerror.Error("varPath is empty")
	} else {
		isIdx, idx, key, err := varPath[0].GetIdxKey(pfa)
		if err != nil {
			pfa.err = err
		} else if isIdx {
			if len(pfa.stack) == 0 {
				pfa.err = bwerror.Error("stack is empty")
			} else if 0 > idx || idx >= len(pfa.stack) {
				pfa.err = bwerror.Error("%d is out of range [%d, %d]", idx, 0, len(pfa.stack)-1)
			} else if len(varPath) == 1 {
				pfa.err = bwerror.Error("%#v requires var name", varPath)
			} else {
				onStackItemVar(
					VarValue{pfa.getTopStackItem(uint(idx)).vars, pfa},
					varVal,
				)
			}
		} else {
			if key == "rune" || key == "stackLen" {
				onSpecial(key)
			} else {
				onPfaVar(VarValue{pfa.vars, pfa}, varVal)
			}
		}
	}
}

func (pfa *pfaStruct) setVarVal(varPath VarPath, varVal interface{}) {
	pfa.getSetHelper(varPath, varVal,
		func(stackItemVars VarValue, varVal interface{}) {
			stackItemVars.SetVal(varPath[1:], varVal)
		},
		func(key string) {
			pfa.err = bwerror.Error("<ansiOutline>%s<ansi> is read only", key)
		},
		func(pfaVars VarValue, varVal interface{}) {
			pfaVars.SetVal(varPath, varVal)
		},
	)
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
	if v.runeSet == nil {
		v.intSet = bwset.Int{}
	}
	v.intSet.Add(i)
}

func (v *_varIs) AddStr(s string) {
	if v.runeSet == nil {
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
		if typedArg, ok := arg.(rune); ok {
			getVarIs(varIsMap, "rune").AddRune(typedArg)
		} else if _, ok := arg.(EOF); ok {
			getVarIs(varIsMap, "rune").isEOF = true
		} else if typedArg, ok := arg.(UnicodeCategory); ok {
			getVarIs(varIsMap, "rune").AddValChecker(typedArg)
		} else if typedArg, ok := arg.(VarIs); ok {
			if len(typedArg.VarPathStr) == 0 {
				bwerror.Panic("len(typedArg.VarPathStr) == 0, typedArg: %#v", typedArg)
			}
			varPath := MustVarPathFrom(typedArg.VarPathStr)
			if varPath[0].val == "rune" {
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
					varIs.isEOF = true
				case ValChecker:
					vc, _ := typedArg.VarValue.(ValChecker)
					varIs.AddValChecker(vc)
				default:
					bwerror.Panic("typedArg.VarValue: %#v", typedArg.VarValue)
				}
			} else if varPath[0].val == "stackLen" {
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
			} else {
				varIs := getVarIs(varIsMap, typedArg.VarPathStr)
				switch typedArg.VarValue.(type) {
				case rune:
					r, _ := typedArg.VarValue.(rune)
					varIs.AddRune(r)
				case EOF:
					bwerror.Panic("EOF is appliable only for rune, varPath: %s", typedArg.VarPathStr)
				case int:
					i, _ := typedArg.VarValue.(int)
					varIs.AddInt(i)
				case string:
					s, _ := typedArg.VarValue.(string)
					varIs.AddStr(s)
				case int8, int16 /*int32, */, int64:
					_int64 := reflect.ValueOf(typedArg.VarValue).Int()
					if int64(bwint.MinInt) <= _int64 && _int64 <= int64(bwint.MaxInt) {
						varIs.AddInt(int(_int64))
					} else {
						varIs.AddValChecker(justVal{typedArg.VarValue})
					}
				case uint, uint8, uint16, uint32, uint64:
					_uint64 := reflect.ValueOf(typedArg.VarValue).Uint()
					if _uint64 <= uint64(bwint.MaxInt) {
						varIs.AddInt(int(_uint64))
					} else {
						varIs.AddValChecker(justVal{typedArg.VarValue})
					}
				case ValChecker:
					vc, _ := typedArg.VarValue.(ValChecker)
					varIs.AddValChecker(vc)
				default:
					bwerror.Panic("typedArg.VarValue: %#v", typedArg.VarValue)
				}
			}

			// var val interface{}
			// if varPath[0]
			// if b, ok = typedArg.VarValue.(bool) {

			// }
			// switch typedArg.VarPathStr {
			// case "currRune":
			// 	if r, ok := typedArg.VarValue.(rune); ok {
			// 		currRuneValues.Add(r)
			// 	} else if r, ok := typedArg.VarValue.(hasRune); ok {
			// 		if r.Len() > 0 {
			// 			runeChecker = append(runeChecker, r)
			// 		}
			// 	} else if v, ok := typedArg.VarValue.(Var); ok {
			// 		currRuneVarPathValues = append(currRuneVarPathValues, MustVarPathFrom(v.VarPathStr))
			// 	} else {
			// 		bwerror.Panic("arg: %#v", arg)
			// 	}
			// default:
			// varIs := varIsMap[typedArg.VarPathStr]
			// if varIs == nil {
			// 	varIs = &_varIs{MustVarPathFrom(typedArg.VarPathStr), []ValChecker{}, false, nil}
			// 	varIsMap[typedArg.VarPathStr] = varIs
			// }
			// varIs := getVarIs(varIsMap, typedArg.VarPathStr)
			// varIs.valCheckers = append(varIs.valCheckers, MustValProviderFrom(typedArg.VarValue))
			// }
		} else if typedArg, ok := arg.(ProccessorActionProvider); ok {
			result.processorActions = append(result.processorActions,
				typedArg.GetAction(),
			)
		} else {
			bwerror.Panic("unexpected %#v", arg)
		}
	}
	// if len(currRuneValues) > 0 {
	// 	runeChecker = append(runeChecker, runeSet{currRuneValues})
	// }
	// if len(currRuneVarPathValues) > 0 {
	// 	runeChecker = append(runeChecker, currRuneVarPaths{currRuneVarPathValues})
	// }
	// if len(runeChecker) > 0 {
	// 	result.conditions = append(result.conditions, runeChecker)
	// }
	for _, v := range varIsMap {
		result.conditions = append(result.conditions, v)
	}
	bwerror.Spew.Printf("result: %#v\n", result)
	return result
}

func (pfa *pfaStruct) processRules(def Rules) {
	pfa.err = nil
def:
	for _, r := range def {
		if r.conditions.conformsTo(pfa) {
			for _, pa := range r.processorActions {
				pa.execute(pfa)
				if pfa.err != nil || pfa.vars["error"] != nil {
					break
				}
			}
			break def
		}
	}
	errVal := pfa.vars["error"]
	if errVal != nil {
		if errName, ok := errVal.(string); ok && len(errName) > 0 {
			var errStr string
			if errName == "unexpectedRune" {
				errStr = pfa.p.UnexpectedRuneError().Error()
			} else if errName == "unknownWord" {
				stackItem := pfa.getTopStackItem()
				itemString, _ := stackItem.vars["string"].(string)
				errStr = pfa.p.WordError("unknown word <ansiPrimary>%s<ansi>", itemString, stackItem.start).Error()
			} else {
				bwerror.Unreachable("errName: " + errName)
			}
			pfa.err = pfaError{pfa, errName, errStr, whereami.WhereAmI(2)}
		}
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
	isEOF       bool
	runeSet     bwset.Rune
	intSet      bwset.Int
	strSet      bwset.String
}

func (v *_varIs) ConformsTo(pfa *pfaStruct) (result bool) {
	if v.varPath[0].val == "rune" {
		var ofs int
		if len(v.varPath) > 1 {
			ofs, _ = v.varPath[1].val.(int)
		}
		r, isEOF := pfa.p.Rune(ofs)
		if v.isEOF && isEOF {
			result = true
		} else if !isEOF {
			if v.runeSet != nil && v.runeSet.Has(r) {
				result = true
			} else {
				for _, p := range v.valCheckers {
					if p.conforms(pfa, r) {
						result = true
						break
					}
				}
			}
		}
	} else if v.varPath[0].val == "stackLen" {
		i := len(pfa.stack)
		if v.intSet != nil && v.intSet.Has(i) {
			result = true
		} else {
			for _, p := range v.valCheckers {
				if p.conforms(pfa, r) {
					result = true
					break
				}
			}
		}
	} else {
		varValue := pfa.getVarValue(v.varPath)
		if pfa.err == nil {
			for _, p := range v.valCheckers {
				if p.conforms(pfa, varValue.val) {
					result = true
					break
				} else if pfa.err != nil {
					break
				}
			}
		}
	}
	// tst := pfa.getVarValue(v.varPath)
	// if tst.Err != nil {
	// 	pfa.err = tst.Err
	// } else {
	// 	for _, p := range v.valCheckers {
	// 		etaVal, err := p.GetVal(pfa)
	// 		if err != nil {
	// 			pfa.err = err
	// 			break
	// 		}
	// 		if tst.Val == etaVal {
	// 			result = true
	// 			break
	// 		}
	// 	}
	// }
	return
}

// type hasRuneChecker []hasRune

// func (v hasRuneChecker) ConformsTo(pfa *pfaStruct) (result bool) {
// 	result = false
// 	for _, i := range v {
// 		result = i.HasRune(pfa)
// 		if result {
// 			break
// 		}
// 	}
// 	return result
// }

// =============================================================================
