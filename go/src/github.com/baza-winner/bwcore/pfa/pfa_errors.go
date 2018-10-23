package pfa

// type ErrorType uint16

// const (
// 	// pfaErrorBelo
// 	NoError ErrorType = iota
// 	UnexpectedRune
// 	FailedToGetNumber
// 	UnknownWord
// 	// pfaErrorAbove
// )

type pfaError struct {
	pfa     *pfaStruct
	errName string
	errStr  string
	Where   string
}

// func pfaErrorMake(pfa *pfaStruct, errorType ErrorType, args ...interface{}) (result pfaError) {
// 	// if NoError == errorType && errorType < pfaErrorAbove) {
// 	// 	bwerror.Panic(" errorType == %s", errorType)
// 	// }
// 	var errStr string
// 	switch errorType {
// 	case UnexpectedRune:
// 		// errStr = pfa.p.UnexpectedRuneError(fmt.Sprintf("pfa.state: %s", pfa.state)).Error()
// 		errStr = pfa.p.UnexpectedRuneError().Error()
// 	case FailedToGetNumber:
// 		stackItem := pfa.getTopStackItem()
// 		itemString, _ := stackItem.vars["string"].(string)
// 		errStr = pfa.p.WordError("failed to get number from string <ansiPrimary>%s<ansi>", itemString, stackItem.start).Error()
// 	case UnknownWord:
// 		stackItem := pfa.getTopStackItem()
// 		itemString, _ := stackItem.vars["string"].(string)
// 		errStr = pfa.p.WordError("unknown word <ansiPrimary>%s<ansi>", itemString, stackItem.start).Error()
// 	default:
// 		bwerror.Unreachable()
// 	}
// 	result = pfaError{pfa, errorType, errStr, whereami.WhereAmI(2)}
// 	return
// }

func (err pfaError) Error() string {
	return err.errStr
}

func (v pfaError) DataForJSON() interface{} {
	result := map[string]interface{}{}
	result["pfa"] = v.pfa.DataForJSON()
	result["errName"] = v.errName
	result["errStr"] = v.errStr
	result["Where"] = v.Where
	return result
}

// type pfaErrorValidator func(pfa *pfaStruct, args ...interface{}) (string, []interface{})

// var pfaErrorValidators = map[ErrorType]pfaErrorValidator{
// 	UnexpectedRune:    _unexpectedRuneError,
// 	FailedToGetNumber: _failedToGetNumberError,
// 	UnknownWord:       _unknownWordError,
// }

// func pfaErrorValidatorsCheck() {
// 	ErrorType := pfaErrorBelow + 1
// 	for ErrorType < pfaErrorAbove {
// 		if _, ok := pfaErrorValidators[ErrorType]; !ok {
// 			bwerror.Panic("not defined <ansiOutline>pfaErrorValidators<ansi>[<ansiPrimary>%s<ansi>]", ErrorType)
// 		}
// 		ErrorType += 1
// 	}
// }

// func _unexpectedRuneError(pfa *pfaStruct, args ...interface{}) (fmtString string, fmtArgs []interface{}) {
// 	if args != nil {
// 		bwerror.Panic("does not expect args instead of <ansiSecondary>%#v", args)
// 	}
// 	if pfa.p.Curr.RunePtr == nil {
// 		suffix := pfa.p.GetSuffix(pfa.p.Curr, "")
// 		fmtString = "unexpected end of string (pfa.state: %s)" + suffix
// 		fmtArgs = []interface{}{pfa.state}
// 	} else {
// 		rune := *pfa.p.Curr.RunePtr
// 		suffix := pfa.p.GetSuffix(pfa.p.Curr, string(rune))
// 		fmtString = "unexpected char <ansiPrimary>%q<ansiReset> (charCode: %v, pfa.state: %s)" + suffix
// 		fmtArgs = []interface{}{rune, rune, pfa.state}
// 	}
// 	return
// }

// func _failedToGetNumberError(pfa *pfaStruct, args ...interface{}) (fmtString string, fmtArgs []interface{}) {
// 	if args != nil {
// 		bwerror.Panic("does not expect args instead of <ansiSecondary>%#v", args)
// 	}
// 	stackItem := pfa.getTopStackItemOfType("number")
// 	itemString, _ := stackItem.vars["string"].(string)
// 	suffix := pfa.p.GetSuffix(stackItem.start, itemString)
// 	return "failed to get number from string <ansiPrimary>%s<ansi>" + suffix, []interface{}{itemString}
// }

// func _unknownWordError(pfa *pfaStruct, args ...interface{}) (fmtString string, fmtArgs []interface{}) {
// 	if args != nil {
// 		bwerror.Panic("does not expect args instead of <ansiSecondary>%#v", args)
// 	}
// 	stackItem := pfa.getTopStackItemOfType("word")
// 	itemString, _ := stackItem.vars["string"].(string)
// 	suffix := pfa.p.GetSuffix(stackItem.start, itemString)
// 	return "unknown word <ansiPrimary>%s<ansi>" + suffix, []interface{}{itemString}
// }

// =============

// func getSuffix(pfa *pfaStruct, pos, length uint, opts ...uint) (suffix string) {
// func getSuffix(pfa *pfaStruct, start runePtrStruct, redString string) (suffix string) {
// 	preLineCount := pfa.p.preLineCount
// 	postLineCount := pfa.p.postLineCount
// 	if pfa.p.curr.runePtr == nil {
// 		preLineCount += postLineCount
// 	}

// 	separator := "\n"
// 	if pfa.p.curr.line <= 1 {
// 		suffix += fmt.Sprintf(" at pos <ansiCmd>%d<ansi>", start.pos)
// 		separator = " "
// 	} else {
// 		suffix += fmt.Sprintf(" at line <ansiCmd>%d<ansi>, col <ansiCmd>%d<ansi> (pos <ansiCmd>%d<ansi>)", start.line, start.col, start.pos)
// 	}
// 	suffix += ":" + separator + "<ansiDarkGreen>"

// 	suffix += pfa.p.curr.prefix[0 : start.pos-pfa.p.curr.prefixStart]
// 	if pfa.p.curr.runePtr != nil {
// 		suffix += "<ansiLightRed>"
// 		suffix += redString
// 		suffix += "<ansiReset>"
// 		for pfa.p.curr.runePtr != nil && postLineCount > 0 {
// 			pfa.PullRune()
// 			if pfa.p.curr.runePtr != nil {
// 				suffix += string(*pfa.p.curr.runePtr)
// 				if *pfa.p.curr.runePtr == '\n' {
// 					postLineCount -= 1
// 				}
// 			}
// 		}
// 	}
// 	if byte(suffix[len(suffix)-1]) != '\n' {
// 		suffix += string('\n')
// 	}
// 	return suffix
// }
