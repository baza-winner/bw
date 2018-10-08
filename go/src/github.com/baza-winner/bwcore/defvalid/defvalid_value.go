package defvalid

import (
	"fmt"
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwint"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwset"
	"reflect"
	// "log"
)

func init() {
	valueErrorValidatorsCheck()
}

type value struct {
	what  string
	value interface{}
}

func (v value) GetDataForJson() interface{} {
	result := map[string]interface{}{}
	result["where"] = v.what
	result["value"] = v.value
	return result
}

func (v value) String() string {
	return fmt.Sprintf(v.what+`<ansi> (<ansiSecondaryLiteral>%s<ansi>)`, bwjson.PrettyJson(v.value))
}

func (v value) asMap() (result map[string]interface{}, err error) {
	var ok bool
	if result, ok = v.value.(map[string]interface{}); !ok {
		err = valueErrorMake(v, valueErrorIsNotOfType, "map")
	}
	return
}

func (v value) mustBeMap() (result map[string]interface{}) {
	var err error
	if result, err = v.asMap(); err != nil {
		bwerror.Panic(err.Error())
	}
	return
}

// func (v value) IsMapString(m interface{}) bool {

// }

func (v value) forEachMapString(f func(k string, v interface{}) (err error)) (err error) {
	if !_isOfType(v.value, "map[string]") {
		err = valueErrorMake(v, valueErrorIsNotOfType, "map[string]")
	} else {
		mv := reflect.ValueOf(v.value)
		mk := mv.MapKeys()
		for i := 0; i < len(mk); i++ {
			err = f(mk[i].String(), mv.MapIndex(mk[i]).Interface())
			if err != nil {
				break
			}
		}
	}
	return err
}

// func (v value) asString() (result string, err error) {
// 	var ok bool
// 	if result, ok = v.value.(string); !ok {
// 		err = valueErrorMake(v, valueErrorIsNotOfType, "string")
// 	}
// 	return
// }

// func (v value) mustBeString() (result string) {
// 	var err error
// 	if result, err = v.asString(); err != nil {
// 		bwerror.Panic(err.Error())
// 	}
// 	return
// }

// func (v value) asBool() (result bool, err error) {
// 	var ok bool
// 	if result, ok = v.value.(bool); !ok {
// 		err = valueErrorMake(v, valueErrorIsNotOfType, "bool")
// 	}
// 	return
// }

func (v value) getElem(elemIndex int, opts ...interface{}) (result value, err error) {
	defaultValue, ofType := getDefaultValueAndOfTypeFromOpts(opts)
	if v.value == nil {
		err = valueErrorMake(v, valueErrorIsNotOfType, "array")
	} else {
		vType := reflect.TypeOf(v.value)
		if vType.Kind() != reflect.Slice {
			err = valueErrorMake(v, valueErrorIsNotOfType, "array")
		} else {
			sv := reflect.ValueOf(v.value)
			result.what = v.what + fmt.Sprintf(".#%d", elemIndex)
			if 0 <= elemIndex && elemIndex < sv.Len() {
				err = checkElemIsOfType(&result, sv.Index(elemIndex), ofType)
			} else if defaultValue == nil {
				err = valueErrorMake(v, valueErrorHasNoKey, fmt.Sprintf("#%d", elemIndex))
			} else {
				result.value = *defaultValue
			}
		}
	}
	return
}

func (v value) getKey(keyName string, opts ...interface{}) (result value, err error) {
	defaultValue, ofType := getDefaultValueAndOfTypeFromOpts(opts)
	if !_isOfType(v.value, "map[string]") {
		err = valueErrorMake(v, valueErrorIsNotOfType, "map[string]")
	} else {
		mv := reflect.ValueOf(v.value)
		result.what = v.what + "." + keyName
		elem := mv.MapIndex(reflect.ValueOf(keyName))
		zeroValue := reflect.Value{}
		if elem != zeroValue {
			err = checkElemIsOfType(&result, elem, ofType)
		} else if defaultValue == nil {
			err = valueErrorMake(v, valueErrorHasNoKey, keyName)
		} else {
			result.value = *defaultValue
		}
	}
	// if v.value == nil {
	// 	err = valueErrorMake(v, valueErrorIsNotOfType, "map[string]")
	// } else {
	// 	vType := reflect.TypeOf(v.value)
	// 	if vType.Kind() != reflect.Map || vType.Key().Kind() != reflect.String {
	// 		err = valueErrorMake(v, valueErrorIsNotOfType, "map[string]")
	// 	} else {
	// 		mv := reflect.ValueOf(v.value)
	// 		result.what = v.what + "." + keyName
	// 		elem := mv.MapIndex(reflect.ValueOf(keyName))
	// 		zeroValue := reflect.Value{}
	// 		if elem != zeroValue {
	// 			err = checkElemIsOfType(&result, elem, ofType)
	// 		} else if defaultValue == nil {
	// 			err = valueErrorMake(v, valueErrorHasNoKey, keyName)
	// 		} else {
	// 			result.value = *defaultValue
	// 		}
	// 	}
	// }
	return
}

