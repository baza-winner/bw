package d

import (
	"reflect"
	"unicode"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/pfa/core"
	"github.com/baza-winner/bwcore/pfa/formatted"
)

// ============================================================================

// func valProviderFrom(val interface{}) (result core.ValProvider, err error) {
// 	common.ProcessVal(val,
// 		func(vk common.ValKind, varPath core.VarPath, valueOf reflect.Value) {
// 			switch vk {
// 			case common.VkInvalid:
// 				bwerr.Panic("%#v", val)
// 			default:
// 				result = common.Val{val}
// 			}
// 		},
// 	)
// 	return
// }

// func MustValProviderFrom(i interface{}) (result core.ValProvider) {
// 	var err error
// 	if result, err = valProviderFrom(i); err != nil {
// 		bwerr.PanicErr(err)
// 	}
// 	return
// }

// ============================================================================

type UnicodeCategory uint8

//go:generate stringer -type=UnicodeCategory

const (
	Space UnicodeCategory = iota
	Letter
	Lower
	Upper
	Digit
	OpenBraces
	Punct
	Symbol
)

func (v UnicodeCategory) Conforms(pfa *core.PfaStruct, val interface{}, varPath core.VarPath) (result bool, err error) {
	if r, ok := val.(rune); ok {
		switch v {
		case Space:
			result = unicode.IsSpace(r)
		case Letter:
			result = unicode.IsLetter(r) || r == '_'
		case Lower:
			result = unicode.IsLower(r)
		case Upper:
			result = unicode.IsUpper(r)
		case Digit:
			result = unicode.IsDigit(r)
		case OpenBraces:
			result = r == '(' || r == '{' || r == '[' || r == '<'
		case Punct:
			result = unicode.IsPunct(r)
		case Symbol:
			result = unicode.IsSymbol(r)
		default:
			bwerr.Panic("UnicodeCategory: %s", v)
		}
	}
	if pfa.TraceLevel > core.TraceNone {
		pfa.TraceCondition(varPath, v, result)
	}
	return
}

func (t UnicodeCategory) FormattedString() formatted.String {
	return formatted.StringFrom("<ansiVar>%s", t)
}

// ============================================================================

type EOF struct{}

func (v EOF) String() string {
	return "EOF"
}

func (t EOF) FormattedString() formatted.String {
	return formatted.StringFrom("<ansiVar>%s", t)
}

// ============================================================================

type Var struct {
	VarPathStr string
}

func (t Var) GetChecker() core.ValChecker {
	return Val{core.MustVarPathFrom(t.VarPathStr)}
}

// ============================================================================

type Val struct{ Val interface{} }

func ValFrom(val interface{}) (result Val) {
	ProcessVal(val,
		func(vk ValKind, varPath core.VarPath, valueOf reflect.Value) {
			switch vk {
			case VkInvalid:
				bwerr.Panic("%#v", val)
			case VkVarPath:
				result = Val{varPath}
			default:
				result = Val{valueOf.Interface()}
			}
		},
	)
	return
}

func (v Val) Conforms(pfa *core.PfaStruct, val interface{}, varPath core.VarPath) (result bool, err error) {
	tst, err := v.GetVal(pfa)
	if err != nil {
		return
	}
	result = reflect.DeepEqual(tst, val)
	return
}

func (v Val) GetVal(pfa *core.PfaStruct) (result interface{}, err error) {
	ProcessVal(v.Val,
		func(vk ValKind, varPath core.VarPath, valueOf reflect.Value) {
			switch vk {
			case VkVarPath:
				var varValue core.VarValue
				varValue, err = pfa.VarValue(varPath)
				if pfa.Err == nil {
					result = varValue.Val
				}
			case VkMap:
				m := map[string]interface{}{}
				for _, keyValue := range valueOf.MapKeys() {
					m[keyValue.String()], err = Val{valueOf.MapIndex(keyValue).Interface()}.GetVal(pfa)
					if err != nil {
						return
					}
				}
				result = m
				// pfa.Panic(bw.StructFrom("%#v", result))
			case VkSlice:
				len := valueOf.Len()
				vals := make([]interface{}, 0, len)
				for i := 0; i < len; i++ {
					val, err := Val{valueOf.Index(i).Interface()}.GetVal(pfa)
					if err != nil {
						return
					}
					vals = append(vals, val)
				}
				result = vals
			case VkInvalid:
				pfa.PanicA(bw.Fmt("%#v", v.Val))
			default:
				result = v.Val
			}
		},
	)
	// fmt.Printf("result: %#v", result)
	return
}

