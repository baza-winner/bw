package bwrune_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwtesting"
)

func getTestFileSpec(basename string) string {
	return filepath.Join(os.Getenv("GOPATH"), "src", "github.com/baza-winner/bwcore/bwrune", basename+".test_file")
}

func prepareTestFile(basename string, content []byte) {
	fileSpec := getTestFileSpec(basename)
	err := ioutil.WriteFile(fileSpec, content, 0644)
	if err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	testFiles = append(testFiles, fileSpec)
}

func TestMain(m *testing.M) { // https://stackoverflow.com/questions/23729790/how-can-i-do-test-setup-using-the-testing-package-in-go/34102842#34102842
	mySetupFunction()
	retCode := m.Run()
	myTeardownFunction()
	os.Exit(retCode)
}

var testFiles []string

func mySetupFunction() {
	prepareTestFile("no newline", []byte("Радость"))
	prepareTestFile("newline", []byte("Радость\nСчастье"))
	prepareTestFile("invalid utf8", []byte("Рад\xa0\xa1")) // https://stackoverflow.com/questions/1301402/example-invalid-utf8-string
}

func myTeardownFunction() {
	for _, i := range testFiles {
		os.Remove(i)
	}
}

func getFirstLine(fileSpec string) (result string, err error) {
	var p bwrune.Provider
	p, err = bwrune.ProviderFromFile(fileSpec)
	if err != nil {
		err = bwerr.FromA(bwerr.E{Error: err})
	} else {
		defer p.Close()
		for {
			var currRunePtr *rune
			currRunePtr, err = p.PullRune()
			if currRunePtr == nil || err != nil || *currRunePtr == '\n' {
				break
			} else {
				result += string(*currRunePtr)
			}
		}
	}
	return
}

func TestGetFirstLine(t *testing.T) {
	tests := map[string]bwtesting.Case{
		"no newline": {
			In: []interface{}{getTestFileSpec("no newline")},
			Out: []interface{}{
				"Радость",
				nil,
			},
		},
		"newline": {
			In: []interface{}{getTestFileSpec("newline")},
			Out: []interface{}{
				"Радость",
				nil,
			},
		},
		"invalid utf8": {
			In: []interface{}{getTestFileSpec("invalid utf8")},
			Out: []interface{}{
				"Рад",
				bwerr.From(
					"utf-8 encoding <ansiVal>%#v<ansi> is invalid at pos <ansiPath>%d<ansi> of file <ansiPath>%s",
					5, 2, getTestFileSpec("invalid utf8"),
				),
			},
		},
		"non existent": {
			In: []interface{}{getTestFileSpec("non existent")},
			Out: []interface{}{
				"",
				bwerr.From(
					"open %s: no such file or directory",
					getTestFileSpec("non existent"),
				),
			},
		},
	}
	testsToRun := tests
	bwmap.CropMap(testsToRun)
	// bwmap.CropMap(testsToRun, "[qw/one two three/]")
	bwtesting.BwRunTests(t, getFirstLine, testsToRun)
}
