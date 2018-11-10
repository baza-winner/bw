package deftype

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerr"
)

type Item uint16

const (
	// ItemBelow Item = iota
	Unknown Item = iota
	Bool
	String
	Int
	Number
	Map
	Array
	ArrayOf
	// ItemAbove
)

const (
	_SetTestItemA Item = iota + 1
	_SetTestItemB
)

// func init() {
// 	bwdebug.Print("json", bwjson.Pretty(map[string]interface{}{"set": From(Bool, String)}))

// }

//go:generate stringer -type=Item
//go:generate bwsetter -type=Item -set=Set -omitprefix -test

var mapItemFromString = map[string]Item{}

func init() {
	// for i := ItemBelow + 1; i < ItemAbove; i++ {
	for i := Unknown; i <= ArrayOf; i++ {
		mapItemFromString[i.String()] = i
	}
	return
}

func (v Item) MarshalJSON() ([]byte, error) {
	// return []byte(v.String()), nil
	return json.Marshal(v.String())
}

var ansiUknown string

func init() {
	ansiUknown = ansi.String("<ansiPath>ItemFromString<ansi>: uknown <ansiVal>%s")
}

func ItemFromString(s string) (result Item, err error) {
	var ok bool
	if result, ok = mapItemFromString[s]; !ok {
		err = bwerr.From(ansiUknown, result)
	}
	return
}

// ============================================================================
