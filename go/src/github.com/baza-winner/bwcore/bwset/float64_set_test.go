// Code generated by "setter -type=float64"; DO NOT EDIT; setter: go get github.com/baza-winner/bwcore/setter

package bwset

import (
	"fmt"
	bwtesting "github.com/baza-winner/bwcore/bwtesting"
	"strconv"
	"testing"
)

func TestFloat64Set(t *testing.T) {
	bwtesting.BwRunTests(t, Float64SetFrom, map[string]bwtesting.TestCaseStruct{"Float64SetFrom": {
		In: []interface{}{[]float64{_Float64SetTestItemA, _Float64SetTestItemB}},
		Out: []interface{}{Float64Set{
			_Float64SetTestItemA: struct{}{},
			_Float64SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float64SetFromSlice, map[string]bwtesting.TestCaseStruct{"Float64SetFromSlice": {
		In: []interface{}{[]float64{_Float64SetTestItemA, _Float64SetTestItemB}},
		Out: []interface{}{Float64Set{
			_Float64SetTestItemA: struct{}{},
			_Float64SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float64SetFromSet, map[string]bwtesting.TestCaseStruct{"Float64SetFromSet": {
		In: []interface{}{Float64Set{
			_Float64SetTestItemA: struct{}{},
			_Float64SetTestItemB: struct{}{},
		}},
		Out: []interface{}{Float64Set{
			_Float64SetTestItemA: struct{}{},
			_Float64SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float64Set.Copy, map[string]bwtesting.TestCaseStruct{"Float64Set.Copy": {
		In: []interface{}{Float64Set{
			_Float64SetTestItemA: struct{}{},
			_Float64SetTestItemB: struct{}{},
		}},
		Out: []interface{}{Float64Set{
			_Float64SetTestItemA: struct{}{},
			_Float64SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float64Set.ToSlice, map[string]bwtesting.TestCaseStruct{"Float64Set.ToSlice": {
		In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}},
		Out: []interface{}{[]float64{_Float64SetTestItemA}},
	}})
	bwtesting.BwRunTests(t, _Float64SetToSliceTestHelper, map[string]bwtesting.TestCaseStruct{"_Float64SetToSliceTestHelper": {
		In:  []interface{}{[]float64{_Float64SetTestItemB, _Float64SetTestItemA}},
		Out: []interface{}{[]float64{_Float64SetTestItemA, _Float64SetTestItemB}},
	}})
	bwtesting.BwRunTests(t, Float64Set.String, map[string]bwtesting.TestCaseStruct{"Float64Set.String": {
		In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}},
		Out: []interface{}{fmt.Sprintf("[\n  %q\n]", strconv.FormatFloat(float64(_Float64SetTestItemA), byte(0x66), -1, 64))},
	}})
	bwtesting.BwRunTests(t, Float64Set.GetDataForJson, map[string]bwtesting.TestCaseStruct{"Float64Set.GetDataForJson": {
		In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}},
		Out: []interface{}{[]interface{}{strconv.FormatFloat(float64(_Float64SetTestItemA), byte(0x66), -1, 64)}},
	}})
	bwtesting.BwRunTests(t, Float64Set.ToSliceOfStrings, map[string]bwtesting.TestCaseStruct{"Float64Set.ToSliceOfStrings": {
		In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}},
		Out: []interface{}{[]string{strconv.FormatFloat(float64(_Float64SetTestItemA), byte(0x66), -1, 64)}},
	}})
	bwtesting.BwRunTests(t, Float64Set.Has, map[string]bwtesting.TestCaseStruct{
		"Float64Set.Has: false": {
			In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, _Float64SetTestItemB},
			Out: []interface{}{false},
		},
		"Float64Set.Has: true": {
			In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, _Float64SetTestItemA},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float64Set.HasAny, map[string]bwtesting.TestCaseStruct{
		"Float64Set.HasAny: empty": {
			In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, []float64{}},
			Out: []interface{}{false},
		},
		"Float64Set.HasAny: false": {
			In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, []float64{_Float64SetTestItemB}},
			Out: []interface{}{false},
		},
		"Float64Set.HasAny: true": {
			In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, []float64{_Float64SetTestItemA, _Float64SetTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float64Set.HasAnyOfSlice, map[string]bwtesting.TestCaseStruct{
		"Float64Set.HasAnyOfSlice: empty": {
			In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, []float64{}},
			Out: []interface{}{false},
		},
		"Float64Set.HasAnyOfSlice: false": {
			In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, []float64{_Float64SetTestItemB}},
			Out: []interface{}{false},
		},
		"Float64Set.HasAnyOfSlice: true": {
			In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, []float64{_Float64SetTestItemA, _Float64SetTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float64Set.HasAnyOfSet, map[string]bwtesting.TestCaseStruct{
		"Float64Set.HasAnyOfSet: empty": {
			In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, Float64Set{}},
			Out: []interface{}{false},
		},
		"Float64Set.HasAnyOfSet: false": {
			In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, Float64Set{_Float64SetTestItemB: struct{}{}}},
			Out: []interface{}{false},
		},
		"Float64Set.HasAnyOfSet: true": {
			In: []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, Float64Set{
				_Float64SetTestItemA: struct{}{},
				_Float64SetTestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float64Set.HasEach, map[string]bwtesting.TestCaseStruct{
		"Float64Set.HasEach: empty": {
			In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, []float64{}},
			Out: []interface{}{true},
		},
		"Float64Set.HasEach: false": {
			In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, []float64{_Float64SetTestItemA, _Float64SetTestItemB}},
			Out: []interface{}{false},
		},
		"Float64Set.HasEach: true": {
			In: []interface{}{Float64Set{
				_Float64SetTestItemA: struct{}{},
				_Float64SetTestItemB: struct{}{},
			}, []float64{_Float64SetTestItemA, _Float64SetTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float64Set.HasEachOfSlice, map[string]bwtesting.TestCaseStruct{
		"Float64Set.HasEachOfSlice: empty": {
			In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, []float64{}},
			Out: []interface{}{true},
		},
		"Float64Set.HasEachOfSlice: false": {
			In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, []float64{_Float64SetTestItemA, _Float64SetTestItemB}},
			Out: []interface{}{false},
		},
		"Float64Set.HasEachOfSlice: true": {
			In: []interface{}{Float64Set{
				_Float64SetTestItemA: struct{}{},
				_Float64SetTestItemB: struct{}{},
			}, []float64{_Float64SetTestItemA, _Float64SetTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float64Set.HasEachOfSet, map[string]bwtesting.TestCaseStruct{
		"Float64Set.HasEachOfSet: empty": {
			In:  []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, Float64Set{}},
			Out: []interface{}{true},
		},
		"Float64Set.HasEachOfSet: false": {
			In: []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, Float64Set{
				_Float64SetTestItemA: struct{}{},
				_Float64SetTestItemB: struct{}{},
			}},
			Out: []interface{}{false},
		},
		"Float64Set.HasEachOfSet: true": {
			In: []interface{}{Float64Set{
				_Float64SetTestItemA: struct{}{},
				_Float64SetTestItemB: struct{}{},
			}, Float64Set{
				_Float64SetTestItemA: struct{}{},
				_Float64SetTestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Float64Set._AddTestHelper, map[string]bwtesting.TestCaseStruct{"Float64Set.Add": {
		In: []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, []float64{_Float64SetTestItemB}},
		Out: []interface{}{Float64Set{
			_Float64SetTestItemA: struct{}{},
			_Float64SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float64Set._AddSliceTestHelper, map[string]bwtesting.TestCaseStruct{"Float64Set.AddSlice": {
		In: []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, []float64{_Float64SetTestItemB}},
		Out: []interface{}{Float64Set{
			_Float64SetTestItemA: struct{}{},
			_Float64SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float64Set._AddSetTestHelper, map[string]bwtesting.TestCaseStruct{"Float64Set.AddSet": {
		In: []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, Float64Set{_Float64SetTestItemB: struct{}{}}},
		Out: []interface{}{Float64Set{
			_Float64SetTestItemA: struct{}{},
			_Float64SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float64Set._DelTestHelper, map[string]bwtesting.TestCaseStruct{"Float64Set.Del": {
		In: []interface{}{Float64Set{
			_Float64SetTestItemA: struct{}{},
			_Float64SetTestItemB: struct{}{},
		}, []float64{_Float64SetTestItemB}},
		Out: []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Float64Set._DelSliceTestHelper, map[string]bwtesting.TestCaseStruct{"Float64Set.DelSlice": {
		In: []interface{}{Float64Set{
			_Float64SetTestItemA: struct{}{},
			_Float64SetTestItemB: struct{}{},
		}, []float64{_Float64SetTestItemB}},
		Out: []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Float64Set._DelSetTestHelper, map[string]bwtesting.TestCaseStruct{"Float64Set.DelSet": {
		In: []interface{}{Float64Set{
			_Float64SetTestItemA: struct{}{},
			_Float64SetTestItemB: struct{}{},
		}, Float64Set{_Float64SetTestItemB: struct{}{}}},
		Out: []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Float64Set.Union, map[string]bwtesting.TestCaseStruct{"Float64Set.Union": {
		In: []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}, Float64Set{_Float64SetTestItemB: struct{}{}}},
		Out: []interface{}{Float64Set{
			_Float64SetTestItemA: struct{}{},
			_Float64SetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Float64Set.Intersect, map[string]bwtesting.TestCaseStruct{"Float64Set.Intersect": {
		In: []interface{}{Float64Set{
			_Float64SetTestItemA: struct{}{},
			_Float64SetTestItemB: struct{}{},
		}, Float64Set{_Float64SetTestItemB: struct{}{}}},
		Out: []interface{}{Float64Set{_Float64SetTestItemB: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Float64Set.Subtract, map[string]bwtesting.TestCaseStruct{"Float64Set.Subtract": {
		In: []interface{}{Float64Set{
			_Float64SetTestItemA: struct{}{},
			_Float64SetTestItemB: struct{}{},
		}, Float64Set{_Float64SetTestItemB: struct{}{}}},
		Out: []interface{}{Float64Set{_Float64SetTestItemA: struct{}{}}},
	}})
}
