// Code generated by "setter -type=uint16"; DO NOT EDIT; setter: go get github.com/baza-winner/bwcore/setter

package bwset

import (
	bwjson "github.com/baza-winner/bwcore/bwjson"
	"sort"
	"strconv"
)

// Uint16Set - множество значений типа uint16 с поддержкой интерфейсов Stringer и github.com/baza-winner/bwcore/bwjson.Jsonable
type Uint16Set map[uint16]struct{}

// Uint16SetFrom - конструктор Uint16Set
func Uint16SetFrom(kk ...uint16) Uint16Set {
	result := Uint16Set{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Uint16SetFromSlice - конструктор Uint16Set
func Uint16SetFromSlice(kk []uint16) Uint16Set {
	result := Uint16Set{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Uint16SetFromSet - конструктор Uint16Set
func Uint16SetFromSet(s Uint16Set) Uint16Set {
	result := Uint16Set{}
	for k, _ := range s {
		result[k] = struct{}{}
	}
	return result
}

// Copy - создает независимую копию
func (v Uint16Set) Copy() Uint16Set {
	return Uint16SetFromSet(v)
}

// ToSlice - возвращает в виде []uint16
func (v Uint16Set) ToSlice() []uint16 {
	result := _uint16Slice{}
	for k, _ := range v {
		result = append(result, k)
	}
	sort.Sort(result)
	return result
}

func _Uint16SetToSliceTestHelper(kk []uint16) []uint16 {
	return Uint16SetFromSlice(kk).ToSlice()
}

// String - поддержка интерфейса Stringer
func (v Uint16Set) String() string {
	return bwjson.PrettyJsonOf(v)
}

// GetDataForJson - поддержка интерфейса bwjson.Jsonable
func (v Uint16Set) GetDataForJson() interface{} {
	result := []interface{}{}
	for k, _ := range v {
		result = append(result, strconv.FormatUint(uint64(k), 10))
	}
	return result
}

// ToSliceOfStrings - возвращает []string строковых представлений элементов множества
func (v Uint16Set) ToSliceOfStrings() []string {
	result := []string{}
	for k, _ := range v {
		result = append(result, strconv.FormatUint(uint64(k), 10))
	}
	sort.Strings(result)
	return result
}

// Has - возвращает true, если множество содержит заданный элемент, в противном случае - false
func (v Uint16Set) Has(k uint16) bool {
	_, ok := v[k]
	return ok
}

/*
HasAny - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v Uint16Set) HasAny(kk ...uint16) bool {
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
func (v Uint16Set) HasAnyOfSlice(kk []uint16) bool {
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
func (v Uint16Set) HasAnyOfSet(s Uint16Set) bool {
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
func (v Uint16Set) HasEach(kk ...uint16) bool {
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
func (v Uint16Set) HasEachOfSlice(kk []uint16) bool {
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
func (v Uint16Set) HasEachOfSet(s Uint16Set) bool {
	for k, _ := range s {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

// Add - добавляет элементы в множество v
func (v Uint16Set) Add(kk ...uint16) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Uint16Set) _AddTestHelper(kk ...uint16) Uint16Set {
	result := v.Copy()
	result.Add(kk...)
	return result
}

// AddSlice - добавляет элементы в множество v
func (v Uint16Set) AddSlice(kk []uint16) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Uint16Set) _AddSliceTestHelper(kk []uint16) Uint16Set {
	result := v.Copy()
	result.AddSlice(kk)
	return result
}

// AddSet - добавляет элементы в множество v
func (v Uint16Set) AddSet(s Uint16Set) {
	for k, _ := range s {
		v[k] = struct{}{}
	}
}

func (v Uint16Set) _AddSetTestHelper(s Uint16Set) Uint16Set {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Del - удаляет элементы из множествa v
func (v Uint16Set) Del(kk ...uint16) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Uint16Set) _DelTestHelper(kk ...uint16) Uint16Set {
	result := v.Copy()
	result.Del(kk...)
	return result
}

// DelSlice - удаляет элементы из множествa v
func (v Uint16Set) DelSlice(kk []uint16) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Uint16Set) _DelSliceTestHelper(kk []uint16) Uint16Set {
	result := v.Copy()
	result.DelSlice(kk)
	return result
}

// DelSet - удаляет элементы из множествa v
func (v Uint16Set) DelSet(s Uint16Set) {
	for k, _ := range s {
		delete(v, k)
	}
}

func (v Uint16Set) _DelSetTestHelper(s Uint16Set) Uint16Set {
	result := v.Copy()
	result.DelSet(s)
	return result
}

// Union - возвращает результат объединения двух множеств. Исходные множества остаются без изменений
func (v Uint16Set) Union(s Uint16Set) Uint16Set {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Intersect - возвращает результат пересечения двух множеств. Исходные множества остаются без изменений
func (v Uint16Set) Intersect(s Uint16Set) Uint16Set {
	result := Uint16Set{}
	for k, _ := range v {
		if _, ok := s[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// Subtract - возвращает результат вычитания двух множеств. Исходные множества остаются без изменений
func (v Uint16Set) Subtract(s Uint16Set) Uint16Set {
	result := Uint16Set{}
	for k, _ := range v {
		if _, ok := s[k]; !ok {
			result[k] = struct{}{}
		}
	}
	return result
}

type _uint16Slice []uint16

func (v _uint16Slice) Len() int {
	return len(v)
}

func (v _uint16Slice) Swap(i int, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v _uint16Slice) Less(i int, j int) bool {
	return v[i] < v[j]
}
