package defvalid

import (
	"github.com/baza-winner/bwcore/defvalid/deftype"
)

func init() {
	valueErrorValidators = map[valueErrorType]valueErrorValidator{
		valueErrorIsNotOfType:            _valueErrorIsNotOfType,
		valueErrorHasUnexpectedKeys:      _valueErrorHasUnexpectedKeys,
		valueErrorHasNoKey:               _valueErrorHasNoKey,
		valueErrorHasNonSupportedValue:   _valueErrorHasNonSupportedValue,
		valueErrorValuesCannotBeCombined: _valueErrorValuesCannotBeCombined,
		valueErrorConflictingKeys:        _valueErrorConflictingKeys,
		valueErrorArrayOf:                _valueErrorArrayOf,
		valueErrorOutOfRange:             _valueErrorOutOfRange,
	}
	valueErrorValidatorsCheck()

	getValidValHelpers = map[deftype.Item]getValidValHelper{
		deftype.Bool:    _Bool,
		deftype.String:  _String,
		deftype.Int:     _Int,
		deftype.Number:  _Number,
		deftype.Map:     _Map,
		deftype.Array:   _Array,
		deftype.ArrayOf: _ArrayOf,
	}
	getValidValHelpersCheck()
}
