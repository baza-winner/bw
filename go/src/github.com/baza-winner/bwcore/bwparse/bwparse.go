package bwparse

import (
	"unicode"

	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwstr"
	"github.com/baza-winner/bwcore/runeprovider"
)

// ============================================================================

func ParseWord(p *runeprovider.Proxy) (result string, err error) {

	type State uint8
	const (
		begin State = iota
		expectWord
		end
	)

	r, _, _ := p.Rune()
	state := begin
	for {
		switch state {
		case begin:
			if unicode.IsLetter(r) || r == '_' {
				result = string(r)
				state = expectWord
			} else {
				err = p.Unexpected(p.Curr)
				return
			}
		case expectWord:
			if unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r) {
				result += string(r)
			} else {
				_ = p.PushRune()
				state = end
			}
		}
		if state == end {
			break
		} else if p.Curr.IsEOF {
			err = p.Unexpected(p.Curr)
			return
		} else {
			_ = p.PullRune()
			r, _, _ = p.Rune()
		}
	}
	return
}

// ============================================================================

func ParseString(p *runeprovider.Proxy) (result string, err error) {
	type State uint8
	const (
		begin State = iota
		expectContent
		expectEscapedContent
		end
	)

	var (
		delimiter rune
		r2        rune
		ok        bool
	)
	r, isEOF, _ := p.Rune()
	state := begin
	for {
		switch state {
		case begin:
			switch r {
			case '"', '\'':
				delimiter = r
				state = expectContent
			default:
				err = p.Unexpected(p.Curr)
				return
			}
		case expectContent:
			switch {
			case r == delimiter:
				state = end
			case r == '\\':
				state = expectEscapedContent
			case !isEOF:
				result += string(r)
			}
		case expectEscapedContent:
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
			state = expectContent
		}
		if state == end {
			break
		} else if p.Curr.IsEOF {
			err = p.Unexpected(p.Curr)
			return
		} else {
			_ = p.PullRune()
			r, isEOF, _ = p.Rune()
		}
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

// ============================================================================

func ParseNumber(p *runeprovider.Proxy) (result interface{}, err error) {
	type State uint8
	const (
		begin State = iota
		expectDigitOnly
		expectDigitOrUnderscore
		expectDigitOrUnderscoreOrDot
		end
	)
	var (
		ps runeprovider.PosStruct
		s  string
	)
	r, _, _ := p.Rune()
	state := begin
	for {
		switch state {
		case begin:
			switch r {
			case '-', '+':
				ps = p.Curr
				s = string(r)
				state = expectDigitOnly
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				ps = p.Curr
				s = string(r)
				state = expectDigitOrUnderscoreOrDot
			default:
				err = p.Unexpected(p.Curr)
				return
			}
		case expectDigitOnly:
			switch r {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				s += string(r)
				state = expectDigitOrUnderscoreOrDot
			default:
				err = p.Unexpected(p.Curr)
				return
			}
		case expectDigitOrUnderscoreOrDot:
			switch r {
			case '.':
				s += string(r)
				state = expectDigitOrUnderscore
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '_':
				s += string(r)
			default:
				p.PushRune()
				state = end
			}
		case expectDigitOrUnderscore:
			switch r {
			case '.':
				s += string(r)
				state = expectDigitOrUnderscore
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '_':
				s += string(r)
			default:
				p.PushRune()
				state = end
			}
		}

		if state == end {
			break
		} else if p.Curr.IsEOF {
			err = p.Unexpected(p.Curr)
			return
		} else {
			_ = p.PullRune()
			r, _, _ = p.Rune()
		}
	}
	if result, err = bwstr.ParseNumber(s); err != nil {
		err = p.Unexpected(ps, bwerr.Err(err))
		return
	}
	return
}
