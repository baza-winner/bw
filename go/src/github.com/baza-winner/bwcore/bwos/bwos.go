/*
Предоставялет функцию ShortenFileSpec.
*/
package bwos

import (
	"fmt"
	"os"
	"regexp"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
)

// ============================================================================

// ShortenFileSpec укорачивает строку за счет замены префикса, совпадающиего (если) cо значением
// ${HOME} (переменная среды), на символ `~`
func ShortenFileSpec(s string) (result string) {
	home := os.Getenv(`HOME`)
	result = s
	if len(result) >= len(home) && result[0:len(home)] == home {
		result = `~` + result[len(home):len(result)]
	}
	return
}

func Exit(exitCode int, fmtString string, fmtArgs ...interface{}) {
	ExitA(exitCode, bw.A{fmtString, fmtArgs})
}

// Exit with exitCode and message defined by bw.I
func ExitA(exitCode int, a bw.I) {
	fmt.Print(exitMsg(a))
	os.Exit(exitCode)
}

// ============================================================================

var newlineAtTheEnd, _ = regexp.Compile(`\n\s*$`)

func exitMsg(a bw.I) (result string) {
	err := bwerr.FromA(a)
	result = err.Ansi
	if !newlineAtTheEnd.MatchString(ansi.ChopReset(result)) {
		result += string('\n')
	}
	return
}

// ============================================================================
