package bwtype

import (
	"encoding/json"
	"fmt"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwdebug"
	"github.com/baza-winner/bwcore/bwerr"
)

// ============================================================================

func Int(val interface{}) (result int, ok bool) {
	ok = true
	switch t := val.(type) {
	case int8:
		result = int(t)
	case int16:
		result = int(t)
	case int32:
		result = int(t)
	case int:
		result = t
	case uint8:
		result = int(t)
	case uint16:
		result = int(t)
	case uint64:
		if ok = t <= uint64(bw.MaxInt); ok {
			result = int(t)
		}
	case uint:
		if ok = t <= uint(bw.MaxInt); ok {
			result = int(t)
		}
	case float32:
		result = int(t)
		if ok = t == float32(result); !ok {
			result = 0
		}
	case float64:
		result = int(t)
		if ok = t == float64(result); !ok {
			result = 0
		}
	default:
		result, ok = platformSpecificInt(val)
	}
	return
}

func Uint(val interface{}) (result uint, ok bool) {
	ok = true
	switch t := val.(type) {
	case int8:
		if ok = t >= 0; ok {
			result = uint(t)
		}
	case int16:
		if ok = t >= 0; ok {
			result = uint(t)
		}
	case int32:
		if ok = t >= 0; ok {
			result = uint(t)
		}
	case int:
		if ok = t >= 0; ok {
			result = uint(t)
		}
	case uint8:
		result = uint(t)
	case uint16:
		result = uint(t)
	case uint32:
		result = uint(t)
	case uint:
		result = t
	case float32:
		result = uint(t)
		if ok = t == float32(result); !ok {
			result = 0
		}
	case float64:
		result = uint(t)
		if ok = t == float64(result); !ok {
			result = 0
		}
	default:
		result, ok = platformSpecificUint(val)
	}
	return
}

func Float64(val interface{}) (result float64, ok bool) {
	ok = true
	switch t := val.(type) {
	case int8:
		result = float64(t)
	case int16:
		result = float64(t)
	case int32:
		result = float64(t)
	case int64:
		result = float64(t)
	case int:
		result = float64(t)
	case uint8:
		result = float64(t)
	case uint16:
		result = float64(t)
	case uint32:
		result = float64(t)
	case uint64:
		result = float64(t)
	case uint:
		result = float64(t)
	case float32:
		result = float64(t)
	case float64:
		result = t
	default:
		ok = false
	}
	return
}

// ============================================================================

func MustInt(val interface{}) (result int) {
	var ok bool
	if result, ok = Int(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Int")
	}
	return
}

func MustUint(val interface{}) (result uint) {
	var ok bool
	if result, ok = Uint(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Uint")
	}
	return
}

func MustFloat64(val interface{}) (result float64) {
	var ok bool
	if result, ok = Float64(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Float64")
	}
	return
}

// ============================================================================

type Number struct {
	val interface{}
}

func NumberFrom(val interface{}) (result Number, ok bool) {
	var (
		i  int
		u  uint
		f  float64
		rl RangeLimit
	)
	if i, ok = Int(val); ok {
		result = Number{i}
	} else if u, ok = Uint(val); ok {
		result = Number{u}
	} else if f, ok = Float64(val); ok {
		result = Number{f}
	} else if rl, ok = val.(RangeLimit); ok {
		result, ok = NumberFrom(rl.val)
	} else {
		result, ok = val.(Number)
	}
	return
}

func MustNumberFrom(val interface{}) (result Number) {
	var ok bool
	if result, ok = NumberFrom(val); !ok {
		bwerr.Panic(ansi.String("<ansiVal>%#v<ansi> can not be a <ansiType>Number"), val)
	}
	return
}

func (n Number) Val() interface{} {
	return n.val
}

func (n Number) IsEqualTo(a Number) (result bool) {
	return n.compareTo(a, func(kind compareKind, u, v uint, i, j int, f, g float64) (result bool) {
		switch kind {
		case compareUintUint:
			result = u == v
		case compareIntInt:
			result = i == j
		case compareFloat64Float64:
			result = f == g
		}
		return
	})
}

