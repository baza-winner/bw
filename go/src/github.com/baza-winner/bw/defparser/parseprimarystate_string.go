// Code generated by "stringer -type=parsePrimaryState"; DO NOT EDIT.

package defparser

import "strconv"

const _parsePrimaryState_name = "_expectBelowexpectEOFexpectValueOrSpaceexpectArrayItemSeparatorOrSpaceexpectMapKeySeparatorOrSpaceexpectRocketexpectMapKeyexpectWordexpectDigitexpectContentOfexpectEscapedContentOfexpectSpaceOrMapKeyexpectSpaceOrQwItemOrDelimiterexpectEndOfQwItem_expectAbove"

var _parsePrimaryState_index = [...]uint16{0, 12, 21, 39, 70, 98, 110, 122, 132, 143, 158, 180, 199, 229, 246, 258}

func (i parsePrimaryState) String() string {
	if i >= parsePrimaryState(len(_parsePrimaryState_index)-1) {
		return "parsePrimaryState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _parsePrimaryState_name[_parsePrimaryState_index[i]:_parsePrimaryState_index[i+1]]
}
