package deftype

import (
	"github.com/baza-winner/bwcore/bwerror"
)

type Item uint16

const (
	Item_below_ Item = iota
	Bool
	String
	Int
	Number
	Map
	Array
	ArrayOf
	Item_above_
)

//go:generate stringer -type=Item

var mapItemFromString = map[string]Item{}

func init() {
	for i := Item_below_ + 1; i < Item_above_; i++ {
		mapItemFromString[i.String()] = i
	}
	return
}

func ItemFromString(s string) (result Item, err error) {
	var ok bool
	if result, ok = mapItemFromString[s]; !ok {
		err = bwerror.Error("<ansiCmd>ItemFromString<ansi>: uknown <ansiPrimaryLiteral>%s", result)
	}
	return
}

// ============================================================================

//go:generate bwsetter -type=Item -set=Set -omitprefix
