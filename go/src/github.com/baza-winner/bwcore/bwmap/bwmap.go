/*
Предоставялет функции для работы со структурой map[string]interface{}.
*/
package bwmap

import (
	"reflect"

	"github.com/baza-winner/bwcore/bwerror"
)

func GetUnexpectedKeys(m interface{}, expected ...interface{}) []string {
	if expected == nil {
		return nil
	}
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		bwerror.Panic("<ansiOutline>m<ansi> (<ansiSecondary>%+v<ansi>) is not <ansiPrimary>map", m)
	}
	for _, vk := range v.MapKeys() {
		if vk.Kind() != reflect.String {
			bwerror.Panic("<ansiOutline>m<ansi> (<ansiSecondary>%+v<ansi>) is not <ansiPrimaryLiteralmap[string]", m)
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
			bwerror.Panic("<ansiOutline>expected<ansi> (<ansiPrimary>%+v<ansi>) neither <ansiSecondary>string<ansi>, nor <ansiSecondary>[]string<ansi>, nor <ansiSecondary>map[string]interface", expected)
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
