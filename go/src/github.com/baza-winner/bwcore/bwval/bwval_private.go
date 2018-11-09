package bwval

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/defvalid/deftype"
)

// ============================================================================

//go:generate stringer -type=ValKind

func varsJSON(path bw.ValPath, optVars []map[string]interface{}) (result string) {
	if hasVar(path) {
		var vars map[string]interface{}
		if len(optVars) > 0 {
			vars = optVars[0]
		}
		result = fmt.Sprintf(ansiVars, bwjson.Pretty(vars))
	}
	return
}

type valHolder struct {
	val interface{}
}

var (
	ansiIsNotOfType string
	// ansiValAtPathIsNotOfType       string
	ansiMustPathValFailed string
	ansiType              string
	// ansiValAtPathIsNotOfTypes      string
	ansiValAtPathIsNil string
	// ansisReadOnlyPath              string
	// ansiValAtPathHasNotEnoughRange string
	ansiVars                 string
	ansiVarsIsNil            string
	ansiMustSetPathValFailed string
)

func init() {
	ansiIsNotOfType = ansi.String("<ansiVal>%#v<ansi> is not <ansiType>%s")
	ansiMustSetPathValFailed = ansi.String("Failed to set <ansiPath>%s<ansi> of <ansiVal>%s<ansi>%s: {Error}")
	ansiMustPathValFailed = ansi.String("Failed to get <ansiPath>%s<ansi> of <ansiVal>%s<ansi>%s: {Error}")
	valPathPrefix := "<ansiPath>%s<ansi> "
	// ansiValAtPathIsNotOfType = ansi.String(valPathPrefix + "(<ansiVal>%s<ansi>) is not <ansiType>%s")
	// ansiValAtPathIsNotOfTypes = ansi.String(valPathPrefix + "(<ansiVal>%s<ansi>) is none of %s")
	ansiValAtPathIsNil = ansi.String(valPathPrefix + "is <ansiErr>nil")
	// ansiValAtPathHasNotEnoughRange = ansi.String(valPathPrefix + "(<ansiVal>%s<ansi>) has not enough length (<ansiVal>%d<ansi>) for idx (<ansiVal>%d)")

	ansiType = ansi.String("<ansiType>%s")
	// ansisReadOnlyPath = ansi.String("<ansiPath>%s<ansi> is <ansiErr>readonly path")
	ansiVars = ansi.String(" with <ansiVar>vars<ansi> <ansiVal>%s<ansi>")
	ansiVarsIsNil = ansi.String("<ansiVar>vars<ansi> is <ansiErr>nil")
}

func hasVar(path bw.ValPath) bool {
	for _, vpi := range path {
		switch vpi.Type {
		case bw.ValPathItemVar:
			return true
		case bw.ValPathItemPath:
			if hasVar(vpi.Path) {
				return true
			}
		}
	}
	return false
}

func simplifyPath(v *valHolder, path bw.ValPath, optVars []map[string]interface{}) (result bw.ValPath, err error) {
	result = bw.ValPath{}
	for _, vpi := range path {
		if vpi.Type != bw.ValPathItemPath {
			result = append(result, vpi)
		} else {
			var val interface{}
			val, err = v.PathVal(vpi.Path, optVars...)
			if err == nil {
				switch _, kind := Kind(val); kind {
				case ValString:
					result = append(result, bw.ValPathItem{Type: bw.ValPathItemKey, Key: MustString(val)})
				case ValInt:
					result = append(result, bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: MustInt(val)})
				default:
					// bwdebug.Print("vpi.Path", vpi.Path)
					err = valPath{val, vpi.Path}.notOfTypeError("Int", "String")
					// err = pathValIsNotOfType(vpi.Path, val, "Int", "String")
				}
			}
		}
	}
	return
}

var typeSeparator = map[bool]string{
	true:  " or ",
	false: ", ",
}

// func (v *valHolder) pathValIsNotOfType(path bw.ValPath, val interface{}, expectedType string, optExpectedType ...string) (result error) {
// func pathValIsNotOfType(path bw.ValPath, val interface{}, expectedType string, optExpectedType ...string) (result error) {
// 	if len(optExpectedType) == 0 {
// 		result = bwerr.From(ansiValAtPathIsNotOfType, path, bwjson.Pretty(val), expectedType)
// 	} else {
// 		expectedTypes := fmt.Sprintf(ansiType, expectedType)
// 		for i, elem := range optExpectedType {
// 			expectedTypes += typeSeparator[i == len(optExpectedType)-1] + fmt.Sprintf(ansiType, elem)
// 		}
// 		result = bwerr.From(ansiValAtPathIsNotOfTypes, path, bwjson.Pretty(val), expectedTypes)
// 		// bwdebug.Print("path", path, "ansiValAtPathIsNotOfTypes", ansiValAtPathIsNotOfTypes, "result", result)
// 	}
// 	return
// }

// func (v *valHolder) valAtPathIsNil(path bw.ValPath) error {
func valAtPathIsNil(path bw.ValPath) error {
	return bwerr.From(ansiValAtPathIsNil, path)
}

