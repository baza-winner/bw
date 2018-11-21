package bwtype

import (
	"encoding/json"
	"fmt"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
)

// ============================================================================
type RangeLimit struct {
	val interface{}
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

func RangeLimitFrom(val interface{}) (result RangeLimit, ok bool) {
	var (
		i int
		f float64
	)
	if val == nil {
		result = RangeLimit{}
		ok = true
	} else if i, ok = Int(val); ok {
		result = RangeLimit{val: i}
	} else if f, ok = Float64(val); ok {
		result = RangeLimit{val: f}
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

func (n RangeLimit) Int() (result int, ok bool) {
	result, ok = n.val.(int)
	return
}

func (n RangeLimit) Float64() (result float64, ok bool) {
	switch t := n.val.(type) {
	case int:
		ok = true
		result = float64(t)
	case float64:
		ok = true
		result = t
	}
	return
}

func (n RangeLimit) MustInt() (result int) {
	var ok bool
	if result, ok = n.Int(); !ok {
		bwerr.Panic(ansiIsNotOfType, n.val, "Int")
	}
	return
}

func (n RangeLimit) MustFloat64() (result float64) {
	var ok bool
	if result, ok = n.Float64(); !ok {
		bwerr.Panic(ansiIsNotOfType, n.val, "Float64")
	}
	return
}

func (n RangeLimit) IsNaN() bool {
	return n.val == nil
}

func (n RangeLimit) IsInt() (result bool) {
	_, result = n.val.(int)
	return
}

func (n RangeLimit) IsEqualTo(a RangeLimit) (result bool) {
	return n.compareTo(a, func(isInt bool, i, j int, f, g float64) bool {
		if isInt {
			return i == j
		} else {
			return f == g
		}
	})
}

func (n RangeLimit) IsLessThan(a RangeLimit) bool {
	return n.compareTo(a, func(isInt bool, i, j int, f, g float64) bool {
		if isInt {
			return i < j
		} else {
			return f < g
		}
	})
}

type compareFunc func(isInt bool, i, j int, f, g float64) (result bool)

func (n RangeLimit) compareTo(a RangeLimit, f compareFunc) (result bool) {
	if i, ok := n.Int(); ok {
		if j, ok := a.Int(); ok {
			result = f(true, i, j, 0, 0)
		} else {
			result = f(false, 0, 0, float64(i), a.MustFloat64())
		}
	} else if g, ok := n.Float64(); ok {
		result = f(false, 0, 0, g, a.MustFloat64())
	}
	return
}

func (n RangeLimit) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.val)
}

// ============================================================================

//go:generate stringer -type=RangeKindValue

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
		if max.IsLessThan(min) {
			err = bwerr.From(ansiMaxMustNotBeLessThanMin, bwjson.Pretty(a.Max), bwjson.Pretty(a.Min))
		}
	}
	return
}

func MustRangeFrom(a A) (result Range) {
	var err error
	if result, err = RangeFrom(a); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return
}

func (r Range) Kind() (result RangeKindValue) {
	if !r.min.IsNaN() {
		if !r.max.IsNaN() {
			result = RangeMinMax
		} else {
			result = RangeMin
		}
	} else if !r.max.IsNaN() {
		result = RangeMax
	}
	return
}

func (v Range) String() (result string) {
	switch v.Kind() {
	case RangeMinMax:
		result = fmt.Sprintf("%s..%s", bwjson.Pretty(v.min), bwjson.Pretty(v.max))
	case RangeMin:
		result = fmt.Sprintf("%s..", bwjson.Pretty(v.min))
	case RangeMax:
		result = fmt.Sprintf("..%s", bwjson.Pretty(v.max))
	default:
		result = ".."
	}
	return
}

func (r Range) Contains(val interface{}) (result bool) {
	n := MustRangeLimitFrom(val)
	if n.IsNaN() {
		return false
	}
	var minResult, maxResult bool
	rangeKind := r.Kind()
	switch rangeKind {
	case RangeMin, RangeMinMax:
		minResult = !n.IsLessThan(r.min)
	}
	switch rangeKind {
	case RangeMax, RangeMinMax:
		maxResult = !r.max.IsLessThan(n)
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
