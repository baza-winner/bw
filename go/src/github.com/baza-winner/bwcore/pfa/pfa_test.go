package pfa

import (
	"testing"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/pfa/a"
	"github.com/baza-winner/bwcore/pfa/b"
	"github.com/baza-winner/bwcore/pfa/c"
	"github.com/baza-winner/bwcore/pfa/core"
	"github.com/baza-winner/bwcore/pfa/d"
	"github.com/baza-winner/bwcore/pfa/r"
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
		"0.key": {
			In: []interface{}{"0.key"},
			Out: []interface{}{(core.VarPath)(nil),
				bwerror.Error(
					"unexpected char <ansiPrimary>%q<ansi> (charCode: %d, state = %s) at pos <ansiCmd>%d<ansi>: <ansiDarkGreen>%s<ansiLightRed>%s<ansi>%s\n",
					'0', '0', "vppsBegin", 0, "", "0", ".key",
				),
			},
		},
		"rune": {
			In:  []interface{}{"rune"},
			Out: []interface{}{core.VarPath{core.VarPathItem{Type: core.VarPathItemKey, Key: "rune"}}, nil},
		},
		"rune.2": {
			In: []interface{}{"rune.2"},
			Out: []interface{}{core.VarPath{
				core.VarPathItem{Type: core.VarPathItemKey, Key: "rune"},
				core.VarPathItem{Type: core.VarPathItemIdx, Idx: 2},
			}, nil},
		},
		"stack.#": {
			In: []interface{}{"stack.#"},
			Out: []interface{}{core.VarPath{
				core.VarPathItem{Type: core.VarPathItemKey, Key: "stack"},
				core.VarPathItem{Type: core.VarPathItemHash},
			}, nil},
		},
		"stack.-1.string": {
			In: []interface{}{"stack.-1.string"},
			Out: []interface{}{core.VarPath{
				core.VarPathItem{Type: core.VarPathItemKey, Key: "stack"},
				core.VarPathItem{Type: core.VarPathItemIdx, Idx: -1},
				core.VarPathItem{Type: core.VarPathItemKey, Key: "string"},
			}, nil},
		},
		// "0.key": {
		// 	In:  []interface{}{"0.key"},
		// 	Out: []interface{}{core.VarPath{core.VarPathItem{0}, core.VarPathItem{"key"}}, nil},
		// },
		// "0.key.{1.some}": {
		// 	In: []interface{}{"0.key.{1.some}"},
		// 	Out: []interface{}{
		// 		core.VarPath{
		// 			core.VarPathItem{0},
		// 			core.VarPathItem{"key"},
		// 			core.VarPathItem{core.VarPath{core.VarPathItem{1}, core.VarPathItem{"some"}}},
		// 		},
		// 		nil,
		// 	},
		// },
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "qw ")
	bwtesting.BwRunTests(t, core.VarPathFrom, tests)
}

