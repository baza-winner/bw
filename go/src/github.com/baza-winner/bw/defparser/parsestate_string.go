// Code generated by "stringer -type=parseState"; DO NOT EDIT.

package defparser

import "strconv"

const _parseState_name = "expectSpaceOrValueexpectSpaceOrArrayItemexpectDigitexpectDigitOrUnderscoreOrDotexpectDigitOrUnderscoreexpectDoubleQuotedStringContentexpectSingleQuotedStringContentexpectDoubleQuotedStringEscapedContentexpectSingleQuotedStringEscapedContentexpectWordexpectKeyunexpectedCharexpectArrayItemSeparatorOrSpaceOrArrayValueexpectMapKeySeparatorOrSpaceexpectMapKeySeparatorOrSpaceOrMapValueexpectSpaceOrMapValueexpectRocketexpectMapValueSeparatorOrSpaceOrMapValueexpectSpaceOrMapKeyexpectMapValueSeparatorOrSpaceOrMapKeyexpectDoubleQuotedKeyContentexpectSingleQuotedKeyContentexpectDoubleQuotedKeyEscapedContentexpectSingleQuotedKeyEscapedContentexpectSpaceOrQwItemOrDelimiterexpectEndOfQwItemexpectArrayItemSeparatorOrSpaceexpectEOFexpectSpaceOrEOF"

var _parseState_index = [...]uint16{0, 18, 40, 51, 79, 102, 133, 164, 202, 240, 250, 259, 273, 316, 344, 382, 403, 415, 455, 474, 512, 540, 568, 603, 638, 668, 685, 716, 725, 741}

func (i parseState) String() string {
	if i >= parseState(len(_parseState_index)-1) {
		return "parseState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _parseState_name[_parseState_index[i]:_parseState_index[i+1]]
}
