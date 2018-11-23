package bwstr

import "strconv"

func ParseInt(s string) (result int, err error) {
	if i, err := strconv.ParseInt(s, 10, 64); err != nil {
		return 0, err
	} else {
		return int(i), nil
	}
}

func ParseUint(s string) (result uint, err error) {
	if u, err := strconv.ParseUint(s, 10, 64); err != nil {
		return 0, err
	} else {
		return uint(u), nil
	}
}