func (v *valHolder) getArray(idx int, result interface{}, resultPath bw.ValPath) ([]interface{}, int, error) {
	var err error
	var ok bool
	var vals []interface{}
	if vals, ok = Array(result); !ok {
		err = valPath{result, resultPath}.notOfTypeError("Array")
		// err = pathValIsNotOfType(resultPath, result, "Array")
	} else {
		l := len(vals)
		minIdx := -l
		maxIdx := l - 1
		if !(minIdx <= idx && idx <= maxIdx) {
			idx = -1
		} else if idx < 0 {
			idx = l + idx
		}
	}
	return vals, idx, err
}

func (v *valHolder) getMap(result interface{}, resultPath bw.ValPath) (map[string]interface{}, error) {
	var err error
	var ok bool
	var m map[string]interface{}
	if m, ok = Map(result); !ok {
		err = valPath{result, resultPath}.notOfTypeError("Map")
		// err = pathValIsNotOfType(resultPath, result, "Map")
	}
	return m, err
}

func (v *valHolder) byKey(val interface{}, path bw.ValPath, i int, key string) (result interface{}, err error) {
	result = val
	if result == nil {
		return
	}
	var m map[string]interface{}
	if m, err = v.getMap(result, path[:i+1]); err == nil {
		result = m[key]
	}
	return
}

func (v *valHolder) byIdx(val interface{}, path bw.ValPath, i int, idx int) (result interface{}, err error) {
	result = val
	if result == nil {
		return
	}
	var vals []interface{}
	if vals, idx, err = v.getArray(idx, result, path[:i+1]); err == nil {
		if idx < 0 {
			result = nil
		} else {
			result = vals[idx]
		}
	}
	return
}

// ============================================================================

type valPath struct {
	val  interface{}
	path bw.ValPath
}

// func (v valPath) forEachMapString(f func(k string, v interface{}) (err error)) (err error) {
// 	if !_isOfType(v.val, "map[string]") {
// 		err = valueErrorMake(v, valueErrorIsNotOfType, "map[string]")
// 	} else {
// 		mv := reflect.ValueOf(v.val)
// 		mk := mv.MapKeys()
// 		for i := 0; i < len(mk); i++ {
// 			err = f(mk[i].String(), mv.MapIndex(mk[i]).Interface())
// 			if err != nil {
// 				break
// 			}
// 		}
// 	}
// 	return err
// }

func (v valPath) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{}
	result["val"] = v.val
	result["path"] = v.path
	return json.Marshal(result)
}

func (v valPath) ansiString() string {
	return fmt.Sprintf(ansi.String("<ansiPath>%s<ansi> (<ansiVal>%s<ansi>)"), v.path, bwjson.Pretty(v.val))
	// return fmt.Sprintf(v.what+`<ansi> (<ansiVal>%s<ansi>)`, bwjson.PrettyJson(v.val))
	// return bw.Spew.Sprintf(v.what+`<ansi> (<ansiVal>%#v<ansi>)`, v.val)
}

func (v valPath) String() (result string, err error) {
	var ok bool
	if result, ok = String(v.val); !ok {
		err = v.notOfTypeError("String")
	}
	return
}

func (v valPath) Int() (result int, err error) {
	var ok bool
	if result, ok = Int(v.val); !ok {
		err = v.notOfTypeError("Int")
	}
	return
}

func (v valPath) Number() (result float64, err error) {
	var ok bool
	if result, ok = Number(v.val); !ok {
		err = v.notOfTypeError("Number")
	}
	return
}

func (v valPath) Array() (result []interface{}, err error) {
	var ok bool
	if result, ok = Array(v.val); !ok {
		err = v.notOfTypeError("Array")
	}
	return
}

func (v valPath) ArrayOfString() (result []string, err error) {
	var vals interface{}
	if vals, err = v.Array(); err != nil {
		return
	}
	result = []string{}
	var s string
	for i := range MustArray(vals) {
		vp, _ := v.Idx(i)
		if s, err = vp.String(); err != nil {
			return
		}
		result = append(result, s)
	}
	return
}

func (v valPath) Map() (result map[string]interface{}, err error) {
	var ok bool
	if result, ok = Map(v.val); !ok {
		err = v.notOfTypeError("Map")
	}
	return
}

func (v valPath) Key(key string) (result valPath, err error) {
	var val interface{}
	vpi := bw.ValPathItem{Type: bw.ValPathItemKey, Key: key}
	if val, err = FromVal(v.val).PathVal(bw.ValPath{vpi}); err == nil {
		result = valPath{val, append(v.path, vpi)}
	}
	return
}

func (v valPath) Idx(idx int) (result valPath, err error) {
	var val interface{}
	vpi := bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: idx}
	if val, err = FromVal(v.val).PathVal(bw.ValPath{vpi}); err == nil {
		result = valPath{val, append(v.path, vpi)}
	}
	return
}

