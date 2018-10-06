package defvalid

import (
	// "github.com/baza-winner/bwcore/bwerror"
	// "github.com/baza-winner/bwcore/bwjson"
// "log"
)

func compileDef(def value) (result *Def, err error) {
	if def.value == nil {
		return nil, defErrorMake(def, defErrorHasUnexpectedValue)
	}
	var defType value
	var isSimpleDef bool
	if isSimpleDef = _isOfType(def.value, "string", "[]string"); isSimpleDef {
		defType = def
	} else if !_isOfType(def.value, "map") {
		return nil, valueErrorMake(def, valueErrorIsNotOfType, "string", "[]string", "map")
	} else {
		if defType, err = def.getKey("type", []string{"string", "[]string"}); err != nil {
			return nil, err
		}
		// } else if skipDefault == nil || !skipDefault[0] {
		//   defaultVal, err = def.getKey("default")
		//   if err != nil {
		//     if valErr, ok := err.(valueError); ok && valErr.errorType == valueErrorHasNoKey {
		//       err = nil
		//     } else {
		//       return nil, err
		//     }
		//   } else if defaultVal.value == nil {
		//     return nil, valueErrMake(defaultVal, valueErrorHasNonSupportedValue)
		//   } else if defaultVal.value, err = getValidVal(defaultVal, def, true); err != nil {
		//     return nil, err
		//   }
	}

  var tp deftype
  if tp, err = getDeftype(defType, isSimpleDef); err != nil {
  	return nil, err
  }
	// log.Printf("tp: %s\n", bwjson.PrettyJson(tp.GetDataForJson()))
	result = &Def{tp: tp}
	// log.Printf("result: %s\n", bwjson.PrettyJson(result.GetDataForJson()))
	// if (_isOfType(defType, "string"))
	// bwerror.Noop(defType)
	return
}

func getDeftype(defType value, isSimpleDef bool) (result deftype, err error) {
	var isString bool
	var ss []string
	if _isOfType(defType.value, "string") {
		ss = []string{_mustBeString(defType.value)}
		isString = true
	} else if _isOfType(defType.value, "[]string") {
		ss = _mustBeSliceOfStrings(defType.value)
		isString = false
	} else {
		return nil, valueErrorMake(defType, valueErrorIsNotOfType, "string", "[]string")
	}
	result = deftype{}
	// if err == nil {
	for i, s := range ss {
		switch s {
		case "bool":
			result.Add(deftypeBool)
		case "string":
			result.Add(deftypeString)
		case "int":
			result.Add(deftypeInt)
		case "number":
			result.Add(deftypeNumber)
		case "map":
			result.Add(deftypeMap)
		case "array":
			result.Add(deftypeArray)
		case "orArrayOf":
			result.Add(deftypeOrArrayOf)
		default:
			if isString {
				err = valueErrorMake(defType, valueErrorHasNonSupportedValue)
			} else {
				var elem value
				if elem, err = defType.getElem(i, "string"); err != nil { return nil, err }
				// log.Printf("elem: %s\n", bwjson.PrettyJsonOf(elem))
				err = valueErrorMake(elem, valueErrorHasNonSupportedValue)
			}
		}
	}
	// log.Printf("result: %s\n", bwjson.PrettyJson(result.GetDataForJson()) )
	// log.Printf("result: %s\n", result)
	// result = bwset.FromSlice(ss)
	// if result.Has("enum") && result.Has("string") {
	// 	// log.Printf("%s\n", result)
	// 	err = valueErrorMake(defType, valueErrorValuesCannotBeCombined, "enum", "string")
	// }
	// // }
	return
}
// func getDeftype
