package defparse

import (
	"encoding/json"
	"fmt"
	"github.com/baza-winner/bw/core"
	"strconv"
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
	stack   parseStack
	state   parseState
	result  interface{}
	pos     int
	charPtr *rune
}

func (pfa *pfaStruct) getDataForJson() interface{} {
	result := map[string]interface{}{}
	result["stack"] = pfa.stack.getDataForJson()
	result["state"] = pfa.state.String()
	result["result"] = pfa.result
	result["pos"] = strconv.FormatInt(int64(pfa.pos), 10)
	if pfa.charPtr == nil {
		result["charPtr"] = nil
	} else {
		result["char"] = string(*pfa.charPtr)
	}
	return result
}

func (pfa pfaStruct) String() (result string) {
	bytes, _ := json.MarshalIndent(pfa.getDataForJson(), ``, `  `)
	result = string(bytes[:]) // https://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string/18615786#18615786
	return
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
	fmt.Println(fmtString, fmtArgs)
	core.Panicd(1, fmtString, fmtArgs...)
}

func (pfa *pfaStruct) ifStackLen(minLen int) bool {
	if len(pfa.stack) < minLen {
		return false
	}
	return true
}

func (pfa *pfaStruct) mustStackLen(minLen int) {
	if !pfa.ifStackLen(minLen) {
		pfa.panic("<ansiOutline>minLen <ansiSecondaryLiteral>%d", minLen)
	}
}

func (pfa *pfaStruct) isTopStackItemOfType(itemType parseStackItemType, ofsList ...int) bool {
	ofs := -1
	if ofsList != nil && ofsList[0] < 0 {
		ofs = ofsList[0]
	}
	if pfa.ifStackLen(-ofs) && pfa.getTopStackItem().itemType == itemType {
		return true
	}
	return false
}

func (pfa *pfaStruct) getTopStackItemOfType(itemType parseStackItemType, ofsList ...int) (stackItem *parseStackItem) {
	stackItem = pfa.getTopStackItem(ofsList...)
	if stackItem.itemType != itemType {
		pfa.panic("<ansiOutline>itemType<ansiSecondaryLiteral>%s", itemType)
	}
	return
}

func (pfa *pfaStruct) getTopStackItem(ofsList ...int) (stackItem *parseStackItem) {
	ofs := -1
	if ofsList != nil && ofsList[0] < 0 {
		ofs = ofsList[0]
	}
	pfa.mustStackLen(-ofs)
	stackItem = &pfa.stack[len(pfa.stack)+ofs]
	return
}

func (pfa *pfaStruct) popStackItem() (stackItem parseStackItem) {
	pfa.mustStackLen(1)
	stackItem = pfa.stack[len(pfa.stack)-1]
	pfa.stack = pfa.stack[:len(pfa.stack)-1]
	return
}

func (pfa *pfaStruct) processCharAtPos(char rune, pos int) error {
	pfa.pos = pos
	pfa.charPtr = &char
	return pfa.doProcessCharAtPos()
}

func (pfa *pfaStruct) processEOF() (err error) {
	pfa.pos = -1
	pfa.charPtr = nil
	if err = pfa.doProcessCharAtPos(); err == nil && pfa.state.primary != expectEOF {
		pfa.panic()
	}
	return err
}

func (pfa *pfaStruct) doProcessCharAtPos() (err error) {
	var needFinishTopStackItem bool
	if needFinishTopStackItem, err = pfaPrimaryStateMethods[pfa.state.primary](pfa); err == nil && needFinishTopStackItem {
		err = pfa.finishTopStackItem()
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
						switch {
						case unicode.IsSpace(*pfa.charPtr):
							pfa.state.setSecondary(expectValueOrSpace, orMapKeySeparator)
						case *pfa.charPtr == ':' || *pfa.charPtr == '=':
							pfa.state.setPrimary(expectValueOrSpace)
						default:
							pfa.state.setPrimary(expectMapKeySeparatorOrSpace)
						}
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