func (v valPath) notOfTypeError(expectedType string, optExpectedType ...string) (result error) {
	if len(optExpectedType) == 0 {
		result = bwerr.From(v.ansiString()+ansi.String(" is not <ansiType>%s"), expectedType)
	} else {
		expectedTypes := fmt.Sprintf(ansiType, expectedType)
		for i, elem := range optExpectedType {
			expectedTypes += typeSeparator[i == len(optExpectedType)-1] + fmt.Sprintf(ansiType, elem)
		}
		result = bwerr.From(v.ansiString()+ansi.String(" is none of %s"), expectedTypes)
	}
	return
}

func (v valPath) notEnoughRangeError(l int, idx int) error {
	return bwerr.From(
		v.ansiString()+
			ansi.String(" has not enough length (<ansiVal>%d<ansi>) for idx (<ansiVal>%d)"),
		l, idx,
	)
}

func (v valPath) nonSupportedValueError() error {
	return bwerr.From(v.ansiString() + ansi.String(" is <ansiErr>non supported<ansi> value"))
}

// func nilError(path bw.ValPath) error {
// 	return bwerr.From(ansi.String("<ansiPath>%s<ansi> is <ansiErr>nil"), path)
// }

func readonlyPathError(path bw.ValPath) error {
	return bwerr.From(ansi.String("<ansiPath>%s<ansi> is <ansiErr>readonly path"), path)
}

type RangeKindValue uint8

const (
	RangeNo RangeKindValue = iota
	RangeMin
	RangeMax
	RangeMinMax
)

func RangeKind(v Range) (result RangeKindValue) {
	if v != nil {
		if v.Min() != nil {
			if v.Max() != nil {
				result = RangeMinMax
			} else {
				result = RangeMin
			}
		} else if v.Max() != nil {
			result = RangeMax
		}
	}
	return
}

func RangeString(v Range) (result string) {
	if v != nil {
		switch RangeKind(v) {
		case RangeMinMax:
			result = fmt.Sprintf("%s..%s", bwjson.Pretty(v.Min()), bwjson.Pretty(v.Max()))
		case RangeMin:
			result = fmt.Sprintf("%s..", bwjson.Pretty(v.Min()))
		case RangeMax:
			result = fmt.Sprintf("..%s", bwjson.Pretty(v.Max()))
		}
	}
	return
}

func RangeMarshalJSON(v Range) ([]byte, error) {
	return []byte(RangeString(v)), nil
	// return json.Marshal(RangeString(v))
}

// func ValInRange(val interface)

// Range - интерфейс для IntRange/NumberRange
type Range interface {
	ValKind() ValKind
	Min() interface{}
	Max() interface{}
}

type IntRange struct {
	MinPtr *int
	MaxPtr *int
}

func (v IntRange) ValKind() ValKind {
	return ValInt
}

func (v IntRange) Min() (result interface{}) {
	if v.MinPtr != nil {
		result = *v.MinPtr
	}
	return
}

func (v IntRange) Max() (result interface{}) {
	if v.MaxPtr != nil {
		result = *v.MaxPtr
	}
	return
}

func (v IntRange) MarshalJSON() ([]byte, error) {
	return RangeMarshalJSON(v)
}

type NumberRange struct {
	MinPtr *float64
	MaxPtr *float64
}

func (v NumberRange) ValKind() ValKind {
	return ValNumber
}

func (v NumberRange) Min() (result interface{}) {
	if v.MinPtr != nil {
		result = *v.MinPtr
	}
	return
}

func (v NumberRange) Max() (result interface{}) {
	if v.MaxPtr != nil {
		result = *v.MaxPtr
	}
	return
}

func (v NumberRange) MarshalJSON() ([]byte, error) {
	return RangeMarshalJSON(v)
}

func (v valPath) inRange(rng Range) (result bool) {
	rangeKind := RangeKind(rng)
	if rangeKind == RangeNo {
		result = true
	} else if _, vk := Kind(v.val); vk == ValNumber || rng.ValKind() == ValNumber {
		if n, ok := Number(v.val); ok {
			switch rangeKind {
			case RangeMinMax:
				result = MustNumber(rng.Min()) <= n && n <= MustNumber(rng.Max())
			case RangeMin:
				result = MustNumber(rng.Min()) <= n
			case RangeMax:
				result = n <= MustNumber(rng.Max())
			}
		}
	} else {
		if n, ok := Int(v.val); ok {
			switch rangeKind {
			case RangeMinMax:
				result = MustInt(rng.Min()) <= n && n <= MustInt(rng.Max())
			case RangeMin:
				result = MustInt(rng.Min()) <= n
			case RangeMax:
				result = n <= MustInt(rng.Max())
			}
		}
	}
	return
}

func (v valPath) outOfRangeError(rng Range) (err error) {
	var s string
	switch RangeKind(rng) {
	case RangeMinMax:
		s = ansi.String(" is <ansiErr>out of range<ansi> <ansiVal>%s")
	case RangeMin:
		s = ansi.String(" is <ansiErr>less<ansi> than<ansiVal>%s")
	case RangeMax:
		s = ansi.String(" is <ansiErr>more<ansi> than<ansiVal>%s")
	}
	if len(s) > 0 {
		err = bwerr.From(v.ansiString()+s, RangeString(rng))
	}
	return
}

