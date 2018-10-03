/*
Вспомогательная утилита для тестирования bwexec.ExecCmd.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwexec"
	"strings"
)

func main() {
	exitOnErrorFlag := flag.Bool("e", false, "exitOnError")
	verbosityFlag := flag.String("v", "none", "verbosity: none, err, ok, all, allBrief")
	silentFlag := flag.String("s", "none", "silent: none, stderr, stdout, all")
	displayFlag := flag.Bool("d", false, "display exitCode, stdout, stderr, output")
	noColorFlag := flag.Bool("n", false, "no color")
	flag.Parse()
	if flag.NArg() < 1 {
		bwerror.ExitWithError(1, `<ansiCmd>bwexectesthelper2<ansi> expects at least one arg`)
	}
	argsWithoutProg := flag.Args()

	exitOnError := *exitOnErrorFlag
	verbosity := *verbosityFlag
	silent := *silentFlag
	display := *displayFlag
	ansi.NoColor = *noColorFlag

	ret := bwexec.ExecCmd(map[string]interface{}{`v`: verbosity, `s`: silent, `exitOnError`: exitOnError}, argsWithoutProg[0], argsWithoutProg[1:]...)
	if display {
		exitCode := ret[`exitCode`].(int)
		ansiExitCode := `<ansiOK>`
		if exitCode != 0 {
			ansiExitCode = `<ansiErr>`
		}
		fmt.Printf(ansi.Ansi(`Header`, `===== exitCode: `+ansiExitCode+"%d\n"), exitCode)
		fmt.Println(ansi.Ansi(`Header`, `===== stdout:`))
		fmt.Println(strings.Join(ret[`stdout`].([]string), "\n"))
		fmt.Println(ansi.Ansi(`Header`, `===== stderr:`))
		fmt.Println(strings.Join(ret[`stderr`].([]string), "\n"))

		if output, ok := ret[`output`]; ok {
			fmt.Println(ansi.Ansi(`Header`, `===== output:`))
			fmt.Println(strings.Join(output.([]string), "\n"))
		}
	}
}
