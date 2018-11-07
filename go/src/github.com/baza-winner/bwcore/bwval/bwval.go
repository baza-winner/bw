package bwval

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwval/path"
)

// ============================================================================

//go:generate stringer -type=ValKind

// ============================================================================

func PathFrom(s string) bw.ValPath {
	return path.MustParse(s)
}

// ============================================================================

func MustPathVal(v bw.Val, path bw.ValPath, optVars ...map[string]interface{}) (result interface{}) {
	var err error
	if result, err = v.PathVal(path, optVars...); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return result
}

var (
	ansiIsNotOfType string
)

func init() {
	ansiIsNotOfType = ansi.String("<ansiVal>%#v<ansi> is not <ansiType>%s")
}

func Bool(val interface{}) (result bool, ok bool) {
	if v, kind := Kind(val); kind == ValBool {
		result, ok = v.(bool)
	}
	return
}

func MustBool(val interface{}) (result bool) {
	var ok bool
	if result, ok = Bool(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Bool")
	}
	return
}

func Int(val interface{}) (result int, ok bool) {
	if v, kind := Kind(val); kind == ValInt {
		result, ok = v.(int)
	}
	return
}

func MustInt(val interface{}) (result int) {
	var ok bool
	if result, ok = Int(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Int")
	}
	return
}

func String(val interface{}) (result string, ok bool) {
	if v, kind := Kind(val); kind == ValString {
		result, ok = v.(string)
	}
	return
}

func MustString(val interface{}) (result string) {
	var ok bool
	if result, ok = String(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "String")
	}
	return
}

func Map(val interface{}) (result map[string]interface{}, ok bool) {
	if v, kind := Kind(val); kind == ValMap {
		result, ok = v.(map[string]interface{})
	}
	return
}

func MustMap(val interface{}) (result map[string]interface{}) {
	var ok bool
	if result, ok = Map(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Map")
	}
	return
}

func Array(val interface{}) (result []interface{}, ok bool) {
	if v, kind := Kind(val); kind == ValArray {
		result, ok = v.([]interface{})
	}
	return
}

func MustArray(val interface{}) (result []interface{}) {
	var ok bool
	if result, ok = Array(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Array")
	}
	return result
}

type ValKind uint8

const (
	ValUnknown ValKind = iota
	ValNil
	ValBool
	ValInt
	ValString
	ValMap
	ValArray
)

func (v ValKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

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

type valHolder struct {
	val interface{}
}

func FromVal(val interface{}) (result bw.Val) {
	var ok bool
	if result, ok = val.(bw.Val); !ok {
		result = valHolder{val}
	}
	return
}

// ============================================================================

func (v valHolder) PathVal(path bw.ValPath, optVars ...map[string]interface{}) (result interface{}, err error) {
	defer func() {
		if err != nil {
			result = nil
		}
	}()
	byKey := func(val interface{}, i int, key string) (result interface{}, err error) {
		result = val
		if result == nil {
			return
		}
		if m, ok := result.(map[string]interface{}); !ok {
			err = bwerr.From(ansiValAtPathIsNotOfType, v.val, path[:i+1], result, "Map")
		} else {
			result = m[key]
		}
		return
	}
	byIdx := func(val interface{}, i int, idx int) (result interface{}, err error) {
		result = val
		if result == nil {
			return
		}
		if vals, ok := result.([]interface{}); !ok {
			err = bwerr.From(ansiValAtPathIsNotOfType, v.val, bw.ValPath(path[:i+1]), result, "Array")
		} else {
			l := len(vals)
			minIdx := -l
			maxIdx := l - 1
			if minIdx <= idx && idx <= maxIdx {
				if idx < 0 {
					idx = l + idx
				}
				result = vals[idx]
			} else {
				result = nil
			}
		}
		return
	}
	result = v.val
	for i, vpi := range path {
		switch vpi.Type {
		case bw.ValPathItemKey:
			result, err = byKey(result, i, vpi.Key)
		case bw.ValPathItemIdx:
			result, err = byIdx(result, i, vpi.Idx)
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
					err = bwerr.From(ansiValAtPathIsNorNeither, v.val, bw.ValPath(path[:i+1]), result, "Map", "Array")
				}
			}
		case bw.ValPathItemPath:
			var val interface{}
			val, err = v.PathVal(vpi.Path, optVars...)
			if err == nil {
				switch _, kind := Kind(val); kind {
				case ValString:
					result, err = byKey(result, i, MustString(val))
				case ValInt:
					result, err = byIdx(result, i, MustInt(val))
				default:
					err = bwerr.From(ansiValAtPathIsNorNeither, v.val, vpi.Path, val, "Int", "String")
				}
			}
		case bw.ValPathItemVar:
			if len(optVars) == 0 {
				result = nil
			} else {
				result = optVars[0][vpi.Key]
			}
		}
		if err != nil {
			return
		}
	}
	return
}

var (
	ansiValAtPathIsNotOfType      string
	ansiValAtPathIsNorNeither     string
	ansiValAtPathIsNil            string
	ansiValAtPathIsReadOnly       string
	ansiValAtPathHasNoEnoughRange string
)

