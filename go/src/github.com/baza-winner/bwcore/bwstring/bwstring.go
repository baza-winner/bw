// Предоставялет функции для работы со строками.
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

// func GetPluralWord(count int, word string, word1 string, word2_4 string, _word5more ...string) (result string) {
// 	var word5more string
// 	if _word5more != nil {
// 		word5more = _word5more[0]
// 	}
// 	if len(word5more) == 0 {
// 		word5more = word2_4
// 	}
// 	result = word5more
// 	decimal := count / 10 % 10
// 	if decimal != 1 {
// 		unit := count % 10
// 		if unit == 1 {
// 			result = word1
// 		} else if 2 <= unit && unit <= 4 {
// 			result = word2_4
// 		}
// 	}
// 	return word + result
// }
