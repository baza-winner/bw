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
	Val interface{}
	Err error
}

func (v VarValue) GetVal(pfa *pfaStruct, varPath VarPath) (result VarValue) {
	if v.Err != nil || len(varPath) == 0 {
		result = v
	} else {
		var val interface{}
		err := v.helper(pfa, varPath, nil,
			func(vIndex reflect.Value, varVal interface{}) (err error) {
				val = vIndex.Interface()
				return
			},
			func(vValue reflect.Value, key string, varVal interface{}) (err error) {
				keyValue := reflect.ValueOf(key)
				valueOfKey := vValue.MapIndex(keyValue)
				zeroValue := reflect.Value{}
				if valueOfKey == zeroValue {
					err = bwerror.Error("no key %s", key)
				} else {
					val = valueOfKey.Interface()
				}
				return
			},
		)
		result = VarValue{val, err}
		if err == nil && len(varPath) > 1 {
			result = result.GetVal(pfa, varPath[1:])
		}
	}
	return
}

func (v VarValue) helper(
	pfa *pfaStruct,
	varPath VarPath,
	varVal interface{},
	onSlice func(vIndex reflect.Value, varVal interface{}) error,
	onMap func(vValue reflect.Value, key string, varVal interface{}) error,
) (err error) {
	vType := reflect.TypeOf(v.Val)
	vValue := reflect.ValueOf(v.Val)
	var (
		isIdx bool
		idx   int
		key   string
	)
	isIdx, idx, key, err = pfa.GetIdxKey(varPath[0])

	if err == nil {
		if isIdx {
			if vType.Kind() != reflect.Slice {
				err = bwerror.Error("%#v is not Slice", v.Val)
			} else if 0 > idx || idx >= vValue.Len() {
				err = bwerror.Error("%d is out of range [%d, %d]", idx, 0, vValue.Len()-1)
			} else {
				vIndex := vValue.Index(idx)
				err = onSlice(vIndex, varVal)
			}
		} else if vType.Kind() != reflect.Map || vType.Key().Kind() != reflect.String {
			bwerror.Panic("%#v is not map[string]", v.Val)
			err = bwerror.Error("%#v is not map[string]", v.Val)
		} else {
			err = onMap(vValue, key, varVal)
		}
	}
	return
}

func (v VarValue) SetVal(pfa *pfaStruct, varPath VarPath, varVal interface{}) (err error) {
	if len(varPath) == 0 {
		err = bwerror.Error("varPath is empty")
	} else {
		var val interface{}
		err = v.helper(pfa, varPath, varVal,
			func(vIndex reflect.Value, varVal interface{}) (err error) {
				if len(varPath) == 1 {
					vIndex.Set(reflect.ValueOf(varVal)) // https://stackoverflow.com/questions/18115785/set-slice-index-using-reflect-in-go
				} else {
					val = vIndex.Interface()
				}
				return
			},
			func(vValue reflect.Value, key string, varVal interface{}) (err error) {
				keyValue := reflect.ValueOf(key)
				if len(varPath) == 1 {
					vValue.SetMapIndex(keyValue, reflect.ValueOf(varVal))
				} else {
					val = vValue.MapIndex(keyValue).Interface()
				}
				return
			},
		)
		if err == nil && len(varPath) > 1 {
			err = VarValue{val, err}.SetVal(pfa, varPath[1:], varVal)
		}
	}
	return
}

func (v VarValue) AsRune() (result rune, err error) {
	if v.Err != nil {
		err = v.Err
	} else {
		var ok bool
		if result, ok = v.Val.(rune); !ok {
			err = bwerror.Error("%#v is not rune", v.Val)
		}
	}
	return
}

