package path_test

import (
	"testing"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwval/path"
)

func TestParse(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"some.thing": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				bw.ValPath{
					{Type: bw.ValPathItemKey, Key: "some"},
					{Type: bw.ValPathItemKey, Key: "thing"},
				},
				nil,
			},
		},
		"some.1": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				bw.ValPath{
					{Type: bw.ValPathItemKey, Key: "some"},
					{Type: bw.ValPathItemIdx, Idx: 1},
				},
				nil,
			},
		},
		"some.#": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				bw.ValPath{
					{Type: bw.ValPathItemKey, Key: "some"},
					{Type: bw.ValPathItemHash},
				},
				nil,
			},
		},
		"{some.thing}.good": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				bw.ValPath{
					{Type: bw.ValPathItemPath,
						Path: bw.ValPath{
							{Type: bw.ValPathItemKey, Key: "some"},
							{Type: bw.ValPathItemKey, Key: "thing"},
						},
					},
					{Type: bw.ValPathItemKey, Key: "good"},
				},
				nil,
			},
		},
		"{$some.thing}.good": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				bw.ValPath{
					{Type: bw.ValPathItemPath,
						Path: bw.ValPath{
							{Type: bw.ValPathItemVar, Key: "some"},
							{Type: bw.ValPathItemKey, Key: "thing"},
						},
					},
					{Type: bw.ValPathItemKey, Key: "good"},
				},
				nil,
			},
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "UnexpectedItem")
	bwtesting.BwRunTests(t, path.Parse, tests)
}

// func TestPathString(t *testing.T) {
// 	for _, test := range []struct {
// 		eta string
// 		v   bw.ValPath
// 	}{
// 		{
// 			"some.thing",
// 			bw.ValPath{
// 				{Type: bw.ValPathItemKey, Key: "some"},
// 				{Type: bw.ValPathItemKey, Key: "thing"},
// 			},
// 		},
// 		{
// 			"some.1",
// 			bw.ValPath{
// 				{Type: bw.ValPathItemKey, Key: "some"},
// 				{Type: bw.ValPathItemIdx, Idx: 1},
// 			},
// 		},
// 		{
// 			"some.#",
// 			bw.ValPath{
// 				{Type: bw.ValPathItemKey, Key: "some"},
// 				{Type: bw.ValPathItemHash},
// 			},
// 		},
// 		{
// 			"{some.thing}.good",
// 			bw.ValPath{
// 				{Type: bw.ValPathItemPath,
// 					Path: bw.ValPath{
// 						{Type: bw.ValPathItemKey, Key: "some"},
// 						{Type: bw.ValPathItemKey, Key: "thing"},
// 					},
// 				},
// 				{Type: bw.ValPathItemKey, Key: "good"},
// 			},
// 		},
// 		{
// 			"{$some.thing}.good",
// 			bw.ValPath{
// 				{Type: bw.ValPathItemPath,
// 					Path: bw.ValPath{
// 						{Type: bw.ValPathItemVar, Key: "some"},
// 						{Type: bw.ValPathItemKey, Key: "thing"},
// 					},
// 				},
// 				{Type: bw.ValPathItemKey, Key: "good"},
// 			},
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
