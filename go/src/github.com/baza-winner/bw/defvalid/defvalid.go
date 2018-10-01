package defvalid

import (
	// "bufio"
	// "encoding/json"
	// "errors"
	// "fmt"
	"github.com/baza-winner/bw/core"
  "encoding/json"
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
		if defType == `map` {
			var valMap map[string]interface{}
			if valMap, ok = val.(map[string]interface{}); ok {
				var defKeys map[string]interface{}
				if defKeys, err = GetMapKeyIfExists(whereDef, def, `keys`); defKeys != nil && err == nil {
					for key := range valMap {
						if _, ok := defKeys[key]; !ok {
              valJson, _ := json.MarshalIndent(val, "", "  ")
							err = core.Error(`<ansiOutline>%s<ansi> (<ansiSecondaryLiteral>%s<ansi>) has unexpected key <ansiPrimaryLiteral>%s`, whereVal, valJson, key)
							return
						}
					}
				}
			} else {
				err = core.Error(`<ansiOutline>%s<ansi> (<ansiSecondaryLiteral>%v<ansi>) is not of type <ansiPrimaryLiteral>%s`, whereVal, val, `map`)
			}
		} else {
			err = core.Error(`<ansiOutline>%s<ansi>[<ansiSecondaryLiteral>%s<ansi>] has non supported value <ansiPrimaryLiteral>%s`, whereDef, `type`, defType)
		}
	}
	return val, err
}

func GetStringKey(where string, m map[string]interface{}, keyName string) (result string, err error) {
	if m != nil {
		if val, ok := m[keyName]; ok {
			if typedVal, ok := val.(string); ok {
				result = typedVal
			} else {
				err = core.Error(`<ansiOutline>%s<ansi>[<ansiSecondaryLiteral>%s<ansi>] (<ansiSecondaryLiteral>%+v<ansi>) is not <ansiPrimaryLiteral>%s`, where, keyName, val, `string`)
			}
		} else {
			err = core.Error(`<ansiOutline>%s<ansi> has not key <ansiPrimaryLiteral>%s`, where, keyName)
		}
	} else {
		err = core.Error(`<ansiOutline>%s<ansi> is not <ansiPrimaryLiteral>map`, where)
	}
	return
}

func GetMapKeyIfExists(where string, m map[string]interface{}, keyName string) (result map[string]interface{}, err error) {
	if m != nil {
		if val, ok := m[keyName]; ok {
			if typedVal, ok := val.(map[string]interface{}); ok {
				result = typedVal
			} else {
				err = core.Error(`<ansiOutline>%s<ansi>[<ansiSecondaryLiteral>%s<ansi>] (<ansiSecondaryLiteral>%+v<ansi>) is not <ansiPrimaryLiteral>%s`, where, keyName, val, `map`)
			}
		} else {
			result = nil
			err = nil
		}
	} else {
		err = core.Error(`<ansiOutline>%s<ansi> is not <ansiPrimaryLiteral>map`, where)
	}
	return
}
