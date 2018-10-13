package main

import (
	"flag"

	. "github.com/baza-winner/bwcore/setter/internal/helper"
	. "github.com/dave/jennifer/jen"
)

func main() {
	typeFlag := flag.String("type", "", "item type name; must be set")
	mainFlag := flag.Bool("main", false, "for `package main`")
	setFlag := flag.String("set", "", `Set type name; default is "${type}Set"`)
	nosortFlag := flag.Bool("nosort", false, `when item has no v[i] < v[j] support, as bool"`)
	omitprefixFlag := flag.Bool("omitprefix", false, `omit prefix for From/FromSlice`)
	testFlag := flag.Bool("test", false, `generate tests as well`)
	flag.Parse()

	bwjsonPackageName := "github.com/baza-winner/bwcore/bwjson"
	code := CreateHelper(
		`github.com/baza-winner/bwcore/setter`,
		*mainFlag,
		*nosortFlag,
		*omitprefixFlag,
		*typeFlag,
		*setFlag,
		*testFlag,
	)

	code.DeclareSet(
		"%s - множество значений типа %s с поддержкой интерфейсов Stringer и %s.Jsonable",
		code.IdSet, code.IdItem, bwjsonPackageName,
	)

	code.BunchOf(
		"конструктор "+code.IdSet,
		SimpleFunc, code.FromName, "", ReturnSet,
		func(v *Helper, codeRange *Statement) []*Statement {
			return []*Statement{
				Id("result").Op(":=").Id(v.IdSet).Values(),
				codeRange.Block(
					Id("result[k]").Op("=").Struct().Values(),
				),
				Return(Id("result")),
			}
		},
		TestCase{
			In:  []interface{}{[]TestItem{A, B}},
			Out: []interface{}{[]TestItem{A, B}},
		})

	code.SetMethod(
		"создает независимую копию ",
		"Copy", ParamNone, ReturnSet,
		[]*Statement{
			Return(Id(code.FromName + "Set").Call(Id("v"))),
		},
		TestCase{
			In:  []interface{}{[]TestItem{A, B}},
			Out: []interface{}{[]TestItem{A, B}},
		},
	)

	var testData TestCase

	testData = TestCase{
		In:  []interface{}{[]TestItem{A}},
		Out: []interface{}{[]TestItem{A}},
	}
	code.SetMethod(
		"возвращает в виде []"+code.IdItem,
		"ToSlice", ParamNone, ReturnSlice,
		[]*Statement{
			Id("result").Op(":=").Add(code.Slice).Values(),
			code.RangeSet("v").Block(
				Id("result").Op("=").Append(Id("result"), Id("k")),
			),
			Add(code.Sort),
			Return(Id("result")),
		},
		testData,
	)

	if code.Sort != nil {
		code.Func("",
			"_"+code.IdSet+"ToSliceTestHelper", ParamSlice, ReturnSlice,
			[]*Statement{
				Return(Id(code.FromName + "Slice").Call(Id("kk")).Dot("ToSlice").Call()),
			},
			TestCase{
				In:  []interface{}{[]TestItem{B, A}},
				Out: []interface{}{[]TestItem{A, B}},
			},
		)
	}

	code.SetMethod(
		"поддержка интерфейса Stringer",
		"String", ParamNone, ReturnString,
		[]*Statement{
			Return(Qual(bwjsonPackageName, "PrettyJsonOf").Call(Id("v"))),
		},
		TestCase{
			In: []interface{}{[]TestItem{A}},
			Out: []interface{}{
				Qual("fmt", "Sprintf").Call(Lit("[\n  %q\n]"), code.ToString(code.TestItemString(A))),
			},
		},
	)

	code.SetMethod(
		"поддержка интерфейса bwjson.Jsonable",
		"GetDataForJson", ParamNone, ReturnInterface,
		[]*Statement{
			Id("result").Op(":=").Index().Interface().Values(),
			code.RangeSet("v").Block(
				Id("result").Op("=").Append(Id("result"), code.ToString("k")),
			),
			Return(Id("result")),
		},
		TestCase{
			In: []interface{}{[]TestItem{A}},
			Out: []interface{}{
				Index().Interface().Values(code.ToString(code.TestItemString(A))),
			},
		},
	)

	code.SetMethod(
		"возвращает []string строковых представлений элементов множества",
		"ToSliceOfStrings", ParamNone, ReturnSliceOfStrings,
		[]*Statement{
			Id("result").Op(":=").Index().String().Values(),
			code.RangeSet("v").Block(
				Id("result").Op("=").Append(Id("result"), code.ToString("k")),
			),
			Qual("sort", "Strings").Call(Id("result")),
			Return(Id("result")),
		},
		TestCase{
			In: []interface{}{[]TestItem{A}},
			Out: []interface{}{
				Index().String().Values(code.ToString(code.TestItemString(A))),
			},
		},
	)

	code.SetMethod(
		"возвращает true, если множество содержит заданный элемент, в противном случае - false",
		"Has", ParamArg, ReturnBool,
		[]*Statement{
			List(Id("_"), Id("ok")).Op(":=").Id("v[k]"),
			Return(Id("ok")),
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
		func(v *Helper, codeRange *Statement) []*Statement {
			return []*Statement{
				codeRange.Block(
					If(
						List(Id("_"), Id("ok")).Op(":=").Id("v[k]"),
						Id("ok"),
					).Block(
						Return().True(),
					),
				),
				Return().False(),
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
		func(v *Helper, codeRange *Statement) []*Statement {
			return []*Statement{
				codeRange.Block(
					If(
						List(Id("_"), Id("ok")).Op(":=").Id("v[k]"),
						Op("!").Id("ok"),
					).Block(
						Return().False(),
					),
				),
				Return().True(),
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
		func(v *Helper, codeRange *Statement) []*Statement {
			return []*Statement{
				codeRange.Block(
					Id("v[k]").Op("=").Struct().Values(),
				),
			}
		}, TestCase{
			In:  []interface{}{[]TestItem{A}, []TestItem{B}},
			Out: []interface{}{[]TestItem{A, B}},
		})

	code.BunchOf(
		"удаляет элементы из множествa v",
		SetMethod, "Del", "", ReturnNone,
		func(v *Helper, codeRange *Statement) []*Statement {
			return []*Statement{
				codeRange.Block(
					Delete(Id("v"), Id("k")),
				),
			}
		}, TestCase{
			In:  []interface{}{[]TestItem{A, B}, []TestItem{B}},
			Out: []interface{}{[]TestItem{A}},
		})

	code.SetMethod(
		"возвращает результат объединения двух множеств. Исходные множества остаются без изменений",
		"Union", ParamSet, ReturnSet,
		[]*Statement{
			Id("result").Op(":=").Id("v").Dot("Copy").Call(),
			Id("result").Dot("AddSet").Call(Id("s")),
			Return(Id("result")),
		},
		TestCase{
			In:  []interface{}{[]TestItem{A}, []TestItem{B}},
			Out: []interface{}{[]TestItem{A, B}},
		},
	)

	code.SetMethod(
		"возвращает результат пересечения двух множеств. Исходные множества остаются без изменений",
		"Intersect", ParamSet, ReturnSet,
		[]*Statement{
			Id("result").Op(":=").Id(code.IdSet).Values(),
			code.RangeSet("v",
				If(
					List(Id("_"), Id("ok")).Op(":=").Id("s[k]"),
					Id("ok"),
				).Block(
					Id("result[k]").Op("=").Struct().Values(),
				)),
			Return(Id("result")),
		},
		TestCase{
			In:  []interface{}{[]TestItem{A, B}, []TestItem{B}},
			Out: []interface{}{[]TestItem{B}},
		},
	)

	code.SetMethod(
		"возвращает результат вычитания двух множеств. Исходные множества остаются без изменений",
		"Subtract", ParamSet, ReturnSet,
		[]*Statement{
			Id("result").Op(":=").Id(code.IdSet).Values(),
			code.RangeSet("v",
				If(
					List(Id("_"), Id("ok")).Op(":=").Id("s[k]"),
					Op("!").Id("ok"),
				).Block(
					Id("result[k]").Op("=").Struct().Values(),
				)),
			Return(Id("result")),
		},
		TestCase{
			In:  []interface{}{[]TestItem{A, B}, []TestItem{B}},
			Out: []interface{}{[]TestItem{A}},
		},
	)

	if code.Sort != nil {
		code.DeclareSlice()
		code.SliceMethod("", "Len", ParamNone, ReturnInt,
			[]*Statement{
				Return().Len(Id("v")),
			},
			nil,
		)

		code.SliceMethod("", "Swap", ParamIJ, ReturnNone,
			[]*Statement{
				List(Id("v[i]"), Id("v[j]")).Op("=").List(Id("v[j]"), Id("v[i]")),
			},
			nil,
		)

		code.SliceMethod("", "Less", ParamIJ, ReturnBool,
			[]*Statement{
				Return().Id("v[i]").Op("<").Id("v[j]"),
			},
			nil,
		)
	}

	code.Save()
}
