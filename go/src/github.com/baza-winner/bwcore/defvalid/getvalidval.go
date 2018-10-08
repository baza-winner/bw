package defvalid

import (
	"github.com/baza-winner/bwcore/bwerror"
	// "github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/defvalid/deftype"
	// "log"
	"reflect"
)

func getValidVal(val value, def Def, optSkipArrayOf ...bool) (result interface{}, err error) {
	skipArrayOf := optSkipArrayOf != nil && optSkipArrayOf[0]
	if val.value == nil {
		if !skipArrayOf {
			if def.dflt != nil {
				return def.dflt, nil
			}
			if def.isOptional {
				return nil, nil
			}
		}
		if def.tp.Has(deftype.Map) {
			val.value = map[string]interface{}{}
		} else {
			return nil, valueErrorMake(val, valueErrorIsNotOfType, def.tp)
		}
	}
	var valDeftype deftype.Item
	valType := reflect.TypeOf(val.value)
	switch valType.Kind() {
	case reflect.Bool:
		if def.tp.Has(deftype.Bool) {
			valDeftype = deftype.Bool
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if def.tp.Has(deftype.Int) {
			valDeftype = deftype.Int
		} else if def.tp.Has(deftype.Number) {
			valDeftype = deftype.Number
		}
	case reflect.Float32, reflect.Float64:
		if def.tp.Has(deftype.Number) {
			valDeftype = deftype.Number
		}
	case reflect.Map:
		if valType.Key().Kind() == reflect.String && valType.Elem().Kind() == reflect.Interface && def.tp.Has(deftype.Map) {
			valDeftype = deftype.Map
		}
	case reflect.Slice:
		if def.tp.Has(deftype.Array) {
			valDeftype = deftype.Array
		} else if !skipArrayOf && def.tp.Has(deftype.ArrayOf) {
			valDeftype = deftype.ArrayOf
		}

	case reflect.String:
		if def.tp.Has(deftype.String) {
			valDeftype = deftype.String
		}
	}
	if valDeftype == deftype.Item_below_ {
		return nil, valueErrorMake(val, valueErrorIsNotOfType, def.tp)
	}

	if val.value, err = getValidValHelpers[valDeftype](val, def); err != nil {
		return nil, err
	}

	if valDeftype != deftype.ArrayOf && def.tp.Has(deftype.ArrayOf) {
		val.value = []interface{}{val.value}
	}

	return val.value, nil
}

type getValidValHelper func(val value, def Def) (result interface{}, err error)

var getValidValHelpers map[deftype.Item]getValidValHelper

func getValidValHelpersCheck() {
	deftypeItem := deftype.Item_below_ + 1
	for deftypeItem < deftype.Item_above_ {
		if _, ok := getValidValHelpers[deftypeItem]; !ok {
			bwerror.Panic("not defined <ansiOutline>deftype.ItemValidators<ansi>[<ansiPrimaryLiteral>%s<ansi>]", deftypeItem)
		}
		deftypeItem += 1
	}
}

func _Bool(val value, def Def) (result interface{}, err error) {
	result = val.value
	return
}

func _String(val value, def Def) (result interface{}, err error) {
	if def.enum != nil {
		if !def.enum.Has(_mustBeString(val.value)) {
			err = valueErrorMake(val, valueErrorHasNonSupportedValue)
		}
	}
	result = val.value
	return
}

func _Int(val value, def Def) (result interface{}, err error) {
	if def.minInt != nil || def.maxInt != nil {
		n := _mustBeInt64(val.value)
		var isOutOfRange bool
		if def.minInt != nil {
			if def.maxInt != nil {
				isOutOfRange = !(*(def.minInt) <= n && n <= *(def.maxInt))
			} else {
				isOutOfRange = !(*(def.minInt) <= n)
			}
		} else {
			isOutOfRange = !(n <= *(def.maxInt))
		}
		if isOutOfRange {
			err = valueErrorMake(val, valueErrorOutOfRange, def.minInt, def.maxInt)
		}
	}
	result = val.value
	return
}

func _Number(val value, def Def) (result interface{}, err error) {
	if def.minNumber != nil || def.maxNumber != nil {
		n := _mustBeFloat64(val.value)
		var isOutOfRange bool
		if def.minNumber != nil {
			if def.maxNumber != nil {
				isOutOfRange = !(*(def.minNumber) <= n && n <= *(def.maxNumber))
			} else {
				isOutOfRange = !(*(def.minNumber) <= n)
			}
		} else {
			isOutOfRange = !(n <= *(def.maxNumber))
		}
		if isOutOfRange {
			err = valueErrorMake(val, valueErrorOutOfRange, def.minNumber, def.maxNumber)
		}
	}
	result = val.value
	return
}

func _Map(val value, def Def) (result interface{}, err error) {
	if def.keys != nil {
		unexpectedKeys := bwmap.GetUnexpectedKeys(val.value, def.keys)
		for key, keyDef := range def.keys {
			if err = _MapHelper(val, key, keyDef); err != nil {
				return
			}
		}
		if unexpectedKeys != nil {
			if def.elem == nil {
				return nil, valueErrorMake(val, valueErrorHasUnexpectedKeys, unexpectedKeys)
			} else {
				for _, key := range unexpectedKeys {
					if err = _MapHelper(val, key, *(def.elem)); err != nil {
						return
					}
				}
			}
		}
	} else if def.elem != nil {
		err = val.forEachMapString(func(k string, v interface{}) error {
			return _MapHelper(val, k, *(def.elem))
		})
	}
	result = val.value
	return
}

func _MapHelper(val value, key string, elemDef Def) (err error) {
	var elemVal value
	if elemVal, err = val.getKey(key, "interface{}", nil); err != nil {
		return err
	}
	var elemValIntf interface{}
	if elemValIntf, err = getValidVal(elemVal, elemDef); err != nil {
		return err
	}
	if elemValIntf != nil {
		if err = val.setKey(key, elemValIntf); err != nil {
			return err
		}
	}
	return
}

func _Array(val value, def Def) (result interface{}, err error) {
	elemDef := def.arrayElem
	if elemDef == nil {
		elemDef = def.elem
	}
	if elemDef == nil {
		result = val.value
	} else {
		result, err = _ArrayHelper(val, *elemDef)
	}
	return
}

func _ArrayOf(val value, def Def) (result interface{}, err error) {
	return _ArrayHelper(val, def, true)
}

func _ArrayHelper(val value, elemDef Def, optSkipArrayOf ...bool) (result interface{}, err error) {
	newSliceValue := reflect.MakeSlice(reflect.TypeOf(val.value), 0, reflect.ValueOf(val.value).Len())
	err = val.forEachSlice(func(i int, v interface{}) (err error) {
		var elemVal value
		if elemVal, err = val.getElem(i); err == nil {
			var elemValIntf interface{}
			if elemValIntf, err = getValidVal(elemVal, elemDef, optSkipArrayOf...); err == nil && elemValIntf != nil {
				reflect.Append(newSliceValue, reflect.ValueOf(elemValIntf))
				// newSlice = append(newSlice, elemValIntf)
			}
		}
		// log.Printf("i: %d, v: %#v, err: %#v", i, v, err)
		return
	})
	return newSliceValue.Interface(), err
}

func getExpectedTypes(defType value, isSimpleDef bool) (result bwset.Strings, err error) {
	var isString bool
	var ss []string
	if _isOfType(defType.value, "string") {
		ss = []string{_mustBeString(defType.value)}
		isString = true
	} else if _isOfType(defType.value, "[]string") {
		ss = _mustBeSliceOfStrings(defType.value)
		isString = false
	} else {
		err = valueErrorMake(defType, valueErrorIsNotOfType, "string", "[]string")
	}
	if err == nil {
		for i, s := range ss {
			switch s {
			case "map":
			case "bool":
			case "string":
			case "enum":
				var elem value
				if isSimpleDef {
					if isString {
						err = valueErrorMake(defType, valueErrorHasNonSupportedValue)
					} else if elem, err = defType.getElem(i); err == nil {
						err = valueErrorMake(elem, valueErrorHasNonSupportedValue)
					}
				}
			default:
				var elem value
				if isString {
					err = valueErrorMake(defType, valueErrorHasNonSupportedValue)
				} else if elem, err = defType.getElem(i); err == nil {
					err = valueErrorMake(elem, valueErrorHasNonSupportedValue)
				}
			}
		}
		result = bwset.StringsFromSlice(ss)
		if result.Has("enum") && result.Has("string") {
			// log.Printf("%s\n", result)
			err = valueErrorMake(defType, valueErrorValuesCannotBeCombined, "enum", "string")
		}
	}
	return
}
