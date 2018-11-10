package bwval

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/bw"
)

// ============================================================================

// ValKind - разновидность interface{}-значения
type ValKind uint8

// разновидности interface{}-значения
const (
	ValUnknown ValKind = iota
	ValNil
	ValBool
	ValInt
	ValNumber
	ValString
	ValMap
	ValArray
)

// MarshalJSON encoding/json support
func (v ValKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// Kind - определяет разновидность  interface{}-значения
func Kind(val interface{}) (result interface{}, kind ValKind) {
	if val == nil {
		kind = ValNil
	} else {
		switch t := val.(type) {
		case bool:
			result = t
			kind = ValBool
		case int8:
			result = int(t)
			kind = ValInt
		case int16:
			result = int(t)
			kind = ValInt
		case int32:
			result = int(t)
			kind = ValInt
		case int64:
			if int64(bw.MinInt) <= t && t <= int64(bw.MaxInt) {
				result = int(t)
				kind = ValInt
			}
		case int:
			result = t
			kind = ValInt
		case uint8:
			result = int(t)
			kind = ValInt
		case uint16:
			result = int(t)
			kind = ValInt
		case uint32:
			result = int(t)
			kind = ValInt
		case uint64:
			if t <= uint64(bw.MaxInt) {
				result = int(t)
				kind = ValInt
			}
		case uint:
			if t <= uint(bw.MaxInt) {
				result = int(t)
				kind = ValInt
			}
		case float32:
			result = float64(t)
			kind = ValNumber
		case float64:
			result = t
			kind = ValNumber
		case string:
			result = t
			kind = ValString
		case map[string]interface{}:
			result = t
			kind = ValMap
		case []interface{}:
			result = t
			kind = ValArray
		}
	}
	return
}
