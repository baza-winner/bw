package pfa

import (
	"testing"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/runeprovider"
)

// func getTestFileSpec(basename string) string {
// 	return filepath.Join(os.Getenv("GOPATH"), "src", "github.com/baza-winner/bwcore/pfa", basename)
// }

// func TestParseLogic(t *testing.T) {
// 	p, err := runeprovider.FromFile(getTestFileSpec("logic.pfa"))
// 	if err != nil {
// 		bwerror.PanicErr(err)
// 	}
// 	tests := map[string]bwtesting.TestCaseStruct{
// 		"": {
// 			In:  []interface{}{p},
// 			Out: []interface{}{nil, nil},
// 		},
// 	}
// 	testsToRun := tests
// 	bwmap.CropMap(testsToRun)
// 	// bwmap.CropMap(testsToRun, "qw ")
// 	// bwmap.CropMap(testsToRun, "qw && fa.curr.runePtr == EOF")
// 	bwtesting.BwRunTests(t, ParseLogic, testsToRun)
// }

func TestVarPathFrom(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"": {
			In:  []interface{}{""},
			Out: []interface{}{[]interface{}{}, nil},
		},
		"0": {
			In:  []interface{}{"0"},
			Out: []interface{}{[]interface{}{0}, nil},
		},
		"0.key": {
			In:  []interface{}{"0.key"},
			Out: []interface{}{[]interface{}{0, "key"}, nil},
		},
		"0.key.{1.some}": {
			In: []interface{}{"0.key.{1.some}"},
			Out: []interface{}{
				[]interface{}{
					0,
					"key",
					[]interface{}{1, "some"},
				},
				nil,
			},
		},
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "qw ")
	bwtesting.BwRunTests(t, VarPathFrom, tests)
}

func TestPfa_getVarValue(t *testing.T) {
	p := runeprovider.ProxyFrom(runeprovider.FromString("some"))
	pfa := pfaStruct{
		parseStack{
			parseStackItem{
				vars: map[string]interface{}{
					"type": "map",
					"value": map[string]interface{}{
						"boolKey":   true,
						"numberKey": 273,
						"stringKey": "string value",
						"runeKey":   '\n',
					},
				},
			},
			parseStackItem{
				vars: map[string]interface{}{
					"type":   "key",
					"string": "boolKey",
				},
			},
		},
		p,
		nil,
		map[string]interface{}{
			"primary": "begin",
			"result": map[string]interface{}{
				"some": "thing",
			},
		},
	}
	tests := map[string]bwtesting.TestCaseStruct{
		"primary": {
			In:  []interface{}{MustVarPathFrom("primary")},
			Out: []interface{}{VarValue{"begin", nil}},
		},
		"result.some": {
			In:  []interface{}{MustVarPathFrom("result.some")},
			Out: []interface{}{VarValue{"thing", nil}},
		},
		"stackLen": {
			In:  []interface{}{MustVarPathFrom("stackLen")},
			Out: []interface{}{VarValue{2, nil}},
		},
		"0.type": {
			In:  []interface{}{MustVarPathFrom("0.type")},
			Out: []interface{}{VarValue{"key", nil}},
		},
		"1.value.numberKey": {
			In:  []interface{}{MustVarPathFrom("1.value.numberKey")},
			Out: []interface{}{VarValue{273, nil}},
		},
		"1.value.{0.string}": {
			In:  []interface{}{MustVarPathFrom("1.value.{0.string}")},
			Out: []interface{}{VarValue{true, nil}},
		},
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "0.type")
	bwtesting.BwRunTests(t, pfa.getVarValue, tests)
}

