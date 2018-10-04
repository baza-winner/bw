package defvalid

import (
	"fmt"
	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"sort"
  // "log"
)

type valueErrorType uint16

const (
	_valueErrorBelow valueErrorType = iota
	valueErrorIsNotOfTypes
	valueErrorHasUnexpectedKeys
	valueErrorHasNoKey
	valueErrorHasNonSupportedValue
	_valueErrorAbove
)

//go:generate stringer -type=valueErrorType

type valueErrorValidator func(v *value, args ...interface{}) (string, []interface{})

var valueErrorValidators = map[valueErrorType]valueErrorValidator{
	valueErrorIsNotOfTypes: func(v *value, args ...interface{}) (string, []interface{}) {
		if args == nil {
			bwerror.Panic("expects at least one arg instead of <ansiSecondaryLiteral>%v", args)
		}
		var result string
		for _, i := range args {
			s := _mustBeString(i)
			if len(result) > 0 {
				result += "<ansi> or <ansiPrimaryLiteral>"
			}
			result += s
		}
		return `is not of type <ansiPrimaryLiteral>%s`, []interface{}{result}
	},
	valueErrorHasUnexpectedKeys: func(v *value, args ...interface{}) (string, []interface{}) {
		if args == nil || len(args) != 1 {
			bwerror.Panic("expects 1 arg instead of <ansiSecondaryLiteral>%v", args)
		}
		var fmtString string
		unexpectedKeys := _mustBeSliceOfStrings(args[0])
		switch {
		case len(unexpectedKeys) == 0:
			bwerror.Panic("expects non empty slice as <ansiOutline>unexpectedKeys")
		case len(unexpectedKeys) == 1:
			fmtString = `has unexpected key <ansiPrimaryLiteral>%s`
			args = []interface{}{unexpectedKeys[0]}
		default:
			sort.Strings(unexpectedKeys)
			fmtString = `has unexpected keys <ansiSecondaryLiteral>%s`
			args = []interface{}{bwjson.PrettyJson(unexpectedKeys)}
		}
		return fmtString, args
	},
	valueErrorHasNoKey: func(v *value, args ...interface{}) (string, []interface{}) {
		if args == nil || len(args) != 1 {
			bwerror.Panic("expects 1 arg instead of <ansiSecondaryLiteral>%v", args)
		}
		_ = _mustBeString(args[0])
		return `has no key <ansiPrimaryLiteral>%s`, args
	},
	valueErrorHasNonSupportedValue: func(v *value, args ...interface{}) (string, []interface{}) {
		if args != nil {
			bwerror.Panic("does not expect args instead of <ansiSecondaryLiteral>%v", args)
		}
		return `has non supported value <ansiPrimaryLiteral>%s`, []interface{}{v.value}
	},
}

func valueErrorValidatorsCheck() {
	valueErrorType := _valueErrorBelow + 1
	for valueErrorType < _valueErrorAbove {
		if _, ok := valueErrorValidators[valueErrorType]; !ok {
			bwerror.Panic("not defined <ansiOutline>valueErrorValidators<ansi>[<ansiPrimaryLiteral>%s<ansi>]", valueErrorType)
		}
		valueErrorType += 1
	}
}

func init() {
	valueErrorValidatorsCheck()
}

type valueError struct {
	errorType valueErrorType
	fmtString string
	args      []interface{}
}

func (v value) Error() string {
	if v.error == nil {
		bwerror.Panic(v.String() + " v.error == nil")
	}
  // log.Printf("%s", v)
  // log.Println(v)
	return ansi.Ansi("Err", "ERR: "+fmt.Sprintf(v.String()+` `+v.error.fmtString, v.error.args...))
}

func (v *value) err(errorType valueErrorType, args ...interface{}) error {
	if !(_valueErrorBelow < errorType && errorType < _valueErrorAbove) {
		bwerror.Panic(v.String()+" errorType == %s", errorType)
	}
	var fmtString string
	fmtString, args = valueErrorValidators[errorType](v, args...)
	v.error = &valueError{errorType, fmtString, args}
	return v
}
