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
	p := runeprovider.ProxyFrom(runeprovider.FromString(s))
	if result, err = RunPfa(nil, p, Begin, None); err != nil {
		return
	}
	var (
		isEOF bool
		r     rune
	)
	r, isEOF, _ = p.Rune()
	for !isEOF {
		if !unicode.IsSpace(r) {
			err = p.Unexpected(p.Curr)
			return
		}
		_ = p.PullRune()
		r, isEOF, _ = p.Rune()
	}
	return
}

// ============================================================================

func RunPfa(stack []StackItem, p *runeprovider.Proxy, primary PrimaryState, secondary SecondaryState) (result interface{}, err error) {
	var (
		needFinish      bool
		skipPostProcess bool
		isEOF           bool
		ok              bool
		r               rune
		r2              rune
	)
	for primary != End {
		_ = p.PullRune()
		r, isEOF, _ = p.Rune()
		needFinish = false
		switch primary {
		case ExpectRocket:
			if r == '>' {
				primary = Begin
				secondary = None
			} else {
				err = p.Unexpected(p.Curr)
				return
			}
		case ExpectWord:
			if unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r) {
				stack[len(stack)-1].S += string(r)
			} else {
				_ = p.PushRune()
				needFinish = true
			}
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
		case ExpectEndOfQwItem:
			switch {
			case unicode.IsSpace(r) || r == stack[len(stack)-1].Delimiter:
				_ = p.PushRune()
				needFinish = true
			case !isEOF:
				stack[len(stack)-1].S += string(r)
			}
		case ExpectSpaceOrMapKey:
			switch {
			case unicode.IsSpace(r):
			case unicode.IsLetter(r) || r == '_':
				stack = append(stack, StackItem{
					PosStruct: p.Curr,
					Kind:      ItemKey,
					S:         string(r),
				})
				primary = ExpectWord
				secondary = None
			case r == '"' || r == '\'':
				stack = append(stack, StackItem{
					PosStruct: p.Curr,
					Kind:      ItemKey,
					Delimiter: r,
				})
				if stack[len(stack)-1].S, err = bwparse.ParseString(p); err != nil {
					return
				}
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
				primary = ExpectRocket
				secondary = None
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
				stack = append(stack, StackItem{
					PosStruct: p.Curr,
					Kind:      ItemQw,
					Result:    []interface{}{},
					Delimiter: '>',
				})
				primary = ExpectSpaceOrQwItemOrDelimiter
				secondary = None
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

				stack = append(stack, StackItem{
					PosStruct: p.Curr,
					Kind:      ItemNumber,
				})
				if stack[len(stack)-1].Result, err = bwparse.ParseNumber(p); err != nil {
					return
				}
				needFinish = true

			case r == '"' || r == '\'':
				stack = append(stack, StackItem{
					PosStruct: p.Curr,
					Kind:      ItemString,
					Delimiter: r,
				})
				if stack[len(stack)-1].S, err = bwparse.ParseString(p); err != nil {
					return
				}
				needFinish = true
			case unicode.IsLetter(r) || r == '_':
				stack = append(stack, StackItem{
					PosStruct: p.Curr,
					Kind:      ItemWord,
					S:         string(r),
				})
				primary = ExpectWord
				secondary = None
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
					_ = p.PullRune()
					r, _, _ = p.Rune()
					if r2, ok = Braces[r]; ok || unicode.IsPunct(r) || unicode.IsSymbol(r) {
						primary = ExpectSpaceOrQwItemOrDelimiter
						secondary = None
						stack[len(stack)-1].Kind = ItemQw
						stack[len(stack)-1].Result = []interface{}{}
						if ok {
							stack[len(stack)-1].Delimiter = r2
						} else {
							stack[len(stack)-1].Delimiter = r
						}
					} else {
						err = p.Unexpected(p.Curr)
						return
					}
					skipPostProcess = true
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
					_ = p.PullRune()
					r, isEOF, _ = p.Rune()
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
	ExpectWord
	ExpectRocket
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
