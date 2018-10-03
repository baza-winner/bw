package defvalid

import (
  // "github.com/baza-winner/bwcore/bwerror"
  // "github.com/baza-winner/bwcore/bwjson"
  // "github.com/baza-winner/bwcore/bwmap"
  "go/types"
)

type value struct {
  value interface{}
  where string
}

func (v *value) asMap() (result map[string]interface{}, err error) {
  var ok bool
  if result, ok = v.value.(map[string]interface{}); !ok {
    err = valueIsNotOfTypeError{v: v, ofType: map[string]interface{}}
  }
  return
}

func (v *value) mustBeMap() (result map[string]interface{}) {
  if result, err = v.asMap(); err != nil {
    bwerror.Panic(err.Error())
  }
  return
}

func (v *value) asString() (result string, err error) {
  var ok bool
  if asMap, ok = v.value.(map[string]interface{}); !ok {
    err = v.errorIsNot(string)
  }
  return
}

func (v *value) mustBeString() (result string) {
  if result, err = v.asString(); err != nil {
    bwerror.Panic(err.Error())
  }
  return
}

func (v *value) getKey(keyName string, ofType string, defaultValue ...interface{}) (result value, err error) {
  m := v.mustBeMap()
  var ok bool
  result.where += "." + keyName
  if result.value, ok = m[keyName]; !ok {
    if defaultValue == nil {
      err = valueHasNoKeyError(v, keyName)
    } else {
      result.value = defaultValue[0]
    }
  } else if ! v.is(ofType) {
    err = result.errorIsNot(ofType)
  }
  return
}

func (v *value) is(ofType string) (ok bool) {
  switch ofType {
  case "string":
    _, ok = v.value.(string)
  case "map":
    _, ok = v.value.(map[string]interface{})
  default:
    bwerror.Panic("unsupported type <ansiPrimaryLiteral>%s", ofType)
  }
}

// func (v *value) getStringKey(keyName string) (result value, err error) {
//   var kv value
//   if kv, err = v.getKey(keyName); err == nil {
//     result, err = kv.asString()
//   }
//   return
// }

// func (v *value) getMapKey(keyName string, defaultValue ...map[string]interface{}) (result value, err error) {
//   var kv value
//   if kv, err = v.getKey(keyName); err == nil {
//     result, err = kv.asMap()
//   } else if defaultValue != nil {
//     result = defaultValue[0]
//     err = nil
//   }
//   return
// }
