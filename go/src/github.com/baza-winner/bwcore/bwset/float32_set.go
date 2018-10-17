// Code generated by "bwsetter -type=float32"; DO NOT EDIT; bwsetter: go get -type=float32 -set=Float32Set -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	bwjson "github.com/baza-winner/bwcore/bwjson"
	"sort"
	"strconv"
)

// Float32Set - множество значений типа float32 с поддержкой интерфейсов Stringer и github.com/baza-winner/bwcore/bwjson.Jsonable
type Float32Set map[float32]struct{}

// Float32SetFrom - конструктор Float32Set
func Float32SetFrom(kk ...float32) Float32Set {
	result := Float32Set{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Float32SetFromSlice - конструктор Float32Set
func Float32SetFromSlice(kk []float32) Float32Set {
	result := Float32Set{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Float32SetFromSet - конструктор Float32Set
func Float32SetFromSet(s Float32Set) Float32Set {
	result := Float32Set{}
	for k, _ := range s {
		result[k] = struct{}{}
	}
	return result
}

// Copy - создает независимую копию
func (v Float32Set) Copy() Float32Set {
	return Float32SetFromSet(v)
}

// ToSlice - возвращает в виде []float32
func (v Float32Set) ToSlice() []float32 {
	result := _float32Slice{}
	for k, _ := range v {
		result = append(result, k)
	}
	sort.Sort(result)
	return result
}

func _Float32SetToSliceTestHelper(kk []float32) []float32 {
	return Float32SetFromSlice(kk).ToSlice()
}

// String - поддержка интерфейса Stringer
func (v Float32Set) String() string {
	return bwjson.PrettyJsonOf(v)
}

// DataForJson - поддержка интерфейса bwjson.Jsonable
func (v Float32Set) DataForJson() interface{} {
	result := []interface{}{}
	for k, _ := range v {
		result = append(result, k)
	}
	return result
}

// ToSliceOfStrings - возвращает []string строковых представлений элементов множества
func (v Float32Set) ToSliceOfStrings() []string {
	result := []string{}
	for k, _ := range v {
		result = append(result, strconv.FormatFloat(float64(k), byte(0x66), -1, 64))
	}
	sort.Strings(result)
	return result
}

// Has - возвращает true, если множество содержит заданный элемент, в противном случае - false
func (v Float32Set) Has(k float32) bool {
	_, ok := v[k]
	return ok
}

/*
HasAny - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v Float32Set) HasAny(kk ...float32) bool {
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
func (v Float32Set) HasAnyOfSlice(kk []float32) bool {
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
func (v Float32Set) HasAnyOfSet(s Float32Set) bool {
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
func (v Float32Set) HasEach(kk ...float32) bool {
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
func (v Float32Set) HasEachOfSlice(kk []float32) bool {
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
func (v Float32Set) HasEachOfSet(s Float32Set) bool {
	for k, _ := range s {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

// Add - добавляет элементы в множество v
func (v Float32Set) Add(kk ...float32) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Float32Set) _AddTestHelper(kk ...float32) Float32Set {
	result := v.Copy()
	result.Add(kk...)
	return result
}

// AddSlice - добавляет элементы в множество v
func (v Float32Set) AddSlice(kk []float32) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Float32Set) _AddSliceTestHelper(kk []float32) Float32Set {
	result := v.Copy()
	result.AddSlice(kk)
	return result
}

// AddSet - добавляет элементы в множество v
func (v Float32Set) AddSet(s Float32Set) {
	for k, _ := range s {
		v[k] = struct{}{}
	}
}

func (v Float32Set) _AddSetTestHelper(s Float32Set) Float32Set {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Del - удаляет элементы из множествa v
func (v Float32Set) Del(kk ...float32) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Float32Set) _DelTestHelper(kk ...float32) Float32Set {
	result := v.Copy()
	result.Del(kk...)
	return result
}

// DelSlice - удаляет элементы из множествa v
func (v Float32Set) DelSlice(kk []float32) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Float32Set) _DelSliceTestHelper(kk []float32) Float32Set {
	result := v.Copy()
	result.DelSlice(kk)
	return result
}

// DelSet - удаляет элементы из множествa v
func (v Float32Set) DelSet(s Float32Set) {
	for k, _ := range s {
		delete(v, k)
	}
}

func (v Float32Set) _DelSetTestHelper(s Float32Set) Float32Set {
	result := v.Copy()
	result.DelSet(s)
	return result
}

// Union - возвращает результат объединения двух множеств. Исходные множества остаются без изменений
func (v Float32Set) Union(s Float32Set) Float32Set {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Intersect - возвращает результат пересечения двух множеств. Исходные множества остаются без изменений
func (v Float32Set) Intersect(s Float32Set) Float32Set {
	result := Float32Set{}
	for k, _ := range v {
		if _, ok := s[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// Subtract - возвращает результат вычитания двух множеств. Исходные множества остаются без изменений
func (v Float32Set) Subtract(s Float32Set) Float32Set {
	result := Float32Set{}
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
