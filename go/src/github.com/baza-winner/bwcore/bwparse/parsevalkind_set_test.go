// Code generated by "bwsetter -type=ParseValKind"; DO NOT EDIT; bwsetter: go get -type ParseValKind -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwparse

import (
	"encoding/json"
	bwtesting "github.com/baza-winner/bwcore/bwtesting"
	"testing"
)

func TestParseValKindSet(t *testing.T) {
	bwtesting.BwRunTests(t, ParseValKindSetFrom, map[string]bwtesting.Case{"ParseValKindSetFrom": {
		In: []interface{}{_ParseValKindSetTestItemA, _ParseValKindSetTestItemB},
		Out: []interface{}{ParseValKindSet{
			_ParseValKindSetTestItemA: struct{}{},
			_ParseValKindSetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, ParseValKindSetFromSlice, map[string]bwtesting.Case{"ParseValKindSetFromSlice": {
		In: []interface{}{[]ParseValKind{_ParseValKindSetTestItemA, _ParseValKindSetTestItemB}},
		Out: []interface{}{ParseValKindSet{
			_ParseValKindSetTestItemA: struct{}{},
			_ParseValKindSetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, ParseValKindSetFromSet, map[string]bwtesting.Case{"ParseValKindSetFromSet": {
		In: []interface{}{ParseValKindSet{
			_ParseValKindSetTestItemA: struct{}{},
			_ParseValKindSetTestItemB: struct{}{},
		}},
		Out: []interface{}{ParseValKindSet{
			_ParseValKindSetTestItemA: struct{}{},
			_ParseValKindSetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, ParseValKindSet.Copy, map[string]bwtesting.Case{"ParseValKindSet.Copy": {
		In: []interface{}{ParseValKindSet{
			_ParseValKindSetTestItemA: struct{}{},
			_ParseValKindSetTestItemB: struct{}{},
		}},
		Out: []interface{}{ParseValKindSet{
			_ParseValKindSetTestItemA: struct{}{},
			_ParseValKindSetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, ParseValKindSet.ToSlice, map[string]bwtesting.Case{"ParseValKindSet.ToSlice": {
		In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}},
		Out: []interface{}{[]ParseValKind{_ParseValKindSetTestItemA}},
	}})
	bwtesting.BwRunTests(t, _ParseValKindSetToSliceTestHelper, map[string]bwtesting.Case{"_ParseValKindSetToSliceTestHelper": {
		In:  []interface{}{[]ParseValKind{_ParseValKindSetTestItemB, _ParseValKindSetTestItemA}},
		Out: []interface{}{[]ParseValKind{_ParseValKindSetTestItemA, _ParseValKindSetTestItemB}},
	}})
	bwtesting.BwRunTests(t, ParseValKindSet.String, map[string]bwtesting.Case{"ParseValKindSet.String": {
		In: []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}},
		Out: []interface{}{func() string {
			result, _ := json.Marshal(_ParseValKindSetTestItemA)
			return "[" + string(result) + "]"
		}()},
	}})
	bwtesting.BwRunTests(t, ParseValKindSet.MarshalJSON, map[string]bwtesting.Case{"ParseValKindSet.MarshalJSON": {
		In: []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}},
		Out: []interface{}{(func() []byte {
			result, _ := json.Marshal([]interface{}{_ParseValKindSetTestItemA})
			return result
		})(), nil},
	}})
	bwtesting.BwRunTests(t, ParseValKindSet.ToSliceOfStrings, map[string]bwtesting.Case{"ParseValKindSet.ToSliceOfStrings": {
		In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}},
		Out: []interface{}{[]string{_ParseValKindSetTestItemA.String()}},
	}})
	bwtesting.BwRunTests(t, ParseValKindSet.Has, map[string]bwtesting.Case{
		"ParseValKindSet.Has: false": {
			In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, _ParseValKindSetTestItemB},
			Out: []interface{}{false},
		},
		"ParseValKindSet.Has: true": {
			In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, _ParseValKindSetTestItemA},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, ParseValKindSet.HasAny, map[string]bwtesting.Case{
		"ParseValKindSet.HasAny: empty": {
			In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}},
			Out: []interface{}{false},
		},
		"ParseValKindSet.HasAny: false": {
			In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, _ParseValKindSetTestItemB},
			Out: []interface{}{false},
		},
		"ParseValKindSet.HasAny: true": {
			In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, _ParseValKindSetTestItemA, _ParseValKindSetTestItemB},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, ParseValKindSet.HasAnyOfSlice, map[string]bwtesting.Case{
		"ParseValKindSet.HasAnyOfSlice: empty": {
			In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, []ParseValKind{}},
			Out: []interface{}{false},
		},
		"ParseValKindSet.HasAnyOfSlice: false": {
			In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, []ParseValKind{_ParseValKindSetTestItemB}},
			Out: []interface{}{false},
		},
		"ParseValKindSet.HasAnyOfSlice: true": {
			In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, []ParseValKind{_ParseValKindSetTestItemA, _ParseValKindSetTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, ParseValKindSet.HasAnyOfSet, map[string]bwtesting.Case{
		"ParseValKindSet.HasAnyOfSet: empty": {
			In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, ParseValKindSet{}},
			Out: []interface{}{false},
		},
		"ParseValKindSet.HasAnyOfSet: false": {
			In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, ParseValKindSet{_ParseValKindSetTestItemB: struct{}{}}},
			Out: []interface{}{false},
		},
		"ParseValKindSet.HasAnyOfSet: true": {
			In: []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, ParseValKindSet{
				_ParseValKindSetTestItemA: struct{}{},
				_ParseValKindSetTestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, ParseValKindSet.HasEach, map[string]bwtesting.Case{
		"ParseValKindSet.HasEach: empty": {
			In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}},
			Out: []interface{}{true},
		},
		"ParseValKindSet.HasEach: false": {
			In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, _ParseValKindSetTestItemA, _ParseValKindSetTestItemB},
			Out: []interface{}{false},
		},
		"ParseValKindSet.HasEach: true": {
			In: []interface{}{ParseValKindSet{
				_ParseValKindSetTestItemA: struct{}{},
				_ParseValKindSetTestItemB: struct{}{},
			}, _ParseValKindSetTestItemA, _ParseValKindSetTestItemB},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, ParseValKindSet.HasEachOfSlice, map[string]bwtesting.Case{
		"ParseValKindSet.HasEachOfSlice: empty": {
			In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, []ParseValKind{}},
			Out: []interface{}{true},
		},
		"ParseValKindSet.HasEachOfSlice: false": {
			In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, []ParseValKind{_ParseValKindSetTestItemA, _ParseValKindSetTestItemB}},
			Out: []interface{}{false},
		},
		"ParseValKindSet.HasEachOfSlice: true": {
			In: []interface{}{ParseValKindSet{
				_ParseValKindSetTestItemA: struct{}{},
				_ParseValKindSetTestItemB: struct{}{},
			}, []ParseValKind{_ParseValKindSetTestItemA, _ParseValKindSetTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, ParseValKindSet.HasEachOfSet, map[string]bwtesting.Case{
		"ParseValKindSet.HasEachOfSet: empty": {
			In:  []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, ParseValKindSet{}},
			Out: []interface{}{true},
		},
		"ParseValKindSet.HasEachOfSet: false": {
			In: []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, ParseValKindSet{
				_ParseValKindSetTestItemA: struct{}{},
				_ParseValKindSetTestItemB: struct{}{},
			}},
			Out: []interface{}{false},
		},
		"ParseValKindSet.HasEachOfSet: true": {
			In: []interface{}{ParseValKindSet{
				_ParseValKindSetTestItemA: struct{}{},
				_ParseValKindSetTestItemB: struct{}{},
			}, ParseValKindSet{
				_ParseValKindSetTestItemA: struct{}{},
				_ParseValKindSetTestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, ParseValKindSet._AddTestHelper, map[string]bwtesting.Case{"ParseValKindSet.Add": {
		In: []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, _ParseValKindSetTestItemB},
		Out: []interface{}{ParseValKindSet{
			_ParseValKindSetTestItemA: struct{}{},
			_ParseValKindSetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, ParseValKindSet._AddSliceTestHelper, map[string]bwtesting.Case{"ParseValKindSet.AddSlice": {
		In: []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, []ParseValKind{_ParseValKindSetTestItemB}},
		Out: []interface{}{ParseValKindSet{
			_ParseValKindSetTestItemA: struct{}{},
			_ParseValKindSetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, ParseValKindSet._AddSetTestHelper, map[string]bwtesting.Case{"ParseValKindSet.AddSet": {
		In: []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, ParseValKindSet{_ParseValKindSetTestItemB: struct{}{}}},
		Out: []interface{}{ParseValKindSet{
			_ParseValKindSetTestItemA: struct{}{},
			_ParseValKindSetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, ParseValKindSet._DelTestHelper, map[string]bwtesting.Case{"ParseValKindSet.Del": {
		In: []interface{}{ParseValKindSet{
			_ParseValKindSetTestItemA: struct{}{},
			_ParseValKindSetTestItemB: struct{}{},
		}, _ParseValKindSetTestItemB},
		Out: []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, ParseValKindSet._DelSliceTestHelper, map[string]bwtesting.Case{"ParseValKindSet.DelSlice": {
		In: []interface{}{ParseValKindSet{
			_ParseValKindSetTestItemA: struct{}{},
			_ParseValKindSetTestItemB: struct{}{},
		}, []ParseValKind{_ParseValKindSetTestItemB}},
		Out: []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, ParseValKindSet._DelSetTestHelper, map[string]bwtesting.Case{"ParseValKindSet.DelSet": {
		In: []interface{}{ParseValKindSet{
			_ParseValKindSetTestItemA: struct{}{},
			_ParseValKindSetTestItemB: struct{}{},
		}, ParseValKindSet{_ParseValKindSetTestItemB: struct{}{}}},
		Out: []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, ParseValKindSet.Union, map[string]bwtesting.Case{"ParseValKindSet.Union": {
		In: []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}, ParseValKindSet{_ParseValKindSetTestItemB: struct{}{}}},
		Out: []interface{}{ParseValKindSet{
			_ParseValKindSetTestItemA: struct{}{},
			_ParseValKindSetTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, ParseValKindSet.Intersect, map[string]bwtesting.Case{"ParseValKindSet.Intersect": {
		In: []interface{}{ParseValKindSet{
			_ParseValKindSetTestItemA: struct{}{},
			_ParseValKindSetTestItemB: struct{}{},
		}, ParseValKindSet{_ParseValKindSetTestItemB: struct{}{}}},
		Out: []interface{}{ParseValKindSet{_ParseValKindSetTestItemB: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, ParseValKindSet.Subtract, map[string]bwtesting.Case{"ParseValKindSet.Subtract": {
		In: []interface{}{ParseValKindSet{
			_ParseValKindSetTestItemA: struct{}{},
			_ParseValKindSetTestItemB: struct{}{},
		}, ParseValKindSet{_ParseValKindSetTestItemB: struct{}{}}},
		Out: []interface{}{ParseValKindSet{_ParseValKindSetTestItemA: struct{}{}}},
	}})
}