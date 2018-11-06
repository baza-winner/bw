package bwval_test

import (
	"fmt"
	"testing"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwval"
)

func TestPathVal(t *testing.T) {
	for _, test := range []struct {
		name string
		val  interface{}
		path bw.ValPath
		vars map[string]interface{}
		eta  interface{}
		err  error
	}{
		{
			name: "nil",
			val:  nil,
			path: bwval.MustPathFrom("some"),
			eta:  nil,
		},
		{
			name: "self",
			val:  1,
			path: bwval.MustPathFrom(""),
			eta:  1,
		},
		{
			name: "by key",
			val:  map[string]interface{}{"some": "thing"},
			path: bwval.MustPathFrom("some"),
			eta:  "thing",
		},
		{
			name: "by idx (1)",
			val:  []interface{}{"some", "thing"},
			path: bwval.MustPathFrom("1"),
			eta:  "thing",
		},
		{
			name: "by idx (-1)",
			val:  []interface{}{"some", "thing"},
			path: bwval.MustPathFrom("-1"),
			eta:  "thing",
		},
		{
			name: "by idx (len)",
			val:  []interface{}{"some", "thing"},
			path: bwval.MustPathFrom("2"),
			eta:  nil,
		},
		{
			name: "by idx of nil",
			val:  nil,
			path: bwval.MustPathFrom("2"),
			eta:  nil,
		},
		{
			name: "by # of nil",
			val:  nil,
			path: bwval.MustPathFrom("1.#"),
			eta:  0,
		},
		{
			name: "by # of Array",
			val:  []interface{}{"a", "b"},
			path: bwval.MustPathFrom("#"),
			eta:  2,
		},
		{
			name: "by # of Map",
			val: []interface{}{
				"a",
				map[string]interface{}{"c": "d", "e": "f", "i": "g"},
			},
			path: bwval.MustPathFrom("1.#"),
			eta:  3,
		},
		{
			name: "by # of nil",
			val:  nil,
			path: bwval.MustPathFrom("1.#"),
			eta:  0,
		},
		{
			name: "by path (idx)",
			val:  map[string]interface{}{"some": []interface{}{"good", "thing"}, "idx": 1},
			path: bwval.MustPathFrom("some.{idx}"),
			eta:  "thing",
		},
		{
			name: "by path (key)",
			val:  map[string]interface{}{"some": []interface{}{"good", "thing"}, "key": "some"},
			path: bwval.MustPathFrom("{key}.1"),
			eta:  "thing",
		},
		{
			name: "by var",
			val:  map[string]interface{}{"some": []interface{}{"good", "thing"}},
			vars: map[string]interface{}{"idx": 1},
			path: bwval.MustPathFrom("some.{$idx}"),
			eta:  "thing",
		},
		{
			name: "by var",
			val:  map[string]interface{}{"some": []interface{}{"good", "thing"}},
			vars: nil,
			path: bwval.MustPathFrom("$idx"),
			eta:  nil,
		},
		{
			name: "err: is not Map",
			val:  1,
			path: bwval.MustPathFrom("some.{$idx}"),
			err:  bwerr.Error{S: "\x1b[96;1m(int)1\x1b[0m::\x1b[38;5;252;1msome\x1b[0m (\x1b[96;1m(int)1\x1b[0m) is not \x1b[97;1mMap\x1b[0m"},
		},
		{
			name: "err: is not Array",
			val:  1,
			path: bwval.MustPathFrom("{$idx}"),
			vars: map[string]interface{}{"idx": 1},
			err:  bwerr.Error{S: "\x1b[96;1m(int)1\x1b[0m::\x1b[38;5;252;1m{$idx}\x1b[0m (\x1b[96;1m(int)1\x1b[0m) is not \x1b[97;1mArray\x1b[0m"},
		},
		{
			name: "err: nor Map, neither Array",
			val:  1,
			path: bwval.MustPathFrom("#"),
			err:  bwerr.Error{S: "\x1b[96;1m(int)1\x1b[0m::\x1b[38;5;252;1m.#\x1b[0m (\x1b[96;1m(int)1\x1b[0m) is nor \x1b[97;1mMap\x1b[0m, neither \x1b[97;1mArray\x1b[0m"},
		},
		{
			name: "err: nor Int, neither String",
			val:  1,
			path: bwval.MustPathFrom("{$idx}"),
			vars: map[string]interface{}{"idx": nil},
			err:  bwerr.Error{S: "\x1b[96;1m(int)1\x1b[0m::\x1b[38;5;252;1m.{$idx}\x1b[0m (\x1b[96;1m(interface {})<nil>\x1b[0m) is nor \x1b[97;1mInt\x1b[0m, neither \x1b[97;1mString\x1b[0m"},
		},
	} {
		crop := false
		// crop = true
		if crop && test.name != "by # of Map" {
			continue
		}
		bwtesting.BwRunTests(t,
			bwval.From(test.val).PathVal,
			map[string]bwtesting.TestCaseStruct{
				test.name: {
					In: []interface{}{
						test.path,
						test.vars,
					},
					Out: []interface{}{
						test.eta,
						test.err,
					},
				},
			},
		)
	}
}

