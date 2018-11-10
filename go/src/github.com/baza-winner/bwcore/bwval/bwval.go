// Package bwval реализует интерфейc bw.Val и утилиты для работы с этим интерфейсом.
package bwval

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/bwval/path"
	"github.com/baza-winner/bwcore/bwval/val"
	"github.com/baza-winner/bwcore/defvalid/deftype"
)

// ============================================================================

// PathFrom - конструктор-парсер bw.ValPath из строки
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

// Number - пытается извлечь float64 из interface{}
func Number(val interface{}) (result float64, ok bool) {
	switch v, kind := Kind(val); kind {
	case ValInt:
		var i int
		i, ok = v.(int)
		result = float64(i)
	case ValNumber:
		result, ok = v.(float64)
	}
	return
}

// MustNumber - must-обертка Number()
func MustNumber(val interface{}) (result float64) {
	var ok bool
	if result, ok = Number(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Number")
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
	ValNumber
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
		case float32:
			result = float64(t)
			kind = ValNumber
		case float64:
			result = t
			kind = ValNumber
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
func FromValPath(val interface{}, path bw.ValPath) *Holder {
	return &Holder{val, path}
}

// FromVal - конструктор bw.Val из interface{}-значения
func FromVal(val interface{}, optPathStr ...string) *Holder {
	path := bw.ValPath{}
	if len(optPathStr) > 0 {
		path = PathFrom(optPathStr[0])
	}
	return FromValPath(val, path)
}

// From - конструктор-парсер bw.Val из строки
func From(s string, optVars ...map[string]interface{}) *Holder {
	return FromValPath(val.MustParse(s, optVars...), bw.ValPath{})
}

// ============================================================================

type Holder struct {
	val  interface{}
	path bw.ValPath
}

// PathVal - реализация интерфейса bw.Val
func (v *Holder) PathVal(path bw.ValPath, optVars ...map[string]interface{}) (result interface{}, err error) {
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
	simplePath, err = v.simplifyPath(path, optVars)
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
			result, err = Holder{result, path[:i+1]}.KeyVal(vpi.Key)
		case bw.ValPathItemIdx:
			result, err = Holder{result, path[:i+1]}.IdxVal(vpi.Idx)
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
					err = Holder{result, path[:i]}.notOfTypeError("Map", "Array")
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
func (v *Holder) MarshalJSON() ([]byte, error) {
	if len(v.path) == 0 {
		return json.Marshal(v.val)
	} else {
		result := map[string]interface{}{}
		result["val"] = v.val
		result["path"] = v.path
		return json.Marshal(result)
	}
}

// SetPathVal - реализация интерфейса bw.Val
func (v *Holder) SetPathVal(val interface{}, path bw.ValPath, optVars ...map[string]interface{}) (err error) {
	if len(path) == 0 {
		v.val = val
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

	result := v.val
	if result == nil {
		return v.wrongValError()
		// return valAtPathIsNil(bw.ValPath{})
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
	resultPath := path[:len(path)-1]
	setKeyElem := func(key string) (err error) {
		var m map[string]interface{}
		vh := Holder{result, resultPath}
		if m, err = vh.Map(); err == nil {
			m[key] = val
		}
		return
	}
	setIdxElem := func(idx int) (err error) {
		var vals []interface{}
		var gotIdx int
		rh := Holder{result, resultPath}
		if vals, gotIdx, err = rh.arrayIdx(idx); err == nil {
			if gotIdx < 0 {
				err = Holder{result, resultPath}.notEnoughRangeError(len(vals), idx)
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

func (v Holder) Bool() (result bool, err error) {
	var ok bool
	if result, ok = Bool(v.val); !ok {
		err = v.notOfTypeError("Bool")
	}
	return
}

func (v Holder) String() (result string, err error) {
	var ok bool
	if result, ok = String(v.val); !ok {
		err = v.notOfTypeError("String")
	}
	return
}

func (v Holder) Int() (result int, err error) {
	var ok bool
	if result, ok = Int(v.val); !ok {
		err = v.notOfTypeError("Int")
	}
	return
}

func (v Holder) Number() (result float64, err error) {
	var ok bool
	if result, ok = Number(v.val); !ok {
		err = v.notOfTypeError("Number")
	}
	return
}

func (v Holder) Array() (result []interface{}, err error) {
	var ok bool
	if result, ok = Array(v.val); !ok {
		err = v.notOfTypeError("Array")
	}
	return
}

func (v Holder) ArrayOfString() (result []string, err error) {
	var vals interface{}
	if vals, err = v.Array(); err != nil {
		return
	}
	result = []string{}
	var s string
	for i := range MustArray(vals) {
		vp, _ := v.Idx(i)
		if s, err = vp.String(); err != nil {
			return
		}
		result = append(result, s)
	}
	return
}

func (v Holder) Map() (result map[string]interface{}, err error) {
	var ok bool
	if result, ok = Map(v.val); !ok {
		err = v.notOfTypeError("Map")
	}
	return
}

func (v Holder) Key(key string) (result Holder, err error) {
	var val interface{}
	if val, err = v.KeyVal(key); err == nil {
		result = Holder{val, v.path.AppendKey(key)}
	}
	return
}

func (v Holder) SetKey(val interface{}, key string) (err error) {
	vpi := bw.ValPathItem{Type: bw.ValPathItemKey, Key: key}
	err = FromVal(v.val).SetPathVal(val, bw.ValPath{vpi})
	return
}

func (v Holder) Idx(idx int) (result Holder, err error) {
	var val interface{}
	if val, err = v.IdxVal(idx); err == nil {
		result = Holder{val, v.path.AppendIdx(idx)}
	}
	return
}

// ============================================================================

type Def struct {
	Types      deftype.Set
	IsOptional bool
	Enum       bwset.String
	Range      Range
	Keys       map[string]Def
	Elem       *Def
	ArrayElem  *Def
	Default    interface{}
}

func (v Def) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{}
	result["Types"] = v.Types
	result["IsOptional"] = v.IsOptional
	if v.Enum != nil {
		result["Enum"] = v.Enum
	}
	if RangeKind(v.Range) != RangeNo {
		result["Range"] = v.Range
	}
	if v.Keys != nil {
		result["keys"] = v.Keys
	}
	if v.Elem != nil {
		result["Elem"] = *(v.Elem)
	}
	if v.ArrayElem != nil {
		result["ArrayElem"] = *(v.ArrayElem)
	}
	if v.Default != nil {
		result["Default"] = v.Default
	}
	return json.Marshal(result)
}

func DefFrom(def interface{}) (result Def) {
	var err error
	var compileDefResult *Def
	if compileDefResult, err = compileDef(Holder{
		def,
		bw.ValPath{bw.ValPathItem{
			Type: bw.ValPathItemVar, Key: "def",
		}},
	}); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	} else if compileDefResult == nil {
		bwerr.Panic("Unexpected behavior; def: %s", bwjson.Pretty(def))
	} else {
		result = *compileDefResult
	}
	return
}

// ============================================================================
