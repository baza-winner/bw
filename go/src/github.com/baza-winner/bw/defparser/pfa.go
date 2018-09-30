package defparser

import (
	"encoding/json"
	"github.com/baza-winner/bw/core"
	"unicode"
)

func init() {
	expect := _expectBelow + 1
	for expect < _expectAbove {
		if _, ok := pfaPrimaryStateMethods[expect]; !ok {
			panic(expect)
		}
		expect += 1
	}

	itemType := _parseStackItemBelow + 1
	for itemType < _parseStackItemAbove {
		if _, ok := pfaItemFinishMethods[itemType]; !ok {
			panic(itemType)
		}
		itemType += 1
	}
}

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
	return `unexpectedCharError`
}

type failedToGetNumberError struct{}

func (e failedToGetNumberError) Error() string {
	return `failedToGetNumberError`
}

type unexpectedWordError struct{}

func (e unexpectedWordError) Error() string {
	return `unexpectedWordError`
}

func (pfa *pfaStruct) panic(args ...interface{}) {
	fmtString := "<ansiOutline>pfa<ansi> <ansiSecondaryLiteral>%s<ansi>"
	if args != nil {
		fmtString += " " + args[0].(string)
	}
	fmtArgs := []interface{}{pfa}
	if args != nil && len(args) > 1 {
		fmtArgs = append(fmtArgs, args[1:])
	}
	core.Panicd(1, fmtString, fmtArgs)
}

func (pfa *pfaStruct) getTopStackItemOfType(itemType parseStackItemType, ofsList ...int) (stackItem *parseStackItem) {
	stackItem = pfa.getTopStackItem(ofsList...)
	if stackItem.itemType != itemType {
		pfa.panic("<ansiOutline>itemType<ansiSecondaryLiteral", itemType)
	}
	return
}

func (pfa *pfaStruct) getTopStackItem(ofsList ...int) (stackItem *parseStackItem) {
	ofs := -1
	if ofsList != nil && ofsList[0] < 0 {
		ofs = ofsList[0]
	}
	if len(pfa.stack) < -ofs {
		pfa.panic()
	}
	stackItem = &pfa.stack[len(pfa.stack)+ofs]
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
	pfaPrimaryStateMethods[pfa.state.primary](pfa)
	if pfa.needFinishTopStackItem {
		err = pfa.finishTopStackItem()
	}
	if err == nil && pfa.charPtr == nil && pfa.state.primary != expectEOF {
		pfa.panic()
	}
	return
}

func (pfa *pfaStruct) finishTopStackItem() (err error) {
	stackItem := pfa.getTopStackItem()
	var skipPostProcess bool
	if skipPostProcess, err = pfaItemFinishMethods[stackItem.itemType](pfa); err == nil && !skipPostProcess {
		if len(pfa.stack) == 1 {
			pfa.result = stackItem.value
			pfa.state.setSecondary(expectEOF, orSpace)
		} else if len(pfa.stack) > 1 {
			if pfa.charPtr == nil {
				err = unexpectedCharError{}
			} else {
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
							err = unexpectedCharError{}
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
								err = unexpectedCharError{}
							}
						default:
							pfa.state.setSecondary(expectSpaceOrMapKey, orMapValueSeparator)
						}
					}
				default:
					pfa.panic()
				}
			}
		}
	}

	return
}
