package core

import (
	"reflect"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwdebug"
	"github.com/baza-winner/bwcore/bwerr"
)

// ============================================================================

type VarValue struct {
	Val interface{}
	pfa *PfaStruct
}

func VarValueFrom(val interface{}) VarValue {
	return VarValue{val, nil}
}

// func (v VarValue) Rune() (result rune, err error) {
// 	if v.pfa != nil && v.pfa.Err != nil {
// 		err = v.pfa.Err
// 	} else {
// 		var ok bool
// 		if result, ok = v.Val.(rune); !ok {
// 			err = bwerr.Error("%#v is not rune", v.Val)
// 		}
// 	}
// 	return
// }

func (v VarValue) Int() (result int, ok bool) {
	// if v.pfa != nil && v.pfa.Err != nil {
	// 	err = v.pfa.Err
	// } else {
	vValue := reflect.ValueOf(v.Val)
	switch vValue.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		_int64 := vValue.Int()
		if int64(bw.MinInt) <= _int64 && _int64 <= int64(bw.MaxInt) {
			result = int(_int64)
			ok = true
			// } else {
			// 	err = bwerr.Error("%d is out of range [%d, %d]", _int64, bw.MinInt, bw.MaxInt)
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		_uint64 := vValue.Uint()
		if _uint64 <= uint64(bw.MaxInt) {
			result = int(_uint64)
			ok = true
			// } else {
			// 	err = bwerr.Error("%d is more than %d", _uint64, bw.MaxInt)
		}
		// default:
		// 	err = bwerr.Error("<ansiVal>%#v<ansi> is not of type <ansiVal>int", v)
	}
	// }
	return
}

func (v VarValue) String() (result string, ok bool) {
	result, ok = v.Val.(string)
	// if v.pfa != nil && v.pfa.Err != nil {
	// 	err = v.pfa.Err
	// } else {
	// var ok bool
	// if ; !ok {
	// 	err = bwerr.Error("<ansiVal>%#v<ansi> is not of type <ansiVal>string", v)
	// }
	// }
	return
}

// ============================================================================

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
	onVal func(vt valType, vals []interface{}, m map[string]interface{}, idx int, key string, VarVal interface{}) (err error),
	onLen func(vt valType, vals []interface{}, m map[string]interface{}) (err error),
) (err error) {
	if varPath[0].Type == VarPathItemHash {
		if v.Val == nil {
			err = onLen(valTypeNil, nil, nil)
		} else {
			switch t := v.Val.(type) {
			case []interface{}:
				err = onLen(valTypeSlice, t, nil)
			case map[string]interface{}:
				err = onLen(valTypeSlice, nil, t)
			default:
				err = v.pfa.Error("%s nor <ansiVar>Array, neither <ansiVar>Map", varPath.FormattedString())
			}
		}
	} else if v.Val == nil {
		err = onVal(valTypeNil, nil, nil, 0, "", VarVal)
	} else {
		vt, idx, key, err := varPath[0].TypeIdxKey(v.pfa)
		if err == nil {
			if vt == VarPathItemIdx {
				if vals, ok := v.Val.([]interface{}); !ok {
					// v.pfa.Err = PfaError{formatted.StringFrom("%s is not <ansiVar>Array", varPath.FormattedString())}
					err = v.pfa.Error("%s is not <ansiVar>Array", varPath.FormattedString())
				} else {
					onVal(valTypeSlice, vals, nil, idx, "", VarVal)
				}
			} else if m, ok := v.Val.(map[string]interface{}); !ok {
				err = v.pfa.Error("%s is not <ansiVar>Map", varPath.FormattedString())
			} else {
				onVal(valTypeMap, nil, m, 0, key, VarVal)
			}
		}
	}
	return
}

func (v VarValue) GetVal(varPath VarPath) (result VarValue, err error) {
	// fmt.Printf("GetVal: %s\n", varPath.formatted.String(nil))

	if len(varPath) == 3 &&
		varPath[0].Key == "stack" && varPath[1].Idx == -1 && varPath[2].Key == "itemPos" {
		bwdebug.Print("!here", "varPath", varPath, "v.pfa.Err", v.pfa.Err)
	}
	if len(varPath) == 0 {
		result = v
	} else {
		result = VarValue{nil, v.pfa}
		err = v.helper(varPath, nil,
			func(vt valType, vals []interface{}, m map[string]interface{}, idx int, key string, VarVal interface{}) (err error) {
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
				return
			},
			func(vt valType, vals []interface{}, m map[string]interface{}) (err error) {
				switch vt {
				case valTypeSlice:
					result.Val = len(vals)
				case valTypeMap:
					result.Val = len(m)
				case valTypeNil:
					result.Val = 0
				default:
					bwerr.Unreachable()
				}
				return
			},
		)
		if err == nil && len(varPath) > 1 {
			result, err = result.GetVal(varPath[1:])
		}
	}

	return
}

