package pfa

import (
	"github.com/baza-winner/bwcore/pfa/core"
	"github.com/baza-winner/bwcore/pfa/r"
	"github.com/baza-winner/bwcore/runeprovider"
)

func ParseLogic(p runeprovider.RuneProvider) (interface{}, error) {

	rules := r.RulesFrom(

	// []interface{}{PrimaryIs{"begin"},
	// 	PushItem{},
	// 	SetState{"start", "waitForNonSpace"},
	// 	// SetVar{"waitForSpace", true},
	// },
	// []interface{}{
	// 	PullRune{},
	// 	SubRules{RulesFrom(
	// 		[]interface{}{SecondaryIs{"waitForSpace"}, Space, SetSecondary{"waitForNonSpace"}},
	// 		[]interface{}{SecondaryIs{"waitForNonSpace"}, Space},
	// 		[]interface{}{SecondaryIs{"waitForNonSpace"}, SecondaryIs{"waitForSpace"}, '/', SetSecondary{"waitForSecondCommentRune"}},
	// 		[]interface{}{SecondaryIs{"waitForSecondCommentRune"}, '/', SetSecondary{"singleLineComment"}},
	// 		[]interface{}{SecondaryIs{"waitForSecondCommentRune"}, '*', SetSecondary{"multiLineComment"}},
	// 		[]interface{}{SecondaryIs{"singleLineComment"}, '\n', SetSecondary{"waitForNonSpace"}},
	// 		[]interface{}{SecondaryIs{"singleLineComment"}},
	// 		[]interface{}{SecondaryIs{"multiLineComment"}, '*', SetSecondary{"waitForSlash"}},
	// 		[]interface{}{SecondaryIs{"multiLineComment"}},
	// 		[]interface{}{SecondaryIs{"waitForSlash"}, '/', SetSecondary{"waitForNonSpace"}},

	// 		[]interface{}{PrimaryIs{"start"}, Lower,
	// 			SetState{"entityId", ""}, PushItem{}, SetTopItemType{"key"}, AppendCurrRune{},
	// 		},
	// 		[]interface{}{PrimaryIs{"entityId"}, Letter, Digit,
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
	// 		[]interface{}{StateIs{"rule", "waitForNonSpace"}, Upper,
	// 			PushItem{}, SetTopItemType{"ruleItem"}, AppendCurrRune{}, SetSecondary{"name"},
	// 		},
	// 		[]interface{}{StateIs{"rule", "name"}, Letter, Digit,
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
	// 		[]interface{}{StateIs{"ruleItem", "waitForNonSpace"}, Upper,
	// 			SetState{"ruleItem", "waitForSpaceOr}"},
	// 		},
	// 		[]interface{}{StateIs{"ruleItem", "waitForSpaceOr}"}, Space, '}', TopItemStringIs{"SetError"},
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
	// 	SubRules{RulesFrom(
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
	return Run(p, rules, core.TraceBrief)
}

func CompileLogic(data interface{}) (result r.Rules, err error) {
	return
}
