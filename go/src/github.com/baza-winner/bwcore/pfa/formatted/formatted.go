package formatted

import (
	"fmt"

	"github.com/baza-winner/bwcore/ansi"
)

type String string

func StringFrom(fmtString string, fmtArgs ...interface{}) String {
	return String(ansi.String(fmt.Sprintf(fmtString, fmtArgs...)))
}

func (v String) Concat(s String) String {
	return String(ansi.Concat(string(v), string(s)))
	// return String(string(v) + string(s))
}

type FormattedString interface {
	FormattedString() String
}