package defvalid

import (
	"github.com/baza-winner/bwcore/bwjson"
	"sort"
)

type deftypeItem uint16

const (
	deftype_below_ deftypeItem = iota
	deftypeBool
	deftypeString
	deftypeInt
	deftypeNumber
	deftypeMap
	deftypeArray
  deftypeOrArrayOf
	deftype_above_
)

//go:generate stringer -type=deftypeItem

type deftype map[deftypeItem]struct{}

type deftypeSlice []deftypeItem

func (v deftypeSlice) Len() int           { return len(v) }
func (v deftypeSlice) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v deftypeSlice) Less(i, j int) bool { return v[i] < v[j] }

func (v deftype) GetDataForJson() interface{} {
	result := []interface{}{}
	for _, k := range v.ToSlice() {
		result = append(result, k.String())
	}
	return result
}

func FromArgs(kk ...deftypeItem) deftype {
  return FromSlice(kk)
}

func FromSlice(kk []deftypeItem) deftype {
  result := deftype{}
  result.Add(kk...)
  return result
}

func (v deftype) Has(k deftypeItem) (ok bool) {
	_, ok = v[k]
	return
}

func (v deftype) Add(kk ...deftypeItem) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v deftype) ToSlice() []deftypeItem {
	result := deftypeSlice{}
	for k, _ := range v {
		result = append(result, k)
	}
	sort.Sort(result)
	return result
}

func (v deftype) String() string {
	return bwjson.PrettyJsonOf(v)
}