func (v Val) GetSource(pfa *core.PfaStruct) formatted.String {
	return pfa.TraceVal(v.Val)
}

// ============================================================================

type ValKind uint8

const (
	VkInvalid ValKind = iota
	VkNil
	VkBool
	VkString
	VkInt
	VkUint
	VkFloat
	VkComplex
	VkMap
	VkSlice
	VkVarPath
	VkEOF
)

func ProcessVal(
	val interface{},
	onVal func(vk ValKind, varPath core.VarPath, valueOf reflect.Value),
) {
	vk := VkInvalid
	var varPath core.VarPath
	var valueOf reflect.Value
	var ok bool
	if val == nil {
		vk = VkNil
	} else if varPath, ok = val.(core.VarPath); ok {
		vk = VkVarPath
	} else if t, ok := val.(Var); ok {
		varPath = core.MustVarPathFrom(t.VarPathStr)
		vk = VkVarPath
	} else {
		valueOf = reflect.ValueOf(val)
		switch valueOf.Kind() {
		case reflect.Bool:
			vk = VkBool
		case reflect.String:
			vk = VkString
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vk = VkInt
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vk = VkUint
		case reflect.Float32, reflect.Float64:
			vk = VkFloat
		case reflect.Complex64, reflect.Complex128:
			vk = VkComplex
		case reflect.Map:
			typeOf := reflect.TypeOf(val)
			if typeOf.Key().Kind() == reflect.String && typeOf.Elem().Kind() == reflect.Interface {
				vk = VkMap
			}
		case reflect.Slice:
			typeOf := reflect.TypeOf(val)
			if typeOf.Elem().Kind() == reflect.Interface {
				vk = VkSlice
			}
		}
	}
	onVal(vk, varPath, valueOf)
}

// ============================================================================

type VarIs struct {
	varPath     core.VarPath
	valCheckers []core.ValChecker
	isNil       bool
	// runeSet     bwset.Rune
	intSet bwset.Int
	strSet bwset.String
}

func VarIsFrom(varPathStr string) *VarIs {
	return &VarIs{varPath: core.MustVarPathFrom(varPathStr)}
}

func (v *VarIs) SetIsNil() {
	v.isNil = true
}

// func (v *VarIs) AddRune(r rune) {
// 	if v.runeSet == nil {
// 		v.runeSet = bwset.Rune{}
// 	}
// 	v.runeSet.Add(r)
// }

func (v *VarIs) AddInt(i int) {
	if v.intSet == nil {
		v.intSet = bwset.Int{}
	}
	v.intSet.Add(i)
}

func (v *VarIs) AddStr(s string) {
	if v.strSet == nil {
		v.strSet = bwset.String{}
	}
	v.strSet.Add(s)
}

func (v *VarIs) AddValChecker(vc core.ValChecker) {
	if v.valCheckers == nil {
		v.valCheckers = []core.ValChecker{}
	}
	v.valCheckers = append(v.valCheckers, vc)
}

