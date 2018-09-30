package defparser

import (
	"github.com/baza-winner/bw/core"
	"log"
)

func Parse(source string) (result interface{}, err error) {
	var pos int
	var char rune
	pfa := pfaStruct{stack: parseStack{}, state: expectSpaceOrValue}
	for pos, char = range source {
		wasState := pfa.state
		if err = pfa.processCharAtPos(pos, char); err != nil {
			log.Printf(`%v`, err)
			if _, ok := err.(unexpectedCharError); ok {
				return nil, core.Error("unexpected <ansiOutline>char <ansiPrimaryLiteral>%+q<ansi> (code <ansiSecondaryLiteral>%v<ansi>) at <ansiOutline>pos <ansiSecondaryLiteral>%d<ansi> while <ansiOutline>state <ansiSecondaryLiteral>%s", char, char, pos, wasState)
			} else {
				log.Printf(`%v is NOT unexpectedCharError`, err)
			}
			return nil, err
		}
	}

	switch len(pfa.stack) {
	case 0:
		return nil, nil
	case 1:
		switch pfa.state {
		case expectDigitOrUnderscore, expectDigitOrUnderscoreOrDot, expectWord:
			if err = pfa.finishTopStackItem(nil); err != nil {
				return nil, err
			}
		}
		switch pfa.state {
		case tokenFinished:
			result = pfa.stack[0].value
		default:
			return nil, core.Error("unexpected <ansiOutline>state<ansi> <ansiPrimaryLiteral>%s<ansi> while at end of source", pfa.state)
			return
		}
	default:
		return nil, core.Error("<ansiOutline>stack<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expects to have one item at end of source and while <ansiOutline>state<ansi> <ansiSecondaryLiteral>%s", pfa.stack, pfa.state)
	}

	return
}
