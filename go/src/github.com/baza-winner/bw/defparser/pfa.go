package defparser

import (
	"encoding/json"
	"fmt"
	"github.com/baza-winner/bw/core"
	"regexp"
	// "strconv"
	// "strings"
	"unicode"
)

type pfaStruct struct {
	stack                  parseStack
	state                  parseState
	result                 interface{}
	needFinishTopStackItem bool
}

func (pfa *pfaStruct) getDataForJson() interface{} {
	result := map[string]interface{}{}
	result["stack"] = pfa.stack.getDataForJson()
	result["state"] = pfa.state.String()
	return result
}

func (pfa *pfaStruct) String() (result string) {
	bytes, _ := json.MarshalIndent(pfa.getDataForJson(), ``, `  `)
	result = string(bytes[:]) // https://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string/18615786#18615786
	return
}

type unexpectedCharError struct{}

func (e unexpectedCharError) Error() string {
	return `unexpected char`
}

func getPosTitle(pos int) (posTitle string) {
	if pos < 0 {
		posTitle = "end of source"
	} else {
		posTitle = fmt.Sprintf("<ansiOutline>pos <ansiSecondaryLiteral>%d<ansi>", pos)
	}
	return
}

func (pfa *pfaStruct) getTopStackItem(itemType parseStackItemType, pos int) (stackItem *parseStackItem) {
	if !(len(pfa.stack) >= 1 && pfa.stack[len(pfa.stack)-1].itemType == itemType) {
		core.Panic("<ansiOutline>stack<ansi> (<ansiSecondaryLiteral>%+v<ansi>) expects to have top item of type <ansiPrimaryLiteral>%s<ansi> while at "+getPosTitle(pos)+" and <ansiOutline>state <ansiSecondaryLiteral>%s", pfa.stack, itemType, pfa.state)
	}
	stackItem = &pfa.stack[len(pfa.stack)-1]
	return
}

