package bwval

import (
	"fmt"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/bwtype"
)

// ============================================================================

var (
	ansiIsNotOfType          string
	ansiMustPathValFailed    string
	ansiType                 string
	ansiWrongVal             string
	ansiUnexpectedEnumValue  string
	ansiVars                 string
	ansiNotEnoughRange       string
	ansiCanNotSetNonString   string
	ansiVarsIsNil            string
	ansiMustSetPathValFailed string
	ansiHasNoKey             string
)

func init() {
	ansiIsNotOfType = ansi.String("<ansiVal>%#v<ansi> is not <ansiType>%s")
	ansiMustSetPathValFailed = ansi.String("Failed to set <ansiPath>%s<ansi> of <ansiVal>%s<ansi>%s: {Error}")

	// ansiMustPathValFailed = ansi.String("Failed to get <ansiPath>%s<ansi> of <ansiVal>%s<ansi>%s: {Error}")
	ansiMustPathValFailed = ansi.String("Failed to get <ansiPath>%s<ansi> of <ansiVal>%s<ansi>%s: {Error}")

	valPathPrefix := "<ansiPath>%s<ansi>"
	ansiWrongVal = ansi.String(valPathPrefix + " is <ansiErr>%#v")
	ansiUnexpectedEnumValue = ansi.String(valPathPrefix + ": expected one of <ansiVal>%s<ansi> instead of <ansiErr>%q")

	ansiNotEnoughRange = ansi.String(" has not enough length (<ansiVal>%d<ansi>) for idx (<ansiVal>%d)")
	ansiCanNotSetNonString = ansi.String(" can not set elem[%d] to non string <ansiVal>%#v")

	ansiType = ansi.String("<ansiType>%s")
	ansiVars = ansi.String(" with <ansiVar>vars<ansi> <ansiVal>%s<ansi>")
	ansiVarsIsNil = ansi.String("<ansiVar>vars<ansi> is <ansiErr>nil")
	ansiHasNoKey = ansi.String(" has no key <ansiVal>%s")
}

// ============================================================================

func readonlyPathError(path bw.ValPath) error {
	return bwerr.From(ansi.String("<ansiPath>%s<ansi> is <ansiErr>readonly path"), path)
}

// ============================================================================

func (v Holder) ansiString() (result string) {
	return fmt.Sprintf(ansi.String("<ansiPath>%s<ansi> (<ansiVal>%s<ansi>)"), v.Pth, bwjson.Pretty(v.Val))
}

// ============================================================================

func (v Holder) hasNoKeyError(key string) error {
	return bwerr.From(v.ansiString()+ansiHasNoKey, key)
}

// ============================================================================

func (v Holder) notOfValKindError(vk ValKindSet) (result error) {
	vks := vk.ToSlice()
	expectedTypes := ""
	for i, elem := range vks {
		expectedTypes += notOfValKindItemSeparator[i > 0] + fmt.Sprintf(ansiType, elem)
	}
	result = bwerr.From(v.ansiString()+ansi.String(" "+notOfValKindInfix[len(vks) == 1]+" %s"), expectedTypes)
	return
}

// https://www.quickanddirtytips.com/education/grammar/when-use-nor
var notOfValKindItemSeparator = map[bool]string{
	true:  " nor ",
	false: "",
}

var notOfValKindInfix = map[bool]string{
	true:  "is not",
	false: "neither",
}

// ============================================================================

func (v Holder) notEnoughRangeError(l int, idx int) error {
	return bwerr.From(v.ansiString()+ansiNotEnoughRange, l, idx)
}

// ============================================================================

func (v Holder) canNotSetNonStringError(idx int, val interface{}) error {
	return bwerr.From(v.ansiString()+ansiCanNotSetNonString, idx, val)
}

// ============================================================================

func (v Holder) nonSupportedValueError() error {
	return bwerr.From(v.ansiString() + ansi.String(" is <ansiErr>non supported<ansi> value"))
}

// ============================================================================

func (v Holder) outOfRangeError(rng bwtype.Range) (err error) {
	var s string
	switch rng.Kind() {
	case bwtype.RangeMinMax:
		s = ansi.String(" is <ansiErr>out of range<ansi> <ansiVal>%s")
	case bwtype.RangeMin:
		s = ansi.String(" is <ansiErr>less<ansi> than<ansiVal>%s")
	case bwtype.RangeMax:
		s = ansi.String(" is <ansiErr>more<ansi> than<ansiVal>%s")
	}
	if len(s) > 0 {
		err = bwerr.From(v.ansiString()+s, rng.String())
	}
	return
}

// ============================================================================

func (v Holder) maxLessThanMinError(max, min bwtype.RangeLimit) error {
	return bwerr.From(v.ansiString()+
		": <ansiPath>.max<ansi> (<ansiVal>%s<ansi>) must not be <ansiErr>less<ansi> then <ansiPath>.min<ansi> (<ansiVal>%s<ansi>)",
		bwjson.Pretty(max), bwjson.Pretty(min),
	)
}

// ============================================================================

func (v Holder) defaultNonOptionalError() error {
	return bwerr.From(v.ansiString() +
		": having <ansiPath>.default<ansi> can not have <ansiPath>.isOptional<ansi> <ansiVal>true",
	)
}

// ============================================================================

func (v Holder) unexpectedKeysError(unexpectedKeys bwset.String) (err error) {
	var fmtString string
	var fmtArg interface{}
	ss := unexpectedKeys.ToSlice()
	switch len(ss) {
	case 0:
		return
	case 1:
		fmtString = ansi.String(" has unexpected key <ansiVal>%s")
		fmtArg = bwjson.Pretty(ss[0])
	default:
		fmtString = ansi.String(" has unexpected keys: <ansiVal>%s")
		fmtArg = bwjson.Pretty(ss)
	}
	err = bwerr.From(v.ansiString()+fmtString, fmtArg)
	return
}

// ============================================================================

func (v Holder) arrayOfMustBeFollowedBySomeTypeError() error {
	return bwerr.From(v.ansiString() + ansi.String(": <ansiVal>ArrayOf<ansi> must be followed by some type, can not be <ansiErr>used alone"))
}

func (v Holder) valuesAreMutuallyExclusiveError(valA, valB interface{}) error {
	return bwerr.From(v.ansiString()+ansi.String(": values <ansiVal>%s<ansi> and <ansiVal>%s<ansi> are <ansiErr>mutually exclusive<ansi>, can not be <ansiErr>used both at once"), bwjson.Pretty(valA), bwjson.Pretty(valB))
}

// ============================================================================

func (v Holder) wrongValError() error {
	return bwerr.From(ansiWrongVal, v.Pth, v.Val)
}

// ============================================================================

func (v Holder) unexpectedEnumValueError(enum bwset.String) error {
	return bwerr.From(ansiUnexpectedEnumValue, v.Pth, bwjson.Pretty(enum), v.Val)
}

// ============================================================================
