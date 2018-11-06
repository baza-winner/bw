package path_test

import (
	"testing"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwval/path"
)

func TestParse(t *testing.T) {
	tests := map[string]bwtesting.TestCaseStruct{
		"": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				bw.ValPath{},
				nil,
			},
		},
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
		"$some.thing.{good}": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				bw.ValPath{
					{Type: bw.ValPathItemVar, Key: "some"},
					{Type: bw.ValPathItemKey, Key: "thing"},
					{Type: bw.ValPathItemPath,
						Path: bw.ValPath{
							{Type: bw.ValPathItemKey, Key: "good"},
						},
					},
				},
				nil,
			},
		},
		"1.some": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				bw.ValPath{
					{Type: bw.ValPathItemIdx, Idx: 1},
					{Type: bw.ValPathItemKey, Key: "some"},
				},
				nil,
			},
		},
		"-1.some": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				bw.ValPath{
					{Type: bw.ValPathItemIdx, Idx: -1},
					{Type: bw.ValPathItemKey, Key: "some"},
				},
				nil,
			},
		},
		"1.": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				bw.ValPath{},
				bwerr.Error{S: "unexpected end of string at pos \x1b[38;5;252;1m2\x1b[0m: \x1b[32m1.\n"},
			},
		},
		"1.@": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				bw.ValPath{},
				bwerr.Error{S: "unexpected char \u001b[96;1m'@'\u001b[0m (\u001b[38;5;201;1mcharCode\u001b[0m: \u001b[96;1m64\u001b[0m)\u001b[0m at pos \u001b[38;5;252;1m2\u001b[0m: \u001b[32m1.\u001b[91m@\u001b[0m\n"},
			},
		},
		"-a": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				bw.ValPath{},
				bwerr.Error{S: "unexpected char \u001b[96;1m'a'\u001b[0m (\u001b[38;5;201;1mcharCode\u001b[0m: \u001b[96;1m97\u001b[0m)\u001b[0m at pos \u001b[38;5;252;1m1\u001b[0m: \u001b[32m-\u001b[91ma\u001b[0m\n"},
			},
		},
		"1a": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				bw.ValPath{},
				bwerr.Error{S: "unexpected char \u001b[96;1m'a'\u001b[0m (\u001b[38;5;201;1mcharCode\u001b[0m: \u001b[96;1m97\u001b[0m)\u001b[0m at pos \u001b[38;5;252;1m1\u001b[0m: \u001b[32m1\u001b[91ma\u001b[0m\n"},
			},
		},
		"12.#.4": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				bw.ValPath{},
				bwerr.Error{S: "unexpected char \u001b[96;1m'.'\u001b[0m (\u001b[38;5;201;1mcharCode\u001b[0m: \u001b[96;1m46\u001b[0m)\u001b[0m at pos \u001b[38;5;252;1m4\u001b[0m: \u001b[32m12.#\u001b[91m.\u001b[0m4\n"},
			},
		},
		"12.{4": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				bw.ValPath{},
				bwerr.Error{S: "unexpected end of string at pos \u001b[38;5;252;1m5\u001b[0m: \u001b[32m12.{4\n"},
			},
		},
		"12.$a": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				bw.ValPath{},
				bwerr.Error{S: "unexpected char \u001b[96;1m'$'\u001b[0m (\u001b[38;5;201;1mcharCode\u001b[0m: \u001b[96;1m36\u001b[0m)\u001b[0m at pos \u001b[38;5;252;1m3\u001b[0m: \u001b[32m12.\u001b[91m$\u001b[0ma\n"},
			},
		},
		"$12": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				// bw.ValPath{},
				bw.ValPath{
					// {Type: bw.ValPathItemIdx, Idx: -1},
					{Type: bw.ValPathItemVar, Key: "12"},
				},
				nil,
				// bwerr.Error{S: "unexpected char \u001b[96;1m'$'\u001b[0m (\u001b[38;5;201;1mcharCode\u001b[0m: \u001b[96;1m36\u001b[0m)\u001b[0m at pos \u001b[38;5;252;1m3\u001b[0m: \u001b[32m12.\u001b[91m$\u001b[0ma\n"},
			},
		},
		"some.{$idx}": {
			In: []interface{}{
				func(testName string) string { return testName },
			},
			Out: []interface{}{
				// bw.ValPath{},
				bw.ValPath{
					{Type: bw.ValPathItemKey, Key: "some"},
					// {Type: bw.ValPathItemKey, Key: "thing"},
					{Type: bw.ValPathItemPath,
						Path: bw.ValPath{
							{Type: bw.ValPathItemVar, Key: "idx"},
						},
					},
				},
				// bw.ValPath{
				// 	// {Type: bw.ValPathItemIdx, Idx: -1},
				// 	{Type: bw.ValPathItemVar, Key: "12"},
				// },
				nil,
				// bwerr.Error{S: "unexpected char \u001b[96;1m'$'\u001b[0m (\u001b[38;5;201;1mcharCode\u001b[0m: \u001b[96;1m36\u001b[0m)\u001b[0m at pos \u001b[38;5;252;1m3\u001b[0m: \u001b[32m12.\u001b[91m$\u001b[0ma\n"},
			},
		},
	}

	bwmap.CropMap(tests)
	// bwmap.CropMap(tests, "some.{$idx}")
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
