// Code generated by "stringer -type=ReturnType"; DO NOT EDIT.

package helper

import "strconv"

const _ReturnType_name = "rt_below_rtNonertBoolrtIntrtInterfacertStringrtSetrtSlicertSliceOfStringsrt_above_"

var _ReturnType_index = [...]uint8{0, 9, 15, 21, 26, 37, 45, 50, 57, 73, 82}

func (i ReturnType) String() string {
	if i >= ReturnType(len(_ReturnType_index)-1) {
		return "ReturnType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ReturnType_name[_ReturnType_index[i]:_ReturnType_index[i+1]]
}