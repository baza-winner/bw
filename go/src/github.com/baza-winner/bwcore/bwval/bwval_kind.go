package bwval

// ============================================================================

// Kind - определяет разновидность  interface{}-значения
// func Kind(val interface{}) (result interface{}, kind bwtype.ValKind) {
// 	if val == nil {
// 		kind = bwtype.ValNil
// 	} else {
// 		switch t := val.(type) {
// 		case bool:
// 			result = t
// 			kind = bwtype.ValBool
// 		case string:
// 			result = t
// 			kind = bwtype.ValString
// 		case map[string]interface{}:
// 			result = t
// 			kind = bwtype.ValMap
// 		case []interface{}:
// 			result = t
// 			kind = bwtype.ValArray
// 		case bw.ValPath:
// 			result = t
// 			kind = bwtype.ValPath
// 		case bwtype.Range:
// 			result = t
// 			kind = bwtype.ValRange
// 		default:
// 			if i, ok := bwtype.Int(val); ok {
// 				result = i
// 				kind = bwtype.ValInt
// 			} else if u, ok := bwtype.Uint(val); ok {
// 				result = u
// 				kind = bwtype.ValUint
// 			} else if f, ok := bwtype.Float64(val); ok {
// 				result = f
// 				kind = bwtype.ValFloat64
// 			}
// 		}
// 	}
// 	return
// }

// ============================================================================

// // ValKind - разновидность interface{}-значения
// type ValKind uint8

// // разновидности interface{}-значения
// const (
// 	ValUnknown ValKind = iota
// 	ValNil
// 	ValBool
// 	ValInt
// 	ValUint
// 	ValFloat64
// 	ValString
// 	ValArray
// 	ValArrayOfString
// 	ValArrayOf
// 	ValMap
// 	ValNumber
// 	ValRange
// 	// ValKindAbove
// )

// // ============================================================================

// //go:generate bwsetter -type=ValKind -test

// //go:generate stringer -type ValKind -trimprefix Val

// const (
// 	_ValKindSetTestItemA ValKind = ValNil
// 	_ValKindSetTestItemB ValKind = ValBool
// )

// // ============================================================================

// var (
// 	ansiUknownValKind    string
// 	mapValKindFromString = map[string]ValKind{}
// )

// func init() {
// 	for i := ValUnknown; i < ValKindAbove; i++ {
// 		mapValKindFromString[i.String()] = i
// 	}
// 	ansiUknownValKind = ansi.String("<ansiPath>ValKindFromString<ansi>: uknown <ansiVal>%s")
// }

// func ValKindFromString(s string) (result ValKind, err error) {
// 	var ok bool
// 	if result, ok = mapValKindFromString[s]; !ok {
// 		err = bwerr.From(ansiUknownValKind, result)
// 	}
// 	return
// }

// // ============================================================================

// // MarshalJSON encoding/json support
// func (v ValKind) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(v.String())
// }

// // Kind - определяет разновидность  interface{}-значения
// func Kind(val interface{}) (result interface{}, kind ValKind) {
// 	if val == nil {
// 		kind = ValNil
// 	} else {
// 		switch t := val.(type) {
// 		case bool:
// 			result = t
// 			kind = ValBool
// 		case int8:
// 			result = int(t)
// 			kind = ValInt
// 		case int16:
// 			result = int(t)
// 			kind = ValInt
// 		case int32:
// 			result = int(t)
// 			kind = ValInt
// 		case int64:
// 			if int64(bw.MinInt) <= t && t <= int64(bw.MaxInt) {
// 				result = int(t)
// 				kind = ValInt
// 			}
// 		case int:
// 			result = t
// 			kind = ValInt
// 		case uint8:
// 			result = int(t)
// 			kind = ValInt
// 		case uint16:
// 			result = int(t)
// 			kind = ValInt
// 		case uint32:
// 			result = int(t)
// 			kind = ValInt
// 		case uint64:
// 			if t <= uint64(bw.MaxInt) {
// 				result = int(t)
// 				kind = ValInt
// 			}
// 		case uint:
// 			if t <= uint(bw.MaxInt) {
// 				result = int(t)
// 				kind = ValInt
// 			}
// 		case float32:
// 			result = float64(t)
// 			kind = ValFloat64
// 		case float64:
// 			result = t
// 			kind = ValFloat64
// 		case string:
// 			result = t
// 			kind = ValString
// 		case map[string]interface{}:
// 			result = t
// 			kind = ValMap
// 		case []interface{}:
// 			result = t
// 			kind = ValArray
// 		case []string:
// 			result = t
// 			kind = ValArrayOfString
// 		// case bwtype.Number:
// 		// 	result = t
// 		// 	kind = ValNumber
// 		case bwtype.Range:
// 			result = t
// 			kind = ValRange
// 		}
// 	}
// 	// bwdebug.Print("val:#v", val, "kind", kind, "result", result)
// 	return
// }

// // ============================================================================
