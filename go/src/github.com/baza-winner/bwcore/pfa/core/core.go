package core

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/runeprovider"
)

// ============================================================================

// type PfaErrorProvider interface {
// 	PfaError(pfa *PfaStruct) error
// }

// ============================================================================

type ProcessorAction interface {
	Execute(pfa *PfaStruct) (err error)
}

// ============================================================================

type ProccessorActionProvider interface {
	GetAction() ProcessorAction
}

// ============================================================================

// type ValProvider interface {
// 	GetVal(pfa *PfaStruct) interface{}
// 	GetSource(pfa *PfaStruct) formatted.String
// }

// ============================================================================

type ValChecker interface {
	Conforms(pfa *PfaStruct, val interface{}, varPath VarPath) (bool, error)
}

// ============================================================================

type ValCheckerProvider interface {
	GetChecker() ValChecker
}

// ============================================================================

type TraceLevel uint8

const (
	TraceNone TraceLevel = iota
	TraceBrief
	TraceAll
)

// ============================================================================

// ============================================================================

type PfaStruct struct {
	// Stack ParseStack
	Proxy *runeprovider.Proxy
	// Err          ErrorProvider
	Err             error
	Vars            map[string]interface{}
	TraceLevel      TraceLevel
	traceConditions []string
	ruleLevel       int
}

func PfaFrom(p runeprovider.RuneProvider, TraceLevel TraceLevel) *PfaStruct {
	return &PfaStruct{
		// Stack:      ParseStack{},
		Proxy:      runeprovider.ProxyFrom(p),
		Vars:       map[string]interface{}{},
		TraceLevel: TraceLevel,
	}
}

// func (pfa *PfaStruct) Value(val interface{}) VarValue {
// 	return VarValue{val, pfa}
// }

func (pfa *PfaStruct) indent(indentLevel int) string {
	indentAtom := "  "
	indent := ""
	for i := 0; i <= indentLevel; i++ {
		indent += indentAtom
	}
	return indent
}

// type formatted.String string

// const fmtPanic = "<ansiVar>pfa<ansi> <ansiVal>%s<ansi>"

var ansiPanic string

func init() {
	ansiPanic = ansi.String("<ansiVar>pfa<ansi> <ansiVal>%s<ansi> {Error}")
}

// func (pfa *PfaStruct) Panic(optA ...bw.I) {
// 	fmtString := fmtPanic
// 	fmtArgs := []interface{}{pfa}
// 	if optFmtStruct == nil {
// 		bwerr.Panicd(1, fmtString, fmtArgs...)
// 	} else {
// 		fmtString += " " + optFmtStruct[0].Fmt
// 		fmtArgs = append(fmtArgs, optFmtStruct[0].Args...)
// 	}
// 	bwerr.Panicd(1, fmtString, fmtArgs...)
// 	// bwerr.PanicErr(fmt.Errorf(err.Error()+"\n"+ansi.Ansi("", fmtString), fmtArgs), 1)
// 	// pfa.panicHelper(fmtString, fmtArgs)
// }

func (pfa *PfaStruct) Panic() {
	pfa.PanicA(bwerr.A{Depth: 1})
}

func (pfa *PfaStruct) PanicA(a bw.I) {
	bwerr.PanicA(bwerr.E{Error: bwerr.FromA(bwerr.IncDepth(a)).Refine(ansiPanic)})
	// bwerr.ModifyBy(bwerr.IncDepth(a), bw.A{fmtPanic, bw.Args(pfa)}, true)
}

// func (pfa *PfaStruct) PanicErr(err error) {
// 	// fmtString := fmtPanic
// 	// fmtString := "<ansiVar>pfa<ansi> <ansiVal>%s<ansi>"
// 	// fmtArgs := []interface{}{pfa}
// 	bwerr.PanicA(fmt.Errorf(err.Error()+"\n"+ansi.Ansi("", fmtPanic), pfa), 1)
// 	// bwerr.PanicErr(fmt.Errorf(err.Error()+"\n"+ansi.Ansi("", fmtString), fmtArgs), 1)
// 	// pfa.panicHelper(err.Error()+"\n"+ansi.Ansi("", fmtString),, fmtArgs)
// }

