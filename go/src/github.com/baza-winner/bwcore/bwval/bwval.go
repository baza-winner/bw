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
	p := bwparse.From(bwrune.FromString(v.S))
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

func From(s string, optVars ...map[string]interface{}) (result interface{}) {
	return FromTemplate(TemplateFrom(s), optVars...)
}

type Template struct {
	val interface{}
}

func TemplateFrom(s string) (result Template) {
	var err error
	var val interface{}
	if val, err = func(s string) (result interface{}, err error) {
		defer func() {
			if err != nil {
				result = nil
			}
		}()
		p := bwparse.From(bwrune.FromString(s))
		var st bwparse.Status
		if result, st = bwparse.Val(p); st.IsOK() {
			_, err = bwparse.SkipSpace(p, bwparse.TillEOF)
		} else {
			err = st.Err
		}
		return
	}(s); err != nil {
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