func TestKind(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"Nil": {
			In:  []interface{}{nil},
			Out: []interface{}{bwval.ValNil},
		},
		"Nool": {
			In:  []interface{}{true},
			Out: []interface{}{bwval.ValBool},
		},
		"String": {
			In:  []interface{}{"some"},
			Out: []interface{}{bwval.ValString},
		},
		"Map": {
			In:  []interface{}{map[string]interface{}{}},
			Out: []interface{}{bwval.ValMap},
		},
		"Array": {
			In:  []interface{}{[]interface{}{}},
			Out: []interface{}{bwval.ValArray},
		},
		"Int(int8)": {
			In:  []interface{}{bw.MaxInt8},
			Out: []interface{}{bwval.ValInt},
		},
		"Int(int16)": {
			In:  []interface{}{bw.MaxInt16},
			Out: []interface{}{bwval.ValInt},
		},
		"Int(int32)": {
			In:  []interface{}{bw.MaxInt32},
			Out: []interface{}{bwval.ValInt},
		},
		"Int(int64(bw.MaxInt32))": {
			In:  []interface{}{int64(bw.MaxInt32)},
			Out: []interface{}{bwval.ValInt},
		},
		"Int(bw.MaxInt64)": {
			In: []interface{}{int64(bw.MaxInt64)},
			Out: []interface{}{
				func() (result bwval.ValKind) {
					if bw.MaxInt64 > int64(bw.MaxInt) {
						result = bwval.ValUnknown
					} else {
						result = bwval.ValInt
					}
					return
				},
			},
		},
		"Int(int)": {
			In:  []interface{}{bw.MaxInt},
			Out: []interface{}{bwval.ValInt},
		},
		"Int(uint8)": {
			In:  []interface{}{bw.MaxUint8},
			Out: []interface{}{bwval.ValInt},
		},
		"Int(uint16)": {
			In:  []interface{}{bw.MaxUint16},
			Out: []interface{}{bwval.ValInt},
		},
		"Int(uint32)": {
			In:  []interface{}{bw.MaxUint32},
			Out: []interface{}{bwval.ValInt},
		},
		"Int(uint64(bw.MaxInt32))": {
			In:  []interface{}{uint64(bw.MaxUint32)},
			Out: []interface{}{bwval.ValInt},
		},
		"Int(bw.MaxUint64)": {
			In: []interface{}{bw.MaxUint64},
			Out: []interface{}{
				func() (result bwval.ValKind) {
					if bw.MaxUint64 > uint64(bw.MaxInt) {
						result = bwval.ValUnknown
					} else {
						result = bwval.ValInt
					}
					return
				},
			},
		},
		"Int(bw.MaxUint)": {
			In: []interface{}{bw.MaxUint},
			Out: []interface{}{
				func() (result bwval.ValKind) {
					if bw.MaxUint > uint(bw.MaxInt) {
						result = bwval.ValUnknown
					} else {
						result = bwval.ValInt
					}
					return
				},
			},
		},
		"Int(uint(0))": {
			In:  []interface{}{uint(0)},
			Out: []interface{}{bwval.ValInt},
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, bwval.Kind, tests)
}

func ValKindPretty(vk bwval.ValKind) string {
	return bwjson.Pretty(vk)
}

func TestValKindMarshalJSON(t *testing.T) {
	for i := bwval.ValUnknown; i <= bwval.ValArray; i++ {
		bwtesting.BwRunTests(t,
			ValKindPretty,
			map[string]bwtesting.TestCaseStruct{
				i.String(): {
					In:  []interface{}{i},
					Out: []interface{}{fmt.Sprintf("%q", i.String())},
				},
			},
		)
	}
	// for _, test := range []struct {
	// 	name string
	// 	val  interface{}
	// 	path bw.ValPath
	// 	vars map[string]interface{}
	// 	eta  interface{}
	// 	err  error
	// }{} {
	// 	crop := false
	// 	// crop = true
	// 	if crop && test.name != "by # of Map" {
	// 		continue
	// 	}
	// 	bwtesting.BwRunTests(t,
	// 		bwval.From(test.val).PathVal,
	// 		map[string]bwtesting.TestCaseStruct{
	// 			test.name: {
	// 				In: []interface{}{
	// 					test.path,
	// 					test.vars,
	// 				},
	// 				Out: []interface{}{
	// 					test.eta,
	// 					test.err,
	// 				},
	// 			},
	// 		},
	// 	)
	// }
	// tests := map[string]bwtesting.TestCaseStruct{
	//   "Nil": {
	//     In:  []interface{}{nil},
	//     Out: []interface{}{bwval.ValNil},
	//   },
	// }

	// bwmap.CropMap(tests)
	// // bwmap.CropMap(tests, "UnexpectedItem")
	// bwtesting.BwRunTests(t, bwval.Kind, tests)
}

