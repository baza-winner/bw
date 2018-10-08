package defvalid

import (
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/defparse"
	"github.com/baza-winner/bwcore/defvalid/deftype"
	"testing"
)

// ============================================================================

func TestCompileDef(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"def: nil": {
			In: []interface{}{nil},
			Out: []interface{}{
				(*Def)(nil),
				bwerror.Error("<ansiOutline>def<ansiCmd><ansi> (<ansiSecondaryLiteral>null<ansi>) has non supported value"),
			},
		},
		"def: invalid type": {
			In: []interface{}{defparse.MustParse(`false`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerror.Error("<ansiOutline>def<ansiCmd><ansi> (<ansiSecondaryLiteral>false<ansi>) is not of type <ansiPrimaryLiteral>[]string<ansi>, or <ansiPrimaryLiteral>map[string]<ansi>, or <ansiPrimaryLiteral>string"),
			},
		},
		"def: simple valid": {
			In: []interface{}{defparse.MustParse(`"bool"`)},
			Out: []interface{}{
				&Def{tp: deftype.FromArgs(deftype.Bool)},
				nil,
			},
		},
		"def: invalid deftypeItem": {
			In: []interface{}{defparse.MustParse(`[ qw/ bool int some / ]`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerror.Error("<ansiOutline>def<ansiCmd>.#2<ansi> (<ansiSecondaryLiteral>\"some\"<ansi>) has non supported value"),
			},
		},
		"def: enum": {
			In: []interface{}{defparse.MustParse(`{ type: "string", enum: [qw/one two three/]}`)},
			Out: []interface{}{
				&Def{tp: deftype.FromArgs(deftype.String), enum: bwset.StringsFromArgs("one", "two", "three")},
				nil,
			},
		},
		"def: map with keys": {
			In: []interface{}{defparse.MustParse(`{ type: "map", keys: { keyBool: ['bool'] }}`)},
			Out: []interface{}{
				&Def{
					tp: deftype.FromArgs(deftype.Map),
					keys: map[string]Def{
						"keyBool": Def{tp: deftype.FromArgs(deftype.Bool)},
					}},
				nil,
			},
		},
		"def: unexpected keys": {
			In: []interface{}{defparse.MustParse(`{ type: "map", kyes: { keyBool: ['bool'] }}`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.TestCaseStruct) error {
					return bwerror.Error(
						"<ansiOutline>def<ansiCmd><ansi> (<ansiSecondaryLiteral>" +
							bwjson.PrettyJson(test.In[0]) +
							"<ansi>) has unexpected key <ansiPrimaryLiteral>kyes",
					)
				},
			},
		},
		"def: array with arrayElem": {
			In: []interface{}{defparse.MustParse(`{ type: "array", arrayElem: 'int' }`)},
			Out: []interface{}{
				&Def{
					tp:        deftype.FromArgs(deftype.Array),
					arrayElem: &Def{tp: deftype.FromArgs(deftype.Int)},
				},
				nil,
			},
		},
		"def: array with arrayElem and Elem": {
			In: []interface{}{defparse.MustParse(`{ type: "array", arrayElem: 'int', elem: 'bool' }`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.TestCaseStruct) error {
					return bwerror.Error(
						"<ansiOutline>def<ansiCmd><ansi> (<ansiSecondaryLiteral>" +
							bwjson.PrettyJson(test.In[0]) +
							"<ansi>) has unexpected key <ansiPrimaryLiteral>elem",
					)
				},
			},
		},
		"def: minInt, maxInt": {
			In: []interface{}{defparse.MustParse(`{ type: "int", minInt: -6, maxInt: 10 }`)},
			Out: []interface{}{
				&Def{
					tp:     deftype.FromArgs(deftype.Int),
					minInt: ptrToInt64(-6),
					maxInt: ptrToInt64(10),
				},
				nil,
			},
		},
		"def: minInt > maxInt": {
			In: []interface{}{defparse.MustParse(`{ type: "int", minInt: 6, maxInt: -10 }`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.TestCaseStruct) error {
					return bwerror.Error(
						"<ansiOutline>def<ansiCmd><ansi> (<ansiSecondaryLiteral>" +
							bwjson.PrettyJson(test.In[0]) +
							"<ansi>) has conflicting keys: <ansiSecondaryLiteral>" +
							bwjson.PrettyJson(defparse.MustParse("{ minInt: 6, maxInt: -10 }")),
					)
				},
			},
		},
		"def: minNumber, maxNumber": {
			In: []interface{}{defparse.MustParse(`{ type: "number", minNumber: -6, maxNumber: 10 }`)},
			Out: []interface{}{
				&Def{
					tp:        deftype.FromArgs(deftype.Number),
					minNumber: ptrToFloat64(float64(-6)),
					maxNumber: ptrToFloat64(float64(10)),
				},
				nil,
			},
		},
		"def: default": {
			In: []interface{}{defparse.MustParse(`{ type: "bool", default: true }`)},
			Out: []interface{}{
				&Def{
					tp:         deftype.FromArgs(deftype.Bool),
					dflt:       true,
					isOptional: true,
				},
				nil,
			},
		},
		"def: default 'string' for bool": {
			In: []interface{}{defparse.MustParse(`{ type: "bool", default: "string" }`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerror.Error("<ansiOutline>def<ansiCmd>.default<ansi> (<ansiSecondaryLiteral>\"string\"<ansi>) is not of type <ansiPrimaryLiteral>Bool"),
			},
		},
		"def: isOptional": {
			In: []interface{}{defparse.MustParse(`{ type: "bool", isOptional: true }`)},
			Out: []interface{}{
				&Def{
					tp:         deftype.FromArgs(deftype.Bool),
					isOptional: true,
				},
				nil,
			},
		},
		"def: isOptional = false conflicts with dflt": {
			In: []interface{}{defparse.MustParse(`{ type: "bool", isOptional: false, default: false }`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.TestCaseStruct) error {
					return bwerror.Error(
						"<ansiOutline>def<ansiCmd><ansi> (<ansiSecondaryLiteral>" +
							bwjson.PrettyJson(test.In[0]) +
							"<ansi>) has conflicting keys: <ansiSecondaryLiteral>" +
							bwjson.PrettyJson(defparse.MustParse("{ isOptional: false, default: false }")),
					)
				},
			},
		},
		"def: arrayOf without follower": {
			In: []interface{}{defparse.MustParse(`{ type:  "arrayOf"  }`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerror.Error("<ansiOutline>def<ansiCmd>.type<ansi> (<ansiSecondaryLiteral>\"arrayOf\"<ansi>) must be followed by some type"),
			},
		},
		"simple.def: arrayOf without follower": {
			In: []interface{}{defparse.MustParse(`"arrayOf"`)},
			Out: []interface{}{
				(*Def)(nil),
				bwerror.Error("<ansiOutline>def<ansiCmd><ansi> (<ansiSecondaryLiteral>\"arrayOf\"<ansi>) must be followed by some type"),
			},
		},
		"simple.def: arrayOf can not be combined with array": {
			In: []interface{}{defparse.MustParse(`["arrayOf", "array", "string"]`)},
			Out: []interface{}{
				(*Def)(nil),
				func(test bwtesting.TestCaseStruct) error {
					return bwerror.Error(
						"<ansiOutline>def<ansiCmd><ansi> (<ansiSecondaryLiteral>" +
							bwjson.PrettyJson(test.In[0]) +
							"<ansi>) following values can not be combined: <ansiSecondaryLiteral>" +
							bwjson.PrettyJson(defparse.MustParse("[ 'array', 'arrayOf' ]")),
					)
				},
			},
		},
	}
	testsToRun := tests
	bwmap.CropMap(testsToRun)
	// bwmap.CropMap(testsToRun, "def: arrayOf without follower")
	bwtesting.BwRunTests(t, testsToRun, CompileDef)
}

// ============================================================================

func TestValidateVal(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"val: nil, simple.def: bool": {
			In: []interface{}{
				"val",
				nil,
				MustCompileDef(defparse.MustParse("'bool'")),
			},
			Out: []interface{}{
				nil,
				bwerror.Error("<ansiOutline>val<ansiCmd><ansi> (<ansiSecondaryLiteral>null<ansi>) is not of type <ansiPrimaryLiteral>Bool"),
			},
		},
		"val: nil, def: bool.isOptional": {
			In: []interface{}{
				"val",
				nil,
				MustCompileDef(defparse.MustParse("{ type: 'bool', isOptional: true }")),
			},
			Out: []interface{}{
				nil,
				nil,
			},
		},
		"val: nil, def: bool.default=true": {
			In: []interface{}{
				"val",
				nil,
				MustCompileDef(defparse.MustParse("{ type: 'bool', default: true }")),
			},
			Out: []interface{}{
				true,
				nil,
			},
		},
		"val: nil, simple.def: map": {
			In: []interface{}{
				"val",
				nil,
				MustCompileDef(defparse.MustParse("'map'")),
			},
			Out: []interface{}{
				map[string]interface{}{},
				nil,
			},
		},
		"val: nil, def: map": {
			In: []interface{}{
				"val",
				nil,
				MustCompileDef(defparse.MustParse("{ type: 'map' }")),
			},
			Out: []interface{}{
				map[string]interface{}{},
				nil,
			},
		},
		"val: valid, def: string.enum": {
			In: []interface{}{
				"val",
				"one",
				MustCompileDef(defparse.MustParse("{ type: 'string', enum: [qw/one two three/] }")),
			},
			Out: []interface{}{
				"one",
				nil,
			},
		},
		"val: invalid, def: string.enum": {
			In: []interface{}{
				"val",
				"One",
				MustCompileDef(defparse.MustParse("{ type: 'string', enum: [qw/one two three/] }")),
			},
			Out: []interface{}{
				nil,
				bwerror.Error("<ansiOutline>val<ansiCmd><ansi> (<ansiSecondaryLiteral>\"One\"<ansi>) has non supported value"),
			},
		},
		"val: valid, def: int.min": {
			In: []interface{}{
				"val",
				1,
				MustCompileDef(defparse.MustParse("{ type: 'int', minInt: 0 }")),
			},
			Out: []interface{}{
				1,
				nil,
				// bwerror.Error("<ansiOutline>val<ansiCmd><ansi> (<ansiSecondaryLiteral>\"One\"<ansi>) has non supported value"),
			},
		},
		// "val: unexpected keys": {
		// 	what: "somewhere",
		// 	val:  defparse.MustParseMap(`{ some: true, thing: false }`),
		// 	def:  defparse.MustParseMap(`{ type: 'map', keys: {} }`),
		// 	err: func(testIntf interface{}) (err error) {
		// 		if test, ok := testIntf.(testValidateValStruct); !ok {
		// 			bwerror.Panic("<ansiOutline>testIntf<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expected to be <ansiPrimaryLiteral>testGetValidValStruct")
		// 		} else {
		// 			err = fmt.Errorf(ansi.Ansi(`Err`,
		// 				"ERR: <ansiOutline>somewhere<ansiCmd><ansi> (<ansiSecondaryLiteral>"+
		// 					bwjson.PrettyJson(test.val)+
		// 					"<ansi>) has unexpected keys <ansiSecondaryLiteral>"+
		// 					bwjson.PrettyJson(defparse.MustParse("[ qw/some thing/ ]"))))
		// 		}
		// 		return
		// 	},
		// },
		// "def.keys.KEY: is not map": {
		// 	what: "somewhere",
		// 	val:  defparse.MustParseMap(`{  }`),
		// 	def:  defparse.MustParseMap(`{ type: 'map', keys: { some: true } }`),
		// 	err:  fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd>.keys.some<ansi> (<ansiSecondaryLiteral>true<ansi>) is not of type <ansiPrimaryLiteral>map")),
		// },
		// "val: some == 0 is not bool": {
		// 	what: "somewhere",
		// 	val:  defparse.MustParseMap(`{  some: 0 }`),
		// 	def:  defparse.MustParseMap(`{ type: 'map', keys: { some: { type: 'bool' } } }`),
		// 	err:  fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere<ansiCmd>.some<ansi> (<ansiSecondaryLiteral>0<ansi>) is not of type <ansiPrimaryLiteral>bool")),
		// },
		// "def: .default": {
		// 	what: "somewhere",
		// 	val:  defparse.MustParseMap(`{ }`),
		// 	def: defparse.MustParseMap(`{
		// 			type: 'map',
		// 			keys: {
		// 				boolKey: {
		// 					type: 'bool',
		// 					default: false
		// 				}
		// 				strKey: {
		// 					type: 'string'
		// 					default "something"
		// 				}
		// 			}
		// 		}`),
		// 	result: defparse.MustParseMap(`{
		// 		boolKey: false
		// 		strKey: "something"
		// 	}`),
		// },
		// "def.keys: unexpected key": {
		// 	what: "somewhere",
		// 	val:  defparse.MustParseMap(`{ }`),
		// 	def:  defparse.MustParseMap(`{ type: 'map', keys: { keyOne: { type: 'bool', default: false, some: "thing" } } }`),
		// 	err: func(testIntf interface{}) (err error) {
		// 		if test, ok := testIntf.(testValidateValStruct); !ok {
		// 			bwerror.Panic("<ansiOutline>testIntf<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expected to be <ansiPrimaryLiteral>testGetValidValStruct")
		// 		} else {
		// 			err = fmt.Errorf(ansi.Ansi(`Err`,
		// 				"ERR: <ansiOutline>somewhere::def<ansiCmd>.keys.keyOne<ansi> (<ansiSecondaryLiteral>"+
		// 					bwjson.PrettyJson(MustValOfPath(test.def, ".keys.keyOne"))+
		// 					"<ansi>) has unexpected key <ansiPrimaryLiteral>some"))
		// 		}
		// 		return
		// 	},
		// },
		// "def.type: []string": {
		// 	what:   "somewhere",
		// 	val:    defparse.MustParse(`true`),
		// 	def:    defparse.MustParseMap(`{ type: [ 'map', 'bool' ] }`),
		// 	result: true,
		// },
		// "simple def": {
		// 	what:   "somewhere",
		// 	val:    defparse.MustParse(`{ some: "thing"}`),
		// 	def:    defparse.MustParse(`[ 'map', 'bool' ]`),
		// 	result: defparse.MustParse(`{ some: "thing"}`),
		// },
		// "simple bool at nil": {
		// 	what: "somewhere",
		// 	val:  nil,
		// 	def:  defparse.MustParse(`'bool'`),
		// 	err:  fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere<ansiCmd><ansi> (<ansiSecondaryLiteral>null<ansi>) is not of type <ansiPrimaryLiteral>bool")),
		// },
		// "bool without default at nil": {
		// 	what: "somewhere",
		// 	val:  nil,
		// 	def:  defparse.MustParse(`{ type: 'bool' }`),
		// 	err:  fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere<ansiCmd><ansi> (<ansiSecondaryLiteral>null<ansi>) is not of type <ansiPrimaryLiteral>bool")),
		// },
		// "def.default: .type'bool' with .default'nil'": {
		// 	what: "somewhere",
		// 	val:  true,
		// 	def:  defparse.MustParse(`{ type: 'bool', default: nil }`),
		// 	err:  fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd>.default<ansi> (<ansiSecondaryLiteral>null<ansi>) has non supported value")),
		// },
		// "bool with non bool default": {
		// 	what: "somewhere",
		// 	val:  false,
		// 	def:  defparse.MustParse(`{ type: 'bool', default: "some" }`),
		// 	err:  fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd>.default<ansi> (<ansiSecondaryLiteral>\"some\"<ansi>) is not of type <ansiPrimaryLiteral>bool")),
		// },
		// "def: unexpected key": {
		// 	what: "somewhere",
		// 	val:  false,
		// 	def:  defparse.MustParse(`{ type: 'bool', some: "thing" }`),
		// 	err: func(testIntf interface{}) (err error) {
		// 		if test, ok := testIntf.(testValidateValStruct); !ok {
		// 			bwerror.Panic("<ansiOutline>testIntf<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expected to be <ansiPrimaryLiteral>testGetValidValStruct")
		// 		} else {
		// 			err = fmt.Errorf(ansi.Ansi(`Err`,
		// 				"ERR: <ansiOutline>somewhere::def<ansiCmd><ansi> (<ansiSecondaryLiteral>"+
		// 					bwjson.PrettyJson(test.def)+
		// 					"<ansi>) has unexpected key <ansiPrimaryLiteral>some"))
		// 		}
		// 		return
		// 	},
		// },
		// "def: unexpected keys": {
		// 	what: "somewhere",
		// 	val:  false,
		// 	def:  defparse.MustParse(`{ type: 'bool', some: 0, thing: nil }`),
		// 	err: func(testIntf interface{}) (err error) {
		// 		if test, ok := testIntf.(testValidateValStruct); !ok {
		// 			bwerror.Panic("<ansiOutline>testIntf<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expected to be <ansiPrimaryLiteral>testGetValidValStruct")
		// 		} else {
		// 			err = fmt.Errorf(ansi.Ansi(`Err`,
		// 				"ERR: <ansiOutline>somewhere::def<ansiCmd><ansi> (<ansiSecondaryLiteral>"+
		// 					bwjson.PrettyJson(test.def)+
		// 					"<ansi>) has unexpected keys <ansiSecondaryLiteral>"+
		// 					bwjson.PrettyJson(defparse.MustParse(`[ qw/some thing/ ]`))))
		// 		}
		// 		return
		// 	},
		// },
		// "default bool": {
		// 	what:   "somewhere",
		// 	val:    nil,
		// 	def:    defparse.MustParse(`{ type: 'bool', default: false }`),
		// 	result: false,
		// },
		// "simple def: enum is not supported": {
		// 	what: "somewhere",
		// 	val:  true,
		// 	def:  defparse.MustParse(`'enum'`),
		// 	err:  fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere::def<ansiCmd><ansi> (<ansiSecondaryLiteral>\"enum\"<ansi>) has non supported value")),
		// },
		// "def.type: enum and string can not be combined": {
		// 	what: "somewhere",
		// 	val:  "something",
		// 	def:  defparse.MustParse(`{ type: ['enum', 'string'] }`),
		// 	err: func(testIntf interface{}) (err error) {
		// 		if test, ok := testIntf.(testValidateValStruct); !ok {
		// 			bwerror.Panic("<ansiOutline>testIntf<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expected to be <ansiPrimaryLiteral>testGetValidValStruct")
		// 		} else {
		// 			err = fmt.Errorf(ansi.Ansi(`Err`,
		// 				"ERR: <ansiOutline>somewhere::def<ansiCmd>.type<ansi> (<ansiSecondaryLiteral>"+
		// 					bwjson.PrettyJson(MustValOfPath(test.def, ".type"))+
		// 					"<ansi>) following values can not be combined: <ansiSecondaryLiteral>"+
		// 					bwjson.PrettyJson(defparse.MustParse(`[ qw/enum string/ ]`))))
		// 		}
		// 		return
		// 	},
		// },
		// "def.enum: is not []string": {
		// 	what: "somewhere",
		// 	val:  "one",
		// 	def:  defparse.MustParse(`{ type: 'enum', enum: [ "one", true, 3 ] }`),
		// 	err: func(testIntf interface{}) (err error) {
		// 		if test, ok := testIntf.(testValidateValStruct); !ok {
		// 			bwerror.Panic("<ansiOutline>testIntf<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expected to be <ansiPrimaryLiteral>testGetValidValStruct")
		// 		} else {
		// 			err = fmt.Errorf(ansi.Ansi(`Err`,
		// 				"ERR: <ansiOutline>somewhere::def<ansiCmd>.enum<ansi> (<ansiSecondaryLiteral>"+
		// 					bwjson.PrettyJson(MustValOfPath(test.def, ".enum"))+
		// 					"<ansi>) is not of type <ansiPrimaryLiteral>[]string"))
		// 		}
		// 		return
		// 	},
		// },
		// "val: is not supported by def.enum": {
		// 	what: "somewhere",
		// 	val:  "One",
		// 	def:  defparse.MustParse(`{ type: 'enum', enum: [ "one", "two", "tree" ] }`),
		// 	err:  fmt.Errorf(ansi.Ansi(`Err`, "ERR: <ansiOutline>somewhere<ansiCmd><ansi> (<ansiSecondaryLiteral>\"One\"<ansi>) has non supported value")),
		// },
		// "ExecCmd opt": {
		// 	what: "ExecCmd.opt",
		// 	val:  nil,
		// 	def: defparse.MustParse(`{
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
		// 	result: defparse.MustParse(`{
		// 		v: 'none'
		// 		s: 'all'
		// 		exitOnError: false
		// 	}`),
		// },
	}
	testsToRun := tests
	bwmap.CropMap(testsToRun)
	// bwmap.CropMap(testsToRun, "def.type: invalid type")
	bwtesting.BwRunTests(t, testsToRun, ValidateVal)
}
