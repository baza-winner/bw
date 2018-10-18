// Code generated by "stringer -type=UnicodeCategory,ErrorType,ruleKind"; DO NOT EDIT.

package pfa

import "strconv"

const _UnicodeCategory_name = "UnicodeSpaceUnicodeLetterUnicodeDigitUnicodeOpenBracesUnicodePunctUnicodeSymbol"

var _UnicodeCategory_index = [...]uint8{0, 12, 25, 37, 54, 66, 79}

func (i UnicodeCategory) String() string {
	if i >= UnicodeCategory(len(_UnicodeCategory_index)-1) {
		return "UnicodeCategory(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _UnicodeCategory_name[_UnicodeCategory_index[i]:_UnicodeCategory_index[i+1]]
}

const _ErrorType_name = "pfaErrorBelowUnexpectedRuneFailedToGetNumberUnknownWordErrorpfaErrorAbove"

var _ErrorType_index = [...]uint8{0, 15, 29, 46, 62, 77}

func (i ErrorType) String() string {
	if i >= ErrorType(len(_ErrorType_index)-1) {
		return "ErrorType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ErrorType_name[_ErrorType_index[i]:_ErrorType_index[i+1]]
}

const _ruleKind_name = "ruleNormalruleDefaultruleEof"

var _ruleKind_index = [...]uint8{0, 10, 21, 28}

func (i ruleKind) String() string {
	if i >= ruleKind(len(_ruleKind_index)-1) {
		return "ruleKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ruleKind_name[_ruleKind_index[i]:_ruleKind_index[i+1]]
}
