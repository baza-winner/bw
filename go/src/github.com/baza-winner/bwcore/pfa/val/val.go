package val

import (
	"unicode"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/pfa/common"
	"github.com/baza-winner/bwcore/pfa/core"
	"github.com/baza-winner/bwcore/pfa/formatted"
)

// ============================================================================

func valProviderFrom(i interface{}) (result core.ValProvider, err error) {
	switch t := i.(type) {
	case Array:
		result = t
	case Map:
		result = t
	case Var:
		// result = common.VarIsFrom(t.VarPathStr)
		result = common.VarVal{core.MustVarPathFrom(t.VarPathStr)}
	default:
		result = common.JustVal{i}
	}
	return
}

func MustValProviderFrom(i interface{}) (result core.ValProvider) {
	var err error
	if result, err = valProviderFrom(i); err != nil {
		bwerror.PanicErr(err)
	}
	return
}

// ============================================================================

type UnicodeCategory uint8

//go:generate stringer -type=UnicodeCategory

const (
	UnicodeSpace UnicodeCategory = iota
	UnicodeLetter
	UnicodeLower
	UnicodeUpper
	UnicodeDigit
	UnicodeOpenBraces
	UnicodePunct
	UnicodeSymbol
)

func (v UnicodeCategory) Conforms(pfa *core.PfaStruct, val interface{}, varPath core.VarPath) (result bool) {
	if r, ok := val.(rune); ok {
		switch v {
		case UnicodeSpace:
			result = unicode.IsSpace(r)
		case UnicodeLetter:
			result = unicode.IsLetter(r) || r == '_'
		case UnicodeLower:
			result = unicode.IsLower(r)
		case UnicodeUpper:
			result = unicode.IsUpper(r)
		case UnicodeDigit:
			result = unicode.IsDigit(r)
		case UnicodeOpenBraces:
			result = r == '(' || r == '{' || r == '[' || r == '<'
		case UnicodePunct:
			result = unicode.IsPunct(r)
		case UnicodeSymbol:
			result = unicode.IsSymbol(r)
		default:
			bwerror.Panic("UnicodeCategory: %s", v)
		}
	}
	if pfa.TraceLevel > core.TraceNone {
		pfa.TraceCondition(varPath, v, result)
	}
	return
}

func (t UnicodeCategory) FormattedString() formatted.String {
	return formatted.StringFrom("<ansiOutline>%s", t)
}

// ============================================================================

type EOF struct{}

func (v EOF) String() string {
	return "EOF"
}

func (t EOF) FormattedString() formatted.String {
	return formatted.StringFrom("<ansiOutline>%s", t)
}

// ============================================================================

type Var struct {
	VarPathStr string
}

func (t Var) GetChecker() core.ValChecker {
	return common.VarVal{core.MustVarPathFrom(t.VarPathStr)}
}

// ============================================================================

type Map struct{}

func (v Map) GetVal(pfa *core.PfaStruct) interface{} {
	return map[string]interface{}{}
}

func (v Map) GetSource(pfa *core.PfaStruct) formatted.String {
	return pfa.TraceVal(v)
}

func (v Map) String() string {
	return "Map"
}

func (t Map) FormattedString() formatted.String {
	return formatted.StringFrom("<ansiOutline>%s", t)
}

// ============================================================================

type Array struct{}

func (v Array) GetVal(pfa *core.PfaStruct) interface{} {
	return []interface{}{}
}

func (v Array) GetSource(pfa *core.PfaStruct) formatted.String {
	return pfa.TraceVal(v)
}

func (v Array) String() string {
	return "Array"
}

func (t Array) FormattedString() formatted.String {
	return formatted.StringFrom("<ansiOutline>%s", t)
}
