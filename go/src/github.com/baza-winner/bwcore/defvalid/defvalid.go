// Предоставляет функции для валидации interface{}.
package defvalid

import (
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"

	// "github.com/baza-winner/bwcore/bwmap"
	// "github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/defparse"
	// "log"
	// "reflect"
)

// ==================================================================================

func CompileDef(def interface{}) (result *Def, err error) {
	result, err = compileDef(value{"<ansiVar>def<ansiPath>", def})
	// log.Printf("result: %#v", result)
	// var compileDefResult *Def
	// compileDefResult, err = compileDef(value{"<ansiVar>def<ansiPath>", def})
	// if err == nil {
	// 	if compileDefResult == nil {
	// 		bwerr.Panic("Unexpected behavior; def: %s", bwjson.PrettyJson(def))
	// 	} else {
	// 		result = *compileDefResult
	// 	}
	// }
	return
}

func MustCompileDef(def interface{}) (result Def) {
	var err error
	var compileDefResult *Def
	if compileDefResult, err = CompileDef(def); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	} else if compileDefResult == nil {
		bwerr.Panic("Unexpected behavior; def: %s", bwjson.Pretty(def))
	} else {
		result = *compileDefResult
	}
	return
}

// ==================================================================================

func ValidateVal(what string, val interface{}, def Def) (result interface{}, err error) {
	return getValidVal(
		value{
			value: val,
			what:  "<ansiVar>" + what + "<ansiPath>",
		},
		def,
	)
}

func MustValidVal(what string, val interface{}, def Def) (result interface{}) {
	var err error
	if result, err = ValidateVal(what, val, def); err != nil {
		bwerr.PanicA(bwerr.Err(err))
	}
	return
}

// ==================================================================================

func GetValOfPath(val interface{}, path string) (result interface{}, valueError error) {
	switch path {
	case ".keys.keyOne":
		result = defparse.MustParse("{ type: 'bool', default: false, some: \"thing\" }")
	case ".type":
		result = defparse.MustParse("['enum', 'string']")
	case ".enum":
		result = defparse.MustParse("['one', true, 3 ]")
	}
	return
}

func MustValOfPath(val interface{}, path string) (result interface{}) {
	var err error
	if result, err = GetValOfPath(val, path); err != nil {
		bwerr.Panic("path <ansiPath>%s<ansi> not found in <ansiVal>%s", path, bwjson.Pretty(val))
	}
	return
}
