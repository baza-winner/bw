package pkgnamegetter

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
)

func getTestFileSpec(basename string) string {
	return filepath.Join(os.Getenv("GOPATH"), "src", "github.com/baza-winner/bwcore/bwsetter/internal/pkgnamegetter", basename+".test_file")
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

func myTeardownFunction() {
	for _, i := range testFiles {
		os.Remove(i)
	}
}

func mySetupFunction() {
	prepareTestFile("package main", []byte("package main"))
	prepareTestFile("package some", []byte("  \n// single line comment \n /* some comment */  \n// single line comment\n package /* infix comment */ some;"))
	prepareTestFile("pack age main", []byte("  pac kage main"))
	prepareTestFile("some", []byte("  \n some"))
	prepareTestFile("package 3", []byte("  \n package 3"))
	prepareTestFile("invalid comment start", []byte("  \n /= some"))
	prepareTestFile("invalid infix comment", []byte("\n\npackage /"))
	// prepareTestFile("invalid comment end", []byte("  \n /* some"))
}

func TestGetPackageNameFromFile(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"package main": {
			In: []interface{}{getTestFileSpec("package main")},
			Out: []interface{}{
				"main",
				nil,
			},
		},
		"package some": {
			In: []interface{}{getTestFileSpec("package some")},
			Out: []interface{}{
				"some",
				nil,
			},
		},
		"pack age main": {
			In: []interface{}{getTestFileSpec("pack age main")},
			Out: []interface{}{
				"",
				bwerr.From(
					"unexpected word <ansiVal>%s<ansi> at line <ansiPath>%d<ansi>, col <ansiPath>%d<ansi> (pos <ansiPath>%d<ansi>) in file <ansiPath>%s",
					"pac", 1, 3, 2, getTestFileSpec("pack age main"),
				),
			},
		},
		"some": {
			In: []interface{}{getTestFileSpec("some")},
			Out: []interface{}{
				"",
				bwerr.From(
					"unexpected rune <ansiVal>%q<ansi> at line <ansiPath>%d<ansi>, col <ansiPath>%d<ansi> (pos <ansiPath>%d<ansi>) in file <ansiPath>%s",
					's', 2, 2, 4, getTestFileSpec("some"),
				),
			},
		},
		"package 3": {
			In: []interface{}{getTestFileSpec("package 3")},
			Out: []interface{}{
				"",
				bwerr.From(
					"unexpected rune <ansiVal>%q<ansi> at line <ansiPath>%d<ansi>, col <ansiPath>%d<ansi> (pos <ansiPath>%d<ansi>) in file <ansiPath>%s",
					'3', 2, 10, 12, getTestFileSpec("package 3"),
				),
			},
		},
		"invalid comment start": {
			In: []interface{}{getTestFileSpec("invalid comment start")},
			Out: []interface{}{
				"",
				bwerr.From(
					"unexpected rune <ansiVal>%q<ansi> at line <ansiPath>%d<ansi>, col <ansiPath>%d<ansi> (pos <ansiPath>%d<ansi>) in file <ansiPath>%s",
					'=', 2, 3, 5, getTestFileSpec("invalid comment start"),
				),
			},
		},
		"invalid infix comment": {
			In: []interface{}{getTestFileSpec("invalid infix comment")},
			Out: []interface{}{
				"",
				bwerr.From(
					"unexpected end of file at line <ansiPath>%d<ansi>, col <ansiPath>%d<ansi> (pos <ansiPath>%d<ansi>) in file <ansiPath>%s",
					3, 9, 10, getTestFileSpec("invalid infix comment"),
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
	bwtesting.BwRunTests(t, GetPackageName, testsToRun)
}
