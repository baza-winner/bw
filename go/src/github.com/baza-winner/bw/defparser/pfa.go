package defparser

import (
	"encoding/json"
	"github.com/baza-winner/bw/core"
	"regexp"
	"unicode"
)

type pfaStruct struct {
	stack                  parseStack
	state                  parseState
	result                 interface{}
	needFinishTopStackItem bool
	pos                    int
	charPtr                *rune
}

func (pfa *pfaStruct) getDataForJson() interface{} {
	result := map[string]interface{}{}
	result["stack"] = pfa.stack.getDataForJson()
	result["state"] = pfa.state.String()
	result["result"] = pfa.result
	result["needFinishTopStackItem"] = pfa.needFinishTopStackItem
	result["pos"] = string(pfa.pos)
	if pfa.charPtr == nil {
		result["charPtr"] = nil
	} else {
		result["char"] = string(*pfa.charPtr)
	}
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

func (pfa *pfaStruct) getTopStackItemOfType(itemType parseStackItemType) (stackItem *parseStackItem) {
	if !(len(pfa.stack) >= 1 && pfa.stack[len(pfa.stack)-1].itemType == itemType) {
		pfa.panic()
	}
	stackItem = &pfa.stack[len(pfa.stack)-1]
	return
}

func (pfa *pfaStruct) panic() {
	core.Panicd(1, "<ansiOutline>pfa<ansi> <ansiSecondaryLiteral>%s<ansi>", pfa)
}

func (pfa *pfaStruct) getTopStackItem() (stackItem *parseStackItem) {
	if !(len(pfa.stack) >= 1) {
		pfa.panic()
	}
	stackItem = &pfa.stack[len(pfa.stack)-1]
	return
}

func (pfa *pfaStruct) popStackItem() (stackItem parseStackItem) {
	if !(len(pfa.stack) >= 1) {
		pfa.panic()
	}
	stackItem = pfa.stack[len(pfa.stack)-1]
	pfa.stack = pfa.stack[:len(pfa.stack)-1]
	return
}

func (pfa *pfaStruct) processCharAtPos() (err error) {
	pfa.needFinishTopStackItem = false
	if method, ok := pfaPrimaryStateMethods[pfa.state.primary]; ok {
		err = method(pfa)
	} else {
		pfa.panic()
	}
	if pfa.needFinishTopStackItem {
		if err = pfa.finishTopStackItem(); err != nil {
			return
		}
	}

	if pfa.charPtr == nil && pfa.state.primary != expectEOF {
		pfa.panic()
	}

	return
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

var underscoreRegexp = regexp.MustCompile("[_]+")

func (pfa *pfaStruct) finishTopStackItem() (err error) {
	if len(pfa.stack) < 1 {
		core.Panic("<ansiOutline>pfa<ansi> <ansiSecondaryLiteral>%s<ansi>", pfa)
	}
	stackItem := &pfa.stack[len(pfa.stack)-1]
	if method, ok := pfaItemFinishMethods[stackItem.itemType]; ok {
		var skipPostProcess bool
		if skipPostProcess, err = method(pfa); err != nil || skipPostProcess {
			return err
		}
	} else {
		core.Panic("<ansiOutline>pfa<ansi> <ansiSecondaryLiteral>%s<ansi>", pfa)
	}

	if len(pfa.stack) == 1 {
		pfa.result = stackItem.value
		pfa.state.setSecondary(expectEOF, orSpace)
		return
	} else if len(pfa.stack) > 1 && pfa.charPtr != nil {
		stackSubItem := pfa.popStackItem()
		stackItem = pfa.getTopStackItem()
		switch stackItem.itemType {
		case parseStackItemArray:
			stackItem.itemArray = append(stackItem.itemArray, stackSubItem.value)
			switch stackSubItem.itemType {
			case parseStackItemNumber, parseStackItemWord:
				switch {
				case unicode.IsSpace(*pfa.charPtr):
					pfa.state.setSecondary(expectValueOrSpace, orArrayItemSeparator)
				case *pfa.charPtr == ',':
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
					case unicode.IsSpace(*pfa.charPtr):
						pfa.state.setSecondary(expectSpaceOrMapKey, orMapValueSeparator)
					case *pfa.charPtr == ',':
						pfa.state.setPrimary(expectSpaceOrMapKey)
					default:
						return unexpectedCharError{}
					}
				default:
					pfa.state.setSecondary(expectSpaceOrMapKey, orMapValueSeparator)
				}
			}
		default:
			core.Panic("<ansiOutline>pfa<ansi> <ansiSecondaryLiteral>%s<ansi>", pfa)
		}
	}

	return
}
