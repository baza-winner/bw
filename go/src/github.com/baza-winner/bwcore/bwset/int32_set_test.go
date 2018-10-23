// Code generated by "bwsetter -type=int32"; DO NOT EDIT; bwsetter: go get -type=int32 -set=Int32 -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	bwjson "github.com/baza-winner/bwcore/bwjson"
	bwtesting "github.com/baza-winner/bwcore/bwtesting"
	"strconv"
	"testing"
)

func TestInt32(t *testing.T) {
	bwtesting.BwRunTests(t, Int32From, map[string]bwtesting.TestCaseStruct{"Int32From": {
		In: []interface{}{[]int32{_Int32TestItemA, _Int32TestItemB}},
		Out: []interface{}{Int32{
			_Int32TestItemA: struct{}{},
			_Int32TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int32FromSlice, map[string]bwtesting.TestCaseStruct{"Int32FromSlice": {
		In: []interface{}{[]int32{_Int32TestItemA, _Int32TestItemB}},
		Out: []interface{}{Int32{
			_Int32TestItemA: struct{}{},
			_Int32TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int32FromSet, map[string]bwtesting.TestCaseStruct{"Int32FromSet": {
		In: []interface{}{Int32{
			_Int32TestItemA: struct{}{},
			_Int32TestItemB: struct{}{},
		}},
		Out: []interface{}{Int32{
			_Int32TestItemA: struct{}{},
			_Int32TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int32.Copy, map[string]bwtesting.TestCaseStruct{"Int32.Copy": {
		In: []interface{}{Int32{
			_Int32TestItemA: struct{}{},
			_Int32TestItemB: struct{}{},
		}},
		Out: []interface{}{Int32{
			_Int32TestItemA: struct{}{},
			_Int32TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int32.ToSlice, map[string]bwtesting.TestCaseStruct{"Int32.ToSlice": {
		In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}},
		Out: []interface{}{[]int32{_Int32TestItemA}},
	}})
	bwtesting.BwRunTests(t, _Int32ToSliceTestHelper, map[string]bwtesting.TestCaseStruct{"_Int32ToSliceTestHelper": {
		In:  []interface{}{[]int32{_Int32TestItemB, _Int32TestItemA}},
		Out: []interface{}{[]int32{_Int32TestItemA, _Int32TestItemB}},
	}})
	bwtesting.BwRunTests(t, Int32.String, map[string]bwtesting.TestCaseStruct{"Int32.String": {
		In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}},
		Out: []interface{}{bwjson.PrettyJson([]int32{_Int32TestItemA})},
	}})
	bwtesting.BwRunTests(t, Int32.DataForJSON, map[string]bwtesting.TestCaseStruct{"Int32.DataForJSON": {
		In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}},
		Out: []interface{}{[]interface{}{_Int32TestItemA}},
	}})
	bwtesting.BwRunTests(t, Int32.ToSliceOfStrings, map[string]bwtesting.TestCaseStruct{"Int32.ToSliceOfStrings": {
		In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}},
		Out: []interface{}{[]string{strconv.FormatInt(int64(_Int32TestItemA), 10)}},
	}})
	bwtesting.BwRunTests(t, Int32.Has, map[string]bwtesting.TestCaseStruct{
		"Int32.Has: false": {
			In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}, _Int32TestItemB},
			Out: []interface{}{false},
		},
		"Int32.Has: true": {
			In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}, _Int32TestItemA},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int32.HasAny, map[string]bwtesting.TestCaseStruct{
		"Int32.HasAny: empty": {
			In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}, []int32{}},
			Out: []interface{}{false},
		},
		"Int32.HasAny: false": {
			In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}, []int32{_Int32TestItemB}},
			Out: []interface{}{false},
		},
		"Int32.HasAny: true": {
			In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}, []int32{_Int32TestItemA, _Int32TestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int32.HasAnyOfSlice, map[string]bwtesting.TestCaseStruct{
		"Int32.HasAnyOfSlice: empty": {
			In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}, []int32{}},
			Out: []interface{}{false},
		},
		"Int32.HasAnyOfSlice: false": {
			In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}, []int32{_Int32TestItemB}},
			Out: []interface{}{false},
		},
		"Int32.HasAnyOfSlice: true": {
			In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}, []int32{_Int32TestItemA, _Int32TestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int32.HasAnyOfSet, map[string]bwtesting.TestCaseStruct{
		"Int32.HasAnyOfSet: empty": {
			In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}, Int32{}},
			Out: []interface{}{false},
		},
		"Int32.HasAnyOfSet: false": {
			In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}, Int32{_Int32TestItemB: struct{}{}}},
			Out: []interface{}{false},
		},
		"Int32.HasAnyOfSet: true": {
			In: []interface{}{Int32{_Int32TestItemA: struct{}{}}, Int32{
				_Int32TestItemA: struct{}{},
				_Int32TestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int32.HasEach, map[string]bwtesting.TestCaseStruct{
		"Int32.HasEach: empty": {
			In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}, []int32{}},
			Out: []interface{}{true},
		},
		"Int32.HasEach: false": {
			In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}, []int32{_Int32TestItemA, _Int32TestItemB}},
			Out: []interface{}{false},
		},
		"Int32.HasEach: true": {
			In: []interface{}{Int32{
				_Int32TestItemA: struct{}{},
				_Int32TestItemB: struct{}{},
			}, []int32{_Int32TestItemA, _Int32TestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int32.HasEachOfSlice, map[string]bwtesting.TestCaseStruct{
		"Int32.HasEachOfSlice: empty": {
			In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}, []int32{}},
			Out: []interface{}{true},
		},
		"Int32.HasEachOfSlice: false": {
			In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}, []int32{_Int32TestItemA, _Int32TestItemB}},
			Out: []interface{}{false},
		},
		"Int32.HasEachOfSlice: true": {
			In: []interface{}{Int32{
				_Int32TestItemA: struct{}{},
				_Int32TestItemB: struct{}{},
			}, []int32{_Int32TestItemA, _Int32TestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int32.HasEachOfSet, map[string]bwtesting.TestCaseStruct{
		"Int32.HasEachOfSet: empty": {
			In:  []interface{}{Int32{_Int32TestItemA: struct{}{}}, Int32{}},
			Out: []interface{}{true},
		},
		"Int32.HasEachOfSet: false": {
			In: []interface{}{Int32{_Int32TestItemA: struct{}{}}, Int32{
				_Int32TestItemA: struct{}{},
				_Int32TestItemB: struct{}{},
			}},
			Out: []interface{}{false},
		},
		"Int32.HasEachOfSet: true": {
			In: []interface{}{Int32{
				_Int32TestItemA: struct{}{},
				_Int32TestItemB: struct{}{},
			}, Int32{
				_Int32TestItemA: struct{}{},
				_Int32TestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int32._AddTestHelper, map[string]bwtesting.TestCaseStruct{"Int32.Add": {
		In: []interface{}{Int32{_Int32TestItemA: struct{}{}}, []int32{_Int32TestItemB}},
		Out: []interface{}{Int32{
			_Int32TestItemA: struct{}{},
			_Int32TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int32._AddSliceTestHelper, map[string]bwtesting.TestCaseStruct{"Int32.AddSlice": {
		In: []interface{}{Int32{_Int32TestItemA: struct{}{}}, []int32{_Int32TestItemB}},
		Out: []interface{}{Int32{
			_Int32TestItemA: struct{}{},
			_Int32TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int32._AddSetTestHelper, map[string]bwtesting.TestCaseStruct{"Int32.AddSet": {
		In: []interface{}{Int32{_Int32TestItemA: struct{}{}}, Int32{_Int32TestItemB: struct{}{}}},
		Out: []interface{}{Int32{
			_Int32TestItemA: struct{}{},
			_Int32TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int32._DelTestHelper, map[string]bwtesting.TestCaseStruct{"Int32.Del": {
		In: []interface{}{Int32{
			_Int32TestItemA: struct{}{},
			_Int32TestItemB: struct{}{},
		}, []int32{_Int32TestItemB}},
		Out: []interface{}{Int32{_Int32TestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Int32._DelSliceTestHelper, map[string]bwtesting.TestCaseStruct{"Int32.DelSlice": {
		In: []interface{}{Int32{
			_Int32TestItemA: struct{}{},
			_Int32TestItemB: struct{}{},
		}, []int32{_Int32TestItemB}},
		Out: []interface{}{Int32{_Int32TestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Int32._DelSetTestHelper, map[string]bwtesting.TestCaseStruct{"Int32.DelSet": {
		In: []interface{}{Int32{
			_Int32TestItemA: struct{}{},
			_Int32TestItemB: struct{}{},
		}, Int32{_Int32TestItemB: struct{}{}}},
		Out: []interface{}{Int32{_Int32TestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Int32.Union, map[string]bwtesting.TestCaseStruct{"Int32.Union": {
		In: []interface{}{Int32{_Int32TestItemA: struct{}{}}, Int32{_Int32TestItemB: struct{}{}}},
		Out: []interface{}{Int32{
			_Int32TestItemA: struct{}{},
			_Int32TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int32.Intersect, map[string]bwtesting.TestCaseStruct{"Int32.Intersect": {
		In: []interface{}{Int32{
			_Int32TestItemA: struct{}{},
			_Int32TestItemB: struct{}{},
		}, Int32{_Int32TestItemB: struct{}{}}},
		Out: []interface{}{Int32{_Int32TestItemB: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Int32.Subtract, map[string]bwtesting.TestCaseStruct{"Int32.Subtract": {
		In: []interface{}{Int32{
			_Int32TestItemA: struct{}{},
			_Int32TestItemB: struct{}{},
		}, Int32{_Int32TestItemB: struct{}{}}},
		Out: []interface{}{Int32{_Int32TestItemA: struct{}{}}},
	}})
}
