package bwval

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwparse"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/bwtype"
)

// ============================================================================

type Def struct {
	Types      bwtype.ValKindSet
	IsOptional bool
	Enum       bwset.String
	Range      bwtype.Range
	Keys       map[string]Def
	Elem       *Def
	ArrayElem  *Def
	Default    interface{}
	IsArrayOf  bool
}

func (v Def) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{}
	result["Types"] = v.Types
	result["IsOptional"] = v.IsOptional
	if v.IsArrayOf {
		result["IsArrayOf"] = v.IsArrayOf
	}
	if v.Enum != nil {
		result["Enum"] = v.Enum
	}
	if v.Range.Kind() != bwtype.RangeNo {
		result["Range"] = v.Range
	}
	if v.Keys != nil {
		result["Keys"] = v.Keys
	}
	if v.Elem != nil {
		result["Elem"] = *(v.Elem)
	}
	if v.ArrayElem != nil {
		result["ArrayElem"] = *(v.ArrayElem)
	}
	if v.Default != nil {
		result["Default"] = v.Default
	}
	return json.Marshal(result)
}

// ============================================================================

var ansiCanNotBeMixed string

func init() {
	ansiCanNotBeMixed = ansi.String("<ansiErr>%s<ansi> can not be mixed with <ansiVal>%s")
}

var defSupportedTypes bwtype.ValKindSet

func init() {
	defSupportedTypes = bwtype.ValKindSetFrom(
		bwtype.ValString,
		bwtype.ValBool,
		bwtype.ValInt,
		bwtype.ValNumber,
		bwtype.ValMap,
		bwtype.ValArray,
	)
}

