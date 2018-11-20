// Package bwval реализует интерфейc bw.Val и утилиты для работы с этим интерфейсом.
package bwval

import (
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwparse"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwtype"
)

// ============================================================================

// PathFrom - конструктор-парсер bw.ValPath из строки
func PathFrom(s string, optBases ...[]bw.ValPath) (result bw.ValPath, err error) {
	p := bwparse.From(bwrune.ProviderFromString(s))
	a := bwparse.PathA{}
	if len(optBases) > 0 {
		a.Bases = optBases[0]
	}
	if result, err = p.PathContent(a); err == nil {
		err = p.SkipSpace(bwparse.TillEOF)
	}
	return
	return
}

// PathFrom - конструктор-парсер bw.ValPath из строки
func MustPathFrom(s string, optBases ...[]bw.ValPath) (result bw.ValPath) {
	var err error
	if result, err = PathFrom(s, optBases...); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return
}

// ============================================================================

// MustPathVal - must-обертка bw.Val.PathVal()
func MustPathVal(v bw.Val, path bw.ValPath, optVars ...map[string]interface{}) (result interface{}) {
	var err error
	if result, err = v.PathVal(path, optVars...); err != nil {
		bwerr.PanicA(bwerr.Err(bwerr.Refine(err,
			ansiMustPathValFailed,
			path, bwjson.Pretty(v), varsJSON(path, optVars),
		)))
	}
	return result
}

// MustSetPathVal - must-обертка bw.Val.SetPathVal()
func MustSetPathVal(val interface{}, v bw.Val, path bw.ValPath, optVars ...map[string]interface{}) {
	var err error
	if err = v.SetPathVal(val, path, optVars...); err != nil {
		bwerr.PanicA(bwerr.Err(bwerr.Refine(err,
			ansiMustSetPathValFailed,
			path, bwjson.Pretty(v), varsJSON(path, optVars),
		)))
	}
}

// ============================================================================

type ValByOpt struct {
	Bases []bw.ValPath
	Vars  map[string]interface{}
}

// ValBy - обертка над PathFrom()+bw.Val.PathVal()
func ValBy(v bw.Val, pathStr string, opt ...ValByOpt) (result interface{}, err error) {
	var (
		optBases [][]bw.ValPath
		optVars  []map[string]interface{}
		path     bw.ValPath
	)
	if len(opt) > 0 {
		optBases = append(optBases, opt[0].Bases)
		optVars = append(optVars, opt[0].Vars)
	}
	if path, err = PathFrom(pathStr, optBases...); err != nil {
		return
	}
	if result, err = v.PathVal(path, optVars...); err != nil {
		err = bwerr.Refine(err,
			ansiMustPathValFailed,
			path, bwjson.Pretty(v), varsJSON(path, optVars),
		)
	}
	return
}

// MustValBy - must-обертка ValBy()
func MustValBy(v bw.Val, pathStr string, opt ...ValByOpt) (result interface{}) {
	var err error
	if result, err = ValBy(v, pathStr, opt...); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return
}

// SetValBy - обертка PathFrom()+bw.Val.SetPathVal()
func SetValBy(val interface{}, v bw.Val, pathStr string, opt ...ValByOpt) (err error) {
	var (
		optBases [][]bw.ValPath
		optVars  []map[string]interface{}
		path     bw.ValPath
	)
	if len(opt) > 0 {
		optBases = append(optBases, opt[0].Bases)
		optVars = append(optVars, opt[0].Vars)
	}
	if path, err = PathFrom(pathStr, optBases...); err != nil {
		return
	}
	if err = v.SetPathVal(val, path, optVars...); err != nil {
		err = bwerr.Refine(err,
			ansiMustSetPathValFailed,
			path, bwjson.Pretty(v), varsJSON(path, optVars),
		)
	}
	return
}

// MustSetPathVal - обертка PathFrom()+bw.Val.SetPathVal()
func MustSetValBy(val interface{}, v bw.Val, pathStr string, opt ...ValByOpt) {
	var err error
	if err = SetValBy(val, v, pathStr, opt...); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return
}

// ============================================================================

func From(s string, optVars ...map[string]interface{}) (result interface{}) {
	return FromTemplate(TemplateFrom(s), optVars...)
}

type Template struct {
	val interface{}
}

func TemplateFrom(s string) (result Template) {
	var err error
	var val interface{}
	if val, err = func(s string) (result interface{}, err error) {
		defer func() {
			if err != nil {
				result = nil
			}
		}()
		p := bwparse.From(bwrune.ProviderFromString(s))
		var ok bool
		if result, _, ok, err = p.Val(); err == nil && ok {
			err = p.SkipSpace(bwparse.TillEOF)
		} else if err == nil {
			err = p.Unexpected(p.Curr)
		}
		return
	}(s); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return Template{val: val}
}

func FromTemplate(template Template, optVars ...map[string]interface{}) (result interface{}) {
	var err error
	if result, err = expandPaths(template.val, template.val, true, optVars...); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return
}

