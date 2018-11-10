// Package bwval реализует интерфейc bw.Val и утилиты для работы с этим интерфейсом.
package bwval

import (
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwval/path"
	"github.com/baza-winner/bwcore/bwval/val"
)

// ============================================================================

//go:generate stringer -type=ValKind,RangeKindValue

// ============================================================================

// PathFrom - конструктор-парсер bw.ValPath из строки
func PathFrom(s string) bw.ValPath {
	return path.MustParse(s)
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
		// bwerr.PanicA(bwerr.Err(err))
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

func From(s string, optVars ...map[string]interface{}) interface{} {
	return val.MustParse(s, optVars...)
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
	if v, kind := Kind(val); kind == ValInt {
		result, ok = v.(int)
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

// Number - пытается извлечь float64 из interface{}
func Number(val interface{}) (result float64, ok bool) {
	switch v, kind := Kind(val); kind {
	case ValInt:
		var i int
		i, ok = v.(int)
		result = float64(i)
	case ValNumber:
		result, ok = v.(float64)
	}
	return
}

// MustNumber - must-обертка Number()
func MustNumber(val interface{}) (result float64) {
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

// ============================================================================
