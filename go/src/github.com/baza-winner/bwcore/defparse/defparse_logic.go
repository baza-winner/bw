package defparse

import (
	"github.com/baza-winner/bwcore/pfa"
	"github.com/baza-winner/bwcore/pfa/a"
	"github.com/baza-winner/bwcore/pfa/b"
	"github.com/baza-winner/bwcore/pfa/c"
	"github.com/baza-winner/bwcore/pfa/e"
	"github.com/baza-winner/bwcore/pfa/r"
	"github.com/baza-winner/bwcore/pfa/val"
)

var rules r.Rules

func init() {
	rules = prepareRules()
}

func prepareRules() r.Rules {

	unexpectedEOF := []interface{}{val.EOF{}, e.UnexpectedRune{}}

	unexpectedRune := []interface{}{e.UnexpectedRune{}}

	finishLogic := r.RulesFrom(
		[]interface{}{c.VarIs{"0.type", "string"}, c.VarIs{"0.type", "qwItem"},
			a.SetVar{"0.result", val.Var{"0.string"}},
		},
		[]interface{}{c.VarIs{"0.type", "map"}},

		[]interface{}{c.VarIs{"0.type", "array"}, c.VarIs{"0.type", "qw"}},

		[]interface{}{c.VarIs{"0.type", "number"},
			a.SetVarBy{"0.result", val.Var{"0.string"}, b.By{b.ParseNumber{}}},
		},
		[]interface{}{c.VarIs{"0.type", "word"},
			pfa.SubRules{r.RulesFrom(
				[]interface{}{c.VarIs{"0.string", "true"},
					a.SetVar{"0.result", true},
				},
				[]interface{}{c.VarIs{"0.string", "false"},
					a.SetVar{"0.result", false},
				},
				[]interface{}{
					c.VarIs{"0.string", "nil"}, c.VarIs{"0.string", "null"},
				},
				[]interface{}{
					c.VarIs{"0.string", "Bool"},
					c.VarIs{"0.string", "String"},
					c.VarIs{"0.string", "Int"},
					c.VarIs{"0.string", "Number"},
					c.VarIs{"0.string", "Map"},
					c.VarIs{"0.string", "Array"},
					c.VarIs{"0.string", "ArrayOf"},
					a.SetVar{"0.result", val.Var{"0.string"}},
				},
				[]interface{}{c.VarIs{"0.string", "qw"},
					a.PullRune{},
					pfa.SubRules{r.RulesFrom(
						[]interface{}{val.UnicodeOpenBraces, val.UnicodePunct, val.UnicodeSymbol,
							a.SetVar{"primary", "expectSpaceOrQwItemOrDelimiter"},
							a.SetVar{"secondary", ""},
							a.SetVar{"0.type", "qw"},
							a.SetVar{"0.result", val.Array{}},
							pfa.SubRules{r.RulesFrom(
								[]interface{}{'<', a.SetVar{"0.delimiter", '>'}},
								[]interface{}{'[', a.SetVar{"0.delimiter", ']'}},
								[]interface{}{'(', a.SetVar{"0.delimiter", ')'}},
								[]interface{}{'{', a.SetVar{"0.delimiter", '}'}},
								[]interface{}{a.SetVar{"0.delimiter", val.Var{"rune"}}},
							)},
						},
						unexpectedRune,
					)},
					a.SetVar{"skipPostProcess", true},
				},
				[]interface{}{
					e.UnexpectedItem{},
				},
			)},
		},
	)

	postProcessLogic := r.RulesFrom(
		[]interface{}{c.VarIs{"stackLen", 0}},
		[]interface{}{c.VarIs{"stackLen", 1},
			a.SetVar{"primary", "end"}, a.SetVar{"secondary", "orSpace"},
		},
		[]interface{}{
			pfa.SubRules{r.RulesFrom(
				[]interface{}{c.VarIs{"1.type", "qw"},
					a.SetVarBy{"1.result", val.Var{"0.result"}, b.By{b.Append{}}},
					a.SetVar{"primary", "expectSpaceOrQwItemOrDelimiter"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{c.VarIs{"1.type", "array"},
					pfa.SubRules{r.RulesFrom(
						[]interface{}{c.VarIs{"0.type", "qw"},
							a.SetVarBy{"1.result", val.Var{"0.result"}, b.By{b.AppendSlice{}}},
						},
						[]interface{}{
							a.SetVarBy{"1.result", val.Var{"0.result"}, b.By{b.Append{}}},
						},
					)},
					a.SetVar{"primary", "begin"}, a.SetVar{"secondary", "orArrayItemSeparator"},
				},
				[]interface{}{c.VarIs{"1.type", "map"},
					pfa.SubRules{r.RulesFrom(
						[]interface{}{c.VarIs{"0.type", "key"},
							a.SetVar{"1.key", val.Var{"0.string"}},
							a.SetVar{"primary", "begin"}, a.SetVar{"secondary", "orMapKeySeparator"},
						},
						[]interface{}{
							a.SetVar{"1.result.{1.key}", val.Var{"0.result"}},
							a.SetVar{"primary", "expectSpaceOrMapKey"}, a.SetVar{"secondary", "orMapValueSeparator"},
						},
					)},
				},
				[]interface{}{
					a.Panic{"unreachable", nil},
				},
			)},
			a.PopItem{},
		},
	)

	primaryStateLogic := r.RulesFrom(
		[]interface{}{c.VarIs{"primary", "end"},
			pfa.SubRules{r.RulesFrom(
				[]interface{}{val.EOF{}, a.SetVar{"primary", "end"}, a.SetVar{"secondary", ""}},
				[]interface{}{val.UnicodeSpace},
				unexpectedRune,
			)},
		},
		[]interface{}{c.VarIs{"primary", "expectRocket"},
			pfa.SubRules{r.RulesFrom(
				[]interface{}{'>', a.SetVar{"primary", "begin"}, a.SetVar{"secondary", ""}},
				unexpectedRune,
			)},
		},
		[]interface{}{c.VarIs{"primary", "expectWord"},
			pfa.SubRules{r.RulesFrom(
				[]interface{}{val.UnicodeLetter, val.UnicodeDigit,
					a.SetVarBy{"0.string", val.Var{"rune"}, b.By{b.Append{}}},
				},
				[]interface{}{
					a.PushRune{},
					a.SetVar{"needFinish", true},
				},
			)},
		},
		[]interface{}{c.VarIs{"primary", "expectSpaceOrQwItemOrDelimiter"},
			pfa.SubRules{r.RulesFrom(
				unexpectedEOF,
				[]interface{}{val.UnicodeSpace},
				[]interface{}{c.VarIs{"rune", val.Var{"0.delimiter"}},
					a.SetVar{"needFinish", true},
				},
				[]interface{}{
					a.PushItem{},
					a.SetVar{"0.type", "qwItem"},
					a.SetVar{"0.delimiter", val.Var{"1.delimiter"}},
					a.SetVarBy{"0.string", val.Var{"rune"}, b.By{b.Append{}}},
					a.SetVar{"primary", "expectEndOfQwItem"}, a.SetVar{"secondary", ""},
				},
			)},
		},
		[]interface{}{c.VarIs{"primary", "expectEndOfQwItem"},
			pfa.SubRules{r.RulesFrom(
				unexpectedEOF,
				[]interface{}{val.UnicodeSpace, c.VarIs{"rune", val.Var{"0.delimiter"}},
					a.PushRune{},
					a.SetVar{"needFinish", true},
				},
				[]interface{}{
					a.SetVarBy{"0.string", val.Var{"rune"}, b.By{b.Append{}}},
				},
			)},
		},
		[]interface{}{c.VarIs{"primary", "expectContentOf"},
			pfa.SubRules{r.RulesFrom(
				unexpectedEOF,
				[]interface{}{c.VarIs{"rune", val.Var{"0.delimiter"}},
					a.SetVar{"needFinish", true},
				},
				[]interface{}{'\\',
					a.SetVar{"primary", "expectEscapedContentOf"},
				},
				[]interface{}{
					a.SetVarBy{"0.string", val.Var{"rune"}, b.By{b.Append{}}},
				},
			)},
		},
		[]interface{}{c.VarIs{"primary", "expectDigit"},
			pfa.SubRules{r.RulesFrom(
				[]interface{}{val.UnicodeDigit, c.VarIs{"secondary", ""},
					a.SetVarBy{"0.string", val.Var{"rune"}, b.By{b.Append{}}},
					a.SetVar{"secondary", "orUnderscoreOrDot"},
				},
				[]interface{}{'.', c.VarIs{"secondary", "orUnderscoreOrDot"},
					a.SetVarBy{"0.string", val.Var{"rune"}, b.By{b.Append{}}},
					a.SetVar{"secondary", "orUnderscore"},
				},
				[]interface{}{'_', val.UnicodeDigit, c.VarIs{"secondary", "orUnderscoreOrDot"}, c.VarIs{"secondary", "orUnderscore"},
					a.SetVarBy{"0.string", val.Var{"rune"}, b.By{b.Append{}}},
				},
				[]interface{}{c.VarIs{"secondary", ""},
					e.UnexpectedRune{},
				},
				[]interface{}{
					a.PushRune{},
					a.SetVar{"needFinish", true},
				},
			)},
		},
		[]interface{}{c.VarIs{"primary", "expectSpaceOrMapKey"},
			pfa.SubRules{r.RulesFrom(
				[]interface{}{val.UnicodeSpace},
				[]interface{}{val.UnicodeLetter,
					a.PushItem{},
					a.SetVar{"0.type", "key"},
					a.SetVarBy{"0.string", val.Var{"rune"}, b.By{b.Append{}}},
					a.SetVar{"primary", "expectWord"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{'"', '\'',
					a.PushItem{},
					a.SetVar{"0.type", "key"},
					a.SetVar{"0.string", ""},
					a.SetVar{"0.delimiter", val.Var{"rune"}},
					a.SetVar{"primary", "expectContentOf"}, a.SetVar{"secondary", "keyToken"},
				},
				[]interface{}{',', c.VarIs{"secondary", "orMapValueSeparator"},
					a.SetVar{"primary", "expectSpaceOrMapKey"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{c.VarIs{"rune", val.Var{"0.delimiter"}}, c.VarIs{"0.type", "map"},
					a.SetVar{"needFinish", true},
				},
				unexpectedRune,
			)},
		},
		[]interface{}{c.VarIs{"primary", "expectEscapedContentOf"},
			pfa.SubRules{r.RulesFrom(
				[]interface{}{'"', '\'', '\\',
					a.SetVarBy{"0.string", val.Var{"rune"}, b.By{b.Append{}}},
					a.SetVar{"primary", "expectContentOf"},
				},
				[]interface{}{c.VarIs{"0.delimiter", '"'},
					pfa.SubRules{r.RulesFrom(
						[]interface{}{'a', a.SetVarBy{"0.string", '\a', b.By{b.Append{}}}},
						[]interface{}{'b', a.SetVarBy{"0.string", '\b', b.By{b.Append{}}}},
						[]interface{}{'f', a.SetVarBy{"0.string", '\f', b.By{b.Append{}}}},
						[]interface{}{'n', a.SetVarBy{"0.string", '\n', b.By{b.Append{}}}},
						[]interface{}{'r', a.SetVarBy{"0.string", '\r', b.By{b.Append{}}}},
						[]interface{}{'t', a.SetVarBy{"0.string", '\t', b.By{b.Append{}}}},
						[]interface{}{'v', a.SetVarBy{"0.string", '\v', b.By{b.Append{}}}},
						unexpectedRune,
					)},
					a.SetVar{"primary", "expectContentOf"},
				},
				unexpectedRune,
			)},
		},
		[]interface{}{c.VarIs{"primary", "begin"},
			pfa.SubRules{r.RulesFrom(
				[]interface{}{val.EOF{}, c.VarIs{"stackLen", 0},
					a.SetVar{"primary", "end"}, a.SetVar{"secondary", ""},
				},
				unexpectedEOF,
				[]interface{}{'=', c.VarIs{"secondary", "orMapKeySeparator"},
					a.SetVar{"primary", "expectRocket"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{':', c.VarIs{"secondary", "orMapKeySeparator"},
					a.SetVar{"primary", "begin"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{',', c.VarIs{"secondary", "orArrayItemSeparator"},
					a.SetVar{"primary", "begin"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{val.UnicodeSpace},
				[]interface{}{'{',
					a.PushItem{},
					a.SetVar{"0.result", val.Map{}},
					a.SetVar{"0.type", "map"},
					a.SetVar{"0.delimiter", '}'},
					a.SetVar{"primary", "expectSpaceOrMapKey"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{'<',
					a.PushItem{},
					a.SetVar{"0.result", val.Array{}},
					a.SetVar{"0.type", "qw"},
					a.SetVar{"0.delimiter", '>'},
					a.SetVar{"primary", "expectSpaceOrQwItemOrDelimiter"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{'[',
					a.PushItem{},
					a.SetVar{"0.result", val.Array{}},
					a.SetVar{"0.type", "array"},
					a.SetVar{"0.delimiter", ']'},
					a.SetVar{"primary", "begin"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{c.VarIs{"0.type", "array"}, c.VarIs{"rune", val.Var{"0.delimiter"}},
					a.SetVar{"needFinish", true},
				},
				[]interface{}{'-', '+',
					a.PushItem{},
					a.SetVar{"0.type", "number"},
					a.SetVarBy{"0.string", val.Var{"rune"}, b.By{b.Append{}}},
					a.SetVar{"primary", "expectDigit"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{val.UnicodeDigit,
					a.PushItem{},
					a.SetVar{"0.type", "number"},
					a.SetVarBy{"0.string", val.Var{"rune"}, b.By{b.Append{}}},
					a.SetVar{"primary", "expectDigit"}, a.SetVar{"secondary", "orUnderscoreOrDot"},
				},
				[]interface{}{'"', '\'',
					a.PushItem{},
					a.SetVar{"0.type", "string"},
					a.SetVar{"0.string", ""},
					a.SetVar{"0.delimiter", val.Var{"rune"}},
					a.SetVar{"primary", "expectContentOf"}, a.SetVar{"secondary", "stringToken"},
				},
				[]interface{}{val.UnicodeLetter,
					a.PushItem{},
					a.SetVar{"0.type", "word"},
					a.SetVarBy{"0.string", val.Var{"rune"}, b.By{b.Append{}}},
					a.SetVar{"primary", "expectWord"}, a.SetVar{"secondary", ""},
				},
				unexpectedRune,
			)},
		},
		[]interface{}{c.VarIs{"primary", "begin"},
			a.Panic{"unreachable", nil},
		},
	)

	result := r.RulesFrom(
		[]interface{}{
			a.PullRune{},
			pfa.SubRules{r.RulesFrom(
				[]interface{}{c.VarIs{"primary", nil},
					a.SetVar{"primary", "begin"},
				},
			)},
			a.SetVar{"needFinish", false},
			pfa.SubRules{primaryStateLogic},
			pfa.SubRules{r.RulesFrom(
				[]interface{}{c.VarIs{"needFinish", true},
					a.SetVar{"skipPostProcess", false},
					pfa.SubRules{finishLogic},
					pfa.SubRules{r.RulesFrom(
						[]interface{}{c.VarIs{"skipPostProcess", false},
							pfa.SubRules{postProcessLogic},
						},
					)},
				},
			)},
		},
	)

	return result
}
