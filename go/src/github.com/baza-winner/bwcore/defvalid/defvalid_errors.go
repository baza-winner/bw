package defvalid

import (
  "github.com/baza-winner/bwcore/bwerror"
  "github.com/baza-winner/bwcore/bwjson"
  "go/types"
  // "github.com/baza-winner/bwcore/bwmap"
)

type valueIsNotOfTypeError struct {
  v value
  ofType string
}

type valueError interface {
  getValue() value
  Error() string
}

func getValueError(e valueError, fmtString string, args ...interface{}) string {
  return bwerror.Error(e.getValue().where+`<ansi> (<ansiSecondaryLiteral>%s<ansi>) `+ fmtString, bwjson.PrettyJson(e.v.value), args...).Error()
}

func (e valueIsNotOfTypeError) getValue() value {
  return e.v
}

func (e valueIsNotOfTypeError) Error() string {
  return getValueError(`is not of type <ansiPrimaryLiteral>%s`, e.ofType)
}

type valueHasUnexpectedKeysError struct {
  v value
  nonExpectedKeys []string
}

func (e valueHasUnexpectedKeysError) getValue() value {
  return e.v
}

func (e valueHasUnexpectedKeysError) Error() string {
  return getValueError(`has unexpected keys <ansiSecondaryLiteral>%v`, e.nonExpectedKeys)
}

type valueHasNoKeyError struct {
  v value
  keyName string
}

func (e valueHasNoKeyError) getValue() Value {
  return e.v
}

func (e valueHasNoKeyError) Error() string {
  return getValueError(`has no key <ansiPrimaryLiteral>%s`, e.keyName)
}

