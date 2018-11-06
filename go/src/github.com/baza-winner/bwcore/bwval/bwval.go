package bwval

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwdebug"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwval/path"
)

// ============================================================================

//go:generate stringer -type=ValKind

// ============================================================================

func PathFrom(s string) (result bw.ValPath, err error) {
	return path.Parse(s)
}

func MustPathFrom(s string) bw.ValPath {
	result, err := PathFrom(s)
	if err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return result
}

// ============================================================================

func MustPathVal(v bw.Val, path bw.ValPath, vars map[string]interface{}) (result interface{}) {
	var err error
	if result, err = v.PathVal(path, vars); err != nil {
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

func Bool(val interface{}) (result bool, err error) {
	var ok bool
	if result, ok = val.(bool); !ok {
		err = bwerr.FromA(bwerr.A{1, ansiIsNotOfType, bw.Args(val, "Bool")})
	}
	return
}

func MustBool(val interface{}) (result bool) {
	var err error
	if result, err = Bool(val); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return result
}

var (
	ansiIsOutOfRange string
)

func init() {
	ansiIsOutOfRange = ansi.String("<ansiVal>%#v<ansi> is out of range <ansiVal>%d..%d")
}

func Int(val interface{}) (result int, err error) {
	switch t := val.(type) {
	case int8:
		result = int(t)
	case int16:
		result = int(t)
	case int32:
		result = int(t)
	case int64:
		if int64(bw.MinInt) <= t && t <= int64(bw.MaxInt) {
			result = int(t)
		} else {
			err = bwerr.FromA(bwerr.A{1, ansiIsOutOfRange, bw.Args(val, bw.MinInt, bw.MaxInt)})
		}
	case int:
		result = t
	case uint8:
		result = int(t)
	case uint16:
		result = int(t)
	case uint32:
		result = int(t)
	case uint64:
		if t <= uint64(bw.MaxInt) {
			result = int(t)
		} else {
			err = bwerr.FromA(bwerr.A{1, ansiIsOutOfRange, bw.Args(val, bw.MinInt, bw.MaxInt)})
		}
	case uint:
		if t <= uint(bw.MaxInt) {
			result = int(t)
		} else {
			err = bwerr.FromA(bwerr.A{1, ansiIsOutOfRange, bw.Args(val, bw.MinInt, bw.MaxInt)})
		}
	default:
		err = bwerr.FromA(bwerr.A{1, ansiIsNotOfType, bw.Args(val, "Int")})
	}
	return
}

func MustInt(val interface{}) (result int) {
	var err error
	if result, err = Int(val); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return result
}

func String(val interface{}) (result string, err error) {
	var ok bool
	if result, ok = val.(string); !ok {
		err = bwerr.FromA(bwerr.A{1, ansiIsNotOfType, bw.Args(val, "String")})
	}
	return
}

func MustString(val interface{}) (result string) {
	var err error
	if result, err = String(val); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return result
}

func Map(val interface{}) (result map[string]interface{}, err error) {
	var ok bool
	if result, ok = val.(map[string]interface{}); !ok {
		err = bwerr.FromA(bwerr.A{1, ansiIsNotOfType, bw.Args(val, "Map")})
	}
	return
}

func MustMap(val interface{}) (result map[string]interface{}) {
	var err error
	if result, err = Map(val); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return result
}

func Array(val interface{}) (result []interface{}, err error) {
	var ok bool
	if result, ok = val.([]interface{}); !ok {
		err = bwerr.FromA(bwerr.A{1, ansiIsNotOfType, bw.Args(val, "Array")})
	}
	return
}

func MustArray(val interface{}) (result []interface{}) {
	var err error
	if result, err = Array(val); err != nil {
		bwerr.PanicA(bwerr.Err(err))
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

func Kind(val interface{}) (result ValKind) {
	if val == nil {
		result = ValNil
	} else {
		switch t := val.(type) {
		case bool:
			result = ValBool
		case int, int8, int16, int32, uint8, uint16, uint32:
			result = ValInt
		case int64:
			if int64(bw.MinInt) <= t && t <= int64(bw.MaxInt) {
				result = ValInt
			}
		case uint64:
			if t <= uint64(bw.MaxInt) {
				result = ValInt
			}
		case uint:
			if t <= uint(bw.MaxInt) {
				result = ValInt
			}
		case string:
			result = ValString
		case map[string]interface{}:
			result = ValMap
		case []interface{}:
			result = ValArray
		}
	}
	return
}

// ============================================================================

type valHolder struct {
	val interface{}
}

func From(val interface{}) (result bw.Val) {
	var ok bool
	if result, ok = val.(bw.Val); !ok {
		result = valHolder{val}
	}
	return
}

// ============================================================================

var (
	ansiValAtPathIsNotOfType  string
	ansiValAtPathIsNorNeither string
)

func init() {
	ansiValAtPathIsNotOfType = ansi.String("<ansiVal>%#v<ansi>::<ansiPath>%s<ansi> (<ansiVal>%#v<ansi>) is not <ansiType>%s")
	ansiValAtPathIsNorNeither = ansi.String("<ansiVal>%#v<ansi>::<ansiPath>.%s<ansi> (<ansiVal>%#v<ansi>) is nor <ansiType>%s<ansi>, neither <ansiType>%s")
}

func (v valHolder) PathVal(path bw.ValPath, vars map[string]interface{}) (result interface{}, err error) {
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
					// bwdebug.Print("!HERE", "result", result, "t", t)
				case []interface{}:
					result = len(t)
				default:
					err = bwerr.From(ansiValAtPathIsNorNeither, v.val, bw.ValPath(path[:i+1]), result, "Map", "Array")
				}
			}
		case bw.ValPathItemPath:
			var val interface{}
			val, err = v.PathVal(vpi.Path, vars)
			if err == nil {
				switch Kind(val) {
				case ValString:
					result, err = byKey(result, i, MustString(val))
				case ValInt:
					result, err = byIdx(result, i, MustInt(val))
				default:
					err = bwerr.From(ansiValAtPathIsNorNeither, v.val, bw.ValPath(path[:i+1]), val, "Int", "String")
				}
			}
		case bw.ValPathItemVar:
			if vars == nil {
				result = nil
			} else {
				bwdebug.Print("vpi", vpi)
				result = vars[vpi.Key]
			}
		}
		if err != nil {
			return
		}
	}
	return
}

func (v valHolder) SetValToPath(val []interface{}, path bw.ValPath, vars map[string]interface{}) (err error) {
	return
}

// ============================================================================
