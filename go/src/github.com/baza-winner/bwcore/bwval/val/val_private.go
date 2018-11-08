package val

import "github.com/baza-winner/bwcore/runeprovider"

//go:generate stringer -type=primaryState,secondaryState,itemType

var trace func(r rune, primary primaryState, secondary secondaryState, stack []stackItem)

type primaryState uint8

const (
	begin primaryState = iota
	expectSpaceOrQwItemOrDelimiter
	expectSpaceOrMapKey
	expectEndOfQwItem
	expectContentOf
	expectWord
	expectEscapedContentOf
	expectRocket
	expectDigit
	end
)

type secondaryState uint8

const (
	none secondaryState = iota
	orArrayItemSeparator
	orMapKeySeparator
	orMapValueSeparator
	orUnderscoreOrDot
	orUnderscore
)

type itemType uint8

const (
	itemString itemType = iota
	itemQw
	itemQwItem
	itemNumber
	itemWord
	itemKey
	itemMap
	itemArray
)

type stackItem struct {
	ps        runeprovider.PosStruct
	it        itemType
	s         string
	result    interface{}
	delimiter rune
}

var escapeRunes = map[rune]rune{
	'a': '\a',
	'b': '\b',
	'f': '\f',
	'n': '\n',
	'r': '\r',
	't': '\t',
	'v': '\v',
}
var braces = map[rune]rune{
	'(': ')',
	'{': '}',
	'<': '>',
	'[': ']',
}
