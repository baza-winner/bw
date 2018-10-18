package defparse

import (
	"github.com/baza-winner/bwcore/bwset"
	. "github.com/baza-winner/bwcore/pfa"
)

var pfaStateDef *LogicDef

func init() {
	pfaStateDef = prepareLogicDef()
}

func prepareLogicDef() *LogicDef {

	unexpectedEOF := []interface{}{IsEOF{}, SetError{UnexpectedRune}}

	unexpectedRune := []interface{}{SetError{UnexpectedRune}}

	finishLogic := CreateLogicDef(
		[]interface{}{TopItemIs{"string"}, TopItemIs{"qwItem"},
			SetTopItemValueAsString{},
		},
		[]interface{}{TopItemIs{"map"},
			SetTopItemValueAsMap{},
		},
		[]interface{}{TopItemIs{"array"}, TopItemIs{"qw"},
			SetTopItemValueAsArray{},
		},
		[]interface{}{TopItemIs{"number"},
			SetTopItemValueAsNumber{},
		},
		[]interface{}{TopItemIs{"word"},
			SubLogic{CreateLogicDef(
				[]interface{}{TopItemStringIs{"true"},
					SetTopItemValueAsBool{true},
				},
				[]interface{}{TopItemStringIs{"false"},
					SetTopItemValueAsBool{false},
				},
				[]interface{}{TopItemStringIsOneOf{bwset.StringSetFrom("nil", "null")}},
				[]interface{}{TopItemStringIsOneOf{bwset.StringSetFrom("Bool", "String", "Int", "Number", "Map", "Array", "ArrayOf")},
					SetTopItemValueAsString{},
				},
				[]interface{}{TopItemStringIs{"qw"},
					PullRune{},
					SubLogic{CreateLogicDef(
						[]interface{}{IsUnicodeOpenBraces, IsUnicodePunct, IsUnicodeSymbol,
							SetPrimary{"expectSpaceOrQwItemOrDelimiter"},
							SetTopItemDelimiter{PairForCurrRune{}},
							SetTopItemType{"qw"},
						},
						unexpectedRune,
					)},
					SetVar{"skipPostProcess", true},
				},
				[]interface{}{
					SetError{UnknownWord},
				},
			)},
		},
	)

	postProcessLogic := CreateLogicDef(
		[]interface{}{StackLenIs{0}},
		[]interface{}{StackLenIs{1},
			SetSecondary{"expectEOF", "orSpace"},
		},
		[]interface{}{
			PopSubItem{},
			SubLogic{CreateLogicDef(
				[]interface{}{TopItemIs{"qw"},
					AppendItemArray{FromSubItemValue{}},
					SetPrimary{"expectSpaceOrQwItemOrDelimiter"},
				},
				[]interface{}{TopItemIs{"array"},
					SubLogic{CreateLogicDef(
						[]interface{}{SubItemIs{"qw"},
							AppendItemArray{FromSubItemArray{}},
						},
						[]interface{}{
							AppendItemArray{FromSubItemValue{}},
						},
					)},
					SetSecondary{"expectValueOrSpace", "orArrayItemSeparator"},
				},
				[]interface{}{TopItemIs{"map"},
					SubLogic{CreateLogicDef(
						[]interface{}{SubItemIs{"key"},
							SetTopItemStringFromSubItem{},
							SetSecondary{"expectValueOrSpace", "orMapKeySeparator"},
						},
						[]interface{}{
							SetTopItemMapKeyValueFromSubItem{},
							SetSecondary{"expectSpaceOrMapKey", "orMapValueSeparator"},
						},
					)},
				},
				[]interface{}{
					Unreachable{},
				},
			)},
		},
	)

	primaryStateLogic := CreateLogicDef(
		[]interface{}{PrimaryIs{"expectEOF"},
			SubLogic{CreateLogicDef(
				[]interface{}{IsEOF{}, SetPrimary{"expectEOF"}},
				[]interface{}{IsUnicodeSpace},
				unexpectedRune,
			)},
		},
		[]interface{}{PrimaryIs{"expectRocket"},
			SubLogic{CreateLogicDef(
				[]interface{}{'>', SetPrimary{"expectValueOrSpace"}},
				unexpectedRune,
			)},
		},
		[]interface{}{PrimaryIs{"expectWord"},
			SubLogic{CreateLogicDef(
				[]interface{}{IsUnicodeLetter, IsUnicodeDigit,
					AppendCurrRune{},
				},
				[]interface{}{
					PushRune{},
					SetVar{"needFinish", true},
				},
			)},
		},
		[]interface{}{PrimaryIs{"expectSpaceOrQwItemOrDelimiter"},
			SubLogic{CreateLogicDef(
				unexpectedEOF,
				[]interface{}{IsUnicodeSpace},
				[]interface{}{IsDelimiterRune{},
					SetVar{"needFinish", true},
				},
				[]interface{}{
					PushItem{ItemType: "qwItem", ItemString: FromCurrRune{}, Delimiter: FromParentItem{}},
					SetPrimary{"expectEndOfQwItem"},
				},
			)},
		},
		[]interface{}{PrimaryIs{"expectEndOfQwItem"},
			SubLogic{CreateLogicDef(
				unexpectedEOF,
				[]interface{}{IsUnicodeSpace, IsDelimiterRune{},
					PushRune{},
					SetVar{"needFinish", true},
				},
				[]interface{}{
					AppendCurrRune{},
				},
			)},
		},
		[]interface{}{PrimaryIs{"expectContentOf"},
			SubLogic{CreateLogicDef(
				unexpectedEOF,
				[]interface{}{IsDelimiterRune{},
					SetVar{"needFinish", true},
				},
				[]interface{}{'\\',
					ChangePrimary{"expectEscapedContentOf"},
				},
				[]interface{}{
					AppendCurrRune{},
				},
			)},
		},
		[]interface{}{PrimaryIs{"expectDigit"},
			SubLogic{CreateLogicDef(
				[]interface{}{IsUnicodeDigit, SecondaryIs{""},
					AppendCurrRune{},
					ChangeSecondary{"orUnderscoreOrDot"},
				},
				[]interface{}{'.', SecondaryIs{"orUnderscoreOrDot"},
					AppendCurrRune{},
					ChangeSecondary{"orUnderscore"},
				},
				[]interface{}{'_', IsUnicodeDigit, SecondaryIs{"orUnderscoreOrDot"}, SecondaryIs{"orUnderscore"},
					AppendCurrRune{},
				},
				[]interface{}{SecondaryIs{""},
					SetError{UnexpectedRune},
				},
				[]interface{}{
					PushRune{},
					SetVar{"needFinish", true},
				},
			)},
		},
		[]interface{}{PrimaryIs{"expectSpaceOrMapKey"},
			SubLogic{CreateLogicDef(
				[]interface{}{IsUnicodeSpace},
				[]interface{}{IsUnicodeLetter,
					PushItem{ItemType: "key", ItemString: FromCurrRune{}},
					SetPrimary{"expectWord"},
				},
				[]interface{}{'"', '\'',
					PushItem{ItemType: "key", Delimiter: FromCurrRune{}},
					SetSecondary{"expectContentOf", "keyToken"},
				},
				[]interface{}{',', SecondaryIs{"orMapValueSeparator"},
					SetPrimary{"expectSpaceOrMapKey"},
				},
				[]interface{}{IsDelimiterRune{}, TopItemIs{"map"},
					SetVar{"needFinish", true},
				},
				unexpectedRune,
			)},
		},
		[]interface{}{PrimaryIs{"expectEscapedContentOf"},
			SubLogic{CreateLogicDef(
				[]interface{}{'"', '\'', '\\',
					AppendCurrRune{},
					ChangePrimary{"expectContentOf"},
				},
				[]interface{}{DelimiterIs{'"'},
					SubLogic{CreateLogicDef(
						[]interface{}{'a', AppendRune{'\a'}},
						[]interface{}{'b', AppendRune{'\b'}},
						[]interface{}{'f', AppendRune{'\f'}},
						[]interface{}{'n', AppendRune{'\n'}},
						[]interface{}{'r', AppendRune{'\r'}},
						[]interface{}{'t', AppendRune{'\t'}},
						[]interface{}{'v', AppendRune{'\v'}},
						unexpectedRune,
					)},
					ChangePrimary{"expectContentOf"},
				},
				unexpectedRune,
			)},
		},
		[]interface{}{PrimaryIs{"expectValueOrSpace"},
			SubLogic{CreateLogicDef(
				[]interface{}{IsEOF{}, StackLenIs{0},
					SetPrimary{"expectEOF"},
				},
				unexpectedEOF,
				[]interface{}{'=', SecondaryIs{"orMapKeySeparator"},
					SetPrimary{"expectRocket"},
				},
				[]interface{}{':', SecondaryIs{"orMapKeySeparator"},
					SetPrimary{"expectValueOrSpace"},
				},
				[]interface{}{',', SecondaryIs{"orArrayItemSeparator"},
					SetPrimary{"expectValueOrSpace"},
				},
				[]interface{}{IsUnicodeSpace},
				[]interface{}{'{',
					PushItem{ItemType: "map", Delimiter: PairForCurrRune{}},
					SetPrimary{"expectSpaceOrMapKey"},
				},
				[]interface{}{'<',
					PushItem{ItemType: "qw", Delimiter: PairForCurrRune{}},
					SetPrimary{"expectSpaceOrQwItemOrDelimiter"},
				},
				[]interface{}{'[',
					PushItem{ItemType: "array", Delimiter: PairForCurrRune{}},
					SetPrimary{"expectValueOrSpace"},
				},
				[]interface{}{TopItemIs{"array"}, IsDelimiterRune{},
					SetVar{"needFinish", true},
				},
				[]interface{}{'-', '+',
					PushItem{ItemType: "number", ItemString: FromCurrRune{}},
					SetPrimary{"expectDigit"},
				},
				[]interface{}{IsUnicodeDigit,
					PushItem{ItemType: "number", ItemString: FromCurrRune{}},
					SetSecondary{"expectDigit", "orUnderscoreOrDot"},
				},
				[]interface{}{'"', '\'',
					PushItem{ItemType: "string", Delimiter: FromCurrRune{}},
					SetSecondary{"expectContentOf", "stringToken"},
				},
				[]interface{}{IsUnicodeLetter,
					PushItem{ItemType: "word", ItemString: FromCurrRune{}},
					SetPrimary{"expectWord"},
				},
				unexpectedRune,
			)},
		},
		[]interface{}{PrimaryIs{"expectValueOrSpace"},
			Unreachable{},
		},
	)

	result := CreateLogicDef(
		[]interface{}{
			PullRune{},
			SetVar{"needFinish", false},
			SubLogic{primaryStateLogic},
			SubLogic{CreateLogicDef(
				[]interface{}{VarIs{"needFinish", true},
					SetVar{"skipPostProcess", false},
					SubLogic{finishLogic},
					SubLogic{CreateLogicDef(
						[]interface{}{VarIs{"skipPostProcess", false},
							SubLogic{postProcessLogic},
						},
					)},
				},
			)},
		},
	)

	return result
}
