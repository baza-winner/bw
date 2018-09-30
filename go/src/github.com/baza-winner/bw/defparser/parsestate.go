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

	expectContentOfDoubleQuotedString
	expectContentOfSingleQuotedString
	expectEscapedContentOfDoubleQuotedString
	expectSingleQuotedStringEscapedContent

	expectContentOfDoubleQuotedKey
	expectContentOfSingleQuotedKey
	expectEscapedContentOfDoubleQuotedKey
	expectEscapedContentOfSingleQuotedKey

	expectSpaceOrMapKey
	expectSpaceOrMapKeyOrMapValueSeparator

	expectSpaceOrQwItemOrDelimiter

	expectEndOfQwItem
)

//go:generate stringer -type=parseState