func init() {
	valPathPrefix := "<ansiVal>%#v<ansi><ansiPath>.%s<ansi> "
	ansiValAtPathIsNotOfType = ansi.String(valPathPrefix + "(<ansiVal>%#v<ansi>) is not <ansiType>%s")
	ansiValAtPathIsNorNeither = ansi.String(valPathPrefix + "(<ansiVal>%#v<ansi>) is nor <ansiType>%s<ansi>, neither <ansiType>%s")
	ansiValAtPathIsNil = ansi.String(valPathPrefix + "is <ansiErr>nil")
	ansiValAtPathIsReadOnly = ansi.String(valPathPrefix + "is <ansiErr>readonly")
	ansiValAtPathHasNoEnoughRange = ansi.String(valPathPrefix + "(<ansiVal>%#v<ansi>) has no enough length (<ansiVal>%d<ansi>) for idx (<ansiVal>%d)")
}

func (v valHolder) SetValToPath(val []interface{}, path bw.ValPath, optVars ...map[string]interface{}) (err error) {
	if len(path) == 0 {
		v.val = val
		return
	}
	if path[len(path)-1].Type == bw.ValPathItemHash {
		err = bwerr.From(ansiValAtPathIsReadOnly, v.val, path)
		return
	}

	byKey := func(val interface{}, i int, key string) (result interface{}, err error) {
		result = val
		if result == nil {
			return
		}
		if m, ok := result.(map[string]interface{}); !ok {
			err = bwerr.From(ansiValAtPathIsNotOfType, v.val, path[:i+1], result, "Map")
		} else {
			result = m[key]
		}
		return
	}
	byIdx := func(val interface{}, i int, idx int) (result interface{}, err error) {
		result = val
		if result == nil {
			return
		}
		if vals, ok := result.([]interface{}); !ok {
			err = bwerr.From(ansiValAtPathIsNotOfType, v.val, bw.ValPath(path[:i+1]), result, "Array")
		} else {
			l := len(vals)
			minIdx := -l
			maxIdx := l - 1
			if minIdx <= idx && idx <= maxIdx {
				if idx < 0 {
					idx = l + idx
				}
				result = vals[idx]
			} else {
				result = nil
			}
		}
		return
	}

	result := v.val
	if result == nil {
		err = bwerr.From(ansiValAtPathIsNil, v.val, bw.ValPath{})
		return
	}
	if len(path) > 1 {
		for i, vpi := range path[:len(path)-1] {
			switch vpi.Type {
			case bw.ValPathItemKey:
				result, err = byKey(result, i, vpi.Key)
			case bw.ValPathItemIdx:
				result, err = byIdx(result, i, vpi.Idx)
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
						err = bwerr.From(ansiValAtPathIsNorNeither, v.val, bw.ValPath(path[:i+1]), result, "Map", "Array")
					}
				}
			case bw.ValPathItemPath:
				var val interface{}
				val, err = v.PathVal(vpi.Path, optVars...)
				if err == nil {
					switch _, kind := Kind(val); kind {
					case ValString:
						result, err = byKey(result, i, MustString(val))
					case ValInt:
						result, err = byIdx(result, i, MustInt(val))
					default:
						err = bwerr.From(ansiValAtPathIsNorNeither, v.val, bw.ValPath(path[:i+1]), val, "Int", "String")
					}
				}
			case bw.ValPathItemVar:
				if len(optVars) == 0 {
					result = nil
				} else {
					result = optVars[0][vpi.Key]
				}
			}
			if err == nil && result == nil {
				err = bwerr.From(ansiValAtPathIsNil, v.val, bw.ValPath(path[0:i+1]))
			}
			if err != nil {
				return
			}
		}
	}
	resultPath := path[:len(path)-1]
	setKeyElem := func(key string) (err error) {
		if m, ok := Map(result); !ok {
			err = bwerr.From(ansiValAtPathIsNotOfType, v.val, resultPath, result, "Map")
		} else {
			m[key] = val
		}
		return
	}
	setIdxElem := func(idx int) (err error) {
		if vals, ok := Array(result); !ok {
			err = bwerr.From(ansiValAtPathIsNotOfType, v.val, resultPath, result, "Array")
		} else {
			l := len(vals)
			minIdx := -l
			maxIdx := l - 1
			if !(minIdx <= idx && idx <= maxIdx) {
				err = bwerr.From(ansiValAtPathHasNoEnoughRange, v.val, resultPath, result, l, idx)
			} else {
				if idx < 0 {
					idx = l + idx
				}
				vals[idx] = val
			}
		}
		return
	}
	vpi := path[len(path)-1]
	switch vpi.Type {
	case bw.ValPathItemKey:
		err = setKeyElem(vpi.Key)
	case bw.ValPathItemIdx:
		err = setIdxElem(vpi.Idx)
	case bw.ValPathItemPath:
		var val interface{}
		val, err = v.PathVal(vpi.Path, optVars...)
		if err == nil {
			switch _, kind := Kind(val); kind {
			case ValString:
				err = setKeyElem(MustString(val))
			case ValInt:
				err = setIdxElem(MustInt(val))
			default:
				err = bwerr.From(ansiValAtPathIsNorNeither, v.val, vpi.Path, val, "Int", "String")
			}
		}
	}
	return
}

// ============================================================================
