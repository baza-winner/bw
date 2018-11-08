package bw

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

//go:generate stringer -type=ValPathItemType

// ============================================================================

var Spew spew.ConfigState

func init() {
	Spew = spew.ConfigState{SortKeys: true}
}

// ============================================================================

func Args(args ...interface{}) []interface{} {
	return args
}

// ============================================================================

type I interface {
	FmtString() string
	FmtArgs() []interface{}
}

// ============================================================================

type A struct {
	Fmt  string
	Args []interface{}
}

func (v A) FmtString() string {
	return v.Fmt
}

func (v A) FmtArgs() []interface{} {
	return v.Args
}

// ============================================================================

func Fmt(fmtString string, fmtArg ...interface{}) A {
	return A{fmtString, fmtArg}
}

// ============================================================================

func PluralWord(count int, word string, word1 string, word2_4 string, _word5more ...string) (result string) {
	var word5more string
	if _word5more != nil {
		word5more = _word5more[0]
	}
	if len(word5more) == 0 {
		word5more = word2_4
	}
	result = word5more
	decimal := count / 10 % 10
	if decimal != 1 {
		unit := count % 10
		if unit == 1 {
			result = word1
		} else if 2 <= unit && unit <= 4 {
			result = word2_4
		}
	}
	return word + result
}

// ============================================================================

// https://stackoverflow.com/questions/6878590/the-maximum-value-for-an-int-type-in-go

const (
	MaxUint8 = ^uint8(0)
	MinUint8 = 0
	MaxInt8  = int8(MaxUint8 >> 1)
	MinInt8  = -MaxInt8 - 1

	MaxUint16 = ^uint16(0)
	MinUint16 = 0
	MaxInt16  = int16(MaxUint16 >> 1)
	MinInt16  = -MaxInt16 - 1

	MaxUint32 = ^uint32(0)
	MinUint32 = 0
	MaxInt32  = int32(MaxUint32 >> 1)
	MinInt32  = -MaxInt32 - 1

	MaxUint64 = ^uint64(0)
	MinUint64 = 0
	MaxInt64  = int64(MaxUint64 >> 1)
	MinInt64  = -MaxInt64 - 1

	MaxUint = ^uint(0)
	MinUint = 0
	MaxInt  = int(MaxUint >> 1)
	MinInt  = -MaxInt - 1
)

// ============================================================================

type ValPathItemType uint8

const (
	ValPathItemHash ValPathItemType = iota
	ValPathItemIdx
	ValPathItemKey
	ValPathItemPath
	ValPathItemVar
)

func (v ValPathItemType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

type ValPathItem struct {
	Type ValPathItemType
	Idx  int
	Key  string
	Path ValPath
}

type Path []ValPathItem

func (v Path) MarshalJSON() ([]byte, error) {
	result := []interface{}{}
	for _, i := range v {
		result = append(result, i)
	}
	return json.Marshal(result)
}

type ValPath []ValPathItem

func (v ValPath) String() (result string) {
	ss := []string{}
	if len(v) == 0 {
		result = "."
	} else {
		for _, vpi := range v {
			switch vpi.Type {
			case ValPathItemPath:
				ss = append(ss, "{"+vpi.Path.String()+"}")
			case ValPathItemKey:
				ss = append(ss, vpi.Key)
			case ValPathItemVar:
				ss = append(ss, "$"+vpi.Key)
			case ValPathItemIdx:
				ss = append(ss, strconv.FormatInt(int64(vpi.Idx), 10))
			case ValPathItemHash:
				ss = append(ss, "#")
			default:
				panic(Spew.Sprintf("%#v", vpi.Type))
			}
		}
		result = strings.Join(ss, ".")
	}
	return
}

// ============================================================================

type Val interface {
	PathVal(path ValPath, optVars ...map[string]interface{}) (result interface{}, err error)
	SetPathVal(val interface{}, path ValPath, optVars ...map[string]interface{}) (err error)
	MarshalJSON() ([]byte, error)
}

// ============================================================================
