// Предоставялет функции для работы со строками.
package bwstr

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
)

// ============================================================================

func SmartQuote(ss ...string) (result string) {
	result = ``
	for i, s := range ss {
		if i > 0 {
			result += ` `
		}
		if strings.ContainsAny(s, ` "`) {
			result += fmt.Sprintf(`%q`, s)
		} else {
			result += s
		}
	}
	return
}

type MultiLineMode uint8

const (
	Auto MultiLineMode = iota
	SingleLine
	MultiLine
)

type I interface {
	Strings() []string
}

type SS struct {
	SS        []string
	Preformat func(s string) string
}

func (v SS) Strings() (result []string) {
	if v.Preformat == nil {
		result = v.SS
	} else {
		for _, s := range v.SS {
			result = append(result, v.Preformat(s))
		}
	}
	return
}

type Vals struct {
	Vals      []interface{}
	Preformat func(val interface{}) string
}

func (v Vals) Strings() (result []string) {
	if v.Preformat == nil {
		v.Preformat = func(val interface{}) (result string) {
			if s, ok := val.(fmt.Stringer); ok {
				result = s.String()
			} else {
				if bytes, err := json.Marshal(val); err == nil {
					result = string(bytes)
				} else {
					result = err.Error()
				}
			}
			return
		}
	}
	for _, val := range v.Vals {
		result = append(result, v.Preformat(val))
	}
	return
}

type A struct {
	Source              I
	SingleJoiner        string
	MultiJoiner         string
	NoJoinerOnMutliline bool
	MultiPrefix         string
	MultiSuffix         string
	MaxLen              uint
	InitialIndent       string
	AdditionalIndent    string
	Mode                MultiLineMode
}

func SmartJoin(a A) (result string) {
	if a.Source == nil {
		return
	}
	source := a.Source.Strings()
	if len(source) == 0 {
		return
	}
	if len(a.MultiPrefix) == 0 {
		a.MultiPrefix = "one of ["
	}
	if len(a.MultiSuffix) == 0 {
		r := []rune(a.MultiPrefix[len(a.MultiPrefix)-1:])[0]
		if r2, ok := bw.Braces()[r]; ok {
			r = r2
		}
		a.MultiSuffix = string(r)
	}
	if len(a.MultiJoiner) == 0 {
		a.MultiJoiner = ", "
	}
	if len(a.SingleJoiner) == 0 {
		a.SingleJoiner = " or "
	}
	if len(a.AdditionalIndent) == 0 {
		a.AdditionalIndent = "  "
	}
	var (
		isMultiline bool
		joinerLen   int
		sumLen      int
		ss          []string
	)
	if a.Mode == MultiLine {
		isMultiline = true
	} else if a.Mode == Auto && a.MaxLen > 0 {
		joinerLen = ansi.Len(a.MultiJoiner)
		sumLen = ansi.Len(a.MultiPrefix) + ansi.Len(a.MultiSuffix) - joinerLen
		// bwdebug.Print("sumLen", sumLen, "a.MaxLen", a.MaxLen)
	}
	for _, s := range source {
		if !isMultiline && a.Mode == Auto && strings.IndexRune(s, '\n') >= 0 {
			isMultiline = true
		}
		if !isMultiline && a.Mode == Auto && a.MaxLen > 0 {
			sumLen += ansi.Len(s) + joinerLen
			// bwdebug.Print("sumLen", sumLen, "ansi.Len(s)", ansi.Len(s), "s", s, "a.MaxLen", a.MaxLen)
			if sumLen > int(a.MaxLen) {
				isMultiline = true
			}
		}
		ss = append(ss, s)
	}
	if len(ss) <= 2 {
		result = strings.Join(ss, a.SingleJoiner)
	} else {
		if isMultiline && !a.NoJoinerOnMutliline {
			a.MultiJoiner = strings.TrimRightFunc(a.MultiJoiner, func(r rune) bool { return unicode.IsSpace(r) })
		}
		result = a.MultiPrefix
		for i, s := range ss {
			if isMultiline {
				result += "\n"
				result += a.InitialIndent + a.AdditionalIndent
			} else if i > 0 {
				result += a.MultiJoiner
			}
			result += s
			if isMultiline && !a.NoJoinerOnMutliline {
				result += a.MultiJoiner
			}
		}
		if isMultiline {
			result += "\n" + a.InitialIndent
		}
		result += a.MultiSuffix
	}
	return
}

// ============================================================================
