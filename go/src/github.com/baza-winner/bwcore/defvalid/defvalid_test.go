package defvalid

import (
	"fmt"
	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/defparse"
	"testing"
)

type testGetValidValStruct struct {
	val    value
	def    value
	result interface{}
	err    interface{}
}

func (test testGetValidValStruct) GetTstResultErr() (interface{}, error) {
	return GetValidVal(test.val, test.def)
}

func (test testGetValidValStruct) GetTitle() string {
	return fmt.Sprintf("GetValidVal(%s, %s)\n", bwjson.PrettyJsonOf(test.val), bwjson.PrettyJsonOf(test.def))
}

func (test testGetValidValStruct) GetEtaErr() interface{} {
	return test.err
}

func (test testGetValidValStruct) GetEtaResult() interface{} {
	return test.result
}

func TestGetValidVal(t *testing.T) {
	tests := map[string]testGetValidValStruct{
		"def nil": {
			val: value{
				value: defparse.MustParseMap(`{ }`),
				where: "<ansiOutline>somewhere<ansiCmd>",
			},
			def: value{
				where: "<ansiOutline>somewhere::def<ansiCmd>",
				value: nil,
			},
			err: fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd><ansi> (<ansiSecondaryLiteral>null<ansi>) is not of type <ansiPrimaryLiteral>map")),
		},
		"no type": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: defparse.MustParseMap(`{ }`),
			},
			def: value{
				value: defparse.MustParseMap(`{ }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			err: fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd><ansi> (<ansiSecondaryLiteral>{}<ansi>) has no key <ansiPrimaryLiteral>type")),
		},
		"type is not string": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: defparse.MustParseMap(`{ }`),
			},
			def: value{
				value: defparse.MustParseMap(`{ type: false }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			err: fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd>.type<ansi> (<ansiSecondaryLiteral>false<ansi>) is not of type <ansiPrimaryLiteral>string")),
		},
		".type has non supported value": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: defparse.MustParseMap(`{  }`),
			},
			def: value{
				value: defparse.MustParseMap(`{ type: 'some' }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			err: fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd>.type<ansi> (<ansiSecondaryLiteral>\"some\"<ansi>) has non supported value <ansiPrimaryLiteral>some")),
		},
		"unexpected keys": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: defparse.MustParseMap(`{ some: true, thing: false }`),
			},
			def: value{
				value: defparse.MustParseMap(`{ type: 'map', keys: {} }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			// err: fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere<ansiCmd><ansi> (<ansiSecondaryLiteral>"+bwjson.PrettyJson(defparse.MustParseMap("{some: true}"))+"<ansi>) has unexpected keys <ansiSecondaryLiteral>" + bwjson.PrettyJson(defparse.MustParse("[ qw/some thing/ ]")))),
			err: func(testIntf interface{}) (err error) {
				if test, ok := testIntf.(testGetValidValStruct); !ok {
					bwerror.Panic("<ansiOutline>testIntf<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expected to be <ansiPrimaryLiteral>testGetValidValStruct")
				} else {
					err = fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere<ansiCmd><ansi> (<ansiSecondaryLiteral>"+bwjson.PrettyJson(test.val.value)+"<ansi>) has unexpected keys <ansiSecondaryLiteral>"+bwjson.PrettyJson(defparse.MustParse("[ qw/some thing/ ]"))))
					// err = fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd>.keys.keyOne<ansi> (<ansiSecondaryLiteral>"+bwjson.PrettyJson(MustValOfPath(test.def.value, ".keys.keyOne"))+"<ansi>) has unexpected key <ansiPrimaryLiteral>some"))
				}
				return
			},
		},
		"keys.KEY is not map": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: defparse.MustParseMap(`{  }`),
			},
			def: value{
				value: defparse.MustParseMap(`{ type: 'map', keys: { some: true } }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			err: fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd>.keys.some<ansi> (<ansiSecondaryLiteral>true<ansi>) is not of type <ansiPrimaryLiteral>map")),
		},
		"some == 0 is not bool": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: defparse.MustParseMap(`{  some: 0 }`),
			},
			def: value{
				value: defparse.MustParseMap(`{ type: 'map', keys: { some: { type: 'bool' } } }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			err: fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere<ansiCmd>.some<ansi> (<ansiSecondaryLiteral>0<ansi>) is not of type <ansiPrimaryLiteral>bool")),
		},
		"default value": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: defparse.MustParseMap(`{ }`),
			},
			def: value{
				value: defparse.MustParseMap(`{ type: 'map', keys: { some: { type: 'bool', default: false } } }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			result: defparse.MustParseMap(`{ some: false }`),
		},
		"unexpected key in keyDef": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: defparse.MustParseMap(`{ }`),
			},
			def: value{
				value: defparse.MustParseMap(`{ type: 'map', keys: { keyOne: { type: 'bool', default: false, some: "thing" } } }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			err: func(testIntf interface{}) (err error) {
				if test, ok := testIntf.(testGetValidValStruct); !ok {
					bwerror.Panic("<ansiOutline>testIntf<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expected to be <ansiPrimaryLiteral>testGetValidValStruct")
				} else {
					err = fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd>.keys.keyOne<ansi> (<ansiSecondaryLiteral>"+bwjson.PrettyJson(MustValOfPath(test.def.value, ".keys.keyOne"))+"<ansi>) has unexpected key <ansiPrimaryLiteral>some"))
				}
				return
			},
		},
		"type: []string": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: defparse.MustParse(`true`),
			},
			def: value{
				value: defparse.MustParseMap(`{ type: [ 'map', 'bool' ] }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			// err: func(testIntf interface{}) (err error) {
			// 	if test, ok := testIntf.(testGetValidValStruct); !ok {
			// 		bwerror.Panic("<ansiOutline>testIntf<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expected to be <ansiPrimaryLiteral>testGetValidValStruct")
			// 	} else {
			// 		err = fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd>.keys.keyOne<ansi> (<ansiSecondaryLiteral>"+bwjson.PrettyJson(MustValOfPath(test.def.value, ".keys.keyOne"))+"<ansi>) has unexpected key <ansiPrimaryLiteral>some"))
			// 	}
			// 	return
			// },
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
	bwmap.CropMap(testsToRun)
	// bwmap.CropMap(testsToRun, "def nil", "no type", "type is not string", ".type has non supported value", "unexpected keys", "keys.KEY is not map", "current")
	bwmap.CropMap(testsToRun, "type: []string")
	// bwmap.CropMap(testsToRun, "some == 0 is not bool")
	// testsToRun = map[string]testGetValidValStruct{"wrong type": tests["wrong type"]}
	for testName, test := range testsToRun {
		bwtesting.Debug(test)
		bwtesting.BtRunTest(t, testName, test)
	}
}
