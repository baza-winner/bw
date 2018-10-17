package defparse

func prepareStateDef() *stateDef {
	var unexpectedEof, unexpectedRune []interface{}
	unexpectedEof = []interface{}{eofRune{}, setError{unexpectedRuneError}}
	unexpectedRune = []interface{}{setError{unexpectedRuneError}}

	finishItemStateDef := createStateDef(
		[]interface{}{topItem{parseStackItemString}, topItem{parseStackItemQwItem},
			setTopItemValueAsString{},
		},
		[]interface{}{topItem{parseStackItemMap},
			setTopItemValueAsMap{},
		},
		[]interface{}{topItem{parseStackItemArray}, topItem{parseStackItemQw},
			setTopItemValueAsArray{},
		},
		[]interface{}{topItem{parseStackItemNumber},
			setTopItemValueAsNumber{},
		},
		[]interface{}{topItem{parseStackItemWord},
			processStateDef{createStateDef(
				[]interface{}{"true",
					setTopItemValueAsBool{true},
				},
				[]interface{}{"false",
					setTopItemValueAsBool{false},
				},
				[]interface{}{"nil", "null"},
				[]interface{}{"Bool", "String", "Int", "Number", "Map", "Array", "ArrayOf",
					setTopItemValueAsString{},
				},
				[]interface{}{"qw",
					pullRune{},
					processStateDef{createStateDef(
						[]interface{}{unicodeOpenBraces, unicodePunct, unicodeSymbol,
							setPrimary{expectSpaceOrQwItemOrDelimiter},
							setTopItemDelimiter{pairForCurrRune{}},
							setTopItemType{parseStackItemQw},
						},
						[]interface{}{
							setError{unexpectedRuneError},
						},
					)},
					setVar{"skipPostProcess", true},
				},
				[]interface{}{
					setError{unknownWordError},
				},
			)},
		},
	)

	postProcessStateDef := createStateDef(
		[]interface{}{stackLenIs{0}},
		[]interface{}{stackLenIs{1},
			setSecondary{expectEOF, orSpace},
		},
		[]interface{}{
			popSubItem{},
			processStateDef{createStateDef(
				[]interface{}{topItem{parseStackItemQw},
					appendItemArray{fromSubItemValue{}},
					setPrimary{expectSpaceOrQwItemOrDelimiter},
				},
				[]interface{}{topItem{parseStackItemArray},
					processStateDef{createStateDef(
						[]interface{}{subItem{parseStackItemQw},
							appendItemArray{fromSubItemArray{}},
						},
						[]interface{}{
							appendItemArray{fromSubItemValue{}},
						},
					)},
					setSecondary{expectValueOrSpace, orArrayItemSeparator},
				},
				[]interface{}{topItem{parseStackItemMap},
					processStateDef{createStateDef(
						[]interface{}{subItem{parseStackItemKey},
							setTopItemStringFromSubItem{},
							setSecondary{expectValueOrSpace, orMapKeySeparator},
						},
						[]interface{}{
							setTopItemMapKeyValueFromSubItem{},
							setSecondary{expectSpaceOrMapKey, orMapValueSeparator},
						},
					)},
				},
				[]interface{}{
					unreachable{},
				},
			)},
		},
	)

	primaryStateDef := createStateDef(
		[]interface{}{expectEOF,
			processStateDef{createStateDef(
				[]interface{}{eofRune{}, setPrimary{expectEOF}},
				[]interface{}{unicodeSpace},
				unexpectedRune,
			)},
		},
		[]interface{}{expectRocket,
			processStateDef{createStateDef(
				[]interface{}{'>', setPrimary{expectValueOrSpace}},
				unexpectedRune,
			)},
		},
		[]interface{}{expectWord,
			processStateDef{createStateDef(
				[]interface{}{unicodeLetter, unicodeDigit,
					appendCurrRune{},
				},
				[]interface{}{
					pushRune{},
					setVar{"needFinish", true},
				},
			)},
		},
		[]interface{}{expectSpaceOrQwItemOrDelimiter,
			processStateDef{createStateDef(
				unexpectedEof,
				[]interface{}{unicodeSpace},
				[]interface{}{delimiterRune{},
					setVar{"needFinish", true},
				},
				[]interface{}{
					pushItem{itemType: parseStackItemQwItem, itemString: fromCurrRune{}, delimiter: fromParentItem{}},
					setPrimary{expectEndOfQwItem},
				},
			)},
		},
		[]interface{}{expectEndOfQwItem,
			processStateDef{createStateDef(
				unexpectedEof,
				[]interface{}{unicodeSpace, delimiterRune{},
					pushRune{},
					setVar{"needFinish", true},
				},
				[]interface{}{
					appendCurrRune{},
				},
			)},
		},
		[]interface{}{expectContentOf,
			processStateDef{createStateDef(
				unexpectedEof,
				[]interface{}{delimiterRune{},
					setVar{"needFinish", true},
				},
				[]interface{}{'\\',
					changePrimary{expectEscapedContentOf},
				},
				[]interface{}{
					appendCurrRune{},
				},
			)},
		},
		[]interface{}{expectDigit,
			processStateDef{createStateDef(
				[]interface{}{unicodeDigit, noSecondaryState,
					appendCurrRune{},
					changeSecondary{orUnderscoreOrDot},
				},
				[]interface{}{'.', orUnderscoreOrDot,
					appendCurrRune{},
					changeSecondary{orUnderscore},
				},
				[]interface{}{'_', unicodeDigit, orUnderscoreOrDot, orUnderscore,
					appendCurrRune{},
				},
				[]interface{}{noSecondaryState,
					setError{unexpectedRuneError},
				},
				[]interface{}{
					pushRune{},
					setVar{"needFinish", true},
				},
			)},
		},
		[]interface{}{expectSpaceOrMapKey,
			processStateDef{createStateDef(
				[]interface{}{unicodeSpace},
				[]interface{}{unicodeLetter,
					pushItem{itemType: parseStackItemKey, itemString: fromCurrRune{}},
					setPrimary{expectWord},
				},
				[]interface{}{'"', '\'',
					pushItem{itemType: parseStackItemKey, delimiter: fromCurrRune{}},
					setSecondary{expectContentOf, keyToken},
				},
				[]interface{}{',', orMapValueSeparator,
					setPrimary{expectSpaceOrMapKey},
				},
				[]interface{}{delimiterRune{}, topItem{parseStackItemMap},
					setVar{"needFinish", true},
				},
				unexpectedRune,
			)},
		},
		[]interface{}{expectEscapedContentOf,
			processStateDef{createStateDef(
				[]interface{}{'"', '\'', '\\',
					appendCurrRune{},
					changePrimary{expectContentOf},
				},
				[]interface{}{delimiterIs{'"'},
					processStateDef{createStateDef(
						[]interface{}{'a', appendRune{'\a'}},
						[]interface{}{'b', appendRune{'\b'}},
						[]interface{}{'f', appendRune{'\f'}},
						[]interface{}{'n', appendRune{'\n'}},
						[]interface{}{'r', appendRune{'\r'}},
						[]interface{}{'t', appendRune{'\t'}},
						[]interface{}{'v', appendRune{'\v'}},
						unexpectedRune,
					)},
					changePrimary{expectContentOf},
				},
				unexpectedRune,
			)},
		},
		[]interface{}{expectValueOrSpace,
			processStateDef{createStateDef(
				[]interface{}{eofRune{}, stackLenIs{0},
					setPrimary{expectEOF},
				},
				unexpectedEof,
				[]interface{}{'=', orMapKeySeparator,
					setPrimary{expectRocket},
				},
				[]interface{}{':', orMapKeySeparator,
					setPrimary{expectValueOrSpace},
				},
				[]interface{}{',', orArrayItemSeparator,
					setPrimary{expectValueOrSpace},
				},
				[]interface{}{unicodeSpace},
				[]interface{}{'{',
					pushItem{itemType: parseStackItemMap, delimiter: pairForCurrRune{}},
					setPrimary{expectSpaceOrMapKey},
				},
				[]interface{}{'<',
					pushItem{itemType: parseStackItemQw, delimiter: pairForCurrRune{}},
					setPrimary{expectSpaceOrQwItemOrDelimiter},
				},
				[]interface{}{'[',
					pushItem{itemType: parseStackItemArray, delimiter: pairForCurrRune{}},
					setPrimary{expectValueOrSpace},
				},
				[]interface{}{topItem{parseStackItemArray}, delimiterRune{},
					setVar{"needFinish", true},
				},
				[]interface{}{'-', '+',
					pushItem{itemType: parseStackItemNumber, itemString: fromCurrRune{}},
					setPrimary{expectDigit},
				},
				[]interface{}{unicodeDigit,
					pushItem{itemType: parseStackItemNumber, itemString: fromCurrRune{}},
					setSecondary{expectDigit, orUnderscoreOrDot},
				},
				[]interface{}{'"', '\'',
					pushItem{itemType: parseStackItemString, delimiter: fromCurrRune{}},
					setSecondary{expectContentOf, stringToken},
				},
				[]interface{}{unicodeLetter,
					pushItem{itemType: parseStackItemWord, itemString: fromCurrRune{}},
					setPrimary{expectWord},
				},
				unexpectedRune,
			)},
		},
		[]interface{}{expectValueOrSpace,
			unreachable{},
		},
	)

	result := createStateDef(
		[]interface{}{
			pullRune{},
			setVar{"needFinish", false},
			processStateDef{primaryStateDef},
			processStateDef{createStateDef(
				[]interface{}{varIs{"needFinish", true},
					setVar{"skipPostProcess", false},
					processStateDef{finishItemStateDef},
					processStateDef{createStateDef(
						[]interface{}{varIs{"skipPostProcess", false},
							processStateDef{postProcessStateDef},
						},
					)},
				},
			)},
		},
	)
	return result
}

const (
	parseStackItem_below_ parseStackItemType = iota
	parseStackItemKey
	parseStackItemString
	parseStackItemMap
	parseStackItemArray
	parseStackItemQw
	parseStackItemQwItem
	parseStackItemNumber
	parseStackItemWord
	parseStackItem_above_
)

const (
	parsePrimaryState_below_ parsePrimaryState = iota
	expectEOF
	expectValueOrSpace
	expectRocket
	expectWord
	expectDigit
	expectContentOf
	expectEscapedContentOf
	expectSpaceOrMapKey
	expectSpaceOrQwItemOrDelimiter
	expectEndOfQwItem
	parsePrimaryState_above_
)

const (
	anySecondaryState parseSecondaryState = iota
	noSecondaryState
	orSpace

	orMapKeySeparator
	orArrayItemSeparator

	orUnderscoreOrDot
	orUnderscore

	stringToken
	keyToken

	orMapValueSeparator
)
