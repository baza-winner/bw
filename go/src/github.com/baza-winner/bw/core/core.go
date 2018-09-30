// Предоставляет базовые функции
package core

import (
	"bufio"
	// "encoding/json"
	"errors"
	"fmt"
	"github.com/baza-winner/bw/ansi"
	"github.com/iancoleman/strcase"
	"github.com/jimlawless/whereami"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func CamelCaseToKebabCase(s string) string {
	return strcase.ToKebab(s)
}

func KebabCaseToCamelCase(s string) string {
	return strcase.ToLowerCamel(s)
}

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

// func GetValidMap(where string, val interface{}, def map[string]interface{}) (result map[string]interface{}) {
// 	var validVal interface{}
// 	var err error
// 	if validVal, err = GetValidVal(where, val, def, where+`::def`); err != nil {
// 		defJson, _ := json.MarshalIndent(def, "", "  ")
// 		log.Panicf(ansi.Ansi(`Err`, `ERR: <ansiOutline>%s<ansi> value ( <ansiPrimaryLiteral>%+v<ansi> ) does not fit def <ansiSecondaryLiteral>%+v<ansi>: %v`), where, val, defJson, err)
// 	}
// 	var ok bool
// 	if result, ok = validVal.(map[string]interface{}); !ok {
// 		log.Panicf(ansi.Ansi(`Err`, `ERR: <ansiOutline>%s<ansi> value ( <ansiPrimaryLiteral>%+v<ansi> ) is not <ansiSecondaryLiteral>map<ansi>: %v`), where, val, err)
// 	}
// 	return
// }

func GetValidVal(whereVal string, val interface{}, def map[string]interface{}, whereDef string) (result interface{}, err error) {
	var defType string
	var ok bool
	if defType, err = GetStringKey(whereDef, def, `type`); err == nil {
		if defType == `map` {
			var valMap map[string]interface{}
			if valMap, ok = val.(map[string]interface{}); ok {
				var defKeys map[string]interface{}
				if defKeys, err = GetMapKeyIfExists(whereDef, def, `keys`); defKeys != nil && err == nil {
					for key := range valMap {
						// if keyVal, ok := defKeys[key]; !ok {
						if _, ok := defKeys[key]; !ok {
							err = Error(`<ansiOutline>%s<ansi> (<ansiSecondaryLiteral>%v<ansi>) has unexpected key <ansiPrimaryLiteral>%s`, whereVal, val, key)
							return
						}
					}
				}
			} else {
				err = Error(`<ansiOutline>%s<ansi> (<ansiSecondaryLiteral>%v<ansi>) is not of type <ansiPrimaryLiteral>%s`, whereVal, val, `map`)
			}
		} else {
			err = Error(`<ansiOutline>%s<ansi>[<ansiSecondaryLiteral>%s<ansi>] has non supported value <ansiPrimaryLiteral>%s`, whereDef, `type`, defType)
		}
	}
	return val, err
}

func GetMapKeyIfExists(where string, m map[string]interface{}, keyName string) (result map[string]interface{}, err error) {
	if m != nil {
		if val, ok := m[keyName]; ok {
			if typedVal, ok := val.(map[string]interface{}); ok {
				result = typedVal
			} else {
				err = Error(`<ansiOutline>%s<ansi>[<ansiSecondaryLiteral>%s<ansi>] (<ansiSecondaryLiteral>%+v<ansi>) is not <ansiPrimaryLiteral>%s`, where, keyName, val, `map`)
			}
		} else {
			result = nil
			err = nil
		}
	} else {
		err = Error(`<ansiOutline>%s<ansi> is not <ansiPrimaryLiteral>map`, where)
	}
	return
}

func GetStringKey(where string, m map[string]interface{}, keyName string) (result string, err error) {
	if m != nil {
		if val, ok := m[keyName]; ok {
			if typedVal, ok := val.(string); ok {
				result = typedVal
			} else {
				err = Error(`<ansiOutline>%s<ansi>[<ansiSecondaryLiteral>%s<ansi>] (<ansiSecondaryLiteral>%+v<ansi>) is not <ansiPrimaryLiteral>%s`, where, keyName, val, `string`)
			}
		} else {
			err = Error(`<ansiOutline>%s<ansi> has not key <ansiPrimaryLiteral>%s`, where, keyName)
		}
	} else {
		err = Error(`<ansiOutline>%s<ansi> is not <ansiPrimaryLiteral>map`, where)
	}
	return
}

const defaultFailedCode = 1

func ExecCmd(opt map[string]interface{}, cmdName string, cmdArgs ...string) (result map[string]interface{}) {
	// _ = GetValidMap(`ExecCmd.opt`, opt, map[string]interface{}{
	// 	`type`: `map`,
	// 	`keys`: map[string]interface{}{
	// 		`v`: map[string]interface{}{
	// 			`type`:    `enum`,
	// 			`enum`:    []string{`all`, `err`, `ok`, `none`},
	// 			`default`: `none`,
	// 		},
	// 		`s`: map[string]interface{}{
	// 			`type`:    `enum`,
	// 			`enum`:    []string{`none`, `stderr`, `stdout`, `all`},
	// 			`default`: `all`,
	// 		},
	// 		`exitOnError`: map[string]interface{}{
	// 			`type`:    `bool`,
	// 			`default`: false,
	// 		},
	// 	},
	// })

	result = map[string]interface{}{}

	cmdTitle := SmartQuote(append([]string{cmdName}, cmdArgs...)...)
	sOpt := GetStringKeyOrDefault(opt, `s`, `all`)
	vOpt := GetStringKeyOrDefault(opt, `v`, `none`)
	if vOpt == `all` || vOpt == `allBrief` {
		fmt.Println(ansi.Ansi(``, `<ansiCmd>`+cmdTitle+`<ansi> . . .`))
	}
	cmd := exec.Command(cmdName, cmdArgs...)

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		ExitWithError(1, "Error creating StdoutPipe for Cmd: %v", err)
	}
	stdoutScanner := bufio.NewScanner(cmdStdout)

	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		ExitWithError(1, "Error creating StderrPipe for Cmd: %v", err)
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
		ExitWithError(1, "Error starting Cmd: %v", err)
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
			ExitWithError(1, "cmd.Wait: %v", err)
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
	if GetBoolKeyOrDefault(opt, `exitOnError`, false) && exitCode != 0 {
		os.Exit(exitCode)
	}
	result[`exitCode`] = exitCode
	result[`stdout`] = stdout
	result[`stderr`] = stderr
	// result[`output`] = output
	return result
}

