// Предоставляет функции для валидации interface{}.
package defvalid

import (
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	// "github.com/baza-winner/bwcore/bwmap"
	// "github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/defparse"
	// "log"
	// "reflect"
)

// ==================================================================================

func CompileDef(def interface{}) (result *Def, err error) {
	result, err = compileDef(value{"<ansiOutline>def<ansiCmd>", def})
	// log.Printf("CompileDef result: %s\n", bwjson.PrettyJson(result.GetDataForJson()))
	return
}

func MustCompileDef(def interface{}) (result *Def) {
	var err error
	if result, err = CompileDef(def); err != nil {
		bwerror.PanicErr(err)
	}
	return
}

// ==================================================================================

func ValidateVal(what string, val, def interface{}) (result interface{}, err error) {
	return getValidVal(
		value{
			value: val,
			what:  "<ansiOutline>" + what + "<ansiCmd>",
		},
		value{
			value: def,
			what:  "<ansiOutline>" + what + "::def<ansiCmd>",
		},
	)
}

func MustValidVal(what string, val, def interface{}) (result interface{}) {
	var err error
	if result, err = ValidateVal(what, val, def); err != nil {
		bwerror.PanicErr(err)
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
		bwerror.Panic("path <ansiCmd>%s<ansi> not found in <ansiSecondaryLiteral>%s", path, bwjson.PrettyJson(val))
	}
	return
}
