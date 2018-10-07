package defparse

import (
	"fmt"
	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/jimlawless/whereami"
)

type pfaErrorType uint16

const (
	pfaError_below_ pfaErrorType = iota
	unexpectedCharError
	failedToGetNumberError
	unknownWordError
	unexpectedWordError
	pfaError_above_
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
	if !(pfaError_below_ < errorType && errorType < pfaError_above_) {
		bwerror.Panic(" errorType == %s", errorType)
	}
	fmtString, fmtArgs := pfaErrorValidators[errorType](pfa, args...)
	result = pfaError{pfa, errorType, fmtString, fmtArgs, whereami.WhereAmI(2)}
	return
}

func (err pfaError) Error() string {
	return bwerror.Error(err.fmtString, err.fmtArgs...).Error()
}

// func (v pfaError) WhereError() (result string) {
// 	result = v.where
// 	return
// }

func (v pfaError) GetDataForJson() interface{} {
	result := map[string]interface{}{}
	result["pfa"] = v.pfa.GetDataForJson()
	result["errorType"] = v.errorType.String()
	result["Where"] = v.Where
	return result
}

type pfaErrorValidator func(pfa *pfaStruct, args ...interface{}) (string, []interface{})

var pfaErrorValidators = map[pfaErrorType]pfaErrorValidator{
	unexpectedCharError:    _unexpectedCharError,
	failedToGetNumberError: _failedToGetNumberError,
	unknownWordError:       _unknownWordError,
	unexpectedWordError:    _unexpectedWordError,
}

func pfaErrorValidatorsCheck() {
	pfaErrorType := pfaError_below_ + 1
	for pfaErrorType < pfaError_above_ {
		if _, ok := pfaErrorValidators[pfaErrorType]; !ok {
			bwerror.Panic("not defined <ansiOutline>pfaErrorValidators<ansi>[<ansiPrimaryLiteral>%s<ansi>]", pfaErrorType)
		}
		pfaErrorType += 1
	}
}

func _unexpectedCharError(pfa *pfaStruct, args ...interface{}) (fmtString string, fmtArgs []interface{}) {
	if args != nil {
		bwerror.Panic("does not expect args instead of <ansiSecondaryLiteral>%#v", args)
	}
	if pfa.charPtr == nil {
		suffix := getSuffix(pfa.source, 0, 0)
		fmtString = "unexpected end of string (pfa.state: %s)" + suffix
		fmtArgs = []interface{}{pfa.state}
	} else {
		suffix := getSuffix(pfa.source, uint(pfa.pos), 1)
		fmtString = "unexpected char <ansiPrimaryLiteral>%q<ansiReset> (charCode: %v, pfa.state: %s)" + suffix
		fmtArgs = []interface{}{
			*pfa.charPtr,
			*pfa.charPtr,
			pfa.state,
		}
	}
	return
}

func _failedToGetNumberError(pfa *pfaStruct, args ...interface{}) (fmtString string, fmtArgs []interface{}) {
	if args != nil {
		bwerror.Panic("does not expect args instead of <ansiSecondaryLiteral>%#v", args)
	}
	stackItem := pfa.getTopStackItemOfType(parseStackItemNumber)
	suffix := getSuffix(pfa.source, uint(stackItem.pos), uint(len(stackItem.itemString)))
	return "failed to get number from string <ansiPrimaryLiteral>%s" + suffix, []interface{}{stackItem.itemString}
}

func _unknownWordError(pfa *pfaStruct, args ...interface{}) (fmtString string, fmtArgs []interface{}) {
	if args != nil {
		bwerror.Panic("does not expect args instead of <ansiSecondaryLiteral>%#v", args)
	}
	stackItem := pfa.getTopStackItemOfType(parseStackItemWord)
	suffix := getSuffix(pfa.source, uint(stackItem.pos), uint(len(stackItem.itemString)))
	return "unknown word <ansiPrimaryLiteral>%s" + suffix, []interface{}{stackItem.itemString}
}

func _unexpectedWordError(pfa *pfaStruct, args ...interface{}) (fmtString string, fmtArgs []interface{}) {
	if args != nil {
		bwerror.Panic("does not expect args instead of <ansiSecondaryLiteral>%#v", args)
	}
	stackItem := pfa.getTopStackItemOfType(parseStackItemWord)
	suffix := getSuffix(pfa.source, uint(stackItem.pos), uint(len(stackItem.itemString)))
	return "unexpected word <ansiPrimaryLiteral>%s" + suffix, []interface{}{stackItem.itemString}
}

// =============

func getSuffix(source string, pos, length uint, opts ...uint) (suffix string) {
	preLineCount := uint(3)
	postLineCount := uint(3)
	if opts != nil {
		preLineCount = opts[0]
		if len(opts) >= 2 {
			postLineCount = opts[1]
		}
	}
	if length == 0 {
		pos = uint(len(source))
		preLineCount += postLineCount
	}

	foundPreBreak := false
	fromPos := pos
	for int(fromPos) >= 1 {
		if source[fromPos-1] == byte('\n') {
			foundPreBreak = true
			preLineCount -= 1
			if preLineCount <= 0 {
				break
			}
		}
		fromPos -= 1
	}
	toPos := pos
	for int(toPos) < len(source) {
		if source[toPos] == byte('\n') {
			postLineCount -= 1
			if postLineCount <= 0 {
				break
			}
		}
		toPos += 1
	}
	separator := "\n"
	if !foundPreBreak {
		suffix += fmt.Sprintf(" at pos <ansiCmd>%d<ansi>", pos)
		separator = " "
	} else {
		fromPos := pos
		line := 1
		col := 1
		foundPreBreak := false
		for int(fromPos) >= 1 {
			if source[fromPos-1] == byte('\n') {
				foundPreBreak = true
				line += 1
			} else if !foundPreBreak {
				col += 1
			}
			fromPos -= 1
		}
		suffix += fmt.Sprintf(" at line <ansiCmd>%d<ansi>, col <ansiCmd>%d<ansi> (pos <ansiCmd>%d<ansi>)", line, col, pos)
	}
	suffix += ":" + separator + "<ansiDarkGreen>"
	suffix += source[fromPos:pos]
	if length > 0 {
		suffix += "<ansiLightRed>"
		suffix += source[pos : pos+length]
		suffix += "<ansiReset>"
		suffix += source[pos+length : toPos]
	}
	if byte(suffix[len(suffix)-1]) != '\n' {
		suffix += string('\n')
	}
	return ansi.Ansi("Reset", suffix)
}
