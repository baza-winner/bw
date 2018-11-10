package bwval

import (
	"encoding/json"
	"fmt"

	"github.com/baza-winner/bwcore/bwjson"
)

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
		default:
			result = ".."
		}
	}
	return
}

func RangeContains(v Range, val interface{}) (result bool) {
	rangeKind := RangeKind(v)
	if rangeKind == RangeNo {
		result = true
	} else if _, vk := Kind(val); vk == ValNumber || v.ValKind() == ValNumber {
		if n, ok := Number(val); ok {
			switch rangeKind {
			case RangeMinMax:
				result = MustNumber(v.Min()) <= n && n <= MustNumber(v.Max())
			case RangeMin:
				result = MustNumber(v.Min()) <= n
			case RangeMax:
				result = n <= MustNumber(v.Max())
			}
		}
	} else {
		if n, ok := Int(val); ok {
			switch rangeKind {
			case RangeMinMax:
				result = MustInt(v.Min()) <= n && n <= MustInt(v.Max())
			case RangeMin:
				result = MustInt(v.Min()) <= n
			case RangeMax:
				result = n <= MustInt(v.Max())
			}
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

// ============================================================================

func PtrToInt(i int) *int {
	return &i
}

func PtrToNumber(i float64) *float64 {
	return &i
}

// ============================================================================