func ParseDef(p bwparse.I, optBaseProvider ...bw.ValPathProvider) (result Def, status bwparse.Status) {

	var base bw.ValPath
	if len(optBaseProvider) > 0 {
		base = MustPath(optBaseProvider[0])
	} else {
		base = bw.ValPath{bw.ValPathItem{Type: bw.ValPathItemVar, Key: "Def"}}
	}
	string2type := map[string]bwtype.ValKind{}
	for _, t := range defSupportedTypes.ToSlice() {
		s := t.String()
		string2type[s] = t
	}
	var (
		tp                  bwtype.ValKind
		ok                  bool
		val                 interface{}
		isOptionalSpecified bool
	)
	result.Types = bwtype.ValKindSet{}
	addType := func(on bwparse.On, tp bwtype.ValKind) (err error) {
		switch tp {
		case bwtype.ValNumber, bwtype.ValInt:
			dontMix := bwtype.ValKindSetFrom(bwtype.ValInt, bwtype.ValNumber)
			for _, vk := range dontMix.ToSlice() {
				if result.Types.Has(vk) {
					err = p.Error(bwparse.E{Start: on.Start, Fmt: bw.Fmt(ansiCanNotBeMixed, tp, vk)})
					return
				}
			}
		case bwtype.ValArray:
			if result.IsArrayOf {
				err = p.Error(bwparse.E{Start: on.Start, Fmt: bw.Fmt(ansiCanNotBeMixed, bwtype.ValArray, "ArrayOf")})
				return
			}
		}
		result.Types.Add(tp)
		return
	}
	processType := func(on bwparse.On, s string) (ok bool, err error) {
		if tp, ok = string2type[s]; ok {
			err = addType(on, tp)
		} else if ok = s == "ArrayOf"; ok {
			if result.Types.Has(bwtype.ValArray) {
				err = p.Error(bwparse.E{Start: on.Start, Fmt: bw.Fmt(ansiCanNotBeMixed, "ArrayOf", bwtype.ValArray)})
			} else {
				result.IsArrayOf = true
			}
		}
		return
	}
	validKeys := bwset.StringFrom("type", "range", "enum", "keys", "elem", "arrayElem", "default", "isOptional")
	setForType := func(opt bwparse.Opt) bwparse.Opt {
		opt.KindSet.Add(bwtype.ValId, bwtype.ValString, bwtype.ValArray)
		opt.OnId = func(on bwparse.On, s string) (val interface{}, ok bool, err error) {
			ok, err = processType(on, s)
			return
		}
		opt.OnValidateString = func(on bwparse.On, s string) (err error) {
			if ok, err = processType(on, s); !ok && err == nil {
				err = bwparse.Unexpected(p, on.Start)
			}
			return
		}
		opt.OnParseArrayElem = func(on bwparse.On, vals []interface{}) (outVals []interface{}, status bwparse.Status) {
			var val interface{}
			kindSet := on.Opt.KindSet
			defer func() { on.Opt.KindSet = kindSet }()
			on.Opt.KindSet = bwtype.ValKindSetFrom(bwtype.ValId, bwtype.ValString)
			if val, status = bwparse.Val(p, *on.Opt); status.IsOK() {
				outVals = append(vals, val)
			}
			return
		}
		opt.OnValidateArray = func(on bwparse.On, vals []interface{}) (err error) {
			if len(vals) == 0 {
				err = bwparse.Expects(p, bwparse.Unexpected(p, on.Start), "non empty <ansiType>Array<ansi>")
			}
			return
		}
		return opt
	}
	checkArrayOf := func(start *bwparse.Start) (err error) {
		if result.IsArrayOf && len(result.Types) == 0 {
			err = p.Error(bwparse.E{
				Start: start,
				Fmt:   bw.A{Fmt: "<ansiType>ArrayOf<ansi> must be followed by another <ansiVar>Type<ansi>"},
			})
		}
		return err
	}
	if val, status = bwparse.Val(p, setForType(bwparse.Opt{
		ValFalse: true,
		KindSet:  bwtype.ValKindSetFrom(bwtype.ValMap),
		OnValidateMapKey: func(on bwparse.On, m map[string]interface{}, key string) (err error) {
			if !validKeys.Has(key) {
				err = p.Error(bwparse.E{
					Start: on.Start,
					Fmt:   bw.Fmt(ansi.String("unexpected key `<ansiErr>%s<ansi>`"), on.Start.Suffix()),
				})
			} else {
				switch key {
				case "default":
					if len(result.Types) == 0 {
						err = p.Error(bwparse.E{
							Start: on.Start,
							Fmt:   bw.A{Fmt: "key <ansiVar>type<ansi> must be specified first"},
						})
					} else if isOptionalSpecified && !result.IsOptional {
						err = p.Error(bwparse.E{
							Start: on.Start,
							Fmt:   bw.Fmt("while <ansiVar>isOptional <ansiVal>false<ansi>, key <ansiVar>%s<ansi> can not be specified", key),
						})
					}
				default:
					if result.Default != nil {
						err = p.Error(bwparse.E{
							Start: on.Start,
							Fmt:   bw.Fmt("key <ansiVar>default<ansi> must be LAST specified, but found key <ansiErr>%s<ansi> after it", key),
						})
					}
				}
			}
			return
		},
		OnParseMapElem: func(on bwparse.On, m map[string]interface{}, key string) (status bwparse.Status) {
			switch key {
			case "type":
				if _, status = bwparse.Val(p, setForType(bwparse.Opt{KindSet: bwtype.ValKindSet{}})); status.IsOK() {
					validKeys = bwset.StringFrom("isOptional", "default")
					for _, vk := range result.Types.ToSlice() {
						switch vk {
						case bwtype.ValString:
							validKeys.Add("enum")
						case bwtype.ValInt, bwtype.ValNumber:
							validKeys.Add("range")
						case bwtype.ValMap:
							validKeys.Add("keys")
							validKeys.Add("elem")
						case bwtype.ValArray:
							validKeys.Add("arrayElem")
							validKeys.Add("elem")
						}
					}
					if status.Err = checkArrayOf(status.Start); status.Err == nil {
						if result.Enum != nil && !result.Types.Has(bwtype.ValString) {
							status.Err = p.Error(bwparse.E{
								Start: status.Start,
								Fmt:   bw.A{Fmt: "key <ansiVar>enum<ansi> is specified, so value of key <ansiVar>type<ansi> expects to have <ansiVal>String<ansi>"},
							})
						} else if result.Range.Kind() != bwtype.RangeNo && !(result.Types.Has(bwtype.ValInt) || result.Types.Has(bwtype.ValNumber)) {
							status.Err = p.Error(bwparse.E{
								Start: status.Start,
								Fmt:   bw.A{Fmt: "key <ansiVar>range<ansi> is specified, so value of key <ansiVar>type<ansi> expects to have <ansiVal>Int<ansi> or <ansiVal>Number<ansi>"},
							})
						} else if result.Keys != nil && !result.Types.Has(bwtype.ValMap) {
							status.Err = p.Error(bwparse.E{
								Start: status.Start,
								Fmt:   bw.A{Fmt: "key <ansiVar>keys<ansi> is specified, so value of key <ansiVar>type<ansi> expects to have <ansiVal>Map<ansi>"},
							})
						} else if result.ArrayElem != nil && !result.Types.Has(bwtype.ValArray) {
							status.Err = p.Error(bwparse.E{
								Start: status.Start,
								Fmt:   bw.A{Fmt: "key <ansiVar>arrayElem<ansi> is specified, so value of key <ansiVar>type<ansi> expects to have <ansiVal>Array<ansi>"},
							})
						} else if result.Elem != nil && !(result.Types.Has(bwtype.ValMap) || result.ArrayElem == nil && result.Types.Has(bwtype.ValArray)) {
							var suffix string
							if result.ArrayElem == nil {
								suffix = " or <ansiVal>Array<ansi>"
							}
							status.Err = p.Error(bwparse.E{
								Start: status.Start,
								Fmt:   bw.A{Fmt: "key <ansiVar>elem<ansi> is specified, so value of key <ansiVar>type<ansi> expects to have <ansiVal>Map<ansi>" + suffix},
							})
						}
					}
				}
			case "enum":
				_, status = bwparse.Val(p, bwparse.Opt{
					KindSet: bwtype.ValKindSetFrom(bwtype.ValString, bwtype.ValArray),
					OnValidateString: func(on bwparse.On, s string) (err error) {
						if result.Enum == nil {
							result.Enum = bwset.String{}
						}
						result.Enum.Add(s)
						return
					},
					OnParseArrayElem: func(on bwparse.On, vals []interface{}) (outVals []interface{}, status bwparse.Status) {
						var s string
						if s, status = bwparse.String(p, *on.Opt); status.IsOK() {
							if result.Enum == nil {
								result.Enum = bwset.String{}
							}
							result.Enum.Add(s)
							outVals = append(vals, s)
						} else if status.Err == nil {
							status.Err = bwparse.Expects(p, nil, "<ansiType>String<ansi>")
						}
						return
					},
				})
			case "range":
				if result.Range, status = bwparse.Range(p, bwparse.Opt{KindSet: bwtype.ValKindSetFrom(bwtype.ValNumber)}); !status.OK && status.Err == nil {
					status.Err = bwparse.Expects(p, nil, "<ansiType>Range<ansi>")
				}
			case "keys":
				if _, status = bwparse.Map(p, bwparse.Opt{
					OnParseMapElem: func(on bwparse.On, m map[string]interface{}, key string) (status bwparse.Status) {
						if result.Keys == nil {
							result.Keys = map[string]Def{}
						}
						if result.Keys[key], status = ParseDef(p, base.AppendKey(key)); status.IsOK() {
							m[key] = nil
						} else if status.Err == nil {
							status.Err = bwparse.Expects(p, nil, "<ansiType>Def<ansi>")
						}
						return
					},
				}); !status.OK && status.Err == nil {
					status.Err = bwparse.Expects(p, nil, "<ansiType>Map<ansi>")
				}
			case "elem", "arrayElem":
				var def Def
				if def, status = ParseDef(p, base.AppendKey(key)); status.IsOK() {
					if key == "elem" {
						result.Elem = &def
					} else {
						result.ArrayElem = &def
					}
				} else if status.Err == nil {
					status.Err = bwparse.Expects(p, nil, "<ansiType>Def<ansi>")
				}
			case "isOptional":
				result.IsOptional, status = bwparse.Bool(p)
				isOptionalSpecified = true
			case "default":
				result.Default, status = ParseValByDef(p, result, base.AppendKey(key))
				if result.Default != nil {
					result.IsOptional = true
				}
			}
			m[key] = nil
			return
		},
	})); status.IsOK() {
		if _, ok := val.(map[string]interface{}); !ok {
			status.Err = checkArrayOf(status.Start)
		}
	}
	return
}

