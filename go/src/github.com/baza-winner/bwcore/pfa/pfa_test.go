package pfa

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/pfa/a"
	"github.com/baza-winner/bwcore/pfa/b"
	"github.com/baza-winner/bwcore/pfa/c"
	"github.com/baza-winner/bwcore/pfa/core"
	"github.com/baza-winner/bwcore/pfa/r"
	"github.com/baza-winner/bwcore/pfa/val"
	"github.com/baza-winner/bwcore/runeprovider"
)

func getTestFileSpec(basename string) string {
	return filepath.Join(os.Getenv("GOPATH"), "src", "github.com/baza-winner/bwcore/pfa", basename)
}

func TestParseLogic(t *testing.T) {
	p, err := runeprovider.FromFile(getTestFileSpec("logic.pfa"))
	if err != nil {
		bwerror.PanicErr(err)
	}
	tests := map[string]bwtesting.TestCaseStruct{
		"": {
			In:  []interface{}{p},
			Out: []interface{}{nil, nil},
		},
	}
	testsToRun := tests
	bwmap.CropMap(testsToRun)
	// bwmap.CropMap(testsToRun, "qw ")
	// bwmap.CropMap(testsToRun, "qw && fa.curr.runePtr == EOF")
	bwtesting.BwRunTests(t, ParseLogic, testsToRun)
}

func TestVarPathFrom(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"": {
			In:  []interface{}{""},
			Out: []interface{}{[]interface{}{}, nil},
		},
		"0": {
			In:  []interface{}{"0"},
			Out: []interface{}{core.VarPath{core.VarPathItem{0}}, nil},
		},
		"0.key": {
			In:  []interface{}{"0.key"},
			Out: []interface{}{core.VarPath{core.VarPathItem{0}, core.VarPathItem{"key"}}, nil},
		},
		"0.key.{1.some}": {
			In: []interface{}{"0.key.{1.some}"},
			Out: []interface{}{
				core.VarPath{
					core.VarPathItem{0},
					core.VarPathItem{"key"},
					core.VarPathItem{core.VarPath{core.VarPathItem{1}, core.VarPathItem{"some"}}},
				},
				nil,
			},
		},
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "qw ")
	bwtesting.BwRunTests(t, core.VarPathFrom, tests)
}

