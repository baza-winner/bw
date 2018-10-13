// Code generated by "setter -type=uint32"; DO NOT EDIT; setter: go get github.com/baza-winner/bwcore/setter

package bwset

import (
	"fmt"
	bwtesting "github.com/baza-winner/bwcore/bwtesting"
	"strconv"
	"testing"
)

func TestUint32Set(t *testing.T) {
	bwtesting.BwRunTests(t, Uint32SetFrom, map[string]bwtesting.TestCaseStruct{"Uint32SetFrom": {
		In: []interface{}{[]uint32{_Uint32SetTestItemA, _Uint32SetTestItemB}},
		Out: []interface{}{Uint32Set{
			_Uint32SetTestItemA: struct{}{},
			_Uint32SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint32SetFromSlice, map[string]bwtesting.TestCaseStruct{"Uint32SetFromSlice": {
		In: []interface{}{[]uint32{_Uint32SetTestItemA, _Uint32SetTestItemB}},
		Out: []interface{}{Uint32Set{
			_Uint32SetTestItemA: struct{}{},
			_Uint32SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint32SetFromSet, map[string]bwtesting.TestCaseStruct{"Uint32SetFromSet": {
		In: []interface{}{Uint32Set{
			_Uint32SetTestItemA: struct{}{},
			_Uint32SetTestItemB: struct{}{},
		}},
		Out: []interface{}{Uint32Set{
			_Uint32SetTestItemA: struct{}{},
			_Uint32SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint32Set.Copy, map[string]bwtesting.TestCaseStruct{"Uint32Set.Copy": {
		In: []interface{}{Uint32Set{
			_Uint32SetTestItemA: struct{}{},
			_Uint32SetTestItemB: struct{}{},
		}},
		Out: []interface{}{Uint32Set{
			_Uint32SetTestItemA: struct{}{},
			_Uint32SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint32Set.ToSlice, map[string]bwtesting.TestCaseStruct{"Uint32Set.ToSlice": {
		In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}},
		Out: []interface{}{[]uint32{_Uint32SetTestItemA}},
	}})
	bwtesting.BwRunTests(t, _Uint32SetToSliceTestHelper, map[string]bwtesting.TestCaseStruct{"_Uint32SetToSliceTestHelper": {
		In:  []interface{}{[]uint32{_Uint32SetTestItemB, _Uint32SetTestItemA}},
		Out: []interface{}{[]uint32{_Uint32SetTestItemA, _Uint32SetTestItemB}},
	}})
	bwtesting.BwRunTests(t, Uint32Set.String, map[string]bwtesting.TestCaseStruct{"Uint32Set.String": {
		In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}},
		Out: []interface{}{fmt.Sprintf("[\n  %q\n]", strconv.FormatUint(uint64(_Uint32SetTestItemA), 10))},
	}})
	bwtesting.BwRunTests(t, Uint32Set.GetDataForJson, map[string]bwtesting.TestCaseStruct{"Uint32Set.GetDataForJson": {
		In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}},
		Out: []interface{}{[]interface{}{strconv.FormatUint(uint64(_Uint32SetTestItemA), 10)}},
	}})
	bwtesting.BwRunTests(t, Uint32Set.ToSliceOfStrings, map[string]bwtesting.TestCaseStruct{"Uint32Set.ToSliceOfStrings": {
		In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}},
		Out: []interface{}{[]string{strconv.FormatUint(uint64(_Uint32SetTestItemA), 10)}},
	}})
	bwtesting.BwRunTests(t, Uint32Set.Has, map[string]bwtesting.TestCaseStruct{
		"Uint32Set.Has: false": {
			In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, _Uint32SetTestItemB},
			Out: []interface{}{false},
		},
		"Uint32Set.Has: true": {
			In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, _Uint32SetTestItemA},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint32Set.HasAny, map[string]bwtesting.TestCaseStruct{
		"Uint32Set.HasAny: empty": {
			In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, []uint32{}},
			Out: []interface{}{false},
		},
		"Uint32Set.HasAny: false": {
			In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, []uint32{_Uint32SetTestItemB}},
			Out: []interface{}{false},
		},
		"Uint32Set.HasAny: true": {
			In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, []uint32{_Uint32SetTestItemA, _Uint32SetTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint32Set.HasAnyOfSlice, map[string]bwtesting.TestCaseStruct{
		"Uint32Set.HasAnyOfSlice: empty": {
			In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, []uint32{}},
			Out: []interface{}{false},
		},
		"Uint32Set.HasAnyOfSlice: false": {
			In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, []uint32{_Uint32SetTestItemB}},
			Out: []interface{}{false},
		},
		"Uint32Set.HasAnyOfSlice: true": {
			In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, []uint32{_Uint32SetTestItemA, _Uint32SetTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint32Set.HasAnyOfSet, map[string]bwtesting.TestCaseStruct{
		"Uint32Set.HasAnyOfSet: empty": {
			In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, Uint32Set{}},
			Out: []interface{}{false},
		},
		"Uint32Set.HasAnyOfSet: false": {
			In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, Uint32Set{_Uint32SetTestItemB: struct{}{}}},
			Out: []interface{}{false},
		},
		"Uint32Set.HasAnyOfSet: true": {
			In: []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, Uint32Set{
				_Uint32SetTestItemA: struct{}{},
				_Uint32SetTestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint32Set.HasEach, map[string]bwtesting.TestCaseStruct{
		"Uint32Set.HasEach: empty": {
			In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, []uint32{}},
			Out: []interface{}{true},
		},
		"Uint32Set.HasEach: false": {
			In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, []uint32{_Uint32SetTestItemA, _Uint32SetTestItemB}},
			Out: []interface{}{false},
		},
		"Uint32Set.HasEach: true": {
			In: []interface{}{Uint32Set{
				_Uint32SetTestItemA: struct{}{},
				_Uint32SetTestItemB: struct{}{},
			}, []uint32{_Uint32SetTestItemA, _Uint32SetTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint32Set.HasEachOfSlice, map[string]bwtesting.TestCaseStruct{
		"Uint32Set.HasEachOfSlice: empty": {
			In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, []uint32{}},
			Out: []interface{}{true},
		},
		"Uint32Set.HasEachOfSlice: false": {
			In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, []uint32{_Uint32SetTestItemA, _Uint32SetTestItemB}},
			Out: []interface{}{false},
		},
		"Uint32Set.HasEachOfSlice: true": {
			In: []interface{}{Uint32Set{
				_Uint32SetTestItemA: struct{}{},
				_Uint32SetTestItemB: struct{}{},
			}, []uint32{_Uint32SetTestItemA, _Uint32SetTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint32Set.HasEachOfSet, map[string]bwtesting.TestCaseStruct{
		"Uint32Set.HasEachOfSet: empty": {
			In:  []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, Uint32Set{}},
			Out: []interface{}{true},
		},
		"Uint32Set.HasEachOfSet: false": {
			In: []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, Uint32Set{
				_Uint32SetTestItemA: struct{}{},
				_Uint32SetTestItemB: struct{}{},
			}},
			Out: []interface{}{false},
		},
		"Uint32Set.HasEachOfSet: true": {
			In: []interface{}{Uint32Set{
				_Uint32SetTestItemA: struct{}{},
				_Uint32SetTestItemB: struct{}{},
			}, Uint32Set{
				_Uint32SetTestItemA: struct{}{},
				_Uint32SetTestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint32Set._AddTestHelper, map[string]bwtesting.TestCaseStruct{"Uint32Set.Add": {
		In: []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, []uint32{_Uint32SetTestItemB}},
		Out: []interface{}{Uint32Set{
			_Uint32SetTestItemA: struct{}{},
			_Uint32SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint32Set._AddSliceTestHelper, map[string]bwtesting.TestCaseStruct{"Uint32Set.AddSlice": {
		In: []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, []uint32{_Uint32SetTestItemB}},
		Out: []interface{}{Uint32Set{
			_Uint32SetTestItemA: struct{}{},
			_Uint32SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint32Set._AddSetTestHelper, map[string]bwtesting.TestCaseStruct{"Uint32Set.AddSet": {
		In: []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, Uint32Set{_Uint32SetTestItemB: struct{}{}}},
		Out: []interface{}{Uint32Set{
			_Uint32SetTestItemA: struct{}{},
			_Uint32SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint32Set._DelTestHelper, map[string]bwtesting.TestCaseStruct{"Uint32Set.Del": {
		In: []interface{}{Uint32Set{
			_Uint32SetTestItemA: struct{}{},
			_Uint32SetTestItemB: struct{}{},
		}, []uint32{_Uint32SetTestItemB}},
		Out: []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Uint32Set._DelSliceTestHelper, map[string]bwtesting.TestCaseStruct{"Uint32Set.DelSlice": {
		In: []interface{}{Uint32Set{
			_Uint32SetTestItemA: struct{}{},
			_Uint32SetTestItemB: struct{}{},
		}, []uint32{_Uint32SetTestItemB}},
		Out: []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Uint32Set._DelSetTestHelper, map[string]bwtesting.TestCaseStruct{"Uint32Set.DelSet": {
		In: []interface{}{Uint32Set{
			_Uint32SetTestItemA: struct{}{},
			_Uint32SetTestItemB: struct{}{},
		}, Uint32Set{_Uint32SetTestItemB: struct{}{}}},
		Out: []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Uint32Set.Union, map[string]bwtesting.TestCaseStruct{"Uint32Set.Union": {
		In: []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}, Uint32Set{_Uint32SetTestItemB: struct{}{}}},
		Out: []interface{}{Uint32Set{
			_Uint32SetTestItemA: struct{}{},
			_Uint32SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint32Set.Intersect, map[string]bwtesting.TestCaseStruct{"Uint32Set.Intersect": {
		In: []interface{}{Uint32Set{
			_Uint32SetTestItemA: struct{}{},
			_Uint32SetTestItemB: struct{}{},
		}, Uint32Set{_Uint32SetTestItemB: struct{}{}}},
		Out: []interface{}{Uint32Set{_Uint32SetTestItemB: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Uint32Set.Subtract, map[string]bwtesting.TestCaseStruct{"Uint32Set.Subtract": {
		In: []interface{}{Uint32Set{
			_Uint32SetTestItemA: struct{}{},
			_Uint32SetTestItemB: struct{}{},
		}, Uint32Set{_Uint32SetTestItemB: struct{}{}}},
		Out: []interface{}{Uint32Set{_Uint32SetTestItemA: struct{}{}}},
	}})
}