// ============================================================================

func ParseValByDef(p bwparse.I, def Def, optBase ...bw.ValPath) (result interface{}, status bwparse.Status) {
	var base bw.ValPath
	if len(optBase) > 0 {
		base = optBase[0]
	}
	return parseValByDef(p, def, base, false)
}

func parseValByDef(p bwparse.I, def Def, base bw.ValPath, skipArrayOf bool) (result interface{}, status bwparse.Status) {

	opt := bwparse.Opt{KindSet: bwtype.ValKindSet{}}
	opt.KindSet.AddSet(def.Types)

	if opt.KindSet.Has(bwtype.ValMap) {
		opt.OnValidateMapKey = func(on bwparse.On, m map[string]interface{}, key string) (err error) {
			if def.Elem == nil && def.Keys != nil {
				if _, ok := def.Keys[key]; !ok {
					err = p.Error(bwparse.E{
						Start: on.Start,
						Fmt:   bw.Fmt(ansi.String("unexpected key `<ansiErr>%s<ansi>`"), on.Start.Suffix()),
					})
				}
			}
			return
		}
		opt.OnParseMapElem = func(on bwparse.On, m map[string]interface{}, key string) (status bwparse.Status) {
			var keyDef *Def
			if def.Keys != nil {
				if v, ok := def.Keys[key]; ok {
					keyDef = &v
				}
			}
			if keyDef == nil {
				keyDef = def.Elem
			}
			if keyDef == nil {
				keyDef = &Def{Types: defSupportedTypes}
			}
			var val interface{}
			if val, status = parseValByDef(p, *keyDef, base.AppendKey(key), skipArrayOf); status.IsOK() {
				m[key] = val
			}
			return
		}
	}

	if hasArray := opt.KindSet.Has(bwtype.ValArray); hasArray || !skipArrayOf && def.IsArrayOf {
		elemDef := func() Def {
			if !hasArray {
				return def
			} else {
				var elemDef *Def
				if def.ArrayElem != nil {
					elemDef = def.ArrayElem
				}
				if elemDef == nil {
					elemDef = def.Elem
				}
				if elemDef == nil {
					elemDef = &Def{Types: defSupportedTypes}
				}
				return *elemDef
			}
		}
		subSkipArrayOf := skipArrayOf
		if !hasArray {
			opt.KindSet.Add(bwtype.ValArray)
			subSkipArrayOf = true
		}
		opt.OnParseArrayElem = func(on bwparse.On, vals []interface{}) (outVals []interface{}, status bwparse.Status) {
			var val interface{}
			if val, status = parseValByDef(p, elemDef(), base.AppendIdx(len(vals)), subSkipArrayOf); status.IsOK() {
				outVals = append(vals, val)
			}
			return
		}
		opt.OnValidateArrayOfStringElem = func(on bwparse.On, ss []string, s string) (err error) {
			h := Holder{Val: s, Pth: base.AppendIdx(len(ss))}
			if _, err = h.validVal(elemDef(), subSkipArrayOf); err != nil {
				// bwdebug.Print("err", err)
				err = p.Error(bwparse.E{
					Start: on.Start,
					Fmt:   bwerr.Err(err),
				})
			}
			return
		}
	}

	if result, status = bwparse.Val(p, opt); status.IsOK() {
		h := Holder{Val: result, Pth: base}
		if result, status.Err = h.validVal(def, skipArrayOf); status.Err != nil {
			status.Err = p.Error(bwparse.E{
				Start: status.Start,
				Fmt:   bwerr.Err(status.Err),
			})
		}
	}
	return
}

// ============================================================================
