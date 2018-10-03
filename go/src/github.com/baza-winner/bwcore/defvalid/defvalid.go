/*
Предоставляет функции для валидации interface{}.
*/
package defvalid

import (
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
)

func GetValidVal(whereVal string, val interface{}, def map[string]interface{}, whereDef string) (result interface{}, err error) {
	var defType string
	var ok bool
	if defType, err = getStringKey(whereDef, def, `type`); err == nil {
		switch defType {
		case "map":
			var valMap map[string]interface{}
			if valMap, ok = val.(map[string]interface{}); !ok {
				err = bwerror.Error(whereVal+`<ansi> (<ansiSecondaryLiteral>%v<ansi>) is not of type <ansiPrimaryLiteral>%s`, bwjson.PrettyJson(val), `map`)
			} else {
				var defKeys map[string]interface{}
				if defKeys, err = getMapKey(whereDef, def, `keys`, true); defKeys != nil && err == nil {
					if nonExpectedKeys := bwmap.GetUnexpectedKeys(valMap, defKeys); nonExpectedKeys != nil {
						err = bwerror.Error(whereVal+`<ansi> (<ansiSecondaryLiteral>%s<ansi>) has unexpected key(s) <ansiPrimaryLiteral>%v`, bwjson.PrettyJson(val), nonExpectedKeys)
					}
					for defKeysKey, _ := range defKeys {
						var defKeysKeyValue map[string]interface{}
						if defKeysKeyValue, err = getMapKey(whereDef+".keys", defKeys, defKeysKey); err == nil {
							bwerror.Noop(defKeysKeyValue)
						}
					}
				}
			}
		default:
			err = bwerror.Error(whereDef+`.type<ansi> has non supported value <ansiPrimaryLiteral>%s`, defType)
		}
	}
	return val, err
}

func getKey(where string, m map[string]interface{}, keyName string, returnNilIfKeyNotExists ...bool) (result interface{}, err error) {
	if m == nil {
		err = bwerror.Error(where + `<ansi> is <ansiPrimaryLiteral>nil`)
	} else {
		var ok bool
		if result, ok = m[keyName]; !ok {
			if returnNilIfKeyNotExists == nil || !returnNilIfKeyNotExists[0] {
				err = bwerror.Error(where+`<ansi> has no key <ansiPrimaryLiteral>%s`, keyName)
			}
		}
	}
	return
}

func errorKeyValueIsNot(typeName string, where string, keyName string, keyValue interface{}) error {
	return bwerror.Error(where+`.%s<ansi> (<ansiSecondaryLiteral>%+v<ansi>) is not <ansiPrimaryLiteral>%s`, keyName, keyValue, typeName)
}

func getStringKey(where string, m map[string]interface{}, keyName string, returnNilIfKeyNotExists ...bool) (result string, err error) {
	var val interface{}
	if val, err = getKey(where, m, keyName, returnNilIfKeyNotExists...); err == nil {
		if typedVal, ok := val.(string); ok {
			result = typedVal
		} else {
			err = errorKeyValueIsNot("string", where, keyName, val)
		}
	}
	return
}

func getMapKey(where string, m map[string]interface{}, keyName string, returnNilIfKeyNotExists ...bool) (result map[string]interface{}, err error) {
	var val interface{}
	if val, err = getKey(where, m, keyName, returnNilIfKeyNotExists...); err == nil {
		if typedVal, ok := val.(map[string]interface{}); ok {
			result = typedVal
		} else {
			err = errorKeyValueIsNot("map", where, keyName, val)
		}
	}
	return
}
