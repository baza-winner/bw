package pfa

import (
	"reflect"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwint"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/pfa/core"
	"github.com/baza-winner/bwcore/pfa/e"
	"github.com/baza-winner/bwcore/runeprovider"
)

func init() {
}

// ============================================================================

// ============================================================================

func Run(p runeprovider.RuneProvider, rules Rules, optTraceLevel ...core.TraceLevel) (result interface{}, err error) {
	traceLevel := core.TraceNone
	if optTraceLevel != nil {
		traceLevel = optTraceLevel[0]
	}
	pfa := core.PfaFrom(p, traceLevel)
	for {
		rules.process(pfa)
		if pfa.Err != nil || pfa.Proxy.Curr.IsEOF {
			break
		}
	}
	if pfa.Err != nil {
		err = pfa.Err
	} else {
		if len(pfa.Stack) > 1 {
			pfa.Panic("len(pfa.Stack) > 1")
		} else if len(pfa.Stack) > 0 {
			result = pfa.GetTopStackItem().Vars["result"]
		}
	}
	return
}

func runePtr(r rune) *rune {
	return &r
}

// ============================================================================

type rule struct {
	conditions       ruleConditions
	processorActions []core.ProcessorAction
}

type Rules []rule

func CreateRules(args ...[]interface{}) Rules {
	result := Rules{}
	for _, arg := range args {
		result = append(result, createRule(arg))
	}
	return result
}

func (rules Rules) process(pfa *core.PfaStruct) {
	pfa.Err = nil
	pfa.ErrVal = nil
rules:
	for _, rule := range rules {

		if pfa.TraceLevel > core.TraceNone {
			pfa.TraceBeginConditions()
		}
		if !rule.conditions.conformsTo(pfa) {
			if pfa.TraceLevel >= core.TraceAll {
				pfa.TraceFailedConditions()
			}
		} else {
			if pfa.TraceLevel > core.TraceNone {
				pfa.TraceBeginActions()
			}
			if pfa.Err != nil {
				bwerror.Panic("r.conditions: %#v, pfa.Err: %#v", rule.conditions, pfa.Err)
			}
			for _, pa := range rule.processorActions {
				pa.Execute(pfa)
				if pfa.Err != nil || pfa.ErrVal != nil {
					break
				}
			}
			break rules
		}
	}
	if pfa.Err == nil {
		if pfa.ErrVal != nil {
			var err error
			switch t := pfa.ErrVal.(type) {
			case e.UnexpectedRune:
				err = pfa.Proxy.UnexpectedRuneError()
			case e.UnexpectedItem:
				stackItem := pfa.GetTopStackItem()
				start := stackItem.Start
				item := pfa.Proxy.Curr.Prefix[start.Pos-pfa.Proxy.Curr.PrefixStart:]
				err = pfa.Proxy.ItemError(start, "unexpected \"<ansiPrimary>%s<ansi>\"", item)
			case failedToTransformToNumber:
				err = pfa.Proxy.ItemError(pfa.GetTopStackItem().Start, t.s)
			default:
				pfa.Panic("pfa.ErrVal: %#v", pfa.ErrVal)
			}
			pfa.Err = pfa.error(err)
		}
	}
	return
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
		[]ProcessorAction{},
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

// ========================= ruleCondition =====================================

type ruleConditions []ruleCondition

func (v ruleConditions) conformsTo(pfa *core.PfaStruct) (result bool) {
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
	ConformsTo(pfa *core.PfaStruct) bool
}

type _varIs struct {
	varPath     core.VarPath
	valCheckers []ValChecker
	isNil       bool
	runeSet     bwset.Rune
	intSet      bwset.Int
	strSet      bwset.String
}

func (v *_varIs) ConformsTo(pfa *core.PfaStruct) (result bool) {
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
		i := len(pfa.Stack)
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
		if !result && pfa.Err == nil {
			for _, p := range v.valCheckers {
				if p.conforms(pfa, varValue.val, v.varPath) {
					result = true
					break
				} else if pfa.Err != nil {
					break
				}
			}
		}
	}
	return
}

// =============================================================================