func expandPaths(val interface{}, rootVal interface{}, isRoot bool, optVars ...map[string]interface{}) (result interface{}, err error) {
	var path bw.ValPath
	var ok bool
	if path, ok = val.(bw.ValPath); ok {
		var h Holder
		if isRoot {
			h = Holder{}
		} else {
			h = Holder{Val: rootVal}
		}
		result, err = h.PathVal(path, optVars...)
	} else {
		result = val
		switch _, kind := Kind(val); kind {
		case ValMap:
			m := result.(map[string]interface{})
			for key, val := range m {
				if val, err = expandPaths(val, rootVal, false, optVars...); err != nil {
					return
				}
				m[key] = val
			}
		case ValArray:
			vals := result.([]interface{})
			for i, val := range vals {
				if val, err = expandPaths(val, rootVal, false, optVars...); err != nil {
					return
				}
				vals[i] = val
			}
		}
	}
	return
}

// ============================================================================

// Bool - пытается извлечь bool из interface{}
func Bool(val interface{}) (result bool, ok bool) {
	if v, kind := Kind(val); kind == ValBool {
		result, ok = v.(bool)
	}
	return
}

// MustBool - must-обертка Bool()
func MustBool(val interface{}) (result bool) {
	var ok bool
	if result, ok = Bool(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Bool")
	}
	return
}

// Int - пытается извлечь int из interface{}
func Int(val interface{}) (result int, ok bool) {
	switch v, kind := Kind(val); kind {
	case ValInt:
		result, _ = v.(int)
		ok = true
	case ValNumber:
		n, _ := v.(bwtype.Number)
		result, ok = n.Int()
	}
	return
}

// MustInt - must-обертка Int()
func MustInt(val interface{}) (result int) {
	var ok bool
	if result, ok = Int(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Int")
	}
	return
}

// Float64 - пытается извлечь float64 из interface{}
func Float64(val interface{}) (result float64, ok bool) {
	switch v, kind := Kind(val); kind {
	case ValInt:
		var i int
		i, _ = v.(int)
		result = float64(i)
		ok = true
	case ValFloat64:
		result, _ = v.(float64)
		ok = true
	case ValNumber:
		n, _ := v.(bwtype.Number)
		result, ok = n.Float64()
	}
	return
}

// MustFloat64 - must-обертка Float64()
func MustFloat64(val interface{}) (result float64) {
	var ok bool
	if result, ok = Float64(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Float64")
	}
	return
}

// bwtype.Number - пытается извлечь bwtype.Number из interface{}
func Number(val interface{}) (result bwtype.Number, ok bool) {
	result, ok = bwtype.NumberFrom(val)
	// var (
	//   kind ValKind
	//   t    interface{}
	// )
	// switch t, kind = Kind(val); kind {
	// case ValInt:
	//   i, _ := t.(int)
	//   result = NumberFromInt(i)
	//   ok = true
	// case ValFloat64:
	//   f, _ := t.(float64)
	//   result = NumberFromFloat64(f)
	//   ok = true
	// case ValNumber:
	//   result, _ = t.(bwtype.Number)
	//   ok = true
	// }
	// bwdebug.Print("kind", kind, "t", t, "val:#v", val)
	return
}

// MustNumber - must-обертка bwtype.Number()
func MustNumber(val interface{}) (result bwtype.Number) {
	var ok bool
	if result, ok = Number(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Number")
	}
	return
}

// String - пытается извлечь string из interface{}
func String(val interface{}) (result string, ok bool) {
	if v, kind := Kind(val); kind == ValString {
		result, ok = v.(string)
	}
	return
}

// MustString - must-обертка String()
func MustString(val interface{}) (result string) {
	var ok bool
	if result, ok = String(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "String")
	}
	return
}

// Map - пытается извлечь map[string]interface{} из interface{}
func Map(val interface{}) (result map[string]interface{}, ok bool) {
	if v, kind := Kind(val); kind == ValMap {
		result, ok = v.(map[string]interface{})
	}
	return
}

// MustMap - must-обертка Map()
func MustMap(val interface{}) (result map[string]interface{}) {
	var ok bool
	if result, ok = Map(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Map")
	}
	return
}

// Map - пытается извлечь []interface{} из interface{}
func Array(val interface{}) (result []interface{}, ok bool) {
	if v, kind := Kind(val); kind == ValArray {
		result, ok = v.([]interface{})
	}
	return
}

// MustArray - must-обертка Array()
func MustArray(val interface{}) (result []interface{}) {
	var ok bool
	if result, ok = Array(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Array")
	}
	return result
}

// Map - пытается извлечь []string из interface{}
func ArrayOfString(val interface{}) (result []string, ok bool) {
	switch v, kind := Kind(val); kind {
	case ValArrayOfString:
		result, ok = v.([]string)
	case ValArray:
		vals, _ := v.([]interface{})
		result = []string{}
		var s string
		for _, val := range vals {
			if s, ok = val.(string); !ok {
				return
			}
			result = append(result, s)
		}
		ok = true
	}
	return
}

// MustArray - must-обертка Array()
func MustArrayOfString(val interface{}) (result []string) {
	var ok bool
	if result, ok = ArrayOfString(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "ArrayOfString")
	}
	return result
}

// ============================================================================

// ============================================================================
