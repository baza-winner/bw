package defparse

import (
	"strconv"
	"unicode"

	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
)

func init() {
	pfaPrimaryStateMethodsCheck()
	pfaItemFinishMethodsCheck()
	pfaErrorValidatorsCheck()
}

type pfaStruct struct {
	stack        parseStack
	state        parseState
	result       interface{}
	source       string
	pos          int
	prevRunePtr  *rune
	runePtr      *rune
	nextRunePtr  *rune
	runeProvider PfaRuneProvider
}

func (pfa pfaStruct) GetDataForJson() interface{} {
	result := map[string]interface{}{}
	result["stack"] = pfa.stack.GetDataForJson()
	result["state"] = pfa.state.String()
	result["result"] = pfa.result
	result["pos"] = strconv.FormatInt(int64(pfa.pos), 10)
	if pfa.runePtr == nil {
		result["runePtr"] = nil
	} else {
		result["rune"] = string(*pfa.runePtr)
	}
	return result
}

func (pfa pfaStruct) String() string {
	return bwjson.PrettyJsonOf(pfa)
}

type PfaRuneProvider interface {
	PullRune() *rune
}

func pfaRun(runeProvider PfaRuneProvider, source string) (interface{}, error) {
	pfa := pfaStruct{stack: parseStack{}, state: parseState{primary: expectValueOrSpace}, pos: -1, runeProvider: runeProvider, source: source}
	var err error
	for {
		pfa.pullRune()
		var needFinishTopStackItem bool
		if needFinishTopStackItem, err = pfaPrimaryStateMethods[pfa.state.primary](&pfa); err == nil && needFinishTopStackItem {
			err = pfa.finishTopStackItem()
		}
		if err != nil {
			break
		}
		if pfa.runePtr == nil {
			if pfa.state.primary != expectEOF {
				pfa.panic("pfa.state.primary != expectEOF")
			}
			break
		}
	}
	return pfa.result, err
}

func (pfa *pfaStruct) pullRune() {
	pfa.prevRunePtr = pfa.runePtr
	if pfa.nextRunePtr != nil {
		pfa.runePtr = pfa.nextRunePtr
		pfa.nextRunePtr = nil
	} else {
		pfa.runePtr = pfa.runeProvider.PullRune()
	}
	pfa.pos += 1
}

func (pfa *pfaStruct) pushRune() {
	if pfa.prevRunePtr == nil {
		pfa.panic("pfa.prevRunePtr == nil")
	} else {
		pfa.nextRunePtr = pfa.runePtr
		pfa.runePtr = pfa.prevRunePtr
		pfa.pos -= 1
	}
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
	bwerror.Panicd(1, fmtString, fmtArgs...)
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

func (pfa *pfaStruct) finishTopStackItem() (err error) {
	stackItem := pfa.getTopStackItem()
	var skipPostProcess bool
	if skipPostProcess, err = pfaItemFinishMethods[stackItem.itemType](pfa); err == nil && !skipPostProcess {
		if len(pfa.stack) == 1 {
			pfa.result = stackItem.value
			pfa.state.setSecondary(expectEOF, orSpace)
		} else if len(pfa.stack) > 1 {
			if pfa.runePtr == nil {
				err = pfaErrorMake(pfa, unexpectedRuneError)
			} else {
				stackSubItem := pfa.popStackItem()
				stackItem = pfa.getTopStackItem()
				switch stackItem.itemType {
				case parseStackItemQw:
					stackItem.itemArray = append(stackItem.itemArray, stackSubItem.value)
					pfa.state.setPrimary(expectSpaceOrQwItemOrDelimiter)

				case parseStackItemArray:
					if stackSubItem.itemType == parseStackItemQw {
						stackItem.itemArray = append(stackItem.itemArray, stackSubItem.itemArray...)
					} else {
						stackItem.itemArray = append(stackItem.itemArray, stackSubItem.value)
					}
					pfa.state.setSecondary(expectValueOrSpace, orArrayItemSeparator)

				case parseStackItemMap:
					switch stackSubItem.itemType {
					case parseStackItemKey:
						stackItem.currentKey = stackSubItem.itemString
						switch {
						case unicode.IsSpace(*pfa.runePtr):
							pfa.state.setSecondary(expectValueOrSpace, orMapKeySeparator)
						case *pfa.runePtr == ':' || *pfa.runePtr == '=':
							pfa.state.setPrimary(expectValueOrSpace)
						default:
							pfa.state.setPrimary(expectMapKeySeparatorOrSpace)
						}
					default:
						stackItem.itemMap[stackItem.currentKey] = stackSubItem.value
						pfa.state.setSecondary(expectSpaceOrMapKey, orMapValueSeparator)
					}
				default:
					pfa.panic()
				}
			}
		}
	}
	return
}
