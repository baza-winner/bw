package bwtype

import (
	"encoding/json"
	"fmt"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
)

// ============================================================================
type Number struct {
	val interface{}
}

func Int(val interface{}) (result int, ok bool) {
	ok = true
	switch t := val.(type) {
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
		} else {
			ok = false
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
		} else {
			ok = false
		}
	default:
		ok = false
	}
	return
}

func Float64(val interface{}) (result float64, ok bool) {

}

func NumberFrom(val interface{}) (result Number, err error) {
	switch t, kind := val.(type); kind {
	case ValInt:
		i, _ := t.(int)
		result = NumberFromInt(i)
	case ValFloat64:
		f, _ := t.(float64)
		result = NumberFromFloat64(f)
	case ValNumber:
		n, _ := t.(Number)
		result = n
	default:
		err = bwerr.From("<ansiVal>%#v<ansi> can not be a <ansiType>Number")
	}
	return
}

func MustNumberFrom(val interface{}) (result Number) {
	var err error
	if result, err = NumberFrom(val); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return
}

func NumberFromInt(i int) Number {
	return Number{val: i}
}

func NumberFromFloat64(f float64) Number {
	return Number{val: f}
}

func (n Number) Int() (result int, ok bool) {
	result, ok = n.val.(int)
	return
}

func (n Number) Float64() (result float64, ok bool) {
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

func (n Number) MustInt() (result int) {
	var ok bool
	if result, ok = n.Int(); !ok {
		bwerr.Panic(ansiIsNotOfType, n.val, "Int")
	}
	return
}

func (n Number) MustFloat64() (result float64) {
	var ok bool
	if result, ok = n.Float64(); !ok {
		bwerr.Panic(ansiIsNotOfType, n.val, "Float64")
	}
	return
}

func (n Number) IsNaN() bool {
	return n.val == nil
}

func (n Number) IsInt() (result bool) {
	_, result = n.val.(int)
	return
}

func (n Number) IsEqualTo(a Number) (result bool) {
	return n.compareTo(a, func(isInt bool, i, j int, f, g float64) bool {
		if isInt {
			return i == j
		} else {
			return f == g
		}
	})
}

func (n Number) IsLessThan(a Number) bool {
	return n.compareTo(a, func(isInt bool, i, j int, f, g float64) bool {
		if isInt {
			return i < j
		} else {
			return f < g
		}
	})
}

type numberCompare func(isInt bool, i, j int, f, g float64) (result bool)

func (n Number) compareTo(a Number, f numberCompare) (result bool) {
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

func (n Number) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.val)
}

// var ansiIsNotOfType string

// func init() {
//  ansiIsNotOfType = ansi.String("<ansiVal>%#v<ansi> is not <ansiType>%s")
// }

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
	min, max Number
}

func (r Range) Min() Number {
	return r.min
}

func (r Range) Max() Number {
	return r.max
}

func RangeFrom(min, max Number) (result Range, err error) {
	result = Range{min: min, max: max}
	if result.Kind() == RangeMinMax {
		var isValidRange bool
		if max.IsLessThan(min) {
			err = bwerr.From(
				"<ansiVar>max<ansi> (<ansiVal>%s<ansi>) must not be <ansiErr>less<ansi> then <ansiVar>min<ansi> (<ansiVal>%s<ansi>)",
				bwjson.Pretty(max), bwjson.Pretty(min),
			)
		}
	}
	return
}

func MustRangeFrom(min, max Number) (result Range) {
	var err error
	if result, err = RangeFrom(min, max); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return
}

func (r Range) Kind() (result RangeKindValue) {
	if !r.Min.IsNaN() {
		if !r.Max.IsNaN() {
			result = RangeMinMax
		} else {
			result = RangeMin
		}
	} else if !r.Max.IsNaN() {
		result = RangeMax
	}
	return
}

func (v Range) String() (result string) {
	switch v.Kind() {
	case RangeMinMax:
		result = fmt.Sprintf("%s..%s", bwjson.Pretty(v.Min), bwjson.Pretty(v.Max))
	case RangeMin:
		result = fmt.Sprintf("%s..", bwjson.Pretty(v.Min))
	case RangeMax:
		result = fmt.Sprintf("..%s", bwjson.Pretty(v.Max))
	default:
		result = ".."
	}
	return
}

func (r Range) Contains(val interface{}) (result bool) {
	var (
		valKind ValKind
		i       int
		f       float64
		t       interface{}
		n       Number
	)
	switch t, valKind = Kind(val); valKind {
	case ValInt:
		i, _ = t.(int)
		n = NumberFromInt(i)
	case ValFloat64:
		f, _ = t.(float64)
		n = NumberFromFloat64(f)
	case ValNumber:
		n, _ = t.(Number)
	default:
		return false
	}
	var (
		minResult, maxResult bool
	)
	rangeKind := r.Kind()
	switch rangeKind {
	case RangeMin, RangeMinMax:
		minResult = !n.IsLessThan(r.Min)
	}
	switch rangeKind {
	case RangeMax, RangeMinMax:
		maxResult = !r.Max.IsLessThan(n)
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
	ansiIsNotOfType string
)

func init() {
	ansiIsNotOfType = ansi.String("<ansiVal>%#v<ansi> is not <ansiType>%s")
}

// ============================================================================
