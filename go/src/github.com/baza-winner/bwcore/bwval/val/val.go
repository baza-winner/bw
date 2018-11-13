package val

import (
	"unicode"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwparse"
	"github.com/baza-winner/bwcore/runeprovider"
)

//go:generate stringer -type=PrimaryState,SecondaryState,ItemKind

// ============================================================================

// MustParse - must-обертка Parse()
func MustParse(s string, optVars ...map[string]interface{}) (result interface{}) {
	var err error
	if result, err = Parse(s, optVars...); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return result
}

// Parse превращает строку в interface{}-значение одной из следующих разновидностей (bwval.Kind()): Nil, Bool, Int, String, Map, Array
func Parse(s string, optVars ...map[string]interface{}) (result interface{}, err error) {
	defer func() {
		if err != nil {
			result = nil
		}
	}()
	var (
		isEOF bool
		r     rune
		ok    bool
	)
	p := runeprovider.ProxyFrom(runeprovider.FromString(s))

	if r, err = p.PullNonEOFRune(); err != nil {
		return
	}
	if result, err = ParseVal(p, r); err != nil {
		return
	}
	if r, isEOF, err = p.PullRuneOrEOF(); err != nil || isEOF {
		return
	}
	if _, ok, err = bwparse.ParseSpace(p, r); err != nil {
		return
	} else if !ok {
		err = p.Unexpected(p.Curr)
		return
	}
	return
}

// ============================================================================

