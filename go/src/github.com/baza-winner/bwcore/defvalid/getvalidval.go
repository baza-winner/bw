package defvalid

import (
	"github.com/baza-winner/bwcore/bwerror"
	// "github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/defvalid/deftype"
	// "log"
	"reflect"
)

func getValidVal(val value, def Def) (result interface{}, err error) {
	// var defType value
	// var isSimpleDef bool
	// var defaultVal value
	// if isSimpleDef = _isOfType(def.value, "string", "[]string"); isSimpleDef {
	// 	defType = def
	// } else if defType, err = def.getKey("type", []string{"string", "[]string"}); err != nil {
	// 	return nil, err
	// } else if skipDefault == nil || !skipDefault[0] {
	// 	defaultVal, err = def.getKey("default")
	// 	if err != nil {
	// 		if valErr, ok := err.(valueError); ok && valErr.errorType == valueErrorHasNoKey {
	// 			err = nil
	// 		} else {
	// 			return nil, err
	// 		}
	// 	} else if defaultVal.value == nil {
	// 		return nil, valueErrorMake(defaultVal, valueErrorHasNonSupportedValue)
	// 	} else if defaultVal.value, err = getValidVal(defaultVal, def, true); err != nil {
	// 		return nil, err
	// 	}
	// }
	// var def.tp bwset.Strings
	// if def.tp, err = getExpectedTypes(defType, isSimpleDef); err != nil {
	// 	return nil, err
	// }

	var typeName string
	if val.value == nil {
		if def.isSimple {
			return nil, valueErrorMake(val, valueErrorIsNotOfType, def.tp)
		} else if def.dflt != nil {
			return def.dflt, nil
		} else if def.tp.Has(deftype.Map) {
			val.value = map[string]interface{}{}
		} else {
			return nil, valueErrorMake(val, valueErrorIsNotOfType, def.tp)
		}
	}
	valType := reflect.TypeOf(val.value)
	switch valType.Kind() {
	case reflect.Bool:
		if def.tp.Has(deftype.Bool) {
			typeName = "bool"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if def.tp.Has(deftype.Int) {
			typeName = "int"
		} else if def.tp.Has(deftype.Number) {
			typeName = "number"
		}
	case reflect.Float32, reflect.Float64:
		if def.tp.Has(deftype.Number) {
			typeName = "number"
		}

	case reflect.Map:
		if valType.Key().Kind() == reflect.String && valType.Elem().Kind() == reflect.Interface && def.tp.Has(deftype.Map) {
			typeName = "map"
		}
	case reflect.Slice:
		if def.tp.Has(deftype.Array) {
			typeName = "array"
		}
	case reflect.String:
		if def.tp.Has(deftype.String) {
			typeName = "string"
		}
	}
	if len(typeName) == 0 {
		return nil, valueErrorMake(val, valueErrorIsNotOfType, def.tp)
	}

	if !def.isSimple {

		// var validDefKeys bwset.Strings
		// validDefKeys = bwset.StringsFromArgs("type", "default")
		switch typeName {
		case "map":
			// validDefKeys.Add("keys")
			// var defKeys value
			// if defKeys, err = def.getKey("keys", "map", nil); err != nil {
			// 	return nil, err
			// } else
			if def.keys != nil {
				unexpectedKeys := bwmap.GetUnexpectedKeys(val.value, def.keys)
				for key, keyDef := range def.keys {
					if err = helper(val, key, keyDef); err != nil {
						return nil, err
					}
				}
				if unexpectedKeys != nil {
					if def.elem == nil {
						return nil, valueErrorMake(val, valueErrorHasUnexpectedKeys, unexpectedKeys)
					} else {
						for _, key := range unexpectedKeys {
							if err = helper(val, key, *(def.elem)); err != nil {
								return nil, err
							}
						}
					}
				}
			} else if def.elem != nil {
				val.forEachMapString(func(k string, v interface{}) (err error) {
					err = helper(val, k, *(def.elem))
					return
				})
			}

		case "bool":
		case "string":
			// case "enum":
			// validDefKeys.Add("enum")
			// var enumValues value
			// if enumValues, err = def.getKey("enum", "[]string"); err != nil {
			// 	return nil, err
			// }
			// enumSet := bwset.StringsFromSlice(_mustBeSliceOfStrings(enumValues.value))
			if def.enum != nil {
				if !def.enum.Has(_mustBeString(val.value)) {
					return nil, valueErrorMake(val, valueErrorHasNonSupportedValue)
				}
			}

		default:
			bwerror.Panic("typeName: %s", typeName)
		}
		// if unexpectedKeys := bwmap.GetUnexpectedKeys(def.mustBeMap(), validDefKeys); unexpectedKeys != nil {
		// 	return nil, valueErrorMake(def, valueErrorHasUnexpectedKeys, unexpectedKeys)
		// }
	}

	return val.value, nil
}

func helper(val value, key string, keyDef Def) (err error) {
	var valKeyVal value
	if valKeyVal, err = val.getKey(key, "interface{}", nil); err != nil {
		return err
	}
	var valKeyValIntf interface{}
	if valKeyValIntf, err = getValidVal(valKeyVal, keyDef); err != nil {
		return err
	}
	if valKeyValIntf != nil {
		if err = val.setKey(key, valKeyValIntf); err != nil {
			return err
		}
	}
	return
}

func getExpectedTypes(defType value, isSimpleDef bool) (result bwset.Strings, err error) {
	var isString bool
	var ss []string
	if _isOfType(defType.value, "string") {
		ss = []string{_mustBeString(defType.value)}
		isString = true
	} else if _isOfType(defType.value, "[]string") {
		ss = _mustBeSliceOfStrings(defType.value)
		isString = false
	} else {
		err = valueErrorMake(defType, valueErrorIsNotOfType, "string", "[]string")
	}
	if err == nil {
		for i, s := range ss {
			switch s {
			case "map":
			case "bool":
			case "string":
			case "enum":
				var elem value
				if isSimpleDef {
					if isString {
						err = valueErrorMake(defType, valueErrorHasNonSupportedValue)
					} else if elem, err = defType.getElem(i); err == nil {
						err = valueErrorMake(elem, valueErrorHasNonSupportedValue)
					}
				}
			default:
				var elem value
				if isString {
					err = valueErrorMake(defType, valueErrorHasNonSupportedValue)
				} else if elem, err = defType.getElem(i); err == nil {
					err = valueErrorMake(elem, valueErrorHasNonSupportedValue)
				}
			}
		}
		result = bwset.StringsFromSlice(ss)
		if result.Has("enum") && result.Has("string") {
			// log.Printf("%s\n", result)
			err = valueErrorMake(defType, valueErrorValuesCannotBeCombined, "enum", "string")
		}
	}
	return
}
