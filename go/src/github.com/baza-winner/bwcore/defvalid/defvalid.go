/*
Предоставляет функции для валидации interface{}.
*/
package defvalid

import (
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/defparse"
)

func GetValidVal(val, def value) (result interface{}, err error) {
	var defType value
	if defType, err = def.getKey(`type`, `string`); err == nil {
		switch defType.mustBeString() {
		case "map":
			var valAsMap map[string]interface{}
			if valAsMap, err = val.asMap(); err == nil {
				var defKeys value
				if defKeys, err = def.getKey(`keys`, `map`, nil); defKeys.value != nil && err == nil {
					if unexpectedKeys := bwmap.GetUnexpectedKeys(valAsMap, defKeys.mustBeMap()); unexpectedKeys != nil {
						err = val.err(valueErrorHasUnexpectedKeys, unexpectedKeys)
					} else {
						for defKeysKey, _ := range defKeys.mustBeMap() {
							var defKeysKeyVal, valMapKeyValTypeVal, valMapKeyDefaultVal, valMapKeyVal value
							if defKeysKeyVal, err = defKeys.getKey(defKeysKey, `map`); err == nil {
								var defKeysKeyValValidKeys = []string{"type", "default"}
								if valMapKeyValTypeVal, err = defKeysKeyVal.getKey(`type`, []string{`string`, `[]string`}); err != nil {
									break
								} else {
									valMapKeyDefaultVal, err = defKeysKeyVal.getKey(`default`, valMapKeyValTypeVal.value)
									if err != nil {
										if valErr, ok := err.(*value); ok && valErr.error.errorType == valueErrorHasNoKey {
											valMapKeyDefaultVal.value = nil
											err = nil
										} else {
											break
										}
									}
									if valMapKeyVal, err = val.getKey(defKeysKey); err == nil {
										if valAsMap[defKeysKey], err = GetValidVal(valMapKeyVal, defKeysKeyVal); err != nil {
											break
										}
									} else if valMapKeyDefaultVal.value != nil {
										valAsMap[defKeysKey] = valMapKeyDefaultVal.value
										err = nil
									}
									if unexpectedKeys := bwmap.GetUnexpectedKeys(defKeysKeyVal.mustBeMap(), defKeysKeyValValidKeys); unexpectedKeys != nil {
										err = defKeysKeyVal.err(valueErrorHasUnexpectedKeys, unexpectedKeys)
										break
									}
								}
							}
						}
					}
				}
			}

		case "bool":
			_, err = val.asBool()

		case "string":
			_, err = val.asString()

		default:
			err = defType.err(valueErrorHasNonSupportedValue)
		}
	}
	return val.value, err
}

func GetValOfPath(val interface{}, path string) (result interface{}, valueError error) {
	result = defparse.MustParse(`{ type: 'bool', default: false, some: "thing" }`)
	return
}

func MustValOfPath(val interface{}, path string) (result interface{}) {
	var err error
	if result, err = GetValOfPath(val, path); err != nil {
		bwerror.Panic("path <ansiCmd>%s<ansi> not found in <ansiSecondaryLiteral>%s", path, bwjson.PrettyJson(val))
	}
	return
}
