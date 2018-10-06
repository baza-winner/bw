// Предоставляет функции для валидации interface{}.
package defvalid

import (
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/defparse"
	// "log"
	"reflect"
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
		err = valueErrMake(defType, valueErrorIsNotOfType, "string", "[]string")
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
						err = valueErrMake(defType, valueErrorHasNonSupportedValue)
					} else {
						elem := defType.getElem(i)
						err = valueErrMake(elem, valueErrorHasNonSupportedValue)
					}
				}
			default:
				if isString {
					err = valueErrMake(defType, valueErrorHasNonSupportedValue)
				} else {
					elem := defType.getElem(i)
					err = valueErrMake(elem, valueErrorHasNonSupportedValue)
				}
			}
		}
		result = bwset.FromSliceOfStrings(ss)
		if result.Has("enum") && result.Has("string") {
			// log.Printf("%s\n", result)
			err = valueErrMake(defType, valueErrorValuesCannotBeCombined, "enum", "string")
		}
	}
	return
}

// type DeftypeType uint16

// const (
// 	deftype_below_ DeftypeType = iota
// 	deftypeBool
// 	deftypeString
// 	deftypeInt
// 	deftypeNumber
// 	deftypeMap
// 	deftypeArray
// 	deftype_above_
// )

// type Def struct {
// 	deftype   []DefTypeType
// 	enum      bwset.Strings
// 	minInt    int
// 	maxInt    int
// 	minNumber float64
// 	maxNumber float64
// 	keys      map[string]Def
// 	elem      Def
// 	arrayItem Def
// }

// func CompileDef(def value) (result Def, err error) {
// 	if def.value == nil {

// 	}

// }

func ValidateVal(what string, val, def interface{}) (result interface{}, err error) {
	return getValidVal(
		value{
			value: val,
			what:  "<ansiOutline>" + what + "<ansiCmd>",
		},
		value{
			value: def,
			what:  "<ansiOutline>" + what + "::def<ansiCmd>",
		},
	)
}

func MustValidVal(what string, val, def interface{}) (result interface{}) {
	var err error
	if result, err = ValidateVal(what, val, def); err != nil {
		bwerror.PanicErr(err)
	}
	return
}

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
			return nil, valueErrMake(defaultVal, valueErrorHasNonSupportedValue)
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
			return nil, valueErrMake(val, valueErrorIsNotOfType, expectedTypes)
		} else if defaultVal.value != nil {
			return defaultVal.value, nil
		} else if expectedTypes.Has("map") {
			val.value = map[string]interface{}{}
		} else {
			return nil, valueErrMake(val, valueErrorIsNotOfType, expectedTypes)
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
		return nil, valueErrMake(val, valueErrorIsNotOfType, expectedTypes)
	}
	var validDefKeys bwset.Strings
	if !isSimpleDef {
		validDefKeys = bwset.FromSliceOfStrings([]string{"type", "default"})
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
					return nil, valueErrMake(val, valueErrorHasUnexpectedKeys, unexpectedKeys)
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
									return nil, valueErrMake(val, valueErrorHasNoKey, defKeysKey)
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
			enumSet := bwset.FromSliceOfStrings(_mustBeSliceOfStrings(enumValues.value))
			if !enumSet.Has(_mustBeString(val.value)) {
				return nil, valueErrMake(val, valueErrorHasNonSupportedValue)
			}

		default:
			bwerror.Panic("typeName: %s", typeName)
		}
	}

	if !isSimpleDef {
		if unexpectedKeys := bwmap.GetUnexpectedKeys(def.mustBeMap(), validDefKeys); unexpectedKeys != nil {
			return nil, valueErrMake(def, valueErrorHasUnexpectedKeys, unexpectedKeys)
		}
	}
	return val.value, nil
}

func GetValOfPath(val interface{}, path string) (result interface{}, valueError error) {
	switch path {
	case ".keys.keyOne":
		result = defparse.MustParse("{ type: 'bool', default: false, some: \"thing\" }")
	case ".type":
		result = defparse.MustParse("['enum', 'string']")
	case ".enum":
		result = defparse.MustParse("['one', true, 3 ]")
	}
	return
}

func MustValOfPath(val interface{}, path string) (result interface{}) {
	var err error
	if result, err = GetValOfPath(val, path); err != nil {
		bwerror.Panic("path <ansiCmd>%s<ansi> not found in <ansiSecondaryLiteral>%s", path, bwjson.PrettyJson(val))
	}
	return
}