func (v *VarIs) ConformsTo(pfa *core.PfaStruct) (result bool, err error) {
	if v.varPath[0].Type == core.VarPathItemKey && (v.varPath[0].Key == "rune" || v.varPath[0].Key == "runePos") {
		var ofs int
		if len(v.varPath) > 2 {
			err = pfa.Error("len(varPath) > 2, varPath: %s", v.varPath.FormattedString())
			// bwerr.Panic("len(varPath) > 2, varPath: %s", typedArg.VarPathStr)
		} else if len(v.varPath) > 1 {
			vt, idx, key, err := v.varPath[1].TypeIdxKey(pfa)
			if err != nil {
				pfa.Err = err
			} else if vt != core.VarPathItemIdx {
				err = pfa.Error("<ansiVal>%s<ansi> path expects <ansiVar>idx<ansi> as second item", key)
			} else {
				ofs = idx
			}
		}
		// if len(varPath) > 2 {
		// bwerr.Panic
		// }
		// if len(v.varPath) > 1 {
		// 	ofs, _ = core.VarValueFrom(v.varPath[1].Val).Int()
		// }
		switch v.varPath[0].Key {
		case "rune":
			r, isEOF := pfa.Proxy.Rune(ofs)
			if v.isNil {
				result = isEOF
				if pfa.TraceLevel > core.TraceNone {
					pfa.TraceCondition(v.varPath, formatted.StringFrom("<ansiVar>EOF"), result)
				}
			}
			if !result && !isEOF {
				if v.intSet != nil {
					result = v.intSet.Has(int(r))
					if pfa.TraceLevel > core.TraceNone {
						pfa.TraceCondition(v.varPath, v.intSet, result)
					}
				}
				if !result {
					var ok bool
					for _, p := range v.valCheckers {
						if ok, err = p.Conforms(pfa, r, v.varPath); err != nil {
							return
						} else if ok {
							result = true
							break
						}
					}
				}
			}
		case "runePos":
			// i := len(pfa.Stack)
			ps := pfa.Proxy.PosStruct(ofs)
			if v.intSet != nil {
				result = v.intSet.Has(ps.Pos)
				if pfa.TraceLevel > core.TraceNone {
					pfa.TraceCondition(v.varPath, v.intSet, result)
				}
			}
			if !result {
				var ok bool
				for _, p := range v.valCheckers {
					if ok, err = p.Conforms(pfa, ps, v.varPath); err != nil {
						return
					} else if ok {
						result = true
						break
					}
				}
			}
		}
	} else if v.varPath[len(v.varPath)-1].Type == core.VarPathItemHash {
		var varValue core.VarValue
		if varValue, err = pfa.VarValue(v.varPath); err != nil {
			return
		}
		i, _ := varValue.Int()
		// i := len(pfa.Stack)
		if v.intSet != nil {
			result = v.intSet.Has(i)
			if pfa.TraceLevel > core.TraceNone {
				pfa.TraceCondition(v.varPath, v.intSet, result)
			}
		}
		if !result {
			var ok bool
			for _, p := range v.valCheckers {
				if ok, err = p.Conforms(pfa, i, v.varPath); err != nil {
					return
				} else if ok {
					result = true
					break
				}
			}
			// for _, p := range v.valCheckers {
			// 	if p.Conforms(pfa, i, v.varPath) {
			// 		result = true
			// 		break
			// 	}
			// }
		}
	} else {
		var varValue core.VarValue
		if varValue, err = pfa.VarValue(v.varPath); err != nil {
			return
		}
		// varValue := pfa.VarValue(v.varPath)
		if varValue.Val == nil {
			result = v.isNil
			if pfa.TraceLevel > core.TraceNone {
				pfa.TraceCondition(v.varPath, nil, result)
			}
			// } else if r, err := varValue.Rune(); err == nil {
			// 	if v.runeSet != nil {
			// 		result = v.runeSet.Has(r)
			// 		if pfa.TraceLevel > core.TraceNone {
			// 			pfa.TraceCondition(v.varPath, v.runeSet, result)
			// 		}
			// 	}
		} else if s, ok := varValue.String(); ok {
			if v.strSet != nil {
				result = v.strSet.Has(s)
				if pfa.TraceLevel > core.TraceNone {
					pfa.TraceCondition(v.varPath, v.strSet, result)
				}
			}
		} else if i, ok := varValue.Int(); ok {
			if v.intSet != nil {
				result = v.intSet.Has(i)
				if pfa.TraceLevel > core.TraceNone {
					pfa.TraceCondition(v.varPath, v.intSet, result)
				}
			}
		}
		if !result && pfa.Err == nil {
			var ok bool
			for _, p := range v.valCheckers {
				if ok, err = p.Conforms(pfa, varValue.Val, v.varPath); err != nil {
					return
				} else if ok {
					result = true
					break
				}
			}
			// for _, p := range v.valCheckers {
			// 	if p.Conforms(pfa, varValue.Val, v.varPath) {
			// 		result = true
			// 		break
			// 	} else if pfa.Err != nil {
			// 		break
			// 	}
			// }
		}
	}
	return
}

// ============================================================================
