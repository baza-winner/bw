package defparser

import (
	"github.com/baza-winner/bw/core"
	// "log"
)

func Parse(source string) (result interface{}, err error) {
	var pos int
	var char rune
	pfa := pfaStruct{stack: parseStack{}, state: parseState{primary: expectValueOrSpace}}
	for pos, char = range source {
		pfa.pos = pos
		pfa.charPtr = &char
		if err = pfa.processCharAtPos(); err != nil {
			if _, ok := err.(unexpectedCharError); ok {
				err = core.Error("unexpected <ansiOutline>char <ansiPrimaryLiteral>%+q<ansi> (code <ansiSecondaryLiteral>%v<ansi>) at <ansiOutline>pos <ansiSecondaryLiteral>%d<ansi> while <ansiOutline>state <ansiSecondaryLiteral>%s", char, char, pos, pfa.state)
			}
			return nil, err
		}
	}
	pfa.pos = -1
	pfa.charPtr = nil
	if err = pfa.processCharAtPos(); err == nil {
		return pfa.result, nil
	} else {
		if _, ok := err.(unexpectedCharError); ok {
			return nil, core.Error("unexpected end of string while <ansiOutline>state <ansiSecondaryLiteral>%s", pfa.state)
		}
		return nil, err
	}
	return
}
