package bwtype

import "github.com/baza-winner/bwcore/bw"

func platformSpecificInt(val interface{}) (result int, ok bool) {
	ok = true
	switch t := val.(type) {
	case int64:
		if int64(bw.MinInt) <= t && t <= int64(bw.MaxInt) {
			result = int(t)
		} else {
			ok = false
		}
	case uint32:
		if t <= uint32(bw.MaxInt) {
			result = int(t)
		} else {
			ok = false
		}
	default:
		ok = false
	}
	return
}

func platformSpecificUint(val interface{}) (result uint, ok bool) {
	ok = true
	switch t := val.(type) {
	case int64:
		if ok = 0 <= t && t <= int64(bw.MaxUint); ok {
			result = uint(t)
		}
	case uint64:
		if ok = t <= uint64(bw.MaxUint); ok {
			result = uint(t)
		}
	}
	return
}
