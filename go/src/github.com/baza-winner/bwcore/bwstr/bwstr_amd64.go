package bwstr

import "strconv"

func ParseInt(s string) (result int, err error) {
	if _int64, err := strconv.ParseInt(s, 10, 64); err != nil {
		return 0, err
	} else {
		return int(_int64), nil
	}
}
