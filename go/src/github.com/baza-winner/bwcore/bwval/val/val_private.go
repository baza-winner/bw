package val

import "github.com/baza-winner/bwcore/runeprovider"

//go:generate stringer -type=PrimaryState,SecondaryState,ItemKind

var trace func(r rune, primary PrimaryState, secondary SecondaryState, stack []StackItem)

type PrimaryState uint8

const (
	Begin PrimaryState = iota
	ExpectSpaceOrQwItemOrDelimiter
	ExpectSpaceOrMapKey
	ExpectEndOfQwItem
	ExpectContentOf
	ExpectWord
	ExpectEscapedContentOf
	ExpectRocket
	ExpectDigit
	End
)

type SecondaryState uint8

const (
	None SecondaryState = iota
	orArrayItemSeparator
	orMapKeySeparator
	orMapValueSeparator
	orUnderscoreOrDot
	orUnderscore
)

type ItemKind uint8

const (
	ItemString ItemKind = iota
	ItemQw
	ItemQwItem
	ItemNumber
	ItemWord
	ItemKey
	ItemMap
	ItemArray
)

type StackItem struct {
	PosStruct runeprovider.PosStruct
	Kind      ItemKind
	S         string
	Result    interface{}
	Delimiter rune
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
var Braces = map[rune]rune{
	'(': ')',
	'{': '}',
	'<': '>',
	'[': ']',
}
