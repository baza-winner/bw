package bwexec_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/baza-winner/bwcore/bwexec"
	"github.com/baza-winner/bwcore/bwtesting"
)

func TestMain(m *testing.M) { // https://stackoverflow.com/questions/23729790/how-can-i-do-test-setup-using-the-testing-package-in-go/34102842#34102842
	mySetupFunction()
	retCode := m.Run()
	// myTeardownFunction()
	os.Exit(retCode)
}

// var installOpt = defparse.MustParseMap("{v: 'err', exitOnError: true, s: 'none'}")
// var installOpt map[string]interface{}

// func init() {
// 	installOpt = bwval.From(bwval.S{S: `
// 		{
// 			v "err"
// 			exitOnError true
// 			s "none"
// 		}
// 		`}).MustMap()
// }

func mySetupFunction() {
	bwexec.Cmd(bwexec.Args(`go`, `install`, `github.com/baza-winner/bwcore/bwexec/bwexectesthelper`))
	bwexec.Cmd(bwexec.Args(`go`, `install`, `github.com/baza-winner/bwcore/bwexec/bwexectesthelper2`))
}

func ExampleCmd() {
	ret := bwexec.Cmd(bwexec.Args(`bwexectesthelper`, `-exit`, `2`, `<stdout>some<stderr>thing`))
	for k, v := range ret {
		fmt.Printf("%s: %v\n", k, v)
	}
	// Unordered ouput:
	// - stdout:[some]
	// - stderr:[thing]
	// - output:[some thing]
	// - exitCode:2
}

func TestCmd(t *testing.T) {
	bwtesting.BwRunTests(t,
		bwexec.Cmd,
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{
				"test1": {
					In: []interface{}{
						bwexec.Args(
							`bwexectesthelper2`,
							`-v`, `none`, `-s`, `all`, `-d`, `-n`, `bwexectesthelper`, `-exit`, `2`, `<stdout>some<stderr>thing`,
						),
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
						bwexec.Args(
							`bwexectesthelper2`,
							`-v`, `none`, `-s`, `all`, `-d`, `-n`, `-e`, `bwexectesthelper`, `-exit`, `2`, `<stdout>some<stderr>thing`,
						),
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
						bwexec.Args(
							`bwexectesthelper2`,
							`-v`, `all`, `-d`, `-n`, `bwexectesthelper`, `-exit`, `2`, `<stdout>some<stderr>thing`,
						),
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
			return tests
		}(),
		// "test1",
	)
}
