package defvalid

import (
	// "github.com/baza-winner/bwcore/bwint"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/defvalid/deftype"
	// "log"
)

func compileDef(def value) (result *Def, err error) {
	// log.Printf("compileDef::result: %#v", result)
	if def.value == nil {
		return nil, valueErrorMake(def, valueErrorHasNonSupportedValue)
	}
	var defType value
	var isSimpleDef bool
	if isSimpleDef = _isOfType(def.value, "string", "[]string"); isSimpleDef {
		defType = def
	} else if !_isOfType(def.value, "map[string]") {
		return nil, valueErrorMake(def, valueErrorIsNotOfType, "string", "[]string", "map[string]")
	} else {
		if defType, err = def.getKey("type", []string{"string", "[]string"}); err != nil {
			return nil, err
		}
	}

	var tp deftype.Set
	if tp, err = getDeftype(defType, isSimpleDef); err != nil {
		return nil, err
	}

	result = &Def{tp: tp}
	if !isSimpleDef {
		validDefKeys := bwset.StringsFromArgs("type")
		if tp.Has(deftype.String) {
			var enumVal value
			if enumVal, err = def.getKey("enum", "[]string", nil); err != nil {
				return nil, err
			}
			if enumVal.value != nil {
				validDefKeys.Add("enum")
				result.enum = bwset.StringsFromSlice(_mustBeSliceOfStrings(enumVal.value))
			}
		}
		if tp.Has(deftype.Map) {
			var keysVal value
			if keysVal, err = def.getKey("keys", "map[string]", nil); err != nil {
				return nil, err
			}
			if keysVal.value != nil {
				validDefKeys.Add("keys")
				result.keys = map[string]Def{}
				if err = keysVal.forEachMapString(func(k string, v interface{}) (err error) {
					var keyDef *Def
					if keyDef, err = compileDef(value{keysVal.what + "." + k, v}); err == nil {
						result.keys[k] = *keyDef
					}
					return
				}); err != nil {
					return nil, err
				}
			}
		}
		if tp.Has(deftype.Array) {
			var arrayElemVal value
			if arrayElemVal, err = def.getKey("arrayElem", "interface{}", nil); err != nil {
				return nil, err
			}
			if arrayElemVal.value != nil {
				validDefKeys.Add("arrayElem")
				if result.arrayElem, err = compileDef(arrayElemVal); err != nil {
					return nil, err
				}
			}
		}
		if tp.Has(deftype.Map) || tp.Has(deftype.Array) && result.arrayElem == nil {
			var elemVal value
			if elemVal, err = def.getKey("elem", "interface{}", nil); err != nil {
				return nil, err
			}
			if elemVal.value != nil {
				validDefKeys.Add("elem")
				if result.elem, err = compileDef(elemVal); err != nil {
					return nil, err
				}
			}
		}
		if tp.Has(deftype.Int) {
			var minIntVal value
			if minIntVal, err = def.getKey("minInt", "int64", nil); err != nil {
				return nil, err
			}
			if minIntVal.value != nil {
				validDefKeys.Add("minInt")
				result.minInt = ptrToInt64(_mustBeInt64(minIntVal.value))
			}
			var maxIntVal value
			if maxIntVal, err = def.getKey("maxInt", "int64", nil); err != nil {
				return nil, err
			}
			if maxIntVal.value != nil {
				validDefKeys.Add("maxInt")
				result.maxInt = ptrToInt64(_mustBeInt64(maxIntVal.value))
			}
			if result.minInt != nil && result.maxInt != nil && *(result.minInt) > *(result.maxInt) {
				return nil, valueErrorMake(def, valueErrorValuesCannotBeCombined, *(result.minInt), *(result.maxInt))
			}
		}
		if tp.Has(deftype.Number) {
			var minNumberVal value
			if minNumberVal, err = def.getKey("minNumber", "float64", nil); err != nil {
				return nil, err
			}
			if minNumberVal.value != nil {
				validDefKeys.Add("minNumber")
				result.minNumber = ptrToFloat64(_mustBeFloat64(minNumberVal.value))
			}
			var maxNumberVal value
			if maxNumberVal, err = def.getKey("maxNumber", "float64", nil); err != nil {
				return nil, err
			}
			if maxNumberVal.value != nil {
				validDefKeys.Add("maxNumber")
				result.maxNumber = ptrToFloat64(_mustBeFloat64(maxNumberVal.value))
			}
			if result.minNumber != nil && result.maxNumber != nil && *(result.minNumber) > *(result.maxNumber) {
				return nil, valueErrorMake(def, valueErrorValuesCannotBeCombined, *(result.minNumber), *(result.maxNumber))
			}
		}
		if unexpectedKeys := bwmap.GetUnexpectedKeys(def.value, validDefKeys); unexpectedKeys != nil {
			return nil, valueErrorMake(def, valueErrorHasUnexpectedKeys, unexpectedKeys)
		}
		// var dfltVal value
		// if dfltVal, err = def.getKey("dflt", "interface{}", nil); err != nil {
		// 	return nil, err
		// }
		// if dfltVal.value != nil {
		// 	if result.dflt, err = getValidVal(dfltVal, result); err != nil {
		// 		return nil, err
		// 	}
		// }

	}
	return
}

func getDeftype(defType value, isSimpleDef bool) (result deftype.Set, err error) {
	var isString bool
	var ss []string
	if _isOfType(defType.value, "string") {
		ss = []string{_mustBeString(defType.value)}
		isString = true
	} else if _isOfType(defType.value, "[]string") {
		ss = _mustBeSliceOfStrings(defType.value)
		isString = false
	} else {
		return nil, valueErrorMake(defType, valueErrorIsNotOfType, "string", "[]string")
	}
	result = deftype.Set{}
	for i, s := range ss {
		switch s {
		case "bool":
			result.Add(deftype.Bool)
		case "string":
			result.Add(deftype.String)
		case "int":
			result.Add(deftype.Int)
		case "number":
			result.Add(deftype.Number)
		case "map":
			result.Add(deftype.Map)
		case "array":
			result.Add(deftype.Array)
		case "orArrayOf":
			result.Add(deftype.OrArrayOf)
		default:
			if isString {
				err = valueErrorMake(defType, valueErrorHasNonSupportedValue)
			} else {
				var elem value
				if elem, err = defType.getElem(i, "string"); err != nil {
					return nil, err
				}
				err = valueErrorMake(elem, valueErrorHasNonSupportedValue)
			}
		}
	}
	return
}
