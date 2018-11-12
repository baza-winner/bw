// Code generated by "bwsetter -type=float32"; DO NOT EDIT; bwsetter: go get -type=float32 -set=Float32 -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	"encoding/json"
	bwjson "github.com/baza-winner/bwcore/bwjson"
	"sort"
	"strconv"
)

// Float32 - множество значений типа float32 с поддержкой интерфейсов Stringer и MarshalJSON
type Float32 map[float32]struct{}

// Float32From - конструктор Float32
func Float32From(kk ...float32) Float32 {
	result := Float32{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Float32FromSlice - конструктор Float32
func Float32FromSlice(kk []float32) Float32 {
	result := Float32{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Float32FromSet - конструктор Float32
func Float32FromSet(s Float32) Float32 {
	result := Float32{}
	for k, _ := range s {
		result[k] = struct{}{}
	}
	return result
}

// Copy - создает независимую копию
func (v Float32) Copy() Float32 {
	return Float32FromSet(v)
}

// ToSlice - возвращает в виде []float32
func (v Float32) ToSlice() []float32 {
	result := _float32Slice{}
	for k, _ := range v {
		result = append(result, k)
	}
	sort.Sort(result)
	return result
}

func _Float32ToSliceTestHelper(kk []float32) []float32 {
	return Float32FromSlice(kk).ToSlice()
}

// String - поддержка интерфейса Stringer
func (v Float32) String() string {
	return bwjson.Pretty(v)
}

// MarshalJSON - поддержка интерфейса MarshalJSON
func (v Float32) MarshalJSON() ([]byte, error) {
	result := []interface{}{}
	for k, _ := range v {
		result = append(result, k)
	}
	return json.Marshal(result)
}

// ToSliceOfStrings - возвращает []string строковых представлений элементов множества
func (v Float32) ToSliceOfStrings() []string {
	result := []string{}
	for k, _ := range v {
		result = append(result, strconv.FormatFloat(float64(k), byte(0x66), -1, 64))
	}
	sort.Strings(result)
	return result
}

// Has - возвращает true, если множество содержит заданный элемент, в противном случае - false
func (v Float32) Has(k float32) bool {
	_, ok := v[k]
	return ok
}

/*
HasAny - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v Float32) HasAny(kk ...float32) bool {
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
func (v Float32) HasAnyOfSlice(kk []float32) bool {
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
func (v Float32) HasAnyOfSet(s Float32) bool {
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
func (v Float32) HasEach(kk ...float32) bool {
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
func (v Float32) HasEachOfSlice(kk []float32) bool {
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
func (v Float32) HasEachOfSet(s Float32) bool {
	for k, _ := range s {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

// Add - добавляет элементы в множество v
func (v Float32) Add(kk ...float32) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Float32) _AddTestHelper(kk ...float32) Float32 {
	result := v.Copy()
	result.Add(kk...)
	return result
}

// AddSlice - добавляет элементы в множество v
func (v Float32) AddSlice(kk []float32) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Float32) _AddSliceTestHelper(kk []float32) Float32 {
	result := v.Copy()
	result.AddSlice(kk)
	return result
}

// AddSet - добавляет элементы в множество v
func (v Float32) AddSet(s Float32) {
	for k, _ := range s {
		v[k] = struct{}{}
	}
}

func (v Float32) _AddSetTestHelper(s Float32) Float32 {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Del - удаляет элементы из множествa v
func (v Float32) Del(kk ...float32) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Float32) _DelTestHelper(kk ...float32) Float32 {
	result := v.Copy()
	result.Del(kk...)
	return result
}

// DelSlice - удаляет элементы из множествa v
func (v Float32) DelSlice(kk []float32) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Float32) _DelSliceTestHelper(kk []float32) Float32 {
	result := v.Copy()
	result.DelSlice(kk)
	return result
}

// DelSet - удаляет элементы из множествa v
func (v Float32) DelSet(s Float32) {
	for k, _ := range s {
		delete(v, k)
	}
}

func (v Float32) _DelSetTestHelper(s Float32) Float32 {
	result := v.Copy()
	result.DelSet(s)
	return result
}

// Union - возвращает результат объединения двух множеств. Исходные множества остаются без изменений
func (v Float32) Union(s Float32) Float32 {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Intersect - возвращает результат пересечения двух множеств. Исходные множества остаются без изменений
func (v Float32) Intersect(s Float32) Float32 {
	result := Float32{}
	for k, _ := range v {
		if _, ok := s[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// Subtract - возвращает результат вычитания двух множеств. Исходные множества остаются без изменений
func (v Float32) Subtract(s Float32) Float32 {
	result := Float32{}
	for k, _ := range v {
		if _, ok := s[k]; !ok {
			result[k] = struct{}{}
		}
	}
	return result
}

type _float32Slice []float32

func (v _float32Slice) Len() int {
	return len(v)
}

func (v _float32Slice) Swap(i int, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v _float32Slice) Less(i int, j int) bool {
	return v[i] < v[j]
}
