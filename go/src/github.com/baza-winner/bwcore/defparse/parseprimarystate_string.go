// Code generated by "stringer -type=parsePrimaryState"; DO NOT EDIT.

package defparse

import "strconv"

const _parsePrimaryState_name = "parsePrimaryState_below_expectEOFexpectValueOrSpaceexpectRocketexpectMapKeyexpectWordexpectDigitexpectContentOfexpectEscapedContentOfexpectSpaceOrMapKeyexpectSpaceOrQwItemOrDelimiterexpectEndOfQwItemparsePrimaryState_above_"

var _parsePrimaryState_index = [...]uint8{0, 24, 33, 51, 63, 75, 85, 96, 111, 133, 152, 182, 199, 223}

func (i parsePrimaryState) String() string {
	if i >= parsePrimaryState(len(_parsePrimaryState_index)-1) {
		return "parsePrimaryState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _parsePrimaryState_name[_parsePrimaryState_index[i]:_parsePrimaryState_index[i+1]]
}
