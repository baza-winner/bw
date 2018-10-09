package defvalid

import (
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
	var isSimple bool
	validDefKeys := bwset.Strings{}
	if isSimple = _isOfType(def.value, "string", "[]string"); isSimple {
		defType = def
	} else if !_isOfType(def.value, "map[string]") {
		return nil, valueErrorMake(def, valueErrorIsNotOfType, "string", "[]string", "map[string]")
	} else {
		if defType, err = getDefKey(def, "type", []string{"string", "[]string"}, nil, &validDefKeys); err != nil {
			return nil, err
		}
	}

	var tp deftype.Set
	if tp, err = getDeftype(defType, isSimple); err != nil {
		return nil, err
	}

	result = &Def{tp: tp}
	if !isSimple {
		if tp.Has(deftype.String) {
			var enumVal value
			if enumVal, err = getDefKey(def, "enum", "[]string", nil, &validDefKeys); err != nil {
				return nil, err
			}
			if enumVal.value != nil {
				result.enum = bwset.StringsFromSlice(_mustBeSliceOfStrings(enumVal.value))
			}
		}
		if tp.Has(deftype.Map) {
			var keysVal value
			if keysVal, err = getDefKey(def, "keys", "map[string]", nil, &validDefKeys); err != nil {
				return nil, err
			}
			if keysVal.value != nil {
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
			if arrayElemVal, err = getDefKey(def, "arrayElem", "interface{}", nil, &validDefKeys); err != nil {
				return nil, err
			}
			if arrayElemVal.value != nil {
				if result.arrayElem, err = compileDef(arrayElemVal); err != nil {
					return nil, err
				}
			}
		}
		if tp.Has(deftype.Map) || tp.Has(deftype.Array) && result.arrayElem == nil {
			var elemVal value
			if elemVal, err = getDefKey(def, "elem", "interface{}", nil, &validDefKeys); err != nil {
				return nil, err
			}
			if elemVal.value != nil {
				if result.elem, err = compileDef(elemVal); err != nil {
					return nil, err
				}
			}
		}
		if tp.Has(deftype.Int) {
			var minIntVal value
			if minIntVal, err = getDefKey(def, "minInt", "int64", nil, &validDefKeys); err != nil {
				return nil, err
			}
			if minIntVal.value != nil {
				result.minInt = ptrToInt64(_mustBeInt64(minIntVal.value))
			}
			var maxIntVal value
			if maxIntVal, err = getDefKey(def, "maxInt", "int64", nil, &validDefKeys); err != nil {
				return nil, err
			}
			if maxIntVal.value != nil {
				result.maxInt = ptrToInt64(_mustBeInt64(maxIntVal.value))
			}
			if result.minInt != nil && result.maxInt != nil && *(result.minInt) > *(result.maxInt) {
				return nil, valueErrorMake(def, valueErrorConflictingKeys, map[string]interface{}{
					"minInt": *(result.minInt),
					"maxInt": *(result.maxInt),
				})
			}
		}
		if tp.Has(deftype.Number) {
			var minNumberVal value
			if minNumberVal, err = getDefKey(def, "minNumber", "float64", nil, &validDefKeys); err != nil {
				return nil, err
			}
			if minNumberVal.value != nil {
				result.minNumber = ptrToFloat64(_mustBeFloat64(minNumberVal.value))
			}
			var maxNumberVal value
			if maxNumberVal, err = getDefKey(def, "maxNumber", "float64", nil, &validDefKeys); err != nil {
				return nil, err
			}
			if maxNumberVal.value != nil {
				result.maxNumber = ptrToFloat64(_mustBeFloat64(maxNumberVal.value))
			}
			if result.minNumber != nil && result.maxNumber != nil && *(result.minNumber) > *(result.maxNumber) {
				return nil, valueErrorMake(def, valueErrorConflictingKeys, map[string]interface{}{
					"minNumber": *(result.minNumber),
					"maxNumber": *(result.maxNumber),
				})
			}
		}
		var dfltVal value
		if dfltVal, err = getDefKey(def, "default", "interface{}", nil, &validDefKeys); err != nil {
			return nil, err
		}
		if dfltVal.value != nil {
			dfltDef := *result
			if result.tp.Has(deftype.ArrayOf) {
				dfltDef = Def{
					tp: deftype.FromArgs(deftype.Array),
					arrayElem: &Def{
						tp:         result.tp.Copy(),
						isOptional: false,
						enum:       result.enum,
						minInt:     result.minInt,
						maxInt:     result.maxInt,
						minNumber:  result.minNumber,
						maxNumber:  result.maxNumber,
						keys:       result.keys,
						elem:       result.elem,
					},
				}
				dfltDef.arrayElem.tp.Del(deftype.ArrayOf)
			}
			if result.dflt, err = getValidVal(dfltVal, dfltDef); err != nil {
				return nil, err
			}
		}
		var boolVal value
		if boolVal, err = getDefKey(def, "isOptional", "bool", result.dflt != nil, &validDefKeys); err != nil {
			return nil, err
		}
		result.isOptional = _mustBeBool(boolVal.value)
		if !result.isOptional && result.dflt != nil {
			return nil, valueErrorMake(def, valueErrorConflictingKeys, map[string]interface{}{
				"isOptional": result.isOptional,
				"default":    result.dflt,
			})
		}
		if unexpectedKeys := bwmap.GetUnexpectedKeys(def.value, validDefKeys); unexpectedKeys != nil {
			return nil, valueErrorMake(def, valueErrorHasUnexpectedKeys, unexpectedKeys)
		}
	}
	return
}

func getDefKey(def value, keyName string, ofType interface{}, defaultValue interface{}, validDefKeys *bwset.Strings) (keyValue value, err error) {
	keyValue, err = def.getKey(keyName, ofType, defaultValue)
	validDefKeys.Add(keyName)
	return
}

func getDeftype(defType value, isSimple bool) (result deftype.Set, err error) {
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
		case "arrayOf":
			result.Add(deftype.ArrayOf)
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
	if result.Has(deftype.ArrayOf) {
		if len(result) < 2 {
			err = valueErrorMake(defType, valueErrorArrayOf)
		} else if result.Has(deftype.Array) {
			err = valueErrorMake(defType, valueErrorValuesCannotBeCombined, "array", "arrayOf")
		}
	}
	return
}
