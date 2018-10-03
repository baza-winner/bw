package defparse

import (
	"fmt"
	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerror"
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

func (pfa *pfaStruct) arrangeError(err error, source string) error {
	if err != nil {
		if _, ok := err.(unexpectedCharError); ok {
			if pfa.charPtr == nil {
				suffix := getSuffix(source, 0, 0)
				err = bwerror.Error("<ansiReset>unexpected end of string (pfa.state: %s)"+suffix, pfa.state)
			} else {
				suffix := getSuffix(source, uint(pfa.pos), 1)
				err = bwerror.Error("<ansiReset>unexpected char <ansiPrimaryLiteral>%q<ansiReset> (charCode: %v, pfa.state: %s)"+suffix, *pfa.charPtr, *pfa.charPtr, pfa.state)
			}
		} else if _, ok := err.(failedToGetNumberError); ok {
			stackItem := pfa.getTopStackItemOfType(parseStackItemNumber)
			suffix := getSuffix(source, uint(stackItem.pos), uint(len(stackItem.itemString)))
			err = bwerror.Error("<ansiReset>failed to get number from string <ansiPrimaryLiteral>" + stackItem.itemString + suffix)
		} else if _, ok := err.(unknownWordError); ok {
			stackItem := pfa.getTopStackItemOfType(parseStackItemWord)
			suffix := getSuffix(source, uint(stackItem.pos), uint(len(stackItem.itemString)))
			err = bwerror.Error("<ansiReset>unknown word <ansiPrimaryLiteral>" + stackItem.itemString + suffix)
		} else if _, ok := err.(unexpectedWordError); ok {
			stackItem := pfa.getTopStackItemOfType(parseStackItemWord)
			suffix := getSuffix(source, uint(stackItem.pos), uint(len(stackItem.itemString)))
			err = bwerror.Error("<ansiReset>unexpected word <ansiPrimaryLiteral>" + stackItem.itemString + suffix)
		}
	}
	return err
}

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
