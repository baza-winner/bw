package bwval

import (
	"encoding/json"
	"fmt"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
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

func MustSetPathVal(val interface{}, v bw.Val, path bw.ValPath, optVars ...map[string]interface{}) {
	var err error
	if err = v.SetPathVal(val, path, optVars...); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
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
		result = &valHolder{val}
	}
	return
}

// ============================================================================

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
	byKey := func(val interface{}, i int, key string) (result interface{}, err error) {
		result = val
		if result == nil {
			return
		}
		if m, ok := result.(map[string]interface{}); !ok {
			err = bwerr.From(ansiValAtPathIsNotOfType, bwjson.Pretty(v.val), path[:i+1], bwjson.Pretty(result), varsJSON(path[:i+1], optVars), "Map")
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
			err = bwerr.From(ansiValAtPathIsNotOfType, bwjson.Pretty(v.val), path[:i+1], bwjson.Pretty(result), varsJSON(path[:i+1], optVars), "Array")
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
					err = v.valAtPathIsNotOfType(path[:i+1], result, optVars, "Map", "Array")
				}
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
	ansiType                      string
	ansiValAtPathIsNotOfTypes     string
	ansiValAtPathIsNorNeither     string
	ansiValAtPathIsNil            string
	ansiValAtPathIsReadOnly       string
	ansiValAtPathHasNoEnoughRange string
	ansiVars                      string
	ansiNoVar                     string
)

func init() {
	valPathPrefix := "<ansiVal>%s<ansi><ansiPath>.%s<ansi> "
	ansiValAtPathIsNotOfType = ansi.String(valPathPrefix + "(<ansiVal>%s<ansi>) %sis not <ansiType>%s")
	ansiValAtPathIsNotOfTypes = ansi.String(valPathPrefix + "(<ansiVal>%s<ansi>) %sis none of %s")
	ansiValAtPathIsNil = ansi.String(valPathPrefix + "%sis <ansiErr>nil")
	ansiValAtPathHasNoEnoughRange = ansi.String(valPathPrefix + "(<ansiVal>%s<ansi>) %shas no enough length (<ansiVal>%d<ansi>) for idx (<ansiVal>%d)")
	ansiNoVar = ansi.String(valPathPrefix + "<ansiVar>vars<ansi> (<ansiVal>%s<ansi>) has no var <ansiVar>%s")

	ansiType = ansi.String("<ansiType>%s")
	ansiValAtPathIsReadOnly = ansi.String("<ansiPath>.%s<ansi> is <ansiErr>readonly")
	ansiVars = ansi.String("(<ansiVar>vars<ansi>: <ansiVal>%s<ansi>) ")
}

func hasVar(path bw.ValPath) bool {
	for _, vpi := range path {
		switch vpi.Type {
		case bw.ValPathItemVar:
			return true
		case bw.ValPathItemPath:
			if hasVar(vpi.Path) {
				return true
			}
		}
	}
	return false
}

func varsJSON(path bw.ValPath, optVars []map[string]interface{}) (result string) {
	if hasVar(path) {
		var vars map[string]interface{}
		if len(optVars) > 0 {
			vars = optVars[0]
		}
		result = fmt.Sprintf(ansiVars, bwjson.Pretty(vars))
	}
	return
}

func simplifyPath(v *valHolder, path bw.ValPath, optVars []map[string]interface{}) (result bw.ValPath, err error) {
	result = bw.ValPath{}
	for _, vpi := range path {
		if vpi.Type != bw.ValPathItemPath {
			result = append(result, vpi)
		} else {
			var val interface{}
			val, err = v.PathVal(vpi.Path, optVars...)
			if err == nil {
				switch _, kind := Kind(val); kind {
				case ValString:
					result = append(result, bw.ValPathItem{Type: bw.ValPathItemKey, Key: MustString(val)})
				case ValInt:
					result = append(result, bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: MustInt(val)})
				default:
					err = v.valAtPathIsNotOfType(vpi.Path, val, optVars, "Int", "String")
					// err = bwerr.From(ansiValAtPathIsNorNeither, bwjson.Pretty(v.val), bw.ValPath(vpi.Path), bwjson.Pretty(val), varsJSON(vpi.Path, optVars), "Int", "String")
				}
			}
		}
	}
	return
}

var typeSeparator = map[bool]string{
	true:  " or ",
	false: ", ",
}

