package bwexec

import (
	"fmt"
	"os"
	"testing"

	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/defparse"
)

func TestMain(m *testing.M) { // https://stackoverflow.com/questions/23729790/how-can-i-do-test-setup-using-the-testing-package-in-go/34102842#34102842
	mySetupFunction()
	retCode := m.Run()
	// myTeardownFunction()
	os.Exit(retCode)
}

var installOpt = defparse.MustParseMap("{v: 'err', exitOnError: true, s: 'none'}")

func mySetupFunction() {
	ExecCmd(installOpt, `go`, `install`, `github.com/baza-winner/bwcore/bwexec/bwexectesthelper`)
	ExecCmd(installOpt, `go`, `install`, `github.com/baza-winner/bwcore/bwexec/bwexectesthelper2`)
}

func ExampleExecCmd() {
	ret := ExecCmd(nil, `bwexectesthelper`, `-exit`, `2`, `<stdout>some<stderr>thing`)
	for k, v := range ret {
		fmt.Printf("%s: %v\n", k, v)
	}
	// Unordered ouput:
	// - stdout:[some]
	// - stderr:[thing]
	// - output:[some thing]
	// - exitCode:2
}

// type testExecCmdStruct struct {
// 	opt     map[string]interface{}
// 	cmdName string
// 	cmdArgs []string
// 	result  map[string]interface{}
// }

func TestExecCmd(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"test1": {
			In: []interface{}{
				map[string]interface{}(nil),
				`bwexectesthelper2`,
				[]string{`-v`, `none`, `-s`, `all`, `-d`, `-n`, `bwexectesthelper`, `-exit`, `2`, `<stdout>some<stderr>thing`},
			},
			Out: []interface{}{
				map[string]interface{}{
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
		},
		"test2": {
			In: []interface{}{
				map[string]interface{}(nil),
				`bwexectesthelper2`,
				[]string{`-v`, `none`, `-s`, `all`, `-d`, `-n`, `-e`, `bwexectesthelper`, `-exit`, `2`, `<stdout>some<stderr>thing`},
			},
			Out: []interface{}{
				map[string]interface{}{
					"stdout":   []string{},
					"stderr":   []string{},
					"exitCode": 2,
				},
			},
		},
		"test3": {
			In: []interface{}{
				map[string]interface{}(nil),
				`bwexectesthelper2`,
				[]string{`-v`, `all`, `-d`, `-n`, `bwexectesthelper`, `-exit`, `2`, `<stdout>some<stderr>thing`},
			},
			Out: []interface{}{
				map[string]interface{}{
					"stdout": []string{
						"bwexectesthelper -exit 2 \u003cstdout\u003esome\u003cstderr\u003ething . . .",
						"some",
						"ERR: bwexectesthelper -exit 2 \u003cstdout\u003esome\u003cstderr\u003ething",
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
		},
	}

	testsToRun := tests
	bwmap.CropMap(testsToRun)
	// bwmap.CropMap(testsToRun, "[qw/one two three/]")
	bwtesting.BwRunTests(t, ExecCmd, testsToRun)
}