func (v VarValue) SetVal(varPath VarPath, VarVal interface{}) (err error) {
	if len(varPath) == 0 {
		v.pfa.PanicA(bw.Fmt("varPath: %#v", varPath))
	} else {
		target := VarValue{nil, v.pfa}
		err = v.helper(varPath, VarVal,
			func(vt valType, vals []interface{}, m map[string]interface{}, idx int, key string, VarVal interface{}) (err error) {
				switch vt {
				case valTypeSlice:
					if len(vals) == 0 {
						// v.pfa.Err = PfaError{formatted.StringFrom("path does not exist (no elem with idx <ansiVal>%d<ansi> at empty Array)", idx)}
						// bwerr.Panic("%#v", varPath)
						err = v.pfa.Error("path does not exist (no elem with idx <ansiVal>%d<ansi> at empty Array)", idx)
					} else {
						minIdx := -len(vals)
						maxIdx := len(vals) - 1
						if !(minIdx <= idx && idx <= maxIdx) {
							// v.pfa.Err = PfaError{formatted.StringFrom("path does not exist (<ansiVal>%d<ansi> is out of range <ansiVal>[%d, %d]<ansi>)", idx, minIdx, maxIdx)}
							err = v.pfa.Error("path does not exist (<ansiVal>%d<ansi> is out of range <ansiVal>[%d, %d]<ansi>)", idx, minIdx, maxIdx)
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
						// v.pfa.Err = PfaError{formatted.StringFrom("path does not exist (no key <ansiVal>%s)", key)}
						err = v.pfa.Error("path does not exist (no key <ansiVal>%s)", key)
					} else {
						target.Val = kv
					}
				case valTypeNil:
					// v.pfa.Err = PfaError{formatted.StringFrom("can not set to nil value")}
					err = v.pfa.Error("can not set to nil value")
				}
				return
			},
			func(vt valType, vals []interface{}, m map[string]interface{}) (err error) {
				return v.pfa.Error("<ansiVar>path.#<ansi> is <ansiPath>readonly")
			},
		)
		if err == nil && len(varPath) > 1 {
			target.SetVal(varPath[1:], VarVal)
		}
	}
	return
}

// ============================================================================

func (pfa *PfaStruct) getSetHelper(
	varPath VarPath,
	VarVal interface{},
	onSpecial func(name string, ofs int) (err error),
	onPfaVar func(pfaVars VarValue, VarVal interface{}) (err error),
) (err error) {
	if len(varPath) == 0 {
		err = pfa.Error("varPath is empty")
	} else {
		vt, _, key, err := varPath[0].TypeIdxKey(pfa)
		if err != nil {
			pfa.Err = err
		} else if vt != VarPathItemKey {
			err = pfa.Error("path must start with key")
		} else if key == "rune" || key == "runePos" {
			var ofs int
			if len(varPath) > 2 {
				// pfa.Err = PfaError{formatted.StringFrom("<ansiVal>%s<ansi> path may have at most 2 items", key)}
				err = pfa.Error("<ansiVal>%s<ansi> path may have at most 2 items", key)
			} else if len(varPath) > 1 {
				vt, idx, key, err := varPath[1].TypeIdxKey(pfa)
				if err != nil {
					pfa.Err = err
				} else if vt != VarPathItemIdx {
					// pfa.Err = PfaError{formatted.StringFrom("<ansiVal>%s<ansi> path expects <ansiVar>idx<ansi> as second item", key)}
					err = pfa.Error("<ansiVal>%s<ansi> path expects <ansiVar>idx<ansi> as second item", key)
				} else {
					ofs = idx
				}
			}
			if pfa.Err == nil {
				err = onSpecial(key, ofs)
			}
		} else {
			err = onPfaVar(VarValue{pfa.Vars, pfa}, VarVal)
		}
	}
	return
}

func (pfa *PfaStruct) VarValue(varPath VarPath) (result VarValue, err error) {
	result = VarValue{nil, pfa}
	if len(varPath) > 0 {
		err = pfa.getSetHelper(varPath, nil,
			func(name string, ofs int) (err error) {
				switch name {
				case "rune":
					currRune, _ := pfa.Proxy.Rune(ofs)
					result.Val = currRune
				case "runePos":
					ps := pfa.Proxy.PosInfo(ofs)
					result.Val = ps
				}
				return
			},
			func(pfaVars VarValue, VarVal interface{}) (err error) {
				result, err = pfaVars.GetVal(varPath)
				if len(varPath) == 3 &&
					varPath[0].Key == "stack" && varPath[1].Idx == -1 && varPath[2].Key == "itemPos" {
					bwdebug.Print("!HERE", "varPath", varPath, "result", result)
				}
				return
			},
		)
		if err != nil {
			if t, ok := pfa.Err.(PfaError); ok && t.State() == PecsNeedPrepare {
				t.PrepareErr("failed to get %s", varPath.FormattedString())
			} else {
				bwerr.Panic("%#v", err)
			}
		}
	}
	// bwerr.Log("<ansiVar>VarValue<ansi>: %s: %#v\n", varPath.FormattedString(), result.Val)
	// if len(varPath) == 3 &&
	// 	varPath[0].Key == "stack" && varPath[1].Idx == -1 && varPath[2].Key == "itemPos" {
	// 	// && reflect.TypeOf(result.Val).Kind() == reflect.Map {
	// 	bwerr.Debug(varPath.FormattedString(), result.Val)
	// 	// pfa.Panic()
	// }
	return
}

func (pfa *PfaStruct) SetVarVal(varPath VarPath, VarVal interface{}) (err error) {
	if len(varPath) == 0 {
		err = pfa.Error("varPath is empty")
	} else {
		err = pfa.getSetHelper(varPath, VarVal,
			func(name string, idx int) (err error) {
				err = pfa.Error("<ansiVar>%s<ansi> is <ansiPath>readonly", name)
				// pfa.Err = bwerr.Error("<ansiVar>%s<ansi> is read only", name)
				return
			},
			func(pfaVars VarValue, VarVal interface{}) (err error) {
				pfaVars.SetVal(varPath, VarVal)
				return
			},
		)
	}
	if err != nil {
		if t, ok := pfa.Err.(PfaError); ok && t.State() == PecsNeedPrepare {
			t.PrepareErr("failed to get %s", varPath.FormattedString())
		} else {
			bwerr.Panic("%#v", err)
		}
	}
	return
}

// ============================================================================
