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