func TestPfa_setVarVal(t *testing.T) {
	p := runeprovider.ProxyFrom(runeprovider.FromString("some"))
	pfa := pfaStruct{
		parseStack{
			parseStackItem{
				vars: map[string]interface{}{
					"type": "map",
					"value": map[string]interface{}{
						"boolKey":   true,
						"numberKey": 273,
						"stringKey": "string value",
						"runeKey":   '\n',
						"arrayKey":  []interface{}{"a", "b"},
					},
				},
			},
			parseStackItem{
				vars: map[string]interface{}{
					"type":   "key",
					"string": "boolKey",
					"array":  []interface{}{"c", "d"},
				},
			},
		},
		p,
		nil,
		map[string]interface{}{
			"primary": "begin",
			"result": map[string]interface{}{
				"some": "thing",
			},
		},
	}
	tests := map[string]bwtesting.TestCaseStruct{
		"primary": {
			In:  []interface{}{MustVarPathFrom("primary"), "end"},
			Out: []interface{}{VarValue{"end", nil}, nil},
		},
		"result.some": {
			In:  []interface{}{MustVarPathFrom("result.some"), "another"},
			Out: []interface{}{VarValue{"another", nil}, nil},
		},
		"stackLen": {
			In:  []interface{}{MustVarPathFrom("stackLen"), 4},
			Out: []interface{}{VarValue{nil, nil}, bwerror.Error("<ansiOutline>stackLen<ansi> is read only")},
		},
		"0.type": {
			In:  []interface{}{MustVarPathFrom("0.type"), "map"},
			Out: []interface{}{VarValue{"map", nil}, nil},
		},
		"0.item": {
			In:  []interface{}{MustVarPathFrom("0.item"), "word"},
			Out: []interface{}{VarValue{"word", nil}, nil},
		},
		"1.value.numberKey": {
			In:  []interface{}{MustVarPathFrom("1.value.numberKey"), "2.71"},
			Out: []interface{}{VarValue{"2.71", nil}, nil},
		},
		"1.value.{0.string}": {
			In:  []interface{}{MustVarPathFrom("1.value.{0.string}"), false},
			Out: []interface{}{VarValue{false, nil}, nil},
		},
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "0.type")
	bwtesting.BwRunTests(t, pfa.setVarValTestHelper, tests)
}

func (pfa *pfaStruct) setVarValTestHelper(varPath VarPath, varVal interface{}) (result VarValue, err error) {
	err = pfa.setVarVal(varPath, varVal)
	if err == nil {
		result = pfa.getVarValue(varPath)
	}
	return
}

func TestPfaActions(t *testing.T) {
	p := runeprovider.ProxyFrom(runeprovider.FromString("some"))
	pfa := pfaStruct{
		parseStack{
			parseStackItem{
				vars: map[string]interface{}{
					"type":  "map",
					"array": []interface{}{"g", "h"},
					"value": map[string]interface{}{
						"boolKey":   true,
						"numberKey": 273,
						"stringKey": "string value",
						"runeKey":   '\n',
						"arrayKey":  []interface{}{"a", "b"},
					},
				},
			},
			parseStackItem{
				vars: map[string]interface{}{
					"type":   "key",
					"string": "boolKey",
					"array":  []interface{}{"c", "d"},
					"value": map[string]interface{}{
						"arrayKey": []interface{}{"e", "f"},
					},
				},
			},
		},
		p,
		nil,
		map[string]interface{}{
			"array":   []interface{}{"i", "j"},
			"primary": "begin",
			"result": map[string]interface{}{
				"some": "thing",
			},
		},
	}
	tests := map[string]bwtesting.TestCaseStruct{
		"primary": {
			In:  []interface{}{MustVarPathFrom("primary"), []interface{}{SetVar{"primary", "end"}}},
			Out: []interface{}{VarValue{"end", nil}},
		},
		"0.type": {
			In:  []interface{}{MustVarPathFrom("0.type"), []interface{}{SetVarBy{"0.type", Var{"1.type"}, By{Append{}}}}},
			Out: []interface{}{VarValue{"keymap", nil}},
		},
		"0.array": {
			In:  []interface{}{MustVarPathFrom("0.array"), []interface{}{SetVarBy{"0.array", Var{"1.value.{0.string}"}, By{Append{}}}}},
			Out: []interface{}{VarValue{[]interface{}{"c", "d", true}, nil}},
		},
		"array": {
			In:  []interface{}{MustVarPathFrom("array"), []interface{}{SetVarBy{"array", Var{"1.array"}, By{Append{}}}}},
			Out: []interface{}{VarValue{[]interface{}{"i", "j", []interface{}{"g", "h"}}, nil}},
		},
		"1.value.arrayKey": {
			In:  []interface{}{MustVarPathFrom("1.value.arrayKey"), []interface{}{SetVarBy{"1.value.arrayKey", Var{"0.value.arrayKey"}, By{AppendSlice{}}}}},
			Out: []interface{}{VarValue{[]interface{}{"a", "b", "e", "f"}, nil}},
		},
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "0.type")
	bwtesting.BwRunTests(t, pfa.pfaActionsTestHelper, tests)
}

