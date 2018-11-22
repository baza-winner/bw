package bwtype

func platformSpecificInt(val interface{}) (result int, ok bool) {
	switch t := val.(type) {
	case int64:
		result = int(t)
		ok = true
	case uint32:
		result = int(t)
		ok = true
	}
	return
}

func platformSpecificUint(val interface{}) (result uint, ok bool) {
	switch t := val.(type) {
	case int64:
		if ok = 0 <= t; ok {
			result = uint(t)
		}
	case uint64:
		result = uint(t)
		ok = true
	}
	return
}

// func Int(val interface{}) (result int, ok bool) {
// 	switch t := val.(type) {
// 	// case int8:
// 	// 	result = int(t)
// 	// case int16:
// 	// 	result = int(t)
// 	// case int32:
// 	// 	result = int(t)
// 	case int64:
// 		result = int(t)
// 		ok = true
// 	// case int:
// 	// 	result = t
// 	// case uint8:
// 	// 	result = int(t)
// 	// case uint16:
// 	// 	result = int(t)
// 	case uint32:
// 		result = int(t)
// 		ok = true
// 		// case uint64:
// 		// 	if t <= uint64(bw.MaxInt) {
// 		// 		result = int(t)
// 		// 	} else {
// 		// 		ok = false
// 		// 	}
// 		// case uint:
// 		// 	if t <= uint(bw.MaxInt) {
// 		// 		result = int(t)
// 		// 	} else {
// 		// 		ok = false
// 		// 	}
// 		// case float32:
// 		// 	result = int(t)
// 		// 	if t == float32(result) {
// 		// 		ok = true
// 		// 	} else {
// 		// 		result = 0
// 		// 	}
// 		// case float64:
// 		// 	result = int(t)
// 		// 	if t == float64(result) {
// 		// 		ok = true
// 		// 	} else {
// 		// 		result = 0
// 		// 	}
// 		// default:
// 		// ok = false
// 	}
// 	return
// }

// func Uint(val interface{}) (result uint, ok bool) {
// 	ok = true
// 	switch t := val.(type) {
// 	case int8:
// 		if ok = t >= 0; ok {
// 			result = uint(t)
// 		}
// 	case int16:
// 		if ok = t >= 0; ok {
// 			result = uint(t)
// 		}
// 	case int32:
// 		if ok = t >= 0; ok {
// 			result = uint(t)
// 		}
// 	case int64:
// 		if ok = t >= 0; ok {
// 			result = uint(t)
// 		}
// 	case int:
// 		if ok = t >= 0; ok {
// 			result = uint(t)
// 		}
// 	case uint8:
// 		result = uint(t)
// 	case uint16:
// 		result = uint(t)
// 	case uint32:
// 		result = uint(t)
// 	case uint64:
// 		result = uint(t)
// 	case uint:
// 		result = t
// 	case float32:
// 		result = uint(t)
// 		if t == float32(result) {
// 			ok = true
// 		} else {
// 			result = 0
// 		}
// 	case float64:
// 		result = uint(t)
// 		if t == float64(result) {
// 			ok = true
// 		} else {
// 			result = 0
// 		}
// 	default:
// 		ok = false
// 	}
// 	return
// }
