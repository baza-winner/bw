package defvalid

import (
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/defvalid/deftype"
	// "github.com/baza-winner/bwcore/bwjson"
	// "github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
	// "github.com/baza-winner/bwcore/defparse"
	// "log"
	// "reflect"
)

type Def struct {
	tp         deftype.Set
	isOptional bool
	isSimple   bool
	enum       bwset.Strings
	minInt     *int64
	maxInt     *int64
	minNumber  *float64
	maxNumber  *float64
	keys       map[string]Def
	elem       *Def
	arrayElem  *Def
	dflt       interface{}
}

func mustDef(v interface{}) (result *Def) {
	if v == nil {
		return nil
	}
	var ok bool
	if result, ok = v.(*Def); !ok {
		bwerror.Panic("%#v is not *Def", v)
	}
	return
}

func (v *Def) GetDataForJson() interface{} {
	if v == nil {
		return nil
	}
	result := map[string]interface{}{}
	result["tp"] = v.tp.GetDataForJson()
	result["isSimple"] = v.isSimple
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
		for k, v := range v.keys {
			keysJsonData[k] = v.GetDataForJson()
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
