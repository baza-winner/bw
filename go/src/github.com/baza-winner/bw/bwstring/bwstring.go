package bwstring

import (
	"fmt"
	"strings"
)

func SmartQuote(ss ...string) (result string) {
	result = ``
	for i, s := range ss {
		if i > 0 {
			result += ` `
		}
		if strings.ContainsAny(s, ` "`) {
			result += fmt.Sprintf(`%q`, s)
		} else {
			result += s
		}
	}
	return
}
