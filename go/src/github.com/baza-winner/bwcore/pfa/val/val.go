package val

import (
	"unicode"

	"github.com/baza-winner/bwcore/bwerror"
)

//go:generate stringer -type=UnicodeCategory

type UnicodeCategory uint8

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

func (v UnicodeCategory) conforms(pfa *PfaStruct, val interface{}, varPath VarPath) (result bool) {
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
	pfa.traceCondition(varPath, v, result)
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

type Map struct{}

func (v Map) GetVal(pfa *PfaStruct) interface{} {
	return map[string]interface{}{}
}

func (v Map) GetSource(pfa *PfaStruct) formatted.String {
	return pfa.traceVal(v)
}

func (v Map) String() string {
	return "Map"
}

func (t Map) FormattedString() formatted.String {
	return formatted.StringFrom("<ansiOutline>%s", t)
}

// ============================================================================

type Array struct{}

func (v Array) GetVal(pfa *PfaStruct) interface{} {
	return []interface{}{}
}

func (v Array) GetSource(pfa *PfaStruct) formatted.String {
	return pfa.traceVal(v)
}

func (v Array) String() string {
	return "Array"
}

func (t Array) FormattedString() formatted.String {
	return formatted.StringFrom("<ansiOutline>%s", t)
}
