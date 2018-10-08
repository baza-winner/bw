package deftype

import (
	"github.com/baza-winner/bwcore/bwjson"
	"sort"
)

type Item uint16

const (
	Item_below_ Item = iota
	Bool
	String
	Int
	Number
	Map
	Array
	ArrayOf
	Item_above_
)

//go:generate stringer -type=Item

type Set map[Item]struct{}

func FromArgs(kk ...Item) Set {
	return FromSlice(kk)
}

func FromSlice(kk []Item) Set {
	result := Set{}
	result.Add(kk...)
	return result
}

func (v Set) Has(k Item) (ok bool) {
	_, ok = v[k]
	return
}

func (v Set) Add(kk ...Item) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Set) String() string { return bwjson.PrettyJsonOf(v) }

func (v Set) GetDataForJson() interface{} {
	result := []interface{}{}
	for _, k := range v.ToSlice() {
		result = append(result, k.String())
	}
	return result
}

func (v Set) ToSlice() []Item {
	s := slice{}
	for k, _ := range v { s = append(s, k) }
	sort.Sort(s)
	return s
}

func (v Set) ToSliceOfStrings() (result []string) {
	result = []string{}
	for k, _ := range v { result = append(result, k.String()) }
	sort.Strings(result)
	return
}

// ============================================================================

type slice []Item

func (v slice) Len() int           { return len(v) }
func (v slice) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v slice) Less(i, j int) bool { return v[i] < v[j] }

