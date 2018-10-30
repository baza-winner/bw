package r

import (
	"reflect"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwint"
	"github.com/baza-winner/bwcore/pfa/c"
	"github.com/baza-winner/bwcore/pfa/core"
	"github.com/baza-winner/bwcore/pfa/d"
)

// ========================= ruleCondition =====================================

type ruleCondition interface {
	ConformsTo(pfa *core.PfaStruct) bool
}

// =============================================================================

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

// =============================================================================

type rule struct {
	conditions       ruleConditions
	processorActions []core.ProcessorAction
}

// =============================================================================

type Rules []rule

func RulesFrom(args ...[]interface{}) Rules {
	result := Rules{}
	for _, arg := range args {
		result = append(result, ruleFrom(arg))
	}
	return result
}

func (rules Rules) Process(pfa *core.PfaStruct) {
	pfa.Err = nil
	pfa.Err = nil
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
				if pfa.Err != nil || pfa.Err != nil {
					break
				}
			}
			break rules
		}
	}
	// if pfa.Err == nil {
	// 	if pfa.Err != nil {
	// 		pfa.Err = pfa.Error(pfa.Err.Error(pfa))
	// 	}
	// }
	return
}

func ruleFrom(args []interface{}) rule {
	result := rule{
		ruleConditions{},
		[]core.ProcessorAction{},
	}
	varIsMap := map[string]*d.VarIs{}
	for _, arg := range args {
		switch typedArg := arg.(type) {
		case rune:
			// getVarIs(varIsMap, "rune").AddRune(typedArg)
			getVarIs(varIsMap, "rune").AddInt(int(typedArg))
		case d.EOF:
			getVarIs(varIsMap, "rune").SetIsNil()
		case d.UnicodeCategory:
			getVarIs(varIsMap, "rune").AddValChecker(typedArg)
		case core.ProccessorActionProvider:
			result.processorActions = append(result.processorActions,
				typedArg.GetAction(),
			)
		case c.VarIs:
			if len(typedArg.VarPathStr) == 0 {
				bwerror.Panic("len(typedArg.VarPathStr) == 0, typedArg: %#v", typedArg)
			}
			varPath := core.MustVarPathFrom(typedArg.VarPathStr)
			key := varPath[0].Key
			if (key == "rune" || key == "runePos") && (len(varPath) > 2) {
				bwerror.Panic("len(varPath) > 2, varPath: %s", typedArg.VarPathStr)
			}
			var needPanic bool
			varIs := getVarIs(varIsMap, typedArg.VarPathStr)

			if typedArg.VarValue == nil {
				if key == "rune" || key == "runePos" {
					needPanic = true
				} else {
					varIs.SetIsNil()
				}
			} else {
				switch t := typedArg.VarValue.(type) {
				case rune:
					if key == "runePos" {
						needPanic = true
					} else {
						// varIs.AddRune(t)
						varIs.AddInt(int(t))
					}
				case d.EOF:
					if key == "rune" {
						varIs.SetIsNil()
					} else {
						needPanic = true
					}
				case string:
					if key == "rune" || key == "runePos" {
						needPanic = true
					} else {
						varIs.AddStr(t)
					}
				case bool:
					if key == "rune" || key == "runePos" {
						needPanic = true
					} else {
						// varIs.AddValChecker(common.JustVal{typedArg.VarValue})
						varIs.AddValChecker(d.ValFrom(typedArg.VarValue))
					}
				case int:
					if key == "rune" {
						needPanic = true
					} else {
						varIs.AddInt(t)
					}
				case int8, int16 /*int32, */, int64:
					if key == "rune" {
						needPanic = true
					} else {
						_int64 := reflect.ValueOf(typedArg.VarValue).Int()
						if int64(bwint.MinInt) <= _int64 && _int64 <= int64(bwint.MaxInt) {
							varIs.AddInt(int(_int64))
						} else {
							// varIs.AddValChecker(common.JustVal{typedArg.VarValue})
							varIs.AddValChecker(d.ValFrom(typedArg.VarValue))
						}
					}
				case uint, uint8, uint16, uint32, uint64:
					needPanic = helperRuleFromUint(key, typedArg.VarValue, varIs)
				case core.ValChecker:
					// if varPathItem == "runePos" {
					// 	needPanic = true
					// } else {
					varIs.AddValChecker(t)
					// }
				case core.ValCheckerProvider:
					// if varPathItem == "runePos" {
					// 	needPanic = true
					// } else {
					varIs.AddValChecker(t.GetChecker())
					// }
				default:
					switch reflect.TypeOf(typedArg.VarValue).Kind() {
					case reflect.Uint8:
						helperRuleFromUint(key, typedArg.VarValue, varIs)
					}
				}
			}
			if needPanic {
				// bwerror.Panic("typedArg.VarValue: %#v", typedArg.VarValue)
				bwerror.Panic("typedArg.VarValue: %#v, %s", typedArg.VarValue, reflect.TypeOf(typedArg.VarValue).Kind())
			}
		default:
			bwerror.Panic("unexpected %#v", arg)
		}
	}
	for _, v := range varIsMap {
		result.conditions = append(result.conditions, v)
	}
	return result
}

func helperRuleFromUint(key, v interface{}, varIs *d.VarIs) (needPanic bool) {
	if key == "rune" {
		needPanic = true
	} else {
		_uint64 := reflect.ValueOf(v).Uint()
		if _uint64 <= uint64(bwint.MaxInt) {
			varIs.AddInt(int(_uint64))
		} else {
			varIs.AddValChecker(d.ValFrom(v))
			// varIs.AddValChecker(common.JustVal{val})
		}
	}
	return
}

func getVarIs(varIsMap map[string]*d.VarIs, varPathStr string) (result *d.VarIs) {
	result = varIsMap[varPathStr]
	if result == nil {
		result = d.VarIsFrom(varPathStr)
		varIsMap[varPathStr] = result
	}
	return
}

// =============================================================================
