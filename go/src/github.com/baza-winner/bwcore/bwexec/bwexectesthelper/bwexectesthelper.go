// Вспомогательная утилита для тестирования bwexec.ExecCmd.
package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/baza-winner/bwcore/ansi"
)

func main() {
	colorOutputFlag := flag.Bool("color", false, "color output")
	exitCodeFlag := flag.Int("exit", 0, "exit code")
	flag.Parse()
	isColorOuput := *colorOutputFlag
	exitCode := *exitCodeFlag

	findStdTagRegexp, _ := regexp.Compile(`<(?:stdout|stderr)>`)

	argsWithoutProg := flag.Args()
	for _, arg := range argsWithoutProg {
		indexes := findStdTagRegexp.FindAllStringIndex(arg, -1)
		if indexes != nil {
			isStdout := true
			lastLineEnd := 0
			for _, beginEnd := range indexes {
				if beginEnd[0] > lastLineEnd {
					helper(arg, lastLineEnd, beginEnd[0], isStdout, isColorOuput)
				}
				tag := arg[beginEnd[0]:beginEnd[1]]
				isStdout = tag == `<stdout>`
				lastLineEnd = beginEnd[1]
			}
			helper(arg, lastLineEnd, len(arg), isStdout, isColorOuput)
		}
	}
	os.Exit(exitCode)
}

var findBrTagRegexp, _ = regexp.Compile(`<br>`)

func helper(arg string, begin, end int, isStdout bool, isColorOuput bool) {
	line := arg[begin:end]
	line = findBrTagRegexp.ReplaceAllLiteralString(line, "\n")
	if isColorOuput {
		defaultAnsi := `ansiOK`
		if !isStdout {
			defaultAnsi = `ansiErr`
		}
		line = ansi.StringA(ansi.A{Default: ansi.MustTag(defaultAnsi), S: line})
	}
	stream := os.Stdout
	if !isStdout {
		stream = os.Stderr
	}
	fmt.Fprint(stream, line)
}
