package defparser

import (
	"log"
)

const (
	expectSpaceOrCurlyBraceOrBracketOrDigitOrSignOrDoubleQuoteOrSingleQuoteOrEOF = iota
)

func Parse(source string) (result map[string]interface{}, err error) {
	state := expectSpaceOrCurlyBraceOrBracketOrDigitOrSignOrDoubleQuoteOrSingleQuoteOrEOF
	for pos, char := range source {
		switch state {
		case expectSpaceOrCurlyBraceOrBracketOrDigitOrSignOrDoubleQuoteOrSingleQuoteOrEOF:
			switch char {
			case ' ', '\t':
			case '{':
			case '[':
			case '-', '+':
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			case '"':
			case '\'':
			}
			log.Printf("%d", pos)
		default:
			log.Panicf("unknown state %d", state)
		}
		// if state == expectSpaceOrCurlyBraceOrBracketOrDigitOrSignOrDoubleQuoteOrSingleQuote {
		// }
	}
	switch state {
	case expectSpaceOrCurlyBraceOrBracketOrDigitOrSignOrDoubleQuoteOrSingleQuoteOrEOF:
	default:
		log.Panicf("unknown state %d", state)
	}
	return
}
