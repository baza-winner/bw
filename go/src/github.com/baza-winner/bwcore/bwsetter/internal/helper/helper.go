package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/baza-winner/bwcore/bwerr"

	. "github.com/dave/jennifer/jen"
)

type ParamType uint8

const (
	ParamBelow ParamType = iota
	ParamNone
	ParamArg
	ParamArgs
	ParamSlice
	ParamSet
	ParamIJ
	ParamAbove
)

//go:generate stringer -type=ParamType

type ReturnType uint8

const (
	ReturnBelow ReturnType = iota
	ReturnNone
	ReturnBool
	ReturnInt
	ReturnInterface
	ReturnJSON
	ReturnString
	ReturnSet
	ReturnSlice
	ReturnSliceOfStrings
	ReturnAbove
)

//go:generate stringer -type=ReturnType

type Helper struct {
	bwjsonPackageName string
	packageDir        string
	packagePath       string
	f                 *File
	t                 *File
	IdItem            string
	IdSet             string
	idSlice           string
	FromName          string
	Sort              *Statement
	Slice             *Statement
	tests             []Code
	needTests         bool
}

func CreateHelper(
	bwjsonPackageName string,
	generatorName string,
	packageDir string,
	packageName string,
	packagePath string,
	noSort bool,
	omitPrefix bool,
	IdItem string,
	IdSet string,
	needTests bool,
) *Helper {
	var codeSort, slice *Statement
	var idSlice string
	if noSort {
		slice = Index().Id(IdItem)
	} else {
		idSlice = "_" + IdItem + "Slice"
		slice = Id(idSlice)
		codeSort = Qual("sort", "Sort").Call(Id("result"))
	}
	if len(IdSet) == 0 {
		IdSet = IdItem + "Set"
	}
	fromName := "From"
	if !omitPrefix {
		fromName = IdSet + fromName
	}
	code := Helper{
		bwjsonPackageName,
		packageDir,
		packagePath,
		NewFilePathName(packagePath, packageName),
		NewFilePathName(packagePath, packageName),
		IdItem,
		IdSet,
		idSlice,
		fromName,
		codeSort,
		slice,
		[]Code{},
		needTests,
	}
	code.codeGenerated(code.f, generatorName)
	code.codeGenerated(code.t, generatorName)

	return &code
}

func (v *Helper) DeclareSet(fmtDescription string, fmtArgs ...interface{}) {
	v.f.Commentf(fmtDescription, fmtArgs...)
	v.f.Type().Id(v.IdSet).Map(Qual(v.packagePath, v.IdItem)).Struct()
}

func (v *Helper) DeclareSlice() {
	v.f.Type().Id(v.idSlice).Index().Id(v.IdItem)
}

func (v *Helper) codeGenerated(f *File, generatorName string) {
	f.HeaderComment(fmt.Sprintf(`Code generated by "bwsetter -type=`+v.IdItem+`"; DO NOT EDIT; bwsetter: go get %s`, strings.Join(os.Args[1:], " "), generatorName))
}

func (v *Helper) RangeSlice(
	id string,
	blockStatements ...Code,
) *Statement {
	result := For().
		List(Id("_"), Id("k")).
		Op(":=").
		Range().Id(id)
	if blockStatements != nil {
		result = result.Block(blockStatements...)
	}
	return result
}

func (v *Helper) RangeSet(
	id string,
	blockStatements ...Code,
) *Statement {
	result := For().
		List(Id("k"), Id("_")).
		Op(":=").
		Range().Id(id)
	if blockStatements != nil {
		result = result.Block(blockStatements...)
	}
	return result
}

func (v *Helper) BunchOf(
	description string,
	fgt FuncGenType,
	name string,
	suffix string,
	rt ReturnType,
	bodyGen func(v *Helper, codeRange *Statement) []*Statement,
	testData interface{},
) {
	FuncGen{v, fgt}.Func()(description, name, ParamArgs, rt,
		bodyGen(v, v.RangeSlice("kk")),
		testData,
	)
	FuncGen{v, fgt}.Func()(description, name+suffix+"Slice", ParamSlice, rt,
		bodyGen(v, v.RangeSlice("kk")),
		testData,
	)
	FuncGen{v, fgt}.Func()(description, name+suffix+"Set", ParamSet, rt,
		bodyGen(v, v.RangeSet("s")),
		testData,
	)
}

type FuncGen struct {
	h   *Helper
	fgt FuncGenType
}

type FuncGenType uint8

const (
	SetMethod FuncGenType = iota
	SliceMethod
	SimpleFunc
)

//go:generate stringer -type FuncGenType

func (v FuncGen) Func() (result func(string, string, ParamType, ReturnType, []*Statement, interface{})) {
	switch v.fgt {
	case SetMethod:
		result = v.h.SetMethod
	case SimpleFunc:
		result = v.h.Func
	default:
		bwerr.Panic("fgt: %s", v.fgt)
	}
	return
}

func (v *Helper) SetMethod(
	description string,
	name string,
	pt ParamType,
	rt ReturnType,
	blockStatements []*Statement,
	testData interface{},
) {
	v.funcHelper(description, SetMethod, Params(Id("v").Id(v.IdSet)), name, pt, rt, blockStatements, testData)
}

func (v *Helper) Func(
	description string,
	name string,
	pt ParamType,
	rt ReturnType,
	blockStatements []*Statement,
	testData interface{},
) {
	v.funcHelper(description, SimpleFunc, nil, name, pt, rt, blockStatements, testData)
}

