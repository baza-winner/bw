// Code generated by "stringer -type=parseSecondaryState"; DO NOT EDIT.

package pathparse

import "strconv"

const _parseSecondaryState_name = "noSecondaryStateorPathSegmentDelimiterorUnderscoredoubleQuotedsingleQuoted"

var _parseSecondaryState_index = [...]uint8{0, 16, 38, 50, 62, 74}

func (i parseSecondaryState) String() string {
	if i >= parseSecondaryState(len(_parseSecondaryState_index)-1) {
		return "parseSecondaryState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _parseSecondaryState_name[_parseSecondaryState_index[i]:_parseSecondaryState_index[i+1]]
}