func (v valPath) maxLessThanMinError() error {
	return bwerr.From(v.ansiString() + "<ansiVar>max<ansi> must not be less then <ansiVar>min")
}

func (v valPath) unexpectedKeysError(unexpectedKeys []string) (err error) {
	var fmtString string
	var fmtArg interface{}
	switch len(unexpectedKeys) {
	case 0:
	case 1:
		fmtString = ansi.String(`has unexpected key <ansiVal>%s`)
		fmtArg = unexpectedKeys[0]
	default:
		fmtString = `has unexpected keys <ansiVal>%s`
		fmtArg = bwjson.Pretty(unexpectedKeys)
	}
	if len(fmtString) > 0 {
		err = bwerr.From(v.ansiString()+fmtString, fmtArg)
	}
	return
}

func (v valPath) forEachMapString(f func(k string, v interface{}) (err error)) (err error) {
	if !_isOfType(v.val, "map[string]") {
		// err =
		return v.notOfTypeError("map[string]")
		// return pathValIsNotOfType(v.path, v.val, "map[string]")
		// return bwerr.From(ansi.String("%s is not of type <ansiType>%s"), v, "map[string]")
		// err = valueErrorMake(v, valueErrorIsNotOfType, "map[string]")
	} else {
		mv := reflect.ValueOf(v.val)
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

func (v valPath) forEachSlice(f func(i int, v interface{}) (err error)) (err error) {
	if !_isOfType(v.val, "[]") {
		// err = valueErrorMake(v, valueErrorIsNotOfType, "[]")
		return bwerr.From(ansi.String("%s is not of type <ansiType>%s"), v, "[]")
	} else {
		sliceValue := reflect.ValueOf(v.val)
		for i := 0; i < sliceValue.Len(); i++ {
			err = f(i, sliceValue.Index(i).Interface())
			if err != nil {
				break
			}
		}
	}
	return err
}

func (v valPath) getElem(elemIndex int, opts ...interface{}) (result valPath, err error) {
	defaultValue, ofType := getDefaultValueAndOfTypeFromOpts(opts)
	if v.val == nil {
		// err = valueErrorMake(v, valueErrorIsNotOfType, "array")
		err = bwerr.From(ansi.String("%s is not of type <ansiType>%s"), v, "array")
		return
	} else {
		vType := reflect.TypeOf(v.val)
		if vType.Kind() != reflect.Slice {
			err = bwerr.From(ansi.String("%s is not of type <ansiType>%s"), v, "array")
			// err = valueErrorMake(v, valueErrorIsNotOfType, "array")
		} else {
			sv := reflect.ValueOf(v.val)
			result.path = append(v.path, bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: elemIndex})
			// result.what = v.what + fmt.Sprintf(".#%d", elemIndex)
			if 0 <= elemIndex && elemIndex < sv.Len() {
				err = checkElemIsOfType(&result, sv.Index(elemIndex), ofType)
			} else if defaultValue == nil {
				err = bwerr.From(ansi.String("%s has not idx <ansiErr>%d"), v, elemIndex)
				// err = valueErrorMake(v, valueErrorHasNoKey, fmt.Sprintf("#%d", elemIndex))
			} else {
				result.val = *defaultValue
			}
		}
	}
	return
}

func (v valPath) setKey(keyName string, keyValue interface{}) (err error) {
	if !_isOfType(v.val, "map[string]") {
		err = bwerr.From(ansi.String("%s is not of type <ansiType>%s"), v, "map[string]")
		// err = valueErrorMake(v, valueErrorIsNotOfType, "map[string]")
	} else {
		mv := reflect.ValueOf(v.val)
		mv.SetMapIndex(reflect.ValueOf(keyName), reflect.ValueOf(keyValue))
	}
	return
}

func (v valPath) getKey(keyName string, opts ...interface{}) (result valPath, err error) {
	defaultValue, ofType := getDefaultValueAndOfTypeFromOpts(opts)
	if !_isOfType(v.val, "map[string]") {
		err = bwerr.From(ansi.String("%s is not of type <ansiType>%s"), v, "map[string]")
		// err = valueErrorMake(v, valueErrorIsNotOfType, "map[string]")
	} else {
		mv := reflect.ValueOf(v.val)
		// result.what = v.what + "." + keyName
		result.path = append(v.path, bw.ValPathItem{Type: bw.ValPathItemKey, Key: keyName})
		elem := mv.MapIndex(reflect.ValueOf(keyName))
		zeroValue := reflect.Value{}
		if elem != zeroValue {
			err = checkElemIsOfType(&result, elem, ofType)
		} else if defaultValue == nil {
			err = bwerr.From(ansi.String("%s has no key <ansiErr>%s"), v, keyName)
			// err = valueErrorMake(v, valueErrorHasNoKey, keyName)
		} else {
			result.val = *defaultValue
		}
	}
	return
}

func checkElemIsOfType(result *valPath, elem reflect.Value, ofType []string) (err error) {
	result.val = elem.Interface()
	if len(ofType) > 0 && !_isOfType(result.val, ofType...) {
		// valPath{result.val, }
		err = result.notOfTypeError(ofType[0], ofType[1:]...)
		// err = pathValIsNotOfType(result.path, result.val, ofType[0], ofType[1:]...)
		// var ofTypeIntfs = []interface{}{}
		// for _, i := range ofType {
		// 	ofTypeIntfs = append(ofTypeIntfs, i)
		// }
		// err = valueErrorMake(*result, valueErrorIsNotOfType, ofTypeIntfs...)
	}
	return
}

func getDefaultValueAndOfTypeFromOpts(opts []interface{}) (defaultValue *interface{}, ofType []string) {
	if opts != nil {
		if _isOfType(opts[0], "string") {
			ofType = []string{MustString(opts[0])}
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
			bwerr.Panic("expects max 2 opts (ofTypes, defaultValue), but found <ansiVal>%v", opts)
		}
	}
	return
}

// ============================================================================

func compileDef(def valPath) (result *Def, err error) {
	// bwdebug.Print("def.val", def.val)
	if def.val == nil {
		err = valAtPathIsNil(def.path)
		return
	}
	var defType valPath
	var isSimple bool
	validDefKeys := bwset.String{}
	switch _, kind := Kind(def.val); kind {
	case ValString, ValArray:
		isSimple = true
		defType = def
	case ValMap:
		validDefKeys.Add("type")
		// bwdebug.Print("def", def)
		if defType, err = def.Key("type"); err != nil {
			return
		}
		// bwdebug.Print("defType", defType)
		switch _, kind := Kind(defType.val); kind {
		case ValString, ValArray:
		default:
			// bwdebug.Print("deftype", defType)
			err = defType.notOfTypeError("String", "Array")
			return
		}
	default:
		err = def.notOfTypeError("String", "Array", "Map")
		return
	}

	var types deftype.Set
	if types, err = getDeftype(defType, isSimple); err != nil {
		return
	}

	result = &Def{Types: types}
	if !isSimple {
		if types.Has(deftype.String) {
			var vp valPath
			var ss []string
			validDefKeys.Add("enum")
			if vp, err = def.Key("enum"); err != nil || vp.val == nil {
				err = nil
			} else if ss, err = vp.ArrayOfString(); err != nil {
				return
			} else {
				result.Enum = bwset.StringFromSlice(ss)
			}
		}
		if types.Has(deftype.Map) {
			var keysVp valPath
			validDefKeys.Add("keys")
			var m map[string]interface{}
			if keysVp, err = def.Key("keys"); err != nil || keysVp.val == nil {
				err = nil
			} else if m, err = keysVp.Map(); err != nil {
				return
			} else {
				result.Keys = map[string]Def{}
				var vp valPath
				var keyDef *Def
				for k := range m {
					vp, _ = keysVp.Key(k)
					if keyDef, err = compileDef(vp); err != nil {
						return
					} else {
						result.Keys[k] = *keyDef
					}
				}
			}
		}
		if types.Has(deftype.Array) {
			var vp valPath
			validDefKeys.Add("arrayElem")
			if vp, err = def.Key("arrayElem"); err != nil || vp.val == nil {
				err = nil
			} else if result.ArrayElem, err = compileDef(vp); err != nil {
				return
			}
		}
		if types.Has(deftype.Map) || types.Has(deftype.Array) && result.ArrayElem == nil {
			if elemVal, _ := getDefKey(&validDefKeys, def, "elem", "interface{}", nil); elemVal.val != nil {
				if result.Elem, err = compileDef(elemVal); err != nil {
					return nil, err
				}
			}
		}
		if types.Has(deftype.Int) {
			var vp valPath
			validDefKeys.Add("min", "max")
			rng := IntRange{}
			var limCount, n int
			if vp, err = def.Key("min"); err != nil || vp.val == nil {
				err = nil
			} else if n, err = vp.Int(); err != nil {
				return
			} else {
				rng.MinPtr = PtrToInt(n)
				limCount++
			}
			if vp, err = def.Key("max"); err != nil || vp.val == nil {
				err = nil
			} else if n, err = vp.Int(); err != nil {
				return
			} else {
				rng.MaxPtr = PtrToInt(n)
				limCount++
			}
			if limCount == 2 && *(rng.MinPtr) > *(rng.MaxPtr) {
				err = def.maxLessThanMinError()
				return
			}
			result.Range = rng
		} else if types.Has(deftype.Number) {
			var vp valPath
			validDefKeys.Add("min", "max")
			rng := NumberRange{}
			var limCount, n float64
			if vp, err = def.Key("min"); err != nil || vp.val == nil {
				err = nil
			} else if n, err = vp.Number(); err != nil {
				return
			} else {
				rng.MinPtr = PtrToNumber(n)
				limCount++
			}
			if vp, err = def.Key("max"); err != nil || vp.val == nil {
				err = nil
			} else if n, err = vp.Number(); err != nil {
				return
			} else {
				rng.MaxPtr = PtrToNumber(n)
				limCount++
			}
			if limCount == 2 && *(rng.MinPtr) > *(rng.MaxPtr) {
				err = def.maxLessThanMinError()
				return
			}
			result.Range = rng
		}
		var vp valPath
		validDefKeys.Add("default")
		if vp, err = def.Key("default"); err != nil || vp.val == nil {
			err = nil
		} else {
			// bwdebug.Print("vp", vp)
			dfltDef := *result
			if types.Has(deftype.ArrayOf) {
				dfltDef = Def{
					Types: deftype.From(deftype.Array),
					ArrayElem: &Def{
						Types:      result.Types.Copy(),
						IsOptional: false,
						Enum:       result.Enum,
						Range:      result.Range,
						Keys:       result.Keys,
						Elem:       result.Elem,
					},
				}
				dfltDef.ArrayElem.Types.Del(deftype.ArrayOf)
			}
			if result.Default, err = vp.ValidVal(dfltDef); err != nil {
				return nil, err
			}
		}

		// if dfltVal, _ := getDefKey(&validDefKeys, def, "default", "interface{}", nil); dfltVal.val != nil {
		// 	dfltDef := *result
		// 	if types.Has(deftype.ArrayOf) {
		// 		dfltDef = Def{
		// 			Types: deftype.From(deftype.Array),
		// 			ArrayElem: &Def{
		// 				Types:      result.Types.Copy(),
		// 				IsOptional: false,
		// 				Enum:       result.Enum,
		// 				Range:      result.Range,
		// 				Keys:       result.Keys,
		// 				Elem:       result.Elem,
		// 			},
		// 		}
		// 		dfltDef.ArrayElem.Types.Del(deftype.ArrayOf)
		// 	}
		// 	if result.Default, err = getValidVal(dfltVal, dfltDef); err != nil {
		// 		return nil, err
		// 	}
		// }
		var boolVal valPath
		if boolVal, err = getDefKey(&validDefKeys, def, "isOptional", "bool", result.Default != nil); err != nil {
			return nil, err
		}
		result.IsOptional = MustBool(boolVal.val)
		if !result.IsOptional && result.Default != nil {
			err = bwerr.From(
				ansi.String(
					"<ansiPath>%s<ansi> (<ansiVal>%s<ansi>) has conflicting keys: <ansiErr>isOptional<ansi> and <ansiErr>default",
				),
				def.path, bwjson.Pretty(result),
			)
			return
			// return nil, valueErrorMake(def, valueErrorConflictingKeys, map[string]interface{}{
			// 	"isOptional": result.IsOptional,
			// 	"default":    result.Default,
			// })
		}
		if unexpectedKeys := bwmap.MustUnexpectedKeys(def.val, validDefKeys); unexpectedKeys != nil {
			err = bwerr.From(
				ansi.String(
					"<ansiPath>%s<ansi> (<ansiVal>%s<ansi>) has conflicting keys: <ansiErr>isOptional<ansi> and <ansiErr>default<ansi>",
				),
				def.path, bwjson.Pretty(result),
			)
			return
			// return nil, valueErrorMake(def, valueErrorHasUnexpectedKeys, unexpectedKeys)
		}
	}
	return
}

func getDefKey(validDefKeys *bwset.String, def valPath, keyName string, ofType interface{}, defaultValue ...interface{}) (keyValue valPath, err error) {
	opts := []interface{}{ofType}
	if defaultValue != nil {
		opts = append(opts, defaultValue[0])
	}
	keyValue, err = def.getKey(keyName, opts...)
	validDefKeys.Add(keyName)
	return
}

func getDeftype(defType valPath, isSimple bool) (result deftype.Set, err error) {
	var ss []string
	var isString bool
	switch val, kind := Kind(defType.val); kind {
	case ValString:
		ss = []string{MustString(val)}
		isString = true
	case ValArray:
		if ss, err = defType.ArrayOfString(); err != nil {
			return
		}
	}
	result = deftype.Set{}
	for i, s := range ss {
		var tpItem deftype.Item
		if tpItem, err = deftype.ItemFromString(s); err == nil {
			result.Add(tpItem)
		} else {
			elem := defType
			if !isString {
				elem, _ = defType.Idx(i)
			}
			err = elem.nonSupportedValueError()
			return
		}
	}
	if result.Has(deftype.ArrayOf) {
		if len(result) < 2 {
			err = bwerr.From(ansi.String("%s: <ansiVal>ArrayOf<ansi> must be followed by some type"), defType)
		} else if result.Has(deftype.Array) {
			err = bwerr.From(ansi.String("%s type values <ansiVal>ArrayOf<ansi> and <ansiVal>Array<ansi> can not be combined"), defType)
		}
	}
	if result.Has(deftype.Int) && result.Has(deftype.Number) {
		err = bwerr.From(ansi.String("%s type values <ansiVal>Int<ansi> and <ansiVal>Number<ansi> can not be combined"), defType)
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
			case "[]":
				ok = vType.Kind() == reflect.Slice
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
				case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
					ok = true
				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
					ok = reflect.ValueOf(v).Uint() <= uint64(bw.MaxInt64)
				default:
					ok = false
				}
			case "float64":
				switch vType.Kind() {
				case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
					ok = true
				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
					ok = true
				case reflect.Float32, reflect.Float64:
					ok = true
				default:
					ok = false
				}
			case "bool":
				ok = vType.Kind() == reflect.Bool
			case "bwset.Strings":
				_, ok = v.(bwset.String)
			case "deftype.Set":
				_, ok = v.(deftype.Set)
			case "interface{}":
				ok = true
			default:
				bwerr.Panic("unsupported type <ansiVal>%s", ofType)
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
		bwerr.Panic("<ansiVal>%+v<ansi> is not of types <ansiVal>%v", v, ofTypes)
	}
	return v
}