func (v *Helper) SliceMethod(
	description string,
	name string,
	pt ParamType,
	rt ReturnType,
	blockStatements []*Statement,
	testData interface{},
) {
	v.funcHelper(description, SliceMethod, Params(Id("v").Id(v.idSlice)), name, pt, rt, blockStatements, testData)
}

func (v *Helper) funcHelper(
	description string,
	fgt FuncGenType,
	preParams *Statement,
	name string,
	pt ParamType,
	rt ReturnType,
	blockStatements []*Statement,
	testData interface{},
) {
	statements := []Code{}
	for _, bs := range blockStatements {
		statements = append(statements, bs)
	}
	if len(description) > 0 {
		v.f.Commentf(name + " - " + description)
	}
	v.f.Func().
		Add(preParams).
		Id(name).
		Params(v.params(pt)...).
		Add(v.returns(rt)).
		Block(statements...).
		Line()
	if fgt != SliceMethod && v.needTests && testData != nil {
		v.GenTestsFor(fgt, name, pt, rt, testData)
	}
}

func (v *Helper) params(pt ParamType) []Code {
	result := []Code{}
	switch pt {
	case ParamNone:
	case ParamArg:
		result = append(result, Id("k").Id(v.IdItem))
	case ParamArgs:
		result = append(result, Id("kk").Id(`...`+v.IdItem))
	case ParamSlice:
		result = append(result, Id("kk").Index().Id(v.IdItem))
	case ParamSet:
		result = append(result, Id("s").Id(v.IdSet))
	case ParamIJ:
		result = append(result, Id("i").Int(), Id("j").Int())
	default:
		bwerr.Panic("pt: %d", pt)
	}
	return result
}

func (v *Helper) returns(rt ReturnType) (result *Statement) {
	switch rt {
	case ReturnNone:
	case ReturnBool:
		result = Bool()
	case ReturnInt:
		result = Int()
	case ReturnInterface:
		result = Interface()
	case ReturnString:
		result = String()
	case ReturnSet:
		result = Id(v.IdSet)
	case ReturnSlice:
		result = Index().Id(v.IdItem)
	case ReturnSliceOfStrings:
		result = Index().String()
	case ReturnJSON:
		result = Params(Index().Byte(), Error())
	default:
		// bwerr.PanicA(bw.A{"rt: %d", bw.Args(rt)})
		bwerr.Panic("rt: %s (%d)", rt, rt)
	}
	return result
}

func (v *Helper) ToString(id string) (codeString *Statement) {
	switch v.IdItem {
	case "string":
		codeString = Id(id)
	case "int8", "int16", "int32", "int64", "int":
		codeString = Qual("strconv", "FormatInt").Call(Id("int64").Call(Id(id)), Lit(10))
	case "uint8", "uint16", "uint32", "uint64", "uint":
		codeString = Qual("strconv", "FormatUint").Call(Id("uint64").Call(Id(id)), Lit(10))
	case "float32", "float64":
		codeString = Qual("strconv", "FormatFloat").Call(Id("float64").Call(Id(id)), LitByte('f'), Lit(-1), Lit(64))
	case "bool":
		codeString = Qual("strconv", "FormatBool").Call(Id(id))
	case "rune":
		codeString = Id("string").Call(Id(id))
	case "interface{}":
		codeString = Qual(v.bwjsonPackageName, "Pretty").Call(Id(id))
	default:
		codeString = Id(id).Dot("String").Call()
	}
	return
}

// func (v *Helper) ToDataForJSON(id string) (codeString *Statement) {
// 	switch v.IdItem {
// 	case "string":
// 		codeString = Id(id)
// 	case "int8", "int16", "int32", "int64", "int":
// 		codeString = Id(id)
// 		// codeString = Qual("strconv", "FormatInt").Call(Id("int64").Call(Id(id)), Lit(10))
// 	case "uint8", "uint16", "uint32", "uint64", "uint":
// 		codeString = Id(id)
// 		// codeString = Qual("strconv", "FormatUint").Call(Id("uint64").Call(Id(id)), Lit(10))
// 	case "float32", "float64":
// 		codeString = Id(id)
// 		// codeString = Qual("strconv", "FormatFloat").Call(Id("float64").Call(Id(id)), LitByte('f'), Lit(-1), Lit(64))
// 	case "bool":
// 		codeString = Id(id)
// 		// codeString = Qual("strconv", "FormatBool").Call(Id(id))
// 	case "rune":
// 		codeString = Id(id)
// 		// codeString = Id("string").Call(Id(id))
// 	case "interface{}":
// 		codeString = Id(id)
// 	default:
// 		codeString = Id(id).Dot("DataForJSON").Call()
// 	}
// 	return
// }

func (v *Helper) Save() {
	v.saveHelper(v.f, "")
	if len(v.tests) > 0 {
		v.t.Func().Id("Test" + strings.Title(v.IdSet)).Params(Id("t").Op("*").Qual("testing", "T")).Block(v.tests...)
		v.saveHelper(v.t, "_test")
	}
}

var nonFileNameRegexp = regexp.MustCompile(`[^_\w\d]+`)

func (v *Helper) saveHelper(f *File, suffix string) {
	name := nonFileNameRegexp.ReplaceAllLiteralString(v.IdItem, ``)
	fileName := name + "_set" + suffix + ".go"
	fileSpec := filepath.Join(v.packageDir, strings.ToLower(fileName))
	if err := f.Save(fileSpec); err != nil {
		bwerr.PanicA(bwerr.E{Error: err})
	}
}

// ============================================================================
