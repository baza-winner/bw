package bwtype

import "github.com/baza-winner/bwcore/bw"

func platformSpecificInt(val interface{}) (result int, ok bool) {
	switch t := val.(type) {
	case int64:
		if ok = int64(bw.MinInt) <= t && t <= int64(bw.MaxInt); ok {
			result = int(t)
		}
	case uint32:
		if ok = t <= uint32(bw.MaxInt); ok {
			result = int(t)
		}
	default:
		result, ok = reflectInt(val)
	}
	return
}

func platformSpecificUint(val interface{}) (result uint, ok bool) {
	switch t := val.(type) {
	case int64:
		if ok = 0 <= t && t <= int64(bw.MaxUint); ok {
			result = uint(t)
		}
	case uint64:
		if ok = t <= uint64(bw.MaxUint); ok {
			result = uint(t)
		}
	default:
		result, ok = reflectUint(val)
	}
	return
}