func ParseVal(p *runeprovider.Proxy, r rune) (result interface{}, err error) {
	var (
		needFinish      bool
		skipPostProcess bool
		isEOF           bool
		ok              bool
		s               string
		start           runeprovider.PosStruct
		val             interface{}
		vals            []interface{}
		stack           []StackItem
		primary         PrimaryState
		secondary       SecondaryState
	)
	for primary != End {
		r, isEOF, _ = p.Rune()
		needFinish = false
		switch primary {
		case ExpectSpaceOrQwItemOrDelimiter:
			switch {
			case unicode.IsSpace(r):
			case r == stack[len(stack)-1].Delimiter:
				needFinish = true
			case !isEOF:
				stack = append(stack, StackItem{
					PosStruct: p.Curr,
					Kind:      ItemQwItem,
					Delimiter: stack[len(stack)-1].Delimiter,
					S:         string(r),
				})
				primary = ExpectEndOfQwItem
				secondary = None
			}
		case ExpectSpaceOrMapKey:
			switch {
			case unicode.IsSpace(r):
			case unicode.IsLetter(r) || r == '_':
				if s, start, _, err = bwparse.ParseWord(p, r); err != nil {
					return
				}
				stack = append(stack, StackItem{
					PosStruct: start,
					Kind:      ItemKey,
					S:         s,
				})
				needFinish = true
			case r == '"' || r == '\'':
				if s, start, _, err = bwparse.ParseString(p, r); err != nil {
					return
				}
				stack = append(stack, StackItem{
					PosStruct: start,
					Kind:      ItemKey,
					Delimiter: r,
					S:         s,
				})
				needFinish = true
			case r == ',' && secondary == orMapValueSeparator:
				primary = ExpectSpaceOrMapKey
				secondary = None
			case r == stack[len(stack)-1].Delimiter && stack[len(stack)-1].Kind == ItemMap:
				needFinish = true
			}
		case Begin:
			switch {
			case r == '=' && secondary == orMapKeySeparator:
				if r, isEOF, err = p.PullRuneOrEOF(); err != nil {
					return
				} else if r == '>' {
					primary = Begin
					secondary = None
				} else {
					err = p.Unexpected(p.Curr)
					return
				}
			case r == ':' && secondary == orMapKeySeparator:
				primary = Begin
				secondary = None
			case r == ',' && secondary == orArrayItemSeparator:
				primary = Begin
				secondary = None
			case unicode.IsSpace(r):
			case r == '{':
				stack = append(stack, StackItem{
					PosStruct: p.Curr,
					Kind:      ItemMap,
					Result:    map[string]interface{}{},
					Delimiter: '}',
				})
				primary = ExpectSpaceOrMapKey
				secondary = None
			case r == '<':
				if vals, start, _, err = bwparse.ParseQw(p, r); err != nil {
					return
				}
				stack = append(stack, StackItem{
					PosStruct: start,
					Kind:      ItemQw,
					Result:    vals,
					Delimiter: '>',
				})
				needFinish = true
			case r == '[':
				stack = append(stack, StackItem{
					PosStruct: p.Curr,
					Kind:      ItemArray,
					Result:    []interface{}{},
					Delimiter: ']',
				})
				primary = Begin
				secondary = None
			case len(stack) > 0 && stack[len(stack)-1].Kind == ItemArray && r == stack[len(stack)-1].Delimiter:
				needFinish = true
			case r == '-' || r == '+' || unicode.IsDigit(r):
				if val, start, _, err = bwparse.ParseNumber(p, r); err != nil {
					return
				}
				stack = append(stack, StackItem{
					PosStruct: start,
					Kind:      ItemNumber,
					Result:    val,
				})
				needFinish = true

			case r == '"' || r == '\'':
				if s, start, _, err = bwparse.ParseString(p, r); err != nil {
					return
				}
				stack = append(stack, StackItem{
					PosStruct: start,
					Kind:      ItemString,
					Delimiter: r,
					S:         s,
				})
				needFinish = true
			case unicode.IsLetter(r) || r == '_':
				if s, start, _, err = bwparse.ParseWord(p, r); err != nil {
					return
				}
				needFinish = true

				stack = append(stack, StackItem{
					PosStruct: start,
					Kind:      ItemWord,
					S:         s,
				})
			default:
				err = p.Unexpected(p.Curr)
				return
			}
		default:
			bwerr.Unreachable()
		}
		if trace != nil {
			trace(r, primary, secondary, stack, "needFinish", needFinish)
		}
		if needFinish {
			skipPostProcess = false
			switch stack[len(stack)-1].Kind {
			case ItemString, ItemQwItem:
				stack[len(stack)-1].Result = stack[len(stack)-1].S
			case ItemNumber:
			case ItemWord:
				switch stack[len(stack)-1].S {
				case "true":
					stack[len(stack)-1].Result = true
				case "false":
					stack[len(stack)-1].Result = false
				case "nil", "null":
					stack[len(stack)-1].Result = nil
				case "Bool", "String", "Int", "Number", "Map", "Array", "ArrayOf":
					stack[len(stack)-1].Result = stack[len(stack)-1].S
				case "qw":
					if r, err = p.PullNonEOFRune(); err != nil {
						return
					}
					if vals, start, ok, err = bwparse.ParseQw(p, r); err != nil {
						return
					}
					if !ok {
						err = p.Unexpected(p.Curr)
						return
					}
					stack[len(stack)-1].Kind = ItemQw
					stack[len(stack)-1].Result = vals
				default:
					err = p.Unexpected(stack[len(stack)-1].PosStruct, bw.Fmt(ansi.String("unexpected <ansiErr>%q<ansi>"), stack[len(stack)-1].S))
					return
				}
			}

			if trace != nil {
				trace(r, primary, secondary, stack, "skipPostProcess", skipPostProcess)
			}
			if !skipPostProcess {
				switch len(stack) {
				case 0:
				case 1:
					primary = End
				default:
					switch stack[len(stack)-2].Kind {
					case ItemQw:
						arr2, _ := stack[len(stack)-2].Result.([]interface{})
						stack[len(stack)-2].Result = append(arr2, stack[len(stack)-1].Result)
						primary = ExpectSpaceOrQwItemOrDelimiter
						secondary = None
					case ItemArray:
						arr2, _ := stack[len(stack)-2].Result.([]interface{})
						if stack[len(stack)-1].Kind == ItemQw {
							arr, _ := stack[len(stack)-1].Result.([]interface{})
							stack[len(stack)-2].Result = append(arr2, arr...)
						} else {
							stack[len(stack)-2].Result = append(arr2, stack[len(stack)-1].Result)
						}
						primary = Begin
						secondary = orArrayItemSeparator
					case ItemMap:
						if stack[len(stack)-1].Kind == ItemKey {
							stack[len(stack)-2].S = stack[len(stack)-1].S
							primary = Begin
							secondary = orMapKeySeparator
						} else {
							m, _ := stack[len(stack)-2].Result.(map[string]interface{})
							m[stack[len(stack)-2].S] = stack[len(stack)-1].Result
							primary = ExpectSpaceOrMapKey
							secondary = orMapValueSeparator
						}
					default:
						bwerr.Unreachable()
					}
					stack = stack[:len(stack)-1]
				}
				if trace != nil {
					trace(r, primary, secondary, stack, "-", false)
				}
			}
		}
		if primary != End && p.Curr.IsEOF {
			err = p.Unexpected(p.Curr)
			return
		}

		_ = p.PullRune()
	}
	result = stack[0].Result
	return
}

//go:generate stringer -type=PrimaryState,SecondaryState,ItemKind

var trace func(r rune, primary PrimaryState, secondary SecondaryState, stack []StackItem, boolVarName string, boolVarVal bool)

type PrimaryState uint8

const (
	Begin PrimaryState = iota
	ExpectSpaceOrQwItemOrDelimiter
	ExpectSpaceOrMapKey
	ExpectEndOfQwItem
	// ExpectWord
	// ExpectRocket
	End
)

type SecondaryState uint8

const (
	None SecondaryState = iota
	orArrayItemSeparator
	orMapKeySeparator
	orMapValueSeparator
)

type ItemKind uint8

const (
	ItemString ItemKind = iota
	ItemQw
	ItemQwItem
	ItemNumber
	ItemWord
	ItemKey
	ItemMap
	ItemArray
)

type StackItem struct {
	PosStruct runeprovider.PosStruct
	Kind      ItemKind
	S         string
	Result    interface{}
	Delimiter rune
}

var Braces = map[rune]rune{
	'(': ')',
	'{': '}',
	'<': '>',
	'[': ']',
}
