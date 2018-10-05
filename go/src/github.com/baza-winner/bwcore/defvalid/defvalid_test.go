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
		".type is invalid type": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: defparse.MustParseMap(`{ }`),
			},
			def: value{
				value: defparse.MustParseMap(`{ type: false }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			err: fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd>.type<ansi> (<ansiSecondaryLiteral>false<ansi>) is not of type <ansiPrimaryLiteral>[]string<ansi> or <ansiPrimaryLiteral>string")),
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
			err: fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd>.type<ansi> (<ansiSecondaryLiteral>\"some\"<ansi>) has non supported value")),
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
				value: defparse.MustParseMap(`{
					type: 'map',
					keys: {
						boolKey: {
							type: 'bool',
							default: false
						}
						strKey: {
							type: 'string'
							default "something"
						}
					}
				}`),
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
		".type: []string": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: defparse.MustParse(`true`),
			},
			def: value{
				value: defparse.MustParseMap(`{ type: [ 'map', 'bool' ] }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			result: true,
		},
		"simple def": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: defparse.MustParse(`{ some: "thing"}`),
			},
			def: value{
				value: defparse.MustParse(`[ 'map', 'bool' ]`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			result: defparse.MustParse(`{ some: "thing"}`),
		},
		"simple bool at nil": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: nil,
			},
			def: value{
				value: defparse.MustParse(`'bool'`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			err: fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere<ansiCmd><ansi> (<ansiSecondaryLiteral>null<ansi>) is not of type <ansiPrimaryLiteral>bool")),
		},
		"bool without default at nil": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: nil,
			},
			def: value{
				value: defparse.MustParse(`{ type: 'bool' }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			err: fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere<ansiCmd><ansi> (<ansiSecondaryLiteral>null<ansi>) is not of type <ansiPrimaryLiteral>bool")),
		},
		"bool with nil default": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: true,
			},
			def: value{
				value: defparse.MustParse(`{ type: 'bool', default: nil }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			err: fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd>.default<ansi> (<ansiSecondaryLiteral>null<ansi>) has non supported value")),
		},
		"bool with non bool default": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: false,
			},
			def: value{
				value: defparse.MustParse(`{ type: 'bool', default: "some" }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			err: fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd>.default<ansi> (<ansiSecondaryLiteral>\"some\"<ansi>) is not of type <ansiPrimaryLiteral>bool")),
		},
		"def with unexpected key": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: false,
			},
			def: value{
				value: defparse.MustParse(`{ type: 'bool', some: "thing" }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			err: func(testIntf interface{}) (err error) {
				if test, ok := testIntf.(testGetValidValStruct); !ok {
					bwerror.Panic("<ansiOutline>testIntf<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expected to be <ansiPrimaryLiteral>testGetValidValStruct")
				} else {
					err = fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd><ansi> (<ansiSecondaryLiteral>"+bwjson.PrettyJson(test.def.value)+"<ansi>) has unexpected key <ansiPrimaryLiteral>some"))
				}
				return
			},
		},
		"def with unexpected keys": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: false,
			},
			def: value{
				value: defparse.MustParse(`{ type: 'bool', some: 0, thing: nil }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			err: func(testIntf interface{}) (err error) {
				if test, ok := testIntf.(testGetValidValStruct); !ok {
					bwerror.Panic("<ansiOutline>testIntf<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expected to be <ansiPrimaryLiteral>testGetValidValStruct")
				} else {
					err = fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd><ansi> (<ansiSecondaryLiteral>"+bwjson.PrettyJson(test.def.value)+"<ansi>) has unexpected keys <ansiSecondaryLiteral>"+bwjson.PrettyJson(defparse.MustParse(`[ qw/some thing/ ]`))))
				}
				return
			},
		},
		"default bool": {
			val: value{
				where: "<ansiOutline>somewhere<ansiCmd>",
				value: nil,
			},
			def: value{
				value: defparse.MustParse(`{ type: 'bool', default: false }`),
				where: "<ansiOutline>somewhere::def<ansiCmd>",
			},
			result: false,
		},
		// "ExecCmd opt": {
		// 	val: value{
		// 		where: "<ansiOutline>ExecCmd.opt<ansiCmd>",
		// 		value: nil,
		// 	},
		// 	def: value{
		// 		where: "<ansiOutline>ExecCmd.opt::def<ansiCmd>",
		// 		value: defparse.MustParse(`{
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
		// 	},
		// 	result: defparse.MustParse(`{
		// 		v: 'none'
		// 		s: 'all'
		// 		exitOnError: false
		// 	}`),
		// },
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
	bwmap.CropMap(testsToRun, "default value")
	// bwmap.CropMap(testsToRun, "bool without default at nil", "default bool", "bool with nil default at nil")
	// bwmap.CropMap(testsToRun, "some == 0 is not bool")
	// testsToRun = map[string]testGetValidValStruct{"wrong type": tests["wrong type"]}
	for testName, test := range testsToRun {
		// bwtesting.Debug(test)
		bwtesting.BtRunTest(t, testName, test)
		// bwtesting.BtRunTest(t, testName, test, map[bwtesting.BtOptFuncType]interface{}{
		// 	bwtesting.WhenDiffErr: func(t *testing.T, bt bwtesting.BT, tstErr, etaErr error) {
		// 		bwtesting.WhenDiffErrFuncDefault(t, bt, tstErr, etaErr)
		// 		if valErr, ok := tstErr.(*value); ok {
		// 			fmt.Printf("err.where: %s\n", valErr.error.where)
		// 		}
		// 	},
		// })
	}
}
