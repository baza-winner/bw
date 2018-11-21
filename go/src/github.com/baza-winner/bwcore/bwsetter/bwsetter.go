package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwos"
	. "github.com/baza-winner/bwcore/bwsetter/internal/helper"
	"github.com/baza-winner/bwcore/bwsetter/internal/pkgnamegetter"
	"github.com/dave/jennifer/jen"
)

var (
	typeFlag     = flag.String("type", "", "item type name; must be set")
	setFlag      = flag.String("set", "", `Set type name; default is "${type}Set"`)
	nosortFlag   = flag.Bool("nosort", false, `when item has no v[i] < v[j] support, as bool"`)
	nostringFlag = flag.Bool("nostring", false, `when item has no .String() method"`)
	// nodataforjsonFlag = flag.Bool("nodataforjson", false, `when item has no .DataForJSON() method"`)
	omitprefixFlag = flag.Bool("omitprefix", false, `omit prefix for From/FromSlice`)
	testFlag       = flag.Bool("test", false, `generate tests as well`)
)

const (
	bwjsonPackageName = "github.com/baza-winner/bwcore/bwjson"
	jsonPackageName   = "encoding/json"
	setterPackageName = "github.com/baza-winner/bwcore/bwsetter"
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of bwsetter:\n")
	fmt.Fprintf(os.Stderr, "\tbwsetter [flags] -type T [directory]\n")
	fmt.Fprintf(os.Stderr, "For more information, see:\n")
	fmt.Fprintf(os.Stderr, "\thttp://github.com/baza-winner/bwcore\n") // TODO: put link to wiki page
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = Usage
	flag.Parse()
	if len(*typeFlag) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// We accept either one directory
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	// Parse the package once.
	g := Generator{}
	if len(args) == 1 && isDirectory(args[0]) {
		g.parsePackageDir(args[0])
	} else {
		flag.Usage()
		os.Exit(2)
	}

	code := CreateHelper(
		bwjsonPackageName,
		setterPackageName,
		g.packageDir,
		g.packageName,
		g.packagePath,
		*nosortFlag,
		*omitprefixFlag,
		*typeFlag,
		*setFlag,
		*testFlag,
	)

	code.DeclareSet(
		"%s - множество значений типа %s с поддержкой интерфейсов Stringer и MarshalJSON",
		code.IdSet, code.IdItem,
	)

	code.BunchOf(
		"конструктор "+code.IdSet,
		SimpleFunc, code.FromName, "", ReturnSet,
		func(v *Helper, codeRange *jen.Statement) []*jen.Statement {
			return []*jen.Statement{
				jen.Id("result").Op(":=").Id(v.IdSet).Values(),
				codeRange.Block(
					jen.Id("result[k]").Op("=").Struct().Values(),
				),
				jen.Return(jen.Id("result")),
			}
		},
		TestCase{
			In: []interface{}{[]TestItem{A, B}},
			// In:  []interface{}{A, B},
			Out: []interface{}{[]TestItem{A, B}},
		})

	code.SetMethod(
		"создает независимую копию ",
		"Copy", ParamNone, ReturnSet,
		[]*jen.Statement{
			jen.Return(jen.Id(code.FromName + "Set").Call(jen.Id("v"))),
		},
		TestCase{
			In:  []interface{}{[]TestItem{A, B}},
			Out: []interface{}{[]TestItem{A, B}},
		},
	)

	// var testData TestCase

	testData := TestCase{
		In:  []interface{}{[]TestItem{A}},
		Out: []interface{}{[]TestItem{A}},
	}
	code.SetMethod(
		"возвращает в виде []"+code.IdItem,
		"ToSlice", ParamNone, ReturnSlice,
		[]*jen.Statement{
			jen.Id("result").Op(":=").Add(code.Slice).Values(),
			code.RangeSet("v").Block(
				jen.Id("result").Op("=").Append(jen.Id("result"), jen.Id("k")),
			),
			jen.Add(code.Sort),
			jen.Return(jen.Id("result")),
		},
		testData,
	)

	if code.Sort != nil {
		code.Func("",
			"_"+code.IdSet+"ToSliceTestHelper", ParamSlice, ReturnSlice,
			[]*jen.Statement{
				jen.Return(jen.Id(code.FromName + "Slice").Call(jen.Id("kk")).Dot("ToSlice").Call()),
			},
			TestCase{
				In:  []interface{}{[]TestItem{B, A}},
				Out: []interface{}{[]TestItem{A, B}},
			},
		)
	}

	// if !*nodataforjsonFlag {
	// prettyJsonArg := jen.Index().Id(code.IdItem).Values(jen.Id(code.TestItemString(A)))
	// if code.IdItem == "uint8" {
	// 	prettyJsonArg = jen.Index().Id("uint16").Values(jen.Id("uint16").Call(jen.Id(code.TestItemString(A))))
	// }
	code.SetMethod(
		"поддержка интерфейса Stringer",
		"String", ParamNone, ReturnString,
		[]*jen.Statement{
			jen.List(jen.Id("result"), jen.Id("_")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Id("v")),
			jen.Return(jen.Id("string").Call(jen.Id("result"))),
			// jen.Return(jen.Qual(bwjsonPackageName, "Pretty").Call(jen.Id("v"))),
		},
		TestCase{
			In: []interface{}{[]TestItem{A}},
			Out: []interface{}{
				jen.Func().Params().String().Block(
					jen.List(jen.Id("result"), jen.Id("_")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Id(code.TestItemString(A))),
					jen.Return(jen.Lit("[").Op("+").String().Call(jen.Id("result")).Op("+").Lit("]")),
				).Call(),
			},
		},
	)

	// func (v ValKindSet) String() string {
	// 	result, _ := json.Marshal(v)
	// 	return string(result)
	// }
	// func() { result, _ := json.Marshal(v); return "[" + string(result) + "]"}

	code.SetMethod(
		"поддержка интерфейса MarshalJSON",
		"MarshalJSON", ParamNone, ReturnJSON,
		[]*jen.Statement{
			jen.Id("result").Op(":=").Index().Interface().Values(),
			// code.RangeSet("v").Block(
			code.RangeSlice(jen.Id("v").Dot("ToSlice").Call()).Block(
				jen.Id("result").Op("=").Append(jen.Id("result"), jen.Id("k")),
			),
			jen.Return(jen.Qual(jsonPackageName, "Marshal").Call(jen.Id("result"))),
		},
		TestCase{
			In: []interface{}{[]TestItem{A}},
			Out: []interface{}{
				jen.Parens(
					jen.Func().Params().Index().Byte().Block(
						jen.List(
							jen.Id("result"),
							jen.Id("_"),
						).Op(":=").Qual(jsonPackageName, "Marshal").Call(
							jen.Index().Interface().Values(
								jen.Id(code.TestItemString(A)),
							),
						),
						jen.Return(jen.Id("result")),
					),
				).Call(),
				jen.Nil(),
			},
		},
	)

	if !*nostringFlag {
		code.SetMethod(
			"возвращает []string строковых представлений элементов множества",
			"ToSliceOfStrings", ParamNone, ReturnSliceOfStrings,
			[]*jen.Statement{
				jen.Id("result").Op(":=").Index().String().Values(),
				code.RangeSet("v").Block(
					jen.Id("result").Op("=").Append(jen.Id("result"), code.ToString("k")),
				),
				jen.Qual("sort", "Strings").Call(jen.Id("result")),
				jen.Return(jen.Id("result")),
			},
			TestCase{
				In: []interface{}{[]TestItem{A}},
				Out: []interface{}{
					jen.Index().String().Values(code.ToString(code.TestItemString(A))),
				},
			},
		)
	}

	code.SetMethod(
		"возвращает true, если множество содержит заданный элемент, в противном случае - false",
		"Has", ParamArg, ReturnBool,
		[]*jen.Statement{
			jen.List(jen.Id("_"), jen.Id("ok")).Op(":=").Id("v[k]"),
			jen.Return(jen.Id("ok")),
		},
		TestCases{
			"true": TestCase{
				In:  []interface{}{[]TestItem{A}, A},
				Out: []interface{}{true},
			},
			"false": TestCase{
				In:  []interface{}{[]TestItem{A}, B},
				Out: []interface{}{false},
			},
		},
	)

	code.BunchOf(
		"возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.\nHasAny(<пустой набор/множесто>) возвращает false",
		SetMethod, "HasAny", "Of", ReturnBool,
		func(v *Helper, codeRange *jen.Statement) []*jen.Statement {
			return []*jen.Statement{
				codeRange.Block(
					jen.If(
						jen.List(jen.Id("_"), jen.Id("ok")).Op(":=").Id("v[k]"),
						jen.Id("ok"),
					).Block(
						jen.Return().True(),
					),
				),
				jen.Return().False(),
			}
		}, TestCases{
			"true": TestCase{
				In:  []interface{}{[]TestItem{A}, []TestItem{A, B}},
				Out: []interface{}{true},
			},
			"false": TestCase{
				In:  []interface{}{[]TestItem{A}, []TestItem{B}},
				Out: []interface{}{false},
			},
			"empty": TestCase{
				In:  []interface{}{[]TestItem{A}, []TestItem{}},
				Out: []interface{}{false},
			},
		})

	code.BunchOf(
		"возвращает true, если множество содержит все заданные элементы, в противном случае - false.\nHasEach(<пустой набор/множесто>) возвращает true",
		SetMethod, "HasEach", "Of", ReturnBool,
		func(v *Helper, codeRange *jen.Statement) []*jen.Statement {
			return []*jen.Statement{
				codeRange.Block(
					jen.If(
						jen.List(jen.Id("_"), jen.Id("ok")).Op(":=").Id("v[k]"),
						jen.Op("!").Id("ok"),
					).Block(
						jen.Return().False(),
					),
				),
				jen.Return().True(),
			}
		}, TestCases{
			"true": TestCase{
				In:  []interface{}{[]TestItem{A, B}, []TestItem{A, B}},
				Out: []interface{}{true},
			},
			"false": TestCase{
				In:  []interface{}{[]TestItem{A}, []TestItem{A, B}},
				Out: []interface{}{false},
			},
			"empty": TestCase{
				In:  []interface{}{[]TestItem{A}, []TestItem{}},
				Out: []interface{}{true},
			},
		})

	code.BunchOf(
		"добавляет элементы в множество v",
		SetMethod, "Add", "", ReturnNone,
		func(v *Helper, codeRange *jen.Statement) []*jen.Statement {
			return []*jen.Statement{
				codeRange.Block(
					jen.Id("v[k]").Op("=").Struct().Values(),
				),
			}
		}, TestCase{
			In:  []interface{}{[]TestItem{A}, []TestItem{B}},
			Out: []interface{}{[]TestItem{A, B}},
		})

	code.BunchOf(
		"удаляет элементы из множествa v",
		SetMethod, "Del", "", ReturnNone,
		func(v *Helper, codeRange *jen.Statement) []*jen.Statement {
			return []*jen.Statement{
				codeRange.Block(
					jen.Delete(jen.Id("v"), jen.Id("k")),
				),
			}
		}, TestCase{
			In:  []interface{}{[]TestItem{A, B}, []TestItem{B}},
			Out: []interface{}{[]TestItem{A}},
		})

	code.SetMethod(
		"возвращает результат объединения двух множеств. Исходные множества остаются без изменений",
		"Union", ParamSet, ReturnSet,
		[]*jen.Statement{
			jen.Id("result").Op(":=").Id("v").Dot("Copy").Call(),
			jen.Id("result").Dot("AddSet").Call(jen.Id("s")),
			jen.Return(jen.Id("result")),
		},
		TestCase{
			In:  []interface{}{[]TestItem{A}, []TestItem{B}},
			Out: []interface{}{[]TestItem{A, B}},
		},
	)

	code.SetMethod(
		"возвращает результат пересечения двух множеств. Исходные множества остаются без изменений",
		"Intersect", ParamSet, ReturnSet,
		[]*jen.Statement{
			jen.Id("result").Op(":=").Id(code.IdSet).Values(),
			code.RangeSet("v",
				jen.If(
					jen.List(jen.Id("_"), jen.Id("ok")).Op(":=").Id("s[k]"),
					jen.Id("ok"),
				).Block(
					jen.Id("result[k]").Op("=").Struct().Values(),
				)),
			jen.Return(jen.Id("result")),
		},
		TestCase{
			In:  []interface{}{[]TestItem{A, B}, []TestItem{B}},
			Out: []interface{}{[]TestItem{B}},
		},
	)

	code.SetMethod(
		"возвращает результат вычитания двух множеств. Исходные множества остаются без изменений",
		"Subtract", ParamSet, ReturnSet,
		[]*jen.Statement{
			jen.Id("result").Op(":=").Id(code.IdSet).Values(),
			code.RangeSet("v",
				jen.If(
					jen.List(jen.Id("_"), jen.Id("ok")).Op(":=").Id("s[k]"),
					jen.Op("!").Id("ok"),
				).Block(
					jen.Id("result[k]").Op("=").Struct().Values(),
				)),
			jen.Return(jen.Id("result")),
		},
		TestCase{
			In:  []interface{}{[]TestItem{A, B}, []TestItem{B}},
			Out: []interface{}{[]TestItem{A}},
		},
	)

	if code.Sort != nil {
		code.DeclareSlice()
		code.SliceMethod("", "Len", ParamNone, ReturnInt,
			[]*jen.Statement{
				jen.Return().Len(jen.Id("v")),
			},
			nil,
		)

		code.SliceMethod("", "Swap", ParamIJ, ReturnNone,
			[]*jen.Statement{
				jen.List(jen.Id("v[i]"), jen.Id("v[j]")).Op("=").List(jen.Id("v[j]"), jen.Id("v[i]")),
			},
			nil,
		)

		code.SliceMethod("", "Less", ParamIJ, ReturnBool,
			[]*jen.Statement{
				jen.Return().Id("v[i]").Op("<").Id("v[j]"),
			},
			nil,
		)
	}

	code.Save()
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	packageDir  string
	packageName string
	packagePath string
}

