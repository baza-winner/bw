package bwval

import (
	"fmt"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
)

var (
	ansiIsNotOfType          string
	ansiMustPathValFailed    string
	ansiType                 string
	ansiWrongVal             string
	ansiVars                 string
	ansiVarsIsNil            string
	ansiMustSetPathValFailed string
)

func init() {
	ansiIsNotOfType = ansi.String("<ansiVal>%#v<ansi> is not <ansiType>%s")
	ansiMustSetPathValFailed = ansi.String("Failed to set <ansiPath>%s<ansi> of <ansiVal>%s<ansi>%s: {Error}")
	ansiMustPathValFailed = ansi.String("Failed to get <ansiPath>%s<ansi> of <ansiVal>%s<ansi>%s: {Error}")
	valPathPrefix := "<ansiPath>%s<ansi> "
	ansiWrongVal = ansi.String(valPathPrefix + "is <ansiErr>%#v")

	ansiType = ansi.String("<ansiType>%s")
	ansiVars = ansi.String(" with <ansiVar>vars<ansi> <ansiVal>%s<ansi>")
	ansiVarsIsNil = ansi.String("<ansiVar>vars<ansi> is <ansiErr>nil")
}

func (v Holder) notOfTypeError(expectedType string, optExpectedType ...string) (result error) {
	if len(optExpectedType) == 0 {
		result = bwerr.From(v.ansiString()+ansi.String(" is not <ansiType>%s"), expectedType)
	} else {
		expectedTypes := fmt.Sprintf(ansiType, expectedType)
		for i, elem := range optExpectedType {
			expectedTypes += typeSeparator[i == len(optExpectedType)-1] + fmt.Sprintf(ansiType, elem)
		}
		result = bwerr.From(v.ansiString()+ansi.String(" is none of %s"), expectedTypes)
	}
	return
}

var typeSeparator = map[bool]string{
	true:  " or ",
	false: ", ",
}

func (v Holder) notEnoughRangeError(l int, idx int) error {
	return bwerr.From(
		v.ansiString()+
			ansi.String(" has not enough length (<ansiVal>%d<ansi>) for idx (<ansiVal>%d)"),
		l, idx,
	)
}

func (v Holder) nonSupportedValueError() error {
	return bwerr.From(v.ansiString() + ansi.String(" is <ansiErr>non supported<ansi> value"))
}

func readonlyPathError(path bw.ValPath) error {
	return bwerr.From(ansi.String("<ansiPath>%s<ansi> is <ansiErr>readonly path"), path)
}

func (v Holder) outOfRangeError(rng Range) (err error) {
	var s string
	switch RangeKind(rng) {
	case RangeMinMax:
		s = ansi.String(" is <ansiErr>out of range<ansi> <ansiVal>%s")
	case RangeMin:
		s = ansi.String(" is <ansiErr>less<ansi> than<ansiVal>%s")
	case RangeMax:
		s = ansi.String(" is <ansiErr>more<ansi> than<ansiVal>%s")
	}
	if len(s) > 0 {
		err = bwerr.From(v.ansiString()+s, RangeString(rng))
	}
	return
}

func (v Holder) maxLessThanMinError() error {
	return bwerr.From(v.ansiString() + "<ansiVar>max<ansi> must not be less then <ansiVar>min")
}

func (v Holder) unexpectedKeysError(unexpectedKeys []string) (err error) {
	var fmtString string
	var fmtArg interface{}
	switch len(unexpectedKeys) {
	case 0:
	case 1:
		fmtString = ansi.String(`has unexpected key <ansiVal>%s`)
		fmtArg = unexpectedKeys[0]
	default:
		fmtString = `has unexpected keys <ansiVal>%s`
		fmtArg = bwjson.Pretty(unexpectedKeys)
	}
	if len(fmtString) > 0 {
		err = bwerr.From(v.ansiString()+fmtString, fmtArg)
	}
	return
}

func (v Holder) wrongValError() error {
	return bwerr.From(ansiWrongVal, v.path, v.val)
}
