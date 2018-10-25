package formatted

import (
	"fmt"

	"github.com/baza-winner/bwcore/ansi"
)

type String string

func StringFrom(fmtString string, fmtArgs ...interface{}) String {
	return String(fmt.Sprintf(ansi.Ansi("", fmtString), fmtArgs...))
}

func (v String) Concat(s String) String {
	return String(string(v) + string(s))
}

type FormattedString interface {
	FormattedString() String
}
