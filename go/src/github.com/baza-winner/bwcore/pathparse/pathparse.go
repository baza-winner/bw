package pathparse

import "github.com/baza-winner/bwcore/bwerr"

type stringRuneProvider struct {
	pos int
	src []rune
}

func (v *stringRuneProvider) PullRune() *rune {
	v.pos += 1
	if v.pos >= len(v.src) {
		return nil
	} else {
		result := v.src[v.pos]
		return &result
	}
}

// Parse - парсит строку
func Parse(source string) ([]interface{}, error) {
	if len(source) == 0 {
		return []interface{}{}, nil
	} else {
		return pfaParse(&stringRuneProvider{pos: -1, src: []rune(source)})
	}
}

// MustParse is like Parse but panics if the expression cannot be parsed.
// It simplifies safe initialization of global variables holding parsed values.
func MustParse(source string) (result []interface{}) {
	var err error
	if result, err = Parse(source); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return result
}
