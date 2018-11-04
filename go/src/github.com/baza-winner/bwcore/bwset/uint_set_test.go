// Code generated by "bwsetter -type=uint"; DO NOT EDIT; bwsetter: go get -type=uint -set=Uint -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	"encoding/json"
	bwjson "github.com/baza-winner/bwcore/bwjson"
	bwtesting "github.com/baza-winner/bwcore/bwtesting"
	"strconv"
	"testing"
)

func TestUint(t *testing.T) {
	bwtesting.BwRunTests(t, UintFrom, map[string]bwtesting.TestCaseStruct{"UintFrom": {
		In: []interface{}{[]uint{_UintTestItemA, _UintTestItemB}},
		Out: []interface{}{Uint{
			_UintTestItemA: struct{}{},
			_UintTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, UintFromSlice, map[string]bwtesting.TestCaseStruct{"UintFromSlice": {
		In: []interface{}{[]uint{_UintTestItemA, _UintTestItemB}},
		Out: []interface{}{Uint{
			_UintTestItemA: struct{}{},
			_UintTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, UintFromSet, map[string]bwtesting.TestCaseStruct{"UintFromSet": {
		In: []interface{}{Uint{
			_UintTestItemA: struct{}{},
			_UintTestItemB: struct{}{},
		}},
		Out: []interface{}{Uint{
			_UintTestItemA: struct{}{},
			_UintTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint.Copy, map[string]bwtesting.TestCaseStruct{"Uint.Copy": {
		In: []interface{}{Uint{
			_UintTestItemA: struct{}{},
			_UintTestItemB: struct{}{},
		}},
		Out: []interface{}{Uint{
			_UintTestItemA: struct{}{},
			_UintTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint.ToSlice, map[string]bwtesting.TestCaseStruct{"Uint.ToSlice": {
		In:  []interface{}{Uint{_UintTestItemA: struct{}{}}},
		Out: []interface{}{[]uint{_UintTestItemA}},
	}})
	bwtesting.BwRunTests(t, _UintToSliceTestHelper, map[string]bwtesting.TestCaseStruct{"_UintToSliceTestHelper": {
		In:  []interface{}{[]uint{_UintTestItemB, _UintTestItemA}},
		Out: []interface{}{[]uint{_UintTestItemA, _UintTestItemB}},
	}})
	bwtesting.BwRunTests(t, Uint.String, map[string]bwtesting.TestCaseStruct{"Uint.String": {
		In:  []interface{}{Uint{_UintTestItemA: struct{}{}}},
		Out: []interface{}{bwjson.Pretty([]uint{_UintTestItemA})},
	}})
	bwtesting.BwRunTests(t, Uint.MarshalJSON, map[string]bwtesting.TestCaseStruct{"Uint.MarshalJSON": {
		In: []interface{}{Uint{_UintTestItemA: struct{}{}}},
		Out: []interface{}{(func() []byte {
			result, _ := json.Marshal([]interface{}{_UintTestItemA})
			return result
		})(), nil},
	}})
	bwtesting.BwRunTests(t, Uint.ToSliceOfStrings, map[string]bwtesting.TestCaseStruct{"Uint.ToSliceOfStrings": {
		In:  []interface{}{Uint{_UintTestItemA: struct{}{}}},
		Out: []interface{}{[]string{strconv.FormatUint(uint64(_UintTestItemA), 10)}},
	}})
	bwtesting.BwRunTests(t, Uint.Has, map[string]bwtesting.TestCaseStruct{
		"Uint.Has: false": {
			In:  []interface{}{Uint{_UintTestItemA: struct{}{}}, _UintTestItemB},
			Out: []interface{}{false},
		},
		"Uint.Has: true": {
			In:  []interface{}{Uint{_UintTestItemA: struct{}{}}, _UintTestItemA},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint.HasAny, map[string]bwtesting.TestCaseStruct{
		"Uint.HasAny: empty": {
			In:  []interface{}{Uint{_UintTestItemA: struct{}{}}, []uint{}},
			Out: []interface{}{false},
		},
		"Uint.HasAny: false": {
			In:  []interface{}{Uint{_UintTestItemA: struct{}{}}, []uint{_UintTestItemB}},
			Out: []interface{}{false},
		},
		"Uint.HasAny: true": {
			In:  []interface{}{Uint{_UintTestItemA: struct{}{}}, []uint{_UintTestItemA, _UintTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint.HasAnyOfSlice, map[string]bwtesting.TestCaseStruct{
		"Uint.HasAnyOfSlice: empty": {
			In:  []interface{}{Uint{_UintTestItemA: struct{}{}}, []uint{}},
			Out: []interface{}{false},
		},
		"Uint.HasAnyOfSlice: false": {
			In:  []interface{}{Uint{_UintTestItemA: struct{}{}}, []uint{_UintTestItemB}},
			Out: []interface{}{false},
		},
		"Uint.HasAnyOfSlice: true": {
			In:  []interface{}{Uint{_UintTestItemA: struct{}{}}, []uint{_UintTestItemA, _UintTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint.HasAnyOfSet, map[string]bwtesting.TestCaseStruct{
		"Uint.HasAnyOfSet: empty": {
			In:  []interface{}{Uint{_UintTestItemA: struct{}{}}, Uint{}},
			Out: []interface{}{false},
		},
		"Uint.HasAnyOfSet: false": {
			In:  []interface{}{Uint{_UintTestItemA: struct{}{}}, Uint{_UintTestItemB: struct{}{}}},
			Out: []interface{}{false},
		},
		"Uint.HasAnyOfSet: true": {
			In: []interface{}{Uint{_UintTestItemA: struct{}{}}, Uint{
				_UintTestItemA: struct{}{},
				_UintTestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint.HasEach, map[string]bwtesting.TestCaseStruct{
		"Uint.HasEach: empty": {
			In:  []interface{}{Uint{_UintTestItemA: struct{}{}}, []uint{}},
			Out: []interface{}{true},
		},
		"Uint.HasEach: false": {
			In:  []interface{}{Uint{_UintTestItemA: struct{}{}}, []uint{_UintTestItemA, _UintTestItemB}},
			Out: []interface{}{false},
		},
		"Uint.HasEach: true": {
			In: []interface{}{Uint{
				_UintTestItemA: struct{}{},
				_UintTestItemB: struct{}{},
			}, []uint{_UintTestItemA, _UintTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint.HasEachOfSlice, map[string]bwtesting.TestCaseStruct{
		"Uint.HasEachOfSlice: empty": {
			In:  []interface{}{Uint{_UintTestItemA: struct{}{}}, []uint{}},
			Out: []interface{}{true},
		},
		"Uint.HasEachOfSlice: false": {
			In:  []interface{}{Uint{_UintTestItemA: struct{}{}}, []uint{_UintTestItemA, _UintTestItemB}},
			Out: []interface{}{false},
		},
		"Uint.HasEachOfSlice: true": {
			In: []interface{}{Uint{
				_UintTestItemA: struct{}{},
				_UintTestItemB: struct{}{},
			}, []uint{_UintTestItemA, _UintTestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint.HasEachOfSet, map[string]bwtesting.TestCaseStruct{
		"Uint.HasEachOfSet: empty": {
			In:  []interface{}{Uint{_UintTestItemA: struct{}{}}, Uint{}},
			Out: []interface{}{true},
		},
		"Uint.HasEachOfSet: false": {
			In: []interface{}{Uint{_UintTestItemA: struct{}{}}, Uint{
				_UintTestItemA: struct{}{},
				_UintTestItemB: struct{}{},
			}},
			Out: []interface{}{false},
		},
		"Uint.HasEachOfSet: true": {
			In: []interface{}{Uint{
				_UintTestItemA: struct{}{},
				_UintTestItemB: struct{}{},
			}, Uint{
				_UintTestItemA: struct{}{},
				_UintTestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint._AddTestHelper, map[string]bwtesting.TestCaseStruct{"Uint.Add": {
		In: []interface{}{Uint{_UintTestItemA: struct{}{}}, []uint{_UintTestItemB}},
		Out: []interface{}{Uint{
			_UintTestItemA: struct{}{},
			_UintTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint._AddSliceTestHelper, map[string]bwtesting.TestCaseStruct{"Uint.AddSlice": {
		In: []interface{}{Uint{_UintTestItemA: struct{}{}}, []uint{_UintTestItemB}},
		Out: []interface{}{Uint{
			_UintTestItemA: struct{}{},
			_UintTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint._AddSetTestHelper, map[string]bwtesting.TestCaseStruct{"Uint.AddSet": {
		In: []interface{}{Uint{_UintTestItemA: struct{}{}}, Uint{_UintTestItemB: struct{}{}}},
		Out: []interface{}{Uint{
			_UintTestItemA: struct{}{},
			_UintTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint._DelTestHelper, map[string]bwtesting.TestCaseStruct{"Uint.Del": {
		In: []interface{}{Uint{
			_UintTestItemA: struct{}{},
			_UintTestItemB: struct{}{},
		}, []uint{_UintTestItemB}},
		Out: []interface{}{Uint{_UintTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Uint._DelSliceTestHelper, map[string]bwtesting.TestCaseStruct{"Uint.DelSlice": {
		In: []interface{}{Uint{
			_UintTestItemA: struct{}{},
			_UintTestItemB: struct{}{},
		}, []uint{_UintTestItemB}},
		Out: []interface{}{Uint{_UintTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Uint._DelSetTestHelper, map[string]bwtesting.TestCaseStruct{"Uint.DelSet": {
		In: []interface{}{Uint{
			_UintTestItemA: struct{}{},
			_UintTestItemB: struct{}{},
		}, Uint{_UintTestItemB: struct{}{}}},
		Out: []interface{}{Uint{_UintTestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Uint.Union, map[string]bwtesting.TestCaseStruct{"Uint.Union": {
		In: []interface{}{Uint{_UintTestItemA: struct{}{}}, Uint{_UintTestItemB: struct{}{}}},
		Out: []interface{}{Uint{
			_UintTestItemA: struct{}{},
			_UintTestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint.Intersect, map[string]bwtesting.TestCaseStruct{"Uint.Intersect": {
		In: []interface{}{Uint{
			_UintTestItemA: struct{}{},
			_UintTestItemB: struct{}{},
		}, Uint{_UintTestItemB: struct{}{}}},
		Out: []interface{}{Uint{_UintTestItemB: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Uint.Subtract, map[string]bwtesting.TestCaseStruct{"Uint.Subtract": {
		In: []interface{}{Uint{
			_UintTestItemA: struct{}{},
			_UintTestItemB: struct{}{},
		}, Uint{_UintTestItemB: struct{}{}}},
		Out: []interface{}{Uint{_UintTestItemA: struct{}{}}},
	}})
}
