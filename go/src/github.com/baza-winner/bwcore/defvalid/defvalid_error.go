package defvalid

import (
	// "fmt"
	// "github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/jimlawless/whereami"

	// "log"
	"reflect"
	"sort"
)

//=============================================================================

type valueErrorType uint16

const (
	valueErrorBelow valueErrorType = iota
	valueErrorIsNotOfType
	valueErrorHasUnexpectedKeys
	valueErrorHasNoKey
	valueErrorHasNonSupportedValue
	valueErrorValuesCannotBeCombined
	valueErrorConflictingKeys
	valueErrorArrayOf
	valueErrorOutOfRange
	valueErrorAbove
)

//go:generate stringer -type=valueErrorType

type valueError struct {
	errorType valueErrorType
	fmtString string
	fmtArgs   []interface{}
	Where     string
}

func valueErrorMake(v value, errorType valueErrorType, args ...interface{}) (result valueError) {
	if !(valueErrorBelow < errorType && errorType < valueErrorAbove) {
		bwerror.Panic(v.String()+" errorType == %s", errorType)
	}
	fmtString, fmtArgs := valueErrorValidators[errorType](v, args...)
	result = valueError{errorType, fmtString, fmtArgs, whereami.WhereAmI(2)}
	return
}

func (err valueError) Error() (result string) {
	result = bwerror.Error(err.fmtString, err.fmtArgs...).Error()
	return
}

type valueErrorValidator func(v value, args ...interface{}) (string, []interface{})

var valueErrorValidators map[valueErrorType]valueErrorValidator

func valueErrorValidatorsCheck() {
	valueErrorType := valueErrorBelow + 1
	for valueErrorType < valueErrorAbove {
		if _, ok := valueErrorValidators[valueErrorType]; !ok {
			bwerror.Panic("not defined <ansiOutline>valueErrorValidators<ansi>[<ansiPrimaryLiteral>%s<ansi>]", valueErrorType)
		}
		valueErrorType += 1
	}
}

func _valueErrorIsNotOfType(v value, args ...interface{}) (string, []interface{}) {
	if args == nil {
		bwerror.Panic("expects at least one arg instead of <ansiSecondaryLiteral>%#v", args)
	}
	var expectedTypes = bwset.StringSet{}
	for _, i := range args {
		if _isOfType(i, "string") {
			expectedTypes.Add(_mustBeString(i))
		} else if _isOfType(i, "[]string") {
			expectedTypes.Add(_mustBeSliceOfStrings(i)...)
		} else if _isOfType(i, "bwset.Strings") {
			expectedTypes.Add(_mustBeBwsetStrings(i).ToSlice()...)
		} else if _isOfType(i, "deftype.Set") {
			expectedTypes.Add(_mustBeDeftypeSet(i).ToSliceOfStrings()...)
		} else {
			bwerror.Panic("args: %#v", args)
		}
	}
	// log.Printf("args: %#v", args)
	var result string
	for _, s := range expectedTypes.ToSlice() {
		if len(result) > 0 {
			result += "<ansi>, or <ansiPrimaryLiteral>"
		}
		result += s
	}
	return v.String() + ` is not of type <ansiPrimaryLiteral>` + result, nil
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
	return v.String() + ` ` + fmtString, args
}

func _valueErrorHasNoKey(v value, args ...interface{}) (string, []interface{}) {
	if args == nil || len(args) != 1 {
		bwerror.Panic("expects 1 arg instead of <ansiSecondaryLiteral>%#v", args)
	}
	_ = _mustBeString(args[0])
	return v.String() + ` has no key <ansiPrimaryLiteral>%s`, args
}

func _valueErrorHasNonSupportedValue(v value, args ...interface{}) (string, []interface{}) {
	if args != nil {
		bwerror.Panic("does not expect args instead of <ansiSecondaryLiteral>%#v", args)
	}
	return v.String() + ` has non supported value`, nil
}

func _valueErrorValuesCannotBeCombined(v value, args ...interface{}) (string, []interface{}) {
	if args == nil || len(args) < 2 {
		bwerror.Panic("expects at least 2 arg instead of <ansiSecondaryLiteral>%#v", args)
	}
	return v.String() + ` following values can not be combined: <ansiSecondaryLiteral>%s`, []interface{}{bwjson.PrettyJson(args)}
}

func _valueErrorConflictingKeys(v value, args ...interface{}) (string, []interface{}) {
	if args == nil || len(args) != 1 {
		bwerror.Panic("expects 1 arg instead of <ansiSecondaryLiteral>%#v", args)
	}
	var ck map[string]interface{}
	var ok bool
	if ck, ok = args[0].(map[string]interface{}); !ok {
		bwerror.Panic("expects map[string]interface{} instead of <ansiSecondaryLiteral>%#v", args[0])
	}
	return v.String() + ` has conflicting keys: <ansiSecondaryLiteral>%s`, []interface{}{bwjson.PrettyJson(ck)}
}

func _valueErrorArrayOf(v value, args ...interface{}) (string, []interface{}) {
	if args != nil {
		bwerror.Panic("does not expect args instead of <ansiSecondaryLiteral>%#v", args)
	}
	return v.String() + ` must be followed by some type`, nil
}

func _valueErrorOutOfRange(v value, args ...interface{}) (fmtString string, fmtArgs []interface{}) {
	if args == nil || len(args) != 2 {
		bwerror.Panic("expects exact 2 arg instead of <ansiSecondaryLiteral>%#v", args)
	}
	minFmt, min := getFmtStringArg(args[0])
	maxFmt, max := getFmtStringArg(args[1])

	if len(minFmt) > 0 {
		if len(maxFmt) > 0 {
			fmtString = v.String() + " is out of <ansiOutline>range <ansiSecondaryLiteral>[" + minFmt + ", " + maxFmt + "]"
			fmtArgs = []interface{}{min, max}
		} else {
			fmtString = v.String() + " is less then <ansiOutline>minLimit <ansiPrimaryLiteral>" + minFmt
			fmtArgs = []interface{}{min}
		}
	} else if len(maxFmt) > 0 {
		fmtString = v.String() + " is greater then <ansiOutline>maxLimit <ansiPrimaryLiteral>" + maxFmt
		fmtArgs = []interface{}{max}
	}
	return
}

func getFmtStringArg(limit interface{}) (fmtString string, fmtArg interface{}) {
	limitValue := reflect.ValueOf(limit).Elem()
	zeroValue := reflect.Value{}
	if limitValue != zeroValue {
		fmtArg = limitValue.Interface()
		if _isOfType(fmtArg, "int64", "float64") {
			fmtString = "%v"
		} else {
			bwerror.Panic("limit %#v is expected to be int64 or float64", fmtArg)
		}
	}
	return
}
