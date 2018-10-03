/*
Предоставялет функции для работы со структурой map[string]interface{}.
*/
package bwmap

import (
  "github.com/baza-winner/bwcore/bwerror"
  "reflect"
)

func GetUnexpectedKeys(m interface{}, expected interface{}) []string {
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
  var unexpectedKeys = []string{}
  if keyName, ok := expected.(string); ok {
    for _, vk := range v.MapKeys() {
      k := vk.String()
      if k != keyName {
        unexpectedKeys = append(unexpectedKeys, k)
      }
    }
  } else if keyNames, ok := expected.([]string); ok {
    keyNameMap := map[string]struct{}{}
    for _, k := range keyNames {
      keyNameMap[k] = struct{}{}
    }
    for _, vk := range v.MapKeys() {
      k := vk.String()
      if _, ok := keyNameMap[k]; !ok {
        unexpectedKeys = append(unexpectedKeys, k)
      }
    }
  } else if keyNameMap, ok := expected.(map[string]interface{}); ok {
    for _, vk := range v.MapKeys() {
      k := vk.String()
      if _, ok := keyNameMap[k]; !ok {
        unexpectedKeys = append(unexpectedKeys, k)
      }
    }
  } else {
    bwerror.Panic("<ansiOutline>expected<ansi> (<ansiPrimaryLiteral>%+v<ansi>) neither <ansiSecondaryLiteral>string<ansi>, nor <ansiSecondaryLiteral>[]string<ansi>, nor <ansiSecondaryLiteral>map[string]interface", expected)
  }
  if len(unexpectedKeys) == 0 {
    return nil
  } else {
    return unexpectedKeys
  }
}

func CropMap(m interface{}, crop interface{}) {
  if unexpectedKeys := GetUnexpectedKeys(m, crop); unexpectedKeys != nil {
    for _, k := range unexpectedKeys {
      v := reflect.ValueOf(m)
      v.SetMapIndex(reflect.ValueOf(k), reflect.Value{})
    }
  }
}

