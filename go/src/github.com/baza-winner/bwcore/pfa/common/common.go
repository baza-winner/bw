package common

import (
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/pfa/core"
	"github.com/baza-winner/bwcore/pfa/formatted"
)

type VarVal struct {
	Path core.VarPath
}

func (v VarVal) Conforms(pfa *core.PfaStruct, val interface{}, Path core.VarPath) (result bool) {
	varValue := pfa.VarValue(v.Path)
	if pfa.Err == nil {
		result = varValue.Val == val
		if pfa.TraceLevel > core.TraceNone {
			pfa.TraceCondition(Path, v.Path, result)
		}
	}
	return
}

func (v VarVal) GetVal(pfa *core.PfaStruct) (result interface{}) {
	varValue := pfa.VarValue(v.Path)
	if pfa.Err == nil {
		result = varValue.Val
	}
	return
}

func (v VarVal) GetSource(pfa *core.PfaStruct) formatted.String {
	return pfa.TraceVal(v.Path)
}

// ============================================================================

type JustVal struct {
	Val interface{}
}

func (v JustVal) Conforms(pfa *core.PfaStruct, val interface{}, varPath core.VarPath) (result bool) {
	result = val == v.Val
	if pfa.TraceLevel > core.TraceNone {
		pfa.TraceCondition(varPath, val, result)
	}
	return
}

func (v JustVal) GetVal(pfa *core.PfaStruct) interface{} {
	return v.Val
}

func (v JustVal) GetSource(pfa *core.PfaStruct) formatted.String {
	return pfa.TraceVal(v.Val)
}

// ============================================================================

type VarIs struct {
	varPath     core.VarPath
	valCheckers []core.ValChecker
	isNil       bool
	runeSet     bwset.Rune
	intSet      bwset.Int
	strSet      bwset.String
}

func VarIsFrom(varPathStr string) *VarIs {
	return &VarIs{varPath: core.MustVarPathFrom(varPathStr)}
}

func (v *VarIs) SetIsNil() {
	v.isNil = true
}

func (v *VarIs) AddRune(r rune) {
	if v.runeSet == nil {
		v.runeSet = bwset.Rune{}
	}
	v.runeSet.Add(r)
}

func (v *VarIs) AddInt(i int) {
	if v.intSet == nil {
		v.intSet = bwset.Int{}
	}
	v.intSet.Add(i)
}

func (v *VarIs) AddStr(s string) {
	if v.strSet == nil {
		v.strSet = bwset.String{}
	}
	v.strSet.Add(s)
}

func (v *VarIs) AddValChecker(vc core.ValChecker) {
	if v.valCheckers == nil {
		v.valCheckers = []core.ValChecker{}
	}
	v.valCheckers = append(v.valCheckers, vc)
}

func (v *VarIs) ConformsTo(pfa *core.PfaStruct) (result bool) {
	if v.varPath[0].Val == "rune" {
		var ofs int
		if len(v.varPath) > 1 {
			ofs, _ = core.VarValueFrom(v.varPath[1].Val).Int()
		}
		r, isEOF := pfa.Proxy.Rune(ofs)
		if v.isNil {
			result = isEOF
			if pfa.TraceLevel > core.TraceNone {
				pfa.TraceCondition(v.varPath, formatted.StringFrom("<ansiOutline>EOF"), result)
			}
		}
		if !result && !isEOF {
			if v.runeSet != nil {
				result = v.runeSet.Has(r)
				if pfa.TraceLevel > core.TraceNone {
					pfa.TraceCondition(v.varPath, v.runeSet, result)
				}
			}
			if !result {
				for _, p := range v.valCheckers {
					if p.Conforms(pfa, r, v.varPath) {
						result = true
						break
					}
				}
			}
		}
	} else if v.varPath[0].Val == "stackLen" {
		i := len(pfa.Stack)
		if v.intSet != nil {
			result = v.intSet.Has(i)
			if pfa.TraceLevel > core.TraceNone {
				pfa.TraceCondition(v.varPath, v.intSet, result)
			}
		}
		if !result {
			for _, p := range v.valCheckers {
				if p.Conforms(pfa, i, v.varPath) {
					result = true
					break
				}
			}
		}
	} else {
		varValue := pfa.VarValue(v.varPath)
		if varValue.Val == nil {
			result = v.isNil
			if pfa.TraceLevel > core.TraceNone {
				pfa.TraceCondition(v.varPath, nil, result)
			}
		} else if r, err := varValue.Rune(); err == nil {
			if v.runeSet != nil {
				result = v.runeSet.Has(r)
				if pfa.TraceLevel > core.TraceNone {
					pfa.TraceCondition(v.varPath, v.runeSet, result)
				}
			}
		} else if i, err := varValue.Int(); err == nil {
			if v.intSet != nil {
				result = v.intSet.Has(i)
				if pfa.TraceLevel > core.TraceNone {
					pfa.TraceCondition(v.varPath, v.intSet, result)
				}
			}
		} else if s, err := varValue.String(); err == nil {
			if v.strSet != nil {
				result = v.strSet.Has(s)
				if pfa.TraceLevel > core.TraceNone {
					pfa.TraceCondition(v.varPath, v.strSet, result)
				}
			}
		}
		if !result && pfa.Err == nil {
			for _, p := range v.valCheckers {
				if p.Conforms(pfa, varValue.Val, v.varPath) {
					result = true
					break
				} else if pfa.Err != nil {
					break
				}
			}
		}
	}
	return
}