// func (pfa *PfaStruct) panicHelper(fmtString string, fmtArgs ...interface{}) {
// 	bwerr.PanicErr(fmt.Errorf(err.Error()+"\n"+ansi.Ansi("", fmtString), fmtArgs), 1)
// 	// bwerr.PanicErr(ansi.Ansi("", fmtString), fmtArgs), 1)
// }

// func (pfa *PfaStruct) ifStackLen(minLen int) bool {
// 	return len(pfa.Stack) >= minLen
// }

// func (pfa *PfaStruct) mustStackLen(minLen int) {
// 	if !pfa.ifStackLen(minLen) {
// 		pfa.Panic(bw.StructFrom("<ansiVar>minLen <ansiVal>%d", minLen))
// 	}
// }

// func (pfa *PfaStruct) GetTopStackItem(optDeep ...uint) *ParseStackItem {
// 	ofs := -1
// 	if optDeep != nil {
// 		ofs = ofs - int(optDeep[0])
// 	}
// 	pfa.mustStackLen(-ofs)
// 	return &pfa.Stack[len(pfa.Stack)+ofs]
// }

// func (pfa *PfaStruct) PopStackItem() {
// 	pfa.mustStackLen(1)
// 	pfa.Stack = pfa.Stack[:len(pfa.Stack)-1]
// }

// func (pfa *PfaStruct) PushStackItem() {
// 	pfa.Stack = append(pfa.Stack, ParseStackItem{
// 		Start: pfa.Proxy.Curr,
// 		Vars:  map[string]interface{}{},
// 	})
// }

// func (pfa PfaStruct) DataForJSON() interface{} {
// 	result := map[string]interface{}{}
// 	// result["Stack"] = pfa.Stack.DataForJSON()
// 	result["Proxy"] = pfa.Proxy.DataForJSON()
// 	if len(pfa.Vars) > 0 {
// 		result["Vars"] = pfa.Vars
// 	}
// 	return result
// }

func (pfa PfaStruct) String() string {
	return bwjson.Pretty(pfa)
}

// ============================================================================

var underscoreRegexp = regexp.MustCompile("[_]+")

func ParseNumber(source string) (value interface{}, err error) {
	source = underscoreRegexp.ReplaceAllLiteralString(source, ``)
	if strings.Contains(source, `.`) {
		var _float64 float64
		if _float64, err = strconv.ParseFloat(source, 64); err == nil {
			value = _float64
		}
	} else {
		var _int64 int64
		if _int64, err = strconv.ParseInt(source, 10, 64); err == nil {
			if int64(bw.MinInt8) <= _int64 && _int64 <= int64(bw.MaxInt8) {
				value = int8(_int64)
			} else if int64(bw.MinInt16) <= _int64 && _int64 <= int64(bw.MaxInt16) {
				value = int16(_int64)
			} else if int64(bw.MinInt32) <= _int64 && _int64 <= int64(bw.MaxInt32) {
				value = int32(_int64)
			} else {
				value = _int64
			}
		}
	}
	return
}

func ParseInt(source string) (value int, err error) {
	source = underscoreRegexp.ReplaceAllLiteralString(source, ``)
	var _int64 int64
	if _int64, err = strconv.ParseInt(source, 10, 64); err == nil {
		if int64(bw.MinInt) <= _int64 && _int64 <= int64(bw.MaxInt) {
			value = int(_int64)
		} else {
			err = bwerr.From("<ansiVal>%d<ansi> is out of range <ansiVal>[%d, %d]", _int64, bw.MinInt, bw.MaxInt)
		}
	}
	return
}

// ============================================================================

type ParseStack []ParseStackItem

func (Stack *ParseStack) MarshalJSON() ([]byte, error) {
	result := []interface{}{}
	for _, item := range *Stack {
		result = append(result, item)
	}
	return json.Marshal(result)
}

func (Stack *ParseStack) String() (result string) {
	return bwjson.Pretty(Stack)
}

// ============================================================================

type ParseStackItem struct {
	Start runeprovider.PosInfo
	Vars  map[string]interface{}
}

func (stackItem *ParseStackItem) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{}
	result["Start"] = stackItem.Start.DataForJSON()
	result["Vars"] = stackItem.Vars
	return json.Marshal(result)
	// return result
}

func (stackItem *ParseStackItem) String() (result string) {
	return bwjson.Pretty(stackItem)
}
