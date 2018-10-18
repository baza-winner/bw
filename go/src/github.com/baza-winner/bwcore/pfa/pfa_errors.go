package pfa

import (
	"fmt"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/jimlawless/whereami"
)

type ErrorType uint16

const (
	pfaErrorBelow ErrorType = iota
	UnexpectedRune
	FailedToGetNumber
	UnknownWord
	pfaErrorAbove
)

type pfaError struct {
	pfa       *pfaStruct
	errorType ErrorType
	fmtString string
	fmtArgs   []interface{}
	Where     string
}

func pfaErrorMake(pfa *pfaStruct, errorType ErrorType, args ...interface{}) (result pfaError) {
	if !(pfaErrorBelow < errorType && errorType < pfaErrorAbove) {
		bwerror.Panic(" errorType == %s", errorType)
	}
	fmtString, fmtArgs := pfaErrorValidators[errorType](pfa, args...)
	result = pfaError{pfa, errorType, fmtString, fmtArgs, whereami.WhereAmI(2)}
	return
}

func (err pfaError) Error() string {
	return bwerror.Error(err.fmtString, err.fmtArgs...).Error()
}

func (v pfaError) DataForJSON() interface{} {
	result := map[string]interface{}{}
	result["pfa"] = v.pfa.DataForJSON()
	result["errorType"] = v.errorType.String()
	result["Where"] = v.Where
	return result
}

type pfaErrorValidator func(pfa *pfaStruct, args ...interface{}) (string, []interface{})

var pfaErrorValidators = map[ErrorType]pfaErrorValidator{
	UnexpectedRune:    _unexpectedRuneError,
	FailedToGetNumber: _failedToGetNumberError,
	UnknownWord:       _unknownWordError,
}

func pfaErrorValidatorsCheck() {
	ErrorType := pfaErrorBelow + 1
	for ErrorType < pfaErrorAbove {
		if _, ok := pfaErrorValidators[ErrorType]; !ok {
			bwerror.Panic("not defined <ansiOutline>pfaErrorValidators<ansi>[<ansiPrimaryLiteral>%s<ansi>]", ErrorType)
		}
		ErrorType += 1
	}
}

func _unexpectedRuneError(pfa *pfaStruct, args ...interface{}) (fmtString string, fmtArgs []interface{}) {
	if args != nil {
		bwerror.Panic("does not expect args instead of <ansiSecondaryLiteral>%#v", args)
	}
	if pfa.curr.runePtr == nil {
		suffix := getSuffix(pfa, pfa.curr, "")
		fmtString = "unexpected end of string (pfa.state: %s)" + suffix
		fmtArgs = []interface{}{pfa.state}
	} else {
		rune := *pfa.curr.runePtr
		suffix := getSuffix(pfa, pfa.curr, string(rune))
		fmtString = "unexpected char <ansiPrimaryLiteral>%q<ansiReset> (charCode: %v, pfa.state: %s)" + suffix
		fmtArgs = []interface{}{rune, rune, pfa.state}
	}
	return
}

func _failedToGetNumberError(pfa *pfaStruct, args ...interface{}) (fmtString string, fmtArgs []interface{}) {
	if args != nil {
		bwerror.Panic("does not expect args instead of <ansiSecondaryLiteral>%#v", args)
	}
	stackItem := pfa.getTopStackItemOfType("number")
	suffix := getSuffix(pfa, stackItem.start, stackItem.itemString)
	return "failed to get number from string <ansiPrimaryLiteral>%s<ansi>" + suffix, []interface{}{stackItem.itemString}
}

func _unknownWordError(pfa *pfaStruct, args ...interface{}) (fmtString string, fmtArgs []interface{}) {
	if args != nil {
		bwerror.Panic("does not expect args instead of <ansiSecondaryLiteral>%#v", args)
	}
	stackItem := pfa.getTopStackItemOfType("word")
	suffix := getSuffix(pfa, stackItem.start, stackItem.itemString)
	return "unknown word <ansiPrimaryLiteral>%s<ansi>" + suffix, []interface{}{stackItem.itemString}
}

// =============

// func getSuffix(pfa *pfaStruct, pos, length uint, opts ...uint) (suffix string) {
func getSuffix(pfa *pfaStruct, start runePtrStruct, redString string) (suffix string) {
	preLineCount := pfa.preLineCount
	postLineCount := pfa.postLineCount
	if pfa.curr.runePtr == nil {
		preLineCount += postLineCount
	}

	separator := "\n"
	if pfa.curr.line <= 1 {
		suffix += fmt.Sprintf(" at pos <ansiCmd>%d<ansi>", start.pos)
		separator = " "
	} else {
		suffix += fmt.Sprintf(" at line <ansiCmd>%d<ansi>, col <ansiCmd>%d<ansi> (pos <ansiCmd>%d<ansi>)", start.line, start.col, start.pos)
	}
	suffix += ":" + separator + "<ansiDarkGreen>"

	suffix += pfa.curr.prefix[0 : start.pos-pfa.curr.prefixStart]
	if pfa.curr.runePtr != nil {
		suffix += "<ansiLightRed>"
		suffix += redString
		suffix += "<ansiReset>"
		for pfa.curr.runePtr != nil && postLineCount > 0 {
			pfa.PullRune()
			if pfa.curr.runePtr != nil {
				suffix += string(*pfa.curr.runePtr)
				if *pfa.curr.runePtr == '\n' {
					postLineCount -= 1
				}
			}
		}
	}
	if byte(suffix[len(suffix)-1]) != '\n' {
		suffix += string('\n')
	}
	return suffix
}
