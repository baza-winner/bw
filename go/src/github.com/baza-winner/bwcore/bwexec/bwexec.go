// Package bwexec предоставялет функцию ExecCmd.
package bwexec

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwos"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwstr"
	"github.com/baza-winner/bwcore/bwval"
)

const defaultFailedCode = 1

var cmdOptDef bwval.Def

func init() {
	cmdOptDef = bwval.MustDef(bwrune.S{`
		{
			type Map
			keys {
				verbosity {
					type String
					enum <all err ok none>
					default "none"
				}
				silent {
					type String
					enum <none stderr stdout all>
					default "all"
				}
				exitOnError {
					type Bool
					default false
				}
				workDir {
					type String
					isOptional true
				}
			}
		}
	`})
}

type A struct {
	Cmd  string
	Args []string
}

func Args(cmdName string, cmdArgs ...string) A {
	return A{Cmd: cmdName, Args: cmdArgs}
}

func MustCmd(a A, optOpt ...interface{}) (result map[string]interface{}) {
	var err error
	if result, err = Cmd(a, optOpt...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func Cmd(a A, optOpt ...interface{}) (result map[string]interface{}, err error) {
	var opt interface{}
	if len(optOpt) > 0 {
		opt = optOpt[0]
	}
	hOpt := bwval.From(bwval.V{opt}, bwval.PathStr{S: "Cmd.opt"}).MustValid(cmdOptDef)

	result = map[string]interface{}{}

	cmdTitle := bwstr.SmartQuote(append([]string{a.Cmd}, a.Args...)...)
	optSilent := hOpt.MustPathStr("silent").MustString()
	optVerbosity := hOpt.MustPathStr("verbosity").MustString()
	optWorkDir := hOpt.MustPathStr("workDir?").MustString("")
	var pwd string
	if optWorkDir != "" {
		if pwd, err = os.Getwd(); err != nil {
			return
		}
		if err = os.Chdir(optWorkDir); err != nil {
			return
		}
	}
	if optVerbosity == `all` || optVerbosity == `allBrief` {
		fmt.Println(ansi.String(`<ansiPath>` + cmdTitle + `<ansi> . . .`))
	}
	cmd := exec.Command(a.Cmd, a.Args...)

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		bwos.Exit(1, "Error creating StdoutPipe for Cmd: %v", err)
	}
	stdoutScanner := bufio.NewScanner(cmdStdout)

	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		bwos.Exit(1, "Error creating StderrPipe for Cmd: %v", err)
	}
	stderrScanner := bufio.NewScanner(cmdStderr)

	stdout := []string{}
	stderr := []string{}
	// output := []string{}

	go func() {
		for stdoutScanner.Scan() {
			s := stdoutScanner.Text()
			stdout = append(stdout, s)
			// output = append(output, s)
			if !(optSilent == `all` || optSilent == `stdout`) {
				fmt.Fprintln(os.Stdout, s)
			}
		}
	}()

	go func() {
		for stderrScanner.Scan() {
			s := stderrScanner.Text()
			stderr = append(stderr, s)
			// output = append(output, s)
			if !(optSilent == `all` || optSilent == `stderr`) {
				fmt.Fprintln(os.Stderr, stderrScanner.Text())
			}
		}
	}()

	if err := cmd.Start(); err != nil {
		bwos.Exit(1, "Error starting Cmd: %v", err)
	}

	// https://stackoverflow.com/questions/10385551/get-exit-code-go
	exitCode := 0
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			} else {
				log.Printf(ansi.String("<ansiWarn>Could not get exit code for failed program: <ansiPath>%s"), cmdTitle)
				exitCode = defaultFailedCode
			}
		} else {
			bwos.Exit(1, "cmd.Wait: %v", err)
		}
	}

	var ansiName, prefix string
	if exitCode == 0 && (optVerbosity == `all` || optVerbosity == `allBrief` || optVerbosity == `ok`) {
		ansiName, prefix = `ansiOK`, `OK`
	}
	if exitCode != 0 && (optVerbosity == `all` || optVerbosity == `allBrief` || optVerbosity == `err`) {
		ansiName, prefix = `ansiErr`, `ERR`
	}
	if len(prefix) > 0 {
		fmt.Println(ansi.StringA(ansi.A{Default: ansi.MustTag(ansiName), S: prefix + `: <ansiPath>` + cmdTitle}))
	}
	if hOpt.MustPathStr("exitOnError").MustBool() && exitCode != 0 {
		os.Exit(exitCode)
	}
	if pwd != "" {
		if err = os.Chdir(pwd); err != nil {
			return
		}
	}
	result[`exitCode`] = exitCode
	result[`stdout`] = stdout
	result[`stderr`] = stderr
	// result[`output`] = output
	return
}
