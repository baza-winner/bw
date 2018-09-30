package defparser

import (
	"fmt"
)

type parsePrimaryState uint16

const (
	_expectBelow parsePrimaryState = iota
	expectEOF
	expectValueOrSpace
	expectArrayItemSeparatorOrSpace
	expectMapKeySeparatorOrSpace
	expectRocket
	expectMapKey
	expectWord
	expectDigit
	expectContentOf
	expectEscapedContentOf
	expectSpaceOrMapKey
	expectSpaceOrQwItemOrDelimiter
	expectEndOfQwItem
	_expectAbove
)

//go:generate stringer -type=parsePrimaryState

type parseSecondaryState uint16

const (
	noSecondaryState parseSecondaryState = iota
	orSpace

	orMapKeySeparator
	orArrayItemSeparator

	orUnderscoreOrDot
	orUnderscore

	doubleQuoted
	singleQuoted

	orMapValueSeparator
)

//go:generate stringer -type=parseSecondaryState

type parseTertiaryState uint16

const (
	noTertiaryState parseTertiaryState = iota
	stringToken
	keyToken
)

//go:generate stringer -type=parseTertiaryState

type parseState struct {
	primary   parsePrimaryState
	secondary parseSecondaryState
	tertiary  parseTertiaryState
}

func (state *parseState) setPrimary(primary parsePrimaryState) {
	state.setSecondary(primary, noSecondaryState)
}

func (state *parseState) setSecondary(primary parsePrimaryState, secondary parseSecondaryState) {
	state.setTertiary(primary, secondary, noTertiaryState)
}

func (state *parseState) setTertiary(primary parsePrimaryState, secondary parseSecondaryState, tertiary parseTertiaryState) {
	state.primary = primary
	state.secondary = secondary
	state.tertiary = tertiary
}

func (state parseState) String() string {
	if state.tertiary != noTertiaryState {
		return fmt.Sprintf(`%s.%s.%s`, state.primary, state.secondary, state.tertiary)
	} else if state.secondary != noSecondaryState {
		return fmt.Sprintf(`%s.%s`, state.primary, state.secondary)
	} else {
		return state.primary.String()
	}
}
