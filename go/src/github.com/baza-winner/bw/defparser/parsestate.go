package defparser

type parseState uint16

const (
	expectSpaceOrValue parseState = iota
	expectDigit
	expectDigitOrUnderscoreOrDot
	expectDigitOrUnderscore
	expectDoubleQuotedStringContent
	expectSingleQuotedStringContent
	expectDoubleQuotedStringEscapedContent
	expectSingleQuotedStringEscapedContent
	expectWord
	expectKey
	unexpectedChar
	expectArrayItemSeparatorOrSpaceOrArrayValue
	expectMapKeySeparatorOrSpace
	expectMapKeySeparatorOrSpaceOrMapValue
	expectRocket
	expectMapValueSeparatorOrSpaceOrMapValue
	expectSpaceOrMapKey
	expectMapValueSeparatorOrSpaceOrMapKey
	expectDoubleQuotedKeyContent
	expectSingleQuotedKeyContent
	expectDoubleQuotedKeyEscapedContent
	expectSingleQuotedKeyEscapedContent
	expectSpaceOrQwItemOrDelimiter
	expectEndOfQwItem
	expectArrayItemSeparatorOrSpace
	expectEOF
	expectSpaceOrEOF
)

//go:generate stringer -type=parseState
