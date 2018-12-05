// Package bwval реализует интерфейc bw.Val и утилиты для работы с этим интерфейсом.
package bwval

import (
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwparse"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwtype"
)

// ============================================================================

type PathStr struct {
	S     string
	Bases []bw.ValPath
}

func (v PathStr) Path() (result bw.ValPath, err error) {
	p := bwparse.From(bwrune.S{v.S})
	opt := bwparse.PathOpt{Bases: v.Bases}
	if result, err = bwparse.PathContent(p, opt); err == nil {
		_, err = bwparse.SkipSpace(p, bwparse.TillEOF)
	}
	return
}

// ============================================================================

func MustPath(pathProvider bw.ValPathProvider) (result bw.ValPath) {
	var err error
	if result, err = pathProvider.Path(); err != nil {
		bwerr.PanicErr(bwerr.Refine(err, "invalid path: {Error}"))
	}
	return
}

// ============================================================================

func MustDef(pp bwrune.ProviderProvider, optBaseProvider ...bw.ValPathProvider) (result Def) {
	var err error
	p := bwparse.From(pp)
	if _, err = bwparse.SkipSpace(p, bwparse.TillNonEOF); err != nil {
		err = bwerr.Refine(err, "failed to MustDef: {Error}")
		return
	}
	var st bwparse.Status
	if result, st = ParseDef(p); !st.IsOK() {
		err = st.Err
	} else {
		_, err = bwparse.SkipSpace(p, bwparse.TillEOF)
	}
	if err != nil {
		bwerr.PanicErr(err)
	}
	return
}

// ============================================================================

// MustPathVal - must-обертка bw.Val.PathVal()
func MustPathVal(v bw.Val, pathProvider bw.ValPathProvider, optVars ...map[string]interface{}) (result interface{}) {
	var err error
	path := MustPath(pathProvider)
	if result, err = v.PathVal(path, optVars...); err != nil {
		bwerr.PanicA(bwerr.Err(bwerr.Refine(err,
			ansiMustPathValFailed,
			path, bwjson.Pretty(v), varsJSON(path, optVars),
		)))
	}
	return result
}

// MustSetPathVal - must-обертка bw.Val.SetPathVal()
func MustSetPathVal(val interface{}, v bw.Val, pathProvider bw.ValPathProvider, optVars ...map[string]interface{}) {
	var err error
	path := MustPath(pathProvider)
	if err = v.SetPathVal(val, path, optVars...); err != nil {
		bwerr.PanicErr(bwerr.Refine(err,
			ansiMustSetPathValFailed,
			path, bwjson.Pretty(v), varsJSON(path, optVars),
		))
	}
}

// ============================================================================

type F struct {
	S    string
	Vars map[string]interface{}
}

func (v F) GetVal() interface{} {
	return FromTemplate(TemplateFrom(bwrune.F{v.S}), v.Vars)
}

type S struct {
	S    string
	Vars map[string]interface{}
}

func (v S) GetVal() interface{} {
	return FromTemplate(TemplateFrom(bwrune.S{v.S}), v.Vars)
}

type T struct {
	T    Template
	Vars map[string]interface{}
}

func (v T) GetVal() interface{} {
	return FromTemplate(v.T, v.Vars)
}

type V struct {
	Val interface{}
}

func (v V) GetVal() interface{} {
	return v.Val
}

type fromProvider interface {
	GetVal() interface{}
}

func From(a fromProvider, optPathProvider ...bw.ValPathProvider) (result Holder) {
	result = Holder{}
	result.Val = a.GetVal()
	// FromTemplate(a.Template(), a.Vars())
	if len(optPathProvider) > 0 {
		result.Pth = MustPath(optPathProvider[0])
	}
	return
}

type Template struct {
	val interface{}
}

func TemplateFrom(pp bwrune.ProviderProvider) (result Template) {
	var err error
	var val interface{}
	p := bwparse.From(pp)
	if _, err = bwparse.SkipSpace(p, bwparse.TillNonEOF); err != nil {
		err = bwerr.Refine(err, "failed to TemplateFrom: {Error}")
		return
	}
	var st bwparse.Status
	if val, st = bwparse.Val(p); !st.IsOK() {
		err = st.Err
	} else {
		_, err = bwparse.SkipSpace(p, bwparse.TillEOF)
	}
	if err != nil {
		bwerr.PanicErr(err)
	}
	return Template{val: val}
}

func FromTemplate(template Template, optVars ...map[string]interface{}) (result interface{}) {
	var err error
	if result, err = expandPaths(template.val, template.val, true, optVars...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func expandPaths(val interface{}, rootVal interface{}, isRoot bool, optVars ...map[string]interface{}) (result interface{}, err error) {
	var path bw.ValPath
	var ok bool
	if path, ok = val.(bw.ValPath); ok {
		var h Holder
		if isRoot {
			h = Holder{}
		} else {
			h = Holder{Val: rootVal}
		}
		result, err = h.PathVal(path, optVars...)
	} else {
		result = val
		switch t, kind := bwtype.Kind(val, bwtype.ValKindSetFrom(bwtype.ValMap, bwtype.ValArray)); kind {
		case bwtype.ValMap:
			m := t.(map[string]interface{})
			for key, val := range m {
				if val, err = expandPaths(val, rootVal, false, optVars...); err != nil {
					return
				}
				m[key] = val
			}
		case bwtype.ValArray:
			vals := t.([]interface{})
			for i, val := range vals {
				if val, err = expandPaths(val, rootVal, false, optVars...); err != nil {
					return
				}
				vals[i] = val
			}
		}
	}
	return
}

// ============================================================================
