/*
Предоставляет функции для валидации interface{}.
*/
package defvalid

import (
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
)

// func getStringKey(where string, m map[string]interface{}, keyName string, defaultValue ...string) (result string, err error) {
// 	var val interface{}
// 	if val, err = getKey(where, m, keyName); err == nil {
// 		if typedVal, ok := val.(string); ok {
// 			result = typedVal
// 		} else {
// 			// err = errorKeyValueIsNot("string", where, keyName, val)
// 			err = errorValueIsNotOfType("string", val, where +"." +keyName )
// 		}
// 			} else if defaultValue != nil {
// 		result = defaultValue[0]
// 		err = nil
// 	}
// 	return
// }

func GetValidVal(val, def value) (result interface{}, err error) {
	var defType value
	var ok bool
	// if defType, err = getStringKey(whereDef, def, `type`); err == nil {
	if defType, err = def.getKey(`type`, `string`); err == nil {
		switch defType.mustBeString() {
		case "map":
			var valAsMap map[string]interface{}
			if valAsMap, err = val.asMap(); err == nil {
				var defKeys value
				if defKeys, err = def.getKey(`keys`, `map`, nil); defKeys.value != nil && err == nil {
					if nonExpectedKeys := bwmap.GetUnexpectedKeys(valAsMap, defKeys); nonExpectedKeys != nil {
						err = valueHasUnexpectedKeysError{val, nonExpectedKeys}
						// err = bwerror.Error(whereVal+`<ansi> (<ansiSecondaryLiteral>%s<ansi>) has unexpected keys <ansiSecondaryLiteral>%v`, bwjson.PrettyJson(val.value), nonExpectedKeys)
					}
					for defKeysKey, _ := range defKeys.mustBeMap() {
						var defKeysKeyValue value
						if defKeysKeyValue, err = defKeys.getKey(defKeysKey, `map`); err == nil {
	 						  if valAsMap[defKeysKey], err = GetValidVal(defKeysKeyValue,
	 						  	value{}
	 						  	whereVal+"."+defKeysKey, valMap[defKeysKey], defKeysKeyValue, whereDef+".keys."+defKeysKey); err != nil {
							 	break
 						 }
						}
					}
				}
			}
		case "bool":
			if _, ok = val.(bool); !ok {
				err = errorValueIsNotOfType(`bool`, val, whereVal)
			}
		default:
			err = bwerror.Error(whereDef+`.type<ansi> has non supported value <ansiPrimaryLiteral>%s`, defType)
		}
	}
	return val, err
}

func getKey(where string, m map[string]interface{}, keyName string) (result interface{}, err error) {
	if m == nil {
		err = bwerror.Error(where + `<ansi> is <ansiPrimaryLiteral>nil`)
	} else {
		var ok bool
		if result, ok = m[keyName]; !ok {
			err = bwerror.Error(where+`<ansi> has no key <ansiPrimaryLiteral>%s`, keyName)
		}
	}
	return
}

func errorValueIsNotOfType(typeName string, value interface{}, where string) error {
	return bwerror.Error(where+`<ansi> (<ansiSecondaryLiteral>%s<ansi>) is not of type <ansiPrimaryLiteral>%s`, bwjson.PrettyJson(value), typeName)
}

func getStringKey(where string, m map[string]interface{}, keyName string, defaultValue ...string) (result string, err error) {
	var val interface{}
	if val, err = getKey(where, m, keyName); err == nil {
		if typedVal, ok := val.(string); ok {
			result = typedVal
		} else {
			// err = errorKeyValueIsNot("string", where, keyName, val)
			err = errorValueIsNotOfType("string", val, where +"." +keyName )
		}
			} else if defaultValue != nil {
		result = defaultValue[0]
		err = nil
	}
	return
}

func getMapKey(where string, m map[string]interface{}, keyName string, defaultValue ...map[string]interface{}) (result map[string]interface{}, err error) {
	var val interface{}
	if val, err = getKey(where, m, keyName); err == nil {
		if typedVal, ok := val.(map[string]interface{}); ok {
			result = typedVal
		} else {
			err = errorValueIsNotOfType("map", val, where +"." +keyName )
		}
	} else if defaultValue != nil {
		result = defaultValue[0]
		err = nil
	}
	return
}

func getBoolKey(where string, m map[string]interface{}, keyName string, defaultValue ...bool) (result bool, err error) {
	var val interface{}
	if val, err = getKey(where, m, keyName); err == nil {
		if typedVal, ok := val.(bool); ok {
			result = typedVal
		} else {
			err = errorValueIsNotOfType("bool", val, where +"." +keyName )
		}
			} else if defaultValue != nil {
		result = defaultValue[0]
		err = nil
	}
	return
}
