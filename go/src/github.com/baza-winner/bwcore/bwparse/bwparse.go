package bwparse

import "github.com/baza-winner/bwcore/runeprovider"

func ParseString(p *runeprovider.Proxy) (result string, err error) {
	type State uint8
	const (
		Begin State = iota
		ExpectContentOf
		ExpectEscapedContentOf
		End
	)

	var (
		delimiter rune
		r2        rune
		ok        bool
	)
	r, isEOF, _ := p.Rune()
	state := Begin
	for {
		switch state {
		case Begin:
			switch r {
			case '"', '\'':
				delimiter = r
				state = ExpectContentOf
			}
		case ExpectContentOf:
			switch {
			case r == delimiter:
				state = End
			case r == '\\':
				state = ExpectEscapedContentOf
			case !isEOF:
				result += string(r)
			}
		case ExpectEscapedContentOf:
			switch {
			case r == '"' || r == '\'' || r == '\\':
				result += string(r)
			case delimiter == '"':
				if r2, ok = EscapeRunes[r]; ok {
					result += string(r2)
				} else {
					err = p.Unexpected(p.Curr)
					return
				}
			}
			state = ExpectContentOf
		}
		if state == End {
			break
		} else if p.Curr.IsEOF {
			err = p.Unexpected(p.Curr)
			return
		} else {
			_ = p.PullRune()
			r, isEOF, _ = p.Rune()
		}
	}
	if state != End && p.Curr.IsEOF {
		err = p.Unexpected(p.Curr)
		return
	}
	return
}

var EscapeRunes = map[rune]rune{
	'a': '\a',
	'b': '\b',
	'f': '\f',
	'n': '\n',
	'r': '\r',
	't': '\t',
	'v': '\v',
}
