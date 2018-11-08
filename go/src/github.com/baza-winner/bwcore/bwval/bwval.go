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
		bwerr.PanicA(bwerr.Err(bwerr.Refine(err,
			ansiMustPathValFailed,
			path, bwjson.Pretty(v), varsJSON(path, optVars),
		)))
		// bwerr.PanicA(bwerr.Err(err))
	}
	return result
}

func MustSetPathVal(val interface{}, v bw.Val, path bw.ValPath, optVars ...map[string]interface{}) {
	var err error
	if err = v.SetPathVal(val, path, optVars...); err != nil {
		bwerr.PanicA(bwerr.Err(bwerr.Refine(err,
			ansiMustSetPathValFailed,
			path, bwjson.Pretty(v), varsJSON(path, optVars),
		)))
	}
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

var (
	ansiValAtPathIsNotOfType       string
	ansiMustPathValFailed          string
	ansiType                       string
	ansiValAtPathIsNotOfTypes      string
	ansiValAtPathIsNorNeither      string
	ansiValAtPathIsNil             string
	ansisReadOnlyPath              string
	ansiValAtPathHasNotEnoughRange string
	ansiVars                       string
	ansiVarsIsNil                  string
	ansiMustSetPathValFailed       string
)

func init() {
	ansiMustSetPathValFailed = ansi.String("Failed to set <ansiPath>%s<ansi> of <ansiVal>%s<ansi>%s: {Error}")
	ansiMustPathValFailed = ansi.String("Failed to get <ansiPath>%s<ansi> of <ansiVal>%s<ansi>%s: {Error}")
	valPathPrefix := "<ansiPath>%s<ansi> "
	ansiValAtPathIsNotOfType = ansi.String(valPathPrefix + "(<ansiVal>%s<ansi>) is not <ansiType>%s")
	ansiValAtPathIsNotOfTypes = ansi.String(valPathPrefix + "(<ansiVal>%s<ansi>) is none of %s")
	ansiValAtPathIsNil = ansi.String(valPathPrefix + "is <ansiErr>nil")
	ansiValAtPathHasNotEnoughRange = ansi.String(valPathPrefix + "(<ansiVal>%s<ansi>) has not enough length (<ansiVal>%d<ansi>) for idx (<ansiVal>%d)")

	ansiType = ansi.String("<ansiType>%s")
	ansisReadOnlyPath = ansi.String("<ansiPath>%s<ansi> is <ansiErr>readonly path")
	ansiVars = ansi.String(" with <ansiVar>vars<ansi> <ansiVal>%s<ansi>")
	ansiVarsIsNil = ansi.String("<ansiVar>vars<ansi> is <ansiErr>nil")
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
					// bwdebug.Print("vpi.Path", vpi.Path)
					err = v.pathValIsNotOfType(vpi.Path, val, "Int", "String")
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

func (v *valHolder) pathValIsNotOfType(path bw.ValPath, val interface{}, expectedType string, optExpectedType ...string) (result error) {
	if len(optExpectedType) == 0 {
		result = bwerr.From(ansiValAtPathIsNotOfType, path, bwjson.Pretty(val), expectedType)
	} else {
		expectedTypes := fmt.Sprintf(ansiType, expectedType)
		for i, elem := range optExpectedType {
			expectedTypes += typeSeparator[i == len(optExpectedType)-1] + fmt.Sprintf(ansiType, elem)
		}
		result = bwerr.From(ansiValAtPathIsNotOfTypes, path, bwjson.Pretty(val), expectedTypes)
		// bwdebug.Print("path", path, "ansiValAtPathIsNotOfTypes", ansiValAtPathIsNotOfTypes, "result", result)
	}
	return
}

func (v *valHolder) valAtPathIsNil(path bw.ValPath) error {
	return bwerr.From(ansiValAtPathIsNil, path)
}

func (v *valHolder) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.val)
}

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

func (v *valHolder) getArray(idx int, result interface{}, resultPath bw.ValPath) ([]interface{}, int, error) {
	var err error
	var ok bool
	var vals []interface{}
	if vals, ok = Array(result); !ok {
		err = v.pathValIsNotOfType(resultPath, result, "Array")
	} else {
		l := len(vals)
		minIdx := -l
		maxIdx := l - 1
		if !(minIdx <= idx && idx <= maxIdx) {
			idx = -1
		} else if idx < 0 {
			idx = l + idx
		}
	}
	return vals, idx, err
}

func (v *valHolder) getMap(result interface{}, resultPath bw.ValPath) (map[string]interface{}, error) {
	var err error
	var ok bool
	var m map[string]interface{}
	if m, ok = Map(result); !ok {
		err = v.pathValIsNotOfType(resultPath, result, "Map")
	}
	return m, err
}

func (v *valHolder) byKey(val interface{}, path bw.ValPath, i int, key string) (result interface{}, err error) {
	result = val
	if result == nil {
		return
	}
	var m map[string]interface{}
	if m, err = v.getMap(result, path[:i+1]); err == nil {
		result = m[key]
	}
	return
}

func (v *valHolder) byIdx(val interface{}, path bw.ValPath, i int, idx int) (result interface{}, err error) {
	result = val
	if result == nil {
		return
	}
	var vals []interface{}
	if vals, idx, err = v.getArray(idx, result, path[:i+1]); err == nil {
		if idx < 0 {
			result = nil
		} else {
			result = vals[idx]
		}
	}
	return
}

// ============================================================================