// func TestPathFrom(t *testing.T) {
// 	tests := map[string]bwtesting.TestCaseStruct{
// 		"some.thing": {
// 			In: []interface{}{
// 				func(testName string) string { return testName },
// 			},
// 			Out: []interface{}{
// 				Path{[]pathItem{
// 					{Type: pathItemKey, Key: "some"},
// 					{Type: pathItemKey, Key: "thing"},
// 				}},
// 				nil,
// 			},
// 		},
// 		"some.1": {
// 			In: []interface{}{
// 				func(testName string) string { return testName },
// 			},
// 			Out: []interface{}{
// 				Path{[]pathItem{
// 					{Type: pathItemKey, Key: "some"},
// 					{Type: pathItemIdx, Idx: 1},
// 				}},
// 				nil,
// 			},
// 		},
// 		"some.#": {
// 			In: []interface{}{
// 				func(testName string) string { return testName },
// 			},
// 			Out: []interface{}{
// 				Path{[]pathItem{
// 					{Type: pathItemKey, Key: "some"},
// 					{Type: pathItemHash},
// 				}},
// 				nil,
// 			},
// 		},
// 		"{some.thing}.good": {
// 			In: []interface{}{
// 				func(testName string) string { return testName },
// 			},
// 			Out: []interface{}{
// 				Path{[]pathItem{
// 					{Type: pathItemPath,
// 						Path: Path{[]pathItem{
// 							{Type: pathItemKey, Key: "some"},
// 							{Type: pathItemKey, Key: "thing"},
// 						}},
// 					},
// 					{Type: pathItemKey, Key: "good"},
// 				}},
// 				nil,
// 			},
// 		},
// 		"{$some.thing}.good": {
// 			In: []interface{}{
// 				func(testName string) string { return testName },
// 			},
// 			Out: []interface{}{
// 				Path{[]pathItem{
// 					{Type: pathItemPath,
// 						Path: Path{[]pathItem{
// 							{Type: pathItemVar, Key: "some"},
// 							{Type: pathItemKey, Key: "thing"},
// 						}},
// 					},
// 					{Type: pathItemKey, Key: "good"},
// 				}},
// 				nil,
// 			},
// 		},
// 	}

// 	bwmap.CropMap(tests)
// 	// bwmap.CropMap(tests, "UnexpectedItem")
// 	bwtesting.BwRunTests(t, PathFrom, tests)
// }

// func TestPathString(t *testing.T) {
// 	for _, test := range []struct {
// 		eta string
// 		v   Path
// 	}{
// 		{
// 			"some.thing",
// 			Path{[]pathItem{
// 				{Type: pathItemKey, Key: "some"},
// 				{Type: pathItemKey, Key: "thing"},
// 			}},
// 		},
// 		{
// 			"some.1",
// 			Path{[]pathItem{
// 				{Type: pathItemKey, Key: "some"},
// 				{Type: pathItemIdx, Idx: 1},
// 			}},
// 		},
// 		{
// 			"some.#",
// 			Path{[]pathItem{
// 				{Type: pathItemKey, Key: "some"},
// 				{Type: pathItemHash},
// 			}},
// 		},
// 		{
// 			"{some.thing}.good",
// 			Path{[]pathItem{
// 				{Type: pathItemPath,
// 					Path: Path{[]pathItem{
// 						{Type: pathItemKey, Key: "some"},
// 						{Type: pathItemKey, Key: "thing"},
// 					}},
// 				},
// 				{Type: pathItemKey, Key: "good"},
// 			}},
// 		},
// 		{
// 			"{$some.thing}.good",
// 			Path{[]pathItem{
// 				{Type: pathItemPath,
// 					Path: Path{[]pathItem{
// 						{Type: pathItemVar, Key: "some"},
// 						{Type: pathItemKey, Key: "thing"},
// 					}},
// 				},
// 				{Type: pathItemKey, Key: "good"},
// 			}},
// 		},
// 	} {
// 		bwtesting.BwRunTests(t,
// 			test.v.String,
// 			map[string]bwtesting.TestCaseStruct{
// 				test.eta: {
// 					In:  []interface{}{},
// 					Out: []interface{}{test.eta},
// 				},
// 			},
// 		)
// 	}
// }
