package defvalid

import (
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/defvalid/deftype"

	// "github.com/baza-winner/bwcore/bwjson"
	// "github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
	// "github.com/baza-winner/bwcore/defparse"
	// "reflect"
)

type Def struct {
	tp         deftype.Set
	isOptional bool
	enum       bwset.String
	minInt     *int64
	maxInt     *int64
	minNumber  *float64
	maxNumber  *float64
	keys       map[string]Def
	elem       *Def
	arrayElem  *Def
	dflt       interface{}
}

// func (v *Def) Copy() (result *Def) {
// 	return &Def{
// 		tp: v.tp.Copy(),
// 		isOptional:v.isOptional,
// 		enum:v.enum.Copy(),
// 		minInt: copyPtrToInt64(v.minInt),
// 		maxInt: copyPtrToInt64(v.maxInt,
// 		minNumber: copyPtrToFloat64(v.minNumber),
// 		maxNumber: copyPtrToFloat64(v.maxNumber),
// 		elem:
// 	}
// }

// func copyPtrToInt64(p *int64) *int64 {
// 	if p == nil {
// 		return nil
// 	} else {
// 		i := *p
// 		return &i
// 	}
// }

// func copyPtrToFloat64(p *float64) *float64 {
// 	if p == nil {
// 		return nil
// 	} else {
// 		i := *p
// 		return &i
// 	}
// }

func MustDef(v interface{}) (result *Def) {
	if v == nil {
		return nil
	}
	var ok bool
	if result, ok = v.(*Def); !ok {
		bwerror.Panic("%#v is not *Def", v)
	}
	return
}

func (v *Def) DataForJSON() interface{} {
	if v == nil {
		return nil
	}
	result := map[string]interface{}{}
	result["tp"] = v.tp.DataForJSON()
	result["isOptional"] = v.isOptional
	if v.enum != nil {
		result["enum"] = v.enum.DataForJSON()
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
		for k, v := range v.keys {
			keysJsonData[k] = v.DataForJSON()
		}
		result["keys"] = keysJsonData
	}
	if v.elem != nil {
		result["elem"] = (*(v.elem)).DataForJSON()
	}
	if v.arrayElem != nil {
		result["arrayElem"] = (*(v.arrayElem)).DataForJSON()
	}
	if v.dflt != nil {
		result["dflt"] = v.dflt
	}
	return result
}
