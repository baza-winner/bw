package bwval

import (
	"encoding/json"
	"strconv"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwdebug"
	"github.com/baza-winner/bwcore/bwerr"
	bwjson "github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwstr"
)

// ============================================================================

type Holder struct {
	Val interface{}
	Pth bw.ValPath
}

func HolderFrom(s string, optVars ...map[string]interface{}) Holder {
	return Holder{Val: From(s, optVars...)}
}

func hasOptional(path bw.ValPath) bool {
	for _, vpi := range path {
		if vpi.IsOptional {
			return true
		}
	}
	return false
}

// PathVal - реализация интерфейса bw.Val
func (v *Holder) PathVal(path bw.ValPath, optVars ...map[string]interface{}) (result interface{}, err error) {
	if len(path) == 0 {
		result = v.Val
		return
	}
	defer func() {
		if err != nil {
			result = nil
		}
	}()

	var simplePath bw.ValPath
	simplePath, err = v.simplifyPath(path, optVars)
	if err != nil {
		return
	}

	if path[0].Type == bw.ValPathItemVar {
		varName := path[0].Key
		var target interface{}
		var ok bool
		if ok = len(optVars) > 0; ok {
			target, ok = optVars[0][varName]
		}
		if !ok && !hasOptional(path) {
			err = bwerr.From(ansi.String("var <ansiVar>%s<ansi> is not defined"), varName)
			return
		}
		h := Holder{Val: target}
		return h.PathVal(simplePath[1:])
	}

	result = v.Val
	for i, vpi := range simplePath {
		switch vpi.Type {
		case bw.ValPathItemKey:
			result, err = Holder{result, path[:i]}.KeyVal(vpi.Key,
				func() (result interface{}, ok bool) {
					ok = hasOptional(path[i:])
					return
				},
			)
		case bw.ValPathItemIdx:
			result, err = Holder{result, path[:i]}.IdxVal(vpi.Idx,
				func() (result interface{}, ok bool) {
					ok = hasOptional(path[i:])
					return
				},
			)
		case bw.ValPathItemHash:
			if result == nil {
				result = 0
			} else {
				switch t := result.(type) {
				case map[string]interface{}:
					result = len(t)
				case []interface{}:
					result = len(t)
				default:
					// err = Holder{result, path[:i]}.notOfValKindError("Map", "Array")
					err = Holder{result, path[:i]}.notOfValKindError(ValKindSetFrom(ValMap, ValArray))
				}
			}
		}
		if err != nil {
			return
		}
	}
	return
}

func (v Holder) Path(pathProvider bw.ValPathProvider, optVars ...map[string]interface{}) (result Holder, err error) {
	var val interface{}
	var path bw.ValPath
	if path, err = pathProvider.Path(); err != nil {
		return
	}
	if val, err = (&v).PathVal(path, optVars...); err != nil {
		return
	}
	result = Holder{val, path}
	return
}

