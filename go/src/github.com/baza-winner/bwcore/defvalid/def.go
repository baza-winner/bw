package defvalid

import (
	// "github.com/baza-winner/bwcore/bwerror"
	// "github.com/baza-winner/bwcore/bwjson"
	// "github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
	// "github.com/baza-winner/bwcore/defparse"
	// "log"
	// "reflect"
)

type Def struct {
	tp        deftype
	enum      *bwset.Strings
	minInt    *int
	maxInt    *int
	minNumber *float64
	maxNumber *float64
	keys      *map[string]*Def
	elem      *Def
	arrayElem *Def
	dflt      interface{}
}

func (v *Def) GetDataForJson() interface{} {
  if v == nil {return nil }
  result := map[string]interface{}{}
	result["tp"] = v.tp.GetDataForJson()
	if v.enum != nil {
		result["enum"] = v.enum.GetDataForJson()
	}
	if v.minInt != nil {
		result["minInt"] = v.minInt
	}
	if v.maxInt != nil {
		result["maxInt"] = v.maxInt
	}
	if v.minNumber != nil {
		result["minNumber"] = v.minNumber
	}
	if v.maxNumber != nil {
		result["maxNumber"] = v.maxNumber
	}
	if v.keys != nil {
		keysJsonData := map[string]interface{}{}
		for k, v := range *(v.keys) {
			keysJsonData[k] = (*v).GetDataForJson()
		}
		result["keys"] = keysJsonData
	}
	if v.elem != nil {
		result["elem"] = (*(v.elem)).GetDataForJson()
	}
	if v.arrayElem != nil {
		result["elem"] = (*(v.arrayElem)).GetDataForJson()
	}
	return result
}
