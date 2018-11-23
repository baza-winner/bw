// Code generated by "bwsetter -type=ParseValKind"; DO NOT EDIT; bwsetter: go get -type ParseValKind -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwparse

import (
	"encoding/json"
	"sort"
)

// ParseValKindSet - множество значений типа ParseValKind с поддержкой интерфейсов Stringer и MarshalJSON
type ParseValKindSet map[ParseValKind]struct{}

// ParseValKindSetFrom - конструктор ParseValKindSet
func ParseValKindSetFrom(kk ...ParseValKind) ParseValKindSet {
	result := ParseValKindSet{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// ParseValKindSetFromSlice - конструктор ParseValKindSet
func ParseValKindSetFromSlice(kk []ParseValKind) ParseValKindSet {
	result := ParseValKindSet{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// ParseValKindSetFromSet - конструктор ParseValKindSet
func ParseValKindSetFromSet(s ParseValKindSet) ParseValKindSet {
	result := ParseValKindSet{}
	for k, _ := range s {
		result[k] = struct{}{}
	}
	return result
}

// Copy - создает независимую копию
func (v ParseValKindSet) Copy() ParseValKindSet {
	return ParseValKindSetFromSet(v)
}

// ToSlice - возвращает в виде []ParseValKind
func (v ParseValKindSet) ToSlice() []ParseValKind {
	result := _ParseValKindSlice{}
	for k, _ := range v {
		result = append(result, k)
	}
	sort.Sort(result)
	return result
}

func _ParseValKindSetToSliceTestHelper(kk []ParseValKind) []ParseValKind {
	return ParseValKindSetFromSlice(kk).ToSlice()
}

// String - поддержка интерфейса Stringer
func (v ParseValKindSet) String() string {
	result, _ := json.Marshal(v)
	return string(result)
}

// MarshalJSON - поддержка интерфейса MarshalJSON
func (v ParseValKindSet) MarshalJSON() ([]byte, error) {
	result := []interface{}{}
	for _, k := range v.ToSlice() {
		result = append(result, k)
	}
	return json.Marshal(result)
}

// ToSliceOfStrings - возвращает []string строковых представлений элементов множества
func (v ParseValKindSet) ToSliceOfStrings() []string {
	result := []string{}
	for k, _ := range v {
		result = append(result, k.String())
	}
	sort.Strings(result)
	return result
}

// Has - возвращает true, если множество содержит заданный элемент, в противном случае - false
func (v ParseValKindSet) Has(k ParseValKind) bool {
	_, ok := v[k]
	return ok
}

/*
HasAny - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v ParseValKindSet) HasAny(kk ...ParseValKind) bool {
	for _, k := range kk {
		if _, ok := v[k]; ok {
			return true
		}
	}
	return false
}

/*
HasAnyOfSlice - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v ParseValKindSet) HasAnyOfSlice(kk []ParseValKind) bool {
	for _, k := range kk {
		if _, ok := v[k]; ok {
			return true
		}
	}
	return false
}

/*
HasAnyOfSet - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v ParseValKindSet) HasAnyOfSet(s ParseValKindSet) bool {
	for k, _ := range s {
		if _, ok := v[k]; ok {
			return true
		}
	}
	return false
}

/*
HasEach - возвращает true, если множество содержит все заданные элементы, в противном случае - false.
HasEach(<пустой набор/множесто>) возвращает true
*/
func (v ParseValKindSet) HasEach(kk ...ParseValKind) bool {
	for _, k := range kk {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

/*
HasEachOfSlice - возвращает true, если множество содержит все заданные элементы, в противном случае - false.
HasEach(<пустой набор/множесто>) возвращает true
*/
func (v ParseValKindSet) HasEachOfSlice(kk []ParseValKind) bool {
	for _, k := range kk {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

/*
HasEachOfSet - возвращает true, если множество содержит все заданные элементы, в противном случае - false.
HasEach(<пустой набор/множесто>) возвращает true
*/
func (v ParseValKindSet) HasEachOfSet(s ParseValKindSet) bool {
	for k, _ := range s {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

// Add - добавляет элементы в множество v
func (v ParseValKindSet) Add(kk ...ParseValKind) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v ParseValKindSet) _AddTestHelper(kk ...ParseValKind) ParseValKindSet {
	result := v.Copy()
	result.Add(kk...)
	return result
}

// AddSlice - добавляет элементы в множество v
func (v ParseValKindSet) AddSlice(kk []ParseValKind) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v ParseValKindSet) _AddSliceTestHelper(kk []ParseValKind) ParseValKindSet {
	result := v.Copy()
	result.AddSlice(kk)
	return result
}

// AddSet - добавляет элементы в множество v
func (v ParseValKindSet) AddSet(s ParseValKindSet) {
	for k, _ := range s {
		v[k] = struct{}{}
	}
}

func (v ParseValKindSet) _AddSetTestHelper(s ParseValKindSet) ParseValKindSet {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Del - удаляет элементы из множествa v
func (v ParseValKindSet) Del(kk ...ParseValKind) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v ParseValKindSet) _DelTestHelper(kk ...ParseValKind) ParseValKindSet {
	result := v.Copy()
	result.Del(kk...)
	return result
}

// DelSlice - удаляет элементы из множествa v
func (v ParseValKindSet) DelSlice(kk []ParseValKind) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v ParseValKindSet) _DelSliceTestHelper(kk []ParseValKind) ParseValKindSet {
	result := v.Copy()
	result.DelSlice(kk)
	return result
}

// DelSet - удаляет элементы из множествa v
func (v ParseValKindSet) DelSet(s ParseValKindSet) {
	for k, _ := range s {
		delete(v, k)
	}
}

func (v ParseValKindSet) _DelSetTestHelper(s ParseValKindSet) ParseValKindSet {
	result := v.Copy()
	result.DelSet(s)
	return result
}

// Union - возвращает результат объединения двух множеств. Исходные множества остаются без изменений
func (v ParseValKindSet) Union(s ParseValKindSet) ParseValKindSet {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Intersect - возвращает результат пересечения двух множеств. Исходные множества остаются без изменений
func (v ParseValKindSet) Intersect(s ParseValKindSet) ParseValKindSet {
	result := ParseValKindSet{}
	for k, _ := range v {
		if _, ok := s[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// Subtract - возвращает результат вычитания двух множеств. Исходные множества остаются без изменений
func (v ParseValKindSet) Subtract(s ParseValKindSet) ParseValKindSet {
	result := ParseValKindSet{}
	for k, _ := range v {
		if _, ok := s[k]; !ok {
			result[k] = struct{}{}
		}
	}
	return result
}

type _ParseValKindSlice []ParseValKind

func (v _ParseValKindSlice) Len() int {
	return len(v)
}

func (v _ParseValKindSlice) Swap(i int, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v _ParseValKindSlice) Less(i int, j int) bool {
	return v[i] < v[j]
}