func TestPfa_getVarValue(t *testing.T) {
	pfa := core.PfaFrom(runeprovider.FromString("some"), core.TraceNone)
	pfa.Stack = core.ParseStack{
		core.ParseStackItem{
			Vars: map[string]interface{}{
				"type": "map",
				"value": map[string]interface{}{
					"boolKey":   true,
					"numberKey": 273,
					"stringKey": "string value",
					"runeKey":   '\n',
				},
			},
		},
		core.ParseStackItem{
			Vars: map[string]interface{}{
				"type":   "key",
				"string": "boolKey",
			},
		},
	}
	pfa.Vars = map[string]interface{}{
		"primary": "begin",
		"result": map[string]interface{}{
			"some": "thing",
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
	bwtesting.BwRunTests(t, TestHelper{pfa}.VarValue, tests)
}

type TestHelper struct {
	pfa *core.PfaStruct
}

func (v TestHelper) VarValue(varPathStr string) (val interface{}, err error) {
	v.pfa.Err = nil
	varValue := v.pfa.VarValue(core.MustVarPathFrom(varPathStr))
	return varValue.Val, v.pfa.Err
}

func TestPfa_setVarVal(t *testing.T) {
	pfa := core.PfaFrom(runeprovider.FromString("some"), core.TraceNone)
	pfa.Stack = core.ParseStack{
		core.ParseStackItem{
			Vars: map[string]interface{}{
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
		core.ParseStackItem{
			Vars: map[string]interface{}{
				"type":   "key",
				"string": "boolKey",
				"array":  []interface{}{"c", "d"},
			},
		},
	}
	pfa.Vars = map[string]interface{}{
		"primary": "begin",
		"result": map[string]interface{}{
			"some": "thing",
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
	bwtesting.BwRunTests(t, TestHelper{pfa}.SetVarValue, tests)
}

func (v TestHelper) SetVarValue(varPathStr string, VarVal interface{}) (interface{}, error) {
	varPath := core.MustVarPathFrom(varPathStr)
	v.pfa.Err = nil
	v.pfa.SetVarVal(varPath, VarVal)
	if v.pfa.Err == nil {
		return v.pfa.VarValue(varPath).Val, v.pfa.Err
	} else {
		return nil, v.pfa.Err
	}
}

func TestPfaActions(t *testing.T) {
	pfa := core.PfaFrom(runeprovider.FromString("some"), core.TraceNone)
	pfa.Stack = core.ParseStack{
		core.ParseStackItem{
			Vars: map[string]interface{}{
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
		core.ParseStackItem{
			Vars: map[string]interface{}{
				"type":   "key",
				"string": "boolKey",
				"array":  []interface{}{"c", "d"},
				"value": map[string]interface{}{
					"arrayKey": []interface{}{"e", "f"},
				},
			},
		},
	}
	pfa.Vars = map[string]interface{}{
		"array":   []interface{}{"i", "j"},
		"primary": "begin",
		"result": map[string]interface{}{
			"some": "thing",
		},
	}
	tests := map[string]bwtesting.TestCaseStruct{
		"primary": {
			In:  []interface{}{"primary", a.SetVar{"primary", "end"}},
			Out: []interface{}{"end", nil},
		},
		"0.type": {
			In:  []interface{}{"0.type", a.SetVarBy{"0.type", val.Var{"1.type"}, b.By{b.Append{}}}},
			Out: []interface{}{"keymap", nil},
		},
		"0.array": {
			In:  []interface{}{"0.array", a.SetVarBy{"0.array", val.Var{"1.value.{0.string}"}, b.By{b.Append{}}}},
			Out: []interface{}{[]interface{}{"c", "d", true}, nil},
		},
		"array": {
			In:  []interface{}{"array", a.SetVarBy{"array", val.Var{"1.array"}, b.By{b.Append{}}}},
			Out: []interface{}{[]interface{}{"i", "j", []interface{}{"g", "h"}}, nil},
		},
		"1.value.arrayKey": {
			In:  []interface{}{"1.value.arrayKey", a.SetVarBy{"1.value.arrayKey", val.Var{"0.value.arrayKey"}, b.By{b.AppendSlice{}}}},
			Out: []interface{}{[]interface{}{"a", "b", "e", "f"}, nil},
		},
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "primary")
	bwtesting.BwRunTests(t, TestHelper{pfa}.Action, tests)
}

func (v TestHelper) Action(varPathStr string, action interface{}) (val interface{}, err error) {
	r.RulesFrom([]interface{}{action}).Process(v.pfa)
	if v.pfa.Err != nil {
		return nil, v.pfa.Err
	} else {
		return v.pfa.VarValue(core.MustVarPathFrom(varPathStr)).Val, v.pfa.Err
	}
}

func TestPfaConditions(t *testing.T) {
	pfa := core.PfaFrom(runeprovider.FromString("something"), core.TraceNone)
	pfa.Proxy.PullRune()
	pfa.Proxy.PullRune()
	pfa.Proxy.PullRune()
	pfa.Proxy.PullRune()
	pfa.Stack = core.ParseStack{
		core.ParseStackItem{
			Vars: map[string]interface{}{
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
		core.ParseStackItem{
			Vars: map[string]interface{}{
				"type":   "key",
				"string": "boolKey",
				"array":  []interface{}{"c", "d"},
				"value": map[string]interface{}{
					"arrayKey": []interface{}{"e", "f"},
				},
			},
		},
	}
	// p,
	// nil,
	pfa.Vars = map[string]interface{}{
		"array":   []interface{}{"i", "j"},
		"primary": "begin",
		"result": map[string]interface{}{
			"some": "thing",
		},
	}
	// 	core.TraceNone,
	// 	nil,
	// 	0,
	// }
	tests := map[string]bwtesting.TestCaseStruct{
		"primary is begin => true": {
			In:  []interface{}{c.VarIs{"primary", "begin"}},
			Out: []interface{}{true, nil},
		},
		"primary is end => false": {
			In:  []interface{}{c.VarIs{"primary", "end"}},
			Out: []interface{}{false, nil},
		},
		"rune is 'e' => true": {
			In:  []interface{}{c.VarIs{"rune", 'e'}},
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
		"is Letter => true": {
			In:  []interface{}{val.Letter},
			Out: []interface{}{true, nil},
		},
		"is Digit => false": {
			In:  []interface{}{val.Digit},
			Out: []interface{}{false, nil},
		},
		"rune.-1 is 'm' => true": {
			In:  []interface{}{c.VarIs{"rune.-1", 'm'}},
			Out: []interface{}{true, nil},
		},
		"rune.-2 is 'o' => true": {
			In:  []interface{}{c.VarIs{"rune.-2", 'o'}},
			Out: []interface{}{true, nil},
		},
		"rune.0 is 'e' => true": {
			In:  []interface{}{c.VarIs{"rune.0", 'e'}},
			Out: []interface{}{true, nil},
		},
		"rune.1 is 't' => true": {
			In:  []interface{}{c.VarIs{"rune.1", 't'}},
			Out: []interface{}{true, nil},
		},
		"rune.2 is 'h' => true": {
			In:  []interface{}{c.VarIs{"rune.2", 'h'}},
			Out: []interface{}{true, nil},
		},
		"rune.3 is 'i' => true": {
			In:  []interface{}{c.VarIs{"rune.3", 'i'}},
			Out: []interface{}{true, nil},
		},
		"rune.5 is 'g' => true": {
			In:  []interface{}{c.VarIs{"rune.5", 'g'}},
			Out: []interface{}{true, nil},
		},
		"rune.5 is Letter => true": {
			In:  []interface{}{c.VarIs{"rune.5", val.Letter}},
			Out: []interface{}{true, nil},
		},
		"rune.5 is Digit => false": {
			In:  []interface{}{c.VarIs{"rune.5", val.Digit}},
			Out: []interface{}{false, nil},
		},
		"rune.6 is EOF => true": {
			In:  []interface{}{c.VarIs{"rune.6", val.EOF{}}},
			Out: []interface{}{true, nil},
		},
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "rune.-2 is 'o' => true")
	// bwerror.Spew.Printf("%#q\n", pfa.VarValue(core.MustVarPathFrom("rune.-2")).val)
	bwtesting.BwRunTests(t, TestHelper{pfa}.Check, tests)
}

func (v TestHelper) Check(arg interface{}) (interface{}, error) {
	v.pfa.Vars["result"] = false
	var args []interface{}
	var ok bool
	if args, ok = arg.([]interface{}); !ok {
		args = []interface{}{arg}
	}
	args = append(args, a.SetVar{"result", true})
	r.RulesFrom(args).Process(v.pfa)
	return v.pfa.Vars["result"], v.pfa.Err
}
