// Code generated by "bwsetter -type=int64"; DO NOT EDIT; bwsetter: go get -type=int64 -set=Int64 -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	"encoding/json"
	bwtesting "github.com/baza-winner/bwcore/bwtesting"
	"strconv"
	"testing"
)

func TestInt64(t *testing.T) {
	bwtesting.BwRunTests(t, Int64From, map[string]bwtesting.Case{"Int64From": {
		In: []interface{}{_Int64TestItemA, _Int64TestItemB},
		Out: []interface{}{Int64{
			_Int64TestItemA: struct{}{},
			_Int64TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int64FromSlice, map[string]bwtesting.Case{"Int64FromSlice": {
		In: []interface{}{[]int64{_Int64TestItemA, _Int64TestItemB}},
		Out: []interface{}{Int64{
			_Int64TestItemA: struct{}{},
			_Int64TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int64FromSet, map[string]bwtesting.Case{"Int64FromSet": {
		In: []interface{}{Int64{
			_Int64TestItemA: struct{}{},
			_Int64TestItemB: struct{}{},
		}},
		Out: []interface{}{Int64{
			_Int64TestItemA: struct{}{},
			_Int64TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int64.Copy, map[string]bwtesting.Case{"Int64.Copy": {
		In: []interface{}{Int64{
			_Int64TestItemA: struct{}{},
			_Int64TestItemB: struct{}{},
		}},
		Out: []interface{}{Int64{
			_Int64TestItemA: struct{}{},
			_Int64TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int64.ToSlice, map[string]bwtesting.Case{"Int64.ToSlice": {
		In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}},
		Out: []interface{}{[]int64{_Int64TestItemA}},
	}})
	bwtesting.BwRunTests(t, _Int64ToSliceTestHelper, map[string]bwtesting.Case{"_Int64ToSliceTestHelper": {
		In:  []interface{}{[]int64{_Int64TestItemB, _Int64TestItemA}},
		Out: []interface{}{[]int64{_Int64TestItemA, _Int64TestItemB}},
	}})
	bwtesting.BwRunTests(t, Int64.String, map[string]bwtesting.Case{"Int64.String": {
		In: []interface{}{Int64{_Int64TestItemA: struct{}{}}},
		Out: []interface{}{func() string {
			result, _ := json.Marshal(_Int64TestItemA)
			return "[" + string(result) + "]"
		}()},
	}})
	bwtesting.BwRunTests(t, Int64.MarshalJSON, map[string]bwtesting.Case{"Int64.MarshalJSON": {
		In: []interface{}{Int64{_Int64TestItemA: struct{}{}}},
		Out: []interface{}{(func() []byte {
			result, _ := json.Marshal([]interface{}{_Int64TestItemA})
			return result
		})(), nil},
	}})
	bwtesting.BwRunTests(t, Int64.ToSliceOfStrings, map[string]bwtesting.Case{"Int64.ToSliceOfStrings": {
		In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}},
		Out: []interface{}{[]string{strconv.FormatInt(int64(_Int64TestItemA), 10)}},
	}})
	bwtesting.BwRunTests(t, Int64.Has, map[string]bwtesting.Case{
		"Int64.Has: false": {
			In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}, _Int64TestItemB},
			Out: []interface{}{false},
		},
		"Int64.Has: true": {
			In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}, _Int64TestItemA},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int64.HasAny, map[string]bwtesting.Case{
		"Int64.HasAny: empty": {
			In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}},
			Out: []interface{}{false},
		},
		"Int64.HasAny: false": {
			In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}, _Int64TestItemB},
			Out: []interface{}{false},
		},
		"Int64.HasAny: true": {
			In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}, _Int64TestItemA, _Int64TestItemB},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int64.HasAnyOfSlice, map[string]bwtesting.Case{
		"Int64.HasAnyOfSlice: empty": {
			In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}, []int64{}},
			Out: []interface{}{false},
		},
		"Int64.HasAnyOfSlice: false": {
			In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}, []int64{_Int64TestItemB}},
			Out: []interface{}{false},
		},
		"Int64.HasAnyOfSlice: true": {
			In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}, []int64{_Int64TestItemA, _Int64TestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int64.HasAnyOfSet, map[string]bwtesting.Case{
		"Int64.HasAnyOfSet: empty": {
			In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}, Int64{}},
			Out: []interface{}{false},
		},
		"Int64.HasAnyOfSet: false": {
			In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}, Int64{_Int64TestItemB: struct{}{}}},
			Out: []interface{}{false},
		},
		"Int64.HasAnyOfSet: true": {
			In: []interface{}{Int64{_Int64TestItemA: struct{}{}}, Int64{
				_Int64TestItemA: struct{}{},
				_Int64TestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int64.HasEach, map[string]bwtesting.Case{
		"Int64.HasEach: empty": {
			In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}},
			Out: []interface{}{true},
		},
		"Int64.HasEach: false": {
			In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}, _Int64TestItemA, _Int64TestItemB},
			Out: []interface{}{false},
		},
		"Int64.HasEach: true": {
			In: []interface{}{Int64{
				_Int64TestItemA: struct{}{},
				_Int64TestItemB: struct{}{},
			}, _Int64TestItemA, _Int64TestItemB},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int64.HasEachOfSlice, map[string]bwtesting.Case{
		"Int64.HasEachOfSlice: empty": {
			In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}, []int64{}},
			Out: []interface{}{true},
		},
		"Int64.HasEachOfSlice: false": {
			In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}, []int64{_Int64TestItemA, _Int64TestItemB}},
			Out: []interface{}{false},
		},
		"Int64.HasEachOfSlice: true": {
			In: []interface{}{Int64{
				_Int64TestItemA: struct{}{},
				_Int64TestItemB: struct{}{},
			}, []int64{_Int64TestItemA, _Int64TestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int64.HasEachOfSet, map[string]bwtesting.Case{
		"Int64.HasEachOfSet: empty": {
			In:  []interface{}{Int64{_Int64TestItemA: struct{}{}}, Int64{}},
			Out: []interface{}{true},
		},
		"Int64.HasEachOfSet: false": {
			In: []interface{}{Int64{_Int64TestItemA: struct{}{}}, Int64{
				_Int64TestItemA: struct{}{},
				_Int64TestItemB: struct{}{},
			}},
			Out: []interface{}{false},
		},
		"Int64.HasEachOfSet: true": {
			In: []interface{}{Int64{
				_Int64TestItemA: struct{}{},
				_Int64TestItemB: struct{}{},
			}, Int64{
				_Int64TestItemA: struct{}{},
				_Int64TestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int64._AddTestHelper, map[string]bwtesting.Case{"Int64.Add": {
		In: []interface{}{Int64{_Int64TestItemA: struct{}{}}, _Int64TestItemB},
		Out: []interface{}{Int64{
			_Int64TestItemA: struct{}{},
			_Int64TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int64._AddSliceTestHelper, map[string]bwtesting.Case{"Int64.AddSlice": {
		In: []interface{}{Int64{_Int64TestItemA: struct{}{}}, []int64{_Int64TestItemB}},
		Out: []interface{}{Int64{
			_Int64TestItemA: struct{}{},
			_Int64TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int64._AddSetTestHelper, map[string]bwtesting.Case{"Int64.AddSet": {
		In: []interface{}{Int64{_Int64TestItemA: struct{}{}}, Int64{_Int64TestItemB: struct{}{}}},
		Out: []interface{}{Int64{
			_Int64TestItemA: struct{}{},
			_Int64TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int64._DelTestHelper, map[string]bwtesting.Case{"Int64.Del": {
		In: []interface{}{Int64{
			_Int64TestItemA: struct{}{},
			_Int64TestItemB: struct{}{},
		}, _Int64TestItemB},
		Out: []interface{}{Int64{_Int64TestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Int64._DelSliceTestHelper, map[string]bwtesting.Case{"Int64.DelSlice": {
		In: []interface{}{Int64{
			_Int64TestItemA: struct{}{},
			_Int64TestItemB: struct{}{},
		}, []int64{_Int64TestItemB}},
		Out: []interface{}{Int64{_Int64TestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Int64._DelSetTestHelper, map[string]bwtesting.Case{"Int64.DelSet": {
		In: []interface{}{Int64{
			_Int64TestItemA: struct{}{},
			_Int64TestItemB: struct{}{},
		}, Int64{_Int64TestItemB: struct{}{}}},
		Out: []interface{}{Int64{_Int64TestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Int64.Union, map[string]bwtesting.Case{"Int64.Union": {
		In: []interface{}{Int64{_Int64TestItemA: struct{}{}}, Int64{_Int64TestItemB: struct{}{}}},
		Out: []interface{}{Int64{
			_Int64TestItemA: struct{}{},
			_Int64TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int64.Intersect, map[string]bwtesting.Case{"Int64.Intersect": {
		In: []interface{}{Int64{
			_Int64TestItemA: struct{}{},
			_Int64TestItemB: struct{}{},
		}, Int64{_Int64TestItemB: struct{}{}}},
		Out: []interface{}{Int64{_Int64TestItemB: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Int64.Subtract, map[string]bwtesting.Case{"Int64.Subtract": {
		In: []interface{}{Int64{
			_Int64TestItemA: struct{}{},
			_Int64TestItemB: struct{}{},
		}, Int64{_Int64TestItemB: struct{}{}}},
		Out: []interface{}{Int64{_Int64TestItemA: struct{}{}}},
	}})
}
