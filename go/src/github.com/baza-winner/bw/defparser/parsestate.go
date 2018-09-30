package defparser

type parseState uint16

const (
	expectEOF parseState = iota
	expectEOFOrSpace

	expectValueOrSpace
	expectValueOrSpaceOrArrayItemSeparator
	expectValueOrSpaceOrMapKeySeparator

	expectArrayItemSeparatorOrSpace

	expectMapKeySeparatorOrSpace

	expectRocket

	expectMapKey

	expectWord

	expectDigit
	expectDigitOrUnderscoreOrDot
	expectDigitOrUnderscore

	expectDoubleQuotedStringContent
	expectSingleQuotedStringContent
	expectDoubleQuotedStringEscapedContent
	expectSingleQuotedStringEscapedContent

	expectContentOfDoubleQuotedKey
	expectSingleQuotedKeyContent
	expectDoubleQuotedKeyEscapedContent
	expectSingleQuotedKeyEscapedContent

	expectSpaceOrMapKey
	expectSpaceOrMapKeyOrMapValueSeparator

	expectSpaceOrQwItemOrDelimiter

	expectEndOfQwItem
)

//go:generate stringer -type=parseState
