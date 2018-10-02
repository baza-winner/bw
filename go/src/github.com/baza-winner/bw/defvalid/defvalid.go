package defvalid

import (
	// "bufio"
	// "encoding/json"
	// "errors"
	// "fmt"
	"github.com/baza-winner/bw/core"
	// "encoding/json"
	// "github.com/iancoleman/strcase"
	// "github.com/jimlawless/whereami"
	// "log"
	// "os"
	// "os/exec"
	// "strings"
	// "syscall"
)

func GetValidVal(whereVal string, val interface{}, def map[string]interface{}, whereDef string) (result interface{}, err error) {
	var defType string
	var ok bool
	if defType, err = GetStringKey(whereDef, def, `type`); err == nil {
		switch defType {
		case "map":
			var valMap map[string]interface{}
			if valMap, ok = val.(map[string]interface{}); !ok {
				err = core.Error(whereVal+`<ansi> (<ansiSecondaryLiteral>%v<ansi>) is not of type <ansiPrimaryLiteral>%s`, core.PrettyJson(val), `map`)
			} else {
				var defKeys map[string]interface{}
				if defKeys, err = GetMapKey(whereDef, def, `keys`, true); defKeys != nil && err == nil {
					for key := range valMap {
						if _, ok := defKeys[key]; !ok {
							err = core.Error(whereVal+`<ansi> (<ansiSecondaryLiteral>%s<ansi>) has unexpected key <ansiPrimaryLiteral>%s`, core.PrettyJson(val), key)
							return
						}
					}
					for defKeysKey, _ := range defKeys {
						var defKeysKeyValue map[string]interface{}
						if defKeysKeyValue, err = GetMapKey(whereDef+".keys", defKeys, defKeysKey); err == nil {
							core.Noop(defKeysKeyValue)
						}
					}
				}
			}
		default:
			err = core.Error(whereDef+`<ansiCmd>.type<ansi> has non supported value <ansiPrimaryLiteral>%s`, defType)
		}
	}
	return val, err
}

func GetKey(where string, m map[string]interface{}, keyName string, returnNilIfKeyNotExists ...bool) (result interface{}, err error) {
	if m == nil {
		err = core.Error(where + `<ansi> is <ansiPrimaryLiteral>nil`)
	} else {
		var ok bool
		if result, ok = m[keyName]; !ok {
			if returnNilIfKeyNotExists == nil || !returnNilIfKeyNotExists[0] {
				err = core.Error(where+`<ansi> has no key <ansiPrimaryLiteral>%s`, keyName)
			}
		}
	}
	return
}

func errorKeyValueIsNot(typeName string, where string, keyName string, keyValue interface{}) error {
	return core.Error(where+`<ansiCmd>.%s<ansi> (<ansiSecondaryLiteral>%+v<ansi>) is not <ansiPrimaryLiteral>%s`, keyName, keyValue, typeName)
}

func GetStringKey(where string, m map[string]interface{}, keyName string, returnNilIfKeyNotExists ...bool) (result string, err error) {
	var val interface{}
	if val, err = GetKey(where, m, keyName, returnNilIfKeyNotExists...); err == nil {
		if typedVal, ok := val.(string); ok {
			result = typedVal
		} else {
			err = errorKeyValueIsNot("string", where, keyName, val)
		}
	}
	return
}

func GetMapKey(where string, m map[string]interface{}, keyName string, returnNilIfKeyNotExists ...bool) (result map[string]interface{}, err error) {
	var val interface{}
	if val, err = GetKey(where, m, keyName, returnNilIfKeyNotExists...); err == nil {
		if typedVal, ok := val.(map[string]interface{}); ok {
			result = typedVal
		} else {
			err = errorKeyValueIsNot("map", where, keyName, val)
		}
	}
	return
}