func TestPfa_getVarValue(t *testing.T) {
	pfa := core.PfaFrom(runeprovider.FromString("some"), core.TraceNone)
	pfa.Vars = map[string]interface{}{
		"primary": "begin",
		"result": map[string]interface{}{
			"some": "thing",
		},
		"stack": []interface{}{
			map[string]interface{}{
				"type": "map",
				"value": map[string]interface{}{
					"boolKey":   true,
					"numberKey": 273,
					"stringKey": "string value",
					"runeKey":   '\n',
				},
			},
			map[string]interface{}{
				"type":   "key",
				"string": "boolKey",
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
		"stack.#": {
			In:  []interface{}{"stack.#"},
			Out: []interface{}{2, nil},
		},
		"stack.-1.type": {
			In:  []interface{}{"stack.-1.type"},
			Out: []interface{}{"key", nil},
		},
		"stack.-2.value.numberKey": {
			In:  []interface{}{"stack.-2.value.numberKey"},
			Out: []interface{}{273, nil},
		},
		"stack.-2.value.{0.string}": {
			In:  []interface{}{"stack.-2.value.{stack.-1.string}"},
			Out: []interface{}{true, nil},
		},
	}
	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "primary")
	bwtesting.BwRunTests(t, TestHelper{pfa}.VarValue, tests)
}

func TestPfa_getVarValue2(t *testing.T) {
	pfa := core.PfaFrom(runeprovider.FromString("some"), core.TraceNone)
	runeA := 'a'
	runeB := 'b'
	itemPos := runeprovider.PosStruct{
		IsEOF:       false,
		RunePtr:     &runeB,
		Pos:         25,
		Line:        0x4,
		Col:         0x3,
		Prefix:      "[\n  qw/one two three/\n  d",
		PrefixStart: 1,
	}
	pfa.Vars = map[string]interface{}{
		"skipPostProcess": false,
		"primary":         "expectWord",
		"needFinish":      true,
		"stack": []interface{}{
			map[string]interface{}{
				"itemPos": runeprovider.PosStruct{
					IsEOF: false,
					// RunePtr: (*int32)(0xc0001b6a60),
					RunePtr:     &runeA,
					Pos:         1,
					Line:        0x2,
					Col:         0x1,
					Prefix:      "\n[",
					PrefixStart: 0,
				},
				"result":    []interface{}{"one", "two", "three"},
				"type":      "array",
				"delimiter": 93,
			},
			map[string]interface{}{
				"itemPos": itemPos,
				"type":    "word",
				"string":  "def",
			},
		},
		"secondary": "",
	}
	tests := map[string]bwtesting.TestCaseStruct{
		"stack.-1.itemPos": {
			In:  []interface{}{"stack.-1.itemPos"},
			Out: []interface{}{itemPos, nil},
		},
		"retry stack.-1.itemPos": {
			In:  []interface{}{"stack.-1.itemPos"},
			Out: []interface{}{itemPos, nil},
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
	pfa.Vars = map[string]interface{}{
		"primary": "begin",
		"result": map[string]interface{}{
			"some": "thing",
		},
		"stack": []interface{}{
			map[string]interface{}{
				"type": "map",
				"value": map[string]interface{}{
					"boolKey":   true,
					"numberKey": 273,
					"stringKey": "string value",
					"runeKey":   '\n',
					"arrayKey":  []interface{}{"a", "b"},
				},
			},
			map[string]interface{}{
				"type":   "key",
				"string": "boolKey",
				"array":  []interface{}{"c", "d"},
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
		"stack.#": {
			In: []interface{}{"stack.#", 4},
			Out: []interface{}{nil,
				bwerror.Error(
					"failed to set <ansiReset><ansiCmd>stack.#<ansi>: <ansiReset><ansiOutline>path.#<ansi> is <ansiCmd>readonly<ansiReset>",
				)},
		},
		"stack.-1.type": {
			In:  []interface{}{"stack.-1.type", "map"},
			Out: []interface{}{"map", nil},
		},
		"stack.-1.item": {
			In:  []interface{}{"stack.-1.item", "word"},
			Out: []interface{}{"word", nil},
		},
		"stack.-2.value.numberKey": {
			In:  []interface{}{"stack.-2.value.numberKey", "2.71"},
			Out: []interface{}{"2.71", nil},
		},
		"stack.-2.value.{stack.-1.string}": {
			In:  []interface{}{"stack.-2.value.{stack.-1.string}", false},
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
	pfa.Vars = map[string]interface{}{
		"array":   []interface{}{"i", "j"},
		"primary": "begin",
		"result": map[string]interface{}{
			"some": "thing",
		},
		"stack": []interface{}{
			map[string]interface{}{
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
			map[string]interface{}{
				"type":   "key",
				"string": "boolKey",
				"array":  []interface{}{"c", "d"},
				"value": map[string]interface{}{
					"arrayKey": []interface{}{"e", "f"},
				},
			},
		},
	}
	tests := map[string]bwtesting.TestCaseStruct{
		"primary": {
			In:  []interface{}{"primary", a.SetVar{"primary", "end"}},
			Out: []interface{}{"end", nil},
		},
		"stack.-1.type": {
			In:  []interface{}{"stack.-1.type", a.SetVarBy{"stack.-1.type", d.Var{"stack.-2.type"}, b.By{b.Append{}}}},
			Out: []interface{}{"keymap", nil},
		},
		"stack.-1.array": {
			In:  []interface{}{"stack.-1.array", a.SetVarBy{"stack.-1.array", d.Var{"stack.-2.value.{stack.-1.string}"}, b.By{b.Append{}}}},
			Out: []interface{}{[]interface{}{"c", "d", true}, nil},
		},
		"array": {
			In:  []interface{}{"array", a.SetVarBy{"array", d.Var{"stack.-2.array"}, b.By{b.Append{}}}},
			Out: []interface{}{[]interface{}{"i", "j", []interface{}{"g", "h"}}, nil},
		},
		"stack.-2.value.arrayKey": {
			In:  []interface{}{"stack.-2.value.arrayKey", a.SetVarBy{"stack.-2.value.arrayKey", d.Var{"stack.-1.value.arrayKey"}, b.By{b.AppendSlice{}}}},
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
	pfa.Vars = map[string]interface{}{
		"array":   []interface{}{"i", "j"},
		"primary": "begin",
		"result": map[string]interface{}{
			"some": "thing",
		},
		"stack": []interface{}{
			map[string]interface{}{
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
			map[string]interface{}{
				"type":   "key",
				"string": "boolKey",
				"array":  []interface{}{"c", "d"},
				"value": map[string]interface{}{
					"arrayKey": []interface{}{"e", "f"},
				},
			},
		},
	}
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
			In:  []interface{}{d.Letter},
			Out: []interface{}{true, nil},
		},
		"is Digit => false": {
			In:  []interface{}{d.Digit},
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
			In:  []interface{}{c.VarIs{"rune.5", d.Letter}},
			Out: []interface{}{true, nil},
		},
		"rune.5 is Digit => false": {
			In:  []interface{}{c.VarIs{"rune.5", d.Digit}},
			Out: []interface{}{false, nil},
		},
		"rune.6 is EOF => true": {
			In:  []interface{}{c.VarIs{"rune.6", d.EOF{}}},
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
