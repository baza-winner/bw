package bwval_test

// func TestValKindMarshalJSON(t *testing.T) {
// 	var m = map[bwval.ValKind]string{}
// 	for vk := bwval.ValUnknown; vk <= bwval.ValArray; vk++ {
// 		m[vk] = vk.String()
// 	}
// 	for k, v := range m {
// 		bwtesting.BwRunTests(t,
// 			ValKindPretty,
// 			map[string]bwtesting.Case{
// 				v: {
// 					In:  []interface{}{k},
// 					Out: []interface{}{fmt.Sprintf("%q", v)},
// 				},
// 			},
// 		)
// 	}
// }

// func ValKindPretty(vk bwval.ValKind) string {
// 	return bwjson.Pretty(vk)
// }

// func TestKind(t *testing.T) {

// 	bwtesting.BwRunTests(t,
// 		bwval.Kind,

// 		map[string]bwtesting.Case{
// 			"Nil": {
// 				In: []interface{}{nil},
// 				Out: []interface{}{
// 					func(test bwtesting.Case) interface{} { return test.In[0] },
// 					bwval.ValNil,
// 				},
// 			},
// 			"Bool": {
// 				In: []interface{}{true},
// 				Out: []interface{}{
// 					func(test bwtesting.Case) interface{} { return test.In[0] },
// 					bwval.ValBool,
// 				},
// 			},
// 			"String": {
// 				In: []interface{}{"some"},
// 				Out: []interface{}{
// 					func(test bwtesting.Case) interface{} { return test.In[0] },
// 					bwval.ValString,
// 				},
// 			},
// 			"Map": {
// 				In: []interface{}{map[string]interface{}{}},
// 				Out: []interface{}{
// 					func(test bwtesting.Case) interface{} { return test.In[0] },
// 					bwval.ValMap,
// 				},
// 			},
// 			"Array": {
// 				In: []interface{}{[]interface{}{}},
// 				Out: []interface{}{
// 					func(test bwtesting.Case) interface{} { return test.In[0] },
// 					bwval.ValArray,
// 				},
// 			},
// 			"Int(int8)": {
// 				In: []interface{}{bw.MaxInt8},
// 				Out: []interface{}{
// 					int(bw.MaxInt8),
// 					bwval.ValInt,
// 				},
// 			},
// 			"Int(int16)": {
// 				In: []interface{}{bw.MaxInt16},
// 				Out: []interface{}{
// 					int(bw.MaxInt16),
// 					bwval.ValInt,
// 				},
// 			},
// 			"Int(int32)": {
// 				In: []interface{}{bw.MaxInt32},
// 				Out: []interface{}{
// 					int(bw.MaxInt32),
// 					bwval.ValInt,
// 				},
// 			},
// 			"Int(int64(bw.MaxInt32))": {
// 				In: []interface{}{int64(bw.MaxInt32)},
// 				Out: []interface{}{
// 					int(bw.MaxInt32),
// 					bwval.ValInt,
// 				},
// 			},
// 			// "Int(bw.MaxInt64)": {
// 			// 	In: []interface{}{int64(bw.MaxInt64)},
// 			// 	Out: []interface{}{
// 			// 		func(test bwtesting.Case) interface{} { return test.In[0] },
// 			// 		func() (result bwval.ValKind) {
// 			// 			if bw.MaxInt64 > int64(bw.MaxInt) {
// 			// 				result = bwval.ValUnknown
// 			// 			} else {
// 			// 				result = bwval.ValInt
// 			// 			}
// 			// 			return
// 			// 		},
// 			// 	},
// 			// },
// 			"Int(int)": {
// 				In: []interface{}{bw.MaxInt},
// 				Out: []interface{}{
// 					bw.MaxInt,
// 					bwval.ValInt,
// 				},
// 			},
// 			"Int(uint8)": {
// 				In: []interface{}{bw.MaxUint8},
// 				Out: []interface{}{
// 					int(bw.MaxUint8),
// 					bwval.ValInt,
// 				},
// 			},
// 			"Int(uint16)": {
// 				In: []interface{}{bw.MaxUint16},
// 				Out: []interface{}{
// 					int(bw.MaxUint16),
// 					bwval.ValInt,
// 				},
// 			},
// 			"Int(uint32)": {
// 				In: []interface{}{bw.MaxUint32},
// 				Out: []interface{}{
// 					int(bw.MaxUint32),
// 					bwval.ValInt,
// 				},
// 			},
// 			"Int(uint64(bw.MaxInt32))": {
// 				In: []interface{}{uint64(bw.MaxUint32)},
// 				Out: []interface{}{
// 					int(bw.MaxUint32),
// 					bwval.ValInt,
// 				},
// 			},
// 			// "Int(bw.MaxUint64)": {
// 			// 	In: []interface{}{bw.MaxUint64},
// 			// 	Out: []interface{}{
// 			// 		func(test bwtesting.Case) (result interface{}) {
// 			// 			if bw.MaxUint64 <= uint64(bw.MaxInt) {
// 			// 				result = test.In[0]
// 			// 			}
// 			// 			return
// 			// 		},
// 			// 		func() (result bwval.ValKind) {
// 			// 			if bw.MaxUint64 > uint64(bw.MaxInt) {
// 			// 				result = bwval.ValUnknown
// 			// 			} else {
// 			// 				result = bwval.ValInt
// 			// 			}
// 			// 			return
// 			// 		},
// 			// 	},
// 			// },
// 			"Int(bw.MaxUint)": {
// 				In: []interface{}{bw.MaxUint},
// 				Out: []interface{}{
// 					nil,
// 					bwval.ValUnknown,
// 				},
// 			},
// 			"Int(uint(0))": {
// 				In: []interface{}{uint(0)},
// 				Out: []interface{}{
// 					int(0),
// 					bwval.ValInt,
// 				},
// 			},
// 		},
// 		// tests,
// 		// nil,
// 	)
// }
