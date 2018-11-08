package val

import (
	"unicode"

	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwstr"
	"github.com/baza-winner/bwcore/runeprovider"
)

//go:generate stringer -type=primaryState,secondaryState,itemType

var trace func(r rune, primary primaryState, secondary secondaryState, stack []stackItem)

func MustParse(s string) (result interface{}) {
	var err error
	if result, err = Parse(s); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return result
}

type primaryState uint8

const (
	begin primaryState = iota
	expectSpaceOrQwItemOrDelimiter
	expectSpaceOrMapKey
	expectEndOfQwItem
	expectContentOf
	expectWord
	expectEscapedContentOf
	expectRocket
	expectDigit
	end
)

type secondaryState uint8

const (
	none secondaryState = iota
	orArrayItemSeparator
	orMapKeySeparator
	orMapValueSeparator
	orUnderscoreOrDot
	orUnderscore
)

type itemType uint8

const (
	itemString itemType = iota
	itemQw
	itemQwItem
	itemNumber
	itemWord
	itemKey
	itemMap
	itemArray
)

// const trace = true

type stackItem struct {
	ps        runeprovider.PosStruct
	it        itemType
	s         string
	result    interface{}
	delimiter rune
}

// Parse превращает строку в interface{}-значение одной из следующих разновидностей (bwval.Kind()): Nil, Bool, Int, String, Map, Array
func Parse(s string, optVars ...map[string]interface{}) (result interface{}, err error) {
	var (
		primary         primaryState
		secondary       secondaryState
		needFinish      bool
		skipPostProcess bool
		isEOF           bool
		ok              bool
		r               rune
		r2              rune
		stack           []stackItem
	)
	p := runeprovider.ProxyFrom(runeprovider.FromString(s))
	for {
		_ = p.PullRune()
		r, isEOF, _ = p.Rune()
		needFinish = false
		switch primary {
		case end:
			switch {
			case isEOF:
				primary = end
				secondary = none
			case unicode.IsSpace(r):
			default:
				err = p.Unexpected(p.Curr)
				return
			}
		case expectRocket:
			if r == '>' {
				primary = begin
				secondary = none
			} else {
				err = p.Unexpected(p.Curr)
				return
			}
		case expectWord:
			if unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r) {
				stack[len(stack)-1].s += string(r)
			} else {
				_ = p.PushRune()
				needFinish = true
			}
		case expectSpaceOrQwItemOrDelimiter:
			switch {
			case isEOF:
				err = p.Unexpected(p.Curr)
				return
			case unicode.IsSpace(r):
			case r == stack[len(stack)-1].delimiter:
				needFinish = true
			default:
				stack = append(stack, stackItem{
					ps:        p.Curr,
					it:        itemQwItem,
					delimiter: stack[len(stack)-1].delimiter,
					s:         string(r),
				})
				primary = expectEndOfQwItem
				secondary = none
			}
		case expectEndOfQwItem:
			switch {
			case isEOF:
				err = p.Unexpected(p.Curr)
				return
			case unicode.IsSpace(r) || r == stack[len(stack)-1].delimiter:
				_ = p.PushRune()
				needFinish = true
			default:
				stack[len(stack)-1].s += string(r)
			}
		case expectContentOf:
			switch {
			case isEOF:
				err = p.Unexpected(p.Curr)
				return
			case r == stack[len(stack)-1].delimiter:
				needFinish = true
			case r == '\\':
				primary = expectEscapedContentOf
			default:
				stack[len(stack)-1].s += string(r)
			}
		case expectDigit:
			switch {
			case unicode.IsDigit(r) && secondary == none:
				stack[len(stack)-1].s += string(r)
				secondary = orUnderscoreOrDot
			case r == '.' && secondary == orUnderscoreOrDot:
				stack[len(stack)-1].s += string(r)
				secondary = orUnderscore
			case (r == '_' || unicode.IsDigit(r)) && (secondary == orUnderscoreOrDot || secondary == orUnderscore):
				stack[len(stack)-1].s += string(r)
			case secondary == none:
				err = p.Unexpected(p.Curr)
				return
			default:
				p.PushRune()
				needFinish = true
			}
		case expectSpaceOrMapKey:
			switch {
			case unicode.IsSpace(r):
			case unicode.IsLetter(r) || r == '_':
				stack = append(stack, stackItem{
					ps: p.Curr,
					it: itemKey,
					s:  string(r),
				})
				primary = expectWord
				secondary = none
			case r == '"' || r == '\'':
				stack = append(stack, stackItem{
					ps:        p.Curr,
					it:        itemKey,
					delimiter: r,
				})
				primary = expectContentOf
			case r == ',' && secondary == orMapValueSeparator:
				primary = expectSpaceOrMapKey
				secondary = none
			case r == stack[len(stack)-1].delimiter && stack[len(stack)-1].it == itemMap:
				needFinish = true
			}
		case expectEscapedContentOf:
			switch {
			case r == '"' || r == '\'' || r == '\\':
				stack[len(stack)-1].s += string(r)
			case stack[len(stack)-1].delimiter == '"':
				if r2, ok = escapeRunes[r]; ok {
					stack[len(stack)-1].s += string(r2)
				} else {
					err = p.Unexpected(p.Curr)
					return
				}
			}
			primary = expectContentOf
		case begin:
			switch {
			case isEOF && len(stack) == 0:
				primary = end
				secondary = none
			case isEOF:
				err = p.Unexpected(p.Curr)
				return
			case r == '=' && secondary == orMapKeySeparator:
				primary = expectRocket
				secondary = none
			case r == ':' && secondary == orMapKeySeparator:
				primary = begin
				secondary = none
			case r == ',' && secondary == orArrayItemSeparator:
				primary = begin
				secondary = none
			case unicode.IsSpace(r):
			case r == '{':
				stack = append(stack, stackItem{
					ps:        p.Curr,
					it:        itemMap,
					result:    map[string]interface{}{},
					delimiter: '}',
				})
				primary = expectSpaceOrMapKey
				secondary = none
			case r == '<':
				stack = append(stack, stackItem{
					ps:        p.Curr,
					it:        itemQw,
					result:    []interface{}{},
					delimiter: '>',
				})
				primary = expectSpaceOrQwItemOrDelimiter
				secondary = none
			case r == '[':
				stack = append(stack, stackItem{
					ps:        p.Curr,
					it:        itemArray,
					result:    []interface{}{},
					delimiter: ']',
				})
				primary = begin
				secondary = none
			case len(stack) > 0 && stack[len(stack)-1].it == itemArray && r == stack[len(stack)-1].delimiter:
				needFinish = true
			case r == '-' || r == '+':
				stack = append(stack, stackItem{
					ps: p.Curr,
					it: itemNumber,
					s:  string(r),
				})
				primary = expectDigit
				secondary = none
			case unicode.IsDigit(r):
				stack = append(stack, stackItem{
					ps: p.Curr,
					it: itemNumber,
					s:  string(r),
				})
				primary = expectDigit
				secondary = orUnderscoreOrDot
			case r == '"' || r == '\'':
				stack = append(stack, stackItem{
					ps:        p.Curr,
					it:        itemString,
					delimiter: r,
				})
				primary = expectContentOf
			case unicode.IsLetter(r) || r == '_':
				stack = append(stack, stackItem{
					ps: p.Curr,
					it: itemWord,
					s:  string(r),
				})
				primary = expectWord
				secondary = none
			default:
				err = p.Unexpected(p.Curr)
				return
			}
		default:
			bwerr.Unreachable()
		}
		if trace != nil {
			trace(r, primary, secondary, stack)
		}
		// if trace {
		// 	bwdebug.Print(
		// 		"r", string(r),
		// 		"primary", primary.String(),
		// 		"secondary", secondary.String(),
		// 		"stack", bwjson.Pretty(stack),
		// 	)
		// }
		if needFinish {
			skipPostProcess = false
			switch stack[len(stack)-1].it {
			case itemString, itemQwItem:
				stack[len(stack)-1].result = stack[len(stack)-1].s
			case itemNumber:
				stack[len(stack)-1].result, err = bwstr.ParseNumber(stack[len(stack)-1].s)
				if err != nil {
					err = p.Unexpected(stack[len(stack)-1].ps, err.Error())
					return
				}
			case itemWord:
				switch stack[len(stack)-1].s {
				case "true":
					stack[len(stack)-1].result = true
				case "false":
					stack[len(stack)-1].result = false
				case "nil", "null":
					stack[len(stack)-1].result = nil
				case "Bool", "String", "Int", "Number", "Map", "Array", "ArrayOf":
					stack[len(stack)-1].result = stack[len(stack)-1].s
				case "qw":
					_ = p.PullRune()
					r, isEOF, _ = p.Rune()
					if r2, ok = braces[r]; ok || unicode.IsPunct(r) || unicode.IsSymbol(r) {
						primary = expectSpaceOrQwItemOrDelimiter
						secondary = none
						stack[len(stack)-1].it = itemQw
						stack[len(stack)-1].result = []interface{}{}
						if ok {
							stack[len(stack)-1].delimiter = r2
						} else {
							stack[len(stack)-1].delimiter = r
						}
					} else {
						err = p.Unexpected(p.Curr)
						return
					}
					skipPostProcess = true
				}
			}

			// if trace {
			// 	bwdebug.Print(
			// 		"skipPostProcess", skipPostProcess,
			// 		"primary", primary.String(),
			// 		"secondary", secondary.String(),
			// 		"stack", bwjson.Pretty(stack),
			// 	)
			// }
			if !skipPostProcess {
				switch len(stack) {
				case 0:
				case 1:
					primary = end
				default:
					switch stack[len(stack)-2].it {
					case itemQw:
						arr2, _ := stack[len(stack)-2].result.([]interface{})
						stack[len(stack)-2].result = append(arr2, stack[len(stack)-1].result)
						primary = expectSpaceOrQwItemOrDelimiter
						secondary = none
					case itemArray:
						arr2, _ := stack[len(stack)-2].result.([]interface{})
						if stack[len(stack)-1].it == itemQw {
							arr, _ := stack[len(stack)-1].result.([]interface{})
							stack[len(stack)-2].result = append(arr2, arr...)
						} else {
							stack[len(stack)-2].result = append(arr2, stack[len(stack)-1].result)
						}
						primary = begin
						secondary = orArrayItemSeparator
					case itemMap:
						if stack[len(stack)-1].it == itemKey {
							stack[len(stack)-2].s = stack[len(stack)-1].s
							primary = begin
							secondary = orMapKeySeparator
						} else {
							m, _ := stack[len(stack)-2].result.(map[string]interface{})
							m[stack[len(stack)-2].s] = stack[len(stack)-1].result
							primary = expectSpaceOrMapKey
							secondary = orMapValueSeparator
						}
					default:
						bwerr.Unreachable()
					}
					stack = stack[:len(stack)-1]
				}
				// if trace {
				// 	bwdebug.Print(
				// 		"primary", primary.String(),
				// 		"secondary", secondary.String(),
				// 		"stack", bwjson.Pretty(stack),
				// 	)
				// }
			}
		}
		if p.Curr.IsEOF {
			break
		}
	}
	if len(stack) != 1 {
		err = p.Unexpected(p.Curr)
	} else {
		result = stack[0].result
	}
	return
}

var escapeRunes = map[rune]rune{
	'a': '\a',
	'b': '\b',
	'f': '\f',
	'n': '\n',
	'r': '\r',
	't': '\t',
	'v': '\v',
}
var braces = map[rune]rune{
	'(': ')',
	'{': '}',
	'<': '>',
	'[': ']',
}
