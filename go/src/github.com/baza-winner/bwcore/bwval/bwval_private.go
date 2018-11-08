package bwval

import (
	"fmt"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
)

// ============================================================================

//go:generate stringer -type=ValKind

func varsJSON(path bw.ValPath, optVars []map[string]interface{}) (result string) {
	if hasVar(path) {
		var vars map[string]interface{}
		if len(optVars) > 0 {
			vars = optVars[0]
		}
		result = fmt.Sprintf(ansiVars, bwjson.Pretty(vars))
	}
	return
}

type valHolder struct {
	val interface{}
}

var (
	ansiIsNotOfType                string
	ansiValAtPathIsNotOfType       string
	ansiMustPathValFailed          string
	ansiType                       string
	ansiValAtPathIsNotOfTypes      string
	ansiValAtPathIsNil             string
	ansisReadOnlyPath              string
	ansiValAtPathHasNotEnoughRange string
	ansiVars                       string
	ansiVarsIsNil                  string
	ansiMustSetPathValFailed       string
)

func init() {
	ansiIsNotOfType = ansi.String("<ansiVal>%#v<ansi> is not <ansiType>%s")
	ansiMustSetPathValFailed = ansi.String("Failed to set <ansiPath>%s<ansi> of <ansiVal>%s<ansi>%s: {Error}")
	ansiMustPathValFailed = ansi.String("Failed to get <ansiPath>%s<ansi> of <ansiVal>%s<ansi>%s: {Error}")
	valPathPrefix := "<ansiPath>%s<ansi> "
	ansiValAtPathIsNotOfType = ansi.String(valPathPrefix + "(<ansiVal>%s<ansi>) is not <ansiType>%s")
	ansiValAtPathIsNotOfTypes = ansi.String(valPathPrefix + "(<ansiVal>%s<ansi>) is none of %s")
	ansiValAtPathIsNil = ansi.String(valPathPrefix + "is <ansiErr>nil")
	ansiValAtPathHasNotEnoughRange = ansi.String(valPathPrefix + "(<ansiVal>%s<ansi>) has not enough length (<ansiVal>%d<ansi>) for idx (<ansiVal>%d)")

	ansiType = ansi.String("<ansiType>%s")
	ansisReadOnlyPath = ansi.String("<ansiPath>%s<ansi> is <ansiErr>readonly path")
	ansiVars = ansi.String(" with <ansiVar>vars<ansi> <ansiVal>%s<ansi>")
	ansiVarsIsNil = ansi.String("<ansiVar>vars<ansi> is <ansiErr>nil")
}

func hasVar(path bw.ValPath) bool {
	for _, vpi := range path {
		switch vpi.Type {
		case bw.ValPathItemVar:
			return true
		case bw.ValPathItemPath:
			if hasVar(vpi.Path) {
				return true
			}
		}
	}
	return false
}

func simplifyPath(v *valHolder, path bw.ValPath, optVars []map[string]interface{}) (result bw.ValPath, err error) {
	result = bw.ValPath{}
	for _, vpi := range path {
		if vpi.Type != bw.ValPathItemPath {
			result = append(result, vpi)
		} else {
			var val interface{}
			val, err = v.PathVal(vpi.Path, optVars...)
			if err == nil {
				switch _, kind := Kind(val); kind {
				case ValString:
					result = append(result, bw.ValPathItem{Type: bw.ValPathItemKey, Key: MustString(val)})
				case ValInt:
					result = append(result, bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: MustInt(val)})
				default:
					// bwdebug.Print("vpi.Path", vpi.Path)
					err = v.pathValIsNotOfType(vpi.Path, val, "Int", "String")
				}
			}
		}
	}
	return
}

var typeSeparator = map[bool]string{
	true:  " or ",
	false: ", ",
}

func (v *valHolder) pathValIsNotOfType(path bw.ValPath, val interface{}, expectedType string, optExpectedType ...string) (result error) {
	if len(optExpectedType) == 0 {
		result = bwerr.From(ansiValAtPathIsNotOfType, path, bwjson.Pretty(val), expectedType)
	} else {
		expectedTypes := fmt.Sprintf(ansiType, expectedType)
		for i, elem := range optExpectedType {
			expectedTypes += typeSeparator[i == len(optExpectedType)-1] + fmt.Sprintf(ansiType, elem)
		}
		result = bwerr.From(ansiValAtPathIsNotOfTypes, path, bwjson.Pretty(val), expectedTypes)
		// bwdebug.Print("path", path, "ansiValAtPathIsNotOfTypes", ansiValAtPathIsNotOfTypes, "result", result)
	}
	return
}

func (v *valHolder) valAtPathIsNil(path bw.ValPath) error {
	return bwerr.From(ansiValAtPathIsNil, path)
}

func (v *valHolder) getArray(idx int, result interface{}, resultPath bw.ValPath) ([]interface{}, int, error) {
	var err error
	var ok bool
	var vals []interface{}
	if vals, ok = Array(result); !ok {
		err = v.pathValIsNotOfType(resultPath, result, "Array")
	} else {
		l := len(vals)
		minIdx := -l
		maxIdx := l - 1
		if !(minIdx <= idx && idx <= maxIdx) {
			idx = -1
		} else if idx < 0 {
			idx = l + idx
		}
	}
	return vals, idx, err
}

func (v *valHolder) getMap(result interface{}, resultPath bw.ValPath) (map[string]interface{}, error) {
	var err error
	var ok bool
	var m map[string]interface{}
	if m, ok = Map(result); !ok {
		err = v.pathValIsNotOfType(resultPath, result, "Map")
	}
	return m, err
}

func (v *valHolder) byKey(val interface{}, path bw.ValPath, i int, key string) (result interface{}, err error) {
	result = val
	if result == nil {
		return
	}
	var m map[string]interface{}
	if m, err = v.getMap(result, path[:i+1]); err == nil {
		result = m[key]
	}
	return
}

func (v *valHolder) byIdx(val interface{}, path bw.ValPath, i int, idx int) (result interface{}, err error) {
	result = val
	if result == nil {
		return
	}
	var vals []interface{}
	if vals, idx, err = v.getArray(idx, result, path[:i+1]); err == nil {
		if idx < 0 {
			result = nil
		} else {
			result = vals[idx]
		}
	}
	return
}
