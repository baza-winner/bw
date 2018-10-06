package defvalid

import (
	"github.com/baza-winner/bwcore/bwerror"
	// "github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
	// "github.com/baza-winner/bwcore/defparse"
	// "log"
	"reflect"
)

func getValidVal(val, def value, skipDefault ...bool) (result interface{}, err error) {
	var defType value
	var isSimpleDef bool
	var defaultVal value
	if isSimpleDef = _isOfType(def.value, "string", "[]string"); isSimpleDef {
		defType = def
	} else if defType, err = def.getKey("type", []string{"string", "[]string"}); err != nil {
		return nil, err
	} else if skipDefault == nil || !skipDefault[0] {
		defaultVal, err = def.getKey("default")
		if err != nil {
			if valErr, ok := err.(valueError); ok && valErr.errorType == valueErrorHasNoKey {
				err = nil
			} else {
				return nil, err
			}
		} else if defaultVal.value == nil {
			return nil, valueErrorMake(defaultVal, valueErrorHasNonSupportedValue)
		} else if defaultVal.value, err = getValidVal(defaultVal, def, true); err != nil {
			return nil, err
		}
	}
	var expectedTypes bwset.Strings
	if expectedTypes, err = getExpectedTypes(defType, isSimpleDef); err != nil {
		return nil, err
	}

	var typeName string
	if val.value == nil {
		if isSimpleDef {
			return nil, valueErrorMake(val, valueErrorIsNotOfType, expectedTypes)
		} else if defaultVal.value != nil {
			return defaultVal.value, nil
		} else if expectedTypes.Has("map") {
			val.value = map[string]interface{}{}
		} else {
			return nil, valueErrorMake(val, valueErrorIsNotOfType, expectedTypes)
		}
	}
	valType := reflect.TypeOf(val.value)
	switch valType.Kind() {
	case reflect.Bool:
		if expectedTypes.Has("bool") {
			typeName = "bool"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if expectedTypes.Has("int") {
			typeName = "int"
		} else if expectedTypes.Has("number") {
			typeName = "number"
		}
	case reflect.Float32, reflect.Float64:
		if expectedTypes.Has("number") {
			typeName = "number"
		}
	case reflect.Map:
		if valType.Key().Kind() == reflect.String && valType.Elem().Kind() == reflect.Interface && expectedTypes.Has("map") {
			typeName = "map"
		}
	case reflect.Slice:
		if expectedTypes.Has("array") {
			typeName = "array"
		}
	case reflect.String:
		if expectedTypes.Has("string") {
			typeName = "string"
		} else if expectedTypes.Has("enum") {
			typeName = "enum"
		}
	}
	if len(typeName) == 0 {
		return nil, valueErrorMake(val, valueErrorIsNotOfType, expectedTypes)
	}
	var validDefKeys bwset.Strings
	if !isSimpleDef {
		validDefKeys = bwset.FromArgs("type", "default")
	}

	if !isSimpleDef {
		switch typeName {
		case "map":
			validDefKeys.Add("keys")
			var defKeys value
			if defKeys, err = def.getKey("keys", "map", nil); err != nil {
				return nil, err
			} else if defKeys.value != nil {
				valAsMap, _ := val.asMap()
				if unexpectedKeys := bwmap.GetUnexpectedKeys(valAsMap, defKeys.mustBeMap()); unexpectedKeys != nil {
					return nil, valueErrorMake(val, valueErrorHasUnexpectedKeys, unexpectedKeys)
				}
				for defKeysKey, _ := range defKeys.mustBeMap() {
					var defKeysKeyVal value
					if defKeysKeyVal, err = defKeys.getKey(defKeysKey); err != nil {
						return nil, err
					}
					var valKeyVal value
					if valKeyVal, err = val.getKey(defKeysKey); err != nil {
						if valErr, ok := err.(valueError); !ok || valErr.errorType != valueErrorHasNoKey {
							return nil, err
						} else {
							valKeyVal = value{what: val.what + "." + defKeysKey, value: nil}
						}
					}
					if valKeyVal.value == nil {
						if _isOfType(defKeysKeyVal.value, "string", "[]string") {
							if _, err = val.getKey(defKeysKey); err != nil {
								if valErr, ok := err.(valueError); ok && valErr.errorType == valueErrorHasNoKey {
									return nil, valueErrorMake(val, valueErrorHasNoKey, defKeysKey)
								} else {
									return nil, err
								}
							}
						} else if _isOfType(defKeysKeyVal.value, "map") {
							if _, err = defKeysKeyVal.getKey("default"); err != nil {
							}
						}
					}
					if valAsMap[defKeysKey], err = getValidVal(valKeyVal, defKeysKeyVal); err != nil {
						return nil, err
					}
				}
			}

		case "bool":
		case "string":
		case "enum":
			validDefKeys.Add("enum")
			var enumValues value
			if enumValues, err = def.getKey("enum", "[]string"); err != nil {
				return nil, err
			}
			enumSet := bwset.FromSlice(_mustBeSliceOfStrings(enumValues.value))
			if !enumSet.Has(_mustBeString(val.value)) {
				return nil, valueErrorMake(val, valueErrorHasNonSupportedValue)
			}

		default:
			bwerror.Panic("typeName: %s", typeName)
		}
	}

	if !isSimpleDef {
		if unexpectedKeys := bwmap.GetUnexpectedKeys(def.mustBeMap(), validDefKeys); unexpectedKeys != nil {
			return nil, valueErrorMake(def, valueErrorHasUnexpectedKeys, unexpectedKeys)
		}
	}
	return val.value, nil
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
		result = bwset.FromSlice(ss)
		if result.Has("enum") && result.Has("string") {
			// log.Printf("%s\n", result)
			err = valueErrorMake(defType, valueErrorValuesCannotBeCombined, "enum", "string")
		}
	}
	return
}
