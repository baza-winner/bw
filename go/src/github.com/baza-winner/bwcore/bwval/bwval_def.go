package bwval

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
)

// ============================================================================

type Def struct {
	Types      ValKindSet
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
		// } else if compileDefResult == nil {
		//  bwerr.Panic("Unexpected behavior; def: %s", bwjson.Pretty(def))
	} else {
		result = *compileDefResult
	}
	return
}

// ============================================================================

func compileDef(def Holder) (result *Def, err error) {
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
		case ValString, ValArray:
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
						if keyDef, err = compileDef(vp); err != nil {
							return
						} else {
							result.Keys[k] = *keyDef
						}
					}
				}
			}
		}
		if types.Has(ValArray) {
			validDefKeys.Add("arrayElem")
			if vp = def.MustKey("arrayElem", nil); vp.Val != nil {
				if result.ArrayElem, err = compileDef(vp); err != nil {
					return
				}
			}
		}
		if types.Has(ValMap) || types.Has(ValArray) && result.ArrayElem == nil {
			validDefKeys.Add("elem")
			if vp = def.MustKey("elem", nil); vp.Val != nil {
				if result.Elem, err = compileDef(vp); err != nil {
					return
				}
			}
		}
		if types.Has(ValInt) {
			validDefKeys.Add("min", "max")
			rng := IntRange{}
			var limCount, n int
			if vp = def.MustKey("min", nil); vp.Val != nil {
				if n, err = vp.Int(); err != nil {
					return
				} else {
					rng.MinPtr = PtrToInt(n)
					limCount++
				}
			}
			if vp = def.MustKey("max", nil); vp.Val != nil {
				if n, err = vp.Int(); err != nil {
					return
				} else {
					rng.MaxPtr = PtrToInt(n)
					limCount++
				}
			}
			if limCount == 2 && *(rng.MinPtr) > *(rng.MaxPtr) {
				err = def.maxLessThanMinError()
				return
			}
			if limCount > 0 {
				result.Range = rng
			}
		} else if types.Has(ValNumber) {
			validDefKeys.Add("min", "max")
			rng := NumberRange{}
			var limCount, n float64
			if vp, err = def.Key("min"); err != nil || vp.Val == nil {
				err = nil
			} else if n, err = vp.Number(); err != nil {
				return
			} else {
				rng.MinPtr = PtrToNumber(n)
				limCount++
			}
			if vp, err = def.Key("max"); err != nil || vp.Val == nil {
				err = nil
			} else if n, err = vp.Number(); err != nil {
				return
			} else {
				rng.MaxPtr = PtrToNumber(n)
				limCount++
			}
			if limCount == 2 && *(rng.MinPtr) > *(rng.MaxPtr) {
				err = def.maxLessThanMinError()
				return
			}
			if limCount > 0 {
				result.Range = rng
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
		}

		validDefKeys.Add("isOptional")
		if vp, err = def.Key("isOptional"); err != nil || vp.Val == nil {
			err = nil
			result.IsOptional = result.Default != nil
		} else if result.IsOptional, err = vp.Bool(); err != nil {
			return
		}
		if !result.IsOptional && result.Default != nil {
			err = bwerr.From(
				ansi.String(
					"<ansiPath>%s<ansi> (<ansiVal>%s<ansi>) has conflicting keys: <ansiErr>isOptional<ansi> and <ansiErr>default",
				),
				def.Path, bwjson.Pretty(result),
			)
			return
		}
		if unexpectedKeys := bwmap.MustUnexpectedKeys(def.Val, validDefKeys); unexpectedKeys != nil {
			err = bwerr.From(
				ansi.String(
					"<ansiPath>%s<ansi> (<ansiVal>%s<ansi>) has unexpected keys: <ansiVal>%s",
				),
				def.Path, bwjson.Pretty(unexpectedKeys),
			)
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
		ss = []string{MustString(val)}
		isString = true
	case ValArray:
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
			err = bwerr.From(ansi.String("%s: <ansiVal>ArrayOf<ansi> must be followed by some type"), defType)
		} else if result.Has(ValArray) {
			err = bwerr.From(ansi.String("%s type values <ansiVal>ArrayOf<ansi> and <ansiVal>Array<ansi> can not be combined"), defType)
		}
	}
	if result.Has(ValInt) && result.Has(ValNumber) {
		err = bwerr.From(ansi.String("%s type values <ansiVal>Int<ansi> and <ansiVal>Number<ansi> can not be combined"), defType)
	}
	return
}

// ============================================================================

func (v Holder) mapHelper(key string, elemDef Def) (err error) {
	vp, _ := v.Key(key)
	var val interface{}
	if val, err = vp.validVal(elemDef); err != nil {
		return
	} else if val != nil {
		if err = v.SetKeyVal(val, key); err != nil {
			return
		}
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
	arr := MustArray(v.Val)
	newArr := make([]interface{}, 0, len(arr))
	var vp Holder
	for i := range arr {
		vp, _ = v.Idx(i)
		var val interface{}
		if val, err = vp.validVal(elemDef, optSkipArrayOf...); err != nil {
			return
		}
		newArr = append(newArr, val)
	}
	result = newArr
	return
}

func (v Holder) validVal(def Def, optSkipArrayOf ...bool) (result interface{}, err error) {
	skipArrayOf := optSkipArrayOf != nil && optSkipArrayOf[0]
	if v.Val == nil {
		if !skipArrayOf {
			if def.Default != nil {
				return def.Default, nil
			}
			if def.IsOptional {
				return nil, nil
			}
		}
		if def.Types.Has(ValMap) {
			v.Val = map[string]interface{}{}
		} else {
			// ss := def.Types.ToSliceOfStrings()
			err = v.notOfValKindError(def.Types)
			return
		}
	}
	var valDeftype ValKind
	switch _, kind := Kind(v.Val); kind {
	case ValBool:
		if def.Types.Has(ValBool) {
			valDeftype = ValBool
		}
	case ValInt:
		if def.Types.Has(ValInt) {
			valDeftype = ValInt
		} else if def.Types.Has(ValNumber) {
			valDeftype = ValNumber
		}
	case ValNumber:
		if def.Types.Has(ValNumber) {
			valDeftype = ValNumber
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
	case ValString:
		if def.Types.Has(ValString) {
			valDeftype = ValString
		}
	}

	if valDeftype == ValUnknown {
		// ss := def.Types.ToSliceOfStrings()
		err = v.notOfValKindError(def.Types)
		return
	}

	switch valDeftype {
	case ValBool:
	case ValString:
		if def.Enum != nil {
			if !def.Enum.Has(MustString(v.Val)) {
				bwerr.TODO() // enum instead nonSupportedValueError
				err = v.nonSupportedValueError()
				return
			}
		}
	case ValInt, ValNumber:
		// if def.Range != nil && def.Range{}
		if !RangeContains(def.Range, v.Val) {
			// if !v.inRange(def.Range) {
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
					for _, key := range unexpectedKeys {
						if err = v.mapHelper(key, *(def.Elem)); err != nil {
							return
						}
					}
				}
			}
		} else if def.Elem != nil {
			m, _ := Map(v.Val)
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
		}
	case ValArrayOf:
		if v.Val, err = v.arrayHelper(def, true); err != nil {
			return
		}
	}

	if !skipArrayOf && valDeftype != ValArrayOf && def.Types.Has(ValArrayOf) {
		v.Val = []interface{}{v.Val}
	}

	result = v.Val
	return
}

// ============================================================================