func (n Number) IsLessThan(a Number) (result bool) {
	result = n.compareTo(a, func(kind compareKind, u, v uint, i, j int, f, g float64) (result bool) {
		switch kind {
		case compareUintUint:
			result = u < v
		case compareIntUint:
			result = true
		case compareIntInt:
			result = i < j
		case compareFloat64Float64:
			result = f < g
		}
		return
	})
	return
}

type compareKind uint8

const (
	compareUintUint compareKind = iota
	compareUintInt
	compareIntUint
	compareIntInt
	compareFloat64Float64
)

type compareFunc func(kind compareKind, u, v uint, i, j int, f, g float64) (result bool)

func (n Number) compareTo(a Number, fn compareFunc) (result bool) {
	if u, ok := Uint(n.val); ok {
		if v, ok := Uint(a.val); ok {
			result = fn(compareUintUint, u, v, 0, 0, 0, 0)
			bwdebug.Print("u", u, "v", v, "result", result)
		} else if j, ok := Int(a.val); ok {
			result = fn(compareUintInt, u, 0, 0, j, 0, 0)
		} else if g, ok := Float64(a.val); ok {
			result = fn(compareFloat64Float64, 0, 0, 0, 0, float64(u), g)
		}
	} else if i, ok := Int(n.val); ok {
		if v, ok := Uint(a.val); ok {
			result = fn(compareIntUint, 0, v, i, 0, 0, 0)
		} else if j, ok := Int(a.val); ok {
			result = fn(compareIntInt, 0, 0, i, j, 0, 0)
		} else if g, ok := Float64(a.val); ok {
			result = fn(compareFloat64Float64, 0, 0, 0, 0, float64(i), g)
		}
	} else if f, ok := Float64(n.val); ok {
		if g, ok := Float64(a.val); ok {
			result = fn(compareFloat64Float64, 0, 0, 0, 0, f, g)
		}
	}
	return
}

func (n Number) String() string {
	bytes, _ := json.Marshal(n.val)
	return string(bytes)
}

func (n Number) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.val)
}

// ============================================================================

type RangeLimit struct {
	val interface{}
}

func RangeLimitFrom(val interface{}) (result RangeLimit, ok bool) {
	var (
		path bw.ValPath
		n    Number
	)
	if val == nil {
		result = RangeLimit{}
		ok = true
	} else if n, ok = NumberFrom(val); ok {
		result = RangeLimit{val: n}
	} else if path, ok = val.(bw.ValPath); ok {
		result = RangeLimit{val: path}
	} else {
		result, ok = val.(RangeLimit)
	}
	return
}

func MustRangeLimitFrom(val interface{}) (result RangeLimit) {
	var ok bool
	if result, ok = RangeLimitFrom(val); !ok {
		bwerr.Panic(ansiValCanNotBeRangeLimit, val)
	}
	return
}

func (n RangeLimit) Nil() (ok bool) {
	ok = n.val == nil
	return
}

func (n RangeLimit) Number() (result Number, ok bool) {
	result, ok = n.val.(Number)
	return
}

func (n RangeLimit) Path() (result bw.ValPath, ok bool) {
	result, ok = n.val.(bw.ValPath)
	return
}

func (n RangeLimit) MustNumber() (result Number) {
	var ok bool
	if result, ok = n.Number(); !ok {
		bwerr.Panic(ansiIsNotOfType, n.val, "Number")
	}
	return
}

func (n RangeLimit) MustPath() (result bw.ValPath) {
	var ok bool
	if result, ok = n.Path(); !ok {
		bwerr.Panic(ansiIsNotOfType, n.val, "bw.ValPath")
	}
	return
}

