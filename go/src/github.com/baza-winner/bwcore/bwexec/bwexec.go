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
	"github.com/baza-winner/bwcore/bwos"
	"github.com/baza-winner/bwcore/bwstr"
)

const defaultFailedCode = 1

func ExecCmd(opt map[string]interface{}, cmdName string, cmdArgs ...string) (result map[string]interface{}) {
	// _ = GetValidMap(`ExecCmd.opt`, opt, map[string]interface{}{
	//  `type`: `map`,
	//  `keys`: map[string]interface{}{
	//    `v`: map[string]interface{}{
	//      `type`:    `enum`,
	//      `enum`:    []string{`all`, `err`, `ok`, `none`},
	//      `default`: `none`,
	//    },
	//    `s`: map[string]interface{}{
	//      `type`:    `enum`,
	//      `enum`:    []string{`none`, `stderr`, `stdout`, `all`},
	//      `default`: `all`,
	//    },
	//    `exitOnError`: map[string]interface{}{
	//      `type`:    `bool`,
	//      `default`: false,
	//    },
	//  },
	// })

	result = map[string]interface{}{}

	cmdTitle := bwstr.SmartQuote(append([]string{cmdName}, cmdArgs...)...)
	sOpt := getStringKeyOrDefault(opt, `s`, `all`)
	vOpt := getStringKeyOrDefault(opt, `v`, `none`)
	if vOpt == `all` || vOpt == `allBrief` {
		fmt.Println(ansi.String(`<ansiPath>` + cmdTitle + `<ansi> . . .`))
	}
	cmd := exec.Command(cmdName, cmdArgs...)

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
			if !(sOpt == `all` || sOpt == `stdout`) {
				fmt.Fprintln(os.Stdout, s)
			}
		}
	}()

	go func() {
		for stderrScanner.Scan() {
			s := stderrScanner.Text()
			stderr = append(stderr, s)
			// output = append(output, s)
			if !(sOpt == `all` || sOpt == `stderr`) {
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
	if exitCode == 0 && (vOpt == `all` || vOpt == `allBrief` || vOpt == `ok`) {
		ansiName, prefix = `ansiOK`, `OK`
	}
	if exitCode != 0 && (vOpt == `all` || vOpt == `allBrief` || vOpt == `err`) {
		ansiName, prefix = `ansiErr`, `ERR`
	}
	if len(prefix) > 0 {
		fmt.Println(ansi.StringA(ansi.A{Default: ansi.MustTag(ansiName), S: prefix + `: <ansiPath>` + cmdTitle}))
	}
	if getBoolKeyOrDefault(opt, `exitOnError`, false) && exitCode != 0 {
		os.Exit(exitCode)
	}
	result[`exitCode`] = exitCode
	result[`stdout`] = stdout
	result[`stderr`] = stderr
	// result[`output`] = output
	return result
}

func getBoolKeyOrDefault(m map[string]interface{}, keyName string, defaultValue bool) (result bool) {
	result = defaultValue
	if m != nil {
		if val, ok := m[keyName]; ok {
			if typedVal, ok := val.(bool); ok {
				result = typedVal
			}
		}
	}
	return
}

func getStringKeyOrDefault(m map[string]interface{}, keyName string, defaultValue string) (result string) {
	result = defaultValue
	if m != nil {
		if val, ok := m[keyName]; ok {
			if typedVal, ok := val.(string); ok {
				result = typedVal
			}
		}
	}
	return
}
