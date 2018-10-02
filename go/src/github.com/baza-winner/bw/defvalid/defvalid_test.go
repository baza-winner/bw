package defvalid

import (
	// "bufio"
	// "encoding/json"
	// "errors"
	"fmt"
	// "encoding/json"
	"github.com/baza-winner/bw/ansi"
	"github.com/baza-winner/bw/bwtesting"
	"github.com/baza-winner/bw/core"
	"github.com/baza-winner/bw/defparse"
	// "reflect"
	"testing"
	// "github.com/iancoleman/strcase"
	// "github.com/jimlawless/whereami"
	// "log"
	// "os"
	// "os/exec"
	// "strings"
	// "syscall"
)

type testGetValidValStruct struct {
	where    string
	val      interface{}
	def      map[string]interface{}
	whereDef string
	result   interface{}
	err      error
}

func TestGetValidVal(t *testing.T) {
	tests := map[string]testGetValidValStruct{
		"def nil": {
			where:    "<ansiOutline>somewhere<ansiCmd>",
			val:      defparse.MustParseMap(`{ }`),
			def:      nil,
			whereDef: "<ansiOutline>somewhere::def<ansiCmd>",
			err:      fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansi> is <ansiPrimaryLiteral>nil")),
		},
		"no type": {
			where:    "<ansiOutline>somewhere<ansiCmd>",
			val:      defparse.MustParseMap(`{ }`),
			def:      defparse.MustParseMap(`{ }`),
			whereDef: "<ansiOutline>somewhere::def<ansiCmd>",
			err:      fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansi> has no key <ansiPrimaryLiteral>type")),
		},
		"type is not string": {
			where:    "<ansiOutline>somewhere<ansiCmd>",
			val:      defparse.MustParseMap(`{ }`),
			def:      defparse.MustParseMap(`{ type: false }`),
			whereDef: "<ansiOutline>somewhere::def<ansiCmd>",
			err:      fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd>.type<ansi> (<ansiSecondaryLiteral>false<ansi>) is not <ansiPrimaryLiteral>string")),
		},
		"wrong type": {
			where:    "<ansiOutline>somewhere<ansiCmd>",
			val:      defparse.MustParseMap(`{  }`),
			def:      defparse.MustParseMap(`{ type: 'some' }`),
			whereDef: "<ansiOutline>somewhere::def<ansiCmd>",
			err:      fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd>.type<ansi> has non supported value <ansiPrimaryLiteral>some")),
		},
		"unexpected key": {
			where:    "<ansiOutline>somewhere<ansiCmd>",
			val:      defparse.MustParseMap(`{ some: true }`),
			def:      defparse.MustParseMap(`{ type: 'map', keys: {} }`),
			whereDef: "<ansiOutline>somewhere::def<ansiCmd>",
			err:      fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere<ansi> (<ansiSecondaryLiteral>"+core.PrettyJson(defparse.MustParseMap("{some: true}"))+"<ansi>) has unexpected key <ansiPrimaryLiteral>some")),
		},
		"keys.KEY is not map": {
			where:    "<ansiOutline>somewhere<ansiCmd>",
			val:      defparse.MustParseMap(`{  }`),
			def:      defparse.MustParseMap(`{ type: 'map', keys: { some: true } }`),
			whereDef: "<ansiOutline>somewhere::def<ansiCmd>",
			err:      fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd>.keys.some<ansi> (<ansiSecondaryLiteral>true<ansi>) is not <ansiPrimaryLiteral>map")),
		},
		// ""
		// where: "<ansiOutline>somewhere<ansiCmd>",
		// // val: defparse.MustParseMap(`{ exitCode: nil, s: 1, v: "ALL", some: 0 }`),
		//    val: defparse.MustParseMap(`{ type: 'some' }`),
		// def: defparse.MustParseMap(`{
		//      type: 'map',
		//      keys: {
		//        v: {
		//          type: 'enum'
		//          enum: [ qw/all err ok none/ ]
		//          default: 'none'
		//        }
		//        s: {
		//          type: 'enum'
		//          enum: [ qw/none stderr stdout all/ ]
		//          default: 'all'
		//        }
		//        exitOnError: {
		//          type: 'bool'
		//          default: false
		//        }
		//      }
		//    }`),
		// whereDef: "<ansiOutline>somewhere::def<ansiCmd>",
		// result:   defparse.MustParseMap(`{ v: 'enum', s: 'all', exitOnError: false }`),
		// err:      nil,
		// },
	}
	testsToRun := tests
	for testName, test := range testsToRun {
		t.Logf(ansi.Ansi(`Header`, "Running test case <ansiPrimaryLiteral>%s"), testName)
		result, err := GetValidVal(test.where, test.val, test.def, test.whereDef)
		testTitle := fmt.Sprintf("GetValidVal(%s, %+s, %+s, %s)\n", test.where, core.PrettyJson(test.val), core.PrettyJson(test.def), test.whereDef)
		bwtesting.CheckTestErrResult(t, err, test.err, result, test.result, testTitle)
		//   if bwtesting.CompareErrors(t, err, test.err, testTitle) && !reflect.DeepEqual(result, test.result) { // https://stackoverflow.com/questions/18208394/testing-equivalence-of-maps-golang
		// 	t.Errorf(ansi.Ansi("", testTitle + "    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), core.PrettyJson(result), core.PrettyJson(test.result))
		// }
	}
}