func (rl RangeLimit) String() (result string) {
	var (
		n    Number
		path bw.ValPath
		ok   bool
	)
	if n, ok = rl.Number(); ok {
		result = n.String()
	} else if path, ok = rl.Path(); ok {
		result = path.String()
		if !(len(path) > 0 && path[0].Type == bw.ValPathItemVar) {
			result = "{{" + result + "}}"
		}
	}
	return
}

// ============================================================================

type RangeKindValue uint8

const (
	RangeNo RangeKindValue = iota
	RangeMin
	RangeMax
	RangeMinMax
)

type Range struct {
	min, max RangeLimit
}

func (r Range) Min() RangeLimit {
	return r.min
}

func (r Range) Max() RangeLimit {
	return r.max
}

type A struct {
	Min, Max interface{}
}

func RangeFrom(a A) (result Range, err error) {
	var min, max RangeLimit
	var ok bool
	if min, ok = RangeLimitFrom(a.Min); !ok {
		err = bwerr.From(ansiVarValCanNotBeRangeLimit, "a.Min", a.Min)
		return
	}
	if max, ok = RangeLimitFrom(a.Max); !ok {
		err = bwerr.From(ansiVarValCanNotBeRangeLimit, "a.Max", a.Max)
		return
	}
	result = Range{min: min, max: max}
	if result.Kind() == RangeMinMax {
		if min, ok := NumberFrom(a.Min); ok {
			if max, ok := NumberFrom(a.Max); ok {
				if max.IsLessThan(min) {
					err = bwerr.From(ansiMaxMustNotBeLessThanMin, max, min)
				}
			}
		}
	}
	return
}

func MustRangeFrom(a A) (result Range) {
	var err error
	if result, err = RangeFrom(a); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (r Range) Kind() (result RangeKindValue) {
	if r.min.val != nil {
		if r.max.val != nil {
			result = RangeMinMax
		} else {
			result = RangeMin
		}
	} else if r.max.val != nil {
		result = RangeMax
	}
	return
}

func (v Range) String() (result string) {
	result = fmt.Sprintf("%s..%s", v.min, v.max)
	return
}

func (r Range) Contains(val interface{}) (result bool) {
	var n Number
	var ok bool
	if n, ok = NumberFrom(val); !ok {
		return false
	}
	var minResult, maxResult bool
	rangeKind := r.Kind()
	switch rangeKind {
	case RangeMin, RangeMinMax:
		if min, ok := NumberFrom(r.min.val); ok {
			if n.IsEqualTo(min) {
				return true
			} else {
				minResult = !n.IsLessThan(min)
			}
		}
	}
	switch rangeKind {
	case RangeMax, RangeMinMax:
		if max, ok := NumberFrom(r.max.val); ok {
			if n.IsEqualTo(max) {
				return true
			} else {
				maxResult = !max.IsLessThan(n)
			}
		}
	}
	switch rangeKind {
	case RangeMinMax:
		result = minResult && maxResult
	case RangeMax:
		result = maxResult
	case RangeMin:
		result = minResult
	default:
		result = true
	}
	return
}

func (r Range) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// ============================================================================

var (
	ansiVarValCanNotBeRangeLimit string
	ansiValCanNotBeRangeLimit    string
	ansiIsNotOfType              string
	ansiMaxMustNotBeLessThanMin  string
)

func init() {
	ansiVarValCanNotBeRangeLimit = ansi.String("<ansiVar>%s<ansi> (<ansiVal>%#v<ansi>) can not be a <ansiType>RangeLimit")
	ansiValCanNotBeRangeLimit = ansi.String("<ansiVal>%#v<ansi> can not be a <ansiType>RangeLimit")
	ansiIsNotOfType = ansi.String("<ansiVal>%#v<ansi> is not <ansiType>%s")
	ansiMaxMustNotBeLessThanMin = ansi.String("<ansiVar>a.Max<ansi> (<ansiVal>%s<ansi>) must not be <ansiErr>less<ansi> then <ansiVar>a.Min<ansi> (<ansiVal>%s<ansi>)")
}

// ============================================================================
