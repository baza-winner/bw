/*
Предоставялет функции для работы со структурой map[string]interface{}.
*/
package bwmap

import (
	"github.com/baza-winner/bwcore/bwerror"
	"reflect"
)

func GetUnexpectedKeys(m interface{}, expected ...interface{}) []string {
	if expected == nil {
		return nil
	}
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		bwerror.Panic("<ansiOutline>m<ansi> (<ansiSecondaryLiteral>%+v<ansi>) is not <ansiPrimaryLiteral>map", m)
	}
	for _, vk := range v.MapKeys() {
		if vk.Kind() != reflect.String {
			bwerror.Panic("<ansiOutline>m<ansi> (<ansiSecondaryLiteral>%+v<ansi>) is not <ansiPrimaryLiteralmap[string]", m)
		}
		break
	}
	expectedKeys := map[string]struct{}{}
	for _, item := range expected {
		if item == nil {
			continue
		}
		if s, ok := item.(string); ok {
			expectedKeys[s] = struct{}{}
		} else if ss, ok := item.([]string); ok {
			for _, k := range ss {
				expectedKeys[k] = struct{}{}
			}
		} else if reflect.TypeOf(item).Kind() == reflect.Map {
			v := reflect.ValueOf(item)
			for _, vk := range v.MapKeys() {
				k := vk.String()
				expectedKeys[k] = struct{}{}
			}
		} else {
			bwerror.Panic("<ansiOutline>expected<ansi> (<ansiPrimaryLiteral>%+v<ansi>) neither <ansiSecondaryLiteral>string<ansi>, nor <ansiSecondaryLiteral>[]string<ansi>, nor <ansiSecondaryLiteral>map[string]interface", expected)
		}
	}
	var unexpectedKeys = []string{}
	for _, vk := range v.MapKeys() {
		k := vk.String()
		if _, ok := expectedKeys[k]; !ok {
			unexpectedKeys = append(unexpectedKeys, k)
		}
	}
	if len(unexpectedKeys) == 0 {
		return nil
	} else {
		return unexpectedKeys
	}
}

func CropMap(m interface{}, crop ...interface{}) {
	if unexpectedKeys := GetUnexpectedKeys(m, crop...); unexpectedKeys != nil {
		for _, k := range unexpectedKeys {
			v := reflect.ValueOf(m)
			v.SetMapIndex(reflect.ValueOf(k), reflect.Value{})
		}
	}
}
