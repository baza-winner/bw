// Code generated by "bwsetter -type=uint16"; DO NOT EDIT; bwsetter: go get -type=uint16 -set=Uint16 -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	"encoding/json"
	bwtesting "github.com/baza-winner/bwcore/bwtesting"
	"strconv"
	"testing"
)

func TestUint16(t *testing.T) {
	bwtesting.BwRunTests(t, Uint16From, map[string]bwtesting.Case{"Uint16From": {
		In: []interface{}{_Uint16TestItemA, _Uint16TestItemB},
		Out: []interface{}{Uint16{
			_Uint16TestItemA: struct{}{},
			_Uint16TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint16FromSlice, map[string]bwtesting.Case{"Uint16FromSlice": {
		In: []interface{}{[]uint16{_Uint16TestItemA, _Uint16TestItemB}},
		Out: []interface{}{Uint16{
			_Uint16TestItemA: struct{}{},
			_Uint16TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint16FromSet, map[string]bwtesting.Case{"Uint16FromSet": {
		In: []interface{}{Uint16{
			_Uint16TestItemA: struct{}{},
			_Uint16TestItemB: struct{}{},
		}},
		Out: []interface{}{Uint16{
			_Uint16TestItemA: struct{}{},
			_Uint16TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint16.Copy, map[string]bwtesting.Case{"Uint16.Copy": {
		In: []interface{}{Uint16{
			_Uint16TestItemA: struct{}{},
			_Uint16TestItemB: struct{}{},
		}},
		Out: []interface{}{Uint16{
			_Uint16TestItemA: struct{}{},
			_Uint16TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint16.ToSlice, map[string]bwtesting.Case{"Uint16.ToSlice": {
		In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}},
		Out: []interface{}{[]uint16{_Uint16TestItemA}},
	}})
	bwtesting.BwRunTests(t, _Uint16ToSliceTestHelper, map[string]bwtesting.Case{"_Uint16ToSliceTestHelper": {
		In:  []interface{}{[]uint16{_Uint16TestItemB, _Uint16TestItemA}},
		Out: []interface{}{[]uint16{_Uint16TestItemA, _Uint16TestItemB}},
	}})
	bwtesting.BwRunTests(t, Uint16.String, map[string]bwtesting.Case{"Uint16.String": {
		In: []interface{}{Uint16{_Uint16TestItemA: struct{}{}}},
		Out: []interface{}{func() string {
			result, _ := json.Marshal(_Uint16TestItemA)
			return "[" + string(result) + "]"
		}()},
	}})
	bwtesting.BwRunTests(t, Uint16.MarshalJSON, map[string]bwtesting.Case{"Uint16.MarshalJSON": {
		In: []interface{}{Uint16{_Uint16TestItemA: struct{}{}}},
		Out: []interface{}{(func() []byte {
			result, _ := json.Marshal([]interface{}{_Uint16TestItemA})
			return result
		})(), nil},
	}})
	bwtesting.BwRunTests(t, Uint16.ToSliceOfStrings, map[string]bwtesting.Case{"Uint16.ToSliceOfStrings": {
		In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}},
		Out: []interface{}{[]string{strconv.FormatUint(uint64(_Uint16TestItemA), 10)}},
	}})
	bwtesting.BwRunTests(t, Uint16.Has, map[string]bwtesting.Case{
		"Uint16.Has: false": {
			In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, _Uint16TestItemB},
			Out: []interface{}{false},
		},
		"Uint16.Has: true": {
			In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, _Uint16TestItemA},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint16.HasAny, map[string]bwtesting.Case{
		"Uint16.HasAny: empty": {
			In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}},
			Out: []interface{}{false},
		},
		"Uint16.HasAny: false": {
			In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, _Uint16TestItemB},
			Out: []interface{}{false},
		},
		"Uint16.HasAny: true": {
			In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, _Uint16TestItemA, _Uint16TestItemB},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint16.HasAnyOfSlice, map[string]bwtesting.Case{
		"Uint16.HasAnyOfSlice: empty": {
			In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, []uint16{}},
			Out: []interface{}{false},
		},
		"Uint16.HasAnyOfSlice: false": {
			In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, []uint16{_Uint16TestItemB}},
			Out: []interface{}{false},
		},
		"Uint16.HasAnyOfSlice: true": {
			In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, []uint16{_Uint16TestItemA, _Uint16TestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint16.HasAnyOfSet, map[string]bwtesting.Case{
		"Uint16.HasAnyOfSet: empty": {
			In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, Uint16{}},
			Out: []interface{}{false},
		},
		"Uint16.HasAnyOfSet: false": {
			In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, Uint16{_Uint16TestItemB: struct{}{}}},
			Out: []interface{}{false},
		},
		"Uint16.HasAnyOfSet: true": {
			In: []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, Uint16{
				_Uint16TestItemA: struct{}{},
				_Uint16TestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint16.HasEach, map[string]bwtesting.Case{
		"Uint16.HasEach: empty": {
			In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}},
			Out: []interface{}{true},
		},
		"Uint16.HasEach: false": {
			In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, _Uint16TestItemA, _Uint16TestItemB},
			Out: []interface{}{false},
		},
		"Uint16.HasEach: true": {
			In: []interface{}{Uint16{
				_Uint16TestItemA: struct{}{},
				_Uint16TestItemB: struct{}{},
			}, _Uint16TestItemA, _Uint16TestItemB},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint16.HasEachOfSlice, map[string]bwtesting.Case{
		"Uint16.HasEachOfSlice: empty": {
			In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, []uint16{}},
			Out: []interface{}{true},
		},
		"Uint16.HasEachOfSlice: false": {
			In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, []uint16{_Uint16TestItemA, _Uint16TestItemB}},
			Out: []interface{}{false},
		},
		"Uint16.HasEachOfSlice: true": {
			In: []interface{}{Uint16{
				_Uint16TestItemA: struct{}{},
				_Uint16TestItemB: struct{}{},
			}, []uint16{_Uint16TestItemA, _Uint16TestItemB}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint16.HasEachOfSet, map[string]bwtesting.Case{
		"Uint16.HasEachOfSet: empty": {
			In:  []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, Uint16{}},
			Out: []interface{}{true},
		},
		"Uint16.HasEachOfSet: false": {
			In: []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, Uint16{
				_Uint16TestItemA: struct{}{},
				_Uint16TestItemB: struct{}{},
			}},
			Out: []interface{}{false},
		},
		"Uint16.HasEachOfSet: true": {
			In: []interface{}{Uint16{
				_Uint16TestItemA: struct{}{},
				_Uint16TestItemB: struct{}{},
			}, Uint16{
				_Uint16TestItemA: struct{}{},
				_Uint16TestItemB: struct{}{},
			}},
			Out: []interface{}{true},
		},
	})
	bwtesting.BwRunTests(t, Uint16._AddTestHelper, map[string]bwtesting.Case{"Uint16.Add": {
		In: []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, _Uint16TestItemB},
		Out: []interface{}{Uint16{
			_Uint16TestItemA: struct{}{},
			_Uint16TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint16._AddSliceTestHelper, map[string]bwtesting.Case{"Uint16.AddSlice": {
		In: []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, []uint16{_Uint16TestItemB}},
		Out: []interface{}{Uint16{
			_Uint16TestItemA: struct{}{},
			_Uint16TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint16._AddSetTestHelper, map[string]bwtesting.Case{"Uint16.AddSet": {
		In: []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, Uint16{_Uint16TestItemB: struct{}{}}},
		Out: []interface{}{Uint16{
			_Uint16TestItemA: struct{}{},
			_Uint16TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint16._DelTestHelper, map[string]bwtesting.Case{"Uint16.Del": {
		In: []interface{}{Uint16{
			_Uint16TestItemA: struct{}{},
			_Uint16TestItemB: struct{}{},
		}, _Uint16TestItemB},
		Out: []interface{}{Uint16{_Uint16TestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Uint16._DelSliceTestHelper, map[string]bwtesting.Case{"Uint16.DelSlice": {
		In: []interface{}{Uint16{
			_Uint16TestItemA: struct{}{},
			_Uint16TestItemB: struct{}{},
		}, []uint16{_Uint16TestItemB}},
		Out: []interface{}{Uint16{_Uint16TestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Uint16._DelSetTestHelper, map[string]bwtesting.Case{"Uint16.DelSet": {
		In: []interface{}{Uint16{
			_Uint16TestItemA: struct{}{},
			_Uint16TestItemB: struct{}{},
		}, Uint16{_Uint16TestItemB: struct{}{}}},
		Out: []interface{}{Uint16{_Uint16TestItemA: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Uint16.Union, map[string]bwtesting.Case{"Uint16.Union": {
		In: []interface{}{Uint16{_Uint16TestItemA: struct{}{}}, Uint16{_Uint16TestItemB: struct{}{}}},
		Out: []interface{}{Uint16{
			_Uint16TestItemA: struct{}{},
			_Uint16TestItemB: struct{}{},
		}},
	}})
	bwtesting.BwRunTests(t, Uint16.Intersect, map[string]bwtesting.Case{"Uint16.Intersect": {
		In: []interface{}{Uint16{
			_Uint16TestItemA: struct{}{},
			_Uint16TestItemB: struct{}{},
		}, Uint16{_Uint16TestItemB: struct{}{}}},
		Out: []interface{}{Uint16{_Uint16TestItemB: struct{}{}}},
	}})
	bwtesting.BwRunTests(t, Uint16.Subtract, map[string]bwtesting.Case{"Uint16.Subtract": {
		In: []interface{}{Uint16{
			_Uint16TestItemA: struct{}{},
			_Uint16TestItemB: struct{}{},
		}, Uint16{_Uint16TestItemB: struct{}{}}},
		Out: []interface{}{Uint16{_Uint16TestItemA: struct{}{}}},
	}})
}
