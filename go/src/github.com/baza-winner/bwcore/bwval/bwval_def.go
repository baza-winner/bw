package bwval

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwparse"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/bwtype"
)

// ============================================================================

type Def struct {
	Types      ValKindSet
	IsOptional bool
	Enum       bwset.String
	Range      bwtype.Range
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
	if v.Range.Kind() != bwtype.RangeNo {
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

func DefFrom(src bwrune.Provider, optPathProvider ...bw.ValPathProvider) (result Def, err error) {
	p := bwparse.From(src, map[string]interface{}{
		"idVals": map[string]interface{}{
			"Bool":    "Bool",
			"String":  "String",
			"Int":     "Int",
			"Float64": "Float64",
			"Array":   "Array",
			"ArrayOf": "ArrayOf",
		}},
	)
	if _, err = bwparse.SkipSpace(p, bwparse.TillNonEOF); err != nil {
		return
	}
	var def interface{}
	var st bwparse.Status
	if def, st = bwparse.Val(p); !st.IsOK() {
		if st.Err != nil {
			err = st.Err
		} else {
			err = bwparse.Unexpected(p)
		}
		return
	}
	var compileDefResult *Def
	h := Holder{Val: def}

	if len(optPathProvider) > 0 {
		if h.Pth, err = optPathProvider[0].Path(); err != nil {
			return
		}
	}
	if compileDefResult, err = h.compileDef(); err == nil {
		result = *compileDefResult
	}
	return
}

func MustDefFrom(src bwrune.Provider, optPathProvider ...bw.ValPathProvider) (result Def) {
	var err error
	if result, err = DefFrom(src, optPathProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

// ============================================================================

func (def Holder) compileDef() (result *Def, err error) {
	if def.Val == nil {
		err = def.wrongValError()
		// err = valAtPathIsNil(def.Path)
		return
	}
	var defType Holder
	var isSimple bool
	validDefKeys := bwset.String{}
	switch _, kind := Kind(def.Val); kind {
	case ValString, ValArray:
		isSimple = true
		defType = def
	case ValMap:
		validDefKeys.Add("type")
		if defType, err = def.Key("type"); err != nil {
			return
		}
		switch _, kind := Kind(defType.Val); kind {
		case ValString, ValArray, ValArrayOfString:
		default:
			err = defType.notOfValKindError(ValKindSetFrom(ValString, ValArray))
			return
		}
	default:
		err = def.notOfValKindError(ValKindSetFrom(ValString, ValArray, ValMap))
		return
	}

	var types ValKindSet
	if types, err = getDeftype(defType, isSimple); err != nil {
		return
	}

	result = &Def{Types: types}
	if !isSimple {
		var vp Holder
		if types.Has(ValString) {
			var ss []string
			validDefKeys.Add("enum")
			if vp = def.MustKey("enum", nil); vp.Val != nil {
				if ss, err = vp.ArrayOfString(); err != nil {
					return
				} else {
					result.Enum = bwset.StringFromSlice(ss)
				}
			}
		}
		if types.Has(ValMap) {
			var keysVp Holder
			validDefKeys.Add("keys")
			var m map[string]interface{}
			if keysVp = def.MustKey("keys", nil); keysVp.Val != nil {
				if m, err = keysVp.Map(); err != nil {
					return
				} else {
					result.Keys = map[string]Def{}
					var vp Holder
					var keyDef *Def
					for k := range m {
						vp, _ = keysVp.Key(k)
						if keyDef, err = vp.compileDef(); err != nil {
							return
						} else {
							result.Keys[k] = *keyDef
						}
					}
				}
			}
		}
		// bwdebug.Print("types:s", types)
		if types.Has(ValArray) {
			// bwdebug.Print("!HERE")
			validDefKeys.Add("arrayElem")
			if vp = def.MustKey("arrayElem", nil); vp.Val != nil {
				if result.ArrayElem, err = vp.compileDef(); err != nil {
					return
				}
			}
		}
		if types.Has(ValMap) || types.Has(ValArray) && result.ArrayElem == nil {
			validDefKeys.Add("elem")
			if vp = def.MustKey("elem", nil); vp.Val != nil {
				if result.Elem, err = vp.compileDef(); err != nil {
					return
				}
			}
		}
		hasInt := types.Has(ValInt)
		if hasInt || types.Has(ValFloat64) {
			validDefKeys.Add("min", "max")
			var min, max bwtype.RangeLimit
			var limCount int
			var getLimit func(key string) (result bwtype.RangeLimit, err error)

			if hasInt {
				getLimit = func(key string) (result bwtype.RangeLimit, err error) {
					if vp = def.MustKey(key, nil); vp.Val != nil {
						var i int
						if i, err = vp.Int(); err != nil {
							return
						} else {
							result = bwtype.MustRangeLimitFrom(i)
							limCount++
						}
					}
					return
				}
			} else {
				getLimit = func(key string) (result bwtype.RangeLimit, err error) {
					if vp = def.MustKey(key, nil); vp.Val != nil {
						var ok bool
						if result, ok = bwtype.RangeLimitFrom(vp.Val); !ok {
							err = vp.notOfValKindError(ValKindSetFrom(ValFloat64))
							return
						} else {
							limCount++
						}
					}
					return
				}
			}

			if min, err = getLimit("min"); err != nil {
				return
			}
			if max, err = getLimit("max"); err != nil {
				return
			}
			if limCount == 2 && max.MustNumber().IsLessThan(min.MustNumber()) {
				err = def.maxLessThanMinError(max, min)
				return
			}
			if limCount > 0 {
				result.Range = bwtype.MustRangeFrom(bwtype.A{Min: min, Max: max})
			}
		}

		validDefKeys.Add("default")
		if vp, err = def.Key("default"); err != nil || vp.Val == nil {
			err = nil
		} else {
			dfltDef := *result
			if types.Has(ValArrayOf) {
				dfltDef = Def{
					Types: ValKindSetFrom(ValArray),
					ArrayElem: &Def{
						Types:      result.Types.Copy(),
						IsOptional: false,
						Enum:       result.Enum,
						Range:      result.Range,
						Keys:       result.Keys,
						Elem:       result.Elem,
					},
				}
				dfltDef.ArrayElem.Types.Del(ValArrayOf)
			}
			if result.Default, err = vp.validVal(dfltDef); err != nil {
				return nil, err
			}
			// bwdebug.Print("vp.Val:#v", vp.Val, "result.Default", result.Default, "dfltDef:json", dfltDef)
		}

		validDefKeys.Add("isOptional")
		if vp, err = def.Key("isOptional"); err != nil || vp.Val == nil {
			err = nil
			result.IsOptional = result.Default != nil
			// } else if result.IsOptional, err = vp.Bool(); err != nil {
		} else if t, kind := Kind(vp.Val); kind != ValBool {
			err = vp.notOfValKindError(ValKindSetFrom(ValBool))
			return
		} else {
			result.IsOptional, _ = t.(bool)
		}
		if !result.IsOptional && result.Default != nil {
			err = def.defaultNonOptionalError()
			return
		}
		if unexpectedKeys := bwmap.MustUnexpectedKeys(def.Val, validDefKeys); unexpectedKeys != nil {
			err = def.unexpectedKeysError(unexpectedKeys)
			return
		}
	}
	return
}

func getDeftype(defType Holder, isSimple bool) (result ValKindSet, err error) {
	var ss []string
	var isString bool
	switch val, kind := Kind(defType.Val); kind {
	case ValString:
		s, _ := val.(string)
		ss = []string{s}
		isString = true
	case ValArray, ValArrayOfString:
		if ss, err = defType.ArrayOfString(); err != nil {
			return
		}
	}
	result = ValKindSet{}
	for i, s := range ss {
		var tpItem ValKind
		if tpItem, err = ValKindFromString(s); err == nil {
			result.Add(tpItem)
		} else {
			elem := defType
			if !isString {
				elem, _ = defType.Idx(i)
			}
			err = elem.nonSupportedValueError()
			return
		}
	}
	if result.Has(ValArrayOf) {
		if len(result) < 2 {
			err = defType.arrayOfMustBeFollowedBySomeTypeError()
		} else if result.Has(ValArray) {
			err = defType.valuesAreMutuallyExclusiveError("ArrayOf", "Array")
		}
	}
	if err == nil && result.Has(ValInt) && result.Has(ValFloat64) {
		err = defType.valuesAreMutuallyExclusiveError("Int", "Float64")
	}
	return
}

// ============================================================================

func (v Holder) validVal(def Def, optSkipArrayOf ...bool) (result interface{}, err error) {
	skipArrayOf := optSkipArrayOf != nil && optSkipArrayOf[0]
	if v.Val == nil {
		if !skipArrayOf {
			if def.Default != nil {
				result = def.Default
				return
			}
			if def.IsOptional {
				return
			}
		}
		if def.Types.Has(ValMap) {
			v.Val = map[string]interface{}{}
		} else {
			err = v.notOfValKindError(def.Types)
			return
		}
	}
	var valDeftype ValKind
	var kind ValKind
	var val interface{}
	switch val, kind = Kind(v.Val); kind {
	case ValBool:
		if def.Types.Has(ValBool) {
			valDeftype = ValBool
		}
	case ValInt:
		if def.Types.Has(ValInt) {
			valDeftype = ValInt
		} else if def.Types.Has(ValFloat64) {
			valDeftype = ValFloat64
		}
	case ValFloat64:
		if def.Types.Has(ValFloat64) {
			valDeftype = ValFloat64
		}

	case ValMap:
		if def.Types.Has(ValMap) {
			valDeftype = ValMap
		}
	case ValArray:
		if def.Types.Has(ValArray) {
			valDeftype = ValArray
		} else if !skipArrayOf && def.Types.Has(ValArrayOf) {
			valDeftype = ValArrayOf
		}
	case ValArrayOfString:
		if def.Types.Has(ValArray) {
			valDeftype = ValArray
		} else if !skipArrayOf && def.Types.Has(ValArrayOf) {
			valDeftype = ValArrayOf
		}
	case ValString:
		if def.Types.Has(ValString) {
			valDeftype = ValString
		}
	}

	if valDeftype == ValUnknown {
		types := def.Types
		if skipArrayOf {
			types = types.Copy()
			types.Del(ValArrayOf)
		}
		err = v.notOfValKindError(types)
		return
	}

	switch valDeftype {
	case ValBool:
	case ValString:
		if def.Enum != nil {
			s, _ := val.(string)
			if !def.Enum.Has(s) {
				err = v.unexpectedEnumValueError(def.Enum)
				return
			}
		}
	case ValInt, ValFloat64:
		if !def.Range.Contains(v.Val) {
			err = v.outOfRangeError(def.Range)
			return
		}

	case ValMap:
		if def.Keys != nil {
			unexpectedKeys := bwmap.MustUnexpectedKeys(v.Val, def.Keys)
			for key, keyDef := range def.Keys {
				if err = v.mapHelper(key, keyDef); err != nil {
					return
				}
			}
			if unexpectedKeys != nil {
				if def.Elem == nil {
					err = v.unexpectedKeysError(unexpectedKeys)
					return
				} else {
					for _, key := range unexpectedKeys.ToSlice() {
						if err = v.mapHelper(key, *(def.Elem)); err != nil {
							return
						}
					}
				}
			}
		} else if def.Elem != nil {
			m, _ := val.(map[string]interface{})
			for k := range m {
				if err = v.mapHelper(k, *(def.Elem)); err != nil {
					return
				}
			}
		}
	case ValArray:
		elemDef := def.ArrayElem
		if elemDef == nil {
			elemDef = def.Elem
		}
		if elemDef != nil {
			if v.Val, err = v.arrayHelper(*elemDef); err != nil {
				return
			}
		} else if kind == ValArrayOfString {
			ss, _ := v.Val.([]string)
			newArr := make([]interface{}, 0, len(ss))
			for _, s := range ss {
				newArr = append(newArr, s)
			}
			v.Val = newArr
		}
	case ValArrayOf:
		if v.Val, err = v.arrayHelper(def, true); err != nil {
			return
		}
	}

	if !skipArrayOf && valDeftype != ValArrayOf && def.Types.Has(ValArrayOf) {
		v.Val = []interface{}{v.Val}
	}

	// bwdebug.Print("v.Val:#v", v.Val, "valDeftype", valDeftype, "kind:s", kind, "def:json", def)
	result = v.Val
	return
}

func (v Holder) mapHelper(key string, elemDef Def) (err error) {
	vp, _ := v.Key(key)
	var val interface{}
	// // // // bwdebug.Print("vp:json", vp, "elemDef:json", elemDef)
	if val, err = vp.validVal(elemDef); err != nil {
		return
	} else if val != nil {
		v.SetKeyVal(val, key)
		// if _ = ; err != nil {
		// 	return
		// }
	}
	return
}

// func (v Holder) arrayHelper(elemDef Def, optSkipArrayOf ...bool) (err error) {
//  var vp Holder
//  for i := range MustArray(v.Val) {
//    vp, _ = v.Idx(i)
//    if _, err = vp.ValidVal(elemDef, optSkipArrayOf...); err != nil {
//      return
//    }
//  }
//  return
// }

func (v Holder) arrayHelper(elemDef Def, optSkipArrayOf ...bool) (result interface{}, err error) {
	appendIdxVal := func(newArr []interface{}, i int) (result []interface{}, err error) {
		var vp Holder
		if vp, err = v.Idx(i); err == nil {
			var val interface{}
			if val, err = vp.validVal(elemDef, optSkipArrayOf...); err == nil {
				result = append(newArr, val)
			}
		}
		// bwdebug.Print("v:#v", v, "result", result, "elemDef:json", elemDef)
		return
	}
	switch val, kind := Kind(v.Val); kind {
	case ValArray:
		arr, _ := val.([]interface{})
		newArr := make([]interface{}, 0, len(arr))
		for i := range arr {
			if newArr, err = appendIdxVal(newArr, i); err != nil {
				return
			}
		}
		result = newArr
	case ValArrayOfString:
		ss, _ := val.([]string)
		newArr := make([]interface{}, 0, len(ss))
		for i := range ss {
			if newArr, err = appendIdxVal(newArr, i); err != nil {
				return
			}
		}
		result = newArr
	}
	return
}

// ============================================================================
