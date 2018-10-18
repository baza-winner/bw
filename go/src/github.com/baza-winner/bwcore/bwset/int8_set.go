// Code generated by "bwsetter -type=int8"; DO NOT EDIT; bwsetter: go get -type=int8 -set=Int8Set -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	bwjson "github.com/baza-winner/bwcore/bwjson"
	"sort"
	"strconv"
)

// Int8Set - множество значений типа int8 с поддержкой интерфейсов Stringer и github.com/baza-winner/bwcore/bwjson.Jsonable
type Int8Set map[int8]struct{}

// Int8SetFrom - конструктор Int8Set
func Int8SetFrom(kk ...int8) Int8Set {
	result := Int8Set{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Int8SetFromSlice - конструктор Int8Set
func Int8SetFromSlice(kk []int8) Int8Set {
	result := Int8Set{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Int8SetFromSet - конструктор Int8Set
func Int8SetFromSet(s Int8Set) Int8Set {
	result := Int8Set{}
	for k, _ := range s {
		result[k] = struct{}{}
	}
	return result
}

// Copy - создает независимую копию
func (v Int8Set) Copy() Int8Set {
	return Int8SetFromSet(v)
}

// ToSlice - возвращает в виде []int8
func (v Int8Set) ToSlice() []int8 {
	result := _int8Slice{}
	for k, _ := range v {
		result = append(result, k)
	}
	sort.Sort(result)
	return result
}

func _Int8SetToSliceTestHelper(kk []int8) []int8 {
	return Int8SetFromSlice(kk).ToSlice()
}

// String - поддержка интерфейса Stringer
func (v Int8Set) String() string {
	return bwjson.PrettyJsonOf(v)
}

// DataForJSON - поддержка интерфейса bwjson.Jsonable
func (v Int8Set) DataForJSON() interface{} {
	result := []interface{}{}
	for k, _ := range v {
		result = append(result, k)
	}
	return result
}

// ToSliceOfStrings - возвращает []string строковых представлений элементов множества
func (v Int8Set) ToSliceOfStrings() []string {
	result := []string{}
	for k, _ := range v {
		result = append(result, strconv.FormatInt(int64(k), 10))
	}
	sort.Strings(result)
	return result
}

// Has - возвращает true, если множество содержит заданный элемент, в противном случае - false
func (v Int8Set) Has(k int8) bool {
	_, ok := v[k]
	return ok
}

/*
HasAny - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v Int8Set) HasAny(kk ...int8) bool {
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
func (v Int8Set) HasAnyOfSlice(kk []int8) bool {
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
func (v Int8Set) HasAnyOfSet(s Int8Set) bool {
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
func (v Int8Set) HasEach(kk ...int8) bool {
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
func (v Int8Set) HasEachOfSlice(kk []int8) bool {
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
func (v Int8Set) HasEachOfSet(s Int8Set) bool {
	for k, _ := range s {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

// Add - добавляет элементы в множество v
func (v Int8Set) Add(kk ...int8) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Int8Set) _AddTestHelper(kk ...int8) Int8Set {
	result := v.Copy()
	result.Add(kk...)
	return result
}

// AddSlice - добавляет элементы в множество v
func (v Int8Set) AddSlice(kk []int8) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Int8Set) _AddSliceTestHelper(kk []int8) Int8Set {
	result := v.Copy()
	result.AddSlice(kk)
	return result
}

// AddSet - добавляет элементы в множество v
func (v Int8Set) AddSet(s Int8Set) {
	for k, _ := range s {
		v[k] = struct{}{}
	}
}

func (v Int8Set) _AddSetTestHelper(s Int8Set) Int8Set {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Del - удаляет элементы из множествa v
func (v Int8Set) Del(kk ...int8) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Int8Set) _DelTestHelper(kk ...int8) Int8Set {
	result := v.Copy()
	result.Del(kk...)
	return result
}

// DelSlice - удаляет элементы из множествa v
func (v Int8Set) DelSlice(kk []int8) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Int8Set) _DelSliceTestHelper(kk []int8) Int8Set {
	result := v.Copy()
	result.DelSlice(kk)
	return result
}

// DelSet - удаляет элементы из множествa v
func (v Int8Set) DelSet(s Int8Set) {
	for k, _ := range s {
		delete(v, k)
	}
}

func (v Int8Set) _DelSetTestHelper(s Int8Set) Int8Set {
	result := v.Copy()
	result.DelSet(s)
	return result
}

// Union - возвращает результат объединения двух множеств. Исходные множества остаются без изменений
func (v Int8Set) Union(s Int8Set) Int8Set {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Intersect - возвращает результат пересечения двух множеств. Исходные множества остаются без изменений
func (v Int8Set) Intersect(s Int8Set) Int8Set {
	result := Int8Set{}
	for k, _ := range v {
		if _, ok := s[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// Subtract - возвращает результат вычитания двух множеств. Исходные множества остаются без изменений
func (v Int8Set) Subtract(s Int8Set) Int8Set {
	result := Int8Set{}
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
