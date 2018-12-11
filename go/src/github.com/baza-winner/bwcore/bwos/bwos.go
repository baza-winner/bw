// Package bwos содержит bw-дополнение для package os.
package bwos

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
)

// ============================================================================

// ShortenFileSpec укорачивает строку за счет замены префикса, совпадающиего (если) cо значением
// ${HOME} (переменная среды), на символ `~`
func ShortenFileSpec(fileSpec string) string {
	fileSpec = filepath.Clean(fileSpec)
	if homeDir := os.Getenv("HOME"); homeDir != "" {
		homeDir = filepath.Clean(homeDir)
		if homeDir[len(homeDir)-1] != '/' {
			homeDir += string('/')
		}
		if len(fileSpec) >= len(homeDir) && fileSpec[:len(homeDir)] == homeDir {
			fileSpec = "~/" + fileSpec[len(homeDir):]
		}
	}
	return fileSpec
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
	result = err.S
	if !newlineAtTheEnd.MatchString(ansi.ChopReset(result)) {
		result += string('\n')
	}
	return
}

// ============================================================================
