package pfa

import "github.com/baza-winner/bwcore/runeprovider"

func ParseLogic(p runeprovider.RuneProvider) (data interface{}, errinitialState error) {

	logicParseDef := CreateRules(

	// []interface{}{PrimaryIs{"begin"},
	// 	PushItem{},
	// 	SetState{"start", "waitForNonSpace"},
	// 	// SetVar{"waitForSpace", true},
	// },
	// []interface{}{
	// 	PullRune{},
	// 	SubRules{CreateRules(
	// 		[]interface{}{SecondaryIs{"waitForSpace"}, UnicodeSpace, SetSecondary{"waitForNonSpace"}},
	// 		[]interface{}{SecondaryIs{"waitForNonSpace"}, UnicodeSpace},
	// 		[]interface{}{SecondaryIs{"waitForNonSpace"}, SecondaryIs{"waitForSpace"}, '/', SetSecondary{"waitForSecondCommentRune"}},
	// 		[]interface{}{SecondaryIs{"waitForSecondCommentRune"}, '/', SetSecondary{"singleLineComment"}},
	// 		[]interface{}{SecondaryIs{"waitForSecondCommentRune"}, '*', SetSecondary{"multiLineComment"}},
	// 		[]interface{}{SecondaryIs{"singleLineComment"}, '\n', SetSecondary{"waitForNonSpace"}},
	// 		[]interface{}{SecondaryIs{"singleLineComment"}},
	// 		[]interface{}{SecondaryIs{"multiLineComment"}, '*', SetSecondary{"waitForSlash"}},
	// 		[]interface{}{SecondaryIs{"multiLineComment"}},
	// 		[]interface{}{SecondaryIs{"waitForSlash"}, '/', SetSecondary{"waitForNonSpace"}},

	// 		[]interface{}{PrimaryIs{"start"}, UnicodeLower,
	// 			SetState{"entityId", ""}, PushItem{}, SetTopItemType{"key"}, AppendCurrRune{},
	// 		},
	// 		[]interface{}{PrimaryIs{"entityId"}, UnicodeLetter, UnicodeDigit,
	// 			AppendCurrRune{},
	// 		},
	// 		[]interface{}{PrimaryIs{"entityId"},
	// 			PopItem{}, SetTopItemStringFromSubItem{}, SetState{"expression", "waitForNonSpace"},
	// 		},
	// 		[]interface{}{PrimaryIs{"expression"}, '(',
	// 			PushItem{}, SetTopItemType{"rule"}, SetState{"rule", "waitForNonSpace"},
	// 		},
	// 		[]interface{}{PrimaryIs{"expression"},
	// 			SetError{UnexpectedRune},
	// 		},
	// 		[]interface{}{StateIs{"rule", "waitForNonSpace"}, UnicodeUpper,
	// 			PushItem{}, SetTopItemType{"ruleItem"}, AppendCurrRune{}, SetSecondary{"name"},
	// 		},
	// 		[]interface{}{StateIs{"rule", "name"}, UnicodeLetter, UnicodeDigit,
	// 			AppendCurrRune{},
	// 		},
	// 		[]interface{}{StateIs{"rule", "name"}, TopItemStringIs{"IsEOF"},
	// 			PushRune{},
	// 			SetTopItemValueAsString{},
	// 			PopItem{},
	// 			AppendItemArray{FromSubItemValue{}},
	// 			SetState{"rule", "waitForSpace"},
	// 		},
	// 		[]interface{}{StateIs{"rule", "name"}, TopItemStringIs{"SetError"},
	// 			SetTopItemValueAsString{},
	// 			// PushItem{}, SetTopItemType{"ruleItemParams"},
	// 			// PopItem{},
	// 			// AppendItemArray{FromSubItemValue{}},
	// 			PushRune{},
	// 			SetState{"ruleItem", "{"},
	// 		},
	// 		[]interface{}{StateIs{"ruleItem", "{"}, '{',
	// 			SetState{"ruleItem", "waitForNonSpace"},
	// 		},
	// 		[]interface{}{StateIs{"ruleItem", "waitForNonSpace"}, UnicodeUpper,
	// 			SetState{"ruleItem", "waitForSpaceOr}"},
	// 		},
	// 		[]interface{}{StateIs{"ruleItem", "waitForSpaceOr}"}, UnicodeSpace, '}', TopItemStringIs{"SetError"},
	// 			SetState{"ruleItem", "waitForSpace"},
	// 		},
	// 		// []interface{}{StateIs{"ruleItem", "waitForSpaceOr}"},
	// 		//   SetState{"ruleItem", "waitForSpace"},
	// 		// },
	// 		// []interface{}{StateIs{"rule", "{"}, '{',
	// 		//   SetState{"ruleItem", "waitForNonSpace"},
	// 		// },

	// 		// []interface{}{PrimaryIs{"expression"}, '(', PushItem{}, SetTopItemType{"rule"}, SetState{"rule", "waitForSpace"}},
	// 		[]interface{}{SetError{UnexpectedRune}},
	// 	)},
	// },

	// []interface{}{SecondaryIs{"waitForSecondCommentRune"}, '*',
	// 	SetState{"multiLineComment"},
	// },
	// []interface{}{PrimaryIs{"start"},
	// 	PullRune{},
	// 	SubRules{CreateRules(
	// 		[]interface{}{'/',
	// 			SetState{"waitForSecondCommentRune"},
	// 		},
	// 	)},
	// 	// Debug{"THERE"},
	// 	// SetError{UnexpectedRune},
	// 	// PushItem{},
	// 	// SetState{"start"},
	// },
	)
	return Run(p, logicParseDef, TraceBrief)
}

func CompileLogic(data interface{}) (result Rules, err error) {
	return
}
