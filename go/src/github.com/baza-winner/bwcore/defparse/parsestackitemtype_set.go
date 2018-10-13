// Code generated by "setter -type=parseStackItemType"; DO NOT EDIT github.com/baza-winner/bwcore/setter

package defparse

import (
	bwjson "github.com/baza-winner/bwcore/bwjson"
	"sort"
)

type parseStackItemTypeSet map[parseStackItemType]struct{}

func parseStackItemTypeSetFromArgs(kk ...parseStackItemType) parseStackItemTypeSet {
	return parseStackItemTypeSetFromSlice(kk)
}

func parseStackItemTypeSetFromSlice(kk []parseStackItemType) parseStackItemTypeSet {
	result := parseStackItemTypeSet{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

func (v parseStackItemTypeSet) Copy() parseStackItemTypeSet {
	return parseStackItemTypeSetFromSlice(v.ToSlice())
}

func (v parseStackItemTypeSet) ToSlice() []parseStackItemType {
	result := _parseStackItemTypeSlice{}
	for k, _ := range v {
		result = append(result, k)
	}
	sort.Sort(result)
	return result
}

func (v parseStackItemTypeSet) String() string {
	return bwjson.PrettyJsonOf(v)
}

func (v parseStackItemTypeSet) GetDataForJson() interface{} {
	result := []interface{}{}
	for _, k := range v.ToSlice() {
		result = append(result, k.String())
	}
	return result
}

func (v parseStackItemTypeSet) ToSliceOfStrings() (result []string) {
	result = []string{}
	for k, _ := range v {
		result = append(result, k.String())
	}
	sort.Strings(result)
	return result
}

func (v parseStackItemTypeSet) ContainesEach(s parseStackItemTypeSet) bool {
	for k, _ := range s {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

func (v parseStackItemTypeSet) ContainsAny(s parseStackItemTypeSet) bool {
	for k, _ := range s {
		if _, ok := v[k]; ok {
			return true
		}
	}
	return false
}

func (v parseStackItemTypeSet) Union(s parseStackItemTypeSet) {
	for k, _ := range s {
		v[k] = struct{}{}
	}
}

func (v parseStackItemTypeSet) Intersect(s parseStackItemTypeSet) {
	for k, _ := range v {
		if _, ok := s[k]; !ok {
			delete(v, k)
		}
	}
	for k, _ := range s {
		if _, ok := v[k]; !ok {
			delete(v, k)
		}
	}
}

func (v parseStackItemTypeSet) Subtract(s parseStackItemTypeSet) {
	for k, _ := range v {
		if _, ok := s[k]; ok {
			delete(v, k)
		}
	}
}

func _parseStackItemTypeSetUnion(s1 parseStackItemTypeSet, s2 parseStackItemTypeSet) parseStackItemTypeSet {
	s1.Union(s2)
	return s1
}

func _parseStackItemTypeSetIntersect(s1 parseStackItemTypeSet, s2 parseStackItemTypeSet) parseStackItemTypeSet {
	s1.Intersect(s2)
	return s1
}

func _parseStackItemTypeSetSubtract(s1 parseStackItemTypeSet, s2 parseStackItemTypeSet) parseStackItemTypeSet {
	s1.Subtract(s2)
	return s1
}

type _parseStackItemTypeSlice []parseStackItemType

func (v _parseStackItemTypeSlice) Len() int {
	return len(v)
}

func (v _parseStackItemTypeSlice) Swap(i int, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v _parseStackItemTypeSlice) Less(i int, j int) bool {
	return v[i] < v[j]
}
