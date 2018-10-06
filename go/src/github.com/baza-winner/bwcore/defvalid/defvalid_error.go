package defvalid

import (
	"fmt"
	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/jimlawless/whereami"
	"sort"
)

// type defErrorType uint16

// const (
// 	defError_below_ defErrorType = iota
// 	defErrorIsNil
// 	// defErrorIsNotOfType
// 	// defErrorHasUnexpectedKeys
// 	// defErrorHasNoKey
// 	// defErrorHasNonSupportedValue
// 	// defErrorValuesCannotBeCombined
// 	defError_above_
// )

// //go:generate stringer -type=defErrorType

// type DefError

// func (v value) defError(errorType defErrorType, args ...interface{}) error {
// 	if !(defError_below_ < errorType && errorType < defError_above_) {
// 		bwerror.Panic(v.String()+" errorType == %s", errorType)
// 	}
// 	var fmtString string
// 	fmtString, args = defErrorValidators[errorType](v, args...)
// 	v.error = &defError{errorType, fmtString, args, whereami.WhereAmI(2)}
// 	return v
// }

// type defErrorValidator func(v value, args ...interface{}) (string, []interface{})

// var defErrorValidators = map[defErrorType]defErrorValidator{
// 	defErrorIsNil: _defErrorIsNil,
// 	// defErrorIsNotOfType:            _defErrorIsNotOfType,
// 	// defErrorHasUnexpectedKeys:      _defErrorHasUnexpectedKeys,
// 	// defErrorHasNoKey:               _defErrorHasNoKey,
// 	// defErrorHasNonSupportedValue:   _defErrorHasNonSupportedValue,
// 	// defErrorValuesCannotBeCombined: _defErrorValuesCannotBeCombined,
// }

// //=============================================================================

type valueErrorType uint16

const (
	valueError_below_ valueErrorType = iota
	valueErrorIsNotOfType
	valueErrorHasUnexpectedKeys
	valueErrorHasNoKey
	valueErrorHasNonSupportedValue
	valueErrorValuesCannotBeCombined
	valueError_above_
)

//go:generate stringer -type=valueErrorType

type valueErrorValidator func(v value, args ...interface{}) (string, []interface{})

var valueErrorValidators = map[valueErrorType]valueErrorValidator{
	valueErrorIsNotOfType:            _valueErrorIsNotOfType,
	valueErrorHasUnexpectedKeys:      _valueErrorHasUnexpectedKeys,
	valueErrorHasNoKey:               _valueErrorHasNoKey,
	valueErrorHasNonSupportedValue:   _valueErrorHasNonSupportedValue,
	valueErrorValuesCannotBeCombined: _valueErrorValuesCannotBeCombined,
}

func _valueErrorIsNotOfType(v value, args ...interface{}) (string, []interface{}) {
	if args == nil {
		bwerror.Panic("expects at least one arg instead of <ansiSecondaryLiteral>%#v", args)
	}
	var expectedTypes = bwset.Strings{}
	for _, i := range args {
		if _isOfType(i, "string") {
			expectedTypes.Add(_mustBeString(i))
		} else if _isOfType(i, "[]string") {
			ss := _mustBeSliceOfStrings(i)
			expectedTypes.Add(ss...)
		} else if _isOfType(i, "bwset.Strings") {
			ss := _mustBeSetOfStrings(i).ToSliceOfStrings()
			expectedTypes.Add(ss...)
		}
	}
	var result string
	for _, s := range expectedTypes.ToSliceOfStrings() {
		if len(result) > 0 {
			result += "<ansi> or <ansiPrimaryLiteral>"
		}
		result += s
	}
	return `is not of type <ansiPrimaryLiteral>%s`, []interface{}{result}
}

func _valueErrorHasUnexpectedKeys(v value, args ...interface{}) (string, []interface{}) {
	if args == nil || len(args) != 1 {
		bwerror.Panic("expects 1 arg instead of <ansiSecondaryLiteral>%#v", args)
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
}

func _valueErrorHasNoKey(v value, args ...interface{}) (string, []interface{}) {
	if args == nil || len(args) != 1 {
		bwerror.Panic("expects 1 arg instead of <ansiSecondaryLiteral>%#v", args)
	}
	_ = _mustBeString(args[0])
	return `has no key <ansiPrimaryLiteral>%s`, args
}

func _valueErrorHasNonSupportedValue(v value, args ...interface{}) (string, []interface{}) {
	if args != nil {
		bwerror.Panic("does not expect args instead of <ansiSecondaryLiteral>%#v", args)
	}
	return `has non supported value`, nil
}

func _valueErrorValuesCannotBeCombined(v value, args ...interface{}) (string, []interface{}) {
	if args == nil || len(args) < 2 {
		bwerror.Panic("expects at least 2 arg instead of <ansiSecondaryLiteral>%#v", args)
	}
	return `following values can not be combined: <ansiSecondaryLiteral>%s`, []interface{}{bwjson.PrettyJson(args)}
}

func valueErrorValidatorsCheck() {
	valueErrorType := valueError_below_ + 1
	for valueErrorType < valueError_above_ {
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
	val       value
	errorType valueErrorType
	fmtString string
	args      []interface{}
	where     string
}

func (v valueError) GetDataForJson() interface{} {
	result := map[string]interface{}{}
	result["errorType"] = v.errorType.String()
	result["fmtString"] = v.fmtString
	result["args"] = v.args
	return result
}

func (err valueError) Error() (result string) {
	// if v.error != nil {
	result = ansi.Ansi("Err", "ERR: "+fmt.Sprintf(err.val.String()+` `+err.fmtString, err.args...))
	// }
	return
}

func (v valueError) WhereError() (result string) {
	// if v.error != nil {
	result = v.where
	// }
	return
}

func getValueErr(v value, errorType valueErrorType, args ...interface{}) (result valueError) {
	if !(valueError_below_ < errorType && errorType < valueError_above_) {
		bwerror.Panic(v.String()+" errorType == %s", errorType)
	}
	var fmtString string
	fmtString, args = valueErrorValidators[errorType](v, args...)
	result = valueError{v, errorType, fmtString, args, whereami.WhereAmI(2)}
	// v.error = &valueError{errorType, fmtString, args, whereami.WhereAmI(2)}
	return
}

// func (v value) err(errorType valueErrorType, args ...interface{}) error {
// 	if !(valueError_below_ < errorType && errorType < valueError_above_) {
// 		bwerror.Panic(v.String()+" errorType == %s", errorType)
// 	}
// 	var fmtString string
// 	fmtString, args = valueErrorValidators[errorType](v, args...)
// 	v.error = &valueError{errorType, fmtString, args, whereami.WhereAmI(2)}
// 	return v
// }
