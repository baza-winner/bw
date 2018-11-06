/*
Предоставялет функции для работы со структурой map[string]interface{}.
*/
package bwmap

import (
	"reflect"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerr"
)

var (
	ansiMustBeMap       string
	ansiMustBeMapString string
)

func init() {
	ansiMustBeMap = ansi.String("<ansiVar>m<ansi> (<ansiVal>%#v<ansi>) must be <ansiType>map")
	ansiMustBeMapString = ansi.String("<ansiVar>m<ansi> (<ansiVal>%#v<ansi>) must be <ansiType>map[string]")
}

func UnexpectedKeys(m interface{}, expected ...interface{}) (result []string, err error) {
	if expected == nil {
		return
	}
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		err = bwerr.From(ansiMustBeMap, m)
		// bwerr.Panic(bw.Fmt(ansiMustBeMap, m))
	}
	for _, vk := range v.MapKeys() {
		if vk.Kind() != reflect.String {
			err = bwerr.From(ansiMustBeMapString, m)
			// bwerr.Panic(bw.Fmt(ansiMustBeMapString, m))
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
			bwerr.Panic("<ansiVar>expected<ansi> (<ansiVal>%+v<ansi>) neither <ansiVal>string<ansi>, nor <ansiVal>[]string<ansi>, nor <ansiVal>map[string]interface", expected)
		}
	}
	result = []string{}
	for _, vk := range v.MapKeys() {
		k := vk.String()
		if _, ok := expectedKeys[k]; !ok {
			result = append(result, k)
		}
	}
	if len(result) == 0 {
		result = nil
	}
	return
}

func MustUnexpectedKeys(m interface{}, expected ...interface{}) (result []string) {
	var err error
	if result, err = UnexpectedKeys(m, expected...); err != nil {
		bwerr.PanicA(bwerr.E{Depth: 1, Error: err})
	}
	return
}

func CropMap(m interface{}, crop ...interface{}) {
	if unexpectedKeys, err := UnexpectedKeys(m, crop...); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	} else if unexpectedKeys != nil {
		for _, k := range unexpectedKeys {
			v := reflect.ValueOf(m)
			v.SetMapIndex(reflect.ValueOf(k), reflect.Value{})
		}
	}
}
