package bwval

import (
	"encoding/json"
	"fmt"

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
	ansiValAtPathIsNil = ansi.String(valPathPrefix + "is <ansiErr>nil")

	ansiType = ansi.String("<ansiType>%s")
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
					err = valPath{val, vpi.Path}.notOfTypeError("Int", "String")
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

func valAtPathIsNil(path bw.ValPath) error {
	return bwerr.From(ansiValAtPathIsNil, path)
}

func (v *valHolder) getArray(idx int, result interface{}, resultPath bw.ValPath) ([]interface{}, int, error) {
	var err error
	var ok bool
	var vals []interface{}
	if vals, ok = Array(result); !ok {
		err = valPath{result, resultPath}.notOfTypeError("Array")
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

func (v valPath) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{}
	result["val"] = v.val
	result["path"] = v.path
	return json.Marshal(result)
}

func (v valPath) ansiString() string {
	return fmt.Sprintf(ansi.String("<ansiPath>%s<ansi> (<ansiVal>%s<ansi>)"), v.path, bwjson.Pretty(v.val))
}

func (v valPath) Bool() (result bool, err error) {
	var ok bool
	if result, ok = Bool(v.val); !ok {
		err = v.notOfTypeError("Bool")
	}
	return
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

func (v valPath) SetKey(val interface{}, key string) (err error) {
	vpi := bw.ValPathItem{Type: bw.ValPathItemKey, Key: key}
	err = FromVal(v.val).SetPathVal(val, bw.ValPath{vpi})
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
	return json.Marshal(RangeString(v))
}

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

// ============================================================================

func compileDef(def valPath) (result *Def, err error) {
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
		if defType, err = def.Key("type"); err != nil {
			return
		}
		switch _, kind := Kind(defType.val); kind {
		case ValString, ValArray:
		default:
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
		var vp valPath
		if types.Has(deftype.String) {
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
			validDefKeys.Add("arrayElem")
			if vp, err = def.Key("arrayElem"); err != nil || vp.val == nil {
				err = nil
			} else if result.ArrayElem, err = compileDef(vp); err != nil {
				return
			}
		}
		if types.Has(deftype.Map) || types.Has(deftype.Array) && result.ArrayElem == nil {
			validDefKeys.Add("elem")
			if vp, err = def.Key("elem"); err != nil || vp.val == nil {
				err = nil
			} else if result.Elem, err = compileDef(vp); err != nil {
				return
			}
		}
		if types.Has(deftype.Int) {
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
			if limCount > 0 {
				result.Range = rng
			}
		} else if types.Has(deftype.Number) {
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
			if limCount > 0 {
				result.Range = rng
			}
		}
		validDefKeys.Add("default")
		if vp, err = def.Key("default"); err != nil || vp.val == nil {
			err = nil
		} else {
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

		validDefKeys.Add("isOptional")
		if vp, err = def.Key("isOptional"); err != nil || vp.val == nil {
			err = nil
			result.IsOptional = result.Default != nil
		} else if result.IsOptional, err = vp.Bool(); err != nil {
			return
		}
		if !result.IsOptional && result.Default != nil {
			err = bwerr.From(
				ansi.String(
					"<ansiPath>%s<ansi> (<ansiVal>%s<ansi>) has conflicting keys: <ansiErr>isOptional<ansi> and <ansiErr>default",
				),
				def.path, bwjson.Pretty(result),
			)
			return
		}
		if unexpectedKeys := bwmap.MustUnexpectedKeys(def.val, validDefKeys); unexpectedKeys != nil {
			err = bwerr.From(
				ansi.String(
					"<ansiPath>%s<ansi> (<ansiVal>%s<ansi>) has conflicting keys: <ansiErr>isOptional<ansi> and <ansiErr>default<ansi>",
				),
				def.path, bwjson.Pretty(result),
			)
			return
		}
	}
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
		}
	}
	var valDeftype deftype.Item
	switch _, kind := Kind(v.val); kind {
	case ValBool:
		if def.Types.Has(deftype.Bool) {
			valDeftype = deftype.Bool
		}
	case ValInt:
		if def.Types.Has(deftype.Int) {
			valDeftype = deftype.Int
		} else if def.Types.Has(deftype.Number) {
			valDeftype = deftype.Number
		}
	case ValNumber:
		if def.Types.Has(deftype.Number) {
			valDeftype = deftype.Number
		}
	case ValMap:
		if def.Types.Has(deftype.Map) {
			valDeftype = deftype.Map
		}
	case ValArray:
		if def.Types.Has(deftype.Array) {
			valDeftype = deftype.Array
		} else if !skipArrayOf && def.Types.Has(deftype.ArrayOf) {
			valDeftype = deftype.ArrayOf
		}
	case ValString:
		if def.Types.Has(deftype.String) {
			valDeftype = deftype.String
		}
	}

	if valDeftype == deftype.Unknown {
		ss := def.Types.ToSliceOfStrings()
		err = v.notOfTypeError(ss[0], ss[1:]...)
		return
	}

	switch valDeftype {
	case deftype.Bool:
	case deftype.String:
		if def.Enum != nil {
			if !def.Enum.Has(MustString(v.val)) {
				bwerr.TODO() // enum instead nonSupportedValueError
				err = v.nonSupportedValueError()
				return
			}
		}
	case deftype.Int, deftype.Number:
		if !v.inRange(def.Range) {
			err = v.outOfRangeError(def.Range)
			return
		}
	case deftype.Map:
		if def.Keys != nil {
			unexpectedKeys := bwmap.MustUnexpectedKeys(v.val, def.Keys)
			for key, keyDef := range def.Keys {
				if err = v.mapHelper(key, keyDef); err != nil {
					return
				}
			}
			if unexpectedKeys != nil {
				if def.Elem == nil {
					err = v.unexpectedKeysError(unexpectedKeys)
					return
				} else {
					for _, key := range unexpectedKeys {
						if err = v.mapHelper(key, *(def.Elem)); err != nil {
							return
						}
					}
				}
			}
		} else if def.Elem != nil {
			m, _ := Map(v.val)
			for k := range m {
				if err = v.mapHelper(k, *(def.Elem)); err != nil {
					return
				}
			}
		}
	case deftype.Array:
		elemDef := def.ArrayElem
		if elemDef == nil {
			elemDef = def.Elem
		}
		if elemDef != nil {
			if v.val, err = v.arrayHelper(*elemDef); err != nil {
				return
			}
		}
	case deftype.ArrayOf:
		if v.val, err = v.arrayHelper(def, true); err != nil {
			return
		}
	}

	if !skipArrayOf && valDeftype != deftype.ArrayOf && def.Types.Has(deftype.ArrayOf) {
		v.val = []interface{}{v.val}
	}

	result = v.val
	return
}

func (v valPath) mapHelper(key string, elemDef Def) (err error) {
	vp, _ := v.Key(key)
	var val interface{}
	if val, err = vp.ValidVal(elemDef); err != nil {
		return
	} else if val != nil {
		if err = v.SetKey(val, key); err != nil {
			return
		}
	}
	return
}

// func (v valPath) arrayHelper(elemDef Def, optSkipArrayOf ...bool) (err error) {
// 	var vp valPath
// 	for i := range MustArray(v.val) {
// 		vp, _ = v.Idx(i)
// 		if _, err = vp.ValidVal(elemDef, optSkipArrayOf...); err != nil {
// 			return
// 		}
// 	}
// 	return
// }

func (v valPath) arrayHelper(elemDef Def, optSkipArrayOf ...bool) (result interface{}, err error) {
	arr := MustArray(v.val)
	newArr := make([]interface{}, 0, len(arr))
	var vp valPath
	for i := range arr {
		vp, _ = v.Idx(i)
		var val interface{}
		if val, err = vp.ValidVal(elemDef, optSkipArrayOf...); err != nil {
			return
		}
		newArr = append(newArr, val)
	}
	result = newArr
	return
}

// ============================================================================
