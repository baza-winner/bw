package pathparse

import (
	"fmt"

	"github.com/baza-winner/bwcore/bwerr"
	"github.com/jimlawless/whereami"
)

type pfaErrorType uint16

const (
	pfaErrorBelow pfaErrorType = iota
	unexpectedRuneError
	failedToGetNumberError
	// unknownWordError
	pfaErrorAbove
)

//go:generate stringer -type=pfaErrorType

type pfaError struct {
	pfa       *pfaStruct
	errorType pfaErrorType
	fmtString string
	fmtArgs   []interface{}
	Where     string
}

func pfaErrorMake(pfa *pfaStruct, errorType pfaErrorType, args ...interface{}) (result pfaError) {
	if !(pfaErrorBelow < errorType && errorType < pfaErrorAbove) {
		bwerr.Panic(" errorType == %s", errorType)
	}
	fmtString, fmtArgs := pfaErrorValidators[errorType](pfa, args...)
	result = pfaError{pfa, errorType, fmtString, fmtArgs, whereami.WhereAmI(2)}
	return
}

func (err pfaError) Error() string {
	return bwerr.From(err.fmtString, err.fmtArgs...).Error()
}

func (v pfaError) DataForJSON() interface{} {
	result := map[string]interface{}{}
	result["pfa"] = v.pfa.DataForJSON()
	result["errorType"] = v.errorType.String()
	result["Where"] = v.Where
	return result
}

type pfaErrorValidator func(pfa *pfaStruct, args ...interface{}) (string, []interface{})

var pfaErrorValidators = map[pfaErrorType]pfaErrorValidator{
	unexpectedRuneError:    _unexpectedRuneError,
	failedToGetNumberError: _failedToGetNumberError,
	// unknownWordError:       _unknownWordError,
}

func pfaErrorValidatorsCheck() {
	pfaErrorType := pfaErrorBelow + 1
	for pfaErrorType < pfaErrorAbove {
		if _, ok := pfaErrorValidators[pfaErrorType]; !ok {
			bwerr.Panic("not defined <ansiVar>pfaErrorValidators<ansi>[<ansiVal>%s<ansi>]", pfaErrorType)
		}
		pfaErrorType += 1
	}
}

func _unexpectedRuneError(pfa *pfaStruct, args ...interface{}) (fmtString string, fmtArgs []interface{}) {
	if args != nil {
		bwerr.Panic("does not expect args instead of <ansiVal>%#v", args)
	}
	if pfa.curr.runePtr == nil {
		suffix := getSuffix(pfa, pfa.curr, "")
		fmtString = "unexpected end of string (pfa.state: %s)" + suffix
		fmtArgs = []interface{}{pfa.state}
	} else {
		rune := *pfa.curr.runePtr
		suffix := getSuffix(pfa, pfa.curr, string(rune))
		fmtString = "unexpected char <ansiVal>%q<ansiReset> (charCode: %v, pfa.state: %s)" + suffix
		fmtArgs = []interface{}{rune, rune, pfa.state}
	}
	return
}

func _failedToGetNumberError(pfa *pfaStruct, args ...interface{}) (fmtString string, fmtArgs []interface{}) {
	if args != nil {
		bwerr.Panic("does not expect args instead of <ansiVal>%#v", args)
	}
	stackItem := pfa.getTopStackItemOfType(parseStackItemNumber)
	suffix := getSuffix(pfa, stackItem.start, stackItem.itemString)
	return "failed to get number from string <ansiVal>%s" + suffix, []interface{}{stackItem.itemString}
}

// func _unknownWordError(pfa *pfaStruct, args ...interface{}) (fmtString string, fmtArgs []interface{}) {
// 	if args != nil {
// 		bwerr.Panic("does not expect args instead of <ansiVal>%#v", args)
// 	}
// 	stackItem := pfa.getTopStackItemOfType(parseStackItemWord)
// 	suffix := getSuffix(pfa, stackItem.start, stackItem.itemString)
// 	return "unknown word <ansiVal>%s" + suffix, []interface{}{stackItem.itemString}
// }

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
		suffix += fmt.Sprintf(" at pos <ansiPath>%d<ansi>", start.pos)
		separator = " "
	} else {
		suffix += fmt.Sprintf(" at line <ansiPath>%d<ansi>, col <ansiPath>%d<ansi> (pos <ansiPath>%d<ansi>)", start.line, start.col, start.pos)
	}
	suffix += ":" + separator + "<ansiDarkGreen>"

	suffix += pfa.curr.prefix[0 : start.pos-pfa.curr.prefixStart]
	if pfa.curr.runePtr != nil {
		suffix += "<ansiLightRed>"
		suffix += redString
		suffix += "<ansiReset>"
		for pfa.curr.runePtr != nil && postLineCount > 0 {
			pfa.pullRune()
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
