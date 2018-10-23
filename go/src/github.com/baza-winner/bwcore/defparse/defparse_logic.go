package defparse

import (
	. "github.com/baza-winner/bwcore/pfa"
)

var pfaStateDef Rules

func init() {
	pfaStateDef = prepareLogicDef()
}

func prepareLogicDef() Rules {

	unexpectedEOF := []interface{}{EOF{}, SetVar{"error", "unexpectedRune"}}

	unexpectedRune := []interface{}{SetVar{"error", "unexpectedRune"}}

	finishLogic := CreateRules(
		[]interface{}{VarIs{"0.type", "string"}, VarIs{"0.type", "qwItem"},
			SetVar{"0.result", Var{"0.string"}},
			// SetTopItemValueAsString{},
		},
		[]interface{}{VarIs{"0.type", "map"},
			SetVar{"0.result", Var{"0.map"}},
		},
		[]interface{}{VarIs{"0.type", "array"}, VarIs{"0.type", "qw"},
			SetVar{"0.result", Var{"0.array"}},
			// SetTopItemValueAsArray{},
		},
		[]interface{}{VarIs{"0.type", "number"},
			SetVarBy{"0.result", Var{"0.string"}, By{ParseNumber{}}},
			// ChangeVar{"0.result", ParseNumber{}},
			// SetTopItemValueAsNumber{},
		},
		[]interface{}{VarIs{"0.type", "word"},
			SubRules{CreateRules(
				[]interface{}{VarIs{"0.string", "true"},
					SetVar{"0.result", true},
					// SetTopItemValueAsBool{true},
				},
				[]interface{}{VarIs{"0.string", "false"},
					SetVar{"0.result", true},
					// SetTopItemValueAsBool{false},
				},
				[]interface{}{
					VarIs{"0.string", "nil"}, VarIs{"0.string", "nul"},
					// TopItemStringIsOneOf{bwset.StringSetFrom("nil", "null")}
				},
				[]interface{}{
					VarIs{"0.string", "Bool"},
					VarIs{"0.string", "String"},
					VarIs{"0.string", "Int"},
					VarIs{"0.string", "Number"},
					VarIs{"0.string", "Map"},
					VarIs{"0.string", "Array"},
					VarIs{"0.string", "ArrayOf"},
					// TopItemStringIsOneOf{bwset.StringSetFrom("Bool", "String", "Int", "Number", "Map", "Array", "ArrayOf")},
					SetVar{"0.result", Var{"0.string"}},
					// SetTopItemValueAsString{},
				},
				[]interface{}{VarIs{"0.string", "qw"},
					PullRune{},
					SubRules{CreateRules(
						[]interface{}{UnicodeOpenBraces, UnicodePunct, UnicodeSymbol,
							SetVar{"primary", "expectSpaceOrQwItemOrDelimiter"},
							SetVar{"secondary", ""},
							// SetVar{"primary", "expectSpaceOrQwItemOrDelimiter"}, SetVar{"secondary", ""},
							// SetVarBy{"0.delimiter", Var{"rune"}, []interface{PairForCurrRune{}}},
							SubRules{CreateRules(
								[]interface{}{'<', SetVar{"0.delimiter", '>'}},
								[]interface{}{'[', SetVar{"0.delimiter", ']'}},
								[]interface{}{'(', SetVar{"0.delimiter", ')'}},
								[]interface{}{'{', SetVar{"0.delimiter", '}'}},
								[]interface{}{SetVar{"0.delimiter", Var{"rune"}}},
							)},
							// SetTopItemDelimiter{PairForCurrRune{}},
							SetVar{"0.type", "qw"},
						},
						unexpectedRune,
					)},
					SetVar{"skipPostProcess", true},
				},
				[]interface{}{
					SetVar{"error", "unknownWord"},
				},
			)},
		},
	)

	postProcessLogic := CreateRules(
		[]interface{}{VarIs{"stackLen", 0}},
		[]interface{}{VarIs{"stackLen", 1},
			SetVar{"primary", "end"}, SetVar{"secondary", "orSpace"},
		},
		[]interface{}{
			// PopItem{},
			SubRules{CreateRules(
				[]interface{}{VarIs{"1.type", "qw"},
					SetVarBy{"1.array", Var{"0.value"}, By{Append{}}},
					// AppendItemArray{FromSubItemValue{}},
					SetVar{"primary", "expectSpaceOrQwItemOrDelimiter"}, SetVar{"secondary", ""},
				},
				[]interface{}{VarIs{"1.type", "array"},
					SubRules{CreateRules(
						[]interface{}{VarIs{"0.type", "qw"}, //SubItemIs{"qw"},
							SetVarBy{"1.array", Var{"0.array"}, By{Append{}}},
							// AppendItemArray{FromSubItemArray{}},
						},
						[]interface{}{
							SetVarBy{"1.array", Var{"0.value"}, By{Append{}}},
							// AppendItemArray{FromSubItemValue{}},
							// AppendItemArray{FromSubItemValue{}},
						},
					)},
					SetVar{"primary", "begin"}, SetVar{"secondary", "orArrayItemSeparator"},
				},
				[]interface{}{VarIs{"1.type", "map"},
					SubRules{CreateRules(
						[]interface{}{VarIs{"0.type", "key"}, //SubItemIs{"key"},
							SetVar{"1.string", Var{"0.string"}},
							// SetTopItemStringFromSubItem{},
							SetVar{"primary", "begin"}, SetVar{"secondary", "orMapKeySeparator"},
						},
						[]interface{}{
							SetVar{"1.map.{1.key}", Var{"0.value"}},
							// SetVarKeyFrom{"1.map", Var{"1.key"}, Var{"0.value"}},
							// SetTopItemMapKeyValueFromSubItem{},
							SetVar{"primary", "expectSpaceOrMapKey"}, SetVar{"secondary", "orMapValueSeparator"},
						},
					)},
				},
				[]interface{}{
					SetVar{"error", "unreachable"},
				},
			)},
			PopItem{},
		},
	)

	primaryStateLogic := CreateRules(
		[]interface{}{VarIs{"primary", "end"},
			SubRules{CreateRules(
				[]interface{}{EOF{}, SetVar{"primary", "end"}, SetVar{"secondary", ""}},
				[]interface{}{UnicodeSpace},
				unexpectedRune,
			)},
		},
		[]interface{}{VarIs{"primary", "expectRocket"},
			SubRules{CreateRules(
				[]interface{}{'>', SetVar{"primary", "begin"}, SetVar{"secondary", ""}},
				unexpectedRune,
			)},
		},
		[]interface{}{VarIs{"primary", "expectWord"},
			SubRules{CreateRules(
				[]interface{}{UnicodeLetter, UnicodeDigit,
					SetVarBy{"0.string", Var{"rune"}, By{Append{}}},
					// AppendToVar{"0.string", Var{"rune"}},
				},
				[]interface{}{
					PushRune{},
					SetVar{"needFinish", true},
				},
			)},
		},
		[]interface{}{VarIs{"primary", "expectSpaceOrQwItemOrDelimiter"},
			SubRules{CreateRules(
				unexpectedEOF,
				[]interface{}{UnicodeSpace},
				[]interface{}{VarIs{"rune", Var{"0.delimiter"}},
					SetVar{"needFinish", true},
				},
				[]interface{}{
					PushItem{},
					SetVar{"0.type", "qwItem"},
					SetVar{"0.delimiter", "1.delimiter"},
					// SetTopItemDelimiter{FromParentItem{}},
					SetVarBy{"0.string", Var{"rune"}, By{Append{}}},
					// AppendToVar{"0.string", Var{"rune"}},
					SetVar{"primary", "expectEndOfQwItem"}, SetVar{"secondary", ""},
				},
			)},
		},
		[]interface{}{VarIs{"primary", "expectEndOfQwItem"},
			SubRules{CreateRules(
				unexpectedEOF,
				[]interface{}{UnicodeSpace, VarIs{"rune", Var{"0.delimiter"}},
					PushRune{},
					SetVar{"needFinish", true},
				},
				[]interface{}{
					SetVarBy{"0.string", Var{"rune"}, By{Append{}}},
					// AppendToVar{"0.string", Var{"rune"}},
				},
			)},
		},
		[]interface{}{VarIs{"primary", "expectContentOf"},
			SubRules{CreateRules(
				unexpectedEOF,
				[]interface{}{VarIs{"rune", Var{"0.delimiter"}},
					// VarIs{"rune", Var{"0.delimiter"}},
					SetVar{"needFinish", true},
				},
				[]interface{}{'\\',
					SetVar{"primary", "expectEscapedContentOf"},
					// SetPrimary{"expectEscapedContentOf"},
				},
				[]interface{}{
					SetVarBy{"0.string", Var{"rune"}, By{Append{}}},
					// AppendToVar{"0.string", Var{"rune"}},
				},
			)},
		},
		[]interface{}{VarIs{"primary", "expectDigit"},
			SubRules{CreateRules(
				[]interface{}{UnicodeDigit, VarIs{"secondary", ""},
					SetVarBy{"0.string", Var{"rune"}, By{Append{}}},
					// AppendToVar{"0.string", Var{"rune"}},
					// SetSecondary{"orUnderscoreOrDot"},
					SetVar{"secondary", "orUnderscoreOrDot"},
				},
				[]interface{}{'.', VarIs{"secondary", "orUnderscoreOrDot"},
					SetVarBy{"0.string", Var{"rune"}, By{Append{}}},
					// AppendToVar{"0.string", Var{"rune"}},
					SetVar{"secondary", "orUnderscore"},
					// SetSecondary{"orUnderscore"},
				},
				[]interface{}{'_', UnicodeDigit, VarIs{"secondary", "orUnderscoreOrDot"}, VarIs{"secondary", "orUnderscore"},
					SetVarBy{"0.string", Var{"rune"}, By{Append{}}},
					// AppendToVar{"0.string", Var{"rune"}},
				},
				[]interface{}{VarIs{"secondary", ""},
					SetVar{"error", "unexpectedRune"},
				},
				[]interface{}{
					PushRune{},
					SetVar{"needFinish", true},
				},
			)},
		},
		[]interface{}{VarIs{"primary", "expectSpaceOrMapKey"},
			SubRules{CreateRules(
				[]interface{}{UnicodeSpace},
				[]interface{}{UnicodeLetter,
					PushItem{},
					SetVar{"0.type", "key"},
					SetVarBy{"0.string", Var{"rune"}, By{Append{}}},
					SetVar{"primary", "expectWord"}, SetVar{"secondary", ""},
				},
				[]interface{}{'"', '\'',
					PushItem{},
					SetVar{"0.type", "key"},
					SetVar{"0.delimiter", Var{"rune"}},
					// SetTopItemDelimiter{FromCurrRune{}},
					SetVar{"primary", "expectContentOf"}, SetVar{"secondary", "keyToken"},
				},
				[]interface{}{',', VarIs{"secondary", "orMapValueSeparator"},
					SetVar{"primary", "expectSpaceOrMapKey"}, SetVar{"secondary", ""},
				},
				[]interface{}{VarIs{"rune", Var{"0.delimiter"}}, VarIs{"0.type", "map"},
					SetVar{"needFinish", true},
				},
				unexpectedRune,
			)},
		},
		[]interface{}{VarIs{"primary", "expectEscapedContentOf"},
			SubRules{CreateRules(
				[]interface{}{'"', '\'', '\\',
					SetVarBy{"0.string", Var{"rune"}, By{Append{}}},
					SetVar{"primary", "expectContentOf"},
					// SetPrimary{"expectContentOf"},
				},
				[]interface{}{VarIs{"0.delimiter", '"'}, // DelimiterIs{'"'},
					SubRules{CreateRules(
						[]interface{}{'a', SetVarBy{"0.string", '\a', By{Append{}}}},
						[]interface{}{'b', SetVarBy{"0.string", '\b', By{Append{}}}},
						[]interface{}{'f', SetVarBy{"0.string", '\f', By{Append{}}}},
						[]interface{}{'n', SetVarBy{"0.string", '\n', By{Append{}}}},
						[]interface{}{'r', SetVarBy{"0.string", '\r', By{Append{}}}},
						[]interface{}{'t', SetVarBy{"0.string", '\t', By{Append{}}}},
						[]interface{}{'v', SetVarBy{"0.string", '\v', By{Append{}}}},
						unexpectedRune,
					)},
					SetVar{"primary", "expectContentOf"},
					// SetPrimary{"expectContentOf"},
				},
				unexpectedRune,
			)},
		},
		[]interface{}{VarIs{"primary", "begin"},
			SubRules{CreateRules(
				[]interface{}{EOF{}, VarIs{"stackLen", 0},
					SetVar{"primary", "end"}, SetVar{"secondary", ""},
				},
				unexpectedEOF,
				[]interface{}{'=', VarIs{"secondary", "orMapKeySeparator"},
					SetVar{"primary", "expectRocket"}, SetVar{"secondary", ""},
				},
				[]interface{}{':', VarIs{"secondary", "orMapKeySeparator"},
					SetVar{"primary", "begin"}, SetVar{"secondary", ""},
				},
				[]interface{}{',', VarIs{"secondary", "orArrayItemSeparator"},
					SetVar{"primary", "begin"}, SetVar{"secondary", ""},
				},
				[]interface{}{UnicodeSpace},
				[]interface{}{'{',
					PushItem{},
					SetVar{"0.type", "map"},
					SetVar{"0.delimiter", '}'},
					SetVar{"primary", "expectSpaceOrMapKey"}, SetVar{"secondary", ""},
				},
				[]interface{}{'<',
					PushItem{},
					SetVar{"0.type", "qw"},
					SetVar{"0.delimiter", '>'},
					SetVar{"primary", "expectSpaceOrQwItemOrDelimiter"}, SetVar{"secondary", ""},
				},
				[]interface{}{'[',
					PushItem{},
					SetVar{"0.type", "array"},
					SetVar{"0.delimiter", ']'},
					SetVar{"primary", "begin"}, SetVar{"secondary", ""},
				},
				[]interface{}{VarIs{"0.type", "array"}, VarIs{"rune", Var{"0.delimiter"}},
					SetVar{"needFinish", true},
				},
				[]interface{}{'-', '+',
					PushItem{},
					SetVar{"0.type", "number"},
					SetVarBy{"0.string", Var{"rune"}, By{Append{}}},
					SetVar{"primary", "expectDigit"}, SetVar{"secondary", ""},
				},
				[]interface{}{UnicodeDigit,
					PushItem{},
					SetVar{"0.type", "number"},
					SetVarBy{"0.string", Var{"rune"}, By{Append{}}},
					SetVar{"primary", "expectDigit"}, SetVar{"secondary", "orUnderscoreOrDot"},
				},
				[]interface{}{'"', '\'',
					PushItem{},
					SetVar{"0.type", "string"},
					SetVar{"0.delimiter", Var{"rune"}},
					// SetTopItemDelimiter{FromCurrRune{}},
					SetVar{"primary", "expectContentOf"}, SetVar{"secondary", "stringToken"},
				},
				[]interface{}{UnicodeLetter,
					PushItem{},
					SetVar{"0.type", "word"},
					SetVarBy{"0.string", Var{"rune"}, By{Append{}}},
					SetVar{"primary", "expectWord"}, SetVar{"secondary", ""},
				},
				unexpectedRune,
			)},
		},
		[]interface{}{VarIs{"primary", "begin"},
			SetVar{"error", "unreachable"},
		},
	)

	result := CreateRules(
		[]interface{}{
			PullRune{},
			SubRules{CreateRules(
				[]interface{}{VarIs{"primary", nil},
					SetVar{"primary", "begin"},
				},
			)},
			SetVar{"needFinish", false},
			SubRules{primaryStateLogic},
			SubRules{CreateRules(
				[]interface{}{VarIs{"needFinish", true},
					SetVar{"skipPostProcess", false},
					SubRules{finishLogic},
					SubRules{CreateRules(
						[]interface{}{VarIs{"skipPostProcess", false},
							SubRules{postProcessLogic},
						},
					)},
				},
			)},
		},
	)

	return result
}
