package deftype

import (
	"github.com/baza-winner/bwcore/bwerr"
)

type Item uint16

const (
	ItemBelow Item = iota
	Bool
	String
	Int
	Number
	Map
	Array
	ArrayOf
	ItemAbove
)

//go:generate stringer -type=Item

var mapItemFromString = map[string]Item{}

func init() {
	for i := ItemBelow + 1; i < ItemAbove; i++ {
		mapItemFromString[i.String()] = i
	}
	return
}

func (v Item) DataForJSON() interface{} {
	return v.String()
}

func ItemFromString(s string) (result Item, err error) {
	var ok bool
	if result, ok = mapItemFromString[s]; !ok {
		err = bwerr.From("<ansiPath>ItemFromString<ansi>: uknown <ansiVal>%s", result)
	}
	return
}

// ============================================================================

//go:generate bwsetter -type=Item -set=Set -omitprefix
