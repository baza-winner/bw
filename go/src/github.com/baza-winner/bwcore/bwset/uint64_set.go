// Code generated by "setter -type=uint64"; DO NOT EDIT; setter: go get github.com/baza-winner/bwcore/setter

package bwset

import (
	bwjson "github.com/baza-winner/bwcore/bwjson"
	"sort"
	"strconv"
)

// Uint64Set - множество значений типа uint64 с поддержкой интерфейсов Stringer и github.com/baza-winner/bwcore/bwjson.Jsonable
type Uint64Set map[uint64]struct{}

// Uint64SetFrom - конструктор Uint64Set
func Uint64SetFrom(kk ...uint64) Uint64Set {
	result := Uint64Set{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Uint64SetFromSlice - конструктор Uint64Set
func Uint64SetFromSlice(kk []uint64) Uint64Set {
	result := Uint64Set{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Uint64SetFromSet - конструктор Uint64Set
func Uint64SetFromSet(s Uint64Set) Uint64Set {
	result := Uint64Set{}
	for k, _ := range s {
		result[k] = struct{}{}
	}
	return result
}

// Copy - создает независимую копию
func (v Uint64Set) Copy() Uint64Set {
	return Uint64SetFromSet(v)
}

// ToSlice - возвращает в виде []uint64
func (v Uint64Set) ToSlice() []uint64 {
	result := _uint64Slice{}
	for k, _ := range v {
		result = append(result, k)
	}
	sort.Sort(result)
	return result
}

func _Uint64SetToSliceTestHelper(kk []uint64) []uint64 {
	return Uint64SetFromSlice(kk).ToSlice()
}

// String - поддержка интерфейса Stringer
func (v Uint64Set) String() string {
	return bwjson.PrettyJsonOf(v)
}

// GetDataForJson - поддержка интерфейса bwjson.Jsonable
func (v Uint64Set) GetDataForJson() interface{} {
	result := []interface{}{}
	for k, _ := range v {
		result = append(result, strconv.FormatUint(uint64(k), 10))
	}
	return result
}

// ToSliceOfStrings - возвращает []string строковых представлений элементов множества
func (v Uint64Set) ToSliceOfStrings() []string {
	result := []string{}
	for k, _ := range v {
		result = append(result, strconv.FormatUint(uint64(k), 10))
	}
	sort.Strings(result)
	return result
}

// Has - возвращает true, если множество содержит заданный элемент, в противном случае - false
func (v Uint64Set) Has(k uint64) bool {
	_, ok := v[k]
	return ok
}

/*
HasAny - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v Uint64Set) HasAny(kk ...uint64) bool {
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
func (v Uint64Set) HasAnyOfSlice(kk []uint64) bool {
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
func (v Uint64Set) HasAnyOfSet(s Uint64Set) bool {
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
func (v Uint64Set) HasEach(kk ...uint64) bool {
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
func (v Uint64Set) HasEachOfSlice(kk []uint64) bool {
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
func (v Uint64Set) HasEachOfSet(s Uint64Set) bool {
	for k, _ := range s {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

// Add - добавляет элементы в множество v
func (v Uint64Set) Add(kk ...uint64) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Uint64Set) _AddTestHelper(kk ...uint64) Uint64Set {
	result := v.Copy()
	result.Add(kk...)
	return result
}

// AddSlice - добавляет элементы в множество v
func (v Uint64Set) AddSlice(kk []uint64) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Uint64Set) _AddSliceTestHelper(kk []uint64) Uint64Set {
	result := v.Copy()
	result.AddSlice(kk)
	return result
}

// AddSet - добавляет элементы в множество v
func (v Uint64Set) AddSet(s Uint64Set) {
	for k, _ := range s {
		v[k] = struct{}{}
	}
}

func (v Uint64Set) _AddSetTestHelper(s Uint64Set) Uint64Set {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Del - удаляет элементы из множествa v
func (v Uint64Set) Del(kk ...uint64) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Uint64Set) _DelTestHelper(kk ...uint64) Uint64Set {
	result := v.Copy()
	result.Del(kk...)
	return result
}

// DelSlice - удаляет элементы из множествa v
func (v Uint64Set) DelSlice(kk []uint64) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Uint64Set) _DelSliceTestHelper(kk []uint64) Uint64Set {
	result := v.Copy()
	result.DelSlice(kk)
	return result
}

// DelSet - удаляет элементы из множествa v
func (v Uint64Set) DelSet(s Uint64Set) {
	for k, _ := range s {
		delete(v, k)
	}
}

func (v Uint64Set) _DelSetTestHelper(s Uint64Set) Uint64Set {
	result := v.Copy()
	result.DelSet(s)
	return result
}

// Union - возвращает результат объединения двух множеств. Исходные множества остаются без изменений
func (v Uint64Set) Union(s Uint64Set) Uint64Set {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Intersect - возвращает результат пересечения двух множеств. Исходные множества остаются без изменений
func (v Uint64Set) Intersect(s Uint64Set) Uint64Set {
	result := Uint64Set{}
	for k, _ := range v {
		if _, ok := s[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// Subtract - возвращает результат вычитания двух множеств. Исходные множества остаются без изменений
func (v Uint64Set) Subtract(s Uint64Set) Uint64Set {
	result := Uint64Set{}
	for k, _ := range v {
		if _, ok := s[k]; !ok {
			result[k] = struct{}{}
		}
	}
	return result
}

type _uint64Slice []uint64

func (v _uint64Slice) Len() int {
	return len(v)
}

func (v _uint64Slice) Swap(i int, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v _uint64Slice) Less(i int, j int) bool {
	return v[i] < v[j]
}