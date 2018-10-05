/*
Предоставляет функции для валидации interface{}.
*/
package defvalid

import (
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/defparse"
	"reflect"
	"log"
)

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
		err = defType.err(valueErrorIsNotOfType, "string", "[]string")
	}
	if err == nil {
		for i, s := range ss {
			switch s {
			case "map":
			case "bool":
			case "string":
			case "enum":
				if isSimpleDef {
					if isString {
						err = defType.err(valueErrorHasNonSupportedValue)
					} else {
						elem := defType.getElem(i)
						err = elem.err(valueErrorHasNonSupportedValue)
					}
				}
			default:
				if isString {
					err = defType.err(valueErrorHasNonSupportedValue)
				} else {
					elem := defType.getElem(i)
					err = elem.err(valueErrorHasNonSupportedValue)
				}
			}
		}
		result = bwset.FromSliceOfStrings(ss)
	}
	return
}

func GetValidVal(val, def value, skipDefault ...bool) (result interface{}, err error) {
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
			if valErr, ok := err.(*value); ok && valErr.error.errorType == valueErrorHasNoKey {
				err = nil
			} else {
				return nil, err
			}
		} else if defaultVal.value == nil {
			return nil, defaultVal.err(valueErrorHasNonSupportedValue)
		} else if defaultVal.value, err = GetValidVal(defaultVal, def, true); err != nil {
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
			return nil, val.err(valueErrorIsNotOfType, expectedTypes)
		} else if defaultVal.value != nil {
			return defaultVal.value, nil
		} else if expectedTypes.Has("map") {
			val.value = map[string]interface{}{}
		} else {
			return nil, val.err(valueErrorIsNotOfType, expectedTypes)
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
		}
	}
	if len(typeName) == 0 {
		bwerror.Panic("HERE")
		return nil, val.err(valueErrorIsNotOfType, expectedTypes)
	}
	var validDefKeys bwset.Strings
	if !isSimpleDef {
		validDefKeys = bwset.FromSliceOfStrings([]string{ "type", "default" })
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
					return nil, val.err(valueErrorHasUnexpectedKeys, unexpectedKeys)
				}
				for defKeysKey, _ := range defKeys.mustBeMap() {
					var defKeysKeyVal value
					if defKeysKeyVal, err = defKeys.getKey(defKeysKey); err != nil {
						return nil, err
					}
					var valKeyVal value
					if valKeyVal, err = val.getKey(defKeysKey); err != nil {
						if valErr, ok := err.(*value); !ok || valErr.error.errorType != valueErrorHasNoKey {
							return nil, err
						} else {
							valKeyVal = *valErr
							valKeyVal.error = nil
						}
					}
					if valKeyVal.value == nil {
						if _isOfType(defKeysKeyVal.value, "string", "[string]") {
							if _, err = val.getKey(defKeysKey); err != nil {
								if valErr, ok := err.(*value); ok && valErr.error.errorType == valueErrorHasNoKey {
									return nil, val.err(valueErrorHasNoKey, defKeysKey)
								} else {
									return nil, err
								}
							}
						} else if _isOfType(defKeysKeyVal.value, "map") {
							if _, err = defKeysKeyVal.getKey("default"); err != nil {
							}
						}
					}
					log.Printf("defKeysKey: %s, valAsMap: %s, valKeyVal: %s, defKeysKeyVal: %s\n", defKeysKey, valAsMap, bwjson.PrettyJsonOf(valKeyVal), bwjson.PrettyJsonOf(defKeysKeyVal))
					if valAsMap[defKeysKey], err = GetValidVal(valKeyVal, defKeysKeyVal); err != nil {
						return nil, err
					}



					// var defKeysKeyVal, valMapKeyValTypeVal, valMapKeyDefaultVal, valMapKeyVal value
					// if defKeysKeyVal, err = defKeys.getKey(defKeysKey, "map"); err != nil {
					// 	return nil, err
					// }

					// var defKeysKeyValValidKeys = []string{"type", "default"}
					// if valMapKeyValTypeVal, err = defKeysKeyVal.getKey("type", []string{"string", "[]string"}); err != nil {
					// 	return nil, err
					// }
					// valMapKeyDefaultVal, err = defKeysKeyVal.getKey("default", valMapKeyValTypeVal.value)
					// if err != nil {
					// 	if valErr, ok := err.(*value); ok && valErr.error.errorType == valueErrorHasNoKey {
					// 		valMapKeyDefaultVal.value = nil
					// 		err = nil
					// 	} else {
					// 		break
					// 	}
					// }
					// if valMapKeyVal, err = val.getKey(defKeysKey); err == nil {
					// 	var validVal interface{}
					// 	if validVal, err = GetValidVal(valMapKeyVal, defKeysKeyVal); err != nil {
					// 		return nil, err
					// 	}
					// 	valAsMap[defKeysKey] = validVal
					// } else if valMapKeyDefaultVal.value != nil {
					// 	valAsMap[defKeysKey] = valMapKeyDefaultVal.value
					// 	err = nil
					// }
					// if unexpectedKeys := bwmap.GetUnexpectedKeys(defKeysKeyVal.mustBeMap(), defKeysKeyValValidKeys); unexpectedKeys != nil {
					// 	return nil, defKeysKeyVal.err(valueErrorHasUnexpectedKeys, unexpectedKeys)
					// }
				}
			}

		case "bool":
		case "string":
		case "enum":
			validDefKeys.Add("enum")

		default:
			bwerror.Panic("typeName: %s", typeName)
		}
	}

	if !isSimpleDef {
		if unexpectedKeys := bwmap.GetUnexpectedKeys(def.mustBeMap(), validDefKeys); unexpectedKeys != nil {
			return nil, def.err(valueErrorHasUnexpectedKeys, unexpectedKeys)
		}
	}
	return val.value, nil
}

func GetValOfPath(val interface{}, path string) (result interface{}, valueError error) {
	result = defparse.MustParse("{ type: 'bool', default: false, some: \"thing\" }")
	return
}

func MustValOfPath(val interface{}, path string) (result interface{}) {
	var err error
	if result, err = GetValOfPath(val, path); err != nil {
		bwerror.Panic("path <ansiCmd>%s<ansi> not found in <ansiSecondaryLiteral>%s", path, bwjson.PrettyJson(val))
	}
	return
}
