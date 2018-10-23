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
			Out: []interface{}{VarPath{VarPathItem{0}}, nil},
		},
		"0.key": {
			In:  []interface{}{"0.key"},
			Out: []interface{}{VarPath{VarPathItem{0}, VarPathItem{"key"}}, nil},
		},
		"0.key.{1.some}": {
			In: []interface{}{"0.key.{1.some}"},
			Out: []interface{}{
				VarPath{
					VarPathItem{0},
					VarPathItem{"key"},
					VarPathItem{VarPath{VarPathItem{1}, VarPathItem{"some"}}},
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
	pfa := &pfaStruct{
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
			In:  []interface{}{"primary"},
			Out: []interface{}{"begin", nil},
		},
		"result.some": {
			In:  []interface{}{"result.some"},
			Out: []interface{}{"thing", nil},
		},
		"stackLen": {
			In:  []interface{}{"stackLen"},
			Out: []interface{}{2, nil},
		},
		"0.type": {
			In:  []interface{}{"0.type"},
			Out: []interface{}{"key", nil},
		},
		"1.value.numberKey": {
			In:  []interface{}{"1.value.numberKey"},
			Out: []interface{}{273, nil},
		},
		"1.value.{0.string}": {
			In:  []interface{}{"1.value.{0.string}"},
			Out: []interface{}{true, nil},
		},
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "primary")
	bwtesting.BwRunTests(t, pfa.TestPfa_getVarValueTestHelper, tests)
}

func (pfa *pfaStruct) TestPfa_getVarValueTestHelper(varPathStr string) (val interface{}, err error) {
	pfa.err = nil
	varValue := pfa.getVarValue(MustVarPathFrom(varPathStr))
	return varValue.val, pfa.err
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
			In:  []interface{}{"primary", "end"},
			Out: []interface{}{"end", nil},
		},
		"result.some": {
			In:  []interface{}{"result.some", "another"},
			Out: []interface{}{"another", nil},
		},
		"stackLen": {
			In:  []interface{}{"stackLen", 4},
			Out: []interface{}{nil, bwerror.Error("<ansiOutline>stackLen<ansi> is read only")},
		},
		"0.type": {
			In:  []interface{}{"0.type", "map"},
			Out: []interface{}{"map", nil},
		},
		"0.item": {
			In:  []interface{}{"0.item", "word"},
			Out: []interface{}{"word", nil},
		},
		"1.value.numberKey": {
			In:  []interface{}{"1.value.numberKey", "2.71"},
			Out: []interface{}{"2.71", nil},
		},
		"1.value.{0.string}": {
			In:  []interface{}{"1.value.{0.string}", false},
			Out: []interface{}{false, nil},
		},
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "result.some")
	bwtesting.BwRunTests(t, pfa.setVarValTestHelper, tests)
}

func (pfa *pfaStruct) setVarValTestHelper(varPathStr string, varVal interface{}) (interface{}, error) {
	varPath := MustVarPathFrom(varPathStr)
	pfa.err = nil
	pfa.setVarVal(varPath, varVal)
	if pfa.err == nil {
		return pfa.getVarValue(varPath).val, pfa.err
	} else {
		return nil, pfa.err
	}
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
			In:  []interface{}{"primary", SetVar{"primary", "end"}},
			Out: []interface{}{"end", nil},
		},
		"0.type": {
			In:  []interface{}{"0.type", SetVarBy{"0.type", Var{"1.type"}, By{Append{}}}},
			Out: []interface{}{"keymap", nil},
		},
		"0.array": {
			In:  []interface{}{"0.array", SetVarBy{"0.array", Var{"1.value.{0.string}"}, By{Append{}}}},
			Out: []interface{}{[]interface{}{"c", "d", true}, nil},
		},
		"array": {
			In:  []interface{}{"array", SetVarBy{"array", Var{"1.array"}, By{Append{}}}},
			Out: []interface{}{[]interface{}{"i", "j", []interface{}{"g", "h"}}, nil},
		},
		"1.value.arrayKey": {
			In:  []interface{}{"1.value.arrayKey", SetVarBy{"1.value.arrayKey", Var{"0.value.arrayKey"}, By{AppendSlice{}}}},
			Out: []interface{}{[]interface{}{"a", "b", "e", "f"}, nil},
		},
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "primary")
	bwtesting.BwRunTests(t, pfa.pfaActionsTestHelper, tests)
}

