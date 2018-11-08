// Package bwval реализует интерфейc bw.Val и утилиты для работы с этим интерфейсом.
package bwval

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwval/path"
)

// ============================================================================

func PathFrom(s string) bw.ValPath {
	return path.MustParse(s)
}

// ============================================================================

// MustPathVal - must-обертка bw.Val.PathVal()
func MustPathVal(v bw.Val, path bw.ValPath, optVars ...map[string]interface{}) (result interface{}) {
	var err error
	if result, err = v.PathVal(path, optVars...); err != nil {
		bwerr.PanicA(bwerr.Err(bwerr.Refine(err,
			ansiMustPathValFailed,
			path, bwjson.Pretty(v), varsJSON(path, optVars),
		)))
		// bwerr.PanicA(bwerr.Err(err))
	}
	return result
}

// MustSetPathVal - must-обертка bw.Val.SetPathVal()
func MustSetPathVal(val interface{}, v bw.Val, path bw.ValPath, optVars ...map[string]interface{}) {
	var err error
	if err = v.SetPathVal(val, path, optVars...); err != nil {
		bwerr.PanicA(bwerr.Err(bwerr.Refine(err,
			ansiMustSetPathValFailed,
			path, bwjson.Pretty(v), varsJSON(path, optVars),
		)))
	}
}

// ============================================================================

// Bool - пытается извлечь bool из interface{}
func Bool(val interface{}) (result bool, ok bool) {
	if v, kind := Kind(val); kind == ValBool {
		result, ok = v.(bool)
	}
	return
}

// MustBool - must-обертка Bool()
func MustBool(val interface{}) (result bool) {
	var ok bool
	if result, ok = Bool(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Bool")
	}
	return
}

// Int - пытается извлечь int из interface{}
func Int(val interface{}) (result int, ok bool) {
	if v, kind := Kind(val); kind == ValInt {
		result, ok = v.(int)
	}
	return
}

// MustInt - must-обертка Int()
func MustInt(val interface{}) (result int) {
	var ok bool
	if result, ok = Int(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Int")
	}
	return
}

// String - пытается извлечь string из interface{}
func String(val interface{}) (result string, ok bool) {
	if v, kind := Kind(val); kind == ValString {
		result, ok = v.(string)
	}
	return
}

// MustString - must-обертка String()
func MustString(val interface{}) (result string) {
	var ok bool
	if result, ok = String(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "String")
	}
	return
}

// Map - пытается извлечь map[string]interface{} из interface{}
func Map(val interface{}) (result map[string]interface{}, ok bool) {
	if v, kind := Kind(val); kind == ValMap {
		result, ok = v.(map[string]interface{})
	}
	return
}

// MustMap - must-обертка Map()
func MustMap(val interface{}) (result map[string]interface{}) {
	var ok bool
	if result, ok = Map(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Map")
	}
	return
}

// Map - пытается извлечь []interface{} из interface{}
func Array(val interface{}) (result []interface{}, ok bool) {
	if v, kind := Kind(val); kind == ValArray {
		result, ok = v.([]interface{})
	}
	return
}

// MustArray - must-обертка Array()
func MustArray(val interface{}) (result []interface{}) {
	var ok bool
	if result, ok = Array(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Array")
	}
	return result
}

// ValKind - разновидность interface{}-значения
type ValKind uint8

// разновидности interface{}-значения
const (
	ValUnknown ValKind = iota
	ValNil
	ValBool
	ValInt
	ValString
	ValMap
	ValArray
)