var goFileRegexp = regexp.MustCompile(`^.+\.go$`)
var goTestFileRegexp = regexp.MustCompile(`^.+_test\.go$`)

// parsePackageDir parses the package residing in the directory.
func (g *Generator) parsePackageDir(directory string) {
	var err error
	g.packageDir, err = filepath.Abs(directory)
	if err != nil {
		bwerr.PanicErr(err)
	}
	var files []os.FileInfo
	files, err = ioutil.ReadDir(g.packageDir)
	if err != nil {
		bwerr.PanicErr(err)
	}
	for _, f := range files {
		if !f.Mode().IsRegular() || !goFileRegexp.Match([]byte(f.Name())) || goTestFileRegexp.Match([]byte(f.Name())) {
			continue
		}
		if packageName, err := pkgnamegetter.GetPackageName(filepath.Join(g.packageDir, f.Name())); err != nil {
			bwerr.PanicErr(err)
		} else if len(g.packageName) == 0 {
			g.packageName = packageName
		} else if packageName != g.packageName {
			bwerr.Panic("packageName: %s or %s?", packageName, g.packageName)
		}
	}
	g.packagePath = getPackagePath(g.packageDir)
}

func getPackagePath(packageDir string) string {
	var err error
	gopath := os.Getenv("GOPATH")
	if len(gopath) == 0 {
		if gopath, err = filepath.Abs(os.Getenv("HOME")); err != nil {
			bwerr.PanicErr(err)
		}
		gopath += "/.go"
	} else {
		if gopath, err = filepath.Abs(gopath); err != nil {
			bwerr.PanicErr(err)
		}
	}
	srcSuffix := "/src/"
	gopathsrc := gopath + srcSuffix

	if len(packageDir) < len(gopathsrc) || packageDir[0:len(gopathsrc)] != gopathsrc {
		bwos.Exit(1,
			"<ansiVar>pwd<ansi> (<ansiPath>%s<ansi>) not in <ansiPath>$GOPATH%s<ansi> (<ansiVar>$GOPATH<ansi>: <ansiPath>%s<ansi>)",
			packageDir, srcSuffix, gopath,
		)
	}
	return packageDir[len(gopathsrc):]
}
