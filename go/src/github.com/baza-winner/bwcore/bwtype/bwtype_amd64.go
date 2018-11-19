package bwtype

import "github.com/baza-winner/bwcore/bw"

func Int(val interface{}) (result int, ok bool) {
	ok = true
	switch t := val.(type) {
	case int8:
		result = int(t)
	case int16:
		result = int(t)
	case int32:
		result = int(t)
	case int64:
		result = int(t)
	case int:
		result = t
	case uint8:
		result = int(t)
	case uint16:
		result = int(t)
	case uint32:
		result = int(t)
	case uint64:
		if t <= uint64(bw.MaxInt) {
			result = int(t)
		} else {
			ok = false
		}
	case uint:
		if t <= uint(bw.MaxInt) {
			result = int(t)
		} else {
			ok = false
		}
	default:
		ok = false
	}
	return
}
