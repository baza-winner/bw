package core

import (
	"reflect"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwfmt"
	"github.com/baza-winner/bwcore/bwint"
)

// ============================================================================

type VarValue struct {
	Val interface{}
	pfa *PfaStruct
}

func VarValueFrom(val interface{}) VarValue {
	return VarValue{val, nil}
}

func (pfa *PfaStruct) VarValue(varPath VarPath) (result VarValue) {
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
				if t.State() == PecsNeedPrepare {
					t.PrepareErr("failed to get %s", varPath.FormattedString())
				}
			}
		}
	}
	return
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
						if idx < 0 {
							idx = len(vals) + idx
						}
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
			onVal(valTypeMap, nil, m, 0, key, VarVal)
		}
	}
}

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
						// bwerror.Panic("%#v", varPath)
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

// ============================================================================

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
				t.PrepareErr("failed to set %s, %#v", varPath.FormattedString(), varPath)
			}
			// pfa.Err = pfa.Error(bwerror.Error("failed to set %s: "+string(t.s), varPath.FormattedString(nil)))
		}
	}
}

// ============================================================================
