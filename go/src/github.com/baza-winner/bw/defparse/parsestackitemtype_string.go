// Code generated by "stringer -type=parseStackItemType"; DO NOT EDIT.

package defparse

import "strconv"

const _parseStackItemType_name = "_parseStackItemBelowparseStackItemKeyparseStackItemStringparseStackItemMapparseStackItemArrayparseStackItemQwparseStackItemQwItemparseStackItemNumberparseStackItemWord_parseStackItemAbove"

var _parseStackItemType_index = [...]uint8{0, 20, 37, 57, 74, 93, 109, 129, 149, 167, 187}

func (i parseStackItemType) String() string {
	if i >= parseStackItemType(len(_parseStackItemType_index)-1) {
		return "parseStackItemType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _parseStackItemType_name[_parseStackItemType_index[i]:_parseStackItemType_index[i+1]]
}