package bwval_test

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