func (v *valHolder) valAtPathIsNotOfType(path bw.ValPath, val interface{}, optVars []map[string]interface{}, expectedType string, optExpectedType ...string) (result error) {
	if len(optExpectedType) == 0 {
		result = bwerr.From(ansiValAtPathIsNotOfType, bwjson.Pretty(v.val), path, bwjson.Pretty(val), varsJSON(path, optVars), expectedType)
	} else {
		expectedTypes := fmt.Sprintf(ansiType, expectedType)
		for i, elem := range optExpectedType {
			expectedTypes += typeSeparator[i == len(optExpectedType)-1] + fmt.Sprintf(ansiType, elem)
		}
		result = bwerr.From(ansiValAtPathIsNotOfTypes, bwjson.Pretty(v.val), path, bwjson.Pretty(val), varsJSON(path, optVars), expectedTypes)
	}
	return
}

func (v *valHolder) valAtPathIsNil(path bw.ValPath) error {
	return bwerr.From(ansiValAtPathIsNil, bwjson.Pretty(v.val), path)
}

func (v *valHolder) SetPathVal(val interface{}, path bw.ValPath, optVars ...map[string]interface{}) (err error) {
	// bwdebug.Print("path", path, "len(path)", len(path))
	if len(path) == 0 {
		v.val = val
		return
	}
	if path[len(path)-1].Type == bw.ValPathItemHash {
		return bwerr.From(ansiValAtPathIsReadOnly, path)
	}

	var simplePath bw.ValPath
	simplePath, err = simplifyPath(v, path, optVars)
	if err != nil {
		return
	}

	result := v.val
	if result == nil {
		return v.valAtPathIsNil(bw.ValPath{})
		// return bwerr.From(ansiValAtPathIsNil, v.val, path, varsJSON(path, optVars))
	}

	if path[0].Type == bw.ValPathItemVar {
		var foundVar bool
		var target interface{}
		var vars map[string]interface{}
		if len(optVars) > 0 {
			vars = optVars[0]
			target, foundVar = vars[path[0].Key]
		}
		if !foundVar {
			return bwerr.From(ansiNoVar, v.val, path, bwjson.Pretty(vars), path[0].Key)
		}
		return FromVal(target).SetPathVal(val, simplePath[1:])
	}

	byKey := func(val interface{}, i int, key string) (result interface{}, err error) {
		result = val
		if result == nil {
			return
		}
		if m, ok := result.(map[string]interface{}); !ok {
			err = v.valAtPathIsNotOfType(path[:i+1], result, optVars, "Map")
			// err = bwerr.From(ansiValAtPathIsNotOfType, bwjson.Pretty(v.val), path[:i+1], bwjson.Pretty(result), varsJSON(path[:i+1], optVars), "Map")
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
			err = v.valAtPathIsNotOfType(path[:i+1], result, optVars, "Array")
			// err = bwerr.From(ansiValAtPathIsNotOfType, bwjson.Pretty(v.val), bw.ValPath(path[:i+1]), bwjson.Pretty(result), varsJSON(path[:i+1], optVars), "Array")
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

	if len(simplePath) > 1 {
		for i, vpi := range simplePath[:len(simplePath)-1] {
			switch vpi.Type {
			case bw.ValPathItemKey:
				result, err = byKey(result, i, vpi.Key)
			case bw.ValPathItemIdx:
				result, err = byIdx(result, i, vpi.Idx)
			}
			if err == nil && result == nil {
				err = v.valAtPathIsNil(path[:i+1])
				// err = bwerr.From(ansiValAtPathIsNil, bwjson.Pretty(v.val), bw.ValPath(path[0:i+1]))
			}
			if err != nil {
				return
			}
		}
	}
	resultPath := path[:len(path)-1]
	setKeyElem := func(key string) (err error) {
		if m, ok := Map(result); !ok {
			err = v.valAtPathIsNotOfType(resultPath, result, optVars, "Map")
			// err = bwerr.From(ansiValAtPathIsNotOfType, bwjson.Pretty(v.val), resultPath, bwjson.Pretty(result), varsJSON(resultPath, optVars), "Map")
		} else {
			m[key] = val
		}
		return
	}
	setIdxElem := func(idx int) (err error) {
		if vals, ok := Array(result); !ok {
			err = v.valAtPathIsNotOfType(resultPath, result, optVars, "Array")
			// err = bwerr.From(ansiValAtPathIsNotOfType, bwjson.Pretty(v.val), resultPath, bwjson.Pretty(result), varsJSON(resultPath, optVars), "Array")
		} else {
			l := len(vals)
			minIdx := -l
			maxIdx := l - 1
			if !(minIdx <= idx && idx <= maxIdx) {
				err = bwerr.From(ansiValAtPathHasNoEnoughRange, bwjson.Pretty(v.val), resultPath, bwjson.Pretty(result), l, idx)
			} else {
				if idx < 0 {
					idx = l + idx
				}
				vals[idx] = val
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