// func GetBoolKey(m map[string]interface{}, keyName string, defaultValue bool) (result bool) {
//   result = defaultValue
//   if m != nil {
//     if val, ok := m[keyName]; ok {
//       if typedVal, ok := val.(bool); ok {
//         result = typedVal
//       }
//     }
//   }
//   return
// }

func GetBoolKeyOrDefault(m map[string]interface{}, keyName string, defaultValue bool) (result bool) {
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

func GetStringKeyOrDefault(m map[string]interface{}, keyName string, defaultValue string) (result string) {
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

func ExitWithError(exitCode int, fmtString string, fmtArgs ...interface{}) {
	log.Print(ansi.Ansi(`Err`, fmt.Sprintf(fmtString, fmtArgs...)))
	os.Exit(exitCode)
}

func Error(msgFmt string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(ansi.Ansi(`Err`, msgFmt), args...))
}

func Panic(msgFmt string, args ...interface{}) {
	log.Panicf(ansi.Ansi(`Err`, msgFmt+` <ansiCmd>at `+whereami.WhereAmI(2)), args...)
}

func Panicd(depth uint, msgFmt string, args ...interface{}) {
	log.Panicf(ansi.Ansi(`Err`, msgFmt+` <ansiCmd>at `+whereami.WhereAmI(int(depth)+2)), args...)
}

// // https://lawlessguy.wordpress.com/2016/04/17/display-file-function-and-line-number-in-go-golang/
// func WhereAmI(depthList ...int) string {
// 	var depth int
// 	if depthList == nil {
// 		depth = 1
// 	} else {
// 		depth = depthList[0]
// 	}
// 	function, file, line, _ := runtime.Caller(depth)
// 	return fmt.Sprintf("File: %s  Function: %s Line: %d", chopPath(file), runtime.FuncForPC(function).Name(), line)
// }
