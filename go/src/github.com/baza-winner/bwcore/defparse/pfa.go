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

type runePtrStruct struct {
	runePtr *rune
	line    uint
	col     uint
	pos     int
}

func (v runePtrStruct) copyPtr() *runePtrStruct {
	return &runePtrStruct{v.runePtr, v.line, v.col, v.pos}
}

func (v runePtrStruct) GetDataForJson() interface{} {
	result := map[string]interface{}{}
	if v.runePtr == nil {
		result["rune"] = "EOF"
	} else {
		result["rune"] = string(*(v.runePtr))
	}
	result["line"] = v.line
	result["col"] = v.col
	result["pos"] = v.pos
	return result
}

type pfaStruct struct {
	stack        parseStack
	state        parseState
	result       interface{}
	source       string
	prev         *runePtrStruct
	curr         runePtrStruct
	next         *runePtrStruct
	runeProvider PfaRuneProvider
}

func (pfa pfaStruct) GetDataForJson() interface{} {
	result := map[string]interface{}{}
	result["stack"] = pfa.stack.GetDataForJson()
	result["state"] = pfa.state.String()
	result["result"] = pfa.result
	result["pos"] = strconv.FormatInt(int64(pfa.curr.pos), 10)
	result["curr"] = pfa.curr.GetDataForJson()
	if pfa.prev != nil {
		result["prev"] = pfa.prev.GetDataForJson()
	}
	if pfa.next != nil {
		result["next"] = pfa.prev.GetDataForJson()
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
	pfa := pfaStruct{stack: parseStack{}, state: parseState{primary: expectValueOrSpace}, runeProvider: runeProvider, curr: runePtrStruct{pos: -1, line: 1}, source: source}
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
		if pfa.curr.runePtr == nil {
			if pfa.state.primary != expectEOF {
				pfa.panic("pfa.state.primary != expectEOF")
			}
			break
		}
	}
	return pfa.result, err
}

func (pfa *pfaStruct) pullRune() {
	pfa.prev = pfa.curr.copyPtr()
	if pfa.next == nil {
		line := pfa.prev.line
		col := pfa.prev.col
		if pfa.prev.runePtr != nil && *(pfa.prev.runePtr) == '\n' {
			line += 1
			col = 1
		} else {
			col += 1
		}
		pfa.curr = runePtrStruct{
			runePtr: pfa.runeProvider.PullRune(),
			pos:     pfa.prev.pos + 1,
			line:    line,
			col:     col,
		}
	} else {
		pfa.curr = *(pfa.next)
		pfa.next = nil
	}
}

func (pfa *pfaStruct) pushRune() {
	if pfa.prev == nil {
		pfa.panic("pfa.prev == nil")
	} else {
		pfa.next = pfa.curr.copyPtr()
		pfa.curr = *(pfa.prev)
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
			if pfa.curr.runePtr == nil {
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
						case unicode.IsSpace(*pfa.curr.runePtr):
							pfa.state.setSecondary(expectValueOrSpace, orMapKeySeparator)
						case *pfa.curr.runePtr == ':' || *pfa.curr.runePtr == '=':
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
