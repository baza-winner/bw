// Code generated by "stringer -type=PrimaryState,SecondaryState,ItemKind"; DO NOT EDIT.

package val

import "strconv"

const _PrimaryState_name = "BeginExpectSpaceOrQwItemOrDelimiterExpectSpaceOrMapKeyExpectEndOfQwItemExpectContentOfExpectWordExpectEscapedContentOfExpectRocketExpectDigitEnd"

var _PrimaryState_index = [...]uint8{0, 5, 35, 54, 71, 86, 96, 118, 130, 141, 144}

func (i PrimaryState) String() string {
	if i >= PrimaryState(len(_PrimaryState_index)-1) {
		return "PrimaryState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PrimaryState_name[_PrimaryState_index[i]:_PrimaryState_index[i+1]]
}

const _SecondaryState_name = "NoneorArrayItemSeparatororMapKeySeparatororMapValueSeparatororUnderscoreOrDotorUnderscore"

var _SecondaryState_index = [...]uint8{0, 4, 24, 41, 60, 77, 89}

func (i SecondaryState) String() string {
	if i >= SecondaryState(len(_SecondaryState_index)-1) {
		return "SecondaryState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SecondaryState_name[_SecondaryState_index[i]:_SecondaryState_index[i+1]]
}

const _ItemKind_name = "ItemStringItemQwItemQwItemItemNumberItemWordItemKeyItemMapItemArray"

var _ItemKind_index = [...]uint8{0, 10, 16, 26, 36, 44, 51, 58, 67}

func (i ItemKind) String() string {
	if i >= ItemKind(len(_ItemKind_index)-1) {
		return "ItemKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ItemKind_name[_ItemKind_index[i]:_ItemKind_index[i+1]]
}
