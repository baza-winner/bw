package bwset

import  (
  "sort"
  "github.com/baza-winner/bwcore/bwjson"
)

type Strings map[string]struct{}


func StringsFromSlice(kk []string) Strings {
  result := Strings{}
  result.Add(kk...)
  return result
}

func StringsFromArgs(kk ...string) Strings {
  return StringsFromSlice(kk)
}

func (v Strings) Has(k string) (ok bool) {
  _, ok = v[k]
  return
}

func (v Strings) Add(kk ...string) {
  for _, s := range kk {
    v[s] = struct{}{}
  }
}

func (v Strings) ToSlice() []string {
  result := []string{}
  for k, _ := range v {
    result = append(result, k)
  }
  sort.Strings(result)
  return result
}

func (v Strings) GetDataForJson() interface{} {
  result := []interface{}{}
  for _, s := range v.ToSlice() {
    result = append(result, s)
  }
  return result
}

func (v Strings) String() string {
  return bwjson.PrettyJsonOf(v)
}