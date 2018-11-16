package bwval

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	bwjson "github.com/baza-winner/bwcore/bwjson"
)

// ============================================================================

type Holder struct {
	Val interface{}
	Pth bw.ValPath
}

func HolderFrom(s string, optVars ...map[string]interface{}) Holder {
	return Holder{Val: From(s, optVars...)}
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
	}

	result = v.Val
	for i, vpi := range simplePath {
		switch vpi.Type {
		case bw.ValPathItemKey:
			result, err = Holder{result, path[:i]}.KeyVal(vpi.Key,
				func() (result interface{}, ok bool) {
					for _, vpi := range path[i:] {
						if vpi.IsOptional {
							ok = true
							break
						}
					}
					return
				},
			)
		case bw.ValPathItemIdx:
			// bwdebug.Print("path[:i]:json", path[:i], "path[]:json", path, "path[i:]:json", path[i:])
			result, err = Holder{result, path[:i]}.IdxVal(vpi.Idx,
				func() (result interface{}, ok bool) {
					for _, vpi := range path[i:] {
						if vpi.IsOptional {
							ok = true
							break
						}
					}
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

func (v Holder) Path(path bw.ValPath, optVars ...map[string]interface{}) (result Holder, err error) {
	var val interface{}
	if val, err = (&v).PathVal(path, optVars...); err != nil {
		return
	}
	result = Holder{val, path}
	return
}

func (v Holder) MustPath(path bw.ValPath, optVars ...map[string]interface{}) (result Holder) {
	var err error
	if result, err = v.Path(path, optVars...); err != nil {
		bwerr.PanicA(bwerr.Err(err))
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

func (v Holder) Bool() (result bool, err error) {
	var ok bool
	if result, ok = Bool(v.Val); !ok {
		err = v.notOfValKindError(ValKindSetFrom(ValBool))
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
		err = v.notOfValKindError(ValKindSetFrom(ValString))
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
		err = v.notOfValKindError(ValKindSetFrom(ValInt))
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
		err = v.notOfValKindError(ValKindSetFrom(ValNumber))
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
		err = v.notOfValKindError(ValKindSetFrom(ValArray))
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
	default:
		err = v.notOfValKindError(ValKindSetFrom(ValArray))
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
		err = v.notOfValKindError(ValKindSetFrom(ValMap))
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
		bwerr.PanicA(bwerr.Err(err))
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

func (v Holder) KeyVal(key string, optDefaultValProvider ...defaultValProvider) (result interface{}, err error) {
	if v.Val == nil {
		var ok bool
		if result, ok = defaultVal(optDefaultValProvider); !ok {
			err = v.wrongValError()
		}
		return
	}
	var m map[string]interface{}
	if m, err = v.Map(); err != nil {
		return
	}
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
		bwerr.PanicA(bwerr.Err(err))
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
	// bwdebug.Print("v:json", v)
	if err != nil {
		err = bwerr.Refine(err, ansi.String("<ansiVal>%s<ansi>::{Error}"), bwjson.Pretty(v.Val))
	}
	return
}

func (v Holder) MustValidVal(def Def) (result interface{}) {
	var err error
	// bwdebug.Print("v:json", v, "def:json", def)
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
					err = Holder{val, vpi.Path}.notOfValKindError(ValKindSetFrom(ValInt, ValString))
				}
			}
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
	switch val, kind := Kind(v.Val); kind {
	case ValArray:
		vals, _ := val.([]interface{})
		nidx, ok = bw.NormalIdx(idx, len(vals))
		err = onArray(vals, nidx, ok)
	case ValArrayOfString:
		ss, _ := val.([]string)
		nidx, ok = bw.NormalIdx(idx, len(ss))
		err = onArrayOfString(ss, nidx, ok)
	default:
		err = v.notOfValKindError(ValKindSetFrom(ValArray))
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
