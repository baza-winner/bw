/*
Предоставялет функцию ShortenFileSpec.
*/
package bwos

import (
	"os"
)

/*
Укорачивает строку за счет замены префикса, совпадающиего (если) cо значением
${HOME} (переменная среды), на символ `~`
*/
func ShortenFileSpec(s string) (result string) {
	home := os.Getenv(`HOME`)
	result = s
	if len(result) >= len(home) && result[0:len(home)] == home {
		result = `~` + result[len(home):len(result)]
	}
	return
}
