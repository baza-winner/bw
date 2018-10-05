package bwset

import  (
  "sort"
  "github.com/baza-winner/bwcore/bwjson"
)
type Strings map[string]struct{}

func FromSliceOfStrings(ss []string) Strings {
  result := Strings{}
  result.Add(ss...)
  return result
}

func (v Strings) Has(s string) (ok bool) {
  _, ok = v[s]
  return
}

func (v Strings) Add(ss ...string) {
  for _, s := range ss {
    v[s] = struct{}{}
  }
}

func (v Strings) ToSliceOfStrings() []string {
  result := []string{}
  for s, _ := range v {
    result = append(result, s)
  }
  sort.Strings(result)
  return result
}

func (v Strings) GetDataForJson() interface{} {
  result := []interface{}{}
  ss := v.ToSliceOfStrings()
  for _, s := range ss {
    result = append(result, s)
  }
  return result
}

func (v Strings) String() string {
  return bwjson.PrettyJsonOf(v)
}