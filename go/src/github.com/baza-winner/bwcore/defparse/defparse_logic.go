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

	unexpectedEOF := []interface{}{val.EOF{}, e.Unexpected{"runePos"}}

	unexpectedRune := []interface{}{e.Unexpected{"runePos"}}

	finishLogic := r.RulesFrom(
		[]interface{}{c.VarIs{"stack.-1.type", "string"}, c.VarIs{"stack.-1.type", "qwItem"},
			a.SetVar{"stack.-1.result", val.Var{"stack.-1.string"}},
		},
		// []interface{}{c.VarIs{"stack.-1.type", "map"}},

		// []interface{}{c.VarIs{"stack.-1.type", "array"}, c.VarIs{"stack.-1.type", "qw"}},

		[]interface{}{c.VarIs{"stack.-1.type", "number"},
			a.SetVarBy{"stack.-1.result", val.Var{"stack.-1.string"}, b.By{b.ParseNumber{}}},
		},

		[]interface{}{c.VarIs{"stack.-1.type", "word"},
			pfa.SubRules{r.RulesFrom(
				[]interface{}{c.VarIs{"stack.-1.string", "true"},
					a.SetVar{"stack.-1.result", true},
				},
				[]interface{}{c.VarIs{"stack.-1.string", "false"},
					a.SetVar{"stack.-1.result", false},
				},
				[]interface{}{
					c.VarIs{"stack.-1.string", "nil"}, c.VarIs{"stack.-1.string", "null"},
				},
				[]interface{}{
					c.VarIs{"stack.-1.string", "Bool"},
					c.VarIs{"stack.-1.string", "String"},
					c.VarIs{"stack.-1.string", "Int"},
					c.VarIs{"stack.-1.string", "Number"},
					c.VarIs{"stack.-1.string", "Map"},
					c.VarIs{"stack.-1.string", "Array"},
					c.VarIs{"stack.-1.string", "ArrayOf"},
					a.SetVar{"stack.-1.result", val.Var{"stack.-1.string"}},
				},
				[]interface{}{c.VarIs{"stack.-1.string", "qw"},
					a.PullRune{},
					pfa.SubRules{r.RulesFrom(
						[]interface{}{val.OpenBraces, val.Punct, val.Symbol,
							a.SetVar{"primary", "expectSpaceOrQwItemOrDelimiter"},
							a.SetVar{"secondary", ""},
							a.SetVar{"stack.-1.type", "qw"},
							a.SetVar{"stack.-1.result", val.Array{}},
							a.SetVarBy{"stack.-1.delimiter", val.Var{"rune"}, b.By{b.PairBrace{}}},
						},
						unexpectedRune,
					)},
					a.SetVar{"skipPostProcess", true},
				},
				[]interface{}{
					e.Unexpected{"stack.-1.itemPos"},
				},
			)},
		},
	)

	postProcessLogic := r.RulesFrom(
		[]interface{}{c.VarIs{"stack.#", 0}},
		[]interface{}{c.VarIs{"stack.#", 1},
			a.SetVar{"primary", "end"}, a.SetVar{"secondary", "orSpace"},
		},
		[]interface{}{
			pfa.SubRules{r.RulesFrom(
				[]interface{}{c.VarIs{"stack.-2.type", "qw"},
					a.SetVarBy{"stack.-2.result", val.Var{"stack.-1.result"}, b.By{b.Append{}}},
					a.SetVar{"primary", "expectSpaceOrQwItemOrDelimiter"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{c.VarIs{"stack.-2.type", "array"},
					pfa.SubRules{r.RulesFrom(
						[]interface{}{c.VarIs{"stack.-1.type", "qw"},
							a.SetVarBy{"stack.-2.result", val.Var{"stack.-1.result"}, b.By{b.AppendSlice{}}},
						},
						[]interface{}{
							a.SetVarBy{"stack.-2.result", val.Var{"stack.-1.result"}, b.By{b.Append{}}},
						},
					)},
					a.SetVar{"primary", "begin"}, a.SetVar{"secondary", "orArrayItemSeparator"},
				},
				[]interface{}{c.VarIs{"stack.-2.type", "map"},
					pfa.SubRules{r.RulesFrom(
						[]interface{}{c.VarIs{"stack.-1.type", "key"},
							a.SetVar{"stack.-2.key", val.Var{"stack.-1.string"}},
							a.SetVar{"primary", "begin"}, a.SetVar{"secondary", "orMapKeySeparator"},
						},
						[]interface{}{
							a.SetVar{"stack.-2.result.{1.key}", val.Var{"stack.-1.result"}},
							a.SetVar{"primary", "expectSpaceOrMapKey"}, a.SetVar{"secondary", "orMapValueSeparator"},
						},
					)},
				},
				[]interface{}{
					e.Unexpected{""},
					// a.Panic{"unreachable", nil},
				},
			)},
			a.SetVarBy{"stack", val.Var{"stack"}, b.By{b.ButLast{}}},
			// a.PopItem{},
		},
	)

	primaryStateLogic := r.RulesFrom(
		[]interface{}{c.VarIs{"primary", "end"},
			pfa.SubRules{r.RulesFrom(
				[]interface{}{val.EOF{}, a.SetVar{"primary", "end"}, a.SetVar{"secondary", ""}},
				[]interface{}{val.Space},
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
				[]interface{}{val.Letter, val.Digit,
					a.SetVarBy{"stack.-1.string", val.Var{"rune"}, b.By{b.Append{}}},
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
				[]interface{}{val.Space},
				[]interface{}{c.VarIs{"rune", val.Var{"stack.-1.delimiter"}},
					a.SetVar{"needFinish", true},
				},
				[]interface{}{
					a.SetVarBy{"stack", map[string]interface{}{"itemPos": val.Var{"runePos"}}, b.By{b.Append{}}},
					a.SetVar{"stack.-1.type", "qwItem"},
					a.SetVar{"stack.-1.delimiter", val.Var{"stack.-2.delimiter"}},
					a.SetVarBy{"stack.-1.string", val.Var{"rune"}, b.By{b.Append{}}},
					a.SetVar{"primary", "expectEndOfQwItem"}, a.SetVar{"secondary", ""},
				},
			)},
		},
		[]interface{}{c.VarIs{"primary", "expectEndOfQwItem"},
			pfa.SubRules{r.RulesFrom(
				unexpectedEOF,
				[]interface{}{val.Space, c.VarIs{"rune", val.Var{"stack.-1.delimiter"}},
					a.PushRune{},
					a.SetVar{"needFinish", true},
				},
				[]interface{}{
					a.SetVarBy{"stack.-1.string", val.Var{"rune"}, b.By{b.Append{}}},
				},
			)},
		},
		[]interface{}{c.VarIs{"primary", "expectContentOf"},
			pfa.SubRules{r.RulesFrom(
				unexpectedEOF,
				[]interface{}{c.VarIs{"rune", val.Var{"stack.-1.delimiter"}},
					a.SetVar{"needFinish", true},
				},
				[]interface{}{'\\',
					a.SetVar{"primary", "expectEscapedContentOf"},
				},
				[]interface{}{
					a.SetVarBy{"stack.-1.string", val.Var{"rune"}, b.By{b.Append{}}},
				},
			)},
		},
		[]interface{}{c.VarIs{"primary", "expectDigit"},
			pfa.SubRules{r.RulesFrom(
				[]interface{}{val.Digit, c.VarIs{"secondary", ""},
					a.SetVarBy{"stack.-1.string", val.Var{"rune"}, b.By{b.Append{}}},
					a.SetVar{"secondary", "orUnderscoreOrDot"},
				},
				[]interface{}{'.', c.VarIs{"secondary", "orUnderscoreOrDot"},
					a.SetVarBy{"stack.-1.string", val.Var{"rune"}, b.By{b.Append{}}},
					a.SetVar{"secondary", "orUnderscore"},
				},
				[]interface{}{'_', val.Digit, c.VarIs{"secondary", "orUnderscoreOrDot"}, c.VarIs{"secondary", "orUnderscore"},
					a.SetVarBy{"stack.-1.string", val.Var{"rune"}, b.By{b.Append{}}},
				},
				[]interface{}{c.VarIs{"secondary", ""},
					e.Unexpected{"runePos"},
				},
				[]interface{}{
					a.PushRune{},
					a.SetVar{"needFinish", true},
				},
			)},
		},
		[]interface{}{c.VarIs{"primary", "expectSpaceOrMapKey"},
			pfa.SubRules{r.RulesFrom(
				[]interface{}{val.Space},
				[]interface{}{val.Letter,
					a.SetVarBy{"stack", map[string]interface{}{"itemPos": val.Var{"runePos"}}, b.By{b.Append{}}},
					a.SetVar{"stack.-1.type", "key"},
					a.SetVarBy{"stack.-1.string", val.Var{"rune"}, b.By{b.Append{}}},
					a.SetVar{"primary", "expectWord"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{'"', '\'',
					a.SetVarBy{"stack", map[string]interface{}{"itemPos": val.Var{"runePos"}}, b.By{b.Append{}}},
					a.SetVar{"stack.-1.type", "key"},
					a.SetVar{"stack.-1.string", ""},
					a.SetVar{"stack.-1.delimiter", val.Var{"rune"}},
					a.SetVar{"primary", "expectContentOf"}, a.SetVar{"secondary", "keyToken"},
				},
				[]interface{}{',', c.VarIs{"secondary", "orMapValueSeparator"},
					a.SetVar{"primary", "expectSpaceOrMapKey"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{c.VarIs{"rune", val.Var{"stack.-1.delimiter"}}, c.VarIs{"stack.-1.type", "map"},
					a.SetVar{"needFinish", true},
				},
				unexpectedRune,
			)},
		},
		[]interface{}{c.VarIs{"primary", "expectEscapedContentOf"},
			pfa.SubRules{r.RulesFrom(
				[]interface{}{'"', '\'', '\\',
					a.SetVarBy{"stack.-1.string", val.Var{"rune"}, b.By{b.Append{}}},
					a.SetVar{"primary", "expectContentOf"},
				},
				[]interface{}{c.VarIs{"stack.-1.delimiter", '"'},
					pfa.SubRules{r.RulesFrom(
						[]interface{}{'a', 'b', 'f', 'n', 'r', 't', 'v', a.SetVarBy{"stack.-1.string", val.Var{"rune"}, b.By{b.Escape{}, b.Append{}}}},
						unexpectedRune,
					)},
					a.SetVar{"primary", "expectContentOf"},
				},
				unexpectedRune,
			)},
		},
		[]interface{}{c.VarIs{"primary", "begin"},
			pfa.SubRules{r.RulesFrom(
				[]interface{}{val.EOF{}, c.VarIs{"stack.#", 0},
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
				[]interface{}{val.Space},
				[]interface{}{'{',
					a.SetVarBy{"stack", map[string]interface{}{"itemPos": val.Var{"runePos"}}, b.By{b.Append{}}},
					a.SetVar{"stack.-1.result", val.Map{}},
					a.SetVar{"stack.-1.type", "map"},
					a.SetVar{"stack.-1.delimiter", '}'},
					a.SetVar{"primary", "expectSpaceOrMapKey"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{'<',
					a.SetVarBy{"stack", map[string]interface{}{"itemPos": val.Var{"runePos"}}, b.By{b.Append{}}},
					a.SetVar{"stack.-1.result", val.Array{}},
					a.SetVar{"stack.-1.type", "qw"},
					a.SetVar{"stack.-1.delimiter", '>'},
					a.SetVar{"primary", "expectSpaceOrQwItemOrDelimiter"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{'[',
					a.SetVarBy{"stack", map[string]interface{}{"itemPos": val.Var{"runePos"}}, b.By{b.Append{}}},
					a.SetVar{"stack.-1.result", val.Array{}},
					a.SetVar{"stack.-1.type", "array"},
					a.SetVar{"stack.-1.delimiter", ']'},
					a.SetVar{"primary", "begin"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{c.VarIs{"stack.-1.type", "array"}, c.VarIs{"rune", val.Var{"stack.-1.delimiter"}},
					a.SetVar{"needFinish", true},
				},
				[]interface{}{'-', '+',
					// a.SetVarBy{"stack", map[string]interface{}{"itemPos": val.Var{"runePos"}}, b.By{b.Append{}}},
					// a.SetVarBy{"stack", val.Map, b.By{b.Append{}}}, a.SetVar{"stack.-1.itemPos", val.Var{"runePos"}},
					a.SetVarBy{"stack", map[string]interface{}{"itemPos": val.Var{"runePos"}}, b.By{b.Append{}}},
					a.SetVar{"stack.-1.type", "number"},
					a.SetVarBy{"stack.-1.string", val.Var{"rune"}, b.By{b.Append{}}},
					a.SetVar{"primary", "expectDigit"}, a.SetVar{"secondary", ""},
				},
				[]interface{}{val.Digit,
					a.SetVarBy{"stack", map[string]interface{}{"itemPos": val.Var{"runePos"}}, b.By{b.Append{}}},
					a.SetVar{"stack.-1.type", "number"},
					a.SetVarBy{"stack.-1.string", val.Var{"rune"}, b.By{b.Append{}}},
					a.SetVar{"primary", "expectDigit"}, a.SetVar{"secondary", "orUnderscoreOrDot"},
				},
				[]interface{}{'"', '\'',
					a.SetVarBy{"stack", map[string]interface{}{"itemPos": val.Var{"runePos"}}, b.By{b.Append{}}},
					a.SetVar{"stack.-1.type", "string"},
					a.SetVar{"stack.-1.string", ""},
					a.SetVar{"stack.-1.delimiter", val.Var{"rune"}},
					a.SetVar{"primary", "expectContentOf"}, a.SetVar{"secondary", "stringToken"},
				},
				[]interface{}{val.Letter,
					a.SetVarBy{"stack", map[string]interface{}{"itemPos": val.Var{"runePos"}}, b.By{b.Append{}}},
					a.SetVar{"stack.-1.type", "word"},
					a.SetVarBy{"stack.-1.string", val.Var{"rune"}, b.By{b.Append{}}},
					a.SetVar{"primary", "expectWord"}, a.SetVar{"secondary", ""},
				},
				unexpectedRune,
			)},
		},
		[]interface{}{
			e.Unexpected{""},
			// a.Panic{"unreachable", nil},
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