func (v Holder) MustPath(pathProvider bw.ValPathProvider, optVars ...map[string]interface{}) (result Holder) {
	var err error
	if result, err = v.Path(pathProvider, optVars...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

// MarshalJSON - реализация интерфейса bw.Val
func (v Holder) MarshalJSON() ([]byte, error) {
	if len(v.Pth) == 0 {
		return json.Marshal(v.Val)
	} else {
		result := map[string]interface{}{}
		result["val"] = v.Val
		result["path"] = v.Pth.String()
		return json.Marshal(result)
	}
}

// SetPathVal - реализация интерфейса bw.Val
func (v *Holder) SetPathVal(val interface{}, path bw.ValPath, optVars ...map[string]interface{}) (err error) {
	if len(path) == 0 {
		v.Val = val
		return
	}
	if path[len(path)-1].Type == bw.ValPathItemHash {
		return readonlyPathError(path)
	}

	var simplePath bw.ValPath
	simplePath, err = v.simplifyPath(path, optVars)
	if err != nil {
		return
	}

	result := v.Val
	if result == nil {
		return v.wrongValError()
	}

	if path[0].Type == bw.ValPathItemVar {
		var vars map[string]interface{}
		if len(optVars) > 0 {
			vars = optVars[0]
		}
		if vars == nil {
			return bwerr.From(ansiVarsIsNil)
		}
		simplePath[0].Type = bw.ValPathItemKey
		h := Holder{Val: vars}
		return h.SetPathVal(val, simplePath)
		// return FromVal(vars).SetPathVal(val, simplePath)
	}

	if len(simplePath) > 1 {
		for i, vpi := range simplePath[:len(simplePath)-1] {
			switch vpi.Type {
			case bw.ValPathItemKey:
				result, err = Holder{result, path[:i+1]}.KeyVal(vpi.Key)
			case bw.ValPathItemIdx:
				result, err = Holder{result, path[:i+1]}.IdxVal(vpi.Idx)
			}
			if err != nil {
				return
			} else if result == nil {
				return Holder{nil, path[:i+1]}.wrongValError()
			}
		}
	}
	rh := Holder{result, path[:len(path)-1]}
	vpi := simplePath[len(simplePath)-1]
	switch vpi.Type {
	case bw.ValPathItemKey:
		err = rh.SetKeyVal(val, vpi.Key)
	case bw.ValPathItemIdx:
		err = rh.SetIdxVal(val, vpi.Idx)
	}
	return
}

// func (v Holder) Bool() (result bool, err error) {
func (v Holder) Bool() (result bool) {
	switch t, kind := Kind(v.Val); kind {
	case ValInt:
		i, _ := t.(int)
		result = i != 0
	case ValFloat64:
		f, _ := t.(float64)
		result = f != 0
	case ValString:
		s, _ := t.(string)
		result = len(s) > 0
	case ValMap:
		m, _ := t.(map[string]interface{})
		result = len(m) > 0
	case ValArray:
		vals, _ := t.([]interface{})
		result = len(vals) > 0
	case ValArrayOfString:
		ss, _ := t.([]string)
		result = len(ss) > 0
	}
	return
}

// func (v Holder) MustBool() (result bool) {
// 	var err error
// 	if result, err = v.Bool(); err != nil {
// 		bwerr.PanicErr(err)
// 	}
// 	return
// }

// func (v Holder) String() (result string, err error) {
func (v Holder) String() (result string) {
	switch t, kind := Kind(v.Val); kind {
	case ValNil:
		result = ""
	case ValInt:
		i, _ := t.(int)
		result = strconv.Itoa(i)
	case ValString:
		result, _ = t.(string)
	default:
		bytes, _ := json.Marshal(t)
		result = string(bytes)
	}
	return
}

// func (v Holder) MustString() (result string) {
// 	var err error
// 	if result, err = v.String(); err != nil {
// 		bwerr.PanicErr(err)
// 	}
// 	return
// }

func (v Holder) Int() (result int, err error) {
	switch t, kind := Kind(v.Val); kind {
	case ValInt:
		result, _ = t.(int)
	case ValString:
		s, _ := t.(string)
		result, err = bwstr.ParseInt(s)
	default:
		err = v.notOfValKindError(ValKindSetFrom(ValInt, ValString))
	}
	return
}

func (v Holder) MustInt() (result int) {
	var err error
	if result, err = v.Int(); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) Float64() (result float64, err error) {
	switch t, kind := Kind(v.Val); kind {
	case ValInt:
		var i int
		i, _ = t.(int)
		result = float64(i)
	case ValFloat64:
		result, _ = t.(float64)
	case ValNil:
		result = 0
	case ValString:
		s, _ := t.(string)
		result, err = strconv.ParseFloat(s, 64)
	default:
		err = v.notOfValKindError(ValKindSetFrom(ValFloat64, ValInt, ValNil))
	}
	return
}

func (v Holder) MustFloat64() (result float64) {
	var err error
	if result, err = v.Float64(); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

// func (v Holder) Number() (result bwtype.RangeLimit, err error) {
// 	var ok bool
// 	if result, ok = Number(v.Val); !ok {
// 		err = v.notOfValKindError(ValKindSetFrom(ValNumber))
// 	}
// 	return
// }

// func (v Holder) MustNumber() (result bwtype.RangeLimit) {
// 	var err error
// 	if result, err = v.Number(); err != nil {
// 		bwerr.PanicErr(err)
// 	}
// 	return
// }

func (v Holder) Array() (result []interface{}, err error) {
	switch t, kind := Kind(v.Val); kind {
	case ValArray:
		result, _ = t.([]interface{})
	case ValArrayOfString:
		ss, _ := t.([]string)
		result = []interface{}{}
		for _, s := range ss {
			result = append(result, s)
		}
	case ValNil:
		result = []interface{}{}
	default:
		err = v.notOfValKindError(ValKindSetFrom(ValArray, ValArrayOfString, ValNil))
	}
	// var ok bool
	// if result, ok = Array(v.Val); !ok {
	// 	err = v.notOfValKindError(ValKindSetFrom(ValArray))
	// }
	return
}

func (v Holder) MustArray() (result []interface{}) {
	var err error
	if result, err = v.Array(); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) ArrayOfString() (result []string, err error) {
	switch t, kind := Kind(v.Val); kind {
	case ValArrayOfString:
		result, _ = t.([]string)
	case ValArray:
		vals, _ := t.([]interface{})
		result = []string{}
		for i, val := range vals {
			if s, ok := val.(string); ok {
				result = append(result, s)
			} else {
				h := Holder{val, v.Pth.AppendIdx(i)}
				err = h.notOfValKindError(ValKindSetFrom(ValString))
				return
			}
		}
	case ValNil:
		result = []string{}
	default:
		err = v.notOfValKindError(ValKindSetFrom(ValArray, ValArrayOfString, ValNil))
	}
	return
}

func (v Holder) MustArrayOfString() (result []string) {
	var err error
	if result, err = v.ArrayOfString(); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) Map() (result map[string]interface{}, err error) {
	switch t, kind := Kind(v.Val); kind {
	case ValMap:
		result, _ = t.(map[string]interface{})
	case ValNil:
		result = map[string]interface{}{}
	default:
		err = v.notOfValKindError(ValKindSetFrom(ValMap, ValNil))
	}
	return
}

func (v Holder) MustMap() (result map[string]interface{}) {
	var err error
	if result, err = v.Map(); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) Key(key string, optDefaultValProvider ...defaultValProvider) (result Holder, err error) {
	var val interface{}
	if val, err = v.KeyVal(key, optDefaultValProvider...); err == nil {
		result = Holder{val, v.Pth.AppendKey(key)}
	}
	return
}

func (v Holder) Idx(idx int, optDefaultValProvider ...defaultValProvider) (result Holder, err error) {
	var val interface{}
	if val, err = v.IdxVal(idx, optDefaultValProvider...); err == nil {
		result = Holder{val, v.Pth.AppendIdx(idx)}
	}
	return
}

func (v Holder) MustKey(key string, optDefaultValProvider ...defaultValProvider) (result Holder) {
	var err error
	if result, err = v.Key(key, optDefaultValProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) SetKeyVal(val interface{}, key string) (err error) {
	var m map[string]interface{}
	if m, err = v.Map(); err == nil {
		bwdebug.Print("m", m, "key", key)
		m[key] = val
	}
	return
}

func (v Holder) KeyVal(key string, optDefaultValProvider ...defaultValProvider) (result interface{}, err error) {
	if v.Val == nil {
		var ok bool
		if result, ok = defaultVal(optDefaultValProvider); !ok {
			err = v.wrongValError()
		}
		return
	}
	var val interface{}
	var kind ValKind
	if val, kind = Kind(v.Val); kind != ValMap {
		err = v.notOfValKindError(ValKindSetFrom(ValMap))
		return
	}
	m, _ := val.(map[string]interface{})
	var ok bool
	if result, ok = m[key]; !ok {
		if result, ok = defaultVal(optDefaultValProvider); !ok {
			err = v.hasNoKeyError(key)
		}
	}
	return
}

func (v Holder) IdxVal(idx int, optDefaultValProvider ...defaultValProvider) (result interface{}, err error) {
	if v.Val == nil {
		var ok bool
		if len(optDefaultValProvider) > 0 {
			result, ok = optDefaultValProvider[0]()
		}
		if !ok {
			err = v.wrongValError()
		}
		return
	}
	err = v.idxHelper(idx,
		func(vals []interface{}, nidx int, ok bool) (err error) {
			if ok {
				result = vals[nidx]
			} else {
				result, ok = defaultVal(optDefaultValProvider)
			}
			if !ok {
				err = v.notEnoughRangeError(len(vals), idx)
			}
			return
		},
		func(ss []string, nidx int, ok bool) (err error) {
			if ok {
				result = ss[nidx]
			} else {
				result, ok = defaultVal(optDefaultValProvider)
			}
			if !ok {
				err = v.notEnoughRangeError(len(ss), idx)
			}
			return
		},
	)
	return
}

func (v Holder) MustKeyVal(key string, optDefaultValProvider ...defaultValProvider) (result interface{}) {
	var err error
	if result, err = v.KeyVal(key, optDefaultValProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) SetIdxVal(val interface{}, idx int) (err error) {
	err = v.idxHelper(idx,
		func(vals []interface{}, nidx int, ok bool) (err error) {
			if !ok {
				err = v.notEnoughRangeError(len(vals), idx)
			} else {
				vals[nidx] = val
			}
			return
		},
		func(ss []string, nidx int, ok bool) (err error) {
			if !ok {
				err = v.notEnoughRangeError(len(ss), idx)
			} else if s, ok := val.(string); !ok {
				err = v.canNotSetNonStringError(idx, val)
			} else {
				ss[nidx] = s
			}
			return
		},
	)
	return
}

func (v Holder) ValidVal(def Def) (result interface{}, err error) {
	result, err = v.validVal(def)
	// bwdebug.Print("v:#v", v, "result:#v", result)
	if err != nil {
		err = bwerr.Refine(err, ansi.String("<ansiVal>%s<ansi>::{Error}"), bwjson.Pretty(v.Val))
	}
	return
}

func (v Holder) MustValidVal(def Def) (result interface{}) {
	var err error
	if result, err = v.ValidVal(def); err != nil {
		bwerr.PanicErr(err)
	}
	// bwdebug.Print("v:#v", v, "result:#v", result, "def:json", def)
	return
}

// ============================================================================

func (v *Holder) simplifyPath(path bw.ValPath, optVars []map[string]interface{}) (result bw.ValPath, err error) {
	result = bw.ValPath{}
	for _, vpi := range path {
		if vpi.Type != bw.ValPathItemPath {
			result = append(result, vpi)
		} else {
			var val interface{}
			if val, err = v.PathVal(vpi.Path, optVars...); err != nil {
				return
			}
			h := Holder{Val: val}
			if _, err = h.KindSwitch(map[ValKind]KindCase{
				ValString: func(val interface{}) (err error) {
					s, _ := val.(string)
					result = append(result, bw.ValPathItem{Type: bw.ValPathItemKey, Key: s})
					return
				},
				ValInt: func(val interface{}) (err error) {
					i, _ := val.(int)
					result = append(result, bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: i})
					return
				},
			}); err != nil {
				return
			}
			// switch _, kind := Kind(val); kind {
			// case ValString:
			// 	result = append(result, bw.ValPathItem{Type: bw.ValPathItemKey, Key: MustString(val)})
			// case ValInt:
			// 	result = append(result, bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: MustInt(val)})
			// default:
			// 	err = Holder{val, vpi.Path}.notOfValKindError(ValKindSetFrom(ValInt, ValString))
			// }
			// }
		}
	}
	return
}

func (v Holder) idxHelper(
	idx int,
	onArray func(vals []interface{}, nidx int, ok bool) error,
	onArrayOfString func(ss []string, nidx int, ok bool) error,
) (err error) {
	var nidx int
	var ok bool
	_, err = v.KindSwitch(map[ValKind]KindCase{
		ValArray: func(val interface{}) (err error) {
			vals, _ := val.([]interface{})
			nidx, ok = bw.NormalIdx(idx, len(vals))
			err = onArray(vals, nidx, ok)
			return
		},
		ValArrayOfString: func(val interface{}) (err error) {
			ss, _ := val.([]string)
			nidx, ok = bw.NormalIdx(idx, len(ss))
			err = onArrayOfString(ss, nidx, ok)
			return
		},
	})
	// switch val, kind := Kind(v.Val); kind {
	// case ValArray:
	// 	vals, _ := val.([]interface{})
	// 	nidx, ok = bw.NormalIdx(idx, len(vals))
	// 	err = onArray(vals, nidx, ok)
	// case ValArrayOfString:
	// 	ss, _ := val.([]string)
	// 	nidx, ok = bw.NormalIdx(idx, len(ss))
	// 	err = onArrayOfString(ss, nidx, ok)
	// default:
	// 	err = v.notOfValKindError(ValKindSetFrom(ValArray, ValArrayOfString))
	// }
	return
}

type KindCase func(val interface{}) (err error)

func (v Holder) KindSwitch(kindCases map[ValKind]KindCase, optDefaultCase ...KindCase) (val interface{}, err error) {
	val, kind := Kind(v.Val)
	if KindCase, ok := kindCases[kind]; ok {
		err = KindCase(val)
	} else if len(optDefaultCase) == 0 {
		vkSet := ValKindSet{}
		for vk := range kindCases {
			vkSet.Add(vk)
		}
		err = v.notOfValKindError(vkSet)
	} else if optDefaultCase[0] != nil {
		err = optDefaultCase[0](val)
	}
	return
}

type defaultValProvider func() (interface{}, bool)

func defaultVal(optDefaultValProvider []defaultValProvider) (result interface{}, ok bool) {
	if len(optDefaultValProvider) > 0 {
		if optDefaultValProvider[0] == nil {
			result = nil
			ok = true
		} else {
			result, ok = optDefaultValProvider[0]()
		}
	}
	return
}

// ============================================================================