func (v VarValue) AsInt() (result int, err error) {
	if v.Err != nil {
		err = v.Err
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

func (v VarValue) AsString() (result string, err error) {
	if v.Err != nil {
		err = v.Err
	} else {
		var ok bool
		if result, ok = v.Val.(string); !ok {
			err = bwerror.Error("<ansiPrimary>%#v<ansi> is not of type <ansiSecondary>string", v)
		}
	}
	return
}

func (pfa *pfaStruct) getVarValue(varPath VarPath) (result VarValue) {
	err := pfa.getSetHelper(varPath, nil,
		func(stackItemVars VarValue, varVal interface{}) (err error) {
			result = stackItemVars.GetVal(pfa, varPath[1:])
			return
		},
		func(key string) (err error) {
			switch key {
			case "currRune":
				currRune, _ := pfa.p.Rune()
				result = VarValue{currRune, nil}
			case "stackLen":
				result.Val = len(pfa.stack)
			}
			return
		},
		func(pfaVars VarValue, varVal interface{}) (err error) {
			result = pfaVars.GetVal(pfa, varPath)
			return
		},
	)
	if err != nil {
		result = VarValue{nil, err}
	}
	return
}

func (pfa *pfaStruct) getSetHelper(
	varPath VarPath,
	varVal interface{},
	onStackItemVar func(stackItemVars VarValue, varVal interface{}) error,
	onSpecial func(key string) error,
	onPfaVar func(pfaVars VarValue, varVal interface{}) error,
) (err error) {
	if len(varPath) == 0 {
		err = bwerror.Error("varPath is empty")
	} else {
		var (
			isIdx bool
			idx   int
			key   string
		)
		isIdx, idx, key, err = pfa.GetIdxKey(varPath[0])
		if err == nil {
			if isIdx {
				if len(pfa.stack) == 0 {
					err = bwerror.Error("stack is empty")
				} else if 0 > idx || idx >= len(pfa.stack) {
					err = bwerror.Error("%d is out of range [%d, %d]", idx, 0, len(pfa.stack)-1)
				} else if len(varPath) == 1 {
					err = bwerror.Error("%#v requires var name", varPath)
				} else {
					err = onStackItemVar(
						VarValue{pfa.getTopStackItem(uint(idx)).vars, nil},
						varVal,
					)
				}
			} else {
				if key == "currRune" || key == "stackLen" {
					if len(varPath) > 1 {
						err = bwerror.Error("%#v requires no additional VarPathItem", varPath)
					} else {
						err = onSpecial(key)
					}
				} else {
					err = onPfaVar(VarValue{pfa.vars, nil}, varVal)
				}
			}
		}
	}

	return
}

func (pfa *pfaStruct) setVarVal(varPath VarPath, varVal interface{}) (err error) {
	err = pfa.getSetHelper(varPath, varVal,
		func(stackItemVars VarValue, varVal interface{}) (err error) {
			err = stackItemVars.SetVal(pfa, varPath[1:], varVal)
			return
		},
		func(key string) (err error) {
			err = bwerror.Error("<ansiOutline>%s<ansi> is read only", key)
			return
		},
		func(pfaVars VarValue, varVal interface{}) (err error) {
			err = pfaVars.SetVal(pfa, varPath, varVal)
			return
		},
	)
	return
}

// ============================================================================

func (pfa *pfaStruct) GetIdxKey(value interface{}) (isIdx bool, idx int, key string, err error) {
	var varValue VarValue
	if varPath, ok := value.(VarPath); ok {
		varValue = pfa.getVarValue(varPath)
	} else {
		varValue = VarValue{value, nil}
	}
	if varValue.Err != nil {
		err = varValue.Err
	} else if idx, err = varValue.AsInt(); err == nil {
		isIdx = true
	} else {
		key, err = varValue.AsString()
	}
	return
}

type VarPath []interface{}

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
				} else if unicode.IsLetter(currRune) || currRune == '_' {
					item = string(currRune)
					state = "key"
				} else if currRune == '{' {
					stack = append(stack, VarPath{})
					state = "begin"
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
					stack[len(stack)-2] = append(stack[len(stack)-2], stack[len(stack)-1])
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
						stack[len(stack)-1] = append(stack[len(stack)-1], i)
					}
					p.PushRune()
					state = "end"
				}
			case "key":
				if unicode.IsLetter(currRune) || currRune == '_' || unicode.IsDigit(currRune) {
					item += string(currRune)
				} else {
					stack[len(stack)-1] = append(stack[len(stack)-1], item)
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

func createRule(args []interface{}) rule {
	result := rule{
		ruleConditions{},
		[]processorAction{},
	}
	runeChecker := hasRuneChecker{}
	currRuneValues := bwset.RuneSet{}
	currRuneVarPathValues := []VarPath{}
	varIsMap := map[string]*_varIs{}
	for _, arg := range args {
		if typedArg, ok := arg.(rune); ok {
			currRuneValues.Add(typedArg)
		} else if typedArg, ok := arg.(hasRune); ok {
			if typedArg.Len() > 0 {
				runeChecker = append(runeChecker, typedArg)
			}
		} else if typedArg, ok := arg.(VarIs); ok {
			switch typedArg.VarPathStr {
			case "currRune":
				if r, ok := typedArg.VarValue.(rune); ok {
					currRuneValues.Add(r)
				} else if r, ok := typedArg.VarValue.(hasRune); ok {
					if r.Len() > 0 {
						runeChecker = append(runeChecker, r)
					}
				} else if v, ok := typedArg.VarValue.(Var); ok {
					currRuneVarPathValues = append(currRuneVarPathValues, MustVarPathFrom(v.VarPathStr))
				} else {
					bwerror.Panic("arg: %#v", arg)
				}
			default:
				varIs := varIsMap[typedArg.VarPathStr]
				if varIs == nil {
					varIs = &_varIs{MustVarPathFrom(typedArg.VarPathStr), []ValProvider{}}
					varIsMap[typedArg.VarPathStr] = varIs
				}
				varIs.valProviders = append(varIs.valProviders, MustValProviderFrom(typedArg.VarValue))
			}
		} else if typedArg, ok := arg.(ProccessorActionProvider); ok {
			result.processorActions = append(result.processorActions,
				typedArg.GetAction(),
			)
		} else {
			bwerror.Panic("unexpected %#v", arg)
		}
	}
	if len(currRuneValues) > 0 {
		runeChecker = append(runeChecker, runeSet{currRuneValues})
	}
	if len(currRuneVarPathValues) > 0 {
		runeChecker = append(runeChecker, currRuneVarPaths{currRuneVarPathValues})
	}
	if len(runeChecker) > 0 {
		result.conditions = append(result.conditions, runeChecker)
	}
	for _, v := range varIsMap {
		result.conditions = append(result.conditions, v)
	}
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
	varPath      VarPath
	valProviders []ValProvider
}

func (v *_varIs) ConformsTo(pfa *pfaStruct) (result bool) {
	tst := pfa.getVarValue(v.varPath)
	if tst.Err != nil {
		pfa.err = tst.Err
	} else {
		for _, p := range v.valProviders {
			etaVal, err := p.GetVal(pfa)
			if err != nil {
				pfa.err = err
				break
			}
			if tst.Val == etaVal {
				result = true
				break
			}
		}
	}
	return
}

type hasRuneChecker []hasRune

func (v hasRuneChecker) ConformsTo(pfa *pfaStruct) (result bool) {
	result = false
	for _, i := range v {
		result = i.HasRune(pfa)
		if result {
			break
		}
	}
	return result
}

// =============================================================================
