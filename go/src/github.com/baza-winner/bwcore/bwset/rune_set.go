// Code generated by "bwsetter -type=rune"; DO NOT EDIT; bwsetter: go get -type=rune -set=Rune -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	"encoding/json"
	bwjson "github.com/baza-winner/bwcore/bwjson"
	"sort"
)

// Rune - множество значений типа rune с поддержкой интерфейсов Stringer и MarshalJSON
type Rune map[rune]struct{}

// RuneFrom - конструктор Rune
func RuneFrom(kk ...rune) Rune {
	result := Rune{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// RuneFromSlice - конструктор Rune
func RuneFromSlice(kk []rune) Rune {
	result := Rune{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// RuneFromSet - конструктор Rune
func RuneFromSet(s Rune) Rune {
	result := Rune{}
	for k, _ := range s {
		result[k] = struct{}{}
	}
	return result
}

// Copy - создает независимую копию
func (v Rune) Copy() Rune {
	return RuneFromSet(v)
}

// ToSlice - возвращает в виде []rune
func (v Rune) ToSlice() []rune {
	result := _runeSlice{}
	for k, _ := range v {
		result = append(result, k)
	}
	sort.Sort(result)
	return result
}

func _RuneToSliceTestHelper(kk []rune) []rune {
	return RuneFromSlice(kk).ToSlice()
}

// String - поддержка интерфейса Stringer
func (v Rune) String() string {
	return bwjson.Pretty(v)
}

// MarshalJSON - поддержка интерфейса MarshalJSON
func (v Rune) MarshalJSON() ([]byte, error) {
	result := []interface{}{}
	for _, k := range v.ToSlice() {
		result = append(result, k)
	}
	return json.Marshal(result)
}

// ToSliceOfStrings - возвращает []string строковых представлений элементов множества
func (v Rune) ToSliceOfStrings() []string {
	result := []string{}
	for k, _ := range v {
		result = append(result, string(k))
	}
	sort.Strings(result)
	return result
}

// Has - возвращает true, если множество содержит заданный элемент, в противном случае - false
func (v Rune) Has(k rune) bool {
	_, ok := v[k]
	return ok
}

/*
HasAny - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v Rune) HasAny(kk ...rune) bool {
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
func (v Rune) HasAnyOfSlice(kk []rune) bool {
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
func (v Rune) HasAnyOfSet(s Rune) bool {
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
func (v Rune) HasEach(kk ...rune) bool {
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
func (v Rune) HasEachOfSlice(kk []rune) bool {
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
func (v Rune) HasEachOfSet(s Rune) bool {
	for k, _ := range s {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

// Add - добавляет элементы в множество v
func (v Rune) Add(kk ...rune) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Rune) _AddTestHelper(kk ...rune) Rune {
	result := v.Copy()
	result.Add(kk...)
	return result
}

// AddSlice - добавляет элементы в множество v
func (v Rune) AddSlice(kk []rune) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Rune) _AddSliceTestHelper(kk []rune) Rune {
	result := v.Copy()
	result.AddSlice(kk)
	return result
}

// AddSet - добавляет элементы в множество v
func (v Rune) AddSet(s Rune) {
	for k, _ := range s {
		v[k] = struct{}{}
	}
}

func (v Rune) _AddSetTestHelper(s Rune) Rune {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Del - удаляет элементы из множествa v
func (v Rune) Del(kk ...rune) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Rune) _DelTestHelper(kk ...rune) Rune {
	result := v.Copy()
	result.Del(kk...)
	return result
}

// DelSlice - удаляет элементы из множествa v
func (v Rune) DelSlice(kk []rune) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Rune) _DelSliceTestHelper(kk []rune) Rune {
	result := v.Copy()
	result.DelSlice(kk)
	return result
}

// DelSet - удаляет элементы из множествa v
func (v Rune) DelSet(s Rune) {
	for k, _ := range s {
		delete(v, k)
	}
}

func (v Rune) _DelSetTestHelper(s Rune) Rune {
	result := v.Copy()
	result.DelSet(s)
	return result
}

// Union - возвращает результат объединения двух множеств. Исходные множества остаются без изменений
func (v Rune) Union(s Rune) Rune {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Intersect - возвращает результат пересечения двух множеств. Исходные множества остаются без изменений
func (v Rune) Intersect(s Rune) Rune {
	result := Rune{}
	for k, _ := range v {
		if _, ok := s[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// Subtract - возвращает результат вычитания двух множеств. Исходные множества остаются без изменений
func (v Rune) Subtract(s Rune) Rune {
	result := Rune{}
	for k, _ := range v {
		if _, ok := s[k]; !ok {
			result[k] = struct{}{}
		}
	}
	return result
}

type _runeSlice []rune

func (v _runeSlice) Len() int {
	return len(v)
}

func (v _runeSlice) Swap(i int, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v _runeSlice) Less(i int, j int) bool {
	return v[i] < v[j]
}
