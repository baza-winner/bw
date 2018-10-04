package bwset

type Strings map[string]struct{}

func getStrings(ss []string) Strings {
  result := Strings{}
  for _, s := range ss {
    result[s] = struct{}{}
  }
  return result
}

func (ss Strings) has(s string) (ok bool) {
  _, ok = ss[s]
  return
}
