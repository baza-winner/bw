// Code generated by "bwsetter -type=uint"; DO NOT EDIT; bwsetter: go get -type=uint -set=Uint -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	"encoding/json"
	bwjson "github.com/baza-winner/bwcore/bwjson"
	"sort"
	"strconv"
)

// Uint - множество значений типа uint с поддержкой интерфейсов Stringer и github.com/baza-winner/bwcore/bwjson.Jsonable
type Uint map[uint]struct{}

// UintFrom - конструктор Uint
func UintFrom(kk ...uint) Uint {
	result := Uint{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// UintFromSlice - конструктор Uint
func UintFromSlice(kk []uint) Uint {
	result := Uint{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// UintFromSet - конструктор Uint
func UintFromSet(s Uint) Uint {
	result := Uint{}
	for k, _ := range s {
		result[k] = struct{}{}
	}
	return result
}

// Copy - создает независимую копию
func (v Uint) Copy() Uint {
	return UintFromSet(v)
}

// ToSlice - возвращает в виде []uint
func (v Uint) ToSlice() []uint {
	result := _uintSlice{}
	for k, _ := range v {
		result = append(result, k)
	}
	sort.Sort(result)
	return result
}

func _UintToSliceTestHelper(kk []uint) []uint {
	return UintFromSlice(kk).ToSlice()
}

// String - поддержка интерфейса Stringer
func (v Uint) String() string {
	return bwjson.Pretty(v)
}

// MarshalJSON - поддержка интерфейса MarshalJSON
func (v Uint) MarshalJSON() ([]byte, error) {
	result := []interface{}{}
	for k, _ := range v {
		result = append(result, k)
	}
	return json.Marshal(result)
}

// ToSliceOfStrings - возвращает []string строковых представлений элементов множества
func (v Uint) ToSliceOfStrings() []string {
	result := []string{}
	for k, _ := range v {
		result = append(result, strconv.FormatUint(uint64(k), 10))
	}
	sort.Strings(result)
	return result
}

// Has - возвращает true, если множество содержит заданный элемент, в противном случае - false
func (v Uint) Has(k uint) bool {
	_, ok := v[k]
	return ok
}

/*
HasAny - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v Uint) HasAny(kk ...uint) bool {
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
func (v Uint) HasAnyOfSlice(kk []uint) bool {
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
func (v Uint) HasAnyOfSet(s Uint) bool {
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
func (v Uint) HasEach(kk ...uint) bool {
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
func (v Uint) HasEachOfSlice(kk []uint) bool {
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
func (v Uint) HasEachOfSet(s Uint) bool {
	for k, _ := range s {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

// Add - добавляет элементы в множество v
func (v Uint) Add(kk ...uint) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Uint) _AddTestHelper(kk ...uint) Uint {
	result := v.Copy()
	result.Add(kk...)
	return result
}

// AddSlice - добавляет элементы в множество v
func (v Uint) AddSlice(kk []uint) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Uint) _AddSliceTestHelper(kk []uint) Uint {
	result := v.Copy()
	result.AddSlice(kk)
	return result
}

// AddSet - добавляет элементы в множество v
func (v Uint) AddSet(s Uint) {
	for k, _ := range s {
		v[k] = struct{}{}
	}
}

func (v Uint) _AddSetTestHelper(s Uint) Uint {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Del - удаляет элементы из множествa v
func (v Uint) Del(kk ...uint) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Uint) _DelTestHelper(kk ...uint) Uint {
	result := v.Copy()
	result.Del(kk...)
	return result
}

// DelSlice - удаляет элементы из множествa v
func (v Uint) DelSlice(kk []uint) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Uint) _DelSliceTestHelper(kk []uint) Uint {
	result := v.Copy()
	result.DelSlice(kk)
	return result
}

// DelSet - удаляет элементы из множествa v
func (v Uint) DelSet(s Uint) {
	for k, _ := range s {
		delete(v, k)
	}
}

func (v Uint) _DelSetTestHelper(s Uint) Uint {
	result := v.Copy()
	result.DelSet(s)
	return result
}

// Union - возвращает результат объединения двух множеств. Исходные множества остаются без изменений
func (v Uint) Union(s Uint) Uint {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Intersect - возвращает результат пересечения двух множеств. Исходные множества остаются без изменений
func (v Uint) Intersect(s Uint) Uint {
	result := Uint{}
	for k, _ := range v {
		if _, ok := s[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// Subtract - возвращает результат вычитания двух множеств. Исходные множества остаются без изменений
func (v Uint) Subtract(s Uint) Uint {
	result := Uint{}
	for k, _ := range v {
		if _, ok := s[k]; !ok {
			result[k] = struct{}{}
		}
	}
	return result
}

type _uintSlice []uint

func (v _uintSlice) Len() int {
	return len(v)
}

func (v _uintSlice) Swap(i int, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v _uintSlice) Less(i int, j int) bool {
	return v[i] < v[j]
}
