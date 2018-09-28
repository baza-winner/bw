package core

import (
	"fmt"
	"github.com/baza-winner/bw/ansi"
	// "log"
	"os"
	// "os/exec"
	"reflect"
	// "syscall"
	"encoding/json"
	"testing"
)

func TestMain(m *testing.M) { // https://stackoverflow.com/questions/23729790/how-can-i-do-test-setup-using-the-testing-package-in-go/34102842#34102842
	mySetupFunction()
	retCode := m.Run()
	// myTeardownFunction()
	os.Exit(retCode)
}

func mySetupFunction() {
	ExecCmd(map[string]interface{}{`v`: `err`, `exitOnError`: true}, `go`, `install`, `github.com/baza-winner/bw/bwcoretestinghelper`)
	ExecCmd(map[string]interface{}{`v`: `err`, `exitOnError`: true}, `go`, `install`, `github.com/baza-winner/bw/bwcoretestinghelper2`)
}

func ExampleCamelCaseToKebabCase() {
	fmt.Printf(`%q`, CamelCaseToKebabCase(`someThing`))
	// Output: "some-thing"
}

func ExampleCamelCaseToKebabCase_2() {
	fmt.Printf(`%q`, CamelCaseToKebabCase(`SomeThing`))
	// Output: "some-thing"
}

func ExampleKebabCaseToCamelCase() {
	fmt.Printf(`%q`, KebabCaseToCamelCase(`some-thing`))
	// Output: "someThing"
}

func ExampleKebabCaseToCamelCase_2() {
	fmt.Printf(`%q`, KebabCaseToCamelCase(`SomeThing`))
	// Output: "someThing"
}

func ExampleShortenFileSpec() {
	fmt.Printf(`%q`, ShortenFileSpec(os.Getenv(`HOME`)+`/bw`))
	// Output: "~/bw"
}

func ExampleShortenFileSpec_2() {
	fmt.Printf(`%q`, ShortenFileSpec(`/lib/bw`))
	// Output: "/lib/bw"
}

func ExampleExecCmd() {
	ret := ExecCmd(nil, `bwcoretestinghelper`, `-exit`, `2`, `<stdout>some<stderr>thing`)
	for k, v := range ret {
		fmt.Printf("%s: %v\n", k, v)
	}
	// Unordered ouput:
	// - stdout:[some]
	// - stderr:[thing]
	// - output:[some thing]
	// - exitCode:2
}

func TestExecCmd(t *testing.T) {
	cases := []struct {
		opt     map[string]interface{}
		cmdName string
		cmdArgs []string
		result  map[string]interface{}
	}{
		{
			opt:     nil,
			cmdName: `bwcoretestinghelper2`,
			cmdArgs: []string{`-v`, `none`, `-s`, `all`, `-d`, `-n`, `bwcoretestinghelper`, `-exit`, `2`, `<stdout>some<stderr>thing`},
			result: map[string]interface{}{
				"stdout": []string{
					"===== exitCode: 2",
					"===== stdout:",
					"some",
					"===== stderr:",
					"thing",
				},
				"stderr":   []string{},
				"exitCode": 0,
			},
		},
		{
			opt:     nil,
			cmdName: `bwcoretestinghelper2`,
			cmdArgs: []string{`-v`, `none`, `-s`, `all`, `-d`, `-n`, `-e`, `bwcoretestinghelper`, `-exit`, `2`, `<stdout>some<stderr>thing`},
			result: map[string]interface{}{
				"stdout":   []string{},
				"stderr":   []string{},
				"exitCode": 2,
			},
		},
		{
			opt:     nil,
			cmdName: `bwcoretestinghelper2`,
			cmdArgs: []string{`-v`, `all`, `-d`, `-n`, `bwcoretestinghelper`, `-exit`, `2`, `<stdout>some<stderr>thing`},
			result: map[string]interface{}{
				"stdout": []string{
					"bwcoretestinghelper -exit 2 \u003cstdout\u003esome\u003cstderr\u003ething . . .",
					"some",
					"ERR: bwcoretestinghelper -exit 2 \u003cstdout\u003esome\u003cstderr\u003ething",
					"===== exitCode: 2",
					"===== stdout:",
					"some",
					"===== stderr:",
					"thing",
				},
				"stderr": []string{
					"thing",
				},
				"exitCode": 0,
			},
		},
	}
	for _, c := range cases {
		got := ExecCmd(c.opt, c.cmdName, c.cmdArgs...)
		if !reflect.DeepEqual(got, c.result) { // https://stackoverflow.com/questions/18208394/testing-equivalence-of-maps-golang
			gotJson, _ := json.MarshalIndent(got, "", "  ")
			resultJson, _ := json.MarshalIndent(c.result, "", "  ")
			t.Errorf(ansi.Ansi("", "ExecCmd(%+v, %+v, %+v)\n    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), c.opt, c.cmdName, c.cmdArgs, gotJson, resultJson)
		}
	}
}

func TestGetValidVal(t *testing.T) {
	cases := []struct {
		where    string
		val      interface{}
		def      map[string]interface{}
		whereDef string
		result   interface{}
		err      error
	}{
		{
			where: "somewhere",
			// val:   0,
			val: map[string]interface{}{
				"exitCode": nil,
				"s":        1,
				"v":        "ALL",
				"some":     0,
			},
			def: map[string]interface{}{
				`type`: `map`,
				// `keys`: 0,
				`keys`: map[string]interface{}{
					`v`: map[string]interface{}{
						`type`:    `enum`,
						`enum`:    []string{`all`, `err`, `ok`, `none`},
						`default`: `none`,
					},
					`s`: map[string]interface{}{
						`type`:    `enum`,
						`enum`:    []string{`none`, `stderr`, `stdout`, `all`},
						`default`: `all`,
					},
					`exitOnError`: map[string]interface{}{
						`type`:    `bool`,
						`default`: false,
					},
				},
			},
			whereDef: "somewhere::def",
			result:   nil,
			err:      nil,
		},
	}
	for _, c := range cases {
		got, err := GetValidVal(c.where, c.val, c.def, c.whereDef)
		if err != nil {
			if err != c.err {
				t.Errorf(ansi.Ansi("", "GetValidVal(%s, %+v, %+v, %s)\n    => err:<ansiErr> %v<ansi>\n, want err:<ansiOK>%v"), c.where, c.val, c.def, c.whereDef, err, c.err)
			}
		} else if !reflect.DeepEqual(got, c.result) { // https://stackoverflow.com/questions/18208394/testing-equivalence-of-maps-golang
			gotJson, _ := json.MarshalIndent(got, "", "  ")
			resultJson, _ := json.MarshalIndent(c.result, "", "  ")
			t.Errorf(ansi.Ansi("", "GetValidVal(%s, %+v, %+v, %s)\n    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), c.where, c.val, c.def, c.whereDef, gotJson, resultJson)
		}
	}
}
