// Code generated by "bwsetter -type=uint8"; DO NOT EDIT; bwsetter: go get -type=uint8 -set=Uint8 -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	"encoding/json"
	bwjson "github.com/baza-winner/bwcore/bwjson"
	"sort"
	"strconv"
)

// Uint8 - множество значений типа uint8 с поддержкой интерфейсов Stringer и MarshalJSON
type Uint8 map[uint8]struct{}

// Uint8From - конструктор Uint8
func Uint8From(kk ...uint8) Uint8 {
	result := Uint8{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Uint8FromSlice - конструктор Uint8
func Uint8FromSlice(kk []uint8) Uint8 {
	result := Uint8{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Uint8FromSet - конструктор Uint8
func Uint8FromSet(s Uint8) Uint8 {
	result := Uint8{}
	for k, _ := range s {
		result[k] = struct{}{}
	}
	return result
}

// Copy - создает независимую копию
func (v Uint8) Copy() Uint8 {
	return Uint8FromSet(v)
}

// ToSlice - возвращает в виде []uint8
func (v Uint8) ToSlice() []uint8 {
	result := _uint8Slice{}
	for k, _ := range v {
		result = append(result, k)
	}
	sort.Sort(result)
	return result
}

func _Uint8ToSliceTestHelper(kk []uint8) []uint8 {
	return Uint8FromSlice(kk).ToSlice()
}

// String - поддержка интерфейса Stringer
func (v Uint8) String() string {
	return bwjson.Pretty(v)
}

// MarshalJSON - поддержка интерфейса MarshalJSON
func (v Uint8) MarshalJSON() ([]byte, error) {
	result := []interface{}{}
	for k, _ := range v {
		result = append(result, k)
	}
	return json.Marshal(result)
}

// ToSliceOfStrings - возвращает []string строковых представлений элементов множества
func (v Uint8) ToSliceOfStrings() []string {
	result := []string{}
	for k, _ := range v {
		result = append(result, strconv.FormatUint(uint64(k), 10))
	}
	sort.Strings(result)
	return result
}

// Has - возвращает true, если множество содержит заданный элемент, в противном случае - false
func (v Uint8) Has(k uint8) bool {
	_, ok := v[k]
	return ok
}

/*
HasAny - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v Uint8) HasAny(kk ...uint8) bool {
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
func (v Uint8) HasAnyOfSlice(kk []uint8) bool {
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
func (v Uint8) HasAnyOfSet(s Uint8) bool {
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
func (v Uint8) HasEach(kk ...uint8) bool {
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
func (v Uint8) HasEachOfSlice(kk []uint8) bool {
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
func (v Uint8) HasEachOfSet(s Uint8) bool {
	for k, _ := range s {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

// Add - добавляет элементы в множество v
func (v Uint8) Add(kk ...uint8) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Uint8) _AddTestHelper(kk ...uint8) Uint8 {
	result := v.Copy()
	result.Add(kk...)
	return result
}

// AddSlice - добавляет элементы в множество v
func (v Uint8) AddSlice(kk []uint8) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Uint8) _AddSliceTestHelper(kk []uint8) Uint8 {
	result := v.Copy()
	result.AddSlice(kk)
	return result
}

// AddSet - добавляет элементы в множество v
func (v Uint8) AddSet(s Uint8) {
	for k, _ := range s {
		v[k] = struct{}{}
	}
}

func (v Uint8) _AddSetTestHelper(s Uint8) Uint8 {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Del - удаляет элементы из множествa v
func (v Uint8) Del(kk ...uint8) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Uint8) _DelTestHelper(kk ...uint8) Uint8 {
	result := v.Copy()
	result.Del(kk...)
	return result
}

// DelSlice - удаляет элементы из множествa v
func (v Uint8) DelSlice(kk []uint8) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Uint8) _DelSliceTestHelper(kk []uint8) Uint8 {
	result := v.Copy()
	result.DelSlice(kk)
	return result
}

// DelSet - удаляет элементы из множествa v
func (v Uint8) DelSet(s Uint8) {
	for k, _ := range s {
		delete(v, k)
	}
}

func (v Uint8) _DelSetTestHelper(s Uint8) Uint8 {
	result := v.Copy()
	result.DelSet(s)
	return result
}

// Union - возвращает результат объединения двух множеств. Исходные множества остаются без изменений
func (v Uint8) Union(s Uint8) Uint8 {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Intersect - возвращает результат пересечения двух множеств. Исходные множества остаются без изменений
func (v Uint8) Intersect(s Uint8) Uint8 {
	result := Uint8{}
	for k, _ := range v {
		if _, ok := s[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// Subtract - возвращает результат вычитания двух множеств. Исходные множества остаются без изменений
func (v Uint8) Subtract(s Uint8) Uint8 {
	result := Uint8{}
	for k, _ := range v {
		if _, ok := s[k]; !ok {
			result[k] = struct{}{}
		}
	}
	return result
}

type _uint8Slice []uint8

func (v _uint8Slice) Len() int {
	return len(v)
}

func (v _uint8Slice) Swap(i int, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v _uint8Slice) Less(i int, j int) bool {
	return v[i] < v[j]
}
