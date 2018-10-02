package bwexec

import (
	"bufio"
	// "encoding/json"
	// "errors"
	"fmt"
	"github.com/baza-winner/bw/ansi"
	"github.com/baza-winner/bw/bwerror"
	"github.com/baza-winner/bw/bwstring"
	// "github.com/iancoleman/strcase"
	// "github.com/jimlawless/whereami"
	"log"
	"os"
	"os/exec"
	// "strings"
	"syscall"
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

	cmdTitle := bwstring.SmartQuote(append([]string{cmdName}, cmdArgs...)...)
	sOpt := getStringKeyOrDefault(opt, `s`, `all`)
	vOpt := getStringKeyOrDefault(opt, `v`, `none`)
	if vOpt == `all` || vOpt == `allBrief` {
		fmt.Println(ansi.Ansi(``, `<ansiCmd>`+cmdTitle+`<ansi> . . .`))
	}
	cmd := exec.Command(cmdName, cmdArgs...)

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		bwerror.ExitWithError(1, "Error creating StdoutPipe for Cmd: %v", err)
	}
	stdoutScanner := bufio.NewScanner(cmdStdout)

	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		bwerror.ExitWithError(1, "Error creating StderrPipe for Cmd: %v", err)
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
		bwerror.ExitWithError(1, "Error starting Cmd: %v", err)
	}

	// https://stackoverflow.com/questions/10385551/get-exit-code-go
	exitCode := 0
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			} else {
				log.Printf(ansi.Ansi(`Warn`, "Could not get exit code for failed program: <ansiCmd>%s"), cmdTitle)
				exitCode = defaultFailedCode
			}
		} else {
			bwerror.ExitWithError(1, "cmd.Wait: %v", err)
		}
	}

	var ansiName, prefix string
	if exitCode == 0 && (vOpt == `all` || vOpt == `allBrief` || vOpt == `ok`) {
		ansiName, prefix = `OK`, `OK`
	}
	if exitCode != 0 && (vOpt == `all` || vOpt == `allBrief` || vOpt == `err`) {
		ansiName, prefix = `Err`, `ERR`
	}
	if len(prefix) > 0 {
		fmt.Println(ansi.Ansi(ansiName, prefix+`: <ansiCmd>`+cmdTitle))
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