func (pfa *pfaStruct) pfaActionsTestHelper(varPath VarPath, args []interface{}) (result VarValue) {
	pfa.processRules(CreateRules(args))
	result = pfa.getVarValue(varPath)
	return
}

func TestPfaConditions(t *testing.T) {
	p := runeprovider.ProxyFrom(runeprovider.FromString("something"))
	p.PullRune()
	p.PullRune()
	p.PullRune()
	p.PullRune()
	pfa := pfaStruct{
		parseStack{
			parseStackItem{
				vars: map[string]interface{}{
					"type":  "map",
					"array": []interface{}{"g", "h"},
					"value": map[string]interface{}{
						"boolKey":   true,
						"numberKey": 273,
						"stringKey": "string value",
						"runeKey":   '\n',
						"arrayKey":  []interface{}{"a", "b"},
					},
				},
			},
			parseStackItem{
				vars: map[string]interface{}{
					"type":   "key",
					"string": "boolKey",
					"array":  []interface{}{"c", "d"},
					"value": map[string]interface{}{
						"arrayKey": []interface{}{"e", "f"},
					},
				},
			},
		},
		p,
		nil,
		map[string]interface{}{
			"array":   []interface{}{"i", "j"},
			"primary": "begin",
			"result": map[string]interface{}{
				"some": "thing",
			},
		},
	}
	tests := map[string]bwtesting.TestCaseStruct{
		"primary is begin => true": {
			In:  []interface{}{VarIs{"primary", "begin"}},
			Out: []interface{}{true, nil},
		},
		"primary is end => false": {
			In:  []interface{}{VarIs{"primary", "end"}},
			Out: []interface{}{false, nil},
		},
		"curr is 'e' => true": {
			In:  []interface{}{VarIs{"currRune", 'e'}},
			Out: []interface{}{false, nil},
		},
		// "0.type": {
		// 	In:  []interface{}{MustVarPathFrom("0.type"), []interface{}{SetVarBy{"0.type", Var{"1.type"}, By{Append{}}}}},
		// 	Out: []interface{}{VarValue{"keymap", nil}},
		// },
		// "0.array": {
		// 	In:  []interface{}{MustVarPathFrom("0.array"), []interface{}{SetVarBy{"0.array", Var{"1.value.{0.string}"}, By{Append{}}}}},
		// 	Out: []interface{}{VarValue{[]interface{}{"c", "d", true}, nil}},
		// },
		// "array": {
		// 	In:  []interface{}{MustVarPathFrom("array"), []interface{}{SetVarBy{"array", Var{"1.array"}, By{Append{}}}}},
		// 	Out: []interface{}{VarValue{[]interface{}{"i", "j", []interface{}{"g", "h"}}, nil}},
		// },
		// "1.value.arrayKey": {
		// 	In:  []interface{}{MustVarPathFrom("1.value.arrayKey"), []interface{}{SetVarBy{"1.value.arrayKey", Var{"0.value.arrayKey"}, By{AppendSlice{}}}}},
		// 	Out: []interface{}{VarValue{[]interface{}{"a", "b", "e", "f"}, nil}},
		// },
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "0.type")
	bwtesting.BwRunTests(t, pfa.pfaConditionsTestHelper, tests)
}

func (pfa *pfaStruct) pfaConditionsTestHelper(arg interface{}) (interface{}, error) {
	pfa.vars["result"] = false
	args := []interface{}{arg}
	args = append(args, SetVar{"result", true})
	pfa.processRules(CreateRules(args))
	return pfa.vars["result"], pfa.err
}
