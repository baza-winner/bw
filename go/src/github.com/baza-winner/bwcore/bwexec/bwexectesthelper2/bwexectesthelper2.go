// Вспомогательная утилита для тестирования bwexec.ExecCmd.
package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/baza-winner/bwcore/ansi"
	_ "github.com/baza-winner/bwcore/ansi/tags"
	"github.com/baza-winner/bwcore/bwexec"
	"github.com/baza-winner/bwcore/bwos"
)

func init() {
	ansi.MustAddTag("ansiHeader",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorWhite}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdFaint),
	)
}

func main() {
	exitOnErrorFlag := flag.Bool("e", false, "exitOnError")
	verbosityFlag := flag.String("v", "none", "verbosity: none, err, ok, all, allBrief")
	silentFlag := flag.String("s", "none", "silent: none, stderr, stdout, all")
	displayFlag := flag.Bool("d", false, "display exitCode, stdout, stderr, output")
	noColorFlag := flag.Bool("n", false, "no color")
	flag.Parse()
	if flag.NArg() < 1 {
		bwos.Exit(1, `<ansiPath>bwexectesthelper2<ansi> expects at least one arg`)
	}
	argsWithoutProg := flag.Args()

	display := *displayFlag
	ansi.NoColor = *noColorFlag

	ret := bwexec.MustCmd(
		bwexec.Args(argsWithoutProg[0], argsWithoutProg[1:]...),
		map[string]interface{}{
			"verbosity":   *verbosityFlag,
			"exitOnError": *exitOnErrorFlag,
			"silent":      *silentFlag,
		},
	)
	if display {
		exitCode := ret[`exitCode`].(int)
		ansiExitCode := `<ansiOK>`
		if exitCode != 0 {
			ansiExitCode = `<ansiErr>`
		}
		fmt.Printf(ansi.String(`<ansiHeader>===== exitCode: `+ansiExitCode+"%d\n"), exitCode)
		fmt.Println(ansi.String(`<ansiHeader>===== stdout:`))
		fmt.Println(strings.Join(ret[`stdout`].([]string), "\n"))
		fmt.Println(ansi.String(`<ansiHeader>===== stderr:`))
		fmt.Println(strings.Join(ret[`stderr`].([]string), "\n"))

		if output, ok := ret[`output`]; ok {
			fmt.Println(ansi.String(`<ansiHeader>===== output:`))
			fmt.Println(strings.Join(output.([]string), "\n"))
		}
	}
}
