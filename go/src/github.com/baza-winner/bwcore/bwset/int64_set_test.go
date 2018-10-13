// Code generated by "setter -type=int64"; DO NOT EDIT; setter: go get github.com/baza-winner/bwcore/setter

package bwset

import (
	"fmt"
	bwtesting "github.com/baza-winner/bwcore/bwtesting"
	"strconv"
	"testing"
)

func TestInt64Set(t *testing.T) {
	bwtesting.BwRunTests(t, Int64SetFrom, map[string]bwtesting.TestCaseStruct{"Int64SetFrom": {
		In: []interface{}{[]int64{_Int64SetTestItemA, _Int64SetTestItemB}},
		Out: []interface{}{Int64Set{
			_Int64SetTestItemA: struct{}{},
			_Int64SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int64SetFromSlice, map[string]bwtesting.TestCaseStruct{"Int64SetFromSlice": {
		In: []interface{}{[]int64{_Int64SetTestItemA, _Int64SetTestItemB}},
		Out: []interface{}{Int64Set{
			_Int64SetTestItemA: struct{}{},
			_Int64SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int64SetFromSet, map[string]bwtesting.TestCaseStruct{"Int64SetFromSet": {
		In: []interface{}{Int64Set{
			_Int64SetTestItemA: struct{}{},
			_Int64SetTestItemB: struct{}{},
		}},
		Out: []interface{}{Int64Set{
			_Int64SetTestItemA: struct{}{},
			_Int64SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int64Set.Copy, map[string]bwtesting.TestCaseStruct{"Int64Set.Copy": {
		In: []interface{}{Int64Set{
			_Int64SetTestItemA: struct{}{},
			_Int64SetTestItemB: struct{}{},
		}},
		Out: []interface{}{Int64Set{
			_Int64SetTestItemA: struct{}{},
			_Int64SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int64Set.ToSlice, map[string]bwtesting.TestCaseStruct{"Int64Set.ToSlice": {
		In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}},
		Out: []interface{}{[]int64{_Int64SetTestItemA}},
	}})
	bwtesting.BwRunTests(t, _Int64SetToSliceTestHelper, map[string]bwtesting.TestCaseStruct{"_Int64SetToSliceTestHelper": {
		In:  []interface{}{[]int64{_Int64SetTestItemB, _Int64SetTestItemA}},
		Out: []interface{}{[]int64{_Int64SetTestItemA, _Int64SetTestItemB}},
	}})
	bwtesting.BwRunTests(t, Int64Set.String, map[string]bwtesting.TestCaseStruct{"Int64Set.String": {
		In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}},
		Out: []interface{}{fmt.Sprintf("[\n  %q\n]", strconv.FormatInt(int64(_Int64SetTestItemA), 10))},
	}})
	bwtesting.BwRunTests(t, Int64Set.GetDataForJson, map[string]bwtesting.TestCaseStruct{"Int64Set.GetDataForJson": {
		In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}},
		Out: []interface{}{[]interface{}{strconv.FormatInt(int64(_Int64SetTestItemA), 10)}},
	}})
	bwtesting.BwRunTests(t, Int64Set.ToSliceOfStrings, map[string]bwtesting.TestCaseStruct{"Int64Set.ToSliceOfStrings": {
		In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}},
		Out: []interface{}{[]string{strconv.FormatInt(int64(_Int64SetTestItemA), 10)}},
	}})
	bwtesting.BwRunTests(t, Int64Set.Has, map[string]bwtesting.TestCaseStruct{
		"Int64Set.Has: false": {
			In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, _Int64SetTestItemB},
			Out: []interface{}{false},
		},
		"Int64Set.Has: true": {
			In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, _Int64SetTestItemA},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int64Set.HasAny, map[string]bwtesting.TestCaseStruct{
		"Int64Set.HasAny: empty": {
			In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, []int64{}},
			Out: []interface{}{false},
		},
		"Int64Set.HasAny: false": {
			In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, []int64{_Int64SetTestItemB}},
			Out: []interface{}{false},
		},
		"Int64Set.HasAny: true": {
			In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, []int64{_Int64SetTestItemA, _Int64SetTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int64Set.HasAnyOfSlice, map[string]bwtesting.TestCaseStruct{
		"Int64Set.HasAnyOfSlice: empty": {
			In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, []int64{}},
			Out: []interface{}{false},
		},
		"Int64Set.HasAnyOfSlice: false": {
			In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, []int64{_Int64SetTestItemB}},
			Out: []interface{}{false},
		},
		"Int64Set.HasAnyOfSlice: true": {
			In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, []int64{_Int64SetTestItemA, _Int64SetTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int64Set.HasAnyOfSet, map[string]bwtesting.TestCaseStruct{
		"Int64Set.HasAnyOfSet: empty": {
			In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, Int64Set{}},
			Out: []interface{}{false},
		},
		"Int64Set.HasAnyOfSet: false": {
			In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, Int64Set{_Int64SetTestItemB: struct{}{}}},
			Out: []interface{}{false},
		},
		"Int64Set.HasAnyOfSet: true": {
			In: []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, Int64Set{
				_Int64SetTestItemA: struct{}{},
				_Int64SetTestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int64Set.HasEach, map[string]bwtesting.TestCaseStruct{
		"Int64Set.HasEach: empty": {
			In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, []int64{}},
			Out: []interface{}{true},
		},
		"Int64Set.HasEach: false": {
			In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, []int64{_Int64SetTestItemA, _Int64SetTestItemB}},
			Out: []interface{}{false},
		},
		"Int64Set.HasEach: true": {
			In: []interface{}{Int64Set{
				_Int64SetTestItemA: struct{}{},
				_Int64SetTestItemB: struct{}{},
			}, []int64{_Int64SetTestItemA, _Int64SetTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int64Set.HasEachOfSlice, map[string]bwtesting.TestCaseStruct{
		"Int64Set.HasEachOfSlice: empty": {
			In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, []int64{}},
			Out: []interface{}{true},
		},
		"Int64Set.HasEachOfSlice: false": {
			In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, []int64{_Int64SetTestItemA, _Int64SetTestItemB}},
			Out: []interface{}{false},
		},
		"Int64Set.HasEachOfSlice: true": {
			In: []interface{}{Int64Set{
				_Int64SetTestItemA: struct{}{},
				_Int64SetTestItemB: struct{}{},
			}, []int64{_Int64SetTestItemA, _Int64SetTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int64Set.HasEachOfSet, map[string]bwtesting.TestCaseStruct{
		"Int64Set.HasEachOfSet: empty": {
			In:  []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, Int64Set{}},
			Out: []interface{}{true},
		},
		"Int64Set.HasEachOfSet: false": {
			In: []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, Int64Set{
				_Int64SetTestItemA: struct{}{},
				_Int64SetTestItemB: struct{}{},
			}},
			Out: []interface{}{false},
		},
		"Int64Set.HasEachOfSet: true": {
			In: []interface{}{Int64Set{
				_Int64SetTestItemA: struct{}{},
				_Int64SetTestItemB: struct{}{},
			}, Int64Set{
				_Int64SetTestItemA: struct{}{},
				_Int64SetTestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Int64Set._AddTestHelper, map[string]bwtesting.TestCaseStruct{"Int64Set.Add": {
		In: []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, []int64{_Int64SetTestItemB}},
		Out: []interface{}{Int64Set{
			_Int64SetTestItemA: struct{}{},
			_Int64SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int64Set._AddSliceTestHelper, map[string]bwtesting.TestCaseStruct{"Int64Set.AddSlice": {
		In: []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, []int64{_Int64SetTestItemB}},
		Out: []interface{}{Int64Set{
			_Int64SetTestItemA: struct{}{},
			_Int64SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int64Set._AddSetTestHelper, map[string]bwtesting.TestCaseStruct{"Int64Set.AddSet": {
		In: []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, Int64Set{_Int64SetTestItemB: struct{}{}}},
		Out: []interface{}{Int64Set{
			_Int64SetTestItemA: struct{}{},
			_Int64SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int64Set._DelTestHelper, map[string]bwtesting.TestCaseStruct{"Int64Set.Del": {
		In: []interface{}{Int64Set{
			_Int64SetTestItemA: struct{}{},
			_Int64SetTestItemB: struct{}{},
		}, []int64{_Int64SetTestItemB}},
		Out: []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Int64Set._DelSliceTestHelper, map[string]bwtesting.TestCaseStruct{"Int64Set.DelSlice": {
		In: []interface{}{Int64Set{
			_Int64SetTestItemA: struct{}{},
			_Int64SetTestItemB: struct{}{},
		}, []int64{_Int64SetTestItemB}},
		Out: []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Int64Set._DelSetTestHelper, map[string]bwtesting.TestCaseStruct{"Int64Set.DelSet": {
		In: []interface{}{Int64Set{
			_Int64SetTestItemA: struct{}{},
			_Int64SetTestItemB: struct{}{},
		}, Int64Set{_Int64SetTestItemB: struct{}{}}},
		Out: []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Int64Set.Union, map[string]bwtesting.TestCaseStruct{"Int64Set.Union": {
		In: []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}, Int64Set{_Int64SetTestItemB: struct{}{}}},
		Out: []interface{}{Int64Set{
			_Int64SetTestItemA: struct{}{},
			_Int64SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Int64Set.Intersect, map[string]bwtesting.TestCaseStruct{"Int64Set.Intersect": {
		In: []interface{}{Int64Set{
			_Int64SetTestItemA: struct{}{},
			_Int64SetTestItemB: struct{}{},
		}, Int64Set{_Int64SetTestItemB: struct{}{}}},
		Out: []interface{}{Int64Set{_Int64SetTestItemB: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Int64Set.Subtract, map[string]bwtesting.TestCaseStruct{"Int64Set.Subtract": {
		In: []interface{}{Int64Set{
			_Int64SetTestItemA: struct{}{},
			_Int64SetTestItemB: struct{}{},
		}, Int64Set{_Int64SetTestItemB: struct{}{}}},
		Out: []interface{}{Int64Set{_Int64SetTestItemA: struct{}{}}},
	}})
}