// func _mustBeString(v interface{}) (result string) {
// 	result, _ = _mustBeOfType(v, "string").(string)
// 	return
// }

// func _mustBeBool(v interface{}) (result bool) {
// 	result, _ = _mustBeOfType(v, "bool").(bool)
// 	return
// }

// func _mustBeInt(v interface{}) (result int) {
// 	vValue := reflect.ValueOf(v)
// 	switch vValue.Kind() {
// 	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
// 		result = vValue.Int()
// 	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
// 		if reflect.ValueOf(v).Uint() <= uint64(bw.MaxInt64) {
// 			result = int64(vValue.Uint())
// 		} else {
// 			bwerr.Panic("<ansiVal>%+v<ansi> is not of type <ansiVal>int64", v)
// 		}
// 	default:
// 		// log.Printf("vValue.Kind(): %s", vValue.Kind())
// 		bwerr.Panic("<ansiVal>%+v<ansi> is not of type <ansiVal>int64", v)
// 	}
// 	return
// }

// func _mustBeFloat64(v interface{}) (result float64) {
// 	vValue := reflect.ValueOf(v)
// 	switch vValue.Kind() {
// 	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
// 		result = float64(vValue.Int())
// 	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
// 		result = float64(vValue.Uint())
// 	case reflect.Float32, reflect.Float64:
// 		result = vValue.Float()
// 	default:
// 		bwerr.Panic("<ansiVal>%+v<ansi> is not of type <ansiVal>float64", v)
// 	}
// 	return
// }

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

