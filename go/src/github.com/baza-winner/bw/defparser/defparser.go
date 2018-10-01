package defparser

import (
	"fmt"
	"github.com/baza-winner/bw/core"
	// "strconv"
)

type unexpectedCharError struct{}

func (e unexpectedCharError) Error() string {
	return `unexpectedCharError`
}

type failedToGetNumberError struct{}

func (e failedToGetNumberError) Error() string {
	return `failedToGetNumberError`
}

type unknownWordError struct{}

func (e unknownWordError) Error() string {
	return `unknownWordError`
}

type unexpectedWordError struct{}

func (e unexpectedWordError) Error() string {
	return `unexpectedWordError`
}

func Parse(source string) (result interface{}, err error) {
	pfa := pfaStruct{stack: parseStack{}, state: parseState{primary: expectValueOrSpace}}
	for pos, char := range source {
		if err = pfa.processCharAtPos(char, pos); err != nil {
			break
		}
	}
	if err == nil {
		err = pfa.processEOF()
	}
	if err != nil {
		if _, ok := err.(unexpectedCharError); ok {
			charTitle := "unexpected end of string"
			if pfa.charPtr != nil {
				charTitle = fmt.Sprintf("unexpected <ansiOutline>char <ansiPrimaryLiteral>%+q<ansi> (code <ansiSecondaryLiteral>%v<ansi>) at <ansiOutline>pos <ansiSecondaryLiteral>%d<ansi>", *pfa.charPtr, *pfa.charPtr, pfa.pos)
			}
			err = core.Error(charTitle+" while <ansiOutline>state <ansiSecondaryLiteral>%s", pfa.state)
		} else if _, ok := err.(failedToGetNumberError); ok {
			stackItem := pfa.getTopStackItemOfType(parseStackItemNumber)
			err = core.Error("failed to get number from string <ansiPrimaryLiteral>%s<ansi> at pos <ansiSecondaryLiteral>%d", stackItem.itemString, stackItem.pos)
		} else if _, ok := err.(unknownWordError); ok {
			stackItem := pfa.getTopStackItemOfType(parseStackItemWord)
			// err = core.Error("unknown word <ansiPrimaryLiteral>%s<ansi> at pos <ansiSecondaryLiteral>%d", stackItem.itemString, stackItem.pos)
			suffix := getSuffix(source, uint(stackItem.pos), uint(len(stackItem.itemString)))
			err = core.Error("unknown word <ansiPrimaryLiteral>" + stackItem.itemString + suffix)
		} else if _, ok := err.(unexpectedWordError); ok {
			stackItem := pfa.getTopStackItemOfType(parseStackItemWord)
			suffix := getSuffix(source, uint(stackItem.pos), uint(len(stackItem.itemString)))
			err = core.Error("unexpected word <ansiPrimaryLiteral>" + stackItem.itemString + suffix)
		}
	}

	return pfa.result, err
}

func getSuffix(source string, pos, length uint, opts ...uint) (suffix string) {
	preLineCount := uint(2)
	postLineCount := uint(2)
	if opts != nil {
		preLineCount = opts[0]
		if len(opts) >= 2 {
			postLineCount = opts[1]
		}
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
	for int(toPos) < len(source)-1 {
		if source[toPos+1] == byte('\n') {
			postLineCount -= 1
			if postLineCount <= 0 {
				break
			}
		}
		toPos += 1
	}
	suffix = "<ansi>"
	if foundPreBreak {
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
	} else {
		suffix += fmt.Sprintf(" at pos <ansiCmd>%d<ansi>", pos)
	}
	suffix += ":\n<ansiOK>"
	suffix += source[fromPos:pos]
	suffix += "<ansiErr>"
	suffix += source[pos : pos+length]
	suffix += "<ansiReset>"
	suffix += source[pos+length : toPos]
	return
}