func (pfa *pfaStruct) pfaActionsTestHelper(varPathStr string, action interface{}) (val interface{}, err error) {
	pfa.processRules(CreateRules([]interface{}{action}))
	if pfa.err != nil {
		return nil, pfa.err
	} else {
		return pfa.getVarValue(MustVarPathFrom(varPathStr)).val, pfa.err
	}
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
		"rune is 'e' => true": {
			In:  []interface{}{VarIs{"rune", 'e'}},
			Out: []interface{}{true, nil},
		},
		"is 'e' => true": {
			In:  []interface{}{'e'},
			Out: []interface{}{true, nil},
		},
		"is 'e' or 'o' => true": {
			In:  []interface{}{[]interface{}{'e', 'o'}},
			Out: []interface{}{true, nil},
		},
		"is 'm' or 'o' => false": {
			In:  []interface{}{[]interface{}{'m', 'o'}},
			Out: []interface{}{false, nil},
		},
		"is UnicodeLetter => true": {
			In:  []interface{}{UnicodeLetter},
			Out: []interface{}{true, nil},
		},
		"is UnicodeDigit => false": {
			In:  []interface{}{UnicodeDigit},
			Out: []interface{}{false, nil},
		},
		"rune.-1 is 'm' => true": {
			In:  []interface{}{VarIs{"rune.-1", 'm'}},
			Out: []interface{}{true, nil},
		},
		"rune.-2 is 'o' => true": {
			In:  []interface{}{VarIs{"rune.-2", 'o'}},
			Out: []interface{}{true, nil},
		},
		"rune.0 is 'e' => true": {
			In:  []interface{}{VarIs{"rune.0", 'e'}},
			Out: []interface{}{true, nil},
		},
		"rune.1 is 't' => true": {
			In:  []interface{}{VarIs{"rune.1", 't'}},
			Out: []interface{}{true, nil},
		},
		"rune.2 is 'h' => true": {
			In:  []interface{}{VarIs{"rune.2", 'h'}},
			Out: []interface{}{true, nil},
		},
		"rune.3 is 'i' => true": {
			In:  []interface{}{VarIs{"rune.3", 'i'}},
			Out: []interface{}{true, nil},
		},
		"rune.5 is 'g' => true": {
			In:  []interface{}{VarIs{"rune.5", 'g'}},
			Out: []interface{}{true, nil},
		},
		"rune.5 is UnicodeLetter => true": {
			In:  []interface{}{VarIs{"rune.5", UnicodeLetter}},
			Out: []interface{}{true, nil},
		},
		"rune.5 is UnicodeDigit => false": {
			In:  []interface{}{VarIs{"rune.5", UnicodeDigit}},
			Out: []interface{}{false, nil},
		},
		"rune.6 is EOF => true": {
			In:  []interface{}{VarIs{"rune.6", EOF{}}},
			Out: []interface{}{true, nil},
		},
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "rune.-2 is 'o' => true")
	// bwerror.Spew.Printf("%#q\n", pfa.getVarValue(MustVarPathFrom("rune.-2")).val)
	bwtesting.BwRunTests(t, pfa.pfaConditionsTestHelper, tests)
}

func (pfa *pfaStruct) pfaConditionsTestHelper(arg interface{}) (interface{}, error) {
	pfa.vars["result"] = false
	var args []interface{}
	var ok bool
	if args, ok = arg.([]interface{}); !ok {
		args = []interface{}{arg}
	}
	args = append(args, SetVar{"result", true})
	pfa.processRules(CreateRules(args))
	return pfa.vars["result"], pfa.err
}