// func _mustBeBwsetStrings(v interface{}) (result bwset.String) {
// 	result, _ = _mustBeOfType(v, "bwset.Strings").(bwset.String)
// 	return
// }

// func _mustBeDeftypeSet(v interface{}) (result deftype.Set) {
// 	result, _ = _mustBeOfType(v, "deftype.Set").(deftype.Set)
// 	return
// }

func PtrToInt(i int) *int {
	return &i
}

func PtrToNumber(i float64) *float64 {
	return &i
}

// ============================================================================

func (v valPath) ValidVal(def Def, optSkipArrayOf ...bool) (result interface{}, err error) {
	skipArrayOf := optSkipArrayOf != nil && optSkipArrayOf[0]
	if v.val == nil {
		if !skipArrayOf {
			if def.Default != nil {
				return def.Default, nil
			}
			if def.IsOptional {
				return nil, nil
			}
		}
		if def.Types.Has(deftype.Map) {
			v.val = map[string]interface{}{}
		} else {
			ss := def.Types.ToSliceOfStrings()
			err = v.notOfTypeError(ss[0], ss[1:]...)
			return
			// err = pathValIsNotOfType(path, val, expectedType, optExpectedType)
			// return nil, valueErrorMake(val, valueErrorIsNotOfType, def.Types)
		}
	}
	var valDeftype deftype.Item
	valType := reflect.TypeOf(v.val)
	switch valType.Kind() {
	case reflect.Bool:
		if def.Types.Has(deftype.Bool) {
			valDeftype = deftype.Bool
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if def.Types.Has(deftype.Int) {
			valDeftype = deftype.Int
		} else if def.Types.Has(deftype.Number) {
			valDeftype = deftype.Number
		}
	case reflect.Float32, reflect.Float64:
		if def.Types.Has(deftype.Number) {
			valDeftype = deftype.Number
		}
	case reflect.Map:
		if valType.Key().Kind() == reflect.String && valType.Elem().Kind() == reflect.Interface && def.Types.Has(deftype.Map) {
			valDeftype = deftype.Map
		}
	case reflect.Slice:
		if def.Types.Has(deftype.Array) {
			valDeftype = deftype.Array
		} else if !skipArrayOf && def.Types.Has(deftype.ArrayOf) {
			valDeftype = deftype.ArrayOf
		}

	case reflect.String:
		if def.Types.Has(deftype.String) {
			valDeftype = deftype.String
		}
	}
	if valDeftype == deftype.Unknown {
		ss := def.Types.ToSliceOfStrings()
		err = v.notOfTypeError(ss[0], ss[1:]...)
		return
		// return nil, valueErrorMake(val, valueErrorIsNotOfType, def.Types)
	}

	// bwdebug.Print("valDeftype", valDeftype, "getValidValHelpers", getValidValHelpers)
	if v.val, err = getValidValHelpers[valDeftype](v, def); err != nil {
		return nil, err
	}

	if !skipArrayOf && valDeftype != deftype.ArrayOf && def.Types.Has(deftype.ArrayOf) {
		v.val = []interface{}{v.val}
	}

	return v.val, nil
}

type getValidValHelper func(val valPath, def Def) (result interface{}, err error)

var getValidValHelpers map[deftype.Item]getValidValHelper

func init() {
	getValidValHelpers = map[deftype.Item]getValidValHelper{
		deftype.Bool:    _Bool,
		deftype.String:  _String,
		deftype.Int:     _Int,
		deftype.Number:  _Number,
		deftype.Map:     _Map,
		deftype.Array:   _Array,
		deftype.ArrayOf: _ArrayOf,
	}
	getValidValHelpersCheck()
}

func getValidValHelpersCheck() {
	for deftypeItem := deftype.Unknown + 1; deftypeItem <= deftype.ArrayOf; deftypeItem += 1 {
		if _, ok := getValidValHelpers[deftypeItem]; !ok {
			bwerr.Panic("not defined <ansiVar>deftype.ItemValidators<ansi>[<ansiVal>%s<ansi>]", deftypeItem)
		}
		deftypeItem += 1
	}
}

func _Bool(val valPath, def Def) (result interface{}, err error) {
	result = val.val
	return
}

func _String(vp valPath, def Def) (result interface{}, err error) {
	if def.Enum != nil {
		if !def.Enum.Has(MustString(vp.val)) {
			bwerr.TODO() // enum instead nonSupportedValueError
			err = vp.nonSupportedValueError()
			// err = valueErrorMake(val, valueErrorHasNonSupportedValue)
		}
	}
	result = vp.val
	return
}

func _Int(vp valPath, def Def) (result interface{}, err error) {
	if !vp.inRange(def.Range) {
		err = vp.outOfRangeError(def.Range)
	}
	result = vp.val
	return
}

func _Number(vp valPath, def Def) (result interface{}, err error) {
	if !vp.inRange(def.Range) {
		err = vp.outOfRangeError(def.Range)
	}
	result = vp.val
	return
}

func _Map(vp valPath, def Def) (result interface{}, err error) {
	if def.Keys != nil {
		unexpectedKeys := bwmap.MustUnexpectedKeys(vp.val, def.Keys)
		for key, keyDef := range def.Keys {
			if err = _MapHelper(vp, key, keyDef); err != nil {
				return
			}
		}
		if unexpectedKeys != nil {
			if def.Elem == nil {
				err = vp.unexpectedKeysError(unexpectedKeys)
				return
				// return nil, valueErrorMake(val, valueErrorHasUnexpectedKeys, unexpectedKeys)
			} else {
				for _, key := range unexpectedKeys {
					if err = _MapHelper(vp, key, *(def.Elem)); err != nil {
						return
					}
				}
			}
		}
	} else if def.Elem != nil {
		err = vp.forEachMapString(func(k string, v interface{}) error {
			return _MapHelper(vp, k, *(def.Elem))
		})
	}
	result = vp.val
	return
}

func _MapHelper(val valPath, key string, elemDef Def) error {
	elemVal, _ := val.getKey(key)
	if elemValIntf, err := elemVal.ValidVal(elemDef); err != nil {
		return err
	} else if elemValIntf != nil {
		if err := val.setKey(key, elemValIntf); err != nil {
			return err
		}
	}
	return nil
}

func _Array(val valPath, def Def) (result interface{}, err error) {
	elemDef := def.ArrayElem
	if elemDef == nil {
		elemDef = def.Elem
	}
	if elemDef == nil {
		result = val.val
	} else {
		result, err = _ArrayHelper(val, *elemDef)
	}
	return
}

func _ArrayOf(val valPath, def Def) (result interface{}, err error) {
	return _ArrayHelper(val, def, true)
}

func _ArrayHelper(val valPath, elemDef Def, optSkipArrayOf ...bool) (result interface{}, err error) {
	newSliceValue := reflect.MakeSlice(reflect.TypeOf(val.val), 0, reflect.ValueOf(val.val).Len())
	err = val.forEachSlice(func(i int, v interface{}) (err error) {
		var elemVal valPath
		if elemVal, err = val.getElem(i); err == nil {
			var elemValIntf interface{}
			if elemValIntf, err = elemVal.ValidVal(elemDef, optSkipArrayOf...); err == nil && elemValIntf != nil {
				// log.Printf("elemValIntf: %#v, val.val: %#v", elemValIntf, val.val)
				newSliceValue = reflect.Append(newSliceValue, reflect.ValueOf(elemValIntf))
			}
		}
		return
	})
	return newSliceValue.Interface(), err
}

// ============================================================================
