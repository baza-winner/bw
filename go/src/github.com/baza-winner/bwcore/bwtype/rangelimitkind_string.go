// Code generated by "stringer -type RangeLimitKind -trimprefix RangeLimit"; DO NOT EDIT.

package bwtype

import "strconv"

const _RangeLimitKind_name = "NilIntFloat64Path"

var _RangeLimitKind_index = [...]uint8{0, 3, 6, 13, 17}

func (i RangeLimitKind) String() string {
	if i >= RangeLimitKind(len(_RangeLimitKind_index)-1) {
		return "RangeLimitKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RangeLimitKind_name[_RangeLimitKind_index[i]:_RangeLimitKind_index[i+1]]
}