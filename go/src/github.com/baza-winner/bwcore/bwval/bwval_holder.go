package bwval

import (
	"encoding/json"
	"fmt"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
)

// ============================================================================

type Holder struct {
	Val  interface{}
	Path bw.ValPath
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
		var target interface{}
		if len(optVars) > 0 {
			target = optVars[0][path[0].Key]
		}
		h := Holder{Val: target}
		return h.PathVal(simplePath[1:])
		// return FromVal(target).PathVal(simplePath[1:])
	}

	result = v.Val
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
func (v Holder) MarshalJSON() ([]byte, error) {
	if len(v.Path) == 0 {
		return json.Marshal(v.Val)
	} else {
		result := map[string]interface{}{}
		result["val"] = v.Val
		result["path"] = v.Path.String()
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

func (v Holder) Bool() (result bool, err error) {
	var ok bool
	if result, ok = Bool(v.Val); !ok {
		err = v.notOfTypeError("Bool")
	}
	return
}

func (v Holder) MustBool() (result bool) {
	var err error
	if result, err = v.Bool(); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return
}

func (v Holder) String() (result string, err error) {
	var ok bool
	if result, ok = String(v.Val); !ok {
		err = v.notOfTypeError("String")
	}
	return
}

func (v Holder) MustString() (result string) {
	var err error
	if result, err = v.String(); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return
}

func (v Holder) Int() (result int, err error) {
	var ok bool
	if result, ok = Int(v.Val); !ok {
		err = v.notOfTypeError("Int")
	}
	return
}

func (v Holder) MustInt() (result int) {
	var err error
	if result, err = v.Int(); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return
}

func (v Holder) Number() (result float64, err error) {
	var ok bool
	if result, ok = Number(v.Val); !ok {
		err = v.notOfTypeError("Number")
	}
	return
}

func (v Holder) MustNumber() (result float64) {
	var err error
	if result, err = v.Number(); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return
}

func (v Holder) Array() (result []interface{}, err error) {
	var ok bool
	if result, ok = Array(v.Val); !ok {
		err = v.notOfTypeError("Array")
	}
	return
}

func (v Holder) MustArray() (result []interface{}) {
	var err error
	if result, err = v.Array(); err != nil {
		bwerr.PanicA(bwerr.Err(err))
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

func (v Holder) MustArrayOfString() (result []string) {
	var err error
	if result, err = v.ArrayOfString(); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return
}

func (v Holder) Map() (result map[string]interface{}, err error) {
	var ok bool
	if result, ok = Map(v.Val); !ok {
		err = v.notOfTypeError("Map")
	}
	return
}

func (v Holder) MustMap() (result map[string]interface{}) {
	var err error
	if result, err = v.Map(); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return
}

func (v Holder) Key(key string) (result Holder, err error) {
	var val interface{}
	if val, err = v.KeyVal(key); err == nil {
		result = Holder{val, v.Path.AppendKey(key)}
	}
	return
}

func (v Holder) SetKeyVal(val interface{}, key string) (err error) {
	var m map[string]interface{}
	if m, err = v.Map(); err == nil {
		m[key] = val
	}
	return
}

func (v Holder) SetIdxVal(val interface{}, idx int) (err error) {
	var vals []interface{}
	var gotIdx int
	if vals, gotIdx, err = v.arrayIdx(idx); err == nil {
		if gotIdx < 0 {
			err = v.notEnoughRangeError(len(vals), idx)
		} else {
			vals[gotIdx] = val
		}
	}
	return
}

func (v Holder) Idx(idx int) (result Holder, err error) {
	var val interface{}
	if val, err = v.IdxVal(idx); err == nil {
		result = Holder{val, v.Path.AppendIdx(idx)}
	}
	return
}

func (v Holder) KeyVal(key string) (result interface{}, err error) {
	if v.Val == nil {
		return
	}
	var m map[string]interface{}
	if m, err = v.Map(); err == nil {
		result = m[key]
	}
	return
}

func (v Holder) IdxVal(idx int) (result interface{}, err error) {
	if v.Val == nil {
		return
	}
	var vals []interface{}
	if vals, idx, err = v.arrayIdx(idx); err == nil {
		if idx < 0 {
			result = nil
		} else {
			result = vals[idx]
		}
	}
	return
}

func (v Holder) ValidVal(def Def) (result interface{}, err error) {
	return v.validVal(def)
}

func (v Holder) MustValidVal(def Def) (result interface{}) {
	var err error
	if result, err = v.ValidVal(def); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
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
			val, err = v.PathVal(vpi.Path, optVars...)
			if err == nil {
				switch _, kind := Kind(val); kind {
				case ValString:
					result = append(result, bw.ValPathItem{Type: bw.ValPathItemKey, Key: MustString(val)})
				case ValInt:
					result = append(result, bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: MustInt(val)})
				default:
					err = Holder{val, vpi.Path}.notOfTypeError("Int", "String")
				}
			}
		}
	}
	return
}

func (v Holder) arrayIdx(idx int) ([]interface{}, int, error) {
	var err error
	var ok bool
	var vals []interface{}
	if vals, ok = Array(v.Val); !ok {
		err = v.notOfTypeError("Array")
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

func (v Holder) ansiString() (result string) {
	return fmt.Sprintf(ansi.String("<ansiPath>%s<ansi> (<ansiVal>%s<ansi>)"), v.Path, bwjson.Pretty(v.Val))
}

// func (v Holder) inRange(rng Range) (result bool) {
// 	rangeKind := RangeKind(rng)
// 	if rangeKind == RangeNo {
// 		result = true
// 	} else if _, vk := Kind(v.Val); vk == ValNumber || rng.ValKind() == ValNumber {
// 		if n, ok := Number(v.Val); ok {
// 			switch rangeKind {
// 			case RangeMinMax:
// 				result = MustNumber(rng.Min()) <= n && n <= MustNumber(rng.Max())
// 			case RangeMin:
// 				result = MustNumber(rng.Min()) <= n
// 			case RangeMax:
// 				result = n <= MustNumber(rng.Max())
// 			}
// 		}
// 	} else {
// 		if n, ok := Int(v.Val); ok {
// 			switch rangeKind {
// 			case RangeMinMax:
// 				result = MustInt(rng.Min()) <= n && n <= MustInt(rng.Max())
// 			case RangeMin:
// 				result = MustInt(rng.Min()) <= n
// 			case RangeMax:
// 				result = n <= MustInt(rng.Max())
// 			}
// 		}
// 	}
// 	return
// }

// ============================================================================