// MarshalJSON encoding/json support
func (v ValKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// Kind - определяет разновидность  interface{}-значения
func Kind(val interface{}) (result interface{}, kind ValKind) {
	if val == nil {
		kind = ValNil
	} else {
		switch t := val.(type) {
		case bool:
			result = t
			kind = ValBool
		case int8:
			result = int(t)
			kind = ValInt
		case int16:
			result = int(t)
			kind = ValInt
		case int32:
			result = int(t)
			kind = ValInt
		case int64:
			if int64(bw.MinInt) <= t && t <= int64(bw.MaxInt) {
				result = int(t)
				kind = ValInt
			}
		case int:
			result = t
			kind = ValInt
		case uint8:
			result = int(t)
			kind = ValInt
		case uint16:
			result = int(t)
			kind = ValInt
		case uint32:
			result = int(t)
			kind = ValInt
		case uint64:
			if t <= uint64(bw.MaxInt) {
				result = int(t)
				kind = ValInt
			}
		case uint:
			if t <= uint(bw.MaxInt) {
				result = int(t)
				kind = ValInt
			}
		case string:
			result = t
			kind = ValString
		case map[string]interface{}:
			result = t
			kind = ValMap
		case []interface{}:
			result = t
			kind = ValArray
		}
	}
	return
}

// ============================================================================

// FromVal - конструктор bw.Val из interface{}-значения
func FromVal(val interface{}) (result bw.Val) {
	var ok bool
	if result, ok = val.(bw.Val); !ok {
		result = &valHolder{val}
	}
	return
}

// ============================================================================

// PathVal - реализация интерфейса bw.Val
func (v *valHolder) PathVal(path bw.ValPath, optVars ...map[string]interface{}) (result interface{}, err error) {
	if len(path) == 0 {
		result = v.val
		return
	}
	defer func() {
		if err != nil {
			result = nil
		}
	}()

	var simplePath bw.ValPath
	simplePath, err = simplifyPath(v, path, optVars)
	if err != nil {
		return
	}

	if path[0].Type == bw.ValPathItemVar {
		var target interface{}
		if len(optVars) > 0 {
			target = optVars[0][path[0].Key]
		}
		return FromVal(target).PathVal(simplePath[1:])
	}

	result = v.val
	for i, vpi := range simplePath {
		switch vpi.Type {
		case bw.ValPathItemKey:
			result, err = v.byKey(result, path, i, vpi.Key)
		case bw.ValPathItemIdx:
			result, err = v.byIdx(result, path, i, vpi.Idx)
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
					err = v.pathValIsNotOfType(path[:i], result, "Map", "Array")
				}
			}
		}
		if err != nil {
			return
		}
	}
	return
}

// MarshalJSON - реализация интерфейса bw.Val
func (v *valHolder) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.val)
}

// SetPathVal - реализация интерфейса bw.Val
func (v *valHolder) SetPathVal(val interface{}, path bw.ValPath, optVars ...map[string]interface{}) (err error) {
	// bwdebug.Print("path", path, "len(path)", len(path))
	if len(path) == 0 {
		v.val = val
		return
	}
	if path[len(path)-1].Type == bw.ValPathItemHash {
		return bwerr.From(ansisReadOnlyPath, path)
	}

	var simplePath bw.ValPath
	simplePath, err = simplifyPath(v, path, optVars)
	if err != nil {
		return
	}

	result := v.val
	if result == nil {
		return v.valAtPathIsNil(bw.ValPath{})
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
		return FromVal(vars).SetPathVal(val, simplePath)
	}

	if len(simplePath) > 1 {
		for i, vpi := range simplePath[:len(simplePath)-1] {
			switch vpi.Type {
			case bw.ValPathItemKey:
				result, err = v.byKey(result, path, i, vpi.Key)
			case bw.ValPathItemIdx:
				result, err = v.byIdx(result, path, i, vpi.Idx)
			}
			if err != nil {
				return
			} else if result == nil {
				return v.valAtPathIsNil(path[:i+1])
			}
		}
	}
	resultPath := path[:len(path)-1]
	setKeyElem := func(key string) (err error) {
		var m map[string]interface{}
		if m, err = v.getMap(result, resultPath); err == nil {
			m[key] = val
		}
		return
	}
	setIdxElem := func(idx int) (err error) {
		var vals []interface{}
		var gotIdx int
		if vals, gotIdx, err = v.getArray(idx, result, resultPath); err == nil {
			if gotIdx < 0 {
				err = bwerr.From(ansiValAtPathHasNotEnoughRange, resultPath, bwjson.Pretty(result), len(vals), idx)
			} else {
				vals[gotIdx] = val
			}
		}
		return
	}
	vpi := simplePath[len(simplePath)-1]
	switch vpi.Type {
	case bw.ValPathItemKey:
		err = setKeyElem(vpi.Key)
	case bw.ValPathItemIdx:
		err = setIdxElem(vpi.Idx)
	}
	return
}

// ============================================================================

func From(s string) bw.Val {
	return FromVal(val.MustParse(s))
}

// ============================================================================
