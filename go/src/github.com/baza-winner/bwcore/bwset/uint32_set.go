// Code generated by "bwsetter -type=uint32"; DO NOT EDIT; bwsetter: go get -type=uint32 -set=Uint32 -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	"encoding/json"
	bwjson "github.com/baza-winner/bwcore/bwjson"
	"sort"
	"strconv"
)

// Uint32 - множество значений типа uint32 с поддержкой интерфейсов Stringer и github.com/baza-winner/bwcore/bwjson.Jsonable
type Uint32 map[uint32]struct{}

// Uint32From - конструктор Uint32
func Uint32From(kk ...uint32) Uint32 {
	result := Uint32{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Uint32FromSlice - конструктор Uint32
func Uint32FromSlice(kk []uint32) Uint32 {
	result := Uint32{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Uint32FromSet - конструктор Uint32
func Uint32FromSet(s Uint32) Uint32 {
	result := Uint32{}
	for k, _ := range s {
		result[k] = struct{}{}
	}
	return result
}

// Copy - создает независимую копию
func (v Uint32) Copy() Uint32 {
	return Uint32FromSet(v)
}

// ToSlice - возвращает в виде []uint32
func (v Uint32) ToSlice() []uint32 {
	result := _uint32Slice{}
	for k, _ := range v {
		result = append(result, k)
	}
	sort.Sort(result)
	return result
}

func _Uint32ToSliceTestHelper(kk []uint32) []uint32 {
	return Uint32FromSlice(kk).ToSlice()
}

// String - поддержка интерфейса Stringer
func (v Uint32) String() string {
	return bwjson.Pretty(v)
}

// MarshalJSON - поддержка интерфейса MarshalJSON
func (v Uint32) MarshalJSON() ([]byte, error) {
	result := []interface{}{}
	for k, _ := range v {
		result = append(result, k)
	}
	return json.Marshal(result)
}

// ToSliceOfStrings - возвращает []string строковых представлений элементов множества
func (v Uint32) ToSliceOfStrings() []string {
	result := []string{}
	for k, _ := range v {
		result = append(result, strconv.FormatUint(uint64(k), 10))
	}
	sort.Strings(result)
	return result
}

// Has - возвращает true, если множество содержит заданный элемент, в противном случае - false
func (v Uint32) Has(k uint32) bool {
	_, ok := v[k]
	return ok
}

/*
HasAny - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v Uint32) HasAny(kk ...uint32) bool {
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
func (v Uint32) HasAnyOfSlice(kk []uint32) bool {
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
func (v Uint32) HasAnyOfSet(s Uint32) bool {
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
func (v Uint32) HasEach(kk ...uint32) bool {
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
func (v Uint32) HasEachOfSlice(kk []uint32) bool {
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
func (v Uint32) HasEachOfSet(s Uint32) bool {
	for k, _ := range s {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

// Add - добавляет элементы в множество v
func (v Uint32) Add(kk ...uint32) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Uint32) _AddTestHelper(kk ...uint32) Uint32 {
	result := v.Copy()
	result.Add(kk...)
	return result
}

// AddSlice - добавляет элементы в множество v
func (v Uint32) AddSlice(kk []uint32) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Uint32) _AddSliceTestHelper(kk []uint32) Uint32 {
	result := v.Copy()
	result.AddSlice(kk)
	return result
}

// AddSet - добавляет элементы в множество v
func (v Uint32) AddSet(s Uint32) {
	for k, _ := range s {
		v[k] = struct{}{}
	}
}

func (v Uint32) _AddSetTestHelper(s Uint32) Uint32 {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Del - удаляет элементы из множествa v
func (v Uint32) Del(kk ...uint32) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Uint32) _DelTestHelper(kk ...uint32) Uint32 {
	result := v.Copy()
	result.Del(kk...)
	return result
}

// DelSlice - удаляет элементы из множествa v
func (v Uint32) DelSlice(kk []uint32) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Uint32) _DelSliceTestHelper(kk []uint32) Uint32 {
	result := v.Copy()
	result.DelSlice(kk)
	return result
}

// DelSet - удаляет элементы из множествa v
func (v Uint32) DelSet(s Uint32) {
	for k, _ := range s {
		delete(v, k)
	}
}

func (v Uint32) _DelSetTestHelper(s Uint32) Uint32 {
	result := v.Copy()
	result.DelSet(s)
	return result
}

// Union - возвращает результат объединения двух множеств. Исходные множества остаются без изменений
func (v Uint32) Union(s Uint32) Uint32 {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Intersect - возвращает результат пересечения двух множеств. Исходные множества остаются без изменений
func (v Uint32) Intersect(s Uint32) Uint32 {
	result := Uint32{}
	for k, _ := range v {
		if _, ok := s[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// Subtract - возвращает результат вычитания двух множеств. Исходные множества остаются без изменений
func (v Uint32) Subtract(s Uint32) Uint32 {
	result := Uint32{}
	for k, _ := range v {
		if _, ok := s[k]; !ok {
			result[k] = struct{}{}
		}
	}
	return result
}

type _uint32Slice []uint32

func (v _uint32Slice) Len() int {
	return len(v)
}

func (v _uint32Slice) Swap(i int, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v _uint32Slice) Less(i int, j int) bool {
	return v[i] < v[j]
}