func checkElemIsOfType(result *value, elem reflect.Value, ofType []string) (err error) {
	result.value = elem.Interface()
	if ofType != nil && !_isOfType(result.value, ofType...) {
		var ofTypeIntfs = []interface{}{}
		for _, i := range ofType {
			ofTypeIntfs = append(ofTypeIntfs, i)
		}
		err = valueErrorMake(*result, valueErrorIsNotOfType, ofTypeIntfs...)
	}
	return
}

func getDefaultValueAndOfTypeFromOpts(opts []interface{}) (defaultValue *interface{}, ofType []string) {
	if opts != nil {
		if _isOfType(opts[0], "string") {
			ofType = []string{_mustBeString(opts[0])}
		} else if _isOfType(opts[0], "[]string") {
			ofType = _mustBeSliceOfStrings(opts[0])
		} else {
			_ = _mustBeOfType(opts[0], "string", "[]string")
		}
		if len(opts) > 1 {
			defaultValueIntf := opts[1]
			defaultValue = &defaultValueIntf
		}
		if len(opts) > 2 {
			bwerror.Panic("expects max 2 opts (ofTypes, defaultValue), but found <ansiSecondaryLiteral>%v", opts)
		}
	}
	return
}

func _isOfType(v interface{}, ofTypes ...string) (ok bool) {
	if v != nil {
		vType := reflect.TypeOf(v)
		for _, ofType := range ofTypes {
			switch ofType {
			case "string", "enum":
				ok = vType.Kind() == reflect.String
			case "[]string":
				if vType.Kind() == reflect.Slice {
					elemType := vType.Elem()
					if elemType.Kind() == reflect.String || elemType.Kind() == reflect.Interface {
						ok = true
						if elemType.Kind() == reflect.Interface {
							sv := reflect.ValueOf(v)
							for i := 0; i < sv.Len(); i++ {
								if ok = _isOfType(sv.Index(i).Interface(), "string"); !ok {
									break
								}
							}
						}
					}
				}
			case "map[string]":
				if vType.Kind() == reflect.Map {
					keyType := vType.Key()
					if keyType.Kind() == reflect.String || keyType.Kind() == reflect.Interface {
						ok = true
						if keyType.Kind() == reflect.Interface {
							mk := reflect.ValueOf(v).MapKeys()
							for i := 0; i < len(mk); i++ {
								if ok = _isOfType(mk[i].Interface(), "string"); !ok {
									break
								}
							}
						}
					}
				}
			case "int64":
				switch vType.Kind() {
				case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					ok = true
				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					ok = reflect.ValueOf(v).Uint() <= uint64(bwint.MaxInt64)
				default:
					ok = false
				}
			case "float64":
				switch vType.Kind() {
				case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					ok = true
				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					ok = true
				case reflect.Float32, reflect.Float64:
					ok = true
				default:
					ok = false
				}
			case "bool":
				ok = vType.Kind() == reflect.Bool
			case "bwset.Strings":
				_, ok = v.(bwset.Strings)
			case "interface{}":
				ok = true
			default:
				bwerror.Panic("unsupported type <ansiPrimaryLiteral>%s", ofType)
			}
			if ok {
				break
			}
		}
	}
	return
}

func _mustBeOfType(v interface{}, ofTypes ...string) (result interface{}) {
	if !_isOfType(v, ofTypes...) {
		bwerror.Panic("<ansiSecondaryLiteral>%+v<ansi> is not of types <ansiSecondaryLiteral>%v", v, ofTypes)
	}
	return v
}

func _mustBeString(v interface{}) (result string) {
	result, _ = _mustBeOfType(v, "string").(string)
	return
}

func _mustBeInt64(v interface{}) (result int64) {
	vValue := reflect.ValueOf(v)
	switch vValue.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = vValue.Int()
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if reflect.ValueOf(v).Uint() <= uint64(bwint.MaxInt64) {
			result = int64(vValue.Uint())
		} else {
			bwerror.Panic("<ansiSecondaryLiteral>%+v<ansi> is not of type <ansiSecondaryLiteral>int64", v)
		}
	default:
		bwerror.Panic("<ansiSecondaryLiteral>%+v<ansi> is not of type <ansiSecondaryLiteral>int64", v)
	}
	return
}

func _mustBeFloat64(v interface{}) (result float64) {
	vValue := reflect.ValueOf(v)
	switch vValue.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = float64(vValue.Int())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = float64(vValue.Uint())
	case reflect.Float32, reflect.Float64:
		result = vValue.Float()
	default:
		bwerror.Panic("<ansiSecondaryLiteral>%+v<ansi> is not of type <ansiSecondaryLiteral>float64", v)
	}
	return
}

func _mustBeSliceOfStrings(v interface{}) (result []string) {
	var ok bool
	if result, ok = _mustBeOfType(v, "[]string").([]string); !ok {
		result = []string{}
		sv := reflect.ValueOf(v)
		for i := 0; i < sv.Len(); i++ {
			s, _ := sv.Index(i).Interface().(string)
			result = append(result, s)
		}
	}
	return
}

func _mustBeSetOfStrings(v interface{}) (result bwset.Strings) {
	result, _ = _mustBeOfType(v, "bwset.Strings").(bwset.Strings)
	return
}

func ptrToInt64(i int64) *int64 {
	return &i
}

func ptrToFloat64(i float64) *float64 {
	return &i
}