func (pfa *pfaStruct) processCharAtPos(pos int, charPtr *rune) (err error) {
	pfa.needFinishTopStackItem = false
	if method, ok := pfaPrimaryStateMethods[pfa.state.primary]; ok {
		err = method(pfa, charPtr, pos)
	} else {
		core.Panic("<ansiOutline>pfa<ansi> <ansiSecondaryLiteral>%s<ansi>", pfa)
	}
	if pfa.needFinishTopStackItem {
		if err = pfa.finishTopStackItem(charPtr); err != nil {
			return
		}
	}

	if charPtr == nil && pfa.state.primary != expectEOF {
		core.Panic("<ansiOutline>pfa<ansi> <ansiSecondaryLiteral>%s<ansi>", pfa)
	}

	return
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

var underscoreRegexp = regexp.MustCompile("[_]+")

func (pfa *pfaStruct) finishTopStackItem(charPtr *rune) (err error) {
	if len(pfa.stack) < 1 {
		core.Panic("<ansiOutline>pfa<ansi> <ansiSecondaryLiteral>%s<ansi>", pfa)
	}
	stackItem := &pfa.stack[len(pfa.stack)-1]
	if method, ok := pfaItemFinishMethods[stackItem.itemType]; ok {
		var skipPostProcess bool
		if skipPostProcess, err = method(pfa, stackItem, charPtr); err != nil || skipPostProcess {
			return err
		}
	} else {
		core.Panic("<ansiOutline>pfa<ansi> <ansiSecondaryLiteral>%s<ansi>", pfa)
	}

	// switch stackItem.itemType {

	// case parseStackItemQwItem:
	// 	var skipPostProcess bool
	// 	if skipPostProcess, err = pfaItemFinishMethods[stackItem.itemType](pfa, charPtr); err != nil || skipPostProcess {
	// 		return err
	// 	}
	// if len(pfa.stack) < 2 {
	// 	core.Panic("len(pfa.stack) < 2")
	// }
	// stackSubItem := pfa.stack[len(pfa.stack)-1]
	// pfa.stack = pfa.stack[:len(pfa.stack)-1]
	// stackItem = pfa.getTopStackItem(parseStackItemQw, stackSubItem.pos)
	// stackItem.itemArray = append(stackItem.itemArray, stackSubItem.itemString)
	// if charPtr == nil {
	// 	core.Panic("charPtr == nil")
	// }
	// if unicode.IsSpace(*charPtr) {
	// 	pfa.state.setPrimary(expectSpaceOrQwItemOrDelimiter)
	// } else {
	// 	if len(pfa.stack) < 2 {
	// 		core.Panic("len(pfa.stack) < 2")
	// 	}
	// 	stackSubItem := pfa.stack[len(pfa.stack)-1]
	// 	pfa.stack = pfa.stack[:len(pfa.stack)-1]
	// 	stackItem = pfa.getTopStackItem(parseStackItemArray, stackSubItem.pos)
	// 	stackItem.itemArray = append(stackItem.itemArray, stackSubItem.itemArray...)
	// 	pfa.state.setPrimary(expectArrayItemSeparatorOrSpace)
	// }
	// return

	// case parseStackItemNumber:
	// 	source := underscoreRegexp.ReplaceAllLiteralString(stackItem.itemString, ``)
	// 	if strings.Contains(source, `.`) {
	// 		var float64Val float64
	// 		if float64Val, err = strconv.ParseFloat(source, 64); err == nil {
	// 			stackItem.value = float64Val
	// 		}
	// 	} else {
	// 		var int64Val int64
	// 		if int64Val, err = strconv.ParseInt(source, 10, 64); err == nil {
	// 			if int64(MinInt) <= int64Val && int64Val <= int64(MaxInt) {
	// 				stackItem.value = int(int64Val)
	// 			} else {
	// 				stackItem.value = int64Val
	// 			}
	// 		}
	// 	}
	// 	if err != nil {
	// 		err = core.Error("failed to get number from string <ansiPrimaryLiteral>%s<ansi> at "+getPosTitle(stackItem.pos)+": %v", stackItem.itemString, err)
	// 	}

	// case parseStackItemString:
	// 	stackItem.value = stackItem.itemString

	// case parseStackItemWord:
	// 	switch stackItem.itemString {
	// 	case "true":
	// 		stackItem.value = true
	// 	case "false":
	// 		stackItem.value = false
	// 	case "nil":
	// 		stackItem.value = nil
	// 	case "qw":
	// 		if len(pfa.stack) >= 2 && pfa.stack[len(pfa.stack)-2].itemType == parseStackItemArray && charPtr != nil {
	// 			pfa.state.setPrimary(expectSpaceOrQwItemOrDelimiter)
	// 			switch *charPtr {
	// 			case '<':
	// 				stackItem.delimiter = '>'
	// 			case '[':
	// 				stackItem.delimiter = ']'
	// 			case '(':
	// 				stackItem.delimiter = ')'
	// 			case '{':
	// 				stackItem.delimiter = '}'
	// 			default:
	// 				switch {
	// 				case unicode.IsPunct(*charPtr) || unicode.IsSymbol(*charPtr):
	// 					stackItem.delimiter = *charPtr
	// 				default:
	// 					return unexpectedCharError{}
	// 				}
	// 			}
	// 			if pfa.state.primary == expectSpaceOrQwItemOrDelimiter {
	// 				stackItem.itemType = parseStackItemQw
	// 				stackItem.itemArray = []interface{}{}
	// 			}
	// 		} else {
	// 			err = core.Error("unexpected word <ansiPrimaryLiteral>%s<ansi> in non array context at "+getPosTitle(stackItem.pos), stackItem.itemString)
	// 		}

	// 		return
	// 	default:
	// 		err = core.Error("unexpected word <ansiPrimaryLiteral>%s<ansi> at "+getPosTitle(stackItem.pos)+" <ansiOutline>pfa.stack <ansiSecondaryLiteral>%s", stackItem.itemString, pfa.stack)
	// 	}

	// case parseStackItemArray:
	// 	stackItem.value = stackItem.itemArray

	// case parseStackItemMap:
	// 	stackItem.value = stackItem.itemMap

	// case parseStackItemKey:

	// default:
	// 	err = core.Error("can not finish item of type <ansiPrimaryLiteral>%s<ansi> at "+getPosTitle(stackItem.pos), stackItem.itemType)
	// }
	// if err != nil {
	// 	return
	// }

	if len(pfa.stack) == 1 {
		pfa.result = stackItem.value
		pfa.state.setSecondary(expectEOF, orSpace)
		return
	} else if len(pfa.stack) > 1 && charPtr != nil {
		var stackSubItem parseStackItem
		stackSubItem, pfa.stack = pfa.stack[len(pfa.stack)-1], pfa.stack[:len(pfa.stack)-1] // https://github.com/golang/go/wiki/SliceTricks
		stackItem = &pfa.stack[len(pfa.stack)-1]
		switch stackItem.itemType {
		case parseStackItemArray:
			stackItem.itemArray = append(stackItem.itemArray, stackSubItem.value)
			switch stackSubItem.itemType {
			case parseStackItemNumber, parseStackItemWord:
				switch {
				case unicode.IsSpace(*charPtr):
					pfa.state.setSecondary(expectValueOrSpace, orArrayItemSeparator)
				case *charPtr == ',':
					pfa.state.setPrimary(expectValueOrSpace)
				default:
					return unexpectedCharError{}
				}
			default:
				pfa.state.setSecondary(expectValueOrSpace, orArrayItemSeparator)
			}

		case parseStackItemMap:
			switch stackSubItem.itemType {
			case parseStackItemKey:
				stackItem.currentKey = stackSubItem.itemString
				pfa.state.setPrimary(expectMapKeySeparatorOrSpace)
			default:
				stackItem.itemMap[stackItem.currentKey] = stackSubItem.value
				switch stackSubItem.itemType {
				case parseStackItemNumber, parseStackItemWord:
					switch {
					case unicode.IsSpace(*charPtr):
						pfa.state.setSecondary(expectSpaceOrMapKey, orMapValueSeparator)
					case *charPtr == ',':
						pfa.state.setPrimary(expectSpaceOrMapKey)
					default:
						return unexpectedCharError{}
					}
				default:
					pfa.state.setSecondary(expectSpaceOrMapKey, orMapValueSeparator)
				}
			}
		default:
			err = core.Error("<ansiOutline>stackItem <ansiSecondaryLiteral>%s<ansi> can not have subitem <ansiSecondaryLiteral>%s<ansiOutline>pfa.stack<ansiSecondaryLiteral>%s", stackItem, stackSubItem, pfa.stack)
		}
	}

	return
}
