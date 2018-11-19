// Code generated by "bwsetter -type=int8"; DO NOT EDIT; bwsetter: go get -type=int8 -set=Int8 -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	"encoding/json"
	"sort"
	"strconv"
)

// Int8 - множество значений типа int8 с поддержкой интерфейсов Stringer и MarshalJSON
type Int8 map[int8]struct{}

// Int8From - конструктор Int8
func Int8From(kk ...int8) Int8 {
	result := Int8{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Int8FromSlice - конструктор Int8
func Int8FromSlice(kk []int8) Int8 {
	result := Int8{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Int8FromSet - конструктор Int8
func Int8FromSet(s Int8) Int8 {
	result := Int8{}
	for k, _ := range s {
		result[k] = struct{}{}
	}
	return result
}

// Copy - создает независимую копию
func (v Int8) Copy() Int8 {
	return Int8FromSet(v)
}

// ToSlice - возвращает в виде []int8
func (v Int8) ToSlice() []int8 {
	result := _int8Slice{}
	for k, _ := range v {
		result = append(result, k)
	}
	sort.Sort(result)
	return result
}

func _Int8ToSliceTestHelper(kk []int8) []int8 {
	return Int8FromSlice(kk).ToSlice()
}

// String - поддержка интерфейса Stringer
func (v Int8) String() string {
	result, _ := json.Marshal(v)
	return string(result)
}

// MarshalJSON - поддержка интерфейса MarshalJSON
func (v Int8) MarshalJSON() ([]byte, error) {
	result := []interface{}{}
	for _, k := range v.ToSlice() {
		result = append(result, k)
	}
	return json.Marshal(result)
}

// ToSliceOfStrings - возвращает []string строковых представлений элементов множества
func (v Int8) ToSliceOfStrings() []string {
	result := []string{}
	for k, _ := range v {
		result = append(result, strconv.FormatInt(int64(k), 10))
	}
	sort.Strings(result)
	return result
}

// Has - возвращает true, если множество содержит заданный элемент, в противном случае - false
func (v Int8) Has(k int8) bool {
	_, ok := v[k]
	return ok
}

/*
HasAny - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v Int8) HasAny(kk ...int8) bool {
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
func (v Int8) HasAnyOfSlice(kk []int8) bool {
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
func (v Int8) HasAnyOfSet(s Int8) bool {
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
func (v Int8) HasEach(kk ...int8) bool {
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
func (v Int8) HasEachOfSlice(kk []int8) bool {
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
func (v Int8) HasEachOfSet(s Int8) bool {
	for k, _ := range s {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

// Add - добавляет элементы в множество v
func (v Int8) Add(kk ...int8) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Int8) _AddTestHelper(kk ...int8) Int8 {
	result := v.Copy()
	result.Add(kk...)
	return result
}

// AddSlice - добавляет элементы в множество v
func (v Int8) AddSlice(kk []int8) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Int8) _AddSliceTestHelper(kk []int8) Int8 {
	result := v.Copy()
	result.AddSlice(kk)
	return result
}

// AddSet - добавляет элементы в множество v
func (v Int8) AddSet(s Int8) {
	for k, _ := range s {
		v[k] = struct{}{}
	}
}

func (v Int8) _AddSetTestHelper(s Int8) Int8 {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Del - удаляет элементы из множествa v
func (v Int8) Del(kk ...int8) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Int8) _DelTestHelper(kk ...int8) Int8 {
	result := v.Copy()
	result.Del(kk...)
	return result
}

// DelSlice - удаляет элементы из множествa v
func (v Int8) DelSlice(kk []int8) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Int8) _DelSliceTestHelper(kk []int8) Int8 {
	result := v.Copy()
	result.DelSlice(kk)
	return result
}

// DelSet - удаляет элементы из множествa v
func (v Int8) DelSet(s Int8) {
	for k, _ := range s {
		delete(v, k)
	}
}

func (v Int8) _DelSetTestHelper(s Int8) Int8 {
	result := v.Copy()
	result.DelSet(s)
	return result
}

// Union - возвращает результат объединения двух множеств. Исходные множества остаются без изменений
func (v Int8) Union(s Int8) Int8 {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Intersect - возвращает результат пересечения двух множеств. Исходные множества остаются без изменений
func (v Int8) Intersect(s Int8) Int8 {
	result := Int8{}
	for k, _ := range v {
		if _, ok := s[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// Subtract - возвращает результат вычитания двух множеств. Исходные множества остаются без изменений
func (v Int8) Subtract(s Int8) Int8 {
	result := Int8{}
	for k, _ := range v {
		if _, ok := s[k]; !ok {
			result[k] = struct{}{}
		}
	}
	return result
}

type _int8Slice []int8

func (v _int8Slice) Len() int {
	return len(v)
}

func (v _int8Slice) Swap(i int, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v _int8Slice) Less(i int, j int) bool {
	return v[i] < v[j]
}
