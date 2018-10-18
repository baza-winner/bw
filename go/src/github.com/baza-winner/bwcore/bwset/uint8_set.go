// Code generated by "bwsetter -type=uint8"; DO NOT EDIT; bwsetter: go get -type=uint8 -set=Uint8Set -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	bwjson "github.com/baza-winner/bwcore/bwjson"
	"sort"
	"strconv"
)

// Uint8Set - множество значений типа uint8 с поддержкой интерфейсов Stringer и github.com/baza-winner/bwcore/bwjson.Jsonable
type Uint8Set map[uint8]struct{}

// Uint8SetFrom - конструктор Uint8Set
func Uint8SetFrom(kk ...uint8) Uint8Set {
	result := Uint8Set{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Uint8SetFromSlice - конструктор Uint8Set
func Uint8SetFromSlice(kk []uint8) Uint8Set {
	result := Uint8Set{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Uint8SetFromSet - конструктор Uint8Set
func Uint8SetFromSet(s Uint8Set) Uint8Set {
	result := Uint8Set{}
	for k, _ := range s {
		result[k] = struct{}{}
	}
	return result
}

// Copy - создает независимую копию
func (v Uint8Set) Copy() Uint8Set {
	return Uint8SetFromSet(v)
}

// ToSlice - возвращает в виде []uint8
func (v Uint8Set) ToSlice() []uint8 {
	result := _uint8Slice{}
	for k, _ := range v {
		result = append(result, k)
	}
	sort.Sort(result)
	return result
}

func _Uint8SetToSliceTestHelper(kk []uint8) []uint8 {
	return Uint8SetFromSlice(kk).ToSlice()
}

// String - поддержка интерфейса Stringer
func (v Uint8Set) String() string {
	return bwjson.PrettyJsonOf(v)
}

// DataForJSON - поддержка интерфейса bwjson.Jsonable
func (v Uint8Set) DataForJSON() interface{} {
	result := []interface{}{}
	for k, _ := range v {
		result = append(result, k)
	}
	return result
}

// ToSliceOfStrings - возвращает []string строковых представлений элементов множества
func (v Uint8Set) ToSliceOfStrings() []string {
	result := []string{}
	for k, _ := range v {
		result = append(result, strconv.FormatUint(uint64(k), 10))
	}
	sort.Strings(result)
	return result
}

// Has - возвращает true, если множество содержит заданный элемент, в противном случае - false
func (v Uint8Set) Has(k uint8) bool {
	_, ok := v[k]
	return ok
}

/*
HasAny - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v Uint8Set) HasAny(kk ...uint8) bool {
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
func (v Uint8Set) HasAnyOfSlice(kk []uint8) bool {
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
func (v Uint8Set) HasAnyOfSet(s Uint8Set) bool {
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
func (v Uint8Set) HasEach(kk ...uint8) bool {
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
func (v Uint8Set) HasEachOfSlice(kk []uint8) bool {
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
func (v Uint8Set) HasEachOfSet(s Uint8Set) bool {
	for k, _ := range s {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

// Add - добавляет элементы в множество v
func (v Uint8Set) Add(kk ...uint8) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Uint8Set) _AddTestHelper(kk ...uint8) Uint8Set {
	result := v.Copy()
	result.Add(kk...)
	return result
}

// AddSlice - добавляет элементы в множество v
func (v Uint8Set) AddSlice(kk []uint8) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Uint8Set) _AddSliceTestHelper(kk []uint8) Uint8Set {
	result := v.Copy()
	result.AddSlice(kk)
	return result
}

// AddSet - добавляет элементы в множество v
func (v Uint8Set) AddSet(s Uint8Set) {
	for k, _ := range s {
		v[k] = struct{}{}
	}
}

func (v Uint8Set) _AddSetTestHelper(s Uint8Set) Uint8Set {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Del - удаляет элементы из множествa v
func (v Uint8Set) Del(kk ...uint8) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Uint8Set) _DelTestHelper(kk ...uint8) Uint8Set {
	result := v.Copy()
	result.Del(kk...)
	return result
}

// DelSlice - удаляет элементы из множествa v
func (v Uint8Set) DelSlice(kk []uint8) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Uint8Set) _DelSliceTestHelper(kk []uint8) Uint8Set {
	result := v.Copy()
	result.DelSlice(kk)
	return result
}

// DelSet - удаляет элементы из множествa v
func (v Uint8Set) DelSet(s Uint8Set) {
	for k, _ := range s {
		delete(v, k)
	}
}

func (v Uint8Set) _DelSetTestHelper(s Uint8Set) Uint8Set {
	result := v.Copy()
	result.DelSet(s)
	return result
}

// Union - возвращает результат объединения двух множеств. Исходные множества остаются без изменений
func (v Uint8Set) Union(s Uint8Set) Uint8Set {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Intersect - возвращает результат пересечения двух множеств. Исходные множества остаются без изменений
func (v Uint8Set) Intersect(s Uint8Set) Uint8Set {
	result := Uint8Set{}
	for k, _ := range v {
		if _, ok := s[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// Subtract - возвращает результат вычитания двух множеств. Исходные множества остаются без изменений
func (v Uint8Set) Subtract(s Uint8Set) Uint8Set {
	result := Uint8Set{}
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
